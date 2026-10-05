package tools

// Strokes for the laws through the glass (2026-10-05, law.go). law.py is the core's, not this repository's, so these strokes stand a stub in front of it
// (lawRun) that speaks as the real tool does -- the same verdict words, the same anchor on every link it lays -- and measure the DOOR'S OWN judgement: what it
// reads off the chain, what it refuses, what it appends and what it never touches. The real tool is held against this file's reading of it by the core's own
// suite (a stroke that runs law.py's own functions and reads this file's patterns off its source).

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"atlas/line/internal/tenant"
)

const (
	lawFixLedgerSealed = "# THE LAW LEDGER\n\n## 1. The first law\n\nOne, sealed.\n"
	lawFixLedgerDraft  = "\n---\n\n## 2. A drafted law\n\nTwo, still a draft.\n"
)

// lawLine is one chain line as law.py writes it: a link whose doc opens with the canonical anchor.
func lawLine(n int, kind, by, file, content string, extent int) string {
	body := []byte(content)
	words := ""
	if extent > 0 {
		body = body[:extent]
		words = fmt.Sprintf(" bytes:%d", extent)
	}
	sum := sha256.Sum256(body)
	doc := fmt.Sprintf("%s by %s -> law/%s sha256:%s%s\nthe note", kind, by, file, hex.EncodeToString(sum[:]), words)
	row := map[string]any{"ts": "2026-10-05T00:00:00+0000", "kind": "link", "prev": strings.Repeat("0", 64),
		"actor": by, "hash": fmt.Sprintf("%016x%048x", n, 0),
		"payload": map[string]any{"n": n, "says": kind + ":" + by, "mode": "open", "doc": doc, "cites": []string{}}}
	b, _ := json.Marshal(row)
	return string(b)
}

// lawStub is the stand-in for law.py: verify walks chain.jsonl the way law.py does (a whole link binds every byte, a prefix link its first N), and seal lays
// a prefix link over the file's length now. It records every call.
type lawStub struct {
	t           *testing.T
	home        string
	calls       [][]string
	verifies    int
	breakAfter  int  // from this verify on, the chain refuses (0 = never)
	silent      bool // verify exits clean and says nothing: no walk was proved
	sealFails   bool
	sealedAfter int
}

func (s *lawStub) run(py string, args []string, opts spawnOpts) spawnResult {
	s.calls = append(s.calls, append([]string{py}, args...))
	tn := tenant.Tenant{Name: "probe", Home: s.home}
	if len(args) < 2 {
		return spawnResult{Combined: "usage", Err: errors.New("exit status 2")}
	}
	switch args[1] {
	case "verify":
		s.verifies++
		if s.silent {
			return spawnResult{}
		}
		links, head, _ := lawLinks(tn)
		bad := []string{}
		count := 0
		for file, refs := range links {
			for _, r := range refs {
				count++
				if got, ok := lawSHA(filepath.Join(s.home, "law", file), r.extent); !ok || got != r.token {
					bad = append(bad, fmt.Sprintf("  - link #%d fingerprint MISMATCH: %s", r.n, file))
				}
			}
		}
		if s.breakAfter > 0 && s.verifies >= s.breakAfter {
			bad = append(bad, "  - link #9 does not open with the canonical anchor line")
		}
		if len(bad) > 0 {
			return spawnResult{Combined: "THE CHAIN REFUSES:\n" + strings.Join(bad, "\n") + "\n", Err: errors.New("exit status 1")}
		}
		return spawnResult{Combined: fmt.Sprintf("the law chain proves whole: %d links, head %s\n", count, head[:16])}
	case "seal":
		if s.sealFails {
			return spawnResult{Combined: "law/LAW_LEDGER.md shrank while it was being sealed\n", Err: errors.New("exit status 1")}
		}
		name := args[2]
		p := filepath.Join(s.home, "law", name)
		b, _ := os.ReadFile(p)
		links, _, _ := lawLinks(tn)
		n := 1
		for _, refs := range links {
			n += len(refs)
		}
		f, _ := os.OpenFile(filepath.Join(s.home, "law", lawChain), os.O_APPEND|os.O_WRONLY, 0)
		defer f.Close()
		f.WriteString(lawLine(n, "DIRECT", "operator", name, string(b), len(b)) + "\n")
		s.sealedAfter = len(b)
		return spawnResult{Combined: fmt.Sprintf("law/%s sealed to byte %d: link #%d  head %016x\n", name, len(b), n, n)}
	}
	return spawnResult{Combined: "usage", Err: errors.New("exit status 2")}
}

