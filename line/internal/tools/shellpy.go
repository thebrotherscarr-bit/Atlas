package tools

// THE OPERATOR'S PYTHON (WHAT'S LEFT H15, 2026-10-03; his word: "the one place I struggled with that
// system was getting the python loop to work properly").
//
// One Python process per world, started from the ground with the python the door runs the engine with,
// and kept: `x = 5` and then `x * 2` works. The door and the process speak JSON lines over a pipe;
// every request carries an id and only the answer with that id is taken, so nothing a child process
// prints can pose as one (the driver points its descriptor 1 at its descriptor 2, which the door
// collects as "stray" output, and its descriptor 0 at the null device, so input() reads nothing).
//
// WHO JUDGES PYTHON. Only Python can read Python, so the driver judges (ast.parse, never exec): an
// entry is a plain calculation -- numbers, text, containers of them, a short list of pure functions and
// methods, names that hold plain data -- or it is a card. An import, a file, a call to anything not on
// that list, a def, a class, an attribute the list does not know: all cards. The names it trusts are
// read off the live session (a name is plain data if its value is), not remembered, so an approved cell
// that bound a name to a file cannot make a later cell look harmless.
//
// BOUNDED. One entry at a time per world. An entry that runs past its limit ends the whole process:
// the names are gone, and the answer says so. The process dies with the door (its stdin closes).

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

