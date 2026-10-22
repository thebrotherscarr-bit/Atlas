// ATLAS UPDATES ITSELF (2026-10-22, his word: "atlas should auto-update just
// like the secondbrain does").
//
// GET /api/update/check asks GitHub's main whether there is anything newer
// than the commit this glass was built from -- Go writes that commit into the
// program itself (debug.ReadBuildInfo, vcs.revision), so nothing has to be
// stamped by hand. POST /api/update/apply starts ../update.ps1, which takes
// main forward, builds the door and the glass, swaps them in and starts both
// again; the page waits for the glass to answer and reloads.
package handlers

import (
	"context"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// builtFrom is the commit this program was built from, and whether the tree
// held unsaved changes then; "" when it was built outside a repository.
func builtFrom() (rev string, modified bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	return rev, modified
}

// atlasDir is the repository the glass was built in: one up from atlas/webapp.
func atlasDir() string {
	d, _ := filepath.Abs("..")
	return d
}

func gitIn(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(cmd.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := hidden(cmd).Output()
	return strings.TrimSpace(string(out)), err
}

// UpdateCheck answers GET /api/update/check.
func (h *Handlers) UpdateCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	dir := atlasDir()
	running, modified := builtFrom()
	out := map[string]any{"running": short(running), "modified": modified}
	if _, err := gitIn(ctx, dir, "fetch", "-q", "origin", "main"); err != nil {
		out["error"] = "GitHub could not be asked (" + err.Error() + ")"
		jsonResp(w, out)
		return
	}
	latest, err := gitIn(ctx, dir, "rev-parse", "origin/main")
	if err != nil {
		out["error"] = err.Error()
		jsonResp(w, out)
		return
	}
	out["latest"] = short(latest)
	base := running
	if base == "" {
		base, _ = gitIn(ctx, dir, "rev-parse", "HEAD")
	}
	n := 0
	if s, err := gitIn(ctx, dir, "rev-list", "--count", base+"..origin/main"); err == nil {
		n, _ = strconv.Atoi(s)
	}
	out["behind"] = n
	out["available"] = n > 0
	if n > 0 {
		subjects, _ := gitIn(ctx, dir, "log", "--format=%s", "-n", "10", base+"..origin/main")
		out["changes"] = strings.Split(subjects, "\n")
	}
	jsonResp(w, out)
}

// UpdateApply answers POST /api/update/apply: update.ps1 is started on its own
// and outlives this glass, which it stops and starts again.
func (h *Handlers) UpdateApply(w http.ResponseWriter, r *http.Request) {
	script := filepath.Join(atlasDir(), "update.ps1")
	shell := "pwsh"
	if _, err := exec.LookPath(shell); err != nil {
		shell = "powershell"
	}
	cmd := exec.Command(shell, "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", script)
	cmd.Dir = atlasDir()
	if err := hidden(cmd).Start(); err != nil {
		jsonErr(w, 500, "the update did not start: "+err.Error())
		return
	}
	go cmd.Wait()
	jsonResp(w, map[string]any{"started": true, "log": filepath.Join(atlasDir(), "update.log")})
}

func short(rev string) string {
	if len(rev) > 7 {
		return rev[:7]
	}
	return rev
}
