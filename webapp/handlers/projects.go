// THE PAGE ON THE GLASS (2026-09-21, the maker's piece 2; the operator: "go on
// piece 2"). A page the maker made is shown on the Dashboard, beside the run --
// served here, into a frame, from the door's `projects` tool.
//
// THE PAGE IS A MODEL'S CODE, AND THIS IS THE GLASS HE IS SIGNED INTO. So it is
// served under one header that takes the glass's own powers away from it:
//
//	sandbox allow-scripts allow-modals
//
// The page's scripts run -- a game is a script -- in an ORIGIN OF THEIR OWN
// that is not this glass's. Never `allow-same-origin`: with it the page would
// share the glass's origin, and its scripts could reach the glass's API with
// his session behind them. And beside the sandbox, nothing reaches out:
// no fetch, no socket, no form, no script or style from anywhere but the page
// itself (RULE 4 -- the maker already refuses a page that reaches for the
// internet; this is the second wall, and it holds whatever the page says).
//
// THIS HEADER IS THE WALL, AND THE ONLY ONE. The frame carries no `sandbox`
// attribute: the Claude app's browser pane refuses any frame that does, and the
// header alone was measured there doing all of it (projects.js says what was
// measured). "Open it in its own tab" gets the same header, so the page runs
// sandboxed wherever it is opened from here.
//
// KNOWN AND NOT CHASED: a page in its own origin has no storage of its own --
// a game that saves a high score to localStorage forgets it here. Opening the
// file from projects\<name>\ in a browser is unchanged, and keeps it.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// The shape of a name the maker gives a project, the same rule the door holds
// (projects.go). Checked here too, so a bad name never costs a call to the door.
var projectNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// pageSandbox is the header a project's page is served under.
const pageSandbox = "sandbox allow-scripts allow-modals; default-src 'none'; " +
	"script-src 'unsafe-inline'; style-src 'unsafe-inline'; img-src data: blob:; " +
	"media-src data: blob:; font-src data:; connect-src 'none'; " +
	"form-action 'none'; base-uri 'none'; frame-ancestors 'self'"

// ProjectPage serves one project's page: GET /api/projects/{name}/page, as it
// stands, or ?v=N for version N.
func (h *Handlers) ProjectPage(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if !projectNameRe.MatchString(name) {
		http.Error(w, "there is no project by that name", http.StatusNotFound)
		return
	}
	args := map[string]any{"action": "page", "name": name}
	if v := strings.TrimSpace(r.URL.Query().Get("v")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			http.Error(w, "a version is a whole number, from 1", http.StatusBadRequest)
			return
		}
		args["version"] = n
	}
	text, err := h.rpcCallAs(r, "projects", args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var page struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(text), &page); err != nil {
		http.Error(w, "the door answered in a shape this cannot read", http.StatusBadGateway)
		return
	}
	hd := w.Header()
	hd.Set("Content-Security-Policy", pageSandbox)
	hd.Set("Content-Type", "text/html; charset=utf-8")
	hd.Set("X-Content-Type-Options", "nosniff")
	hd.Set("Cache-Control", "no-store")
	hd.Set("Referrer-Policy", "no-referrer")
	_, _ = io.WriteString(w, page.Text)
}
