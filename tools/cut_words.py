#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
cut_words.py -- THE WORDS A CUTTER KNOWS. A word it does not know writes nothing.

Every cutter in tools/ answers two words: none at all (or `--cut`), which CUTS
its fixture, and `--verify`, which cuts again in memory and compares. Until
2026-09-29 each of them read its word as `"--verify" in sys.argv` and cut on
ANYTHING ELSE -- a typo, another tool's flag, `--check` typed by a hand that
meant verify -- and a cut is a write to a tracked fixture. It happened that
day to the flow fixture, which was put back from the last save.

    from cut_words import word
    sys.exit(verify() if word(sys.argv[1:], __file__) == "--verify" else cut())

`word` answers "--verify" or "--cut". For anything else it prints the usage
and exits 2, having written nothing.

RUN AS A SCRIPT, THIS IS THE PROOF OF IT. prove.py runs every file in tools/
that answers --verify, so this is one of its legs:

    python tools/cut_words.py --verify

WHICH SCRIPTS ARE ASKED. Not a list and not a name: every script beside this
file that takes its word off `sys.argv` and knows `--verify` -- a script with
a second mode, and the second mode writes. Twenty-three of them are named
cut_*_vectors.py; cut_chain_verdicts.py and fold_agents.py are not, and had
the same fault. A script that parses its words with argparse is refused an
unknown one by argparse itself and is not asked here.

For each: its __main__ block asks `word` before it does anything -- READ FROM
ITS SOURCE, so a cutter that does not ask is never run with a word that would
make it cut -- and, run with a word nobody knows, it exits 2 and leaves every
file under tests/fixtures and agents as it was, byte for byte.

Zero external dependencies (law 6): stdlib only.
"""

import hashlib
import os
import subprocess
import sys

KNOWN = ("--verify", "--cut")
UNKNOWN = "--a-word-nobody-knows"

TOOLS = os.path.dirname(os.path.abspath(__file__))
ATLAS = os.path.normpath(os.path.join(TOOLS, ".."))
# Where a cut lands: the fixtures, and the declarations fold_agents.py folds.
WRITTEN = (os.path.join(ATLAS, "tests", "fixtures"), os.path.join(ATLAS, "agents"))


def word(argv, script=""):
    """The one word this run was given: "--verify", or "--cut" (said, or
    nothing said). A word that is neither ends the run here, exit 2."""
    args = list(argv or [])
    unknown = [a for a in args if a not in KNOWN]
    if unknown:
        name = os.path.basename(script or "cutter")
        print("usage: %s            cut the fixture\n"
              "       %s --cut      the same\n"
              "       %s --verify   prove it, writing nothing\n"
              "refused: %r is not one of those; nothing was written"
              % (name, name, name, unknown))
        sys.exit(2)
    return "--verify" if "--verify" in args else "--cut"


def cutters():
    """Every script beside this one that takes its word off sys.argv and knows
    `--verify`, whatever it is named."""
    out = []
    for f in sorted(os.listdir(TOOLS)):
        if not f.endswith(".py") or f == os.path.basename(__file__):
            continue
        with open(os.path.join(TOOLS, f), encoding="utf-8", errors="replace") as fh:
            source = fh.read()
        if "sys.argv" in source and "--verify" in source:
            out.append(f)
    return out


def main_block(source):
    """The text of a script's `if __name__ == "__main__":` block."""
    mark = 'if __name__ == "__main__":'
    at = source.rfind(mark)
    return source[at:] if at != -1 else ""


def asks(source):
    """Does this cutter ask `word` before it does anything? Read, not run."""
    block = main_block(source)
    return ("from cut_words import word" in block
            and "word(sys.argv[1:], __file__)" in block
            and " in sys.argv" not in block)


def fingerprint():
    """Every file a cut could land on, by path, with the hash of its bytes."""
    out = {}
    for top in WRITTEN:
        for root, _dirs, files in os.walk(top):
            for f in files:
                path = os.path.join(root, f)
                with open(path, "rb") as fh:
                    out[os.path.relpath(path, ATLAS)] = hashlib.sha256(fh.read()).hexdigest()
    return out


def verify():
    print()
    print("  CUTTER WORDS -- a word a cutter does not know writes nothing")
    print()
    names = cutters()
    width = max(len(n) for n in names) if names else 10
    ok = bool(names)
    before = fingerprint()
    for name in names:
        path = os.path.join(TOOLS, name)
        with open(path, encoding="utf-8", errors="replace") as fh:
            source = fh.read()
        if not asks(source):
            # NOT RUN: it would cut on the word this proof is about to hand it.
            print("    [FAIL]  %-*s  does not ask `word` before it acts; not run" % (width, name))
            ok = False
            continue
        p = subprocess.run([sys.executable, path, UNKNOWN], cwd=ATLAS,
                           stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
                           stderr=subprocess.STDOUT, timeout=120)
        said = p.stdout.decode("utf-8", "replace")
        after = fingerprint()
        moved = sorted(k for k in set(before) | set(after) if before.get(k) != after.get(k))
        if p.returncode == 2 and "nothing was written" in said and not moved:
            print("    [PASS]  %-*s  refused %s, exit 2, nothing moved"
                  % (width, name, UNKNOWN))
        else:
            print("    [FAIL]  %-*s  exit %d%s" % (
                width, name, p.returncode,
                "; MOVED: " + ", ".join(moved[:4]) if moved else ""))
            ok = False
        before = after
    print()
    if ok:
        print("  PROVEN. %d cutters, and none of them cuts on a word it does not know."
              % len(names))
        return 0
    print("  A cutter would write on a word it does not know.")
    return 1


if __name__ == "__main__":
    if word(sys.argv[1:], __file__) == "--verify":
        sys.exit(verify())
    print("cut_words.py cuts nothing: it holds the words the cutters know.\n"
          "       cut_words.py --verify   prove that every cutter refuses a word it does not know")
    sys.exit(0)
