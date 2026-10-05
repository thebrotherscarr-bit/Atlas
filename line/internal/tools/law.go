package tools

// THE LAWS, THROUGH THE GLASS (2026-10-05, his words: "set it as a law all the agents read ... THROUGH THE SYSTEM", then "a LAWS page ... Add it.").
//
// The laws live where they always did: law/*.md, bound to the hash chain by law/law.py, with law/LAW_LEDGER.md the one law that is APPENDABLE
// ("sealed as far as it is written, and grows below it"). What did not exist was a way to see how far each of them is sealed without a shell, or to
// set a law or seal the ledger without typing a command. These three tools are that, and nothing more:
//
//	law_status  a reader: every law file, how far the chain seals it, what is draft, and the chain's own verdict (law.py verify)
//	law_add     the operator's hand: append one numbered entry to the ledger, BELOW the seal, as a draft
//	law_seal    the operator's hand: law.py seal on the ledger -- what is written becomes law, and is frozen
//
// THE CHAIN IS APPENDED ONLY THROUGH law.py ("the chain is appended ONLY through this tool, ONLY by links"), so law_seal runs it and writes nothing of its
// own; law_add writes only the ledger's own tail, which law.py's seal says may grow, and checks the chain before and after. THE VERDICT IS law.py's, never
// this file's: the per-file lines below are read off chain.jsonl for display and are NOT a second verifier -- a file whose link no longer matches is
// said to have changed, and the verdict beside it is the walk's.
//
// BOTH WRITERS ARE ServiceOnly: no seat, agent or role is ever handed the laws (a rule a seat could reword is not a law), and neither runs while a sitting
// is open -- SITTING LAW 5, the same wall every writer of this ground keeps. Neither is held on a card: his click on the glass, behind the glass's own lock,
// is the decision, and the page says what the click does before it is made.
//
// WHAT NEITHER DOES: edit a sealed byte, bind a law whole (`direct`), touch CLAUDE.md, or write anywhere but law/LAW_LEDGER.md and the chain through law.py.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"atlas/line/internal/engine"
	"atlas/line/internal/tenant"
)

const (
	lawLedger   = "LAW_LEDGER.md"
	lawChain    = "chain.jsonl"
	lawMaxTitle = 120
	lawMinText  = 20
	lawMaxText  = 8000
	lawMaxFrom  = 600
	lawWait     = 30 * time.Second
)

// lawAnchorRe is law.py's ANCHOR_RE, the first line of every link's doc; the cross-language stroke in the core's suite holds the two to each other.
var lawAnchorRe = regexp.MustCompile(`^(COUNSEL|RULING|DIRECT|REFUSED) by ([A-Za-z0-9_-]+) -> (?:Archive/law/|law/)([A-Za-z0-9._-]+\.md) sha256:([0-9a-f]{64})(?: bytes:([1-9][0-9]*))?$`)

// lawEntryRe numbers the ledger's entries: "## 6. Title".
var lawEntryRe = regexp.MustCompile(`(?m)^## (\d+)\. `)

// lawHeadingRe is a line that would open a new entry or section if it stood in an entry's own text.
var lawHeadingRe = regexp.MustCompile(`(?m)^(##\s|---\s*$)`)

// lawProvenRe is the line law.py verify prints when the chain walks whole.
var lawProvenRe = regexp.MustCompile(`the law chain proves whole: (\d+) links, head ([0-9a-f]+)`)

// lawRun is how law.py is started. A variable so the strokes can stand a stub in front of the real tool (it is the core's, not this repository's) and still
// measure the door's own judgement; the real tool is held against this file's reading of it by the core's own suite.
var lawRun = func(py string, args []string, opts spawnOpts) spawnResult { return spawn(py, args, opts) }

// lawMu keeps the laws one writer at a time: an append and a seal racing is a seal over bytes that moved.
var lawMu sync.Mutex

// ---- what the chain says --------------------------------------------------------------------------

