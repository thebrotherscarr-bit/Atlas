package tools

// THE SHELL GATE'S MATRIX (WHAT'S LEFT H15, 2026-10-03), in the manner of the law gate's own phrasing
// matrix (tests/matrix.py, D11): each entry says what the gate must call a typed command, and then the
// same command is said several more ways -- padded, tabbed, followed by a comment, chained behind a
// harmless look, piped into one, its stderr folded in -- and none of the ways is allowed to come out
// differently from the plain one. A gate that reads the plain phrasing and not the padded one is a gate
// for people who type carefully.
//
// The canonical table is the SPECIFICATION of the read list: what is in "read" runs without asking, so
// adding to it is a decision, made here, in the same commit as the code.

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type shellCase struct {
	cmd   string
	class ShellClass
	why   string // a fragment the reasons must carry ("" for a read)
}

func reads(cmds ...string) []shellCase {
	var out []shellCase
	for _, c := range cmds {
		out = append(out, shellCase{c, ShellRead, ""})
	}
	return out
}

func cards(why string, cmds ...string) []shellCase {
	var out []shellCase
	for _, c := range cmds {
		out = append(out, shellCase{c, ShellWrite, why})
	}
	return out
}

func refusals(why string, cmds ...string) []shellCase {
	var out []shellCase
	for _, c := range cmds {
		out = append(out, shellCase{c, ShellRefused, why})
	}
	return out
}

