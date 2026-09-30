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
