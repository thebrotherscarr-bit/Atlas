package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"atlas/line/internal/play"
)

// stubEngine answers from a canned map; every call is recorded.
type stubEngine struct {
	answers map[string]string
	calls   []string
}

func (s *stubEngine) Ask(_ context.Context, question, _ string) (string, error) {
	s.calls = append(s.calls, "ask:"+question)
	if a, ok := s.answers[question]; ok {
		return a, nil
	}
	return "stub-answer", nil
}

func (s *stubEngine) RunPrompt(name string, _ int, vars map[string]string, _ string) (play.Run, error) {
	s.calls = append(s.calls, "prompt:"+name)
	return play.Run{Output: "stub-prompt:" + name}, nil
}

func (s *stubEngine) SeatAsk(seat, question, _, _ string) (play.Run, error) {
	s.calls = append(s.calls, "seat:"+seat)
	return play.Run{Output: "stub-seat:" + seat + ":" + question}, nil
}

func (s *stubEngine) Turn(_ context.Context, objective, feed, method string) (string, error) {
	s.calls = append(s.calls, "run:"+objective)
	// A canned answer lets a stroke hand back a turn WITH a verdict block, so
	// the evidence rule can be struck both ways. With none, the old string is
	// returned unchanged and every existing stroke reads as it did.
	if a, ok := s.answers[objective]; ok {
		return a, nil
	}
	return "stub-run:" + objective, nil
}

func (s *stubEngine) Aider(_ context.Context, instruction, files, _ string) (string, error) {
	s.calls = append(s.calls, "aider:"+instruction+"|"+files)
	return "stub-aider:" + instruction, nil
}

func (s *stubEngine) Recall(_, question string) (string, error) {
	s.calls = append(s.calls, "memory:"+question)
	return "stub-memory", nil
}

func branchSpec() Spec {
	return Spec{Name: "branch", BudgetS: 600,
		Nodes: []Node{
			{Name: "a", Kind: "ask", Question: "Q"},
			{Name: "e", Kind: "eval", Ref: "a", Expected: "yes"},
			{Name: "b", Kind: "ask", Question: "B {{out_a}}"},
			{Name: "c", Kind: "ask", Question: "C"},
		},
		Edges: []Edge{
			{From: "a", To: "e", When: "always"},
			{From: "e", To: "b", When: "pass"},
			{From: "e", To: "c", When: "fail"},
		}}
}

// --- retry: it answers an ERROR, and never a verdict ------------------------
//
// Taken from the sovereign-microkernel read on 2026-09-12, which was the one
// thing in that WorkflowEngine this engine genuinely lacked. Everything else
// it offered, this ground already had and had proved.

// flakyEngine fails its first `failures` turns outright -- no answer at all,
// which is what `rerr` means -- and answers normally after that.
type flakyEngine struct {
	stubEngine
	failures int
	turns    int
}

func (f *flakyEngine) Turn(ctx context.Context, objective, feed, method string) (string, error) {
	f.turns++
	if f.turns <= f.failures {
		return "", fmt.Errorf("the door is not open")
	}
	return f.stubEngine.Turn(ctx, objective, feed, method)
}

func retrySpec(retries int) Spec {
	return Spec{Name: "flaky", BudgetS: 600, Nodes: []Node{
		{Name: "work", Kind: "run", Question: "do the thing", Retries: retries},
	}}
}

