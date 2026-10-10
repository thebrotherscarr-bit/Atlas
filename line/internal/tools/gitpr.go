package tools

// PULL REQUESTS, THROUGH THE GLASS (2026-10-07; his ruling on a card: "Open the wall for PRs").
//
// gitctl.go kept GitHub-side objects out of the door: "a separate wall to open deliberately, not a
// dependency to acquire by accident" (RULE 4, atlas law 6). He opened that wall for pull requests and
// for nothing else. git_pr lists a world's open pull requests with their checks, opens one from the line
// of work he stands on into the main line, and merges one on GitHub, each on his click. It runs gh, the
// GitHub CLI already signed in on this machine, behind the same wall as every send (GIT_REMOTE), and it is
// his glass's alone.
//
// IT MERGES NOW (2026-10-09, WHAT'S LEFT I1's fourth step). Until then it never merged: a line landed as it
// always had (Land onto main, fast-forward, then Send main) and GitHub marked the pull request merged once
// the main line held its saves, because the door's founding law kept `merge` out of its table. His word
// superseded that law for this -- "we have to be able to merge, thats before the github functionality was
// complete ... supersede that one" -- and his card said how: "Merge on GitHub when green", as a merge
// commit. The merge is GitHub's act, on a pull request whose checks have all passed; the main line here only
// fast-forwards to it. RULE 6 still holds: a forbidden verb is never free (holds.go), and this one is his
// glass's alone.

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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

// prViewFields is what a merge reads of a pull request before it asks GitHub to merge it.
const prViewFields = "number,title,state,baseRefName,headRefName,headRefOid,mergeable,url,statusCheckRollup"

// prNumber is the pull request a call names: a number as JSON brings it, or its digits as typed, with or without
// '#'; 0 when it names none.
func prNumber(args map[string]any) int {
	switch v := args["number"].(type) {
	case float64:
		if v >= 1 && v == float64(int(v)) {
			return int(v)
		}
	case string:
		if n, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSpace(v), "#")); err == nil && n >= 1 {
			return n
		}
	}
	return 0
}

// prMerge merges one pull request ON GITHUB, as a merge commit, and brings the new main line down to this machine
// (2026-10-09, WHAT'S LEFT I1's fourth step; his card: "Merge on GitHub when green", and "Merge commit"). It reads the
// pull request first and refuses one that is not open, does not go into the main line, has not passed every check or
// does not merge cleanly; then it asks GitHub to merge exactly the head it read, so a save sent in between is never
// merged unproved. Never --admin and never --auto, and the line is not deleted: closing it stays its own act.
func prMerge(t tenant.Tenant, args map[string]any) (string, error) {
	n := prNumber(args)
	if n == 0 {
		return "Refused: name the pull request to merge, by its number.", nil
	}
	num := strconv.Itoa(n)
	out, err := ghRun(t, 30*time.Second, "pr", "view", num, "--json", prViewFields)
	if err != nil {
		return fmt.Sprintf("Refused: gh could not read pull request #%d -- %s", n, firstLine(out)), nil
	}
	var pr struct {
		Title             string    `json:"title"`
		State             string    `json:"state"`
		BaseRefName       string    `json:"baseRefName"`
		HeadRefName       string    `json:"headRefName"`
		HeadRefOid        string    `json:"headRefOid"`
		Mergeable         string    `json:"mergeable"`
		StatusCheckRollup []prCheck `json:"statusCheckRollup"`
	}
	if err := json.Unmarshal([]byte(out), &pr); err != nil || pr.State == "" {
		return "Refused: gh answered in a shape this cannot read -- " + firstLine(out), nil
	}
	if state := strings.ToLower(pr.State); state != "open" {
		return fmt.Sprintf("Refused: pull request #%d is %s, not open: there is nothing to merge.", n, state), nil
	}
	mainName := mainLine(t)
	if mainName == "" {
		return "Refused: this world has no main line -- no `main` or `master` branch -- so there is nowhere here for the merge to come down to.", nil
	}
	if pr.BaseRefName != mainName {
		return fmt.Sprintf("Refused: pull request #%d goes into %q, not the main line %q. Merge on GitHub takes a line of work into the main line and nothing else.", n, pr.BaseRefName, mainName), nil
	}
	switch prChecks(pr.StatusCheckRollup) {
	case "none":
		return fmt.Sprintf("Refused: no checks have run on pull request #%d, so nothing has proved it. Merge on GitHub waits until its checks have passed.", n), nil
	case "pending":
		return fmt.Sprintf("Refused: pull request #%d's checks are still running. Merge on GitHub waits until every one has passed.", n), nil
	case "failed":
		return fmt.Sprintf("Refused: pull request #%d's checks FAILED. Mend the line, send it again, and merge when they pass.", n), nil
	}
	if strings.EqualFold(pr.Mergeable, "CONFLICTING") {
		return fmt.Sprintf("Refused: GitHub says pull request #%d conflicts with %q, so it cannot merge as it stands. Joining the two is your terminal's; send the line again and merge when its checks pass.", n, mainName), nil
	}
	if strings.TrimSpace(pr.HeadRefOid) == "" {
		return fmt.Sprintf("Refused: gh did not say which save pull request #%d stands on, so it is not merged: a merge here is of exactly the save that was proved.", n), nil
	}
	if out, err := ghRun(t, 120*time.Second, "pr", "merge", num, "--merge", "--match-head-commit", pr.HeadRefOid); err != nil {
		return fmt.Sprintf("Refused: GitHub did not merge pull request #%d -- %s", n, firstLine(out)), nil
	}
	return fmt.Sprintf("Merged pull request #%d on GitHub as a merge commit: %q, %s into %s.%s", n, pr.Title, pr.HeadRefName, mainName, mainDown(t, mainName)), nil
}

