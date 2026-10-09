package tools

// PULL REQUESTS, held by strokes (gitpr.go): glass-only, behind the wall, checks read as one word, open refused where it must be, and a merge made on GitHub only when every check has passed, with the new main brought down here. gh is stood in for, so no stroke reaches GitHub, and GitHub's own merge is stood in for by a merge commit made on an origin on this disk.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"atlas/line/internal/tenant"
)

// ghStub stands in for gh for one stroke: every call is recorded, and each answers what it is given.
func ghStub(t *testing.T, answer string, fail error) *[][]string {
	t.Helper()
	oldGhRun := ghRun
	t.Cleanup(func() {
		ghRun = oldGhRun
	})
	calls := &[][]string{}
	ghRun = func(tn tenant.Tenant, timeout time.Duration, args ...string) (string, error) {
		*calls = append(*calls, args)
		return answer, fail
	}
	return calls
}

// THE OPEN PULL REQUESTS ARE LISTED with their line, where they go and their checks as one word, and a gh that fails says so.
func TestPullRequestsAreListedWithTheirChecks(t *testing.T) {
	t.Setenv("MANJUEL_GIT_REMOTE", "1")
	tn := tempWorld(t)

	json := `[
		{
			"number": 7,
			"title": "work on the spur",
			"headRefName": "spur",
			"baseRefName": "main",
			"url": "https://github.com/owner/repo/pull/7",
			"statusCheckRollup": [
				{
					"status": "COMPLETED",
					"conclusion": "SUCCESS"
				}
			]
		}
	]`
	calls := ghStub(t, json, nil)

	for _, action := range []string{"", "list"} {
		out := call(t, toolGitPR, tn, glassArgs(action))
		if !strings.Contains(out, `"number": 7`) ||
			!strings.Contains(out, `"line": "spur"`) ||
			!strings.Contains(out, `"into": "main"`) ||
			!strings.Contains(out, `"checks": "passed"`) {
			t.Fatalf("toolGitPR(%v) = %v; want to contain 'number': 7, 'line': 'spur', 'into': 'main', 'checks': 'passed'", action, out)
		}
	}

	if len(*calls) != 2 || !strings.HasPrefix((*calls)[0][0], "pr") || !strings.HasPrefix((*calls)[0][1], "list") {
		t.Fatalf("gh was asked incorrectly: %v", *calls)
	}

	ghStub(t, "gh: not signed in", errors.New("exit status 1"))
	out := call(t, toolGitPR, tn, glassArgs("list"))
	if !strings.Contains(out, "gh could not list") {
		t.Fatalf("toolGitPR(list) = %v; want to contain 'gh could not list'", out)
	}
}

// sentLine is a world standing on the line "spur", saved and sent to an origin on this machine.
func sentLine(t *testing.T) tenant.Tenant {
	t.Helper()
	tn := tempWorld(t)
	origin := t.TempDir()
	if err := exec.Command("git", "init", "--bare", origin).Run(); err != nil {
		t.Fatalf("git init --bare failed: %v", err)
	}
	if _, err := gitRun(tn, 10*time.Second, "remote", "add", "origin", origin); err != nil {
		t.Fatalf("git remote add failed: %v", err)
	}
	call(t, toolGitBranch, tn, map[string]any{"action": "new", "name": "spur"})
	write(t, tn.Home, "spur.txt", "work on the spur\n")
	call(t, toolGitCommit, tn, map[string]any{"message": "work on the spur"})
	if _, err := gitRun(tn, 30*time.Second, "push", "--set-upstream", "origin", "spur"); err != nil {
		t.Fatalf("git push failed: %v", err)
	}
	return tn
}

// glassArgs is one call's arguments as his glass sends them.
func glassArgs(action string) map[string]any {
	return map[string]any{"action": action, CallerKey: glass}
}

// A PULL REQUEST'S CHECKS READ AS ONE WORD, and a failure outranks anything still running.
func TestAPullRequestsChecksReadAsOneWord(t *testing.T) {
	tests := []struct {
		checks []prCheck
		want   string
	}{
		{nil, "none"},
		{[]prCheck{{Status: "COMPLETED", Conclusion: "SUCCESS"}}, "passed"},
		{[]prCheck{{Status: "IN_PROGRESS"}}, "pending"},
		{[]prCheck{{State: "PENDING"}}, "pending"},
		{[]prCheck{{State: "SUCCESS"}, {Status: "COMPLETED", Conclusion: "NEUTRAL"}}, "passed"},
		{[]prCheck{{Status: "IN_PROGRESS"}, {Status: "COMPLETED", Conclusion: "FAILURE"}}, "failed"},
		{[]prCheck{{State: "ERROR"}}, "failed"},
	}

	for _, tt := range tests {
		got := prChecks(tt.checks)
		if got != tt.want {
			t.Errorf("prChecks(%v) = %v; want %v", tt.checks, got, tt.want)
		}
	}
}

