package tools

// AIDER BEHIND THE GLASS -- aider_run, aider_undo, aider_status (WHAT'S LEFT H12 and B20, 2026-10-04).
//
// His idea of 2026-10-02: a page on the glass that drives Aider, open source and proven at search-and-replace
// edits. His answers of 2026-10-04, on the card: download `aider-chat` from PyPI (yes), into `aider/` at the
// ground's root, driven by the coding seat's own model on Ollama. This file is the door's half; the front
// page's Aider tab is a thin view over it (webapp/static/js/agent.js).
//
//	WHAT A RUN IS. He types an instruction and names the files Aider may change (and, optionally, files it may
//	only read). The door copies exactly those into a scratch folder under aider/work/runs/, runs Aider on the
//	copies, headless, one message and out (`--message-file`), against the seat's model on loopback Ollama, and
//	then -- only if it finished and every Python file it changed still parses -- writes the changed files back,
//	each in its own terminator, and answers with the diff, what Aider said, and a run id. That id is also the
//	undo: the files as they were are kept beside it.
//
//	EVERY GATE THE ESTATE ALREADY HAS STANDS IN FRONT OF IT. This is the same discipline the council's own
//	`ground_edit` keeps, restated for a program that is not ours:
//
//	  THE LINE OF WORK. A write lands only while the repository the file belongs to stands on a branch that is
//	  not main or master and is not detached. The main line is his (RULE 6). The repository is the NEAREST one
//	  (atlas/ carries its own `.git` inside the ground), and a ground with none refuses.
//	  WHAT A MODEL MAY NEVER WRITE. The seats' own list (aiderNeverWritten, below): the law, the seats, the
//	  ledgers, the proof stamps, the governing files, a secret, client material, a binary, a `.git`. A
//	  read-only file is held to the shorter list of what may never be READ. A new folder is his to place
//	  (RULE 8): a new file lands in a folder that exists.
//	  HIS HAND ALONE. Both writing tools are ServiceOnly: no seat, agent or other client is ever handed
//	  Aider, and the council cannot park a call of it for him to approve. Aider as the flow's `attempt` engine
//	  is H12's other half, and it waits for a measurement and for him (see the WHAT'S LEFT line).
//	  THE SUITES AND HIS CLICK are not here and are not skipped: Aider's edit is unsaved work on a line of
//	  work, exactly as `ground_edit`'s is. He runs the suites (suite_run) and saves it (git_commit); only his
//	  Land click moves main.
//
//	AIDER RUNS INSIDE A WALL (aiderguard.go): writes only into the run's own scratch folder, reads only from
//	aider/ and the interpreter's library, loopback only for the network, no other process. Measured on the first
//	probe: Aider's `--yes-always` will otherwise write a file at any path the model names.
//
//	WHAT THE CHILD IS GIVEN. A built environment, never the door's: no key, no service wire (RULE 7). HOME and
//	TEMP are inside aider/work, its history and caches are in the run's folder, its analytics, update checks
//	and URL scraping are off, its model metadata is a file the door writes (so it asks GitHub for none), and
//	its context window is the coding seat's declared Context (the whole of the files has to fit in it, so a
//	file bigger than that is refused up front rather than silently truncated by the model's server).
//
//	WHAT IT DOES NOT DO. It does not run the suites, save, commit, push or land. It does not give an agent
//	Aider. It does not call a hosted model: the seat's model must be on THIS machine's rack, and a model that
//	is not is refused by name (RULE 4, as amended 2026-10-02: a hosted route is his to turn on, per seat, on a
//	measured need). It keeps the last 30 runs' folders; the undo of an older run is git's.
//
//	HOW LONG IT TAKES. Measured on the first runs (2026-10-04): the coding seat's 14b model generates at about 3.6 tokens a second on this
//	machine and reads a prompt at about 50, so an edit of a few hundred tokens against an 8 KB file is four to five minutes. The default wait is
//	therefore ten minutes (twenty at most), counted from the moment Aider starts: the model is loaded before that (aiderWarm), and a load is
//	not counted against the run.
//
//	RECORDED. Every run, refusal and undo is a line in the world's holds log (state/holds.jsonl) -- the log
//	hold_answer already writes, as shell_run does -- with the instruction (keys withheld).

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	gopath "path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"atlas/line/internal/rack"
	"atlas/line/internal/tenant"
)

const (
	aiderDirName       = "aider"
	aiderDefaultWait   = 600 * time.Second
	aiderMaxWait       = 1200 * time.Second
	aiderMaxMessage    = 4000
	aiderMaxEdit       = 4
	aiderMaxRead       = 4
	aiderMaxFileBytes  = 64 * 1024
	aiderKeepRuns      = 30
	aiderOutBytes      = 256 * 1024
	aiderSaidBytes     = 16 * 1024
	aiderDiffBytes     = 64 * 1024
	aiderCharsPerToken = 3
	aiderPromptTokens  = 3500
	aiderGuardMark     = "ATLAS-AIDER-GUARD"
	aiderDefaultCtx    = 8192
	aiderWarmWait      = 180 * time.Second
	aiderSeatFile      = "expert_coder.md"
)

// aiderRunRe is what a run id looks like. Every id that names a folder is matched against it first.
var aiderRunRe = regexp.MustCompile(`^a\d{8}-\d{6}-[0-9a-f]{4}$`)

// aiderTokensRe is the line Aider prints when it has spoken to the model.
var aiderTokensRe = regexp.MustCompile(`(?m)^Tokens:\s*(.+?)\s*$`)

// aiderMu holds one run at a time: one model in one VRAM, and one scratch folder being judged.
var aiderMu sync.Mutex

// aiderStart is how Aider is started. A variable so the strokes can stand a stub in front of the real
// program (the real one is 640 MB, a model and minutes) and still measure the door's own judgement: what is
// copied, what is written back, what is refused. The real program is exercised by the live proof.
var aiderStart = func(py string, args []string, opts spawnOpts) spawnResult { return spawn(py, args, opts) }

// ---- what a run answers --------------------------------------------------------------------------

type aiderFile struct {
	File    string `json:"file"`
	New     bool   `json:"new,omitempty"`
	Added   int    `json:"added"`
	Removed int    `json:"removed"`
}

// aiderAnswer is the one document aider_run and aider_undo answer with, for every outcome that is not a
// malformed call. The glass parses it; its keys are the contract (TestTheAiderTabIsTheDoorsAiderAndNothingElse reads them off these tags).
type aiderAnswer struct {
	State     string      `json:"state"` // ran | refused | undone
	Class     string      `json:"class,omitempty"`
	Why       []string    `json:"why,omitempty"`
	Run       string      `json:"run,omitempty"`
	Model     string      `json:"model,omitempty"`
	Line      string      `json:"line,omitempty"`
	Changed   []aiderFile `json:"changed,omitempty"`
	Diff      string      `json:"diff,omitempty"`
	Withheld  bool        `json:"withheld,omitempty"`
	Said      string      `json:"said,omitempty"`
	Tokens    string      `json:"tokens,omitempty"`
	Ignored   []string    `json:"ignored,omitempty"`
	Guarded   []string    `json:"guarded,omitempty"`
	Exit      *int        `json:"exit,omitempty"`
	Ms        int64       `json:"ms,omitempty"`
	TimedOut  bool        `json:"timed_out,omitempty"`
	Truncated bool        `json:"truncated,omitempty"`
	Note      string      `json:"note,omitempty"`
}

func (a aiderAnswer) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}

