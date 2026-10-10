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

// THE FRONT PAGE'S FOOTER READS THE COVENANT OFF THE RECORD (C29). The sidebar
// carried the house covenant as a literal, one of the 59 copies the door
// stopped minting from on 2026-09-25; the operator's own declaration says it,
// and us_to_vc mints from that declaration in exactly that namespace. Since the
// sidebar retired (2026-10-07, on his card: "Front page footer") the line is
// the front page's footer, drawn every time the page is. The TUI's half is in
// cmd/atlas-tui's own stroke.
func TestTheFooterReadsTheCovenantOffTheRecord(t *testing.T) {
	index := page(t, "index.html")
	app := page(t, "js/app.js")
	src := page(t, "js/agent.js")
	if strings.Contains(index, "1512741580b7239b") || strings.Contains(src, "1512741580b7239b") {
		t.Fatal("the glass carries the covenant as a literal again")
	}
	shell := funcOf(src, "  shell(")
	if foot, el := strings.Index(shell, `<footer class="ag-foot">`), strings.Index(shell, `id="covenant"`); foot < 0 || el < foot {
		t.Fatal("the front page's footer has no element for the covenant to be read into")
	}
	paint := funcOf(app, "  async paintCovenant(")
	for _, want := range []string{`'us_to_vc'`, `'agents/operator.us'`, `project: 'atlas'`, `.covenant`} {
		if !strings.Contains(paint, want) {
			t.Fatalf("the footer no longer reads the covenant off the operator's declaration (wanted %s)", want)
		}
	}
	if !strings.Contains(funcOf(src, "  async render("), "App.paintCovenant()") {
		t.Fatal("nothing calls paintCovenant, so the footer would read 'reading the record...' for good")
	}
}

// UNSENT WORK IS A NUMBER ON THE GITHUB DOOR (2026-10-08, on his card: "A GitHub door in the top bar"; on the Aider
// Pair door from 2026-10-07, on his card: "On the Aider Pair door"). The sidebar's warning badge on Version control
// moved with the page when the sidebar retired, and again when the page became GitHub with a door of its own: the
// glass reads the number as it loads and after every act on the GitHub page, and the front page draws it on the
// GitHub door, beside Aider Pair, hidden while nothing is owed, and on no other door.
func TestUnsentWorkIsANumberOnTheGitHubDoor(t *testing.T) {
	src := page(t, "js/agent.js")
	app := page(t, "js/app.js")
	if !strings.Contains(funcOf(src, "  shell("), `d.id === 'github' ? '<span class="ag-owed" id="ag-owed" hidden></span>'`) || strings.Count(src, `id="ag-owed"`) != 1 {
		t.Fatal("the GitHub door has no place for the count of unsent work, or another door has one too")
	}
	i := strings.Index(src, "  DOORS: [")
	if i < 0 {
		t.Fatal("the bar's doors are gone; this stroke is reading the wrong thing")
	}
	doors := src[i:]
	if j := strings.Index(doors, "  ],"); j >= 0 {
		doors = doors[:j]
	}
	if a, g := strings.Index(doors, "{ id: 'aider',  label: 'Aider Pair' },"), strings.Index(doors, "{ id: 'github', label: 'GitHub' },"); a < 0 || g < a {
		t.Fatal("the GitHub door is not in the bar, beside Aider Pair")
	}
	show := funcOf(app, "  showOwed(")
	if !strings.Contains(show, "getElementById('ag-owed')") || !strings.Contains(show, "el.hidden = true") {
		t.Fatal("the count is not drawn on the door, or is not hidden while nothing is owed")
	}
	if !strings.Contains(funcOf(app, "  async paintOwed("), "this.showOwed()") ||
		!strings.Contains(funcOf(app, "  init("), "this.paintOwed()") ||
		!strings.Contains(page(t, "js/flows.js"), "App.paintOwed()") ||
		!strings.Contains(funcOf(src, "  async render("), "App.showOwed()") {
		t.Fatal("the count is not read as the glass loads and after an act on the GitHub page, or not drawn when the front page is")
	}
}

