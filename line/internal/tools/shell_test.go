package tools

// The operator's typed shell (WHAT'S LEFT H15, 2026-10-03). These are the wires: each test names the
// thing that would let a shell reach a hand it was not given to, or let a key out, and runs the real
// tool -- a real Git Bash, a real Python, the real hold queue -- against a ground it builds in a temp
// folder. Nothing here is a stub of the gate; the gate is what is being measured.
//
// HERMETIC: every ground is t.TempDir(); the only network any command could reach is refused at the
// card or by name. Tests that need bash or python skip when the machine has none, and say so.

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"atlas/line/internal/tenant"
)

func needBash(t *testing.T) {
	t.Helper()
	if _, err := findBash(); err != nil {
		t.Skipf("no bash here: %v", err)
	}
}

func needPython(t *testing.T) {
	t.Helper()
	if out, err := exec.Command("python", "-c", "import ast, json; print('ok')").Output(); err != nil || !strings.Contains(string(out), "ok") {
		t.Skip("no python on the PATH")
	}
}

// shellWorld is a ground with a registry whose holds are armed, as the door runs with --auth.
func shellWorld(t *testing.T, files map[string]string) (*Registry, *tenant.Registry, tenant.Tenant) {
	t.Helper()
	tn := world(t, files)
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
	t.Cleanup(func() { s := pySessionFor(resolved.Home); s.mu.Lock(); s.reset(); s.mu.Unlock() })
	return Build(tr, Options{HoldWrites: true}), tr, resolved
}

func callShell(t *testing.T, reg *Registry, tr *tenant.Registry, who Caller, args map[string]any) shellAnswer {
	t.Helper()
	out, err := reg.Call(tr, "shell_run", args, who)
	if err != nil {
		t.Fatalf("shell_run errored: %v", err)
	}
	var a shellAnswer
	if err := json.Unmarshal([]byte(out), &a); err != nil {
		t.Fatalf("shell_run did not answer in its shape: %v\n%s", err, out)
	}
	return a
}

func bash(cmd string) map[string]any   { return map[string]any{"shell": "bash", "command": cmd} }
func python(cmd string) map[string]any { return map[string]any{"shell": "python", "command": cmd} }

// decide answers a hold as the operator's glass and returns what the replay said, in its shape.
func decide(t *testing.T, reg *Registry, tr *tenant.Registry, id, decision string) (string, shellAnswer) {
	t.Helper()
	out, err := reg.Call(tr, "hold_answer", map[string]any{"id": id, "decision": decision}, glass)
	if err != nil {
		t.Fatalf("hold_answer errored: %v", err)
	}
	var a shellAnswer
	if i := strings.Index(out, "{"); decision == "approve" && i >= 0 {
		if err := json.Unmarshal([]byte(out[i:]), &a); err != nil {
			t.Fatalf("the approved run did not answer in its shape: %v\n%s", err, out)
		}
	}
	return out, a
}

func exists(p string) bool { _, err := os.Stat(p); return err == nil }