// PULL REQUESTS ARE HIS GLASS'S ALONE AND BEHIND THE WALL, and nothing refused ever asks gh.
func TestPullRequestsAreHisGlassAloneAndBehindTheWall(t *testing.T) {
	tn := tempWorld(t)
	calls := ghStub(t, "[]", nil)

	out, err := toolGitPR(tn, map[string]any{"action": "list", CallerKey: agent})
	if err == nil || out != "" || !strings.Contains(err.Error(), "operator's own hand") {
		t.Fatalf("toolGitPR should have failed with 'operator's own hand' error")
	}

	t.Setenv("MANJUEL_GIT_REMOTE", "")
	t.Setenv("CHAINKIT_GIT_REMOTE", "")
	out = call(t, toolGitPR, tn, glassArgs("list"))
	if !strings.Contains(out, "is OFF") {
		t.Fatalf("Expected 'is OFF' in output, got: %v", out)
	}

	bare := tenant.Tenant{Name: "bare", Home: t.TempDir(), Manifest: tenant.DefaultManifest()}
	out = call(t, toolGitPR, bare, glassArgs("list"))
	if !strings.Contains(out, "not a repository") {
		t.Fatalf("Expected 'not a repository' in output, got: %v", out)
	}

	if len(*calls) != 0 {
		t.Fatalf("gh was asked although every call was refused")
	}
}

// OPENING IS REFUSED WHERE IT MUST BE: on the main line, for a line GitHub does not have, over unsaved work; and nothing merges. None of it asks gh.
func TestOpeningAPullRequestIsRefusedWhereItMustBe(t *testing.T) {
	t.Setenv("MANJUEL_GIT_REMOTE", "1")
	tn := tempWorld(t)
	calls := ghStub(t, "", nil)

	out := call(t, toolGitPR, tn, glassArgs("open"))
	if !strings.Contains(out, "the main line") {
		t.Fatalf("Expected 'the main line' in output, got: %v", out)
	}

	call(t, toolGitBranch, tn, map[string]any{"action": "new", "name": "spur"})
	out = call(t, toolGitPR, tn, glassArgs("open"))
	if !strings.Contains(out, "Send it first") {
		t.Fatalf("Expected 'Send it first' in output, got: %v", out)
	}

	write(t, tn.Home, "loose.txt", "unsaved work\n")
	out = call(t, toolGitPR, tn, glassArgs("open"))
	if !strings.Contains(out, "unsaved work") {
		t.Fatalf("Expected 'unsaved work' in output, got: %v", out)
	}

	out = call(t, toolGitPR, tn, glassArgs("merge"))
	if !strings.Contains(out, "name the pull request") {
		t.Fatalf("Expected a merge naming no pull request to be refused, got: %v", out)
	}

	if len(*calls) != 0 {
		t.Fatalf("gh was asked by a refused call")
	}
}

// A SENT LINE OPENS ONE INTO THE MAIN LINE: titled with its newest save, its saves listed, and no name but his on it.
func TestASentLineOpensAPullRequestIntoTheMainLine(t *testing.T) {
	t.Setenv("MANJUEL_GIT_REMOTE", "1")
	tn := sentLine(t)
	calls := ghStub(t, "https://github.com/owner/repo/pull/8", nil)

	out := call(t, toolGitPR, tn, glassArgs("open"))
	if !strings.Contains(out, "Opened a pull request") || !strings.Contains(out, "https://github.com/owner/repo/pull/8") {
		t.Fatalf("Expected 'Opened a pull request' and 'https://github.com/owner/repo/pull/8' in output, got: %v", out)
	}

	if len(*calls) != 1 {
		t.Fatalf("Expected exactly one call to gh, got: %v", len(*calls))
	}
	args := strings.Join((*calls)[0], "|")
	if !strings.Contains(args, "pr|create|--base|main|--head|spur|--title|work on the spur|--body|") || !strings.Contains(args, "- work on the spur") {
		t.Fatalf("Expected args to contain 'pr|create|--base|main|--head|spur|--title|work on the spur|--body|' and '- work on the spur', got: %v", args)
	}

	if strings.Contains(args, "Co-Authored") || strings.Contains(args, "Generated with") {
		t.Fatalf("A pull request carries no name but his")
	}
}