// aiderMFile is one file of a run as the manifest keeps it: enough to undo it, and to know the undo is safe.
type aiderMFile struct {
	File   string `json:"file"`
	New    bool   `json:"new,omitempty"`
	Eol    string `json:"eol"`
	Before string `json:"before,omitempty"` // sha256 of the file as it was
	After  string `json:"after"`            // sha256 of what was written
}

type aiderManifest struct {
	Run     string       `json:"run"`
	When    string       `json:"when"`
	Model   string       `json:"model"`
	Line    string       `json:"line"`
	Message string       `json:"message"`
	Files   []aiderMFile `json:"files"`
	Undone  string       `json:"undone,omitempty"`
}

// ---- where Aider lives ---------------------------------------------------------------------------

type aiderPlace struct {
	Root, Work, Runs, Home, Tmp, Python string
}

// aiderPlaceOf is the folder a world keeps Aider in: `aider/` at its root, with the venv beside a work
// folder (his word, 2026-10-04). Nothing here is created by asking; only a run creates inside work/.
func aiderPlaceOf(home string) aiderPlace {
	root := filepath.Join(home, aiderDirName)
	work := filepath.Join(root, "work")
	py := filepath.Join(root, "venv", "bin", "python")
	if runtime.GOOS == "windows" {
		py = filepath.Join(root, "venv", "Scripts", "python.exe")
	}
	return aiderPlace{Root: root, Work: work, Runs: filepath.Join(work, "runs"), Home: filepath.Join(work, "home"),
		Tmp: filepath.Join(work, "tmp"), Python: py}
}

func aiderInstalled(pl aiderPlace) bool {
	st, err := os.Stat(pl.Python)
	return err == nil && !st.IsDir()
}

// aiderVersion reads the installed version off the package's own metadata folder; it starts no process.
func aiderVersion(pl aiderPlace) string {
	for _, pat := range []string{
		filepath.Join(pl.Root, "venv", "Lib", "site-packages", "aider_chat-*.dist-info"),
		filepath.Join(pl.Root, "venv", "lib", "python3*", "site-packages", "aider_chat-*.dist-info"),
	} {
		if m, _ := filepath.Glob(pat); len(m) > 0 {
			return strings.TrimSuffix(strings.TrimPrefix(filepath.Base(m[0]), "aider_chat-"), ".dist-info")
		}
	}
	return ""
}

// ---- the model: the coding seat's own ------------------------------------------------------------

type aiderSeat struct {
	Model string
	Ctx   int
}

var aiderModelRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]*$`)

// aiderCoder reads the coding seat's declared model and context window from the world's own seat file, at
// the time of the call, so a change the operator makes to the seat is the change Aider follows.
func aiderCoder(home string) (aiderSeat, error) {
	raw, err := os.ReadFile(filepath.Join(home, "agents", aiderSeatFile))
	if err != nil {
		return aiderSeat{}, errors.New("this world declares no coding seat (agents/" + aiderSeatFile + "), so there is no model to name")
	}
	seat := readSeat(aiderSeatFile, strings.ReplaceAll(string(raw), "\r\n", "\n"), nil)
	fields, _ := seat["fields"].(map[string]string)
	model := strings.TrimSpace(fields["Model Target"])
	if !aiderModelRe.MatchString(model) {
		return aiderSeat{}, fmt.Errorf("the coding seat's Model Target is %q, which is not a model name this door will hand to Aider", model)
	}
	ctx := aiderDefaultCtx
	if n, err := strconv.Atoi(strings.TrimSpace(fields["Context"])); err == nil && n >= 2048 && n <= 131072 {
		ctx = n
	}
	return aiderSeat{Model: model, Ctx: ctx}, nil
}

// aiderRack is the loopback Ollama the seat's model must be on. The same host resolution every rack tool uses,
// with a wildcard address mapped to this machine, since Aider has to DIAL it.
func aiderRack() (base string, voices []rack.Voice, err error) {
	host, err := rack.Host()
	if err != nil {
		return "", nil, err
	}
	if u, perr := url.Parse(host); perr == nil {
		if h := u.Hostname(); h == "0.0.0.0" || h == "::" {
			u.Host = "127.0.0.1:" + u.Port()
			host = u.String()
		}
	}
	voices, err = rack.List(host)
	return host, voices, err
}

// aiderWarm loads the seat's model into the rack's memory, at the window Aider will ask for, BEFORE Aider is started. Measured on
// the first batch (2026-10-04): a request that reached Ollama in the moment it was unloading the model hung for four minutes, inside
// Aider's own timeout and with nothing to say for it, until a request of this kind freed it. Asked here, with a bound of its own, a
// rack that will not load the model is a refusal that names it, and Aider is never started. The window is Aider's: a different one
// would make Ollama unload the model and load it again when Aider asked.
func aiderWarm(base, model string, ctx int) error {
	body, _ := json.Marshal(map[string]any{"model": model, "prompt": "", "stream": false, "keep_alive": "10m",
		"options": map[string]any{"num_ctx": ctx}})
	client := &http.Client{Timeout: aiderWarmWait}
	resp, err := client.Post(base+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 300 {
		return fmt.Errorf("the rack answered %d", resp.StatusCode)
	}
	return nil
}

func aiderHasModel(voices []rack.Voice, model string) bool {
	for _, v := range voices {
		if v.Name == model || v.Name == model+":latest" {
			return true
		}
	}
	return false
}

// ---- what a model may never write, and never read ------------------------------------------------

// These are manjuel/skills.py's `_NEVER_WRITTEN_*` restated for Go. They are a COPY, and a copy drifts; the
// core's suite reads this file and holds the two lists equal (test_the_aider_door_keeps_what_the_seats_keep),
// so a seat's list that grows without this one is red where he runs the suites.
var aiderNeverTop = map[string]string{
	"worlds":          "worlds/ is another world's, closed until the operator points at it (ESTATE LAW 2, SITTING LAW 2)",
	"law":             "law/ is sealed, and a hand does not edit the law, so a model does not either",
	"foundation":      "foundation/ holds the founding texts, sealed elders among them, and is the operator's",
	"agents":          "agents/ is the seats' own declarations, the operator's to set, and hot-reloaded into his live sitting",
	"skills":          "skills/ is the estate's own tools, the operator's to set, and hot-reloaded into his live sitting",
	"sessions":        "sessions/ is the ledger, written by the chain alone (LAW 8: one write-path)",
	"logs":            "logs/ is the transcripts, written by the runs alone (LAW 8)",
	"index":           "index/ is derived from the ground and rebuilt, never authored",
	"state":           "state/ is the door's own record (LAW 8)",
	"flows":           "flows/ is the workflow engine's own store (LAW 8)",
	"memory":          "memory/ is the memory's own store (LAW 8)",
	"projects":        "projects/ is the maker's, each project its own repository with its own doors",
	"agent_workspace": "agent_workspace/ has doors of its own: write_file and edit_file",
	"bin":             "bin/ holds binaries, which are placed on the operator's allowance",
	"aider":           "aider/ is Aider's own environment, not part of the work",
}

var aiderNeverFiles = map[string]string{
	"claude.md":       "CLAUDE.md is the operator's standing rules",
	".gitignore":      ".gitignore says what is not part of the work, and is the operator's",
	".gitattributes":  ".gitattributes governs the terminators, and is the operator's",
	"index_roots.txt": "index_roots.txt governs what is indexed, on his ruling",
	"pipelines.md":    "pipelines.md is the pipelines' law, the operator's, and hot-reloaded into his live sitting",
	"commands.md":     "commands.md is the REPL's own commands, hot-reloaded into his live sitting",
	".env.example":    ".env.example names the dials and the keys' shapes, and is the operator's",
	"memory.md":       "memory.md is the memory's own record (LAW 8)",
	"seat_log.md":     "SEAT_LOG.md is the tolls, written by the chain alone (LAW 8)",
	"buildmap.md":     "BUILDMAP.md is generated from the code by tests/buildmap.py, never authored",
}

var aiderNeverStamps = []string{"tests/last_run.json", "tests/last_run.md", "tests/run_history.jsonl", "tests/last_audit.md"}

var aiderNeverSuffixes = []string{".exe", ".dll", ".bin", ".pyc", ".pyd", ".so", ".bundle", ".obj"}

// aiderNeverWritten is why a model may never write THIS path, or "" when it may. Judged on the cleaned
// slash path, by name -- never on what the caller meant. rel is already inside the ground (jailed).
func aiderNeverWritten(rel string) string {
	low := strings.ToLower(rel)
	parts := strings.Split(low, "/")
	for _, p := range parts {
		if p == ".git" {
			return "a `.git/` is the history, and the history is git's to write (LAW 8)"
		}
	}
	if why, ok := aiderNeverTop[parts[0]]; ok {
		return why
	}
	for _, s := range aiderNeverStamps {
		if low == s {
			return rel + " is a proof stamp, written by the suite that ran and by nothing else; a model never writes its own proof"
		}
	}
	if why, ok := aiderNeverFiles[parts[len(parts)-1]]; ok {
		return why
	}
	ext := strings.ToLower(gopath.Ext(rel))
	for _, s := range aiderNeverSuffixes {
		if ext == s {
			return gopath.Base(rel) + " is a binary, and a binary is placed on the operator's allowance, never written by a model"
		}
	}
	return ""
}

// aiderNeverRead is the shorter list: what Aider may not even be shown.
func aiderNeverRead(rel string) string {
	parts := strings.Split(strings.ToLower(rel), "/")
	for _, p := range parts {
		if p == ".git" {
			return "a `.git/` is the history, not a file to read into a model"
		}
	}
	switch parts[0] {
	case "worlds":
		return aiderNeverTop["worlds"]
	case "aider":
		return aiderNeverTop["aider"]
	}
	return ""
}

// aiderNames splits what the glass sent -- one path per line, or commas, or semicolons -- into clean,
// slash-spelled, de-duplicated relative paths.
func aiderNames(raw string) []string {
	var out []string
	seen := map[string]bool{}
	for _, p := range strings.FieldsFunc(raw, func(r rune) bool { return r == '\n' || r == '\r' || r == ',' || r == ';' }) {
		p = strings.Trim(strings.TrimSpace(p), "\"'`")
		p = strings.ReplaceAll(p, "\\", "/")
		for strings.HasPrefix(p, "./") {
			p = strings.TrimPrefix(p, "./")
		}
		if p == "" {
			continue
		}
		if k := strings.ToLower(p); !seen[k] {
			seen[k] = true
			out = append(out, p)
		}
	}
	return out
}