// A WATCHED TURN ENDS THE WAY EVERY TURN ENDS (C27). endWatched emitted 'done'
// while Home and Chat listen for 'end', so a turn started in another window
// kept its bubble marked live after the delivery had landed. Home retired on
// 2026-10-07; the pages that listen now are the Watchboard and the front page.
func TestAWatchedTurnEndsTheWayEveryTurnEnds(t *testing.T) {
	council := page(t, "js/council.js")
	end := funcOf(council, "  endWatched(")
	if !strings.Contains(end, "this.emit('end')") || strings.Contains(end, "this.emit('done')") {
		t.Fatalf("endWatched does not end the turn with 'end', the word the pages listen for:\n%s", end)
	}
	for _, f := range []string{"js/chat.js", "js/agent.js"} {
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

// THE FRONT PAGE SAYS WHAT IDLE COSTS (C17). The Dashboard said it on its hero
// until it retired (2026-10-07, on his card: "Carry it over"); the terminal's
// title on the front page says it now. The Dashboard's live check (D5) is the
// Audit Ledger tab's, held by TestTheAuditLedgerCarriesTheProofAndTheLiveCheck.
func TestTheFrontPageSaysWhatIdleCosts(t *testing.T) {
	src := page(t, "js/agent.js")
	idle := funcOf(src, "  idle(")
	if !strings.Contains(idle, "IDLE_CLOSE_MIN") || !strings.Contains(idle, "closes itself") {
		t.Fatal("the idle line no longer says that the engine closes itself, or at what")
	}
	if !strings.Contains(src, "IDLE_CLOSE_MIN: 30,") {
		t.Fatal("the idle close is not thirty minutes here, or is not named; serve.py's IDLE_CLOSE is thirty")
	}
	if !strings.Contains(funcOf(src, "  paintTitle("), "this.idle()") {
		t.Fatal("the front page's title does not say the idle line")
	}
}

// CLOSING SAYS WHEN NO TOLL IS OWED (C28). Three places said closing "pays its
// toll"; the core pays one only when a turn ran (serve.py: "the toll is paid
// unattended if runs happened"). Since the sidebar and the Dashboard retired
// (2026-10-07) the places that say it are the palette and the front page.
func TestClosingSaysWhenNoTollIsOwed(t *testing.T) {
	if !strings.Contains(page(t, "js/palette.js"), "no toll is owed") {
		t.Fatal("palette.js says closing always pays its toll")
	}
	src := page(t, "js/agent.js")
	if strings.Contains(src, "closing pays its toll.") || !strings.Contains(src, "pays the toll if a turn ran") {
		t.Fatal("agent.js says closing always pays its toll")
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
	// SUPERSEDED IN PART 2026-10-09, on his word: "we have to be able to merge ... supersede that one". The glass asks
	// git_pr to merge a pull request on GitHub (the stroke below holds it); the Lines-of-work box still lands with
	// land, fast-forward, and never asks the door to merge.
	if lines := funcOf(fl, "  async lines(") + funcOf(fl, "  async line("); strings.Contains(lines, "merge") {
		t.Fatal("the Lines-of-work box asks the door to merge; a line lands onto main with land, fast-forward")
	}
}

// PULL REQUESTS ON THE GITHUB TAB (2026-10-07, his ruling on a card: "Open the wall for PRs"; merged there since
// 2026-10-09, his card: "Merge on GitHub when green"). Each world's buttons carry Pull requests beside the marks; the
// panel lists through the door's git_pr and opens one from its own button and nowhere else; each row's Merge on GitHub
// is live only while that pull request's checks have all passed and is greyed with the reason otherwise, and it merges
// from the panel's own merge, on a second press, and nowhere else; and the door carries the tool the glass asks for,
// with list as its reading action, open and merge as its writing ones, and merge judged by the pull request's checks.
func TestVersionControlListsAndOpensPullRequests(t *testing.T) {
	fl := page(t, "js/flows.js")
	if !strings.Contains(funcOf(fl, "  controls("), `data-act="pulls"`) || !strings.Contains(funcOf(fl, "  async act("), "await this.pulls(w)") {
		t.Fatal("the world's buttons offer no Pull requests, or the button does not open its panel")
	}
	pulls := funcOf(fl, "  async pulls(")
	if !strings.Contains(pulls, "App.tool('git_pr', { project: w, action: 'list' })") || !strings.Contains(pulls, `data-pull="open"`) {
		t.Fatal("the Pull requests panel does not read the door's list, or offers no way to open one")
	}
	if strings.Count(fl, "action: 'open' })") != 1 || !strings.Contains(funcOf(fl, "  async pull("), "App.tool('git_pr', { project: w, action: 'open' })") {
		t.Fatal("a pull request can be opened from somewhere other than the panel's own button")
	}
	if !regexp.MustCompile(`(?s)const merge = p\.checks === 'passed'\s+\? `+"`"+`<button class="btn btn-sm" data-pull="merge"`).MatchString(pulls) ||
		strings.Count(pulls, `data-pull="merge"`) != 1 || !strings.Contains(pulls, `disabled title="${escHtml(notYet[p.checks]`) {
		t.Fatal("Merge on GitHub is live on a pull request whose checks have not all passed, or is greyed without its reason")
	}
	merge := funcOf(fl, "  async merge(")
	if strings.Count(fl, "action: 'merge'") != 1 || !strings.Contains(merge, "App.tool('git_pr', { project: w, action: 'merge', number: Number(n) })") ||
		!strings.Contains(merge, "if (this.arming !== key)") || !strings.Contains(merge, "Click again to merge") {
		t.Fatal("a pull request can be merged from somewhere other than the panel's own button, or on one press")
	}
	door, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "gitpr.go"))
	if err != nil {
		t.Fatalf("the door's pull requests are not beside the glass: %v", err)
	}
	for _, want := range []*regexp.Regexp{
		regexp.MustCompile(`Name:\s+"git_pr"`),
		regexp.MustCompile(`Reads:\s+\[\]string\{"list"\}`),
		regexp.MustCompile(`case "open":`),
		regexp.MustCompile(`case "merge":`),
		regexp.MustCompile(`Args:\s+\[\]string\{"action\?", "project\?", "number\?", "line\?", "wait\?"\}`),
		regexp.MustCompile(`(?s)func prMerge\(.*?switch prChecks\(pr\.StatusCheckRollup\)`),
	} {
		if !want.Match(door) {
			t.Fatalf("the door's git_pr no longer carries %s, which the glass is built on", want)
		}
	}
}

// A LINE OF WORK THAT HAS NEVER BEEN SENT CAN BE SENT, AND IS SAID TO BE NEVER SENT (2026-10-09, the core's WHAT'S
// LEFT I1, third step: "Send to GitHub on a line of work that has never been sent, and the GitHub tab's words for it").
// Ahead and behind are counted against an upstream, and a line that has never been sent has none: the page greyed its
// Send as "Nothing to send" over saves GitHub had never seen, and called the line "not linked to GitHub at all" in a
// world that was linked. The `git` tool lists the world's remotes either way, so the page asks whether there is an
// origin before it judges a line with no upstream: with one, Send is offered "for the first time", because the door's
// git_push sends such a line to origin and links it there, and says so; with none, the line is not linked at all.
func TestALineThatHasNeverBeenSentCanBeSent(t *testing.T) {
	fl := page(t, "js/flows.js")
	controls, plain := funcOf(fl, "  controls("), funcOf(fl, "  plain(")
	for _, f := range []string{controls, plain} {
		if !strings.Contains(f, "const origin = (g.remotes || []).includes('origin');") {
			t.Fatal("the GitHub page judges a line with no upstream without asking whether the world has an origin to send to")
		}
	}
	if !strings.Contains(controls, "const nothingToSend = !g.upstream ? !origin : !g.ahead;") || !strings.Contains(controls, "for the first time") {
		t.Fatal("Send is greyed on a line that has never been sent, as if GitHub already had its saves, or does not say it sends the line for the first time")
	}
	if !strings.Contains(controls, "this line has never been sent, so GitHub holds nothing of it") {
		t.Fatal("Take from GitHub tells a never-sent line that it already has every save on GitHub")
	}
	if strings.Count(fl, "not linked to GitHub at all") != 1 ||
		!regexp.MustCompile(`(?s)gh = origin\s+\? 'this line of work has never been sent to GitHub.*?: 'not linked to GitHub at all`).MatchString(plain) {
		t.Fatal("a never-sent line in a world with an origin is called not linked to GitHub at all, or a world with no origin is not")
	}
	// The page judges by two keys the door's `git` writes, and offers a send the door's git_push makes.
	state, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "gitstate.go"))
	if err != nil {
		t.Fatalf("the door's git state is not beside the glass: %v", err)
	}
	for _, key := range []string{`out["upstream"]`, `out["remotes"]`} {
		if !strings.Contains(string(state), key) {
			t.Fatalf("the door's git tool no longer writes %s, which the page judges a never-sent line by", key)
		}
	}
	ctl, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "gitctl.go"))
	if err != nil {
		t.Fatalf("the door's git verbs are not beside the glass: %v", err)
	}
	push := string(ctl)
	i := strings.Index(push, "func toolGitPush(")
	if i < 0 {
		t.Fatal("the door's git_push is not where this stroke reads it")
	}
	push = push[i:]
	if j := strings.Index(push, "\nfunc "); j >= 0 {
		push = push[:j]
	}
	if !strings.Contains(push, `"--set-upstream", "origin", branch`) || !strings.Contains(push, "had never been sent before") {
		t.Fatal("the door's git_push no longer sends a never-sent line to origin and links it there, which the page offers")
	}
}