// pyDriver is the program the session runs (python -c). It is tested as a program, by the same
// judgements and runs a person would make (shellpy_test.go).
const pyDriver = `import ast, builtins, io, json, os, sys, traceback, contextlib

# THE OPERATOR'S PYTHON, one session per world (the door's shellpy.go starts this with -c).
# It speaks JSON lines on a private copy of the pipe the door gave it. Descriptor 1 is then pointed at
# descriptor 2, so anything a child process prints can never be mistaken for an answer, and descriptor 0
# at the null device, so input() reads nothing and cannot eat a request.
_wire_out = os.fdopen(os.dup(1), "w", encoding="utf-8", newline="\n")
_wire_in = os.fdopen(os.dup(0), "r", encoding="utf-8", newline="\n")
os.dup2(2, 1)
os.dup2(os.open(os.devnull, os.O_RDONLY), 0)

G = {"__name__": "__console__"}
CAP = 60000

DATA_ATOMS = (int, float, complex, bool, str, bytes, type(None), range)
SAFE_CALLS = set((
    "abs all any ascii bin bool bytes chr complex dict divmod enumerate filter float frozenset hex int "
    "isinstance len list map max min oct ord pow print range repr reversed round set slice sorted str sum "
    "tuple zip").split())
SAFE_METHODS = set((
    "capitalize casefold center count endswith expandtabs find index isalnum isalpha isascii isdecimal "
    "isdigit islower isnumeric isspace istitle isupper join ljust lower lstrip partition removeprefix "
    "removesuffix replace rfind rindex rjust rpartition rsplit rstrip split splitlines startswith strip "
    "swapcase title upper zfill decode encode hex fromhex "
    "append clear copy extend insert pop remove reverse sort "
    "get items keys popitem setdefault update values "
    "add difference discard intersection isdisjoint issubset issuperset symmetric_difference union "
    "bit_length conjugate is_integer").split())


def plain(v, depth=0, budget=None):
    """A value the gate may treat as plain data: numbers, text, and containers of them."""
    if budget is None:
        budget = [2000]
    budget[0] -= 1
    if budget[0] < 0 or depth > 6:
        return False
    t = type(v)
    if t in DATA_ATOMS:
        return True
    if t in (list, tuple, set, frozenset):
        return all(plain(x, depth + 1, budget) for x in v)
    if t is dict:
        return all(plain(k, depth + 1, budget) and plain(x, depth + 1, budget) for k, x in v.items())
    return False


class Judge:
    """Reads Python and says whether it is a plain calculation. It never runs anything."""

    def __init__(self):
        self.why = []
        self.bound = set(n for n, v in G.items() if not n.startswith("__") and plain(v))

    def no(self, msg):
        if msg not in self.why and len(self.why) < 6:
            self.why.append(msg)
        return False

    def stmts(self, body):
        ok = True
        for s in body:
            ok = self.stmt(s) and ok
        return ok

    def target(self, t):
        if isinstance(t, ast.Name):
            self.bound.add(t.id)
            return True
        if isinstance(t, (ast.Tuple, ast.List)):
            ok = True
            for e in t.elts:
                ok = self.target(e) and ok
            return ok
        if isinstance(t, ast.Starred):
            return self.target(t.value)
        if isinstance(t, ast.Subscript):
            if isinstance(t.value, ast.Name) and t.value.id in self.bound:
                return self.expr(t.slice)
            return self.no("assigns into something that is not plain data")
        return self.no("assigns to an attribute or other target")

    def stmt(self, s):
        if isinstance(s, ast.Expr):
            return self.expr(s.value)
        if isinstance(s, ast.Assign):
            ok = self.expr(s.value)
            for t in s.targets:
                ok = self.target(t) and ok
            return ok
        if isinstance(s, ast.AugAssign):
            if isinstance(s.target, ast.Name) and s.target.id not in self.bound:
                return self.no("changes the name " + s.target.id + ", which is not plain data")
            return self.expr(s.value) and self.target(s.target)
        if isinstance(s, ast.AnnAssign):
            return (s.value is None or self.expr(s.value)) and self.target(s.target)
        if isinstance(s, ast.If):
            return self.expr(s.test) and self.stmts(s.body) and self.stmts(s.orelse)
        if isinstance(s, ast.While):
            return self.expr(s.test) and self.stmts(s.body) and self.stmts(s.orelse)
        if isinstance(s, ast.For):
            return self.expr(s.iter) and self.target(s.target) and self.stmts(s.body) and self.stmts(s.orelse)
        if isinstance(s, (ast.Pass, ast.Break, ast.Continue)):
            return True
        if isinstance(s, ast.Assert):
            return self.expr(s.test) and (s.msg is None or self.expr(s.msg))
        if isinstance(s, ast.Delete):
            ok = True
            for t in s.targets:
                if isinstance(t, ast.Name) and t.id in self.bound:
                    continue
                ok = self.no("deletes something that is not a plain name") and ok
            return ok
        names = {
            ast.Import: "imports a module", ast.ImportFrom: "imports a module",
            ast.FunctionDef: "defines a function", ast.AsyncFunctionDef: "defines a function",
            ast.ClassDef: "defines a class", ast.Return: "uses return", ast.With: "uses with",
            ast.Try: "uses try", ast.Raise: "raises", ast.Global: "uses global", ast.Nonlocal: "uses nonlocal",
        }
        for k, msg in names.items():
            if isinstance(s, k):
                return self.no(msg)
        return self.no("uses " + type(s).__name__)

    def name_ok(self, n):
        if n in self.bound:
            return True
        if n in SAFE_CALLS and G.get(n, getattr(builtins, n, None)) is getattr(builtins, n, None):
            return True
        if n not in G and not hasattr(builtins, n):
            return True  # undefined: running it only raises NameError
        return self.no("uses the name " + n + ", which is not plain data")

    def exprs(self, nodes):
        ok = True
        for n in nodes:
            ok = self.expr(n) and ok
        return ok

    def expr(self, n):
        if n is None:
            return True
        if isinstance(n, ast.Constant):
            return True
        if isinstance(n, ast.Name):
            return self.name_ok(n.id)
        if isinstance(n, ast.JoinedStr):
            return self.exprs(n.values)
        if isinstance(n, ast.FormattedValue):
            return self.expr(n.value) and self.expr(n.format_spec)
        if isinstance(n, (ast.Tuple, ast.List, ast.Set)):
            return self.exprs(n.elts)
        if isinstance(n, ast.Starred):
            return self.expr(n.value)
        if isinstance(n, ast.Dict):
            return self.exprs([k for k in n.keys if k is not None]) and self.exprs(n.values)
        if isinstance(n, ast.UnaryOp):
            return self.expr(n.operand)
        if isinstance(n, ast.BinOp):
            return self.expr(n.left) and self.expr(n.right)
        if isinstance(n, ast.BoolOp):
            return self.exprs(n.values)
        if isinstance(n, ast.Compare):
            return self.expr(n.left) and self.exprs(n.comparators)
        if isinstance(n, ast.IfExp):
            return self.expr(n.test) and self.expr(n.body) and self.expr(n.orelse)
        if isinstance(n, ast.Subscript):
            return self.expr(n.value) and self.expr(n.slice)
        if isinstance(n, ast.Slice):
            return self.expr(n.lower) and self.expr(n.upper) and self.expr(n.step)
        if isinstance(n, ast.NamedExpr):
            ok = self.expr(n.value)
            return self.target(n.target) and ok
        if isinstance(n, ast.Lambda):
            keep = set(self.bound)
            for a in n.args.args + n.args.kwonlyargs + n.args.posonlyargs:
                self.bound.add(a.arg)
            if n.args.vararg:
                self.bound.add(n.args.vararg.arg)
            if n.args.kwarg:
                self.bound.add(n.args.kwarg.arg)
            ok = self.exprs(n.args.defaults) and self.expr(n.body)
            self.bound = keep
            return ok
        if isinstance(n, (ast.ListComp, ast.SetComp, ast.GeneratorExp, ast.DictComp)):
            keep = set(self.bound)
            ok = True
            for g in n.generators:
                ok = self.expr(g.iter) and ok
                ok = self.target(g.target) and ok
                ok = self.exprs(g.ifs) and ok
            if isinstance(n, ast.DictComp):
                ok = self.expr(n.key) and self.expr(n.value) and ok
            else:
                ok = self.expr(n.elt) and ok
            self.bound = keep
            return ok
        if isinstance(n, ast.Call):
            f = n.func
            args_ok = self.exprs(n.args) and self.exprs([k.value for k in n.keywords])
            if isinstance(f, ast.Name):
                return self.name_ok(f.id) and (f.id in SAFE_CALLS or self.no("calls " + f.id + "(), which is not a plain-data function")) and args_ok
            if isinstance(f, ast.Attribute):
                if f.attr.startswith("_") or f.attr not in SAFE_METHODS:
                    return self.no("calls ." + f.attr + "(), which the gate does not read as a plain-data method") and False
                return self.expr(f.value) and args_ok
            return self.no("calls something that is not a name")
        if isinstance(n, ast.Attribute):
            return self.no("reads the attribute ." + n.attr)
        return self.no("uses " + type(n).__name__)


def judge(code):
    try:
        tree = ast.parse(code, "<console>", "exec")
    except SyntaxError as e:
        return {"class": "read", "why": ["a syntax error (" + str(e.msg) + "): nothing would run"]}
    j = Judge()
    ok = j.stmts(tree.body)
    if ok and not j.why:
        return {"class": "read", "why": []}
    return {"class": "write", "why": j.why or ["code the gate cannot read as a plain calculation"]}


def run(code):
    out, err = io.StringIO(), io.StringIO()
    ok = True
    try:
        tree = ast.parse(code, "<console>", "exec")
        last = None
        if tree.body and isinstance(tree.body[-1], ast.Expr):
            last = ast.Expression(tree.body.pop().value)
        with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
            exec(compile(tree, "<console>", "exec"), G)
            if last is not None:
                v = eval(compile(last, "<console>", "eval"), G)
                if v is not None:
                    G["_"] = v
                    print(repr(v))
    except SystemExit as e:
        ok = False
        err.write("SystemExit(" + str(e.code) + ") ignored: the session keeps running\n")
    except SyntaxError as e:
        ok = False
        err.write("".join(traceback.format_exception_only(type(e), e)))
    except BaseException as e:
        ok = False
        frames = [ln for ln in traceback.format_exc().splitlines() if ln.lstrip().startswith('File "<console>"')]
        err.write("Traceback (most recent call last):\n" + "\n".join(frames[-3:]) + ("\n" if frames else "")
                  + "".join(traceback.format_exception_only(type(e), e)))
    o, e = out.getvalue(), err.getvalue()
    trunc = len(o) > CAP or len(e) > CAP
    return {"ok": ok, "out": o[:CAP], "err": e[:CAP], "truncated": trunc}


for line in _wire_in:
    line = line.strip()
    if not line:
        continue
    try:
        req = json.loads(line)
        if req.get("op") == "judge":
            res = judge(req.get("code", ""))
        elif req.get("op") == "run":
            res = run(req.get("code", ""))
        else:
            res = {"error": "unknown op"}
    except BaseException as ex:
        res = {"error": type(ex).__name__ + ": " + str(ex)}
    res["id"] = req.get("id") if isinstance(req, dict) else None
    _wire_out.write(json.dumps(res) + "\n")
    _wire_out.flush()
`

