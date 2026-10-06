// Strokes for the release as one workflow (2026-10-05, his word: "make the
// release one workflow"): the door's release_step, and a flow's `tool` node
// calling the door it was fired through.
package tools

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"atlas/line/internal/flow"
	"atlas/line/internal/tenant"
)

func releaseWorld(t *testing.T) (*tenant.Registry, string) {
	t.Helper()
	tr := tenant.NewRegistry()
	home := t.TempDir()
	if err := tr.Add("w", home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("w"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ATLAS_BIN", filepath.Join(home, "NO-SUCH-SPINE.exe"))
	return tr, home
}

// release_step's refusals by name, its declaration, and a run against a
// stand-in tests/cut.py: a step that answers comes back whole, its mark passed
// on, and a step that refused is an error carrying its words, never an answer.
func TestTheReleaseStepsRunFromTheCore(t *testing.T) {
	tr, home := releaseWorld(t)
	glass := Caller{Name: "glass", Service: true}
	noCore := Build(tr, Options{})
	if _, err := noCore.Call(tr, "release_step", map[string]any{"step": "index"}, glass); err == nil ||
		!strings.Contains(err.Error(), "started without --manjuel") {
		t.Fatalf("with no core command the tool must say so: %v", err)
	}
	reg := Build(tr, Options{CoreCmd: "python"})
	for _, c := range []struct {
		args map[string]any
		want string
	}{
		{map[string]any{"step": "tag"}, `"tag" is none of them`},
		{map[string]any{"step": "bump"}, "names its mark as vX.Y.Z"},
		{map[string]any{"step": "bump", "mark": "0.2.2"}, `"0.2.2" is not one`},
		{map[string]any{"step": "record", "mark": "v0.2"}, `"v0.2" is not one`},
		{map[string]any{"step": "index"}, "carries no tests/cut.py"},
	} {
		if _, err := reg.Call(tr, "release_step", c.args, glass); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: want a refusal naming %q, got %v", c.args, c.want, err)
		}
	}
	if tool, ok := reg.Get("release_step"); !ok || !tool.Writes {
		t.Fatal("release_step must be declared a writer: it moves the pins, the changelogs and the record")
	}
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("no python on the PATH to run the stand-in with")
	}
	if err := os.MkdirAll(filepath.Join(home, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	stand := "import sys\nprint('STEP', ' '.join(sys.argv[1:]))\nsys.exit(1 if sys.argv[1] == 'check' else 0)\n"
	if err := os.WriteFile(filepath.Join(home, "tests", "cut.py"), []byte(stand), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := reg.Call(tr, "release_step", map[string]any{"step": "bump", "mark": "v9.9.9"}, glass)
	if err != nil || !strings.Contains(out, "STEP bump v9.9.9") {
		t.Fatalf("a step that answers must come back whole, its mark passed on: %q %v", out, err)
	}
	out, err = reg.Call(tr, "release_step", map[string]any{"step": "index", "mark": "v9.9.9"}, glass)
	if err != nil || !strings.Contains(out, "STEP index") || strings.Contains(out, "v9.9.9") {
		t.Fatalf("index names no mark, whatever it is handed: %q %v", out, err)
	}
	if _, err := reg.Call(tr, "release_step", map[string]any{"step": "check", "mark": "v9.9.9"}, glass); err == nil ||
		!strings.Contains(err.Error(), "release step check refused") || !strings.Contains(err.Error(), "STEP check v9.9.9") {
		t.Fatalf("a step that refused must be an error carrying its words: %v", err)
	}
}

// The call is judged by the door: a council with no door refuses, and so do a
// tool the door does not carry, a flow verb and a writer with no grant; under
// the gate's grant the call reaches the tool itself.
func TestAToolNodesCallIsJudgedByTheDoor(t *testing.T) {
	tr, home := releaseWorld(t)
	reg := Build(tr, Options{})
	ctx := context.Background()
	if _, err := council(home).Tool(ctx, "muster", nil); err == nil || !strings.Contains(err.Error(), "no door") {
		t.Fatalf("a council with no door must refuse a tool node: %v", err)
	}
	c := councilAt(home, map[string]any{registryKey: reg, CallerKey: Caller{Name: "glass", Service: true}})
	for name, want := range map[string]string{
		"no_such_tool": "not a tool this door carries",
		"flow_run":     "does not fire a flow",
		"release_step": "does not grant it",
	} {
		if _, err := c.Tool(ctx, name, map[string]string{"step": "index"}); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: want a refusal naming %q, got %v", name, want, err)
		}
	}
	granted := c.(councilEngine).WithHand("r", "g", []string{"release_step"})
	if _, err := granted.Tool(ctx, "release_step", map[string]string{"step": "nope"}); err == nil ||
		!strings.Contains(err.Error(), `"nope" is none of them`) {
		t.Fatalf("under the gate's grant the call must reach the tool itself: %v", err)
	}
	out, err := c.Tool(ctx, "muster", nil)
	if err != nil || !strings.Contains(out, "carried projects") {
		t.Fatalf("a reader needs no grant: %q %v", out, err)
	}
}

// THE WIRE: a flow fired through the door calls the door. A reader answers the
// run; a writer reached past no gate that grants it fails the run, even fired
// from his glass. Unplug councilAt from flow_run and the first run fails.
// A DECISION IS REFUSED WHILE A RUN IS MOVING (2026-10-05): waiting on the
// world's flow lock instead, it answered the next gate, unseen.
func TestADecisionIsRefusedWhileARunIsMoving(t *testing.T) {
	tr, _ := releaseWorld(t)
	reg := Build(tr, Options{})
	glass := Caller{Name: "glass", Service: true}
	tn, err := tr.Resolve("")
	if err != nil {
		t.Fatal(err)
	}
	lock := flowLock(tn.Home)
	lock.Lock()
	_, err = reg.Call(tr, "flow_resume", map[string]any{"run": "f-20261006-000000-00000000", "decision": "continue"}, glass)
	lock.Unlock()
	if err == nil || !strings.Contains(err.Error(), "a run is moving on this world") {
		t.Fatalf("a decision made while a run moves must be refused, not held for the next gate: %v", err)
	}
	if _, err := reg.Call(tr, "flow_resume", map[string]any{"run": "f-20261006-000000-00000000", "decision": "continue"}, glass); err == nil ||
		strings.Contains(err.Error(), "a run is moving") {
		t.Fatalf("with no run moving, the decision must reach the run itself, refused there for a run that is not: %v", err)
	}
}

func TestAFlowsToolNodeCallsTheDoorItWasFiredThrough(t *testing.T) {
	tr, home := releaseWorld(t)
	glass := Caller{Name: "glass", Service: true}
	reg := Build(tr, Options{})
	save := func(name string, nodes []flow.Node, edges []flow.Edge) {
		t.Helper()
		if _, err := flow.Save(home, flow.Spec{Name: name, Nodes: nodes, Edges: edges}); err != nil {
			t.Fatal(err)
		}
	}
	save("roll", []flow.Node{
		{Name: "roll", Kind: "tool", Tool: "muster"},
		{Name: "said", Kind: "eval", Ref: "roll", Expected: "carried projects", Match: "contains"},
	}, []flow.Edge{{From: "roll", To: "said", When: "always"}})
	out, err := reg.Call(tr, "flow_run", map[string]any{"name": "roll"}, glass)
	if err != nil || !strings.Contains(out, flow.VerdictComplete) {
		t.Fatalf("a reader called by a tool node must answer the run: %q %v", out, err)
	}
	save("ungranted", []flow.Node{
		{Name: "bump", Kind: "tool", Tool: "release_step", Args: map[string]string{"step": "index"}},
	}, nil)
	out, err = reg.Call(tr, "flow_run", map[string]any{"name": "ungranted"}, glass)
	if err != nil || !strings.Contains(out, flow.VerdictFail) {
		t.Fatalf("a writer past no granting gate must fail the run, even fired from the glass: %q %v", out, err)
	}
}
