package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

// THE BANNER READS THE COVENANT OFF THE RECORD (2026-09-29, the core's WHAT'S
// LEFT C29). It carried the house covenant as a literal; the ground this
// command is pointed at declares it (agents/operator.us), and a ground that
// declares none is said so.
func TestTheBannerReadsTheCovenantOffTheRecord(t *testing.T) {
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	us := "# operator\n\n```json\n" +
		`{"can_approve": false, "covenant": "feedfacefeedface", "id": "operator", "office": "OPERATOR"}` +
		"\n```\n"
	if err := os.WriteFile(filepath.Join(home, "agents", "operator.us"), []byte(us), 0o644); err != nil {
		t.Fatal(err)
	}
	lines := banner(home)
	if !strings.Contains(strings.Join(lines, "\n"), "covenant: feedfacefeedface") {
		t.Fatalf("the banner does not carry the covenant the ground declares:\n%s", strings.Join(lines, "\n"))
	}
	width := utf8.RuneCountInString(lines[0])
	for _, l := range lines {
		if utf8.RuneCountInString(l) != width {
			t.Fatalf("the banner's lines are not one width (%d wanted): %q", width, l)
		}
	}
	bare := strings.Join(banner(t.TempDir()), "\n")
	if !strings.Contains(bare, "none declared") || strings.Contains(bare, "1512741580b7239b") {
		t.Fatalf("a ground with no declaration must say so, never a number from memory:\n%s", bare)
	}
	src, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(src), "1512741580b7239b") {
		t.Fatal("main.go carries the covenant as a literal again")
	}
}
