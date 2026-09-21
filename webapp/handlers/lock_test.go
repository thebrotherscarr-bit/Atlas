// THE LOCK'S STROKES (2026-09-21). Every property lock.go claims, stated so a
// change to one is a decision somebody made rather than a hole nobody saw.
//
// Hermetic: temp files, a fake clock and a fake door. Nothing here touches the
// real glass, its data folder or the door on :8090.
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// here is this computer, as the connection reports it.
const here = "127.0.0.1:50000"

func lockedHandlers(t *testing.T) (*Handlers, string) {
	t.Helper()
	h := newTestHandlers(t)
	dir := t.TempDir()
	h.ConfigureAuth(true, "", filepath.Join(dir, "sessions.json"))
	user := filepath.Join(dir, "user.json")
	h.ConfigureLock(user)
	return h, user
}

func call(fn http.HandlerFunc, method, path, body, remote string, c *http.Cookie) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if remote != "" {
		r.RemoteAddr = remote
	}
	if c != nil {
		r.AddCookie(c)
	}
	w := httptest.NewRecorder()
	fn(w, r)
	return w
}

func sessionCookie(t *testing.T, w *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, c := range w.Result().Cookies() {
		if c.Name == "atlas_session" && c.Value != "" {
			return c
		}
	}
	t.Fatalf("no session cookie was set: %v", w.Header())
	return nil
}

func body(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var d map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &d); err != nil {
		t.Fatalf("not JSON: %q", w.Body.String())
	}
	return d
}

func TestAPINIsFourToEightDigitsAndNothingElse(t *testing.T) {
	for _, ok := range []string{"1234", "0000", "12345678"} {
		if err := validPIN(ok); err != nil {
			t.Fatalf("%q is a PIN and was refused: %v", ok, err)
		}
	}
	// Arabic-Indic digits are digits to unicode and NOT to this rule: a PIN is
	// what a keypad types, and the check is by rune, not by byte length.
	for _, bad := range []string{"", "123", "123456789", "12a4", "12 34", "١٢٣٤"} {
		if validPIN(bad) == nil {
			t.Fatalf("%q is not a PIN and was taken", bad)
		}
	}
}

func TestANameIsCleanedAndBounded(t *testing.T) {
	if n, err := cleanName("  Ada   Lovelace "); err != nil || n != "Ada Lovelace" {
		t.Fatalf("got %q, %v", n, err)
	}
	for _, bad := range []string{"", "   ", strings.Repeat("x", nameMax+1), "a\x00b"} {
		if _, err := cleanName(bad); err == nil {
			t.Fatalf("%q should have been refused as a name", bad)
		}
	}
}

func TestSetupHappensOnceAndOnlyFromThisComputer(t *testing.T) {
	h, user := lockedHandlers(t)

	// NOT FROM ANOTHER MACHINE -- judged from the connection, not a header.
	w := call(h.LockSetup, "POST", "/api/setup", `{"name":"Kyle","pin":"4321"}`, "192.168.1.20:5555", nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("setup from the network: want 403, got %d", w.Code)
	}
	if _, err := os.Stat(user); err == nil {
		t.Fatal("a refused setup wrote the user file")
	}

	// FROM THIS ONE: saved, cleaned, and signed straight in.
	w = call(h.LockSetup, "POST", "/api/setup", `{"name":"  Kyle  ","pin":"4321"}`, here, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("setup: want 200, got %d %s", w.Code, w.Body.String())
	}
	sessionCookie(t, w)
	raw, err := os.ReadFile(user)
	if err != nil {
		t.Fatal(err)
	}
	// THE PIN NEVER REACHES THE DISK: no field holds it, and what is there
	// checks it.
	if strings.Contains(string(raw), `"4321"`) {
		t.Fatalf("the PIN itself reached the disk: %s", raw)
	}
	var f lockFile
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	if f.Name != "Kyle" || f.Rounds != pinRounds || len(f.Hash) != 64 || len(f.Salt) != 32 {
		t.Fatalf("the user file is not what lock.go promises: %+v", f)
	}

	// NOT TWICE: a second setup is refused and changes nothing.
	w = call(h.LockSetup, "POST", "/api/setup", `{"name":"Someone","pin":"9999"}`, here, nil)
	if w.Code != http.StatusConflict {
		t.Fatalf("second setup: want 409, got %d", w.Code)
	}
	if again, _ := os.ReadFile(user); string(again) != string(raw) {
		t.Fatal("a refused second setup changed the user")
	}

	// AND A BAD PIN OR NAME IS REFUSED BEFORE ANYTHING IS WRITTEN.
	h2, user2 := lockedHandlers(t)
	for _, b := range []string{`{"name":"Kyle","pin":"12"}`, `{"name":"","pin":"1234"}`, `not json`} {
		if w := call(h2.LockSetup, "POST", "/api/setup", b, here, nil); w.Code != http.StatusBadRequest {
			t.Fatalf("%s: want 400, got %d", b, w.Code)
		}
	}
	if _, err := os.Stat(user2); err == nil {
		t.Fatal("a refused setup wrote the user file")
	}
}