// aiderJudgePath is why a path may not be handed to Aider, or "". writing says whether Aider may change it.
func aiderJudgePath(t tenant.Tenant, rel string, writing bool) string {
	if why := shellRefusals(rel, t.Home); len(why) > 0 {
		return rel + ": " + why[0]
	}
	if r := jailed(t, rel); r != "" {
		return strings.TrimPrefix(r, "Refused: ")
	}
	clean := gopath.Clean(rel)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return rel + " is not a file inside this world"
	}
	if writing {
		if why := aiderNeverWritten(clean); why != "" {
			return rel + ": " + why
		}
		return ""
	}
	if why := aiderNeverRead(clean); why != "" {
		return rel + ": " + why
	}
	return ""
}

// ---- the line of work ----------------------------------------------------------------------------

// aiderRepoOf is the nearest repository a file belongs to, inside the ground: the first folder walking up from
// the file that carries a `.git`. "" when there is none.
func aiderRepoOf(home, abs string) string {
	home = filepath.Clean(home)
	dir := filepath.Dir(abs)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		if dir == home {
			return ""
		}
		parent := filepath.Dir(dir)
		if parent == dir || !strings.HasPrefix(parent, home) {
			return ""
		}
		dir = parent
	}
}

// aiderBranch is the branch a repository stands on, or "" when it is detached or cannot be asked.
func aiderBranch(repo string) string {
	res := spawn("git", []string{"symbolic-ref", "--short", "-q", "HEAD"}, spawnOpts{Dir: repo, Timeout: 10 * time.Second, Env: gitEnv()})
	if res.Err != nil {
		return ""
	}
	return strings.TrimSpace(res.Stdout)
}

func aiderIsLine(branch string) bool {
	return branch != "" && branch != "main" && branch != "master"
}

// aiderLineRefusal is the refusal for a file whose repository is not on a line of work, in the words the
// council's own tree doors use. "" when it is. branches caches one answer per repository.
func aiderLineRefusal(home, abs string, branches map[string]string) (branch, why string) {
	repo := aiderRepoOf(home, abs)
	if repo == "" {
		return "", "the ground is not under version control here, so a line of work cannot be checked, and a model's write lands only on a line of work"
	}
	b, seen := branches[repo]
	if !seen {
		b = aiderBranch(repo)
		branches[repo] = b
	}
	if !aiderIsLine(b) {
		where := "is detached"
		if b != "" {
			where = "stands on `" + b + "`"
		}
		name := filepath.Base(repo)
		return "", fmt.Sprintf("the repository at `%s` %s, and the main line is the operator's (RULE 6): open a line of work first "+
			"(`git_branch new <name>` through the door, or Lines of work on Version control) and write on it", name, where)
	}
	return b, ""
}

// ---- terminators ---------------------------------------------------------------------------------

// aiderEol is the one terminator a file uses, and false when it holds both.
func aiderEol(b []byte) (string, bool) {
	crlf := bytes.Count(b, []byte("\r\n"))
	lf := bytes.Count(b, []byte("\n")) - crlf
	switch {
	case crlf > 0 && lf > 0:
		return "", false
	case crlf > 0:
		return "\r\n", true
	}
	return "\n", true
}

// aiderSiblingEol is what a NEW file takes: its nearest sibling's of the same suffix, and LF where there is
// none (the terminator ruling of 2026-09-03, as terminator_for makes it).
func aiderSiblingEol(abs string) string {
	entries, err := os.ReadDir(filepath.Dir(abs))
	if err != nil {
		return "\n"
	}
	ext := filepath.Ext(abs)
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ext || e.Name() == filepath.Base(abs) {
			continue
		}
		b, err := os.ReadFile(filepath.Join(filepath.Dir(abs), e.Name()))
		if err != nil || len(b) == 0 || len(b) > 1<<20 {
			continue
		}
		if eol, ok := aiderEol(b); ok && bytes.Contains(b, []byte("\n")) {
			return eol
		}
	}
	return "\n"
}

func aiderEolName(eol string) string {
	if eol == "\r\n" {
		return "crlf"
	}
	return "lf"
}

func aiderApplyEol(flat, eol string) string {
	if eol == "\r\n" {
		return strings.ReplaceAll(flat, "\n", "\r\n")
	}
	return flat
}

func aiderLF(b []byte) string { return strings.ReplaceAll(string(b), "\r\n", "\n") }

