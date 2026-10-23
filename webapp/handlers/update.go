// ATLAS UPDATES ITSELF (2026-10-22, his word: "atlas should auto-update just
// like the secondbrain does"), AND SINCE 2026-10-23 EVERYTHING ELSE ON THE
// GROUND (his word: "is there a better way ... where my shit always just like
// gets the latest version from github so i dont have to keep like going in and
// manually pulling and fast forwarding"; "build the glass button").
//
// Three checkouts, one button: Atlas (main), the Manjuel ground one up from it
// (main), and the notes kernel at worlds/jesster (the branch it tracks).
//
// GET /api/update/check fetches each and says how far behind it is, what it
// would take, and why it would not move if it would not. For Atlas it also asks
// whether the commit this glass was built from -- Go writes it into the program
// (debug.ReadBuildInfo, vcs.revision) -- is behind the checkout.
//
// POST /api/update/apply refuses everything while a sitting is open. Otherwise
// it takes Manjuel and the kernel forward here, fast-forward only, each on its
// own (one that cannot move is skipped and says why, and the rest still move),
// and answers with what moved. Manjuel moving means "restart required": the
// core is his to restart. Atlas is left to ../update.ps1, started only when
// Atlas has something to take or to build: it takes main forward, builds the
// door and the glass, swaps them in and starts both again, and the page waits
// for the glass to answer and reloads.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// builtFrom is the commit this program was built from, and whether the tree
// held unsaved changes then; "" when it was built outside a repository.
func builtFrom() (rev string, modified bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	return rev, modified
}

// atlasDir is the repository the glass was built in: one up from atlas/webapp.
// A variable so the strokes can stand a ground of their own in its place.
var atlasDir = func() string {
	d, _ := filepath.Abs("..")
	return d
}

// gitIn reads: git's answer alone (stdout), so a warning git prints beside it
// -- on Windows, every line-ending one -- is never read as part of it.
func gitIn(ctx context.Context, dir string, args ...string) (string, error) {
	return gitRun(ctx, dir, false, args...)
}

// gitDo acts (fetch, merge): stdout and stderr together, so a refusal says
// git's own words.
func gitDo(ctx context.Context, dir string, args ...string) (string, error) {
	return gitRun(ctx, dir, true, args...)
}

func gitRun(ctx context.Context, dir string, both bool, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(cmd.Environ(), "GIT_TERMINAL_PROMPT=0")
	var out []byte
	var err error
	if both {
		out, err = hidden(cmd).CombinedOutput()
	} else {
		out, err = hidden(cmd).Output()
	}
	return strings.TrimSpace(string(out)), err
}

// checkout is one repository the button keeps up with GitHub. Branch "" means
// whatever branch the checkout is on, as long as it tracks one on GitHub.
type checkout struct {
	Name   string
	Dir    string
	Branch string
}

// checkouts is the whole list, from Atlas's own directory: Atlas, the ground it
// stands in, and the notes kernel's checkout in the ground's worlds.
func checkouts(atlas string) []checkout {
	ground := filepath.Dir(atlas)
	return []checkout{
		{Name: "Atlas", Dir: atlas, Branch: "main"},
		{Name: "Manjuel", Dir: ground, Branch: "main"},
		{Name: "notes kernel", Dir: filepath.Join(ground, "worlds", "jesster")},
	}
}

// repoState is one checkout as GitHub now has it: fetched, then counted.
type repoState struct {
	Name     string   `json:"name"`
	Branch   string   `json:"branch,omitempty"`
	Upstream string   `json:"upstream,omitempty"`
	Head     string   `json:"head,omitempty"`
	Behind   int      `json:"behind"`
	Changes  []string `json:"changes,omitempty"`
	Skip     string   `json:"skip,omitempty"`    // why it will not move; "" when it can
	Unsaved  string   `json:"unsaved,omitempty"` // changed files, named; only Atlas is skipped for them
}

