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
	// A tab that is not wired is marked so and says so when pressed. Aider is the one left: Bash and
	// Python were wired on 2026-10-03 (WHAT'S LEFT H15) and are held by
	// TestTheShellTabsAreTheDoorsShellAndNothingElse.
	notWired := regexp.MustCompile(`(?s)NOT_WIRED: \{(.*?)\n  \},`).FindStringSubmatch(src)
	if notWired == nil {
		t.Fatal("the table of what an unwired tab says is gone")
	}
	for _, id := range []string{"aider"} {
		if !regexp.MustCompile(`id: '` + id + `',[^}]*wired: false`).MatchString(src) {
			t.Fatalf("the %s tab is not marked as unwired", id)
		}
		if !regexp.MustCompile(id + `: '[^']*not wired yet`).MatchString(notWired[1]) {
			t.Fatalf("pressing the %s tab does not say it is not wired yet", id)
		}
	}
	for _, id := range []string{"bash", "python"} {
		if !regexp.MustCompile(`id: '` + id + `',[^}]*wired: true`).MatchString(src) {
			t.Fatalf("the %s tab is not marked as wired", id)
		}
		if strings.Contains(notWired[1], id+": '") {
			t.Fatalf("pressing the %s tab still says it is not wired", id)
		}
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
	if !strings.Contains(funcOf(src, "  async shellRun("), "{ shell: kind, command: text }") {
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
	for _, id := range []string{"agent", "bash", "python"} {
		if !regexp.MustCompile(`(?s)HINTS: \{.*?\b`+id+`: \[`).MatchString(src) || !regexp.MustCompile(`(?s)PLACEHOLDER: \{.*?\b`+id+`: '`).MatchString(src) {
			t.Fatalf("the %s tab has no hints or no placeholder", id)
		}
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
