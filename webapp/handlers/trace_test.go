// Which tool calls the glass keeps.
//
// His ruling, 2026-09-16, on the optimization pass's first decision: "D1 b".
// The Dashboard asks the door for muster, rack_list and proofs when it opens,
// when it is shown again and every fifteen seconds it is on screen, and every
// answer was kept whole in the trace ledger, announced to every open tab and
// toasted there. Asked as background reads, those three are answered and not
// kept. Everything else is kept as before: the same three asked any other way,
// any other tool, and any call that writes, however it is asked.
//
// Hermetic: a temp store and a loopback door that answers every call.
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

// The three reads Home.read asks for.
var dashboardReads = []string{"muster", "rack_list", "proofs"}

// answeringDoor points h at a door that answers every call with the same few
// words, and counts the calls it was asked.
func answeringDoor(t *testing.T, h *Handlers) *atomic.Int32 {
	t.Helper()
	asked := new(atomic.Int32)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		asked.Add(1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  map[string]any{"content": []map[string]string{{"type": "text", "text": "2 carried projects"}}},
		})
	}))
	t.Cleanup(srv.Close)
	h.store.SetSetting("mcp_url", srv.URL)
	return asked
}

// listen subscribes to the bus the way an open tab does.
func listen(h *Handlers) chan event {
	ch := make(chan event, 16)
	h.clientMu.Lock()
	h.clients = append(h.clients, ch)
	h.clientMu.Unlock()
	return ch
}

// ask sends one body to the glass's tool call and returns the reply.
func ask(t *testing.T, h *Handlers, body string) map[string]any {
	t.Helper()
	w := httptest.NewRecorder()
	h.CallTool(w, httptest.NewRequest("POST", "/api/tools/call", strings.NewReader(body)))
	var reply map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &reply); err != nil {
		t.Fatalf("%s: the reply is not JSON: %q", body, w.Body.String())
	}
	return reply
}

func inBackground(tool string) string {
	return `{"tool":"` + tool + `","args":{},"background":true}`
}

// kept reports whether a reply names a trace the ledger holds.
func kept(h *Handlers, reply map[string]any) bool {
	id, _ := reply["trace_id"].(string)
	return id != "" && h.store.Get(id) != nil
}

// THE DASHBOARD'S OWN READS ARE NOT KEPT: no trace in the ledger, and no trace
// named in the reply, because there is none to name.
func TestTheDashboardsOwnReadsAreNotKept(t *testing.T) {
	h := newTestHandlers(t)
	answeringDoor(t, h)
	for _, tool := range dashboardReads {
		if reply := ask(t, h, inBackground(tool)); reply["trace_id"] != nil {
			t.Fatalf("%s asked in the background came back naming trace %v", tool, reply["trace_id"])
		}
	}
	if n := h.store.Count(); n != 0 {
		t.Fatalf("three background reads left %d traces in the ledger", n)
	}
}

// NOR ANNOUNCED: no open tab hears of them, so none shows "New trace recorded".
func TestTheDashboardsOwnReadsAreNotAnnounced(t *testing.T) {
	h := newTestHandlers(t)
	answeringDoor(t, h)
	tab := listen(h)
	for _, tool := range dashboardReads {
		ask(t, h, inBackground(tool))
	}
	if len(tab) != 0 {
		e := <-tab
		t.Fatalf("an open tab was told of a background read: %s", e.Type)
	}
}

// AND STILL ANSWERED, by the door: leaving the ledger out must not leave the
// page without its reading.
func TestTheDashboardsOwnReadsAreStillAnswered(t *testing.T) {
	h := newTestHandlers(t)
	asked := answeringDoor(t, h)
	for _, tool := range dashboardReads {
		reply := ask(t, h, inBackground(tool))
		if out, _ := reply["output"].(string); !strings.Contains(out, "2 carried projects") {
			t.Fatalf("%s asked in the background answered %q, not the door's words", tool, out)
		}
		if reply["hash"] == nil || reply["duration_ms"] == nil {
			t.Fatalf("%s asked in the background lost its hash or its duration: %v", tool, reply)
		}
	}
	if n := asked.Load(); n != int32(len(dashboardReads)) {
		t.Fatalf("the door was asked %d times for %d reads", n, len(dashboardReads))
	}
}

// THE SAME THREE ASKED ANY OTHER WAY ARE KEPT -- Records reading proofs, the end
// of a turn repainting it, a call from the Tools page -- and each open tab is
// told, as before.
func TestTheSameReadsAskedPlainlyAreKept(t *testing.T) {
	h := newTestHandlers(t)
	answeringDoor(t, h)
	tab := listen(h)
	for _, tool := range dashboardReads {
		for _, body := range []string{
			`{"tool":"` + tool + `","args":{}}`,
			`{"tool":"` + tool + `","args":{},"background":false}`,
		} {
			if !kept(h, ask(t, h, body)) {
				t.Fatalf("%s was not kept", body)
			}
			if len(tab) != 1 {
				t.Fatalf("%s was kept and %d events reached an open tab, not one", body, len(tab))
			}
			if e := <-tab; e.Type != "trace_added" {
				t.Fatalf("%s announced %q, not trace_added", body, e.Type)
			}
		}
	}
	if n := h.store.Count(); n != 2*len(dashboardReads) {
		t.Fatalf("the ledger holds %d traces for %d plain calls", n, 2*len(dashboardReads))
	}
}

// NO OTHER TOOL LEAVES THE RECORD BY ASKING TO. A call that writes, marked
// background, is kept and announced; so is a read the Dashboard does not poll.
func TestOnlyTheDashboardsThreeReadsCanBeLeftOut(t *testing.T) {
	h := newTestHandlers(t)
	answeringDoor(t, h)
	tab := listen(h)
	for _, tool := range []string{"git_commit", "env_close", "records", "seats"} {
		if !kept(h, ask(t, h, inBackground(tool))) {
			t.Fatalf("%s marked background was left out of the ledger", tool)
		}
		if len(tab) != 1 {
			t.Fatalf("%s marked background was kept and %d events reached an open tab, not one", tool, len(tab))
		}
		<-tab
	}
}
