package tools

// Aider behind the glass (WHAT'S LEFT H12 and B20, 2026-10-04). Each stroke here names the thing that would
// let Aider reach what it was not given, or let a model's edit land where a model may not write, and runs the
// door's own judgement against a ground it builds in a temp folder. Aider itself is a stub standing where the
// real program stands (aiderStart): the real one is 640 MB and a model, and what is measured here is the DOOR --
// what is copied in, what is written back and where, what is refused, what is recorded. The wall that stands
// around the real program (aiderguard.go) is measured against a REAL Python, with probes that try to leave it.
//
// HERMETIC: every ground is t.TempDir(); the "rack" is a server on loopback that only answers /api/tags; no
// model is asked and no network is reached. Strokes that need a real interpreter skip when the machine has
// none, and say so.

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"atlas/line/internal/tenant"
)

const aiderCoderDoc = "## Expert Coder\n- **Model Target:** %s\n- **Wakes On:** technical\n- **Context:** %d\n- **Timeout:** 600\n- **System Prompt:**\nYou are the smith.\n"

// aiderWarms is what the fake rack was asked to load, in order: the model, and the window it was asked at.
var aiderWarms struct {
	sync.Mutex
	asked []map[string]any
	fail  bool
}

func aiderWarmsSeen() []map[string]any {
	aiderWarms.Lock()
	defer aiderWarms.Unlock()
	return append([]map[string]any(nil), aiderWarms.asked...)
}

// aiderRackServer is a rack that holds the named models and nothing else. It answers /api/tags, and /api/generate (the warm-up) with
// nothing, remembering what it was asked.
func aiderRackServer(t *testing.T, models ...string) {
	t.Helper()
	aiderWarms.Lock()
	aiderWarms.asked, aiderWarms.fail = nil, false
	aiderWarms.Unlock()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/generate" {
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			aiderWarms.Lock()
			defer aiderWarms.Unlock()
			aiderWarms.asked = append(aiderWarms.asked, body)
			if aiderWarms.fail {
				http.Error(w, "the model could not be loaded", http.StatusInternalServerError)
				return
			}
			_, _ = w.Write([]byte(`{"done":true}`))
			return
		}
		if r.URL.Path != "/api/tags" {
			http.NotFound(w, r)
			return
		}
		var rows []map[string]any
		for _, m := range models {
			rows = append(rows, map[string]any{"name": m, "size": 9000000000, "details": map[string]any{"family": "qwen2"}})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"models": rows})
	}))
	t.Cleanup(srv.Close)
	t.Setenv("OLLAMA_HOST", srv.URL)
}

// aiderRealVenv makes a real, empty virtual environment where the door expects Aider's, so the parts of the
// door that start the venv's interpreter (the parse gate, the wall) have one to start.
func aiderRealVenv(t *testing.T, pl aiderPlace) {
	t.Helper()
	needPython(t)
	cmd := exec.Command("python", "-m", "venv", "--without-pip", filepath.Join(pl.Root, "venv"))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("this python cannot make a venv here: %v\n%s", err, out)
	}
	if !aiderInstalled(pl) {
		t.Fatalf("the venv has no interpreter at %s", pl.Python)
	}
}

// aiderWorld is a ground that is a repository standing on a line of work, with a coding seat, an Aider the door
// believes is installed, and a rack that has the seat's model. real says whether the venv is a real one.
func aiderWorld(t *testing.T, real bool, files map[string]string) (*Registry, *tenant.Registry, tenant.Tenant) {
	t.Helper()
	return aiderBuild(t, real, true, files)
}