type lawFile struct {
	Name        string `json:"name"`
	Bytes       int64  `json:"bytes"`
	Fingerprint string `json:"fingerprint"`
	// State is one of: sealed, partly sealed, not on the chain, changed since sealed.
	State       string `json:"state"`
	SealedBytes int64  `json:"sealed_bytes,omitempty"`
	DraftBytes  int64  `json:"draft_bytes,omitempty"`
	Kind        string `json:"kind,omitempty"`
	By          string `json:"by,omitempty"`
	Links       []int  `json:"links,omitempty"`
}

type lawAnswer struct {
	World       string    `json:"world"`
	Ok          bool      `json:"ok"`
	Verdict     string    `json:"verdict"`
	Links       int       `json:"links"`
	Head        string    `json:"head,omitempty"`
	Files       []lawFile `json:"files"`
	DraftBytes  int64     `json:"draft_bytes"`
	SittingOpen bool      `json:"sitting_open"`
	// Why is why the two writing tools would refuse right now, in words; empty when they would not.
	Why string `json:"why,omitempty"`
}

type lawLinkRef struct {
	n      int
	kind   string
	by     string
	extent int64 // 0 = the whole file
	token  string
}

// lawLinks reads chain.jsonl: every link per law file, in order, and the head's hash. A line that does not read is skipped here and named by law.py verify.
func lawLinks(t tenant.Tenant) (map[string][]lawLinkRef, string, error) {
	out := map[string][]lawLinkRef{}
	b, err := os.ReadFile(filepath.Join(t.Home, "law", lawChain))
	if err != nil {
		if os.IsNotExist(err) {
			return out, "", nil
		}
		return nil, "", err
	}
	head := ""
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var row struct {
			Kind    string `json:"kind"`
			Hash    string `json:"hash"`
			Payload struct {
				N   int    `json:"n"`
				Doc string `json:"doc"`
			} `json:"payload"`
		}
		if json.Unmarshal([]byte(line), &row) != nil {
			continue
		}
		head = row.Hash
		if row.Kind != "link" {
			continue
		}
		m := lawAnchorRe.FindStringSubmatch(strings.TrimSpace(strings.SplitN(row.Payload.Doc, "\n", 2)[0]))
		if m == nil {
			continue
		}
		ref := lawLinkRef{n: row.Payload.N, kind: m[1], by: m[2], token: m[4]}
		if m[5] != "" {
			ref.extent, _ = strconv.ParseInt(m[5], 10, 64)
		}
		out[m[3]] = append(out[m[3]], ref)
	}
	return out, head, nil
}

func lawSHA(path string, limit int64) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()
	h := sha256.New()
	if limit > 0 {
		if n, err := io.CopyN(h, f, limit); err != nil || n != limit {
			return "", false // a seal over bytes that are gone cannot hold
		}
	} else if _, err := io.Copy(h, f); err != nil {
		return "", false
	}
	return hex.EncodeToString(h.Sum(nil)), true
}

// lawVerify runs law.py verify: the walk's own words, and whether it proved the chain whole.
func lawVerify(t tenant.Tenant, opts Options) (bool, string) {
	script := filepath.Join(t.Home, "law", "law.py")
	if st, err := os.Stat(script); err != nil || st.IsDir() {
		return false, "law/law.py is not in this world: the chain cannot be walked here"
	}
	res := lawRun(pythonWord(opts.CoreCmd), []string{script, "verify"}, spawnOpts{Dir: t.Home, Timeout: lawWait})
	words := strings.TrimSpace(strings.ReplaceAll(res.Combined, "\r\n", "\n"))
	if words == "" && res.Err != nil {
		words = "law.py did not answer: " + res.Err.Error()
	}
	return res.Err == nil && lawProvenRe.MatchString(words), words
}

