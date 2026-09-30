// Chat streaming face (N1): GET /chat/stream sends Server-Sent Events.
// Tokens stream as `event: token`; the witnessed whole answer closes as
// `event: done` with its receipt; refusals arrive as `event: refused`.
// The handler owns no logic — guard, route, ask, witness and append all
// live in the chat package behind the one-writer lock, shared with the
// tools/call face. Strangers are refused by name before anything streams.
package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"atlas/line/internal/tools"
)

func (s *Server) handleChatStream(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	args := map[string]any{
		"session":  q.Get("session"),
		"question": q.Get("question"),
		"actor":    q.Get("actor"),
		"voice":    q.Get("voice"),
		"project":  q.Get("project"),
	}
	if s.auth.On {
		cred := s.credentialFor(bearer(r.Header.Get("Authorization")))
		if err := s.gateCall("chat_send", args, cred); err != nil {
			s.m.err()
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
	s.m.inc("chat/stream")
	tn, err := s.tenants.Resolve(q.Get("project"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	emit := func(ev string, data any) {
		b, _ := json.Marshal(data)
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev, string(b))
		flusher.Flush()
	}
	// THE SOCKET HAS ONE WRITER (2026-09-29). The token callback runs on the
	// send's goroutine and used to write each token to the socket itself: two
	// goroutines on one ResponseWriter, and -- once the browser had gone and
	// this handler had returned -- writes to a ResponseWriter net/http had
	// already taken back. Tokens are handed to THIS goroutine, which owns the
	// socket, the way /run/stream hands over the engine's events; a token that
	// would wait on a browser that has left is dropped, because there is no
	// one to show it to. The send itself keeps its own counsel: whole or
	// nothing, witnessed or not, whether anyone is watching.
	gone := r.Context().Done()
	tokens := make(chan string, 256)
	done := make(chan struct{})
	var out string
	var callErr error
	go func() {
		defer close(done)
		defer close(tokens)
		out, callErr = tools.ChatSendStream(s.tools, tn, args, func(tok string) {
			select {
			case tokens <- tok:
			case <-gone:
			}
		})
	}()
	for {
		select {
		case tok, ok := <-tokens:
			if !ok {
				tokens = nil
				continue
			}
			emit("token", map[string]any{"token": tok})
		case <-done:
			// Guarded: tokens is nil once drained and closed, and ranging a
			// nil channel blocks forever (runstream.go learned it first).
			if tokens != nil {
				for tok := range tokens {
					emit("token", map[string]any{"token": tok})
				}
			}
			if callErr != nil {
				emit("refused", map[string]any{"error": callErr.Error()})
				return
			}
			emit("done", map[string]any{"text": out})
			return
		case <-gone:
			// The browser went away mid-turn: the send keeps its own
			// counsel (whole or nothing, witnessed or not).
			return
		}
	}
}
