// The glass's batch of 2026-09-29 (WHAT'S LEFT C17, C23, C25, C26, C27, C28,
// C29, D4, D5), held by reading the pages' own source the way the What's-left
// stroke does: the binary embeds these files and builds whether or not a page
// says what it must, so each promise below is pinned to the line that keeps it.
package server

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func page(t *testing.T, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "static", filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// funcOf is one method's source, from its name to the next method at the same
// indent, so a promise about ONE function is asked of that function alone.
func funcOf(src, sig string) string {
	i := strings.Index(src, sig)
	if i < 0 {
		return ""
	}
	body := src[i:]
	if j := strings.Index(body[1:], "\n  },\n"); j >= 0 {
		body = body[:j+1]
	}
	return body
}

// NO PAGE OPENS A BROWSER POP-UP (C26; §4.7 of the spec: "No prompt()
// anywhere"). The desktop app's browser pane dismisses a native prompt(),
// confirm() or alert() unseen, so a button behind one did nothing there --
// the Tools page's Call, and the eval scorer's three questions. The pages'
// own comments already forbade it; this is the comment made a wire. A call
// carries an argument, so the pattern is the paren and an opening quote,
// which a comment saying "never a prompt()" does not match.
func TestNoPageOpensABrowserPopUp(t *testing.T) {
	popup := regexp.MustCompile(`\b(prompt|confirm|alert)\(\s*['"` + "`" + `]`)
	dir := filepath.Join("..", "static", "js")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".js") {
			continue
		}
		src := page(t, "js/"+e.Name())
		for n, line := range strings.Split(src, "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue
			}
			if popup.MatchString(line) {
				t.Fatalf("%s:%d opens a browser pop-up, which the desktop pane dismisses unseen: %s",
					e.Name(), n+1, strings.TrimSpace(line))
			}
		}
	}
	// And the Tools page asks in the page: the modal it already owns.
	app := page(t, "js/app.js")
	call := funcOf(app, "  invokeTool(")
	for _, want := range []string{`tool-call-form`, `tool-call-args`, `modal.style.display = 'flex'`} {
		if !strings.Contains(call, want) {
			t.Fatalf("the Tools page's Call no longer asks in the modal (wanted %s)", want)
		}
	}
	if !strings.Contains(funcOf(app, "  addEval("), `eval-form`) {
		t.Fatal("the eval scorer no longer asks in the modal")
	}
}

// THE SIDEBAR READS THE COVENANT OFF THE RECORD (C29). It carried the house
// covenant as a literal, one of the 59 copies the door stopped minting from on
// 2026-09-25; the operator's own declaration says it, and us_to_vc mints from
// that declaration in exactly that namespace. The TUI's half is in
// cmd/atlas-tui's own stroke.
func TestTheSidebarReadsTheCovenantOffTheRecord(t *testing.T) {
	index := page(t, "index.html")
	app := page(t, "js/app.js")
	if strings.Contains(index, "1512741580b7239b") {
		t.Fatal("index.html carries the covenant as a literal again")
	}
	if !strings.Contains(index, `id="covenant"`) {
		t.Fatal("index.html has no element for the covenant to be read into")
	}
	paint := funcOf(app, "  async paintCovenant(")
	for _, want := range []string{`'us_to_vc'`, `'agents/operator.us'`, `project: 'atlas'`, `.covenant`} {
		if !strings.Contains(paint, want) {
			t.Fatalf("the sidebar no longer reads the covenant off the operator's declaration (wanted %s)", want)
		}
	}
	if !strings.Contains(funcOf(app, "  async loadHealth("), "this.paintCovenant()") {
		t.Fatal("nothing calls paintCovenant, so the sidebar would read 'reading the record...' for good")
	}
}

// A WATCHED TURN ENDS THE WAY EVERY TURN ENDS (C27). endWatched emitted 'done'
// while Home and Chat listen for 'end', so a turn started in another window
// kept its bubble marked live after the delivery had landed.
func TestAWatchedTurnEndsTheWayEveryTurnEnds(t *testing.T) {
	council := page(t, "js/council.js")
	end := funcOf(council, "  endWatched(")
	if !strings.Contains(end, "this.emit('end')") || strings.Contains(end, "this.emit('done')") {
		t.Fatalf("endWatched does not end the turn with 'end', the word the pages listen for:\n%s", end)
	}
	for _, f := range []string{"js/home.js", "js/chat.js"} {
		if !strings.Contains(page(t, f), "what === 'end'") {
			t.Fatalf("%s no longer finishes a bubble on 'end'", f)
		}
	}
}

// A DOOR THAT ANSWERED WITH A REFUSAL IS NOT A WORLD WITH NO ENGINE (C23). The
// glass proxies /council/state and answers 502 with a JSON body when the door
// is silent; that body parsed, `open` was absent, and the hero said "No
// engine" -- boot one -- over a door that was down.
func TestRunCheckTellsASilentDoorFromAWorldWithNoEngine(t *testing.T) {
	check := funcOf(page(t, "js/council.js"), "  async check(")
	if !strings.Contains(check, "if (!x.ok) throw") {
		t.Fatal("Run.check reads a 502's JSON as a world with no engine again")
	}
	if !strings.Contains(check, "this.unreachable = true") {
		t.Fatal("Run.check no longer says the door is unreachable when the read fails")
	}
}

