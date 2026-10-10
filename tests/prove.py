#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
prove.py -- THE BALL. One command that proves the whole of atlas.

atlas proves itself in eight places and they had never been gathered: the
Rust spine, the two Go modules, the two shipped batteries that live inside
the binaries, twenty-six golden verifiers, six workflow scripts and one
end-to-end suite. AGENTS.md named four of the verifiers; nobody ran the
other twenty-two. This rolls every leg into one run and one verdict.

    python tests/prove.py              the hermetic legs (no server, no model)
    python tests/prove.py --check      the fast legs only (no cargo)
    python tests/prove.py --live       add the legs that need :8090 and Ollama
    python tests/prove.py --quiet      one line per leg, nothing else

THREE VERDICTS, NOT TWO. A leg is PASS, FAIL, or ABSENT. ABSENT means the
leg named a dependency this ground does not hold -- the read-only source
grounds (estate\\, secondbrain\\) that fourteen of the cutters were cut
from, a binary not yet built, a door not answering. ABSENT is never counted
as a pass and never silently skipped: it prints with the exact path or
command that would answer it. Only a FAIL sets a red exit.

Law 5 (proves are hermetic) is honoured: nothing here writes to the record.
Every child gets stdin=DEVNULL -- a child that inherits a live engine's
stdin deadlocks it, which cost this estate nine sittings on 2026-09-10.

