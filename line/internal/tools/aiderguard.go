package tools

// THE WALL AIDER RUNS INSIDE (WHAT'S LEFT H12, 2026-10-04).
//
// Aider is somebody else's program, and a MODEL decides what it writes. Measured on the first probe
// (2026-10-04, Aider 0.86.2, qwen2.5-coder:14b): with `--yes-always` it will edit or create a file at
// ANY path the model names, outside its own root included -- and `--yes-always` is not optional, because
// a run with nobody at the keyboard would otherwise wait for a "y" for ever. So the door does not trust
// Aider's own confirmations. It runs Aider inside this wall, and the wall is the same one the council's
// own `run_python` child has stood inside since 2026-09-22 (manjuel/skills.py, `_PY_JAIL`): a PEP 578
// audit hook, installed by THIS source, which is passed on the command line, so nothing in the scratch
// folder Aider works in can edit it.
//
//	WRITES go only where the door put a place for them: the one scratch folder of this run, Aider's own
//	home and temp folders, and the three history files of this run. Not the ground. Not the venv.
//	READS stay out of the ground and the real user profile, except aider/ itself (the venv, the work
//	folder) and the interpreter's own library -- so a `.env`, a client world or Desktop\Archive is not
//	something Aider can be sent to look at.
//	THE NETWORK is loopback only: an address lookup or a connection to anything that is not 127.x / ::1
//	is refused by name (RULE 4). Ollama is on loopback, and that is the only server Aider has.
//	NO OTHER PROCESS: Aider starts nothing (shell suggestions are off, there is no git, there is no
//	linter), so a spawn of any kind is refused.
//
// EVERY REFUSAL IS SAID, on stderr, as a line beginning ATLAS-AIDER-GUARD refused, before the exception is
// raised; the door lifts those lines into the answer, so the card shows what the wall stopped. TWO PROBES ARE THE
// EXCEPTION, refused without a word: the first real run (2026-10-04) showed every healthy run asking for them --
// GitPython runs `git version` when it is imported, and the platform module runs the Windows shell's `ver` -- and a card
// that warns of the wall on every run teaches him to read past the one run where it matters. Both are still refused.
//
// AND THE WALL ANNOUNCES ITSELF. Once the hook is armed it writes `ATLAS-AIDER-GUARD <nonce> armed` to
// stderr, where the nonce is a fresh random value the door put in the child's environment and the guard
// removed from it. A run whose stderr lacks that line did not run behind the wall, and the door throws
// the run away: a wall that silently was not there is believed, which is worse than none.
//
// HONEST LIMITS, copied from the precedent because they are the same:
//
//	AN AUDIT HOOK IS NOT A KERNEL SANDBOX. CPython's own documentation says so. It shuts every route
//	named above; it does not prove no route exists. A C extension already on this machine, or native code
//	loaded by ctypes (which Aider's own terminal libraries import, so it cannot be refused), steps around
//	it. The wall is there to bound an accident -- a model that names the wrong path, a library that
//	phones home -- not to contain a hostile program. The scratch copy, the line of work, the suites and
//	his click are what stand behind it.
//	A SYMLINK OR JUNCTION is judged by the spelled path, and also by its resolved one (both are roots).
//	EXISTENCE IS NOT SECRECY: stat() is not refused; reading a path's bytes is what is shut.

