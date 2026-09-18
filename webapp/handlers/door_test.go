// How long the glass waits on the door.
//
// Every call the glass made to the door went through http.DefaultClient, which
// never times out, so a door that stopped answering held the glass for as long
// as it stayed silent. These strokes stand up a door that never answers -- a
// real loopback server that accepts every request and says nothing -- and hold
// each kind of call to its bound.
//
// Hermetic: a temp store and a loopback server; nothing of the estate's.
package handlers

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// silentDoor points h at a door that accepts every request and answers none,
// until the stroke ends or the caller gives up.
func silentDoor(t *testing.T, h *Handlers) {
	t.Helper()
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	t.Cleanup(srv.Close)
	t.Cleanup(func() { close(release) })
	h.store.SetSetting("mcp_url", srv.URL)
}

// shorten sets a wait for the length of one stroke.
func shorten(t *testing.T, wait *time.Duration, to time.Duration) {
	t.Helper()
	was := *wait
	*wait = to
	t.Cleanup(func() { *wait = was })
}

// within runs serve and fails the stroke if it has not returned by limit.
func within(t *testing.T, limit time.Duration, what string, serve func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		serve()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(limit):
		t.Fatalf("%s was still waiting on a silent door after %s", what, limit)
	}
}

// THE QUICK READS GIVE UP, and say the door did not answer, rather than hold
// the page. /run/state is asked by every page, and it waited as long as the
// door did.
func TestAQuickReadGivesUpOnASilentDoor(t *testing.T) {
	for _, c := range []struct {
		name string
		path string
		call func(h *Handlers) func(http.ResponseWriter, *http.Request)
	}{
		{"run_state", "/api/council/state", func(h *Handlers) func(http.ResponseWriter, *http.Request) { return h.CouncilState }},
		{"tools", "/api/tools", func(h *Handlers) func(http.ResponseWriter, *http.Request) { return h.ListTools }},
	} {
		t.Run(c.name, func(t *testing.T) {
			h := newTestHandlers(t)
			silentDoor(t, h)
			shorten(t, &pollWait, 200*time.Millisecond)
			w := httptest.NewRecorder()
			within(t, 5*time.Second, c.path, func() { c.call(h)(w, httptest.NewRequest("GET", c.path, nil)) })
			if w.Code != http.StatusBadGateway || !strings.Contains(w.Body.String(), "deadline") {
				t.Fatalf("%s answered %d %q; a silent door is a 502 that says so", c.path, w.Code, w.Body.String())
			}
		})
	}
}

// A TOOL CALL IS BOUNDED TOO, and its trace records that the door did not
// answer rather than recording nothing.
func TestAToolCallGivesUpOnASilentDoor(t *testing.T) {
	h := newTestHandlers(t)
	silentDoor(t, h)
	shorten(t, &callWait, 200*time.Millisecond)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/tools/call", strings.NewReader(`{"tool":"muster","args":{}}`))
	within(t, 5*time.Second, "a tool call", func() { h.CallTool(w, req) })
	if !strings.Contains(w.Body.String(), "deadline") {
		t.Fatalf("the tool call answered %q; it must say the door did not answer", w.Body.String())
	}
	if tr := h.store.List("", 1); len(tr) != 1 || !strings.Contains(tr[0].Output, "deadline") {
		t.Fatalf("the trace does not record that the door did not answer: %+v", tr)
	}
}

// AND A PROXIED CALL -- a flow, a prompt, a chat -- through rpcCall.
func TestAProxiedCallGivesUpOnASilentDoor(t *testing.T) {
	h := newTestHandlers(t)
	silentDoor(t, h)
	shorten(t, &callWait, 200*time.Millisecond)
	w := httptest.NewRecorder()
	within(t, 5*time.Second, "a proxied call", func() {
		h.FlowList(w, httptest.NewRequest("GET", "/api/flows", nil))
	})
	if w.Code != http.StatusBadGateway || !strings.Contains(w.Body.String(), "deadline") {
		t.Fatalf("flow_list answered %d %q; a silent door is a 502 that says so", w.Code, w.Body.String())
	}
}

// A STREAM ENDS WHEN ITS BROWSER LEAVES, even before the door has said a word.
// A stream's own bound is callWait, thirty minutes, and it is left at that
// here: the browser leaving is the only thing that can end it sooner.
func TestAStreamEndsWhenItsBrowserLeavesBeforeTheDoorSpeaks(t *testing.T) {
	for _, c := range []struct {
		name string
		path string
		call func(h *Handlers) func(http.ResponseWriter, *http.Request)
	}{
		{"council", "/api/council/stream?objective=x", func(h *Handlers) func(http.ResponseWriter, *http.Request) { return h.StreamCouncil }},
		{"chat", "/api/chat/stream?session=s&question=q", func(h *Handlers) func(http.ResponseWriter, *http.Request) { return h.StreamChat }},
	} {
		t.Run(c.name, func(t *testing.T) {
			h := newTestHandlers(t)
			silentDoor(t, h)
			ctx, leave := context.WithCancel(context.Background())
			defer leave()
			req := httptest.NewRequest("GET", c.path, nil).WithContext(ctx)
			w := httptest.NewRecorder()
			time.AfterFunc(200*time.Millisecond, leave)
			within(t, 5*time.Second, "the "+c.name+" stream", func() { c.call(h)(w, req) })
		})
	}
}

// A TURN ASKED OVER THE WEBSOCKET is bounded by callWait alone, and tells its
// socket the door did not answer.
func TestAWebSocketTurnGivesUpOnASilentDoor(t *testing.T) {
	h := newTestHandlers(t)
	silentDoor(t, h)
	shorten(t, &callWait, 200*time.Millisecond)
	glass, browser := net.Pipe()
	var frames bytes.Buffer
	read := make(chan struct{})
	go func() {
		_, _ = io.Copy(&frames, browser)
		close(read)
	}()
	within(t, 5*time.Second, "a WebSocket turn", func() {
		h.wsChatSend(&wsConn{c: glass}, "", "s", "q", "", "")
	})
	glass.Close()
	<-read
	browser.Close()
	if !strings.Contains(frames.String(), "deadline") {
		t.Fatalf("the socket was told %q; it must hear that the door did not answer", frames.String())
	}
}
