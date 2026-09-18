// The first strokes on the glass's store.
//
// WHY THIS FILE DID NOT EXIST UNTIL 2026-09-12. The `webapp` module is ~2,800
// lines over eight packages and had ONE test function in all of it, in
// handlers/ws_test.go. ADR-006 made the same measurement about the door -- the
// protocol layer and the tenant model were its two least-tested things, and
// both were given first strokes on 2026-09-11. The glass never had that pass,
// and the operator named it the one thing to fix before field testing.
//
// This package is where every trace, eval, agent and message actually lives.
// Nothing here was a bug; what was missing was any statement of what it
// PROMISES, so the blast radius of a change was discoverable only by suffering
// it.
//
// Hermetic: temp directories only, nothing read from the real estate.
package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestARoundTripSurvivesReopening(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	d, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	d.AddTrace(Trace{ID: "t1", AgentID: "a1", Tool: "run_python", Tenant: "research"})
	d.AddEval(Eval{ID: "e1", TraceID: "t1", Name: "verdict", Passed: true})
	d.UpsertAgent(Agent{ID: "a1", Office: "Router"})
	d.AddMessage(Message{ID: "m1", Channel: "ops", Content: "hi"})
	d.SetKey("mcp_url", "http://127.0.0.1:8090")
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}

	// A NEW HANDLE ON THE SAME DIRECTORY IS THE WHOLE PROMISE. Every write
	// above persists -- the trace as a ledger line since 2026-09-14, the rest
	// through Save(); if any of them did not, this is where it shows.
	again, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := again.GetTrace("t1"); got == nil || got.Tool != "run_python" {
		t.Fatalf("trace did not survive: %+v", got)
	}
	if got := again.GetTrace("t1"); got.Tenant != "research" {
		t.Fatalf("the tenant did not survive, and it is what the wall reads: %+v", got)
	}
	if e := again.GetEvals("t1"); len(e) != 1 || !e[0].Passed {
		t.Fatalf("eval did not survive: %+v", e)
	}
	if a := again.GetAgent("a1"); a == nil || a.Office != "Router" {
		t.Fatalf("agent did not survive: %+v", a)
	}
	if m := again.GetMessages("ops", 10); len(m) != 1 {
		t.Fatalf("message did not survive: %+v", m)
	}
	if again.GetKey("mcp_url") == "" {
		t.Fatal("the key did not survive, and the glass finds the door by it")
	}
	// AND A KEY NOBODY SET IS EMPTY, not a zero value that reads as a setting
	if again.GetKey("never_set") != "" {
		t.Fatal("an unset key must answer empty")
	}
}