// aiderGuardSource is the whole guard. Run as: python -I -B -X utf8 -c <this> <module-or-.py> <args...>,
// with the wall in ATLAS_AIDER_WALL (JSON) and the handshake nonce in ATLAS_AIDER_NONCE. It is a raw
// string: there is no backtick in it, and every backslash is Python's own.
const aiderGuardSource = `
import ipaddress
import json
import os
import runpy
import socket
import sys

WALL = json.loads(os.environ.pop('ATLAS_AIDER_WALL'))
NONCE = os.environ.pop('ATLAS_AIDER_NONCE')


def norm(p):
    try:
        if isinstance(p, os.PathLike):
            p = os.fspath(p)
        if isinstance(p, bytes):
            p = p.decode('utf-8')
        if not isinstance(p, str) or not p:
            return None
        if p.startswith('\\\\?\\'):
            p = p[4:]
        return os.path.normcase(os.path.abspath(p))
    except (TypeError, ValueError, UnicodeDecodeError):
        return None


def roots(paths):
    out = []
    for p in paths:
        for q in (p, os.path.realpath(p)):
            n = norm(q)
            if n and n not in out:
                out.append(n)
    return out


GROUND = roots([WALL['ground']])
PROFILE = roots([WALL['profile']]) if WALL.get('profile') else []
AIDER = roots([WALL['aider']])
WRITE_DIRS = roots(WALL['write_dirs'])
WRITE_FILES = set(roots(WALL['write_files']))
LIB = roots([sys.base_prefix, sys.prefix, sys.exec_prefix, sys.base_exec_prefix] + [p for p in sys.path if p])
NULLS = set(filter(None, [norm(os.devnull), os.path.normcase('\\\\.\\nul'), os.path.normcase(os.devnull)]))
WRITE_FLAGS = os.O_WRONLY | os.O_RDWR | os.O_APPEND | os.O_CREAT | os.O_TRUNC

BOTH = ('os.remove', 'os.rmdir', 'os.rename', 'os.replace', 'os.link', 'os.symlink', 'os.chmod',
        'os.chown', 'os.truncate', 'os.utime', 'shutil.move', 'shutil.rmtree', 'shutil.chown')
LISTING = ('os.listdir', 'os.scandir', 'glob.glob')
SPAWN = ('subprocess.Popen', '_winapi.CreateProcess', 'os.system', 'os.exec', 'os.spawn', 'os.posix_spawn',
         'os.fork', 'os.forkpty', 'os.startfile', 'pty.spawn', 'webbrowser.open')
LOOKUP = ('socket.getaddrinfo', 'socket.gethostbyname', 'socket.gethostbyname_ex', 'socket.gethostbyaddr')
REACH = ('socket.connect', 'socket.bind', 'socket.sendto', 'socket.sendmsg')


def under(p, root):
    return p == root or p.startswith(root.rstrip(os.sep) + os.sep)


def inany(p, rs):
    return any(under(p, r) for r in rs)


def can_read(p):
    if p in NULLS or inany(p, AIDER) or inany(p, LIB):
        return True
    return not (inany(p, GROUND) or inany(p, PROFILE))


def can_write(p):
    return p in NULLS or inany(p, WRITE_DIRS) or p in WRITE_FILES


def loopback(h):
    if h is None:
        return True
    if isinstance(h, bytes):
        h = h.decode('ascii', 'replace')
    if not isinstance(h, str):
        return False
    h = h.strip('[]').lower()
    if h == 'localhost':
        return True
    try:
        return ipaddress.ip_address(h.split('%')[0]).is_loopback
    except ValueError:
        return False


def detail(args):
    try:
        a = args[1] if len(args) > 1 and isinstance(args[1], (list, tuple)) else args
        return ' '.join(str(x) for x in a)[:120]
    except Exception:
        return ''


def quiet(args):
    try:
        a = args[1] if len(args) > 1 else ''
        parts = [str(x).strip(chr(34)).lower() for x in (a if isinstance(a, (list, tuple)) else str(a).split())]
    except Exception:
        return False
    if len(parts) == 2 and os.path.basename(parts[0]) in ('git', 'git.exe') and parts[1] == 'version':
        return True
    # platform asks the Windows shell for its version three ways, in turn: ver, command /c ver and cmd /c ver
    return (len(parts) >= 3 and os.path.basename(parts[0]) in ('cmd', 'cmd.exe') and parts[1] == '/c'
            and ' '.join(parts[2:]) in ('ver', 'command /c ver', 'cmd /c ver'))


def refuse(what, why):
    line = what + ' -- ' + why
    try:
        sys.__stderr__.write('ATLAS-AIDER-GUARD refused: ' + line + '\n')
        sys.__stderr__.flush()
    except Exception:
        pass
    raise PermissionError('aider wall: ' + line)


def hook(event, args):
    if event == 'open':
        p = norm(args[0]) if args else None
        if p is None:
            return
        mode = args[1] if len(args) > 1 else None
        flags = args[2] if len(args) > 2 and isinstance(args[2], int) else 0
        writing = (isinstance(mode, str) and any(c in mode for c in 'wax+')) or bool(flags & WRITE_FLAGS)
        if writing:
            if not can_write(p):
                refuse('writing ' + p, 'the only places Aider may write are this run\'s scratch folder and its own home')
        elif not can_read(p):
            refuse('reading ' + p, 'the ground and the user profile are closed to Aider; only aider/ is open')
    elif event in BOTH:
        for a in args:
            p = norm(a)
            if p is not None and not can_write(p):
                refuse(event + ' on ' + p, 'that is outside what Aider may change')
    elif event == 'os.mkdir':
        p = norm(args[0]) if args else None
        # a folder that is already there cannot be made again, so asking is no change: Aider asks for its history files' folder
        # (the run's own) every time it starts, and the answer to that is the operating system's, not the wall's.
        if p is not None and not can_write(p) and not os.path.isdir(p):
            refuse('os.mkdir on ' + p, 'that is outside what Aider may change')
    elif event in LISTING:
        p = norm(args[0]) if args else None
        if p is not None and not can_read(p):
            refuse('listing ' + p, 'the ground and the user profile are closed to Aider')
    elif event in LOOKUP:
        if not loopback(args[0] if args else None):
            refuse('looking up ' + repr(args[0] if args else None), 'the estate is local (RULE 4): only loopback is reachable')
    elif event == 'socket.getnameinfo':
        a = args[0] if args else None
        if not (isinstance(a, tuple) and a and loopback(a[0])):
            refuse('looking up ' + repr(a), 'the estate is local (RULE 4): only loopback is reachable')
    elif event in REACH:
        a = args[1] if len(args) > 1 else None
        if not (isinstance(a, tuple) and a and a[0] != '' and loopback(a[0])):
            refuse(event + ' ' + repr(a), 'the estate is local (RULE 4): only loopback is reachable')
    elif event.startswith(SPAWN):
        if event == 'subprocess.Popen' and quiet(args):
            raise PermissionError('aider wall: no other process')
        refuse(event + ' ' + detail(args), 'Aider starts no other process')


sys.addaudithook(hook)
sys.__stderr__.write('ATLAS-AIDER-GUARD ' + NONCE + ' armed\n')
sys.__stderr__.flush()
TARGET = sys.argv[1]
sys.argv = [TARGET] + sys.argv[2:]
if TARGET.endswith('.py'):
    runpy.run_path(TARGET, run_name='__main__')
else:
    runpy.run_module(TARGET, run_name='__main__', alter_sys=True)
`

