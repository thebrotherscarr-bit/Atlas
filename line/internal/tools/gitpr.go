package tools

// PULL REQUESTS, THROUGH THE GLASS (2026-10-07; his ruling on a card: "Open the wall for PRs").
//
// gitctl.go kept GitHub-side objects out of the door: "a separate wall to open deliberately, not a
// dependency to acquire by accident" (RULE 4, atlas law 6). He opened that wall for pull requests and
// for nothing else. git_pr lists a world's open pull requests with their checks, and opens one from the
// line of work he stands on into the main line, on his click. It never merges: a line still lands as it
// always has (Land onto main, fast-forward, then Send main), and GitHub marks the pull request merged once
// the main line holds its saves. It runs gh, the GitHub CLI already signed in on this machine, behind the
// same wall as every send (GIT_REMOTE), and it is his glass's alone.

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"atlas/line/internal/tenant"
)

// ghRun runs the GitHub CLI in a world's ground: a variable so the strokes can stand in for it, and no stroke reaches GitHub.
var ghRun = func(t tenant.Tenant, timeout time.Duration, args ...string) (string, error) {
	res := spawn("gh", args, spawnOpts{Dir: t.Home, Timeout: timeout, Env: ghEnv()})
	return strings.TrimRight(res.Combined, "\r\n"), res.Err
}

// prCheck is one check on a pull request as gh reports it: a check run carries a status and a
// conclusion, a commit status a state.
type prCheck struct {
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	State      string `json:"state"`
}

// prChecks is a pull request's checks as one word: none, pending, failed or passed.
func prChecks(checks []prCheck) string {
	if len(checks) == 0 {
		return "none"
	}

	pending := false
	for _, c := range checks {
		switch strings.ToUpper(c.Conclusion) {
		case "FAILURE", "CANCELLED", "TIMED_OUT", "ACTION_REQUIRED", "STARTUP_FAILURE":
			return "failed"
		}
		switch strings.ToUpper(c.State) {
		case "FAILURE", "ERROR":
			return "failed"
		case "PENDING", "EXPECTED":
			pending = true
		}
		if c.State == "" && strings.ToUpper(c.Status) != "COMPLETED" {
			pending = true
		}
	}

	if pending {
		return "pending"
	}
	return "passed"
}

// prRow is one open pull request as the glass shows it.
type prRow struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Line   string `json:"line"`
	Into   string `json:"into"`
	URL    string `json:"url"`
	Checks string `json:"checks"`
}

// prList is the world's open pull requests, each with its checks as one word.
func prList(t tenant.Tenant) (string, error) {
	out, err := ghRun(t, 30*time.Second, "pr", "list", "--state", "open", "--json", "number,title,headRefName,baseRefName,url,statusCheckRollup")
	if err != nil {
		return "Refused: gh could not list the pull requests -- " + firstLine(out), nil
	}

	var raw []struct {
		Number            int       `json:"number"`
		Title             string    `json:"title"`
		HeadRefName       string    `json:"headRefName"`
		BaseRefName       string    `json:"baseRefName"`
		URL               string    `json:"url"`
		StatusCheckRollup []prCheck `json:"statusCheckRollup"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return "Refused: gh answered in a shape this cannot read -- " + firstLine(out), nil
	}

	rows := []prRow{}
	for _, p := range raw {
		rows = append(rows, prRow{
			Number: p.Number,
			Title:  p.Title,
			Line:   p.HeadRefName,
			Into:   p.BaseRefName,
			URL:    p.URL,
			Checks: prChecks(p.StatusCheckRollup),
		})
	}

	b, err := json.MarshalIndent(map[string]any{"pull_requests": rows}, "", " ")
	return string(b), err
}

// prOpen opens a pull request from the line of work he stands on into the main line: refused on the main line, over unsaved work, and for a line GitHub does not have exactly as it stands here.
func prOpen(t tenant.Tenant) (string, error) {
	line, err := currentBranch(t)
	if err != nil || line == "" {
		return "Refused: this world stands on no line of work, so there is nothing to ask to land.", nil
	}
	mainName := mainLine(t)
	if line == mainName {
		return fmt.Sprintf("Refused: you are standing on %q, the main line. A pull request takes a line of work into it: move to one first.", mainName), nil
	}
	if dirty(t) {
		return "Refused: there is unsaved work here. Save it first; a pull request carries saves, not a tree.", nil
	}
	here, _ := gitRun(t, 10*time.Second, "rev-parse", "HEAD")
	there, err := gitRun(t, 10*time.Second, "rev-parse", "@{upstream}")
	if err != nil || strings.TrimSpace(here) != strings.TrimSpace(there) {
		return fmt.Sprintf("Refused: GitHub does not have %q exactly as it stands here. Send it first, then open the pull request.", line), nil
	}
	title, _ := gitRun(t, 10*time.Second, "log", "-1", "--format=%s")
	saves, _ := gitRun(t, 10*time.Second, "log", "--reverse", "--format=- %s", line, "--not", mainName)
	body := "The saves on " + line + ":\n\n" + saves
	out, err := ghRun(t, 60*time.Second, "pr", "create", "--base", mainName, "--head", line, "--title", strings.TrimSpace(title), "--body", body)
	if err != nil {
		return "Refused: gh could not open the pull request -- " + firstLine(out), nil
	}
	url := strings.TrimSpace(out)
	url = strings.Split(url, "\n")[len(strings.Split(url, "\n"))-1]
	return fmt.Sprintf("Opened a pull request for %q into %q: %s. CI runs on it; when its checks pass, land it as always (move to %q, Land onto main, Send) and GitHub marks it merged.", line, mainName, url, mainName), nil
}

// toolGitPR is git_pr: list (a read, and the default) or open (his click).
func toolGitPR(t tenant.Tenant, args map[string]any) (string, error) {
	if !callerOf(args).Service {
		return "", errors.New("git_pr is the operator's own hand: his glass alone may call it")
	}
	if !isRepo(t) {
		return notARepo, nil
	}
	if !dial(t.Home, "GIT_REMOTE") {
		return walled("Pull requests"), nil
	}
	switch action := actionOf(args); action {
	case "", "list":
		return prList(t)
	case "open":
		return prOpen(t)
	default:
		return fmt.Sprintf("Refused: %q is not something this does. It lists the open pull requests (list) or opens one from the line you stand on (open). It never merges: land the line as always.", action), nil
	}
}

// gitPRTool is git_pr's declaration; tools.go registers it with one line, r.add(gitPRTool()).
func gitPRTool() Tool {
	return Tool{
		Name:        "git_pr",
		Description: "pull requests, the one GitHub-side object his ruling of 2026-10-07 let into the door: list a world's open pull requests with their checks as one word (none, pending, passed, failed), or open one from the line of work he stands on into the main line, refused on the main line, over unsaved work, and for a line GitHub does not have exactly as it stands. It never merges: a line lands as always (Land onto main, then Send). Behind the same wall as every send (GIT_REMOTE). His glass alone may call it",
		Writes:      true,
		Reads:       []string{"list"},
		ServiceOnly: true,
		Args:        []string{"action?", "project?"},
		Fn:          toolGitPR,
	}
}

// ghEnv keeps gh from stopping for a human, paging, or painting colour.
func ghEnv() []string {
	return []string{
		"GH_PROMPT_DISABLED=1",
		"GH_PAGER=cat",
		"NO_COLOR=1",
		"GH_NO_UPDATE_NOTIFIER=1",
	}
}