// state fetches c's upstream and says how far behind it is and whether it may
// be moved. Every refusal is a sentence he can act on.
func state(ctx context.Context, c checkout) repoState {
	s := repoState{Name: c.Name}
	if _, err := os.Stat(c.Dir); err != nil {
		s.Skip = "not on this machine (" + c.Dir + ")"
		return s
	}
	if _, err := gitIn(ctx, c.Dir, "rev-parse", "--git-dir"); err != nil {
		s.Skip = "not a git checkout (" + c.Dir + ")"
		return s
	}
	s.Branch, _ = gitIn(ctx, c.Dir, "rev-parse", "--abbrev-ref", "HEAD")
	s.Head, _ = gitIn(ctx, c.Dir, "rev-parse", "--short", "HEAD")
	if c.Branch != "" && s.Branch != c.Branch {
		s.Skip = fmt.Sprintf("on %q, not %s; only %s is kept up with GitHub", s.Branch, c.Branch, c.Branch)
		return s
	}
	up, err := gitIn(ctx, c.Dir, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{upstream}")
	if err != nil || up == "" {
		s.Skip = fmt.Sprintf("%q follows no branch on GitHub", s.Branch)
		return s
	}
	s.Upstream = up
	remote, branch, ok := strings.Cut(up, "/")
	if !ok {
		s.Skip = "its upstream " + up + " is not on a remote"
		return s
	}
	if out, err := gitDo(ctx, c.Dir, "fetch", "-q", remote, branch); err != nil {
		s.Skip = "GitHub could not be asked: " + firstLineOf(out, err)
		return s
	}
	if n, err := gitIn(ctx, c.Dir, "rev-list", "--count", "HEAD.."+up); err == nil {
		s.Behind, _ = strconv.Atoi(n)
	}
	if s.Behind > 0 {
		subjects, _ := gitIn(ctx, c.Dir, "log", "--format=%s", "-n", "10", "HEAD.."+up)
		s.Changes = strings.Split(subjects, "\n")
	}
	if ahead, _ := gitIn(ctx, c.Dir, "rev-list", "--count", up+"..HEAD"); ahead != "" && ahead != "0" {
		s.Skip = fmt.Sprintf("has %s save(s) GitHub does not; joining the two is yours", ahead)
		return s
	}
	if lock, _ := gitIn(ctx, c.Dir, "rev-parse", "--git-path", "index.lock"); lock != "" {
		if !filepath.IsAbs(lock) {
			lock = filepath.Join(c.Dir, lock)
		}
		if _, err := os.Stat(lock); err == nil {
			s.Skip = "git's lock file is there (" + lock + "): git is running, or one that crashed left it"
			return s
		}
	}
	// UNSAVED WORK. Atlas is rebuilt from its checkout, and update.ps1 refuses over unsaved work, so Atlas is skipped
	// here first. Manjuel and the kernel are not: the ground's suites rewrite tracked stamps (tests/last_run.*) every
	// run, and a fast-forward over a change it does not touch is safe. Git's own guard decides -- merge --ff-only
	// refuses to overwrite a local change, in its own words -- and the files are named beside the count.
	if names, _ := gitIn(ctx, c.Dir, "diff", "--name-only", "HEAD"); names != "" {
		files := strings.Split(names, "\n")
		more := ""
		if len(files) > 3 {
			more = fmt.Sprintf(" and %d more", len(files)-3)
			files = files[:3]
		}
		s.Unsaved = strings.Join(files, ", ") + more
		if c.Name == "Atlas" {
			s.Skip = "has unsaved work (" + s.Unsaved + "); save it or put it aside first"
		}
	}
	return s
}

// forward takes a checkout that state let through up to its upstream,
// fast-forward only, and says what happened.
func forward(ctx context.Context, c checkout, s repoState) string {
	if out, err := gitDo(ctx, c.Dir, "merge", "-q", "--ff-only", s.Upstream); err != nil {
		return "did not move: " + firstLineOf(out, err)
	}
	now, _ := gitIn(ctx, c.Dir, "rev-parse", "--short", "HEAD")
	return fmt.Sprintf("moved %s -> %s (%d save%s)", s.Head, now, s.Behind, plural(s.Behind))
}

func firstLineOf(out string, err error) string {
	if out = strings.TrimSpace(out); out != "" {
		return strings.SplitN(out, "\n", 2)[0]
	}
	return err.Error()
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// sittingOpen is why nothing may move now, or "": the last sitting on the
// ground's ledger has no end (RULE 9: nothing moves under a sitting).
func sittingOpen(atlas string) string {
	b, err := os.ReadFile(filepath.Join(filepath.Dir(atlas), "sessions", "sessions.jsonl"))
	if err != nil {
		return ""
	}
	lines := strings.Split(strings.TrimRight(string(b), "\r\n"), "\n")
	var row struct {
		N     any    `json:"n"`
		Ended string `json:"ended"`
	}
	if last := strings.TrimSpace(lines[len(lines)-1]); last != "" && json.Unmarshal([]byte(last), &row) == nil && strings.TrimSpace(row.Ended) == "" {
		return fmt.Sprintf("sitting %v is open; close it first (nothing moves under a sitting)", row.N)
	}
	return ""
}

// atlasBuild says whether this glass is older than Atlas's checkout, and by how
// much: the update rebuilds Atlas when either it or the checkout is behind.
func atlasBuild(ctx context.Context, dir string) (running string, behind int) {
	rev, _ := builtFrom()
	if rev == "" {
		return "", 0
	}
	if n, err := gitIn(ctx, dir, "rev-list", "--count", rev+"..HEAD"); err == nil {
		behind, _ = strconv.Atoi(n)
	}
	return short(rev), behind
}

// UpdateCheck answers GET /api/update/check.
func (h *Handlers) UpdateCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 150*time.Second)
	defer cancel()
	atlas := atlasDir()
	running, built := atlasBuild(ctx, atlas)
	out := map[string]any{"running": running, "built_behind": built}
	if why := sittingOpen(atlas); why != "" {
		out["refused"] = why
	}
	var repos []repoState
	total := 0
	for _, c := range checkouts(atlas) {
		s := state(ctx, c)
		if s.Skip == "" {
			total += s.Behind
		}
		repos = append(repos, s)
	}
	if repos[0].Skip == "" && repos[0].Behind == 0 {
		total += built
	}
	out["repos"] = repos
	out["behind"] = total
	out["available"] = total > 0
	jsonResp(w, out)
}