// lawWriteRefusal is why a writing tool may not run now, or "". SITTING LAW 5 is the whole of it, plus the tool the chain is appended through.
func lawWriteRefusal(t tenant.Tenant) string {
	if e, open := engines.Get(t.Home); open {
		return fmt.Sprintf("an engine is open on %q (sitting %s): nothing in this ground is edited while the operator's sitting is open (SITTING LAW 5) -- close it first", t.Name, e.Opened().Str("sitting"))
	}
	if n, started, sat := engine.SittingOpen(t.Home); sat {
		return fmt.Sprintf("%q has an open sitting (%d, opened %s): nothing in this ground is edited while it is open (SITTING LAW 5)", t.Name, int(n), started)
	}
	if st, err := os.Stat(filepath.Join(t.Home, "law", "law.py")); err != nil || st.IsDir() {
		return fmt.Sprintf("%q carries no law/law.py, so there is no chain to append to", t.Name)
	}
	if st, err := os.Stat(filepath.Join(t.Home, "law", lawLedger)); err != nil || st.IsDir() {
		return fmt.Sprintf("%q carries no law/%s", t.Name, lawLedger)
	}
	return ""
}

// lawState is the whole answer law_status gives, and what the writers judge by.
func lawState(t tenant.Tenant, opts Options) (lawAnswer, error) {
	ans := lawAnswer{World: t.Name, Files: []lawFile{}}
	ents, err := os.ReadDir(filepath.Join(t.Home, "law"))
	if err != nil {
		return ans, fmt.Errorf("this world carries no law/ folder: %w", err)
	}
	links, head, err := lawLinks(t)
	if err != nil {
		return ans, err
	}
	ans.Ok, ans.Verdict = lawVerify(t, opts)
	if m := lawProvenRe.FindStringSubmatch(ans.Verdict); m != nil {
		ans.Links, _ = strconv.Atoi(m[1])
		ans.Head = m[2]
	}
	if ans.Head == "" && len(head) >= 16 {
		ans.Head = head[:16]
	}
	if ans.Links == 0 {
		for _, refs := range links {
			ans.Links += len(refs)
		}
	}
	for _, e := range ents {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".md") {
			continue
		}
		p := filepath.Join(t.Home, "law", e.Name())
		st, err := os.Stat(p)
		if err != nil {
			continue
		}
		f := lawFile{Name: e.Name(), Bytes: st.Size()}
		if full, ok := lawSHA(p, 0); ok {
			f.Fingerprint = full[:16]
		}
		refs := links[e.Name()]
		switch {
		case len(refs) == 0:
			f.State = "not on the chain"
		default:
			intact := true
			for _, r := range refs {
				f.Links = append(f.Links, r.n)
				f.Kind, f.By = r.kind, r.by
				if r.extent > f.SealedBytes {
					f.SealedBytes = r.extent
				}
				if got, ok := lawSHA(p, r.extent); !ok || got != r.token {
					intact = false
				}
			}
			switch {
			case !intact:
				f.State = "changed since sealed"
			case f.SealedBytes > 0 && st.Size() > f.SealedBytes:
				f.State = "partly sealed"
				f.DraftBytes = st.Size() - f.SealedBytes
				ans.DraftBytes += f.DraftBytes
			default:
				f.State = "sealed"
			}
		}
		ans.Files = append(ans.Files, f)
	}
	sort.Slice(ans.Files, func(i, j int) bool { return ans.Files[i].Name < ans.Files[j].Name })
	_, _, ans.SittingOpen = engine.SittingOpen(t.Home)
	if _, open := engines.Get(t.Home); open {
		ans.SittingOpen = true
	}
	ans.Why = lawWriteRefusal(t)
	return ans, nil
}

func toolLawStatus(opts Options) Fn {
	return func(t tenant.Tenant, args map[string]any) (string, error) {
		ans, err := lawState(t, opts)
		if err != nil {
			return "", err
		}
		return marshal(ans)
	}
}

// ---- the operator's hand ---------------------------------------------------------------------------