func TestTheSaveIsAtomicAndLeavesNoScratch(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	d, _ := Open(dir)
	d.AddTrace(Trace{ID: "t1"})
	// A TRACE NO LONGER SAVES THE STORE (2026-09-14): it is one ledger line, and
	// the strokes at the foot of this file pin that. This stroke drove the save
	// through AddTrace until then; it is driven the way every other write
	// drives it now, and the guard is unchanged.
	d.SetKey("mcp_url", "http://127.0.0.1:8090")

	// tmp-then-rename, so a reader never sees a half-written store. The scratch
	// file must not be left behind either: `store.json.tmp` sitting in the
	// directory is how a later reader learns a write died midway.
	if _, err := os.Stat(filepath.Join(dir, "store.json")); err != nil {
		t.Fatalf("no store written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "store.json.tmp")); err == nil {
		t.Fatal("the scratch file was left behind")
	}
	// and it is real JSON a person can read, not an opaque blob
	raw, err := os.ReadFile(filepath.Join(dir, "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fd fileData
	if err := json.Unmarshal(raw, &fd); err != nil {
		t.Fatalf("the store is not readable JSON: %v", err)
	}
	// AND THE TRACE IS NOT IN IT. store.json holds what stays small; the ledger
	// holds what grows.
	if len(fd.Traces) != 0 || bytes.Contains(raw, []byte(`"traces"`)) {
		t.Fatalf("store.json carries traces again: %d of them", len(fd.Traces))
	}
}

func TestNewestFirstAndTheLimitHolds(t *testing.T) {
	d, _ := Open(filepath.Join(t.TempDir(), "store"))
	for _, id := range []string{"t1", "t2", "t3"} {
		d.AddTrace(Trace{ID: id, AgentID: "a1", CreatedAt: time.Now()})
	}
	d.AddTrace(Trace{ID: "other", AgentID: "a2"})

	all := d.GetTraces("", 0)
	if len(all) != 4 {
		t.Fatalf("want 4 traces, got %d", len(all))
	}
	// NEWEST FIRST. A trace log read oldest-first shows a stale top on every
	// page, which is the one thing the page is for.
	if all[0].ID != "other" {
		t.Fatalf("newest must lead: %s", all[0].ID)
	}
	if got := d.GetTraces("a1", 2); len(got) != 2 || got[0].ID != "t3" {
		t.Fatalf("limit + order broke: %d %+v", len(got), got)
	}
	// filtering by agent must not leak another agent's rows
	for _, tr := range d.GetTraces("a1", 0) {
		if tr.AgentID != "a1" {
			t.Fatalf("agent filter leaked %s", tr.AgentID)
		}
	}
	// a trace nobody wrote is absent, not a zero value
	if d.GetTrace("ghost") != nil {
		t.Fatal("a missing trace must be nil, never an empty Trace")
	}
}

func TestUpsertReplacesRatherThanDuplicates(t *testing.T) {
	d, _ := Open(filepath.Join(t.TempDir(), "store"))
	d.UpsertAgent(Agent{ID: "a1", Office: "Router", Role: "route"})
	d.UpsertAgent(Agent{ID: "a1", Office: "Router", Role: "READ ONLY"})
	d.UpsertAgent(Agent{ID: "a2", Office: "Steward"})

	if n := len(d.GetAgents()); n != 2 {
		t.Fatalf("upsert duplicated: %d agents", n)
	}
	if a := d.GetAgent("a1"); a.Role != "READ ONLY" {
		t.Fatalf("upsert did not replace: %+v", a)
	}
}

func TestGetAgentsHandsBackACopy(t *testing.T) {
	d, _ := Open(filepath.Join(t.TempDir(), "store"))
	d.UpsertAgent(Agent{ID: "a1", Office: "Router"})

	got := d.GetAgents()
	got[0].Office = "TAMPERED"

	// A caller that mutates the slice it was handed must not reach the store.
	// `GetAgents` copies on purpose; nothing said so until now.
	if a := d.GetAgent("a1"); a.Office != "Router" {
		t.Fatalf("a caller's write reached the store: %+v", a)
	}
}

// AND THE HONEST LIMIT, NAMED RATHER THAN FOUND LATER. `load()` ignores a
// store.json it cannot read: no error, no log, an empty database. That is a
// deliberate choice for a missing file and it is the SAME path for a corrupt
// one, so a damaged store comes up silently empty rather than refusing to
// start. This stroke does not call that right -- it pins it, so a change is a
// decision instead of a surprise.
func TestACorruptStoreComesUpEmptyAndSaysNothing(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "store.json"),
		[]byte("{not json at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	d, err := Open(dir)
	if err != nil {
		t.Fatalf("Open refused a corrupt store: %v", err)
	}
	if n := len(d.GetTraces("", 0)); n != 0 {
		t.Fatalf("want an empty store, got %d traces", n)
	}
	// and a write over it succeeds, which is how the corruption becomes
	// permanent -- the old bytes are gone after the first Save
	d.AddTrace(Trace{ID: "t1"})
	if d.GetTrace("t1") == nil {
		t.Fatal("a write after a corrupt load must still work")
	}
}

// ---- the ledger (2026-09-14) ----------------------------------------------
//
// The glass rewrote its whole store to record one trace: 29,156 traces in a
// 286 MB store.json, marshalled on every tool call, until the process held
// 17.3 GB and a tool call took 25 s. Each stroke below is one promise of the
// ledger that replaced it.

func ledgerLines(t *testing.T, dir string) []string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, "traces.jsonl"))
	if err != nil {
		t.Fatalf("no ledger: %v", err)
	}
	var out []string
	for _, l := range strings.Split(string(raw), "\n") {
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}

func TestATraceIsOneLineAndTheStoreIsNotRewritten(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	d, _ := Open(dir)
	d.SetKey("mcp_url", "http://127.0.0.1:8090")
	store := filepath.Join(dir, "store.json")
	before, err := os.ReadFile(store)
	if err != nil {
		t.Fatal(err)
	}
	// A REWRITE OF THE SAME BYTES IS STILL A REWRITE, and it was the cost. The
	// file is stamped with a time no write could produce, so any write shows.
	stamp := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(store, stamp, stamp); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 50; i++ {
		d.AddTrace(Trace{ID: fmt.Sprintf("t%d", i), Tool: "git", Output: strings.Repeat("x", 4096)})
	}
	after, _ := os.ReadFile(store)
	info, err := os.Stat(store)
	if err != nil {
		t.Fatal(err)
	}
	// THE WHOLE POINT: fifty traces, and store.json was not touched.
	if !info.ModTime().Equal(stamp) || !bytes.Equal(before, after) {
		t.Fatalf("recording a trace rewrote store.json (%d -> %d bytes, modified %s)",
			len(before), len(after), info.ModTime())
	}
	if n := len(ledgerLines(t, dir)); n != 50 {
		t.Fatalf("want one ledger line per trace, got %d", n)
	}
	if d.TraceCount() != 50 {
		t.Fatalf("count %d, want 50", d.TraceCount())
	}
}

func TestTheWindowIsBoundedAndTheLedgerKeepsEveryTrace(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	d, _ := Open(dir)
	n := TraceWindow + 25
	for i := 0; i < n; i++ {
		d.AddTrace(Trace{ID: fmt.Sprintf("t%d", i), Tool: fmt.Sprintf("tool%d", i)})
	}
	check := func(d *DB, when string) {
		t.Helper()
		all := d.GetTraces("", 0)
		if len(all) != TraceWindow {
			t.Fatalf("%s: memory holds %d traces, want the window of %d", when, len(all), TraceWindow)
		}
		if all[0].ID != fmt.Sprintf("t%d", n-1) {
			t.Fatalf("%s: newest must lead, got %s", when, all[0].ID)
		}
		if d.TraceCount() != n {
			t.Fatalf("%s: count %d, want every trace: %d", when, d.TraceCount(), n)
		}
		// OUT OF THE WINDOW IS NOT OUT OF THE RECORD. An eval can name any trace.
		if got := d.GetTrace("t0"); got == nil || got.Tool != "tool0" {
			t.Fatalf("%s: the oldest trace is not found by its id: %+v", when, got)
		}
	}
	check(d, "while open")
	again, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	check(again, "after reopening")
}

// writeOldStore writes a store.json from before the ledger, in the shape it
// was saved in: the traces inline beside everything else.
func writeOldStore(t *testing.T, dir string, ids ...string) {
	t.Helper()
	var list []Trace
	for _, id := range ids {
		list = append(list, Trace{ID: id, Tool: "tool-" + id, Output: "out-" + id})
	}
	raw, err := json.MarshalIndent(map[string]any{
		"traces":   list,
		"evals":    []Eval{{ID: "e1", TraceID: "t1", Name: "verdict", Passed: true}},
		"agents":   []Agent{},
		"messages": []Message{},
		"keys":     map[string]string{"mcp_url": "http://127.0.0.1:8090"},
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "store.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAStoreFromBeforeTheLedgerFoldsWithoutLosingATrace(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	writeOldStore(t, dir, "t1", "t2", "t3")
	d, err := Open(dir)
	if err != nil {
		t.Fatalf("a store from before the ledger was refused: %v", err)
	}
	lines := ledgerLines(t, dir)
	if len(lines) != 3 || !strings.Contains(lines[0], `"id":"t1"`) || !strings.Contains(lines[2], `"id":"t3"`) {
		t.Fatalf("the fold lost or reordered traces: %v", lines)
	}
	if got := d.GetTrace("t2"); got == nil || got.Output != "out-t2" {
		t.Fatalf("a folded trace lost its output: %+v", got)
	}
	// store.json keeps everything else, and no longer carries the traces
	raw, _ := os.ReadFile(filepath.Join(dir, "store.json"))
	if bytes.Contains(raw, []byte(`"traces"`)) {
		t.Fatal("store.json still carries the traces after the fold")
	}
	if d.GetKey("mcp_url") == "" || len(d.GetEvals("t1")) != 1 {
		t.Fatal("the fold dropped a key or an eval")
	}
	if _, err := os.Stat(filepath.Join(dir, "traces.jsonl.tmp")); err == nil {
		t.Fatal("the fold left its scratch ledger behind")
	}
	// and a second start folds nothing twice
	again, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if again.TraceCount() != 3 || len(ledgerLines(t, dir)) != 3 {
		t.Fatalf("reopening folded again: count %d", again.TraceCount())
	}
}

// A START THAT DIED BETWEEN THE RENAME AND THE REWRITE leaves the traces in
// both places, and a glass may have recorded more since. The next fold writes
// each trace once, in order.
func TestAFoldThatDiedMidwayWritesNoTraceTwice(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	writeOldStore(t, dir, "t1", "t2", "t3")
	var ledger []byte
	for _, id := range []string{"t1", "t2", "t3", "t4"} {
		line, _ := json.Marshal(Trace{ID: id, Tool: "tool-" + id})
		ledger = append(append(ledger, line...), '\n')
	}
	if err := os.WriteFile(filepath.Join(dir, "traces.jsonl"), ledger, 0o644); err != nil {
		t.Fatal(err)
	}
	d, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	lines := ledgerLines(t, dir)
	if len(lines) != 4 || d.TraceCount() != 4 {
		t.Fatalf("a trace was written twice: %d lines, count %d", len(lines), d.TraceCount())
	}
	for i, id := range []string{"t1", "t2", "t3", "t4"} {
		if !strings.Contains(lines[i], `"id":"`+id+`"`) {
			t.Fatalf("line %d is not %s: %s", i, id, lines[i])
		}
	}
}

// A KILLED PROCESS CAN LEAVE HALF A LINE. It is not counted, and the next trace
// does not glue itself onto it and vanish with it.
func TestATornLastLineIsSkippedAndTheNextTraceStandsAlone(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	good, _ := json.Marshal(Trace{ID: "t1", Tool: "git"})
	torn := append(append(good, '\n'), []byte(`{"id":"t2","tool":"gi`)...)
	if err := os.WriteFile(filepath.Join(dir, "traces.jsonl"), torn, 0o644); err != nil {
		t.Fatal(err)
	}
	d, _ := Open(dir)
	if d.TraceCount() != 1 || len(d.GetTraces("", 0)) != 1 {
		t.Fatalf("the torn line was counted: %d", d.TraceCount())
	}
	d.AddTrace(Trace{ID: "t3", Tool: "records"})
	again, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if again.TraceCount() != 2 {
		t.Fatalf("after a trace past the tear, count %d, want 2", again.TraceCount())
	}
	if got := again.GetTrace("t3"); got == nil || got.Tool != "records" {
		t.Fatalf("the trace after the tear was lost: %+v", got)
	}
}

// ---- the save (2026-09-15) --------------------------------------------------

// SAVES THAT RACE WRITES NEITHER END THE PROCESS NOR LOSE A KEY. Save marshalled
// the keys map after letting go of the lock, so a SetKey landing mid-marshal was
// a map write under a live iteration -- a fatal error, the whole glass down --
// and two Saves wrote one scratch file at once. Eight writers set 800 keys, each
// write saving; read back from disk, the store must hold all 800.
func TestSavesThatRaceWritesNeitherCrashNorLoseAKey(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	d, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	const writers, each = 8, 100
	var wg sync.WaitGroup
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < each; i++ {
				d.SetKey(fmt.Sprintf("k%d-%d", w, i), "v")
			}
		}(w)
	}
	wg.Wait()
	// FROM DISK, NOT MEMORY: the Save that renamed last is the store.
	again, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	missing := 0
	for w := 0; w < writers; w++ {
		for i := 0; i < each; i++ {
			if again.GetKey(fmt.Sprintf("k%d-%d", w, i)) == "" {
				missing++
			}
		}
	}
	if missing > 0 {
		t.Fatalf("%d of %d keys are not in store.json: the Save that renamed last wrote an older snapshot",
			missing, writers*each)
	}
}

// SAVES AT THE SAME MOMENT MUST EACH LAND. Every caller of Save discards its
// error, so a Save that failed because another was writing the same scratch
// file was a write nobody heard about. Sixteen Saves at once, twenty times
// each, over a store with some weight in it: every one must succeed, and the
// store must come back whole.
func TestSavesAtTheSameMomentEachLand(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "store")
	d, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	d.SetKey("mcp_url", "http://127.0.0.1:8090")
	for i := 0; i < 50; i++ {
		d.UpsertAgent(Agent{ID: fmt.Sprintf("a%d", i), USContent: strings.Repeat("x", 2048)})
	}
	const savers, each = 16, 20
	errs := make(chan error, savers*each)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for g := 0; g < savers; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for i := 0; i < each; i++ {
				if err := d.Save(); err != nil {
					errs <- err
				}
			}
		}()
	}
	close(start)
	wg.Wait()
	close(errs)
	var failed []string
	for err := range errs {
		failed = append(failed, err.Error())
	}
	if len(failed) > 0 {
		t.Fatalf("%d of %d Saves failed, and no caller would have seen it; the first: %s",
			len(failed), savers*each, failed[0])
	}
	again, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if again.GetKey("mcp_url") == "" || len(again.GetAgents()) != 50 {
		t.Fatalf("the store did not come back whole: key %q, %d agents",
			again.GetKey("mcp_url"), len(again.GetAgents()))
	}
}