// THE DASHBOARD OFFERS THE LIVE CHECK (D5) AND SAYS WHAT IDLE COSTS (C17).
func TestTheDashboardOffersTheLiveCheckAndSaysWhatIdleCosts(t *testing.T) {
	home := page(t, "js/home.js")
	if !strings.Contains(funcOf(home, "  async standup("), `App.tool('standup_run'`) {
		t.Fatal("the Dashboard has no way to run the live check")
	}
	if !strings.Contains(funcOf(home, "  paintEngine("), `id="eng-standup"`) {
		t.Fatal("the live check has no button on the hero")
	}
	age := funcOf(home, "  age(")
	if !strings.Contains(age, "IDLE_CLOSE_MIN") || !strings.Contains(age, "closes itself") {
		t.Fatal("the idle line no longer says that the engine closes itself, or at what")
	}
	if !strings.Contains(home, "IDLE_CLOSE_MIN: 30,") {
		t.Fatal("the idle close is not thirty minutes here, or is not named; serve.py's IDLE_CLOSE is thirty")
	}
}

// CLOSING SAYS WHEN NO TOLL IS OWED (C28). Three places said closing "pays its
// toll"; the core pays one only when a turn ran (serve.py: "the toll is paid
// unattended if runs happened").
func TestClosingSaysWhenNoTollIsOwed(t *testing.T) {
	for _, f := range []string{"js/app.js", "js/palette.js"} {
		if !strings.Contains(page(t, f), "no toll is owed") {
			t.Fatalf("%s says closing always pays its toll", f)
		}
	}
	home := page(t, "js/home.js")
	if strings.Contains(home, "closing pays its toll.") || !strings.Contains(home, "pays the toll if a turn ran") {
		t.Fatal("home.js says closing always pays its toll")
	}
}

// THE BUILDER OFFERS LOOPS (D4) AND FINDS A PAUSED RUN AFTER A RELOAD (C25).
// LANDING IS A BUTTON HE PRESSES HIMSELF (his ruling 2026-09-30, B16): the
// Lines-of-work box offers "Land onto main" for a line that is not the main
// line, only when the main line is the one you stand on, and it calls the
// door's `land` -- never a verb the table does not carry.
func TestVersionControlOffersLandingALineOntoMain(t *testing.T) {
	fl := page(t, "js/flows.js")
	if !strings.Contains(fl, `data-line="land"`) {
		t.Fatal("the Lines-of-work box offers no Land button")
	}
	if !strings.Contains(fl, "const onMain =") || !strings.Contains(fl, "(!b.main && onMain)") {
		t.Fatal("Land is not gated on standing on the main line and the line not being it")
	}
	if strings.Contains(fl, `action: 'merge'`) || strings.Contains(fl, `data-line="merge"`) {
		t.Fatal("the glass must not ask the door for a verb the table does not carry")
	}
}

func TestTheBuilderOffersLoopsAndFindsAPausedRunAfterAReload(t *testing.T) {
	wf := page(t, "js/workflows.js")
	for _, kind := range []string{"ask:", "run:", "seat:", "prompt:", "memory:"} {
		i := strings.Index(wf, "    "+kind)
		if i < 0 {
			t.Fatalf("the builder no longer offers the %s kind", kind)
		}
		line := wf[i:]
		line = line[:strings.Index(line, "\n")]
		if !strings.Contains(line, "'loops'") {
			t.Fatalf("the %s kind offers no loops box: %s", kind, strings.TrimSpace(line))
		}
	}
	for _, kind := range []string{"eval:", "gate:"} {
		i := strings.Index(wf, "    "+kind)
		line := wf[i:]
		line = line[:strings.Index(line, "\n")]
		if strings.Contains(line, "'loops'") {
			t.Fatalf("the %s kind offers a loops box the engine refuses by name", kind)
		}
	}
	if !strings.Contains(wf, "loops: true") {
		t.Fatal("loops is not typed as a number, so the door would refuse the whole save")
	}
	runs := funcOf(wf, "  async runs(")
	if !strings.Contains(runs, "PAUSED") || !strings.Contains(runs, "data-paused") {
		t.Fatal("the runs list no longer finds a paused run in the record")
	}
	if !strings.Contains(funcOf(wf, "  async reopen("), "API.flowStatus(") {
		t.Fatal("a paused run found after a reload cannot be opened")
	}
	if !strings.Contains(funcOf(wf, "  async render("), "await this.runs()") {
		t.Fatal("the page does not read the runs on arrival, so a paused run is found only once a flow is opened")
	}
}

