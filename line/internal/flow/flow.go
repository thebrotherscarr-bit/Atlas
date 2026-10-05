// Package flow is N2's workflow builder: versioned DAG specs with measured
// runs, REVIEW gates, eval branches, compare and replay.
//
// Specs live in <home>/flows/<name>.json — {name, version, budget_s,
// nodes, edges}; history folds as <name>.v<k>.json, never rewritten.
// Node kinds are a closed set: ask | prompt | seat | memory | eval | gate |
// run | aider -- `run` drives a whole Manjuel turn (the council), `aider` is an
// attempt by Aider on the files it is handed (and the council's own turn when
// Aider cannot take it), the others one voice.
// Branches declare parallelism but run sequentially in topo order — one
// rack queue, no interleaved output, the queue visible in the waterfall.
// There is no verb here that finishes a task, lands a memory, or closes a
// gate on the operator's behalf: gate nodes pause, the hand resumes with
// continue|stop, evals only steer edges.
package flow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// NameRe pins the flow/node name law (cutter holds the same pattern).
var NameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

// RunRe pins the run-id shape law (cutter holds the same pattern).
var RunRe = regexp.MustCompile(`^f-\d{8}-\d{6}-[0-9a-f]{8}$`)

// Kinds is the closed node set. Anything else is refused by name.
var Kinds = map[string]bool{
	"ask": true, "prompt": true, "seat": true,
	"memory": true, "eval": true, "gate": true,
	// `run` is the whole council, not one voice: the objective goes through
	// Manjuel, so the law gate stamps it, the Router runs the tools, the dedup
	// refuses a repeat and the recompose puts every failure in the answer. An
	// `ask` node reaches a bare model; a `run` node reaches the estate.
	"run": true,
	// `aider` is an attempt by Aider (2026-10-05, WHAT'S LEFT H17, his ruling on
	// B21): the instruction and the files it names go to the door's Aider -- the
	// operator's own hand, never a seat's -- and when Aider cannot take them (a
	// file past its window, nothing it would change) the SAME instruction is the
	// council's own turn, as a `run` node would have made it. It is reached only
	// through a gate whose `grants` name AiderTool, refused at the save otherwise.
	"aider": true,
}

// AiderTool is the door tool an `aider` node asks. A gate must grant it on every
// path to the node (Validate), and on the door's side only a flow's `aider` node,
// beside his own glass, may reach it -- a stroke there holds the two names to
// each other, so a rename on one side is loud on the other.
const AiderTool = "aider_run"

// Verdicts. COMPLETE means every reached node came back ok; the rest name
// exactly how a run stopped. No verdict here finishes anyone's task.
const (
	VerdictComplete = "COMPLETE"
	VerdictPaused   = "PAUSED"
	VerdictFail     = "FAIL"
	VerdictOverTime = "OUT_OF_TIME"
	VerdictStopped  = "STOPPED"
)

// Node is one step: kind plus the fields its kind reads.
type Node struct {
	Name     string            `json:"name"`
	Kind     string            `json:"kind"`
	Voice    string            `json:"voice,omitempty"`
	Question string            `json:"question,omitempty"`
	Prompt   string            `json:"prompt,omitempty"`
	Version  int               `json:"version,omitempty"`
	Vars     map[string]string `json:"vars,omitempty"`
	Seat     string            `json:"seat,omitempty"`
	Method   string            `json:"method,omitempty"`
	Ref      string            `json:"node,omitempty"`
	Expected string            `json:"expected,omitempty"`
	Match    string            `json:"match,omitempty"`
	Title    string            `json:"title,omitempty"`
	Retries  int               `json:"retries,omitempty"`
	// Grants is a GATE's declaration of what the hand that crosses it
	// authorises (2026-09-28): the writing tools the council may call, from
	// this gate until the next gate or the end of the run, without each call
	// parking for a second decision. The gate's title is the question; the
	// grants are what `continue` answers. Only a gate carries them; THE LINE
	// refuses a grant naming a tool it does not carry or one that does not
	// write, and writes every call that rode a crossing to the holds record.
	Grants []string `json:"grants,omitempty"`
	// Loops is how many times this node may be RETURNED TO (2026-09-28,
	// LAW_003's mechanism, on the operator's word: "create the bounded
	// back-edge looping"). A check downstream that FAILS may send the run back
	// here, and every node between is fired again -- the WORK re-done, never
	// an answer re-scored -- at most this many times. The ceiling is declared
	// where a reader meets it and read by the loop (MaxLoops bounds it); the
	// stop condition is the check, a machine over the machine's own evidence;
	// every pass is a fresh line in the record, with a `loop` line between.
	// Only a check's `fail` edge may return here, only around a body that does
	// work, never around a gate, and a node that declares this with nothing
	// returning to it is refused: a ceiling read by nothing.
	Loops int `json:"loops,omitempty"`
	// Files is an `aider` node's file list (2026-10-05, H17): templated like a
	// question, one path to a line or separated by commas; the door judges every
	// path. Only an aider node carries it.
	Files string `json:"files,omitempty"`
}