func TestTheRightPINOpensAndAWrongOneIsCounted(t *testing.T) {
	h, _ := lockedHandlers(t)
	call(h.LockSetup, "POST", "/api/setup", `{"name":"Kyle","pin":"2468"}`, here, nil)

	w := call(h.Unlock, "POST", "/api/unlock", `{"pin":"1111"}`, here, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("wrong PIN: want 401, got %d", w.Code)
	}
	if d := body(t, w); d["tries_left"] != float64(pinTries-1) {
		t.Fatalf("the refusal must say how many tries are left: %v", d)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == "atlas_session" && c.Value != "" {
			t.Fatal("a wrong PIN was handed a session")
		}
	}

	w = call(h.Unlock, "POST", "/api/unlock", `{"pin":"2468"}`, here, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("right PIN: want 200, got %d %s", w.Code, w.Body.String())
	}
	c := sessionCookie(t, w)
	// THE COOKIE IS HTTP-ONLY AND SAME-SITE STRICT: no script reads it, and no
	// other site's link carries it here (session.go says why).
	if !c.HttpOnly || c.SameSite != http.SameSiteStrictMode {
		t.Fatalf("cookie flags: HttpOnly=%v SameSite=%v", c.HttpOnly, c.SameSite)
	}

	// THE SESSION IS THE OPERATOR'S: his name, and no tenant pinned.
	r := httptest.NewRequest("GET", "/api/me", nil)
	r.AddCookie(c)
	sess, ok := h.sessionOf(r)
	if !ok || sess == nil || sess.User != "Kyle" || sess.Tenant != "" {
		t.Fatalf("the PIN session is not the operator's: %+v", sess)
	}
	me := call(h.Me, "GET", "/api/me", "", here, c)
	if d := body(t, me); d["user"] != "Kyle" || d["auth"] != true {
		t.Fatalf("/api/me must say who is signed in: %v", d)
	}
}

