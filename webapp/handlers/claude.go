// THE CLAUDE TAB (2026-10-22, his word: "I hit inspector, boom, it opens up a
// little side tab. And then I can set the model ... a little box for ...
// instructions for this session. And then ... talk to you directly and then
// watch you work. In that little tab.").
//
// The tab runs the Claude Code harness itself, headless (`claude -p` with
// stream-json), in the ground, so it reads the same CLAUDE.md, laws and
// memory as the Claude he talks to in the desktop app. Every event the harness
// emits -- its words, each tool it calls, each result -- is piped to the page
// as it happens.
//
// ITS HANDS ARE THE DOOR'S. It is given the ground's .mcp.json (Atlas's door
// through bin/atlas_bridge.py, and his notes) and nothing else: no shell, no
// file edits of its own. A write it wants goes through the door as a client and
// parks as an Approve card in his dock, the same as every other caller.
//
// THE MODEL IS HIS PICK. A `claude-*` model runs on the Anthropic key he put
// into SecondBrain (Integrations), read from its settings file into the
// child's environment only -- never printed, logged or put on a command line
// (RULE 7) -- or on his Claude login when no key is there. Any other name is a
// local Ollama model: the harness is pointed at Ollama's Anthropic-shaped API
// on loopback, so the same tab keeps working on his own card.
package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// ClaudeModels are the Claude models the tab offers, newest first.
var ClaudeModels = []string{"claude-opus-5-5", "claude-sonnet-5-5", "claude-haiku-5-5"}

const ollamaURL = "http://127.0.0.1:11434"

var claudeRun struct {
	sync.Mutex
	cmd *exec.Cmd
}

// claudeGround is the ground the glass serves: two up from atlas/webapp.
func claudeGround() string {
	g, _ := filepath.Abs(filepath.Join("..", ".."))
	if _, err := os.Stat(filepath.Join(g, "manjuel.py")); err == nil {
		return g
	}
	return ""
}

// claudeBin finds the Claude Code program: ATLAS_CLAUDE_BIN, else the newest
// copy the Claude desktop app keeps, else `claude` on the PATH.
func claudeBin() string {
	if p := os.Getenv("ATLAS_CLAUDE_BIN"); p != "" {
		return p
	}
	found, _ := filepath.Glob(filepath.Join(os.Getenv("LOCALAPPDATA"), "Claude-Data", "*", "claude-code", "*", "*", "claude.exe"))
	sort.Slice(found, func(i, j int) bool {
		a, _ := os.Stat(found[i])
		b, _ := os.Stat(found[j])
		return a != nil && b != nil && a.ModTime().After(b.ModTime())
	})
	if len(found) > 0 {
		return found[0]
	}
	if p, err := exec.LookPath("claude"); err == nil {
		return p
	}
	return ""
}

// anthropicKey is the key in SecondBrain's Integrations, or "".
func anthropicKey() string {
	b, err := os.ReadFile(filepath.Join(os.Getenv("APPDATA"), "SecondBrain", "settings.json"))
	if err != nil {
		return ""
	}
	var s struct {
		Integrations struct {
			Anthropic struct {
				Key string `json:"key"`
			} `json:"anthropic"`
		} `json:"integrations"`
	}
	if json.Unmarshal(b, &s) != nil {
		return ""
	}
	return strings.TrimSpace(s.Integrations.Anthropic.Key)
}

func isClaudeModel(m string) bool { return strings.HasPrefix(m, "claude-") }

// ollamaModels lists what the local rack holds, or nothing when it is down.
func ollamaModels() []string {
	c := http.Client{Timeout: 2 * time.Second}
	resp, err := c.Get(ollamaURL + "/api/tags")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var t struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if json.NewDecoder(resp.Body).Decode(&t) != nil {
		return nil
	}
	out := []string{}
	for _, m := range t.Models {
		if !strings.Contains(m.Name, "embed") {
			out = append(out, m.Name)
		}
	}
	sort.Strings(out)
	return out
}

// ClaudeInfo answers GET /api/claude/info: what the tab can offer.
func (h *Handlers) ClaudeInfo(w http.ResponseWriter, r *http.Request) {
	bin := claudeBin()
	jsonResp(w, map[string]any{
		"found":   bin != "",
		"claude":  ClaudeModels,
		"local":   ollamaModels(),
		"key":     anthropicKey() != "",
		"ground":  claudeGround() != "",
		"running": claudeBusy(),
	})
}

func claudeBusy() bool {
	claudeRun.Lock()
	defer claudeRun.Unlock()
	return claudeRun.cmd != nil
}

