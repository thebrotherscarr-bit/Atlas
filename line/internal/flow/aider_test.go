package flow

// Strokes for the `aider` step kind (2026-10-05, WHAT'S LEFT H17, his ruling on B21: "yes"). An `aider` node is an attempt by Aider on the files it is
// handed, reached only through a gate whose grants name the door's Aider tool; the engine behind it is THE LINE's (Aider through the door, the council's
// own turn when Aider cannot take it) and is held by the door's own strokes. What is held HERE is the flow's side of the wire: what a node must carry,
// who may reach it, what it asks of the engine and what is scored of its answer.

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"atlas/line/internal/play"
)

// aidingEngine answers an `aider` node from a script and records what the node asked.
type aidingEngine struct {
	stubEngine
	asked    []aiderAsk
	answers_ []string // what Aider answers, one per call; the last repeats
	failures int      // the first calls that answer nothing at all
	notReady error
}

type aiderAsk struct{ Instruction, Files, Feed string }

func (a *aidingEngine) Aider(_ context.Context, instruction, files, feed string) (string, error) {
	a.asked = append(a.asked, aiderAsk{instruction, files, feed})
	if len(a.asked) <= a.failures {
		return "", fmt.Errorf("the door is not open")
	}
	if len(a.answers_) == 0 {
		return "aided", nil
	}
	i := len(a.asked) - 1 - a.failures
	if i >= len(a.answers_) {
		i = len(a.answers_) - 1
	}
	return a.answers_[i], nil
}

func (a *aidingEngine) Ready() error { return a.notReady }

// aiderSpec is the smallest honest shape: a plan, the gate the hand crosses (which grants Aider), and the attempt.
func aiderSpec() Spec {
	return Spec{Name: "aided", BudgetS: 600,
		Nodes: []Node{
			{Name: "brief", Kind: "ask", Question: "plan {{objective}}"},
			{Name: "open", Kind: "gate", Title: "open the line?", Grants: []string{AiderTool}},
			{Name: "attempt", Kind: "aider", Question: "{{objective}}", Files: "{{out_brief}}"},
		},
		Edges: []Edge{
			{From: "brief", To: "open", When: "always"},
			{From: "open", To: "attempt", When: "pass"},
		}}
}

// evidenced is an answer WITH the machine's block, as THE LINE writes one: prose first, the marker, the door's own lines.
func evidenced(prose string, lines ...string) string {
	return prose + "\n\n" + play.ToolVerdictHead + "\n" + strings.Join(lines, "\n")
}

func TestAnAiderNodeNeedsAnInstructionAndFilesAndNothingElseNamesFiles(t *testing.T) {
	if _, err := Validate(aiderSpec()); err != nil {
		t.Fatalf("the smallest honest aider spec was refused: %v", err)
	}
	s := aiderSpec()
	s.Nodes[2].Question = "  "
	if _, err := Validate(s); err == nil || !strings.Contains(err.Error(), "has no instruction") {
		t.Fatalf("an aider node with no instruction was not refused by name: %v", err)
	}
	s = aiderSpec()
	s.Nodes[2].Files = ""
	if _, err := Validate(s); err == nil || !strings.Contains(err.Error(), "names no files") {
		t.Fatalf("an aider node that names no files was not refused by name: %v", err)
	}
	for _, kind := range []string{"run", "ask", "memory"} {
		s = aiderSpec()
		s.Nodes = append(s.Nodes, Node{Name: "x", Kind: kind, Question: "Q", Files: "manjuel/mathkit.py"})
		s.Edges = append(s.Edges, Edge{From: "attempt", To: "x", When: "always"})
		if _, err := Validate(s); err == nil || !strings.Contains(err.Error(), "names no files for Aider") {
			t.Fatalf("a %s node that names files was not refused: %v", kind, err)
		}
	}
	if !Kinds["aider"] {
		t.Fatal("aider is not in the closed set of kinds")
	}
}

