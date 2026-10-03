package tools

// THE SHELL'S GATE (WHAT'S LEFT H15, 2026-10-03; the operator's word of 2026-10-02, on the card:
// "typed by you, gated").
//
// WHAT A TYPED ENTRY IS, decided before anything runs and by reading it. Never by asking a model
// (nothing here is a seat), and never by trusting what the text says it will do.
//
//	read     a plain look -- list, cat, grep, git status/log/diff. Runs at once.
//	write    anything else: it writes, deletes, installs, reaches the network, saves or sends work,
//	         runs a program the gate cannot read, or uses shell syntax the gate does not model.
//	         It waits for his approval, and the card says why.
//	refused  named for what it reaches -- a secret, a path outside the ground, client material, a
//	         key typed into a command. Approval cannot lift a refusal.
//
// THE DEFAULT IS THE CARD. The read list is short and exact and everything not on it is "write", so
// a command the gate has never heard of is asked about, not run. The cost of a wrong guess in that
// direction is one click; in the other it is a deleted file.
//
// WHERE IT IS MIRRORED. The refusals are the law gate's own checks restated for text that is code
// rather than prose (manjuel/lawgate.py: the reach, the secret, the client tag), and they are held
// to the matrix in shellgate_test.go. A command word dressed in quotes or backslashes is still the
// command word (the lexer unquotes before it looks), and every unquoted word is checked again with
// its quoting stripped, so `.e""nv` is `.env`.
//
// WHAT IT DOES NOT DO. It does not prove a "read" cannot hurt. `grep -r` reads every file under a
// folder, hidden ones included; what stands behind it is the OUTPUT being scrubbed of every secret
// value the ground holds (shell.go redactor), and the child's environment holding none.