// claudeArgs is the harness's command line: headless, streaming, the door and
// the notes as its only hands, his instructions appended to its system prompt.
func claudeArgs(message, model, instructions, session, ground string) []string {
	args := []string{"-p", message, "--output-format", "stream-json", "--verbose",
		"--mcp-config", filepath.Join(ground, ".mcp.json"), "--strict-mcp-config",
		"--allowedTools", "mcp__atlas__*", "mcp__notes__*", "Read", "Grep", "Glob",
		"--disallowedTools", "Bash", "PowerShell", "Edit", "Write", "NotebookEdit", "WebFetch", "WebSearch"}
	if model != "" {
		args = append(args, "--model", model)
	}
	if strings.TrimSpace(instructions) != "" {
		args = append(args, "--append-system-prompt", instructions)
	}
	if session != "" {
		args = append(args, "--resume", session)
	}
	return args
}

// claudeEnv is the child's environment for the model picked.
func claudeEnv(model string) []string {
	env := []string{}
	for _, kv := range os.Environ() {
		k := strings.ToUpper(strings.SplitN(kv, "=", 2)[0])
		if k == "ANTHROPIC_API_KEY" || k == "ANTHROPIC_BASE_URL" || k == "ANTHROPIC_AUTH_TOKEN" {
			continue
		}
		env = append(env, kv)
	}
	if isClaudeModel(model) || model == "" {
		if key := anthropicKey(); key != "" {
			env = append(env, "ANTHROPIC_API_KEY="+key)
		}
		return env
	}
	return append(env, "ANTHROPIC_BASE_URL="+ollamaURL, "ANTHROPIC_AUTH_TOKEN=ollama", "ANTHROPIC_API_KEY=")
}

// ClaudeTurn answers POST /api/claude/turn: one message to the harness, every
// event it emits streamed back as server-sent events, one turn at a time.
func (h *Handlers) ClaudeTurn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Message      string `json:"message"`
		Model        string `json:"model"`
		Instructions string `json:"instructions"`
		Session      string `json:"session"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Message) == "" {
		jsonErr(w, 400, "say something")
		return
	}
	bin, ground := claudeBin(), claudeGround()
	if bin == "" {
		jsonErr(w, 503, "Claude Code was not found (set ATLAS_CLAUDE_BIN, or install the Claude desktop app)")
		return
	}
	if ground == "" {
		jsonErr(w, 503, "the ground (the folder with manjuel.py) was not found two up from the glass")
		return
	}
	cmd := exec.Command(bin, claudeArgs(body.Message, body.Model, body.Instructions, body.Session, ground)...)
	cmd.Dir = ground
	cmd.Env = claudeEnv(body.Model)
	out, err := cmd.StdoutPipe()
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	var stderr strings.Builder
	cmd.Stderr = &stderr
	claudeRun.Lock()
	if claudeRun.cmd != nil {
		claudeRun.Unlock()
		jsonErr(w, 409, "Claude is still working on the last message; stop it or wait")
		return
	}
	if err := hidden(cmd).Start(); err != nil {
		claudeRun.Unlock()
		jsonErr(w, 500, "Claude Code did not start: "+err.Error())
		return
	}
	claudeRun.cmd = cmd
	claudeRun.Unlock()
	defer func() {
		claudeRun.Lock()
		claudeRun.cmd = nil
		claudeRun.Unlock()
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	flusher, _ := w.(http.Flusher)
	sc := bufio.NewScanner(out)
	sc.Buffer(make([]byte, 0, 1<<20), 32<<20)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		fmt.Fprintf(w, "data: %s\n\n", line)
		if flusher != nil {
			flusher.Flush()
		}
	}
	werr := cmd.Wait()
	end := map[string]any{"ok": werr == nil}
	if werr != nil {
		tail := strings.TrimSpace(stderr.String())
		if len(tail) > 2000 {
			tail = tail[len(tail)-2000:]
		}
		end["error"] = strings.TrimSpace(werr.Error() + "\n" + tail)
	}
	b, _ := json.Marshal(end)
	fmt.Fprintf(w, "event: done\ndata: %s\n\n", b)
	if flusher != nil {
		flusher.Flush()
	}
}

// ClaudeStop answers POST /api/claude/stop: the turn in flight is ended.
func (h *Handlers) ClaudeStop(w http.ResponseWriter, r *http.Request) {
	claudeRun.Lock()
	cmd := claudeRun.cmd
	claudeRun.Unlock()
	if cmd == nil || cmd.Process == nil {
		jsonResp(w, map[string]any{"stopped": false})
		return
	}
	_ = cmd.Process.Kill()
	jsonResp(w, map[string]any{"stopped": true})
}
