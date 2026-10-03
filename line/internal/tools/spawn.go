// ONE SPAWN CONTRACT. ADR-006 item 5, accepted 2026-09-11.
//
// The door is a process supervisor, and the contract for "how do I correctly
// start a child" was invented at every call site. ADR-006 counted six seams
// and found four different answers to the same four questions: does stdin get
// closed, is there a timeout, does the environment get hardened, and is the
// child's own output kept when it fails.
//
//	gitctl.go     git    devnull,
//	gitstate.go   git    devnull, timeout, no env hardening, stderr DISCARDED
//	tools.go:242  spine  no stdin, NO timeout, combined output
//	tools.go:2276 engine no stdin, timeout, separate buffers
//
// Every difference is a decision somebody made once and nobody reviewed. The
// two that actually bit: gitstate threw stderr away, so a git failure there
// was a shrug; and the spine call had no timeout at all, so a hung Rust binary
// hangs a tool call forever.
//
// WHAT IS DELIBERATELY NOT HERE. Two of the six seams stay as they are, with
// reasons:
//
//   - internal/engine spawns the Manjuel engine over StdinPipe, because that
//     pipe IS the wire. A helper that closes stdin would break it by doing its
//     job. ADR-006 says so in the same table.
//   - cmd/atlas-door carries its own copy in a DIFFERENT BINARY. Sharing with
//     it would mean a new leaf package, because internal/tools already imports
//     internal/engine and a helper both could reach cannot live in either. A
//     new package is a new folder, and folders are the operator's to place
//     (RULE 8). Left, named here, rather than invented.
//
// Stdlib only, per the atlas law.
package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// spawnOpts is what a caller may vary. Everything not named here is fixed on
// purpose -- that is the point of having one contract.
type spawnOpts struct {
	// Dir is the working directory. Always a tenant's Home in practice, and
	// never empty: a child inheriting the door's own cwd is how a tool ends up
	// acting on the wrong ground.
	Dir string
	// Timeout bounds the child. Zero takes spawnDefaultTimeout rather than
	// meaning "forever" -- ESTATE LAW 7 is bounded everything, and the one
	// seam that had no timeout was the one that could hang a tool call.
	Timeout time.Duration
	// Env is appended to the process environment. Used to harden git.
	Env []string
	// CleanEnv, when non-nil, IS the child's whole environment: nothing of the
	// door's own is inherited (Env is still appended). The operator's shell
	// runs on this (shell.go): the door holds the service wire and the keys in
	// its environment, and a child that inherited them could print them.
	CleanEnv []string
	// MaxBytes keeps the first MaxBytes of each stream and drains the rest
	// (spawnResult.Truncated says so). Zero is unbounded, as it always was.
	MaxBytes int
	// Grace bounds how long Wait lingers after the deadline for a child's
	// children to let go of the pipes. Zero is the old behaviour: forever.
	Grace time.Duration
	// KillTree makes the deadline take the child's children with it (taskkill
	// /T on Windows), so a pipeline that outlives its shell does not outlive
	// its limit.
	KillTree bool
}

// spawnDefaultTimeout bounds any child whose caller did not choose. Generous
// enough for a git operation on a large tree, short enough that a wedged
// binary surfaces as a refusal rather than a hang.
const spawnDefaultTimeout = 120 * time.Second

// spawnResult is everything the child said and how it ended.
type spawnResult struct {
	Stdout string
	Stderr string
	// Combined is stdout followed by stderr, which is what a caller shows a
	// human. Kept separately from Err on purpose: ADR-006 item 2 is that a
	// tool's own words outrank the Go error, and a caller cannot honour that
	// if the words were never returned.
	Combined string
	Err      error
	TimedOut bool
	// Truncated is whether MaxBytes cut either stream.
	Truncated bool
}

// capBuffer keeps the first max bytes written to it and swallows the rest, so a
// child that floods its pipe is drained, never blocked, and never held whole in
// memory. max <= 0 keeps everything, as bytes.Buffer does.
type capBuffer struct {
	max     int
	buf     bytes.Buffer
	dropped bool
}

func (c *capBuffer) Write(p []byte) (int, error) {
	if c.max <= 0 {
		return c.buf.Write(p)
	}
	if room := c.max - c.buf.Len(); room > 0 {
		if len(p) <= room {
			return c.buf.Write(p)
		}
		c.buf.Write(p[:room])
	}
	c.dropped = true
	return len(p), nil
}

func (c *capBuffer) String() string { return c.buf.String() }

// killTree ends a child and, on Windows, everything it started.
func killTree(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	if runtime.GOOS == "windows" {
		_ = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
	}
	return cmd.Process.Kill()
}

// spawn runs one child under the single contract:
//
//   - STDIN IS CLOSED. Every child gets os.DevNull. Go already does this for a
//     nil Stdin, but two of the six seams set it explicitly and four did not,
//     and a reader cannot tell "deliberate" from "forgotten" by looking. Now
//     there is one answer and it is written down.
//   - IT IS BOUNDED. No child runs without a deadline (ESTATE LAW 7).
//   - ITS WORDS ARE KEPT, both streams, whether it succeeded or failed. The
//     caller decides what to show; this never decides for it by discarding.
//   - A TIMEOUT SAYS SO. "signal: killed" tells a human nothing; TimedOut and
//     a named error tell them the thing did not finish.
func spawn(name string, args []string, opts spawnOpts) spawnResult {
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = spawnDefaultTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = opts.Dir
	switch {
	case opts.CleanEnv != nil:
		cmd.Env = append(append([]string(nil), opts.CleanEnv...), opts.Env...)
	case len(opts.Env) > 0:
		cmd.Env = append(os.Environ(), opts.Env...)
	}
	if opts.KillTree {
		cmd.Cancel = func() error { return killTree(cmd) }
	}
	if opts.Grace > 0 {
		cmd.WaitDelay = opts.Grace
	}
	// Closed, explicitly, at the one place that decides it.
	if devnull, err := os.Open(os.DevNull); err == nil {
		cmd.Stdin = devnull
		defer devnull.Close()
	}

	out, errBuf := &capBuffer{max: opts.MaxBytes}, &capBuffer{max: opts.MaxBytes}
	cmd.Stdout = out
	cmd.Stderr = errBuf
	err := cmd.Run()

	res := spawnResult{
		Stdout:    out.String(),
		Stderr:    errBuf.String(),
		Err:       err,
		Truncated: out.dropped || errBuf.dropped,
	}
	res.Combined = res.Stdout
	if res.Stderr != "" {
		if res.Combined != "" && !strings.HasSuffix(res.Combined, "\n") {
			res.Combined += "\n"
		}
		res.Combined += res.Stderr
	}
	if ctx.Err() == context.DeadlineExceeded {
		res.TimedOut = true
		res.Err = fmt.Errorf("REFUSED — %s did not finish inside %ds. Nothing is claimed.",
			name, int(timeout.Seconds()))
	}
	return res
}

// gitEnv hardens git so it can never stop for a human this process cannot show
// to anyone. Carried over from gitctl, where it was the only seam that had it.
func gitEnv() []string {
	return []string{"GIT_TERMINAL_PROMPT=0", "GIT_PAGER=cat"}
}
