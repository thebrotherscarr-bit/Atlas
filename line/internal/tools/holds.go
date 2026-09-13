package tools

// holds: the queue a WRITING call waits in when the hand behind it is not the
// operator's, and the answer that lets it through.
//
// WHY THIS EXISTS. RULE 6 is "no agent commits, pushes, lands, approves, or
// authorises a spend", and until now this door enforced that by nobody having
// wired an agent to the writing verbs. That is a convention, not a gate. The
// shape here is taken from the sovereign-microkernel read on 2026-09-12 -- its
// `high_impact` + `pending_approvals` was the one idea in that MCPServer this
// estate did not already have in a better form.
//
// IDENTITY COMES FROM THE TRANSPORT AND NEVER FROM THE MESSAGE. The hold
// decision reads the `Caller` PARAMETER the door built from the connection --
// the service wire on an HTTP request, or the plain fact of being a stdio pipe.
// It never reads the args.
//
// THAT WAS CHECKED BY MUTATION, NOT BELIEVED, and the first check was pointed
// at the wrong line: deleting the reserved key out of the inbound args changes
// no verdict at all, because the decision never consults them and the door
// overwrites them on the way out. The strip is belt and braces. Point the
// decision at the args instead -- `!callerOf(args).Service` -- and a call
// carrying a forged `__caller` writes straight through, which is what the
// stroke reports as `WROTE smuggled`.
//
// The `actor` argument RBAC already reads IS taken from the caller's own args,
// so anyone may declare themselves anyone there. A gate keyed the same way
// would be a gate in name only.
//
// AND IT IS INERT WITHOUT THE DOOR'S AUTH, WHICH IT SAYS OUT LOUD. With
// `--auth` off the door reads no credential at all, so it cannot tell a seat
// from the glass and `Service` is false for everyone. Holding on that would
// stop the operator's own panel while stopping no agent that thought to send a
// header. So holds ARM only when auth is on, `hold_list` says so in plain
// words when they are not, and the boot line names it. A guard that quietly
// does nothing is worse than no guard, because it is believed.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"atlas/line/internal/tenant"
)

// Caller is who the DOOR judged the request to be, from the transport. Never
// from the message body.
type Caller struct {
	// Name is what the client called itself at initialize, or the bearer's
	// key name. It is a label for the record, and it is not trusted.
	Name string
	// Service is the one bit that decides anything: the request carried the
	// service wire, so it is the operator's own glass.
	Service bool
}

// CallerKey is the reserved args key the door writes its judgement into, and
// strips from anything a caller sent. Reserved spelling so a real tool
// argument can never collide with it.
const CallerKey = "__caller"

// Hold is one parked call, kept whole so approving it runs exactly what was
// asked -- not a re-reading of it.
type Hold struct {
	ID      string         `json:"id"`
	Tool    string         `json:"tool"`
	Args    map[string]any `json:"args"`
	Caller  string         `json:"caller"`
	Project string         `json:"project"`
	When    time.Time      `json:"when"`
}

// heldTools never hold, whatever they declare. `hold_answer` writes by
// definition -- it runs the held call -- and holding the answer behind another
// answer is a queue that can never drain.
var heldExempt = map[string]bool{"hold_answer": true, "hold_list": true}

func (r *Registry) parkedPath(home string) string {
	return filepath.Join(home, "state", "holds.jsonl")
}

// park queues the call and writes it to the record before answering. The
// record is written FIRST: a hold the operator cannot find later is a call
// that vanished.
func (r *Registry) park(tn tenant.Tenant, t Tool, args map[string]any, caller Caller) Hold {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.held == nil {
		r.held = map[string]Hold{}
	}
	// A SEQUENCE, NOT ONLY A CLOCK. Two parks inside the same millisecond
	// produced the SAME id, and the second silently REPLACED the first in the
	// map: one call vanished, and approving the surviving id would have run
	// the other one's arguments. Found by its own stroke rather than in the
	// field, because two agent writes in one millisecond is a Tuesday.
	r.heldSeq++
	h := Hold{
		ID:      fmt.Sprintf("hold_%d_%d_%s", time.Now().UnixMilli(), r.heldSeq, t.Name),
		Tool:    t.Name,
		Args:    copyArgs(args),
		Caller:  callerLabel(caller),
		Project: tn.Name,
		When:    time.Now().UTC(),
	}
	r.held[h.ID] = h
	r.record(tn.Home, "held", h, "")
	return h
}

// record appends one line to the ground's own holds log. Append-only, like
// every other ledger here (ESTATE LAW 8).
func (r *Registry) record(home, what string, h Hold, note string) {
	line := map[string]any{
		"ts": time.Now().UTC().Format(time.RFC3339), "what": what,
		"id": h.ID, "tool": h.Tool, "caller": h.Caller, "project": h.Project,
	}
	if note != "" {
		line["note"] = note
	}
	b, err := json.Marshal(line)
	if err != nil {
		return
	}
	p := r.parkedPath(home)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return
	}
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(append(b, '\n'))
}

func copyArgs(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		if k == CallerKey {
			continue
		}
		out[k] = v
	}
	return out
}

func callerLabel(c Caller) string {
	if c.Service {
		return "operator (service wire)"
	}
	if c.Name == "" {
		return "unnamed client"
	}
	return c.Name
}