// A MAIN ON GITHUB IS OFFERED THE PULL REQUEST AND NOTHING ELSE (2026-10-09, the core's WHAT'S LEFT I1, seventh step;
// his card: "Main is on GitHub"). Since the fifth step GitHub takes a change to main only through a pull request, in
// both repositories and from him too, and the page still offered Land onto main and a Send of main: a road that ended
// in GitHub's refusal. The door judges it once (mainByPullRequest: origin has the main line) and carries its sentence on
// the lines list (`why_not_land`) and on the `git` state while he stands on main (`why_not_send`), and land and git_push
// refuse in those words; so the page greys Land onto main and Send to GitHub with the door's sentence, says it once
// under the lines, and judges nothing itself. (The council's Push went with the page's council path, 2026-10-10.)
func TestAMainOnGitHubIsOfferedThePullRequestAndNothingElse(t *testing.T) {
	fl := page(t, "js/flows.js")
	lines, controls := funcOf(fl, "  async lines("), funcOf(fl, "  controls(")
	if !strings.Contains(lines, "d.why_not_land ? `<button class=\"btn btn-sm\" disabled title=\"${escHtml(d.why_not_land)}\">Land onto main</button>`") ||
		!strings.Contains(lines, "<div class=\"muted mt-16\">${escHtml(d.why_not_land)}</div>") {
		t.Fatal("Land onto main is live on a main GitHub has, or is greyed without the door's sentence, or the sentence is not said under the lines")
	}
	if !strings.Contains(controls, "const mainByPr = g.why_not_send || '';") || !strings.Contains(controls, ": mainByPr ? mainByPr") ||
		!strings.Contains(controls, "${walled || nothingToSend || mainByPr ? 'disabled' : ''}") {
		t.Fatal("Send to GitHub is live on a main GitHub has, or is greyed without the door's sentence")
	}
	if strings.Contains(fl, "Land onto main, then Send, still lands a line") {
		t.Fatal("the page still says that Land onto main, then Send, lands a line without a pull request")
	}
	// The keys the page reads are the door's, and the door's land, git_push, lines list and git state ask one judgement.
	read := func(name string) string {
		b, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", name))
		if err != nil {
			t.Fatalf("the door's %s is not beside the glass: %v", name, err)
		}
		return string(b)
	}
	body := func(src, sig string) string {
		i := strings.Index(src, sig)
		if i < 0 {
			t.Fatalf("the door has no %s; this stroke is reading the wrong thing", sig)
		}
		src = src[i:]
		if j := strings.Index(src, "\nfunc "); j >= 0 {
			src = src[:j]
		}
		return src
	}
	ctl, state := read("gitctl.go"), read("gitstate.go")
	if !strings.Contains(body(ctl, "func mainByPullRequest("), `"rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+mainName`) {
		t.Fatal("the door no longer judges a main on GitHub by origin having it")
	}
	if !strings.Contains(body(ctl, "func toolGitBranch("), `if refusal := mainByPullRequest(t); refusal != ""`) ||
		!strings.Contains(body(ctl, "func toolGitPush("), `if refusal := mainByPullRequest(t); refusal != ""`) ||
		!strings.Contains(body(ctl, "func branchList("), `out["why_not_land"] = why`) ||
		!strings.Contains(body(state, "func toolGit("), `out["why_not_send"] = why`) {
		t.Fatal("land, git_push, the lines list and the git state no longer ask one judgement, so a button and its act could disagree")
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

// ONE DECISION PER GATE (2026-10-05): a gate's buttons go grey at the first
// click and come back only if the door refuses the move. Live, a click made
// while the run moved was held by the door and answered the next gate, unseen.
func TestAGatesButtonsGoGreyAtTheFirstClick(t *testing.T) {
	wf := page(t, "js/workflows.js")
	paint := funcOf(wf, "  paint(")
	if !strings.Contains(paint, "b.disabled = true") || strings.Contains(paint, "onclick = () => this.resume(") {
		t.Fatal("a gate's buttons stay live while the run moves, so a second click answers the next gate unseen")
	}
	if !strings.Contains(funcOf(wf, "  async resume("), "b.disabled = false") {
		t.Fatal("a refused move leaves the gate's buttons grey for good")
	}
}

// THE PAGE ASKS FOR EVERY VAR THE ENGINE RENDERS (2026-10-05). The release
// flow's `mark` rides only in its tool nodes' args and its gates' titles, and
// openVars read neither: fired from this page, the flow had no mark and its
// first step refused. Every field the engine renders is scanned.
func TestTheBuilderAsksForEveryVarTheEngineRenders(t *testing.T) {
	scan := funcOf(page(t, "js/workflows.js"), "  openVars(")
	for _, field := range []string{"n.question", "n.expected", "n.vars", "n.files", "n.args", "n.title"} {
		if !strings.Contains(scan, field) {
			t.Fatalf("openVars no longer reads %s, so a var only it names is never asked for", field)
		}
	}
}

// THE BUILDER IS IN THE INSPECTOR (his word, 2026-10-06: the builder goes
// "Inside the Inspector"). The front page's Workflows tab draws the Workflows
// page's own builder into its panel -- one object, so the two cannot drift --
// and no longer sends him to that page for it; the flow in hand outlives the
// redraw the tab gets every time it opens; a run that comes back while the tab
// is shut waits on the record; the builder takes the page's look; and the tab
// is named as the bar's door is. Since 2026-10-10 (his cards: "Every page with a
// tab") the page's address opens the tab, so the builder is drawn there alone,
// and the tab's way to GitHub opens the GitHub tab rather than a page.
func TestTheInspectorDrawsTheBuilderItself(t *testing.T) {
	src := page(t, "js/agent.js")
	wf := page(t, "js/workflows.js")
	css := page(t, "css/agent.css")
	pane := funcOf(src, "  async paintPane(")
	if !strings.Contains(pane, "id === 'flows'") || !strings.Contains(pane, "this.flowsPane(stamp)") {
		t.Fatal("the Workflows tab does not draw the builder")
	}
	if !strings.Contains(funcOf(src, "  async flowsPane("), "Workflows.render(box)") {
		t.Fatal("the Workflows tab draws something other than the Workflows page's own builder")
	}
	if strings.Contains(src, `data-act="goto"`) || strings.Contains(src, "act === 'goto'") ||
		!strings.Contains(funcOf(src, "  async flowsPane("), `data-act="door" data-door="github"`) {
		t.Fatal("a tab still sends him to a page, or the builder's way to GitHub does not open the GitHub tab")
	}
	render := funcOf(wf, "  async render(")
	if strings.Contains(wf, "inPanel") || strings.Contains(render, "page-title") {
		t.Fatal("the builder still asks whether it stands in the panel, or draws a page's header")
	}
	if !strings.Contains(render, "if (this.spec) this.build();") {
		t.Fatal("the flow in hand is lost when the tab is drawn again, and an unsaved edit with it")
	}
	if !strings.Contains(funcOf(wf, "  paint("), "if (!box) return;") {
		t.Fatal("a run that comes back while the tab is shut throws instead of waiting on the record")
	}
	if !regexp.MustCompile(`\.ag-wf \{[^}]*--bg:`).MatchString(css) || !strings.Contains(css, ".ag-wf .btn {") {
		t.Fatal("the builder in the panel does not take the page's look, and .ag button strips its buttons bare")
	}
	if strings.Contains(src, "Docs & Registry") || strings.Count(src, "label: 'Registry & Docs'") != 2 {
		t.Fatal("the Inspector's tab and the bar's door do not both say Registry & Docs")
	}
}

// NOTHING IS FIRED ON AN EMPTY BOX, AND NOTHING TWICE (2026-10-06, the Fire
// hole): an empty box went to the engine as an empty value, which nothing
// refused, and a second press went to the door, which queued it behind the run
// in flight. Both are refused before the door is asked; the door's own half is
// TestAFireIsRefusedWhileARunIsMoving.
func TestFireRefusesAnEmptyBoxAndASecondPress(t *testing.T) {
	fire := funcOf(page(t, "js/workflows.js"), "  async fire(")
	empty, busy, call := strings.Index(fire, "if (empty.length)"), strings.Index(fire, "if (this.busy)"), strings.Index(fire, "API.fireFlow(")
	if empty < 0 || busy < 0 || call < 0 || empty > call || busy > call {
		t.Fatal("Fire asks the door before refusing an empty box or a second press")
	}
	if !strings.Contains(fire, "!i.value.trim()") || !strings.Contains(fire, "Nothing was fired.") {
		t.Fatal("an empty box is not refused, or the refusal does not say that nothing was fired")
	}
}

// A RETIRED FLOW IS SHOWN SO AND OFFERED NO FIRE (2026-10-10, the core's WHAT'S LEFT
// I2; his card: "Retire both, versions kept"). The door names a retired flow on its
// own flow_list line, with why, and refuses to fire it from any version; the list
// marks the line retired, and the build says why and greys Fire it -- a move the
// door refuses is a move the page does not offer. The words the list reads are the
// door's own, and the key the build reads is the spec's own, both held here.
func TestARetiredFlowIsShownSoAndOfferedNoFire(t *testing.T) {
	wf := page(t, "js/workflows.js")
	list, build := funcOf(wf, "  async list("), funcOf(wf, "  build(")
	if !strings.Contains(list, "/ · retired: /.test(l)") || !strings.Contains(list, `<span class="badge badge-muted">retired</span>`) {
		t.Fatal("the list does not mark a retired flow's line")
	}
	if !strings.Contains(build, "const retired = String(s.retired || '').trim();") ||
		!strings.Contains(build, `id="wf-fire" ${s.version && !retired ? '' : 'disabled'}`) ||
		!strings.Contains(build, "Retired: ${esc(retired)}") {
		t.Fatal("the build offers Fire on a retired flow, or does not say why it is retired")
	}
	door, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "tools.go"))
	if err != nil {
		t.Fatalf("the door's tool table is not beside the glass: %v", err)
	}
	body := string(door)
	i := strings.Index(body, "func toolFlowList(")
	if i < 0 {
		t.Fatal("the door's flow_list is not where this stroke reads it")
	}
	body = body[i:]
	if j := strings.Index(body, "\nfunc "); j >= 0 {
		body = body[:j]
	}
	if !strings.Contains(body, `fmt.Fprintf(&b, " · retired: %s", strings.TrimSpace(s.Retired))`) {
		t.Fatal("the door no longer names a retired flow on its line in the words the page reads")
	}
	spec, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "flow", "flow.go"))
	if err != nil || !strings.Contains(string(spec), "Retired string `json:\"retired,omitempty\"`") {
		t.Fatalf("the spec no longer carries the key the build reads, retired (%v)", err)
	}
}

// A REFUSAL REACHES THE PAGE IN WORDS (2026-10-06). The handlers answer a
// failure as {"error": words}, and the page's client threw "HTTP 502" and
// dropped them, so every refusal the door gave a page arrived as a number.
func TestARefusalReachesThePageInWords(t *testing.T) {
	api := page(t, "js/api.js")
	if !strings.Contains(funcOf(api, "  async fail("), ".error") {
		t.Fatal("the page's client no longer reads the words a refusal is answered in")
	}
	for _, sig := range []string{"  async get(", "  async post("} {
		if !strings.Contains(funcOf(api, sig), "throw await this.fail(r)") {
			t.Fatalf("%s throws a number in place of the refusal's words", strings.TrimSpace(sig))
		}
	}
	hand, err := os.ReadFile(filepath.Join("..", "handlers", "handlers.go"))
	if err != nil || !strings.Contains(string(hand), `map[string]string{"error": msg}`) {
		t.Fatalf("the handlers no longer answer a failure as {\"error\": words}, the shape the page reads (%v)", err)
	}
}

