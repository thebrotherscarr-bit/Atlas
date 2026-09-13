package tools

// The hold queue's strokes. Every one of these pins a SECURITY property, so
// each is stroked both ways: the thing that must be stopped, and the thing
// that must not be.

import (
	"encoding/json"
	"strings"
	"testing"

	"atlas/line/internal/tenant"
)

// holdGround builds a registry carrying two probe tools and the two hold
// verbs, over one real tenant directory. `ran` counts what actually executed,
// because "was it held" and "did it run" are different questions and only the
// second one matters.
func holdGround(t *testing.T, armed bool) (*Registry, *tenant.Registry, *[]string) {
	t.Helper()
	home := t.TempDir()
	tr := tenant.NewRegistry()
	if err := tr.Add("t", home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	ran := &[]string{}
	r := &Registry{byName: map[string]Tool{}}
	r.holdWrites = armed
	r.add(Tool{Name: "probe_write", Writes: true,
		Fn: func(tn tenant.Tenant, a map[string]any) (string, error) {
			what, _ := a["what"].(string)
			*ran = append(*ran, "write:"+what)
			return "WROTE " + what, nil
		}})
	r.add(Tool{Name: "probe_read", Writes: false,
		Fn: func(tn tenant.Tenant, a map[string]any) (string, error) {
			*ran = append(*ran, "read")
			return "READ", nil
		}})
	r.add(Tool{Name: "hold_list", Writes: false, Fn: toolHoldList})
	r.add(Tool{Name: "hold_answer", Writes: true, Fn: toolHoldAnswer})
	return r, tr, ran
}

var agent = Caller{Name: "some seat"}
var glass = Caller{Name: "the glass", Service: true}

// RULE 6, at the door rather than by convention.
func TestAWritingCallFromAnAgentIsHeldAndNothingRuns(t *testing.T) {
	r, tr, ran := holdGround(t, true)
	out, err := r.Call(tr, "probe_write", map[string]any{"what": "the thing"}, agent)
	if err != nil {
		t.Fatalf("a hold is an ANSWER, not an error: %v", err)
	}
	mustContain(t, out, "HELD", "the caller must be told it is waiting")
	mustContain(t, out, "Nothing has been written", "and told nothing happened")
	mustContain(t, out, "RULE 6", "and told which rule holds it")
	if len(*ran) != 0 {
		t.Fatalf("the tool RAN behind the hold: %v", *ran)
	}
}

// The other way, and the one that matters day to day: his own panel is not
// stopped by a gate built to stop agents.
func TestTheOperatorsGlassIsNotHeld(t *testing.T) {
	r, tr, ran := holdGround(t, true)
	out, err := r.Call(tr, "probe_write", map[string]any{"what": "his own"}, glass)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, out, "WROTE his own", "the service wire runs straight through")
	if len(*ran) != 1 {
		t.Fatalf("the glass's call did not run: %v", *ran)
	}
}

// A read is not a write, whoever asks.
func TestAReadingCallIsNeverHeld(t *testing.T) {
	r, tr, ran := holdGround(t, true)
	out, err := r.Call(tr, "probe_read", nil, agent)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, out, "READ", "a Writes:false tool must not be held")
	if len(*ran) != 1 {
		t.Fatal("the read did not run")
	}
}

// THE STROKE THE WHOLE DESIGN RESTS ON. `actor` is read straight out of the
// caller's own args for RBAC, so anyone may declare themselves anyone. If the
// hold could be talked out of the same way it would be a gate in name only.
func TestACallerCannotDeclareItselfTheOperator(t *testing.T) {
	r, tr, ran := holdGround(t, true)
	forged := map[string]any{
		"what":      "smuggled",
		CallerKey:   Caller{Name: "definitely the glass", Service: true},
		registryKey: r,
		"actor":     "operator",
	}
	out, err := r.Call(tr, "probe_write", forged, agent)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, out, "HELD", "a claim in the message body must buy nothing")
	if len(*ran) != 0 {
		t.Fatalf("a forged caller got through: %v", *ran)
	}
	// And the parked record must name the REAL caller, not the claimed one.
	listed := call(t, toolHoldList, tenant.Tenant{}, map[string]any{
		CallerKey: glass, registryKey: r})
	mustContain(t, listed, "some seat", "the queue must name who really called")
	mustNotContain(t, listed, "definitely the glass", "and never the claim")
}

