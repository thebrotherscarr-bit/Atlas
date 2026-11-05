package protocol

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"atlas/line/internal/tenant"
	"atlas/line/internal/tools"
)

// serve runs the stdio server over the given lines against one temp tenant and
// returns one decoded answer per answered line.
func serve(t *testing.T, lines ...string) []map[string]any {
	t.Helper()
	tr := tenant.NewRegistry()
	if err := tr.Add("t", t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if err := tr.SetDefault("t"); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := Serve(strings.NewReader(strings.Join(lines, "\n")+"\n"), &out,
		ServerInfo{Name: "atlas-mcp", Version: "test"}, "the law", tools.Build(tr, tools.Options{}), tr); err != nil {
		t.Fatal(err)
	}
	var got []map[string]any
	for _, l := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if l == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			t.Fatalf("not JSON: %q", l)
		}
		got = append(got, m)
	}
	return got
}

const meta = `"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{},"io.modelcontextprotocol/clientInfo":{"name":"t","version":"0"}}`

func res(t *testing.T, m map[string]any) map[string]any {
	t.Helper()
	r, ok := m["result"].(map[string]any)
	if !ok {
		t.Fatalf("expected a result, got %v", m)
	}
	return r
}

func code(t *testing.T, m map[string]any) int {
	t.Helper()
	e, ok := m["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected an error, got %v", m)
	}
	c, _ := e["code"].(float64)
	return int(c)
}

// A 2026-07-28 CLIENT NEEDS NO HANDSHAKE (spec 2026-07-28, changelog major 2
// and 3): server/discover, then tools/list and tools/call with the _meta
// envelope; every result says resultType "complete" and who answered, the
// list carries its cache hints, and ping is gone.
func TestAModernClientListsAndCallsWithNoInitialize(t *testing.T) {
	got := serve(t,
		`{"jsonrpc":"2.0","id":1,"method":"server/discover","params":{`+meta+`}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{`+meta+`}}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"muster","arguments":{},`+meta+`}}`,
		`{"jsonrpc":"2.0","id":4,"method":"ping","params":{`+meta+`}}`,
		`{"jsonrpc":"2.0","method":"notifications/cancelled","params":{`+meta+`}}`,
		`{"jsonrpc":"2.0","id":5,"method":"tools/list","params":{`+meta+`}}`,
	)
	if len(got) != 5 {
		t.Fatalf("want five answers (the notification gets none), got %d: %v", len(got), got)
	}
	d := res(t, got[0])
	if d["resultType"] != "complete" || d["instructions"] != "the law" {
		t.Errorf("discover = %v", d)
	}
	vs, _ := d["supportedVersions"].([]any)
	if len(vs) == 0 || vs[0] != "2026-07-28" || !strings.Contains(strings.Join(toStrings(vs), ","), "2025-06-18") {
		t.Errorf("supportedVersions = %v, want 2026-07-28 first and 2025-06-18 among them", vs)
	}
	info := d["_meta"].(map[string]any)["io.modelcontextprotocol/serverInfo"].(map[string]any)
	if info["name"] != "atlas-mcp" {
		t.Errorf("serverInfo = %v", info)
	}
	l := res(t, got[1])
	if l["resultType"] != "complete" || l["cacheScope"] != "private" || l["ttlMs"] != float64(0) {
		t.Errorf("tools/list must carry resultType, ttlMs and cacheScope: %v", l)
	}
	if names(l) == "" || names(l) != names(res(t, got[4])) {
		t.Errorf("tools/list must come back in the same order every time")
	}
	c := res(t, got[2])
	if c["resultType"] != "complete" || c["isError"] != nil {
		t.Errorf("tools/call = %v", c)
	}
	if _, ok := c["ttlMs"]; ok {
		t.Errorf("a call result is not cacheable: %v", c)
	}
	if code(t, got[3]) != CodeMethodNotFound {
		t.Errorf("ping is removed in 2026-07-28")
	}
}

// AN UNKNOWN VERSION IS -32022 AND NAMES WHAT IS SPOKEN, so a client can retry
// with one (basic/versioning, "Protocol Version Negotiation").
func TestAnUnknownVersionIsRefusedWithTheSupportedList(t *testing.T) {
	got := serve(t,
		`{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{"_meta":{"io.modelcontextprotocol/protocolVersion":"1900-01-01","io.modelcontextprotocol/clientCapabilities":{}}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28"}}}`,
	)
	if code(t, got[0]) != CodeUnsupportedProtocolVersion {
		t.Fatalf("want -32022, got %v", got[0])
	}
	data := got[0]["error"].(map[string]any)["data"].(map[string]any)
	if data["requested"] != "1900-01-01" || len(data["supported"].([]any)) != len(Supported) {
		t.Errorf("data = %v", data)
	}
	if code(t, got[1]) != CodeInvalidParams {
		t.Errorf("a missing clientCapabilities is -32602, got %v", got[1])
	}
}

// THE LEGACY HANDSHAKE STILL WORKS: initialize answers the client's 2025
// version (2025-06-18 when it names none we know), ping answers, and nothing
// new is put on a legacy result.
func TestTheLegacyHandshakeStillWorks(t *testing.T) {
	got := serve(t,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25"}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"ping"}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":4,"method":"initialize","params":{"protocolVersion":"2099-01-01"}}`,
	)
	if len(got) != 4 {
		t.Fatalf("want four answers, got %v", got)
	}
	if res(t, got[0])["protocolVersion"] != "2025-11-25" {
		t.Errorf("initialize = %v", got[0])
	}
	if len(res(t, got[1])) != 0 {
		t.Errorf("ping answers an empty result: %v", got[1])
	}
	if _, ok := res(t, got[2])["resultType"]; ok {
		t.Errorf("a legacy answer is as it was: %v", got[2])
	}
	if res(t, got[3])["protocolVersion"] != LegacyDefault {
		t.Errorf("an unknown initialize version gets %s: %v", LegacyDefault, got[3])
	}
}

func toStrings(vs []any) []string {
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		s, _ := v.(string)
		out = append(out, s)
	}
	return out
}

func names(r map[string]any) string {
	ts, _ := r["tools"].([]any)
	out := make([]string, 0, len(ts))
	for _, t := range ts {
		n, _ := t.(map[string]any)["name"].(string)
		out = append(out, n)
	}
	return strings.Join(out, ",")
}