// THE AUDIT LEDGER TAB CARRIES RECORDS, EVALS AND THE LIVE CHECK (his word,
// 2026-10-06: "go, next: audit ledger"). What the build has proved and the
// estate behind it are drawn by the console's one reader of the proof, from the
// one read the tab makes -- the Dashboard's ledger rows were split out of its
// deck so both drew one function, and since the Dashboard retired (2026-10-07)
// the tab is the rows' one placement; the live check, once the Dashboard's own
// tool, is run only from the button he presses and offered only while no engine
// is open, since the door refuses it then; and the hand-scored evals are read as
// Evals reads them.
func TestTheAuditLedgerCarriesTheProofAndTheLiveCheck(t *testing.T) {
	src := page(t, "js/agent.js")
	app := page(t, "js/app.js")
	css := page(t, "css/agent.css")
	led := funcOf(src, "  async ledgerHtml(")
	if strings.Count(led, "App.tool('proofs'") != 1 || !strings.Contains(led, `id="ag-proof-led"`) || !strings.Contains(led, `id="ag-proof-estate"`) {
		t.Fatal("the ledger does not read the proof once and leave the two places it is drawn")
	}
	pane := funcOf(src, "  async paintPane(")
	for _, want := range []string{"App.paintProof('ag-proof-led', 'ledger', this._proofs)", "App.paintProof('ag-proof-estate', 'estate', this._proofs)"} {
		if !strings.Contains(pane, want) {
			t.Fatalf("the ledger tab does not draw the proof with the console's own reader (wanted %s)", want)
		}
	}
	if !strings.Contains(funcOf(app, "  async paintProof("), "only === 'ledger' ? this.proofLedger(p)") {
		t.Fatal("the ledger tab does not draw the proof's rows with the console's one definition of them, proofLedger")
	}
	if strings.Count(src, "App.tool('standup_run'") != 1 ||
		!regexp.MustCompile(`(?s)act === 'livecheck'.{0,400}App\.tool\('standup_run', \{ set: 'morning' \}\)`).MatchString(funcOf(src, "  async paneClick(")) {
		t.Fatal("the live check can be run from somewhere other than its button, or not as the morning set")
	}
	if !regexp.MustCompile(`data-act="livecheck"\$\{[^}]*Run\.engineOpen`).MatchString(led) {
		t.Fatal("the live check is offered while an engine is open, which the door refuses")
	}
	if !strings.Contains(led, "API.listEvals()") {
		t.Fatal("the hand-scored evals are not read")
	}
	if !regexp.MustCompile(`\.ag-proof \{[^}]*--muted:`).MatchString(css) {
		t.Fatal("the proof in the panel does not take the page's look")
	}
}

// REGISTRY & DOCS CARRIES TOOLS WITH CALL, THE SEATS AND RECORDS' DOCUMENTS (his
// word, 2026-10-06: "go, next: registry & docs"; on the card, the seats as "The
// Agents page's own cards" and the documents with what Records shows). A tool is
// called through the Tools page's own form and from nowhere else in the tab; the
// seats are the Agents page's cards, told they stand in the panel so a name cannot
// leave the front page, with the Agents page's own filter; and a document row says
// when it changed and whether it is sealed, every document of the kind listed.
func TestRegistryAndDocsCarriesTheToolsTheSeatsAndTheDocuments(t *testing.T) {
	src := page(t, "js/agent.js")
	app := page(t, "js/app.js")
	docs := funcOf(src, "  async docsHtml(")
	if !strings.Contains(docs, "App.seatCard(s, i, true)") || !strings.Contains(docs, `id="agent-grid"`) {
		t.Fatal("the tab does not draw the seats with the Agents page's own card")
	}
	card := funcOf(app, "  seatCard(")
	if !regexp.MustCompile(`(?s)const title = inPanel\s+\? `+"`"+`<span class="card-title">`).MatchString(card) || !strings.Contains(card, `data-act="seat-file"`) {
		t.Fatal("a seat card in the panel still links away from the front page, or has no way to its file")
	}
	click := funcOf(src, "  async paneClick(")
	if !strings.Contains(funcOf(src, "  toolRows("), `data-act="tool-call"`) ||
		!strings.Contains(click, "App.invokeTool(b.dataset.tool)") || strings.Count(src, "App.invokeTool(") != 1 {
		t.Fatal("a tool is called from the tab by some way other than the Tools page's own Call")
	}
	if !strings.Contains(src, "else if (e.target.id === 'ag-seats-q') App.filterAgents(e.target.value)") {
		t.Fatal("the seats' filter is not the Agents page's own")
	}
	i := strings.Index(click, "act === 'kind'")
	if i < 0 {
		t.Fatal("the documents' kind list is gone; this stroke is reading the wrong thing")
	}
	kind := click[i:]
	if !strings.Contains(kind, "x.sealed") || !strings.Contains(kind, "when(Date.parse(x.modified)") || strings.Contains(kind, "documents.slice(") {
		t.Fatal("a document row lacks Records' changed date or sealed badge, or the kind is cut short")
	}
}

// SETTINGS IS A DOOR OF ITS OWN (his word, 2026-10-06: "Add another page along the
// top next to guardrails that is for settings. that'll include all the settings
// for messaging integration."): beside Laws (Guardrails until 2026-10-08) in the
// bar and among the Inspector's tabs; the door's address and the pass mark read from the glass's own
// store and saved only by the button beside each, in a shape it can hold; the
// messaging bridge as the door reports it; none of the old page's written-in
// badges; and since 2026-10-10 (his cards: "Every page with a tab", "Remove what
// nothing reads") the page's address opens the tab, and the old page is gone.
func TestSettingsIsADoorOfItsOwn(t *testing.T) {
	src := page(t, "js/agent.js")
	if !regexp.MustCompile(`(?s)DOORS: \[.*?\{ id: 'laws',   label: 'Laws' \},\s+\{ id: 'settings', label: 'Settings' \}\s+\],`).MatchString(src) ||
		!regexp.MustCompile(`(?s)TABS: \[.*?\{ id: 'settings', label: 'Settings' \}`).MatchString(src) {
		t.Fatal("Settings is not a door beside Laws, or not one of the Inspector's tabs")
	}
	if !strings.Contains(funcOf(src, "  async paintPane("), "else if (id === 'settings') put(await this.settingsHtml())") {
		t.Fatal("the Settings tab draws nothing")
	}
	set := funcOf(src, "  async settingsHtml(")
	for _, want := range []string{"read('mcp_url')", "read('eval_threshold')", "API.teamStatus()", `data-act="set"`} {
		if !strings.Contains(set, want) {
			t.Fatalf("the Settings tab lost %s", want)
		}
	}
	for _, bad := range []string{"can_approve", "Append-Only", "SHA-256"} {
		if strings.Contains(set, bad) {
			t.Fatalf("the Settings tab carries the old page's written-in badge %q", bad)
		}
	}
	click := funcOf(src, "  async paneClick(")
	if !strings.Contains(click, "API.setSetting(b.dataset.key, v)") || !strings.Contains(click, "if (!ok)") {
		t.Fatal("a setting is saved without its shape being checked first")
	}
	app := page(t, "js/app.js")
	if strings.Contains(app, "renderSettings") || strings.Contains(app, "saveSetting(") || !strings.Contains(app, "settings: 'settings'") {
		t.Fatal("the old Settings page is drawn still, or its address does not open the Settings tab")
	}
}

// THE LAWS ARE A TAB OF THEIR OWN (his card, 2026-10-08: "It becomes Laws"; from 2026-10-06, his word "go, next:
// guardrails", they stood under the Guardrails tab's hold queue, which is the dock at the foot of the terminal now).
// The tab draws the Laws page's own object -- one Laws, so the two cannot drift -- and its Read it again, Set and Seal
// draw it again where it stands, never sending him to a page; and since 2026-10-10 (his cards: "Every page with a
// tab") the page's address opens the tab, the one place the laws are drawn.
func TestTheLawsTabCarriesTheLaws(t *testing.T) {
	src := page(t, "js/agent.js")
	laws := page(t, "js/laws.js")
	app := page(t, "js/app.js")
	css := page(t, "css/agent.css")
	pane := funcOf(src, "  async paintPane(")
	at := strings.Index(pane, "else if (id === 'laws') {")
	if at < 0 {
		t.Fatal("the Laws tab draws nothing")
	}
	branch := pane[at:]
	if k := strings.Index(branch, "else if (id === 'docs')"); k >= 0 {
		branch = branch[:k]
	}
	if !strings.Contains(branch, `put('<div class="ag-laws" id="ag-laws"></div>')`) ||
		!strings.Contains(branch, "await Laws.render($ag('ag-laws'))") || strings.Count(src, "Laws.render(") != 1 {
		t.Fatal("the Laws tab does not draw the Laws page's own object, or something else on the front page draws it too")
	}
	if strings.Contains(laws, "inPanel") || strings.Contains(funcOf(laws, "  paint("), "page-header") {
		t.Fatal("the laws still ask whether they stand in the panel, or draw a page's title")
	}
	if strings.Count(laws, "this.render(el)") != 3 {
		t.Fatal("Read it again, Set or Seal draws the laws somewhere other than where they stand")
	}
	if !strings.Contains(app, "laws: 'laws'") || strings.Contains(app, "case 'laws':") {
		t.Fatal("the Laws page's address does not open the tab, or the router draws the page again")
	}
	if !regexp.MustCompile(`\.ag-laws \{[^}]*--muted:`).MatchString(css) || !strings.Contains(css, ".ag-laws .btn {") {
		t.Fatal("the laws in the panel do not take the page's look")
	}
}

