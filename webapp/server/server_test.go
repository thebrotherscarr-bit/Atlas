// The first strokes on the glass's front door.
//
// `webapp` had one test function in ~2,800 lines. This package holds the
// session gate and the cache validators, and neither had a statement of what
// it promises.
//
// WHAT CANNOT BE REACHED FROM HERE, said plainly rather than worked around:
// the route table is built inside `ListenAndServe`, which then binds a port,
// so the mux itself is not testable without refactoring that function. The two
// things that decide whether a request is ANSWERED AT ALL -- `gated` and
// `etags` -- are both reachable, and they are what these strokes hold.
package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"atlas/webapp/agents"
	"atlas/webapp/db"
	"atlas/webapp/evals"
	"atlas/webapp/handlers"
	"atlas/webapp/messaging"
	"atlas/webapp/search"
	"atlas/webapp/traces"
)

func newHandlers(t *testing.T) *handlers.Handlers {
	t.Helper()
	database, err := db.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	store := traces.NewStore(database)
	reg := agents.NewRegistry(database)
	return handlers.New(store, reg, evals.NewEngine(database, store),
		search.NewEngine(store, reg), messaging.NewBus())
}

// ok is the handler the gate wraps; it records that the request got through.
func ok(hit *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*hit = true
		w.WriteHeader(http.StatusOK)
	})
}

func TestEveryEmbeddedFileGetsAQuotedValidator(t *testing.T) {
	fsys := fstest.MapFS{
		"index.html":  {Data: []byte("<html>")},
		"js/home.js":  {Data: []byte("console.log(1)")},
		"js/same.js":  {Data: []byte("console.log(1)")},
		"css/app.css": {Data: []byte("body{}")},
	}
	tags := etags(fsys)

	if len(tags) != 4 {
		t.Fatalf("want a tag per file, got %d: %v", len(tags), tags)
	}
	for path, tag := range tags {
		// AN ETAG MUST BE QUOTED. An unquoted one is not a valid entity-tag
		// and a browser is free to ignore it, which puts us back where
		// 2026-09-09 started: a tab holding stale bytes across four rebuilds.
		if !strings.HasPrefix(tag, `"`) || !strings.HasSuffix(tag, `"`) {
			t.Fatalf("%s: tag %s is not quoted", path, tag)
		}
	}
	// SAME BYTES, SAME TAG -- it is a content hash, not a path hash, so a
	// rebuild that does not change a file does not bust its cache.
	if tags["js/home.js"] != tags["js/same.js"] {
		t.Fatal("identical bytes must hash to the same tag")
	}
	// and different bytes differ, which is what ends the 304 after a rebuild
	if tags["index.html"] == tags["css/app.css"] {
		t.Fatal("different bytes must hash differently")
	}

	// a directory is not a file and gets no tag
	nested := fstest.MapFS{"a/b/c.js": {Data: []byte("x")}}
	if n := len(etags(nested)); n != 1 {
		t.Fatalf("directories must not be tagged: %v", etags(nested))
	}
}

func TestWithAuthOffTheGateIsInert(t *testing.T) {
	// THIS WAS PRODUCTION UNTIL 2026-09-21: `ConfigureAuth` had no caller, so
	// `authOn` was false for the life of every process. main.go now closes the
	// gate at every start (THE LOCK) -- see the two strokes at the end of this
	// file. The inert path is kept and still pinned: it is what a Handlers
	// nobody configured does, and a test harness builds exactly that.
	s := &Server{handlers: newHandlers(t)}
	for _, path := range []string{"/api/traces", "/ws", "/metrics", "/", "/records"} {
		hit := false
		w := httptest.NewRecorder()
		s.gated(ok(&hit)).ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if !hit {
			t.Fatalf("auth off: %s was gated, and nothing should be", path)
		}
	}
}

