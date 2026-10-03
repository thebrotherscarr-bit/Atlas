package tools

// THE OPERATOR'S TYPED SHELL -- shell_run (WHAT'S LEFT H15, 2026-10-03). His word of 2026-10-02, on the
// card: "typed by you, gated". Bash (Git Bash) and Python, run in a world, from the front page's two tabs.
//
//	WHAT YOU TYPE RUNS AS YOU, and a plain look (ls, cat, grep, git status/log/diff, a calculation) runs at
//	once. Anything that writes, deletes, installs, reaches the network, saves or sends work, or that the
//	gate cannot read, is PARKED on the same queue every other writing call waits in (holds.go) and answers
//	"held" with its id; the glass shows a card, and hold_answer runs EXACTLY what was parked. A secret, a
//	path outside the ground, client material or a key typed into a command is refused by name, and no
//	approval lifts a refusal (shellgate.go).
//
//	NO AGENT IS EVER GIVEN THE SHELL. The tool is ServiceOnly: Registry.Call refuses every caller but the
//	glass before RBAC or the holds look at it, so a role cannot grant it and an agent cannot park a call
//	for him to approve. Approval is a typed field of the Caller (Hold), set by hold_answer alone.
//
//	THE LAW GATE. Nothing runs on a law that does not verify: the sealed chain is walked with the law's
//	own tool (law/law.py verify), cached against what law/ holds, exactly as the engine's gate does for
//	every council turn (manjuel/lawgate.py).
//
//	KEYS ARE SILENT. The child's environment is built, not inherited: the door holds the service wire and
//	the keys in its own, and a shell that inherited them could print them. What it prints is then scrubbed
//	of every secret VALUE the ground holds (.env) or the door's environment names, so a command that
//	reaches a secret by a route the gate did not see still cannot show it.
//
//	RECORDED. Every run, hold and refusal is a line in the ground's holds log (state/holds.jsonl) with the
//	command (keys withheld) -- the log hold_answer already writes, not a new one.

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"atlas/line/internal/tenant"
)

const (
	shellDefaultWait = 30 * time.Second
	shellMaxWait     = 300 * time.Second
	shellMaxBytes    = 64 * 1024
	shellMaxCommand  = 8000
	shellWithheld    = "[withheld: a secret]"
)

// shellAnswer is the one document shell_run answers with, for every outcome that is not a malformed call.
// The glass parses it; its keys are the contract (TestTheShellAnswersInOneShape).
type shellAnswer struct {
	Shell     string   `json:"shell"`
	State     string   `json:"state"` // ran | held | refused
	Class     string   `json:"class,omitempty"`
	Why       []string `json:"why,omitempty"`
	Hold      string   `json:"hold,omitempty"`
	Approved  bool     `json:"approved,omitempty"`
	Exit      *int     `json:"exit,omitempty"`
	Ms        int64    `json:"ms,omitempty"`
	TimedOut  bool     `json:"timed_out,omitempty"`
	Truncated bool     `json:"truncated,omitempty"`
	Output    string   `json:"output,omitempty"`
	Note      string   `json:"note,omitempty"`
	Command   string   `json:"command,omitempty"`
}

func (a shellAnswer) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}

// ---- the child's world --------------------------------------------------------------------------

