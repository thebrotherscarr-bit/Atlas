package engine

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// I8: this package had no tests at all. SittingOpen is pure and it is the
// guard that stands between a second engine and a forked ledger, so it is
// where coverage starts.
func ledger(t *testing.T, body string) string {
	t.Helper()
	g := t.TempDir()
	if err := os.MkdirAll(filepath.Join(g, "sessions"), 0o755); err != nil {
		t.Fatal(err)
	}
	if body != "" {
		if err := os.WriteFile(filepath.Join(g, "sessions", "sessions.jsonl"),
			[]byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return g
}

func TestSittingOpen(t *testing.T) {
	cases := []struct {
		name string
		body string
		open bool
	}{
		{"no ledger at all is not an open sitting", "", false},
		{"a closed last line is not open",
			`{"n":1,"started":"a","ended":"b"}` + "\n", false},
		{"an empty ended IS open",
			`{"n":99,"started":"2026-09-09T06:36:55","ended":""}` + "\n", true},
		{"a missing ended IS open",
			`{"n":7,"started":"x"}` + "\n", true},
		{"only the LAST line decides",
			`{"n":1,"started":"a","ended":""}` + "\n" +
				`{"n":2,"started":"b","ended":"c"}` + "\n", false},
		{"trailing blank lines are skipped",
			`{"n":3,"started":"a","ended":"b"}` + "\n\n\n", false},
		{"CRLF is the ruling and must parse",
			`{"n":4,"started":"a","ended":""}` + "\r\n", true},
		// C5: the whole point. A killed engine leaves a half-written line, and
		// reporting that as "no sitting" invites a second engine into a world
		// whose record is already damaged.
		{"a corrupt last line REFUSES rather than inviting a second engine",
			`{"n":5,"started":"a","ended":"b"}` + "\n" + `{"n":6,"star`, true},
		{"garbage refuses", "not json at all\n", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, _, open := SittingOpen(ledger(t, c.body))
			if open != c.open {
				t.Fatalf("open = %v, want %v", open, c.open)
			}
		})
	}
}

func TestSittingOpenNamesTheSitting(t *testing.T) {
	n, started, open := SittingOpen(ledger(t,
		`{"n":99,"started":"2026-09-09T06:36:55","ended":""}`+"\n"))
	if !open || int(n) != 99 || started != "2026-09-09T06:36:55" {
		t.Fatalf("got (%v,%q,%v); the refusal must be able to NAME the sitting",
			n, started, open)
	}
}

func TestSplitCommand(t *testing.T) {
	for _, c := range []struct {
		in   string
		want []string
	}{
		{`python manjuel.py`, []string{"python", "manjuel.py"}},
		{`"C:\Program Files\py.exe" x.py`, []string{`C:\Program Files\py.exe`, "x.py"}},
		{`  python   x.py  `, []string{"python", "x.py"}},
		{``, nil},
	} {
		got := splitCommand(c.in)
		if len(got) != len(c.want) {
			t.Fatalf("%q -> %v, want %v", c.in, got, c.want)
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Fatalf("%q -> %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestOpenRefusesWithoutACommand(t *testing.T) {
	if _, err := Open("w", t.TempDir(), "  "); err == nil {
		t.Fatal("an unwired engine must refuse, not spawn")
	}
}

func TestOpenRefusesAWorldBeingSatIn(t *testing.T) {
	g := ledger(t, `{"n":12,"started":"now","ended":""}`+"\n")
	_, err := Open("research", g, "python manjuel.py")
	if err == nil {
		t.Fatal("a world with an open sitting must be refused")
	}
	// the refusal has to name it, or the operator cannot act on it
	for _, want := range []string{"research", "12"} {
		if !contains(err.Error(), want) {
			t.Fatalf("refusal %q does not name %q", err, want)
		}
	}
}

// ---- a sitting whose process is gone (2026-09-22) ---------------------------

// deadPid is the pid of a process that has come and gone: this test binary run
// again with nothing to run, waited for.
func deadPid(t *testing.T) int {
	t.Helper()
	c := exec.Command(os.Args[0], "-test.run=^$")
	if err := c.Run(); err != nil {
		t.Fatalf("the short-lived child did not run: %v", err)
	}
	return c.ProcessState.Pid()
}

func openLine(n, pid int) string {
	return fmt.Sprintf(`{"n":%d,"started":"2026-09-22T10:00:00","ended":"","pid":%d}`+"\n", n, pid)
}

// ONLY A PROVABLE DEATH COUNTS. A crashed engine's line names a pid that is
// gone, and that line alone reads as orphaned; a live pid, no pid at all, a
// closed line and an unreadable ledger all leave the refusal standing.
func TestAnOrphanedSittingIsOnlyOneWhoseProcessIsProvablyGone(t *testing.T) {
	gone := deadPid(t)
	cases := []struct {
		name     string
		body     string
		orphaned bool
	}{
		{"an open line whose process is gone is orphaned", openLine(7, gone), true},
		{"an open line whose process is alive is not", openLine(7, os.Getpid()), false},
		{"an open line with no pid cannot be judged, so it is not",
			`{"n":7,"started":"x","ended":""}` + "\n", false},
		{"a closed line is not, whatever its pid",
			`{"n":7,"started":"x","ended":"y","pid":` + fmt.Sprint(gone) + `}` + "\n", false},
		{"an unreadable last line is not -- the C5 refusal keeps it",
			openLine(6, gone) + `{"n":7,"star`, false},
		{"no ledger is not", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pid, orphaned := SittingOrphaned(ledger(t, c.body))
			if orphaned != c.orphaned {
				t.Fatalf("orphaned = %v (pid %d), want %v", orphaned, pid, c.orphaned)
			}
			if orphaned && pid != gone {
				t.Fatalf("named pid %d, want %d -- the note must name the process", pid, gone)
			}
		})
	}
}

// THE WORLD A CRASH USED TO LOCK. Its last line is open and its process is
// gone, so the door opens it -- and the same line with a live process behind
// it is refused, naming the sitting, exactly as before.
func TestOpenOpensAWorldWhoseSittingsProcessIsGone(t *testing.T) {
	g := ledger(t, openLine(41, deadPid(t)))
	r := NewRegistry()
	e, err := r.Open("w", g, stubEngine(t))
	if err != nil {
		t.Fatalf("a world held by a dead process was refused: %v", err)
	}
	t.Cleanup(func() { _, _ = r.CloseOne(g) })
	if !e.Alive() {
		t.Fatal("the engine opened on it is not standing")
	}

	live := ledger(t, openLine(42, os.Getpid()))
	_, err = Open("w", live, stubEngine(t))
	if err == nil || !contains(err.Error(), "42") {
		t.Fatalf("a world whose sitting's process is alive must still be refused by name; got %v", err)
	}
}

func TestOpenNeverCreatesAWorld(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	if _, err := Open("w", missing, "python manjuel.py"); err == nil {
		t.Fatal("a missing ground must refuse")
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Fatal("a world was created; SITTING LAW 4 says a hand never invents one")
	}
}

func TestRingBufferKeepsTheTail(t *testing.T) {
	r := newRing(8)
	r.Write([]byte("0123456789abcdef"))
	if got := r.String(); got != "9abcdef" && got != "89abcdef" {
		t.Fatalf("ring kept %q; it must keep the TAIL, which is where the error is", got)
	}
}

// ---- the stale check: the engine's code, and nothing else -------------------

// What a running engine holds: the script, the package, and the law the gate
// loads -- law/pen/jesster.py is imported by name and kept for the process.
var engineLoads = []string{
	"manjuel.py",
	"manjuel/pipeline.py",
	"law/law.py",
	"law/pen/jesster.py",
}

// .py a ground carries that no engine imports.
var notEngineCode = []string{
	"atlas/tools/cut_vectors.py",
	"atlas/line/engine/manjuel_ask.py",
	"skills/helper.py",
	"agent_workspace/scratch.py",
	"worlds/w/tool.py",
	"manjuel/__pycache__/pipeline.py",
}

var codeBase = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// codeGround writes every file above, all stamped codeBase, and returns the
// command that runs its manjuel.py. Quoted, so a temp path with a space holds.
func codeGround(t *testing.T) (root, cmd string) {
	t.Helper()
	root = t.TempDir()
	for _, rel := range append(append([]string{}, engineLoads...), notEngineCode...) {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("# code\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		stamp(t, p, codeBase)
	}
	return root, `python "` + filepath.Join(root, "manjuel.py") + `"`
}

func stamp(t *testing.T, p string, at time.Time) {
	t.Helper()
	if err := os.Chtimes(p, at, at); err != nil {
		t.Fatal(err)
	}
}

// 30 of the 64 .py the old walk counted on this ground were atlas/'s tools and
// tests. An edit to code no engine runs must never read as "restart".
func TestTheStaleCheckIgnoresCodeNoEngineLoads(t *testing.T) {
	for _, rel := range notEngineCode {
		t.Run(rel, func(t *testing.T) {
			root, cmd := codeGround(t)
			stamp(t, filepath.Join(root, filepath.FromSlash(rel)), codeBase.Add(time.Hour))
			if when, what := CodeChanged(cmd); !when.Equal(codeBase) {
				t.Fatalf("reported %s changed at %s; an engine never runs %s", what, when, rel)
			}
		})
	}
}

// law/ is what a walk of manjuel/ alone would have dropped, and the pen inside
// it is one level down -- the folder rule is for the top of the ground only.
func TestTheStaleCheckSeesEveryPlaceAnEngineLoadsCodeFrom(t *testing.T) {
	for _, rel := range engineLoads {
		t.Run(rel, func(t *testing.T) {
			root, cmd := codeGround(t)
			p := filepath.Join(root, filepath.FromSlash(rel))
			stamp(t, p, codeBase.Add(time.Hour))
			if when, what := CodeChanged(cmd); what != p || !when.Equal(codeBase.Add(time.Hour)) {
				t.Fatalf("reported %q at %s; wanted %s, which a running engine holds", what, when, rel)
			}
		})
	}
}

// The row the operator reads: stale over the engine's code, not over the ground.
func TestAnEngineIsStaleOverItsCodeNotOverTheGround(t *testing.T) {
	root, cmd := codeGround(t)
	e := &Engine{Started: codeBase.Add(time.Hour), CoreCmd: cmd}
	stamp(t, filepath.Join(root, "atlas", "tools", "cut_vectors.py"), codeBase.Add(2*time.Hour))
	if stale, _, what := e.Stale(); stale {
		t.Fatalf("stale over %s, which the engine never loaded", what)
	}
	stamp(t, filepath.Join(root, "law", "pen", "jesster.py"), codeBase.Add(2*time.Hour))
	if stale, _, what := e.Stale(); !stale || filepath.Base(what) != "jesster.py" {
		t.Fatalf("stale=%v over %q; the pen changed after the engine started", stale, what)
	}
}

// ---- a real process: whether the engine is still there --------------------

const stubEngineEnv = "ATLAS_STUB_ENGINE"

// TestStubEngineProcess is not a stroke. It is the engine the strokes below
// spawn: this test binary run again as a child, speaking serve.py's wire --
// `opened` at start, a delivery per objective, `closed` on close. The
// objective "die" ends the process mid-turn with no terminal event, the way a
// crash does. Run any other way, it skips.
func TestStubEngineProcess(t *testing.T) {
	if os.Getenv(stubEngineEnv) != "1" {
		t.Skip("the stub engine runs only as the child of a stroke that spawns it")
	}
	out := bufio.NewWriter(os.Stdout)
	say := func(ev map[string]any) {
		b, _ := json.Marshal(ev)
		out.Write(append(b, '\n'))
		out.Flush()
	}
	say(map[string]any{"event": "opened", "sitting": "1", "session": "stub"})
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		var row map[string]any
		if json.Unmarshal(in.Bytes(), &row) != nil {
			continue
		}
		switch row["cmd"] {
		case "close":
			say(map[string]any{"event": "closed"})
			os.Exit(0)
		case "objective":
			if row["text"] == "die" {
				os.Exit(3)
			}
			say(map[string]any{"event": "delivery", "text": row["text"]})
		}
	}
	os.Exit(0)
}

// stubEngine is the command that runs TestStubEngineProcess as a child.
func stubEngine(t *testing.T) string {
	t.Helper()
	t.Setenv(stubEngineEnv, "1")
	return `"` + os.Args[0] + `" -test.run=^TestStubEngineProcess$ --`
}

// I5 SAID A DEAD ENGINE MUST NOT BE HANDED BACK, and the check written for it
// could not see a death. A real engine is told to die mid-turn; it must read
// as gone to Alive, and to the registry every /run/state asks.
func TestAnEngineThatDiesIsNotHandedBack(t *testing.T) {
	ground := t.TempDir()
	r := NewRegistry()
	e, err := r.Open("w", ground, stubEngine(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = r.CloseOne(ground) })
	if !e.Alive() {
		t.Fatal("a standing engine reads as dead")
	}
	if _, ok := r.Get(ground); !ok {
		t.Fatal("the registry does not hand back a standing engine")
	}
	if _, err := e.Run("die", "", "", Head{}, nil); err == nil {
		t.Fatal("a turn whose engine died reported success")
	}
	deadline := time.Now().Add(10 * time.Second)
	for e.Alive() {
		if time.Now().After(deadline) {
			t.Fatal("the engine's process is gone and Alive still says it is standing")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if _, ok := r.Get(ground); ok {
		t.Fatal("the registry handed back an engine whose process is gone")
	}
}

// A CLOSE IS STILL WAITED FOR, by the same waiter Alive reads: the engine is
// asked to close, goes, and the close reports its toll paid inside the grace
// rather than a minute later.
func TestACloseIsSeenByTheProcesssOneWaiter(t *testing.T) {
	e, err := Open("w", t.TempDir(), stubEngine(t))
	if err != nil {
		t.Fatal(err)
	}
	type closing struct {
		paid bool
		err  error
	}
	got := make(chan closing, 1)
	go func() {
		paid, err := e.Close()
		got <- closing{paid, err}
	}()
	select {
	case c := <-got:
		if c.err != nil || !c.paid {
			t.Fatalf("close = (%v, %v); an engine that went when asked has paid its toll", c.paid, c.err)
		}
	case <-time.After(20 * time.Second):
		t.Fatal("the close never returned: nothing saw the process go")
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
