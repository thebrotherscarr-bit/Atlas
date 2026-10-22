package tools

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"atlas/line/internal/tenant"
)

// notesTenant is a bare world with the notes kernel declared in its .env, and
// the kernel itself stood in for: every call is recorded and answered with reply.
func notesTenant(t *testing.T, declared string, reply spawnResult) (*tenant.Registry, *[][]string) {
	t.Helper()
	t.Setenv(notesDeclared, "")
	home := t.TempDir()
	if declared != "" {
		if err := os.WriteFile(filepath.Join(home, ".env"), []byte("# his\r\nOTHER=1\r\n"+notesDeclared+"="+declared+"\r\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	tr := tenant.NewRegistry()
	if err := tr.Add("t", home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	calls := &[][]string{}
	old := notesRun
	notesRun = func(argv []string, dir string) spawnResult {
		if dir != home {
			t.Errorf("the kernel ran in %q, not the world's Home", dir)
		}
		*calls = append(*calls, argv)
		return reply
	}
	t.Cleanup(func() { notesRun = old })
	return tr, calls
}

const kernelCmd = `python worlds/jesster/notes/cli.py --vault "worlds/my notes" mcp`

func TestNotesReadsThroughTheDeclaredKernel(t *testing.T) {
	tr, calls := notesTenant(t, kernelCmd, spawnResult{Stdout: `[{"id":"whats-left-a","title":"What's left: A"}]`})
	reg := Build(tr, Options{})
	out, err := reg.Call(tr, "notes", map[string]any{"action": "list", "tag": "whats-left", "limit": float64(3), "project": "t"}, Caller{Name: "prove"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"python", "worlds/jesster/notes/cli.py", "--vault", "worlds/my notes", "call", "notes.list", `{"limit":3,"tag":"whats-left"}`}
	if !reflect.DeepEqual((*calls)[0], want) {
		t.Fatalf("argv\n got %q\nwant %q", (*calls)[0], want)
	}
	var got struct {
		Action string
		Result []map[string]string
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil || got.Action != "list" || got.Result[0]["id"] != "whats-left-a" {
		t.Fatalf("the kernel's answer did not come back whole: %v %s", err, out)
	}
}

func TestNotesActionsMapOnlyToReaders(t *testing.T) {
	tr, calls := notesTenant(t, kernelCmd, spawnResult{Stdout: `{}`})
	reg := Build(tr, Options{})
	call := func(a map[string]any) error {
		a["project"] = "t"
		_, err := reg.Call(tr, "notes", a, Caller{Name: "prove"})
		return err
	}
	for _, c := range []struct {
		args   map[string]any
		method string
	}{
		{map[string]any{"q": "voice memos"}, "notes.search"},
		{map[string]any{}, "notes.list"},
		{map[string]any{"action": "list", "tag": "whats-left", "full": true}, "notes.list"},
		{map[string]any{"action": "get", "id": "handoff-2026-10-21"}, "notes.get"},
		{map[string]any{"action": "context", "q": "what next"}, "context.for"},
	} {
		if err := call(c.args); err != nil {
			t.Fatalf("%v: %v", c.args, err)
		}
		if got := (*calls)[len(*calls)-1][5]; got != c.method {
			t.Errorf("%v ran %s, want %s", c.args, got, c.method)
		}
	}
	if got := (*calls)[2][6]; got != `{"full":true,"limit":50,"tag":"whats-left"}` {
		t.Errorf("full did not reach the kernel: %s", got)
	}
	for args, why := range map[string]string{"create": "no notes action", "notes.remove": "no notes action"} {
		if err := call(map[string]any{"action": args}); err == nil || !strings.Contains(err.Error(), why) {
			t.Errorf("action %q was not refused by name: %v", args, err)
		}
	}
	if err := call(map[string]any{"action": "get"}); err == nil || !strings.Contains(err.Error(), "needs id") {
		t.Errorf("get with no id: %v", err)
	}
	tool, _ := reg.Get("notes")
	if tool.Writes || tool.ServiceOnly {
		t.Errorf("notes must read and be open to every reader: %+v", tool)
	}
}

func TestNotesUndeclaredOrNotTheKernelIsRefused(t *testing.T) {
	tr, calls := notesTenant(t, "", spawnResult{})
	reg := Build(tr, Options{})
	if _, err := reg.Call(tr, "notes", map[string]any{"project": "t"}, Caller{Name: "prove"}); err == nil || !strings.Contains(err.Error(), "not declared") {
		t.Fatalf("an undeclared second brain: %v", err)
	}
	tr2, _ := notesTenant(t, "python cli.py rpc", spawnResult{})
	reg2 := Build(tr2, Options{})
	if _, err := reg2.Call(tr2, "notes", map[string]any{"project": "t"}, Caller{Name: "prove"}); err == nil || !strings.Contains(err.Error(), "will not guess") {
		t.Fatalf("a declaration that is not the MCP command: %v", err)
	}
	if len(*calls) != 0 {
		t.Fatal("something ran with no kernel declared")
	}
}

func TestNotesSaysWhatTheKernelSaid(t *testing.T) {
	tr, _ := notesTenant(t, kernelCmd, spawnResult{Stderr: "  refused (not_found): no note 'x'\n", Err: errors.New("exit status 2")})
	reg := Build(tr, Options{})
	_, err := reg.Call(tr, "notes", map[string]any{"action": "get", "id": "x", "project": "t"}, Caller{Name: "prove"})
	if err == nil || !strings.Contains(err.Error(), "no note 'x'") {
		t.Fatalf("the kernel's own words were lost: %v", err)
	}
	tr2, _ := notesTenant(t, kernelCmd, spawnResult{Stdout: "  vault C:\\notes\n"})
	reg2 := Build(tr2, Options{})
	if _, err := reg2.Call(tr2, "notes", map[string]any{"project": "t"}, Caller{Name: "prove"}); err == nil || !strings.Contains(err.Error(), "not JSON") {
		t.Fatalf("a non-JSON answer was passed off: %v", err)
	}
}

func TestNotesDeclarationFromTheGroundOrTheEnvironment(t *testing.T) {
	t.Setenv(notesDeclared, "")
	ground := t.TempDir()
	world := filepath.Join(ground, "world")
	if err := os.MkdirAll(world, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ground, ".env"), []byte(notesDeclared+"='py k.py mcp'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if argv, why := notesCommand(world); why != "" || !reflect.DeepEqual(argv, []string{"py", "k.py"}) {
		t.Fatalf("one level up: %q %s", argv, why)
	}
	if argv, _ := notesCommand(filepath.Join(world, "deeper")); argv != nil {
		t.Fatalf("two levels up is outside the bound: %q", argv)
	}
	t.Setenv(notesDeclared, `C:\Python\python.exe C:\R\worlds\jesster\notes\cli.py --vault C:\R\worlds\notes mcp`)
	argv, _ := notesCommand(world)
	if argv[0] != `C:\Python\python.exe` || argv[len(argv)-1] != `C:\R\worlds\notes` {
		t.Fatalf("the environment wins, backslashes kept: %q", argv)
	}
}