Zero external dependencies (law 6): stdlib only.
"""

import argparse
import os
import re
import subprocess
import sys
import time

HERE = os.path.dirname(os.path.abspath(__file__))
ATLAS = os.path.normpath(os.path.join(HERE, ".."))
GROUND = os.path.normpath(os.path.join(ATLAS, ".."))
TOOLS = os.path.join(ATLAS, "tools")

PASS, FAIL, ABSENT, SKIPPED = "PASS", "FAIL", "ABSENT", "SKIPPED"

# The atlas binary the door battery shells. Never built here: a prove that
# builds its own subject is not a prove.
ATLAS_BIN = os.path.join(ATLAS, "target", "debug",
                         "atlas.exe" if os.name == "nt" else "atlas")

MISSING_PATH = re.compile(r"No such file or directory: '([^']+)'")
NOT_FOUND = re.compile(r"panicked at [^\n]*\n *([^\n:]+): The system cannot find")


class Leg(object):
    """One provable thing: its name, its verdict, and what it would take."""

    def __init__(self, group, name, verdict, detail="", need="", secs=0.0):
        self.group = group
        self.name = name
        self.verdict = verdict
        self.detail = detail
        self.need = need
        self.secs = secs


def run(cmd, cwd=None, timeout=1800):
    """Run a child with its stdin closed. Returns (code, text)."""
    try:
        p = subprocess.run(cmd, cwd=cwd or ATLAS,
                           stdin=subprocess.DEVNULL,
                           stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                           timeout=timeout)
        return p.returncode, p.stdout.decode("utf-8", "replace")
    except FileNotFoundError as ex:
        return 127, "toolchain absent: %s" % ex
    except subprocess.TimeoutExpired:
        return 124, "timed out after %ds" % timeout


def tail(text, n=6):
    lines = [l.rstrip() for l in text.strip().splitlines() if l.strip()]
    return "\n".join(lines[-n:])


def absent_dependency(text):
    """Name the missing dependency if a refusal is an absence, not a break.

    A refusal counts as ABSENT only when it names a path that (a) does not
    exist and (b) lies outside the atlas tree -- one of the read-only source
    grounds this repo was folded from. A missing file INSIDE atlas is a real
    break and stays a FAIL.
    """
    for path in _candidates(text):
        if os.path.exists(path):
            continue
        if os.path.normpath(path).lower().startswith(
                os.path.normpath(ATLAS).lower()):
            continue
        return path
    return ""


def _candidates(text):
    """Paths a refusal named, absolute-ised against the ground."""
    out = []
    m = MISSING_PATH.search(text)
    if m:
        # A traceback prints the path as a Python literal, so every
        # separator arrives doubled. Put it back the way the disk spells it.
        out.append(m.group(1).replace("\\\\", "\\"))
    m = NOT_FOUND.search(text)
    if m:
        frag = m.group(1).strip()
        if frag:
            out.append(os.path.join(GROUND, frag.replace("/", os.sep)))
    return out


def leg_toolchains():
    """Nothing proves without these; name each one before leaning on it."""
    out = []
    for name, cmd in (("rust", ["cargo", "--version"]),
                      ("go", ["go", "version"]),
                      ("python", [sys.executable, "--version"])):
        code, text = run(cmd, timeout=60)
        first = text.strip().splitlines()[0] if text.strip() else ""
        if code == 0:
            out.append(Leg("TOOLCHAIN", name, PASS, first))
        else:
            out.append(Leg("TOOLCHAIN", name, ABSENT, first,
                           "install it, or prepend %USERPROFILE%\\.cargo\\bin"))
    if os.path.exists(ATLAS_BIN):
        out.append(Leg("TOOLCHAIN", "atlas binary", PASS,
                       os.path.relpath(ATLAS_BIN, ATLAS)))
    else:
        out.append(Leg("TOOLCHAIN", "atlas binary", ABSENT,
                       "the door battery shells this",
                       "cargo build -p atlas"))
    return out


def leg_spine():
    """The Rust spine: core, store, apps/atlas."""
    t = time.time()
    code, text = run(["cargo", "test", "--workspace"])
    secs = time.time() - t
    total = sum(int(n) for n in
                re.findall(r"^test result: \w+\. (\d+) passed", text, re.M))
    detail = "%d strokes" % total
    if code == 0:
        return [Leg("SPINE", "cargo test --workspace", PASS, detail, "", secs)]
    need = absent_dependency(text)
    if need:
        return [Leg("SPINE", "cargo test --workspace", ABSENT,
                    detail + " held; one leg wants the read-only oracle",
                    need, secs)]
    return [Leg("SPINE", "cargo test --workspace", FAIL, tail(text), "", secs)]


def leg_go(module, label):
    """One Go module: it must build before it is asked to prove."""
    out = []
    cwd = os.path.join(ATLAS, module)
    t = time.time()
    code, text = run(["go", "build", "./..."], cwd=cwd)
    if code != 0:
        return [Leg(label, "go build ./...", FAIL, tail(text), "",
                    time.time() - t)]
    out.append(Leg(label, "go build ./...", PASS, "", "", time.time() - t))
    t = time.time()
    code, text = run(["go", "test", "./...", "-count=1"], cwd=cwd)
    secs = time.time() - t
    oks = len(re.findall(r"^ok\s", text, re.M))
    bare = len(re.findall(r"\[no test files\]", text, re.M))
    detail = "%d packages proven, %d carry no prover" % (oks, bare)
    if code == 0:
        out.append(Leg(label, "go test ./...", PASS, detail, "", secs))
    elif "no working atlas binary" in text and not os.path.exists(ATLAS_BIN):
        out.append(Leg(label, "go test ./...", ABSENT,
                       "the door battery has no binary to shell",
                       "cargo build -p atlas", secs))
    else:
        out.append(Leg(label, "go test ./...", FAIL, tail(text), "", secs))
    return out


def leg_mcp():
    """THE LINE's shipped battery -- hermetic, temp grounds, loopback."""
    t = time.time()
    code, text = run(["go", "run", "./cmd/atlas-mcp", "--prove"],
                     cwd=os.path.join(ATLAS, "line"))
    secs = time.time() - t
    strokes = len(re.findall(r"\[PASS\]", text))
    broke = len(re.findall(r"\[FAIL\]", text))
    detail = "%d strokes" % strokes
    if code == 0 and not broke:
        return [Leg("BATTERY", "atlas-mcp --prove", PASS, detail, "", secs)]
    return [Leg("BATTERY", "atlas-mcp --prove", FAIL,
                "%s, %d broke\n%s" % (detail, broke, tail(text)), "", secs)]


def verifiers():
    """Every cutter that answers --verify, in name order."""
    out = []
    for fn in sorted(os.listdir(TOOLS)):
        if not fn.endswith(".py"):
            continue
        path = os.path.join(TOOLS, fn)
        with open(path, encoding="utf-8", errors="replace") as f:
            if "--verify" not in f.read():
                continue
        out.append((fn[:-3], path))
    return out


