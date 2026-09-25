package protocol

import "testing"

// P0-12 (2026-09-25). An argument declared `name?` is optional; `name` is
// required; and the rule is read by both schema builders, so `project?` no
// longer reaches a client as required. This was the package's first test.
func TestAnOptionalArgumentIsOptional(t *testing.T) {
	for _, c := range []struct {
		arg  string
		want bool
	}{{"project?", true}, {"path", false}, {"?", false}, {"x??", true}, {"", false}} {
		if got := Optional(c.arg); got != c.want {
			t.Errorf("Optional(%q) = %v, want %v", c.arg, got, c.want)
		}
	}
	if TrimOptional("project?") != "project" || TrimOptional("path") != "path" {
		t.Fatal("TrimOptional must strip only the marker")
	}
}