func aiderBuild(t *testing.T, real, withGit bool, files map[string]string) (*Registry, *tenant.Registry, tenant.Tenant) {
	t.Helper()
	all := map[string]string{
		"agents/expert_coder.md": fmt.Sprintf(aiderCoderDoc, "qwen2.5-coder:14b", 8192),
		"notes.txt":              "one\ntwo\nthree\n",
	}
	for k, v := range files {
		all[k] = v
	}
	tn := world(t, all)
	if withGit {
		mustGit(t, tn, "init", "-b", "main")
		mustGit(t, tn, "config", "user.name", "prove")
		mustGit(t, tn, "config", "user.email", "prove@localhost")
		mustGit(t, tn, "config", "core.autocrlf", "false")
		mustGit(t, tn, "add", "-A")
		mustGit(t, tn, "commit", "-m", "the first save")
		mustGit(t, tn, "checkout", "-b", "line-of-work")
	}
	pl := aiderPlaceOf(tn.Home)
	if real {
		aiderRealVenv(t, pl)
	} else {
		if err := os.MkdirAll(filepath.Dir(pl.Python), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(pl.Python, nil, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	aiderRackServer(t, "qwen2.5-coder:14b")
	tr := tenant.NewRegistry()
	if err := tr.Add("probe", tn.Home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("probe"); err != nil {
		t.Fatal(err)
	}
	resolved, err := tr.Resolve("")
	if err != nil {
		t.Fatal(err)
	}
	return Build(tr, Options{HoldWrites: true}), tr, resolved
}

// aiderCall is what the door started Aider with, as the stub reads it back.
type aiderCall struct {
	Dir     string
	Args    []string
	Env     []string
	Edit    []string
	Read    []string
	Message string
	Nonce   string
}

func aiderCallOf(args []string, opts spawnOpts) aiderCall {
	c := aiderCall{Dir: opts.Dir, Args: args, Env: opts.CleanEnv}
	for _, kv := range opts.CleanEnv {
		if strings.HasPrefix(kv, "ATLAS_AIDER_NONCE=") {
			c.Nonce = strings.TrimPrefix(kv, "ATLAS_AIDER_NONCE=")
		}
	}
	for i := 0; i < len(args)-1; i++ {
		switch args[i] {
		case "--file":
			c.Edit = append(c.Edit, args[i+1])
		case "--read":
			c.Read = append(c.Read, args[i+1])
		case "--message-file":
			b, _ := os.ReadFile(args[i+1])
			c.Message = strings.TrimSpace(string(b))
		}
	}
	return c
}

// armedStderr is what the wall writes when it is up.
func (c aiderCall) armedStderr() string { return aiderGuardMark + " " + c.Nonce + " armed\n" }

// said is a finished run: the wall armed, Aider's words on stdout, exit 0.
func (c aiderCall) said(words string) spawnResult {
	return spawnResult{Stdout: words, Stderr: c.armedStderr()}
}

// stubAider stands fn where Aider stands, and counts the times it was started.
func stubAider(t *testing.T, fn func(c aiderCall) spawnResult) *int {
	t.Helper()
	old := aiderStart
	calls := 0
	aiderStart = func(py string, args []string, opts spawnOpts) spawnResult {
		calls++
		return fn(aiderCallOf(args, opts))
	}
	t.Cleanup(func() { aiderStart = old })
	return &calls
}

// editScratch writes a file in the run's scratch folder, as Aider would, creating its folder.
func (c aiderCall) editScratch(t *testing.T, rel, body string) {
	t.Helper()
	p := filepath.Join(c.Dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func callAider(t *testing.T, reg *Registry, tr *tenant.Registry, tool string, args map[string]any) aiderAnswer {
	t.Helper()
	out, err := reg.Call(tr, tool, args, glass)
	if err != nil {
		t.Fatalf("%s errored: %v", tool, err)
	}
	var a aiderAnswer
	if err := json.Unmarshal([]byte(out), &a); err != nil {
		t.Fatalf("%s did not answer in its shape: %v\n%s", tool, err, out)
	}
	return a
}

func aiderAsk(files, message string) map[string]any {
	return map[string]any{"files": files, "message": message}
}

func readFile(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func hasAnyOf(list []string, frag string) bool {
	for _, s := range list {
		if strings.Contains(s, frag) {
			return true
		}
	}
	return false
}

// A RUN THAT WORKS: the named file is copied in, Aider's edit is written back, and the answer says what changed.
func TestAiderEditsTheNamedFileAndAnswersInOneShape(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	var seen aiderCall
	calls := stubAider(t, func(c aiderCall) spawnResult {
		seen = c
		if got := readFile(t, filepath.Join(c.Dir, "notes.txt")); got != "one\ntwo\nthree\n" {
			t.Errorf("Aider was handed %q, not the file as it stands", got)
		}
		c.editScratch(t, "notes.txt", "one\nTWO\nthree\n")
		return c.said("notes.txt\n```\n<<<<<<< SEARCH\ntwo\n=======\nTWO\n>>>>>>> REPLACE\n```\n\nTokens: 2.0k sent, 39 received.\nApplied edit to notes.txt\n")
	})
	a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "make the second line shout"))
	if *calls != 1 || a.State != "ran" || a.Exit == nil || *a.Exit != 0 {
		t.Fatalf("the run: %+v (started %d times)", a, *calls)
	}
	if got := readFile(t, filepath.Join(tn.Home, "notes.txt")); got != "one\nTWO\nthree\n" {
		t.Fatalf("the edit was not written back: %q", got)
	}
	if a.Line != "line-of-work" || a.Model != "qwen2.5-coder:14b" || !aiderRunRe.MatchString(a.Run) {
		t.Fatalf("line/model/run: %+v", a)
	}
	if len(a.Changed) != 1 || a.Changed[0].File != "notes.txt" || a.Changed[0].Added != 1 || a.Changed[0].Removed != 1 {
		t.Fatalf("what changed: %+v", a.Changed)
	}
	if !strings.Contains(a.Diff, "-two") || !strings.Contains(a.Diff, "+TWO") || !strings.Contains(a.Diff, "a/notes.txt") || strings.Contains(a.Diff, "before/") {
		t.Fatalf("the diff: %q", a.Diff)
	}
	if a.Tokens != "2.0k sent, 39 received." || !strings.Contains(a.Said, "Applied edit") {
		t.Fatalf("what Aider said: tokens %q said %q", a.Tokens, a.Said)
	}
	if !strings.Contains(a.Note, "line of work `line-of-work`") || !strings.Contains(a.Note, "unsaved") {
		t.Fatalf("the note must say the work is unsaved on its line: %q", a.Note)
	}
	// the answer's keys are the contract the glass reads
	out, _ := reg.Call(tr, "aider_status", map[string]any{}, agent)
	var st map[string]any
	if err := json.Unmarshal([]byte(out), &st); err != nil {
		t.Fatalf("status: %v\n%s", err, out)
	}
	for _, k := range []string{"installed", "folder", "model", "context", "rack", "on_rack", "busy", "line", "on_line", "max_files", "max_context", "budget_kb", "runs"} {
		if _, ok := st[k]; !ok {
			t.Errorf("aider_status lacks %q: %s", k, out)
		}
	}
	if runs, _ := st["runs"].([]any); len(runs) != 1 {
		t.Fatalf("the status should list the run: %s", out)
	}
	allowed := map[string]bool{"state": true, "class": true, "why": true, "run": true, "model": true, "line": true, "changed": true,
		"diff": true, "withheld": true, "said": true, "tokens": true, "ignored": true, "guarded": true, "exit": true, "ms": true,
		"timed_out": true, "truncated": true, "note": true}
	raw, _ := json.Marshal(a)
	var keys map[string]any
	_ = json.Unmarshal(raw, &keys)
	for k := range keys {
		if !allowed[k] {
			t.Errorf("aider_run answers with a key the glass does not know: %q", k)
		}
	}
	_ = seen
}

// THE ARGUMENTS ARE THE ONES THE PROBE PROVED. A flag list is a convention, and a convention is not a wire; this is
// the stroke that makes a dropped flag loud. Each one is here because the probe (Aider 0.86.2, 2026-10-04) showed
// what happens without it.
func TestAiderIsStartedWithTheFlagsThatKeepItOnTheLeash(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	stubAider(t, func(c aiderCall) spawnResult {
		joined := " " + strings.Join(c.Args, " ") + " "
		for _, want := range []string{" --no-git ", " --yes-always ", " --no-analytics ", " --analytics-disable ", " --no-check-update ",
			" --no-suggest-shell-commands ", " --no-auto-lint ", " --no-auto-test ", " --map-tokens 0 ", " --no-detect-urls ",
			" --no-auto-commits ", " --no-dirty-commits ", " --edit-format diff ", " --line-endings lf ", " --encoding utf-8 ",
			" --model ollama_chat/qwen2.5-coder:14b ", " --message-file ", " --env-file ", " --model-settings-file ", " --model-metadata-file ",
			" --chat-history-file ", " --llm-history-file ", " --no-stream ", " --no-pretty ", " --no-fancy-input ", " --file notes.txt ",
			" --read ctx.txt "} {
			if !strings.Contains(joined, want) {
				t.Errorf("Aider was started without %q:\n%s", strings.TrimSpace(want), joined)
			}
		}
		// the program and the wall: -I -B -X utf8 -c <guard> aider, in that order, before any Aider flag
		if len(c.Args) < 8 || strings.Join(c.Args[:4], " ") != "-I -B -X utf8" || c.Args[4] != "-c" || c.Args[5] != aiderGuardSource || c.Args[6] != "aider" {
			t.Errorf("Aider was not started behind the wall: %.120v", c.Args)
		}
		// every path it was handed is inside aider/ (the run's own folder or the work folder)
		for i, a := range c.Args {
			if strings.HasPrefix(a, "--") && a != "--file" && strings.HasSuffix(a, "-file") && i+1 < len(c.Args) {
				if rel, err := filepath.Rel(aiderPlaceOf(tn.Home).Root, c.Args[i+1]); err != nil || strings.HasPrefix(rel, "..") {
					t.Errorf("%s names a path outside aider/: %s", a, c.Args[i+1])
				}
			}
		}
		return c.said("")
	})
	write(t, tn.Home, "ctx.txt", "context\n")
	a := callAider(t, reg, tr, "aider_run", map[string]any{"files": "notes.txt", "context": "ctx.txt", "message": "nothing"})
	if a.State != "ran" {
		t.Fatalf("%+v", a)
	}
}

// AIDER CANNOT FIND THE WORLD'S OWN REPOSITORY FROM ITS SCRATCH FOLDER. It looks up from where it stands for a git repository, even
// with --no-git, and the scratch folder is inside the world: without a repository of its own there, it reached the world's `.git`
// and the wall refused it, which killed the first real run (2026-10-04). The stub asks git what repository the scratch folder is in.
func TestAiderScratchFolderIsItsOwnRepository(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	var top string
	stubAider(t, func(c aiderCall) spawnResult {
		res := spawn("git", []string{"rev-parse", "--show-toplevel"}, spawnOpts{Dir: c.Dir, Timeout: 20 * time.Second, Env: gitEnv()})
		top = strings.TrimSpace(strings.ReplaceAll(res.Stdout, "/", string(filepath.Separator)))
		c.editScratch(t, "notes.txt", "one\nTWO\nthree\n")
		return c.said("")
	})
	a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x"))
	if a.State != "ran" || len(a.Changed) != 1 {
		t.Fatalf("%+v", a)
	}
	// compared by its tail, so a short (8.3) spelling of the temp folder on one side and a long one on the other cannot matter
	if want := "/aider/work/runs/" + a.Run + "/scratch"; !strings.HasSuffix(strings.ToLower(filepath.ToSlash(top)), want) {
		t.Fatalf("Aider's search for a repository from its scratch folder ends at %q, not at the scratch folder (...%s)", top, want)
	}
	_ = tn
	if hasAnyOf(a.Ignored, ".git") {
		t.Fatalf("the scratch folder's own repository was reported as something Aider made: %v", a.Ignored)
	}
}

// THE MODEL IS LOADED BEFORE AIDER IS STARTED, at the window Aider will ask for, and a rack that will not load it is a refusal that says so. A
// first measured run hung four minutes inside Aider's own timeout because its request reached Ollama while it was unloading the model.
func TestAiderLoadsTheModelAtTheSeatsWindowBeforeItStartsAider(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	var warmedBefore int
	calls := stubAider(t, func(c aiderCall) spawnResult {
		warmedBefore = len(aiderWarmsSeen())
		c.editScratch(t, "notes.txt", "one\nTWO\nthree\n")
		return c.said("")
	})
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x")); a.State != "ran" {
		t.Fatalf("%+v", a)
	}
	warms := aiderWarmsSeen()
	if warmedBefore != 1 || len(warms) != 1 || warms[0]["model"] != "qwen2.5-coder:14b" || warms[0]["keep_alive"] != "10m" {
		t.Fatalf("the model was not loaded exactly once, before Aider: %d before, %v", warmedBefore, warms)
	}
	if opts, _ := warms[0]["options"].(map[string]any); opts["num_ctx"] != float64(8192) {
		t.Fatalf("the model was loaded at a window other than the one Aider asks for: %v", warms[0])
	}
	// the seat's window moves, and so does the warm-up
	write(t, tn.Home, "agents/expert_coder.md", fmt.Sprintf(aiderCoderDoc, "qwen2.5-coder:14b", 16384))
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "again")); a.State != "ran" {
		t.Fatalf("%+v", a)
	}
	if w := aiderWarmsSeen(); len(w) != 2 || func() any { o, _ := w[1]["options"].(map[string]any); return o["num_ctx"] }() != float64(16384) {
		t.Fatalf("the warm-up did not follow the seat's window: %v", w)
	}
	// a rack that will not load it: refused by name, and Aider is never started
	aiderWarms.Lock()
	aiderWarms.fail = true
	aiderWarms.Unlock()
	before := *calls
	a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "a third time"))
	if a.State != "refused" || !hasAnyOf(a.Why, "did not load qwen2.5-coder:14b") || *calls != before {
		t.Fatalf("a rack that would not load the model: %+v (Aider started %d times, was %d)", a, *calls, before)
	}
}