// Two dead turns, two retries allowed, and the run lands -- with BOTH failures
// on the record, because a retry nobody can see is flakiness being hidden.
func TestANodeThatErrorsIsRetriedAndTheRunSurvives(t *testing.T) {
	home := t.TempDir()
	eng := &flakyEngine{failures: 2}
	res, err := Run(home, eng, retrySpec(2), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s, want COMPLETE", res.Verdict)
	}
	if eng.turns != 3 {
		t.Fatalf("engine was called %d times, want 3 (one try, two retries)", eng.turns)
	}

	lines, err := runLog(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	retries, okLines := 0, 0
	for _, l := range lines {
		if l["kind"] != "node" {
			continue
		}
		switch l["status"] {
		case "retry":
			retries++
			if e, _ := l["error"].(string); !strings.Contains(e, "door is not open") {
				t.Fatalf("a retry line must carry the reason, got %q", e)
			}
			if a, ok := l["attempt"].(float64); !ok || a < 1 {
				t.Fatalf("a retry line must number the attempt, got %v", l["attempt"])
			}
			// THE PAUSE IS COUNTED, so the budget sees what retry spends.
			if ms, ok := l["latency_ms"].(float64); !ok || ms < 900 {
				t.Fatalf("a retry line must carry the attempt AND its pause, got %v ms", l["latency_ms"])
			}
		case "ok":
			okLines++
		}
	}
	if retries != 2 {
		t.Fatalf("%d retries written down, want 2", retries)
	}
	if okLines != 1 {
		t.Fatalf("%d ok lines, want exactly 1 -- the attempt that worked", okLines)
	}
}

// The other way: the count is a ceiling, not a promise.
func TestRetryStopsAtTheCountAndTheRunStillFails(t *testing.T) {
	home := t.TempDir()
	eng := &flakyEngine{failures: 99}
	res, err := Run(home, eng, retrySpec(1), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictFail {
		t.Fatalf("verdict = %s, want FAIL", res.Verdict)
	}
	if eng.turns != 2 {
		t.Fatalf("engine was called %d times, want 2 (one try, one retry)", eng.turns)
	}
}

// AND WITH NO RETRIES DECLARED, NOTHING CHANGED. Every flow folded before this
// existed behaves exactly as it did: one call, one failure, one FAIL.
func TestWithoutRetriesTheEngineIsCalledExactlyOnce(t *testing.T) {
	home := t.TempDir()
	eng := &flakyEngine{failures: 99}
	res, err := Run(home, eng, retrySpec(0), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictFail {
		t.Fatalf("verdict = %s, want FAIL", res.Verdict)
	}
	if eng.turns != 1 {
		t.Fatalf("engine was called %d times, want 1", eng.turns)
	}
	lines, err := runLog(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lines {
		if l["status"] == "retry" {
			t.Fatal("a flow that asked for no retries must write no retry line")
		}
	}
}

// THE REFUSAL THAT MATTERS MOST. Retrying an eval until it agrees is the
// laundering path by another name, so the shape is not offered at all.
func TestAVerdictIsNeverRetried(t *testing.T) {
	for _, kind := range []string{"eval", "gate"} {
		s := Spec{Name: "argue", Nodes: []Node{
			{Name: "a", Kind: "ask", Question: "Q"},
			{Name: "v", Kind: kind, Ref: "a", Expected: "yes", Title: "decide", Retries: 2},
		}, Edges: []Edge{{From: "a", To: "v"}}}
		_, err := Validate(s)
		if err == nil {
			t.Fatalf("a %s with retries must be refused", kind)
		}
		if !strings.Contains(err.Error(), "never a verdict") {
			t.Fatalf("the refusal must say why, got: %v", err)
		}
	}
	// Both ways: the kinds that DO call out are allowed to ask.
	for _, kind := range []string{"run", "ask", "seat"} {
		n := Node{Name: "n", Kind: kind, Question: "Q", Seat: "s", Retries: 2}
		if _, err := Validate(Spec{Name: "fine", Nodes: []Node{n}}); err != nil {
			t.Fatalf("a %s may ask for retries: %v", kind, err)
		}
	}
}

func TestRetriesOutsideTheRangeAreRefused(t *testing.T) {
	for _, n := range []int{-1, MaxRetries + 1} {
		s := Spec{Name: "greedy", Nodes: []Node{
			{Name: "a", Kind: "run", Question: "Q", Retries: n}}}
		_, err := Validate(s)
		if err == nil {
			t.Fatalf("retries=%d must be refused", n)
		}
		if !strings.Contains(err.Error(), "the range is 0 to") {
			t.Fatalf("the refusal must name the range, got: %v", err)
		}
	}
	for _, n := range []int{0, MaxRetries} {
		s := Spec{Name: "fine", Nodes: []Node{
			{Name: "a", Kind: "run", Question: "Q", Retries: n}}}
		if _, err := Validate(s); err != nil {
			t.Fatalf("retries=%d must be allowed: %v", n, err)
		}
	}
}

// Bounded, and it says so in numbers rather than in a comment.
func TestTheBackoffIsBoundedAndClimbs(t *testing.T) {
	if retryWait(1) >= retryWait(2) || retryWait(2) >= retryWait(3) {
		t.Fatal("the pause must grow with the attempt")
	}
	for _, a := range []int{5, 50, 5000} {
		if retryWait(a) > 8*time.Second {
			t.Fatalf("retryWait(%d) = %v, past the cap", a, retryWait(a))
		}
	}
	if retryWait(0) != time.Second {
		t.Fatalf("a nonsense attempt must still return a sane pause, got %v", retryWait(0))
	}
}

func TestBranchPass(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{answers: map[string]string{"Q": "yes", "B yes": "b-out"}}
	res, err := Run(home, eng, branchSpec(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s", res.Verdict)
	}
	if strings.Join(res.Fired, ",") != "a,b,e" {
		t.Fatalf("fired = %v (c must stay silent)", res.Fired)
	}
}

func TestBranchFail(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{answers: map[string]string{"Q": "no", "C": "c-out"}}
	res, err := Run(home, eng, branchSpec(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s", res.Verdict)
	}
	if strings.Join(res.Fired, ",") != "a,c,e" {
		t.Fatalf("fired = %v (b must stay silent)", res.Fired)
	}
}

func TestEvalFailNoBranch(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{answers: map[string]string{"Q": "no"}}
	s := Spec{Name: "strict", Nodes: []Node{
		{Name: "a", Kind: "ask", Question: "Q"},
		{Name: "e", Kind: "eval", Ref: "a", Expected: "yes"},
	}, Edges: []Edge{{From: "a", To: "e"}}}
	res, err := Run(home, eng, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictFail {
		t.Fatalf("verdict = %s, want FAIL", res.Verdict)
	}
}

func gateSpec() Spec {
	return Spec{Name: "gated", BudgetS: 600,
		Nodes: []Node{
			{Name: "a", Kind: "ask", Question: "Q"},
			{Name: "g", Kind: "gate", Title: "human review"},
			{Name: "b", Kind: "ask", Question: "B"},
		},
		Edges: []Edge{
			{From: "a", To: "g", When: "always"},
			{From: "g", To: "b", When: "pass"},
			{From: "g", To: "c", When: "fail"},
		}}
}

func TestGatePausesAndResumes(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{answers: map[string]string{}}
	// gate spec needs its fail target to exist for validation
	s := gateSpec()
	s.Nodes = append(s.Nodes, Node{Name: "c", Kind: "ask", Question: "C"})
	res, err := Run(home, eng, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused {
		t.Fatalf("verdict = %s, want PAUSED", res.Verdict)
	}
	st, err := Status(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(st, "PAUSED") || !strings.Contains(st, "budget") {
		t.Fatalf("waterfall must name PAUSED + budget:\n%s", st)
	}
	res2, err := Resume(home, eng, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res2.Verdict != VerdictComplete {
		t.Fatalf("after continue verdict = %s", res2.Verdict)
	}
	if strings.Join(res2.Fired, ",") != "a,b,g" {
		t.Fatalf("fired after resume = %v", res2.Fired)
	}
}

func TestGateStop(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{}
	s := Spec{Name: "stopper", Nodes: []Node{
		{Name: "a", Kind: "ask", Question: "Q"},
		{Name: "g", Kind: "gate", Title: "t"},
	}, Edges: []Edge{{From: "a", To: "g"}}}
	res, err := Run(home, eng, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	res2, err := Resume(home, eng, res.Run, "stop")
	if err != nil {
		t.Fatal(err)
	}
	if res2.Verdict != VerdictStopped {
		t.Fatalf("verdict = %s, want STOPPED", res2.Verdict)
	}
	if _, err := Resume(home, eng, res.Run, "continue"); err == nil {
		t.Fatal("a stopped run must not resume")
	}
}

func TestMissingVarFailsNode(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{}
	s := Spec{Name: "varless", Nodes: []Node{
		{Name: "a", Kind: "ask", Question: "Hi {{who}}"},
	}, Edges: []Edge{}}
	res, err := Run(home, eng, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictFail {
		t.Fatalf("verdict = %s, want FAIL on missing var", res.Verdict)
	}
	if len(eng.calls) != 0 {
		t.Fatal("a failed render must never reach a voice")
	}
}

func TestDownstreamTemplating(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{answers: map[string]string{"Q": "yes", "B yes": "b-out"}}
	if _, err := Run(home, eng, branchSpec(), nil); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, c := range eng.calls {
		if c == "ask:B yes" {
			found = true
		}
	}
	if !found {
		t.Fatalf("downstream must render {{out_a}}: %v", eng.calls)
	}
}

func TestCompareAndReplay(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{answers: map[string]string{"Q": "same"}}
	s := Spec{Name: "twice", Nodes: []Node{{Name: "a", Kind: "ask", Question: "Q"}},
		Edges: []Edge{}}
	r1, err := Run(home, eng, s, map[string]string{"in": "1"})
	if err != nil {
		t.Fatal(err)
	}
	eng2 := &stubEngine{answers: map[string]string{"Q": "different"}}
	r2, err := Run(home, eng2, s, map[string]string{"in": "1"})
	if err != nil {
		t.Fatal(err)
	}
	out, err := Compare(home, r1.Run, r2.Run)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "DIFFER") {
		t.Fatalf("compare must name the difference:\n%s", out)
	}
	r3, err := Replay(home, eng, r1.Run)
	if err != nil {
		t.Fatal(err)
	}
	if r3.Run == r1.Run || r3.Verdict != VerdictComplete {
		t.Fatalf("replay must fire fresh: %+v", r3)
	}
	if _, err := Status(home, r3.Run); err != nil {
		t.Fatal(err)
	}
	if _, err := Compare(home, r1.Run, "f-20260909-120000-deadbeef"); err == nil {
		t.Fatal("compare with a ghost run must refuse")
	}
}

func TestRunIsolation(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{}
	s := Spec{Name: "iso", Nodes: []Node{{Name: "a", Kind: "ask", Question: "Q"}},
		Edges: []Edge{}}
	r1, _ := Run(home, eng, s, nil)
	r2, _ := Run(home, eng, s, nil)
	if r1.Run == r2.Run {
		t.Fatal("run ids must be unique")
	}
	o1, _ := nodeOutputs(home, r1.Run)
	o2, _ := nodeOutputs(home, r2.Run)
	if len(o1) != 1 || len(o2) != 1 {
		t.Fatal("runs must not share outputs")
	}
}

// TestNoFinishPath is the kept static self-check for this package: no
// verb here finishes, lands, closes, merges or pushes anything — gates
// pause, evals steer, the hand moves. Concat-split literals are joined
// before matching so the check cannot be dodged by string-splitting.
func TestNoFinishPath(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		flat := strings.ReplaceAll(string(src), "\"+\"", "")
		flat = strings.ReplaceAll(flat, "\" + \"", "")
		lowered := strings.ToLower(flat)
		for _, bad := range []string{"approve(", "ascend(", "resolve(",
			"promote(", "reject(", "\"done\"", "approved"} {
			if strings.Contains(lowered, bad) {
				t.Fatalf("%s: forbidden path %q", f, bad)
			}
		}
	}
}

// THE PASS BRANCH THAT HAD NEVER BEEN REACHABLE (2026-09-12).
//
// `play.Score` is exact match after trim and casefold, and it also scores
// prompt-eval datasets -- so it could not be loosened without rescoring saved
// runs. Meanwhile the builder's label for this field read "what the answer
// should carry", which is `contains` in words. The coder flow believed the
// label: it checked a `run` node, whose answer is the council's prose, for
// "RAN". `verify` came back "RAN: fizz_buzz.py" over correct FizzBuzz and the
// check failed anyway, every time, since the flow was first folded.
func TestEvalMatchModes(t *testing.T) {
	answer := "RAN: fizz_buzz.py\n--- stdout ---\n1\n2\nFizz"

	// contains: the test the label always described
	if !scoreNode(Node{Kind: "eval", Expected: "RAN", Match: "contains"}, "RAN", answer) {
		t.Fatal("contains did not find RAN in an answer that carries it")
	}
	if !scoreNode(Node{Kind: "eval", Expected: "RAN", Match: "CONTAINS"}, "RAN", answer) {
		t.Fatal("the MODE is case-blind even though the needle is not")
	}
	if scoreNode(Node{Kind: "eval", Expected: "FAILED", Match: "contains"}, "FAILED", answer) {
		t.Fatal("contains passed on a string the answer does not carry")
	}

	// THE PASS THAT WAS A LIE, measured 2026-09-12 on the first live run after
	// `contains` landed. calculate_sum.py died of a SyntaxError, run_python
	// said so, and the check passed anyway -- because the delivery contained
	// the English word "ran". A word-boundary test would not have caught it:
	// that match WAS a whole word. Case is what separates a verdict from prose.
	failed := `FAILED (exit 1): calculate_sum.py
--- stderr ---
SyntaxError: invalid syntax

The tools that actually ran this turn were write_file, run_python.`
	if scoreNode(Node{Kind: "eval", Expected: "RAN", Match: "contains"}, "RAN", failed) {
		t.Fatal("a FAILED run scored as a pass; a gate must not tell that lie")
	}
	if !scoreNode(Node{Kind: "eval", Expected: "FAILED", Match: "contains"}, "FAILED", failed) {
		t.Fatal("contains could not find the verdict that is actually there")
	}
	// and the verdict token run_python really emits, on a real success
	if !scoreNode(Node{Kind: "eval", Expected: "RAN:", Match: "contains"}, "RAN:", answer) {
		t.Fatal("contains missed the RAN: verdict line run_python emits")
	}
	if scoreNode(Node{Kind: "eval", Expected: "RAN:", Match: "contains"}, "RAN:", failed) {
		t.Fatal("RAN: matched a run that failed")
	}

	// equals: UNCHANGED, which is the whole reason the mode exists
	if scoreNode(Node{Kind: "eval", Expected: "RAN"}, "RAN", answer) {
		t.Fatal("an unmarked eval stopped being exact match; saved specs moved")
	}
	if !scoreNode(Node{Kind: "eval"}, " ran ", "RAN") {
		t.Fatal("equals must still trim and casefold, as play.Score does")
	}

	// AN EMPTY EXPECTED NEVER PASSES -- every string contains ""
	for _, m := range []string{"", "equals", "contains"} {
		if scoreNode(Node{Kind: "eval", Expected: "  ", Match: m}, "  ", answer) {
			t.Fatalf("match %q passed on a blank expected; that is a green light nobody set", m)
		}
	}
}

// --- the outcome is written down, and read back (2026-09-22) ----------------
//
// The handoff's second piece, his word: "2. Flows report truthfully". A node
// line carried `status`, which is "ok" for an eval that ANSWERED -- pass or
// fail alike -- so nothing on disk said which way a check went, and Resume
// rebuilt every ok line as a pass. A check that FAILED came back from a gate
// as one that passed: the fail-branch went quiet and the pass-branch fired.

// evalGateSpec: a check fires, a gate pauses, and two nodes downstream wait on
// the CHECK's own outcome -- the shape a resume used to lie about.
func evalGateSpec() Spec {
	return Spec{Name: "judged", BudgetS: 600,
		Nodes: []Node{
			{Name: "a", Kind: "ask", Question: "Q"},
			{Name: "e", Kind: "eval", Ref: "a", Expected: "yes"},
			{Name: "g", Kind: "gate", Title: "look at it"},
			{Name: "x", Kind: "ask", Question: "X"},
			{Name: "y", Kind: "ask", Question: "Y"},
		},
		Edges: []Edge{
			{From: "a", To: "e", When: "always"},
			{From: "a", To: "g", When: "always"},
			{From: "e", To: "x", When: "pass"},
			{From: "e", To: "y", When: "fail"},
		}}
}

func TestAResumeKeepsTheCheckThatFailed(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{answers: map[string]string{"Q": "no"}}
	res, err := Run(home, eng, evalGateSpec(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused {
		t.Fatalf("verdict = %s, want PAUSED", res.Verdict)
	}

	// The check's own verdict is ON the line, not left to be guessed later.
	lines, err := runLog(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	seen := false
	for _, l := range lines {
		if l["kind"] != "node" || l["node"] != "e" {
			continue
		}
		seen = true
		p, ok := l["pass"].(bool)
		if !ok {
			t.Fatalf("the check's line does not say which way it went: %v", l)
		}
		if p {
			t.Fatalf("a check that failed was written down as a pass: %v", l)
		}
	}
	if !seen {
		t.Fatal("no line at all for the check")
	}

	res2, err := Resume(home, eng, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(res2.Fired, ",") != "a,e,g,y" {
		t.Fatalf("fired after the resume = %v; the fail branch must fire and the "+
			"pass branch must stay silent", res2.Fired)
	}
}

// AND THE OTHER WAY, or the stroke above would pass on a runner that called
// every check a failure.
func TestAResumeKeepsTheCheckThatPassed(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{answers: map[string]string{"Q": "yes"}}
	res, err := Run(home, eng, evalGateSpec(), nil)
	if err != nil {
		t.Fatal(err)
	}
	res2, err := Resume(home, eng, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(res2.Fired, ",") != "a,e,g,x" {
		t.Fatalf("fired after the resume = %v; the pass branch must fire", res2.Fired)
	}
}

// A RUN LOGGED BEFORE THE OUTCOME WAS WRITTEN DOWN still resumes correctly: an
// eval's own answer says which way it went, and every other kind passed by
// answering at all. Nothing in the record is rewritten to make this true.
func TestAnOlderRunsCheckIsReadFromItsOwnAnswer(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{answers: map[string]string{"Q": "no"}}
	res, err := Run(home, eng, evalGateSpec(), nil)
	if err != nil {
		t.Fatal(err)
	}
	// Strip the field a run written today carries, leaving yesterday's shape.
	raw, err := os.ReadFile(runsPath(home))
	if err != nil {
		t.Fatal(err)
	}
	var kept []string
	for _, ln := range strings.Split(strings.TrimRight(string(raw), "\n"), "\n") {
		var doc map[string]any
		if json.Unmarshal([]byte(ln), &doc) == nil {
			delete(doc, "pass")
			b, _ := json.Marshal(doc)
			ln = string(b)
		}
		kept = append(kept, ln)
	}
	if err := os.WriteFile(runsPath(home), []byte(strings.Join(kept, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res2, err := Resume(home, eng, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(res2.Fired, ",") != "a,e,g,y" {
		t.Fatalf("fired = %v; an older run's failed check was read as a pass", res2.Fired)
	}
}

// --- every gate, not just the first (2026-09-22) -----------------------------
//
// Any `resumed` line used to close the door on the whole run, so a flow with
// two gates could never pass the second -- while Status and flow_runs went on
// showing PAUSED and the tool went on telling the operator to resume it.
func twoGateSpec() Spec {
	return Spec{Name: "twogates", BudgetS: 600,
		Nodes: []Node{
			{Name: "a", Kind: "ask", Question: "Q"},
			{Name: "g1", Kind: "gate", Title: "first look"},
			{Name: "b", Kind: "ask", Question: "B"},
			{Name: "g2", Kind: "gate", Title: "second look"},
			{Name: "c", Kind: "ask", Question: "C"},
		},
		Edges: []Edge{
			{From: "a", To: "g1", When: "always"},
			{From: "g1", To: "b", When: "pass"},
			{From: "b", To: "g2", When: "always"},
			{From: "g2", To: "c", When: "pass"},
		}}
}

func TestEveryGateCanBeResumed(t *testing.T) {
	home := t.TempDir()
	eng := &stubEngine{}
	res, err := Run(home, eng, twoGateSpec(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused || res.PausedNode != "g1" {
		t.Fatalf("first pause = %s at %q", res.Verdict, res.PausedNode)
	}
	res2, err := Resume(home, eng, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res2.Verdict != VerdictPaused || res2.PausedNode != "g2" {
		t.Fatalf("second pause = %s at %q", res2.Verdict, res2.PausedNode)
	}
	res3, err := Resume(home, eng, res.Run, "continue")
	if err != nil {
		t.Fatalf("the second gate refused to resume: %v", err)
	}
	if res3.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s after both gates were answered", res3.Verdict)
	}
	if _, err := Resume(home, eng, res.Run, "continue"); err == nil {
		t.Fatal("a run that has finished resumed again")
	}
}

// --- a cancel is not a clock (2026-09-22) ------------------------------------

// cancellingEngine ends the run from inside its own first node, the way
// flow_cancel ends it from outside.
type cancellingEngine struct{ stubEngine }

func (c *cancellingEngine) Ask(ctx context.Context, question, voice string) (string, error) {
	flightsMu.Lock()
	for _, stop := range flights {
		stop()
	}
	flightsMu.Unlock()
	<-ctx.Done()
	return "", ctx.Err()
}

func TestACancelledRunIsStoppedAndNotOutOfTime(t *testing.T) {
	if verdictFor(context.Canceled) != VerdictStopped {
		t.Fatal("a cancelled context must read as STOPPED")
	}
	if verdictFor(context.DeadlineExceeded) != VerdictOverTime {
		t.Fatal("a spent budget must still read as OUT_OF_TIME")
	}
	home := t.TempDir()
	res, err := Run(home, &cancellingEngine{}, Spec{Name: "cancelme", BudgetS: 600,
		Nodes: []Node{{Name: "a", Kind: "ask", Question: "Q"}}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictStopped {
		t.Fatalf("verdict = %s; a run the hand cancelled did not run out of time", res.Verdict)
	}
}

// --- a resume with no engine burns nothing (2026-09-22) ----------------------

// notReadyEngine answers everything, and says no engine is standing -- the
// shape of one that idled out while a gate waited for the operator.
type notReadyEngine struct{ stubEngine }

func (n *notReadyEngine) Ready() error {
	return fmt.Errorf("no engine is open on this world (env_open first)")
}

func TestAResumeWithNoEngineRefusesAndLeavesTheRunAtItsGate(t *testing.T) {
	home := t.TempDir()
	s := Spec{Name: "guarded", BudgetS: 600,
		Nodes: []Node{
			{Name: "g", Kind: "gate", Title: "look before it runs"},
			{Name: "w", Kind: "run", Question: "do the thing"},
		},
		Edges: []Edge{{From: "g", To: "w", When: "pass"}}}
	res, err := Run(home, &stubEngine{}, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused {
		t.Fatalf("verdict = %s, want PAUSED", res.Verdict)
	}
	before, err := runLog(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Resume(home, &notReadyEngine{}, res.Run, "continue")
	if err == nil {
		t.Fatal("a resume with no engine was allowed to burn the run")
	}
	for _, want := range []string{"no engine is open", "stands at its gate"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("the refusal must say %q: %v", want, err)
		}
	}
	after, err := runLog(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatalf("the refused resume wrote %d line(s) into the run", len(after)-len(before))
	}
	// And with an engine standing, the same gate resumes and the run lands.
	res2, err := Resume(home, &stubEngine{}, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res2.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s once an engine was open", res2.Verdict)
	}
}

func TestValidateGuardsTheMatchMode(t *testing.T) {
	spec := func(n Node) Spec {
		return Spec{Name: "m", Nodes: []Node{{Name: "a", Kind: "ask", Question: "Q"}, n},
			Edges: []Edge{{From: "a", To: "e", When: "always"}}}
	}
	if _, err := Validate(spec(Node{Name: "e", Kind: "eval", Ref: "a", Expected: "x", Match: "roughly"})); err == nil {
		t.Fatal("Validate accepted an unknown match mode")
	} else if !strings.Contains(err.Error(), "equals, contains") {
		t.Fatalf("the refusal must name the set: %v", err)
	}
	// a mode on something with no answer to test is a field that means nothing
	s := Spec{Name: "m", Nodes: []Node{{Name: "a", Kind: "ask", Question: "Q", Match: "contains"}}}
	if _, err := Validate(s); err == nil {
		t.Fatal("Validate accepted `match` on an ask node")
	}
	// and both real modes pass validation
	for _, m := range []string{"", "equals", "contains"} {
		if _, err := Validate(spec(Node{Name: "e", Kind: "eval", Ref: "a", Expected: "x", Match: m})); err != nil {
			t.Fatalf("Validate refused match %q: %v", m, err)
		}
	}
}

// --- the head belongs to the run, not the spec (2026-09-23) ------------------
//
// A model could only be named per NODE before this, so measuring one flow on
// two models meant folding two specs -- and two specs stop being one experiment
// the moment either is edited. These strike the other reading: one spec, the
// head named at the fire, and the record saying which head answered.

// headLog is shared by every copy of a headEngine, so a stroke can read what
// each node was actually measured on after WithVoice has handed back a copy.
type headLog struct {
	asks  []string // "<what>@<voice>" -- a node's own voice, or the run's
	turns []string // "<objective>@<head>" -- the council's, off the engine
}

// headEngine takes a run-level head exactly as THE LINE's council does: by
// value, so WithHead binds a COPY and the engine handed in stays unbound.
type headEngine struct {
	log  *headLog
	head Head
}

func (h headEngine) WithHead(head Head) Engine { h.head = head; return h }

func (h headEngine) Ask(_ context.Context, q, voice string) (string, error) {
	h.log.asks = append(h.log.asks, q+"@"+voice)
	return "answered", nil
}

func (h headEngine) RunPrompt(name string, _ int, _ map[string]string, voice string) (play.Run, error) {
	h.log.asks = append(h.log.asks, "prompt:"+name+"@"+voice)
	return play.Run{Output: "prompted"}, nil
}

func (h headEngine) SeatAsk(seat, _, voice, _ string) (play.Run, error) {
	h.log.asks = append(h.log.asks, "seat:"+seat+"@"+voice)
	return play.Run{Output: "seated"}, nil
}

func (h headEngine) Recall(voice, q string) (string, error) {
	h.log.asks = append(h.log.asks, "memory:"+q+"@"+voice)
	return "recalled", nil
}

func (h headEngine) Aider(_ context.Context, instruction, _, _ string) (string, error) {
	h.log.turns = append(h.log.turns, "aider:"+instruction+"@"+headName(h.head))
	return "aided " + instruction, nil
}

// The council's head is written down whole: the roster-wide voice, then the
// seats named one by one, which is the only record of what a `run` node was
// actually fired on.
func (h headEngine) Turn(_ context.Context, objective, _, _ string) (string, error) {
	h.log.turns = append(h.log.turns, objective+"@"+headName(h.head))
	return "ran " + objective, nil
}

// One node that names no voice, one that pins its own, a seat, and a `run`
// node whose head can only come off the engine.
func headSpec() Spec {
	return Spec{Name: "heads", BudgetS: 600,
		Nodes: []Node{
			{Name: "free", Kind: "ask", Question: "Q1"},
			{Name: "pinned", Kind: "ask", Question: "Q2", Voice: "pinned:latest"},
			{Name: "seated", Kind: "seat", Seat: "steward", Question: "Q3"},
			{Name: "council", Kind: "run", Question: "do the thing"},
		},
		Edges: []Edge{
			{From: "free", To: "pinned", When: "always"},
			{From: "pinned", To: "seated", When: "always"},
			{From: "seated", To: "council", When: "always"},
		}}
}

func startLineOf(t *testing.T, home, run string) map[string]any {
	t.Helper()
	lines, err := runLog(home, run)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lines {
		if l["kind"] == "start" {
			return l
		}
	}
	t.Fatalf("run %s carries no start line", run)
	return nil
}

func hasCall(calls []string, want string) bool {
	for _, c := range calls {
		if c == want {
			return true
		}
	}
	return false
}

// The head reaches every node that names none, leaves alone the one that does,
// reaches the council through the engine, and is written down.
func TestTheHeadBelongsToTheRun(t *testing.T) {
	home := t.TempDir()
	log := &headLog{}
	res, err := RunOn(home, headEngine{log: log}, headSpec(), nil,
		Head{Voice: "head-a:latest"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s", res.Verdict)
	}
	if !hasCall(log.asks, "Q1@head-a:latest") {
		t.Fatalf("a node naming no voice must answer on the run's head: %v", log.asks)
	}
	if !hasCall(log.asks, "seat:steward@head-a:latest") {
		t.Fatalf("a seat node naming no voice must answer on the run's head: %v", log.asks)
	}
	// A PINNED NODE IS PINNED. The spec said this one aloud; a run-level head
	// overriding it would make the spec on disk a lie about what fires.
	if !hasCall(log.asks, "Q2@pinned:latest") {
		t.Fatalf("a node that names its own voice must keep it: %v", log.asks)
	}
	// The council's head cannot ride on a node -- a `run` node is the whole
	// estate, not one voice -- so it comes off the engine or nowhere.
	if !hasCall(log.turns, "do the thing@head-a:latest") {
		t.Fatalf("the council must be fired on the run's head: %v", log.turns)
	}
	if got := startHead(startLineOf(t, home, res.Run)).Voice; got != "head-a:latest" {
		t.Fatalf("the start line must carry the head, got %q", got)
	}
	st, err := Status(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(st, "head: head-a:latest") {
		t.Fatalf("the waterfall must name the head:\n%s", st)
	}
}

// And with no head named, nothing moves: no field on the line, no word on the
// waterfall, every node on the ground's declared targets as before.
func TestARunWithNoHeadNamesNoneAtAll(t *testing.T) {
	home := t.TempDir()
	log := &headLog{}
	res, err := Run(home, headEngine{log: log}, headSpec(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !hasCall(log.asks, "Q1@") {
		t.Fatalf("an unheaded run must leave the voice empty: %v", log.asks)
	}
	if !hasCall(log.turns, "do the thing@the declared targets") {
		t.Fatalf("an unheaded run must leave the council on the ground's own "+
			"targets: %v", log.turns)
	}
	if _, ok := startLineOf(t, home, res.Run)["voice"]; ok {
		t.Fatal("an unheaded run wrote a voice field into its record")
	}
	st, err := Status(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(st, "head:") {
		t.Fatalf("the waterfall named a head nobody set:\n%s", st)
	}
}

func headGateSpec() Spec {
	return Spec{Name: "headgate", BudgetS: 600,
		Nodes: []Node{
			{Name: "before", Kind: "ask", Question: "Q1"},
			{Name: "g", Kind: "gate", Title: "look at it"},
			{Name: "after", Kind: "ask", Question: "Q2"},
		},
		Edges: []Edge{
			{From: "before", To: "g", When: "always"},
			{From: "g", To: "after", When: "pass"},
		}}
}

// A gate is FOR walking away. Coming back must not finish the run somewhere
// else: the head is read off the run's own start line, not off the engine the
// hand happens to be holding when it answers.
func TestAResumedRunFinishesOnTheHeadItBeganOn(t *testing.T) {
	home := t.TempDir()
	res, err := RunOn(home, headEngine{log: &headLog{}}, headGateSpec(), nil,
		Head{Voice: "head-a:latest"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused || res.PausedNode != "g" {
		t.Fatalf("the gate did not pause: %s at %q", res.Verdict, res.PausedNode)
	}
	// A FRESH, UNBOUND engine -- the shape of coming back hours later.
	after := &headLog{}
	res2, err := Resume(home, headEngine{log: after}, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res2.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s", res2.Verdict)
	}
	if !hasCall(after.asks, "Q2@head-a:latest") {
		t.Fatalf("the second half ran on a different head: %v", after.asks)
	}
}

// Same spec, same inputs, same model. A replay that dropped the head would
// answer a different question and still report COMPLETE.
func TestAReplayRefiresTheHeadAndStampsIt(t *testing.T) {
	home := t.TempDir()
	res, err := RunOn(home, headEngine{log: &headLog{}}, headSpec(), nil,
		Head{Voice: "head-a:latest"})
	if err != nil {
		t.Fatal(err)
	}
	again := &headLog{}
	rep, err := Replay(home, headEngine{log: again}, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Run == res.Run {
		t.Fatal("a replay must take a fresh id")
	}
	if !hasCall(again.asks, "Q1@head-a:latest") || !hasCall(again.turns, "do the thing@head-a:latest") {
		t.Fatalf("the replay ran on a different head: %v / %v", again.asks, again.turns)
	}
	if got := startHead(startLineOf(t, home, rep.Run)).Voice; got != "head-a:latest" {
		t.Fatalf("the replay's own start line must carry the head, got %q", got)
	}
}

// The parity reader: two columns, and which model produced each.
func TestCompareNamesTheTwoHeadsWhenTheyDiffer(t *testing.T) {
	home := t.TempDir()
	s := Spec{Name: "one", BudgetS: 600,
		Nodes: []Node{{Name: "a", Kind: "ask", Question: "Q"}}, Edges: []Edge{}}
	a, err := RunOn(home, headEngine{log: &headLog{}}, s, nil, Head{Voice: "head-a:latest"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := RunOn(home, headEngine{log: &headLog{}}, s, nil, Head{Voice: "head-b:latest"})
	if err != nil {
		t.Fatal(err)
	}
	out, err := Compare(home, a.Run, b.Run)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "head-a:latest") || !strings.Contains(out, "head-b:latest") {
		t.Fatalf("a compare across two heads must name both:\n%s", out)
	}
	// Two runs on ONE head is the model's own variance, not a comparison, and
	// is not dressed as one.
	c, err := RunOn(home, headEngine{log: &headLog{}}, s, nil, Head{Voice: "head-a:latest"})
	if err != nil {
		t.Fatal(err)
	}
	same, err := Compare(home, a.Run, c.Run)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(same, "A: head-a:latest") {
		t.Fatalf("two runs on one head must not render a head row:\n%s", same)
	}
}

// An engine with no WithVoice is left exactly as it was: the bare prodEngine
// has none, and a stub in another stroke must not start answering differently
// because a head was named.
func TestAnEngineThatTakesNoHeadIsLeftAlone(t *testing.T) {
	eng := &stubEngine{}
	if got := onHead(eng, Head{Voice: "head-a:latest"}); got != Engine(eng) {
		t.Fatal("onHead replaced an engine that cannot take a head")
	}
	bound := headEngine{log: &headLog{}}
	if got := onHead(bound, Head{}); got.(headEngine).head.Named() {
		t.Fatal("an empty head must bind nothing")
	}
	// A HEAD OF NOTHING BUT WHITESPACE IS NO HEAD. tidy drops it, so a caller
	// that passed a blank field cannot produce a run whose record says it was
	// headed -- and the record is the whole of what two runs are compared on.
	if got := onHead(bound, Head{Voice: "  ", Voices: map[string]string{"Steward": " "}}); got.(headEngine).head.Named() {
		t.Fatal("a head of blanks bound something")
	}
}

// --- and a head per seat (2026-09-23, "then B underneath it") ----------------
//
// The run-level voice above answers "is this flow better on that model". It
// cannot answer "does the STEWARD raise a flag where it used to announce",
// because it moves the whole roster in the same breath and the answer becomes
// a fact about two changes at once. `Voices` is that narrower question.

// A per-seat head reaches the COUNCIL and nothing else, because a `run` node is
// the only kind with a roster for it to name. An `ask` node measures against
// one model; there is nothing for a seat map to say to it, and pretending
// otherwise would quietly pick one of its entries.
func TestAPerSeatHeadReachesTheCouncilAndLeavesTheNodesAlone(t *testing.T) {
	home := t.TempDir()
	log := &headLog{}
	head := Head{Voices: map[string]string{"Steward": "phi4-mini:latest"}}
	res, err := RunOn(home, headEngine{log: log}, headSpec(), nil, head)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s", res.Verdict)
	}
	if !hasCall(log.turns, "do the thing@Steward on phi4-mini:latest") {
		t.Fatalf("the council must be fired on the per-seat head: %v", log.turns)
	}
	if !hasCall(log.asks, "Q1@") {
		t.Fatalf("a seat map must not be handed to a node that measures one "+
			"voice: %v", log.asks)
	}
	if !hasCall(log.asks, "Q2@pinned:latest") {
		t.Fatalf("a node that names its own voice must still keep it: %v", log.asks)
	}
	got := startHead(startLineOf(t, home, res.Run))
	if got.Voices["Steward"] != "phi4-mini:latest" || got.Voice != "" {
		t.Fatalf("the start line must carry the seat map and nothing it was not "+
			"given, got %+v", got)
	}
	st, err := Status(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(st, "seats: Steward on phi4-mini:latest") {
		t.Fatalf("the waterfall must name the seats, or a run that moved one is "+
			"indistinguishable from one that moved none:\n%s", st)
	}
}

// BOTH AT ONCE: everything on one head, one seat over it. That is the shape a
// parity of a single voice actually needs, and the two must ride together or
// the record says one thing while the run did another.
func TestAHeadCanNameTheRosterAndOneSeatOverIt(t *testing.T) {
	home := t.TempDir()
	log := &headLog{}
	head := Head{Voice: "head-a:latest",
		Voices: map[string]string{"Steward": "phi4-mini:latest"}}
	res, err := RunOn(home, headEngine{log: log}, headSpec(), nil, head)
	if err != nil {
		t.Fatal(err)
	}
	if !hasCall(log.turns, "do the thing@head-a:latest, then Steward on phi4-mini:latest") {
		t.Fatalf("the council must get both halves, in that order: %v", log.turns)
	}
	// The roster-wide half still defaults the nodes; the seat map still does not.
	if !hasCall(log.asks, "Q1@head-a:latest") {
		t.Fatalf("the roster-wide half must still default a node: %v", log.asks)
	}
	got := startHead(startLineOf(t, home, res.Run))
	if got.Voice != "head-a:latest" || got.Voices["Steward"] != "phi4-mini:latest" {
		t.Fatalf("the start line must carry both, got %+v", got)
	}
	st, err := Status(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(st, "head: head-a:latest") ||
		!strings.Contains(st, "seats: Steward on phi4-mini:latest") {
		t.Fatalf("the waterfall must name both:\n%s", st)
	}
}

// A seat map survives the gate and the replay for the same reason the voice
// does: a parity carried on hours later, or repeated, must be the same run.
func TestAPerSeatHeadSurvivesResumeAndReplay(t *testing.T) {
	home := t.TempDir()
	head := Head{Voices: map[string]string{"Steward": "phi4-mini:latest"}}
	res, err := RunOn(home, headEngine{log: &headLog{}}, headGateSpec(), nil, head)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused {
		t.Fatalf("the gate did not pause: %s", res.Verdict)
	}
	res2, err := Resume(home, headEngine{log: &headLog{}}, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res2.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s", res2.Verdict)
	}
	full, err := RunOn(home, headEngine{log: &headLog{}}, headSpec(), nil, head)
	if err != nil {
		t.Fatal(err)
	}
	again := &headLog{}
	rep, err := Replay(home, headEngine{log: again}, full.Run)
	if err != nil {
		t.Fatal(err)
	}
	if !hasCall(again.turns, "do the thing@Steward on phi4-mini:latest") {
		t.Fatalf("the replay dropped the seat map: %v", again.turns)
	}
	if got := startHead(startLineOf(t, home, rep.Run)); got.Voices["Steward"] != "phi4-mini:latest" {
		t.Fatalf("the replay's own start line must carry it, got %+v", got)
	}
}

// Two runs that differ ONLY in one seat are the parity this exists for, and
// the compare has to say which seat, or the two columns are unlabelled.
func TestCompareNamesTwoHeadsThatDifferByOneSeat(t *testing.T) {
	home := t.TempDir()
	s := Spec{Name: "one", BudgetS: 600,
		Nodes: []Node{{Name: "c", Kind: "run", Question: "do it"}}, Edges: []Edge{}}
	a, err := RunOn(home, headEngine{log: &headLog{}}, s, nil,
		Head{Voices: map[string]string{"Steward": "llama3.2:latest"}})
	if err != nil {
		t.Fatal(err)
	}
	b, err := RunOn(home, headEngine{log: &headLog{}}, s, nil,
		Head{Voices: map[string]string{"Steward": "phi4-mini:latest"}})
	if err != nil {
		t.Fatal(err)
	}
	out, err := Compare(home, a.Run, b.Run)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "A: Steward on llama3.2:latest") ||
		!strings.Contains(out, "B: Steward on phi4-mini:latest") {
		t.Fatalf("a compare across two seat heads must name both:\n%s", out)
	}
	// AND THE SEATS ARE NAMED IN SEAT ORDER, not in Go's map order, or two
	// runs' lines cannot be read against each other at all.
	many := Head{Voices: map[string]string{
		"Router": "r:latest", "Steward": "s:latest", "Delivery Agent": "d:latest"}}
	if got := headName(many); got != "Delivery Agent on d:latest, Router on r:latest, Steward on s:latest" {
		t.Fatalf("seats must render in seat order, got %q", got)
	}
}

// --- a crossed gate carries its grants (2026-09-28) ---------------------------
//
// Under --auth the council's writes parked even after the operator resumed a
// gate that asked exactly that question. A gate may now declare `grants`: what
// its `continue` authorises, from that gate until the next gate or the end of
// the run. These strike the flow's half -- the grants ride to the nodes after
// the gate and no further, are written on the pause and on the resume, and
// belong to gates alone. THE LINE's half, the hand over the door, is in tools.

// handLog is shared by every copy of a handEngine, so a stroke can read what
// hand each turn ran under after WithHand has handed back a copy.
type handLog struct{ turns []string }

// handEngine takes a crossing exactly as THE LINE's council does: by value, so
// WithHand binds a COPY and the engine handed in stays unbound.
type handEngine struct {
	log    *handLog
	run    string
	gate   string
	grants []string
}

func (h handEngine) WithHand(run, gate string, grants []string) Engine {
	h.run, h.gate, h.grants = run, gate, grants
	return h
}

func (h handEngine) Turn(_ context.Context, objective, _, _ string) (string, error) {
	h.log.turns = append(h.log.turns, objective+"@"+h.gate+":"+strings.Join(h.grants, ","))
	return "ran " + objective, nil
}

func (h handEngine) Ask(_ context.Context, _, _ string) (string, error) { return "answered", nil }

func (h handEngine) RunPrompt(_ string, _ int, _ map[string]string, _ string) (play.Run, error) {
	return play.Run{Output: "prompted"}, nil
}

func (h handEngine) SeatAsk(_, _, _, _ string) (play.Run, error) {
	return play.Run{Output: "seated"}, nil
}

func (h handEngine) Recall(_, _ string) (string, error) { return "recalled", nil }

func (h handEngine) Aider(_ context.Context, instruction, _, _ string) (string, error) {
	h.log.turns = append(h.log.turns, "aider:"+instruction+"@"+h.gate+":"+strings.Join(h.grants, ","))
	return "aided " + instruction, nil
}

// One gate that grants, a run node after it, a second gate that grants
// nothing, and a run node after that.
func grantSpec() Spec {
	return Spec{Name: "granted", BudgetS: 600,
		Nodes: []Node{
			{Name: "a", Kind: "ask", Question: "Q"},
			{Name: "g", Kind: "gate", Title: "cut it?", Grants: []string{"git_tag"}},
			{Name: "w", Kind: "run", Question: "cut it"},
			{Name: "g2", Kind: "gate", Title: "send it?"},
			{Name: "w2", Kind: "run", Question: "send it"},
		},
		Edges: []Edge{
			{From: "a", To: "g", When: "always"},
			{From: "g", To: "w", When: "pass"},
			{From: "w", To: "g2", When: "always"},
			{From: "g2", To: "w2", When: "pass"},
		}}
}

func TestACrossedGateCarriesItsGrantsToTheNodesAfterItAndNoFurther(t *testing.T) {
	home := t.TempDir()
	res, err := Run(home, handEngine{log: &handLog{}}, grantSpec(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused || res.PausedNode != "g" {
		t.Fatalf("the first gate did not pause: %s at %q", res.Verdict, res.PausedNode)
	}
	// The pause says what continue will authorise, beside the question.
	lines, err := runLog(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	pausedGrants := ""
	for _, l := range lines {
		if l["kind"] == "node" && l["node"] == "g" && l["status"] == "paused" {
			b, _ := json.Marshal(l["grants"])
			pausedGrants = string(b)
		}
	}
	if pausedGrants != `["git_tag"]` {
		t.Fatalf("the pause must carry the grants, got %s", pausedGrants)
	}
	st, err := Status(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(st, "continue grants the council: git_tag") {
		t.Fatalf("the waterfall must say what continue authorises:\n%s", st)
	}

	// Crossing the first gate: the node after it runs under the grants.
	first := &handLog{}
	res2, err := Resume(home, handEngine{log: first}, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res2.Verdict != VerdictPaused || res2.PausedNode != "g2" {
		t.Fatalf("the second gate did not pause: %s at %q", res2.Verdict, res2.PausedNode)
	}
	if strings.Join(first.turns, "|") != "cut it@g:git_tag" {
		t.Fatalf("the node after the crossed gate must run under its grants: %v", first.turns)
	}
	lines, _ = runLog(home, res.Run)
	resumedGrants := ""
	for _, l := range lines {
		if l["kind"] == "resumed" && l["node"] == "g" {
			b, _ := json.Marshal(l["grants"])
			resumedGrants = string(b)
		}
	}
	if resumedGrants != `["git_tag"]` {
		t.Fatalf("the resume must say what the hand authorised, got %s", resumedGrants)
	}

	// Crossing the second gate, which grants nothing: the hand is gone.
	second := &handLog{}
	res3, err := Resume(home, handEngine{log: second}, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res3.Verdict != VerdictComplete {
		t.Fatalf("verdict = %s after both gates", res3.Verdict)
	}
	if strings.Join(second.turns, "|") != "send it@:" {
		t.Fatalf("a gate that grants nothing must carry no hand past it: %v", second.turns)
	}

	// An engine that takes no hand is left alone; a gate with no grants binds
	// nothing; and a bound copy leaves the engine handed in unbound.
	eng := &stubEngine{}
	if got := onHand(eng, "r", Node{Kind: "gate", Grants: []string{"git_tag"}}); got != Engine(eng) {
		t.Fatal("onHand replaced an engine that cannot take a hand")
	}
	plain := handEngine{log: &handLog{}}
	if got := onHand(plain, "r", Node{Kind: "gate", Title: "t"}); len(got.(handEngine).grants) != 0 {
		t.Fatal("a gate with no grants bound a hand")
	}
	bound := onHand(plain, "r", Node{Kind: "gate", Grants: []string{"git_tag"}})
	if len(plain.grants) != 0 || len(bound.(handEngine).grants) != 1 {
		t.Fatal("WithHand must bind a copy and leave the engine handed in unbound")
	}
}

// ONLY A GATE GRANTS: a grant on any other kind is refused at the save, and an
// empty grant name is refused with it.
func TestOnlyAGateGrants(t *testing.T) {
	for _, n := range []Node{
		{Name: "n", Kind: "ask", Question: "Q", Grants: []string{"git_tag"}},
		{Name: "n", Kind: "run", Question: "Q", Grants: []string{"git_tag"}},
		{Name: "n", Kind: "eval", Ref: "a", Expected: "x", Grants: []string{"git_tag"}},
	} {
		s := Spec{Name: "g", Nodes: []Node{{Name: "a", Kind: "ask", Question: "Q"}, n},
			Edges: []Edge{{From: "a", To: "n"}}}
		_, err := Validate(s)
		if err == nil || !strings.Contains(err.Error(), "only a gate") {
			t.Fatalf("a %s with grants must be refused by name: %v", n.Kind, err)
		}
	}
	empty := Spec{Name: "g", Nodes: []Node{{Name: "g", Kind: "gate", Title: "t", Grants: []string{" "}}}}
	if _, err := Validate(empty); err == nil || !strings.Contains(err.Error(), "empty name") {
		t.Fatalf("an empty grant must be refused: %v", err)
	}
	for _, grants := range [][]string{nil, {"git_tag"}, {"git_tag", "git_branch"}} {
		s := Spec{Name: "g", Nodes: []Node{{Name: "g", Kind: "gate", Title: "t", Grants: grants}}}
		if _, err := Validate(s); err != nil {
			t.Fatalf("a gate granting %v must validate: %v", grants, err)
		}
	}
}

// --- a bounded return (2026-09-28, LAW_003's mechanism) -----------------------
//
// The operator: "create the bounded back-edge looping". A check whose fail-edge
// points back at a node that declares `loops` sends the run back there, the
// work between is fired again, at most that many times, with every pass and
// every return on the record. The ceiling is on the node returned to (his
// card), it is read by the loop and nowhere else, and the check is the stop.

// countingEngine answers a run node by how many times it has been asked, so a
// check can fail the first passes and pass a later one: the work re-done, the
// score never re-rolled.
type countingEngine struct {
	stubEngine
	turns  int
	passAt int
	asked  []string
}

func (c *countingEngine) Turn(_ context.Context, objective, _, _ string) (string, error) {
	c.turns++
	c.asked = append(c.asked, objective)
	if c.turns >= c.passAt {
		return fmt.Sprintf("attempt %d %s RAN: fine", c.turns, play.ToolVerdictHead), nil
	}
	return fmt.Sprintf("attempt %d %s FAILED: not yet", c.turns, play.ToolVerdictHead), nil
}

// loopSpec: brief -> attempt(loops) -> verify -> verdict; the verdict returns
// to attempt on fail and reaches a gate on pass. `attempt` is told its pass
// and why the last one came back, so the retry is not blind.
func loopSpec(loops int) Spec {
	return Spec{Name: "looped", BudgetS: 600,
		Nodes: []Node{
			{Name: "brief", Kind: "ask", Question: "Q"},
			{Name: "attempt", Kind: "run", Question: "do it (pass {{pass_attempt}}) {{fail_attempt}}", Loops: loops},
			{Name: "verify", Kind: "ask", Question: "V {{out_attempt}}"},
			{Name: "verdict", Kind: "eval", Ref: "attempt", Expected: "RAN:", Match: "contains"},
			{Name: "land", Kind: "gate", Title: "land after {{pass_attempt}} pass(es)?"},
		},
		Edges: []Edge{
			{From: "brief", To: "attempt", When: "always"},
			{From: "attempt", To: "verify", When: "always"},
			{From: "verify", To: "verdict", When: "always"},
			{From: "verdict", To: "attempt", When: "fail"},
			{From: "verdict", To: "land", When: "pass"},
		}}
}

// feedEngine records what each `run` turn is handed as its FEED beside its
// objective, so a stroke can hold the carried pass to the channel it rides.
type feedEngine struct {
	stubEngine
	turns int
	asked []string
	feeds []string
}

func (f *feedEngine) Turn(_ context.Context, objective, feed, _ string) (string, error) {
	f.turns++
	f.asked = append(f.asked, objective)
	f.feeds = append(f.feeds, feed)
	if f.turns >= 2 {
		return fmt.Sprintf("attempt %d %s RAN: fine", f.turns, play.ToolVerdictHead), nil
	}
	return fmt.Sprintf("attempt %d %s FAILED: not yet", f.turns, play.ToolVerdictHead), nil
}

// The operator, 2026-09-29: "carry the failed pass without the door's name".
// The fifth firing of coder-tree measured why a question cannot carry it: a
// sent-back pass quotes the door's own reply, and the council reads a tool's
// name in the OBJECTIVE as a request for that tool, shutting the Coder's
// window on every retry. So the words stay the words, and the pass rides as
// the turn's feed -- the council's own second channel, shown to every seat
// and routed on by none.
func TestTheFailedPassRidesAsTheFeedNeverInTheObjective(t *testing.T) {
	home := t.TempDir()
	spec := loopSpec(2)
	spec.Nodes[1].Question = "do it" // the words name no {{fail_attempt}}
	eng := &feedEngine{}
	res, err := Run(home, eng, spec, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused || res.PausedNode != "land" {
		t.Fatalf("the loop must reach the gate on the second pass: %s at %q", res.Verdict, res.PausedNode)
	}
	if len(eng.feeds) != 2 || eng.feeds[0] != "" {
		t.Fatalf("a first pass is handed no feed, and there were two passes: %q", eng.feeds)
	}
	if !strings.HasPrefix(eng.feeds[1], "pass 1 of `attempt` was sent back by `verdict`: fail: expected contains \"RAN:\"") ||
		!strings.Contains(eng.feeds[1], "what `attempt` answered on that pass:\nattempt 1 ") {
		t.Fatalf("the retry's feed must carry the first pass's failure and its answer: %q", eng.feeds[1])
	}
	for i, obj := range eng.asked {
		if obj != "do it" {
			t.Fatalf("the objective must stay the words alone on pass %d: %q", i+1, obj)
		}
	}
}

// recordOf renders a run's node and loop lines as one line, for a golden.
func recordOf(t *testing.T, home, run string) string {
	t.Helper()
	lines, err := runLog(home, run)
	if err != nil {
		t.Fatal(err)
	}
	var seq []string
	for _, l := range lines {
		switch l["kind"] {
		case "node":
			seq = append(seq, fmt.Sprint(l["node"], ":", l["status"]))
		case "loop":
			seq = append(seq, fmt.Sprint("loop:", l["node"], ":", l["status"], ":", l["return"], "/", l["ceiling"]))
		}
	}
	return strings.Join(seq, " ")
}

func TestACheckThatFailsReturnsTheRunToTheWorkAndTheCeilingHolds(t *testing.T) {
	home := t.TempDir()
	// Passes on the third attempt, under a ceiling of two returns.
	eng := &countingEngine{passAt: 3}
	res, err := Run(home, eng, loopSpec(2), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused || res.PausedNode != "land" {
		t.Fatalf("a loop that passes must reach the gate: %s at %q", res.Verdict, res.PausedNode)
	}
	if eng.turns != 3 {
		t.Fatalf("the work must be re-done, not the score: %d turns", eng.turns)
	}
	// The retry is told its pass and why the last one came back.
	if !strings.Contains(eng.asked[0], "(pass 1) ") || strings.Contains(eng.asked[0], "sent back") {
		t.Fatalf("the first pass must see pass 1 and no failure: %q", eng.asked[0])
	}
	if !strings.Contains(eng.asked[1], "(pass 2) pass 1 of `attempt` was sent back by `verdict`: fail: expected contains \"RAN:\"") ||
		!strings.Contains(eng.asked[1], "what `attempt` answered on that pass:\nattempt 1 ") {
		t.Fatalf("the second pass must carry the first's failure and answer: %q", eng.asked[1])
	}
	if !strings.Contains(eng.asked[2], "(pass 3) pass 2 of `attempt` was sent back") {
		t.Fatalf("the third pass must carry the second's failure: %q", eng.asked[2])
	}
	// Every pass is on the record, with a loop line between each.
	want := "brief:ok attempt:ok verify:ok verdict:ok loop:attempt:returned:1/2 " +
		"attempt:ok verify:ok verdict:ok loop:attempt:returned:2/2 " +
		"attempt:ok verify:ok verdict:ok land:paused"
	if got := recordOf(t, home, res.Run); got != want {
		t.Fatalf("the record must carry every pass and every return:\n got %s\nwant %s", got, want)
	}
	// The waterfall shows the returns, and the gate's title saw the pass count.
	st, err := Status(home, res.Run)
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{
		"loop     return 1 of 2 -> back to attempt; the work runs again",
		"loop     return 2 of 2 -> back to attempt",
		`why: fail: expected contains "RAN:"`,
		"land after 3 pass(es)?",
	} {
		if !strings.Contains(st, needle) {
			t.Fatalf("waterfall lacks %q:\n%s", needle, st)
		}
	}
	// The final pass's answer is the run's.
	outs, _ := nodeOutputs(home, res.Run)
	if !strings.HasPrefix(outs["attempt"], "attempt 3 ") {
		t.Fatalf("the last pass must stand as the node's answer: %q", outs["attempt"])
	}

	// The ceiling holds: never passes, two returns, then the fail is a fail
	// and the run stops FAIL with the spent ceiling on the record.
	home2 := t.TempDir()
	eng2 := &countingEngine{passAt: 99}
	res2, err := Run(home2, eng2, loopSpec(2), nil)
	if err != nil {
		t.Fatal(err)
	}
	if res2.Verdict != VerdictFail {
		t.Fatalf("a spent ceiling with no forward fail-edge must end FAIL, got %s", res2.Verdict)
	}
	if eng2.turns != 3 {
		t.Fatalf("a ceiling of 2 returns is 3 passes, got %d", eng2.turns)
	}
	if got := recordOf(t, home2, res2.Run); !strings.HasSuffix(got, "verdict:ok loop:attempt:spent:2/2") {
		t.Fatalf("the spent ceiling must close the record: %s", got)
	}
	st2, _ := Status(home2, res2.Run)
	if !strings.Contains(st2, "loop     ceiling 2 of 2 spent; attempt is not returned to again") {
		t.Fatalf("the spent ceiling must be on the waterfall:\n%s", st2)
	}
	// Zero returns declared is no loop at all: the same spec without `loops`
	// is refused as the cycle it is.
	if _, err := Validate(loopSpec(0)); err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("a back-edge to a node with no ceiling must be refused as a cycle: %v", err)
	}
}

// A spent ceiling with a forward fail-edge goes down it: the loop is exhausted
// and the run reports, rather than stopping FAIL.
func TestASpentCeilingFollowsAForwardFailEdge(t *testing.T) {
	s := loopSpec(1)
	s.Nodes = append(s.Nodes, Node{Name: "report", Kind: "ask", Question: "R {{fail_attempt}}"})
	s.Edges = append(s.Edges, Edge{From: "verdict", To: "report", When: "fail"})
	home := t.TempDir()
	eng := &countingEngine{passAt: 99}
	res, err := Run(home, eng, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictComplete {
		t.Fatalf("a spent ceiling must take the forward fail-edge to the end: %s", res.Verdict)
	}
	if eng.turns != 2 {
		t.Fatalf("one return is two passes, got %d", eng.turns)
	}
	outs, _ := nodeOutputs(home, res.Run)
	if _, ok := outs["report"]; !ok {
		t.Fatal("the forward fail-edge must have fired")
	}
	if _, ok := outs["land"]; ok {
		t.Fatal("the gate on the pass-edge must not fire on a fail")
	}
}

// A resume after a looped segment rebuilds the count and what the node was
// told, off the run's own loop lines: a node after the gate still reads
// pass_attempt and fail_attempt, and the looped work is not fired again.
func TestAResumeRebuildsTheLoopCount(t *testing.T) {
	s := loopSpec(2)
	s.Nodes = append(s.Nodes, Node{Name: "after", Kind: "ask", Question: "A {{pass_attempt}} {{fail_attempt}}"})
	s.Edges = append(s.Edges, Edge{From: "land", To: "after", When: "pass"})
	home := t.TempDir()
	eng := &countingEngine{passAt: 2}
	res, err := Run(home, eng, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Verdict != VerdictPaused {
		t.Fatalf("expected a pause at the gate, got %s", res.Verdict)
	}
	after := &stubEngine{}
	res2, err := Resume(home, after, res.Run, "continue")
	if err != nil {
		t.Fatal(err)
	}
	if res2.Verdict != VerdictComplete {
		t.Fatalf("the resumed run must complete: %s", res2.Verdict)
	}
	if len(after.calls) != 1 || !strings.HasPrefix(after.calls[0], "ask:A 2 pass 1 of `attempt` was sent back by `verdict`") {
		t.Fatalf("the node after the gate must read the rebuilt count and carry: %v", after.calls)
	}
}

// What refuses by name (LAW_003 section 5), at the save and not at the run.
func TestAReturnIsRefusedWhereItIsNotALoop(t *testing.T) {
	work := func(loops int) Node {
		return Node{Name: "w", Kind: "run", Question: "do", Loops: loops}
	}
	check := Node{Name: "c", Kind: "eval", Ref: "w", Expected: "RAN:", Match: "contains"}
	back := []Edge{{From: "w", To: "c"}, {From: "c", To: "w", When: "fail"}}
	cases := []struct {
		why   string
		spec  Spec
		wants string
	}{
		{"a back-edge to a node with no ceiling is the cycle it always was",
			Spec{Name: "p", Nodes: []Node{work(0), check}, Edges: back}, "cycle"},
		{"a ceiling above MaxLoops",
			Spec{Name: "p", Nodes: []Node{work(MaxLoops + 1), check}, Edges: back}, "the range is 0 to 5"},
		{"a ceiling nothing returns to is read by nothing",
			Spec{Name: "p", Nodes: []Node{work(2), check}, Edges: back[:1]}, "nothing returns to it"},
		{"a pass-edge back is not a return",
			Spec{Name: "p", Nodes: []Node{work(2), check},
				Edges: []Edge{{From: "w", To: "c"}, {From: "c", To: "w", When: "pass"}}}, "cycle"},
		{"an always-edge back is not a return",
			Spec{Name: "p", Nodes: []Node{work(2), {Name: "x", Kind: "ask", Question: "Q"}},
				Edges: []Edge{{From: "w", To: "x"}, {From: "x", To: "w"}}}, "cycle"},
		{"a gate's fail-edge is a stop, not a return",
			Spec{Name: "p", Nodes: []Node{work(2), {Name: "g", Kind: "gate", Title: "t"}},
				Edges: []Edge{{From: "w", To: "g"}, {From: "g", To: "w", When: "fail"}}}, "cycle"},
		{"a gate inside the loop",
			Spec{Name: "p", Nodes: []Node{work(2), {Name: "g", Kind: "gate", Title: "t"}, check},
				Edges: []Edge{{From: "w", To: "g"}, {From: "g", To: "c", When: "pass"}, {From: "c", To: "w", When: "fail"}}},
			"gate \"g\" stands inside the return"},
		{"loops on a check",
			Spec{Name: "p", Nodes: []Node{{Name: "a", Kind: "ask", Question: "Q"},
				{Name: "c1", Kind: "eval", Ref: "a", Expected: "x", Loops: 1},
				{Name: "c2", Kind: "eval", Ref: "a", Expected: "y"}},
				Edges: []Edge{{From: "a", To: "c1"}, {From: "c1", To: "c2", When: "pass"}, {From: "c2", To: "c1", When: "fail"}}},
			"is not returned to"},
		{"loops on a gate",
			Spec{Name: "p", Nodes: []Node{{Name: "g", Kind: "gate", Title: "t", Loops: 1}}}, "is not returned to"},
		{"a check that returns twice",
			Spec{Name: "p", Nodes: []Node{{Name: "a", Kind: "run", Question: "A", Loops: 1}, work(1), check},
				Edges: []Edge{{From: "a", To: "w"}, {From: "w", To: "c"}, {From: "c", To: "w", When: "fail"}, {From: "c", To: "a", When: "fail"}}},
			"returns twice"},
	}
	for _, c := range cases {
		_, err := Validate(c.spec)
		// A back-edge that is not a return leaves the spec with no start node
		// at all when the loop is the whole flow, and Validate says THAT first;
		// either refusal is the cycle refused, and neither is a loop.
		if err == nil || !(strings.Contains(err.Error(), c.wants) ||
			(c.wants == "cycle" && strings.Contains(err.Error(), "want exactly one start"))) {
			t.Errorf("%s: want a refusal naming %q, got %v", c.why, c.wants, err)
		}
	}
	// The lawful shapes validate: the loop above, and a return to the start node.
	if _, err := Validate(loopSpec(2)); err != nil {
		t.Fatalf("the lawful loop must validate: %v", err)
	}
	if _, err := Validate(Spec{Name: "p", Nodes: []Node{work(1), check}, Edges: back}); err != nil {
		t.Fatalf("a return to the start node must validate: %v", err)
	}
}