// MaxLoops caps how many times a node may be returned to. BOUNDED EVERYTHING
// (ESTATE LAW 7), as a number in the spec rather than a prohibition: LAW_003
// §6.
const MaxLoops = 5

// MaxRetries caps what a node may ask for. BOUNDED EVERYTHING (ESTATE LAW 7):
// an unbounded retry is an indefinite ticker wearing a different hat, and the
// budget is the only other thing standing between a stuck engine and a run
// that never ends.
const MaxRetries = 5

// Matches is the closed set of tests an eval node may make of the answer it
// checks. Empty means `equals`, so every spec folded before this existed keeps
// the verdict it already had -- a scoring rule that changes under saved runs is
// a rewritten record.
//
// WHY `contains` HAD TO EXIST (2026-09-12). `play.Score` is exact match after
// trim and casefold, and it is also what scores prompt-eval datasets, so it
// could not simply be loosened. Meanwhile the builder's own label for this
// field read "what the answer should carry" -- which is `contains`, in words.
// The coder flow believed the label: its check expected RAN against a `run`
// node, whose answer is the council's prose. `verify` came back
// "RAN: fizz_buzz.py" over correct FizzBuzz and the check failed anyway. The
// pass branch had never once been reachable. The label was right about what
// the field is for; the engine offered no way to mean it.
var Matches = map[string]bool{"equals": true, "contains": true}

// MatchMode is the test this node makes, defaulted and lowercased.
func (n Node) MatchMode() string {
	if m := strings.ToLower(strings.TrimSpace(n.Match)); m != "" {
		return m
	}
	return "equals"
}

// Edge steers: always fires from a fired source; pass/fail follow the
// source's check outcome (eval score, gate choice).
type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
	When string `json:"when"`
}

// Spec is one versioned workflow.
type Spec struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
	BudgetS int    `json:"budget_s"`
	Nodes   []Node `json:"nodes"`
	Edges   []Edge `json:"edges"`
}

func flowsDir(home string) string { return filepath.Join(home, "flows") }
func specPath(home, name string, v int) string {
	if v <= 0 {
		return filepath.Join(flowsDir(home), name+".json")
	}
	return filepath.Join(flowsDir(home), fmt.Sprintf("%s.v%d.json", name, v))
}
func runsPath(home string) string { return filepath.Join(flowsDir(home), "runs.jsonl") }

// RunID mints f-YYYYMMDD-HHMMSS-<8hex>.
func RunID() (string, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("f-%s-%s", time.Now().UTC().Format("20060102-150405"), hex.EncodeToString(b[:])), nil
}

// Receipt binds one node outcome (cutter reproduces this byte-for-byte).
func Receipt(run, node, output, ts string) string {
	h := sha256.Sum256([]byte(run + "\n" + node + "\n" + output + "\n" + ts))
	return hex.EncodeToString(h[:])
}

// OverBudget is pure: wall milliseconds summed past the ceiling in seconds.
func OverBudget(elapsedMs []int64, budgetS int) bool {
	var sum int64
	for _, e := range elapsedMs {
		sum += e
	}
	return sum > int64(budgetS)*1000
}