// GITHUB IS A TAB OF ITS OWN, DRAWN BY THE PAGE'S OWN OBJECT (his word, 2026-10-08: "it should be in its own github tab
// on the inspector so its less confusing", "let's make sure we are using the existing infrastructure to underpin the new
// tab", and "let's rename it to Github"). Version control stood under Aider Pair from 2026-10-06 ("go, next: aider
// pair"); now the page is GitHub, under Pages and in the Inspector's own tab beside Aider Pair, and the tab draws the
// page's own object -- the repositories alone, one object so the two cannot drift -- with
// the whole-world diff under it. Aider Pair draws Aider alone and points the way in one line, nothing on the front page
// sends him to Version control by its old name, and since 2026-10-10 (his cards: "Every page with a tab", "Remove
// what nothing reads") the page's address opens the tab: the object is drawn there alone, and what only the page drew
// went with it.
func TestGitHubIsATabOfItsOwn(t *testing.T) {
	src := page(t, "js/agent.js")
	fl := page(t, "js/flows.js")
	app := page(t, "js/app.js")
	css := page(t, "css/agent.css")
	i := strings.Index(src, "  TABS: [")
	if i < 0 {
		t.Fatal("the Inspector's tabs are gone; this stroke is reading the wrong thing")
	}
	tabs := src[i:]
	if j := strings.Index(tabs, "  ],"); j >= 0 {
		tabs = tabs[:j]
	}
	if a, g := strings.Index(tabs, "{ id: 'aider',  label: 'Aider Pair' },"), strings.Index(tabs, "{ id: 'github', label: 'GitHub' },"); a < 0 || g < a {
		t.Fatal("GitHub is not one of the Inspector's tabs, beside Aider Pair")
	}
	pane := funcOf(src, "  async paintPane(")
	at := strings.Index(pane, "else if (id === 'github') {")
	if at < 0 {
		t.Fatal("the GitHub tab draws nothing")
	}
	branch := pane[at:]
	if k := strings.Index(branch, "else if (id === 'ledger')"); k >= 0 {
		branch = branch[:k]
	}
	gh := funcOf(src, "  async githubHtml(")
	if !strings.Contains(branch, "put(await this.githubHtml())") || !strings.Contains(branch, "await Flows.render($ag('ag-vc'))") ||
		!strings.Contains(gh, `id="ag-vc"`) || strings.Count(src, "Flows.render(") != 1 {
		t.Fatal("the GitHub tab does not draw the GitHub page's own object, or something else on the front page draws it too")
	}
	if !strings.Contains(gh, `data-act="diff"`) || !strings.Contains(funcOf(src, "  async paneClick("), "App.tool('git_diff', { project: b.dataset.world })") {
		t.Fatal("the whole-world diff he kept is gone")
	}
	aider := funcOf(src, "  async aiderHtml(")
	if strings.Contains(aider, `id="ag-vc"`) || strings.Contains(aider, `data-act="diff"`) || !strings.Contains(aider, `data-act="door" data-door="github"`) ||
		!strings.Contains(funcOf(src, "  async paneClick("), "if (act === 'door') { this.door(b.dataset.door); return; }") {
		t.Fatal("Aider Pair still draws GitHub, or its one line no longer opens the GitHub tab")
	}
	for _, old := range []string{"open a line of work (Version control)", "save it in Version control", "until Version control saves it", ">Version control</button>"} {
		if strings.Contains(src, old) {
			t.Fatalf("the front page still sends him to Version control by its old name (%s)", old)
		}
	}
	render := funcOf(fl, "  async render(")
	if !strings.Contains(render, "await this.repos();") || strings.Contains(fl, "inPanel") || strings.Contains(render, "page-title") {
		t.Fatal("GitHub draws more than its repositories, or still asks whether it stands in the panel")
	}
	if !strings.Contains(app, "{ page: 'flows', href: '/flows', label: 'GitHub' },") || !strings.Contains(app, "flows: 'github'") ||
		strings.Contains(app, "case 'flows':") {
		t.Fatal("GitHub is not on the page list, or its address does not open the tab, or the router draws the page again")
	}
	// WHAT ONLY THE PAGE DREW WENT WITH IT (2026-10-10), each to the front page's own: the hold queue to the dock at the
	// foot of the terminal, the way through the council to the terminal itself, Recent to the terminal's own history.
	for _, gone := range []string{"async holds(", "answerHold(", "async readGit(", "paintRecent(", "flow-council", "holds-card"} {
		if strings.Contains(fl, gone) {
			t.Fatalf("the GitHub object still carries %s, which only the page drew", gone)
		}
	}
	if strings.Contains(page(t, "js/chat.js"), "  pending: ''") || strings.Contains(src, "Chat.pending") ||
		strings.Contains(page(t, "css/app.css"), ".home-recent") {
		t.Fatal("Recent went with the page, and what it left behind is still here: Chat.pending, or its styles")
	}
	if !regexp.MustCompile(`\.ag-vc \{[^}]*--muted:`).MatchString(css) || !strings.Contains(css, ".ag-vc .btn {") {
		t.Fatal("GitHub in the panel does not take the page's look")
	}
}

// RUN CARRIES THE WATCHBOARD'S WIRE AND THE DASHBOARD'S PROJECTS (his word, 2026-10-06: "go, next: run"; on his
// cards, 2026-10-07, the wire replaces the tab's own list of events, and a pick plays "Beside the terminal"). The
// turn's events are drawn by the Watchboard's own wire, handed the tab's box, with its one raw switch; the projects
// are the Projects card's own object, with no frame of its own -- a pick plays in the front page's one play frame --
// and its words go to the council through the terminal's Agent tab; an event draws the turn again alone, never the
// projects; and the Watchboard still draws its wire whole. Since the Dashboard retired (2026-10-07) the Run tab is
// the one place the Projects card stands.
func TestRunCarriesTheWireAndTheProjects(t *testing.T) {
	src := page(t, "js/agent.js")
	chat := page(t, "js/chat.js")
	proj := page(t, "js/projects.js")
	css := page(t, "css/agent.css")
	if !strings.Contains(funcOf(src, "  paintRun("), "Chat.paintWire($ag('ag-wire'), $ag('ag-wire-n'))") || strings.Contains(src, "evRow") {
		t.Fatal("the Run tab draws the turn's events its own way, not with the Watchboard's own wire")
	}
	if !strings.Contains(chat, "paintWire(box = document.getElementById('wb-wire'), n = document.getElementById('wb-wire-n')) {") {
		t.Fatal("the Watchboard's wire cannot be handed a box, or no longer draws its own when it is handed none")
	}
	if !strings.Contains(funcOf(src, "  async paneClick("), "if (act === 'wire-raw') { Chat.raw = !Chat.raw; this.paintRun(); return; }") ||
		!strings.Contains(chat, "type=\"button\">${this.raw ? 'folded' : 'raw'}</button>") {
		t.Fatal("the Run tab's raw switch is not the Watchboard's own, or the Watchboard's button does not say which way it stands")
	}
	pane := funcOf(src, "  async paintPane(")
	if !strings.Contains(pane, `id="ag-proj"`) || !strings.Contains(pane, "await Projects.render($ag('ag-proj'))") {
		t.Fatal("the Run tab does not draw the Projects card's own object")
	}
	sched := funcOf(src, "  schedulePane(")
	if !strings.Contains(sched, "this.paintRun()") || strings.Contains(sched, "this.paintPane()") ||
		!strings.Contains(funcOf(src, "  onRun("), "if (this.open && this.tab === 'run') this.paintRun();") {
		t.Fatal("a turn's event draws the whole Run tab again, so the projects would be read again with every token")
	}
	if strings.Contains(proj, "inPanel") || strings.Contains(proj, "<iframe") || !strings.Contains(funcOf(proj, "  paint("), "it plays beside the terminal") {
		t.Fatal("the Projects card draws a frame of its own, or still asks where it stands, where a pick plays beside the terminal")
	}
	if !strings.Contains(funcOf(proj, "  say("), "Agent.council(words);") ||
		!regexp.MustCompile(`(?s)this\.mode = 'agent';.*this\.go\(words\);`).MatchString(funcOf(src, "  council(")) {
		t.Fatal("Work on this and Put it down in the Run tab do not go to the council as if typed in the Agent tab")
	}
	if !strings.Contains(funcOf(proj, "  play("), "if (this.selected) Agent.pickProject(this.selected, this.version);") ||
		strings.Count(proj, "Agent.pickProject(") != 1 || strings.Count(src, "pickProject(") != 1 ||
		strings.Count(src, "this.play = { name, version, closed: false };") != 2 {
		t.Fatal("a project reaches the play frame from something other than a turn's delivery or his pick in the Run tab")
	}
	if !regexp.MustCompile(`\.ag-proj \{[^}]*--muted:`).MatchString(css) || !strings.Contains(css, ".ag-proj .btn {") ||
		!strings.Contains(css, ".ag-wire {") || strings.Contains(css, ".ag-ev") {
		t.Fatal("the projects or the wire in the panel do not take the page's look, or the old list of events left its rules behind")
	}
}