// UpdateApply answers POST /api/update/apply.
func (h *Handlers) UpdateApply(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 300*time.Second)
	defer cancel()
	atlas := atlasDir()
	if why := sittingOpen(atlas); why != "" {
		jsonErr(w, 409, why)
		return
	}
	var said []map[string]string
	restart := []string{}
	atlasGoes := false
	for _, c := range checkouts(atlas) {
		s := state(ctx, c)
		row := map[string]string{"name": c.Name}
		switch {
		case s.Skip != "":
			row["result"] = "skipped: " + s.Skip
		case c.Name == "Atlas":
			_, built := atlasBuild(ctx, atlas)
			if s.Behind > 0 || built > 0 {
				atlasGoes = true
				row["result"] = fmt.Sprintf("taking %d save%s and rebuilding; Atlas restarts itself", s.Behind, plural(s.Behind))
			} else {
				row["result"] = "up to date"
			}
		case s.Behind == 0:
			row["result"] = "up to date"
		default:
			row["result"] = forward(ctx, c, s)
			if strings.HasPrefix(row["result"], "moved") {
				switch c.Name {
				case "Manjuel":
					restart = append(restart, "Manjuel moved: restart required (the core runs the old code until it is restarted)")
				case "notes kernel":
					restart = append(restart, "the notes kernel moved: Manjuel's notes calls use it from the next call; restart SecondBrain for its own")
				}
			}
		}
		said = append(said, row)
		updateLog(atlas, c.Name+": "+row["result"])
	}
	out := map[string]any{"repos": said, "restart": restart, "atlas": atlasGoes}
	if atlasGoes {
		script := filepath.Join(atlas, "update.ps1")
		shell := "pwsh"
		if _, err := exec.LookPath(shell); err != nil {
			shell = "powershell"
		}
		cmd := exec.Command(shell, "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", script)
		cmd.Dir = atlas
		if err := hidden(cmd).Start(); err != nil {
			out["atlas"] = false
			out["atlas_error"] = "the Atlas update did not start: " + err.Error()
		} else {
			go cmd.Wait()
		}
	}
	out["log"] = filepath.Join(atlas, "update.log")
	jsonResp(w, out)
}

// updateLog adds one line to atlas/update.log, the same log update.ps1 writes.
func updateLog(atlas, line string) {
	f, err := os.OpenFile(filepath.Join(atlas, "update.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s  update: %s\r\n", time.Now().Format("2006-01-02 15:04:05"), line)
}

func short(rev string) string {
	if len(rev) > 7 {
		return rev[:7]
	}
	return rev
}