// heldAnswer is what the agent gets back. A REFUSAL STRING, not an error: the
// call did not fail, it is waiting, and the difference matters to whatever is
// reading the answer. It says who must act and where, because a seat told only
// "held" will invent the rest.
func heldAnswer(h Hold) string {
	return fmt.Sprintf("HELD: %q writes, and this call did not come from the "+
		"operator's own glass (%s). Nothing has been written.\n\n"+
		"It is parked as %s and waits for his hand -- RULE 6: no agent commits, "+
		"pushes, lands, approves, or authorises a spend. He answers it on the "+
		"Version control page, or with hold_answer.\n\n"+
		"Do not retry this call and do not work around it. Say plainly that it "+
		"is waiting.", h.Tool, h.Caller, h.ID)
}

// callerOf reads the door's own judgement back out of the args it wrote.
func callerOf(args map[string]any) Caller {
	if c, ok := args[CallerKey].(Caller); ok {
		return c
	}
	return Caller{}
}

// --- the two tools ----------------------------------------------------------

func toolHoldList(t tenant.Tenant, args map[string]any) (string, error) {
	reg := regOf(args)
	if reg == nil {
		return "Refused: the hold queue is not reachable from here.", nil
	}
	if !callerOf(args).Service {
		return "Refused: the hold queue is the operator's. This asks who is " +
			"calling, and the answer was not his glass.", nil
	}
	reg.mu.Lock()
	defer reg.mu.Unlock()
	out := map[string]any{
		"armed": reg.holdWrites,
		"held":  []map[string]any{},
	}
	if !reg.holdWrites {
		out["why_not"] = "The door was started WITHOUT --auth, so it cannot tell " +
			"a seat from the glass and nothing is being held. Restart it with " +
			"--auth and a service wire to arm this."
	}
	ids := make([]string, 0, len(reg.held))
	for id := range reg.held {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	rows := []map[string]any{}
	for _, id := range ids {
		h := reg.held[id]
		rows = append(rows, map[string]any{
			"id": h.ID, "tool": h.Tool, "caller": h.Caller,
			"project": h.Project, "when": h.When.Format(time.RFC3339),
			"args": h.Args,
		})
	}
	out["held"] = rows
	b, err := json.MarshalIndent(out, "", " ")
	return string(b), err
}

func toolHoldAnswer(t tenant.Tenant, args map[string]any) (string, error) {
	reg := regOf(args)
	if reg == nil {
		return "Refused: the hold queue is not reachable from here.", nil
	}
	// THE WHOLE GATE IS THIS LINE. If anything but the operator's glass can
	// answer a hold, the hold protects nothing -- an agent parks its own write
	// and then waves it through.
	if !callerOf(args).Service {
		return "Refused: answering a hold is the operator's act alone, and this " +
			"call did not come from his glass. RULE 6.", nil
	}
	id, _ := args["id"].(string)
	decision, _ := args["decision"].(string)
	if id == "" {
		return "Refused: name the hold to answer, e.g. id=hold_1757...", nil
	}
	if decision != "approve" && decision != "deny" {
		return fmt.Sprintf("Refused: a hold is answered approve or deny, got %q.", decision), nil
	}

	reg.mu.Lock()
	h, ok := reg.held[id]
	if ok {
		delete(reg.held, id)
	}
	reg.mu.Unlock()
	if !ok {
		return fmt.Sprintf("Refused: no hold called %s is waiting. It may have "+
			"been answered already, or the door restarted -- a restart drops the "+
			"queue, and a dropped hold means the call never ran.", id), nil
	}

	if decision == "deny" {
		reg.record(t.Home, "denied", h, "")
		return fmt.Sprintf("Denied %s. %q was not run and nothing was written.", h.ID, h.Tool), nil
	}

	tool, exists := reg.Get(h.Tool)
	if !exists {
		reg.record(t.Home, "vanished", h, "the tool no longer exists")
		return fmt.Sprintf("Refused: %q no longer exists at this door, so the "+
			"held call cannot be run.", h.Tool), nil
	}
	// Run it AS ITSELF, with the args it was parked with -- not re-read, not
	// re-judged. The operator approved that call, not a fresh one.
	run := copyArgs(h.Args)
	run[CallerKey] = Caller{Name: "operator via hold " + h.ID, Service: true}
	out, err := tool.Fn(t, run)
	if err != nil {
		reg.record(t.Home, "approved_errored", h, err.Error())
		return fmt.Sprintf("Approved %s and %q errored: %v\n\n%s", h.ID, h.Tool, err, out), nil
	}
	reg.record(t.Home, "approved_ran", h, "")
	return fmt.Sprintf("Approved %s. %q ran:\n\n%s", h.ID, h.Tool, out), nil
}

// regOf hands a tool the registry it lives in. Set by Call beside the caller,
// for the two tools whose whole job is the registry's own queue.
func regOf(args map[string]any) *Registry {
	r, _ := args[registryKey].(*Registry)
	return r
}

const registryKey = "__registry"

// mu and held live on Registry; declared here beside everything that uses them.
type holdState struct {
	mu         sync.Mutex
	held       map[string]Hold
	heldSeq    int
	holdWrites bool
}