// WHAT AIDER IS GIVEN IS BUILT, NOT INHERITED (RULE 7). The door holds the service wire and keys in its own
// environment; none of them may reach the child, nor a proxy or a path variable that could point it elsewhere.
func TestAiderGetsABuiltEnvironmentWithNoKeyInIt(t *testing.T) {
	t.Setenv("ATLAS_SERVICE_KEY", "svc_abcdefghijklmnopqrstuv")
	t.Setenv("OPENAI_API_KEY", "sk-abcdefghijklmnopqrstuvwxyz")
	t.Setenv("ANTHROPIC_API_KEY", "ant_abcdefghijklmnopqrstuv")
	t.Setenv("HTTPS_PROXY", "http://proxy.example:3128")
	t.Setenv("PYTHONPATH", "elsewhere")
	reg, tr, tn := aiderWorld(t, false, nil)
	pl := aiderPlaceOf(tn.Home)
	stubAider(t, func(c aiderCall) spawnResult {
		all := strings.Join(c.Env, "\n")
		for _, bad := range []string{"svc_abcdefghijkl", "sk-abcdefghijkl", "ant_abcdefghijkl", "HTTPS_PROXY", "HTTP_PROXY", "PYTHONPATH", "OPENAI", "ANTHROPIC"} {
			if strings.Contains(strings.ToUpper(all), strings.ToUpper(bad)) {
				t.Errorf("the child's environment carries %q", bad)
			}
		}
		want := map[string]string{"HOME": pl.Home, "USERPROFILE": pl.Home, "TEMP": pl.Tmp, "TMP": pl.Tmp, "OLLAMA_API_BASE": os.Getenv("OLLAMA_HOST"),
			"LITELLM_LOCAL_MODEL_COST_MAP": "True", "GIT_PYTHON_REFRESH": "quiet", "PYTHONDONTWRITEBYTECODE": "1"}
		for k, v := range want {
			found := false
			for _, kv := range c.Env {
				if kv == k+"="+v {
					found = true
				}
			}
			if !found {
				t.Errorf("the child's environment lacks %s=%s", k, v)
			}
		}
		// the wall it carries names this run's scratch as a place to write, and neither the ground nor the undo copies
		var wall aiderWall
		for _, kv := range c.Env {
			if strings.HasPrefix(kv, "ATLAS_AIDER_WALL=") {
				if err := json.Unmarshal([]byte(strings.TrimPrefix(kv, "ATLAS_AIDER_WALL=")), &wall); err != nil {
					t.Errorf("the wall is not JSON: %v", err)
				}
			}
		}
		if wall.Ground != tn.Home || wall.Aider != pl.Root || len(wall.WriteDirs) != 3 || len(wall.WriteFiles) != 3 {
			t.Errorf("the wall: %+v", wall)
		}
		for _, d := range append(append([]string{}, wall.WriteDirs...), wall.WriteFiles...) {
			if rel, err := filepath.Rel(tn.Home, d); err != nil || !strings.HasPrefix(filepath.ToSlash(rel), "aider/work/") || strings.Contains(filepath.ToSlash(rel), "/before") {
				t.Errorf("the wall lets Aider write at %s", d)
			}
		}
		return c.said("")
	})
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "nothing")); a.State != "ran" {
		t.Fatalf("%+v", a)
	}
}

// NO AGENT IS EVER GIVEN AIDER, and a tool reached any other way still refuses for itself.
func TestAiderIsHisAloneByEveryRoute(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	calls := stubAider(t, func(c aiderCall) spawnResult {
		c.editScratch(t, "notes.txt", "stolen\n")
		return c.said("")
	})
	for _, who := range []Caller{agent, {Name: "k-agent"}, {}, {Name: "operator via hold hold_1"}} {
		for tool, args := range map[string]map[string]any{"aider_run": aiderAsk("notes.txt", "x"), "aider_undo": {"run": "a20261004-120000-abcd"}} {
			if _, err := reg.Call(tr, tool, args, who); err == nil || !strings.Contains(err.Error(), "operator's own hand") {
				t.Fatalf("%q was handed %s: %v", who.Name, tool, err)
			}
		}
		for _, name := range []string{"aider_run", "aider_undo"} {
			tool, _ := reg.Get(name)
			if out, err := tool.Fn(tn, map[string]any{"message": "x", "files": "notes.txt", "run": "a20261004-120000-abcd", CallerKey: who}); err == nil || out != "" || !strings.Contains(err.Error(), "operator's own hand") {
				t.Fatalf("%s answered %q directly: %q / %v", name, who.Name, out, err)
			}
		}
	}
	if *calls != 0 || readFile(t, filepath.Join(tn.Home, "notes.txt")) != "one\ntwo\nthree\n" {
		t.Fatalf("a refused caller reached Aider (%d starts)", *calls)
	}
	for _, name := range []string{"aider_run", "aider_undo"} {
		tool, _ := reg.Get(name)
		if !tool.ServiceOnly || !tool.Writes {
			t.Errorf("%s must be declared a writing, service-only tool: %+v", name, tool)
		}
	}
	if tool, _ := reg.Get("aider_status"); tool.ServiceOnly || tool.Writes {
		t.Errorf("aider_status is the page's reading side: %+v", tool)
	}
}

// A WRITE LANDS ONLY ON A LINE OF WORK. The main line is his (RULE 6): main, master, a detached head, no repository,
// and a nested repository on its own main line all refuse, by name, and Aider is never started.
func TestAiderWritesOnlyOnALineOfWork(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, map[string]string{"sub/inner.txt": "inner\n"})
	calls := stubAider(t, func(c aiderCall) spawnResult {
		c.editScratch(t, "notes.txt", "EDITED\n")
		return c.said("")
	})
	ask := aiderAsk("notes.txt", "x")
	for _, branch := range []string{"main", "master"} {
		mustGit(t, tn, "checkout", "-B", branch)
		a := callAider(t, reg, tr, "aider_run", ask)
		if a.State != "refused" || !hasAnyOf(a.Why, "main line is the operator's") || !hasAnyOf(a.Why, "`"+branch+"`") {
			t.Fatalf("on %s: %+v", branch, a)
		}
	}
	mustGit(t, tn, "checkout", "--detach")
	if a := callAider(t, reg, tr, "aider_run", ask); a.State != "refused" || !hasAnyOf(a.Why, "detached") {
		t.Fatalf("detached: %+v", a)
	}
	mustGit(t, tn, "checkout", "-B", "line-of-work")
	// the NEAREST repository decides: a folder with its own .git on its own main line refuses even under a root on a line
	inner := filepath.Join(tn.Home, "sub")
	for _, args := range [][]string{{"init", "-b", "main"}, {"config", "user.name", "prove"}, {"config", "user.email", "p@l"}, {"add", "-A"}, {"commit", "-m", "x"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = inner
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("sub/inner.txt", "x")); a.State != "refused" || !hasAnyOf(a.Why, "`sub`") {
		t.Fatalf("a nested repository on main: %+v", a)
	}
	if *calls != 0 || readFile(t, filepath.Join(tn.Home, "notes.txt")) != "one\ntwo\nthree\n" {
		t.Fatalf("Aider ran, or the file moved, on a main line (%d starts)", *calls)
	}
	// and on the line it writes
	if a := callAider(t, reg, tr, "aider_run", ask); a.State != "ran" || len(a.Changed) != 1 {
		t.Fatalf("on a line of work: %+v", a)
	}
	// a ground that is not a repository at all
	reg2, tr2, _ := aiderBuild(t, false, false, nil)
	if a := callAider(t, reg2, tr2, "aider_run", ask); a.State != "refused" || !hasAnyOf(a.Why, "not under version control") {
		t.Fatalf("no repository: %+v", a)
	}
}

// A MODEL NEVER WRITES WHAT THE SEATS NEVER WRITE, and is never shown what no hand reads. By name; no approval
// lifts it; Aider is never started.
func TestAiderRefusesWhatTheSeatsAreRefused(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	calls := stubAider(t, func(c aiderCall) spawnResult { return c.said("") })
	for _, rel := range []string{
		"law/LAW_004.md", "agents/expert_coder.md", "skills/x.py", "sessions/sessions.jsonl", "logs/pack.md", "index/vectors.db",
		"state/holds.jsonl", "flows/coder-tree.json", "memory/x.md", "projects/p/index.html", "agent_workspace/x.py", "foundation/05.md",
		"bin/tool.txt", "worlds/someone/notes.txt", "aider/work/x.txt", "CLAUDE.md", "claude.md", ".gitignore", ".gitattributes",
		"pipelines.md", "commands.md", "index_roots.txt", "memory.md", "SEAT_LOG.md", "BUILDMAP.md", "tests/last_run.json",
		"tests/run_history.jsonl", "tool.exe", "lib.dll", "x.pyc", ".git/config", "sub/.git/HEAD", ".env", ".env.local", "id_rsa",
		"vault/notes.txt", "a.client.md", "../outside.txt", "..\\outside.txt", "C:/Windows/x.txt", "/etc/passwd", "~/x.txt",
	} {
		a := callAider(t, reg, tr, "aider_run", aiderAsk(rel, "change it"))
		if a.State != "refused" || a.Class != "refused" || len(a.Why) == 0 {
			t.Errorf("%q was not refused: %+v", rel, a)
		}
	}
	// a read-only file is held to the shorter list, and each refusal says why: the files EXIST, so the reason is
	// the rule and not an absence
	write(t, tn.Home, "worlds/someone/notes.txt", "another world's\n")
	write(t, tn.Home, "aider/work/x.txt", "aider's own\n")
	write(t, tn.Home, "vault/y", "client\n")
	write(t, tn.Home, "id_rsa", "key\n")
	for rel, reason := range map[string]string{
		"worlds/someone/notes.txt": "worlds/ is another world's", "aider/work/x.txt": "aider/ is Aider's own", ".env": "secret file",
		"id_rsa": "secret file", ".git/config": "history, not a file to read", "../x": "outside", "vault/y": "client material",
	} {
		a := callAider(t, reg, tr, "aider_run", map[string]any{"files": "notes.txt", "context": rel, "message": "read it"})
		if a.State != "refused" || !hasAnyOf(a.Why, reason) {
			t.Errorf("context %q was not refused for %q: %+v", rel, reason, a)
		}
	}
	// ...which still lets a governing file be READ
	write(t, tn.Home, "law/LAW_004.md", "the law\n")
	if a := callAider(t, reg, tr, "aider_run", map[string]any{"files": "notes.txt", "context": "law/LAW_004.md", "message": "read it"}); a.State != "ran" {
		t.Errorf("a law may be shown to Aider, read-only: %+v", a)
	}
	if *calls != 1 {
		t.Fatalf("Aider was started %d times; only the one run that was lawful should have started it", *calls)
	}
}