// THE FRONT PAGE (2026-10-02). The operator, shown a page out of his own AI
// Studio project: "that's what I am looking for", then, asked where it should
// go, "replace the Dashboard now". The strokes below hold what makes it his
// page and not that project's: it is routed at / with the old launchpad kept
// whole at /dashboard; it asks the door only for tools the door carries, and
// asks quietly only where the glass will let it; it fakes nothing; and a tab
// that is not wired says so when it is pressed. The binary embeds these files
// and builds whether or not a page keeps any of this.
func TestTheFrontPageIsRoutedAndTheLaunchpadIsKept(t *testing.T) {
	index := page(t, "index.html")
	app := page(t, "js/app.js")
	pal := page(t, "js/palette.js")
	for _, want := range []string{
		`<link rel="stylesheet" href="/css/agent.css">`,
		`<script src="/js/agent.js"></script>`,
		`href="/" class="nav-link active" data-page="agent"`,
		`href="/dashboard" class="nav-link" data-page="dashboard"`,
	} {
		if !strings.Contains(index, want) {
			t.Fatalf("index.html no longer carries %s", want)
		}
	}
	router := funcOf(app, "  router(")
	for _, want := range []string{`|| 'agent'`, `classList.toggle('is-agent'`} {
		if !strings.Contains(router, want) {
			t.Fatalf("the router no longer sends / to the front page and takes the sidebar off it (wanted %s)", want)
		}
	}
	render := funcOf(app, "  async render(")
	for _, want := range []string{`case 'agent': await Agent.render(el)`, `case 'dashboard': await Home.render(el)`} {
		if !strings.Contains(render, want) {
			t.Fatalf("the page table lost %s", want)
		}
	}
	// BOOT AND CLOSE ARE THE FRONT PAGE'S. They navigate to / and then act; if
	// they went on calling Home, the engine would open under a page that has
	// no boot report to show, and say nothing.
	side := funcOf(app, "  bindSidebar(")
	for name, src := range map[string]string{"the sidebar button": side, "the palette": pal} {
		for _, want := range []string{"Agent.boot()", "Agent.closeSitting()"} {
			if !strings.Contains(src, want) {
				t.Fatalf("%s no longer calls %s", name, want)
			}
		}
		for _, bad := range []string{"Home.boot()", "Home.closeSitting()"} {
			if strings.Contains(src, bad) {
				t.Fatalf("%s still boots through %s, on a page that has no Dashboard under it", name, bad)
			}
		}
	}
}

func TestTheFrontPageFakesNothing(t *testing.T) {
	src := page(t, "js/agent.js")
	// A figure the record cannot prove is not shown (SPEC 3 invariant 10). The
	// page this one came from drew its tokens-per-second from a random number.
	if strings.Contains(src, "Math.random") {
		t.Fatal("the front page draws a figure from a random number")
	}
	// One lock: the PIN. No second sign-in, no token the page invents for itself.
	for _, bad := range []string{`type="password"`, `/api/login`, `agentos_token`} {
		if strings.Contains(src, bad) {
			t.Fatalf("the front page carries its own sign-in (%s); the lock is the PIN", bad)
		}
	}
	// A tab that is not wired is marked so and says so when pressed. None is left: Bash and Python were wired on
	// 2026-10-03 (WHAT'S LEFT H15) and Aider on 2026-10-04 (H12), held by the two strokes beside this one. The
	// machinery stays, and is held generically: a tab marked unwired must have its words, and words must belong to
	// an unwired tab.
	notWired := regexp.MustCompile(`(?s)NOT_WIRED: \{(.*?)\n  \},`).FindStringSubmatch(src)
	if notWired == nil {
		t.Fatal("the table of what an unwired tab says is gone")
	}
	modes := regexp.MustCompile(`id: '([a-z]+)',[^}]*wired: (true|false)`).FindAllStringSubmatch(src, -1)
	if len(modes) != 5 {
		t.Fatalf("read %d tabs; this stroke is reading the wrong thing", len(modes))
	}
	says := map[string]bool{}
	for _, m := range regexp.MustCompile(`(?m)^\s+([a-z]+): '[^']*not wired yet`).FindAllStringSubmatch(notWired[1], -1) {
		says[m[1]] = true
	}
	for _, m := range modes {
		if unwired := m[2] == "false"; unwired != says[m[1]] {
			t.Fatalf("the %s tab is marked wired: %s and its words say %v when pressed", m[1], m[2], says[m[1]])
		}
	}
	for _, id := range []string{"agent", "bash", "python", "pwsh", "aider"} {
		if !regexp.MustCompile(`id: '` + id + `',[^}]*wired: true`).MatchString(src) {
			t.Fatalf("the %s tab is not marked as wired", id)
		}
	}
	if strings.Contains(src, "Aider is not wired yet") {
		t.Fatal("the banner still says Aider is not wired")
	}
	// The failures go on the face of the answer, wherever the answer is shown
	// (ESTATE LAW 5).
	if face := funcOf(src, "  entryHtml("); !strings.Contains(face, "Run.failures(") || !strings.Contains(face, "NOT EVERYTHING RAN") {
		t.Fatal("the terminal can show an answer without showing what did not run")
	}
	// The gate is his. A question from the council is answered in a field he
	// types into; while it is open the prompt line takes nothing else; and the
	// hold queue is answered from a button he presses, nowhere else.
	if !strings.Contains(funcOf(src, "  answer("), "Run.start({ answer: text })") {
		t.Fatal("the gate's field does not answer the council")
	}
	if !strings.Contains(funcOf(src, "  async go("), "this.asking()") {
		t.Fatal("the prompt line takes an objective while the council is waiting on an answer")
	}
	// Two places answer a hold, and each is a button he presses: the Guardrails tab (paneClick) and a
	// shell card in the thread (decide, called from the terminal's click handler for the card's two
	// buttons and from nowhere else).
	if strings.Count(src, "App.tool('hold_answer'") != 2 ||
		!strings.Contains(funcOf(src, "  async paneClick("), "App.tool('hold_answer'") ||
		!strings.Contains(funcOf(src, "  async decide("), "App.tool('hold_answer'") {
		t.Fatal("a held call can be answered from somewhere other than a button he presses")
	}
	if strings.Count(src, "this.decide(") != 1 || !regexp.MustCompile(`(?s)act === 'shell-yes' \|\| act === 'shell-no'.{0,400}this\.decide\(`).MatchString(src) {
		t.Fatal("a shell card can be answered from somewhere other than the click handler of its two buttons")
	}
}

