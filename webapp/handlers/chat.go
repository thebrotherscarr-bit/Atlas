// Chat faces (N1): whole-answer send over MCP /rpc, token streaming by
// proxying the MCP door's /chat/stream SSE 1:1. The webapp owns no chat
// logic — guard, route, ask, witness and append live in the LINE's chat
// package; every turn below receipts through it.
package handlers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// THE GLASS WAITS ON THE DOOR FOR A STATED TIME, AND NO LONGER (2026-09-15).
// Every call to the door went through http.DefaultClient, which has no timeout
// at all, so a door that stopped answering held each caller -- a goroutine and a
// connection -- for as long as it stayed silent. ESTATE LAW 7: bounded
// everything, timeouts on calls. Variables, not constants, only so a stroke can
// shorten them.
var (
	// pollWait bounds the two quick reads: /run/state, which every page polls,
	// and /tools. Both answer in milliseconds, and the Dashboard asks for
	// /run/state again every fifteen seconds.
	pollWait = 15 * time.Second
	// callWait bounds everything else: a tool call, a proxied call, a turn's
	// stream. It is the longest fixed bound the door sets on any tool of its own
	// -- rack_pull's download, thirty minutes -- so nothing the door itself
	// promises to finish is cut short here. A flow's budget_s and a prompt_eval's
	// cases are bounded by what the operator wrote, not by the door; past this
	// the glass stops waiting, and the door still finishes them and writes their
	// runs.
	callWait = 30 * time.Minute
)

func (h *Handlers) mcpURL() string {
	u := h.store.GetSetting("mcp_url")
	if u == "" {
		u = "http://127.0.0.1:8090"
	}
	return strings.TrimRight(u, "/")
}

func (h *Handlers) rpcCall(tool string, args map[string]any) (string, error) {
	payload, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0", "id": 1, "method": "tools/call",
		"params": map[string]any{"name": tool, "arguments": args},
	})
	// Not tied to the asker: callers announce the answer to every open window,
	// and it is still announced after the one that asked has gone.
	ctx, cancel := context.WithTimeout(context.Background(), callWait)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", h.mcpURL()+"/rpc", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if h.service != "" {
		req.Header.Set("Authorization", "Bearer "+h.service)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return "", err
	}
	var doc struct {
		Result struct {
			IsError bool `json:"isError"`
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return "", fmt.Errorf("mcp spoke unparsably")
	}
	text := ""
	if len(doc.Result.Content) > 0 {
		text = doc.Result.Content[0].Text
	}
	if doc.Result.IsError {
		return "", fmt.Errorf("%s", text)
	}
	return text, nil
}

// ListSessions proxies chat_sessions.
func (h *Handlers) ListSessions(w http.ResponseWriter, r *http.Request) {
	text, err := h.rpcCallAs(r, "chat_sessions", map[string]any{})
	if err != nil {
		jsonErr(w, 502, err.Error())
		return
	}
	jsonResp(w, map[string]string{"sessions": text})
}

// GetSession proxies chat_list for one session.
func (h *Handlers) GetSession(w http.ResponseWriter, r *http.Request) {
	session := r.URL.Query().Get("session")
	text, err := h.rpcCallAs(r, "chat_list", map[string]any{"session": session})
	if err != nil {
		jsonErr(w, 502, err.Error())
		return
	}
	jsonResp(w, map[string]string{"turns": text})
}

// StartSession proxies chat_start and announces the session on both buses.
func (h *Handlers) StartSession(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Actor string `json:"actor"`
		Voice string `json:"voice"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	text, err := h.rpcCallAs(r, "chat_start", map[string]any{"actor": body.Actor, "voice": body.Voice})
	if err != nil {
		jsonErr(w, 502, err.Error())
		return
	}
	h.broadcast("chat_opened", map[string]string{"text": text})
	jsonResp(w, map[string]string{"status": "ok", "text": text})
}

// SendChat proxies chat_send (whole answer) and announces completion.
func (h *Handlers) SendChat(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session  string `json:"session"`
		Question string `json:"question"`
		Actor    string `json:"actor"`
		Voice    string `json:"voice"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr(w, 400, "invalid json")
		return
	}
	if strings.TrimSpace(body.Session) == "" || strings.TrimSpace(body.Question) == "" {
		jsonErr(w, 400, "session and question are required")
		return
	}
	text, err := h.rpcCallAs(r, "chat_send", map[string]any{
		"session": body.Session, "question": body.Question,
		"actor": body.Actor, "voice": body.Voice,
	})
	if err != nil {
		jsonErr(w, 502, err.Error())
		return
	}
	h.broadcast("chat_done", map[string]string{"session": body.Session})
	jsonResp(w, map[string]string{"status": "ok", "text": text})
}

// StreamChat proxies the MCP door's /chat/stream SSE byte-for-byte: the
// browser's EventSource sees tokens as the voice speaks them, then done.
func (h *Handlers) StreamChat(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	// Re-encode query safely through the URL parser.
	req, err := http.NewRequest("GET", h.mcpURL()+"/chat/stream", nil)
	if err != nil {
		jsonErr(w, 502, err.Error())
		return
	}
	qq := req.URL.Query()
	qq.Set("session", q.Get("session"))
	qq.Set("question", q.Get("question"))
	qq.Set("voice", q.Get("voice"))
	qq.Set("actor", q.Get("actor"))
	if h.authOn {
		sess, ok := h.sessionOf(r)
		if !ok || sess == nil {
			jsonErr(w, 401, "login required")
			return
		}
		qq.Set("project", sess.Tenant)
	}
	req.URL.RawQuery = qq.Encode()
	if h.service != "" {
		req.Header.Set("Authorization", "Bearer "+h.service)
	}
	// A stream has one reader, the browser that opened it, so it also ends when
	// that browser leaves -- including while the door has yet to say a word.
	ctx, cancel := context.WithTimeout(r.Context(), callWait)
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		jsonErr(w, 502, fmt.Sprintf("mcp unreachable: %v", err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		jsonErr(w, 502, strings.TrimSpace(string(raw)))
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 0, 1<<20), 8<<20)
	notify := r.Context().Done()
	type line struct {
		s   string
		err error
	}
	lines := make(chan line, 64)
	go func() {
		defer close(lines)
		for sc.Scan() {
			select {
			case lines <- line{s: sc.Text()}:
			case <-notify:
				return
			}
		}
	}()
	for {
		select {
		case l, ok := <-lines:
			if !ok {
				return
			}
			fmt.Fprintf(w, "%s\n", l.s)
			flusher.Flush()
		case <-notify:
			return
		}
	}
}