// Validate judges a spec: names, kinds, refs, one start, full reach, and no
// cycle but a check's bounded RETURN (LAW_003: a fail-edge back to a node that
// declares `loops`). It returns the deterministic fire order (Kahn over the
// forward edges, name-sorted); a return is judged against that order.
func Validate(s Spec) ([]string, error) {
	if !NameRe.MatchString(s.Name) {
		return nil, fmt.Errorf("refused: flow name %q breaks the name law", s.Name)
	}
	if len(s.Nodes) == 0 {
		return nil, fmt.Errorf("refused: flow %q carries no nodes", s.Name)
	}
	byName := map[string]Node{}
	for _, n := range s.Nodes {
		if !NameRe.MatchString(n.Name) {
			return nil, fmt.Errorf("refused: node name %q breaks the name law", n.Name)
		}
		if byName[n.Name].Name != "" {
			return nil, fmt.Errorf("refused: duplicate node name %q", n.Name)
		}
		if !Kinds[n.Kind] {
			return nil, fmt.Errorf("refused: node %q carries unknown kind %q", n.Name, n.Kind)
		}
		if n.Kind == "run" && strings.TrimSpace(n.Question) == "" {
			return nil, fmt.Errorf("refused: run node %q has no objective", n.Name)
		}
		if n.Kind == "aider" && strings.TrimSpace(n.Question) == "" {
			return nil, fmt.Errorf("refused: aider node %q has no instruction", n.Name)
		}
		if n.Kind == "aider" && strings.TrimSpace(n.Files) == "" {
			return nil, fmt.Errorf("refused: aider node %q names no files; Aider works on the files it is given", n.Name)
		}
		if n.Kind != "aider" && strings.TrimSpace(n.Files) != "" {
			return nil, fmt.Errorf("refused: node %q is a %s and names no files for Aider, so "+
				"`files` means nothing on it", n.Name, n.Kind)
		}
		if n.Kind == "eval" && strings.TrimSpace(n.Ref) == "" {
			return nil, fmt.Errorf("refused: eval node %q names no node to check", n.Name)
		}
		if n.Kind == "eval" && !Matches[n.MatchMode()] {
			return nil, fmt.Errorf("refused: eval node %q carries unknown match %q; "+
				"the set is equals, contains", n.Name, n.Match)
		}
		if n.Kind != "eval" && strings.TrimSpace(n.Match) != "" {
			return nil, fmt.Errorf("refused: node %q is a %s and has no answer to "+
				"test, so `match` means nothing on it", n.Name, n.Kind)
		}
		if n.Retries < 0 || n.Retries > MaxRetries {
			return nil, fmt.Errorf("refused: node %q asks for %d retries; the range "+
				"is 0 to %d", n.Name, n.Retries, MaxRetries)
		}
		// A VERDICT IS NOT RETRIED. An eval scores the same answer the same way
		// every time, so a second run of it can only return what the first did --
		// and if it ever did not, retrying until the check agrees is the exact
		// laundering path this estate spent 2026-09-12 closing. A gate never
		// executes at all; it pauses. Both are refused by name so nobody reads
		// `retries` on them as a way to argue with a FAIL.
		if n.Retries > 0 && (n.Kind == "eval" || n.Kind == "gate") {
			return nil, fmt.Errorf("refused: node %q is a %s, and a %s is not "+
				"retried -- retry answers an ERROR, never a verdict", n.Name, n.Kind, n.Kind)
		}
		// ONLY A GATE GRANTS. A grant is what a hand's `continue` authorises,
		// and nothing but a gate is ever answered by a hand.
		if len(n.Grants) > 0 && n.Kind != "gate" {
			return nil, fmt.Errorf("refused: node %q is a %s and grants nothing -- only "+
				"a gate is crossed by a hand, so only a gate carries `grants`", n.Name, n.Kind)
		}
		for _, g := range n.Grants {
			if strings.TrimSpace(g) == "" {
				return nil, fmt.Errorf("refused: gate %q grants an empty name", n.Name)
			}
		}
		if n.Loops < 0 || n.Loops > MaxLoops {
			return nil, fmt.Errorf("refused: node %q may be returned to %d times; the range "+
				"is 0 to %d", n.Name, n.Loops, MaxLoops)
		}
		// THE WORK IS RETURNED TO, NEVER THE VERDICT. A check scores the same
		// answer the same way every time, so returning to one re-scores what it
		// already has (LAW_003 §4); a gate is never inside a loop at all (§3).
		if n.Loops > 0 && (n.Kind == "eval" || n.Kind == "gate") {
			return nil, fmt.Errorf("refused: node %q is a %s and is not returned to -- a loop "+
				"re-does WORK, so `loops` goes on the node that does it", n.Name, n.Kind)
		}
		byName[n.Name] = n
	}
	for _, n := range s.Nodes {
		if n.Kind == "eval" {
			if _, ok := byName[n.Ref]; !ok {
				return nil, fmt.Errorf("refused: eval node %q checks unknown node %q", n.Name, n.Ref)
			}
		}
	}
	incoming := map[string]int{}
	adj := map[string][]Edge{}
	for name := range byName {
		incoming[name] = 0
	}
	returning, forward, err := loopsOf(s, byName)
	if err != nil {
		return nil, err
	}
	for _, e := range forward {
		switch e.When {
		case "", "always":
			e.When = "always"
		case "pass", "fail":
		default:
			return nil, fmt.Errorf("refused: edge %s->%s carries bad when %q", e.From, e.To, e.When)
		}
		if e.When == "fail" && byName[e.From].Kind != "eval" && byName[e.From].Kind != "gate" {
			return nil, fmt.Errorf("refused: fail-edges leave eval/gate nodes only (%s is %s)",
				e.From, byName[e.From].Kind)
		}
		incoming[e.To]++
		adj[e.From] = append(adj[e.From], e)
	}
	var starts []string
	for name, n := range incoming {
		if n == 0 {
			starts = append(starts, name)
		}
	}
	sort.Strings(starts)
	if len(starts) != 1 {
		return nil, fmt.Errorf("refused: want exactly one start, got %d", len(starts))
	}
	ready := append([]string{}, starts...)
	var order []string
	indeg := map[string]int{}
	for k, v := range incoming {
		indeg[k] = v
	}
	for len(ready) > 0 {
		n := ready[0]
		ready = ready[1:]
		order = append(order, n)
		out := append([]Edge{}, adj[n]...)
		sort.Slice(out, func(i, j int) bool { return out[i].To < out[j].To })
		for _, e := range out {
			indeg[e.To]--
			if indeg[e.To] == 0 {
				ready = append(ready, e.To)
			}
		}
		sort.Strings(ready)
	}
	if len(order) != len(byName) {
		return nil, fmt.Errorf("refused: cycle or unreachable node in flow %q -- a node is "+
			"returned to only by a check's fail-edge, and only when it declares `loops`", s.Name)
	}
	// A RETURN IS JUDGED AGAINST THE FORWARD ORDER: it must go BACK, around a
	// body that does work, with no gate inside; and a declared ceiling must be
	// reached by something.
	if err := lawfulReturns(s, byName, order, returning, forward); err != nil {
		return nil, err
	}
	if err := aiderGranted(s, byName, forward); err != nil {
		return nil, err
	}
	return order, nil
}