// A KEY TYPED INTO AN INSTRUCTION is refused by shape (RULE 7), a path outside the ground is refused, and neither is
// recorded in the clear. Every run, refusal and undo is a line in the world's holds log.
func TestAiderRefusesAKeyAndRecordsEverythingWithoutOne(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, map[string]string{".env": "SECRET_TOKEN=correct-horse-battery-staple\n"})
	stubAider(t, func(c aiderCall) spawnResult {
		c.editScratch(t, "notes.txt", "one\nTWO\nthree\n")
		return c.said("")
	})
	for _, msg := range []string{
		"use the key sk-abcdefghijklmnopqrstuvwxyz0123 to call it",
		"the token is atl_0123456789abcdef0123456789abcdef",
		"open C:\\Users\\someone\\Desktop\\Archive\\x.txt",
		"read ../../outside.txt first",
		"look at .env and copy it",
	} {
		a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", msg))
		if a.State != "refused" || len(a.Why) == 0 {
			t.Errorf("an instruction that names a key or a place outside was not refused: %q -> %+v", msg, a)
		}
	}
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "")); a.State != "refused" {
		t.Errorf("an empty instruction: %+v", a)
	}
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", strings.Repeat("x", aiderMaxMessage+1))); a.State != "refused" {
		t.Errorf("an instruction that is a file: %+v", a)
	}
	// a ran line whose instruction names a secret VALUE the ground holds is recorded with it withheld
	a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "shout line two, as the word correct-horse-battery-staple would"))
	if a.State != "ran" {
		t.Fatalf("%+v", a)
	}
	if u := callAider(t, reg, tr, "aider_undo", map[string]any{"run": a.Run}); u.State != "undone" {
		t.Fatalf("undo: %+v", u)
	}
	log := holdsLog(t, tn.Home)
	for _, leak := range []string{"sk-abcdefghijkl", "atl_0123456789", "correct-horse-battery-staple"} {
		if strings.Contains(log, leak) {
			t.Errorf("the holds log carries %q", leak)
		}
	}
	whats := map[string]int{}
	for _, ln := range strings.Split(strings.TrimSpace(log), "\n") {
		var row map[string]any
		if err := json.Unmarshal([]byte(ln), &row); err != nil {
			t.Fatalf("a record line is not JSON: %q", ln)
		}
		if tool, _ := row["tool"].(string); strings.HasPrefix(tool, "aider_") {
			whats[row["what"].(string)]++
		}
	}
	if whats["aider_refused"] < 7 || whats["aider_ran"] != 1 || whats["aider_undone"] != 1 {
		t.Fatalf("what the log recorded: %v", whats)
	}
}

// THE FILE'S OWN TERMINATOR comes back (the ruling of 2026-09-03): Aider writes whatever its platform writes, and
// the door puts each file's own back. A MIXED file is refused up front; a NEW file takes its nearest sibling's.
func TestAiderWritesEachFileInItsOwnTerminator(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, map[string]string{
		"lf.txt": "a\nb\nc\n", "crlf.txt": "a\r\nb\r\nc\r\n", "mixed.txt": "a\r\nb\nc\n",
		"sub/old.txt": "x\r\ny\r\n",
	})
	stubAider(t, func(c aiderCall) spawnResult {
		// the worst case: everything written the Windows way, whatever it was given
		c.editScratch(t, "lf.txt", "a\r\nB\r\nc\r\n")
		c.editScratch(t, "crlf.txt", "a\r\nB\r\nc\r\n")
		c.editScratch(t, "sub/new.txt", "fresh\r\nfile\r\n")
		return c.said("")
	})
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("mixed.txt", "x")); a.State != "refused" || !hasAnyOf(a.Why, "MIXED") {
		t.Fatalf("a mixed file: %+v", a)
	}
	a := callAider(t, reg, tr, "aider_run", aiderAsk("lf.txt\ncrlf.txt\nsub/new.txt", "capitalise b"))
	if a.State != "ran" || len(a.Changed) != 3 {
		t.Fatalf("%+v", a)
	}
	if got := readFile(t, filepath.Join(tn.Home, "lf.txt")); got != "a\nB\nc\n" {
		t.Errorf("an LF file came back %q", got)
	}
	if got := readFile(t, filepath.Join(tn.Home, "crlf.txt")); got != "a\r\nB\r\nc\r\n" {
		t.Errorf("a CRLF file came back %q", got)
	}
	if got := readFile(t, filepath.Join(tn.Home, "sub", "new.txt")); got != "fresh\r\nfile\r\n" {
		t.Errorf("a new file beside CRLF siblings came back %q", got)
	}
	// and a new file beside no sibling is LF
	stubAider(t, func(c aiderCall) spawnResult { c.editScratch(t, "alone.md", "x\r\ny\r\n"); return c.said("") })
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("alone.md", "make it")); a.State != "ran" || readFile(t, filepath.Join(tn.Home, "alone.md")) != "x\ny\n" {
		t.Fatalf("a lone new file: %+v / %q", a, readFile(t, filepath.Join(tn.Home, "alone.md")))
	}
	// a new file lands in a folder that exists; a model makes no folder (RULE 8)
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("nowhere/new.txt", "make it")); a.State != "refused" || !hasAnyOf(a.Why, "RULE 8") {
		t.Fatalf("a new folder: %+v", a)
	}
	if _, err := os.Stat(filepath.Join(tn.Home, "nowhere")); err == nil {
		t.Fatal("a folder was made")
	}
}

// THE SCRATCH IS ALL IT CAN CHANGE. Only the files he named are written back; whatever else Aider made is reported and
// dropped, and a read-only file it changed is never written.
func TestAiderOnlyWritesBackTheFilesHeNamed(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, map[string]string{"ctx.txt": "context\n", "other.txt": "other\n"})
	stubAider(t, func(c aiderCall) spawnResult {
		c.editScratch(t, "notes.txt", "one\nTWO\nthree\n")
		c.editScratch(t, "extra.txt", "made by Aider\n")
		c.editScratch(t, "ctx.txt", "TAMPERED\n")
		c.editScratch(t, ".aider.tags.cache.v4/cache.db", "junk")
		return c.said("")
	})
	a := callAider(t, reg, tr, "aider_run", map[string]any{"files": "notes.txt", "context": "ctx.txt", "message": "x"})
	if a.State != "ran" || len(a.Changed) != 1 {
		t.Fatalf("%+v", a)
	}
	if exists(filepath.Join(tn.Home, "extra.txt")) || readFile(t, filepath.Join(tn.Home, "ctx.txt")) != "context\n" || readFile(t, filepath.Join(tn.Home, "other.txt")) != "other\n" {
		t.Fatal("something Aider was not given to change was changed")
	}
	if !hasAnyOf(a.Ignored, "extra.txt") || !hasAnyOf(a.Ignored, "ctx.txt (read-only") || hasAnyOf(a.Ignored, ".aider") {
		t.Fatalf("what was ignored: %v", a.Ignored)
	}
}