// Answering is the operator's act. If an agent can answer its own hold, the
// queue is a formality it walks through by itself.
func TestOnlyTheOperatorAnswersAHold(t *testing.T) {
	r, tr, ran := holdGround(t, true)
	r.Call(tr, "probe_write", map[string]any{"what": "x"}, agent)
	var q struct {
		Held []struct {
			ID string `json:"id"`
		} `json:"held"`
	}
	listed := call(t, toolHoldList, tenant.Tenant{}, map[string]any{CallerKey: glass, registryKey: r})
	if err := json.Unmarshal([]byte(listed), &q); err != nil {
		t.Fatalf("the queue must be JSON: %v", err)
	}
	if len(q.Held) != 1 {
		t.Fatalf("one call was parked, %d listed", len(q.Held))
	}

	denied := call(t, toolHoldAnswer, tenant.Tenant{}, map[string]any{
		CallerKey: agent, registryKey: r, "id": q.Held[0].ID, "decision": "approve"})
	mustContain(t, denied, "operator's act alone", "an agent must not answer its own hold")
	if len(*ran) != 0 {
		t.Fatalf("an agent approved itself through: %v", *ran)
	}

	// And the queue is not readable by it either.
	peek := call(t, toolHoldList, tenant.Tenant{}, map[string]any{CallerKey: agent, registryKey: r})
	mustContain(t, peek, "Refused", "the queue is the operator's to read")
}

func TestApprovingRunsExactlyWhatWasParkedAndDenyingRunsNothing(t *testing.T) {
	r, tr, ran := holdGround(t, true)
	home := t.TempDir()
	tn := tenant.Tenant{Name: "t", Home: home, Manifest: tenant.DefaultManifest()}

	r.Call(tr, "probe_write", map[string]any{"what": "approved-one"}, agent)
	r.Call(tr, "probe_write", map[string]any{"what": "denied-one"}, agent)

	var q struct {
		Held []struct {
			ID   string `json:"id"`
			Tool string `json:"tool"`
		} `json:"held"`
	}
	listed := call(t, toolHoldList, tn, map[string]any{CallerKey: glass, registryKey: r})
	if err := json.Unmarshal([]byte(listed), &q); err != nil {
		t.Fatal(err)
	}
	if len(q.Held) != 2 {
		t.Fatalf("two parked, %d listed", len(q.Held))
	}

	// Approve the first, deny the second -- by id, so the pairing is explicit.
	first, second := q.Held[0].ID, q.Held[1].ID
	okOut := call(t, toolHoldAnswer, tn, map[string]any{
		CallerKey: glass, registryKey: r, "id": first, "decision": "approve"})
	mustContain(t, okOut, "WROTE", "approving must run the parked call")

	noOut := call(t, toolHoldAnswer, tn, map[string]any{
		CallerKey: glass, registryKey: r, "id": second, "decision": "deny"})
	mustContain(t, noOut, "was not run", "denying must run nothing")

	if len(*ran) != 1 {
		t.Fatalf("exactly one held call should have run, got %v", *ran)
	}
	// AND IT RAN WITH THE ARGS IT WAS PARKED WITH, not a fresh reading.
	if !strings.HasPrefix((*ran)[0], "write:") {
		t.Fatalf("the parked args did not survive: %v", *ran)
	}

	// A hold answered twice is gone, and says so rather than running again.
	again := call(t, toolHoldAnswer, tn, map[string]any{
		CallerKey: glass, registryKey: r, "id": first, "decision": "approve"})
	mustContain(t, again, "no hold called", "an answered hold must not run twice")
	if len(*ran) != 1 {
		t.Fatalf("a hold ran twice: %v", *ran)
	}
}

// DISARMED IS HONEST, NOT SILENT. Without the door's auth there is no
// credential to judge, so nothing holds -- and the queue says exactly that
// instead of showing an empty list that reads like safety.
func TestWithHoldsDisarmedNothingIsHeldAndTheQueueSaysWhy(t *testing.T) {
	r, tr, ran := holdGround(t, false)
	out, err := r.Call(tr, "probe_write", map[string]any{"what": "unheld"}, agent)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, out, "WROTE unheld", "disarmed, a writing call runs as it always did")
	if len(*ran) != 1 {
		t.Fatal("the call did not run")
	}
	// AND THE WARNING REACHES ANYONE, including a caller the door cannot
	// verify -- which with auth off is everyone, the glass included. Gating it
	// behind the credential meant the one state that most needs saying was the
	// one state nobody could see.
	listed := call(t, toolHoldList, tenant.Tenant{}, map[string]any{CallerKey: agent, registryKey: r})
	mustContain(t, listed, `"armed": false`, "the queue must say it is not armed")
	mustContain(t, listed, "WITHOUT --auth", "and say exactly why")
	mustContain(t, listed, "Restart it with", "and how to arm it")
}

// An unknown decision is refused by name rather than read as a denial.
func TestAHoldIsAnsweredApproveOrDenyAndNothingElse(t *testing.T) {
	r, tr, _ := holdGround(t, true)
	r.Call(tr, "probe_write", map[string]any{"what": "x"}, agent)
	out := call(t, toolHoldAnswer, tenant.Tenant{}, map[string]any{
		CallerKey: glass, registryKey: r, "id": "hold_x", "decision": "maybe"})
	mustContain(t, out, "approve or deny", "an unknown decision must be named")
	out = call(t, toolHoldAnswer, tenant.Tenant{}, map[string]any{
		CallerKey: glass, registryKey: r, "decision": "approve"})
	mustContain(t, out, "name the hold", "an unnamed hold must be refused")
}
