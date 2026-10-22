package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"atlas/line/internal/tenant"
)

// notes reads HIS NOTES -- SecondBrain, the place he keeps everything he knows
// -- through the door, so the glass and the seats can read what the record has
// become instead of the files it is leaving.
//
// WHY THIS EXISTS. His ruling, 2026-10-22: "simplify the fuck out of
// atlas/manjuel so that it just does the coding/workflows/agentic stuff there
// and uses the secondbrain as essentially a repository for any and all
// information relating to what i am working on." Stage 1 mirrored the record
// into his notes (the core's /mirror sync). Stage 2, set going 2026-10-23 on
// his word ("let's do A and then B as proposed"), is that what reads the files
// reads the notes: the glass's What's left, Documents and Audit Ledger, and
// Steward. Nothing in this door could reach a note. This is the one reader they
// all share.
//
// HOW IT REACHES THEM. The notes kernel he already declared for the core, in
// the world's .env (or the door's environment):
//
//	MANJUEL_MCP_NOTES_CMD=python worlds/jesster/notes/cli.py --vault worlds/notes mcp
//
// The core speaks to it as an MCP server over stdin. The door's one spawn
// contract closes every child's stdin (spawn.go), so this does not: the same
// command, its trailing "mcp" replaced by the kernel's own one-shot form
// ("call <method> <json>"), run in the world's Home with a deadline. One
// declaration, two readers, no second place to configure.
//
// WHAT IT WILL NOT DO, and these are the whole design:
//
//   - IT READS. `Writes: false`. Four actions, each mapped to one kernel
//     method that only reads: search, list, get, context. A caller's string
//     never names a method; it names an action, and an action not in the
//     table is refused by name. The record's writers stay where they are (LAW 8).
//   - NOTHING IS DECLARED HERE. No second variable, no path guessed. With no
//     MANJUEL_MCP_NOTES_CMD the answer is a refusal that says so.
//   - RULE 7. The declaration is a command, not a key; it is still never
//     returned, logged or echoed, and a value from .env goes nowhere but argv.
//   - ONE LEVEL UP, AND NO FURTHER: the world's .env, then its ground's, the
//     same bound dial() keeps (gitstate.go).
//
// Stdlib only, per the atlas law.

const (
	notesTimeout  = 30 * time.Second
	notesMaxBytes = 4 << 20
	notesDeclared = "MANJUEL_MCP_NOTES_CMD"
)

// notesMethods is the whole surface: an action, and the kernel method that
// answers it. Every one of them only reads.
var notesMethods = map[string]string{
	"search":  "notes.search", // q, tag?, limit?   -- best matches first
	"list":    "notes.list",   // tag?, limit?, full? -- newest first; full gives each note whole
	"get":     "notes.get",    // id                -- one note, whole, with its links
	"context": "context.for",  // q                 -- the briefing the core hands a turn
}

// notesRun is the child, behind a seam so the strokes can stand in for the
// kernel without a Python on the machine.
var notesRun = func(argv []string, dir string) spawnResult {
	return spawn(argv[0], argv[1:], spawnOpts{Dir: dir, Timeout: notesTimeout, MaxBytes: notesMaxBytes, KillTree: true})
}

func toolNotes(t tenant.Tenant, args map[string]any) (string, error) {
	action := strings.ToLower(strings.TrimSpace(str(args, "action")))
	q := strings.TrimSpace(str(args, "q"))
	if action == "" {
		action = "list"
		if q != "" {
			action = "search"
		}
	}
	method, ok := notesMethods[action]
	if !ok {
		return "", fmt.Errorf("no notes action %q. Actions: search (q, tag?, limit?), list (tag?, limit?, full?), get (id), context (q)", action)
	}

	params := map[string]any{}
	tag := strings.TrimSpace(str(args, "tag"))
	limit := notesLimit(args["limit"])
	switch action {
	case "search":
		if q == "" {
			return "", fmt.Errorf("notes search needs q, the words to look for")
		}
		params["q"] = q
		if tag != "" {
			params["tag"] = tag
		}
		params["limit"] = limit
	case "list":
		if tag != "" {
			params["tag"] = tag
		}
		params["limit"] = limit
		if full, _ := args["full"].(bool); full {
			params["full"] = true // a kernel from before 2026-10-23 refuses it, by name
		}
	case "get":
		id := strings.TrimSpace(str(args, "id"))
		if id == "" {
			return "", fmt.Errorf("notes get needs id, the note's id")
		}
		params["id"] = id
	case "context":
		if q == "" {
			return "", fmt.Errorf("notes context needs q, the question to brief for")
		}
		params["q"] = q
	}

	argv, why := notesCommand(t.Home)
	if why != "" {
		return "", fmt.Errorf("%s", why)
	}
	p, _ := json.Marshal(params)
	argv = append(argv, "call", method, string(p))

	res := notesRun(argv, t.Home)
	out := strings.TrimSpace(res.Stdout)
	if res.Err != nil {
		// The kernel's own words outrank the Go error (ADR-006 item 2).
		said := strings.TrimSpace(res.Stderr)
		if len(said) > 600 {
			said = said[len(said)-600:]
		}
		if said == "" {
			said = res.Err.Error()
		}
		return "", fmt.Errorf("his notes did not answer %s: %s", action, said)
	}
	if !json.Valid([]byte(out)) {
		return "", fmt.Errorf("his notes answered %s with something that is not JSON; nothing is claimed", action)
	}
	return marshal(map[string]any{
		"world": t.Name, "action": action, "result": json.RawMessage(out),
	})
}

func notesLimit(v any) int {
	n := 50
	switch x := v.(type) {
	case float64:
		n = int(x)
	case int:
		n = x
	case string:
		fmt.Sscan(x, &n)
	}
	if n < 1 {
		n = 1
	}
	if n > 500 {
		n = 500
	}
	return n
}

// notesCommand is the declared kernel, as argv, with its trailing "mcp" taken
// off so a one-shot call can take its place.
func notesCommand(home string) ([]string, string) {
	raw := strings.TrimSpace(os.Getenv(notesDeclared))
	if raw == "" {
		for _, dir := range []string{home, filepath.Dir(home)} {
			if v := envFileValue(filepath.Join(dir, ".env"), notesDeclared); v != "" {
				raw = v
				break
			}
		}
	}
	if raw == "" {
		return nil, "his notes are not declared (" + notesDeclared + " in this world's .env): the door has no second brain to read"
	}
	argv := splitCommand(raw) // tools.go: the door's own split, quotes kept whole, backslashes as they are
	if len(argv) < 2 || !strings.EqualFold(argv[len(argv)-1], "mcp") {
		return nil, notesDeclared + " is not the notes kernel's MCP command (it should end in \"mcp\"); the door will not guess another"
	}
	return argv[:len(argv)-1], ""
}

// envFileValue reads one variable's value out of one .env. Absent file, key or
// value all answer "". RULE 7: the value goes back to the caller and nowhere else.
func envFileValue(path, name string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "export "))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) != name {
			continue
		}
		val = strings.TrimSpace(val)
		if len(val) >= 2 && (val[0] == '"' || val[0] == '\'') && val[len(val)-1] == val[0] {
			val = val[1 : len(val)-1]
		}
		return val
	}
	return ""
}
