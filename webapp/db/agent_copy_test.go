package db

import (
	"path/filepath"
	"testing"
)

// GetAgent HANDS BACK A COPY, NOT THE STORE'S OWN MEMORY (2026-09-29, the
// core's WHAT'S LEFT C23). It returned a pointer into d.agents and the caller
// read it with the lock gone, while UpsertAgent may overwrite that element or
// grow the slice under it. Two ways to see it without a race detector: a
// caller's write must not reach the store, and an upsert must not move under
// a reader's copy.
func TestGetAgentHandsBackACopyAndNotTheStoresOwnMemory(t *testing.T) {
	d, err := Open(filepath.Join(t.TempDir(), "store"))
	if err != nil {
		t.Fatal(err)
	}
	d.UpsertAgent(Agent{ID: "a1", Office: "Router"})

	p := d.GetAgent("a1")
	if p == nil {
		t.Fatal("premise: the agent is there")
	}
	p.Office = "changed by the caller"
	if again := d.GetAgent("a1"); again.Office != "Router" {
		t.Fatalf("a caller's write reached the store: GetAgent handed back its own memory (%q)", again.Office)
	}

	held := d.GetAgent("a1")
	d.UpsertAgent(Agent{ID: "a1", Office: "Steward"})
	if held.Office != "Router" {
		t.Fatalf("an upsert moved under a reader's copy: %q", held.Office)
	}
	if now := d.GetAgent("a1"); now.Office != "Steward" {
		t.Fatalf("the upsert itself did not land: %q", now.Office)
	}
	if d.GetAgent("nobody") != nil {
		t.Fatal("an unknown id must answer nil")
	}
}