// shellEnv is the whole environment a typed command gets. Built from a list of names, never by
// removing the ones that look secret: a variable nobody thought of is not passed.
func shellEnv() []string {
	keep := map[string]bool{}
	for _, k := range []string{"PATH", "PATHEXT", "SYSTEMROOT", "SYSTEMDRIVE", "COMSPEC", "WINDIR", "TEMP", "TMP",
		"USERPROFILE", "HOMEDRIVE", "HOMEPATH", "USERNAME", "USERDOMAIN", "COMPUTERNAME", "LOCALAPPDATA", "APPDATA",
		"PROGRAMFILES", "PROGRAMFILES(X86)", "PROGRAMW6432", "PROGRAMDATA", "COMMONPROGRAMFILES",
		"NUMBER_OF_PROCESSORS", "PROCESSOR_ARCHITECTURE", "OS", "HOME", "USER", "LOGNAME", "LANG", "LC_ALL",
		"TMPDIR", "TERM"} {
		keep[k] = true
	}
	var out []string
	for _, kv := range os.Environ() {
		i := strings.IndexByte(kv, '=')
		if i <= 0 {
			continue // Windows keeps per-drive cwd entries named "=C:"
		}
		if keep[strings.ToUpper(kv[:i])] {
			out = append(out, kv)
		}
	}
	return append(out, "GIT_OPTIONAL_LOCKS=0", "GIT_TERMINAL_PROMPT=0", "GIT_PAGER=cat",
		"PYTHONDONTWRITEBYTECODE=1", "PYTHONUTF8=1", "PYTHONIOENCODING=utf-8", "NO_COLOR=1")
}

// shellCommandVar is where a bash entry travels. As an argument of `bash -c` it went through the Windows
// command line, and Git Bash's runtime halves a doubled backslash that is not inside double quotes on the
// way in: the gate lexed `cat .\\env` as the file `.\env` (not a secret's name) and bash was handed
// `cat .\env`, which is `.env` (found 2026-10-03 by TestTheLexerSeesTheWordsBashSees, an hour after the
// first placement). An environment variable crosses unchanged, and `eval` reads it exactly as `bash -c`
// would have read it on any other machine -- which is what the lexer models.
const shellCommandVar = "ATLAS_SHELL_COMMAND"

// bashInvocation is the arguments and the extra environment that run a typed entry in bash.
func bashInvocation(command string) (args, env []string) {
	return []string{"-c", `eval "$` + shellCommandVar + `"`}, []string{shellCommandVar + "=" + command}
}

// findBash is Git Bash on Windows -- found by the git on the PATH, never by the word `bash`, which on
// this machine is the WSL stub in WindowsApps -- and bash on the PATH anywhere else.
func findBash() (string, error) {
	if runtime.GOOS != "windows" {
		return exec.LookPath("bash")
	}
	var cands []string
	if g, err := exec.LookPath("git"); err == nil {
		d := filepath.Dir(g)
		for i := 0; i < 3; i++ {
			cands = append(cands, filepath.Join(d, "bin", "bash.exe"))
			d = filepath.Dir(d)
		}
	}
	for _, pf := range []string{os.Getenv("ProgramFiles"), os.Getenv("ProgramW6432"), os.Getenv("ProgramFiles(x86)")} {
		if pf != "" {
			cands = append(cands, filepath.Join(pf, "Git", "bin", "bash.exe"))
		}
	}
	for _, c := range cands {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}
	return "", errors.New("Git Bash is not installed here (Git for Windows ships it); a plain `bash` on this PATH may be the WSL stub, which the door never uses")
}

// ---- keys are silent ----------------------------------------------------------------------------

var shEnvSecretName = regexp.MustCompile(`(?i)KEY|TOKEN|SECRET|PASS|SERVICE|CREDENTIAL|BEARER|AUTH`)

// shellSecretValues are the values to withhold from anything a command prints: every value in the
// ground's .env, and every value of a variable in the door's own environment whose name says it is one.
func shellSecretValues(home string) []string {
	seen := map[string]bool{}
	var vals []string
	add := func(v string) {
		v = strings.TrimSpace(v)
		if len(v) >= 8 && !seen[v] {
			seen[v] = true
			vals = append(vals, v)
		}
	}
	if b, err := os.ReadFile(filepath.Join(home, ".env")); err == nil {
		if len(b) > 256*1024 {
			b = b[:256*1024]
		}
		for _, ln := range strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n") {
			ln = strings.TrimPrefix(strings.TrimSpace(ln), "export ")
			if ln == "" || strings.HasPrefix(ln, "#") {
				continue
			}
			i := strings.IndexByte(ln, '=')
			if i <= 0 {
				continue
			}
			v := strings.TrimSpace(ln[i+1:])
			if len(v) >= 2 && (v[0] == '"' || v[0] == '\'') && v[len(v)-1] == v[0] {
				v = v[1 : len(v)-1]
			}
			add(v)
		}
	}
	for _, kv := range os.Environ() {
		if i := strings.IndexByte(kv, '='); i > 0 && shEnvSecretName.MatchString(kv[:i]) {
			add(kv[i+1:])
		}
	}
	sort.Slice(vals, func(a, b int) bool { return len(vals[a]) > len(vals[b]) })
	return vals
}