// aiderParseSource is the structural gate for a Python file Aider changed: parsed, never run, by the venv's
// own interpreter (3.11 here, so a construct the core's older CI legs cannot read is refused here too).
// Prints `line N: message` and exits 1 when it does not parse.
const aiderParseSource = `
import ast
import sys

p = sys.argv[1]
try:
    ast.parse(open(p, 'rb').read(), p)
except SyntaxError as e:
    print('line %s: %s' % (e.lineno, e.msg))
    sys.exit(1)
except ValueError as e:
    print('not text: %s' % e)
    sys.exit(1)
`

// aiderGuardArgs is the interpreter's argument list for one walled run: -I (no environment variables, no user
// site, no script folder on the path), -B (no bytecode is written, so the wall never has to reason about a
// cache), -X utf8 (a flag, because -I ignores PYTHONUTF8), then the guard and the program it runs.
func aiderGuardArgs(target string, rest ...string) []string {
	return append([]string{"-I", "-B", "-X", "utf8", "-c", aiderGuardSource, target}, rest...)
}

// aiderWall is the JSON the guard reads: where the ground is, which folder is open, and the only places a
// write may land.
type aiderWall struct {
	Ground     string   `json:"ground"`
	Profile    string   `json:"profile,omitempty"`
	Aider      string   `json:"aider"`
	WriteDirs  []string `json:"write_dirs"`
	WriteFiles []string `json:"write_files"`
}