// THE SHELL TABS (WHAT'S LEFT H15, 2026-10-03). The page decides nothing about a command: the door's
// shell_run reads it, runs it, parks it or refuses it, and the page draws the one document the door
// answers in. These hold the page's half of that wire against the door's own text, so renaming a key
// at the door, or letting the page grow an opinion of its own, goes red here.
func TestTheShellTabsAreTheDoorsShellAndNothingElse(t *testing.T) {
	src := page(t, "js/agent.js")
	door, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "shell.go"))
	if err != nil {
		t.Fatalf("the door's shell is not beside the glass: %v", err)
	}
	holds, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "holds.go"))
	if err != nil {
		t.Fatal(err)
	}

	// The keys the page reads off the answer are keys the door's answer has.
	keys := map[string]bool{}
	for _, m := range regexp.MustCompile("`json:\"([a-z_]+)[,\"]").FindAllStringSubmatch(string(door), -1) {
		keys[m[1]] = true
	}
	if len(keys) < 10 {
		t.Fatalf("read %d keys off the door's answer; this stroke is reading the wrong thing", len(keys))
	}
	read := map[string]bool{}
	for _, fn := range []string{"  shellEntryHtml(", "  shellTake(", "  async decide("} {
		for _, m := range regexp.MustCompile(`\b(?:a|ran)\.([a-z_]+)\b`).FindAllStringSubmatch(funcOf(src, fn), -1) {
			read[m[1]] = true
			if !keys[m[1]] {
				t.Fatalf("the page reads .%s off the shell's answer and the door's answer has no such key", m[1])
			}
		}
	}
	for _, must := range []string{"state", "why", "hold", "exit", "class", "approved", "timed_out", "note", "output", "truncated"} {
		if !keys[must] {
			t.Fatalf("the door's answer lost the key %q the page is built on", must)
		}
	}
	for _, must := range []string{"state", "why", "hold", "exit", "approved", "timed_out", "note", "truncated"} {
		if !read[must] {
			t.Fatalf("the page no longer reads .%s, so what the door says in it never reaches him", must)
		}
	}

	// A card is drawn for a held answer and for nothing else; a refusal draws no button; a failure is
	// on the face (a non-zero exit is a bad chip and a bad box).
	face := funcOf(src, "  shellEntryHtml(")
	card := regexp.MustCompile(`(?s)if \(en\.state === 'held'\) \{(.*?)\n    \}`).FindStringSubmatch(face)
	if card == nil || !strings.Contains(card[1], "shell-yes") || !strings.Contains(card[1], "shell-no") {
		t.Fatal("the card is not drawn for a held answer, or does not carry its two buttons")
	}
	if strings.Count(face, "shell-yes") != 1 {
		t.Fatal("a shell button is drawn somewhere other than the card")
	}
	for _, want := range []string{"REFUSED BY NAME", "no click lifts a refusal", "a.exit !== 0", "ag-chip bad", "ag-pre${bad ? ' bad' : ''}"} {
		if !strings.Contains(face, want) {
			t.Fatalf("the shell's entry no longer carries %q", want)
		}
	}

	// The page asks the door's shell and nothing else for a command: one call, with the shell named,
	// and `reset` only from the session's own reset.
	if strings.Count(src, "App.tool('shell_run'") != 1 || !strings.Contains(funcOf(src, "  async shellAsk("), "App.tool('shell_run'") {
		t.Fatal("a command can reach the door by some way other than shellAsk")
	}
	if !strings.Contains(funcOf(src, "  async shellReset("), "reset: true") || strings.Count(src, "reset: true") != 1 {
		t.Fatal("the session can be ended from somewhere other than its own reset")
	}
	if !strings.Contains(funcOf(src, "  async shellRun("), "{ shell: kind, command: text, timeout: 300 }") {
		t.Fatal("an entry is not sent as the shell it was typed in, as typed")
	}
	// The page holds no list of commands or words of its own to judge by.
	for _, bad := range []string{"'rm'", "'curl'", "'git push'", "rm -rf", "writes|deletes", "isWrite", "looksLike"} {
		if strings.Contains(src, bad) {
			t.Fatalf("the page carries its own opinion of what a command does (%s); the door judges, the page draws", bad)
		}
	}
	// A card answered somewhere else (the Guardrails tab) stops offering buttons that could only be refused.
	if !strings.Contains(funcOf(src, "  async readHolds("), "this.settleCards(") || !strings.Contains(funcOf(src, "  settleCards("), "en.state = 'gone'") {
		t.Fatal("a shell card answered elsewhere keeps its buttons")
	}
	// A shell needs no engine: the council path is where the engine is asked for.
	if g := funcOf(src, "  async go("); !regexp.MustCompile(`(?s)this\.shellRun\(.*?Run\.engineOpen`).MatchString(g) {
		t.Fatal("a shell line is held up behind the engine check")
	}
	// The replay's words are the ones hold_answer says: `... ran:\n\n<document>`.
	if !strings.Contains(string(holds), `%q ran:\n\n%s`) || !strings.Contains(funcOf(src, "  async decide("), "ran:") {
		t.Fatal("the page and hold_answer no longer agree on where an approved run's document begins")
	}
	// Every wired tab offers something to try and says what it is for.
	for _, id := range []string{"agent", "bash", "python", "pwsh", "aider"} {
		if !regexp.MustCompile(`(?s)HINTS: \{.*?\b`+id+`: \[`).MatchString(src) || !regexp.MustCompile(`(?s)PLACEHOLDER: \{.*?\b`+id+`: '`).MatchString(src) {
			t.Fatalf("the %s tab has no hints or no placeholder", id)
		}
	}
}

