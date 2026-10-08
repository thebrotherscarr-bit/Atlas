package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"atlas/line/internal/tenant"
)

// THE OPERATOR'S EDIT, held by strokes (edit.go): it parks and writes only on his Approve, writes exactly the change his card showed, keeps the file's own line endings, and refuses what it must.

func editCall(t *testing.T, reg *Registry, tr *tenant.Registry, args map[string]any) fileEditAnswer {
	t.Helper()
	out, err := reg.Call(tr, "file_edit", args, glass)
	if err != nil {
		t.Fatalf("reg.Call failed: %v", err)
	}
	var answer fileEditAnswer
	if err := json.Unmarshal([]byte(out), &answer); err != nil {
		t.Fatalf("json.Unmarshal failed: %v, output: %s", err, out)
	}
	return answer
}

func editApprove(t *testing.T, reg *Registry, tr *tenant.Registry, id string) fileEditAnswer {
	t.Helper()
	out, err := reg.Call(tr, "hold_answer", map[string]any{"id": id, "decision": "approve"}, glass)
	if err != nil {
		t.Fatalf("reg.Call failed: %v", err)
	}
	i := strings.Index(out, "{")
	if i < 0 {
		t.Fatalf("invalid output: %s", out)
	}
	var answer fileEditAnswer
	if err := json.Unmarshal([]byte(out[i:]), &answer); err != nil {
		t.Fatalf("json.Unmarshal failed: %v, output: %s", err, out[i:])
	}
	return answer
}

// AN EDIT KEEPS THE FILE'S OWN LINE ENDINGS: the old and new text come with line feeds, and a CRLF file stays CRLF.
func TestAnEditKeepsTheFilesOwnLineEndings(t *testing.T) {
	reg, tr, tn := shellWorld(t, map[string]string{"dos.txt": "a\r\nb\r\n"})
	a := editCall(t, reg, tr, map[string]any{"filepath": "dos.txt", "old": "a\nb\n", "new": "A\nB\n"})
	if a.State != "held" {
		t.Fatalf("Unexpected state: %v", a)
	}
	w := editApprove(t, reg, tr, a.Hold)
	if w.State != "written" {
		t.Fatalf("Unexpected state after approval: %v", w)
	}
	if readText(t, tn.Home, "dos.txt") != "A\r\nB\r\n" {
		t.Fatalf("File content not updated after approval: %q", readText(t, tn.Home, "dos.txt"))
	}
}

// AN EDIT IS REFUSED WHERE IT MUST BE, and a refusal parks nothing and writes nothing.
func TestAnEditIsRefusedWhereItMustBe(t *testing.T) {
	reg, tr, tn := shellWorld(t, map[string]string{
		"notes.txt": "one\ntwo\nthree\n",
		"twice.txt": "x\nx\n",
		"mixed.txt": "a\r\nb\n",
		"sub/f.txt": "f\n",
	})

	cases := []map[string]any{
		{"filepath": ""},
		{"filepath": "law/x.md", "old": "a", "new": "b"},
		{"filepath": "notes.txt", "old": "", "new": "x"},
		{"filepath": "notes.txt", "old": "two", "new": "two"},
		{"filepath": "notes.txt", "old": "absent", "new": "x"},
		{"filepath": "twice.txt", "old": "x", "new": "y"},
		{"filepath": "mixed.txt", "old": "a", "new": "b"},
		{"filepath": "sub", "old": "a", "new": "b"},
		{"filepath": "missing.txt", "old": "a", "new": "b"},
	}

	for _, c := range cases {
		a := editCall(t, reg, tr, c)
		if a.State != "refused" || len(a.Why) == 0 {
			t.Errorf("Unexpected state or reason for case %v: %v", c, a)
		}
	}

	if len(heldIDs(t, reg, tr)) != 0 {
		t.Fatalf("Held IDs should be empty after all refusals")
	}

	if readText(t, tn.Home, "notes.txt") != "one\ntwo\nthree\n" {
		t.Fatalf("File content was changed after refusals")
	}

}