// strayBuf collects what the session's descriptors 1 and 2 carried: output a child process wrote. Capped.
type strayBuf struct {
	mu sync.Mutex
	b  []byte
}

const strayCap = 32 * 1024

func (s *strayBuf) Write(p []byte) (int, error) {
	n := len(p)
	s.mu.Lock()
	defer s.mu.Unlock()
	if room := strayCap - len(s.b); room > 0 {
		if len(p) > room {
			p = p[:room]
		}
		s.b = append(s.b, p...)
	}
	return n, nil
}

func (s *strayBuf) take() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := string(s.b)
	s.b = nil
	return out
}

type pySession struct {
	mu      sync.Mutex // one entry at a time
	cmd     *exec.Cmd
	in      io.WriteCloser
	out     *bufio.Reader
	stray   *strayBuf
	alive   bool
	started bool // started since the last reset: a start after that means the names were lost
	seq     int
}

var pySessions = struct {
	mu sync.Mutex
	by map[string]*pySession
}{by: map[string]*pySession{}}

func pySessionFor(home string) *pySession {
	pySessions.mu.Lock()
	defer pySessions.mu.Unlock()
	s, ok := pySessions.by[home]
	if !ok {
		s = &pySession{stray: &strayBuf{}}
		pySessions.by[home] = s
	}
	return s
}