// THE AIDER TAB (WHAT'S LEFT H12 and B20, 2026-10-04). The page decides nothing about an edit: the door's aider_run copies the
// files, runs Aider behind its wall, judges the edit, and writes it back or does not; the page draws the one document the
// door answers in. These hold the page's half of that wire against the door's own text, so renaming a key at the door,
// letting the page grow an opinion of its own, or letting a button reach something it must not goes red here.
func TestTheAiderTabIsTheDoorsAiderAndNothingElse(t *testing.T) {
	src := page(t, "js/agent.js")
	door, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "aider.go"))
	if err != nil {
		t.Fatalf("the door's Aider is not beside the glass: %v", err)
	}
	fn := func(sig string) string {
		body := funcOf(src, sig)
		if body == "" {
			t.Fatalf("the page has no %s", strings.TrimSpace(sig))
		}
		return body
	}
	asks := []string{"  async readAider(", "  aiderTitleState(", "  aiderBannerHtml(", "  aiderStatusText(", "  aiderEntryHtml(", "  aiderTake(", "  async aiderHtml("}
	var all strings.Builder
	for _, sig := range append([]string{"  async aiderGo(", "  aiderNames(", "  aiderDrop(", "  async aiderRun(", "  async aiderUndo(", "  aiderSay("}, asks...) {
		all.WriteString(fn(sig))
	}
	aider := all.String()

	// The keys the page reads off the door's documents are keys the door's documents have.
	keys := map[string]bool{}
	for _, m := range regexp.MustCompile("`json:\"([a-z_]+)[,\"]").FindAllStringSubmatch(string(door), -1) {
		keys[m[1]] = true
	}
	if len(keys) < 25 {
		t.Fatalf("read %d keys off the door's documents; this stroke is reading the wrong thing", len(keys))
	}
	read := map[string]bool{}
	for _, sig := range asks {
		for _, m := range regexp.MustCompile(`\b(?:a|f|r|s)\.([a-z_]+)\b`).FindAllStringSubmatch(fn(sig), -1) {
			read[m[1]] = true
			// `error` is the page's own word for a status it could not read (readAider), not a key of the door's.
			if !keys[m[1]] && m[1] != "error" {
				t.Fatalf("the page reads .%s off an Aider document and the door's documents have no such key", m[1])
			}
		}
	}
	for _, must := range []string{"state", "why", "changed", "diff", "withheld", "said", "ignored", "guarded", "exit", "timed_out", "note", "truncated", "run", "line", "model", "tokens", "ms",
		"installed", "on_line", "on_rack", "busy", "runs", "budget_kb", "undone", "new", "file", "added", "removed"} {
		if !keys[must] {
			t.Fatalf("the door's documents lost the key %q the page is built on", must)
		}
		if !read[must] {
			t.Fatalf("the page no longer reads .%s, so what the door says in it never reaches him", must)
		}
	}
	if !strings.Contains(string(door), `// ran | refused | undone`) {
		t.Fatal("the door's answer no longer says which states it has; this stroke is reading the wrong thing")
	}
	face := fn("  aiderEntryHtml(")
	for _, state := range []string{"'running'", "'note'", "'refused'", "'error'", "'undone'", "'ran'"} {
		if !strings.Contains(face, "en.state === "+state) {
			t.Fatalf("the Aider entry no longer draws the %s state", state)
		}
	}

	// A failure is on the face of the answer: a refusal, a wall that stopped something, an edit that was withheld,
	// a run that was ended or failed. A refusal draws no button.
	for _, want := range []string{"REFUSED BY NAME", "no click lifts a refusal", "THE WALL STOPPED", "NOT WRITTEN", "ENDED", "FAILED", "ag-chip bad", "COULD NOT RUN"} {
		if !strings.Contains(face, want) {
			t.Fatalf("the Aider entry no longer carries %q", want)
		}
	}
	// The Undo button is drawn for a run that wrote, and for nothing else: not a refusal, not a withheld edit, not a run
	// that changed nothing, not one already taken back.
	if strings.Count(face, "aider-undo") != 1 || !regexp.MustCompile(`if \(en\.state === 'ran' && a && a\.run && \(a\.changed \|\| \[\]\)\.length && !a\.withheld && !en\.undone\) \{\s*lines\.push\(`+"`"+`<div class="ag-btns"><button type="button" class="ag-btn no" data-act="aider-undo"`).MatchString(face) {
		t.Fatal("the Undo button is drawn for something other than a run that wrote")
	}

	// The page asks the door's Aider and nothing else: one call each, with what was typed as typed, and a status read
	// asked quietly. It reads two other things for the Inspector's world cards (the muster and each world's git), and
	// it never reaches a verb that saves, sends, lands, runs the suites or answers a hold.
	if strings.Count(src, "App.tool('aider_run'") != 1 || !strings.Contains(fn("  async aiderRun("), "App.tool('aider_run', { message: text, files: files.join('\\n'), context: read.join('\\n') })") {
		t.Fatal("an instruction can reach the door by some way other than aiderRun, or is not sent as typed with the names of the files in the chat")
	}
	if strings.Count(src, "App.tool('aider_undo'") != 1 || !strings.Contains(fn("  async aiderUndo("), "App.tool('aider_undo', { run: id })") {
		t.Fatal("a run can be taken back by some way other than aiderUndo")
	}
	if strings.Count(src, "App.tool('aider_status'") != 1 || !strings.Contains(fn("  async readAider("), "App.tool('aider_status', {}, true)") {
		t.Fatal("Aider's state is read from somewhere other than readAider, or not quietly")
	}
	for _, m := range regexp.MustCompile(`App\.tool\('([a-z_]+)'`).FindAllStringSubmatch(aider, -1) {
		switch m[1] {
		case "aider_run", "aider_undo", "aider_status", "muster", "git":
		default:
			t.Fatalf("the Aider tab asks the door for %s", m[1])
		}
	}
	// Taking a run back is a button he presses or a line he types: the terminal's click, the Inspector's click, and /undo.
	if strings.Count(src, "this.aiderUndo(") != 3 ||
		!regexp.MustCompile(`(?s)act === 'aider-undo'\) this\.aiderUndo\(`).MatchString(src) ||
		!regexp.MustCompile(`(?s)act === 'aider-undo'\) \{ this\.aiderUndo\(`).MatchString(funcOf(src, "  async paneClick(")) ||
		!strings.Contains(fn("  async aiderGo("), "this.aiderUndo(rest || this.aiderLast, t)") {
		t.Fatal("a run can be taken back from somewhere other than his button or his /undo")
	}
	// The page holds no opinion of what may be written: no list of names, no word of the door's rules.
	for _, bad := range []string{"'law/'", "worlds/", "CLAUDE.md", "'.env'", "isProtected", "neverWritten", "isMain", "'main'", "'master'"} {
		if strings.Contains(aider, bad) {
			t.Fatalf("the Aider tab carries its own opinion of what may be written (%s); the door judges, the page draws", bad)
		}
	}
	// The tab's own commands, and an instruction that is anything else.
	for _, cmd := range []string{"'add'", "'read'", "'drop'", "'files'", "'status'", "'undo'", "'help'"} {
		if !strings.Contains(fn("  async aiderGo("), "cmd === "+cmd) {
			t.Fatalf("the Aider tab no longer answers /%s", strings.Trim(cmd, "'"))
		}
	}
	// /drop says what it took out and what was never there: an answer that is only the chat's state afterwards reads as nothing done.
	if d := fn("  aiderDrop("); !strings.Contains(d, "'dropped '") || !strings.Contains(d, "' not in the chat'") || !strings.Contains(d, "this.aiderSay(typed, said.join(") {
		t.Fatal("/drop no longer says what it dropped, or what was never in the chat")
	}
	// Aider needs the door, not an engine: its line is sent, from the tab or by name, before the council's engine check is
	// made, and a line typed in the tab is routed there by the tab itself.
	g := funcOf(src, "  async go(")
	at := strings.Index(g, "this.aiderGo(")
	if at < 0 || strings.Contains(g[:at], "Run.engineOpen") || !strings.Contains(g, "via ? via.mode === 'aider' : this.isAider()") {
		t.Fatal("an Aider line is held up behind the engine check, or a line typed in the tab is not sent to Aider")
	}
	// A line may name the tab from any other, with the prefix stripped.
	if !strings.Contains(funcOf(src, "  prefixMode("), "mode: 'aider', text: m[1]") {
		t.Fatal("`aider ...` no longer names the tab, or sends the prefix to Aider as part of the instruction")
	}
	// The thread says a run is waiting, so a working Aider is not read as a hung one.
	if !strings.Contains(face, "A run takes minutes") || !strings.Contains(fn("  aiderTick("), "setInterval") {
		t.Fatal("a waiting run says nothing of how long it takes, or its clock does not move")
	}
}