import (
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// ShellClass is what the gate says a typed entry is.
type ShellClass string

const (
	ShellRead    ShellClass = "read"
	ShellWrite   ShellClass = "write"
	ShellRefused ShellClass = "refused"
)

// ShellVerdict is the gate's whole answer for one entry: the class and the reasons for it. A read
// carries none; a write names what made it one; a refusal names only what was refused.
type ShellVerdict struct {
	Class ShellClass
	Why   []string
}

type shellGate struct {
	refusals []string
	writes   []string
}

func addOnce(list *[]string, s string) {
	for _, have := range *list {
		if have == s {
			return
		}
	}
	if len(*list) < 8 {
		*list = append(*list, s)
	}
}

func (g *shellGate) refuse(s string) { addOnce(&g.refusals, s) }
func (g *shellGate) write(s string)  { addOnce(&g.writes, s) }

// refusedFor is whether a refusal that says this has been made already.
func (g *shellGate) refusedFor(frag string) bool {
	for _, r := range g.refusals {
		if strings.Contains(r, frag) {
			return true
		}
	}
	return false
}

func (g *shellGate) verdict() ShellVerdict {
	switch {
	case len(g.refusals) > 0:
		return ShellVerdict{ShellRefused, g.refusals}
	case len(g.writes) > 0:
		return ShellVerdict{ShellWrite, g.writes}
	}
	return ShellVerdict{ShellRead, nil}
}

// ---- refused by name ----------------------------------------------------------------------------

// shQ are the marks that open a quoted word: both quotes and the backtick.
const shQ = "\"'`"

var (
	// A key typed into a command: the door's own key shape, and the shapes of the vendors' (the same
	// ones guard.Redact strips). RULE 7: a key never rides a command line, where a process list reads it.
	shKeyRe = regexp.MustCompile(`atl_[0-9a-f]{32}|\bsk-[A-Za-z0-9_-]{16,}|\bBearer\s+[A-Za-z0-9._~+/-]{16,}|` +
		`\bxox[bap]-[A-Za-z0-9-]{10,}|\bghp_[A-Za-z0-9]{20,}|\bAKIA[0-9A-Z]{16}\b`)
	// is_secret(), restated (manjuel/vectors.py): any name that STARTS with .env; the key files by
	// name; and the generic names only when they are a file (a data extension), so that a variable
	// called `token` or the word in a grep pattern is not a secret.
	shEnvFileRe = regexp.MustCompile(`(?i)(?:^|[^\w])\.env[\w.-]*`)
	shKeyFileRe = regexp.MustCompile(`(?i)(?:^|[^\w])(?:id_rsa|id_ed25519|\.netrc|\.npmrc|\.pypirc)(?:$|[^\w])`)
	shStemRe    = regexp.MustCompile(`(?i)(?:^|[^\w])(?:secrets?|credentials?|creds|tokens?|keyfile|apikey|api_key|` +
		`[\w-]*_secrets|[\w-]*_key)\.(?:json|ya?ml|toml|ini|cfg|conf|pem|key|crt|p12|pfx|txt|csv|db|sqlite|env|enc|gpg|asc)\b`)
	// SITTING LAW 2: client material by tag (vectors.is_protected).
	shClientRe = regexp.MustCompile(`(?i)(?:^|[\\/\s` + shQ + `])vault[\\/]|\.client\.`)
	// A home directory by name: it is not the ground.
	shHomeVarRe = regexp.MustCompile(`(?i)\$\{?(?:HOME|USERPROFILE|APPDATA|LOCALAPPDATA|HOMEPATH|HOMEDRIVE)\b|` +
		`%(?:USERPROFILE|APPDATA|LOCALAPPDATA|HOMEPATH|HOMEDRIVE|HOME)%`)
	// A reach outside the ground (lawgate._REACH_RE, widened for code): a parent step, a drive, a
	// rooted path to a place that is not the ground, a home, a UNC share. The ground's own path is
	// allowed, and checked after, by insideOwn.
	shReachRe = regexp.MustCompile(`(?i)(?:^|[\s` + shQ + `(=])(?:\.[\\/])*(?:` +
		`\.\.(?:[\\/]+|$|[\s` + shQ + `);|&])|` +
		`[A-Za-z]:[\\/]|` +
		`~(?:[\\/]|$|[\s` + shQ + `)])|` +
		`\\\\[A-Za-z0-9_.-]+[\\/]|` +
		`[\\/]+(?:(?:home|users|etc|tmp|var|windows|usr|bin|sbin|lib|dev|proc|sys|root|opt|mnt)(?:[\\/]|$|[\s` + shQ + `);|&,])|[a-z][\\/]))`)
)

// ownSpellings are the ways the ground's own path is written: Windows, forward-slashed, and the
// MSYS spelling Git Bash prints (/c/Users/...). Lower-cased, since Windows is.
func ownSpellings(home string) []string {
	s := strings.ToLower(strings.TrimRight(filepath.ToSlash(home), "/"))
	out := []string{s}
	if len(s) >= 3 && s[1] == ':' && s[2] == '/' {
		out = append(out, "/"+s[:1]+s[2:])
	}
	return out
}

// insideOwn is whether a path, once its dots are resolved, is the ground or something under it.
func insideOwn(rest string, own []string) bool {
	tok := rest
	if i := strings.IndexAny(tok, " \t\r\n"+shQ+";|&<>()"); i >= 0 {
		tok = tok[:i]
	}
	cl := path.Clean(tok)
	for _, o := range own {
		if cl == o || strings.HasPrefix(cl, o+"/") {
			return true
		}
	}
	return false
}

func reachesOutside(text, home string) (string, bool) {
	t := strings.ReplaceAll(text, "/dev/null", " ")
	own := ownSpellings(home)
	for _, loc := range shReachRe.FindAllStringIndex(t, -1) {
		start := loc[0]
		for start < loc[1] && strings.ContainsRune(" \t\r\n"+shQ+"(=", rune(t[start])) {
			start++
		}
		rest := strings.ToLower(strings.ReplaceAll(t[start:], "\\", "/"))
		if insideOwn(rest, own) {
			continue
		}
		hit := t[start:loc[1]]
		if len(hit) > 24 {
			hit = hit[:24]
		}
		return strings.TrimSpace(hit), true
	}
	return "", false
}

// shellRefusals names, for a typed text, what is refused outright. It never echoes a key.
func shellRefusals(text, home string) []string {
	var out []string
	if strings.ContainsRune(text, 0) {
		addOnce(&out, "a NUL byte in the entry")
	}
	if shKeyRe.MatchString(text) {
		addOnce(&out, "RULE 7: a key typed into a command; keys are silent and never ride a command line")
	}
	if m := shEnvFileRe.FindString(text); m != "" {
		addOnce(&out, "LAW 9 / RULE 7: names a secret file ("+strings.TrimSpace(m)+"); keys are silent")
	} else if m := shKeyFileRe.FindString(text); m != "" {
		addOnce(&out, "LAW 9 / RULE 7: names a secret file ("+strings.Trim(m, " \t\"'`/\\=:(,")+"); keys are silent")
	} else if m := shStemRe.FindString(text); m != "" {
		addOnce(&out, "LAW 9 / RULE 7: names a secret file ("+strings.Trim(m, " \t\"'`/\\=:(,")+"); keys are silent")
	}
	if shClientRe.MatchString(text) {
		addOnce(&out, "SITTING LAW 2: names client material by tag; it is opened only when the operator points at it, by hand")
	}
	if shHomeVarRe.MatchString(text) {
		addOnce(&out, "RULE 1: names a home directory; every read and write stays inside the ground")
	}
	if hit, ok := reachesOutside(text, home); ok {
		addOnce(&out, "RULE 1 / LAW 8: names a path outside the ground ("+hit+"...); every read and write stays inside Research")
	}
	return out
}

// ---- a small bash lexer -------------------------------------------------------------------------

type shKind int

const (
	shWordK  shKind = iota
	shOpK           // ; && || |  (a newline is ;)
	shRedirK        // > >> >| &> &>> <  and descriptor duplications (2>&1)
)

type shTok struct {
	kind    shKind
	text    string // the word with its quoting stripped, the operator, or the redirection operator
	fd      string // a redirection's descriptor prefix
	dup     bool   // a redirection that duplicates a descriptor (2>&1), and so names no file
	glob    bool   // an unquoted * ? [ in the word
	brace   bool   // an unquoted { or } in the word
	expands bool   // a $variable the gate cannot read
	tilde   bool   // starts with an unquoted ~
}

func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isVarStart(n rune) bool {
	return n == '{' || n == '_' || unicode.IsLetter(n) || unicode.IsDigit(n) || strings.ContainsRune("?$!@*#-", n)
}

// shLex reads a command line the way bash would split it, far enough to see the command words, and
// reports every construct it does not model. It never evaluates anything.
func shLex(s string) ([]shTok, []string) {
	rs := []rune(s)
	var toks []shTok
	var bad []string
	unsupported := func(m string) { addOnce(&bad, m) }
	var cur strings.Builder
	var w shTok
	have := false
	flush := func() {
		if have {
			w.kind = shWordK
			w.text = cur.String()
			toks = append(toks, w)
		}
		cur.Reset()
		w = shTok{}
		have = false
	}
	peek := func(i int) rune {
		if i >= 0 && i < len(rs) {
			return rs[i]
		}
		return 0
	}
	i := 0
	for i < len(rs) {
		c := rs[i]
		switch {
		case c == ' ' || c == '\t' || c == '\r':
			flush()
			i++
		case c == '\n':
			flush()
			toks = append(toks, shTok{kind: shOpK, text: ";"})
			i++
		case c == '#' && !have:
			for i < len(rs) && rs[i] != '\n' {
				i++
			}
		case c == '\'':
			j := i + 1
			for j < len(rs) && rs[j] != '\'' {
				j++
			}
			if j >= len(rs) {
				unsupported("an unterminated quote")
				i = len(rs)
				break
			}
			cur.WriteString(string(rs[i+1 : j]))
			have = true
			i = j + 1
		case c == '"':
			have = true
			i++
			closed := false
			for i < len(rs) {
				d := rs[i]
				if d == '"' {
					closed = true
					i++
					break
				}
				if d == '\\' && i+1 < len(rs) {
					n := rs[i+1]
					switch n {
					case '$', '`', '"', '\\':
						cur.WriteRune(n)
						i += 2
						continue
					case '\n':
						i += 2
						continue
					}
					cur.WriteRune(d)
					i++
					continue
				}
				if d == '`' {
					unsupported("command substitution (backticks)")
				}
				if d == '$' {
					n := peek(i + 1)
					if n == '(' {
						unsupported("command substitution $(...)")
					}
					if isVarStart(n) {
						w.expands = true
					}
				}
				cur.WriteRune(d)
				i++
			}
			if !closed {
				unsupported("an unterminated quote")
			}
		case c == '\\':
			have = true
			if i+1 < len(rs) {
				if rs[i+1] == '\n' {
					i += 2
					continue
				}
				cur.WriteRune(rs[i+1])
				i += 2
			} else {
				i++
			}
		case c == '$':
			n := peek(i + 1)
			switch {
			case n == '(':
				unsupported("command substitution $(...)")
			case n == '\'':
				unsupported("$'...' quoting")
			case isVarStart(n):
				w.expands = true
			}
			cur.WriteRune(c)
			have = true
			i++
		case c == '`':
			unsupported("command substitution (backticks)")
			cur.WriteRune(c)
			have = true
			i++
		case c == '(' || c == ')':
			unsupported("a subshell or group ( )")
			flush()
			i++
		case c == '{' || c == '}':
			w.brace = true
			cur.WriteRune(c)
			have = true
			i++
		case c == '*' || c == '?' || c == '[':
			w.glob = true
			cur.WriteRune(c)
			have = true
			i++
		case c == '~' && !have:
			w.tilde = true
			cur.WriteRune(c)
			have = true
			i++
		case c == ';':
			flush()
			if peek(i+1) == ';' {
				unsupported("a case terminator ;;")
				i += 2
			} else {
				toks = append(toks, shTok{kind: shOpK, text: ";"})
				i++
			}
		case c == '&':
			flush()
			switch peek(i + 1) {
			case '&':
				toks = append(toks, shTok{kind: shOpK, text: "&&"})
				i += 2
			case '>':
				op, j := "&>", i+2
				if peek(j) == '>' {
					op, j = "&>>", j+1
				}
				toks = append(toks, shTok{kind: shRedirK, text: op})
				i = j
			default:
				unsupported("a background job (&)")
				i++
			}
		case c == '|':
			flush()
			switch peek(i + 1) {
			case '|':
				toks = append(toks, shTok{kind: shOpK, text: "||"})
				i += 2
			case '&':
				toks = append(toks, shTok{kind: shOpK, text: "|"})
				i += 2
			default:
				toks = append(toks, shTok{kind: shOpK, text: "|"})
				i++
			}
		case c == '<' || c == '>':
			fd := ""
			if have && allDigits(cur.String()) && !w.glob && !w.brace && !w.expands {
				fd = cur.String()
				cur.Reset()
				w = shTok{}
				have = false
			} else {
				flush()
			}
			op := string(c)
			j := i + 1
			dup := false
			if c == '>' {
				switch peek(j) {
				case '>':
					op, j = ">>", j+1
				case '|':
					op, j = ">|", j+1
				case '(':
					unsupported("process substitution >( )")
				case '&':
					k := j + 1
					for k < len(rs) && unicode.IsDigit(rs[k]) {
						k++
					}
					if k > j+1 {
						op, dup, j = ">&", true, k
					} else if peek(j+1) == '-' {
						op, dup, j = ">&", true, j+2
					} else {
						op, j = ">&", j+1
					}
				}
			} else {
				switch peek(j) {
				case '<':
					if peek(j+1) == '<' {
						unsupported("a here-string (<<<)")
						j += 2
					} else {
						unsupported("a here-document (<<)")
						j++
					}
				case '(':
					unsupported("process substitution <( )")
				case '>':
					unsupported("a read-write redirection (<>)")
					j++
				case '&':
					k := j + 1
					for k < len(rs) && unicode.IsDigit(rs[k]) {
						k++
					}
					op, dup, j = "<&", true, k
				}
			}
			toks = append(toks, shTok{kind: shRedirK, text: op, fd: fd, dup: dup})
			i = j
		default:
			cur.WriteRune(c)
			have = true
			i++
		}
	}
	flush()
	return toks, bad
}

// ---- the read list -------------------------------------------------------------------------------

// denyFlags makes a rule that is a plain look unless one of these options is given. Short options may
// be clustered (-ro), long ones may carry =value.
func denyFlags(deny ...string) func([]string) string {
	set := map[string]bool{}
	for _, d := range deny {
		set[d] = true
	}
	return func(args []string) string {
		for _, a := range args {
			if a == "--" {
				break
			}
			switch {
			case strings.HasPrefix(a, "--"):
				name := a
				if i := strings.Index(a, "="); i >= 0 {
					name = a[:i]
				}
				if set[name] {
					return "uses the option " + name + ", which writes or runs something"
				}
			case strings.HasPrefix(a, "-") && len(a) > 1:
				for _, ch := range a[1:] {
					if set["-"+string(ch)] {
						return "uses the option -" + string(ch) + ", which writes or runs something"
					}
				}
			}
		}
		return ""
	}
}

func operandsOf(args []string) []string {
	var out []string
	after := false
	for _, a := range args {
		if !after && a == "--" {
			after = true
			continue
		}
		if !after && strings.HasPrefix(a, "-") && len(a) > 1 {
			continue
		}
		out = append(out, a)
	}
	return out
}

func findRule(args []string) string {
	for _, a := range args {
		switch a {
		case "-delete", "-exec", "-execdir", "-ok", "-okdir", "-fprint", "-fprint0", "-fprintf", "-fls":
			return "uses find " + a + ", which deletes, runs a program or writes a file"
		}
	}
	return ""
}

func uniqRule(args []string) string {
	if len(operandsOf(args)) > 1 {
		return "names a second file, which uniq writes"
	}
	return ""
}

func dateRule(args []string) string {
	if r := denyFlags("-s", "--set")(args); r != "" {
		return "sets the clock"
	}
	for _, o := range operandsOf(args) {
		if !strings.HasPrefix(o, "+") {
			return "sets the clock"
		}
	}
	return ""
}

func noOperands(what string) func([]string) string {
	return func(args []string) string {
		if len(operandsOf(args)) > 0 {
			return what
		}
		return ""
	}
}

// gitRule judges `git ...`: a few subcommands that only read, and the options of those that make them
// write or run a program. Everything else is the repository changing or the network, and is a card. Its
// reasons begin with "git", so the card reads as a sentence.
func gitRule(args []string) string {
	i := 0
	for i < len(args) && strings.HasPrefix(args[i], "-") {
		switch args[i] {
		case "--no-pager", "--no-optional-locks", "--no-replace-objects":
			i++
		default:
			return "git option " + args[i] + " changes what git runs or where"
		}
	}
	if i >= len(args) {
		return ""
	}
	sub, rest := args[i], args[i+1:]
	readsOnly := denyFlags("--output", "--ext-diff", "--open-files-in-pager", "-O")
	switch sub {
	case "status", "log", "diff", "show", "ls-files", "ls-tree", "rev-parse", "rev-list", "describe",
		"blame", "shortlog", "grep", "cat-file", "show-ref", "for-each-ref", "name-rev", "merge-base",
		"check-ignore", "check-attr", "diff-tree", "diff-index", "diff-files", "count-objects",
		"whatchanged", "version", "var", "verify-commit", "verify-tag", "annotate":
		if why := readsOnly(rest); why != "" {
			return "git " + sub + " " + why
		}
		return ""
	case "branch":
		list := false
		for _, a := range rest {
			if a == "--list" || a == "-l" {
				list = true
			}
		}
		if r := denyFlags("-d", "-D", "-m", "-M", "-c", "-C", "-u", "-f", "-t", "--delete", "--move", "--copy",
			"--set-upstream-to", "--unset-upstream", "--edit-description", "--force", "--track", "--no-track")(rest); r != "" {
			return "git branch changes the repository (" + r + ")"
		}
		if len(operandsOf(rest)) > 0 && !list {
			return "git branch with a name changes the repository (it makes a branch)"
		}
		return ""
	case "tag":
		list := false
		for _, a := range rest {
			if a == "--list" || a == "-l" {
				list = true
			}
		}
		if len(rest) == 0 || list {
			if r := denyFlags("-d", "-f", "-s", "-u", "-a", "-m", "--delete", "--force", "--sign")(rest); r != "" {
				return "git tag changes the repository (" + r + ")"
			}
			return ""
		}
		return "git tag with a name or an option changes the repository (it makes or removes a tag)"
	case "remote":
		for _, a := range rest {
			if a != "-v" && a != "--verbose" {
				return "git remote " + a + " changes the repository's remotes or contacts one"
			}
		}
		return ""
	case "stash":
		if len(rest) > 0 && (rest[0] == "list" || rest[0] == "show") {
			return ""
		}
		return "git stash changes the repository (the working tree)"
	case "worktree":
		if len(rest) > 0 && rest[0] == "list" {
			return ""
		}
		return "git worktree changes the repository"
	case "reflog":
		if len(rest) == 0 || rest[0] == "show" {
			return ""
		}
		return "git reflog " + rest[0] + " changes the repository (it rewrites the log)"
	case "config":
		for _, a := range rest {
			switch a {
			case "--get", "--get-all", "--get-regexp", "--list", "-l", "--show-origin":
				return ""
			}
		}
		return "git config changes the repository (it writes the configuration)"
	case "push":
		return "git push sends saved work to the remote"
	case "pull", "fetch", "clone", "ls-remote":
		return "git " + sub + " reaches the network"
	case "commit", "merge", "rebase", "reset", "checkout", "switch", "restore", "add", "mv", "rm", "cherry-pick",
		"revert", "am", "apply", "init", "clean", "gc", "prune":
		return "git " + sub + " changes the repository"
	}
	return "git " + sub + " is not one of the plain looks"
}

// shReadCommands is the whole list of what runs at once. A nil rule means the command is a plain look
// with any options; a rule returns the reason it is not one, or "".
var shReadCommands = map[string]func([]string) string{
	"ls": nil, "dir": nil, "pwd": nil, "wc": nil, "nl": nil, "stat": nil, "du": nil, "df": nil,
	"basename": nil, "dirname": nil, "realpath": nil, "readlink": nil, "which": nil, "type": nil,
	"whoami": nil, "uname": nil, "echo": nil, "printf": nil, "true": nil, "false": nil, ":": nil,
	"cut": nil, "tr": nil, "cat": nil, "diff": nil, "cmp": nil, "grep": nil, "egrep": nil, "fgrep": nil,
	"sha256sum": nil, "sha1sum": nil, "md5sum": nil, "test": nil,
	"file":     denyFlags("-C", "--compile"),
	"head":     denyFlags("-f", "-F", "--follow"),
	"tail":     denyFlags("-f", "-F", "--follow"),
	"sort":     denyFlags("-o", "--output", "--compress-program"),
	"uniq":     uniqRule,
	"date":     dateRule,
	"hostname": noOperands("sets the host name"),
	"find":     findRule,
	"rg":       denyFlags("--pre", "--pre-glob", "--hostname-bin", "-z", "--search-zip"),
	"git":      gitRule,
}

var (
	shDeletes  = setOf("rm", "rmdir", "del", "erase", "rd", "unlink", "shred")
	shInstalls = setOf("pip", "pip3", "pipx", "npm", "npx", "pnpm", "yarn", "bun", "choco", "winget", "scoop",
		"apt", "apt-get", "brew", "gem", "conda", "uv")
	shNetwork = setOf("curl", "wget", "ssh", "scp", "sftp", "ftp", "nc", "ncat", "netcat", "telnet", "ping",
		"nslookup", "dig", "tracert", "traceroute", "rsync")
	shWrites = setOf("touch", "mkdir", "cp", "mv", "ln", "tee", "sed", "awk", "chmod", "chown", "truncate", "dd",
		"install", "tar", "zip", "unzip", "7z", "patch", "xargs")
	shRuns = setOf("python", "python3", "py", "node", "deno", "bash", "sh", "zsh", "pwsh", "powershell", "cmd",
		"make", "go", "cargo", "dotnet", "java", "ruby", "perl", "php", "mvn", "gradle", "eval", "exec", "source",
		".", "env", "command", "sudo", "nohup", "time", "nice", "start")
	shReserved = setOf("if", "then", "else", "elif", "fi", "for", "while", "until", "do", "done", "case", "esac",
		"function", "select", "!", "coproc", "[[", "{", "}")
)

func setOf(names ...string) map[string]bool {
	m := make(map[string]bool, len(names))
	for _, n := range names {
		m[n] = true
	}
	return m
}

// whyNotRead is the card's sentence for a command that is not on the read list.
func whyNotRead(cmd string, args []string) string {
	switch {
	case shDeletes[cmd]:
		return "deletes (" + cmd + ")"
	case shInstalls[cmd]:
		return "installs software (" + cmd + ")"
	case shNetwork[cmd]:
		return "reaches the network (" + cmd + ")"
	case shWrites[cmd]:
		return "writes files (" + cmd + ")"
	case shRuns[cmd]:
		return "runs a program the gate cannot read (" + cmd + ")"
	}
	return "runs " + cmd + ", which the gate cannot read as a plain look"
}

func globReachesHidden(w string) bool {
	for _, seg := range strings.FieldsFunc(w, func(r rune) bool { return r == '/' || r == '\\' }) {
		if (strings.HasPrefix(seg, ".") || strings.HasPrefix(seg, "[")) && strings.ContainsAny(seg, "*?[") {
			return true
		}
	}
	return false
}

var shAssignRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)