// A WALL THAT DID NOT ANNOUNCE ITSELF WAS NOT THERE. A run whose stderr lacks this run's nonce is thrown away whole,
// however well it edited -- and the nonce of another run is no better.
func TestAiderIsBelievedOnlyBehindAWallThatSaidSo(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	edit := func(c aiderCall) { c.editScratch(t, "notes.txt", "one\nTWO\nthree\n") }
	for name, stderr := range map[string]func(c aiderCall) string{
		"silent":      func(aiderCall) string { return "" },
		"wrong nonce": func(aiderCall) string { return aiderGuardMark + " 0123456789abcdef armed\n" },
		"half":        func(c aiderCall) string { return aiderGuardMark + " " + c.Nonce + "\n" },
	} {
		stubAider(t, func(c aiderCall) spawnResult {
			edit(c)
			return spawnResult{Stdout: "Applied edit to notes.txt\n", Stderr: stderr(c)}
		})
		a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x"))
		if a.State != "refused" || !hasAnyOf(a.Why, "THE WALL") || len(a.Changed) != 0 {
			t.Errorf("%s: a run behind no wall was believed: %+v", name, a)
		}
		if readFile(t, filepath.Join(tn.Home, "notes.txt")) != "one\ntwo\nthree\n" {
			t.Fatalf("%s: a run behind no wall wrote", name)
		}
	}
	// what the wall refused rides the answer, so the card says what it stopped
	stubAider(t, func(c aiderCall) spawnResult {
		edit(c)
		return spawnResult{Stdout: "ok\n", Stderr: c.armedStderr() + aiderGuardMark + " refused: socket.getaddrinfo 'raw.githubusercontent.com' -- the estate is local\nsome warning\n"}
	})
	a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x"))
	if a.State != "ran" || !hasAnyOf(a.Guarded, "raw.githubusercontent.com") || !strings.Contains(a.Said, "some warning") || strings.Contains(a.Said, aiderGuardMark) {
		t.Fatalf("the wall's refusal: %+v", a)
	}
}

// NOTHING IS WRITTEN FROM A RUN THAT DID NOT FINISH: a timeout, a non-zero exit, a refused spawn.
func TestAiderWritesNothingFromARunThatDidNotFinish(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	for name, mk := range map[string]func(c aiderCall) spawnResult{
		"timeout": func(c aiderCall) spawnResult {
			r := c.said("half")
			r.TimedOut, r.Err = true, fmt.Errorf("did not finish")
			return r
		},
		"exit 1": func(c aiderCall) spawnResult {
			r := c.said("the model server did not answer")
			r.Err = &exec.ExitError{ProcessState: nil}
			return r
		},
		"spawn failed": func(c aiderCall) spawnResult {
			r := c.said("")
			r.Err = fmt.Errorf("the interpreter could not be started")
			return r
		},
	} {
		stubAider(t, func(c aiderCall) spawnResult {
			c.editScratch(t, "notes.txt", "one\nTWO\nthree\n")
			return mk(c)
		})
		a := callAider(t, reg, tr, "aider_run", map[string]any{"files": "notes.txt", "message": "x", "timeout": 1})
		if a.State != "ran" || len(a.Changed) != 0 || !strings.Contains(a.Note, "nothing was written") && !strings.Contains(a.Note, "could not be started") {
			t.Errorf("%s: %+v", name, a)
		}
		if name == "timeout" && (!a.TimedOut || !strings.Contains(a.Note, "did not finish inside 1 s")) {
			t.Errorf("a timeout must say so, as a timeout: %+v", a)
		}
		if readFile(t, filepath.Join(tn.Home, "notes.txt")) != "one\ntwo\nthree\n" {
			t.Fatalf("%s: an unfinished run wrote", name)
		}
	}
	// a run that finished and changed nothing says so, with Aider's words
	stubAider(t, func(c aiderCall) spawnResult { return c.said("I do not see a change to make.") })
	a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x"))
	if a.State != "ran" || len(a.Changed) != 0 || !strings.Contains(a.Note, "changed nothing") || !strings.Contains(a.Said, "do not see") {
		t.Fatalf("a run that changed nothing: %+v", a)
	}
}

// THE GROUND MUST BE AS IT WAS WHEN AIDER WAS STARTED. An edit made against a file that moved while Aider worked is
// withheld, whole; the operator's own change stands.
func TestAiderWithholdsAnEditToAFileThatMovedWhileItWorked(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, map[string]string{"other.txt": "other\n"})
	stubAider(t, func(c aiderCall) spawnResult {
		c.editScratch(t, "notes.txt", "one\nTWO\nthree\n")
		c.editScratch(t, "other.txt", "OTHER\n")
		write(t, tn.Home, "notes.txt", "the operator, meanwhile\n")
		return c.said("")
	})
	a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt\nother.txt", "x"))
	if a.State != "ran" || !a.Withheld || !strings.Contains(a.Note, "changed while Aider was working") {
		t.Fatalf("%+v", a)
	}
	if readFile(t, filepath.Join(tn.Home, "notes.txt")) != "the operator, meanwhile\n" || readFile(t, filepath.Join(tn.Home, "other.txt")) != "other\n" {
		t.Fatal("a moved file was overwritten, or a half of the run landed")
	}
}

// A PYTHON FILE THAT DOES NOT PARSE IS NOT WRITTEN. The venv's own interpreter reads it, never runs it; the edit is
// shown, the ground is untouched, and every file of the run lands or none does.
func TestAiderWithholdsAPythonFileThatDoesNotParse(t *testing.T) {
	reg, tr, tn := aiderWorld(t, true, map[string]string{"calc.py": "def add(a, b):\n    return a - b\n", "other.txt": "other\n"})
	bad := true
	stubAider(t, func(c aiderCall) spawnResult {
		if bad {
			c.editScratch(t, "calc.py", "def add(a, b)\n    return a + b\n")
		} else {
			c.editScratch(t, "calc.py", "def add(a, b):\n    return a + b\n")
		}
		c.editScratch(t, "other.txt", "OTHER\n")
		return c.said("")
	})
	a := callAider(t, reg, tr, "aider_run", aiderAsk("calc.py\nother.txt", "fix add"))
	if a.State != "ran" || !a.Withheld || !strings.Contains(a.Note, "does not parse") || !strings.Contains(a.Note, "calc.py") || !strings.Contains(a.Diff, "+def add(a, b)") {
		t.Fatalf("a file that does not parse: %+v", a)
	}
	if readFile(t, filepath.Join(tn.Home, "calc.py")) != "def add(a, b):\n    return a - b\n" || readFile(t, filepath.Join(tn.Home, "other.txt")) != "other\n" {
		t.Fatal("part of a withheld run was written")
	}
	bad = false
	a = callAider(t, reg, tr, "aider_run", aiderAsk("calc.py\nother.txt", "fix add"))
	if a.State != "ran" || a.Withheld || len(a.Changed) != 2 || readFile(t, filepath.Join(tn.Home, "calc.py")) != "def add(a, b):\n    return a + b\n" {
		t.Fatalf("a file that parses: %+v", a)
	}
}

// THE UNDO puts the files back exactly, removes what the run made, refuses a run twice, refuses what moved since, and
// takes only a run's own id.
func TestAiderUndoPutsBackExactlyAndRefusesWhatMovedSince(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, map[string]string{"crlf.txt": "a\r\nb\r\n", "sub/keep.txt": "keep\n"})
	stubAider(t, func(c aiderCall) spawnResult {
		c.editScratch(t, "crlf.txt", "a\r\nB\r\n")
		c.editScratch(t, "sub/made.txt", "made\n")
		return c.said("")
	})
	before := readFile(t, filepath.Join(tn.Home, "crlf.txt"))
	a := callAider(t, reg, tr, "aider_run", aiderAsk("crlf.txt\nsub/made.txt", "x"))
	if a.State != "ran" || len(a.Changed) != 2 || !exists(filepath.Join(tn.Home, "sub", "made.txt")) {
		t.Fatalf("%+v", a)
	}
	u := callAider(t, reg, tr, "aider_undo", map[string]any{"run": a.Run})
	if u.State != "undone" || len(u.Changed) != 2 || readFile(t, filepath.Join(tn.Home, "crlf.txt")) != before || exists(filepath.Join(tn.Home, "sub", "made.txt")) {
		t.Fatalf("the undo: %+v", u)
	}
	if again := callAider(t, reg, tr, "aider_undo", map[string]any{"run": a.Run}); again.State != "refused" || !hasAnyOf(again.Why, "already taken back") {
		t.Fatalf("a second undo: %+v", again)
	}
	// a run whose file has changed since is not taken back
	b := callAider(t, reg, tr, "aider_run", aiderAsk("crlf.txt", "again"))
	if b.State != "ran" {
		t.Fatalf("%+v", b)
	}
	write(t, tn.Home, "crlf.txt", "the operator moved on\n")
	if u := callAider(t, reg, tr, "aider_undo", map[string]any{"run": b.Run}); u.State != "refused" || !hasAnyOf(u.Why, "not as Aider left it") {
		t.Fatalf("an undo over later work: %+v", u)
	}
	if readFile(t, filepath.Join(tn.Home, "crlf.txt")) != "the operator moved on\n" {
		t.Fatal("an undo lost later work")
	}
	// only a run's own id names a folder: a record that a traversal would reach, and a copy beside it, is not read
	decoy := filepath.Join(aiderPlaceOf(tn.Home).Work, "decoy")
	cur := readFile(t, filepath.Join(tn.Home, "crlf.txt"))
	write(t, decoy, "before/crlf.txt", "a\r\nb\r\n")
	write(t, decoy, "manifest.json", fmt.Sprintf(`{"run":"../decoy","files":[{"file":"crlf.txt","eol":"crlf","before":%q,"after":%q}]}`,
		aiderSum([]byte("a\r\nb\r\n")), aiderSum([]byte(cur))))
	if u := callAider(t, reg, tr, "aider_undo", map[string]any{"run": "../decoy"}); u.State != "refused" || readFile(t, filepath.Join(tn.Home, "crlf.txt")) != cur {
		t.Fatalf("an id that is a path reached a record outside runs/: %+v", u)
	}
	for _, bad := range []string{"", "../../x", "a20261004-120000-abcd/../..", "a20261004-120000-ABCD", "x", "a1", "..\\..\\x"} {
		if u := callAider(t, reg, tr, "aider_undo", map[string]any{"run": bad}); u.State != "refused" {
			t.Errorf("run %q was not refused: %+v", bad, u)
		}
	}
	if u := callAider(t, reg, tr, "aider_undo", map[string]any{"run": "a20261004-120000-abcd"}); u.State != "refused" || !hasAnyOf(u.Why, "no record") {
		t.Errorf("a run that never was: %+v", u)
	}
	// an undo is not made on the main line either
	stubAider(t, func(c aiderCall) spawnResult { c.editScratch(t, "sub/keep.txt", "KEEP\n"); return c.said("") })
	c := callAider(t, reg, tr, "aider_run", aiderAsk("sub/keep.txt", "x"))
	if c.State != "ran" || len(c.Changed) != 1 {
		t.Fatalf("%+v", c)
	}
	mustGit(t, tn, "checkout", "-B", "main")
	if u := callAider(t, reg, tr, "aider_undo", map[string]any{"run": c.Run}); u.State != "refused" || !hasAnyOf(u.Why, "main line") {
		t.Fatalf("an undo on the main line: %+v", u)
	}
	if readFile(t, filepath.Join(tn.Home, "sub", "keep.txt")) != "KEEP\n" {
		t.Fatal("an undo wrote on the main line")
	}
}