func aiderSum(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// ---- one file of a run ---------------------------------------------------------------------------

type aiderPlan struct {
	Rel    string
	Abs    string
	New    bool
	Eol    string
	Before []byte
	Result []byte
	Branch string
}

func aiderAbs(t tenant.Tenant, rel string) string {
	return filepath.Join(t.Home, filepath.FromSlash(rel))
}

func aiderRecord(args map[string]any, t tenant.Tenant, c Caller, tool, what, note string) {
	reg := regOf(args)
	if reg == nil {
		return
	}
	reg.record(t.Home, what, Hold{
		ID: fmt.Sprintf("aider_%d", time.Now().UnixMilli()), Tool: tool,
		Caller: callerLabel(c), Project: t.Name, RBAC: rbacState(t, c),
	}, note)
}

func aiderWait(args map[string]any) time.Duration {
	d := aiderDefaultWait
	switch x := args["timeout"].(type) {
	case int:
		if x > 0 {
			d = time.Duration(x) * time.Second
		}
	case float64:
		if x > 0 {
			d = time.Duration(x * float64(time.Second))
		}
	case string:
		if f, err := strconv.ParseFloat(strings.TrimSpace(x), 64); err == nil && f > 0 {
			d = time.Duration(f * float64(time.Second))
		}
	}
	if d > aiderMaxWait {
		d = aiderMaxWait
	}
	return d
}

func aiderNewID() string {
	var b [2]byte
	_, _ = rand.Read(b[:])
	return "a" + time.Now().Format("20060102-150405") + "-" + hex.EncodeToString(b[:])
}

func aiderNonce() string {
	var b [12]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// aiderWriteFile writes a file inside aider/work, making its folder if it is one of the run's own.
func aiderWriteFile(p string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o644)
}

// aiderKeep writes a file only when its bytes differ, so a file Aider reads every run is not touched every run.
func aiderKeep(p string, b []byte) error {
	if old, err := os.ReadFile(p); err == nil && bytes.Equal(old, b) {
		return nil
	}
	return aiderWriteFile(p, b)
}

// ---- Aider's own world ---------------------------------------------------------------------------

// aiderEnv is the whole environment Aider gets. Built from a list of names, never by removing the ones that
// look secret: a variable nobody thought of is not passed. No key and no service wire reach it (RULE 7).
// There is deliberately no proxy variable: a library that tries to leave the machine reaches the wall's own
// address lookup first, and the wall says so by name.
func aiderEnv(pl aiderPlace, ollama, nonce, wall string) []string {
	keep := map[string]bool{"SYSTEMROOT": true, "SYSTEMDRIVE": true, "COMSPEC": true, "WINDIR": true, "PATHEXT": true,
		"NUMBER_OF_PROCESSORS": true, "PROCESSOR_ARCHITECTURE": true, "OS": true, "USERNAME": true, "COMPUTERNAME": true}
	var env []string
	for _, kv := range os.Environ() {
		if i := strings.IndexByte(kv, '='); i > 0 && keep[strings.ToUpper(kv[:i])] {
			env = append(env, kv)
		}
	}
	scripts := filepath.Dir(pl.Python)
	pathv := scripts + ":/usr/bin:/bin"
	if runtime.GOOS == "windows" {
		pathv = scripts + string(os.PathListSeparator) + filepath.Join(os.Getenv("SYSTEMROOT"), "System32")
	}
	vol := filepath.VolumeName(pl.Home)
	return append(env,
		"PATH="+pathv,
		"HOME="+pl.Home, "USERPROFILE="+pl.Home, "HOMEDRIVE="+vol, "HOMEPATH="+strings.TrimPrefix(pl.Home, vol),
		"APPDATA="+filepath.Join(pl.Home, "AppData", "Roaming"), "LOCALAPPDATA="+filepath.Join(pl.Home, "AppData", "Local"),
		"TEMP="+pl.Tmp, "TMP="+pl.Tmp, "TMPDIR="+pl.Tmp,
		"OLLAMA_API_BASE="+ollama,
		"LITELLM_LOCAL_MODEL_COST_MAP=True", "LITELLM_LOG=ERROR",
		"PYTHONDONTWRITEBYTECODE=1", "PYTHONUTF8=1", "PYTHONIOENCODING=utf-8",
		"NO_COLOR=1", "TERM=dumb", "GIT_PYTHON_REFRESH=quiet",
		"ATLAS_AIDER_WALL="+wall, "ATLAS_AIDER_NONCE="+nonce)
}

// aiderArgs is Aider's own command line, one message and out. Every flag is here for a reason the probe found:
// --no-git so it never looks for or writes a repository (the door saves, he lands); --yes-always because
// nobody is there to answer (the wall is what stops it writing where it should not); no analytics, update
// check, URL scraping, shell suggestions, lint or test, because each reaches for something (the network, a
// process) that the wall would only refuse; --map-tokens 0 because there is no repo to map; LF endings, the
// door puts each file's own back.
func aiderArgs(pl aiderPlace, model, runDir string, edit, read []string) []string {
	a := []string{
		"--model", "ollama_chat/" + model, "--edit-format", "diff",
		"--no-git", "--no-gitignore", "--no-auto-commits", "--no-dirty-commits",
		"--yes-always", "--message-file", filepath.Join(runDir, "message.txt"),
		"--no-stream", "--no-pretty", "--no-fancy-input",
		"--no-analytics", "--analytics-disable", "--no-check-update", "--no-show-release-notes", "--no-show-model-warnings",
		"--no-suggest-shell-commands", "--no-auto-lint", "--no-auto-test", "--map-tokens", "0", "--no-detect-urls",
		"--chat-history-file", filepath.Join(runDir, "chat.md"),
		"--input-history-file", filepath.Join(runDir, "input.history"),
		"--llm-history-file", filepath.Join(runDir, "llm.history"),
		"--env-file", filepath.Join(pl.Work, "empty.env"),
		"--model-settings-file", filepath.Join(pl.Work, "model-settings.yml"),
		"--model-metadata-file", filepath.Join(pl.Work, "model-metadata.json"),
		"--encoding", "utf-8", "--line-endings", "lf",
	}
	for _, f := range edit {
		a = append(a, "--file", f)
	}
	for _, f := range read {
		a = append(a, "--read", f)
	}
	return a
}

// aiderWorkFiles writes the three small files Aider is pointed at: an empty .env (so it looks for none), the
// model's settings (its context window, the seat's own) and the model's metadata (so it asks no server).
func aiderWorkFiles(pl aiderPlace, model string, ctx int) error {
	if err := os.MkdirAll(pl.Tmp, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(pl.Home, 0o755); err != nil {
		return err
	}
	if err := aiderKeep(filepath.Join(pl.Work, "empty.env"), nil); err != nil {
		return err
	}
	settings := fmt.Sprintf("- name: ollama_chat/%s\n  extra_params:\n    num_ctx: %d\n", model, ctx)
	if err := aiderKeep(filepath.Join(pl.Work, "model-settings.yml"), []byte(settings)); err != nil {
		return err
	}
	meta, _ := json.Marshal(map[string]any{"ollama_chat/" + model: map[string]any{
		"max_input_tokens": ctx, "max_output_tokens": 4096, "max_tokens": 4096,
		"input_cost_per_token": 0, "output_cost_per_token": 0, "litellm_provider": "ollama_chat", "mode": "chat"}})
	return aiderKeep(filepath.Join(pl.Work, "model-metadata.json"), meta)
}

func aiderWallFor(t tenant.Tenant, pl aiderPlace, runDir string) aiderWall {
	profile, _ := os.UserHomeDir()
	return aiderWall{
		Ground: t.Home, Profile: profile, Aider: pl.Root,
		WriteDirs:  []string{filepath.Join(runDir, "scratch"), pl.Home, pl.Tmp},
		WriteFiles: []string{filepath.Join(runDir, "chat.md"), filepath.Join(runDir, "input.history"), filepath.Join(runDir, "llm.history")},
	}
}

// ---- reading what Aider left ---------------------------------------------------------------------

// aiderGuardLines splits the wall's own lines out of stderr: what it refused, and whether it announced itself.
func aiderGuardLines(stderr, nonce string) (rest string, armed bool, refused []string) {
	var keep []string
	for _, ln := range strings.Split(strings.ReplaceAll(stderr, "\r\n", "\n"), "\n") {
		switch {
		case ln == aiderGuardMark+" "+nonce+" armed":
			armed = true
		case strings.HasPrefix(ln, aiderGuardMark+" refused: "):
			if len(refused) < 6 {
				refused = append(refused, clip(strings.TrimPrefix(ln, aiderGuardMark+" refused: "), 240))
			}
		default:
			keep = append(keep, ln)
		}
	}
	return strings.Join(keep, "\n"), armed, refused
}

// aiderDiff is the change one file made, as git writes it: the file as it was against the file as it stands,
// compared where they sit in the run's folder so the header names the file and not a scratch path.
func aiderDiff(runDir string, f aiderPlan) (string, aiderFile) {
	info := aiderFile{File: f.Rel, New: f.New}
	before := gopath.Join("before", f.Rel)
	if f.New {
		before = "/dev/null" // git's own name for nothing, on every platform
	}
	res := spawn("git", []string{"diff", "--no-index", "--no-color", "--unified=3", "--", before, gopath.Join("after", f.Rel)},
		spawnOpts{Dir: runDir, Timeout: 20 * time.Second, Env: gitEnv(), MaxBytes: aiderDiffBytes})
	var ee *exec.ExitError
	if res.Err != nil && !(errors.As(res.Err, &ee) && ee.ExitCode() == 1) {
		return "(the diff of " + f.Rel + " could not be made: " + firstLine(res.Combined+" "+res.Err.Error()) + ")\n", info
	}
	text := strings.ReplaceAll(res.Stdout, "\r\n", "\n")
	text = strings.ReplaceAll(text, "a/before/", "a/")
	text = strings.ReplaceAll(text, "b/after/", "b/")
	for _, ln := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(ln, "+++") || strings.HasPrefix(ln, "---"):
		case strings.HasPrefix(ln, "+"):
			info.Added++
		case strings.HasPrefix(ln, "-"):
			info.Removed++
		}
	}
	return strings.ReplaceAll(text, "\r", ""), info
}