func cat(parts ...[]shellCase) []shellCase {
	var out []shellCase
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func bashMatrix(home string) []shellCase {
	slashed := filepath.ToSlash(home)
	msys := slashed
	if len(slashed) > 2 && slashed[1] == ':' {
		msys = "/" + strings.ToLower(slashed[:1]) + slashed[2:]
	}
	return cat(
		// ---- a plain look runs at once ---------------------------------------------------------
		reads("ls", "ls -la", "ls -R docs", "ls | head -3", "dir", "pwd", "cat README.md", "cat docs/a.md docs/b.md",
			"head -n 5 README.md", "tail -n 20 log.txt", "wc -l README.md", "grep -rn foo .", `grep -n "a b" README.md`,
			"rg foo", "rg -n foo src", "find . -name '*.go'", "find . -maxdepth 2 -type f", "echo hello",
			`echo "a > b"`, `echo 'rm -rf /'`, `printf '%s\n' x`, "date", "date +%Y", "whoami", "hostname", "uname -a",
			"which git", "type ls", "stat README.md", "du -sh .", "df -h", "basename a/b.txt", "dirname a/b.txt",
			"diff a.txt b.txt", "cmp a b", "sha256sum README.md", "cut -d, -f1 data.csv", "sort data.csv",
			"sort -r data.csv", "uniq data.txt", "tr a-z A-Z", "cat *.md", "ls *.go", "ls -d */", "echo hi 2>/dev/null",
			"ls nosuchthing 2>&1", "ls > /dev/null", "test -f README.md", "head -n1 README.md # a comment", "true", ":",
			"file README.md", "nl README.md", "realpath .", "readlink -f .", "cat < README.md", `ls "my file.txt"`,
			"cat < /dev/null", "ls /dev/null", "l's'",
			// the ground's own path is not a reach, in either spelling
			"ls "+slashed, "ls "+msys+"/docs", "cat "+slashed+"/README.md"),
		reads("git status", "git status --short", "git log --oneline -5", "git diff", "git diff HEAD~1", "git show HEAD",
			"git branch", "git branch -a", "git branch --list", "git tag", "git tag -l", "git ls-files | wc -l",
			"git remote -v", "git stash list", "git rev-parse HEAD", "git --no-pager log -1", "git blame README.md",
			"git config --get user.name", "git grep foo", "git worktree list", "git reflog",
			"cat a.txt | grep x | sort | uniq -c", "ls && pwd", "ls; pwd", "ls || echo no"),

		// ---- anything else asks ------------------------------------------------------------------
		cards("deletes", "rm x", "rm -rf d", "del x", "rmdir d", "unlink x", "'rm' x", `"rm" x`, "r''m x", `r\m x`,
			`\rm x`, "ls; rm x", "ls && rm x", "ls | rm", "find . -delete"),
		cards("writes files", "mkdir d", "touch f", "cp a b", "mv a b", "ln -s a b", "tee f", "echo hi | tee f",
			"sed -i s/a/b/ f", "sed s/a/b/ f", "awk '{print}' f", "chmod +x f", "xargs echo", "tar -xf a.tar"),
		cards("redirects output", "echo hi > f", "echo hi >> f", "cat a > b", "ls >f", "echo hi &> f", "ls >| f", "ls 2> err.txt"),
		cards("installs", "pip install x", "pip3 install x", "npm install", "npx x", "pnpm add x", "choco install x",
			"winget install x", "apt install x", "brew install x"),
		cards("reaches the network", "curl http://example.com", "wget http://example.com", "ssh host", "scp a h:b",
			"nc host 80", "ping host", "git pull", "git fetch", "git clone x", "git ls-remote"),
		cards("sends saved work", "git push", "git push origin main"),
		cards("changes the repository", "git commit -m x", "git add .", "git checkout main", "git switch x", "git reset --hard",
			"git clean -fd", "git tag v1", "git branch newbranch", "git branch -d old", "git stash", "git stash pop",
			"git config user.name x", "git rebase main", "git merge x", "git restore f", "git rm f", "git mv a b"),
		cards("git option", "git -c core.pager=x log", "git --git-dir=x log", "git --paginate log"),
		cards("uses the option", "git log --output=f", "git diff --ext-diff"),
		cards("remotes", "git remote add x y", "git remote show origin"),
		cards("rewrites the log", "git reflog expire"),
		cards("runs a program", "python x.py", `python -c "print(1)"`, "node x.js", "bash x.sh", "sh -c ls", "make",
			"go test ./...", "env rm x", "command rm x", "sudo ls", "nohup ls", "time ls", "exec ls", "eval ls", "source x", ". x"),
		cards("runs a program by path", "./run.sh", "./x"),
		cards("which the gate cannot read", "cd docs", "cd", "zzz"),
		cards("sets a variable", "FOO=bar ls", "X=1", "A=b C=d ls"),
		cards("shell variable", "echo $PATH", "ls $PWD", `echo "$PATH"`, "echo ${X}", "echo $?"),
		cards("shell syntax the gate does not read", "echo $(ls)", "echo `ls`", `echo "$(ls)"`, "ls &", "(ls)",
			"cat <<EOF", "cat <<< hi", "diff <(ls) <(ls)", "ls ;; ls", `echo "unterminated`, "echo 'unterminated"),
		cards("brace", "echo {a,b}", "ls {1..3}", "{ ls; }"),
		cards("an operator", "&& ls", "|| ls", "| ls"),
		cards("the command word is built by the shell", "$X ls", "l* x"),
		cards("shell keyword", "if true; then ls; fi", "for x in a b; do ls; done", "while true; do ls; done", "! ls"),
		cards("option", "sort -o out in", "sort --output=out in", "head -f x", "tail -f x", "rg --pre cat foo", "rg -z foo",
			"file -C x"),
		cards("names a second file", "uniq in out"),
		cards("sets the clock", "date 010112002020", "date -s 12:00"),
		cards("sets the host name", "hostname newname"),
		cards("deletes, runs", "find . -exec rm {} \\;", "find . -ok rm {} \\;"),
		cards("writes", "find . -name x -fprint out"),

		// ---- refused by name; approval does not lift it -------------------------------------------
		refusals("secret", "cat .env", "cat ./.env", "cat .env.local", "cat '.env'", `cat ".env"`, `cat .e""nv`, `cat .\env`,
			"cat docs/.env", "grep KEY .env", "head -1 .env", "source .env", "echo x > .env", "cp .env x", "rm .env",
			"git add .env", "git show HEAD:.env", "find . -name .env", "cat .envrc", "cat .environment",
			"cat id_rsa", "cat ~/.ssh/id_rsa", "cat credentials.json", "cat secrets.yaml", "cat server_key.pem",
			"cat token.txt", "cat .netrc", "cat .npmrc", "cat .pypirc", "cat sub/api_key.txt"),
		refusals("hidden files", "cat .e*", "cat .*", "ls .[e]nv", "cat [.]env", "cat .e?v", "cat docs/.*"),
		refusals("client", "cat vault/notes.md", "ls vault/", "cat a.client.md", "cat x/vault/y"),
		refusals("outside the ground", "ls ..", "ls ../", "cat ../x", "ls ./..", "cat a/../../x", "cd ..", "ls /", "ls /etc",
			"cat /etc/passwd", "ls /c/Users", "cat C:/Windows/win.ini", `cat C:\Windows\win.ini`, "ls ~", "cat ~/x",
			"ls /tmp", "cat /proc/self/environ", "ls /mnt/c", "cat /dev/tcp/1.2.3.4/80", "ls //server/share",
			"ls /usr/bin", "cat /bin/ls", "git -C .. log", "cat "+slashed+"/../x", "cat "+msys+"/../x", "find / -name x", "ls /home"),
		refusals("home directory", "echo $HOME", "ls $HOME/x", "ls ${HOME}", "ls %USERPROFILE%", "cat $USERPROFILE/x"),
		refusals("key", "echo atl_0123456789abcdef0123456789abcdef", "echo sk-abcdefghijklmnopqrstuvwxyz",
			`curl -H "Authorization: Bearer abcdefghijklmnopqrstuvwxyz0123" http://x`, "export X=ghp_abcdefghijklmnopqrstuvwxyz0123",
			"echo AKIAABCDEFGHIJKLMNOP"),
	)
}

// variants say a command the ways a hurried person would, plus the ways that put it behind or beside
// something harmless. None may change the verdict.
func variants(cmd string) []string {
	return []string{
		cmd,
		"  " + cmd,
		cmd + "   ",
		"\t" + cmd,
		cmd + "\n",
		"\n" + cmd,
		cmd + " # a trailing comment",
		"ls && " + cmd,
		"ls; " + cmd,
		cmd + " | cat",
		cmd + " 2>&1",
	}
}

func hasWhy(v ShellVerdict, frag string) bool {
	for _, w := range v.Why {
		if strings.Contains(strings.ToLower(w), strings.ToLower(frag)) {
			return true
		}
	}
	return false
}

func TestEveryTypedCommandIsJudgedTheSameHoweverItIsSaid(t *testing.T) {
	home := t.TempDir()
	cases := bashMatrix(home)
	if len(cases) < 250 {
		t.Fatalf("the matrix holds %d commands; it was built with more than 250 and a shrunk table is a weaker gate", len(cases))
	}
	said := 0
	for _, c := range cases {
		plain := JudgeBash(c.cmd, home)
		if plain.Class != c.class {
			t.Errorf("%-52q judged %s %v, want %s", c.cmd, plain.Class, plain.Why, c.class)
			continue
		}
		if c.why != "" && !hasWhy(plain, c.why) {
			t.Errorf("%-52q judged %s but its reasons %v never say %q", c.cmd, plain.Class, plain.Why, c.why)
		}
		if c.class == ShellRead && len(plain.Why) != 0 {
			t.Errorf("%-52q is a read and still carries reasons %v", c.cmd, plain.Why)
		}
		for _, v := range variants(c.cmd) {
			said++
			got := JudgeBash(v, home)
			// A read said behind a worse thing is the worse thing; every other way keeps the class.
			want := c.class
			if strings.HasPrefix(v, "ls &&") || strings.HasPrefix(v, "ls;") {
				want = worse(ShellRead, c.class)
			}
			if got.Class != want {
				t.Errorf("%q said as %q came out %s %v, want %s", c.cmd, v, got.Class, got.Why, want)
			}
		}
	}
	t.Logf("%d commands, said %d ways", len(cases), said)
}

func worse(a, b ShellClass) ShellClass {
	rank := map[ShellClass]int{ShellRead: 0, ShellWrite: 1, ShellRefused: 2}
	if rank[a] >= rank[b] {
		return a
	}
	return b
}

// A REFUSAL NEVER ECHOES A KEY, and a refused text carries only refusals -- the write reasons of a
// command that was refused are noise on a card nobody can approve.
func TestARefusalNamesWhatWasRefusedAndNeverEchoesAKey(t *testing.T) {
	home := t.TempDir()
	key := "atl_0123456789abcdef0123456789abcdef"
	v := JudgeBash("curl -H 'X: "+key+"' http://example.com", home)
	if v.Class != ShellRefused {
		t.Fatalf("a key on a command line was %s", v.Class)
	}
	for _, w := range v.Why {
		if strings.Contains(w, key) || strings.Contains(w, key[:12]) {
			t.Errorf("the refusal repeats the key: %q", w)
		}
	}
	v = JudgeBash("rm .env", home)
	if v.Class != ShellRefused || hasWhy(v, "deletes") {
		t.Errorf("a refused delete carries its write reasons too: %v", v.Why)
	}
	// A reach is refused once, in one sentence, however many of the gate's passes see it.
	for _, cmd := range []string{"cat ../x", "ls /etc", "ls ..", "cat /etc/passwd"} {
		if v := JudgeBash(cmd, home); len(v.Why) != 1 {
			t.Errorf("%q is refused %d times over: %v", cmd, len(v.Why), v.Why)
		}
	}
}

// A WORD NAMED BY A VARIABLE, A SUBSTITUTION OR A BRACE IS NEVER A READ: the gate cannot see what it
// becomes, so it is asked about, whatever the command word is.
func TestNothingTheShellBuildsAtRunTimeIsARead(t *testing.T) {
	home := t.TempDir()
	for _, cmd := range []string{"cat $F", "cat ${F}", "cat $(echo x)", "cat `echo x`", "cat {a,b}", `cat "$F"`, "cat x$F", "echo $((1+1))"} {
		if v := JudgeBash(cmd, home); v.Class == ShellRead {
			t.Errorf("%q was read although the shell builds part of it: %v", cmd, v.Why)
		}
	}
}

// THE LEXER READS QUOTES THE WAY BASH DOES, which is how `r”m` is `rm` and `".env"` is `.env`.
func TestTheLexerUnquotesBeforeItLooks(t *testing.T) {
	toks, bad := shLex(`r''m 'a b' "c d" e\ f $'x' ` + "`y`")
	if len(bad) == 0 {
		t.Error("$'...' and backticks went unreported")
	}
	var words []string
	for _, tk := range toks {
		if tk.kind == shWordK {
			words = append(words, tk.text)
		}
	}
	if len(words) < 4 || words[0] != "rm" || words[1] != "a b" || words[2] != "c d" || words[3] != "e f" {
		t.Errorf("words = %q", words)
	}
	toks, _ = shLex("a 2>&1 >out b >>log <in")
	var redirs []string
	for _, tk := range toks {
		if tk.kind == shRedirK {
			d := ""
			if tk.dup {
				d = "(dup)"
			}
			redirs = append(redirs, tk.fd+tk.text+d)
		}
	}
	if fmt.Sprint(redirs) != "[2>&(dup) > >> <]" {
		t.Errorf("redirections = %v", redirs)
	}
}

// PYTHON'S HALF THAT IS TEXT: what is refused by name is refused before the session reads a line, and
// ordinary code that merely contains slashes or the word `token` is not a secret.
func TestPythonTextIsRefusedByNameAndOrdinaryCodeIsNot(t *testing.T) {
	home := t.TempDir()
	for _, code := range []string{
		`open(".env")`, `x = ".env.local"`, `open("../secret.txt")`, `import os; os.listdir("/etc")`, `Path("C:/Users/x")`,
		`open("credentials.json")`, `k = "atl_0123456789abcdef0123456789abcdef"`, `open("vault/x")`, `open("~/x")`,
		"open('C:\\\\Windows\\\\x')", `s = "$HOME"`,
	} {
		if v := JudgePythonText(code, home); v.Class != ShellRefused {
			t.Errorf("%q was %s, want refused", code, v.Class)
		}
	}
	for _, code := range []string{
		`print(os.environ)`, `tokens.append(1)`, `import secrets`, `10 / 2 / 5`, `"a/b/c".split("/")`, `s = "/"`,
		`y = 5 // 2`, `u = "http://example.com/x"`, `os.getenv("HOME")`, `d["env"]`, `x = a/b/c`, `1 if x else 2`,
		`r = re.sub("a", "b", s)`, `print("token:", token)`,
	} {
		if v := JudgePythonText(code, home); v.Class != ShellRead {
			t.Errorf("%q was %s %v, want it left to the session", code, v.Class, v.Why)
		}
	}
}

// THE GATE AND BASH MUST SEE THE SAME WORDS. The lexer reads a command the way bash splits it; if the two
// ever parted, the gate would judge one command and bash would run another. `printf '%s|'` over the same text
// prints the words bash made of it, one per argument, and the lexer's words must be those -- through the real
// Git Bash and the real command line Windows hands it, which is where a quote can be eaten or doubled.
func TestTheLexerSeesTheWordsBashSees(t *testing.T) {
	needBash(t)
	bashPath, _ := findBash()
	cases := []string{
		`a b  c`, `'a b' "c d" e\ f`, `"say \"hi\""`, `'it'"'"'s'`, `a\\b`, `"a\\b"`, `'a\\b'`, `"$"`, `"a\$b"`,
		`x"y"'z'w`, "tab\tsep", "\"multi\nline\"", `a"b c"d`, `\"quoted\"`, `"\\"`, `''`, `""`, `a''b`,
		`'#notcomment' #a comment`, "\"back\\`tick\"", `'it''s'`, `"%s" '%d'`, `a^b`, `"a&b" 'c|d' e\;f`, `'"' "'"`,
		`\\\\`, `"\\\\"`, `'C:\Users\x' "D:\y z"`, `"tail\\"`, `"é ü 日本"`, `a=b`, `--opt="v w"`,
	}
	for _, text := range cases {
		line := "printf '%s|' " + text
		toks, bad := shLex(line)
		if len(bad) != 0 {
			t.Errorf("%q: the lexer reports it cannot read it: %v", text, bad)
			continue
		}
		var words []string
		for _, tk := range toks {
			if tk.kind == shWordK {
				words = append(words, tk.text)
			}
		}
		if len(words) < 2 {
			t.Errorf("%q: the lexer found %d words", text, len(words))
			continue
		}
		lexer := strings.Join(words[2:], "|")
		if len(words) > 2 {
			lexer += "|"
		}
		args, env := bashInvocation(line)
		res := spawn(bashPath, args, spawnOpts{Dir: t.TempDir(), CleanEnv: shellEnv(), Env: env, Timeout: 20 * time.Second})
		if res.Err != nil {
			t.Errorf("%q: bash said %v / %q", text, res.Err, res.Stderr)
			continue
		}
		if res.Stdout != lexer {
			t.Errorf("%q: the gate sees %q and bash sees %q", text, lexer, res.Stdout)
		}
	}
}