// ONE RUN AT A TIME: one model in one memory. A second call while one runs is refused, not queued, and the status says busy.
func TestAiderRunsOneAtATime(t *testing.T) {
	reg, tr, _ := aiderWorld(t, false, nil)
	started, release := make(chan struct{}), make(chan struct{})
	stubAider(t, func(c aiderCall) spawnResult {
		close(started)
		<-release
		return c.said("")
	})
	done := make(chan string, 1)
	go func() {
		out, err := reg.Call(tr, "aider_run", aiderAsk("notes.txt", "first"), glass)
		if err != nil {
			out = "ERROR " + err.Error()
		}
		done <- out
	}()
	select {
	case <-started:
	case <-time.After(20 * time.Second):
		t.Fatal("the first run never started")
	}
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "second")); a.State != "refused" || !hasAnyOf(a.Why, "already going") {
		t.Fatalf("a second run: %+v", a)
	}
	out, _ := reg.Call(tr, "aider_status", map[string]any{}, glass)
	if !strings.Contains(out, `"busy":true`) {
		t.Fatalf("the status does not say busy: %s", out)
	}
	close(release)
	var first aiderAnswer
	if out := <-done; json.Unmarshal([]byte(out), &first) != nil || first.State != "ran" {
		t.Fatalf("the first run: %s", out)
	}
}

// THE WINDOW IS THE SEAT'S. The model and its context window are read off the coding seat's declaration at the time of
// the call; files that cannot fit are refused up front rather than silently cut by the model's server; a model that
// is not on this rack is refused by name (RULE 4) -- a hosted route never gets here.
func TestAiderFollowsTheCodingSeatAndRefusesWhatCannotFit(t *testing.T) {
	big := strings.Repeat("a line of code that is somewhat long\n", 600) // ~22 KB
	reg, tr, tn := aiderWorld(t, false, map[string]string{"big.txt": big, "huge.txt": strings.Repeat("x", aiderMaxFileBytes+10)})
	var model, settings, meta string
	stubAider(t, func(c aiderCall) spawnResult {
		model = strings.Join(c.Args, " ")
		pl := aiderPlaceOf(tn.Home)
		settings = readFile(t, filepath.Join(pl.Work, "model-settings.yml"))
		meta = readFile(t, filepath.Join(pl.Work, "model-metadata.json"))
		return c.said("")
	})
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("big.txt", "x")); a.State != "refused" || !hasAnyOf(a.Why, "window (8192 tokens)") {
		t.Fatalf("a file that cannot fit an 8192 window: %+v", a)
	}
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("huge.txt", "x")); a.State != "refused" || !hasAnyOf(a.Why, "largest this door will give") {
		t.Fatalf("a file over the cap: %+v", a)
	}
	// the operator raises the seat's window and names another model; Aider follows the record
	write(t, tn.Home, "agents/expert_coder.md", fmt.Sprintf(aiderCoderDoc, "deepseek-coder:6.7b", 32768))
	aiderRackServer(t, "qwen2.5-coder:14b", "deepseek-coder:6.7b")
	a := callAider(t, reg, tr, "aider_run", aiderAsk("big.txt", "x"))
	if a.State != "ran" || a.Model != "deepseek-coder:6.7b" {
		t.Fatalf("after the seat moved: %+v", a)
	}
	if !strings.Contains(model, "--model ollama_chat/deepseek-coder:6.7b") || !strings.Contains(settings, "name: ollama_chat/deepseek-coder:6.7b") || !strings.Contains(settings, "num_ctx: 32768") {
		t.Fatalf("Aider did not follow the seat:\n%s\n%s", model, settings)
	}
	var doc map[string]map[string]any
	if err := json.Unmarshal([]byte(meta), &doc); err != nil || doc["ollama_chat/deepseek-coder:6.7b"]["max_input_tokens"] != float64(32768) {
		t.Fatalf("the model's metadata: %v / %s", err, meta)
	}
	// a model the rack does not hold -- a hosted route's name is one -- is refused by name
	write(t, tn.Home, "agents/expert_coder.md", fmt.Sprintf(aiderCoderDoc, "anthropic:claude-opus", 8192))
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x")); a.State != "refused" || !hasAnyOf(a.Why, "not on this machine's rack") {
		t.Fatalf("a hosted model: %+v", a)
	}
	write(t, tn.Home, "agents/expert_coder.md", fmt.Sprintf(aiderCoderDoc, "bad model;rm -rf", 8192))
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x")); a.State != "refused" || !hasAnyOf(a.Why, "not a model name") {
		t.Fatalf("a model name that is a command: %+v", a)
	}
	if err := os.Remove(filepath.Join(tn.Home, "agents", "expert_coder.md")); err != nil {
		t.Fatal(err)
	}
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x")); a.State != "refused" || !hasAnyOf(a.Why, "no coding seat") {
		t.Fatalf("no seat: %+v", a)
	}
}

// THE RACK MUST BE UP, and Aider must be installed. Each refuses by name, and the status says why.
func TestAiderSaysWhyItCannotRun(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	stubAider(t, func(c aiderCall) spawnResult { return c.said("") })
	t.Setenv("OLLAMA_HOST", "http://127.0.0.1:1")
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x")); a.State != "refused" || !hasAnyOf(a.Why, "rack is not reachable") {
		t.Fatalf("a silent rack: %+v", a)
	}
	out, _ := reg.Call(tr, "aider_status", map[string]any{}, agent)
	var st aiderStatus
	if err := json.Unmarshal([]byte(out), &st); err != nil || st.Rack != "silent" || !hasAnyOf(st.Why, "rack is not reachable") || !st.Installed {
		t.Fatalf("status on a silent rack: %v %+v", err, st)
	}
	if err := os.RemoveAll(filepath.Join(tn.Home, "aider")); err != nil {
		t.Fatal(err)
	}
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x")); a.State != "refused" || !hasAnyOf(a.Why, "not installed in this world") {
		t.Fatalf("no aider/venv: %+v", a)
	}
	out, _ = reg.Call(tr, "aider_status", map[string]any{}, agent)
	st = aiderStatus{}
	if err := json.Unmarshal([]byte(out), &st); err != nil || st.Installed || !hasAnyOf(st.Why, "not installed") || st.Runs == nil {
		t.Fatalf("status with no install: %v %+v", err, st)
	}
}

// A WORLD KEEPS THE LAST THIRTY RUNS and never removes a folder that is not a run's own.
func TestAiderKeepsTheLastThirtyRunsAndNothingElse(t *testing.T) {
	reg, tr, tn := aiderWorld(t, false, nil)
	pl := aiderPlaceOf(tn.Home)
	var olds []string
	for i := 0; i < 33; i++ {
		id := fmt.Sprintf("a2026010%d-0000%02d-0000", 1+i/10, i)
		olds = append(olds, id)
		write(t, pl.Runs, filepath.Join(id, "manifest.json"), `{"run":"`+id+`","files":[]}`)
	}
	write(t, pl.Runs, "not-a-run/keep.txt", "keep\n")
	write(t, pl.Work, "chat.md", "keep\n")
	stubAider(t, func(c aiderCall) spawnResult { c.editScratch(t, "notes.txt", "x\n"); return c.said("") })
	if a := callAider(t, reg, tr, "aider_run", aiderAsk("notes.txt", "x")); a.State != "ran" {
		t.Fatalf("%+v", a)
	}
	entries, _ := os.ReadDir(pl.Runs)
	var runs []string
	for _, e := range entries {
		if aiderRunRe.MatchString(e.Name()) {
			runs = append(runs, e.Name())
		}
	}
	sort.Strings(runs)
	if len(runs) != aiderKeepRuns || runs[0] != olds[4] || !exists(filepath.Join(pl.Runs, "not-a-run", "keep.txt")) || !exists(filepath.Join(pl.Work, "chat.md")) {
		t.Fatalf("the runs kept: %d, first %s", len(runs), runs[0])
	}
}