def leg_goldens():
    """Law 2: goldens before assertions. These re-cut and compare."""
    out = []
    for name, path in verifiers():
        t = time.time()
        code, text = run([sys.executable, path, "--verify"], timeout=600)
        secs = time.time() - t
        if code == 0:
            out.append(Leg("GOLDENS", name, PASS, "", "", secs))
            continue
        # A CUTTER THAT SHELLS THE SPINE IS ABSENT WITHOUT IT, NOT BROKEN.
        # absent_dependency only recognises an absence that NAMES A PATH. A
        # cutter whose subject is the unbuilt binary names a COMMAND instead,
        # so check_trade_parity fell through to FAIL and set a RED EXIT on any
        # machine that had not yet run cargo build -- which is every fresh
        # clone, and the first thing a second machine does. Found 2026-09-11 in
        # the packaging run by parking the binary and re-running: 20 held, 14
        # absent, 1 broke, exit 1, on a tree where nothing was wrong.
        #
        # This is the doctrine leg_go already applies two functions up, and the
        # one cmd/atlas-door/prove_test.go learned the same day: ABSENT names
        # what would answer it and is never a pass; only FAIL is red.
        #
        # Gated on the binary being GENUINELY ABSENT so this can never turn a
        # real break into an absence -- with the binary on disk the branch is
        # unreachable. The cutter says the same thing from the other side with
        # exit 2, and it is the only cutter in tools/ that uses 2 at all.
        if not os.path.exists(ATLAS_BIN) and "cargo build -p atlas" in text:
            out.append(Leg("GOLDENS", name, ABSENT,
                           "the spine this cutter shells is not built",
                           "cargo build -p atlas", secs))
            continue
        need = absent_dependency(text)
        if need:
            # NAME THE DIAL, not just the missing file. The goldens themselves
            # are committed and travel with the repo; only the RE-CUT needs the
            # oracle it was cut from, and that ground is private and does not
            # ship. Until 2026-09-11 this said "oracle not in this ground" and
            # pointed at a path that has never existed on any machine — the
            # cutters resolved it relative to atlas's parent, which stopped
            # being the oracle ground at the 2026-09-10 split. ABSENT was the
            # right verdict for the wrong reason, with no way to act on it.
            out.append(Leg("GOLDENS", name, ABSENT,
                           "oracle not in this ground "
                           "(set ATLAS_ORACLE_ROOT to the ground holding it)",
                           need, secs))
        else:
            out.append(Leg("GOLDENS", name, FAIL, tail(text), "", secs))
    return out


def leg_chain_verify():
    """THE CHAIN VERDICT, IN PYTHON (2026-10-01, his ruling: the verdict in
    Python first, proved against the goldens, then the Rust folded). The
    verifier is not a cutter -- it writes nothing and takes no cutter's word
    -- so it is not in leg_goldens' list; this is its own leg: every chain
    and injection in tests/fixtures/canon/chain_verdicts.json re-verdicted by
    tools/chain_verify.py and compared. The goldens are committed, so this
    leg is never ABSENT: it holds or it broke."""
    t = time.time()
    code, text = run([sys.executable, os.path.join(TOOLS, "chain_verify.py"), "--goldens"], timeout=300)
    secs = time.time() - t
    last = text.strip().splitlines()[-1] if text.strip() else ""
    if code == 0:
        return [Leg("GOLDENS", "chain_verify (python spine)", PASS, last, "", secs)]
    return [Leg("GOLDENS", "chain_verify (python spine)", FAIL, tail(text), "", secs)]


def mcp_online():
    import urllib.request
    try:
        urllib.request.urlopen("http://127.0.0.1:8090/health", timeout=3)
        return True
    except Exception:
        return False