// aiderGranted refuses an `aider` node that can be reached without the hand
// having crossed a gate that grants AiderTool (2026-10-05, H17: "on a grant from
// the gate you click to open the line and nowhere else").
//
// THE NEAREST GATE WINS, on every path. A crossing carries from its gate to the
// NEXT gate or the end of the run (runFrom), so a gate that grants nothing
// between the granting one and the node ends the hand's reach, and a path that
// meets the start before any gate was never authorised at all. The walk is over
// the forward edges only: a check's return re-fires work inside a crossing and
// opens no path of its own. Judged at the save, never discovered when it runs.
func aiderGranted(s Spec, byName map[string]Node, forward []Edge) error {
	in := map[string][]string{}
	for _, e := range forward {
		in[e.To] = append(in[e.To], e.From)
	}
	memo := map[string]bool{}
	var covered func(name string) bool
	covered = func(name string) bool {
		if v, ok := memo[name]; ok {
			return v
		}
		ups := in[name]
		ok := len(ups) > 0 // the start reached with no gate between: never authorised
		for _, up := range ups {
			if byName[up].Kind == "gate" {
				granted := false
				for _, g := range byName[up].Grants {
					if g == AiderTool {
						granted = true
					}
				}
				if !granted {
					ok = false
				}
				continue
			}
			if !covered(up) {
				ok = false
			}
		}
		memo[name] = ok
		return ok
	}
	for _, n := range s.Nodes {
		if n.Kind == "aider" && !covered(n.Name) {
			return fmt.Errorf("refused: aider node %q can be reached without passing a gate that grants %q "+
				"-- Aider is the operator's own hand, so the gate he crosses must say so, and the nearest "+
				"gate before the node on every path is the one that counts", n.Name, AiderTool)
		}
	}
	return nil
}

