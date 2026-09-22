// THE PAGE ON THE GLASS (2026-09-21, the maker's piece 2): every property
// projects.go claims. Hermetic: a temp store and a fake door on loopback.
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// pageDoor stands up a door whose `projects` tool answers with page, or
// refuses with refusal, and records what it was asked.
func pageDoor(t *testing.T, h *Handlers, page, refusal string) (*int, *map[string]any) {
	t.Helper()
	calls := 0
	var got map[string]any
	door := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var req struct {
			Params struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments"`
			} `json:"params"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		got = req.Params.Arguments
		got["_tool"] = req.Params.Name
		if refusal != "" {
			raw, _ := json.Marshal(refusal)
			w.Write([]byte(`{"result":{"isError":true,"content":[{"text":` + string(raw) + `}]}}`))
			return
		}
		inner, _ := json.Marshal(map[string]any{"name": "snake-game", "version": 2, "text": page})
		outer, _ := json.Marshal(string(inner))
		w.Write([]byte(`{"result":{"content":[{"text":` + string(outer) + `}]}}`))
	}))
	t.Cleanup(door.Close)
	h.store.SetSetting("mcp_url", door.URL)
	return &calls, &got
}

func pageRequest(name, query string) *http.Request {
	r := httptest.NewRequest("GET", "/api/projects/"+url.PathEscape(name)+"/page"+query, nil)
	r.SetPathValue("name", name)
	return r
}

const aGame = "<!DOCTYPE html>\n<html><body><script>alert('game over')</script></body></html>\n"

// THE PAGE IS SERVED SANDBOXED: its scripts run, in an origin that is not the
// glass's, and nothing on it reaches out.
func TestAProjectPageIsServedSandboxed(t *testing.T) {
	h := newTestHandlers(t)
	calls, got := pageDoor(t, h, aGame, "")
	w := httptest.NewRecorder()
	h.ProjectPage(w, pageRequest("snake-game", "?v=2"))

	if w.Code != http.StatusOK || w.Body.String() != aGame {
		t.Fatalf("the page was not served whole: %d %q", w.Code, w.Body.String())
	}
	csp := w.Header().Get("Content-Security-Policy")
	for _, want := range []string{"sandbox allow-scripts", "connect-src 'none'",
		"form-action 'none'", "default-src 'none'"} {
		if !strings.Contains(csp, want) {
			t.Fatalf("the page's header does not carry %q: %q", want, csp)
		}
	}
	// NEVER THE GLASS'S OWN ORIGIN: with it, the page's scripts would reach the
	// glass's API with his session behind them.
	if strings.Contains(csp, "allow-same-origin") {
		t.Fatalf("a model's page was served into the glass's own origin: %q", csp)
	}
	if !strings.HasPrefix(w.Header().Get("Content-Type"), "text/html") ||
		w.Header().Get("X-Content-Type-Options") != "nosniff" ||
		w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("the page's headers are not what projects.go promises: %v", w.Header())
	}
	if *calls != 1 || (*got)["_tool"] != "projects" || (*got)["action"] != "page" ||
		(*got)["name"] != "snake-game" || (*got)["version"] != float64(2) {
		t.Fatalf("the door was not asked for that page: %d %v", *calls, *got)
	}
}

// A NAME THAT IS NOT ONE, OR A VERSION THAT IS NOT A NUMBER, NEVER COSTS A CALL.
func TestABadProjectNameOrVersionNeverReachesTheDoor(t *testing.T) {
	h := newTestHandlers(t)
	calls, _ := pageDoor(t, h, aGame, "")
	for _, name := range []string{"..", "Snake-Game", "", "a/b", "-x", "snake game"} {
		w := httptest.NewRecorder()
		h.ProjectPage(w, pageRequest(name, ""))
		if w.Code != http.StatusNotFound {
			t.Fatalf("%q: want 404, got %d", name, w.Code)
		}
	}
	for _, v := range []string{"?v=x", "?v=0", "?v=-1", "?v=1.5"} {
		w := httptest.NewRecorder()
		h.ProjectPage(w, pageRequest("snake-game", v))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("%s: want 400, got %d", v, w.Code)
		}
	}
	if *calls != 0 {
		t.Fatalf("a refused name or version still asked the door %d time(s)", *calls)
	}
}

// THE DOOR'S REFUSAL IS PASSED ON IN ITS OWN WORDS, AND NOTHING IS SANDBOXED
// BECAUSE NOTHING IS SERVED.
func TestADoorRefusalIsSaidAndServesNoPage(t *testing.T) {
	h := newTestHandlers(t)
	pageDoor(t, h, "", `there is no version 9 -- snake-game has versions 1 to 2`)
	w := httptest.NewRecorder()
	h.ProjectPage(w, pageRequest("snake-game", "?v=9"))
	if w.Code != http.StatusNotFound || !strings.Contains(w.Body.String(), "versions 1 to 2") {
		t.Fatalf("the door's refusal was not passed on: %d %q", w.Code, w.Body.String())
	}
	if strings.HasPrefix(w.Header().Get("Content-Type"), "text/html") {
		t.Fatal("a refusal was served as a page")
	}
}

// THE LIST IS ONE OF THE DASHBOARD'S OWN READS: answered, and not kept as a
// trace (the ruling of 2026-09-16, "D1 b").
func TestTheProjectListIsABackgroundRead(t *testing.T) {
	h := newTestHandlers(t)
	pageDoor(t, h, aGame, "")
	before := len(h.store.List("", 100))
	w := httptest.NewRecorder()
	h.CallTool(w, httptest.NewRequest("POST", "/api/tools/call",
		strings.NewReader(`{"tool":"projects","args":{"action":"list"},"background":true}`)))
	if w.Code != http.StatusOK {
		t.Fatalf("the list call failed: %d %s", w.Code, w.Body.String())
	}
	if after := len(h.store.List("", 100)); after != before {
		t.Fatalf("the Dashboard's project list was kept as a trace (%d -> %d)", before, after)
	}
}