def leg_workflows(live):
    """The six seat workflows. They drive the live door; without it they are
    ABSENT, never passed."""
    wfdir = os.path.join(HERE, "workflows")
    scripts = sorted(f for f in os.listdir(wfdir)
                     if f.startswith("wf_") and f.endswith(".py"))
    if not live:
        return [Leg("WORKFLOWS", f[:-3], SKIPPED, "pass --live to run")
                for f in scripts]
    if not mcp_online():
        return [Leg("WORKFLOWS", f[:-3], ABSENT, "the door is not answering",
                    "start atlas-mcp on :8090") for f in scripts]
    out = []
    for f in scripts:
        t = time.time()
        code, text = run([sys.executable, os.path.join(wfdir, f)], timeout=900)
        secs = time.time() - t
        line = ""
        for ln in reversed(text.strip().splitlines()):
            if "PASS" in ln:
                line = ln.strip()
                break
        if code == 0:
            out.append(Leg("WORKFLOWS", f[:-3], PASS, line, "", secs))
        elif code == 3:
            out.append(Leg("WORKFLOWS", f[:-3], ABSENT,
                           "the door is not answering",
                           "start atlas-mcp on :8090", secs))
        else:
            out.append(Leg("WORKFLOWS", f[:-3], FAIL, line or tail(text),
                           "", secs))
    return out


def leg_e2e(live):
    """The end-to-end suite: the door AND a local model."""
    path = os.path.join(HERE, "e2e", "test_suite.py")
    if not live:
        return [Leg("E2E", "test_suite", SKIPPED, "pass --live to run")]
    if not mcp_online():
        return [Leg("E2E", "test_suite", ABSENT, "the door is not answering",
                    "start atlas-mcp on :8090")]
    t = time.time()
    code, text = run([sys.executable, path], timeout=3600)
    secs = time.time() - t
    if code == 0:
        return [Leg("E2E", "test_suite", PASS, "", "", secs)]
    return [Leg("E2E", "test_suite", FAIL, tail(text), "", secs)]


# THE GITHUB WORKFLOWS THEMSELVES (2026-10-09, the core's WHAT'S LEFT I1, sixth
# step). Read as text, the way the core's release gate reads its own (its
# tests/release.py, `workflows`), and held to the same two rules. An action from
# outside GitHub's own `actions/` is named by its 40-hex commit: a tag is a
# pointer its owner can move, and release.yml's is the one step that holds
# `contents: write`. And a line is proved once: a workflow that runs on pull
# requests does not also run on a push to every line, as prove.yml's bare
# `push:` did -- the same commit twice on every line with a pull request open.
USES = re.compile(r"^\s*(?:-\s+)?uses:\s*[\"']?([^\"'\s#]+)")
SHA_PIN = re.compile(r"^[0-9a-f]{40}$")
ON_KEY = re.compile(r"^[\"']?on[\"']?:\s*(.*)$")
SUB_KEY = re.compile(r"^([\w-]+):\s*(.*)$")


def triggers(text):
    """Each trigger in a workflow's top-level `on:`, with the lines under it
    stripped, read by indentation; `on: push` and `on: [push, pull_request]`
    as well as the block."""
    lines = text.splitlines()
    for i, line in enumerate(lines):
        m = ON_KEY.match(line)
        if not m:
            continue
        rest = m.group(1).split("#")[0].strip()
        if rest:
            return dict((t.strip(" \"'"), []) for t in rest.strip("[]").split(",")
                        if t.strip())
        out, key, depth = {}, "", 0
        for body in lines[i + 1:]:
            if not body.strip() or body.lstrip().startswith("#"):
                continue
            ind = len(body) - len(body.lstrip())
            if not ind:
                break
            if not depth or ind <= depth:
                depth, key = ind, body.strip().lstrip("- ").split(":")[0].strip()
                out[key] = []
            else:
                out[key].append(body.strip())
        return out
    return {}


def pushes_every_line(on):
    """Whether a push to any line of work runs the workflow: a `push` with no
    `branches` to limit it (or only a bare `*`), and not one limited to marks."""
    body = on.get("push")
    if body is None:
        return False
    for j, line in enumerate(body):
        m = SUB_KEY.match(line)
        if m and m.group(1) == "branches":
            named = m.group(2).split("#")[0].strip()
            items = named.strip("[]").split(",") if named else []
            for item in ([] if named else body[j + 1:]):
                if not item.startswith("-"):
                    break
                items.append(item)
            return any(t.strip(" -\"'") and not t.strip(" -\"'*") for t in items)
    keys = set(m.group(1) for m in map(SUB_KEY.match, body) if m)
    return "branches-ignore" in keys or not keys & set(["tags", "tags-ignore"])