func shellRedact(home, s string) string {
	for _, v := range shellSecretValues(home) {
		s = strings.ReplaceAll(s, v, shellWithheld)
	}
	return shKeyRe.ReplaceAllString(s, shellWithheld)
}

// shellShape is what a command printed, as the glass may see it: one terminator, valid text, secrets
// withheld, and no more than max bytes.
func shellShape(home, s string, max int) (string, bool) {
	s = strings.ToValidUTF8(strings.ReplaceAll(s, "\r\n", "\n"), "�")
	s = shellRedact(home, s)
	trunc := false
	if len(s) > max {
		s = strings.ToValidUTF8(s[:max], "")
		trunc = true
	}
	return strings.TrimRight(s, "\n"), trunc
}

// ---- the law gate -------------------------------------------------------------------------------

type lawVerdict struct {
	stamp string
	ok    bool
	note  string
}

var lawCache = struct {
	mu sync.Mutex
	by map[string]lawVerdict
}{by: map[string]lawVerdict{}}

// lawStamp is what the verdict was walked OVER: every file in law/, by name, size and time -- not the
// newest time alone, which a deleted law does not move (the engine's gate learned that on 2026-09-29).
func lawStamp(home string) string {
	entries, err := os.ReadDir(filepath.Join(home, "law"))
	if err != nil {
		return ""
	}
	var b strings.Builder
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if info, err := e.Info(); err == nil {
			fmt.Fprintf(&b, "%s|%d|%d;", e.Name(), info.Size(), info.ModTime().UnixNano())
		}
	}
	return b.String()
}

// shellLawGate is whether the sealed law verifies in this ground. No ledger is not a broken chain: the
// gate says so and runs on the rules alone, as the engine's does.
func shellLawGate(py, home string) (bool, string) {
	if st, err := os.Stat(filepath.Join(home, "law", "chain.jsonl")); err != nil || st.IsDir() {
		return true, "no law ledger in this ground; the gate ran on the rules alone"
	}
	tool := filepath.Join(home, "law", "law.py")
	if st, err := os.Stat(tool); err != nil || st.IsDir() {
		return true, "no law tool in this ground; the gate ran on the rules alone"
	}
	stamp := lawStamp(home)
	lawCache.mu.Lock()
	if c, ok := lawCache.by[home]; ok && c.stamp == stamp {
		lawCache.mu.Unlock()
		return c.ok, c.note
	}
	lawCache.mu.Unlock()
	res := spawn(py, []string{tool, "verify"}, spawnOpts{Dir: home, Timeout: 60 * time.Second,
		CleanEnv: shellEnv(), MaxBytes: 8 * 1024, Grace: 2 * time.Second})
	v := lawVerdict{stamp: stamp}
	switch {
	case res.TimedOut:
		v.note = "the law could not be walked in time"
	case res.Err != nil:
		v.note = strings.TrimSpace(strings.ReplaceAll(res.Combined, "\r\n", "\n"))
		if v.note == "" {
			v.note = "the law could not be walked: " + res.Err.Error()
		}
	default:
		v.ok = true
		v.note = strings.TrimSpace(strings.ReplaceAll(res.Combined, "\r\n", "\n"))
	}
	lawCache.mu.Lock()
	lawCache.by[home] = v
	lawCache.mu.Unlock()
	return v.ok, v.note
}

