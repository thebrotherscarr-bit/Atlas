package server

// Strokes for the Laws page (2026-10-05, his word: "a LAWS page! great idea, Add it."), held the way the What's-left page's are: by reading the page's own
// source, because the binary embeds these files and builds whether or not a page says what it must. The door's side (what law_status reads, what law_add
// and law_seal refuse) is held in line/internal/tools/law_test.go; the real law.py is held against the door's reading of it by the core's own suite.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// THE LAWS PAGE IS FOUR THINGS THAT MUST AGREE, as What's left is: a line in the one page list, the tab its address opens (a case in
// the router until 2026-10-10, when the page retired into its tab; his cards: "Every page with a tab"), a script loaded BEFORE the
// router, and what the page reads and asks. (Until 2026-10-07 the line was the sidebar's, with a glyph of its own; the sidebar retired, and its
// links became App.PAGES, which the front page's Pages menu, the palette and each page's crumb read.)
func TestTheLawsPageIsWired(t *testing.T) {
	index := page(t, "index.html")
	app := page(t, "js/app.js")
	laws := page(t, "js/laws.js")
	for _, c := range []struct{ in, want, why string }{
		{app, `{ page: 'laws', href: '/laws', label: 'Laws' }`, "the page list has no line for Laws, so neither Pages nor the palette offers it"},
		{index, `<script src="/js/laws.js"></script>`, "index.html does not load the page's script"},
		{app, "laws: 'laws'", "the Laws page's address opens no tab, so /laws answers Page not found"},
		{laws, `App.tool('law_status', {})`, "the page no longer reads how far each law is sealed from the door's law_status"},
		{laws, `RULES: 'CLAUDE.md'`, "the page no longer names the standing rules it shows"},
		{laws, `LEDGER: 'law/LAW_LEDGER.md'`, "the page no longer names the ledger it shows"},
	} {
		if !strings.Contains(c.in, c.want) {
			t.Fatalf("%s -- wanted %s", c.why, c.want)
		}
	}
	if l, a := strings.Index(index, "/js/laws.js"), strings.Index(index, "/js/app.js"); l < 0 || l > a {
		t.Fatal("index.html loads js/laws.js after js/app.js -- the router would call Laws before it exists")
	}
}

// EVERY TOOL THE PAGE ASKS FOR IS CARRIED BY THE DOOR, and NOTHING IT ASKS IS QUIET. The glass answers a read off the record only for the front page's own
// polling reads (the 2026-09-16 ruling: the Dashboard and the front page ask every fifteen to thirty seconds); this page reads when it opens and after an
// action, as What's left does, so its reads are kept like any page's. The two that WRITE are never on the quiet list.
func TestEveryToolTheLawsPageAsksForIsCarriedByTheDoorAndNothingItAsksIsQuiet(t *testing.T) {
	src := page(t, "js/laws.js")
	carried := map[string]bool{}
	for _, f := range []string{"tools.go"} {
		reg, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", f))
		if err != nil {
			t.Fatalf("the door's tool table is not beside the glass: %v", err)
		}
		for _, m := range regexp.MustCompile(`Name:\s+"([a-z_]+)"`).FindAllStringSubmatch(string(reg), -1) {
			carried[m[1]] = true
		}
	}
	asked := regexp.MustCompile(`App\.tool\('([a-z_]+)'`).FindAllStringSubmatch(src, -1)
	want := map[string]bool{"law_status": false, "records": false, "law_add": false, "law_seal": false}
	for _, m := range asked {
		if !carried[m[1]] {
			t.Fatalf("the Laws page asks for %s and the door does not carry it", m[1])
		}
		if _, known := want[m[1]]; !known {
			t.Fatalf("the Laws page asks for %s, a tool it was never meant to reach", m[1])
		}
		want[m[1]] = true
	}
	for tool, seen := range want {
		if !seen {
			t.Fatalf("the Laws page no longer asks for %s", tool)
		}
	}
	hand, err := os.ReadFile(filepath.Join("..", "handlers", "handlers.go"))
	if err != nil {
		t.Fatal(err)
	}
	block := regexp.MustCompile(`(?s)var backgroundReads = map\[string\]bool\{(.*?)\}`).FindStringSubmatch(string(hand))
	if block == nil {
		t.Fatal("the glass's list of quiet reads is gone")
	}
	quiet := map[string]bool{}
	for _, m := range regexp.MustCompile(`"([a-z_]+)":\s*true`).FindAllStringSubmatch(block[1], -1) {
		quiet[m[1]] = true
	}
	for _, m := range regexp.MustCompile(`App\.tool\('([a-z_]+)',[^\n]*,\s*true\)`).FindAllStringSubmatch(src, -1) {
		t.Fatalf("the Laws page asks for %s quietly: it does not poll, so its reads are kept like What's left's", m[1])
	}
	for _, tool := range []string{"law_status", "records", "law_add", "law_seal"} {
		if quiet[tool] {
			t.Fatalf("%s is on the glass's quiet list: the Laws page's reads are kept, and a write must never be answerable off the record", tool)
		}
	}
}