// THE FRONT PAGE PLAYS WHAT THE MAKER MADE (WHAT'S LEFT E1, 2026-10-05). A person asked for a game and was handed a file path to open: the page never read the
// delivery's `project`, so the game she asked for was on another page. The delivery names the project in hand, read off the maker and never off a seat's
// words; the page asks the door's read-only `projects` tool which version that is and frames the page the glass already serves, sandboxed by the header it is
// served under. Each promise is pinned to the line that keeps it, because the binary embeds the page and builds whether or not it says what it must.
func TestTheFrontPagePlaysWhatTheMakerMade(t *testing.T) {
	src := page(t, "js/agent.js")
	css := page(t, "css/agent.css")
	fn := func(sig string) string {
		body := funcOf(src, sig)
		if body == "" {
			t.Fatalf("the page has no %s", strings.TrimSpace(sig))
		}
		return body
	}
	shell, play, paint, onRun := fn("  shell("), fn("  async playProject("), fn("  paintPlay("), fn("  onRun(")

	// Its place is a sibling of the output, between the output and the input line: #ag-out is rewritten whole on most events, and a frame inside it would
	// restart the game on every line she types.
	out, box, hints := strings.Index(shell, `id="ag-out"`), strings.Index(shell, `id="ag-play"`), strings.Index(shell, `id="ag-hints"`)
	if out < 0 || box < out || hints < box {
		t.Fatal("the play frame is not between the terminal's output and its input line, as a sibling of the output")
	}
	for _, sig := range []string{"  entryHtml(", "  shellEntryHtml(", "  aiderEntryHtml(", "  bannerHtml(", "  gateHtml(", "  liveHtml(", "  paintOut(", "  paintLive("} {
		if strings.Contains(fn(sig), "<iframe") {
			t.Fatalf("%s builds a frame inside what the terminal rewrites whole, so the game would restart on every repaint", strings.TrimSpace(sig))
		}
	}
	if strings.Count(src, "<iframe") != 1 || !strings.Contains(paint, "<iframe") {
		t.Fatal("the page builds a frame somewhere other than paintPlay, or more than one")
	}

	// The frame is the glass's own sandboxed page, framed as the Dashboard frames it: the url is Projects.pageUrl's, there is no `sandbox` attribute (the app's
	// own browser pane refuses any frame that carries one, measured 2026-09-21; the header the page is served under is the wall), and no referrer goes out.
	if !strings.Contains(paint, "Projects.pageUrl(s.name, s.version)") || strings.Contains(paint, "sandbox") || !strings.Contains(paint, `referrerpolicy="no-referrer"`) {
		t.Fatal("the frame is not the glass's own sandboxed project page, framed as the Dashboard frames it")
	}

	// The project named is the delivery's and nobody else's, as a string, only when it names one.
	if !regexp.MustCompile(`const made = Run\.turn && Run\.turn\.delivery;\s+if \(made && typeof made\.project === 'string' && made\.project\) this\.playProject\(made\.project\);`).MatchString(onRun) {
		t.Fatal("what is played is not read off the turn's delivery (`project`), or is played when the delivery names none")
	}
	if strings.Count(src, "this.playProject(") != 1 {
		t.Fatal("something other than a delivery puts a project on the play frame")
	}

	// Which version it is comes from the door: a quiet read of the read-only `projects` tool (the glass keeps no trace of a quiet read).
	if !strings.Contains(play, "App.tool('projects', { action: 'list' }, true)") {
		t.Fatal("the version of the project is not read off the door's `projects` tool, quietly")
	}
	bg, err := os.ReadFile(filepath.Join("..", "handlers", "handlers.go"))
	if err != nil || !strings.Contains(string(bg), `"projects": true`) {
		t.Fatalf("`projects` is not among the glass's quiet reads (%v)", err)
	}

	// The same page is never loaded twice, so a game she is playing keeps playing while she talks; a new version or another project is loaded, and a page she
	// put away comes back only then.
	if !regexp.MustCompile(`if \(s && s\.name === name && s\.version === version\) return;`).MatchString(play) || !strings.Contains(play, "closed: false") {
		t.Fatal("a page that is already shown would be loaded again, or a new version would not bring a closed page back")
	}
	if !strings.Contains(play, "this.stick();") {
		t.Fatal("showing the frame takes room from the terminal and the reply she asked for would scroll out of view")
	}
	if !strings.Contains(paint, "s.closed = true") || !strings.Contains(paint, "box.hidden = true") {
		t.Fatal("she cannot put the page away")
	}

	// The WHOLE game is on screen: a game is drawn for a window of its own size and the page it is served as cannot be scrolled to, so the frame is given a
	// window of PLAY_W by PLAY_H and shrunk (never enlarged) to the room the stage has, now and whenever that room changes.
	if !regexp.MustCompile(`PLAY_W: \d+,\s+PLAY_H: \d+,`).MatchString(src) || !strings.Contains(paint, `class="ag-play-body"`) || !strings.Contains(paint, "this.fitPlay();") ||
		!strings.Contains(paint, "new ResizeObserver(() => this.fitPlay())") || !strings.Contains(fn("  fitPlay("), "Math.min(w / this.PLAY_W, h / this.PLAY_H, 1)") {
		t.Fatal("the frame is not fitted to the room the stage has, so a game made for a bigger window would be cut off with no way to scroll to the rest")
	}

	// It has room and cannot crowd the terminal out: its own rule caps its height, and the window keeps the height it starts at.
	if !regexp.MustCompile(`\.ag-play \{[^}]*flex: none;[^}]*max-height: \d+px;`).MatchString(css) || !strings.Contains(css, ".ag-play[hidden] { display: none; }") {
		t.Fatal("the play frame has no height cap, or is not hidden when there is nothing to play")
	}
}

