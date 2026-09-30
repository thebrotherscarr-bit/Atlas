// Strokes for what THE LINE puts on the engine's wire (2026-09-29).
//
// A hermetic stroke cannot stand up Manjuel, and a mock of the council would
// prove the mock. What it CAN stand up is the other end of the pipe: this test
// binary run again as a child that speaks serve.py's wire and hands every
// objective row back as its delivery. So these strokes prove what the door
// SENDS -- the head, and whether anybody is at the prompt -- and nothing about
// what an engine does with it, which is the core's own suite.
package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"atlas/line/internal/tenant"
)

const stubEngineEnv = "ATLAS_STUB_ENGINE"

// TestStubEngineProcess is not a stroke. It is the engine the strokes below
// spawn: this test binary run again as a child, speaking serve.py's wire --
// `opened` at start, and for every objective a delivery whose text is THE ROW
// AS IT ARRIVED, so a stroke can read what crossed the wire. The objective
// "ask-me" stops on a question instead, the way a core that needs the hand
// does. Run any other way, it skips.
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
			if row["text"] == "ask-me" {
				say(map[string]any{"event": "needs_answer", "prompt": "which one?"})
				continue
			}
			b, _ := json.Marshal(row)
			say(map[string]any{"event": "delivery", "text": string(b)})
		case "answer":
			say(map[string]any{"event": "delivery", "text": "answered"})
		}
	}
	os.Exit(0)
}

// standing carries one fresh world in a fresh registry and opens the stub
// engine on it; the engine is closed with the stroke.
func standing(t *testing.T) (*tenant.Registry, tenant.Tenant) {
	t.Helper()
	tr := tenant.NewRegistry()
	if err := tr.Add("t", t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	tn, err := tr.Resolve("t")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv(stubEngineEnv, "1")
	if _, err := engines.Open("t", tn.Home,
		`"`+os.Args[0]+`" -test.run=^TestStubEngineProcess$ --`); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = engines.CloseOne(tn.Home) })
	return tr, tn
}

// sent is the row the door put on the wire, read back out of what the stub
// handed up. json.Marshal writes a map's keys in order, so the row begins at
// its `cmd`.
func sent(t *testing.T, out string) map[string]any {
	t.Helper()
	i := strings.Index(out, `{"cmd"`)
	if i < 0 {
		t.Fatalf("the row the door sent is not in what came back:\n%s", out)
	}
	var row map[string]any
	if err := json.NewDecoder(strings.NewReader(out[i:])).Decode(&row); err != nil {
		t.Fatalf("the row the door sent does not read back: %v\n%s", err, out)
	}
	return row
}

// A FLOW'S TURN IS SENT UNATTENDED, AND HIS OWN IS NOT (2026-09-29). The core
// asks "retry / skip / abort?" when a seat marked `On Fail: prompt` fails, and
// a flow cannot answer; told that nobody is at the prompt it skips the seat,
// the question's own default, and says so. Until this the door sent every
// turn as though his hand were on the keyboard, so a flow died at the first
// such question. The word rides ONLY the flow's turn: run_start and the
// glass's stream are his, and say nothing.
func TestAFlowsTurnIsSentUnattendedAndHisOwnIsNot(t *testing.T) {
	tr, tn := standing(t)
	out, err := council(tn.Home).Turn(context.Background(), "echo", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if row := sent(t, out); row["unattended"] != true {
		t.Fatalf("a flow's turn crossed the wire as though someone were at the prompt: %v", row)
	}

	reg := Build(tr, Options{})
	out, err = reg.Call(tr, "run_start", map[string]any{"objective": "echo"}, Caller{Service: true})
	if err != nil {
		t.Fatal(err)
	}
	if row := sent(t, out); row["unattended"] != nil {
		t.Fatalf("run_start, his own turn, said nobody was at the prompt: %v", row)
	}
	res, err := RunStream(tn, "echo", "", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if row := sent(t, res.Final.Str("text")); row["unattended"] != nil {
		t.Fatalf("the glass's stream, his own turn, said nobody was at the prompt: %v", row)
	}

	// AND A QUESTION THE CORE STILL ASKS IS STILL HIS (RULE 6): the flow's turn
	// stops on it and is refused in words, never answered for him.
	_, err = council(tn.Home).Turn(context.Background(), "ask-me", "", "")
	if err == nil || !strings.Contains(err.Error(), "stopped on a question") {
		t.Fatalf("a turn that stopped on a question was not refused as one: %v", err)
	}
	if _, err := AnswerStream(tn, "this one", nil); err != nil {
		t.Fatalf("the question could not be answered afterwards: %v", err)
	}
}

// run_start TAKES A HEAD THE WAY flow_run DOES (2026-09-29): `voice`, every
// seat on one model for the turn, and `voices`, a model per seat over it --
// the same two words, read by the same reader, so a head that flow_run would
// refuse is refused here in run_start's own name and nothing runs. No head
// named sends nothing: the ground's declared targets, the wire as it was.
func TestRunStartTakesAHeadTheWayFlowRunDoes(t *testing.T) {
	tr, tn := standing(t)
	reg := Build(tr, Options{})
	tool, ok := reg.Get("run_start")
	if !ok {
		t.Fatal("run_start is gone from the registry")
	}
	for _, a := range []string{"voice?", "voices?"} {
		if !hasArg(tool.Args, a) {
			t.Fatalf("run_start does not offer %s over the wire: %v", a, tool.Args)
		}
	}
	call := func(args map[string]any) (string, error) {
		return reg.Call(tr, "run_start", args, Caller{Service: true})
	}

	out, err := call(map[string]any{"objective": "echo", "voice": " head-a:latest ",
		"voices": `{"Steward":"phi4-mini:latest"}`})
	if err != nil {
		t.Fatal(err)
	}
	row := sent(t, out)
	if row["model"] != "head-a:latest" {
		t.Fatalf("the roster-wide head did not reach the wire: %v", row)
	}
	if v, _ := row["voices"].(map[string]any); v["Steward"] != "phi4-mini:latest" {
		t.Fatalf("the seat map did not reach the wire: %v", row)
	}

	out, err = call(map[string]any{"objective": "echo"})
	if err != nil {
		t.Fatal(err)
	}
	row = sent(t, out)
	if _, named := row["model"]; named {
		t.Fatalf("no head was named and one crossed the wire: %v", row)
	}
	if _, named := row["voices"]; named {
		t.Fatalf("no seat map was named and one crossed the wire: %v", row)
	}

	e, ok := engines.Get(tn.Home)
	if !ok {
		t.Fatal("the stub engine is gone")
	}
	before, _ := e.Runs()
	_, err = call(map[string]any{"objective": "echo", "voices": "Steward"})
	if err == nil || !strings.Contains(err.Error(), "run_start voices") {
		t.Fatalf("a seat map that cannot be read must be refused in run_start's own name: %v", err)
	}
	if after, _ := e.Runs(); after != before {
		t.Fatal("a refused head still ran a turn")
	}
	if _, err := flowVoices(map[string]any{"voices": 42}); err == nil ||
		!strings.Contains(err.Error(), "flow voices") {
		t.Fatalf("flow_run's refusal no longer names the flow: %v", err)
	}
}