// lawWorld is a ground with a law folder: two laws bound whole, the ledger sealed to its first entry (and, if asked, a draft below it), and a stand-in law.py.
func lawWorld(t *testing.T, draft bool) (tenant.Tenant, *tenant.Registry, *Registry, *lawStub) {
	t.Helper()
	ledger := lawFixLedgerSealed
	if draft {
		ledger += lawFixLedgerDraft
	}
	estate, sitting := "# the ten\nfold, never delete\n", "# sitting laws\n1. read it\n"
	chain := strings.Join([]string{
		lawLine(1, "DIRECT", "operator", "ESTATE_LAWS.md", estate, 0),
		lawLine(2, "DIRECT", "operator", "SITTING_LAWS.md", sitting, 0),
		lawLine(3, "DIRECT", "operator", lawLedger, ledger, len(lawFixLedgerSealed)),
	}, "\n") + "\n"
	tn := world(t, map[string]string{
		"law/law.py": "# stand-in\n", "law/ESTATE_LAWS.md": estate, "law/SITTING_LAWS.md": sitting,
		"law/" + lawLedger: ledger, "law/" + lawChain: chain,
	})
	tr := tenant.NewRegistry()
	if err := tr.Add("probe", tn.Home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("probe"); err != nil {
		t.Fatal(err)
	}
	stub := &lawStub{t: t, home: tn.Home}
	was := lawRun
	lawRun = stub.run
	t.Cleanup(func() { lawRun = was })
	return tn, tr, Build(tr, Options{HoldWrites: true}), stub
}

func lawCall(t *testing.T, reg *Registry, tr *tenant.Registry, tool string, args map[string]any, service bool) (map[string]any, error) {
	t.Helper()
	out, err := reg.Call(tr, tool, args, Caller{Name: "glass", Service: service})
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if jerr := json.Unmarshal([]byte(out), &m); jerr != nil {
		t.Fatalf("%s did not answer JSON: %v\n%s", tool, jerr, out)
	}
	return m, nil
}

func lawBytes(t *testing.T, tn tenant.Tenant, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(tn.Home, "law", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func holdsText(t *testing.T, tn tenant.Tenant) string {
	t.Helper()
	b, _ := os.ReadFile(filepath.Join(tn.Home, "state", "holds.jsonl"))
	return string(b)
}

func TestLawStatusSaysHowFarEachLawIsSealedAndLetsTheChainSpeakForItself(t *testing.T) {
	tn, tr, reg, _ := lawWorld(t, true)
	if err := os.WriteFile(filepath.Join(tn.Home, "law", "NEW_LAW.md"), []byte("# not bound\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	st, err := lawCall(t, reg, tr, "law_status", map[string]any{}, false)
	if err != nil {
		t.Fatalf("law_status is a reader and was refused: %v", err)
	}
	if st["ok"] != true || !strings.Contains(fmt.Sprint(st["verdict"]), "the law chain proves whole: 3 links") {
		t.Fatalf("the chain's own verdict is not what law_status shows: %v", st)
	}
	if st["links"] != float64(3) || st["head"] != "0000000000000003" {
		t.Fatalf("the links and the head are not read off the walk: %v %v", st["links"], st["head"])
	}
	files := map[string]map[string]any{}
	for _, f := range st["files"].([]any) {
		m := f.(map[string]any)
		files[m["name"].(string)] = m
	}
	if len(files) != 4 {
		t.Fatalf("law_status lists %d files, want the four in law/ (law.py is not a law): %v", len(files), files)
	}
	for name, want := range map[string]string{"ESTATE_LAWS.md": "sealed", "SITTING_LAWS.md": "sealed",
		lawLedger: "partly sealed", "NEW_LAW.md": "not on the chain"} {
		if got := files[name]["state"]; got != want {
			t.Fatalf("%s is %v, want %s: %v", name, got, want, files[name])
		}
	}
	led := files[lawLedger]
	if led["sealed_bytes"] != float64(len(lawFixLedgerSealed)) || led["draft_bytes"] != float64(len(lawFixLedgerDraft)) {
		t.Fatalf("the ledger's sealed and draft bytes are not its own: %v", led)
	}
	if st["draft_bytes"] != float64(len(lawFixLedgerDraft)) {
		t.Fatalf("the total of draft bytes is not the ledger's tail: %v", st["draft_bytes"])
	}
	if led["kind"] != "DIRECT" || led["by"] != "operator" {
		t.Fatalf("a law does not say who laid its link: %v", led)
	}
	if len(files["ESTATE_LAWS.md"]["fingerprint"].(string)) != 16 {
		t.Fatalf("a fingerprint is the first sixteen hex of the sha256, as law.py status prints: %v", files["ESTATE_LAWS.md"])
	}
	if st["sitting_open"] != false || st["why"] != nil {
		t.Fatalf("no sitting is open and the writers would run: %v", st)
	}

	// A LAW CHANGED UNDER ITS SEAL IS SAID, and the walk's refusal is shown beside it.
	if err := os.WriteFile(filepath.Join(tn.Home, "law", "SITTING_LAWS.md"), []byte("# sitting laws\n1. skip it\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	st, _ = lawCall(t, reg, tr, "law_status", map[string]any{}, false)
	if st["ok"] != false || !strings.Contains(fmt.Sprint(st["verdict"]), "THE CHAIN REFUSES") || !strings.Contains(fmt.Sprint(st["verdict"]), "SITTING_LAWS.md") {
		t.Fatalf("a changed law did not turn the verdict red and name the file: %v", st)
	}
	for _, f := range st["files"].([]any) {
		if m := f.(map[string]any); m["name"] == "SITTING_LAWS.md" && m["state"] != "changed since sealed" {
			t.Fatalf("a law changed under its seal is %v, want 'changed since sealed'", m["state"])
		}
	}
}

func TestTheTwoWritersAreHisAloneAndTheReaderIsAnyones(t *testing.T) {
	tn, tr, reg, _ := lawWorld(t, true)
	for name, writes := range map[string]bool{"law_status": false, "law_add": true, "law_seal": true} {
		tool, ok := reg.Get(name)
		if !ok {
			t.Fatalf("the door does not carry %s", name)
		}
		if tool.Writes != writes || tool.ServiceOnly != writes {
			t.Fatalf("%s declares Writes=%v ServiceOnly=%v, want %v both: the laws are his hand's, never a seat's", name, tool.Writes, tool.ServiceOnly, writes)
		}
	}
	before := lawBytes(t, tn, lawLedger)
	for _, tool := range []string{"law_add", "law_seal"} {
		_, err := lawCall(t, reg, tr, tool, map[string]any{"title": "Smuggled", "text": strings.Repeat("a law of no standing ", 3)}, false)
		if err == nil || !strings.Contains(err.Error(), "operator's own hand") {
			t.Fatalf("%s answered a caller that is not his glass: %v", tool, err)
		}
	}
	// EVEN UNDER A HAND A CROSSED GATE RAISED: a grant never reaches a ServiceOnly tool for an agent.
	hands.raise(tn.Home, hand{Run: "r", Gate: "g", Tools: []string{"law_add", "law_seal"}})
	t.Cleanup(func() { hands.lower(tn.Home) })
	for _, tool := range []string{"law_add", "law_seal"} {
		_, err := reg.Call(tr, tool, map[string]any{"title": "Granted?", "text": strings.Repeat("a law of no standing ", 3)}, Caller{Name: "council"})
		if err == nil || !strings.Contains(err.Error(), "operator's own hand") {
			t.Fatalf("%s answered an agent standing under a crossed gate's grant: %v", tool, err)
		}
	}
	if string(lawBytes(t, tn, lawLedger)) != string(before) {
		t.Fatal("the ledger changed while every caller was refused")
	}
	if !strings.Contains(holdsText(t, tn), `"what":"refused"`) {
		t.Fatal("a refused caller left nothing in the holds record")
	}
}

func TestLawAddAppendsOneNumberedEntryBelowTheSealAndTouchesNothingAbove(t *testing.T) {
	tn, tr, reg, stub := lawWorld(t, true)
	before := lawBytes(t, tn, lawLedger)
	sealed := before[:len(lawFixLedgerSealed)]
	text := "9. **Nothing here is edited without his word.** A second paragraph\n   that wraps.\n\n### a sub heading is fine\nand so is this."
	got, err := lawCall(t, reg, tr, "law_add", map[string]any{"title": "SITTING LAW 9 -- a test law", "text": text,
		"from": "his words in chat, 2026-10-05\nand a second line of where it came from"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if got["state"] != "added" || got["entry"] != float64(3) {
		t.Fatalf("the entry was not numbered after the ledger's last (2): %v", got)
	}
	after := lawBytes(t, tn, lawLedger)
	if string(after[:len(before)]) != string(before) {
		t.Fatal("law_add changed bytes that were already in the ledger: it must only append")
	}
	if sha256.Sum256(after[:len(sealed)]) != sha256.Sum256(sealed) {
		t.Fatal("law_add touched the sealed prefix")
	}
	entry := string(after[len(before):])
	for _, want := range []string{"\n---\n\n## 3. SITTING LAW 9 -- a test law\n\n```\nentered:  ", ", on the operator's word\nfrom:     his words in chat, 2026-10-05\n          and a second line of where it came from\n```\n\n",
		"9. **Nothing here is edited without his word.** A second paragraph\n   that wraps.\n\n### a sub heading is fine\nand so is this.\n"} {
		if !strings.Contains(entry, want) {
			t.Fatalf("the entry is not in the ledger's own shape; missing %q in:\n%s", want, entry)
		}
	}
	if strings.Contains(string(after), "\r") || !strings.HasSuffix(string(after), "\n") {
		t.Fatal("the entry carries a carriage return or does not end on a line break")
	}
	st, _ := lawCall(t, reg, tr, "law_status", map[string]any{}, false)
	if st["draft_bytes"] != float64(len(after)-len(lawFixLedgerSealed)) || st["ok"] != true {
		t.Fatalf("the draft below the seal is not what was appended: %v", st)
	}
	if !strings.Contains(holdsText(t, tn), `"what":"law_added"`) || !strings.Contains(holdsText(t, tn), "entry 3") {
		t.Fatal("setting a law left nothing in the holds record, under his own hand")
	}
	if stub.verifies < 2 {
		t.Fatalf("the chain was walked %d times; it is walked before the append and again after it", stub.verifies)
	}
	// and the default of where it came from, and the numbering from nothing but the ledger itself
	got, _ = lawCall(t, reg, tr, "law_add", map[string]any{"title": "The next one", "text": "Another law, long enough to be one."}, true)
	if got["entry"] != float64(4) || !strings.Contains(string(lawBytes(t, tn, lawLedger)), "from:     his word, through the glass's Laws page") {
		t.Fatalf("the fourth entry or the default 'from' is wrong: %v", got)
	}
}

func TestLawAddRefusesWhatWouldBreakTheLedgerAndLeavesItByteForByte(t *testing.T) {
	long := strings.Repeat("x", lawMaxText+1)
	cases := []struct {
		why, find string
		args      map[string]any
	}{
		{"a title of two characters", "title is one line", map[string]any{"title": "ab", "text": strings.Repeat("a law ", 5)}},
		{"a title over a line", "title is one line", map[string]any{"title": "one\ntwo", "text": strings.Repeat("a law ", 5)}},
		{"a law too short to be one", "characters", map[string]any{"title": "Short", "text": "no"}},
		{"a law longer than the ledger takes", "characters", map[string]any{"title": "Long", "text": long}},
		{"a line that opens a new entry", "split the entry", map[string]any{"title": "Heading", "text": "a law that\n## 9. pretends to be an entry\nand goes on"}},
		{"a bare rule that ends an entry", "split the entry", map[string]any{"title": "Rule", "text": "a law that\n---\nand goes on, so it is two"}},
		{"a control character", "control character", map[string]any{"title": "Control", "text": "a law with a bell \x07 in it, long enough"}},
	}
	for _, c := range cases {
		tn, tr, reg, _ := lawWorld(t, true)
		before := lawBytes(t, tn, lawLedger)
		got, err := lawCall(t, reg, tr, "law_add", c.args, true)
		if err != nil || got["state"] != "refused" || !strings.Contains(fmt.Sprint(got["why"]), c.find) {
			t.Fatalf("%s: not refused in the right words: %v %v", c.why, got, err)
		}
		if string(lawBytes(t, tn, lawLedger)) != string(before) {
			t.Fatalf("%s: the ledger moved on a refusal", c.why)
		}
		if !strings.Contains(holdsText(t, tn), `"what":"law_refused"`) {
			t.Fatalf("%s: a refusal left no record", c.why)
		}
	}
	good := map[string]any{"title": "A fine law", "text": "A law of the right length to be set down here."}

	// A SITTING IS OPEN: SITTING LAW 5, read off the world's own ledger.
	tn, tr, reg, _ := lawWorld(t, true)
	before := lawBytes(t, tn, lawLedger)
	if err := os.MkdirAll(filepath.Join(tn.Home, "sessions"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tn.Home, "sessions", "sessions.jsonl"), []byte(`{"n":9,"started":"2026-10-05T09:00:00","ended":""}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	st, _ := lawCall(t, reg, tr, "law_status", map[string]any{}, false)
	if st["sitting_open"] != true || !strings.Contains(fmt.Sprint(st["why"]), "SITTING LAW 5") {
		t.Fatalf("law_status does not say a sitting is open, or why the writers would refuse: %v", st)
	}
	for _, tool := range []string{"law_add", "law_seal"} {
		got, err := lawCall(t, reg, tr, tool, good, true)
		if err != nil || got["state"] != "refused" || !strings.Contains(fmt.Sprint(got["why"]), "SITTING LAW 5") {
			t.Fatalf("%s ran with a sitting open: %v %v", tool, got, err)
		}
	}
	if string(lawBytes(t, tn, lawLedger)) != string(before) {
		t.Fatal("the ledger moved with a sitting open")
	}

	// A CHAIN THAT REFUSES: nothing is laid over it.
	tn, tr, reg, stub := lawWorld(t, true)
	stub.breakAfter = 1
	before = lawBytes(t, tn, lawLedger)
	for _, tool := range []string{"law_add", "law_seal"} {
		got, _ := lawCall(t, reg, tr, tool, good, true)
		if got["state"] != "refused" || !strings.Contains(fmt.Sprint(got["why"]), "does not prove whole") {
			t.Fatalf("%s worked over a chain that refuses: %v", tool, got)
		}
	}
	if string(lawBytes(t, tn, lawLedger)) != string(before) || len(stub.calls) == 0 {
		t.Fatal("the ledger moved over a refusing chain")
	}

	// A LEDGER A HAND CANNOT APPEND TO CLEANLY (a carriage return, or no final line break) is said to be a file to fix by hand.
	for _, body := range []string{lawFixLedgerSealed + "\r\n", strings.TrimSuffix(lawFixLedgerSealed, "\n")} {
		tn, tr, reg, _ := lawWorld(t, false)
		if err := os.WriteFile(filepath.Join(tn.Home, "law", lawLedger), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		got, _ := lawCall(t, reg, tr, "law_add", good, true)
		if got["state"] != "refused" {
			t.Fatalf("an entry was appended to a ledger that is not a clean file: %v", got)
		}
	}
}

// THE VERDICT IS THE WALK'S OWN WORDS, NOT ITS EXIT CODE: a law.py that exits clean without saying the chain proves whole (a stand-in, the wrong
// interpreter, a page of usage) has proved nothing, so law_status does not show it whole and nothing is set or sealed over it.
func TestTheVerdictIsTheWalksOwnWordsNotItsExitCode(t *testing.T) {
	tn, tr, reg, stub := lawWorld(t, true)
	stub.silent = true
	before, links := lawBytes(t, tn, lawLedger), lawBytes(t, tn, lawChain)
	st, err := lawCall(t, reg, tr, "law_status", map[string]any{}, false)
	if err != nil || st["ok"] != false {
		t.Fatalf("a walk that said nothing was shown as whole: %v %v", st, err)
	}
	for _, tool := range []string{"law_add", "law_seal"} {
		got, _ := lawCall(t, reg, tr, tool, map[string]any{"title": "Over silence", "text": "A law laid over a walk that said nothing."}, true)
		if got["state"] != "refused" || !strings.Contains(fmt.Sprint(got["why"]), "does not prove whole") {
			t.Fatalf("%s worked over a walk that never said the chain was whole: %v", tool, got)
		}
	}
	if string(lawBytes(t, tn, lawLedger)) != string(before) || string(lawBytes(t, tn, lawChain)) != string(links) {
		t.Fatal("the ledger or the chain moved over a walk that said nothing")
	}
}

// AN ENTRY THAT LEFT THE CHAIN REFUSING IS TAKEN BACK, to the byte.
func TestLawAddTakesBackAnEntryThatLeftTheChainRefusing(t *testing.T) {
	tn, tr, reg, stub := lawWorld(t, true)
	stub.breakAfter = 2 // the walk before the append is clean; the walk after it refuses
	before := lawBytes(t, tn, lawLedger)
	got, err := lawCall(t, reg, tr, "law_add", map[string]any{"title": "Taken back", "text": "A law that the chain would not take."}, true)
	if err != nil || got["state"] != "refused" || !strings.Contains(fmt.Sprint(got["why"]), "taken back") {
		t.Fatalf("an entry the chain refused was not reported as taken back: %v %v", got, err)
	}
	if string(lawBytes(t, tn, lawLedger)) != string(before) {
		t.Fatal("the ledger is not byte for byte what it was before the entry that was taken back")
	}
}

func TestLawSealRunsLawPyOnTheLedgerAndNothingElse(t *testing.T) {
	tn, tr, reg, stub := lawWorld(t, true)
	got, err := lawCall(t, reg, tr, "law_seal", map[string]any{"note": "sealed on his word, in chat"}, true)
	if err != nil || got["state"] != "sealed" {
		t.Fatalf("the ledger was not sealed: %v %v", got, err)
	}
	var seal []string
	for _, c := range stub.calls {
		if len(c) > 2 && c[2] == "seal" {
			seal = c
		}
	}
	if len(seal) != 6 || seal[1] != filepath.Join(tn.Home, "law", "law.py") || seal[3] != lawLedger || seal[4] != "--note" || seal[5] != "sealed on his word, in chat" {
		t.Fatalf("law.py was not asked to seal the ledger, and only the ledger, with his note: %v", seal)
	}
	size := len(lawFixLedgerSealed) + len(lawFixLedgerDraft)
	if got["sealed_to"] != float64(size) || got["draft_bytes"] != float64(0) {
		t.Fatalf("the ledger is not sealed to its last byte: %v", got)
	}
	st, _ := lawCall(t, reg, tr, "law_status", map[string]any{}, false)
	for _, f := range st["files"].([]any) {
		if m := f.(map[string]any); m["name"] == lawLedger && m["state"] != "sealed" {
			t.Fatalf("after the seal the ledger reads %v", m["state"])
		}
	}
	if !strings.Contains(holdsText(t, tn), `"what":"law_sealed"`) {
		t.Fatal("sealing the ledger left nothing in the holds record")
	}

	// NOTHING BELOW THE SEAL, NOTHING TO SEAL; and a seal law.py would not make says so and writes nothing.
	got, _ = lawCall(t, reg, tr, "law_seal", map[string]any{}, true)
	if got["state"] != "refused" || !strings.Contains(fmt.Sprint(got["why"]), "nothing is written below") {
		t.Fatalf("a seal was asked for with nothing below the seal: %v", got)
	}
	tn2, tr2, reg2, stub2 := lawWorld(t, true)
	stub2.sealFails = true
	links := string(lawBytes(t, tn2, lawChain))
	got, _ = lawCall(t, reg2, tr2, "law_seal", map[string]any{}, true)
	if got["state"] != "refused" || !strings.Contains(fmt.Sprint(got["why"]), "would not seal") || string(lawBytes(t, tn2, lawChain)) != links {
		t.Fatalf("a seal law.py refused was not reported, or the chain moved: %v", got)
	}
	got, _ = lawCall(t, reg2, tr2, "law_seal", map[string]any{"note": "two\nlines"}, true)
	if got["state"] != "refused" || !strings.Contains(fmt.Sprint(got["why"]), "one line") {
		t.Fatalf("a note of two lines was taken: %v", got)
	}
}

func TestTheLawToolsRefuseWhileAnEngineIsOpenOnTheWorld(t *testing.T) {
	_, tn := standing(t)
	for rel, body := range map[string]string{"law/law.py": "# stand-in\n", "law/" + lawLedger: lawFixLedgerSealed,
		"law/" + lawChain: lawLine(1, "DIRECT", "operator", lawLedger, lawFixLedgerSealed, len(lawFixLedgerSealed)) + "\n"} {
		p := filepath.Join(tn.Home, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	stub := &lawStub{t: t, home: tn.Home}
	was := lawRun
	lawRun = stub.run
	t.Cleanup(func() { lawRun = was })
	tr := tenant.NewRegistry()
	_ = tr.Add("t", tn.Home)
	_ = tr.SetDefault("t")
	reg := Build(tr, Options{})
	got, err := lawCall(t, reg, tr, "law_add", map[string]any{"title": "Engine open", "text": "A law that must wait for the engine to close."}, true)
	if err != nil || got["state"] != "refused" || !strings.Contains(fmt.Sprint(got["why"]), "engine is open") {
		t.Fatalf("a law was set while an engine stood open on the world: %v %v", got, err)
	}
	if st, _ := lawCall(t, reg, tr, "law_status", map[string]any{}, false); st["sitting_open"] != true {
		t.Fatalf("law_status does not say the engine's sitting is open: %v", st["sitting_open"])
	}
}

// THE ANCHOR THE DOOR READS IS THE SHAPE LAW.PY WRITES, and a forged one is not read: the cross-language stroke in the core's suite holds this pattern to
// law.py's own source; this holds what the pattern does.
func TestTheAnchorTheDoorReadsIsLawPysAndForgeriesAreNotRead(t *testing.T) {
	fp := strings.Repeat("ab", 32)
	for _, ok := range []string{
		"DIRECT by operator -> law/LAW_LEDGER.md sha256:" + fp + " bytes:15509",
		"RULING by manjuel -> law/LAW_002_THE_TWELVE.md sha256:" + fp,
		"DIRECT by operator -> Archive/law/SITTING_LAWS.md sha256:" + fp,
	} {
		if !lawAnchorRe.MatchString(ok) {
			t.Fatalf("an anchor law.py writes is not read: %q", ok)
		}
	}
	for _, bad := range []string{
		"cf. -> law/LAW_T_DECOY.md sha256:" + fp,
		"DIRECT by operator -> law/../.env sha256:" + fp,
		"DIRECT by operator -> law/X.md sha256:" + fp + " bytes:0",
		"DIRECT by operator -> law/X.md sha256:" + fp[:60],
	} {
		if lawAnchorRe.MatchString(bad) {
			t.Fatalf("a forged anchor was read: %q", bad)
		}
	}
	if !regexp.MustCompile(`proves whole`).MatchString("the law chain proves whole: 6 links, head 07491469cd7d6d6c") || !lawProvenRe.MatchString("the law chain proves whole: 6 links, head 07491469cd7d6d6c") {
		t.Fatal("the walk's own verdict line is not recognised")
	}
	if lawProvenRe.MatchString("THE CHAIN REFUSES:\n  - link #3 fingerprint MISMATCH: X.md") {
		t.Fatal("a refusing walk was read as a proof")
	}
}