// mainDown brings the main line GitHub now holds down to this machine, as far as it can without touching unsaved work
// or joining anything: taken where he stands on a clean main line, moved where he stands on another line, and fetched
// but left for him over unsaved work. Fast-forward only, as Take from GitHub is; it says which, and where main stands.
func mainDown(t tenant.Tenant, mainName string) string {
	cur, _ := currentBranch(t)
	switch {
	case cur == mainName && dirty(t):
		if out, err := gitRun(t, 120*time.Second, "fetch", "origin", mainName); err != nil {
			return fmt.Sprintf("\n\nThe new %s was not fetched -- %s", mainName, firstLine(out))
		}
		return fmt.Sprintf("\n\nThe new %s is fetched but not taken: there is unsaved work here. Save it, then Take from GitHub.", mainName)
	case cur == mainName:
		if out, err := gitRun(t, 120*time.Second, "pull", "--ff-only", "origin", mainName); err != nil {
			return fmt.Sprintf("\n\n%s here did not move -- %s. It moves fast-forward only; if both have moved, joining them is your terminal's.", mainName, firstLine(out))
		}
	default:
		if out, err := gitRun(t, 120*time.Second, "fetch", "origin", mainName+":"+mainName); err != nil {
			return fmt.Sprintf("\n\n%s here did not move -- %s. It moves fast-forward only; if both have moved, joining them is your terminal's.", mainName, firstLine(out))
		}
	}
	short, _ := gitRun(t, 10*time.Second, "rev-parse", "--short", "refs/heads/"+mainName)
	subject, _ := gitRun(t, 10*time.Second, "log", "-1", "--format=%s", "refs/heads/"+mainName)
	return fmt.Sprintf("\n\nThe new %s is down on this machine: %s %q.", mainName, strings.TrimSpace(short), strings.TrimSpace(subject))
}

// toolGitPR is git_pr: list (a read, and the default), open or merge (his click).
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
	case "merge":
		return prMerge(t, args)
	default:
		return fmt.Sprintf("Refused: %q is not something this does. It lists the open pull requests (list), opens one from the line you stand on (open), or merges one on GitHub once every check on it has passed (merge, with its number).", action), nil
	}
}

// gitPRTool is git_pr's declaration; tools.go registers it with one line, r.add(gitPRTool()).
func gitPRTool() Tool {
	return Tool{
		Name:        "git_pr",
		Description: "pull requests, the one GitHub-side object his ruling of 2026-10-07 let into the door: list a world's open pull requests with their checks as one word (none, pending, passed, failed); open one from the line of work he stands on into the main line, refused on the main line, over unsaved work, and for a line GitHub does not have exactly as it stands; or merge one on GitHub by its number, as a merge commit, refused unless it is open, goes into the main line, merges cleanly and every check on it has passed, and then bring the new main line down to this machine (since 2026-10-09). Behind the same wall as every send (GIT_REMOTE). His glass alone may call it",
		Writes:      true,
		Reads:       []string{"list"},
		ServiceOnly: true,
		Args:        []string{"action?", "project?", "number?"},
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