// ---- the record ---------------------------------------------------------------------------------

func shellRecord(args map[string]any, t tenant.Tenant, c Caller, what, note string) {
	reg := regOf(args)
	if reg == nil {
		return
	}
	reg.record(t.Home, what, Hold{
		ID: fmt.Sprintf("shell_%d", time.Now().UnixMilli()), Tool: "shell_run",
		Caller: callerLabel(c), Project: t.Name, RBAC: rbacState(t, c),
	}, note)
}

func shellTruthy(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case string:
		s := strings.ToLower(strings.TrimSpace(x))
		return s == "true" || s == "1" || s == "yes"
	case float64:
		return x != 0
	}
	return false
}

func waitArg(args map[string]any) time.Duration {
	d := shellDefaultWait
	switch x := args["timeout"].(type) {
	case int:
		if x > 0 {
			d = time.Duration(x) * time.Second
		}
	case float64:
		if x > 0 {
			d = time.Duration(x * float64(time.Second))
		}
	case string:
		if f, err := strconv.ParseFloat(strings.TrimSpace(x), 64); err == nil && f > 0 {
			d = time.Duration(f * float64(time.Second))
		}
	}
	if d > shellMaxWait {
		d = shellMaxWait
	}
	return d
}

// ---- the tool -----------------------------------------------------------------------------------