func TestFiveWrongPINsCloseTheLockForAMinute(t *testing.T) {
	h, _ := lockedHandlers(t)
	now := time.Date(2026, 9, 21, 15, 0, 0, 0, time.UTC)
	h.lock.now = func() time.Time { return now }
	call(h.LockSetup, "POST", "/api/setup", `{"name":"Kyle","pin":"1357"}`, here, nil)

	wrong := func() *httptest.ResponseRecorder {
		return call(h.Unlock, "POST", "/api/unlock", `{"pin":"0000"}`, here, nil)
	}
	right := func() *httptest.ResponseRecorder {
		return call(h.Unlock, "POST", "/api/unlock", `{"pin":"1357"}`, here, nil)
	}
	for i := 1; i < pinTries; i++ {
		if w := wrong(); w.Code != http.StatusUnauthorized {
			t.Fatalf("wrong PIN %d: want 401, got %d", i, w.Code)
		}
	}
	w := wrong()
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("wrong PIN %d: want 429, got %d", pinTries, w.Code)
	}
	if w.Header().Get("Retry-After") != "60" || body(t, w)["wait"] != float64(60) {
		t.Fatalf("the closed lock must say how long: %v %s", w.Header(), w.Body.String())
	}

	// THE RIGHT PIN WAITS IT OUT TOO, or the minute would stop nothing.
	if w := right(); w.Code != http.StatusTooManyRequests {
		t.Fatalf("right PIN while closed: want 429, got %d", w.Code)
	}
	now = now.Add(30 * time.Second)
	if w := right(); w.Code != http.StatusTooManyRequests || body(t, w)["wait"] != float64(30) {
		t.Fatalf("half a minute in: want 429 with 30 left, got %d %s", w.Code, w.Body.String())
	}
	now = now.Add(31 * time.Second)
	if w := right(); w.Code != http.StatusOK {
		t.Fatalf("after the minute the right PIN opens: got %d", w.Code)
	}
	// and the count starts over
	if d := body(t, wrong()); d["tries_left"] != float64(pinTries-1) {
		t.Fatalf("the count must start over after the lock opens: %v", d)
	}
}