func lawRecord(args map[string]any, t tenant.Tenant, c Caller, tool, what, note string) {
	reg := regOf(args)
	if reg == nil {
		return
	}
	reg.record(t.Home, what, Hold{
		ID: fmt.Sprintf("law_%d", time.Now().UnixMilli()), Tool: tool,
		Caller: callerLabel(c), Project: t.Name, RBAC: rbacState(t, c),
	}, note)
}

type lawDone struct {
	State      string `json:"state"` // added | sealed | refused
	Entry      int    `json:"entry,omitempty"`
	Title      string `json:"title,omitempty"`
	Why        string `json:"why,omitempty"`
	SealedTo   int64  `json:"sealed_to,omitempty"`
	DraftBytes int64  `json:"draft_bytes"`
	Verdict    string `json:"verdict,omitempty"`
}

func (d lawDone) String() string { s, _ := marshal(d); return s }

// lawText is what a person typed, made safe to append: LF only, trimmed, valid UTF-8, no control bytes but the tab.
func lawText(s string) (string, bool) {
	s = strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
	if !utf8.ValidString(s) {
		return "", false
	}
	for _, r := range s {
		if r < 0x20 && r != '\n' && r != '\t' {
			return "", false
		}
	}
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	return strings.Trim(strings.Join(lines, "\n"), "\n"), true
}

func toolLawAdd(opts Options) Fn {
	return func(t tenant.Tenant, args map[string]any) (string, error) {
		c := callerOf(args)
		// Registry.Call has refused every other caller already (ServiceOnly); a tool reached any other way refuses for itself.
		if !c.Service {
			return "", errors.New("refused: the laws are the operator's own hand and are offered to his glass alone")
		}
		refuse := func(why string) (string, error) {
			lawRecord(args, t, c, "law_add", "law_refused", why)
			return lawDone{State: "refused", Why: why}.String(), nil
		}
		title, okT := lawText(str(args, "title"))
		text, okX := lawText(str(args, "text"))
		from, okF := lawText(str(args, "from"))
		switch {
		case !okT || !okX || !okF:
			return refuse("the words carry a control character or are not text")
		case strings.Contains(title, "\n") || len([]rune(title)) < 3 || len([]rune(title)) > lawMaxTitle:
			return refuse(fmt.Sprintf("the title is one line of 3 to %d characters", lawMaxTitle))
		case len([]rune(text)) < lawMinText || len([]rune(text)) > lawMaxText:
			return refuse(fmt.Sprintf("the law is %d to %d characters; this one is %d", lawMinText, lawMaxText, len([]rune(text))))
		case len([]rune(from)) > lawMaxFrom:
			return refuse(fmt.Sprintf("where it came from is at most %d characters", lawMaxFrom))
		case lawHeadingRe.MatchString(text) || lawHeadingRe.MatchString(from):
			return refuse("a line that opens like a heading (## ...) or is a bare --- would split the entry in two; indent it or reword it")
		}
		lawMu.Lock()
		defer lawMu.Unlock()
		if why := lawWriteRefusal(t); why != "" {
			return refuse(why)
		}
		if ok, words := lawVerify(t, opts); !ok {
			return refuse("the chain does not prove whole right now, and nothing is laid on a chain that refuses: " + clip(words, 400))
		}
		path := filepath.Join(t.Home, "law", lawLedger)
		before, err := os.ReadFile(path)
		if err != nil {
			return refuse("the ledger could not be read: " + err.Error())
		}
		if len(before) == 0 || before[len(before)-1] != '\n' || strings.Contains(string(before), "\r") {
			return refuse("the ledger does not end on a line break, or carries a carriage return: it is a file to fix by hand, not by appending")
		}
		next := 1
		for _, m := range lawEntryRe.FindAllStringSubmatch(string(before), -1) {
			if n, _ := strconv.Atoi(m[1]); n >= next {
				next = n + 1
			}
		}
		if from == "" {
			from = "his word, through the glass's Laws page"
		}
		fromLines := strings.Split(from, "\n")
		var b strings.Builder
		b.WriteString("\n---\n\n")
		fmt.Fprintf(&b, "## %d. %s\n\n```\n", next, title)
		fmt.Fprintf(&b, "entered:  %s, on the operator's word\n", time.Now().Format("2006-01-02"))
		b.WriteString("from:     " + fromLines[0] + "\n")
		for _, l := range fromLines[1:] {
			b.WriteString("          " + l + "\n")
		}
		b.WriteString("```\n\n" + text + "\n")
		entry := b.String()

		// APPEND ONLY: the file is opened for append, so no byte already in it can be written over.
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			return refuse("the ledger could not be opened to append: " + err.Error())
		}
		_, werr := f.WriteString(entry)
		cerr := f.Close()
		if werr != nil || cerr != nil {
			_ = os.Truncate(path, int64(len(before)))
			return refuse("the entry could not be written, and the ledger is as it was")
		}
		// THE CHAIN IS WALKED AGAIN. An entry that left it refusing is taken back, to the byte.
		if ok, words := lawVerify(t, opts); !ok {
			_ = os.Truncate(path, int64(len(before)))
			return refuse("the chain refused after the entry was written, so it was taken back: " + clip(words, 400))
		}
		lawRecord(args, t, c, "law_add", "law_added", fmt.Sprintf("entry %d :: %s", next, title))
		ans, _ := lawState(t, opts)
		return lawDone{State: "added", Entry: next, Title: title, DraftBytes: ans.DraftBytes}.String(), nil
	}
}