func toolShellRun(coreCmd string) Fn {
	return func(t tenant.Tenant, args map[string]any) (string, error) {
		c := callerOf(args)
		// Registry.Call has refused every other caller already (ServiceOnly); a tool that is
		// reached any other way refuses for itself.
		if !c.Service {
			return "", errors.New("refused: the shell is the operator's own hand and is offered to his glass alone")
		}
		shell := strings.ToLower(strings.TrimSpace(str(args, "shell")))
		if shell != "bash" && shell != "python" {
			return "", fmt.Errorf("refused: name the shell, bash or python (got %q)", shell)
		}
		py := pythonWord(coreCmd)
		if shell == "python" && shellTruthy(args["reset"]) {
			s := pySessionFor(t.Home)
			s.mu.Lock()
			s.reset()
			s.mu.Unlock()
			return shellAnswer{Shell: shell, State: "ran", Note: "the Python session was ended: its names are gone, and the next entry starts a fresh one"}.String(), nil
		}
		command := strings.TrimRight(strings.TrimLeft(str(args, "command"), "\r\n"), " \t\r\n")
		if strings.TrimSpace(command) == "" {
			return "", errors.New("refused: nothing to run -- type something")
		}
		wait := waitArg(args)
		ans := shellAnswer{Shell: shell, Command: clip(shellRedact(t.Home, command), 500)}
		record := func(what, note string) {
			shellRecord(args, t, c, what, shell+" "+note+": "+ans.Command)
		}
		refuse := func(why ...string) (string, error) {
			ans.State, ans.Class, ans.Why = "refused", string(ShellRefused), why
			record("shell_refused", "refused")
			return ans.String(), nil
		}
		if len(command) > shellMaxCommand {
			return refuse(fmt.Sprintf("longer than %d characters; it is not an entry, it is a file", shellMaxCommand))
		}
		if ok, note := shellLawGate(py, t.Home); !ok {
			return refuse("THE CHAIN: the law does not verify -- " + clip(note, 300) + ". No run proceeds on a law that cannot be trusted")
		}

		// stop is the common end of the gate: a refusal, or a card.
		stop := func(v ShellVerdict) (string, bool) {
			ans.Class, ans.Why = string(v.Class), v.Why
			switch {
			case v.Class == ShellRefused:
				out, _ := refuse(v.Why...)
				return out, true
			case v.Class == ShellWrite && c.Hold == "":
				reg := regOf(args)
				self, found := Tool{}, false
				if reg != nil {
					self, found = reg.Get("shell_run")
				}
				if !found {
					out, _ := refuse("the hold queue is not reachable from here, so this cannot wait for his hand")
					return out, true
				}
				keep := copyArgs(args)
				delete(keep, registryKey)
				h := reg.park(t, self, keep, c)
				ans.State, ans.Hold = "held", h.ID
				record("shell_held", "held ("+strings.Join(v.Why, "; ")+")")
				return ans.String(), true
			}
			ans.Approved = v.Class == ShellWrite // a run that is not a read ran on his word
			return "", false
		}

		if shell == "bash" {
			if out, done := stop(JudgeBash(command, t.Home)); done {
				return out, nil
			}
			bash, err := findBash()
			if err != nil {
				return refuse(err.Error())
			}
			began := time.Now()
			bashArgs, bashEnv := bashInvocation(command)
			res := spawn(bash, bashArgs, spawnOpts{Dir: t.Home, Timeout: wait, CleanEnv: shellEnv(), Env: bashEnv,
				MaxBytes: shellMaxBytes, Grace: 2 * time.Second, KillTree: true})
			exit := 0
			var ee *exec.ExitError
			switch {
			case res.TimedOut:
				exit = -1
				ans.TimedOut = true
				ans.Note = fmt.Sprintf("did not finish inside %d s and was ended", int(wait.Seconds()))
			case res.Err == nil:
			case errors.As(res.Err, &ee):
				exit = ee.ExitCode()
			default:
				exit = -1
				ans.Note = res.Err.Error()
			}
			var trunc bool
			ans.Output, trunc = shellShape(t.Home, res.Combined, shellMaxBytes)
			ans.Truncated = trunc || res.Truncated
			ans.State, ans.Exit, ans.Ms = "ran", &exit, time.Since(began).Milliseconds()
			record("shell_ran", fmt.Sprintf("%s exit=%d %dms", ans.Class, exit, ans.Ms))
			return ans.String(), nil
		}

		// Python: the text is refused by name first, then the session reads it.
		if v := JudgePythonText(command, t.Home); v.Class == ShellRefused {
			if out, done := stop(v); done {
				return out, nil
			}
		}
		s := pySessionFor(t.Home)
		s.mu.Lock()
		defer s.mu.Unlock()
		lost, err := s.start(py, t.Home)
		if err != nil {
			return refuse(err.Error())
		}
		if lost {
			ans.Note = "the earlier Python session had ended, so this one starts fresh: its names are gone"
		}
		v, err := s.judge(command)
		if err != nil {
			if errors.Is(err, errPyTimeout) {
				return refuse("the Python session did not answer; it was ended")
			}
			return refuse("the Python session could not read the entry: " + err.Error())
		}
		if out, done := stop(v); done {
			return out, nil
		}
		began := time.Now()
		r, err := s.run(command, wait)
		ans.State, ans.Ms = "ran", time.Since(began).Milliseconds()
		exit := 0
		switch {
		case errors.Is(err, errPyTimeout):
			exit = -1
			ans.TimedOut = true
			ans.Note = fmt.Sprintf("ran past %d s: the session was ended, and its names are gone", int(wait.Seconds()))
		case err != nil:
			exit = -1
			ans.Note = err.Error()
		default:
			if !r.OK {
				exit = 1
			}
			text := r.Out + r.Err
			if strings.TrimSpace(r.Stray) != "" {
				if text != "" && !strings.HasSuffix(text, "\n") {
					text += "\n"
				}
				text += r.Stray
			}
			var trunc bool
			ans.Output, trunc = shellShape(t.Home, text, shellMaxBytes)
			ans.Truncated = trunc || r.Truncated
		}
		ans.Exit = &exit
		record("shell_ran", fmt.Sprintf("%s exit=%d %dms", ans.Class, exit, ans.Ms))
		return ans.String(), nil
	}
}
