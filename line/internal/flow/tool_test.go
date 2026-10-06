// Strokes for the `tool` node (2026-10-05, his word: "make the release one
// workflow"). A tool node calls one of the door's own tools by name, its
// arguments rendered like a question; the flow law here judges its shape, and
// the door judges the call (internal/tools, councilEngine.Tool).
package flow

import (
	"context"
	"errors"
	"sort"
	"strings"
	"testing"
)

// The stand-in engines take a tool node's call as they take the others: the
// stub writes it down with its arguments in name order and answers from its
// canned map (an answer that begins "ERR " is a refusal), and the head and hand
// engines write it down beside the head or the crossing it ran under.
func (s *stubEngine) Tool(_ context.Context, name string, args map[string]string) (string, error) {
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	call := name
	for _, k := range keys {
		call += " " + k + "=" + args[k]
	}
	s.calls = append(s.calls, "tool:"+call)
	if a, ok := s.answers["tool:"+name]; ok {
		if strings.HasPrefix(a, "ERR ") {
			return "", errors.New(strings.TrimPrefix(a, "ERR "))
		}
		return a, nil
	}
	return "stub-tool:" + call, nil
}

func (h headEngine) Tool(_ context.Context, name string, _ map[string]string) (string, error) {
	h.log.turns = append(h.log.turns, "tool:"+name+"@"+headName(h.head))
	return "called " + name, nil
}

func (h handEngine) Tool(_ context.Context, name string, _ map[string]string) (string, error) {
	h.log.turns = append(h.log.turns, "tool:"+name+"@"+h.gate+":"+strings.Join(h.grants, ","))
	return "called " + name, nil
}

func toolSpec(nodes []Node, edges []Edge) Spec {
	return Spec{Name: "tooled", BudgetS: 600, Nodes: nodes, Edges: edges}
}

// The call carries every argument rendered -- the run's inputs and an earlier
// node's answer -- and the run reads what the tool said.
func TestAToolNodeCallsItsToolWithItsArgumentsRendered(t *testing.T) {
	eng := &stubEngine{answers: map[string]string{}}
	s := toolSpec([]Node{
		{Name: "plan", Kind: "ask", Question: "Q"},
		{Name: "bump", Kind: "tool", Tool: "release_step",
			Args: map[string]string{"step": "bump", "mark": "{{mark}}", "why": "{{out_plan}}"}},
		{Name: "said", Kind: "eval", Ref: "bump", Expected: "mark=v9.9.9", Match: "contains"},
	}, []Edge{{From: "plan", To: "bump", When: "always"}, {From: "bump", To: "said", When: "always"}})
	res, err := Run(t.TempDir(), eng, s, map[string]string{"mark": "v9.9.9"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictComplete {
		t.Fatalf("a tool node whose call answers must complete the run: %s (calls %v)", res.Verdict, eng.calls)
	}
	want := "tool:release_step mark=v9.9.9 step=bump why=stub-answer"
	for _, c := range eng.calls {
		if c == want {
			return
		}
	}
	t.Fatalf("the call must carry its arguments rendered: want %q in %v", want, eng.calls)
}

// A refused call is the node failing, never an answer: the run stops there.
func TestAToolNodeThatIsRefusedFailsTheRun(t *testing.T) {
	eng := &stubEngine{answers: map[string]string{"tool:git_push": "ERR refused: the wall is shut"}}
	s := toolSpec([]Node{
		{Name: "send", Kind: "tool", Tool: "git_push", Args: map[string]string{"project": "research"}},
		{Name: "after", Kind: "ask", Question: "Q"},
	}, []Edge{{From: "send", To: "after", When: "always"}})
	res, err := Run(t.TempDir(), eng, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictFail {
		t.Fatalf("a refused call must fail the run: %s", res.Verdict)
	}
	for _, c := range eng.calls {
		if strings.HasPrefix(c, "ask:") {
			t.Fatalf("nothing after a refused call may fire: %v", eng.calls)
		}
	}
}

// The call rides the crossing of the gate before it, as a `run` node's turn
// does: the door reads that hand to judge a tool that writes.
func TestAToolNodeRidesTheNearestGatesHand(t *testing.T) {
	home := t.TempDir()
	log := &handLog{}
	s := toolSpec([]Node{
		{Name: "open", Kind: "gate", Title: "save?", Grants: []string{"git_commit"}},
		{Name: "save", Kind: "tool", Tool: "git_commit", Args: map[string]string{"message": "m"}},
	}, []Edge{{From: "open", To: "save", When: "pass"}})
	res, err := Run(home, handEngine{log: log}, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused {
		t.Fatalf("the gate must pause the run first: %s", res.Verdict)
	}
	res, err = Resume(home, handEngine{log: log}, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictComplete || len(log.turns) != 1 || log.turns[0] != "tool:git_commit@open:git_commit" {
		t.Fatalf("the call must ride the gate's crossing: %s %v", res.Verdict, log.turns)
	}
}

// The flow law judges the node's shape at the save; which tools exist, and
// whether one that writes was granted, is the door's to say when it runs.
func TestTheToolLawIsJudgedAtTheSave(t *testing.T) {
	ask := Node{Name: "a", Kind: "ask", Question: "Q"}
	for _, c := range []struct {
		why  string
		node Node
		want string
	}{
		{"a tool node that names no tool", Node{Name: "t", Kind: "tool"}, "names no tool the door could carry"},
		{"a tool named outside the door's shape", Node{Name: "t", Kind: "tool", Tool: "Git Push"}, "names no tool the door could carry"},
		{"a tool on a node that is not a tool node", Node{Name: "t", Kind: "ask", Question: "Q", Tool: "git_push"}, "calls no tool, so `tool` and"},
		{"arguments on a node that is not a tool node", Node{Name: "t", Kind: "ask", Question: "Q", Args: map[string]string{"x": "y"}}, "calls no tool, so `tool` and"},
		{"an argument named as the door's own keys are", Node{Name: "t", Kind: "tool", Tool: "git_push", Args: map[string]string{"__caller": "me"}}, "an argument is named as a tool is"},
	} {
		s := toolSpec([]Node{ask, c.node}, []Edge{{From: "a", To: "t", When: "always"}})
		if _, err := Validate(s); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: want a refusal naming %q, got %v", c.why, c.want, err)
		}
	}
	ok := toolSpec([]Node{ask, {Name: "t", Kind: "tool", Tool: "git_push", Args: map[string]string{"project": "{{world}}"}}},
		[]Edge{{From: "a", To: "t", When: "always"}})
	if _, err := Validate(ok); err != nil {
		t.Fatalf("a lawful tool node is refused: %v", err)
	}
}

func TestTheBareEngineRefusesAToolNode(t *testing.T) {
	_, err := Production(t.TempDir()).Tool(context.Background(), "git_push", nil)
	if err == nil || !strings.Contains(err.Error(), "no door wired") {
		t.Fatalf("the bare engine must refuse a tool node by name: %v", err)
	}
}