// ghPRStub stands in for gh for a merge: `pr view` answers what it is given, and `pr merge` does GitHub's own act, onMerge, when one is given. Every call is recorded.
func ghPRStub(t *testing.T, view string, onMerge func() error) *[][]string {
	t.Helper()
	oldGhRun := ghRun
	t.Cleanup(func() {
		ghRun = oldGhRun
	})
	calls := &[][]string{}
	ghRun = func(tn tenant.Tenant, timeout time.Duration, args ...string) (string, error) {
		*calls = append(*calls, args)
		if len(args) > 1 && args[0] == "pr" && args[1] == "merge" {
			if onMerge != nil {
				if err := onMerge(); err != nil {
					return "GitHub refused the merge: " + err.Error(), err
				}
			}
			return "Merged pull request #9 (work on the spur)", nil
		}
		return view, nil
	}
	return calls
}

// The checks a pull request can carry, as gh reports them.
const (
	checksPassed  = `[{"status": "COMPLETED", "conclusion": "SUCCESS"}]`
	checksRunning = `[{"status": "IN_PROGRESS"}]`
	checksFailed  = `[{"status": "COMPLETED", "conclusion": "FAILURE"}]`
	checksNone    = `[]`
)

// prView is gh's answer to `pr view` for pull request 9 from spur: its state, where it goes, the head it stands on, its checks and whether it merges cleanly.
func prView(state, into, head, checks, mergeable string) string {
	return fmt.Sprintf(`{"number": 9, "title": "work on the spur", "state": %q, "baseRefName": %q, "headRefName": "spur", "headRefOid": %q, "mergeable": %q, "url": "https://github.com/owner/repo/pull/9", "statusCheckRollup": %s}`,
		state, into, head, mergeable, checks)
}

// mergeArgs is a merge as his glass sends it, the number as JSON brings it.
func mergeArgs(number any) map[string]any {
	args := glassArgs("merge")
	args["number"] = number
	return args
}

// A PULL REQUEST MERGES ONLY WHEN IT IS NAMED, OPEN, INTO THE MAIN LINE, GREEN AND CLEAN: anything else is refused by name, and nothing refused asks gh to merge.
func TestAPullRequestMergesOnlyWhenEveryCheckHasPassed(t *testing.T) {
	t.Setenv("MANJUEL_GIT_REMOTE", "1")
	tn := tempWorld(t)

	calls := ghPRStub(t, prView("OPEN", "main", "abc123", checksPassed, "MERGEABLE"), nil)
	for _, number := range []any{nil, "", "nine", float64(0), float64(-3), float64(9.5)} {
		out := call(t, toolGitPR, tn, mergeArgs(number))
		if !strings.Contains(out, "name the pull request") {
			t.Fatalf("a merge naming %v = %v; want it refused: 'name the pull request'", number, out)
		}
	}
	if len(*calls) != 0 {
		t.Fatalf("gh was asked though no pull request was named: %v", *calls)
	}

	for _, tc := range []struct{ view, want string }{
		{prView("MERGED", "main", "abc123", checksPassed, "UNKNOWN"), "is merged, not open"},
		{prView("CLOSED", "main", "abc123", checksPassed, "MERGEABLE"), "is closed, not open"},
		{prView("OPEN", "spur-two", "abc123", checksPassed, "MERGEABLE"), "not the main line"},
		{prView("OPEN", "main", "abc123", checksNone, "MERGEABLE"), "no checks have run"},
		{prView("OPEN", "main", "abc123", checksRunning, "MERGEABLE"), "still running"},
		{prView("OPEN", "main", "abc123", checksFailed, "MERGEABLE"), "FAILED"},
		{prView("OPEN", "main", "abc123", checksPassed, "CONFLICTING"), "conflicts with"},
		{prView("OPEN", "main", "", checksPassed, "MERGEABLE"), "which save"},
		{"not the shape gh answers in", "shape this cannot read"},
	} {
		calls := ghPRStub(t, tc.view, nil)
		out := call(t, toolGitPR, tn, mergeArgs("9"))
		if !strings.HasPrefix(out, "Refused") || !strings.Contains(out, tc.want) {
			t.Fatalf("merging %s = %v; want it refused: %q", tc.view, out, tc.want)
		}
		if len(*calls) != 1 || !strings.HasPrefix(strings.Join((*calls)[0], "|"), "pr|view|9|--json|") {
			t.Fatalf("gh was asked for something other than the pull request itself: %v", *calls)
		}
	}

	ghStub(t, "no pull requests found", errors.New("exit status 1"))
	out := call(t, toolGitPR, tn, mergeArgs(float64(9)))
	if !strings.Contains(out, "gh could not read pull request #9") {
		t.Fatalf("a merge with gh failing = %v; want 'gh could not read pull request #9'", out)
	}
}

