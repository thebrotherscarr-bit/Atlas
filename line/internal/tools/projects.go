package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"atlas/line/internal/tenant"
)

// projects reads what THE MAKER made (the core's manjuel/maker.py), for the
// glass. The maker's piece 2, 2026-09-21, on the operator's word: "go on piece
// 2" -- the page on the glass, and a project list to pick one from.
//
// A PROJECT IS A FOLDER UNDER projects/ WITH A HISTORY OF ITS OWN. The maker
// writes one page, projects/<name>/index.html, and saves every change as a
// commit INSIDE the project, never in the ground. So the list is read off each
// project's own history -- the versions and their notes, as the maker wrote
// them -- and a page can be served as it stands now or as any version was.
//
// WHAT IT WILL NOT DO, and these are the whole design:
//
//   - IT READS. `Writes: false`, and nothing here opens a file for writing or
//     runs a git verb that changes anything. The maker is the only writer of a
//     project (LAW 8, one write-path); a second one in the glass is what this
//     estate refuses.
//   - A NAME IS A NAME, NOT A PATH. It must look like one the maker makes --
//     lowercase letters, digits and hyphens -- so "../.env", a drive letter or
//     a separator never reaches the disk. And every folder on the way -- the
//     projects/ folder, the project, its `.git` -- must be a PLAIN folder, not
//     a link or a junction to somewhere else, and the page a plain file.
//   - GIT NEVER ANSWERS FOR THE GROUND. projects/ sits inside the ground, and
//     the ground is itself a repository: git asked about a folder with no
//     history of its own walks UP and answers for the ground. So a folder with
//     no `.git` is not a project, and every git call names the project's own
//     `.git` outright rather than letting git look for one.
//   - BOUNDED. Every git call has a deadline (ESTATE LAW 7), and a page is
//     served only up to projectPageCap.
func toolProjects(t tenant.Tenant, args map[string]any) (string, error) {
	switch action := strings.ToLower(strings.TrimSpace(str(args, "action"))); action {
	case "", "list":
		return projectsList(t)
	case "page":
		return projectsPage(t, strings.TrimSpace(str(args, "name")), args["version"])
	default:
		return "", fmt.Errorf("projects takes action list or page, not %q", action)
	}
}

const (
	projectsDir = "projects"
	projectPage = "index.html"
	// A page the maker made is one self-contained file, and the Coder's window
	// keeps it small. A megabyte is far past anything it writes, and small
	// enough that the answer crosses the glass's wire whole.
	projectPageCap = 1 << 20
	projectGitWait = 10 * time.Second
)

// The shape of a name the maker gives a folder (maker.name_for): letters,
// digits and hyphens, and `-2`, `-3` for a second of the same name.
var projectNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// The maker's own subject line, "Version 3: make it faster". The note is what
// follows the colon; a subject that does not match is shown whole.
var projectVersionRe = regexp.MustCompile(`^Version \d+: `)

type projectVersion struct {
	N    int    `json:"n"`
	Sha  string `json:"sha"`
	Note string `json:"note"`
	When string `json:"when"`
}

type projectRow struct {
	Name     string           `json:"name"`
	Versions []projectVersion `json:"versions"`
	Bytes    int64            `json:"bytes"`
	Lines    int              `json:"lines"`
	Updated  string           `json:"updated"`
	Error    string           `json:"error,omitempty"`
}

// plainDir: a real folder, not a link or a junction to somewhere else.
//
// LOAD-BEARING ON WINDOWS, AND MEASURED, NOT ASSUMED (2026-09-21, this machine,
// Go 1.26). The first cut asked filepath.EvalSymlinks whether a project
// resolved inside projects/ -- and EvalSymlinks does not follow a junction, so
// a junction under projects/ pointing out of the ground read as inside, and
// its `.git` was found straight through it. A junction needs no privilege to
// make. Lstat reports it as irregular (and a symlink as a symlink), and
// neither is a plain folder, so neither is walked into.
func plainDir(p string) bool {
	fi, err := os.Lstat(p)
	return err == nil && fi.IsDir() && fi.Mode()&(os.ModeSymlink|os.ModeIrregular) == 0
}

// plainFile: a real file, so reading a page never follows a link out.
func plainFile(p string) (os.FileInfo, bool) {
	fi, err := os.Lstat(p)
	return fi, err == nil && fi.Mode().IsRegular()
}

// projectDir answers where a project is, or why the name is not one.
func projectDir(t tenant.Tenant, name string) (string, error) {
	if !projectNameRe.MatchString(name) {
		return "", fmt.Errorf("no project named %q -- a project's name is lowercase "+
			"letters, digits and hyphens, as the list shows it", name)
	}
	root := filepath.Join(t.Home, projectsDir)
	if !plainDir(root) {
		return "", fmt.Errorf("no project named %q -- %s has no projects/ folder of "+
			"its own", name, t.Name)
	}
	dir := filepath.Join(root, name)
	if _, err := os.Lstat(dir); err != nil {
		return "", fmt.Errorf("no project named %q in %s", name, t.Name)
	}
	if !plainDir(dir) {
		return "", fmt.Errorf("REFUSED -- %q is a link, not a folder inside "+
			"projects/. %s", name, WallLaw)
	}
	if !plainDir(filepath.Join(dir, ".git")) {
		return "", fmt.Errorf("%q is a folder with no history of its own, so it is "+
			"not a project", name)
	}
	return dir, nil
}

