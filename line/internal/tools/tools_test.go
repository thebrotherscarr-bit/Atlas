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
	"os"
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