// githubMerges stands in for GitHub's own act on pull request 9: a merge commit on the origin's main, whose parents are main and spur.
func githubMerges(origin string) func() error {
	return func() error {
		git := func(args ...string) (string, error) {
			cmd := exec.Command("git", append([]string{"-C", origin}, args...)...)
			cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=github", "GIT_AUTHOR_EMAIL=github@localhost",
				"GIT_COMMITTER_NAME=github", "GIT_COMMITTER_EMAIL=github@localhost")
			out, err := cmd.Output()
			return strings.TrimSpace(string(out)), err
		}
		tree, err := git("rev-parse", "spur^{tree}")
		if err != nil {
			return err
		}
		merge, err := git("commit-tree", tree, "-p", "main", "-p", "spur", "-m", "Merge pull request #9 from owner/spur")
		if err != nil {
			return err
		}
		_, err = git("update-ref", "refs/heads/main", merge)
		return err
	}
}

// A GREEN PULL REQUEST MERGES ON GITHUB AS A MERGE COMMIT, AND THE NEW MAIN COMES DOWN: taken where he stands on the line or on a clean main line, fetched and left for him over unsaved work. gh is asked to merge exactly the head it read, and nothing more.
func TestAGreenPullRequestMergesOnGitHubAndTheNewMainComesDown(t *testing.T) {
	t.Setenv("MANJUEL_GIT_REMOTE", "1")
	for _, stand := range []string{"on the line", "on the main line", "on the main line over unsaved work"} {
		t.Run(stand, func(t *testing.T) {
			tn := sentLine(t)
			if _, err := gitRun(tn, 30*time.Second, "push", "origin", "main"); err != nil {
				t.Fatalf("git push main failed: %v", err)
			}
			origin, _ := gitRun(tn, 10*time.Second, "remote", "get-url", "origin")
			head, _ := gitRun(tn, 10*time.Second, "rev-parse", "spur")
			before, _ := gitRun(tn, 10*time.Second, "rev-parse", "main")
			if stand != "on the line" {
				call(t, toolGitBranch, tn, map[string]any{"action": "switch", "name": "main"})
			}
			if stand == "on the main line over unsaved work" {
				write(t, tn.Home, "loose.txt", "unsaved work\n")
			}
			calls := ghPRStub(t, prView("OPEN", "main", head, checksPassed, "MERGEABLE"), githubMerges(origin))

			out := call(t, toolGitPR, tn, mergeArgs(float64(9)))
			if !strings.Contains(out, "Merged pull request #9 on GitHub") {
				t.Fatalf("merge = %v; want 'Merged pull request #9 on GitHub'", out)
			}
			if len(*calls) != 2 || !strings.HasPrefix(strings.Join((*calls)[0], "|"), "pr|view|9|--json|") {
				t.Fatalf("gh was not asked to read the pull request and then merge it: %v", *calls)
			}
			if got := strings.Join((*calls)[1], "|"); got != "pr|merge|9|--merge|--match-head-commit|"+head {
				t.Fatalf("gh was asked to merge as %q; want a merge commit of exactly the head it read, and nothing more", got)
			}

			merged, err := exec.Command("git", "-C", origin, "rev-parse", "main").Output()
			if err != nil {
				t.Fatalf("the origin's main could not be read: %v", err)
			}
			want := strings.TrimSpace(string(merged))
			local, _ := gitRun(tn, 10*time.Second, "rev-parse", "main")
			if stand == "on the main line over unsaved work" {
				tracking, _ := gitRun(tn, 10*time.Second, "rev-parse", "origin/main")
				if local != before || tracking != want || !strings.Contains(out, "fetched but not taken") {
					t.Fatalf("over unsaved work, main here is %v (was %v) and GitHub's is fetched as %v (want %v): %v", local, before, tracking, want, out)
				}
				return
			}
			if local != want || !strings.Contains(out, "is down on this machine") {
				t.Fatalf("main here is %v; want GitHub's new main %v: %v", local, want, out)
			}
			if stand == "on the main line" {
				if _, err := os.Stat(filepath.Join(tn.Home, "spur.txt")); err != nil {
					t.Fatalf("the merged work is not on the disk: %v", err)
				}
			}
		})
	}
}
