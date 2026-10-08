package tools

// THE OPERATOR'S EDIT (2026-10-07; his word: "go, build the edit tool in the door").
//
// One exact edit of one existing text file. The path, the exact old text and the new text travel
// as plain fields, never packed into a command. It never writes on the first call: it parks itself
// as shell_run does, with the change set out line by line for his card (Inspector, Guardrails,
// "Waiting for your hand"), and it writes only when hold_answer replays it on his Approve and the
// file still gives exactly the change he was shown. It keeps what the seats and Aider keep
// (aiderNeverWritten), refuses a secret, client material or a path outside the world (shellRefusals
// and the door's jail), and edits nothing while a sitting is open (SITTING LAW 5).

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"atlas/line/internal/engine"
	"atlas/line/internal/tenant"
)

// fileEditMay is why an edit may not touch rel now, or nil.
func fileEditMay(t tenant.Tenant, rel string) []string {
	if rel == "" || rel == "." {
		return []string{"name the file as filepath, relative to this world"}
	}
	if s := jailed(t, rel); s != "" {
		return []string{s}
	}
	if rs := shellRefusals(rel, t.Home); len(rs) != 0 {
		return rs
	}
	if s := aiderNeverWritten(rel); s != "" {
		return []string{s}
	}
	if e, open := engines.Get(t.Home); open {
		return []string{fmt.Sprintf("an engine is open on %q (sitting %s): nothing in this ground is edited while the operator's sitting is open (SITTING LAW 5): close it first", t.Name, e.Opened().Str("sitting"))}
	}
	if n, started, sat := engine.SittingOpen(t.Home); sat {
		return []string{fmt.Sprintf("%q has an open sitting (%d, opened %s): nothing in this ground is edited while the operator's sitting is open (SITTING LAW 5): close it first", t.Name, int(n), started)}
	}
	return nil
}

// fileEditLoad is the file at rel, read whole, or why it cannot be edited.
func fileEditLoad(t tenant.Tenant, rel string) (abs string, st os.FileInfo, b []byte, why []string) {
	abs = filepath.Join(t.Home, filepath.FromSlash(rel))
	var err error
	st, err = os.Stat(abs)
	if err != nil {
		return "", nil, nil, []string{"there is no such file: file_edit changes a file that is already there"}
	}
	if st.IsDir() {
		return "", nil, nil, []string{"that is a folder, not a file"}
	}
	if st.Size() > fileEditMaxBytes {
		return "", nil, nil, []string{fmt.Sprintf("the file is %d bytes, over the %d an edit takes", st.Size(), fileEditMaxBytes)}
	}
	b, err = os.ReadFile(abs)
	if err != nil {
		return "", nil, nil, []string{err.Error()}
	}
	if !utf8.Valid(b) {
		return "", nil, nil, []string{"the file is not UTF-8 text"}
	}
	return abs, st, b, nil
}

// fileEditPark parks the call for his card with the change it showed him: the hold's id, or why not.
func fileEditPark(t tenant.Tenant, args map[string]any, c Caller, change []string) (string, string) {
	reg := regOf(args)
	if reg == nil {
		return "", "the hold queue is not reachable from here, so this cannot wait for his hand"
	}
	reg.mu.Lock()
	armed := reg.holdWrites
	reg.mu.Unlock()
	if !armed {
		return "", "the door was started without --auth, so no card can be shown: nothing was written"
	}
	self, found := reg.Get("file_edit")
	if !found {
		return "", "the hold queue is not reachable from here, so this cannot wait for his hand"
	}
	keep := copyArgs(args)
	delete(keep, registryKey)
	keep["change"] = change
	return reg.park(t, self, keep, c).ID, ""
}

// fileEditPlan is one edit, checked against the file as it is now: where, the file's text and the change.
type fileEditPlan struct {
	rel, abs, eol, lf, oldLF, newLF string
	st                              os.FileInfo
	change                          []string
}

// fileEditPlanFor checks one call's edit, or says why it cannot be made.
func fileEditPlanFor(t tenant.Tenant, args map[string]any) (fileEditPlan, []string) {
	p := fileEditPlan{rel: filepath.ToSlash(filepath.Clean(strings.TrimSpace(str(args, "filepath"))))}
	old := str(args, "old")
	nw := str(args, "new")
	if old == "" {
		return p, []string{"the old text is empty: name the exact passage to replace"}
	}
	if old == nw {
		return p, []string{"the old and the new text are the same: there is nothing to change"}
	}
	why := fileEditMay(t, p.rel)
	if why != nil {
		return p, why
	}
	abs, st, b, why := fileEditLoad(t, p.rel)
	if why != nil {
		return p, why
	}
	eol, ok := aiderEol(b)
	if !ok {
		return p, []string{"the file mixes two line endings, so an edit cannot keep what it has"}
	}
	p.abs, p.st, p.eol = abs, st, eol
	p.lf, p.oldLF, p.newLF = aiderLF(b), aiderLF([]byte(old)), aiderLF([]byte(nw))
	n := strings.Count(p.lf, p.oldLF)
	if n == 0 {
		return p, []string{"the old text is not in the file exactly as given: nothing was written"}
	}
	if n > 1 {
		return p, []string{fmt.Sprintf("the old text appears %d times: give a longer passage that appears once", n)}
	}
	p.change = fileEditChange(p.lf, p.oldLF, p.newLF)
	return p, nil
}