// aiderParses is the structural gate for a Python file: parsed by the venv's own interpreter, never run.
func aiderParses(pl aiderPlace, file string) (bool, string) {
	res := spawn(pl.Python, []string{"-I", "-B", "-X", "utf8", "-c", aiderParseSource, file},
		spawnOpts{Dir: filepath.Dir(file), Timeout: 30 * time.Second, MaxBytes: 4096, Grace: 2 * time.Second,
			CleanEnv: aiderEnv(pl, "", "", "")})
	if res.Err == nil {
		return true, ""
	}
	if msg := strings.TrimSpace(res.Stdout); msg != "" {
		return false, clip(msg, 200)
	}
	return false, "the parse check could not run: " + clip(firstLine(res.Combined+" "+res.Err.Error()), 200)
}

// ---- aider_run -----------------------------------------------------------------------------------

func toolAiderRun(t tenant.Tenant, args map[string]any) (string, error) {
	c := callerOf(args)
	// Registry.Call has refused every other caller already (ServiceOnly); a tool that is reached any other
	// way refuses for itself.
	if !c.Service {
		return "", errors.New("refused: Aider is the operator's own hand and is offered to his glass alone")
	}
	message := strings.TrimSpace(strings.ReplaceAll(str(args, "message"), "\r\n", "\n"))
	ans := aiderAnswer{}
	refuse := func(why ...string) (string, error) {
		ans.State, ans.Class, ans.Why = "refused", string(ShellRefused), why
		aiderRecord(args, t, c, "aider_run", "aider_refused", strings.Join(why, "; ")+" :: "+clip(shellRedact(t.Home, message), 300))
		return ans.String(), nil
	}

	// -- what he asked, judged before anything is touched
	if message == "" {
		return refuse("nothing to ask: write what Aider should change")
	}
	if len(message) > aiderMaxMessage {
		return refuse(fmt.Sprintf("longer than %d characters; it is not an instruction, it is a file", aiderMaxMessage))
	}
	if why := shellRefusals(message, t.Home); len(why) > 0 {
		return refuse(why...)
	}
	editNames := aiderNames(str(args, "files"))
	if len(editNames) == 0 {
		return refuse("name the files Aider may change: it works on the files it is given, and none were")
	}
	if len(editNames) > aiderMaxEdit {
		return refuse(fmt.Sprintf("%d files to change; Aider takes at most %d at a time, and a small local model does better with fewer", len(editNames), aiderMaxEdit))
	}
	var readNames []string
	editSet := map[string]bool{}
	for _, n := range editNames {
		editSet[strings.ToLower(n)] = true
	}
	for _, n := range aiderNames(str(args, "context")) {
		if !editSet[strings.ToLower(n)] {
			readNames = append(readNames, n)
		}
	}
	if len(readNames) > aiderMaxRead {
		return refuse(fmt.Sprintf("%d read-only files; at most %d", len(readNames), aiderMaxRead))
	}
	for _, n := range editNames {
		if why := aiderJudgePath(t, n, true); why != "" {
			return refuse(why)
		}
	}
	for _, n := range readNames {
		if why := aiderJudgePath(t, n, false); why != "" {
			return refuse(why)
		}
	}

	// -- the files, as they stand now
	branches := map[string]string{}
	var plans []*aiderPlan
	flatBytes := len(message)
	for _, n := range editNames {
		p := &aiderPlan{Rel: gopath.Clean(n), Abs: aiderAbs(t, gopath.Clean(n))}
		branch, why := aiderLineRefusal(t.Home, p.Abs, branches)
		if why != "" {
			return refuse(n + ": " + why)
		}
		p.Branch = branch
		st, err := os.Stat(p.Abs)
		switch {
		case err != nil && os.IsNotExist(err):
			p.New = true
			if pst, perr := os.Stat(filepath.Dir(p.Abs)); perr != nil || !pst.IsDir() {
				return refuse(n + ": its folder does not exist, and a model makes no folder; a new folder is his to place (RULE 8)")
			}
			p.Eol = aiderSiblingEol(p.Abs)
		case err != nil:
			return refuse(n + ": " + err.Error())
		case st.IsDir():
			return refuse(n + " is a folder, not a file")
		case st.Size() > aiderMaxFileBytes:
			return refuse(fmt.Sprintf("%s is %d KB; Aider is shown a whole file at a time and the largest this door will give it is %d KB", n, st.Size()/1024, aiderMaxFileBytes/1024))
		default:
			b, rerr := os.ReadFile(p.Abs)
			if rerr != nil {
				return refuse(n + ": " + rerr.Error())
			}
			eol, ok := aiderEol(b)
			if !ok {
				return refuse(n + " has MIXED line endings, so a write cannot keep what it has; that is a file to fix by hand, not through Aider")
			}
			if !utf8.Valid(b) || bytes.IndexByte(b, 0) >= 0 {
				return refuse(n + " is not text Aider can read")
			}
			p.Before, p.Eol = b, eol
			flatBytes += len(aiderLF(b))
		}
		plans = append(plans, p)
	}
	type readFile struct {
		Rel  string
		Text string
	}
	var reads []readFile
	for _, n := range readNames {
		abs := aiderAbs(t, gopath.Clean(n))
		st, err := os.Stat(abs)
		if err != nil || st.IsDir() {
			return refuse(n + " is not a file in this world, so it cannot be read into the chat")
		}
		if st.Size() > aiderMaxFileBytes {
			return refuse(fmt.Sprintf("%s is %d KB; the largest file this door will show Aider is %d KB", n, st.Size()/1024, aiderMaxFileBytes/1024))
		}
		b, rerr := os.ReadFile(abs)
		if rerr != nil || !utf8.Valid(b) || bytes.IndexByte(b, 0) >= 0 {
			return refuse(n + " is not text Aider can read")
		}
		flatBytes += len(aiderLF(b))
		reads = append(reads, readFile{Rel: gopath.Clean(n), Text: aiderLF(b)})
	}

	// -- the machine: Aider installed, the seat's model on this rack, the window big enough
	pl := aiderPlaceOf(t.Home)
	if !aiderInstalled(pl) {
		return refuse("Aider is not installed in this world: there is no " + aiderDirName + "/venv here (the install is in RUNBOOK.md)")
	}
	seat, err := aiderCoder(t.Home)
	if err != nil {
		return refuse(err.Error())
	}
	ans.Model = seat.Model
	ollama, voices, err := aiderRack()
	if err != nil {
		return refuse("the rack is not reachable, so there is no model to drive Aider: " + clip(err.Error(), 200))
	}
	if !aiderHasModel(voices, seat.Model) {
		return refuse(fmt.Sprintf("the coding seat's model %q is not on this machine's rack; Aider drives local models only (RULE 4)", seat.Model))
	}
	budget := (seat.Ctx - aiderPromptTokens) * aiderCharsPerToken
	if budget < 4096 {
		budget = 4096
	}
	if flatBytes > budget {
		return refuse(fmt.Sprintf("the files and the instruction come to about %d KB and the coding seat's window (%d tokens) holds about %d KB "+
			"beside Aider's own prompt: name smaller files, or raise the seat's Context", (flatBytes+1023)/1024, seat.Ctx, budget/1024))
	}

	// -- one run at a time
	if !aiderMu.TryLock() {
		return refuse("an Aider run is already going; wait for it, one model is in the memory")
	}
	defer aiderMu.Unlock()

	if err := aiderWarm(ollama, seat.Model, seat.Ctx); err != nil {
		return refuse(fmt.Sprintf("the rack did not load %s within %d s (%s); nothing was started. Try again: a rack caught in unloading a model can need a minute",
			seat.Model, int(aiderWarmWait.Seconds()), clip(err.Error(), 160)))
	}

	id := aiderNewID()
	runDir := filepath.Join(pl.Runs, id)
	scratch := filepath.Join(runDir, "scratch")
	ans.Run = id
	fail := func(what string, err error) (string, error) {
		return refuse(what + ": " + err.Error())
	}
	if err := aiderWorkFiles(pl, seat.Model, seat.Ctx); err != nil {
		return fail("Aider's work folder could not be made ready", err)
	}
	if err := os.MkdirAll(scratch, 0o755); err != nil {
		return fail("the run's scratch folder could not be made", err)
	}
	// THE SCRATCH FOLDER IS ITS OWN REPOSITORY. Aider looks UP from where it stands for a git repository (and for the
	// `.aider.conf.yml` beside it), even with --no-git; the scratch folder is inside the world, so without this it finds the
	// world's own `.git` and the wall rightly refuses to let it read it -- measured on the first real run, 2026-10-04, as an
	// uncaught PermissionError at start. An empty repository of its own stops the search where it is allowed to look.
	if res := spawn("git", []string{"init", "-q"}, spawnOpts{Dir: scratch, Timeout: 20 * time.Second, Env: gitEnv()}); res.Err != nil {
		return fail("the scratch folder could not be made its own repository, so Aider would find this world's", errors.New(firstLine(res.Combined+" "+res.Err.Error())))
	}
	if err := aiderWriteFile(filepath.Join(runDir, "message.txt"), []byte(message+"\n")); err != nil {
		return fail("the instruction could not be written", err)
	}
	var editArgs, readArgs []string
	for _, p := range plans {
		editArgs = append(editArgs, p.Rel)
		target := filepath.Join(scratch, filepath.FromSlash(p.Rel))
		if p.New {
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fail("the scratch folder for "+p.Rel+" could not be made", err)
			}
			continue
		}
		if err := aiderWriteFile(target, []byte(aiderLF(p.Before))); err != nil {
			return fail("the scratch copy of "+p.Rel+" could not be made", err)
		}
	}
	for _, r := range reads {
		readArgs = append(readArgs, r.Rel)
		if err := aiderWriteFile(filepath.Join(scratch, filepath.FromSlash(r.Rel)), []byte(r.Text)); err != nil {
			return fail("the scratch copy of "+r.Rel+" could not be made", err)
		}
	}

	nonce := aiderNonce()
	wall, _ := json.Marshal(aiderWallFor(t, pl, runDir))
	argv := aiderGuardArgs("aider", aiderArgs(pl, seat.Model, runDir, editArgs, readArgs)...)
	wait := aiderWait(args)
	began := time.Now()
	res := aiderStart(pl.Python, argv, spawnOpts{Dir: scratch, Timeout: wait, CleanEnv: aiderEnv(pl, ollama, nonce, string(wall)),
		MaxBytes: aiderOutBytes, Grace: 3 * time.Second, KillTree: true})
	ans.Ms = time.Since(began).Milliseconds()

	errText, armed, refused := aiderGuardLines(res.Stderr, nonce)
	ans.Guarded = refused
	said := res.Stdout
	if strings.TrimSpace(errText) != "" {
		if said != "" && !strings.HasSuffix(said, "\n") {
			said += "\n"
		}
		said += errText
	}
	var trunc bool
	ans.Said, trunc = shellShape(t.Home, said, aiderSaidBytes)
	ans.Truncated = trunc || res.Truncated
	if m := aiderTokensRe.FindStringSubmatch(ans.Said); m != nil {
		ans.Tokens = m[1]
	}
	exit := 0
	var ee *exec.ExitError
	switch {
	case res.TimedOut:
		exit = -1
		ans.TimedOut = true
	case res.Err == nil:
	case errors.As(res.Err, &ee):
		exit = ee.ExitCode()
	default:
		exit = -1
		ans.Note = res.Err.Error() + "; nothing was written"
	}
	ans.Exit = &exit
	note := func(what string) string {
		return fmt.Sprintf("%s run=%s model=%s exit=%d %dms :: %s", what, id, seat.Model, exit, ans.Ms, clip(shellRedact(t.Home, message), 300))
	}

	// A wall that did not announce itself was not there: nothing of the run is believed.
	if !armed {
		ans.State, ans.Class = "refused", string(ShellRefused)
		ans.Why = []string{"THE WALL did not announce itself, so whatever ran did not run behind it; the run is thrown away and nothing was written"}
		aiderRecord(args, t, c, "aider_run", "aider_refused", note("the wall was not armed"))
		return ans.String(), nil
	}
	ans.State = "ran"
	if ans.TimedOut {
		ans.Note = fmt.Sprintf("did not finish inside %d s and was ended; nothing was written", int(wait.Seconds()))
		aiderRecord(args, t, c, "aider_run", "aider_ran", note("timed out"))
		return ans.String(), nil
	}
	if exit != 0 {
		if ans.Note == "" {
			ans.Note = fmt.Sprintf("Aider exited with code %d; nothing was written", exit)
		}
		aiderRecord(args, t, c, "aider_run", "aider_ran", note("failed"))
		return ans.String(), nil
	}

	// -- what it left. Only the files he named are looked at; anything else it made is reported, never kept.
	var changed []*aiderPlan
	for _, p := range plans {
		raw, rerr := os.ReadFile(filepath.Join(scratch, filepath.FromSlash(p.Rel)))
		if rerr != nil {
			ans.Ignored = append(ans.Ignored, p.Rel+" (gone from the scratch folder; left as it is)")
			continue
		}
		if !utf8.Valid(raw) || bytes.IndexByte(raw, 0) >= 0 || len(raw) > 4*aiderMaxFileBytes {
			ans.Note = "NOT written: Aider's version of " + p.Rel + " is not text of a size this door will write"
			ans.Withheld = true
			aiderRecord(args, t, c, "aider_run", "aider_ran", note("withheld"))
			return ans.String(), nil
		}
		out := aiderApplyEol(strings.ReplaceAll(string(raw), "\r\n", "\n"), p.Eol)
		if (!p.New && out == string(p.Before)) || (p.New && len(raw) == 0) {
			continue
		}
		p.Result = []byte(out)
		changed = append(changed, p)
	}
	known := map[string]bool{}
	for _, p := range plans {
		known[strings.ToLower(p.Rel)] = true
	}
	for _, r := range reads {
		known[strings.ToLower(r.Rel)] = true
		if raw, rerr := os.ReadFile(filepath.Join(scratch, filepath.FromSlash(r.Rel))); rerr == nil && aiderLF(raw) != r.Text && len(ans.Ignored) < 8 {
			ans.Ignored = append(ans.Ignored, r.Rel+" (read-only; Aider changed its copy, which is never written back)")
		}
	}
	_ = filepath.WalkDir(scratch, func(p string, d os.DirEntry, werr error) error {
		if werr == nil && d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if werr != nil || d.IsDir() || len(ans.Ignored) >= 8 {
			return nil
		}
		rel, rerr := filepath.Rel(scratch, p)
		if rerr != nil {
			return nil
		}
		if rel = filepath.ToSlash(rel); !known[strings.ToLower(rel)] && !strings.Contains("/"+rel, "/.aider") {
			ans.Ignored = append(ans.Ignored, rel+" (made by Aider; not one of the files he named, so not kept)")
		}
		return nil
	})
	if len(changed) == 0 {
		ans.Note = "Aider changed nothing: it answered without an edit, or the edit it wrote matched nothing. What it said is below."
		aiderRecord(args, t, c, "aider_run", "aider_ran", note("no change"))
		aiderPrune(pl)
		return ans.String(), nil
	}

	// -- the copies, the structural gate, the diff
	for _, p := range changed {
		if !p.New {
			if err := aiderWriteFile(filepath.Join(runDir, "before", filepath.FromSlash(p.Rel)), p.Before); err != nil {
				return fail("the undo copy of "+p.Rel+" could not be kept, so nothing was written", err)
			}
		}
		if err := aiderWriteFile(filepath.Join(runDir, "after", filepath.FromSlash(p.Rel)), p.Result); err != nil {
			return fail("the copy of Aider's "+p.Rel+" could not be kept, so nothing was written", err)
		}
	}
	var diffs []string
	for _, p := range changed {
		d, info := aiderDiff(runDir, *p)
		diffs = append(diffs, d)
		ans.Changed = append(ans.Changed, info)
	}
	ans.Diff, _ = shellShape(t.Home, strings.Join(diffs, ""), aiderDiffBytes)
	for _, p := range changed {
		if strings.ToLower(filepath.Ext(p.Rel)) != ".py" {
			continue
		}
		if ok, why := aiderParses(pl, filepath.Join(runDir, "after", filepath.FromSlash(p.Rel))); !ok {
			ans.Withheld = true
			ans.Note = "NOT written: Aider's version of " + p.Rel + " does not parse (" + why + "). Nothing in the world changed; its edit is above for you to read."
			aiderRecord(args, t, c, "aider_run", "aider_ran", note("withheld: "+p.Rel+" does not parse"))
			aiderPrune(pl)
			return ans.String(), nil
		}
	}

	// -- the write: the ground must be as it was when Aider was started, and every file lands or none does
	for _, p := range changed {
		now, rerr := os.ReadFile(p.Abs)
		if (p.New && rerr == nil) || (!p.New && (rerr != nil || !bytes.Equal(now, p.Before))) {
			ans.Withheld = true
			ans.Note = "NOT written: " + p.Rel + " changed while Aider was working, and its edit was made against the old file. Nothing was written; the edit is above, and in " + aiderDirName + "/work/runs/" + id + "/after/."
			aiderRecord(args, t, c, "aider_run", "aider_ran", note("withheld: "+p.Rel+" moved"))
			aiderPrune(pl)
			return ans.String(), nil
		}
	}
	var written []*aiderPlan
	rollback := func() {
		for i := len(written) - 1; i >= 0; i-- {
			w := written[i]
			if w.New {
				_ = os.Remove(w.Abs)
			} else {
				_ = os.WriteFile(w.Abs, w.Before, 0o644)
			}
		}
	}
	m := aiderManifest{Run: id, When: time.Now().UTC().Format(time.RFC3339), Model: seat.Model, Message: clip(shellRedact(t.Home, message), 500)}
	for _, p := range changed {
		mode := os.FileMode(0o644)
		if st, serr := os.Stat(p.Abs); serr == nil {
			mode = st.Mode().Perm()
		}
		if werr := os.WriteFile(p.Abs, p.Result, mode); werr != nil {
			rollback()
			return fail("the write of "+p.Rel+" failed, and what was written has been put back", werr)
		}
		written = append(written, p)
		mf := aiderMFile{File: p.Rel, New: p.New, Eol: aiderEolName(p.Eol), After: aiderSum(p.Result)}
		if !p.New {
			mf.Before = aiderSum(p.Before)
		}
		m.Files = append(m.Files, mf)
		m.Line = p.Branch
	}
	mb, _ := json.MarshalIndent(m, "", " ")
	if werr := aiderWriteFile(filepath.Join(runDir, "manifest.json"), mb); werr != nil {
		rollback()
		return fail("the run's record could not be kept, so what was written has been put back", werr)
	}
	ans.Line = m.Line
	ans.Note = "Written on the line of work `" + m.Line + "`; unsaved until git_commit saves it there. The suites are his to run, and the main line moves only by his hand."
	aiderRecord(args, t, c, "aider_run", "aider_ran", note(fmt.Sprintf("wrote %d file(s) on %s", len(changed), m.Line)))
	aiderPrune(pl)
	return ans.String(), nil
}