// HOW FAR A LAW IS SEALED COMES FROM THE CHAIN AND NEVER FROM THE FOLDER. `records` marks the whole law/ folder sealed, which stopped being true the day
// the ledger grew a draft below its seal; the page that said "sealed" off that flag would show a draft as law.
func TestTheLawsPageSaysHowFarALawIsSealedFromTheChainAndNotFromTheFolder(t *testing.T) {
	src := page(t, "js/laws.js")
	if m := regexp.MustCompile(`\b(rules|ledger|doc|d)\.sealed\b`).FindString(src); m != "" {
		t.Fatalf("the Laws page reads a sealed flag off a document (%s): that is the folder's word, not the chain's", m)
	}
	for _, want := range []string{"led.sealed_bytes", "f.state", "st.draft_bytes", "st.verdict", "st.ok", "st.why"} {
		if !strings.Contains(src, want) {
			t.Fatalf("the Laws page no longer takes %s from law_status", want)
		}
	}
	// An entry is a draft or sealed by where its heading stands against the seal's byte, counted in bytes (not characters).
	parse := funcOf(src, "  parseLedger(")
	if !strings.Contains(parse, "TextEncoder") || !strings.Contains(parse, "< sealedBytes") {
		t.Fatal("the ledger's entries are no longer judged against the seal's byte, in bytes")
	}
}

// THE TWO ACTIONS ARE HIS, AND ASKED FOR ONCE: law_add only from the form's own button, law_seal only after the inline question, and the page never binds a
// law whole (`direct`) or writes anywhere itself.
func TestTheLawsPageSetsOnlyWhatWasTypedAndSealsOnlyAfterAsking(t *testing.T) {
	src := page(t, "js/laws.js")
	if n := strings.Count(src, "App.tool('law_add'"); n != 1 {
		t.Fatalf("law_add is asked for %d times on the page, want once (the form's own button)", n)
	}
	if n := strings.Count(src, "App.tool('law_seal'"); n != 1 {
		t.Fatalf("law_seal is asked for %d times on the page, want once (after the question)", n)
	}
	bind := funcOf(src, "  bind(el, st) {")
	if bind == "" {
		t.Fatal("the page's bind is gone")
	}
	add := bind[strings.Index(bind, "$('laws-add').onclick"):]
	if i := strings.Index(add, "App.tool('law_add'"); i < 0 || !strings.Contains(add[:i], "title") {
		t.Fatal("law_add is no longer asked for from the form's own button, with what was typed")
	}
	// The seal button only REVEALS the question; the call lives behind the question's own confirm.
	reveal := bind[strings.Index(bind, "$('laws-seal').onclick"):]
	reveal = reveal[:strings.Index(reveal, "\n")]
	if strings.Contains(reveal, "law_seal") || !strings.Contains(reveal, "laws-seal-ask") {
		t.Fatalf("the seal button calls law_seal itself instead of asking first: %s", reveal)
	}
	goes := bind[strings.Index(bind, "$('laws-seal-go').onclick"):]
	if !strings.Contains(goes, "App.tool('law_seal'") {
		t.Fatal("law_seal is not asked for from the question's own confirm")
	}
	for _, want := range []string{`id="laws-seal-ask" hidden`, "After it a byte of it cannot", `id="laws-seal-no"`, "Not now"} {
		if !strings.Contains(src, want) {
			t.Fatalf("the question before sealing no longer says what it does and lets him say no (wanted %s)", want)
		}
	}
	for _, bad := range []string{"'direct'", `"direct"`, "fetch(", "'/api/", "localStorage"} {
		if strings.Contains(src, bad) {
			t.Fatalf("the Laws page reaches around the door's two tools (%s)", bad)
		}
	}
	// A page that asks nothing it knows would be refused: the buttons are disabled with the reason, read off law_status.
	if strings.Count(src, "${blocked ? 'disabled title=") != 2 {
		t.Fatal("the buttons that write are no longer disabled, with the reason, while the door would refuse them")
	}
}