func TestEveryToolTheFrontPageAsksForIsCarriedByTheDoor(t *testing.T) {
	src := page(t, "js/agent.js")
	reg, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "tools.go"))
	if err != nil {
		t.Fatalf("the door's tool table is not beside the glass: %v", err)
	}
	carried := map[string]bool{}
	for _, m := range regexp.MustCompile(`Name:\s+"([a-z_]+)"`).FindAllStringSubmatch(string(reg), -1) {
		carried[m[1]] = true
	}
	asked := regexp.MustCompile(`App\.tool\('([a-z_]+)'`).FindAllStringSubmatch(src, -1)
	if len(asked) < 8 {
		t.Fatalf("the front page asks the door for %d things; this stroke is reading the wrong thing", len(asked))
	}
	for _, m := range asked {
		if !carried[m[1]] {
			t.Fatalf("the front page asks for %s and the door does not carry it", m[1])
		}
	}
	// A READ ASKED QUIETLY MUST BE ONE THE GLASS WILL ANSWER QUIETLY. Marked
	// `background` and not on the list, it is kept whole as a trace and
	// announced to every open tab -- the noise the 2026-09-16 ruling ended.
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
	n := 0
	for _, m := range regexp.MustCompile(`App\.tool\('([a-z_]+)',[^\n]*,\s*true\)`).FindAllStringSubmatch(src, -1) {
		n++
		if !quiet[m[1]] {
			t.Fatalf("the front page asks for %s quietly and the glass would keep it as a trace", m[1])
		}
	}
	if n == 0 {
		t.Fatal("the front page asks for nothing quietly, so it would toast a new trace at every refresh")
	}
	// AND THE ROUTES IT FETCHES ARE ROUTES THE GLASS SERVES.
	srv, err := os.ReadFile("server.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ js, route string }{
		{"'/api/me'", `"GET /api/me"`},
		{"'/council/stream?'", `"GET /api/council/stream"`},
	} {
		if !strings.Contains(src, c.js) {
			t.Fatalf("the front page no longer fetches %s", c.js)
		}
		if !strings.Contains(string(srv), c.route) {
			t.Fatalf("the front page fetches %s and the glass serves no %s", c.js, c.route)
		}
	}
}