// aiderPrune keeps the newest aiderKeepRuns run folders. The names are matched against the id shape first,
// so nothing that is not a run's own folder is ever removed.
func aiderPrune(pl aiderPlace) {
	entries, err := os.ReadDir(pl.Runs)
	if err != nil {
		return
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() && aiderRunRe.MatchString(e.Name()) {
			ids = append(ids, e.Name())
		}
	}
	sort.Strings(ids)
	for len(ids) > aiderKeepRuns {
		_ = os.RemoveAll(filepath.Join(pl.Runs, ids[0]))
		ids = ids[1:]
	}
}

// ---- aider_undo ----------------------------------------------------------------------------------

func toolAiderUndo(t tenant.Tenant, args map[string]any) (string, error) {
	c := callerOf(args)
	if !c.Service {
		return "", errors.New("refused: Aider is the operator's own hand and is offered to his glass alone")
	}
	ans := aiderAnswer{}
	id := strings.TrimSpace(str(args, "run"))
	refuse := func(why ...string) (string, error) {
		ans.State, ans.Class, ans.Why = "refused", string(ShellRefused), why
		aiderRecord(args, t, c, "aider_undo", "aider_refused", "undo "+clip(id, 40)+" :: "+strings.Join(why, "; "))
		return ans.String(), nil
	}
	if !aiderRunRe.MatchString(id) {
		return refuse("name the run to take back, as the run's own id (a1234...); got " + strconv.Quote(clip(id, 40)))
	}
	ans.Run = id
	pl := aiderPlaceOf(t.Home)
	mpath := filepath.Join(pl.Runs, id, "manifest.json")
	raw, err := os.ReadFile(mpath)
	if err != nil {
		return refuse("no record of run " + id + " (a run's record is kept for the last " + strconv.Itoa(aiderKeepRuns) + " runs; an older one is git's to take back)")
	}
	var m aiderManifest
	if err := json.Unmarshal(raw, &m); err != nil || m.Run != id {
		return refuse("the record of run " + id + " cannot be read")
	}
	if m.Undone != "" {
		return refuse("run " + id + " was already taken back at " + m.Undone)
	}
	branches := map[string]string{}
	type undo struct {
		f       aiderMFile
		abs     string
		restore []byte
	}
	var todo []undo
	for _, f := range m.Files {
		if why := aiderJudgePath(t, f.File, true); why != "" {
			return refuse(why)
		}
		abs := aiderAbs(t, gopath.Clean(f.File))
		if _, why := aiderLineRefusal(t.Home, abs, branches); why != "" {
			return refuse(f.File + ": " + why)
		}
		now, rerr := os.ReadFile(abs)
		if rerr != nil || aiderSum(now) != f.After {
			return refuse(f.File + " is not as Aider left it (it has changed since, or is gone), so taking the run back would lose what came after; take back the later run first, or use git")
		}
		u := undo{f: f, abs: abs}
		if !f.New {
			b, berr := os.ReadFile(filepath.Join(pl.Runs, id, "before", filepath.FromSlash(f.File)))
			if berr != nil || aiderSum(b) != f.Before {
				return refuse("the copy of " + f.File + " as it was, kept beside run " + id + ", is missing or does not match its record")
			}
			u.restore = b
		}
		todo = append(todo, u)
	}
	for _, u := range todo {
		if u.f.New {
			if err := os.Remove(u.abs); err != nil {
				return refuse("could not remove " + u.f.File + ": " + err.Error())
			}
			ans.Changed = append(ans.Changed, aiderFile{File: u.f.File, New: true})
			continue
		}
		if err := os.WriteFile(u.abs, u.restore, 0o644); err != nil {
			return refuse("could not put " + u.f.File + " back: " + err.Error())
		}
		ans.Changed = append(ans.Changed, aiderFile{File: u.f.File})
	}
	m.Undone = time.Now().UTC().Format(time.RFC3339)
	mb, _ := json.MarshalIndent(m, "", " ")
	_ = os.WriteFile(mpath, mb, 0o644)
	ans.State = "undone"
	ans.Line = m.Line
	ans.Note = "Run " + id + " is taken back: the files are as they were before Aider. Nothing was saved or sent."
	aiderRecord(args, t, c, "aider_undo", "aider_undone", "undo "+id)
	return ans.String(), nil
}

