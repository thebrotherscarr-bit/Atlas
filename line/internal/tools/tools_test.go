// Strokes for the readers: records, proofs and seats.
//
// `internal/tools` carries every MCP tool handler and had none of these until
// 2026-09-10. gitctl_test.go covered the verbs that WRITE; these cover the three
// that READ the estate's own record, which is where a quiet wrong answer does
// the most damage -- a page that misreports what was proven is worse than a page
// that shows nothing, and all three of these exist because an earlier page did
// exactly that.
//
// HERMETIC BY LAW 5: every stroke builds its own ground in t.TempDir(). Nothing
// reads the estate's real record, nothing writes, nothing reaches a network.
//
// runstream.go is NOT covered here and that is deliberate: RunStream, Answer-
// Stream and ListenStream all require a live engine on an open sitting, which is
// not a thing a hermetic stroke can stand up. It is named in tests/PROVING.md as
// still uncovered rather than papered over with a mock that would prove the mock.
package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"atlas/line/internal/auth"
	"atlas/line/internal/engine"
	"atlas/line/internal/flow"
	"atlas/line/internal/tenant"
)

// --- ground ----------------------------------------------------------------

func world(t *testing.T, files map[string]string) tenant.Tenant {
	t.Helper()
	home := t.TempDir()
	for rel, body := range files {
		p := filepath.Join(home, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return tenant.Tenant{Name: "probe", Home: home, Manifest: tenant.DefaultManifest()}
}

// jsonOf runs a tool and parses what it answered, failing on a hard error.
func jsonOf(t *testing.T, fn func(tenant.Tenant, map[string]any) (string, error),
	tn tenant.Tenant, args map[string]any) map[string]any {
	t.Helper()
	out, err := fn(tn, args)
	if err != nil {
		t.Fatalf("tool errored: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("tool did not answer JSON: %v\n%s", err, out)
	}
	return m
}

func has(t *testing.T, got, want, why string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("%s\n  wanted %q in:\n%s", why, want, got)
	}
}

// --- records: what the estate carries ---------------------------------------

func TestRecordsSortsByWhatADocumentIs(t *testing.T) {
	tn := world(t, map[string]string{
		"CLAUDE.md":            "the rules",
		"CHANGELOG.md":         "what changed",
		"SPEC.md":              "what it is",
		"pipelines.md":         "## Pipeline: default\n1. Steward\n",
		"law/ESTATE_LAWS.md":   "the ten",
		"agents/steward.md":    "## Steward\n",
		"skills/git_status.md": "# Skill\n",
		"logs/a_run.md":        "a transcript",
		"WANDERING.md":         "claimed by nothing",
	})
	d := jsonOf(t, toolRecords, tn, nil)

	kinds := map[string]int{}
	for _, k := range d["kinds"].([]any) {
		m := k.(map[string]any)
		kinds[m["kind"].(string)] = int(m["count"].(float64))
	}
	for kind, n := range map[string]int{
		"doctrine": 2, "record": 1, "spec": 1, "commands": 1,
		"agents": 1, "skills": 1, "logs": 1, "other": 1,
	} {
		if kinds[kind] != n {
			t.Fatalf("kind %q holds %d, wanted %d (%v)", kind, kinds[kind], n, kinds)
		}
	}
	// ALPHABETICAL WOULD PUT AGENTS ABOVE THE LAW. The order is the estate's
	// furniture: what binds a hand first, what happened second.
	first := d["kinds"].([]any)[0].(map[string]any)["kind"]
	if first != "doctrine" {
		t.Fatalf("the law does not come first: %v", first)
	}
	// A root file no set claims is still SHOWN -- nothing is hidden.
	check := jsonOf(t, toolRecords, tn, map[string]any{"kind": "other"})
	has(t, check["documents"].([]any)[0].(map[string]any)["name"].(string),
		"WANDERING.md", "an unclaimed root document must still be listed")
}

func TestLawIsMarkedSealed(t *testing.T) {
	tn := world(t, map[string]string{
		"law/SITTING_LAWS.md": "the laws",
		"CLAUDE.md":           "the rules",
	})
	d := jsonOf(t, toolRecords, tn, map[string]any{"kind": "doctrine"})
	sealed := map[string]bool{}
	for _, x := range d["documents"].([]any) {
		m := x.(map[string]any)
		s, _ := m["sealed"].(bool)
		sealed[m["name"].(string)] = s
	}
	if !sealed["law/SITTING_LAWS.md"] {
		t.Fatal("a law is not marked sealed")
	}
	if sealed["CLAUDE.md"] {
		t.Fatal("a root document is marked sealed and is not")
	}
}

// IT IS A LIST, NOT A PATH. Nothing is joined onto Home from the caller's
// string, so there is no traversal to defend against -- the escape simply is
// not IN the listing, and so cannot resolve.
func TestARecordNameIsMatchedAgainstTheListingNotTheDisk(t *testing.T) {
	tn := world(t, map[string]string{"CLAUDE.md": "the rules"})
	if err := os.WriteFile(filepath.Join(tn.Home, ".env"),
		[]byte("MANJUEL_GIT_REMOTE=1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"../.env", ".env", "/etc/passwd", "C:\\keys.txt",
		"law/../.env", "../../CLAUDE.md",
	} {
		out, err := toolRecords(tn, map[string]any{"name": name})
		if err == nil {
			t.Fatalf("%q was served: %s", name, out[:120])
		}
		has(t, err.Error(), "no record named", "the refusal must say it is not in the list")
	}
}

func TestOnlyMarkdownIsADocument(t *testing.T) {
	tn := world(t, map[string]string{
		"CLAUDE.md":        "the rules",
		"secrets.txt":      "not a document",
		"index/vectors.db": "not a document",
		"notes.json":       "not a document",
	})
	d := jsonOf(t, toolRecords, tn, nil)
	body, _ := json.Marshal(d)
	for _, gone := range []string{"secrets.txt", "vectors.db", "notes.json"} {
		if strings.Contains(string(body), gone) {
			t.Fatalf("%s was listed as a document", gone)
		}
	}
	if int(d["count"].(float64)) != 1 {
		t.Fatalf("wanted 1 document, got %v", d["count"])
	}
}

// The receipt is the point: a page showing a document can be checked against
// the disk without trusting the page.
func TestADocumentComesBackWithTheShaOfTheBytesServed(t *testing.T) {
	body := "# The rules\nevery turn.\n"
	tn := world(t, map[string]string{"CLAUDE.md": body})
	d := jsonOf(t, toolRecords, tn, map[string]any{"name": "CLAUDE.md"})
	sum := sha256.Sum256([]byte(body))
	if d["sha256"] != hex.EncodeToString(sum[:]) {
		t.Fatalf("sha256 does not match the bytes served:\n  got  %v\n  want %s",
			d["sha256"], hex.EncodeToString(sum[:]))
	}
	if d["text"] != body {
		t.Fatal("the text served is not the file")
	}
	if int(d["bytes"].(float64)) != len(body) {
		t.Fatalf("byte count %v != %d", d["bytes"], len(body))
	}
}

func TestAnAbsentNameIsDeniedBySayingWhatThereIs(t *testing.T) {
	tn := world(t, map[string]string{"CLAUDE.md": "x", "CHANGELOG.md": "y"})
	_, err := toolRecords(tn, map[string]any{"name": "NOPE.md"})
	if err == nil {
		t.Fatal("an absent name was not refused")
	}
	has(t, err.Error(), "doctrine (1)", "the refusal must name the kinds carried")
	has(t, err.Error(), "record (1)", "the refusal must name the kinds carried")
}

// The count reported is the TRUE count, not the shown one: a listing that
// silently truncated would be a page lying about the record.
func TestTheTranscriptsAreCappedNewestFirst(t *testing.T) {
	files := map[string]string{"CLAUDE.md": "x"}
	for i := 0; i < LOGS_SHOWN+7; i++ {
		files[filepath.ToSlash(filepath.Join("logs",
			"run_"+string(rune('a'+i/26))+string(rune('a'+i%26))+".md"))] = "t"
	}
	tn := world(t, files)
	d := jsonOf(t, toolRecords, tn, map[string]any{"kind": "logs"})
	if n := len(d["documents"].([]any)); n != LOGS_SHOWN {
		t.Fatalf("logs returned %d, wanted the %d cap", n, LOGS_SHOWN)
	}
	docs := d["documents"].([]any)
	for i := 1; i < len(docs); i++ {
		a := docs[i-1].(map[string]any)["modified"].(string)
		b := docs[i].(map[string]any)["modified"].(string)
		if a < b {
			t.Fatal("the transcripts are not newest first")
		}
	}
}

// --- proofs: what a world has actually proved -------------------------------

// A world that never ran a suite has not FAILED at anything, and one missing
// file must never blank the other two.
func TestAProofPageWithAGapBeatsNoProofPage(t *testing.T) {
	tn := world(t, map[string]string{
		"tests/last_run.json": `{"strokes":{"passed":2106,"total":2106,"green":true}}`,
	})
	d := jsonOf(t, toolProofs, tn, nil)
	if d["suites"] == nil {
		t.Fatal("the suite that IS on disk was not read")
	}
	for _, absent := range []string{"standups_error", "parity_error"} {
		s, _ := d[absent].(string)
		if s != "never run in this world" {
			t.Fatalf("%s says %q, wanted the honest sentence", absent, s)
		}
	}
}

// sessions.jsonl is append-only and a CLOSING line supersedes its opening one.
// Keyed by n, last line wins -- get this wrong and every sitting counts twice.
func TestAClosingLineSupersedesItsOpening(t *testing.T) {
	tn := world(t, map[string]string{
		"sessions/sessions.jsonl": strings.Join([]string{
			`{"n":1,"started":"a","runs":[]}`,
			`{"n":2,"started":"b","runs":[{}]}`,
			`{"n":1,"started":"a","ended":"z","toll_paid":true,"runs":[{},{}]}`,
			`  `,
			`{"n":3,"started":"c","runs":[]}`,
		}, "\n"),
	})
	rec := jsonOf(t, toolProofs, tn, nil)["record"].(map[string]any)
	if int(rec["sittings"].(float64)) != 3 {
		t.Fatalf("sittings counted %v, wanted 3 -- a superseded line was double counted", rec["sittings"])
	}
	if int(rec["tolled"].(float64)) != 1 {
		t.Fatalf("tolled %v, wanted 1", rec["tolled"])
	}
	// 2 from the closing line of sitting 1, 1 from sitting 2, 0 from sitting 3.
	if int(rec["runs"].(float64)) != 3 {
		t.Fatalf("runs %v, wanted 3", rec["runs"])
	}
	// Sittings 2 and 3 never ended.
	if int(rec["still_open"].(float64)) != 2 {
		t.Fatalf("still_open %v, wanted 2", rec["still_open"])
	}
}

// A half-written last line is what a KILLED PROCESS leaves. The lines before it
// are still true, so it is dropped rather than failing the whole read.
func TestAHalfWrittenLineDoesNotLoseTheLinesBeforeIt(t *testing.T) {
	tn := world(t, map[string]string{
		"sessions/sessions.jsonl": "{\"n\":1,\"started\":\"a\",\"ended\":\"z\",\"runs\":[]}\n{\"n\":2,\"star",
	})
	rec := jsonOf(t, toolProofs, tn, nil)["record"].(map[string]any)
	if rec["error"] != nil {
		t.Fatalf("a torn last line failed the whole read: %v", rec["error"])
	}
	if int(rec["sittings"].(float64)) != 1 {
		t.Fatalf("sittings %v, wanted the 1 whole line", rec["sittings"])
	}
}

// Some numbers are the CORE's to define. A second definition here would drift
// from his the first time it changed, so they are named as not counted.
func TestTheNumbersTheCoreOwnsAreNamedNotRecounted(t *testing.T) {
	tn := world(t, map[string]string{"CLAUDE.md": "x"})
	rec := jsonOf(t, toolProofs, tn, nil)["record"].(map[string]any)
	named := []string{}
	for _, x := range rec["counted_by_the_engine"].([]any) {
		named = append(named, x.(string))
	}
	joined := strings.Join(named, " | ")
	has(t, joined, "memory entries", "the core's own counts must be named")
	has(t, joined, "index docs", "the core's own counts must be named")
}

// --- proofs: the record is read again only when it changes ------------------

var recordAt = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

// recordWorld carries the four files fromRecord keeps, each stamped recordAt
// so a stroke decides when that time moves.
func recordWorld(t *testing.T) tenant.Tenant {
	t.Helper()
	tn := world(t, map[string]string{
		"tests/run_history.jsonl":       `{"suite":"standup","run":"a"}` + "\n",
		"sessions/parity_history.jsonl": `{"p":"a"}` + "\n",
		"sessions/sessions.jsonl":       `{"n":1,"started":"a","ended":"z","runs":[]}` + "\n",
		"SEAT_LOG.md":                   "## 1 — sitting 1\n## 2 — sitting 2\n",
	})
	for _, rel := range []string{"tests/run_history.jsonl", "sessions/parity_history.jsonl",
		"sessions/sessions.jsonl", "SEAT_LOG.md"} {
		stampAt(t, tn, rel, recordAt)
	}
	return tn
}

func stampAt(t *testing.T, tn tenant.Tenant, rel string, at time.Time) {
	t.Helper()
	if err := os.Chtimes(filepath.Join(tn.Home, filepath.FromSlash(rel)), at, at); err != nil {
		t.Fatal(err)
	}
}

// readsCounted swaps readRecord for one that counts by file name.
func readsCounted(t *testing.T) map[string]int {
	t.Helper()
	n := map[string]int{}
	real := readRecord
	readRecord = func(p string) ([]byte, error) {
		n[filepath.Base(p)]++
		return real(p)
	}
	t.Cleanup(func() { readRecord = real })
	return n
}

// The Dashboard asks every 15 seconds for an answer that changes once a
// sitting. An unchanged record is read once, and every later answer is the
// first one, byte for byte.
func TestAnUnchangedRecordIsReadOnceAndAnsweredTheSame(t *testing.T) {
	tn := recordWorld(t)
	reads := readsCounted(t)
	first, err := toolProofs(tn, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 2; i <= 4; i++ {
		if again, _ := toolProofs(tn, nil); again != first {
			t.Fatalf("call %d answered differently from the first over an unchanged record", i)
		}
	}
	for _, f := range []string{"run_history.jsonl", "parity_history.jsonl", "sessions.jsonl", "SEAT_LOG.md"} {
		if reads[f] != 1 {
			t.Fatalf("%s read %d times over four calls; an unchanged file is read once (%v)", f, reads[f], reads)
		}
	}
}

// An append moves the size. Held at the same time, a ledger that grew is still
// read again, and the sitting appended to it is counted.
func TestALedgerThatGrewIsReadAgainUnderAnUnmovedTime(t *testing.T) {
	tn := recordWorld(t)
	if n := jsonOf(t, toolProofs, tn, nil)["record"].(map[string]any)["sittings"]; n != 1.0 {
		t.Fatalf("sittings %v, wanted 1", n)
	}
	p := filepath.Join(tn.Home, "sessions", "sessions.jsonl")
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(`{"n":2,"started":"b","ended":"y","runs":[]}` + "\n")
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		t.Fatal(err)
	}
	stampAt(t, tn, "sessions/sessions.jsonl", recordAt)
	if n := jsonOf(t, toolProofs, tn, nil)["record"].(map[string]any)["sittings"]; n != 2.0 {
		t.Fatalf("sittings %v after a sitting was appended: the kept answer was served over a longer file", n)
	}
}

// A rewrite can keep the size. It cannot keep the time, and a moved time is
// read again.
func TestAFileRewrittenToTheSameSizeIsReadAgainWhenItsTimeMoves(t *testing.T) {
	tn := recordWorld(t)
	tolls := func() any {
		return jsonOf(t, toolProofs, tn, nil)["record"].(map[string]any)["seat_log_tolls"]
	}
	if n := tolls(); n != 2.0 {
		t.Fatalf("tolls %v, wanted 2", n)
	}
	p := filepath.Join(tn.Home, "SEAT_LOG.md")
	before, _ := os.Stat(p)
	write(t, tn.Home, "SEAT_LOG.md", "## 1 — sitting 1\n## 2 — Sitting 2\n") // one heading stops matching
	if after, _ := os.Stat(p); after.Size() != before.Size() {
		t.Fatalf("the rewrite changed the size (%d -> %d); this stroke must hold it", before.Size(), after.Size())
	}
	stampAt(t, tn, "SEAT_LOG.md", recordAt.Add(time.Second))
	if n := tolls(); n != 1.0 {
		t.Fatalf("tolls %v: SEAT_LOG was rewritten and its time moved, and the old count was served", n)
	}
}

// A file held for a moment -- a scanner, a writer mid-flush -- fails one read.
// The page says so once; the failure is not kept, and the next call reads.
func TestAReadThatFailedIsNotKept(t *testing.T) {
	tn := recordWorld(t)
	real := readRecord
	t.Cleanup(func() { readRecord = real })
	readRecord = func(p string) ([]byte, error) {
		if filepath.Base(p) == "sessions.jsonl" {
			return nil, errors.New("held by another process")
		}
		return real(p)
	}
	rec := jsonOf(t, toolProofs, tn, nil)["record"].(map[string]any)
	if s, _ := rec["error"].(string); !strings.Contains(s, "held by another process") {
		t.Fatalf("the failed read was not reported: %v", rec)
	}
	readRecord = real
	rec = jsonOf(t, toolProofs, tn, nil)["record"].(map[string]any)
	if rec["error"] != nil || rec["sittings"] != 1.0 {
		t.Fatalf("the failure was kept over an unchanged file: %v", rec)
	}
}

// --- seats: the declarations, read as a SHAPE -------------------------------

// The core parses a declaration with two regexes and NEITHER NAMES A FIELD.
// So this returns whatever keys a declaration carries -- including one added
// tomorrow. Encoding a field list here is what would drift.
func TestASeatsFieldsAreWhateverItDeclares(t *testing.T) {
	tn := world(t, map[string]string{
		"agents/steward.md": "## Steward\n" +
			"- **Model Target:** llama3.2\n" +
			"- **Stage:** deliver\n" +
			"- **Invented Tomorrow:** and still read\n",
		"pipelines.md": "## Pipeline: default\n1. Steward\n2. Router (when: needs_tool)\n",
	})
	d := jsonOf(t, toolSeats, tn, nil)
	seat := d["seats"].([]any)[0].(map[string]any)
	if seat["name"] != "Steward" {
		t.Fatalf("name read as %v", seat["name"])
	}
	fields := seat["fields"].(map[string]any)
	// THE COLON LIVES INSIDE THE BOLD MARKERS in this estate's files, so a key
	// that kept it would label every field with a trailing colon.
	for k, want := range map[string]string{
		"Model Target": "llama3.2", "Stage": "deliver",
		"Invented Tomorrow": "and still read",
	} {
		if fields[k] != want {
			t.Fatalf("field %q read as %v, wanted %q (%v)", k, fields[k], want, fields)
		}
	}
	stands := seat["stands_in"].([]any)
	if len(stands) != 1 || stands[0].(map[string]any)["pipeline"] != "default" {
		t.Fatalf("the seat does not know where it stands: %v", stands)
	}
	if int(stands[0].(map[string]any)["step"].(float64)) != 1 {
		t.Fatal("the step is wrong")
	}
}

// A field with an empty value and a body beneath it IS the body's heading --
// the core's own shape. The prompt must come back whole, and not as a field.
func TestASystemPromptIsTheBodyBeneathIt(t *testing.T) {
	body := "You are the Steward. Speak plainly.\nNever invent a tool result."
	tn := world(t, map[string]string{
		"agents/steward.md": "## Steward\n- **Model Target:** llama3.2\n" +
			"- **System Prompt:**\n" + body + "\n",
	})
	d := jsonOf(t, toolSeats, tn, nil)
	seat := d["seats"].([]any)[0].(map[string]any)
	if got, _ := seat["prompt"].(string); got != body {
		t.Fatalf("prompt read as %q", got)
	}
	if _, still := seat["fields"].(map[string]any)["System Prompt"]; still {
		t.Fatal("the prompt is also sitting in fields, counted twice")
	}
	if int(seat["prompt_chars"].(float64)) != len(body) {
		t.Fatal("prompt_chars does not match the prompt")
	}
}

// A world with no agents/ has not failed at anything either.
func TestAWorldWithNoSeatsSaysSoRatherThanCrashing(t *testing.T) {
	tn := world(t, map[string]string{"CLAUDE.md": "x"})
	d := jsonOf(t, toolSeats, tn, nil)
	if d["seats_error"] == nil {
		t.Fatal("an absent agents/ was not reported")
	}
	if d["world"] != "probe" {
		t.Fatal("the answer does not name its world")
	}
}

// The registry is the contract the door is: a tool that writes must SAY it
// writes, because that flag is what the read-only table is refused by.
func TestTheRegistryAgreesWithWhatEachToolDoes(t *testing.T) {
	// Built the way the door builds it, against an empty tenant registry: the
	// wiring is what is under test, not any world's contents.
	r := Build(tenant.NewRegistry(), Options{})
	writes := map[string]bool{}
	for _, name := range r.Names() {
		writes[name] = r.byName[name].Writes
	}
	for _, name := range []string{"git", "git_diff", "git_remote", "records", "proofs", "seats", "projects"} {
		if v, ok := writes[name]; !ok {
			t.Fatalf("%s is not registered", name)
		} else if v {
			t.Fatalf("%s is declared as writing and does not", name)
		}
	}
	for _, name := range []string{"git_commit", "git_push", "git_pull", "git_branch", "git_tag"} {
		if v, ok := writes[name]; !ok {
			t.Fatalf("%s is not registered", name)
		} else if !v {
			t.Fatalf("%s writes and does not declare it", name)
		}
	}
}

// ---- ADR-006: the tier contract ------------------------------------------

// EVERY CORE TOOL STANDS ALONE. ADR-006's whole claim is that the door is a
// product: point it at any directory on any machine and 75 of its 78 tools
// answer. This stroke is what makes that a fact rather than a sentence.
//
// It calls every tool declaring TierCore against a bare temp tenant with NO
// engine wired and NO Rust binary configured, and refuses to accept a failure
// whose CAUSE is one of those two. A core tool may absolutely refuse -- for a
// missing argument, an absent file, a path outside the wall -- it just may not
// refuse because the estate's own machinery is not there.
//
// The tier is the zero value, so a new tool is CORE until it says otherwise.
// This is the thing that catches the one that should have said otherwise.
func TestEveryCoreToolStandsAlone(t *testing.T) {
	home := t.TempDir()
	tr := tenant.NewRegistry()
	if err := tr.Add("t", home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	// THE SPINE IS DENIED, NOT MERELY UNCONFIGURED. Options{} leaves
	// AtlasBin empty, but findAtlas then walks the tree and FINDS the real
	// binary on any machine where cargo has run — so verify_chain succeeded
	// here and the stroke reported it core. Pointing ATLAS_BIN at a path that
	// cannot exist is what makes "needs nothing" mean it.
	t.Setenv("ATLAS_BIN", filepath.Join(home, "NO-SUCH-SPINE.exe"))

	// Options{} on purpose: no CoreCmd, no AtlasBin. A tenant that looks
	// nothing like manjuel -- no agents/, no pipelines.md, no sessions/.
	reg := Build(tr, Options{})

	// A real file, so a tool given a path reaches its body rather than
	// bouncing off "no such file".
	if err := os.WriteFile(filepath.Join(home, "probe.txt"),
		[]byte("a probe\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// The cause a core tool may never fail for. Matched on the message the
	// caller would actually see.
	forbidden := []string{
		"no engine wired",
		"--manjuel",
		"was started without",
		`exec: "atlas"`,
		"atlas-bin",
		"ATLAS_BIN",
		"NO-SUCH-SPINE",
		// run_start refuses with "no engine is open ... env_open first",
		// which names none of the above. A first cut missed it entirely and
		// reported an engine tool as core.
		"no engine is open",
		"env_open first",
		"nothing to wake",
	}

	for _, tool := range reg.All() {
		if tool.Tier != TierCore {
			continue
		}
		t.Run(tool.Name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("a core tool panicked against a bare tenant: %v", r)
				}
			}()
			_, err := reg.Call(tr, tool.Name, argsFor(tool, home), Caller{Name: "prove"})
			if err == nil {
				return // answered; nothing more to ask of it
			}
			msg := err.Error()
			for _, bad := range forbidden {
				if strings.Contains(msg, bad) {
					t.Errorf("declares TierCore but failed because the estate's "+
						"machinery is absent.\n  refusal: %s\n  "+
						"Either mark it TierEngine/TierSpine, or stop it "+
						"reaching for what a core tool may not need.", msg)
					return
				}
			}
		})
	}
}

// A TIER IS A PROMISE AND MUST BE ONE OF THE THREE. A tool carrying a tier
// nobody defined would sort into "core" by String()'s default and quietly
// claim a contract it never made.
func TestEveryToolCarriesAKnownTier(t *testing.T) {
	reg := Build(tenant.NewRegistry(), Options{})
	for _, tool := range reg.All() {
		switch tool.Tier {
		case TierCore, TierSpine, TierEngine:
		default:
			t.Errorf("%s carries tier %d, which is not one of core/spine/engine",
				tool.Name, int(tool.Tier))
		}
	}
}

// argsFor fills a tool's OWN declared arguments so the call reaches the tool's
// body instead of bouncing off its argument check. A first cut passed only
// {"project": "t"} -- and run_start and verify_chain both refused for a
// missing argument BEFORE they ever reached for the engine or the spine, so
// the stroke reported them clean when they were not. Reading the declaration
// is what makes the classification honest.
func argsFor(tool Tool, home string) map[string]any {
	args := map[string]any{"project": "t"}
	for _, a := range tool.Args {
		name := strings.TrimRight(a, "?")
		if name == "" || name == "project" {
			continue
		}
		switch {
		case strings.Contains(name, "path") || strings.Contains(name, "file"):
			args[name] = "probe.txt" // a real file, written below
		case name == "name" || name == "document" || name == "doc":
			args[name] = "probe.txt"
		case name == "kind":
			args[name] = "record"
		case strings.Contains(name, "content") || strings.Contains(name, "text") ||
			strings.Contains(name, "message") || strings.Contains(name, "body"):
			args[name] = "a probe from TestEveryCoreToolStandsAlone"
		default:
			args[name] = "probe"
		}
	}
	return args
}

// --- the council seam: what a flow's `run` node is handed (2026-09-22) -------

// A TURN THAT DID NOT DELIVER IS NOT AN ANSWER. Every terminal event but
// `delivery` used to be handed to the flow as the node's OUTPUT, so the run
// log wrote it "ok" and a flow could walk its whole happy path having done
// nothing. `command` bites hardest: a runtime error inside the engine leaves
// that event carrying the objective's own words, so the node's recorded answer
// was its own question.
func TestOnlyADeliveryIsATurnsAnswer(t *testing.T) {
	out, err := deliveryOf(engine.Event{"event": "delivery", "text": "the work, done"})
	if err != nil || out != "the work, done" {
		t.Fatalf("a delivery must be the answer: %q, %v", out, err)
	}
	for _, kind := range []string{"refused", "aborted", "cancelled", "unreachable", "command"} {
		out, err := deliveryOf(engine.Event{"event": kind, "text": "do the thing"})
		if err == nil {
			t.Fatalf("%q was handed up as an answer: %q", kind, out)
		}
		if out != "" {
			t.Fatalf("%q refused and still returned text: %q", kind, out)
		}
		for _, want := range []string{"did not deliver", kind} {
			if !strings.Contains(err.Error(), want) {
				t.Fatalf("the refusal over %q must name %q: %v", kind, want, err)
			}
		}
	}
	// A turn that ended with nothing to say still refuses, and says that.
	if _, err := deliveryOf(engine.Event{"event": "refused"}); err == nil ||
		!strings.Contains(err.Error(), "nothing further") {
		t.Fatalf("an empty refusal must still refuse in words: %v", err)
	}
}

// ONE FLOW MUST NOT FREEZE EVERY WORLD. flow_run, flow_resume and flow_replay
// held askLock -- one mutex across every tenant -- for a whole run, so a flow
// froze the glass's chat, every prompt run and every other world's flow.
func TestAFlowLocksItsOwnWorldAndNoOther(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	if flowLock(a) != flowLock(a) {
		t.Fatal("one world was given two locks; two flows on it would run at once")
	}
	if flowLock(a) == flowLock(b) {
		t.Fatal("two worlds share one lock; a flow on one freezes the other")
	}
	flowLock(a).Lock()
	defer flowLock(a).Unlock()

	free := make(chan struct{})
	go func() {
		flowLock(b).Lock()
		flowLock(b).Unlock()
		askLock.Lock()
		askLock.Unlock()
		close(free)
	}()
	select {
	case <-free:
	case <-time.After(5 * time.Second):
		t.Fatal("a flow in flight holds another world's lock, or askLock itself -- " +
			"the glass's chat and every other world wait on it")
	}

	// AND THE THREE TOOLS TAKE IT. The lock above is only a lock; what was
	// wrong was which one flow_run, flow_resume and flow_replay reached for,
	// and no hermetic stroke can hold a real run open to watch. The source
	// says it, the way internal/flow's own no-finish stroke reads source.
	src, err := os.ReadFile("tools.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, fn := range []string{"func toolFlowRun(", "func toolFlowResume(", "func toolFlowReplay("} {
		i := strings.Index(string(src), fn)
		if i < 0 {
			t.Fatalf("%s is gone from tools.go", fn)
		}
		body := string(src)[i:]
		if j := strings.Index(body[1:], "\nfunc "); j >= 0 {
			body = body[:j+1]
		}
		if !strings.Contains(body, "flowLock(t.Home)") {
			t.Fatalf("%s does not take its world's own flow lock", fn)
		}
		if strings.Contains(body, "askLock.Lock()") {
			t.Fatalf("%s holds askLock for a whole run -- every world waits on it", fn)
		}
	}
}

// --- the head is named at the fire (2026-09-23) ------------------------------

// THE COUNCIL TAKES A HEAD WITHOUT KEEPING IT. flow hands the run's voice down
// with WithVoice on the way into Run, Resume and Replay alike; if that bound
// the engine itself rather than a copy, the head would outlive its run and the
// next flow on that world -- or a second flow in flight beside it -- would
// answer on a model nobody named for it. Which is the parity failing silently:
// both columns on one model, both labelled honestly, and the number meaning
// nothing.
func TestTheCouncilTakesAHeadByCopyAndNotByKeeping(t *testing.T) {
	home := t.TempDir()
	base := council(home)
	bound := base.(councilEngine).WithHead(flow.Head{Voice: "head-a:latest"})
	if bound.(councilEngine).head.Voice != "head-a:latest" {
		t.Fatal("WithHead did not bind the head it was handed")
	}
	if base.(councilEngine).head.Named() {
		t.Fatal("WithHead bound the engine itself -- the head outlives its run")
	}
	other := base.(councilEngine).WithHead(flow.Head{Voice: "head-b:latest"})
	if bound.(councilEngine).head.Voice != "head-a:latest" {
		t.Fatal("two flows in flight crossed heads")
	}
	if other.(councilEngine).head.Voice != "head-b:latest" {
		t.Fatal("the second binding did not take")
	}
	// And a fresh council is unheaded: the ground's declared targets.
	if council(home).(councilEngine).head.Named() {
		t.Fatal("a council was born holding a head")
	}
	// A PER-SEAT HEAD RIDES THE SAME WAY (2026-09-23, "then B underneath it").
	seated := base.(councilEngine).WithHead(flow.Head{
		Voices: map[string]string{"Steward": "phi4-mini:latest"}})
	if seated.(councilEngine).head.Voices["Steward"] != "phi4-mini:latest" {
		t.Fatal("WithHead dropped the per-seat map")
	}
	if base.(councilEngine).head.Named() {
		t.Fatal("a per-seat head bound the registered engine")
	}
}

// flow_run says the head over the wire, and the two that CARRY one do not:
// resume and replay read it off the run's own start line, so a run cannot be
// finished or repeated on a head other than the one it began on.
func TestFlowRunTakesTheHeadAndResumeReplayDoNot(t *testing.T) {
	r := Build(tenant.NewRegistry(), Options{})
	run, ok := r.Get("flow_run")
	if !ok {
		t.Fatal("flow_run is gone from the registry")
	}
	if !hasArg(run.Args, "voice?") {
		t.Fatalf("flow_run does not offer the head over the wire: %v", run.Args)
	}
	for _, name := range []string{"flow_resume", "flow_replay"} {
		tl, ok := r.Get(name)
		if !ok {
			t.Fatalf("%s is gone from the registry", name)
		}
		if hasArg(tl.Args, "voice?") || hasArg(tl.Args, "voice") {
			t.Fatalf("%s takes a head from the caller; it must read the run's own: %v",
				name, tl.Args)
		}
	}
	// AND IT IS ACTUALLY PASSED. No hermetic stroke can fire a real flow here
	// -- a `run` node needs a standing engine -- so the source says it, the way
	// the flow-lock stroke above does.
	src, err := os.ReadFile("tools.go")
	if err != nil {
		t.Fatal(err)
	}
	body := funcBody(string(src), "func toolFlowRun(")
	if body == "" {
		t.Fatal("func toolFlowRun( is gone from tools.go")
	}
	if !strings.Contains(body, "flow.RunOn(") {
		t.Fatal("toolFlowRun fires through flow.Run -- the head it was given goes nowhere")
	}
	if !strings.Contains(body, `args["voice"]`) {
		t.Fatal("toolFlowRun never reads the head off the call")
	}
}

func hasArg(args []string, want string) bool {
	for _, a := range args {
		if a == want {
			return true
		}
	}
	return false
}

// funcBody is one function's source, from its signature to the next one.
func funcBody(src, sig string) string {
	i := strings.Index(src, sig)
	if i < 0 {
		return ""
	}
	body := src[i:]
	if j := strings.Index(body[1:], "\nfunc "); j >= 0 {
		body = body[:j+1]
	}
	return body
}

// A PER-SEAT HEAD COMES OFF THE CALL OR IT IS REFUSED (2026-09-23). A caller
// that meant to vary one seat and silently varied none would get a parity of a
// model against itself, and BOTH columns would look honest -- which is the one
// answer a parity must never be able to give. So every shape that is not a
// seat -> model object refuses in words rather than being dropped.
func TestAPerSeatHeadIsReadOffTheCallOrRefused(t *testing.T) {
	got, err := flowVoices(map[string]any{
		"voices": map[string]any{"Steward": "phi4-mini:latest"}})
	if err != nil || got["Steward"] != "phi4-mini:latest" {
		t.Fatalf("an object must be read straight: %v / %v", got, err)
	}
	// The same object as a string, exactly as `inputs` is accepted -- an MCP
	// schema hands every arg over as a string, so refusing this would put the
	// whole feature out of reach over the wire it was built for.
	got, err = flowVoices(map[string]any{"voices": `{"Steward":"phi4-mini:latest"}`})
	if err != nil || got["Steward"] != "phi4-mini:latest" {
		t.Fatalf("an object in a string must be read too: %v / %v", got, err)
	}
	for _, empty := range []any{nil, "", "   ", map[string]any{}} {
		got, err = flowVoices(map[string]any{"voices": empty})
		if err != nil || got != nil {
			t.Fatalf("%#v must name no head at all: %v / %v", empty, got, err)
		}
	}
	if got, err = flowVoices(map[string]any{}); err != nil || got != nil {
		t.Fatalf("a call with no voices must name no head: %v / %v", got, err)
	}
	for _, bad := range []struct {
		v    any
		says string
	}{
		{"Steward", "JSON object"},
		{42, "refused rather than run with a head nobody named"},
		{map[string]any{"Steward": 7}, "not a model"},
		{map[string]any{"Steward": "  "}, "names no model"},
	} {
		if _, err := flowVoices(map[string]any{"voices": bad.v}); err == nil {
			t.Fatalf("%#v was accepted", bad.v)
		} else if !strings.Contains(err.Error(), bad.says) {
			t.Fatalf("the refusal of %#v must say why: %v", bad.v, err)
		}
	}
}

// The wire offers it where a run is FIRED and nowhere else: resume and replay
// read the head off the run's own record, so a parity cannot be finished or
// repeated on a different one.
func TestTheWireOffersASeatMapOnlyWhereARunIsFired(t *testing.T) {
	r := Build(tenant.NewRegistry(), Options{})
	run, ok := r.Get("flow_run")
	if !ok {
		t.Fatal("flow_run is gone from the registry")
	}
	if !hasArg(run.Args, "voices?") {
		t.Fatalf("flow_run does not offer a per-seat head: %v", run.Args)
	}
	for _, name := range []string{"flow_resume", "flow_replay"} {
		tl, _ := r.Get(name)
		if hasArg(tl.Args, "voices?") || hasArg(tl.Args, "voices") {
			t.Fatalf("%s takes a head from the caller; it must read the run's "+
				"own: %v", name, tl.Args)
		}
	}
	src, err := os.ReadFile("tools.go")
	if err != nil {
		t.Fatal(err)
	}
	body := funcBody(string(src), "func toolFlowRun(")
	if !strings.Contains(body, "flowVoices(args)") {
		t.Fatal("toolFlowRun never reads the seat map off the call")
	}
	if !strings.Contains(body, "flow.Head{") {
		t.Fatal("toolFlowRun does not hand the flow a head")
	}
	// AND THE COUNCIL CARRIES IT TO THE ENGINE. flow's head and the engine's
	// are separate types on purpose -- `Voice` means "one model" to a flow and
	// "every seat" to the council -- and councilEngine.Turn is the one place
	// that knows both. A translation that dropped the map would leave a parity
	// varying nothing while its record named a seat.
	turn := funcBody(string(src), "func (c councilEngine) Turn(")
	if !strings.Contains(turn, "Voices: c.head.Voices") {
		t.Fatal("the council does not carry the seat map down to the engine")
	}
	if !strings.Contains(turn, "Model: c.head.Voice") {
		t.Fatal("the council does not carry the roster-wide head down to the engine")
	}
}

// P0-13 (2026-09-25): RBAC judges the CALLER the door built from the transport,
// never the args, on every call; open mode is said, not silent.
func TestRBACJudgesTheTransportNotTheArgs(t *testing.T) {
	home := t.TempDir()
	policy := `{"roles": {"open": {"name": "open", "permissions": {"*": "allow"}},
	                      "shut": {"name": "shut", "permissions": {"*": "deny"}}},
	            "assign": {"k-open": "open", "k-shut": "shut"}}`
	if err := os.WriteFile(filepath.Join(home, "rbac.json"), []byte(policy), 0o644); err != nil {
		t.Fatal(err)
	}
	tr := tenant.NewRegistry()
	if err := tr.Add("t", home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ATLAS_BIN", filepath.Join(home, "NO-SUCH-SPINE.exe"))
	reg := Build(tr, Options{})
	call := func(c Caller, args map[string]any) error {
		_, err := reg.Call(tr, "muster", args, c)
		return err
	}
	if err := call(Caller{Name: "k-shut"}, map[string]any{}); err == nil || !strings.Contains(err.Error(), `rbac: caller "k-shut"`) {
		t.Fatalf("a key whose role denies must be refused by its own name: %v", err)
	}
	if err := call(Caller{Name: "k-open"}, map[string]any{}); err != nil && strings.Contains(err.Error(), "rbac:") {
		t.Fatalf("a key whose role allows must pass rbac: %v", err)
	}
	if err := call(Caller{Name: "k-nobody"}, map[string]any{}); err == nil || !strings.Contains(err.Error(), "no assigned role") {
		t.Fatalf("with roles assigned, an unassigned caller is refused: %v", err)
	}
	if err := call(Caller{Name: "k-shut"}, map[string]any{"actor": "k-open"}); err == nil || !strings.Contains(err.Error(), `rbac: caller "k-shut"`) {
		t.Fatalf("the args' actor decides nothing -- the forgery path: %v", err)
	}
	if err := call(Caller{Name: "anyone", Service: true}, map[string]any{}); err != nil && strings.Contains(err.Error(), "rbac:") {
		t.Fatalf("the glass is the operator's hand and is never judged by rbac: %v", err)
	}
	tn, _ := tr.Resolve("t")
	if tn.RBACOpen() {
		t.Fatal("roles are assigned here; this tenant is not open")
	}

	// OPEN MODE IS SAID. A bare tenant assigns nothing: every caller passes,
	// the boot line names it, and a parked hold carries it.
	bare := t.TempDir()
	tr2 := tenant.NewRegistry()
	if err := tr2.Add("bare", bare); err != nil {
		t.Fatal(err)
	}
	if err := tr2.SetDefault("bare"); err != nil {
		t.Fatal(err)
	}
	reg2 := Build(tr2, Options{HoldWrites: true})
	if _, err := reg2.Call(tr2, "muster", map[string]any{}, Caller{Name: "k-nobody"}); err != nil && strings.Contains(err.Error(), "rbac:") {
		t.Fatalf("open mode passes everyone: %v", err)
	}
	if line := RBACLine(tr2); !strings.Contains(line, "open (no roles assigned) on: bare") {
		t.Fatalf("open mode must be said by name: %q", line)
	}
	if line := RBACLine(tr); line != "" {
		t.Fatalf("a registry whose tenants all assign roles says nothing: %q", line)
	}
	out, err := reg2.Call(tr2, "git_commit", map[string]any{"message": "x"}, Caller{Name: "k-nobody"})
	if err != nil || !strings.HasPrefix(out, "HELD:") {
		t.Fatalf("a writing call from a key is held: %q %v", out, err)
	}
	for _, h := range reg2.held {
		if h.RBAC != "open (no roles assigned)" {
			t.Fatalf("the hold record must say the rbac state: %+v", h)
		}
	}
}

// P0-14 (2026-09-25): a forbidden verb is never free.
func TestAForbiddenVerbIsNeverFree(t *testing.T) {
	if ForbiddenWord("git_commit") != "commit" || ForbiddenWord("commitment") != "" ||
		ForbiddenWord("PUSH_it") != "push" || ForbiddenWord("hold_answer") != "" {
		t.Fatal("the verb is matched as a word, never as a substring")
	}
	home := t.TempDir()
	tr := tenant.NewRegistry()
	if err := tr.Add("t", home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	reg := Build(tr, Options{})
	carried := 0
	for _, tool := range reg.All() {
		v := ForbiddenWord(tool.Name)
		if v == "" {
			continue
		}
		carried++
		if !tool.Writes {
			t.Errorf("%q carries %q and does not declare Writes -- a forbidden verb, free", tool.Name, v)
		}
		if HeldExempt(tool.Name) {
			t.Errorf("%q carries %q and is exempt from the holds", tool.Name, v)
		}
	}
	if carried == 0 {
		t.Fatal("no tool carries a forbidden verb -- this test would be the vacuous one it replaced")
	}
}

// THE ROLE MODEL (2026-09-26): a SHIPPED role assigned to a key decides by
// what the tool declares. Until this day every shipped role denied every tool
// -- `DefaultPolicy` spoke in kinds (`read`, `edit`, `bash`, `net`, `tools`)
// while `rbac.Can` looked up tool names and `*` -- and no stroke had ever
// assigned one. This one assigns all four through a real registry and calls a
// reader and a writer as each. The door hands `Tool.Writes` down with the name;
// that is the wire, and passing `false` there lets the agent write.
func TestAShippedRoleAssignedToAKeyDecidesByWhatTheToolDeclares(t *testing.T) {
	home := t.TempDir()
	// Only the assignments: the roles are the shipped ones, merged beneath.
	if err := os.WriteFile(filepath.Join(home, "rbac.json"), []byte(
		`{"assign": {"k-op": "operator", "k-steward": "steward", "k-agent": "agent", "k-guest": "guest"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	tr := tenant.NewRegistry()
	if err := tr.Add("t", home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ATLAS_BIN", filepath.Join(home, "NO-SUCH-SPINE.exe"))
	reg := Build(tr, Options{})
	rbacErr := func(key, tool string, args map[string]any) string {
		_, err := reg.Call(tr, tool, args, Caller{Name: key})
		if err == nil || !strings.Contains(err.Error(), "rbac:") {
			return ""
		}
		return err.Error()
	}
	reads := func(key string) string { return rbacErr(key, "muster", map[string]any{}) }
	writes := func(key string) string {
		return rbacErr(key, "remember", map[string]any{"text": "a line from " + key})
	}
	for _, key := range []string{"k-op", "k-steward"} {
		if why := reads(key); why != "" {
			t.Fatalf("%s may read and was refused: %s", key, why)
		}
		if why := writes(key); why != "" {
			t.Fatalf("%s may write and was refused: %s", key, why)
		}
	}
	if why := reads("k-agent"); why != "" {
		t.Fatalf("agent may read and was refused: %s", why)
	}
	why := writes("k-agent")
	for _, want := range []string{`rbac: caller "k-agent" (role "agent") denied tool "remember"`, "remember writes"} {
		if !strings.Contains(why, want) {
			t.Fatalf("agent's write must be refused by kind, saying %q: %q", want, why)
		}
	}
	why = reads("k-guest")
	for _, want := range []string{`(role "guest")`, "denies tools"} {
		if !strings.Contains(why, want) {
			t.Fatalf("guest must be refused calling at all, saying %q: %q", want, why)
		}
	}
	if why := writes("k-guest"); why == "" {
		t.Fatal("guest wrote")
	}
	// The two that may write did, and the two refusals wrote nothing.
	body, err := os.ReadFile(filepath.Join(home, "state", "remembered.jsonl"))
	if err != nil {
		t.Fatal("operator and steward may write and nothing landed:", err)
	}
	if n := strings.Count(strings.TrimRight(string(body), "\n"), "\n") + 1; n != 2 {
		t.Fatalf("remembered.jsonl holds %d lines; the two allowed writes and no other", n)
	}
	if strings.Contains(string(body), "k-agent") || strings.Contains(string(body), "k-guest") {
		t.Fatal("a refused write landed")
	}

	// tenant_rbac_check asks the same question with the same declaration, and
	// refuses a tool the door does not carry rather than judging a word.
	ask := func(actor, tool string) (string, error) {
		return reg.Call(tr, "tenant_rbac_check", map[string]any{"actor": actor, "tool": tool},
			Caller{Name: "glass", Service: true})
	}
	if out, err := ask("k-agent", "remember"); err != nil || !strings.HasPrefix(out, "DENIED") ||
		!strings.Contains(out, "remember writes") {
		t.Fatalf("the check must deny the agent's write by kind: %q %v", out, err)
	}
	if out, err := ask("k-agent", "muster"); err != nil || !strings.HasPrefix(out, "ALLOWED") {
		t.Fatalf("the check must allow the agent's read: %q %v", out, err)
	}
	if _, err := ask("k-agent", "no_such_tool"); !errors.Is(err, ErrUnknownTool) {
		t.Fatalf("a tool the door does not carry is refused by name, not judged: %v", err)
	}
}

// --- a reading action of a writing tool (2026-09-28) ---------------------------

// A READING ACTION IS DECLARED ONLY WHERE IT CAN BE READ. `Reads` narrows a
// writing tool by its `action`; on a tool that does not write, or takes no
// action, it would be a declaration connected to nothing. And the two tools
// that earned it keep it.
func TestAReadingActionIsDeclaredOnlyWhereItCanBeRead(t *testing.T) {
	reg := Build(tenant.NewRegistry(), Options{})
	declared := 0
	for _, tool := range reg.All() {
		if len(tool.Reads) == 0 {
			continue
		}
		declared++
		if !tool.Writes {
			t.Errorf("%s declares reading actions and does not write", tool.Name)
		}
		if !hasArg(tool.Args, "action?") && !hasArg(tool.Args, "action") {
			t.Errorf("%s declares reading actions and takes no action", tool.Name)
		}
		for _, a := range tool.Reads {
			if a == "" || a != strings.ToLower(strings.TrimSpace(a)) {
				t.Errorf("%s declares %q; a reading action is one lower-case word", tool.Name, a)
			}
		}
	}
	for _, name := range []string{"git_tag", "git_branch"} {
		tool, ok := reg.Get(name)
		if !ok || !hasArg(tool.Reads, "list") {
			t.Errorf("%s no longer declares list as a reading action", name)
		}
	}
	if declared == 0 {
		t.Fatal("no tool declares a reading action; this stroke would be vacuous")
	}
	// WritesFor, on the declaration alone: the default and the declared word
	// read; anything else writes; a tool with no reading actions is what it says.
	tag, _ := reg.Get("git_tag")
	for _, args := range []map[string]any{nil, {}, {"action": ""}, {"action": "list"}, {"action": " LIST "}} {
		if tag.WritesFor(args) {
			t.Errorf("git_tag %v must read", args)
		}
	}
	for _, a := range []string{"cut", "send", "remove", "new", "nonsense"} {
		if !tag.WritesFor(map[string]any{"action": a}) {
			t.Errorf("git_tag %s must write", a)
		}
	}
	commit, _ := reg.Get("git_commit")
	if !commit.WritesFor(map[string]any{"action": "list"}) {
		t.Error("a tool with no reading actions writes whatever the action says")
	}
	muster, _ := reg.Get("muster")
	if muster.WritesFor(map[string]any{"action": "cut"}) {
		t.Error("a reading tool reads whatever the action says")
	}
}

// A READING ACTION OF A WRITING TOOL IS NOT HELD, AND IS A READ TO RBAC.
// version-tag's read step was parked on 2026-09-26: `git_tag list` from the
// council's key, held because Writes is one flag for the whole tool. Measured
// here on a real repository with holds armed: the list answers in three
// spellings, the cut and the open are parked and land nothing, the glass is
// never held -- and a role that may not write lists the marks and is refused
// the cut by kind, before any hold.
func TestAReadingActionOfAWritingToolIsNotHeld(t *testing.T) {
	tn := tempWorld(t)
	t.Setenv("MANJUEL_GIT_REMOTE", "")
	t.Setenv("CHAINKIT_GIT_REMOTE", "")
	t.Setenv("ATLAS_BIN", filepath.Join(tn.Home, "NO-SUCH-SPINE.exe"))
	tr := tenant.NewRegistry()
	if err := tr.Add("t", tn.Home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	reg := Build(tr, Options{HoldWrites: true})
	key := Caller{Name: "k-council"}

	for _, args := range []map[string]any{{}, {"action": "list"}, {"action": " LIST "}} {
		out, err := reg.Call(tr, "git_tag", args, key)
		if err != nil {
			t.Fatalf("git_tag %v errored: %v", args, err)
		}
		if strings.HasPrefix(out, "HELD:") {
			t.Fatalf("git_tag %v was held: %s", args, out)
		}
		if !strings.Contains(out, `"tags"`) {
			t.Fatalf("git_tag %v did not list the marks: %s", args, out)
		}
	}
	out, err := reg.Call(tr, "git_branch", map[string]any{"action": "list"}, key)
	if err != nil || strings.HasPrefix(out, "HELD:") || !strings.Contains(out, `"branches"`) {
		t.Fatalf("git_branch list must answer, not park: %q %v", out, err)
	}

	// The writing actions are parked as before, and nothing landed.
	out, err = reg.Call(tr, "git_tag", map[string]any{"action": "cut", "name": "v0.1.5", "message": "m"}, key)
	if err != nil || !strings.HasPrefix(out, "HELD:") {
		t.Fatalf("git_tag cut from a key must be held: %q %v", out, err)
	}
	out, err = reg.Call(tr, "git_branch", map[string]any{"action": "new", "name": "spur"}, key)
	if err != nil || !strings.HasPrefix(out, "HELD:") {
		t.Fatalf("git_branch new from a key must be held: %q %v", out, err)
	}
	if _, err := gitRun(tn, 10*time.Second, "rev-parse", "--verify", "--quiet", "refs/tags/v0.1.5"); err == nil {
		t.Fatal("the held cut landed")
	}
	if _, err := gitRun(tn, 10*time.Second, "rev-parse", "--verify", "--quiet", "refs/heads/spur"); err == nil {
		t.Fatal("the held open landed")
	}
	// The glass is never held, either way.
	out, err = reg.Call(tr, "git_tag", map[string]any{"action": "list"}, Caller{Name: "glass", Service: true})
	if err != nil || !strings.Contains(out, `"tags"`) {
		t.Fatalf("the glass lists: %q %v", out, err)
	}

	// AND THE KIND FOLLOWS THE ACTION. A role that may not write (the shipped
	// agent) lists the marks, and is refused the cut by kind before any hold.
	if err := os.WriteFile(filepath.Join(tn.Home, "rbac.json"),
		[]byte(`{"assign": {"k-agent": "agent"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	tr2 := tenant.NewRegistry()
	if err := tr2.Add("t", tn.Home); err != nil {
		t.Fatal(err)
	}
	if err := tr2.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	reg2 := Build(tr2, Options{HoldWrites: true})
	agent := Caller{Name: "k-agent"}
	out, err = reg2.Call(tr2, "git_tag", map[string]any{"action": "list"}, agent)
	if err != nil || !strings.Contains(out, `"tags"`) {
		t.Fatalf("agent may list the marks: %q %v", out, err)
	}
	_, err = reg2.Call(tr2, "git_tag", map[string]any{"action": "cut", "name": "v0.1.5", "message": "m"}, agent)
	if err == nil || !strings.Contains(err.Error(), `role "agent" denies edit (git_tag writes)`) {
		t.Fatalf("agent's cut is refused by kind, before the hold: %v", err)
	}
}

// --- a key's scope (2026-09-28) -----------------------------------------------

// A KEY'S SCOPE IS MOVED BY A PROVED HAND, ONTO CARRIED GROUND, OR HELD. The
// verb re-proves possession like create and revoke, refuses a tenant the door
// does not carry by name, refuses an unknown id, and from any hand but the
// glass parks in the holds -- where the operator's approval runs exactly the
// call that was parked, and the same plaintext verifies wider afterwards.
func TestAKeysScopeIsMovedByAProvedHandOntoCarriedGroundOrHeld(t *testing.T) {
	home, other := t.TempDir(), t.TempDir()
	key, rec, err := auth.Create(home, "council", []string{"t"}, "test")
	if err != nil {
		t.Fatal(err)
	}
	tr := tenant.NewRegistry()
	for name, h := range map[string]string{"t": home, "u": other} {
		if err := tr.Add(name, h); err != nil {
			t.Fatal(err)
		}
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ATLAS_BIN", filepath.Join(home, "NO-SUCH-SPINE.exe"))
	reg := Build(tr, Options{HoldWrites: true})
	glass := Caller{Name: "glass", Service: true}
	scope := func(c Caller, args map[string]any) (string, error) {
		return reg.Call(tr, "auth_key_scope", args, c)
	}

	// The refusals, as the glass: no re-proof, a stranger tenant, an unknown id.
	if _, err := scope(glass, map[string]any{"id": rec.ID, "tenants": "t,u"}); err == nil ||
		!strings.Contains(err.Error(), "live key") {
		t.Fatalf("a move with no live key named is refused (re-proof): %v", err)
	}
	if _, err := scope(glass, map[string]any{"id": rec.ID, "tenants": "t,stranger", "key": key}); err == nil ||
		!strings.Contains(err.Error(), `"stranger" is not a carried tenant`) {
		t.Fatalf("a tenant the door does not carry is refused by name: %v", err)
	}
	if _, err := scope(glass, map[string]any{"id": "k-00000000", "tenants": "t,u", "key": key}); err == nil ||
		!strings.Contains(err.Error(), `no live key "k-00000000"`) {
		t.Fatalf("an unknown id is refused by id: %v", err)
	}
	if _, err := scope(glass, map[string]any{"id": rec.ID, "key": key}); err == nil ||
		!strings.Contains(err.Error(), "needs the tenants") {
		t.Fatalf("no tenants named is refused: %v", err)
	}
	if got, _ := auth.Verify(home, key); strings.Join(got.Tenants, ",") != "t" {
		t.Fatalf("a refused move must move nothing: %+v", got)
	}

	// From a key, it parks; the approval runs exactly the parked call.
	out, err := scope(Caller{Name: rec.ID}, map[string]any{"id": rec.ID, "tenants": "t,u", "key": key})
	if err != nil || !strings.HasPrefix(out, "HELD:") {
		t.Fatalf("a move from a key must be held for the operator: %q %v", out, err)
	}
	if got, _ := auth.Verify(home, key); strings.Join(got.Tenants, ",") != "t" {
		t.Fatalf("a held move must move nothing yet: %+v", got)
	}
	var id string
	for h := range reg.held {
		id = h
	}
	out, err = reg.Call(tr, "hold_answer", map[string]any{"id": id, "decision": "approve"}, glass)
	if err != nil || !strings.Contains(out, "Approved") || !strings.Contains(out, "SCOPED "+rec.ID) {
		t.Fatalf("the approved hold must run the move: %q %v", out, err)
	}
	got, err := auth.Verify(home, key)
	if err != nil || strings.Join(got.Tenants, ",") != "t,u" {
		t.Fatalf("the same plaintext must carry both grounds now: %+v %v", got, err)
	}
	// And the door's own gate reads it: the widened credential carries u.
	if !auth.ScopeOK(append([]string{"t"}, got.Tenants...), "u") {
		t.Fatal("the widened key does not carry the second ground")
	}
	// As the glass, with re-proof, it moves at once and narrows the same way.
	out, err = scope(glass, map[string]any{"id": rec.ID, "tenants": "u", "key": key})
	if err != nil || !strings.HasPrefix(out, "SCOPED") {
		t.Fatalf("the glass with a live key moves the scope directly: %q %v", out, err)
	}
	if got, _ := auth.Verify(home, key); strings.Join(got.Tenants, ",") != "u" {
		t.Fatalf("narrowing replaces the list whole: %+v", got)
	}
}

// --- a secret in the hold queue (2026-09-28) ------------------------------------

// A SECRET ARGUMENT IS DECLARED WHEREVER A KEY IS TAKEN. Every tool that takes
// `key` says so in Secrets, every Secrets entry names an argument the tool
// takes, and `shown` withholds by the declaration first and by the door's own
// key shape second -- even on a tool the door no longer carries.
func TestASecretArgumentIsDeclaredWhereverAKeyIsTaken(t *testing.T) {
	reg := Build(tenant.NewRegistry(), Options{})
	declared := 0
	for _, tool := range reg.All() {
		takes := map[string]bool{}
		for _, a := range tool.Args {
			takes[strings.TrimRight(a, "?")] = true
		}
		for _, s := range tool.Secrets {
			declared++
			if !takes[s] {
				t.Errorf("%s declares %q secret and takes no such argument", tool.Name, s)
			}
		}
		if takes["key"] && !hasArg(tool.Secrets, "key") {
			t.Errorf("%s takes a key and does not declare it secret", tool.Name)
		}
	}
	if declared == 0 {
		t.Fatal("no tool declares a secret; this stroke would be vacuous")
	}
	shaped := "atl_" + strings.Repeat("ab", 16)
	// By the declaration: a value of any shape, in a declared argument.
	scope, _ := reg.Get("auth_key_scope")
	got := shown(scope, true, map[string]any{"id": "k-1", "tenants": "t", "key": "not-a-key-shape"})
	if got["key"] != Withheld || got["id"] != "k-1" || got["tenants"] != "t" {
		t.Fatalf("a declared secret is withheld and nothing else is: %v", got)
	}
	// By the shape: a key where nobody declared one.
	rem, _ := reg.Get("remember")
	got = shown(rem, true, map[string]any{"text": shaped, "project": "t"})
	if got["text"] != Withheld || got["project"] != "t" {
		t.Fatalf("a key-shaped value is withheld wherever it sits: %v", got)
	}
	// And on a tool the door does not carry, the shape still holds.
	got = shown(Tool{}, false, map[string]any{"key": shaped, "note": "plain"})
	if got["key"] != Withheld || got["note"] != "plain" {
		t.Fatalf("with no declaration the shape still withholds: %v", got)
	}
	// A plain word in a declared-elsewhere name is not a secret here.
	got = shown(rem, true, map[string]any{"key": "the key to the shed"})
	if got["key"] != "the key to the shed" {
		t.Fatalf("a plain value on an undeclared tool is shown: %v", got)
	}
}

// A PARKED SECRET IS WITHHELD WHERE IT IS SHOWN AND RUNS WHOLE. A scope move
// from a key parks with the re-proof key in its arguments; the queue shows the
// call without it, the holds record never carried it, and approving runs the
// parked call with the real key.
func TestAParkedSecretIsWithheldWhereItIsShownAndRunsWhole(t *testing.T) {
	home, other := t.TempDir(), t.TempDir()
	key, rec, err := auth.Create(home, "council", []string{"t"}, "test")
	if err != nil {
		t.Fatal(err)
	}
	tr := tenant.NewRegistry()
	for name, h := range map[string]string{"t": home, "u": other} {
		if err := tr.Add(name, h); err != nil {
			t.Fatal(err)
		}
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ATLAS_BIN", filepath.Join(home, "NO-SUCH-SPINE.exe"))
	reg := Build(tr, Options{HoldWrites: true})
	glass := Caller{Name: "glass", Service: true}

	out, err := reg.Call(tr, "auth_key_scope",
		map[string]any{"id": rec.ID, "tenants": "t,u", "key": key}, Caller{Name: rec.ID})
	if err != nil || !strings.HasPrefix(out, "HELD:") {
		t.Fatalf("the move from a key must park: %q %v", out, err)
	}
	if strings.Contains(out, key) {
		t.Fatal("the held answer carries the key")
	}
	list, err := reg.Call(tr, "hold_list", map[string]any{}, glass)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(list, key) {
		t.Fatalf("the plaintext key is shown on the queue:\n%s", list)
	}
	if !strings.Contains(list, Withheld) || !strings.Contains(list, `"tenants": "t,u"`) ||
		!strings.Contains(list, rec.ID) {
		t.Fatalf("the queue must show the call with its secret withheld and the rest whole:\n%s", list)
	}
	raw, err := os.ReadFile(filepath.Join(home, "state", "holds.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), key) {
		t.Fatal("the holds record carries the key")
	}
	var id string
	for h := range reg.held {
		id = h
	}
	out, err = reg.Call(tr, "hold_answer", map[string]any{"id": id, "decision": "approve"}, glass)
	if err != nil || !strings.Contains(out, "SCOPED "+rec.ID) {
		t.Fatalf("approving must run the parked call whole, with the real key: %q %v", out, err)
	}
	if got, _ := auth.Verify(home, key); strings.Join(got.Tenants, ",") != "t,u" {
		t.Fatalf("the parked call ran with its key and moved the scope: %+v", got)
	}
}

// --- the operator's hand carried past a crossed gate (2026-09-28) ---------------

// A CROSSED GATE CARRIES THE HAND TO THE TOOLS IT GRANTS AND NO OTHER. With the
// hand raised over git_tag for a world, the council's cut reaches the tool
// (which speaks for itself) and the holds record says `crossed`, naming the run
// and the gate; an ungranted writer still parks; the hand is lowered after the
// turn and the same cut parks again. WithHand binds a copy, and a council with
// no hand raises nothing.
func TestACrossedGateCarriesTheHandToTheToolsItGrantsAndNoOther(t *testing.T) {
	tn := tempWorld(t)
	t.Setenv("MANJUEL_GIT_REMOTE", "")
	t.Setenv("CHAINKIT_GIT_REMOTE", "")
	t.Setenv("ATLAS_BIN", filepath.Join(tn.Home, "NO-SUCH-SPINE.exe"))
	tr := tenant.NewRegistry()
	if err := tr.Add("t", tn.Home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	reg := Build(tr, Options{HoldWrites: true})
	key := Caller{Name: "k-council"}
	cut := func() map[string]any {
		return map[string]any{"action": "cut", "name": "v0.1.5", "message": "m"}
	}

	out, err := reg.Call(tr, "git_tag", cut(), key)
	if err != nil || !strings.HasPrefix(out, "HELD:") {
		t.Fatalf("with no hand the cut parks: %q %v", out, err)
	}

	base := council(tn.Home).(councilEngine)
	if len(base.hand.Tools) != 0 {
		t.Fatal("a council is born with no hand")
	}
	bound := base.WithHand("f-20260928-000000-deadbeef", "judge", []string{"git_tag"}).(councilEngine)
	if len(base.hand.Tools) != 0 {
		t.Fatal("WithHand bound the registered engine -- the hand would outlive its segment")
	}
	if strings.Join(bound.hand.Tools, ",") != "git_tag" || bound.hand.Gate != "judge" {
		t.Fatalf("WithHand did not bind the crossing it was handed: %+v", bound.hand)
	}

	err = bound.underHand(func() error {
		out, err := reg.Call(tr, "git_tag", cut(), key)
		if err != nil {
			return err
		}
		if strings.HasPrefix(out, "HELD:") {
			t.Fatalf("under the hand the cut parked: %s", out)
		}
		if !strings.HasPrefix(out, "Refused:") {
			t.Fatalf("under the hand the tool itself must answer (here, refusing on its own law): %s", out)
		}
		out, err = reg.Call(tr, "remember", map[string]any{"text": "x"}, key)
		if err != nil || !strings.HasPrefix(out, "HELD:") {
			t.Fatalf("an ungranted writer must still park under the hand: %q %v", out, err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// Lowered after the turn: the same cut parks again.
	out, err = reg.Call(tr, "git_tag", cut(), key)
	if err != nil || !strings.HasPrefix(out, "HELD:") {
		t.Fatalf("after the turn the hand must be lowered: %q %v", out, err)
	}
	// The record says whose hand it was, and where.
	raw, err := os.ReadFile(filepath.Join(tn.Home, "state", "holds.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	crossed := 0
	for _, ln := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		if !strings.Contains(ln, `"what":"crossed"`) {
			continue
		}
		crossed++
		for _, want := range []string{`"tool":"git_tag"`, `"caller":"k-council"`,
			"run f-20260928-000000-deadbeef gate judge"} {
			if !strings.Contains(ln, want) {
				t.Fatalf("the crossed line must say %s: %s", want, ln)
			}
		}
	}
	if crossed != 1 {
		t.Fatalf("%d crossed lines, wanted exactly one -- the cut that rode the hand", crossed)
	}
	// A council with no hand raises nothing.
	ran := false
	if err := base.underHand(func() error {
		ran = true
		if _, ok := hands.covers(tn.Home, "git_tag"); ok {
			t.Fatal("a handless turn raised a hand")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if !ran {
		t.Fatal("underHand did not run the turn")
	}
}

// A GRANT NAMES A WRITING TOOL THIS DOOR CARRIES, or the save is refused by
// name -- at the save, not at the crossing.
func TestFlowSaveRefusesAGrantOnAStrangerOrAReader(t *testing.T) {
	home := t.TempDir()
	tr := tenant.NewRegistry()
	if err := tr.Add("t", home); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ATLAS_BIN", filepath.Join(home, "NO-SUCH-SPINE.exe"))
	reg := Build(tr, Options{})
	glass := Caller{Name: "glass", Service: true}
	spec := func(grant string) string {
		return `{"nodes":[{"name":"a","kind":"ask","question":"Q"},` +
			`{"name":"g","kind":"gate","title":"cut?","grants":["` + grant + `"]},` +
			`{"name":"w","kind":"run","question":"cut it"}],` +
			`"edges":[{"from":"a","to":"g"},{"from":"g","to":"w","when":"pass"}]}`
	}
	save := func(name, spec string) (string, error) {
		return reg.Call(tr, "flow_save", map[string]any{"name": name, "spec": spec}, glass)
	}
	if _, err := save("granted", spec("muster")); err == nil || !strings.Contains(err.Error(), `grants "muster", which does not write`) {
		t.Fatalf("a grant on a reader must be refused by name: %v", err)
	}
	if _, err := save("granted", spec("no_such_tool")); err == nil || !strings.Contains(err.Error(), "not a tool this door carries") {
		t.Fatalf("a grant on a stranger must be refused by name: %v", err)
	}
	out, err := save("granted", spec("git_tag"))
	if err != nil || !strings.Contains(out, "SAVED flow granted v1") {
		t.Fatalf("a grant on a writing tool the door carries must fold: %q %v", out, err)
	}
	got, err := reg.Call(tr, "flow_get", map[string]any{"name": "granted"}, glass)
	if err != nil || !strings.Contains(got, `"grants"`) || !strings.Contains(got, `"git_tag"`) {
		t.Fatalf("the folded spec must carry the grant: %q %v", got, err)
	}
	// The flow package's own refusal reaches the door too.
	bad := `{"nodes":[{"name":"w","kind":"run","question":"cut it","grants":["git_tag"]}],"edges":[]}`
	if _, err := save("wrong", bad); err == nil || !strings.Contains(err.Error(), "only a gate") {
		t.Fatalf("a grant on a run node must be refused: %v", err)
	}
}

// --- the standup, fired from the glass (2026-09-28) -----------------------------
//
// His word: "fire the standup through the glass". `standup_run` runs the
// estate's own tests/standup.py in the world, with the python the door runs the
// engine with, and hands back the tally line and the whole of what it printed.
// Its refusals are the door's own, and each is struck here by name; the run
// itself is struck against a stand-in script where a python is on the PATH.
func TestTheStandupRunsInTheWorldAndRefusesByName(t *testing.T) {
	tr := tenant.NewRegistry()
	bare := t.TempDir()
	if err := tr.Add("bare", bare); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("bare"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ATLAS_BIN", filepath.Join(bare, "NO-SUCH-SPINE.exe"))
	glass := Caller{Name: "glass", Service: true}

	// No core command: nothing to run the standup with.
	noCore := Build(tr, Options{})
	if _, err := noCore.Call(tr, "standup_run", map[string]any{}, glass); err == nil ||
		!strings.Contains(err.Error(), "started without --manjuel") {
		t.Fatalf("with no core command the tool must say so: %v", err)
	}
	// A world with no standup of its own.
	reg := Build(tr, Options{CoreCmd: "python"})
	if _, err := reg.Call(tr, "standup_run", map[string]any{}, glass); err == nil ||
		!strings.Contains(err.Error(), "carries no tests/standup.py") {
		t.Fatalf("a world without the standup must be refused by name: %v", err)
	}
	// A set that is not one of the three, refused before anything is looked up.
	if _, err := reg.Call(tr, "standup_run", map[string]any{"set": "evening"}, glass); err == nil ||
		!strings.Contains(err.Error(), `"evening" is none of them`) {
		t.Fatalf("an unknown set must be refused by name: %v", err)
	}

	// A world that carries a standup: a stand-in script that prints what the
	// real one prints and exits 1, the way a missed expectation does.
	script := "import sys\n" +
		"print('  manjuel -- the standup')\n" +
		"print('  ' + ' '.join(sys.argv[1:]) if len(sys.argv) > 1 else '  (morning)')\n" +
		"print('  8/9 cases met their expectations.  report: logs/standup_x.md')\n" +
		"sys.exit(1)\n"
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "tests", "standup.py"), []byte(script), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := tr.Add("w", home); err != nil {
		t.Fatal(err)
	}
	// An open sitting on the world refuses the run: the standup opens its own.
	os.MkdirAll(filepath.Join(home, "sessions"), 0o755)
	os.WriteFile(filepath.Join(home, "sessions", "sessions.jsonl"),
		[]byte(`{"n": 7, "started": "2026-09-28T09:00:00", "pid": 1}`+"\n"), 0o644)
	if _, err := reg.Call(tr, "standup_run", map[string]any{"project": "w"}, glass); err == nil ||
		!strings.Contains(err.Error(), "has an open sitting (7, opened 2026-09-28T09:00:00)") {
		t.Fatalf("an open sitting must refuse the standup by name: %v", err)
	}
	os.WriteFile(filepath.Join(home, "sessions", "sessions.jsonl"),
		[]byte(`{"n": 7, "started": "2026-09-28T09:00:00", "ended": "2026-09-28T09:10:00"}`+"\n"), 0o644)

	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("no python on the PATH to run the stand-in with")
	}
	out, err := reg.Call(tr, "standup_run", map[string]any{"project": "w", "set": "Court"}, glass)
	if err != nil {
		t.Fatalf("the standup must answer, exit code and all: %v", err)
	}
	head := strings.SplitN(out, "\n", 2)[0]
	if head != `STANDUP court on "w" · exit status 1 · 8/9 cases met their expectations. · report logs/standup_x.md` {
		t.Fatalf("the head must carry the set, the exit, the tally and the report: %q", head)
	}
	if !strings.Contains(out, "  --court\n") || !strings.Contains(out, "manjuel -- the standup") {
		t.Fatalf("the set's flag must reach the script and its words must come back whole:\n%s", out)
	}
	// The morning set is the bare command, and its head says exit 0.
	script0 := "print('  9/9 cases met their expectations.  report: logs/standup_y.md')\n"
	os.WriteFile(filepath.Join(home, "tests", "standup.py"), []byte(script0), 0o644)
	out, err = reg.Call(tr, "standup_run", map[string]any{"project": "w"}, glass)
	if err != nil || !strings.HasPrefix(out, `STANDUP morning on "w" · exit 0 · 9/9 cases met their expectations. · report logs/standup_y.md`) {
		t.Fatalf("the morning set must run bare and report exit 0: %q %v", out, err)
	}
	// A script that prints no tally says so rather than inventing one.
	os.WriteFile(filepath.Join(home, "tests", "standup.py"), []byte("print('the rack is unreachable')\nraise SystemExit(2)\n"), 0o644)
	out, err = reg.Call(tr, "standup_run", map[string]any{"project": "w"}, glass)
	if err != nil || !strings.Contains(out, "the script printed no tally -- read what it said below") ||
		!strings.Contains(out, "the rack is unreachable") {
		t.Fatalf("no tally must be said out loud, with the script's words: %q %v", out, err)
	}
}

// --- the suites, from the glass (2026-09-28) ------------------------------------
//
// His ruling: "if its on the glass, and the record matches, id call it proof."
// `suite_run` runs the world's own strokes and smoke, one after the other,
// with the python the door runs the engine with; the suites stamp their own
// proof and the head is read back off that stamp. Its refusals are struck by
// name; the run itself against stand-in scripts where a python is on the PATH.
func TestTheSuitesRunFromTheGlassAndTheHeadIsReadOffTheStamp(t *testing.T) {
	tr := tenant.NewRegistry()
	bare := t.TempDir()
	if err := tr.Add("bare", bare); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("bare"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ATLAS_BIN", filepath.Join(bare, "NO-SUCH-SPINE.exe"))
	glass := Caller{Name: "glass", Service: true}

	noCore := Build(tr, Options{})
	if _, err := noCore.Call(tr, "suite_run", map[string]any{}, glass); err == nil ||
		!strings.Contains(err.Error(), "started without --manjuel") {
		t.Fatalf("with no core command the tool must say so: %v", err)
	}
	reg := Build(tr, Options{CoreCmd: "python"})
	if _, err := reg.Call(tr, "suite_run", map[string]any{}, glass); err == nil ||
		!strings.Contains(err.Error(), "carries no tests/test_manjuel.py") {
		t.Fatalf("a world without the suites must be refused by name: %v", err)
	}
	if _, err := reg.Call(tr, "suite_run", map[string]any{"set": "parity"}, glass); err == nil ||
		!strings.Contains(err.Error(), `"parity" is none of them`) {
		t.Fatalf("an unknown set must be refused by name: %v", err)
	}
	// A reader in the door's eyes, so the coder's loop can ask it without a hand.
	if tool, ok := reg.Get("suite_run"); !ok || tool.Writes {
		t.Fatal("suite_run must be declared a reader: it writes the suites' own stamps and nothing of the work")
	}

	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("no python on the PATH to run the stand-ins with")
	}
	// A world that carries the suites: stand-ins that stamp the way the real
	// ones do -- `running` first, the tally after -- and print what they print.
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	stamp := func(suite string, passed, total int, green bool, extra string) string {
		return "import json, os, time\n" +
			"p = os.path.join('tests', 'last_run.json')\n" +
			"book = json.load(open(p)) if os.path.exists(p) else {}\n" +
			"book['" + suite + "'] = {'state': 'running', 'at': time.time(), 'passed': None, 'total': None, 'green': False}\n" +
			"json.dump(book, open(p, 'w'))\n" +
			extra +
			"book['" + suite + "'] = {'state': 'finished', 'at': time.time(), 'passed': " + fmt.Sprint(passed) +
			", 'total': " + fmt.Sprint(total) + ", 'green': " + map[bool]string{true: "True", false: "False"}[green] + "}\n" +
			"json.dump(book, open(p, 'w'))\n" +
			"open(os.path.join('tests', '" + suite + ".ran'), 'w').write('ran')\n"
	}
	strokes := stamp("strokes", 12, 12, true, "") +
		"print('    [PASS]  a probe')\n" +
		"print('  12/12 strokes.  PROVEN.')\n"
	smoke := stamp("smoke", 2, 3, false, "") +
		"print('    [FAIL]  a check that went red    the reason')\n" +
		"print('  2/3 checks.  RED')\n" +
		"raise SystemExit(1)\n"
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(home, "tests", name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("test_manjuel.py", strokes)
	write("smoke_cli.py", smoke)
	if err := tr.Add("w", home); err != nil {
		t.Fatal(err)
	}

	out, err := reg.Call(tr, "suite_run", map[string]any{"project": "w"}, glass)
	if err != nil {
		t.Fatalf("both suites must answer, exit codes and all: %v", err)
	}
	head := strings.SplitN(out, "\n", 2)[0]
	if head != `SUITES both on "w" · strokes: 12/12 green · exit 0 · smoke: 2/3 RED · exit status 1` {
		t.Fatalf("the head must be read off the stamp, suite by suite: %q", head)
	}
	for _, needle := range []string{"--- strokes (tests/test_manjuel.py) ---", "12/12 strokes.  PROVEN.",
		"--- smoke (tests/smoke_cli.py) ---", "[FAIL]  a check that went red", "2/3 checks.  RED"} {
		if !strings.Contains(out, needle) {
			t.Fatalf("the body must carry each suite's failures and tally; lacks %q:\n%s", needle, out)
		}
	}
	if strings.Contains(out, "[PASS]") {
		t.Fatalf("a pass is a count, not a line in the answer:\n%s", out)
	}
	if n := strings.Count(out, "[FAIL]  a check that went red"); n != 1 {
		t.Fatalf("a failure is said once, got %d times:\n%s", n, out)
	}
	// The order is the gate's: the strokes ran before the smoke.
	if i, j := strings.Index(out, "--- strokes"), strings.Index(out, "--- smoke"); i < 0 || j < 0 || j < i {
		t.Fatalf("the strokes must run before the smoke:\n%s", out)
	}

	// ONE SET RUNS ONE SUITE: `smoke` leaves the strokes untouched.
	os.Remove(filepath.Join(home, "tests", "strokes.ran"))
	os.Remove(filepath.Join(home, "tests", "smoke.ran"))
	out, err = reg.Call(tr, "suite_run", map[string]any{"project": "w", "set": "Smoke"}, glass)
	if err != nil || !strings.HasPrefix(out, `SUITES smoke on "w" · smoke: 2/3 RED · exit status 1`) {
		t.Fatalf("the smoke set must run the smoke alone: %q %v", out, err)
	}
	if _, err := os.Stat(filepath.Join(home, "tests", "strokes.ran")); err == nil {
		t.Fatal("the smoke set ran the strokes")
	}

	// A SUITE THAT NEVER FINISHED is said so, off its own stamp.
	write("test_manjuel.py", "import json, os, time\n"+
		"p = os.path.join('tests', 'last_run.json')\n"+
		"book = json.load(open(p)) if os.path.exists(p) else {}\n"+
		"book['strokes'] = {'state': 'running', 'at': time.time(), 'passed': None, 'total': None, 'green': False}\n"+
		"json.dump(book, open(p, 'w'))\n"+
		"raise SystemExit(2)\n")
	out, err = reg.Call(tr, "suite_run", map[string]any{"project": "w", "set": "strokes"}, glass)
	if err != nil || !strings.Contains(out, "strokes: DID NOT FINISH (the stamp still says running) · exit status 2") {
		t.Fatalf("a crashed suite must be said off its stamp: %q %v", out, err)
	}

	// ONE AT A TIME, MACHINE-WIDE: a second call while one runs is refused.
	write("test_manjuel.py", "import time\ntime.sleep(3)\nprint('  0/0 strokes.  PROVEN.')\n")
	done := make(chan error, 1)
	go func() {
		_, err := reg.Call(tr, "suite_run", map[string]any{"project": "w", "set": "strokes"}, glass)
		done <- err
	}()
	time.Sleep(700 * time.Millisecond)
	if _, err := reg.Call(tr, "suite_run", map[string]any{"project": "w", "set": "strokes"}, glass); err == nil ||
		!strings.Contains(err.Error(), "already running") {
		t.Fatalf("a second suite while one runs must be refused by name: %v", err)
	}
	if err := <-done; err != nil {
		t.Fatalf("the first run must still finish: %v", err)
	}
}
