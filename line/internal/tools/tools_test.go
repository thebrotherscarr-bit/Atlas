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

	"atlas/line/internal/engine"
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