func toolLawSeal(opts Options) Fn {
	return func(t tenant.Tenant, args map[string]any) (string, error) {
		c := callerOf(args)
		if !c.Service {
			return "", errors.New("refused: the laws are the operator's own hand and are offered to his glass alone")
		}
		refuse := func(why string) (string, error) {
			lawRecord(args, t, c, "law_seal", "law_refused", why)
			return lawDone{State: "refused", Why: why}.String(), nil
		}
		note, ok := lawText(str(args, "note"))
		if !ok || strings.Contains(note, "\n") || len([]rune(note)) > lawMaxTitle*2 {
			return refuse("the note is one line of plain text, at most 240 characters")
		}
		lawMu.Lock()
		defer lawMu.Unlock()
		if why := lawWriteRefusal(t); why != "" {
			return refuse(why)
		}
		st, err := lawState(t, opts)
		if err != nil {
			return refuse(err.Error())
		}
		if !st.Ok {
			return refuse("the chain does not prove whole right now, and nothing is sealed over a chain that refuses: " + clip(st.Verdict, 400))
		}
		var ledger *lawFile
		for i := range st.Files {
			if st.Files[i].Name == lawLedger {
				ledger = &st.Files[i]
			}
		}
		if ledger == nil || ledger.DraftBytes == 0 {
			return refuse("nothing is written below the ledger's seal, so there is nothing to seal")
		}
		if note == "" {
			note = "sealed from the glass's Laws page, " + time.Now().Format("2006-01-02")
		}
		script := filepath.Join(t.Home, "law", "law.py")
		res := lawRun(pythonWord(opts.CoreCmd), []string{script, "seal", lawLedger, "--note", note}, spawnOpts{Dir: t.Home, Timeout: lawWait})
		if res.Err != nil {
			return refuse("law.py would not seal: " + clip(strings.TrimSpace(res.Combined+" "+res.Err.Error()), 400))
		}
		after, _ := lawState(t, opts)
		var sealedTo int64
		for _, f := range after.Files {
			if f.Name == lawLedger {
				sealedTo = f.SealedBytes
			}
		}
		lawRecord(args, t, c, "law_seal", "law_sealed", fmt.Sprintf("ledger sealed to byte %d :: %s", sealedTo, note))
		return lawDone{State: "sealed", SealedTo: sealedTo, DraftBytes: after.DraftBytes, Verdict: after.Verdict}.String(), nil
	}
}
