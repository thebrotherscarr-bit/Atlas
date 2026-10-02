#!/usr/bin/env python3
"""chain_verify.py -- the chain verdict EMPTY | INTACT | FLIP | TAMPER, in Python.

    python tools/chain_verify.py chain verify <path>    one line, the spine's shape; exit 0 while
                                                        the chain stays appendable, 1 when not,
                                                        2 when refused
    python tools/chain_verify.py --json <path>          the verdict as the oracle's dict
    python tools/chain_verify.py --goldens              every chain and injection in
                                                        tests/fixtures/canon/chain_verdicts.json,
                                                        re-verdicted here and compared; exit 1 on
                                                        any difference (prove.py's own leg for it;
                                                        this is a verifier, not a cutter -- it
                                                        writes nothing and takes no cutter's word)

WHY (2026-10-01; his ruling of the same morning, the core's WHAT'S LEFT B14). The door's
`verify_chain` shelled the Rust spine (`atlas chain verify <path>`), and that spine is itself "a
byte-faithful mirror of the oracle walk" -- `us_chain.py :: verify`, a Python prover that lives
outside this ground. He ruled: the verdict in Python first, proved against the 21 golden
masters and the 8 injections the oracle cut, then the Rust folded to an attic. This is the
verdict, written from the Rust in the ground (`core/src/chain.rs`, `canon.rs`), which is the
oracle's shape line for line:

  * a chain file is non-empty lines; an unparsable line is TAMPER at its index;
  * each entry's hashed view is BODY_KEYS -- ts, kind, n, payload, prev, actor, body_v -- where
    present; its hash is sha256(prev_hex + canon(view)) under the entry's own `body_v`
    (1 the board shape, `json.dumps(sort_keys=True)`; 2 the links shape, the same with raw
    UTF-8; 3 JCS, RFC 8785 -- UTF-16 key order, no whitespace, floats and out-of-range
    integers refused), defaulting to 3 when unstamped;
  * the weld follows each entry's STORED hash: a stored `prev` that is not the previous stored
    hash is TAMPER at that index (the weld broke); a hash that does not recompute while the
    weld holds is a FLIP, named, and the chain stays appendable;
  * an empty file is EMPTY, head GENESIS, appendable.

ONE PLACE THE RUST AND THE ORACLE PART, and the oracle is followed: a chain whose parsed rows
carry no `hash` field at all is not verified -- the oracle's walk was written for the stamped
world and the cutter recorded those three chains as SKIPPED with the reason. The Rust would call
them TAMPER at entry 1. The goldens say skipped; so does this.

Stdlib only. Reads the chain; writes nothing.
"""

from __future__ import annotations

import hashlib
import json
import sys
from pathlib import Path

GENESIS = "0" * 64
BODY_KEYS = ("ts", "kind", "n", "payload", "prev", "actor", "body_v")
BODY_V_BOARD, BODY_V_LINKS, BODY_V_JCS = 1, 2, 3
SAFE_INT_MAX = 9_007_199_254_740_991
SKIPPED_WHY = "rows without hash fields (us_chain.verify requires one per entry)"


class CanonRefused(Exception):
    """A body that cannot be canonicalised under its form (floats and big integers under JCS,
    an unknown body_v). The oracle raises mid-verify; the walk reads it as TAMPER at the entry."""


# ---- canon ------------------------------------------------------------------------------------

def _jcs_string(s: str) -> str:
    out = ['"']
    for ch in s:
        o = ord(ch)
        if ch == '"':
            out.append('\\"')
        elif ch == "\\":
            out.append("\\\\")
        elif ch == "\b":
            out.append("\\b")
        elif ch == "\f":
            out.append("\\f")
        elif ch == "\n":
            out.append("\\n")
        elif ch == "\r":
            out.append("\\r")
        elif ch == "\t":
            out.append("\\t")
        elif o < 0x20:
            out.append("\\u%04x" % o)
        else:
            out.append(ch)          # U+007F and every non-ASCII character stays literal
    out.append('"')
    return "".join(out)