// ---- the wall, against a real Python -------------------------------------------------------------

// probeWall runs a probe script behind the REAL guard, in a real venv, and reports what each attempt came to. The
// ground is the temp world; outside is a second folder the wall is told is the user's profile.
type wallRun struct {
	Out, Err string
	Code     int
	Armed    bool
	Refused  []string
	Ground   string
	Scratch  string
	Home     string
	Outside  string
	RunDir   string
}

func probeWall(t *testing.T, probe string, extra ...string) wallRun {
	t.Helper()
	tn := world(t, map[string]string{".env": "KEY=value\n", "victim.txt": "must survive\n"})
	pl := aiderPlaceOf(tn.Home)
	aiderRealVenv(t, pl)
	runDir := filepath.Join(pl.Runs, "a20261004-120000-abcd")
	scratch := filepath.Join(runDir, "scratch")
	for _, d := range []string{scratch, pl.Home, pl.Tmp} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(t, runDir, "message.txt", "the instruction\n")
	write(t, runDir, "probe.py", probe)
	outside := t.TempDir()
	write(t, outside, "profile-file.txt", "the user's own\n")
	wall := aiderWallFor(tn, pl, runDir)
	wall.Profile = outside
	raw, _ := json.Marshal(wall)
	nonce := aiderNonce()
	args := append([]string{tn.Home, scratch, pl.Home, outside, runDir}, extra...)
	res := spawn(pl.Python, aiderGuardArgs(filepath.Join(runDir, "probe.py"), args...), spawnOpts{Dir: scratch, Timeout: 90 * time.Second,
		CleanEnv: aiderEnv(pl, "http://127.0.0.1:11434", nonce, string(raw)), Grace: 2 * time.Second})
	rest, armed, refused := aiderGuardLines(res.Stderr, nonce)
	code := 0
	var ee *exec.ExitError
	if res.Err != nil {
		code = -1
		if asExit(res.Err, &ee) {
			code = ee.ExitCode()
		}
	}
	return wallRun{Out: strings.ReplaceAll(res.Stdout, "\r\n", "\n"), Err: rest, Code: code, Armed: armed, Refused: refused,
		Ground: tn.Home, Scratch: scratch, Home: pl.Home, Outside: outside, RunDir: runDir}
}

func asExit(err error, target **exec.ExitError) bool {
	ee, ok := err.(*exec.ExitError)
	if ok {
		*target = ee
	}
	return ok
}

const wallProbe = `
import os, shutil, socket, subprocess, sys

ground, scratch, home, outside, rundir = sys.argv[1:6]
port = int(sys.argv[6])
lan = sys.argv[7] if len(sys.argv) > 7 else ''
results = {}

def attempt(name, fn):
    try:
        fn()
        results[name] = 'allowed'
    except PermissionError:
        results[name] = 'refused'
    except Exception as e:
        results[name] = 'error:' + type(e).__name__

def w(p, text='x'):
    with open(p, 'w') as f:
        f.write(text)

attempt('write_scratch', lambda: w(os.path.join(scratch, 'ok.txt')))
attempt('write_home', lambda: w(os.path.join(home, 'h.txt')))
attempt('write_devnull', lambda: w(os.devnull))
attempt('write_ground', lambda: w(os.path.join(ground, 'evil.txt')))
attempt('write_run_folder', lambda: w(os.path.join(rundir, 'planted.txt')))
attempt('mkdir_existing', lambda: os.mkdir(rundir))
attempt('mkdir_new_in_run_folder', lambda: os.mkdir(os.path.join(rundir, 'newdir')))
attempt('mkdir_new_in_scratch', lambda: os.mkdir(os.path.join(scratch, 'sub')))
attempt('write_venv', lambda: w(os.path.join(ground, 'aider', 'venv', 'planted.txt')))
attempt('write_outside', lambda: w(os.path.join(outside, 'planted.txt')))
attempt('read_env', lambda: open(os.path.join(ground, '.env')).read())
attempt('read_ground_file', lambda: open(os.path.join(ground, 'victim.txt')).read())
attempt('read_profile', lambda: open(os.path.join(outside, 'profile-file.txt')).read())
attempt('read_message', lambda: open(os.path.join(rundir, 'message.txt')).read())
attempt('read_own_library', lambda: open(os.__file__).read(10))
attempt('list_ground', lambda: os.listdir(ground))
attempt('list_scratch', lambda: os.listdir(scratch))
attempt('remove_ground', lambda: os.remove(os.path.join(ground, 'victim.txt')))
attempt('rename_into_ground', lambda: os.rename(os.path.join(scratch, 'ok.txt'), os.path.join(ground, 'moved.txt')))
attempt('copy_out_of_ground', lambda: shutil.copyfile(os.path.join(ground, 'victim.txt'), os.path.join(scratch, 'copied.txt')))
attempt('copy_within_scratch', lambda: shutil.copyfile(os.path.join(scratch, 'ok.txt'), os.path.join(scratch, 'copy.txt')))
attempt('connect_loopback', lambda: socket.create_connection(('127.0.0.1', port), 5).close())
if lan:
    attempt('connect_lan', lambda: socket.create_connection((lan, port), 2).close())
    attempt('lookup_lan', lambda: socket.getaddrinfo(lan, port))
attempt('lookup_localhost', lambda: socket.getaddrinfo('localhost', port))
attempt('bind_everywhere', lambda: socket.socket().bind(('0.0.0.0', 0)))
attempt('popen', lambda: subprocess.Popen([sys.executable, '-c', 'pass']))
attempt('system', lambda: os.system('echo hi'))
for k in sorted(results):
    print(k + '=' + results[k])
`

// THE WALL SHUTS WHAT IT NAMES, and leaves open what Aider needs. Measured with a real interpreter behind the real
// guard: every attempt comes back as the wall judged it, and what the refused ones would have done did not happen.
func TestTheWallShutsWhatItNamesAndLeavesOpenWhatAiderNeeds(t *testing.T) {
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		t.Skipf("no listener: %v", err)
	}
	defer ln.Close()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()
	port := ln.Addr().(*net.TCPAddr).Port
	lan := ""
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, a := range addrs {
			if ipn, ok := a.(*net.IPNet); ok && !ipn.IP.IsLoopback() && ipn.IP.To4() != nil {
				lan = ipn.IP.String()
				break
			}
		}
	}
	r := probeWall(t, wallProbe, fmt.Sprint(port), lan)
	if !r.Armed || r.Code != 0 {
		t.Fatalf("the probe did not run behind an announced wall: armed=%v code=%d\nout: %s\nerr: %s", r.Armed, r.Code, r.Out, r.Err)
	}
	got := map[string]string{}
	for _, ln := range strings.Split(strings.TrimSpace(r.Out), "\n") {
		if k, v, ok := strings.Cut(ln, "="); ok {
			got[k] = v
		}
	}
	want := map[string]string{
		"write_scratch": "allowed", "write_home": "allowed", "write_devnull": "allowed", "read_message": "allowed",
		"mkdir_existing": "error:FileExistsError", "mkdir_new_in_scratch": "allowed", "mkdir_new_in_run_folder": "refused",
		"read_own_library": "allowed", "list_scratch": "allowed", "copy_within_scratch": "allowed",
		"connect_loopback": "allowed", "lookup_localhost": "allowed",
		"write_ground": "refused", "write_run_folder": "refused", "write_venv": "refused", "write_outside": "refused",
		"read_env": "refused", "read_ground_file": "refused", "read_profile": "refused", "list_ground": "refused",
		"remove_ground": "refused", "rename_into_ground": "refused", "copy_out_of_ground": "refused",
		"bind_everywhere": "refused", "popen": "refused", "system": "refused",
	}
	if lan != "" { // this machine's own address: anything that is not loopback is "outside", and nothing can leave the machine
		want["connect_lan"], want["lookup_lan"] = "refused", "refused"
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s: the wall said %q, want %q", k, got[k], v)
		}
	}
	if len(got) != len(want) {
		t.Errorf("the probe made %d attempts and the stroke expects %d: %v", len(got), len(want), got)
	}
	// and what the refused ones would have done did not happen
	if exists(filepath.Join(r.RunDir, "newdir")) {
		t.Error("a folder was made in the run's own folder, where the undo copies and the record are kept")
	}
	for _, rel := range []string{"evil.txt", "moved.txt"} {
		if exists(filepath.Join(r.Ground, rel)) {
			t.Errorf("%s was made in the ground", rel)
		}
	}
	for _, p := range []string{filepath.Join(r.RunDir, "planted.txt"), filepath.Join(r.Outside, "planted.txt"), filepath.Join(r.Scratch, "copied.txt")} {
		if exists(p) {
			t.Errorf("%s exists", p)
		}
	}
	if readFile(t, filepath.Join(r.Ground, "victim.txt")) != "must survive\n" {
		t.Error("a file in the ground was removed or changed")
	}
	// every refusal is said by name, on stderr, for the card
	if len(r.Refused) != 6 { // the door lifts at most six; the probe made seventeen
		t.Errorf("what the wall said it refused (%d): %v", len(r.Refused), r.Refused)
	}
}