// loopsOf splits a spec's edges into the ones that RETURN -- a check's
// fail-edge into a node that declares `loops` -- and the forward rest, keyed
// by the check that returns. ONE READING for Validate and the runner, so what
// is refused at the save and what is re-fired at run time cannot disagree.
// Every other edge is forward, and a forward edge that closes a cycle is still
// refused by Kahn's count: a node is returned to by a check's fail-edge, and
// only when it declares `loops`.
func loopsOf(s Spec, byName map[string]Node) (returning map[string]Edge, forward []Edge, err error) {
	returning = map[string]Edge{}
	for _, e := range s.Edges {
		from, okf := byName[e.From]
		to, okt := byName[e.To]
		if !okf {
			return nil, nil, fmt.Errorf("refused: edge from unknown node %q", e.From)
		}
		if !okt {
			return nil, nil, fmt.Errorf("refused: edge to unknown node %q", e.To)
		}
		if to.Loops > 0 && from.Kind == "eval" && e.When == "fail" {
			if _, twice := returning[e.From]; twice {
				return nil, nil, fmt.Errorf("refused: check %q returns twice; one return per check", e.From)
			}
			returning[e.From] = e
			continue
		}
		forward = append(forward, e)
	}
	return returning, forward, nil
}

// lawfulReturns judges every return against the forward order, and names the
// ceiling nobody reaches.
func lawfulReturns(s Spec, byName map[string]Node, order []string, returning map[string]Edge, forward []Edge) error {
	pos := map[string]int{}
	for i, n := range order {
		pos[n] = i
	}
	reached := map[string]bool{}
	checks := make([]string, 0, len(returning))
	for check := range returning {
		checks = append(checks, check)
	}
	sort.Strings(checks) // one refusal, the same one every time
	for _, check := range checks {
		e := returning[check]
		if pos[e.To] >= pos[check] {
			return fmt.Errorf("refused: edge %s->%s does not return: %s does not stand before %s",
				check, e.To, e.To, check)
		}
		body := loopBody(e.To, check, forward)
		work := false
		for _, name := range body {
			switch byName[name].Kind {
			case "gate":
				return fmt.Errorf("refused: gate %q stands inside the return from %s to %s; a "+
					"gate stands at the end of a loop, never inside it (LAW_003 §3)", name, check, e.To)
			case "eval":
			default:
				work = true
			}
		}
		if !work {
			return fmt.Errorf("refused: the return from %s to %s re-does no work -- it would score "+
				"the same answer again until the score agrees (LAW_003 §4)", check, e.To)
		}
		reached[e.To] = true
	}
	for _, n := range s.Nodes {
		if n.Loops > 0 && !reached[n.Name] {
			return fmt.Errorf("refused: node %q declares `loops` and nothing returns to it -- "+
				"a ceiling read by nothing", n.Name)
		}
	}
	return nil
}

// loopBody is the nodes on forward paths from `to` to `check`, both included:
// what a return re-fires.
func loopBody(to, check string, forward []Edge) []string {
	down := reach(to, forward, false)
	up := reach(check, forward, true)
	body := []string{to, check}
	seen := map[string]bool{to: true, check: true}
	for name := range down {
		if up[name] && !seen[name] {
			body = append(body, name)
			seen[name] = true
		}
	}
	sort.Strings(body)
	return body
}

