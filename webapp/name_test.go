package main

import (
	"strings"
	"testing"
)

// The glass is Aurora (the operator's word, 2026-11-04: "Atlas PWA says Aurora on the name"). Its app window and
// any installed shortcut take their name from the page's <title> (there is no web manifest), and the top bar
// shows the brand; both are built into the binary (staticFiles), so this reads what the binary serves.
func TestTheGlassIsNamedAurora(t *testing.T) {
	page, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(page), "<title>Aurora") {
		t.Fatalf("the page's title does not name Aurora")
	}
	bar, err := staticFiles.ReadFile("static/js/agent.js")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(bar), "AURORA") || strings.Contains(string(bar), ">ATLAS<") {
		t.Fatalf("the top bar's brand is not AURORA")
	}
}