// JudgeBash is the gate for one bash entry typed in the ground at home.
func JudgeBash(cmd, home string) ShellVerdict {
	g := &shellGate{}
	for _, r := range shellRefusals(cmd, home) {
		g.refuse(r)
	}
	toks, bad := shLex(cmd)
	for _, b := range bad {
		g.write("uses shell syntax the gate does not read: " + b)
	}
	// Every unquoted word is looked at again with its quoting stripped, so a name built out of quotes
	// and backslashes is the name it spells.
	for _, t := range toks {
		if t.kind != shWordK {
			continue
		}
		for _, r := range shellRefusals(t.text, home) {
			g.refuse(r)
		}
		if t.tilde {
			g.refuse("RULE 1: ~ names a home directory; every read and write stays inside the ground")
		}
		if t.glob && globReachesHidden(t.text) {
			g.refuse("LAW 9 / RULE 7: a pattern that can reach hidden files, where the secrets live")
		}
		if t.expands {
			g.write("uses a shell variable ($...) whose value the gate cannot read")
		}
		if t.brace {
			g.write("uses a brace group or brace expansion")
		}
		// The text of the whole entry has been read for a reach already; a word is read again here
		// for what only a word shows (`a/../../x`), and a refusal for it is said once.
		if !strings.HasPrefix(t.text, "-") && !g.refusedFor("outside the ground") {
			if p := path.Clean(strings.ReplaceAll(t.text, "\\", "/")); p == ".." || strings.HasPrefix(p, "../") {
				g.refuse("RULE 1 / LAW 8: names a path outside the ground, above it (" + clip(t.text, 24) + "); every read and write stays inside Research")
			}
			// A word that starts at the root names a place, and only the ground's own is lawful.
			if (strings.HasPrefix(t.text, "/") || strings.HasPrefix(t.text, "\\")) && t.text != "/dev/null" {
				rest := strings.ToLower(strings.ReplaceAll(t.text, "\\", "/"))
				if !insideOwn(rest, ownSpellings(home)) {
					g.refuse("RULE 1 / LAW 8: names a rooted path outside the ground (" + clip(t.text, 24) + "); every read and write stays inside Research")
				}
			}
		}
	}

	var words []shTok
	type redirect struct {
		tok    shTok
		target string
		have   bool
	}
	var redirs []redirect
	judgeSimple := func() {
		defer func() { words, redirs = nil, nil }()
		for _, r := range redirs {
			if r.tok.dup {
				continue
			}
			if !r.have {
				g.write("a redirection with no file after it")
				continue
			}
			if r.tok.text == "<" {
				continue // reading a file: its name was checked with every other word
			}
			if r.target != "/dev/null" {
				g.write("redirects output into a file (" + clip(r.target, 40) + ")")
			}
		}
		if len(words) == 0 {
			if len(redirs) > 0 {
				g.write("a redirection with no command")
			}
			return
		}
		first := words[0]
		args := make([]string, 0, len(words)-1)
		for _, a := range words[1:] {
			args = append(args, a.text)
		}
		cmdWord := first.text
		switch {
		case shReserved[cmdWord]:
			g.write("uses the shell keyword " + cmdWord)
		case first.expands || first.glob || first.brace || first.tilde:
			g.write("the command word is built by the shell (" + clip(cmdWord, 30) + ")")
		case shAssignRe.MatchString(cmdWord):
			g.write("sets a variable before the command")
		case strings.ContainsAny(cmdWord, `/\`):
			g.write("runs a program by path (" + clip(cmdWord, 30) + ")")
		default:
			rule, listed := shReadCommands[cmdWord]
			if !listed {
				g.write(whyNotRead(cmdWord, args))
			} else if rule != nil {
				if why := rule(args); why != "" {
					if !strings.HasPrefix(why, cmdWord+" ") {
						why = cmdWord + " " + why
					}
					g.write(why)
				}
			}
		}
	}
	for i := 0; i < len(toks); i++ {
		t := toks[i]
		switch t.kind {
		case shWordK:
			words = append(words, t)
		case shRedirK:
			r := redirect{tok: t}
			if !t.dup && i+1 < len(toks) && toks[i+1].kind == shWordK {
				r.target, r.have = toks[i+1].text, true
				i++
			}
			redirs = append(redirs, r)
		case shOpK:
			if len(words) == 0 && len(redirs) == 0 && t.text != ";" {
				g.write("an operator (" + t.text + ") with no command before it")
			}
			judgeSimple()
		}
	}
	judgeSimple()
	return g.verdict()
}

func clip(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	return string([]rune(s)[:n]) + "..."
}

// JudgePythonText is the half of the Python gate that is text: what is refused by name. Whether the
// code is a plain calculation or something to ask about is the session's to say (shellpy.go), because
// only Python can read Python.
func JudgePythonText(code, home string) ShellVerdict {
	g := &shellGate{}
	for _, r := range shellRefusals(code, home) {
		g.refuse(r)
	}
	return g.verdict()
}