func TestTheLockStateTellsThePageWhatToDraw(t *testing.T) {
	h, user := lockedHandlers(t)
	st := body(t, call(h.LockState, "GET", "/api/lock", "", here, nil))
	if st["setup"] != true || st["name"] != "" || st["signed_in"] != false {
		t.Fatalf("before setup the page must draw the welcome: %v", st)
	}
	w := call(h.LockSetup, "POST", "/api/setup", `{"name":"Kyle","pin":"8642"}`, here, nil)
	c := sessionCookie(t, w)
	st = body(t, call(h.LockState, "GET", "/api/lock", "", here, nil))
	if st["setup"] != false || st["name"] != "Kyle" || st["signed_in"] != false {
		t.Fatalf("after setup, a stranger's page draws the lock with his name: %v", st)
	}
	if st = body(t, call(h.LockState, "GET", "/api/lock", "", here, c)); st["signed_in"] != true {
		t.Fatalf("a signed-in page is told so: %v", st)
	}

	// A DAMAGED FILE IS SAID, AND NEVER READ AS "NOBODY IS SET UP" -- that
	// would hand the lock to whoever opened the page next.
	if err := os.WriteFile(user, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	st = body(t, call(h.LockState, "GET", "/api/lock", "", here, nil))
	if st["setup"] != false || st["damaged"] == nil {
		t.Fatalf("a damaged file must be named, not treated as empty: %v", st)
	}
	if w := call(h.LockSetup, "POST", "/api/setup", `{"name":"Mallory","pin":"1234"}`, here, nil); w.Code != http.StatusConflict {
		t.Fatalf("setup over a damaged file: want 409, got %d", w.Code)
	}
	w = call(h.Unlock, "POST", "/api/unlock", `{"pin":"8642"}`, here, nil)
	if w.Code != http.StatusInternalServerError || !strings.Contains(w.Body.String(), "user.json") {
		t.Fatalf("unlock over a damaged file must say how to reset: %d %s", w.Code, w.Body.String())
	}
}

func TestThePINSessionKeepsThePagesProjectAndAKeyStaysPinned(t *testing.T) {
	var got map[string]any
	var gotQuery map[string][]string
	door := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/chat/stream" {
			gotQuery = r.URL.Query()
			w.Header().Set("Content-Type", "text/event-stream")
			w.Write([]byte("data: done\n\n"))
			return
		}
		var req struct {
			Params struct {
				Arguments map[string]any `json:"arguments"`
			} `json:"params"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		got = req.Params.Arguments
		w.Write([]byte(`{"result":{"content":[{"text":"ok"}]}}`))
	}))
	defer door.Close()

	h, _ := lockedHandlers(t)
	h.store.SetSetting("mcp_url", door.URL)
	c := sessionCookie(t, call(h.LockSetup, "POST", "/api/setup", `{"name":"Kyle","pin":"1122"}`, here, nil))
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(c)

	// THE OPERATOR NAMES ANY WORLD, as the glass did with the gate off.
	if _, err := h.rpcCallAs(r, "flow_list", map[string]any{"project": "atlas"}); err != nil {
		t.Fatal(err)
	}
	if got["project"] != "atlas" {
		t.Fatalf("the PIN session overwrote the project the page asked for: %v", got)
	}
	qq := map[string][]string{}
	if err := h.scopeProject(r, qq, "atlas"); err != nil || len(qq["project"]) != 1 || qq["project"][0] != "atlas" {
		t.Fatalf("scopeProject must keep the asked project for the operator: %v %v", qq, err)
	}
	// CallTool is the road every Dashboard tool call takes, so it is held too.
	if w := call(h.CallTool, "POST", "/api/tools/call", `{"tool":"flow_list","args":{"project":"atlas"}}`, here, c); w.Code != http.StatusOK || got["project"] != "atlas" {
		t.Fatalf("CallTool overwrote the project the page asked for: %d %v", w.Code, got)
	}
	// and a chat stream names no project of its own, as with the gate off
	call(h.StreamChat, "GET", "/api/chat/stream?question=hi", "", here, c)
	if gotQuery == nil || gotQuery["project"] != nil {
		t.Fatalf("the operator's chat stream named a project: %v", gotQuery)
	}

	// AND A KEY'S SESSION IS STILL PINNED TO ITS TENANT -- the wall did not move.
	key := h.sessions.create("research", []string{"research"}, "k1")
	r2 := httptest.NewRequest("GET", "/", nil)
	r2.AddCookie(&http.Cookie{Name: "atlas_session", Value: key.Token})
	if _, err := h.rpcCallAs(r2, "flow_list", map[string]any{"project": "atlas"}); err != nil {
		t.Fatal(err)
	}
	if got["project"] != "research" {
		t.Fatalf("a key session named another tenant's world: %v", got)
	}
	qq = map[string][]string{}
	if err := h.scopeProject(r2, qq, "atlas"); err != nil || qq["project"][0] != "research" {
		t.Fatalf("scopeProject let a key name another world: %v", qq)
	}
	kc := &http.Cookie{Name: "atlas_session", Value: key.Token}
	call(h.CallTool, "POST", "/api/tools/call", `{"tool":"flow_list","args":{"project":"atlas"}}`, here, kc)
	if got["project"] != "research" {
		t.Fatalf("CallTool let a key name another world: %v", got)
	}
	call(h.StreamChat, "GET", "/api/chat/stream?question=hi", "", here, kc)
	if p := gotQuery["project"]; len(p) != 1 || p[0] != "research" {
		t.Fatalf("a key's chat stream was not pinned to its tenant: %v", gotQuery)
	}

	// WITH NO SESSION AT ALL, nothing goes to the door.
	got = nil
	if _, err := h.rpcCallAs(httptest.NewRequest("GET", "/", nil), "flow_list", map[string]any{}); err == nil || got != nil {
		t.Fatalf("no session reached the door: %v %v", err, got)
	}
}

func TestTheLockFacesSayWhenNothingIsConfigured(t *testing.T) {
	h := newTestHandlers(t)
	for _, fn := range []http.HandlerFunc{h.LockState, h.LockSetup, h.Unlock} {
		if w := call(fn, "POST", "/api/lock", `{}`, here, nil); w.Code != http.StatusServiceUnavailable {
			t.Fatalf("an unconfigured lock must say so: %d", w.Code)
		}
	}
}
