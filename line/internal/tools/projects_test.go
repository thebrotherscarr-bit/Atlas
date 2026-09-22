// Strokes for `projects`: the maker's projects, read for the glass (the
// maker's piece 2, 2026-09-21).
//
// HERMETIC BY LAW 5: every stroke builds its own ground in t.TempDir() -- and
// the ground is itself a git repository, as the operator's is, because the one
// fault this tool must never have is git answering for the GROUND when asked
// about a folder with no history of its own. A stroke on a bare temp dir could
// not tell the two apart.
package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"atlas/line/internal/tenant"
)

var makerEpoch = time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)

// makerGround is a ground that is a repository of its own, as his is.
func makerGround(t *testing.T) tenant.Tenant {
	t.Helper()
	home := t.TempDir()
	runGit(t, home, nil, "init", "-b", "main")
	runGit(t, home, nil, "config", "user.name", "prove")
	runGit(t, home, nil, "config", "user.email", "prove@localhost")
	write(t, home, "ground.txt", "the ground's own file\n")
	runGit(t, home, nil, "add", "-A")
	runGit(t, home, nil, "commit", "-m", "the ground's own save")
	return tenant.Tenant{Name: "probe", Home: home, Manifest: tenant.DefaultManifest()}
}

func runGit(t *testing.T, dir string, env []string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(append(os.Environ(), "GIT_TERMINAL_PROMPT=0"), env...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// makerProject builds projects/<name> the way the maker does: a repository of
// its own, one page, and one commit per version titled "Version N: note".
func makerProject(t *testing.T, home, name string, at time.Time, pages ...string) string {
	t.Helper()
	dir := filepath.Join(home, "projects", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, nil, "init", "-b", "main")
	runGit(t, dir, nil, "config", "user.name", "manjuel (the maker)")
	runGit(t, dir, nil, "config", "user.email", "maker@localhost")
	for i, page := range pages {
		write(t, dir, "index.html", page)
		stamp := at.Add(time.Duration(i) * time.Minute).Format(time.RFC3339)
		env := []string{"GIT_AUTHOR_DATE=" + stamp, "GIT_COMMITTER_DATE=" + stamp}
		runGit(t, dir, env, "add", "-A")
		runGit(t, dir, env, "commit", "-m", fmt.Sprintf("Version %d: note %d", i+1, i+1))
	}
	return dir
}

func pageOf(body string) string {
	return "<!DOCTYPE html>\n<html><body>" + body + "</body></html>\n"
}

type listed struct {
	Count    int `json:"count"`
	Projects []struct {
		Name     string `json:"name"`
		Versions []struct {
			N    int    `json:"n"`
			Sha  string `json:"sha"`
			Note string `json:"note"`
			When string `json:"when"`
		} `json:"versions"`
		Bytes   int64  `json:"bytes"`
		Lines   int    `json:"lines"`
		Updated string `json:"updated"`
	} `json:"projects"`
}

func listOf(t *testing.T, tn tenant.Tenant) listed {
	t.Helper()
	out, err := toolProjects(tn, map[string]any{"action": "list"})
	if err != nil {
		t.Fatalf("the list errored: %v", err)
	}
	var l listed
	if err := json.Unmarshal([]byte(out), &l); err != nil {
		t.Fatalf("the list is not JSON: %v\n%s", err, out)
	}
	return l
}

type served struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
	Sha     string `json:"sha"`
	Bytes   int    `json:"bytes"`
	Text    string `json:"text"`
}

func pageFor(t *testing.T, tn tenant.Tenant, args map[string]any) served {
	t.Helper()
	args["action"] = "page"
	out, err := toolProjects(tn, args)
	if err != nil {
		t.Fatalf("the page errored: %v", err)
	}
	var s served
	if err := json.Unmarshal([]byte(out), &s); err != nil {
		t.Fatalf("the page is not JSON: %v\n%s", err, out)
	}
	return s
}

// THE LIST IS EACH PROJECT'S OWN HISTORY, and nothing that is not a project.
func TestTheProjectsAreReadFromTheirOwnHistory(t *testing.T) {
	tn := makerGround(t)
	// NEWEST FIRST IS NOT ALPHABETICAL HERE ON PURPOSE: the calculator was
	// touched last, so a list that merely came back in folder order would
	// put snake-game first and be caught.
	makerProject(t, tn.Home, "snake-game", makerEpoch,
		pageOf("<p>one</p>"), pageOf("<p>two</p>"))
	makerProject(t, tn.Home, "tip-calculator", makerEpoch.Add(time.Hour), pageOf("<p>calc</p>"))
	// Not projects: a folder with no history of its own, a stray file, and a
	// folder whose name the maker would never give one.
	write(t, tn.Home, "projects/plain/index.html", pageOf("<p>plain</p>"))
	write(t, tn.Home, "projects/stray.txt", "not a project\n")
	makerProject(t, tn.Home, "Bad_Name", makerEpoch, pageOf("<p>bad</p>"))

	l := listOf(t, tn)
	names := []string{}
	for _, p := range l.Projects {
		names = append(names, p.Name)
	}
	if strings.Join(names, ",") != "tip-calculator,snake-game" || l.Count != 2 {
		t.Fatalf("the list is %v (count %d); wanted the two projects, newest first", names, l.Count)
	}
	snake := l.Projects[1]
	if len(snake.Versions) != 2 || snake.Versions[0].N != 1 || snake.Versions[1].N != 2 {
		t.Fatalf("snake-game's versions read as %+v", snake.Versions)
	}
	if snake.Versions[0].Note != "note 1" || snake.Versions[1].Note != "note 2" {
		t.Fatalf("the maker's notes are not read off its subjects: %+v", snake.Versions)
	}
	if snake.Versions[1].Sha == "" || snake.Versions[1].When == "" {
		t.Fatalf("a version comes back without its sha or its time: %+v", snake.Versions[1])
	}
	if snake.Bytes == 0 || snake.Lines == 0 || snake.Updated == "" {
		t.Fatalf("the page's size and time are missing: %+v", snake)
	}
	if !strings.HasPrefix(snake.Updated, "2026-09-21T12:01") {
		t.Fatalf("updated is not the last version's time: %q", snake.Updated)
	}
}

// A GROUND WITH NO projects/ HAS NOT FAILED AT ANYTHING.
func TestAWorldWithNoProjectsListsNone(t *testing.T) {
	tn := world(t, map[string]string{"CLAUDE.md": "x"})
	l := listOf(t, tn)
	if l.Count != 0 || len(l.Projects) != 0 {
		t.Fatalf("a ground with no projects listed %+v", l)
	}
}

// A PAGE AS IT STANDS, OR AS ANY VERSION WAS.
func TestAPageIsServedAsItStandsOrAsAVersionWas(t *testing.T) {
	tn := makerGround(t)
	dir := makerProject(t, tn.Home, "snake-game", makerEpoch,
		pageOf("<p>one</p>"), pageOf("<p>two</p>"))

	now := pageFor(t, tn, map[string]any{"name": "snake-game"})
	disk, _ := os.ReadFile(filepath.Join(dir, "index.html"))
	if now.Text != string(disk) || now.Version != 2 || now.Sha == "" || now.Bytes != len(disk) {
		t.Fatalf("the page now is not the file on disk: %+v", now)
	}
	for _, v := range []any{float64(1), "1", 1} {
		first := pageFor(t, tn, map[string]any{"name": "snake-game", "version": v})
		if !strings.Contains(first.Text, "<p>one</p>") || first.Version != 1 {
			t.Fatalf("version 1 (asked as %T) came back as %+v", v, first)
		}
	}
	for _, v := range []any{float64(3), float64(0.5), "x", float64(-1)} {
		if _, err := toolProjects(tn, map[string]any{"action": "page", "name": "snake-game", "version": v}); err == nil {
			t.Fatalf("version %v was served; it is not one of this project's", v)
		}
	}
}

// THE TEXT CROSSES AS IT IS. A page is mostly <, > and &; escaped, each would
// arrive at the glass six characters long.
func TestAPageIsNotEscapedOnTheWay(t *testing.T) {
	tn := makerGround(t)
	makerProject(t, tn.Home, "snake-game", makerEpoch, pageOf("<p>a & b</p>"))
	out, err := toolProjects(tn, map[string]any{"action": "page", "name": "snake-game"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<p>a & b</p>") || strings.Contains(out, `\u003c`) {
		t.Fatalf("the page was escaped on its way out:\n%s", out)
	}
}

// A NAME IS A NAME, NOT A PATH -- and a folder with no history of its own is
// never read, because git asked about it would answer for the ground.
func TestAProjectNameIsANameNotAPath(t *testing.T) {
	tn := makerGround(t)
	makerProject(t, tn.Home, "snake-game", makerEpoch, pageOf("<p>one</p>"))
	write(t, tn.Home, "projects/plain/index.html", pageOf("<p>plain</p>"))
	for _, name := range []string{
		"../ground.txt", "..", "snake-game/../../ground.txt", `C:\keys.txt`,
		"/etc/passwd", "SNAKE-GAME", "", "snake game", "plain", "nope", "-flag",
	} {
		out, err := toolProjects(tn, map[string]any{"action": "page", "name": name})
		if err == nil {
			t.Fatalf("%q was served:\n%.200s", name, out)
		}
	}
	_, err := toolProjects(tn, map[string]any{"action": "page", "name": "plain"})
	if err == nil || !strings.Contains(err.Error(), "no history of its own") {
		t.Fatalf("a folder with no history must be refused for that reason: %v", err)
	}
}

// A PROJECT LINKED IN FROM OUTSIDE projects/ IS NOT ONE.
//
// A JUNCTION, WHERE A SYMLINK NEEDS PRIVILEGE. On Windows a symlink needs a
// privilege this machine's user does not hold, and a junction needs none --
// so the junction is the link anyone here can actually make, and the one the
// first cut of this tool let through: filepath.EvalSymlinks does not follow
// it, so it read as inside projects/, and its history was found through it.
func TestAProjectLinkedFromOutsideIsRefused(t *testing.T) {
	tn := makerGround(t)
	makerProject(t, tn.Home, "snake-game", makerEpoch, pageOf("<p>one</p>"))
	outside := t.TempDir()
	target := makerProject(t, outside, "elsewhere", makerEpoch, pageOf("<p>secret</p>"))
	link := filepath.Join(tn.Home, "projects", "linked")
	if err := os.Symlink(target, link); err != nil {
		if runtime.GOOS != "windows" {
			t.Skipf("this machine will not make a link here (%v)", err)
		}
		out, jerr := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput()
		if jerr != nil {
			t.Skipf("no link of either kind could be made here: %v; %v %s", err, jerr, out)
		}
	}
	out, err := toolProjects(tn, map[string]any{"action": "page", "name": "linked"})
	if err == nil {
		t.Fatalf("a project linked from outside projects/ was served:\n%.200s", out)
	}
	for _, p := range listOf(t, tn).Projects {
		if p.Name == "linked" {
			t.Fatal("a project linked from outside projects/ was listed")
		}
	}
	if s := pageFor(t, tn, map[string]any{"name": "snake-game"}); !strings.Contains(s.Text, "<p>one</p>") {
		t.Fatalf("the project beside the link must still be served: %+v", s)
	}
}

func TestAnUnknownProjectsActionIsRefusedByName(t *testing.T) {
	tn := makerGround(t)
	_, err := toolProjects(tn, map[string]any{"action": "burn"})
	if err == nil || !strings.Contains(err.Error(), "list or page") {
		t.Fatalf("an unknown action must be refused by name: %v", err)
	}
}
