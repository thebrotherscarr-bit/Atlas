package protocol

// BOTH MCP ERAS ON ONE LINE (spec 2026-07-28, basic/versioning, "Backward
// Compatibility with Initialization-Based Versions"): a dual-era server
// answers `initialize` with legacy semantics and serves a request carrying
// `_meta["io.modelcontextprotocol/protocolVersion"]` statelessly.
//
//	modern  2026-07-28: no handshake; version and client capabilities ride
//	        every request's _meta; server/discover is required; every result
//	        carries resultType; list results carry ttlMs and cacheScope; ping
//	        and logging/setLevel are removed.
//	legacy  2025-11-25 and earlier: initialize, then requests without the
//	        envelope -- answered exactly as before, ping included.
//
// Both of the door's transports are newline-delimited JSON-RPC (stdio, and
// httpserver's POST /rpc), neither is Streamable HTTP, so there is no header
// rung here and every answer rides the line it was asked on.

import "encoding/json"

// Modern is the stateless revision; Legacy the handshake revisions answered.
var (
	Modern    = []string{"2026-07-28"}
	Legacy    = []string{"2025-11-25", "2025-06-18", "2025-03-26", "2024-11-05"}
	Supported = append(append([]string{}, Modern...), Legacy...)
)

// LegacyDefault answers an initialize that names no version we know: the
// revision this door was built to.
const LegacyDefault = "2025-06-18"

// Error codes the 2026-07-28 spec reserves, and the JSON-RPC ones used here.
const (
	CodeUnsupportedProtocolVersion = -32022
	CodeInvalidParams              = -32602
	CodeMethodNotFound             = -32601
)

const (
	metaProtocolVersion    = "io.modelcontextprotocol/protocolVersion"
	metaClientCapabilities = "io.modelcontextprotocol/clientCapabilities"
	metaServerInfo         = "io.modelcontextprotocol/serverInfo"
)

// removed are the methods 2026-07-28 took out: method-not-found when modern.
var removed = map[string]bool{"ping": true, "logging/setLevel": true,
	"resources/subscribe": true, "resources/unsubscribe": true}

// cacheable results carry ttlMs and cacheScope (server/utilities/caching).
var cacheable = map[string]bool{"tools/list": true, "prompts/list": true,
	"resources/list": true, "resources/read": true,
	"resources/templates/list": true, "server/discover": true}

// NegotiateLegacy is initialize's answer: the client's version when it is a
// handshake revision we speak, else LegacyDefault.
func NegotiateLegacy(requested string) string {
	for _, v := range Legacy {
		if v == requested {
			return requested
		}
	}
	return LegacyDefault
}

type envelope struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params struct {
		Meta map[string]json.RawMessage `json:"_meta"`
	} `json:"params"`
}

// Dual answers one line in whichever era it was asked: `legacy` is the
// server's own handler, used as it is for a legacy line and, for a modern one,
// finished with resultType, the server's identity and the cache hints.
// A line that does not parse goes to `legacy`, which names the parse error.
func Dual(line []byte, info ServerInfo, instructions string, legacy func([]byte) any) any {
	var env envelope
	if err := json.Unmarshal(line, &env); err != nil || env.Method == "initialize" {
		return legacy(line)
	}
	rawVersion, modern := env.Params.Meta[metaProtocolVersion]
	if !modern {
		return legacy(line)
	}
	if len(env.ID) == 0 || string(env.ID) == "null" {
		return nil // a notification
	}
	if _, ok := env.Params.Meta[metaClientCapabilities]; !ok {
		return errResponse(env.ID, CodeInvalidParams,
			"params._meta is missing the required envelope key "+metaClientCapabilities)
	}
	var version string
	if err := json.Unmarshal(rawVersion, &version); err != nil {
		return errResponse(env.ID, CodeInvalidParams,
			"the protocol-version envelope value must be a string")
	}
	if !contains(Modern, version) {
		return errDataResponse(env.ID, CodeUnsupportedProtocolVersion,
			"Unsupported protocol version",
			map[string]any{"supported": Supported, "requested": version})
	}
	switch {
	case env.Method == "server/discover":
		return finish(resultResponse(env.ID, map[string]any{
			"supportedVersions": Supported,
			"capabilities":      map[string]any{"tools": map[string]any{}},
			"instructions":      instructions,
		}), env.Method, info)
	case removed[env.Method]:
		return errResponse(env.ID, CodeMethodNotFound,
			"unknown method \""+env.Method+"\" (removed in "+Modern[0]+")")
	}
	return finish(legacy(line), env.Method, info)
}

// finish puts what 2026-07-28 requires on a result: resultType "complete",
// the server's identity in its _meta, and on a cacheable one ttlMs 0 and
// cacheScope "private" -- the tool list follows who is asking (the gate) and
// the grounds carried, so nothing here is for a shared cache to keep.
func finish(resp any, method string, info ServerInfo) any {
	m, ok := resp.(map[string]any)
	if !ok {
		return resp
	}
	result, ok := m["result"].(map[string]any)
	if !ok {
		return resp
	}
	out := map[string]any{"resultType": "complete"}
	for k, v := range result {
		out[k] = v
	}
	meta := map[string]any{}
	if old, ok := result["_meta"].(map[string]any); ok {
		for k, v := range old {
			meta[k] = v
		}
	}
	meta[metaServerInfo] = info
	out["_meta"] = meta
	if cacheable[method] {
		if _, ok := out["ttlMs"]; !ok {
			out["ttlMs"] = 0
		}
		if _, ok := out["cacheScope"]; !ok {
			out["cacheScope"] = "private"
		}
	}
	m["result"] = out
	return m
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func errDataResponse(id json.RawMessage, code int, msg string, data any) map[string]any {
	r := errResponse(id, code, msg)
	r["error"] = map[string]any{"code": code, "message": msg, "data": data}
	return r
}
