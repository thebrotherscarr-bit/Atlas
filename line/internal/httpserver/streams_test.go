// The two streams that had kept what /run/stream lost on 2026-09-15 (2026-09-29).
package httpserver

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"atlas/line/internal/chat"
	"atlas/line/internal/engine"
	"atlas/line/internal/tools"
)

// A CAPTURE OUTLIVES THE BROWSER THAT STARTED IT. /run/listen's sink was the
// bare send /run/stream gave up on 2026-09-15: sixty-four events of room, and
// nothing reading them once the tab had closed, so the sixty-fifth stopped the
// engine's reader mid-capture and the world's run lock was held for good. A
// turn asked for after the browser has gone must get its turn.
func TestACaptureOutlivesTheBrowserThatStartedIt(t *testing.T) {
	s, _ := newTestServer(t)
	tn, err := s.tenants.Resolve("t")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv(stubEngineEnv, "1")
	e, err := tools.Engines().Open("t", tn.Home,
		`"`+os.Args[0]+`" -test.run=^TestStubEngineProcess$ --`)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = tools.Engines().CloseOne(tn.Home) })

	srv := httptest.NewServer(s.Handler())
	t.Cleanup(srv.Close)
	ctx, leave := context.WithCancel(context.Background())
	defer leave()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		srv.URL+"/run/listen?project=t&seconds=50000", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	// Read until the capture is under way -- the first of the engine's own
	// frames -- so the run lock is held when the browser leaves.
	rd := bufio.NewReader(resp.Body)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			t.Fatalf("the capture never began: %v", err)
		}
		if strings.Contains(line, "event: engine") {
			break
		}
	}
	leave() // the browser goes
	resp.Body.Close()

	done := make(chan error, 1)
	go func() {
		_, err := e.Run("tokens 1", "", "", engine.Head{}, nil)
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("the turn after the capture did not run: %v", err)
		}
	case <-time.After(20 * time.Second):
		t.Fatal("the browser left and the capture never ended: the engine is stopped " +
			"behind a queue nobody reads, and the world's run lock with it")
	}
}

// watched is a socket that counts what is written to it and WHEN: a frame
// written after the handler has returned is a write to a ResponseWriter
// net/http has taken back, which is the fault this file is about.
type watched struct {
	mu       sync.Mutex
	header   http.Header
	frames   int
	late     atomic.Int64
	returned atomic.Bool
	first    chan struct{}
	once     sync.Once
}

func (w *watched) Header() http.Header { return w.header }
func (w *watched) WriteHeader(int)     {}
func (w *watched) Flush()              {}
func (w *watched) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.returned.Load() {
		w.late.Add(1)
	}
	w.frames++
	if bytes.Contains(b, []byte("event: token")) {
		w.once.Do(func() { close(w.first) })
	}
	return len(b), nil
}

// fakeRack is a loopback door that speaks enough of Ollama for one chat turn:
// one voice on the ladder, which speaks, and an answer streamed a token at a
// time, slowly enough for a browser to leave in the middle of it.
func fakeRack(t *testing.T, tokens int, gap time.Duration) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/tags":
			fmt.Fprint(w, `{"models":[{"name":"stub:latest","size":1000,"details":{"family":"stub"}}]}`)
		case "/api/show":
			fmt.Fprint(w, `{"capabilities":["completion"]}`)
		case "/api/generate":
			fl, _ := w.(http.Flusher)
			for i := 0; i < tokens; i++ {
				fmt.Fprint(w, `{"response":"x","done":false}`+"\n")
				if fl != nil {
					fl.Flush()
				}
				time.Sleep(gap)
			}
			fmt.Fprint(w, `{"response":"","done":true}`+"\n")
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	t.Setenv("OLLAMA_HOST", srv.URL)
}

// THE CHAT STREAM HAS ONE WRITER, AND NONE ONCE THE BROWSER HAS GONE. The
// token callback runs on the send's goroutine and used to write each token to
// the socket itself; when the browser left, this handler returned and the send
// went on writing tokens to a ResponseWriter that was no longer its to write
// to. The send still finishes -- whole or nothing, witnessed -- with nobody
// watching; it just says nothing to a socket that has been taken back.
func TestTheChatStreamHasOneWriterAndNoneOnceTheBrowserHasGone(t *testing.T) {
	s, _ := newTestServer(t)
	tn, err := s.tenants.Resolve("t")
	if err != nil {
		t.Fatal(err)
	}
	fakeRack(t, 300, 2*time.Millisecond)
	session, err := chat.Start(tn.Home, "t", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx, leave := context.WithCancel(context.Background())
	defer leave()
	req := httptest.NewRequest(http.MethodGet,
		"/chat/stream?project=t&session="+session+"&question="+url.QueryEscape("hello there"),
		nil).WithContext(ctx)
	w := &watched{header: http.Header{}, first: make(chan struct{})}
	back := make(chan struct{})
	go func() {
		s.Handler().ServeHTTP(w, req)
		w.returned.Store(true)
		close(back)
	}()
	select {
	case <-w.first:
	case <-time.After(10 * time.Second):
		t.Fatal("no token ever reached the socket")
	}
	leave() // the browser goes, mid-answer
	select {
	case <-back:
	case <-time.After(10 * time.Second):
		t.Fatal("the handler did not return when the browser left")
	}
	// The send keeps its own counsel: the answer is finished and witnessed
	// with nobody watching, which is exactly when the old handler went on
	// writing.
	deadline := time.Now().Add(20 * time.Second)
	for {
		turns, _ := chat.List(tn.Home, session, 0)
		if len(turns) == 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the send did not finish after the browser left; a closed tab must not lose the turn")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if n := w.late.Load(); n != 0 {
		t.Fatalf("%d frame(s) were written to the socket after the handler had returned it", n)
	}
}