func TestWithAuthOnTheGateRefusesTheRightThings(t *testing.T) {
	// THE MECHANISM WORKED BEFORE ANYTHING TURNED IT ON. It was configured by
	// hand here and nowhere else until 2026-09-21, when main.go began closing
	// the gate at every start (THE LOCK). This stroke is what said the gate was
	// already right; it now also holds the lock's own open faces.
	h := newHandlers(t)
	h.ConfigureAuth(true, "", t.TempDir()+"/sessions.json")
	s := &Server{handlers: h}

	// OPEN, because the login page and a health probe must work ungated, and
	// a platform hook is authenticated by its own signature, not a cookie. The
	// lock's three faces are open for the same reason as the login: the lock
	// screen must ask and answer before anyone is signed in.
	for _, path := range []string{"/api/login", "/api/health", "/hooks/slack",
		"/api/lock", "/api/setup", "/api/unlock"} {
		hit := false
		w := httptest.NewRecorder()
		s.gated(ok(&hit)).ServeHTTP(w, httptest.NewRequest("POST", path, nil))
		if !hit {
			t.Fatalf("auth on: %s must stay open", path)
		}
	}

	// THE STATIC FACE IS OPEN TOO -- the login page itself has to load, and it
	// is served by the catch-all, not by an /api route.
	hit := false
	w := httptest.NewRecorder()
	s.gated(ok(&hit)).ServeHTTP(w, httptest.NewRequest("GET", "/login", nil))
	if !hit {
		t.Fatal("auth on: the login page must load ungated")
	}

	// AN API CALL WITH NO SESSION IS REFUSED IN JSON, not redirected -- a
	// fetch() that follows a 302 to an HTML page reports a parse error, and
	// the page then shows nothing rather than "login required".
	hit = false
	w = httptest.NewRecorder()
	s.gated(ok(&hit)).ServeHTTP(w, httptest.NewRequest("GET", "/api/traces", nil))
	if hit {
		t.Fatal("auth on: /api/traces reached the handler with no session")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want a JSON refusal, got %q", ct)
	}
	if !strings.Contains(w.Body.String(), "login required") {
		t.Fatalf("the refusal must say what to do: %q", w.Body.String())
	}

	// /metrics AND /ws ARE NOT OPEN. Metrics carry global counts and the
	// socket carries live engine events; both are behind the gate.
	for _, path := range []string{"/metrics", "/ws"} {
		hit = false
		w = httptest.NewRecorder()
		s.gated(ok(&hit)).ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if hit {
			t.Fatalf("auth on: %s reached the handler with no session", path)
		}
	}

	// a page request redirects to the login page rather than 401-ing a human
	hit = false
	w = httptest.NewRecorder()
	s.gated(ok(&hit)).ServeHTTP(w, httptest.NewRequest("GET", "/ws/thing", nil))
	if hit {
		t.Fatal("auth on: a /ws path reached the handler")
	}
	if w.Code != http.StatusFound {
		t.Fatalf("a non-api path wants a redirect, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/login" {
		t.Fatalf("redirect went to %q, want /login", loc)
	}
}

func TestTheGateCountsEveryRequestIncludingRefusedOnes(t *testing.T) {
	// /metrics is fed from here, and a refused request is still traffic. A
	// counter that only saw successes would hide exactly the burst a person
	// most wants to see.
	h := newHandlers(t)
	h.ConfigureAuth(true, "", t.TempDir()+"/sessions.json")
	s := &Server{handlers: h}

	hit := false
	w := httptest.NewRecorder()
	s.gated(ok(&hit)).ServeHTTP(w, httptest.NewRequest("GET", "/api/traces", nil))
	if hit || w.Code != http.StatusUnauthorized {
		t.Fatal("premise: that request should have been refused")
	}

	m := httptest.NewRecorder()
	h.Metrics(m, httptest.NewRequest("GET", "/metrics", nil))
	if !strings.Contains(m.Body.String(), "GET /api/traces") {
		t.Fatalf("the refused request was not counted:\n%s", m.Body.String())
	}
}

func TestTheGlassListensOnThisComputerOnly(t *testing.T) {
	// "This PC only" (2026-09-21). It listened on ":" + port -- every address
	// the machine has -- with no gate switched on. Loopback is this computer.
	if got := (&Server{port: "8091"}).addr(); got != "127.0.0.1:8091" {
		t.Fatalf("the glass listens on %q; this PC only is 127.0.0.1", got)
	}
}

func TestTheLockIsSwitchedOnWhereTheGlassStarts(t *testing.T) {
	// The gate stood unused for weeks because the ONE line that turns it on
	// was never written, and every stroke above still passed. So the wiring is
	// pinned where it lives: main.go must close the gate and name the lock.
	src, err := os.ReadFile("../main.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`h.ConfigureAuth(true,`, `h.ConfigureLock(`} {
		if !strings.Contains(string(src), want) {
			t.Fatalf("main.go no longer does %s -- the glass would start unlocked", want)
		}
	}

	// AND IT CARRIES A SERVICE WIRE (2026-09-24). The line above passed for
	// weeks while the second argument was the empty string, so the glass sent
	// no Authorization header and the door could not be armed without blanking
	// every page. "Is the call there" and "does the call carry anything" are
	// different questions and only the first was ever asked.
	if strings.Contains(string(src), `h.ConfigureAuth(true, "",`) {
		t.Fatal(`main.go hard-codes an empty service wire -- the glass sends no ` +
			`Authorization header, and arming the door with --auth answers ` +
			`every page 401`)
	}
	if !strings.Contains(string(src), `os.Getenv("ATLAS_SERVICE")`) {
		t.Fatal(`main.go does not read ATLAS_SERVICE -- the door reads that ` +
			`same name (cmd/atlas-mcp/main.go), and the two must agree by ` +
			`reading one place rather than by someone setting two`)
	}
	// RULE 7: a key never rides a command line where `ps` can read it.
	//
	// ASKED OF THE IMPORT BLOCK, NOT OF THE WHOLE FILE. A first cut grepped
	// the source for "auth-service" and went red on the COMMENT above the
	// change, which names the door's flag to say the two read one variable.
	// Prose quotes the thing it is explaining -- the same fault the `contains`
	// eval mode was narrowed for on 2026-09-12. A file that does not import
	// `flag` cannot take a flag, and an import block is structure rather than
	// wording.
	head := string(src)
	if i := strings.Index(head, "import ("); i >= 0 {
		head = head[i : i+strings.Index(head[i:], "\n)")]
	}
	if strings.Contains(head, `"flag"`) {
		t.Fatal("main.go imports flag -- the service wire must not ride a " +
			"command line where ps can read it (RULE 7)")
	}
}

// THE HEADER IS SENT WHEN THERE IS ONE TO SEND, AND NOT OTHERWISE.
//
// Six call sites reach the door and each guards on `h.service != ""`. That
// guard is the whole mechanism: with a wire the glass is the operator's own
// panel and the door's holds let it through; without one it behaves exactly
// as it did before `--auth` existed. A regression either way is silent --
// no header means 401 on an armed door, and a header built from an empty
// string would be a bearer of "" offered to `subtleEqual`.
func TestTheGlassSendsItsServiceWireAndOnlyWhenItHasOne(t *testing.T) {
	for _, f := range []string{"../handlers/handlers.go", "../handlers/chat.go",
		"../handlers/council.go", "../handlers/ws.go"} {
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		body := string(src)
		sends := strings.Count(body, `"Bearer "+h.service`)
		guards := strings.Count(body, `if h.service != ""`)
		if sends != guards {
			t.Fatalf("%s sends the bearer %d time(s) behind %d guard(s) -- every "+
				"send must sit behind `if h.service != \"\"`, or an empty wire "+
				"becomes a bearer of \"\"", f, sends, guards)
		}
	}
}