// ONLY THROUGH A GATE THE HAND CROSSED, AND THE NEAREST GATE ON EVERY PATH IS THE ONE THAT COUNTS.
func TestAnAiderNodeIsReachedOnlyThroughAGateThatGrantsTheTool(t *testing.T) {
	gate := func(name string, grants ...string) Node {
		return Node{Name: name, Kind: "gate", Title: "t", Grants: grants}
	}
	ask := func(name string) Node { return Node{Name: name, Kind: "ask", Question: "Q"} }
	aider := Node{Name: "w", Kind: "aider", Question: "do it", Files: "a.py"}
	edge := func(f, to string) Edge { return Edge{From: f, To: to, When: "always"} }
	refused := func(why string, nodes []Node, edges []Edge) {
		t.Helper()
		_, err := Validate(Spec{Name: "g", BudgetS: 600, Nodes: nodes, Edges: edges})
		if err == nil || !strings.Contains(err.Error(), "without passing a gate that grants") {
			t.Fatalf("%s: not refused as an aider node no gate authorised: %v", why, err)
		}
	}
	allowed := func(why string, nodes []Node, edges []Edge) {
		t.Helper()
		if _, err := Validate(Spec{Name: "g", BudgetS: 600, Nodes: nodes, Edges: edges}); err != nil {
			t.Fatalf("%s: refused: %v", why, err)
		}
	}
	refused("no gate at all", []Node{ask("a"), aider}, []Edge{edge("a", "w")})
	refused("the node is the start", []Node{aider}, nil)
	refused("a gate that grants nothing", []Node{ask("a"), gate("g"), aider}, []Edge{edge("a", "g"), edge("g", "w")})
	refused("a gate that grants another tool only", []Node{ask("a"), gate("g", "git_branch"), aider}, []Edge{edge("a", "g"), edge("g", "w")})
	refused("a gate that grants it, then one that grants nothing", []Node{ask("a"), gate("g", AiderTool), gate("g2"), aider},
		[]Edge{edge("a", "g"), edge("g", "g2"), edge("g2", "w")})
	refused("one path through the granting gate and one around it",
		[]Node{ask("a"), gate("g", AiderTool), ask("c"), ask("d"), aider},
		[]Edge{edge("a", "g"), edge("g", "c"), edge("c", "w"), edge("a", "d"), edge("d", "w")})
	allowed("the granting gate, then work in between",
		[]Node{ask("a"), gate("g", AiderTool), ask("c"), aider}, []Edge{edge("a", "g"), edge("g", "c"), edge("c", "w")})
	allowed("the granting gate among other grants",
		[]Node{ask("a"), gate("g", "git_branch", AiderTool), aider}, []Edge{edge("a", "g"), edge("g", "w")})
	allowed("both paths through a granting gate",
		[]Node{ask("a"), gate("g", AiderTool), gate("h", AiderTool), aider},
		[]Edge{edge("a", "g"), edge("a", "h"), edge("g", "w"), edge("h", "w")})
}