// toolFileEdit is file_edit: the first call parks; his approved replay writes.
func toolFileEdit(t tenant.Tenant, args map[string]any) (string, error) {
	c := callerOf(args)
	if !c.Service {
		return "", errors.New("file_edit is the operator's own hand: his glass alone may call it")
	}
	p, why := fileEditPlanFor(t, args)
	ans := fileEditAnswer{Filepath: p.rel, Change: p.change}

	refuse := func(why ...string) (string, error) {
		ans.State = "refused"
		ans.Why = why
		return ans.String(), nil
	}

	if len(why) > 0 {
		return refuse(why...)
	}

	if c.Hold == "" {
		id, whyNot := fileEditPark(t, args, c, p.change)
		if whyNot != "" {
			return refuse(whyNot)
		}
		ans.State = "held"
		ans.Hold = id
		ans.Note = "Nothing has been written. It waits for his Approve on the card in Guardrails."
		return ans.String(), nil
	}

	if !fileEditSame(args["change"], p.change) {
		return refuse("the file has changed since the card was shown: nothing was written. Ask again")
	}

	err := os.WriteFile(p.abs, []byte(aiderApplyEol(strings.Replace(p.lf, p.oldLF, p.newLF, 1), p.eol)), p.st.Mode().Perm())
	if err != nil {
		return refuse(err.Error())
	}

	ans.State = "written"
	ans.Note = fmt.Sprintf("Wrote %s, keeping its %s line endings. The suites are his to run.", p.rel, strings.ToUpper(aiderEolName(p.eol)))
	return ans.String(), nil
}

// fileEditMaxBytes is the largest file an edit takes.
const fileEditMaxBytes = 4 << 20

// fileEditAnswer is the one shape file_edit answers in.
type fileEditAnswer struct {
	State    string   `json:"state"`
	Filepath string   `json:"filepath,omitempty"`
	Hold     string   `json:"hold,omitempty"`
	Why      []string `json:"why,omitempty"`
	Change   []string `json:"change,omitempty"`
	Note     string   `json:"note,omitempty"`
}

// fileEditChange is the change as his card shows it, one line per entry.
func fileEditChange(lf, old, nw string) []string {
	i := strings.Index(lf, old)
	n := strings.Count(lf[:i], "\n") + 1
	change := []string{fmt.Sprintf("at line %d", n)}
	for _, line := range strings.Split(old, "\n") {
		change = append(change, "- "+line)
	}
	if nw == "" {
		change = append(change, "+ (nothing: the old text is taken out)")
	} else {
		for _, line := range strings.Split(nw, "\n") {
			change = append(change, "+ "+line)
		}
	}
	return change
}

// fileEditSame is whether the change parked with the card is the change the file gives now.
func fileEditSame(parked any, now []string) bool {
	var got []string
	switch v := parked.(type) {
	case []string:
		got = v
	case []any:
		for _, elem := range v {
			str, ok := elem.(string)
			if !ok {
				return false
			}
			got = append(got, str)
		}
	default:
		return false
	}
	if len(got) != len(now) {
		return false
	}
	for i := range got {
		if got[i] != now[i] {
			return false
		}
	}
	return true
}

// String is the answer as the glass reads it.
func (a fileEditAnswer) String() string {
	b, _ := json.MarshalIndent(a, "", " ")
	return string(b)
}

// fileEditTool is file_edit's declaration; tools.go registers it with one line, r.add(fileEditTool()).
func fileEditTool() Tool {
	return Tool{
		Name:        "file_edit",
		Writes:      true,
		ServiceOnly: true,
		Args:        []string{"filepath", "old", "new", "project?"},
		Fn:          toolFileEdit,
		Description: "the operator's own exact edit of one text file in this world: the path, the exact old text and the new text are plain fields; the old text must appear exactly once; the file keeps its own line endings; the first call writes nothing and parks itself with the change line by line, and it writes only when he approves that card and the file still gives exactly that change. Refused for a path outside the world, a secret, client material, what a model never writes, and while a sitting is open. His glass alone may call it",
	}
}