def leg_github():
    """The workflows that prove atlas on GitHub, held to the two rules above.
    Static and hermetic: never ABSENT, it holds or it broke."""
    wfdir = os.path.join(ATLAS, ".github", "workflows")
    names = sorted(f for f in os.listdir(wfdir)
                   if f.endswith((".yml", ".yaml"))) if os.path.isdir(wfdir) else []
    if not names:
        return [Leg("GITHUB", ".github/workflows", FAIL,
                    "no workflow here: nothing proves a line before main takes it")]
    out = []
    for fn in names:
        with open(os.path.join(wfdir, fn), encoding="utf-8", errors="replace") as f:
            text = f.read()
        used = [m.group(1) for m in map(USES.match, text.splitlines())
                if m and not m.group(1).startswith(("./", "actions/"))]
        faults = ["%s names a tag its owner can move; an action from outside "
                  "GitHub's own is named by its 40-hex commit" % a
                  for a in used if not SHA_PIN.match(a.partition("@")[2])]
        on = triggers(text)
        if "pull_request" in on and pushes_every_line(on):
            faults.append("it runs on pull requests AND on a push to every line, so "
                          "a line with one open is proved twice; limit the push to "
                          "main (branches: [main])")
        if faults:
            out.append(Leg("GITHUB", fn, FAIL, "\n".join(faults)))
            continue
        said = ("%d outside action%s, each named by its commit"
                % (len(used), "" if len(used) == 1 else "s")) if used else "no outside action"
        out.append(Leg("GITHUB", fn, PASS, said + "; no line proved twice"))
    return out


def report(legs, quiet):
    width = max(len(l.name) for l in legs)
    group = None
    for l in legs:
        if l.group != group:
            group = l.group
            print()
            print("  %s" % group)
        lines = l.detail.splitlines() if l.detail else [""]
        secs = "  %5.1fs" % l.secs if l.secs >= 0.05 else "        "
        print("    [%-7s]  %-*s%s  %s"
              % (l.verdict, width, l.name, secs, lines[0]))
        if not quiet:
            for ln in lines[1:]:
                print("                 %s" % ln)
        if l.need:
            print("                 needs: %s" % l.need)


def main():
    ap = argparse.ArgumentParser(description="prove the whole of atlas")
    ap.add_argument("--check", action="store_true",
                    help="fast legs only: go + goldens, no cargo")
    ap.add_argument("--live", action="store_true",
                    help="also run the legs that need :8090 and Ollama")
    ap.add_argument("--quiet", action="store_true", help="one line per leg")
    args = ap.parse_args()

    t0 = time.time()
    print()
    print("  ATLAS -- THE BALL (%s)"
          % ("fast" if args.check else "live" if args.live else "hermetic"))

    legs = leg_toolchains()
    if not args.check:
        legs += leg_spine()
    legs += leg_go("line", "LINE")
    legs += leg_go("webapp", "GLASS")
    if not args.check:
        legs += leg_mcp()
    legs += leg_goldens()
    legs += leg_chain_verify()
    legs += leg_github()
    legs += leg_workflows(args.live)
    legs += leg_e2e(args.live)

    report(legs, args.quiet)

    broke = [l for l in legs if l.verdict == FAIL]
    gone = [l for l in legs if l.verdict == ABSENT]
    held = [l for l in legs if l.verdict == PASS]
    print()
    print("  %d held - %d absent - %d broke - %.1fs"
          % (len(held), len(gone), len(broke), time.time() - t0))
    if gone and not args.quiet:
        print()
        print("  ABSENT is not a pass. These legs named something this ground")
        print("  does not hold; nothing was assumed in their place:")
        seen = set()
        for l in gone:
            if not l.need or l.need in seen:
                continue
            seen.add(l.need)
            print("    %s" % l.need)
    print()
    if broke:
        print("  A leg broke. atlas does not claim what it cannot show.")
        return 1
    print("  PROVEN, to the edge of this ground.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
