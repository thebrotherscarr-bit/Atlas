package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Strokes for the Update button's three checkouts (2026-10-23, his word: "build the glass button"). A ground is laid out
// as his is -- Manjuel's checkout, Atlas inside it, the notes kernel at worlds/jesster -- each cloned from a bare
// "GitHub" in the test's own directory, so nothing here reaches a network.

func run(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=Kyle", "-c", "user.email=k@example.invalid",
		"-c", "init.defaultBranch=main", "-c", "advice.detachedHead=false"}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
	return strings.TrimSpace(string(out))
}

// repo makes a bare origin with one save on branch, and a second clone to push new saves from.
func repo(t *testing.T, root, name, branch string) (origin, writer string) {
	t.Helper()
	origin = filepath.Join(root, name+".git")
	writer = filepath.Join(root, name+"-writer")
	run(t, root, "init", "-q", "--bare", "-b", branch, origin)
	run(t, root, "clone", "-q", origin, writer)
	run(t, writer, "checkout", "-q", "-B", branch)
	os.WriteFile(filepath.Join(writer, "a.txt"), []byte("one\n"), 0o644)
	run(t, writer, "add", "a.txt")
	run(t, writer, "commit", "-q", "-m", "first")
	run(t, writer, "push", "-q", "-u", "origin", branch)
	return origin, writer
}

func push(t *testing.T, writer, file, subject string) {
	t.Helper()
	os.WriteFile(filepath.Join(writer, file), []byte(subject+"\n"), 0o644)
	run(t, writer, "add", file)
	run(t, writer, "commit", "-q", "-m", subject)
	run(t, writer, "push", "-q")
}

type ground struct {
	atlas, manjuel, kernel    string
	atlasW, manjuelW, kernelW string
}

func newGround(t *testing.T) ground {
	t.Helper()
	root := t.TempDir()
	var g ground
	mo, mw := repo(t, root, "manjuel", "main")
	ao, aw := repo(t, root, "atlas", "main")
	ko, kw := repo(t, root, "jesster", "secondbrain-app")
	g.manjuel = filepath.Join(root, "Research")
	run(t, root, "clone", "-q", mo, g.manjuel)
	g.atlas = filepath.Join(g.manjuel, "atlas")
	run(t, root, "clone", "-q", ao, g.atlas)
	g.kernel = filepath.Join(g.manjuel, "worlds", "jesster")
	run(t, root, "clone", "-q", "-b", "secondbrain-app", ko, g.kernel)
	g.atlasW, g.manjuelW, g.kernelW = aw, mw, kw
	return g
}

func TestTheUpdateKeepsThreeCheckouts(t *testing.T) {
	cs := checkouts(filepath.Join("R", "atlas"))
	want := []checkout{{"Atlas", filepath.Join("R", "atlas"), "main"}, {"Manjuel", "R", "main"}, {"notes kernel", filepath.Join("R", "worlds", "jesster"), ""}}
	if len(cs) != len(want) {
		t.Fatalf("the Update button keeps %d checkouts, not Atlas, Manjuel and the notes kernel: %+v", len(cs), cs)
	}
	for i := range want {
		if cs[i] != want[i] {
			t.Fatalf("checkout %d is %+v, wanted %+v", i, cs[i], want[i])
		}
	}
}

func TestEachCheckoutSaysHowFarBehindAndWhyItWillNotMove(t *testing.T) {
	g := newGround(t)
	ctx := context.Background()
	cs := checkouts(g.atlas)
	push(t, g.manjuelW, "b.txt", "manjuel two")
	push(t, g.kernelW, "b.txt", "kernel two")
	push(t, g.kernelW, "c.txt", "kernel three")

	if s := state(ctx, cs[1]); s.Skip != "" || s.Behind != 1 || s.Changes[0] != "manjuel two" {
		t.Fatalf("Manjuel should be one save behind and free to move: %+v", s)
	}
	if s := state(ctx, cs[2]); s.Skip != "" || s.Behind != 2 || s.Upstream != "origin/secondbrain-app" {
		t.Fatalf("the kernel should follow the branch it tracks, two behind: %+v", s)
	}

	// A change git would not overwrite does not hold Manjuel back; on Atlas it does.
	os.WriteFile(filepath.Join(g.manjuel, "a.txt"), []byte("stamp\n"), 0o644)
	if s := state(ctx, cs[1]); s.Skip != "" || s.Unsaved != "a.txt" {
		t.Fatalf("Manjuel with a changed file should name it and still be free to move: %+v", s)
	}
	os.WriteFile(filepath.Join(g.atlas, "a.txt"), []byte("mid-edit\n"), 0o644)
	if s := state(ctx, cs[0]); !strings.Contains(s.Skip, "unsaved work (a.txt)") {
		t.Fatalf("Atlas with unsaved work must be skipped, saying which file: %+v", s)
	}
	run(t, g.atlas, "checkout", "-q", "--", "a.txt")

	run(t, g.atlas, "checkout", "-q", "-b", "side")
	if s := state(ctx, cs[0]); !strings.Contains(s.Skip, `on "side", not main`) {
		t.Fatalf("Atlas off main must be skipped: %+v", s)
	}
	run(t, g.atlas, "checkout", "-q", "main")

	lock := filepath.Join(g.atlas, ".git", "index.lock")
	os.WriteFile(lock, nil, 0o644)
	if s := state(ctx, cs[0]); !strings.Contains(s.Skip, "lock file") {
		t.Fatalf("a checkout with git's lock file must be skipped: %+v", s)
	}
	os.Remove(lock)

	os.WriteFile(filepath.Join(g.kernel, "mine.txt"), []byte("local\n"), 0o644)
	run(t, g.kernel, "add", "mine.txt")
	run(t, g.kernel, "-c", "user.name=Kyle", "commit", "-q", "-m", "a local save")
	if s := state(ctx, cs[2]); !strings.Contains(s.Skip, "1 save(s) GitHub does not") {
		t.Fatalf("a kernel that has diverged must be skipped, not joined: %+v", s)
	}

	run(t, g.kernel, "branch", "-q", "--unset-upstream")
	if s := state(ctx, cs[2]); !strings.Contains(s.Skip, "follows no branch on GitHub") {
		t.Fatalf("a kernel that tracks nothing must be skipped: %+v", s)
	}

	missing := checkout{Name: "gone", Dir: filepath.Join(g.manjuel, "nowhere")}
	if s := state(ctx, missing); !strings.Contains(s.Skip, "not on this machine") {
		t.Fatalf("a checkout that is not there must say so: %+v", s)
	}
}

