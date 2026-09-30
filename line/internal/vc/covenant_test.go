package vc

import (
	"os"
	"path/filepath"
	"testing"
)

// THE COVENANT A GROUND DECLARES IS READ OFF THE OPERATOR'S OWN DECLARATION
// (2026-09-29): agents/operator.us, the record that says who holds the gate.
// A ground that carries none answers "", never a number from memory.
func TestTheCovenantIsReadOffTheOperatorsOwnDeclaration(t *testing.T) {
	home := t.TempDir()
	if got := CovenantOf(home); got != "" {
		t.Fatalf("a ground with no declaration answered %q", got)
	}
	if err := os.MkdirAll(filepath.Join(home, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	us := "# operator - declared seat\n\nthe hand\n\n```json\n" +
		`{"can_approve": false, "covenant": "feedfacefeedface", "id": "operator", "kind": "agent", "office": "OPERATOR"}` +
		"\n```\n"
	if err := os.WriteFile(filepath.Join(home, "agents", "operator.us"), []byte(us), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := CovenantOf(home); got != "feedfacefeedface" {
		t.Fatalf("read %q off the declaration, wanted feedfacefeedface", got)
	}
	// A declaration that carries none is "", the same as no declaration.
	bare := "# operator\n\n```json\n{\"can_approve\": false, \"id\": \"operator\"}\n```\n"
	if err := os.WriteFile(filepath.Join(home, "agents", "operator.us"), []byte(bare), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := CovenantOf(home); got != "" {
		t.Fatalf("a declaration with no covenant answered %q", got)
	}
}