// TWO PROBES ARE REFUSED WITHOUT A WORD and everything else about a process is refused with one. GitPython asks `git version` when it
// is imported and the platform module asks the Windows shell for `ver` (three ways); every healthy run makes both, and a card that warns on every
// run teaches him to read past the one that matters. They are still refused (nothing starts), and a process the wall was not told of
// is said.
func TestTheWallRefusesTheTwoKnownProbesQuietlyAndEverythingElseOutLoud(t *testing.T) {
	r := probeWall(t, `
import os, subprocess, sys
results = {}
def run_as(line):
    # a string is a command line only on Windows (the form the platform module and the shell paths use); anywhere else a string names ONE program, so the same words go as a list
    return line if os.name == 'nt' else line.replace('"', '').split()
for name, cmd in (('git_version', run_as('git version')), ('git_version_list', ['git', 'version']), ('cmd_ver', run_as('cmd /c "ver"')),
                  ('cmd_ver2', run_as('cmd /c "command /c ver"')), ('cmd_ver3', run_as('cmd /c "cmd /c ver"')), ('git_status', run_as('git status')),
                  ('echo', run_as('cmd /c echo hi')), ('python', [sys.executable, '-c', 'pass'])):
    try:
        subprocess.Popen(cmd)
        results[name] = 'allowed'
    except PermissionError:
        results[name] = 'refused'
    except Exception as e:
        results[name] = 'error:' + type(e).__name__
for k in sorted(results):
    print(k + '=' + results[k])
`, "0")
	if !r.Armed || r.Code != 0 {
		t.Fatalf("the probe did not run behind an announced wall: %+v", r)
	}
	for _, name := range []string{"git_version", "git_version_list", "cmd_ver", "cmd_ver2", "cmd_ver3", "git_status", "echo", "python"} {
		if !strings.Contains(r.Out, name+"=refused") {
			t.Errorf("%s was not refused: %s", name, r.Out)
		}
	}
	// said: the three that are not the known probes; not said: the two that are
	if len(r.Refused) != 3 || hasAnyOf(r.Refused, "git version") || hasAnyOf(r.Refused, `/c ver`) {
		t.Errorf("what the wall said (%d): %v", len(r.Refused), r.Refused)
	}
	if !hasAnyOf(r.Refused, "git status") || !hasAnyOf(r.Refused, "echo hi") {
		t.Errorf("a process that is not a known probe was refused without a word: %v", r.Refused)
	}
}

// BOUNDED EVERYTHING (ESTATE LAW 7): a run's wait has a default long enough for this machine's measured speed (a few hundred tokens at 3.6 a second is minutes) and a
// ceiling no caller can raise; anything that is not a positive number takes the default.
func TestAiderWaitIsBoundedAndHasADefault(t *testing.T) {
	for in, want := range map[any]time.Duration{nil: 600 * time.Second, 30: 30 * time.Second, 30.5: 30500 * time.Millisecond, "45": 45 * time.Second,
		"junk": 600 * time.Second, 0: 600 * time.Second, -5: 600 * time.Second, 99999: 1200 * time.Second, "99999": 1200 * time.Second} {
		args := map[string]any{}
		if in != nil {
			args["timeout"] = in
		}
		if got := aiderWait(args); got != want {
			t.Errorf("timeout %v: waits %v, want %v", in, got, want)
		}
	}
}

// THE GUARD TRAVELS ON A COMMAND LINE (python -c <it>), and a Windows command line is one string that every program splits for itself: a double
// quote inside it is read by two sets of rules, and the line is at most 32,767 characters. The source carries no double quote and stays far under
// the limit. A convention doing a type's job, held here so that adding either is red at the stroke and not strange behaviour in a run.
func TestTheGuardSourceIsSafeToCarryOnAWindowsCommandLine(t *testing.T) {
	for name, src := range map[string]string{"the guard": aiderGuardSource, "the parse check": aiderParseSource} {
		if strings.ContainsAny(src, "\"\x00\r") {
			t.Errorf("%s carries a double quote, a NUL or a CR, which a Windows command line reads by two sets of rules", name)
		}
		if len(src) > 16000 {
			t.Errorf("%s is %d characters; a Windows command line holds 32,767 in all", name, len(src))
		}
	}
}

// THE WALL FAILS CLOSED: without the wall's own environment the guard does not start the program, and says nothing
// that could be taken for an announcement.
func TestTheWallDoesNotStartWithoutItsEnvironment(t *testing.T) {
	needPython(t)
	home := t.TempDir()
	cmd := exec.Command("python", aiderGuardArgs(filepath.Join(home, "probe.py"))...)
	write(t, home, "probe.py", "print('RAN')\n")
	cmd.Env = []string{"SYSTEMROOT=" + os.Getenv("SYSTEMROOT")}
	out, err := cmd.CombinedOutput()
	if err == nil || strings.Contains(string(out), "RAN") || strings.Contains(string(out), "armed") {
		t.Fatalf("the guard ran its program without a wall: %v\n%s", err, out)
	}
}

// THE WALL'S ANNOUNCEMENT IS THE DOOR'S TO READ and nobody else's to forge: a nonce the program could not have known.
func TestTheWallsAnnouncementNamesTheRunsOwnNonce(t *testing.T) {
	r := probeWall(t, "import os, sys\nprint('nonce in env:', 'ATLAS_AIDER_NONCE' in os.environ, 'ATLAS_AIDER_WALL' in os.environ)\n", "0")
	if !r.Armed || !strings.Contains(r.Out, "nonce in env: False False") {
		t.Fatalf("the program could read the wall's secrets, or the wall did not announce: armed=%v out=%q err=%q", r.Armed, r.Out, r.Err)
	}
	rest, armed, refused := aiderGuardLines(aiderGuardMark+" abc armed\n"+aiderGuardMark+" refused: x -- y\nother\n", "abc")
	if !armed || len(refused) != 1 || rest != "other\n" && rest != "other" {
		t.Fatalf("the lift of the wall's lines: %q %v %v", rest, armed, refused)
	}
	if _, armed, _ := aiderGuardLines(aiderGuardMark+" abc armed\n", "abd"); armed {
		t.Fatal("another run's nonce was believed")
	}
}

// ONE SOURCE FOR THE LAWS OF WHAT MAY BE WRITTEN: the seats' own list (manjuel/skills.py) and this door's copy of it agree.
// The core's suite holds the same two lists equal from the other side, where the ground carries both.
func TestAiderNeverWritesTheSameNamesTheSeatsNeverWrite(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "manjuel", "skills.py"))
	if err != nil {
		t.Skipf("this checkout carries no core beside it (%v); the core's suite holds the lists equal where it does", err)
	}
	src := string(raw)
	block := func(name string) string {
		i := strings.Index(src, name+" = {")
		if i < 0 {
			t.Fatalf("manjuel/skills.py has no %s", name)
		}
		rest := src[i+len(name)+4:]
		return rest[:strings.Index(rest, "}")]
	}
	quoted := regexp.MustCompile(`"([^"]+)"`)
	keys := func(text string, pairs bool) []string {
		var out []string
		for _, ln := range strings.Split(text, "\n") {
			for _, m := range quoted.FindAllStringSubmatch(ln, -1) {
				out = append(out, m[1])
				if pairs {
					break
				}
			}
		}
		return out
	}
	setOf := func(names []string) string { sort.Strings(names); return strings.Join(names, ",") }
	var top, files []string
	for k := range aiderNeverTop {
		if k != "aider" { // aider/ is this door's own addition: the seats' list has no such folder
			top = append(top, k)
		}
	}
	for k := range aiderNeverFiles {
		files = append(files, k)
	}
	if got, want := setOf(keys(block("_NEVER_WRITTEN_TOP"), true)), setOf(top); got != want {
		t.Errorf("the seats' folders %s\nthis door's      %s", got, want)
	}
	if got, want := setOf(keys(block("_NEVER_WRITTEN_FILES"), true)), setOf(files); got != want {
		t.Errorf("the seats' files %s\nthis door's     %s", got, want)
	}
	if got, want := setOf(keys(block("_NEVER_WRITTEN_STAMPS"), false)), setOf(append([]string{}, aiderNeverStamps...)); got != want {
		t.Errorf("the seats' stamps %s\nthis door's      %s", got, want)
	}
	if got, want := setOf(keys(block("_NEVER_WRITTEN_SUFFIXES"), false)), setOf(append([]string{}, aiderNeverSuffixes...)); got != want {
		t.Errorf("the seats' binaries %s\nthis door's        %s", got, want)
	}
}