func TestApplyMovesWhatItCanAndRefusesUnderASitting(t *testing.T) {
	g := newGround(t)
	was := atlasDir
	atlasDir = func() string { return g.atlas }
	defer func() { atlasDir = was }()
	push(t, g.manjuelW, "b.txt", "manjuel two")
	push(t, g.kernelW, "b.txt", "kernel two")
	os.WriteFile(filepath.Join(g.manjuel, "a.txt"), []byte("stamp\n"), 0o644) // a suite's stamp: not touched by the save

	ledger := filepath.Join(g.manjuel, "sessions", "sessions.jsonl")
	os.MkdirAll(filepath.Dir(ledger), 0o755)
	os.WriteFile(ledger, []byte(`{"n": 7, "ended": "2026-10-22T10:00"}`+"\n"+`{"n": 8}`+"\n"), 0o644)
	rec := httptest.NewRecorder()
	(&Handlers{}).UpdateApply(rec, httptest.NewRequest("POST", "/api/update/apply", nil))
	if rec.Code != 409 || !strings.Contains(rec.Body.String(), "sitting 8 is open") {
		t.Fatalf("an open sitting must refuse the whole update: %d %s", rec.Code, rec.Body.String())
	}
	if head := run(t, g.manjuel, "log", "-1", "--format=%s"); head != "first" {
		t.Fatalf("Manjuel moved under an open sitting: %s", head)
	}

	os.WriteFile(ledger, []byte(`{"n": 8, "ended": "2026-10-23T09:00"}`+"\n"), 0o644)
	rec = httptest.NewRecorder()
	(&Handlers{}).UpdateApply(rec, httptest.NewRequest("POST", "/api/update/apply", nil))
	var got struct {
		Repos   []map[string]string `json:"repos"`
		Restart []string            `json:"restart"`
		Atlas   bool                `json:"atlas"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || rec.Code != 200 {
		t.Fatalf("apply did not answer: %d %s", rec.Code, rec.Body.String())
	}
	if got.Atlas || got.Repos[0]["result"] != "up to date" {
		t.Fatalf("an Atlas with nothing to take must not be rebuilt: %+v", got)
	}
	if !strings.HasPrefix(got.Repos[1]["result"], "moved") || !strings.HasPrefix(got.Repos[2]["result"], "moved") {
		t.Fatalf("Manjuel and the kernel should both have moved: %+v", got.Repos)
	}
	if run(t, g.manjuel, "log", "-1", "--format=%s") != "manjuel two" || run(t, g.kernel, "log", "-1", "--format=%s") != "kernel two" {
		t.Fatal("the checkouts did not take GitHub's saves")
	}
	if b, _ := os.ReadFile(filepath.Join(g.manjuel, "a.txt")); string(b) != "stamp\n" {
		t.Fatal("the fast-forward touched a changed file it should have left alone")
	}
	if len(got.Restart) != 2 || !strings.Contains(got.Restart[0], "Manjuel moved: restart required") {
		t.Fatalf("Manjuel moving must say restart required: %+v", got.Restart)
	}
	if log, _ := os.ReadFile(filepath.Join(g.atlas, "update.log")); !strings.Contains(string(log), "Manjuel: moved") {
		t.Fatalf("update.log does not say what moved: %s", log)
	}
}