// start launches the driver if none is running, and says whether the names of an earlier one are lost.
func (s *pySession) start(py, home string) (lost bool, err error) {
	if s.alive {
		return false, nil
	}
	cmd := exec.Command(py, "-u", "-X", "utf8", "-c", pyDriver)
	cmd.Dir = home
	cmd.Env = shellEnv()
	cmd.Stderr = s.stray
	in, err := cmd.StdinPipe()
	if err != nil {
		return false, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	if err := cmd.Start(); err != nil {
		return false, fmt.Errorf("the Python session could not start (%s): %v", py, err)
	}
	s.cmd, s.in, s.out = cmd, in, bufio.NewReaderSize(out, 256*1024)
	s.alive = true
	lost = s.started
	s.started = true
	return lost, nil
}

// kill ends the process and everything it started. The names go with it.
func (s *pySession) kill() {
	if s.cmd != nil && s.cmd.Process != nil {
		_ = killTree(s.cmd)
	}
	if s.in != nil {
		s.in.Close()
	}
	if c := s.cmd; c != nil {
		go c.Wait()
	}
	s.alive = false
}

// reset is the operator's own: end the session and forget that it ever ran.
func (s *pySession) reset() {
	if s.alive {
		s.kill()
	}
	s.started = false
	s.stray.take()
}

var errPyTimeout = errors.New("timed out")

// ask sends one request and waits for the answer that carries its id, up to wait. A session that does
// not answer is killed: a process that cannot be interrupted cannot be trusted with the next entry.
func (s *pySession) ask(op, code string, wait time.Duration) (map[string]any, error) {
	s.seq++
	id := s.seq
	req, err := json.Marshal(map[string]any{"id": id, "op": op, "code": code})
	if err != nil {
		return nil, err
	}
	if _, err := s.in.Write(append(req, '\n')); err != nil {
		s.kill()
		return nil, fmt.Errorf("the Python session is gone: %v", err)
	}
	type reply struct {
		m   map[string]any
		err error
	}
	ch := make(chan reply, 1)
	go func() {
		for {
			line, err := s.out.ReadString('\n')
			if err != nil {
				ch <- reply{nil, err}
				return
			}
			var m map[string]any
			if json.Unmarshal([]byte(line), &m) != nil {
				continue // not an answer
			}
			if n, _ := m["id"].(float64); int(n) == id {
				ch <- reply{m, nil}
				return
			}
		}
	}()
	select {
	case r := <-ch:
		if r.err != nil {
			s.kill()
			return nil, fmt.Errorf("the Python session ended unexpectedly: %v", r.err)
		}
		return r.m, nil
	case <-time.After(wait):
		s.kill()
		return nil, errPyTimeout
	}
}

func anyStrings(v any) []string {
	var out []string
	if xs, ok := v.([]any); ok {
		for _, x := range xs {
			out = append(out, fmt.Sprint(x))
		}
	}
	return out
}

// judge asks the session whether code is a plain calculation. Anything but a plain "read" is a card.
func (s *pySession) judge(code string) (ShellVerdict, error) {
	m, err := s.ask("judge", code, 20*time.Second)
	if err != nil {
		return ShellVerdict{}, err
	}
	if e, _ := m["error"].(string); e != "" {
		return ShellVerdict{}, errors.New(e)
	}
	v := ShellVerdict{Class: ShellWrite, Why: anyStrings(m["why"])}
	if c, _ := m["class"].(string); c == string(ShellRead) {
		v.Class = ShellRead
	}
	return v, nil
}

type pyRun struct {
	OK        bool
	Out, Err  string
	Truncated bool
	Stray     string
}

func (s *pySession) run(code string, wait time.Duration) (pyRun, error) {
	m, err := s.ask("run", code, wait)
	if err != nil {
		return pyRun{}, err
	}
	if e, _ := m["error"].(string); e != "" {
		return pyRun{}, errors.New(e)
	}
	r := pyRun{}
	r.OK, _ = m["ok"].(bool)
	r.Out, _ = m["out"].(string)
	r.Err, _ = m["err"].(string)
	r.Truncated, _ = m["truncated"].(bool)
	time.Sleep(40 * time.Millisecond) // let a finished child's last words reach the stray channel
	r.Stray = s.stray.take()
	return r, nil
}