func TestAnAiderNodeAsksTheEngineWithTheRenderedInstructionAndFiles(t *testing.T) {
	home := t.TempDir()
	eng := &aidingEngine{stubEngine: stubEngine{answers: map[string]string{"plan add clamp": "manjuel/mathkit.py"}}}
	res, err := Run(home, eng, aiderSpec(), map[string]string{"objective": "add clamp"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused || res.PausedNode != "open" {
		t.Fatalf("the gate before Aider did not pause: %s at %q", res.Verdict, res.PausedNode)
	}
	if len(eng.asked) != 0 {
		t.Fatalf("Aider was asked before the hand crossed the gate: %v", eng.asked)
	}
	res, err = Resume(home, eng, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s after the gate was crossed", res.Verdict)
	}
	if len(eng.asked) != 1 || eng.asked[0] != (aiderAsk{"add clamp", "manjuel/mathkit.py", ""}) {
		t.Fatalf("the node did not ask with its rendered instruction and files and an empty failed pass: %v", eng.asked)
	}
	lines, err := runLog(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range lines {
		if l["kind"] == "node" && l["node"] == "attempt" && l["status"] == "ok" {
			found = l["nkind"] == "aider" && l["output"] == "aided"
		}
	}
	if !found {
		t.Fatalf("the aider node's answer is not in the record as an aider node's: %v", lines)
	}
}

// AN EVAL OVER AN AIDER NODE SCORES THE MACHINE'S LINES AND NOTHING ELSE, as over a run node: the prose a fall-back turn wrote can say anything.
func TestAnEvalOverAnAiderNodeScoresTheMachinesLinesAndNothingElse(t *testing.T) {
	spec := func() Spec {
		s := aiderSpec()
		s.Nodes = append(s.Nodes,
			Node{Name: "changed", Kind: "eval", Ref: "attempt", Expected: "line of work `", Match: "contains"},
			Node{Name: "yes", Kind: "ask", Question: "on a line"},
			Node{Name: "no", Kind: "ask", Question: "not on a line"})
		s.Edges = append(s.Edges,
			Edge{From: "attempt", To: "changed", When: "always"},
			Edge{From: "changed", To: "yes", When: "pass"},
			Edge{From: "changed", To: "no", When: "fail"})
		return s
	}
	fire := func(answer string) []string {
		home := t.TempDir()
		eng := &aidingEngine{answers_: []string{answer}}
		res, err := Run(home, eng, spec(), map[string]string{"objective": "add clamp"})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Resume(home, eng, res.Run, "continue"); err != nil {
			t.Fatal(err)
		}
		return eng.calls
	}
	has := func(calls []string, q string) bool {
		for _, c := range calls {
			if c == "ask:"+q {
				return true
			}
		}
		return false
	}
	if c := fire("Edited it on line of work `x`, honestly."); !has(c, "not on a line") {
		t.Fatalf("a seat's own words, with no machine block, passed a check for them: %v", c)
	}
	if c := fire(evidenced("It went fine.", "aider_run: Aider changed 1 file on line of work `x`")); !has(c, "on a line") {
		t.Fatalf("the machine's own line, in the block, did not pass: %v", c)
	}
	if c := fire(evidenced("Edited it on line of work `x`, honestly.", "aider_run: Aider changed nothing")); !has(c, "not on a line") {
		t.Fatalf("the prose above the block passed a check the block fails: %v", c)
	}
}

// A FAILED PASS IS CARRIED TO THE NODE THAT IS RETURNED TO, and the ceiling is read: two returns, then a fail.
func TestAFailedPassIsCarriedToAnAiderNodeThatDeclaresLoopsAndTheCeilingHolds(t *testing.T) {
	s := aiderSpec()
	s.Nodes[2].Loops = 2
	s.Nodes = append(s.Nodes, Node{Name: "changed", Kind: "eval", Ref: "attempt", Expected: "CHANGED", Match: "contains"})
	s.Edges = append(s.Edges,
		Edge{From: "attempt", To: "changed", When: "always"},
		Edge{From: "changed", To: "attempt", When: "fail"})
	home := t.TempDir()
	eng := &aidingEngine{answers_: []string{evidenced("no", "aider_run: Aider changed nothing")}}
	res, err := Run(home, eng, s, map[string]string{"objective": "add clamp"})
	if err != nil {
		t.Fatal(err)
	}
	res, err = Resume(home, eng, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictFail {
		t.Fatalf("a check that never passed ended %s, not FAIL", res.Verdict)
	}
	if len(eng.asked) != 3 {
		t.Fatalf("a node declaring loops: 2 was asked %d times, want 3 (the first pass and two returns)", len(eng.asked))
	}
	if eng.asked[0].Feed != "" {
		t.Fatalf("the first pass was handed a failed pass: %q", eng.asked[0].Feed)
	}
	if !strings.Contains(eng.asked[1].Feed, "Aider changed nothing") || eng.asked[2].Feed == "" {
		t.Fatalf("a sent-back pass did not carry what the last one failed with: %q / %q", eng.asked[1].Feed, eng.asked[2].Feed)
	}
}

// A RETRY ANSWERS AN ERROR: the door not answering at all is retried, the record says so, and the run survives.
func TestAnAiderNodeThatErrorsIsRetriedLikeAnyOtherAndFailsTheRunWhenItKeepsFailing(t *testing.T) {
	s := aiderSpec()
	s.Nodes[2].Retries = 1
	home := t.TempDir()
	eng := &aidingEngine{failures: 1}
	res, err := Run(home, eng, s, map[string]string{"objective": "x"})
	if err != nil {
		t.Fatal(err)
	}
	if res, err = Resume(home, eng, res.Run, "continue"); err != nil || res.Verdict != VerdictComplete || len(eng.asked) != 2 {
		t.Fatalf("one dead call and one retry did not land: %v %v %d", err, res.Verdict, len(eng.asked))
	}
	home2 := t.TempDir()
	dead := &aidingEngine{failures: 5}
	res, err = Run(home2, dead, s, map[string]string{"objective": "x"})
	if err != nil {
		t.Fatal(err)
	}
	if res, err = Resume(home2, dead, res.Run, "continue"); err != nil || res.Verdict != VerdictFail {
		t.Fatalf("a node that kept failing did not fail the run: %v %v", err, res.Verdict)
	}
}

// A RESUME NEEDS THE ENGINE FOR AN AIDER NODE AS FOR A RUN NODE: its fall-back is the council's own turn.
func TestAResumeWithNoEngineStandingIsRefusedForAnAiderNodeToo(t *testing.T) {
	home := t.TempDir()
	eng := &aidingEngine{}
	res, err := Run(home, eng, aiderSpec(), map[string]string{"objective": "x"})
	if err != nil {
		t.Fatal(err)
	}
	eng.notReady = fmt.Errorf("no engine is open")
	_, err = Resume(home, eng, res.Run, "continue")
	if err == nil || !strings.Contains(err.Error(), "`aider` node") || !strings.Contains(err.Error(), "attempt") {
		t.Fatalf("a resume with no engine standing was not refused in the aider node's name: %v", err)
	}
	if len(eng.asked) != 0 {
		t.Fatalf("Aider was asked with no engine standing: %v", eng.asked)
	}
	lines, _ := runLog(home, res.Run)
	for _, l := range lines {
		if l["kind"] == "resumed" {
			t.Fatalf("the refused resume was written down as a resume: %v", l)
		}
	}
}

func TestTheBareEngineRefusesAnAiderNode(t *testing.T) {
	_, err := Production(t.TempDir()).Aider(context.Background(), "do it", "a.py", "")
	if err == nil || !strings.Contains(err.Error(), "no door wired") {
		t.Fatalf("the bare engine answered an aider node, or refused it in the wrong words: %v", err)
	}
}