// reach is every node reachable from `start` over the forward edges, walking
// them backwards when `up` is set; the start itself is not included.
func reach(start string, forward []Edge, up bool) map[string]bool {
	out := map[string]bool{}
	queue := []string{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, e := range forward {
			next := ""
			if !up && e.From == cur {
				next = e.To
			} else if up && e.To == cur {
				next = e.From
			}
			if next == "" || out[next] || next == start {
				continue
			}
			out[next] = true
			queue = append(queue, next)
		}
	}
	return out
}

// Save folds a new spec version; history kept whole.
func Save(home string, s Spec) (Spec, error) {
	if !NameRe.MatchString(s.Name) {
		return Spec{}, fmt.Errorf("refused: flow name %q breaks the name law", s.Name)
	}
	if _, err := Validate(s); err != nil {
		return Spec{}, err
	}
	if err := os.MkdirAll(flowsDir(home), 0o755); err != nil {
		return Spec{}, err
	}
	latest := 0
	if cur, err := Get(home, s.Name, 0); err == nil {
		latest = cur.Version
		old, _ := os.ReadFile(specPath(home, s.Name, 0))
		if err := os.WriteFile(specPath(home, s.Name, latest), old, 0o644); err != nil {
			return Spec{}, err
		}
	}
	s.Version = latest + 1
	if s.BudgetS <= 0 {
		s.BudgetS = 600
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return Spec{}, err
	}
	if err := os.WriteFile(specPath(home, s.Name, 0), append(b, '\n'), 0o644); err != nil {
		return Spec{}, err
	}
	return s, nil
}

// Get reads version v (0 = latest); absence denied honestly.
func Get(home, name string, v int) (Spec, error) {
	if !NameRe.MatchString(name) {
		return Spec{}, fmt.Errorf("refused: flow name %q breaks the name law", name)
	}
	if v > 0 {
		if cur, err := Get(home, name, 0); err == nil && cur.Version == v {
			return cur, nil
		}
		raw, err := os.ReadFile(specPath(home, name, v))
		if err != nil {
			return Spec{}, fmt.Errorf("no such flow version: %s v%d — try flow_list", name, v)
		}
		var s Spec
		if err := json.Unmarshal(raw, &s); err != nil {
			return Spec{}, fmt.Errorf("flow %s v%d is corrupt: %s", name, v, err)
		}
		return s, nil
	}
	raw, err := os.ReadFile(specPath(home, name, 0))
	if err != nil {
		return Spec{}, fmt.Errorf("no such flow: %s — try flow_list", name)
	}
	var s Spec
	if err := json.Unmarshal(raw, &s); err != nil {
		return Spec{}, fmt.Errorf("flow %s is corrupt: %s", name, err)
	}
	return s, nil
}

// historyRe names a folded version -- <name>.v<k>.json, what Save writes when
// it folds -- which List leaves to Get(name, k) and never lists as a flow of
// its own.
var historyRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}\.v[0-9]+$`)

// Unread is a .json under flows/ that List could not hand back as a flow, and
// why. NOTHING IS HIDDEN (operator, 2026-09-25: "the engine shouldn't hide a
// corrupt spec, either"). Until then List skipped a spec that would not parse
// without a word, so a corrupt flow was invisible until somebody fired it --
// the core's release gate found the skip the first day it read this folder,
// and the door itself said nothing.
type Unread struct {
	File string `json:"file"`
	Why  string `json:"why"`
}

// List names every flow with its latest version, and every .json under flows/
// it could NOT read as one -- corrupt, or named outside the name law, which no
// tool can reach. Folded versions (<name>.v<k>.json) are history and are
// neither: Get(name, k) reaches them.
func List(home string) ([]Spec, []Unread, error) {
	entries, err := os.ReadDir(flowsDir(home))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	var out []Spec
	var unread []Unread
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		base := strings.TrimSuffix(e.Name(), ".json")
		if historyRe.MatchString(base) {
			continue
		}
		if !NameRe.MatchString(base) {
			unread = append(unread, Unread{File: e.Name(),
				Why: fmt.Sprintf("name %q breaks the name law, so no tool can reach it", base)})
			continue
		}
		s, err := Get(home, base, 0)
		if err != nil {
			unread = append(unread, Unread{File: e.Name(), Why: err.Error()})
			continue
		}
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	sort.Slice(unread, func(i, j int) bool { return unread[i].File < unread[j].File })
	return out, unread, nil
}