def _jcs(v, out: list) -> None:
    if v is None:
        out.append("null")
    elif v is True:
        out.append("true")
    elif v is False:
        out.append("false")
    elif isinstance(v, int):
        if v > SAFE_INT_MAX or v < -SAFE_INT_MAX:
            raise CanonRefused(f"integer {v} is outside the ECMAScript safe range")
        out.append(str(v))
    elif isinstance(v, float):
        raise CanonRefused("float refused under BODY_V 3")
    elif isinstance(v, str):
        out.append(_jcs_string(v))
    elif isinstance(v, list):
        out.append("[")
        for i, item in enumerate(v):
            if i:
                out.append(",")
            _jcs(item, out)
        out.append("]")
    elif isinstance(v, dict):
        out.append("{")
        # RFC 8785: keys in the order of their UTF-16 code units
        for i, k in enumerate(sorted(v, key=lambda s: s.encode("utf-16-be"))):
            if i:
                out.append(",")
            out.append(_jcs_string(k))
            out.append(":")
            _jcs(v[k], out)
        out.append("}")
    else:
        raise CanonRefused(f"no canon for {type(v).__name__}")


def canon(body, body_v: int) -> str:
    """Canonical TEXT for a body under the named form (the Rust's canon(), the oracle's us_canon)."""
    if body_v == BODY_V_BOARD:
        return json.dumps(body, sort_keys=True)
    if body_v == BODY_V_LINKS:
        return json.dumps(body, sort_keys=True, ensure_ascii=False)
    if body_v == BODY_V_JCS:
        out: list = []
        _jcs(body, out)
        return "".join(out)
    raise CanonRefused(f"no such BODY_V: {body_v} (known: 1, 2, 3)")


# ---- the walk ---------------------------------------------------------------------------------

def read_rows(text: str) -> list:
    """One item per non-empty line: the parsed object, or None where the line is not JSON."""
    rows = []
    for line in text.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            rows.append(json.loads(line))
        except ValueError:
            rows.append(None)
    return rows


def project(entry: dict) -> dict:
    return {k: entry[k] for k in BODY_KEYS if k in entry} if isinstance(entry, dict) else entry


def entry_hash(prev: str, entry: dict) -> str:
    body_v = entry.get("body_v", BODY_V_JCS) if isinstance(entry, dict) else BODY_V_JCS
    if body_v is None or isinstance(body_v, bool) or not isinstance(body_v, int) or not 0 <= body_v <= 255:
        raise CanonRefused("body_v is not a small integer")
    text = canon(project(entry), body_v)
    return hashlib.sha256((prev + text).encode("utf-8")).hexdigest()


def _stored(entry, key: str):
    v = entry.get(key) if isinstance(entry, dict) else None
    return v if isinstance(v, str) else None


def _flip_why(flips: list) -> str:
    n = len(flips)
    return (f"{n} {'entry' if n == 1 else 'entries'} do not hash to their stored value while the weld "
            f"holds end to end — corruption, not a rewrite. Restore "
            + ", ".join(f"#{i}" for i in flips)
            + " from a backup copy; the head is sound and the chain stays appendable.")


def verify_rows(rows: list) -> dict:
    """The oracle's dict: verdict, entries, flips, appendable, head/broke_at where they apply, why."""
    if not rows:
        return {"verdict": "EMPTY", "entries": 0, "flips": [], "appendable": True, "head": GENESIS, "why": None}
    for i, row in enumerate(rows):
        if row is None:
            return {"verdict": "TAMPER", "entries": len(rows), "flips": [], "appendable": False,
                    "broke_at": i, "why": "line is not JSON"}
    if any(not (isinstance(r, dict) and r.get("hash")) for r in rows):
        return {"verify_skipped": SKIPPED_WHY}
    prev = GENESIS
    flips: list = []
    for i, entry in enumerate(rows):
        claimed = _stored(entry, "prev")
        if claimed != prev:
            return {"verdict": "TAMPER", "entries": len(rows), "flips": flips, "appendable": False,
                    "broke_at": i,
                    "why": f"weld broken: entry {i} claims prev {(claimed or '')[:12]}, the chain says {prev[:12]}"}
        try:
            expect = entry_hash(prev, entry)
        except CanonRefused as exc:
            return {"verdict": "TAMPER", "entries": len(rows), "flips": flips, "appendable": False,
                    "broke_at": i, "why": str(exc)}
        stored = _stored(entry, "hash")
        if stored != expect:
            flips.append(i)
        prev = stored            # the weld follows the STORED hash
    if flips:
        return {"verdict": "FLIP", "entries": len(rows), "flips": flips, "appendable": True,
                "head": prev, "why": _flip_why(flips)}
    return {"verdict": "INTACT", "entries": len(rows), "flips": [], "appendable": True, "head": prev, "why": None}


