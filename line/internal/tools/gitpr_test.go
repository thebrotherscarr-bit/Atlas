package tools

// PULL REQUESTS, held by strokes (gitpr.go): glass-only, behind the wall, checks read as one word, and open refused where it must be. gh is stood in for, so no stroke reaches GitHub.

import (
	"errors"
	"os/exec"
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
	if !strings.Contains(out, "never merges") {
		t.Fatalf("Expected 'never merges' in output, got: %v", out)
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