// THE FRONT PAGE (2026-10-02). The operator, shown a page out of his own AI
// Studio project: "that's what I am looking for", then, asked where it should
// go, "replace the Dashboard now". The strokes below hold what makes it his
// page and not that project's: it is routed at /; it asks the door only for
// tools the door carries, and asks quietly only where the glass will let it; it
// fakes nothing; and a tab that is not wired says so when it is pressed. The
// binary embeds these files and builds whether or not a page keeps any of this.
//
// THE SIDEBAR AND THE OLD DASHBOARD RETIRED (his word, 2026-10-07: "go, next:
// retire the sidebar and the old dashboard"). The launchpad kept whole at
// /dashboard is gone and the address opens the front page; the sidebar's links
// became one list (App.PAGES) that the front page's Pages menu, the palette and
// each page's crumb read, and every line on it is a page the router draws or, since
// 2026-10-10, a tab it opens; Boot and Close, which the sidebar carried, act through
// the front page.
func TestTheFrontPageIsRoutedAndTheSidebarIsRetired(t *testing.T) {
	index := page(t, "index.html")
	app := page(t, "js/app.js")
	pal := page(t, "js/palette.js")
	src := page(t, "js/agent.js")
	for _, want := range []string{
		`<link rel="stylesheet" href="/css/agent.css">`,
		`<script src="/js/agent.js"></script>`,
	} {
		if !strings.Contains(index, want) {
			t.Fatalf("index.html no longer carries %s", want)
		}
	}
	for _, bad := range []string{`class="sidebar"`, `class="nav-link`, `/js/home.js`, `/css/icons.css`} {
		if strings.Contains(index, bad) {
			t.Fatalf("index.html still carries %s, which retired with the sidebar and the Dashboard", bad)
		}
	}
	router := funcOf(app, "  router(")
	for _, want := range []string{`|| 'agent'`, `classList.toggle('is-agent'`, `if (location.pathname === '/dashboard') history.replaceState(null, '', '/');`} {
		if !strings.Contains(router, want) {
			t.Fatalf("the router no longer sends / and /dashboard to the front page (wanted %s)", want)
		}
	}
	render := funcOf(app, "  async render(")
	if !strings.Contains(render, `case 'agent': await Agent.render(el)`) {
		t.Fatal("the page table lost the front page")
	}
	if strings.Contains(render, `case 'dashboard'`) || strings.Contains(app, "Home.render(") {
		t.Fatal("the page table still draws the Dashboard")
	}
	// EVERY LINE ON THE LIST IS A PAGE THE ROUTER DRAWS OR A TAB IT OPENS, at the address the router reads it from: a line
	// that is neither goes to Page not found from the Pages menu, the palette and the crumb alike. A PAGE WITH A TAB OPENS
	// ITS TAB (2026-10-10, his word: "everything stays within the new interface"; his cards: "Every page with a tab" and
	// "Keep them, opening the tab"): Workflows, Laws, GitHub and Settings stay on the list, the router draws none of them,
	// and each address opens the front page with that Inspector tab open, the way /dashboard opens the front page.
	list := regexp.MustCompile(`(?s)  PAGES: \[(.*?)\n  \],`).FindStringSubmatch(app)
	if list == nil {
		t.Fatal("the one page list (App.PAGES) is gone")
	}
	pages := regexp.MustCompile(`\{ page: '([a-z]+)', href: '(/[a-z]*)', label: `).FindAllStringSubmatch(list[1], -1)
	if len(pages) < 10 {
		t.Fatalf("read %d pages off App.PAGES; this stroke is reading the wrong thing", len(pages))
	}
	tabbed := map[string]string{}
	if m := regexp.MustCompile(`  TABBED: \{([^}]*)\},`).FindStringSubmatch(app); m != nil {
		for _, kv := range regexp.MustCompile(`([a-z]+): '([a-z]+)'`).FindAllStringSubmatch(m[1], -1) {
			tabbed[kv[1]] = kv[2]
		}
	}
	for pg, tab := range map[string]string{"workflows": "flows", "laws": "laws", "flows": "github", "settings": "settings"} {
		if tabbed[pg] != tab {
			t.Fatalf("/%s does not open the Inspector's %s tab (App.TABBED gives %q)", pg, tab, tabbed[pg])
		}
	}
	tabs := regexp.MustCompile(`(?s)  TABS: \[(.*?)\n  \],`).FindStringSubmatch(src)
	if tabs == nil {
		t.Fatal("the Inspector's tabs are gone; this stroke is reading the wrong thing")
	}
	for _, p := range pages {
		tab, opens := tabbed[p[1]]
		if drawn := strings.Contains(render, "case '"+p[1]+"':"); drawn == opens {
			t.Fatalf("App.PAGES offers %s, and the router draws it and opens a tab for it both, or neither", p[1])
		}
		if opens && !strings.Contains(tabs[1], "{ id: '"+tab+"',") {
			t.Fatalf("/%s opens the Inspector's %s tab, which the Inspector does not have", p[1], tab)
		}
		href := "/" + p[1]
		if p[1] == "agent" {
			href = "/"
		}
		if p[2] != href {
			t.Fatalf("App.PAGES sends %s to %s, which the router reads as another page", p[1], p[2])
		}
	}
	for _, want := range []string{
		"const tab = Object.hasOwn(this.TABBED, first) ? this.TABBED[first] : '';",
		"history.replaceState(null, '', '/');\n      if (this.currentPage === 'agent' && document.getElementById('ag')) { Agent.door(tab); return; }",
		"Agent.open = true;\n      Agent.tab = tab;",
	} {
		if !strings.Contains(router, want) {
			t.Fatalf("the router does not open the front page on a page's tab (wanted %s)", want)
		}
	}
	for _, c := range []struct{ who, in, want string }{
		{"the front page's Pages menu", funcOf(src, "  menu("), "App.PAGES.filter("},
		{"the palette", funcOf(pal, "  build("), "App.PAGES.map("},
		{"each page's crumb", funcOf(app, "  paintCrumb("), "this.PAGES.find("},
	} {
		if !strings.Contains(c.in, c.want) {
			t.Fatalf("%s no longer reads the one page list (wanted %s)", c.who, c.want)
		}
	}
	// BOOT AND CLOSE ARE THE FRONT PAGE'S. They navigate to / and then act; if
	// they went on calling Home, the engine would open under a page that has
	// no boot report to show, and say nothing.
	for _, want := range []string{"Agent.boot()", "Agent.closeSitting()"} {
		if !strings.Contains(pal, want) {
			t.Fatalf("the palette no longer calls %s", want)
		}
	}
	for _, bad := range []string{"Home.boot()", "Home.closeSitting()"} {
		if strings.Contains(pal, bad) || strings.Contains(app, bad) {
			t.Fatalf("%s is still called, and there is no Dashboard under it", bad)
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
	// Two calls answer a hold, and each is a button he presses on a card pinned at the foot of the terminal: a shell
	// card of this thread (decide, which lands the run in its entry) and any other call the door parked (holdAnswer),
	// each called once, from the dock's handler for its card's two buttons (which button reaches which is held by
	// TestWhatWaitsForHisHandIsPinnedAtTheFootOfTheTerminal).
	if strings.Count(src, "App.tool('hold_answer'") != 2 ||
		!strings.Contains(funcOf(src, "  async holdAnswer("), "App.tool('hold_answer'") ||
		!strings.Contains(funcOf(src, "  async decide("), "App.tool('hold_answer'") {
		t.Fatal("a held call can be answered from somewhere other than a button he presses")
	}
	if strings.Count(src, "this.decide(") != 1 || strings.Count(src, "this.holdAnswer(") != 1 {
		t.Fatal("a card can be answered from somewhere other than the click handler of its two buttons")
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
	for _, fn := range []string{"  shellEntryHtml(", "  shellCardHtml(", "  shellTake(", "  async decide("} {
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

	// A card is drawn for a held answer and for nothing else, pinned in the dock and not in the entry, which says where
	// it is; a refusal draws no button; a failure is on the face (a non-zero exit is a bad chip and a bad box).
	face := funcOf(src, "  shellEntryHtml(")
	card := funcOf(src, "  shellCardHtml(")
	if strings.Contains(face, "shell-yes") || strings.Contains(face, "shell-no") || !strings.Contains(face, "its card is pinned at the foot of the window") {
		t.Fatal("the entry draws a card's buttons itself, or does not say where its card is")
	}
	if strings.Count(card, "shell-yes") != 1 || strings.Count(card, "shell-no") != 1 || strings.Count(src, `data-act="shell-yes"`) != 1 ||
		!strings.Contains(funcOf(src, "  paintDock("), "en.state === 'held'") {
		t.Fatal("the card does not carry its two buttons, is drawn for something other than a held answer, or a shell button is drawn outside it")
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
	// A card answered somewhere else (another window) stops offering buttons that could only be refused. The GitHub page
	// answered them too until it retired into its tab (2026-10-10, his cards: "Every page with a tab"), and the words that
	// sent him there went with it.
	if strings.Contains(src, "GitHub page under Pages") {
		t.Fatal("a settled card still sends him to the GitHub page under Pages, whose address opens the GitHub tab now")
	}
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

// A LINE FROM A SCRIPT IS A TYPED LINE (2026-10-07; his word: "write those helpers in so you don't need to rewrite them
// every time"). The hand drives the terminal from the page's own script, and the lines it did that with were lost with every
// reload; Agent.drive keeps them. It turns to the tab as the tab's own button does, hands the line to go() as if typed, and
// waits on the entry that line made, so every gate, card, refusal and record is the typed line's own; it reaches the door by
// no road of its own and answers no card.
func TestALineFromAScriptIsATypedLine(t *testing.T) {
	src := page(t, "js/agent.js")
	drive := funcOf(src, "  async drive(")
	if drive == "" {
		t.Fatal("the page has no drive(); a script would have to carry its own way into the terminal again")
	}
	if !strings.Contains(drive, "this.go(line)") || !strings.Contains(drive, "m.wired") {
		t.Fatal("drive() does not hand the line to go() as typed, or turns to a tab that is not wired")
	}
	for _, bad := range []string{"App.tool(", "API.", "fetch(", "hold_answer", "this.decide(", "this.send(", "this.shellRun(", "this.aiderRun(", "Run.start("} {
		if strings.Contains(drive, bad) {
			t.Fatalf("drive() reaches past the typed line's own road (%s)", bad)
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

	// The frame is the glass's own sandboxed page, framed as the Dashboard framed it: the url is Projects.pageUrl's, there is no `sandbox` attribute (the app's
	// own browser pane refuses any frame that carries one, measured 2026-09-21; the header the page is served under is the wall), and no referrer goes out.
	if !strings.Contains(paint, "Projects.pageUrl(s.name, s.version)") || strings.Contains(paint, "sandbox") || !strings.Contains(paint, `referrerpolicy="no-referrer"`) {
		t.Fatal("the frame is not the glass's own sandboxed project page, framed as the Dashboard framed it")
	}

	// The project named is the delivery's and nobody else's, as a string, only when it names one.
	if !regexp.MustCompile(`const made = Run\.turn && Run\.turn\.delivery;\s+if \(made && typeof made\.project === 'string' && made\.project\) this\.playProject\(made\.project\);`).MatchString(onRun) {
		t.Fatal("what is played is not read off the turn's delivery (`project`), or is played when the delivery names none")
	}
	if strings.Count(src, "this.playProject(") != 1 {
		t.Fatal("playProject is called from somewhere other than a turn's delivery (his pick in the Run tab has its own way, held by TestRunCarriesTheWireAndTheProjects)")
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

// THE FRONT PAGE IS THE WINDOW. His word, 2026-10-03, watching the app's browser pane at 888 px: "the window should be locked to
// that scale it starts at and just scroll the terminal instead of the whole page"; and 2026-10-07, having sized his pane to the
// layout he wanted: "it should fit the browser window it sits in naturally", with "Fill, no scaling" on his card. The page takes
// exactly the window's height and never scrolls, either way; the terminal and the Inspector share what the bar and the footer
// leave, with no floor that could push past the window's edge, and scroll inside themselves; and below 1024 px the Inspector
// lies over the terminal's side instead of stacking under it, which made the page two windows tall. The page itself was measured
// in the browser at his size and two smaller ones (the CHANGELOG has the numbers), because a rule read as text can be in the
// file and still lose.
func TestTheFrontPageIsTheWindow(t *testing.T) {
	css := page(t, "css/agent.css")
	rule := func(sel string) string {
		m := regexp.MustCompile(`(?s)\n` + regexp.QuoteMeta(sel) + ` \{([^}]*)\}`).FindStringSubmatch(css)
		if m == nil {
			t.Fatalf("the stylesheet has no %s rule; this stroke is reading the wrong thing", sel)
		}
		return m[1]
	}
	for _, c := range []struct{ sel, want, why string }{
		{"body.is-agent .main", "height: 100vh", "the front page is not the window's height"},
		{"body.is-agent .main", "overflow: hidden", "the front page can scroll"},
		{".ag", "height: 100vh", "the page's frame is not the window's height"},
		{".ag", "overflow: hidden", "the page's frame can grow past the window"},
		{".ag-main", "min-height: 0", "the terminal's row has a floor and can push the footer past the window"},
		{".ag-main", "align-items: stretch", "the terminal and the Inspector do not fill the row"},
		{".ag-main", "position: relative", "the Inspector has no frame to lie over the terminal in"},
		{".ag-win, .ag-panel", "min-height: 0", "the terminal or the Inspector has a floor"},
	} {
		if !strings.Contains(rule(c.sel), c.want) {
			t.Errorf("%s (%s wants %s)", c.why, c.sel, c.want)
		}
	}
	if b := rule(".ag-win, .ag-panel"); strings.Contains(b, "100vh") || strings.Contains(b, "560px") {
		t.Error("the terminal and the Inspector are given a height of their own again, which a short window cannot hold")
	}
	i := strings.Index(css, "@media (max-width: 1023px) {")
	if i < 0 {
		t.Fatal("the narrow-screen block (max-width: 1023px) is gone")
	}
	block := css[i:]
	if j := strings.Index(block, "\n}\n"); j >= 0 {
		block = block[:j]
	}
	if strings.Contains(block, "flex-direction: column") || !strings.Contains(block, ".ag-panel { position: absolute;") {
		t.Error("below 1024 px the Inspector stacks under the terminal again instead of lying over its side")
	}
	if !regexp.MustCompile(`(?s)\.ag-out \{[^}]*overflow-y: auto;`).MatchString(css) {
		t.Error("the output box no longer scrolls (overflow-y: auto); the window would clip the terminal instead of scrolling it")
	}
}

// WHAT WAITS FOR HIS HAND IS PINNED AT THE FOOT OF THE TERMINAL (his word, 2026-10-08: "Same as the cards in the chat,
// they are perfect, just pin them to the window so they actually stay at the botton if there is an approve card. should
// work on all tabs, agent, bash, python, powershell, aider. Basically replacing the Guardrails tab"; on his card, "It
// becomes Laws"). Read fast, the output scrolled the cards out of his view. The council's question, every shell line
// that asks first and every other call the door has parked are drawn in one dock: a sibling of the output, between it
// and the line he types in, outside what the terminal rewrites and outside any one tab; nothing draws a card in the
// thread any more. The dock draws what was read and asks the door for nothing itself; the queue is read quietly every
// HOLD_POLL_MS while the page is in view, so a call parked by anything but this page reaches it while he looks; its
// buttons and the question's field are answered from the dock's own handler, on his click; the keys it reads off the
// queue are the keys hold_list writes; a refusal or a call that errored is said as a failure; and Guardrails is gone.
func TestWhatWaitsForHisHandIsPinnedAtTheFootOfTheTerminal(t *testing.T) {
	src := page(t, "js/agent.js")
	css := page(t, "css/agent.css")
	fn := func(sig string) string {
		body := funcOf(src, sig)
		if body == "" {
			t.Fatalf("the page has no %s", strings.TrimSpace(sig))
		}
		return body
	}
	shell, dock, wire, answer := fn("  shell("), fn("  paintDock("), fn("  wire("), fn("  async holdAnswer(")

	// A sibling of the output, under it and the play frame and over the line he types in.
	out, play, at, hints, in := strings.Index(shell, `id="ag-out"`), strings.Index(shell, `id="ag-play"`), strings.Index(shell, `id="ag-dock"`),
		strings.Index(shell, `id="ag-hints"`), strings.Index(shell, `<div class="ag-in">`)
	if out < 0 || play < out || at < play || hints < at || in < hints {
		t.Fatal("the dock is not between the terminal's output and the line he types in")
	}
	// Nothing draws a card in what scrolls.
	for _, sig := range []string{"  paintOut(", "  entryHtml(", "  shellEntryHtml(", "  aiderEntryHtml(", "  bannerHtml("} {
		for _, bad := range []string{"this.gateHtml()", "ag-gate", "shell-yes", "hold-yes"} {
			if strings.Contains(fn(sig), bad) {
				t.Fatalf("%s draws a card in the output, which reading fast scrolls out of view (%s)", strings.TrimSpace(sig), bad)
			}
		}
	}
	// The dock draws the question, the thread's held lines and the door's queue, and says when the door holds nothing or
	// its queue could not be read; it asks the door for nothing itself, and is drawn again whenever the thread is and
	// whenever the queue is read.
	for _, want := range []string{"this.gateHtml()", "en.state === 'held'", "this.shellCardHtml(", "this.holds.filter(", "this.holdCardHtml(", "this.holdNoteHtml()", "box.hidden = !html"} {
		if !strings.Contains(dock, want) {
			t.Fatalf("the dock does not draw what waits for his hand (wanted %s)", want)
		}
	}
	for _, bad := range []string{"App.tool(", "API.", "fetch("} {
		if strings.Contains(dock, bad) {
			t.Fatalf("the dock asks for something itself (%s); it draws what was read", bad)
		}
	}
	if !strings.Contains(fn("  paintOut("), "this.paintDock()") || !strings.Contains(fn("  async readHolds("), "this.paintDock()") {
		t.Fatal("the dock is not drawn again when the thread is, or when the queue is read")
	}
	// The queue is read quietly, every HOLD_POLL_MS while the page is in view and again when he comes back to it, and no
	// more once the page has gone.
	if !strings.Contains(src, "HOLD_POLL_MS: 3000,") || !strings.Contains(fn("  async readHolds("), "App.tool('hold_list', {}, true)") {
		t.Fatal("the queue is not read quietly, or how often it is read is not named")
	}
	watch := fn("  watch(")
	for _, want := range []string{"setInterval(holds, this.HOLD_POLL_MS)", "if (!document.hidden) this.readHolds()", "holds(); }", "clearInterval(this._hw)"} {
		if !strings.Contains(watch, want) {
			t.Fatalf("the queue is not read while he looks, or goes on being read after the page has gone (wanted %s)", want)
		}
	}
	// His click, and one handler: a shell card's two buttons go to decide, any other card's to holdAnswer, and the
	// question's field to answer; the output's own handler answers none of them.
	i := strings.Index(wire, "dock.addEventListener('click'")
	if i < 0 {
		t.Fatal("the dock's buttons are not listened for")
	}
	for _, bad := range []string{"shell-yes", "hold-yes", "ag-gate-form"} {
		if strings.Contains(wire[:i], bad) {
			t.Fatalf("the output's own handler still answers a card (%s)", bad)
		}
	}
	for _, c := range []struct{ act, call, not string }{
		{"act === 'shell-yes' || act === 'shell-no'", "this.decide(", "this.holdAnswer("},
		{"act === 'hold-yes' || act === 'hold-no'", "this.holdAnswer(", "this.decide("},
		{"dock.addEventListener('submit'", "this.answer(", "this.decide("},
	} {
		seg := ""
		if j := strings.Index(wire[i:], c.act); j >= 0 {
			seg = wire[i+j:]
			if len(seg) > 240 {
				seg = seg[:240]
			}
		}
		if !strings.Contains(seg, c.call) || strings.Contains(seg, c.not) {
			t.Fatalf("the dock's handler does not send %s to %s and nothing else", c.act, c.call)
		}
	}
	// What the dock reads off the queue are keys hold_list writes, so renaming one at the door goes red here.
	door, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "holds.go"))
	if err != nil {
		t.Fatalf("the door's hold queue is not beside the glass: %v", err)
	}
	holds := string(door)
	a, b := strings.Index(holds, "func toolHoldList("), strings.Index(holds, "func toolHoldAnswer(")
	if a < 0 || b < a {
		t.Fatal("hold_list is not where this stroke reads it")
	}
	keys := map[string]bool{}
	for _, m := range regexp.MustCompile(`"([a-z_]+)"[:\]]`).FindAllStringSubmatch(holds[a:b], -1) {
		keys[m[1]] = true
	}
	read := map[string]bool{}
	for _, sig := range []string{"  async readHolds(", "  holdCardHtml(", "  async holdAnswer("} {
		for _, m := range regexp.MustCompile(`\b(?:d|h)\.([a-z_]+)\b`).FindAllStringSubmatch(fn(sig), -1) {
			read[m[1]] = true
			if !keys[m[1]] {
				t.Fatalf("the dock reads .%s off the hold queue and hold_list writes no such key", m[1])
			}
		}
	}
	for _, must := range []string{"armed", "held", "why_not", "id", "tool", "caller", "project", "when", "args"} {
		if !read[must] {
			t.Fatalf("the dock no longer reads .%s, so what the queue says in it never reaches him", must)
		}
	}
	// A refusal and a call that errored are said as failures, in the words hold_answer begins them with; and a file's
	// change is shown as file_edit sets it out for his card.
	if !strings.Contains(holds, `"Approved %s and %q errored: %v`) || !strings.Contains(holds, `return fmt.Sprintf("Refused: `) ||
		!strings.Contains(answer, `Refused|Approved [^ ]+ and "[^"]*" errored`) || !strings.Contains(answer, "status: bad ? 'ERROR'") {
		t.Fatal("a refusal or a call that errored is not said as a failure, or hold_answer no longer begins them with the words the page reads")
	}
	edit, err := os.ReadFile(filepath.Join("..", "..", "line", "internal", "tools", "edit.go"))
	if err != nil || !strings.Contains(string(edit), `keep["change"] = change`) || !strings.Contains(fn("  holdCardHtml("), "Array.isArray(a.change)") {
		t.Fatalf("a file's change is not shown as file_edit sets it out for his card (%v)", err)
	}
	// Guardrails is gone: no tab, no door, no list of its own, and no dot on the Inspector's button pointing there.
	for _, bad := range []string{"'guards'", "guardsHtml", "ag-holdout", "ag-adot", "label: 'Guardrails'"} {
		if strings.Contains(src, bad) || strings.Contains(css, bad) {
			t.Fatalf("the front page still carries Guardrails (%s)", bad)
		}
	}
	// The dock is capped, scrolls inside itself, gives way before the line he types in does, and is hidden while
	// nothing waits.
	rule := ""
	if k := strings.Index(css, ".ag-dock {"); k >= 0 {
		rule = css[k:]
		if e := strings.Index(rule, "}"); e >= 0 {
			rule = rule[:e]
		}
	}
	for _, want := range []string{"flex: 0 1 auto", "min-height: 0", "max-height: 45%", "overflow-y: auto"} {
		if !strings.Contains(rule, want) {
			t.Fatalf("the dock is not capped, does not scroll inside itself, or could push the line he types in out of the window (wanted %s)", want)
		}
	}
	if !strings.Contains(css, ".ag-dock[hidden] { display: none; }") {
		t.Fatal("the dock is not hidden while nothing waits")
	}
}
