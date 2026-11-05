// ask_steward asks Steward, the MAF agent, over his MCP server (2026-11-04).
// These strokes stand up a fake MCP server with httptest and point the tool at
// it: the answer comes back, an unreachable Steward is said (and nothing else
// takes the question), and an event-stream answer is read the same as JSON.
package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"atlas/line/internal/tenant"
)

type fakeSteward struct {
	sse      bool
	gotTask  string
	gotSess  string
	gotTool  string
	gotKey   string
	methods  []string
	accepted string
}

func (f *fakeSteward) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f.accepted = r.Header.Get("Accept")
		f.gotKey = r.Header.Get("Authorization")
		raw, _ := io.ReadAll(r.Body)
		var req struct {
			Method string         `json:"method"`
			Params map[string]any `json:"params"`
		}
		if err := json.Unmarshal(raw, &req); err != nil {
			t.Errorf("fake Steward: request not JSON: %v", err)
		}
		f.methods = append(f.methods, req.Method)
		var result any
		switch req.Method {
		case "tools/list":
			result = map[string]any{"tools": []any{map[string]any{
				"name": "Steward",
				"inputSchema": map[string]any{
					"type":       "object",
					"properties": map[string]any{"task": map[string]any{}, "session_id": map[string]any{}},
					"required":   []any{"task"},
				},
			}}}
		case "tools/call":
			f.gotTool, _ = req.Params["name"].(string)
			a, _ := req.Params["arguments"].(map[string]any)
			f.gotTask, _ = a["task"].(string)
			f.gotSess, _ = a["session_id"].(string)
			result = map[string]any{"content": []any{map[string]any{"type": "text", "text": "STEWARD SAYS: " + f.gotTask}}}
		default:
			t.Errorf("fake Steward: unexpected method %q", req.Method)
		}
		b, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "result": result})
		if f.sse {
			w.Header().Set("Content-Type", "text/event-stream")
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", b)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func pointStewardAt(t *testing.T, u string) {
	t.Helper()
	old := stewardAt
	stewardAt = func() string { return u }
	t.Cleanup(func() { stewardAt = old })
}

func TestAskStewardAnswersFromStewardOverMCP(t *testing.T) {
	for _, sse := range []bool{false, true} {
		f := &fakeSteward{sse: sse}
		srv := httptest.NewServer(f.handler(t))
		pointStewardAt(t, srv.URL+"/mcp")
		t.Setenv(stewardKeyDial, "")
		out, err := toolAskSteward(tenant.Tenant{Name: "atlas", Home: t.TempDir()},
			map[string]any{"question": "does the door work"})
		srv.Close()
		if err != nil {
			t.Fatalf("sse=%v: ask_steward refused: %v", sse, err)
		}
		if out != "STEWARD SAYS: does the door work" {
			t.Fatalf("sse=%v: answer not Steward's: %q", sse, out)
		}
		if strings.Join(f.methods, ",") != "tools/list,tools/call" || f.gotTool != "Steward" {
			t.Fatalf("sse=%v: wrong conversation: %v tool=%q", sse, f.methods, f.gotTool)
		}
		if f.gotSess != "atlas-atlas" {
			t.Fatalf("sse=%v: session_id = %q, want atlas-atlas", sse, f.gotSess)
		}
		if !strings.Contains(f.accepted, "application/json") || !strings.Contains(f.accepted, "text/event-stream") {
			t.Fatalf("sse=%v: Accept = %q, must list both kinds", sse, f.accepted)
		}
		if f.gotKey != "" {
			t.Fatalf("sse=%v: a bearer rode with no key declared", sse)
		}
	}
}

func TestAskStewardCarriesSessionAndKey(t *testing.T) {
	f := &fakeSteward{}
	srv := httptest.NewServer(f.handler(t))
	defer srv.Close()
	pointStewardAt(t, srv.URL+"/mcp")
	t.Setenv(stewardKeyDial, "k-test")
	out, err := toolAskSteward(tenant.Tenant{Name: "atlas"},
		map[string]any{"question": "hello", "session_id": "s-1"})
	if err != nil {
		t.Fatal(err)
	}
	if f.gotSess != "s-1" || f.gotKey != "Bearer k-test" {
		t.Fatalf("session %q, auth %q", f.gotSess, f.gotKey)
	}
	if strings.Contains(out, "k-test") {
		t.Fatal("the key leaked into the answer")
	}
}

// RED WHEN UNPLUGGED: Steward not listening is said by name and address, and
// the tenant's wired engine is NOT run in his place.
func TestAskStewardUnreachableIsSaidWithNoFallback(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := "http://" + ln.Addr().String() + "/mcp"
	ln.Close() // nothing listens there now
	pointStewardAt(t, addr)
	marker := t.TempDir() + "/engine-ran"
	_, err = toolAskSteward(tenant.Tenant{Name: "atlas", Home: t.TempDir(), Engine: "touch " + marker},
		map[string]any{"question": "anyone home?"})
	if err == nil {
		t.Fatal("ask_steward answered with Steward unplugged")
	}
	if !strings.Contains(err.Error(), "Steward is not reachable at "+addr) {
		t.Fatalf("refusal does not say so: %v", err)
	}
	if strings.Contains(err.Error(), "engine") {
		t.Fatalf("refusal mentions an engine: %v", err)
	}
	if _, serr := os.Stat(marker); serr == nil {
		t.Fatal("the wired engine ran in Steward's place")
	}
}

func TestAskStewardRefusesOffMachineAddress(t *testing.T) {
	pointStewardAt(t, "http://example.com/mcp")
	_, err := toolAskSteward(tenant.Tenant{Name: "atlas"}, map[string]any{"question": "hi"})
	if err == nil || !strings.Contains(err.Error(), "not on this machine") {
		t.Fatalf("off-machine address not refused: %v", err)
	}
}

func TestAskStewardDefaultAddress(t *testing.T) {
	t.Setenv(stewardDial, "")
	if got := stewardAt(); got != "http://127.0.0.1:43121/mcp" {
		t.Fatalf("default = %q", got)
	}
	t.Setenv(stewardDial, "http://127.0.0.1:9/mcp")
	if got := stewardAt(); got != "http://127.0.0.1:9/mcp" {
		t.Fatalf("dial not read: %q", got)
	}
}

func TestAskStewardWaitsPastTheApprovalCard(t *testing.T) {
	if stewardWait != 900*time.Second || askTimeout != 900*time.Second {
		t.Fatalf("Steward's turn is waited %v (askTimeout %v), want 900s past a 660s card", stewardWait, askTimeout)
	}
}