def verify_text(text: str) -> dict:
    return verify_rows(read_rows(text))


def verify_path(path: Path) -> dict:
    return verify_text(path.read_bytes().decode("utf-8", errors="replace"))


# ---- the spine's line ---------------------------------------------------------------------------

def spine_line(label: str, v: dict) -> tuple[bool, str]:
    """`<label>: verdict=... entries=... flips=[...] broke_at=None|Some(i) appendable=true|false`,
    the exact shape the Rust CLI printed, so the door's callers read nothing new. A skipped
    verify prints its own verdict word and reason, and is not appendable until it is read."""
    if "verify_skipped" in v:
        return False, f"{label}: verdict=SKIPPED why={v['verify_skipped']} appendable=false"
    broke = "None" if v.get("broke_at") is None else f"Some({v['broke_at']})"
    flips = "[" + ", ".join(str(i) for i in v["flips"]) + "]"
    ok = bool(v["appendable"])
    return ok, (f"{label}: verdict={v['verdict']} entries={v['entries']} flips={flips} "
                f"broke_at={broke} appendable={'true' if ok else 'false'}")


# ---- the goldens --------------------------------------------------------------------------------

def goldens(atlas: Path) -> int:
    doc = json.loads((atlas / "tests" / "fixtures" / "canon" / "chain_verdicts.json").read_text(encoding="utf-8"))
    sys.path.insert(0, str(atlas / "tools"))
    import cut_chain_verdicts as cutter   # the recipes, not the oracle (which it never loads at import)
    bad, seen = [], 0
    for c in doc["chains"]:
        seen += 1
        mine = verify_path(atlas / "tests" / "fixtures" / c["file"])
        if mine != c["verify"]:
            bad.append((c["file"], c["verify"], mine))
    for inj in doc["injections"]:
        seen += 1
        raw = (atlas / "tests" / "fixtures" / inj["source"]).read_bytes().decode("utf-8")
        lines = [ln for ln in raw.split("\n") if ln.strip()]
        text = "\n".join(cutter.apply_op(list(lines), dict(inj["recipe"]))) + "\n"
        mine = verify_text(text)
        if mine != inj["verify"]:
            bad.append((f"{inj['source']} :: {inj['recipe']['op']}", inj["verify"], mine))
    for name, want, got in bad:
        print(f"DIFFERS {name}\n   golden {json.dumps(want, ensure_ascii=False)[:200]}\n   here   {json.dumps(got, ensure_ascii=False)[:200]}")
    print(f"{seen - len(bad)} of {seen} verdicts match the goldens ({doc['tally']['chains']} chains, "
          f"{doc['tally']['injections']} injections)" + ("" if not bad else f"; {len(bad)} DIFFER"))
    return 1 if bad else 0


def main(argv: list) -> int:
    atlas = Path(__file__).resolve().parent.parent
    if argv[:1] == ["--goldens"]:
        return goldens(atlas)
    if argv[:1] == ["--json"] and len(argv) == 2:
        print(json.dumps(verify_path(Path(argv[1])), ensure_ascii=False))
        return 0
    if argv[:2] == ["chain", "verify"] and len(argv) == 3:
        path = Path(argv[2])
        if not path.is_file():
            print(f"refused: cannot read {path}", file=sys.stderr)
            return 2
        ok, line = spine_line(argv[2], verify_path(path))
        print(line)
        return 0 if ok else 1
    print("refused: chain_verify takes `chain verify <path>`, `--json <path>` or `--goldens`", file=sys.stderr)
    return 2


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