// ---- aider_status --------------------------------------------------------------------------------

type aiderRunInfo struct {
	Run     string   `json:"run"`
	When    string   `json:"when"`
	Line    string   `json:"line,omitempty"`
	Files   []string `json:"files"`
	Message string   `json:"message,omitempty"`
	Undone  bool     `json:"undone,omitempty"`
}

type aiderStatus struct {
	Installed bool           `json:"installed"`
	Version   string         `json:"version,omitempty"`
	Folder    string         `json:"folder"`
	Model     string         `json:"model,omitempty"`
	Context   int            `json:"context,omitempty"`
	Rack      string         `json:"rack"`
	OnRack    bool           `json:"on_rack"`
	Busy      bool           `json:"busy"`
	Line      string         `json:"line,omitempty"`
	OnLine    bool           `json:"on_line"`
	MaxEdit   int            `json:"max_files"`
	MaxRead   int            `json:"max_context"`
	MaxKB     int            `json:"budget_kb,omitempty"`
	Runs      []aiderRunInfo `json:"runs"`
	Why       []string       `json:"why,omitempty"`
}

// toolAiderStatus is the reading side, for the tab and the Inspector: whether Aider can run in this world
// now, with what, and what the last runs were. It starts no process and writes nothing.
func toolAiderStatus(t tenant.Tenant, args map[string]any) (string, error) {
	pl := aiderPlaceOf(t.Home)
	st := aiderStatus{Folder: aiderDirName + "/", Rack: "silent", MaxEdit: aiderMaxEdit, MaxRead: aiderMaxRead, Runs: []aiderRunInfo{}}
	why := func(s string) { st.Why = append(st.Why, s) }
	if st.Installed = aiderInstalled(pl); st.Installed {
		st.Version = aiderVersion(pl)
	} else {
		why("Aider is not installed in this world: there is no " + aiderDirName + "/venv here")
	}
	if seat, err := aiderCoder(t.Home); err != nil {
		why(err.Error())
	} else {
		st.Model, st.Context = seat.Model, seat.Ctx
		st.MaxKB = ((seat.Ctx-aiderPromptTokens)*aiderCharsPerToken + 1023) / 1024
		if _, voices, rerr := aiderRack(); rerr != nil {
			why("the rack is not reachable: " + clip(rerr.Error(), 160))
		} else {
			st.Rack = "up"
			if st.OnRack = aiderHasModel(voices, seat.Model); !st.OnRack {
				why(fmt.Sprintf("the coding seat's model %q is not on this machine's rack", seat.Model))
			}
		}
	}
	if aiderMu.TryLock() {
		aiderMu.Unlock()
	} else {
		st.Busy = true
	}
	if repo := aiderRepoOf(t.Home, filepath.Join(t.Home, "x")); repo != "" {
		st.Line = aiderBranch(repo)
		st.OnLine = aiderIsLine(st.Line)
	}
	if entries, err := os.ReadDir(pl.Runs); err == nil {
		var ids []string
		for _, e := range entries {
			if e.IsDir() && aiderRunRe.MatchString(e.Name()) {
				ids = append(ids, e.Name())
			}
		}
		sort.Sort(sort.Reverse(sort.StringSlice(ids)))
		for _, id := range ids {
			raw, rerr := os.ReadFile(filepath.Join(pl.Runs, id, "manifest.json"))
			if rerr != nil {
				continue
			}
			var m aiderManifest
			if json.Unmarshal(raw, &m) != nil {
				continue
			}
			ri := aiderRunInfo{Run: m.Run, When: m.When, Line: m.Line, Message: clip(m.Message, 120), Undone: m.Undone != "", Files: []string{}}
			for _, f := range m.Files {
				ri.Files = append(ri.Files, f.File)
			}
			st.Runs = append(st.Runs, ri)
			if len(st.Runs) >= 8 {
				break
			}
		}
	}
	b, err := json.Marshal(st)
	return string(b), err
}
