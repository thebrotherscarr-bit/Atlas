package tools

// The Python session's judge, held by a matrix the way the bash gate is (shellgate_test.go): what the
// session calls a plain calculation runs at once, and everything it does not know how to read is a
// card. The driver is the real program, started the way the door starts it; only Python can read
// Python, so the judgement is measured by asking it.

import (
	"testing"
	"time"
)

func TestEveryPythonEntryIsJudgedAsAPersonWouldJudgeIt(t *testing.T) {
	needPython(t)
	home := t.TempDir()
	s := pySessionFor(home)
	s.mu.Lock()
	defer s.mu.Unlock()
	t.Cleanup(func() { s.reset() })
	if _, err := s.start("python", home); err != nil {
		t.Fatal(err)
	}
	// Names the judge relies on are read off the live session, so make them real.
	for _, setup := range []string{"x = 5", "d = {}", "y = 1"} {
		if r, err := s.run(setup, 20*time.Second); err != nil || !r.OK {
			t.Fatalf("%q: %+v %v", setup, r, err)
		}
	}
	plain := []string{
		"1 + 1", "x = 5", "x * 2", "[i*i for i in range(5)]", `"abc".upper()`, "sum(range(10))", `len("x")`,
		`print("hi")`, `{"a": 1}["a"]`, `f"{1+1}"`, "sorted([3, 1, 2])", "d = {}", `d["a"] = 1`, `d.get("a")`,
		"for i in range(3): print(i)", "xs = [3,1,2]; xs.sort(); xs", "(lambda a: a + 1)", "y = x if x else 0",
		"a, b = 1, 2", "z = [k for k in range(3) if k % 2]", "'%d' % 5", "max(1, 2, key=abs)", "x += 1",
		"while False: pass", "del x", "assert 1 == 1", "len = 3", "unknown_name + 1", "s = 'a,b'.split(',')",
		"t = (1, 2)[0]", "n = None", "x if x else y", "{i: i * 2 for i in range(3)}",
		"sorted(['b','a'], key=lambda s: s)", "print('a', 'b', sep='-')", "int('7') + float('1.5')",
	}
	cards := []string{
		"import os", "from os import path", "open('x', 'w')", "__import__('os')", "eval('1')", "exec('x=1')",
		"os.system('dir')", "def f(): pass", "class A: pass", "(lambda: 1)()", "with open('x') as f: pass",
		"try:\n    pass\nexcept Exception:\n    pass", "x.__class__", "getattr(x, 'y')", "type(1)", "input()",
		"compile('1', 'x', 'eval')", "globals()", "vars()", "dir()", "help()", "breakpoint()", "print.__self__",
		"(1).real", "x.imag", "f = open", "f('x')", "import sys; sys.exit()", "yield 1", "return 1",
		"raise ValueError()", "global q", "q.attr = 1", "[x for x in open('f')]", "x.format(1)", "'{}'.format(1)",
		"str.__dict__", "map(open, ['a'])", "set_trace()", "lambda: open('x')()", "unknown_function(1)",
		"print(open('f').read())", "[1, 2].__len__()", "exit()", "quit()", "os.getcwd()", "subprocess.run(['ls'])",
		"x = open('f')", "del d['a']; os", "async def g(): pass", "match x:\n    case 1: pass",
	}
	for _, code := range plain {
		v, err := s.judge(code)
		if err != nil || v.Class != ShellRead {
			t.Errorf("%-48q judged %s %v (%v), want a plain calculation", code, v.Class, v.Why, err)
		}
	}
	for _, code := range cards {
		v, err := s.judge(code)
		if err != nil || v.Class != ShellWrite || len(v.Why) == 0 {
			t.Errorf("%-48q judged %s %v (%v), want a card with a reason", code, v.Class, v.Why, err)
		}
	}
	t.Logf("%d plain, %d cards", len(plain), len(cards))
}

// A NAME IS TRUSTED FOR WHAT IT HOLDS NOW, not for what it held when it was judged: an approved cell
// that bound `x` to a file does not make `x.read()` look like arithmetic.
func TestAPythonNameIsTrustedForWhatItHoldsNow(t *testing.T) {
	needPython(t)
	home := t.TempDir()
	s := pySessionFor(home)
	s.mu.Lock()
	defer s.mu.Unlock()
	t.Cleanup(func() { s.reset() })
	if _, err := s.start("python", home); err != nil {
		t.Fatal(err)
	}
	if v, _ := s.judge("n = 1"); v.Class != ShellRead {
		t.Fatalf("n = 1: %v", v)
	}
	s.run("n = 1", 20*time.Second)
	if v, _ := s.judge("n + 1"); v.Class != ShellRead {
		t.Fatalf("a name holding a number is plain data: %v", v)
	}
	// What an approved cell does: n is now an object, not a number.
	if r, err := s.run("import io; n = io.StringIO('abc')", 20*time.Second); err != nil || !r.OK {
		t.Fatalf("%+v %v", r, err)
	}
	for _, code := range []string{"n + 1", "n.read()", "n.getvalue()", "print(n)", "[n]"} {
		if v, _ := s.judge(code); v.Class != ShellWrite {
			t.Errorf("%q was %s although n now holds an object that is not plain data", code, v.Class)
		}
	}
	// And shadowing a builtin with something else takes it off the list of pure functions.
	if r, err := s.run("import os; len = os.system", 20*time.Second); err != nil || !r.OK {
		t.Fatalf("%+v %v", r, err)
	}
	if v, _ := s.judge(`len("x")`); v.Class != ShellWrite {
		t.Errorf("len now names os.system and was %s", v.Class)
	}
}