// AN APPROVE THAT NO LONGER MATCHES WRITES NOTHING: the file changed after his card was shown.
func TestAnEditThatNoLongerMatchesWritesNothing(t *testing.T) {
	reg, tr, tn := shellWorld(t, map[string]string{"notes.txt": "one\ntwo\nthree\n"})
	a := editCall(t, reg, tr, map[string]any{"filepath": "notes.txt", "old": "two", "new": "TWO"})
	if a.State != "held" {
		t.Fatalf("Unexpected state: %v", a)
	}
	if err := os.WriteFile(filepath.Join(tn.Home, "notes.txt"), []byte("zero\none\ntwo\nthree\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	w := editApprove(t, reg, tr, a.Hold)
	if w.State != "refused" {
		t.Fatalf("Unexpected state after approval: %v", w)
	}
	if readText(t, tn.Home, "notes.txt") != "zero\none\ntwo\nthree\n" {
		t.Fatalf("File content not updated after approval: %q", readText(t, tn.Home, "notes.txt"))
	}
}

// NO EDIT IS MADE WHILE A SITTING IS OPEN (SITTING LAW 5).
func TestNoEditIsMadeWhileASittingIsOpen(t *testing.T) {
	sitting := `{"n": 7, "started": "2026-10-07T09:00:00", "ended": ""}` + "\n"
	reg, tr, tn := shellWorld(t, map[string]string{
		"notes.txt":               "one\ntwo\nthree\n",
		"sessions/sessions.jsonl": sitting,
	})
	a := editCall(t, reg, tr, map[string]any{"filepath": "notes.txt", "old": "two", "new": "TWO"})
	if a.State != "refused" || !strings.Contains(strings.Join(a.Why, " "), "SITTING LAW 5") {
		t.Fatalf("Unexpected state or reason: %v", a)
	}
	if len(heldIDs(t, reg, tr)) != 0 {
		t.Fatalf("Held IDs should be empty after refusal")
	}
	if readText(t, tn.Home, "notes.txt") != "one\ntwo\nthree\n" {
		t.Fatalf("File content was changed after refusal: %q", readText(t, tn.Home, "notes.txt"))
	}
}

// FILE_EDIT IS HIS GLASS'S ALONE: no other hand reaches it, through the door or by calling the tool itself.
func TestFileEditIsHisGlassAlone(t *testing.T) {
	reg, tr, tn := shellWorld(t, map[string]string{"notes.txt": "one\ntwo\nthree\n"})
	_, err := reg.Call(tr, "file_edit", map[string]any{"filepath": "notes.txt", "old": "two", "new": "TWO"}, agent)
	if err == nil || !strings.Contains(err.Error(), "operator's own hand") {
		t.Fatalf("Unexpected error: %v", err)
	}
	tool, ok := reg.Get("file_edit")
	if !ok {
		t.Fatal("Tool not found")
	}
	out, err := tool.Fn(tn, map[string]any{"filepath": "notes.txt", "old": "two", "new": "TWO", CallerKey: agent})
	if err == nil || out != "" {
		t.Fatalf("Unexpected output: %v", out)
	}
	if len(heldIDs(t, reg, tr)) != 0 || readText(t, tn.Home, "notes.txt") != "one\ntwo\nthree\n" {
		t.Fatalf("Unexpected state or file content")
	}
}

func readText(t *testing.T, home, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(home, rel))
	if err != nil {
		t.Fatalf("os.ReadFile failed: %v", err)
	}
	return string(b)
}

// AN EDIT WAITS FOR HIS APPROVE, and then writes exactly the change his card showed.
func TestAnEditWaitsForHisApprovalAndWritesExactlyWhatWasShown(t *testing.T) {
	reg, tr, tn := shellWorld(t, map[string]string{"notes.txt": "one\ntwo\nthree\n"})
	a := editCall(t, reg, tr, map[string]any{"filepath": "notes.txt", "old": "two", "new": "TWO"})
	if a.State != "held" || a.Hold == "" {
		t.Fatalf("Unexpected state or hold: %v", a)
	}
	found := false
	for _, id := range heldIDs(t, reg, tr) {
		if id == a.Hold {
			found = true
		}
	}
	if !found {
		t.Fatalf("Hold ID not found in held IDs: %s", a.Hold)
	}
	if strings.Join(a.Change, "|") != "at line 2|- two|+ TWO" {
		t.Fatalf("Unexpected change: %v", a.Change)
	}
	if readText(t, tn.Home, "notes.txt") != "one\ntwo\nthree\n" {
		t.Fatalf("File was written before approval")
	}
	w := editApprove(t, reg, tr, a.Hold)
	if w.State != "written" {
		t.Fatalf("Unexpected state after approval: %v", w)
	}
	if readText(t, tn.Home, "notes.txt") != "one\nTWO\nthree\n" {
		t.Fatalf("File content not updated after approval")
	}
}
