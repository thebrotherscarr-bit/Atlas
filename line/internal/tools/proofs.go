package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"atlas/line/internal/tenant"
)

// proofs reads what a world has actually PROVED, from the world's own record.
//
// WHY THIS EXISTS. The Evals page showed TOTAL 0 · PASSED 0 · FAILED 0 · PASS
// RATE 0% on an estate with ten live standups and 1870 green strokes on disk,
// because it was counting the webapp's own SQLite store -- which nothing
// writes -- instead of asking the record. That is the same fault the dashboard
// carried until P0-11, and its comment in app.js already says the rule:
// "Every row on this page names the tool it was read from. A number the record
// cannot prove is not shown."
//
// WHY JSON AND NOT PROSE. Every other tool here answers in readable lines and
// the glass parses them with regexes. That is fine for a rack ladder and it is
// wrong for this: the one job of an evals page is to not misreport a number,
// and a regex that drifts one character turns 1870 into 187 silently. The
// numbers cross as numbers. The text is still text, and a model reading it
// loses nothing.
//
// THREE FILES, EACH ALLOWED TO BE ABSENT ON ITS OWN. A world that never ran a
// suite is not an error, and one unreadable file must not blank the other two
// -- a proof page that vanishes when one thing is missing is worse than a
// proof page with a gap in it (boot.py's own rule).
//
// READ AGAIN ONLY WHEN THE RECORD CHANGES (2026-09-15). The three ledgers and
// SEAT_LOG go through fromRecord, below. last_run.json is a few hundred bytes
// and is read every call, as before.
func toolProofs(t tenant.Tenant, _ map[string]any) (string, error) {
	out := map[string]any{"world": t.Name}

	// ---- the suites: tests/last_run.json -------------------------------
	// Stamped by the suites themselves. `green` and the counts are theirs;
	// nothing here re-counts or re-judges them.
	if b, err := os.ReadFile(filepath.Join(t.Home, "tests", "last_run.json")); err != nil {
		out["suites_error"] = readable(err)
	} else {
		var suites map[string]any
		if err := json.Unmarshal(b, &suites); err != nil {
			out["suites_error"] = "tests/last_run.json is not readable JSON: " + err.Error()
		} else {
			out["suites"] = suites
		}
	}

	// ---- the standups: tests/run_history.jsonl -------------------------
	// The live harness's own line per run. Newest last, as the file is.
	if runs, err := fromRecord(filepath.Join(t.Home, "tests", "run_history.jsonl"), standupsIn); err != nil {
		out["standups_error"] = readable(err)
	} else {
		out["standups"] = runs
	}

	// ---- the parity: sessions/parity_history.jsonl ---------------------
	if rows, err := fromRecord(filepath.Join(t.Home, "sessions", "parity_history.jsonl"), linesIn); err != nil {
		out["parity_error"] = readable(err)
	} else {
		out["parity"] = rows
	}

	// ---- the record: what this estate has actually done ----------------
	//
	// sessions.jsonl is pure JSON, one object per line, and a CLOSING line
	// supersedes its opening one (seatlog.record's own rule, and the reason a
	// sitting can appear twice). Keyed by n, last line wins -- so these counts
	// are read exactly, with no format rule to get wrong.
	rec := map[string]any{}
	if sat, err := fromRecord(filepath.Join(t.Home, "sessions", "sessions.jsonl"), sittingsIn); err != nil {
		rec["error"] = readable(err)
	} else {
		for k, v := range sat {
			rec[k] = v
		}
	}

	// The tolls as SEAT_LOG carries them. Matched on " — sitting " alone: it is
	// the one part of that heading the record has never varied.
	if n, err := fromRecord(filepath.Join(t.Home, "SEAT_LOG.md"), tollsIn); err == nil {
		rec["seat_log_tolls"] = n
	} else {
		rec["seat_log_tolls"] = nil
	}

	// NOT COUNTED HERE, and the page says so rather than guessing: a memory
	// entry is whatever memory.py's _ENTRY_RE says it is, and the index count
	// is a SELECT against index/vectors.db. A second definition of either in
	// this file would drift from the core's the first time his changed. Both
	// are in boot.report()'s RECORD block, which the boot panel shows whole.
	rec["counted_by_the_engine"] = []string{"memory entries", "index docs and passages"}
	out["record"] = rec

	// ---- the transcripts a run can be read back from -------------------
	if names, err := os.ReadDir(filepath.Join(t.Home, "logs")); err == nil {
		reports := []string{}
		for _, n := range names {
			if !n.IsDir() && strings.HasPrefix(n.Name(), "standup_") {
				reports = append(reports, "logs/"+n.Name())
			}
		}
		sort.Strings(reports)
		out["standup_reports"] = reports
	}

	b, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ---- the record, read again only when it changes ---------------------------
//
// WHY (2026-09-15). The Dashboard asks for proofs on every refresh, every 15
// seconds while it is on screen, and every answer was built from the files
// again: 15ms and 7.6MB of garbage a call, 11ms and 6.0MB of it parsing
// sessions.jsonl -- 846KB, append-only, and changed once a sitting. Between
// sittings every one of those answers was the same answer.
//
// So what a file was made into is kept, and made again only when the file's
// size or modification time has moved: one stat a file, 17µs. An append moves
// the size and a rewrite moves the time.
//
// THREE RULES:
//   - the stat is taken BEFORE the read. A write landing between the two
//     leaves an older key on newer bytes: one extra read next call, never an
//     old answer.
//   - a read that fails is not kept. A file a scanner holds for a moment
//     would otherwise read "unreadable" until it next changed.
//   - what is kept is what the file was made into -- counts, twelve rows, the
//     standups -- never the bytes, and nothing changes it once kept. Every
//     call builds its own reply around it.
type kept struct {
	size int64
	mod  time.Time
	val  any
}

var (
	keptMu sync.Mutex
	keptBy = map[string]kept{}
	// readRecord is os.ReadFile. The strokes swap it to count reads and to
	// fail one.
	readRecord = os.ReadFile
)

// fromRecord returns in(the file's bytes), reading the file again only when
// its size or modification time has moved since the last read.
func fromRecord[T any](path string, in func([]byte) T) (T, error) {
	var zero T
	st, err := os.Stat(path)
	if err != nil {
		return zero, err
	}
	keptMu.Lock()
	k, ok := keptBy[path]
	keptMu.Unlock()
	if v, same := k.val.(T); ok && same && k.size == st.Size() && k.mod.Equal(st.ModTime()) {
		return v, nil
	}
	b, err := readRecord(path)
	if err != nil {
		return zero, err
	}
	v := in(b)
	keptMu.Lock()
	keptBy[path] = kept{size: st.Size(), mod: st.ModTime(), val: v}
	keptMu.Unlock()
	return v, nil
}

// standupsIn is tests/run_history.jsonl's standup lines, newest last.
func standupsIn(raw []byte) []map[string]any {
	runs := []map[string]any{}
	for _, r := range linesIn(raw) {
		if s, _ := r["suite"].(string); s == "standup" {
			runs = append(runs, r)
		}
	}
	return runs
}

// sittingsIn counts sessions.jsonl by sitting, the last line for each n
// winning.
func sittingsIn(raw []byte) map[string]any {
	byN := map[float64]map[string]any{}
	order := []float64{}
	for _, r := range linesIn(raw) {
		n, ok := r["n"].(float64)
		if !ok {
			continue
		}
		if _, seen := byN[n]; !seen {
			order = append(order, n)
		}
		byN[n] = r
	}
	sort.Float64s(order)
	tolled, runs, open := 0, 0, 0
	recent := []map[string]any{}
	for _, n := range order {
		r := byN[n]
		if b, _ := r["toll_paid"].(bool); b {
			tolled++
		}
		rs, _ := r["runs"].([]any)
		runs += len(rs)
		if e, _ := r["ended"].(string); strings.TrimSpace(e) == "" {
			open++
		}
		recent = append(recent, map[string]any{
			"n": n, "started": r["started"], "ended": r["ended"],
			"toll_paid": r["toll_paid"], "runs": len(rs),
		})
	}
	if len(recent) > 12 {
		// Copied, not resliced: this is kept, and a reslice would keep every
		// sitting's row alive behind the twelve shown.
		recent = append([]map[string]any(nil), recent[len(recent)-12:]...)
	}
	return map[string]any{"sittings": len(order), "tolled": tolled, "runs": runs,
		"still_open": open, "recent": recent}
}

// tollsIn counts SEAT_LOG's toll headings.
func tollsIn(raw []byte) int {
	n := 0
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(line, "## ") && strings.Contains(line, " sitting ") {
			n++
		}
	}
	return n
}

// linesIn reads one JSON object per line, skipping blanks and refusing
// nothing: a half-written last line (what a killed process leaves) is dropped
// rather than failing the whole read, because the lines before it are still
// true.
func linesIn(raw []byte) []map[string]any {
	out := []map[string]any{}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var row map[string]any
		if json.Unmarshal([]byte(line), &row) != nil {
			continue
		}
		out = append(out, row)
	}
	return out
}

// readable turns a missing file into a sentence rather than a stack trace: a
// world that never ran a suite has not failed at anything.
func readable(err error) string {
	if os.IsNotExist(err) {
		return "never run in this world"
	}
	return fmt.Sprintf("unreadable: %v", err)
}