// projectGit runs one read-only git verb against the project's OWN history.
func projectGit(dir string, args ...string) spawnResult {
	full := append([]string{"--git-dir", filepath.Join(dir, ".git"),
		"--work-tree", dir}, args...)
	return spawn("git", full, spawnOpts{Dir: dir, Timeout: projectGitWait, Env: gitEnv()})
}

// projectVersions is the project's history, oldest first, numbered as the
// maker numbers it: the first save is version 1.
func projectVersions(dir string) ([]projectVersion, error) {
	res := projectGit(dir, "log", "--reverse", "--format=%h%x09%ct%x09%s")
	if res.Err != nil {
		// A history with no saves yet is a project with no versions, not a fault.
		if strings.Contains(res.Stderr, "does not have any commits") {
			return []projectVersion{}, nil
		}
		return nil, fmt.Errorf("its history could not be read: %s",
			strings.TrimSpace(firstLines(res.Combined, 2)))
	}
	out := []projectVersion{}
	for _, line := range strings.Split(strings.TrimSpace(res.Stdout), "\n") {
		parts := strings.SplitN(strings.TrimRight(line, "\r"), "\t", 3)
		if len(parts) != 3 {
			continue
		}
		when := ""
		if secs, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
			when = time.Unix(secs, 0).UTC().Format(time.RFC3339)
		}
		out = append(out, projectVersion{
			N: len(out) + 1, Sha: parts[0], When: when,
			Note: projectVersionRe.ReplaceAllString(parts[2], ""),
		})
	}
	return out, nil
}

func projectsList(t tenant.Tenant) (string, error) {
	rows := []projectRow{}
	ents, err := os.ReadDir(filepath.Join(t.Home, projectsDir))
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	for _, e := range ents {
		if !e.IsDir() || !projectNameRe.MatchString(e.Name()) {
			continue
		}
		dir, err := projectDir(t, e.Name())
		if err != nil {
			continue // not a project: nothing of it is shown
		}
		row := projectRow{Name: e.Name(), Versions: []projectVersion{}}
		if vs, err := projectVersions(dir); err != nil {
			row.Error = err.Error()
		} else {
			row.Versions = vs
		}
		if st, ok := plainFile(filepath.Join(dir, projectPage)); ok {
			row.Bytes = st.Size()
			row.Updated = st.ModTime().UTC().Format(time.RFC3339)
			if st.Size() <= projectPageCap {
				if b, err := os.ReadFile(filepath.Join(dir, projectPage)); err == nil {
					row.Lines = bytes.Count(b, []byte("\n"))
				}
			}
		}
		if n := len(row.Versions); n > 0 && row.Versions[n-1].When != "" {
			row.Updated = row.Versions[n-1].When
		}
		rows = append(rows, row)
	}
	// NEWEST FIRST: the project touched last is the one most likely wanted.
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Updated > rows[j].Updated })
	return projectsJSON(map[string]any{"world": t.Name, "count": len(rows), "projects": rows})
}

func projectsPage(t tenant.Tenant, name string, version any) (string, error) {
	dir, err := projectDir(t, name)
	if err != nil {
		return "", err
	}
	want := 0
	switch v := version.(type) {
	case nil:
	case float64:
		want = int(v)
		if float64(want) != v {
			return "", fmt.Errorf("a version is a whole number, not %v", v)
		}
	case int:
		want = v
	case string:
		if s := strings.TrimSpace(v); s != "" {
			if want, err = strconv.Atoi(s); err != nil {
				return "", fmt.Errorf("a version is a number, not %q", s)
			}
		}
	default:
		return "", fmt.Errorf("a version is a number, not %T", version)
	}
	vs, err := projectVersions(dir)
	if err != nil {
		return "", fmt.Errorf("%s: %s", name, err)
	}
	var text, sha string
	if want == 0 {
		// THE PAGE AS IT STANDS -- the file the person opens in a browser.
		path := filepath.Join(dir, projectPage)
		st, ok := plainFile(path)
		if !ok {
			return "", fmt.Errorf("%s has no %s of its own yet", name, projectPage)
		}
		if st.Size() > projectPageCap {
			return "", fmt.Errorf("%s's page is %d bytes, past the %d this serves",
				name, st.Size(), projectPageCap)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		text = string(b)
		want = len(vs)
		if want > 0 {
			sha = vs[want-1].Sha
		}
	} else {
		if want < 1 || want > len(vs) {
			return "", fmt.Errorf("there is no version %d -- %s has versions 1 to %d",
				want, name, len(vs))
		}
		sha = vs[want-1].Sha
		res := projectGit(dir, "show", sha+":"+projectPage)
		if res.Err != nil {
			return "", fmt.Errorf("version %d of %s could not be read: %s", want, name,
				strings.TrimSpace(firstLines(res.Combined, 2)))
		}
		if len(res.Stdout) > projectPageCap {
			return "", fmt.Errorf("version %d of %s is past the %d bytes this serves",
				want, name, projectPageCap)
		}
		text = res.Stdout
	}
	return projectsJSON(map[string]any{
		"world": t.Name, "name": name, "version": want, "sha": sha,
		"bytes": len(text), "text": text,
	})
}

// projectsJSON writes the answer WITHOUT escaping <, > and & -- a page is
// mostly those, and escaping each as < would grow it six times over on
// its way to the glass.
func projectsJSON(v any) (string, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if err := enc.Encode(v); err != nil {
		return "", err
	}
	return strings.TrimRight(b.String(), "\n"), nil
}