func heldIDs(t *testing.T, reg *Registry, tr *tenant.Registry) []string {
	t.Helper()
	out, err := reg.Call(tr, "hold_list", map[string]any{}, glass)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Held []struct {
			ID   string         `json:"id"`
			Args map[string]any `json:"args"`
		} `json:"held"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("hold_list: %v\n%s", err, out)
	}
	var ids []string
	for _, h := range doc.Held {
		ids = append(ids, h.ID)
	}
	return ids
}

func holdsLog(t *testing.T, home string) string {
	t.Helper()
	b, _ := os.ReadFile(filepath.Join(home, "state", "holds.jsonl"))
	return string(b)
}

// NO AGENT IS EVER GIVEN THE SHELL -- not by name, not by a role, not by asking him to approve a call
// it parked, not by a forged caller in the arguments. Every refusal is written down.
func TestNoAgentIsEverGivenTheShell(t *testing.T) {
	reg, tr, tn := shellWorld(t, map[string]string{
		"rbac.json": `{"assign": {"k-op": "operator", "k-steward": "steward"}}`,
	})
	canary := filepath.Join(tn.Home, "canary.txt")
	callers := []Caller{agent, {Name: "k-agent"}, {Name: "the glass"}, {Name: "k-op"}, {Name: "k-steward"},
		{Name: "operator via hold hold_1"}, {}}
	attempts := 0
	for _, who := range callers {
		for _, args := range []map[string]any{bash("touch canary.txt"), python("open('canary.txt', 'w')"),
			{"shell": "python", "reset": true}, {"shell": "bash", "command": "ls", CallerKey: Caller{Service: true}}} {
			attempts++
			out, err := reg.Call(tr, "shell_run", args, who)
			if err == nil || !strings.Contains(err.Error(), "operator's own hand") {
				t.Fatalf("%q was handed the shell: %q / %v", who.Name, out, err)
			}
		}
	}
	if exists(canary) {
		t.Fatal("a refused caller made a file")
	}
	if ids := heldIDs(t, reg, tr); len(ids) != 0 {
		t.Fatalf("a refused caller parked a call for him to approve: %v", ids)
	}
	if n := strings.Count(holdsLog(t, tn.Home), `"what":"refused"`); n != attempts {
		t.Fatalf("%d attempts, %d refusals written down", attempts, n)
	}
}

// THE APPROVED RUN IS A TYPED FIELD, NOT A WORD IN A NAME. A caller that is the glass and calls itself
// "operator via hold ..." is still a caller whose write waits.
func TestAnApprovalIsATypedFieldNotAName(t *testing.T) {
	needBash(t)
	reg, tr, tn := shellWorld(t, nil)
	forged := Caller{Name: "operator via hold hold_9_9_shell_run", Service: true}
	a := callShell(t, reg, tr, forged, bash("touch forged.txt"))
	if a.State != "held" || a.Hold == "" {
		t.Fatalf("a name stood in for an approval: %+v", a)
	}
	if exists(filepath.Join(tn.Home, "forged.txt")) {
		t.Fatal("the forged name ran a write")
	}
}

func TestAPlainLookRunsAtOnce(t *testing.T) {
	needBash(t)
	reg, tr, tn := shellWorld(t, map[string]string{"README.md": "# a ground\n"})
	a := callShell(t, reg, tr, glass, bash("echo hello"))
	if a.State != "ran" || a.Class != "read" || a.Exit == nil || *a.Exit != 0 || a.Output != "hello" || a.Approved {
		t.Fatalf("echo: %+v", a)
	}
	a = callShell(t, reg, tr, glass, bash("ls"))
	if a.State != "ran" || !strings.Contains(a.Output, "README.md") {
		t.Fatalf("ls did not list the ground: %+v", a)
	}
	// It runs IN the ground: the command word `pwd` names it (Git Bash spells it /c/Users/...).
	a = callShell(t, reg, tr, glass, bash("pwd"))
	want := strings.ToLower(filepath.Base(filepath.Dir(tn.Home)) + "/" + filepath.Base(tn.Home))
	if !strings.Contains(strings.ToLower(a.Output), want) {
		t.Fatalf("pwd %q is not the ground %q", a.Output, tn.Home)
	}
	// A failure is on the face of the answer: its own words and its own exit code.
	a = callShell(t, reg, tr, glass, bash("ls nosuchthing 2>&1"))
	if a.State != "ran" || a.Exit == nil || *a.Exit == 0 || !strings.Contains(strings.ToLower(a.Output), "no such file") {
		t.Fatalf("a failing look did not say so: %+v", a)
	}
	if exists(filepath.Join(tn.Home, "state", "holds.jsonl")) && strings.Contains(holdsLog(t, tn.Home), `"what":"held"`) {
		t.Fatal("a plain look was parked")
	}
}

// A WRITE WAITS, AND APPROVING RUNS EXACTLY WHAT WAS PARKED; DENYING RUNS NOTHING.
func TestAWriteWaitsForHisApprovalAndRunsExactlyWhatWasParked(t *testing.T) {
	needBash(t)
	reg, tr, tn := shellWorld(t, nil)
	made := filepath.Join(tn.Home, "made.txt")

	a := callShell(t, reg, tr, glass, bash("touch made.txt"))
	if a.State != "held" || a.Hold == "" || a.Class != "write" || len(a.Why) == 0 || a.Exit != nil {
		t.Fatalf("a write was not parked: %+v", a)
	}
	if exists(made) {
		t.Fatal("a write ran before his hand")
	}
	ids := heldIDs(t, reg, tr)
	if len(ids) != 1 || ids[0] != a.Hold {
		t.Fatalf("the queue holds %v, the answer named %s", ids, a.Hold)
	}
	if out, _ := decide(t, reg, tr, a.Hold, "deny"); !strings.Contains(out, "Denied") || exists(made) {
		t.Fatalf("denying ran it: %s", out)
	}

	a = callShell(t, reg, tr, glass, bash("touch made.txt"))
	_, ran := decide(t, reg, tr, a.Hold, "approve")
	if ran.State != "ran" || !ran.Approved || ran.Class != "write" || ran.Exit == nil || *ran.Exit != 0 {
		t.Fatalf("the approved run: %+v", ran)
	}
	if !exists(made) {
		t.Fatal("approving did not run what was parked")
	}
	if len(heldIDs(t, reg, tr)) != 0 {
		t.Fatal("an answered hold is still waiting")
	}
}

// A REFUSAL IS NOT A QUESTION. Nothing is parked for a refused entry, and a refused entry that gets
// into the queue some other way is refused again at the replay: approval cannot lift it.
func TestApprovalCannotLiftARefusal(t *testing.T) {
	needBash(t)
	reg, tr, tn := shellWorld(t, map[string]string{".env": "ROUTE_KEY=SECRETVALUE_abc12345\n"})

	a := callShell(t, reg, tr, glass, bash("cat .env"))
	if a.State != "refused" || a.Class != "refused" || a.Hold != "" || strings.Contains(a.Output, "SECRETVALUE") {
		t.Fatalf("cat .env: %+v", a)
	}
	if len(heldIDs(t, reg, tr)) != 0 {
		t.Fatal("a refused entry was parked for his approval")
	}

	self, _ := reg.Get("shell_run")
	h := reg.park(tn, self, map[string]any{"shell": "bash", "command": "cat .env"}, glass)
	out, replay := decide(t, reg, tr, h.ID, "approve")
	if replay.State != "refused" || strings.Contains(out, "SECRETVALUE") || strings.Contains(out, "abc12345") {
		t.Fatalf("approving lifted a refusal: %s", out)
	}
	for _, c := range []string{"cat ../x", "ls /etc", "cat credentials.json", "cat vault/x"} {
		if a := callShell(t, reg, tr, glass, bash(c)); a.State != "refused" {
			t.Errorf("%q: %+v", c, a)
		}
	}
}

// KEYS ARE SILENT. The child is handed a built environment, not the door's: the service wire and the
// route key are in the door's, and are in neither a bash command's nor a Python entry's.
func TestTheChildInheritsNoKeys(t *testing.T) {
	needBash(t)
	t.Setenv("ATLAS_SERVICE", "canary-service-wire-0123456789")
	t.Setenv("MANJUEL_ROUTE_ANTHROPIC_KEY", "canary-route-key-0123456789")
	t.Setenv("SOME_SETTING_NOBODY_LISTED", "visible-only-if-inherited")
	reg, tr, _ := shellWorld(t, nil)

	for _, cmd := range []string{
		`echo "svc=[$ATLAS_SERVICE] key=[$MANJUEL_ROUTE_ANTHROPIC_KEY] other=[$SOME_SETTING_NOBODY_LISTED]"`,
		"env", "printenv",
	} {
		a := callShell(t, reg, tr, glass, bash(cmd))
		if a.State != "held" {
			t.Fatalf("%q was %s, want a card", cmd, a.State)
		}
		_, ran := decide(t, reg, tr, a.Hold, "approve")
		if ran.State != "ran" {
			t.Fatalf("%q: %+v", cmd, ran)
		}
		for _, bad := range []string{"canary-", "visible-only-if-inherited", "ATLAS_SERVICE=", "MANJUEL_ROUTE_ANTHROPIC_KEY="} {
			if strings.Contains(ran.Output, bad) {
				t.Errorf("%q printed %q: the child inherited the door's environment\n%s", cmd, bad, ran.Output)
			}
		}
		if strings.HasPrefix(cmd, "echo") && !strings.Contains(ran.Output, "svc=[] key=[] other=[]") {
			t.Errorf("the echo saw a value: %q", ran.Output)
		}
	}
}

func TestPythonInheritsNoKeysEither(t *testing.T) {
	needPython(t)
	t.Setenv("ATLAS_SERVICE", "canary-service-wire-0123456789")
	t.Setenv("SOME_SETTING_NOBODY_LISTED", "visible-only-if-inherited")
	reg, tr, _ := shellWorld(t, nil)
	a := callShell(t, reg, tr, glass, python(`import os; print(os.environ.get("ATLAS_SERVICE"), os.environ.get("SOME_SETTING_NOBODY_LISTED"))`))
	if a.State != "held" {
		t.Fatalf("an import was %s, want a card", a.State)
	}
	_, ran := decide(t, reg, tr, a.Hold, "approve")
	if ran.Output != "None None" {
		t.Fatalf("the session inherited the door's environment: %q", ran.Output)
	}
}

// WHAT IT PRINTS IS SCRUBBED OF EVERY SECRET VALUE THE GROUND HOLDS, whatever route reached it: a
// recursive search is a read, and it reads .env like any file.
func TestWhatItPrintsIsScrubbedOfEverySecretValue(t *testing.T) {
	needBash(t)
	reg, tr, _ := shellWorld(t, map[string]string{
		".env":      "ROUTE_KEY=SECRETVALUE_abc12345\nexport QUOTED=\"quoted-secret-value-99\"\nSHORT=tiny\n",
		"README.md": "mentions SECRETVALUE_ only in passing\n",
	})
	a := callShell(t, reg, tr, glass, bash("grep -r SECRETVALUE_abc ."))
	if a.State != "ran" || a.Class != "read" {
		t.Fatalf("a recursive search is a read: %+v", a)
	}
	if strings.Contains(a.Output, "abc12345") || !strings.Contains(a.Output, shellWithheld) {
		t.Fatalf("a secret value reached the page: %q", a.Output)
	}
	a = callShell(t, reg, tr, glass, bash("echo quoted-secret-value-99"))
	if a.State != "ran" || strings.Contains(a.Output, "quoted-secret-value-99") || a.Output != shellWithheld {
		t.Fatalf("a quoted value reached the page: %+v", a)
	}
	// the command it echoes back is scrubbed the same way
	if strings.Contains(a.Command, "quoted-secret-value-99") {
		t.Fatalf("the command echo carries a secret: %q", a.Command)
	}
}

// THE LAW MUST VERIFY: no run proceeds on a law that does not, and one that does is cached against
// what law/ holds, so a changed law is walked again.
func TestNothingRunsOnALawThatDoesNotVerify(t *testing.T) {
	needBash(t)
	needPython(t)
	bad := "import sys\nprint('THE CHAIN REFUSES:')\nprint('  - link #2 fingerprint MISMATCH')\nsys.exit(1)\n"
	reg, tr, tn := shellWorld(t, map[string]string{"law/chain.jsonl": "{}\n", "law/law.py": bad})
	a := callShell(t, reg, tr, glass, bash("echo hi"))
	if a.State != "refused" || !hasAny(a.Why, "THE CHAIN") || a.Output != "" {
		t.Fatalf("a run proceeded on a broken law: %+v", a)
	}
	a = callShell(t, reg, tr, glass, python("1 + 1"))
	if a.State != "refused" {
		t.Fatalf("python proceeded on a broken law: %+v", a)
	}
	good := "print('the law chain proves whole: 1 links, head abc')\n"
	if err := os.WriteFile(filepath.Join(tn.Home, "law", "law.py"), []byte(good+"# changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	a = callShell(t, reg, tr, glass, bash("echo hi"))
	if a.State != "ran" || a.Output != "hi" {
		t.Fatalf("a mended law did not mend the gate: %+v", a)
	}
	// and a ground with no ledger at all is not a broken one
	reg2, tr2, _ := shellWorld(t, nil)
	if a := callShell(t, reg2, tr2, glass, bash("echo hi")); a.State != "ran" {
		t.Fatalf("a ground with no law ledger was refused: %+v", a)
	}
}

func hasAny(list []string, frag string) bool {
	for _, s := range list {
		if strings.Contains(s, frag) {
			return true
		}
	}
	return false
}

// EVERY RUN, HOLD AND REFUSAL IS A LINE IN THE GROUND'S HOLDS LOG, with the command, keys withheld.
func TestEveryRunHoldAndRefusalIsRecorded(t *testing.T) {
	needBash(t)
	reg, tr, tn := shellWorld(t, nil)
	callShell(t, reg, tr, glass, bash("echo recorded-read"))
	callShell(t, reg, tr, glass, bash("cat .env"))
	a := callShell(t, reg, tr, glass, bash("touch recorded-write.txt"))
	decide(t, reg, tr, a.Hold, "approve")
	log := holdsLog(t, tn.Home)
	for _, want := range []string{`"what":"shell_ran"`, `"what":"shell_refused"`, `"what":"shell_held"`, `"what":"approved_ran"`,
		"echo recorded-read", "cat .env", "touch recorded-write.txt", `"tool":"shell_run"`} {
		if !strings.Contains(log, want) {
			t.Errorf("the holds log never says %q:\n%s", want, log)
		}
	}
	if strings.Count(log, `"what":"shell_ran"`) != 2 {
		t.Errorf("want two runs recorded (the read, and the approved write):\n%s", log)
	}
}

// THE ANSWER'S SHAPE is the contract with the page: these keys, and no others.
func TestTheShellAnswersInOneShape(t *testing.T) {
	needBash(t)
	reg, tr, _ := shellWorld(t, nil)
	allowed := map[string]bool{"shell": true, "state": true, "class": true, "why": true, "hold": true, "approved": true,
		"exit": true, "ms": true, "timed_out": true, "truncated": true, "output": true, "note": true, "command": true}
	for name, args := range map[string]map[string]any{"ran": bash("echo x"), "held": bash("touch s.txt"), "refused": bash("cat .env")} {
		out, err := reg.Call(tr, "shell_run", args, glass)
		if err != nil {
			t.Fatal(err)
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(out), &m); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if m["state"] != name || m["shell"] != "bash" {
			t.Errorf("%s answered %v", name, m)
		}
		for k := range m {
			if !allowed[k] {
				t.Errorf("%s answered a key the page does not know: %q", name, k)
			}
		}
	}
	// a malformed call is an error, not a document
	for _, args := range []map[string]any{{"shell": "zsh", "command": "ls"}, {"shell": "bash", "command": "  "}, {"command": "ls"}} {
		if out, err := reg.Call(tr, "shell_run", args, glass); err == nil {
			t.Errorf("%v answered %q", args, out)
		}
	}
}

// AN ENDLESS COMMAND ENDS, and the answer says it was ended, not that it failed.
func TestAnEndlessCommandIsEndedAndSaysSo(t *testing.T) {
	needBash(t)
	reg, tr, _ := shellWorld(t, nil)
	args := bash("sleep 60 | cat")
	args["timeout"] = float64(1)
	a := callShell(t, reg, tr, glass, args)
	if a.State != "held" {
		t.Fatalf("sleep was %s", a.State)
	}
	began := time.Now()
	_, ran := decide(t, reg, tr, a.Hold, "approve")
	if took := time.Since(began); took > 20*time.Second {
		t.Fatalf("a 1 s limit took %s", took)
	}
	if !ran.TimedOut || ran.Exit == nil || *ran.Exit != -1 || !strings.Contains(ran.Note, "ended") {
		t.Fatalf("the answer does not say it was ended: %+v", ran)
	}
}

// PYTHON KEEPS ITS NAMES, a card is asked for what is not a calculation, an entry that runs away takes
// the session with it and says so, and what a child process prints cannot pose as an answer.
func TestPythonKeepsItsNamesAndAsksAboutTheRest(t *testing.T) {
	needPython(t)
	reg, tr, tn := shellWorld(t, nil)
	for _, c := range []struct{ code, out string }{{"x = 41", ""}, {"x + 1", "42"}, {"print('hi'); x", "hi\n41"}, {"[i * i for i in range(4)]", "[0, 1, 4, 9]"}} {
		a := callShell(t, reg, tr, glass, python(c.code))
		if a.State != "ran" || a.Class != "read" || a.Output != c.out {
			t.Fatalf("%q: %+v", c.code, a)
		}
	}
	a := callShell(t, reg, tr, glass, python("1 / 0"))
	if a.State != "ran" || a.Exit == nil || *a.Exit != 1 || !strings.Contains(a.Output, "ZeroDivisionError") {
		t.Fatalf("an error is not on the face of the answer: %+v", a)
	}
	// a card: an import, and a child process's words arrive without breaking the wire
	a = callShell(t, reg, tr, glass, python(`import os; os.system("echo from-a-child"); open("py-made.txt", "w").close()`))
	if a.State != "held" || exists(filepath.Join(tn.Home, "py-made.txt")) {
		t.Fatalf("an import and a file were not asked about: %+v", a)
	}
	_, ran := decide(t, reg, tr, a.Hold, "approve")
	if ran.State != "ran" || !ran.Approved || !strings.Contains(ran.Output, "from-a-child") || !exists(filepath.Join(tn.Home, "py-made.txt")) {
		t.Fatalf("the approved entry: %+v", ran)
	}
	if a := callShell(t, reg, tr, glass, python("x + 2")); a.Output != "43" {
		t.Fatalf("the session lost its names after a child process spoke: %+v", a)
	}
	// a runaway ends the session, and the next entry is told
	args := python("while True: pass")
	args["timeout"] = float64(1)
	a = callShell(t, reg, tr, glass, args)
	if !a.TimedOut || !strings.Contains(a.Note, "names are gone") {
		t.Fatalf("a runaway entry: %+v", a)
	}
	a = callShell(t, reg, tr, glass, python("1 + 1"))
	if a.Output != "2" || !strings.Contains(a.Note, "starts fresh") {
		t.Fatalf("the entry after a runaway was not told its names are gone: %+v", a)
	}
	// reset is his: names go, and the next start is not reported as a loss
	callShell(t, reg, tr, glass, python("keep = 1"))
	r := callShell(t, reg, tr, glass, map[string]any{"shell": "python", "reset": true})
	if r.State != "ran" || !strings.Contains(r.Note, "ended") {
		t.Fatalf("reset: %+v", r)
	}
	a = callShell(t, reg, tr, glass, python("keep"))
	if a.Exit == nil || *a.Exit != 1 || !strings.Contains(a.Output, "NameError") || a.Note != "" {
		t.Fatalf("a name survived a reset: %+v", a)
	}
}

// EVERY TOOL DECLARED SERVICE-ONLY REFUSES EVERY OTHER CALLER, whatever it is and whatever it is given.
func TestEveryServiceOnlyToolRefusesEveryOtherCaller(t *testing.T) {
	reg, tr, tn := shellWorld(t, nil)
	found := map[string]bool{}
	for _, tool := range reg.All() {
		if !tool.ServiceOnly {
			continue
		}
		found[tool.Name] = true
		if !tool.Writes {
			t.Errorf("%s is the operator's own hand and declares itself a reader", tool.Name)
		}
		for _, who := range []Caller{agent, {Name: "k-agent"}, {}} {
			args := argsFor(tool, tn.Home)
			args["project"] = tn.Name // argsFor names its own world; this one is "probe"
			_, err := reg.Call(tr, tool.Name, args, who)
			if err == nil || !strings.Contains(err.Error(), "operator's own hand") {
				t.Errorf("%s answered %q: %v", tool.Name, who.Name, err)
			}
		}
	}
	if !found["shell_run"] {
		t.Fatal("shell_run is not declared service-only; this stroke would be the vacuous one it replaced")
	}
}

// THE TOOL REFUSES FOR ITSELF TOO. Registry.Call is the first line; a tool reached any other way -- a
// replayed hold, a later caller of the function -- still will not run for a hand that is not his.
func TestTheToolRefusesForItselfToo(t *testing.T) {
	reg, _, tn := shellWorld(t, nil)
	tool, ok := reg.Get("shell_run")
	if !ok {
		t.Fatal("no shell_run")
	}
	for _, who := range []Caller{agent, {}, {Name: "operator via hold hold_1"}} {
		out, err := tool.Fn(tn, map[string]any{"shell": "bash", "command": "touch direct.txt", CallerKey: who})
		if err == nil || out != "" || !strings.Contains(err.Error(), "operator's own hand") {
			t.Fatalf("%q was answered: %q / %v", who.Name, out, err)
		}
	}
	if exists(filepath.Join(tn.Home, "direct.txt")) {
		t.Fatal("the tool ran for a hand that was not his")
	}
}

// AN ENTRY THAT IS A FILE IS REFUSED, and a flood is cut short rather than held: the answer carries the
// first 64 KB and says it was cut.
func TestAnEntryThatIsAFileIsRefusedAndAFloodIsCut(t *testing.T) {
	needBash(t)
	reg, tr, _ := shellWorld(t, nil)
	long := bash("echo " + strings.Repeat("a", shellMaxCommand+10))
	if a := callShell(t, reg, tr, glass, long); a.State != "refused" || !hasAny(a.Why, "longer than") {
		t.Fatalf("an 8000-character entry: %+v", a)
	}
	a := callShell(t, reg, tr, glass, bash("seq 1 200000"))
	if a.State != "held" {
		t.Fatalf("seq was %s", a.State)
	}
	_, ran := decide(t, reg, tr, a.Hold, "approve")
	if !ran.Truncated || len(ran.Output) > shellMaxBytes || !strings.HasPrefix(ran.Output, "1\n2\n3\n") {
		t.Fatalf("a flood of %d bytes, truncated=%v, begins %q", len(ran.Output), ran.Truncated, ran.Output[:min(20, len(ran.Output))])
	}
}

// THE FLOOD HAS TWO CAPS and each is held on its own: the spawn helper keeps the first bytes of a stream
// and drains the rest (so a child that floods its pipe is never blocked and never held whole), and the
// shape of the answer cuts what is left, whole characters only.
func TestSpawnKeepsTheFirstBytesOfAFloodAndSaysSo(t *testing.T) {
	needBash(t)
	bashPath, _ := findBash()
	res := spawn(bashPath, []string{"-c", "seq 1 100000"}, spawnOpts{Dir: t.TempDir(), MaxBytes: 1000, Timeout: 20 * time.Second})
	if res.Err != nil || len(res.Stdout) != 1000 || !res.Truncated || !strings.HasPrefix(res.Stdout, "1\n2\n3\n") {
		t.Fatalf("a flood: %d bytes, truncated=%v, err=%v", len(res.Stdout), res.Truncated, res.Err)
	}
	res = spawn(bashPath, []string{"-c", "seq 1 10"}, spawnOpts{Dir: t.TempDir(), Timeout: 20 * time.Second})
	if res.Truncated || !strings.HasSuffix(strings.TrimSpace(res.Stdout), "10") {
		t.Fatalf("an unbounded stream was cut: %+v", res)
	}
}

func TestShellShapeCutsAndScrubs(t *testing.T) {
	home := t.TempDir()
	out, cut := shellShape(home, strings.Repeat("é", 200), 101)
	if !cut || len(out) > 101 || !utf8Valid(out) {
		t.Fatalf("cut=%v len=%d valid=%v", cut, len(out), utf8Valid(out))
	}
	out, cut = shellShape(home, "line one\r\nline two\r\n", 1000)
	if cut || out != "line one\nline two" {
		t.Fatalf("terminators and trailing newlines: %q", out)
	}
	out, _ = shellShape(home, "a key atl_0123456789abcdef0123456789abcdef in a line", 1000)
	if strings.Contains(out, "atl_0123") || !strings.Contains(out, shellWithheld) {
		t.Fatalf("a key-shaped value was not withheld: %q", out)
	}
}

func utf8Valid(s string) bool { return strings.ToValidUTF8(s, "") == s }

// A DOUBLED BACKSLASH CANNOT SMUGGLE A SECRET'S NAME PAST THE GATE (found 2026-10-03, an hour after the
// first placement). Handed to Git Bash on its command line, a doubled backslash was halved on the way in:
// the gate lexed a file called ".\env" (two backslashes typed) and bash was handed one backslash, which it
// reads as an escaped e, so .env. The entry travels in an environment variable now, unchanged. The ground's .env holds a value
// too short to be scrubbed, so a read of it would show.
func TestADoubledBackslashIsNotASecretName(t *testing.T) {
	needBash(t)
	reg, tr, _ := shellWorld(t, map[string]string{".env": "PLAIN=tiny\n"})
	a := callShell(t, reg, tr, glass, bash(`cat .\env`))
	if a.State != "refused" || strings.Contains(a.Output, "PLAIN") {
		t.Fatalf("one backslash IS .env to bash, and must be refused by name: %+v", a)
	}
	a = callShell(t, reg, tr, glass, bash(`cat .\\env`))
	if strings.Contains(a.Output, "PLAIN") || a.Exit == nil || *a.Exit == 0 || !strings.Contains(strings.ToLower(a.Output), "no such file") {
		t.Fatalf("two backslashes are a literal one in a name that does not exist, and bash must say so: %+v", a)
	}
}