// NO METHOD OF THE FRONT PAGE IS DEFINED TWICE. The page is one object literal, and in an object literal
// the later of two keys with the same name silently replaces the earlier: on 2026-10-03 a new method
// called `shell` replaced the page's own `shell()` (the template the page is drawn from), the page drew
// "[object Promise]", and every stroke that reads the source as text stayed green. Only a run of the
// page showed it. This is the check that would have said so.
func TestNoMethodOfTheFrontPageIsDefinedTwice(t *testing.T) {
	src := page(t, "js/agent.js")
	seen := map[string]int{}
	for i, line := range strings.Split(src, "\n") {
		m := regexp.MustCompile(`^  (?:async )?([A-Za-z_][A-Za-z0-9_]*)\(`).FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if first, dup := seen[m[1]]; dup {
			t.Errorf("agent.js defines %s() on line %d and again on line %d: the later one replaces the first", m[1], first, i+1)
		}
		seen[m[1]] = i + 1
	}
	if len(seen) < 60 {
		t.Fatalf("read %d methods; this stroke is reading the wrong thing", len(seen))
	}
}

// THE FRONT PAGE'S WINDOW STAYS THE SIZE IT STARTS AT, ON A NARROW SCREEN TOO. His word, 2026-10-03, watching the app's browser
// pane at 888 px: "the window should be locked to that scale it starts at and just scroll the terminal instead of the whole
// page." Below 1024 px the stylesheet makes the main area a COLUMN, and in a column the window's `flex: 1` sizes its HEIGHT by
// its content, overriding the `height: calc(100vh - 125px)` it is given: the page grew with every line (3696 px, then 4042) and
// the output box, which scrolls, was never short enough to. This holds the two halves of the cure -- the window keeps a given
// height, and the narrow block tells it not to flex -- and the page itself was measured in the browser at the same 888 px before
// and after (the CHANGELOG has the numbers), because a rule read as text can be in the file and still lose.
func TestTheWindowIsNotSizedByItsContentOnANarrowScreen(t *testing.T) {
	css := page(t, "css/agent.css")
	if !regexp.MustCompile(`(?s)\.ag-win, \.ag-panel \{[^}]*height: calc\(100vh - 125px\);`).MatchString(css) {
		t.Fatal("the window and the panel no longer carry their given height; this stroke is reading the wrong thing")
	}
	i := strings.Index(css, "@media (max-width: 1023px) {")
	if i < 0 {
		t.Fatal("the narrow-screen block (max-width: 1023px) is gone")
	}
	block := css[i:]
	if j := strings.Index(block, "\n}\n"); j >= 0 {
		block = block[:j]
	}
	if !strings.Contains(block, ".ag-main { flex-direction: column;") {
		t.Fatal("the narrow-screen block no longer makes the main area a column; this stroke is reading the wrong thing")
	}
	if !regexp.MustCompile(`\.ag-win \{ flex: none; \}`).MatchString(block) {
		t.Error("below 1024 px the window is a child of a column, and without `.ag-win { flex: none; }` its flex: 1 sizes its height by its content again")
	}
	if !regexp.MustCompile(`(?s)\.ag-out \{[^}]*overflow-y: auto;`).MatchString(css) {
		t.Error("the output box no longer scrolls (overflow-y: auto); the window would clip the terminal instead of scrolling it")
	}
}
