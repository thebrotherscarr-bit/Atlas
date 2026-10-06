#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
cut_flow_vectors.py -- N2 workflow goldens (spec-first).

The flow v1 contract (the oracle; the Go flow package must honor it):

  names         : ^[a-z0-9][a-z0-9_-]{0,63}$ for flows, nodes, runs carry
                  f-YYYYMMDD-HHMMSS-<8hex> (RUN_RE pinned below)
  node kinds    : ask | prompt | seat | memory | eval | gate | run | aider |
                  tool —
                  a closed set. A `run` node drives a whole Manjuel turn (the
                  council), so it must carry an objective or it refuses; the
                  others reach one voice. An `aider` node (2026-10-05, H17) is
                  an attempt by Aider on the files it is handed, and the
                  council's own turn when Aider cannot take it: it carries an
                  instruction (`question`) and `files`, no other kind names
                  files, and it is reached only through a gate whose `grants`
                  name `aider_run` on every path (the nearest gate wins)
                  A `tool` node (2026-10-05) calls one of the door's own
                  tools by name: it names the tool in the door's shape
                  (TOOL_RE), its `args` are named the same way, and no
                  other kind carries `tool` or `args`; which tools exist,
                  and whether one that writes was granted, is the door's
  edges         : {from, to, when: always|pass|fail}; fail-edges only from
                  eval/gate nodes; everything else with when:fail refuses
  validation    : unique names, known kinds, refs resolve, no cycles,
                  exactly one start, all nodes reachable from the start
  the return    : LAW_003, 2026-09-28. A node that does work may declare
                  `loops` (0 to 5); an eval's fail-edge INTO such a node is a
                  RETURN, not a forward edge, and leaves Kahn's count. One
                  return per check; it goes BACK; the body it goes around
                  holds no gate and does work; a declared ceiling is reached
                  by something; `loops` never sits on an eval or a gate
  order         : Kahn over name-sorted ready sets (deterministic)
  branches      : eval pass/fail (exact trim+casefold) and gate
                  continue/stop steer which out-edges fire; parallelism is
                  declared but runs sequentially, topo order (one rack queue)
  templates     : {{var}} from run inputs plus {{out_<node>}} from finished
                  nodes; missing vars fail the node, never guess
  run receipt   : sha256(run + "\\n" + node + "\\n" + output + "\\n" + ts)
  budget        : wall sum vs budget_s (default 600); over() is pure and
                  pinned here; OUT_OF_TIME stops the run honestly

No live Ollama, no network. The fixture file is the oracle; --verify
recomputes every vector deterministically with canned node outputs.
"""

import hashlib
import json
import os
import re
import sys

SCRIPT = os.path.dirname(os.path.abspath(__file__))
ATLAS = os.path.normpath(os.path.join(SCRIPT, ".."))
FIX = os.path.join(ATLAS, "tests", "fixtures")
FLOW = os.path.join(FIX, "flow_vectors.json")

NAME_RE = r"^[a-z0-9][a-z0-9_-]{0,63}$"
RUN_RE = r"^f-\d{8}-\d{6}-[0-9a-f]{8}$"
KINDS = ["ask", "prompt", "seat", "memory", "eval", "gate", "run", "aider", "tool"]
AIDER_TOOL = "aider_run"    # flow.go AiderTool
TOOL_RE = r"^[a-z][a-z0-9_]{0,63}$"    # flow.go ToolRe
MAX_LOOPS = 5          # flow.go MaxLoops


def reach(start, forward, up):
    """Every node reachable from `start` over the forward edges -- walked
    backwards when `up` -- the start itself not included (flow.go's reach)."""
    out, queue = set(), [start]
    while queue:
        cur = queue.pop(0)
        for frm, to in forward:
            nxt = to if (not up and frm == cur) else (frm if (up and to == cur) else None)
            if nxt is None or nxt in out or nxt == start:
                continue
            out.add(nxt)
            queue.append(nxt)
    return out


def topo(nodes, edges):
    names = [n["name"] for n in nodes]
    if len(set(names)) != len(names):
        raise ValueError("duplicate node name")
    if len(names) == 0:
        raise ValueError("empty flow")
    kinds = {n["name"]: n["kind"] for n in nodes}
    loops = {}
    for n in nodes:
        if n["kind"] not in KINDS:
            raise ValueError("unknown kind: " + n["kind"])
        if n["kind"] == "run" and not (n.get("question") or "").strip():
            raise ValueError("run node with no objective")
        if n["kind"] == "aider" and not (n.get("question") or "").strip():
            raise ValueError("aider node with no instruction")
        if n["kind"] == "aider" and not (n.get("files") or "").strip():
            raise ValueError("aider node that names no files")
        if n["kind"] != "aider" and (n.get("files") or "").strip():
            raise ValueError("only an aider node names files")
        # A `tool` NODE (2026-10-05): the tool and its arguments named in the
        # door's shape, and no other kind calls a tool (flow.go Validate).
        if n["kind"] == "tool" and not re.fullmatch(TOOL_RE, n.get("tool") or ""):
            raise ValueError("a tool node that names no tool the door could carry")
        if n["kind"] != "tool" and ((n.get("tool") or "").strip() or n.get("args")):
            raise ValueError("only a tool node calls a tool")
        if any(not re.fullmatch(TOOL_RE, k) for k in (n.get("args") or {})):
            raise ValueError("an argument named outside the door's shape")
        k = n.get("loops") or 0
        if not isinstance(k, int) or k < 0 or k > MAX_LOOPS:
            raise ValueError("loops out of range: %r" % (k,))
        if k > 0 and n["kind"] in ("eval", "gate"):
            raise ValueError("loops on a check or a gate: a loop re-does WORK")
        loops[n["name"]] = k
    incoming = {n: 0 for n in names}
    adj = {n: [] for n in names}
    returning = {}          # check -> the node it returns to
    forward = []
    for e in edges:
        if e["from"] not in incoming or e["to"] not in incoming:
            raise ValueError("edge refs unknown node")
        # THE RETURN (LAW_003). A check's fail-edge into a node that declares
        # `loops` leaves Kahn's count and is judged below; with no ceiling on
        # the node it is the cycle it always was.
        if (e.get("when", "always") == "fail" and kinds[e["from"]] == "eval"
                and loops[e["to"]] > 0):
            if e["from"] in returning:
                raise ValueError("a check returns once")
            returning[e["from"]] = e["to"]
            continue
        if e.get("when", "always") == "fail" and kinds[e["from"]] not in ("eval", "gate"):
            raise ValueError("fail-edge only from eval/gate")
        if e.get("when", "always") not in ("always", "pass", "fail"):
            raise ValueError("bad when: " + str(e.get("when")))
        incoming[e["to"]] += 1
        adj[e["from"]].append(e)
        forward.append((e["from"], e["to"]))
    starts = sorted(n for n in names if incoming[n] == 0)
    if len(starts) != 1:
        raise ValueError("want exactly one start, got %d" % len(starts))
    ready = sorted(starts)
    order = []
    indeg = dict(incoming)
    while ready:
        n = ready.pop(0)
        order.append(n)
        for e in sorted(adj[n], key=lambda x: x["to"]):
            indeg[e["to"]] -= 1
            if indeg[e["to"]] == 0:
                ready.append(e["to"])
        ready.sort()
    if len(order) != len(names):
        raise ValueError("cycle or unreachable node")
    pos = {n: i for i, n in enumerate(order)}
    reached = set()
    for chk in sorted(returning):
        to = returning[chk]
        if pos[to] >= pos[chk]:
            raise ValueError("a return that does not go back")
        body = {to, chk} | (reach(to, forward, False) & reach(chk, forward, True))
        if any(kinds[b] == "gate" for b in body):
            raise ValueError("gate inside the return")
        if not any(kinds[b] not in ("eval", "gate") for b in body):
            raise ValueError("the return re-does no work")
        reached.add(to)
    for n in names:
        if loops[n] > 0 and n not in reached:
            raise ValueError("a ceiling nothing returns to")
    # AIDER IS REACHED ONLY THROUGH A GATE THAT GRANTS IT (H17): the nearest gate
    # before the node on every forward path must grant AIDER_TOOL, and a path that
    # meets the start before any gate was never authorised (flow.go aiderGranted).
    grants = {n["name"]: n.get("grants") or [] for n in nodes}
    ups = {n: [] for n in names}
    for frm, to in forward:
        ups[to].append(frm)
    memo = {}

    def covered(name):
        if name in memo:
            return memo[name]
        ok = len(ups[name]) > 0
        for up in ups[name]:
            if kinds[up] == "gate":
                if AIDER_TOOL not in grants[up]:
                    ok = False
                continue
            if not covered(up):
                ok = False
        memo[name] = ok
        return ok

    for n in names:
        if kinds[n] == "aider" and not covered(n):
            raise ValueError("an aider node no gate authorised")
    return order


def run_outcome(spec, canned, eval_expected, gate_decision):
    """Simulate branch routing: canned outputs per node, eval expected map,
    gate decision map. Returns fired node names in order + verdict."""
    order = topo(spec["nodes"], spec.get("edges", []))
    byname = {n["name"]: n for n in spec["nodes"]}
    edges = spec.get("edges", [])
    fired = []
    verdict = "DONE"
    passed = {}
    for n in order:
        nd = byname[n]
        # gate: does any in-edge fire?
        fire = True
        for e in edges:
            if e["to"] != n:
                continue
            src = e["from"]
            if src not in fired:
                continue
            when = e.get("when", "always")
            if when == "always":
                break
            if when == "pass" and passed.get(src) is True:
                break
            if when == "fail" and passed.get(src) is False:
                break
        else:
            # no in-edge fired (start has none -> fires)
            fire = any(e["to"] == n for e in edges) is False
        if n == order[0]:
            fire = True
        if not fire:
            continue
        fired.append(n)
        if nd["kind"] == "eval":
            exp = (eval_expected or {}).get(n, "")
            got = canned.get(nd.get("node", ""), "")
            ok = exp.strip().casefold() == got.strip().casefold()
            passed[n] = ok
            # fail with no fail-edge -> run FAILs here
            if not ok and not any(e["from"] == n and e.get("when") == "fail" for e in edges):
                verdict = "FAIL"
                break
        elif nd["kind"] == "gate":
            decision = (gate_decision or {}).get(n, "continue")
            passed[n] = (decision == "continue")
            if decision == "stop":
                # stop with no fail-edge -> STOPPED here
                if not any(e["from"] == n and e.get("when") == "fail" for e in edges):
                    verdict = "STOPPED"
                    break
            else:
                verdict = "PAUSED" if n == order[-1] or True else verdict
                # gate pauses the run unless resumed; simulation marks PAUSED
                verdict = "PAUSED"
                break
    return {"fired": fired, "verdict": verdict}


def receipt(run, node, output, ts):
    h = hashlib.sha256()
    h.update((run + "\n" + node + "\n" + output + "\n" + ts).encode("utf-8"))
    return h.hexdigest()


def over(elapsed_ms, budget_s):
    return sum(elapsed_ms) > budget_s * 1000


def spec(nodes, edges):
    return {"nodes": nodes, "edges": edges}


def vectors():
    ask = lambda n, q="Q": {"name": n, "kind": "ask", "question": q}
    ev = lambda n, ref, exp: {"name": n, "kind": "eval", "node": ref, "expected": exp}
    gate = lambda n: {"name": n, "kind": "gate", "title": "review"}
    E = lambda f, t, w="always": {"from": f, "to": t, "when": w}
    linear = spec([ask("a"), ask("b")], [E("a", "b")])
    branch = spec([ask("a"), ev("e", "a", "yes"), ask("b"), ask("c")],
                  [E("a", "e"), E("e", "b", "pass"), E("e", "c", "fail")])
    gated = spec([ask("a"), gate("g"), ask("b")], [E("a", "g"), E("g", "b")])
    # THE RETURN'S VECTORS (LAW_003, 2026-09-28): a node that does work, a
    # check on it that matches by containment, and the check's fail-edge back.
    work = lambda n, k=0: dict({"name": n, "kind": "run", "question": "do"},
                               **({"loops": k} if k else {}))
    chk = lambda n, ref: {"name": n, "kind": "eval", "node": ref,
                          "expected": "RAN:", "match": "contains"}
    back = [E("w", "c"), E("c", "w", "fail")]
    # THE AIDER NODE'S VECTORS (H17, 2026-10-05): the node, the gate that grants it,
    # and what is refused about it.
    aider = lambda n, files="a.py", q="do": dict({"name": n, "kind": "aider", "question": q}, **({"files": files} if files else {}))
    granting = lambda n, *g: {"name": n, "kind": "gate", "title": "t", "grants": list(g)}
    return {
        "name_re": NAME_RE,
        "run_re": RUN_RE,
        "kinds": KINDS,
        "topo_linear": {"order": topo(linear["nodes"], linear["edges"])},
        "branch_pass": run_outcome(branch, {"a": "yes"}, {"e": "yes"}, {}),
        "branch_fail": run_outcome(branch, {"a": "no"}, {"e": "yes"}, {}),
        "gate_pauses": run_outcome(gated, {"a": "hi"}, {}, {}),
        "receipt_example": {
            "run": "f-20260909-120000-01234567",
            "node": "a",
            "output": "the ledger holds",
            "ts": "2026-09-09T12:00:00Z",
            "receipt": receipt("f-20260909-120000-01234567", "a",
                               "the ledger holds", "2026-09-09T12:00:00Z"),
        },
        "budget_cases": [
            {"elapsed_ms": [100, 200], "budget_s": 600, "over": False},
            {"elapsed_ms": [300000, 300001], "budget_s": 600, "over": True},
            {"elapsed_ms": [], "budget_s": 600, "over": False},
        ],
        "refusals": [
            {"why": "cycle",
             "spec": spec([ask("a"), ask("b")], [E("a", "b"), E("b", "a")])},
            {"why": "two starts",
             "spec": spec([ask("a"), ask("b")], [])},
            {"why": "unreachable",
             "spec": spec([ask("a"), ask("b"), ask("c")], [E("a", "b")])},
            {"why": "unknown kind",
             "spec": spec([{"name": "a", "kind": "teleport"}], [])},
            {"why": "fail-edge from ask",
             "spec": spec([ask("a"), ask("b")], [E("a", "b", "fail")])},
            {"why": "edge refs unknown",
             "spec": spec([ask("a")], [E("a", "ghost")])},
            {"why": "run node with no objective",
             "spec": spec([{"name": "a", "kind": "run"}], [])},
            {"why": "back-edge with no ceiling (LAW_003: the cycle it always was)",
             "spec": spec([work("w"), chk("c", "w")], back)},
            {"why": "ceiling above MaxLoops",
             "spec": spec([work("w", 6), chk("c", "w")], back)},
            {"why": "ceiling nothing returns to",
             "spec": spec([work("w", 2), chk("c", "w")], [E("w", "c")])},
            {"why": "gate inside the loop",
             "spec": spec([work("w", 2), {"name": "g", "kind": "gate", "title": "t"},
                           chk("c", "w")],
                          [E("w", "g"), E("g", "c", "pass"), E("c", "w", "fail")])},
            {"why": "aider node with no instruction",
             "spec": spec([ask("a"), granting("g", AIDER_TOOL), aider("w", q="")],
                          [E("a", "g"), E("g", "w")])},
            {"why": "aider node that names no files",
             "spec": spec([ask("a"), granting("g", AIDER_TOOL), aider("w", files="")],
                          [E("a", "g"), E("g", "w")])},
            {"why": "files on a node that is not an aider node",
             "spec": spec([dict(ask("a"), files="a.py")], [])},
            {"why": "aider node with no gate before it",
             "spec": spec([ask("a"), aider("w")], [E("a", "w")])},
            {"why": "aider node behind a gate that grants another tool",
             "spec": spec([ask("a"), granting("g", "git_branch"), aider("w")],
                          [E("a", "g"), E("g", "w")])},
            {"why": "aider node behind the granting gate and then a gate that grants nothing",
             "spec": spec([ask("a"), granting("g", AIDER_TOOL), granting("h"), aider("w")],
                          [E("a", "g"), E("g", "h"), E("h", "w")])},
            {"why": "aider node reached by a path around the granting gate",
             "spec": spec([ask("a"), granting("g", AIDER_TOOL), ask("c"), ask("d"), aider("w")],
                          [E("a", "g"), E("g", "c"), E("c", "w"), E("a", "d"), E("d", "w")])},
            {"why": "tool node that names no tool",
             "spec": spec([ask("a"), {"name": "t", "kind": "tool"}], [E("a", "t")])},
            {"why": "tool node whose name breaks the door's shape",
             "spec": spec([ask("a"), {"name": "t", "kind": "tool", "tool": "Git Push"}], [E("a", "t")])},
            {"why": "tool on a node that is not a tool node",
             "spec": spec([dict(ask("a"), tool="git_push")], [])},
            {"why": "args on a node that is not a tool node",
             "spec": spec([dict(ask("a"), args={"project": "research"})], [])},
            {"why": "tool node passing an argument named as the door's own keys are",
             "spec": spec([ask("a"), {"name": "t", "kind": "tool", "tool": "git_push",
                                      "args": {"__caller": "me"}}], [E("a", "t")])},
            {"why": "loops on a check",
             "spec": spec([ask("a"),
                           {"name": "c1", "kind": "eval", "node": "a",
                            "expected": "x", "loops": 1},
                           {"name": "c2", "kind": "eval", "node": "a", "expected": "y"}],
                          [E("a", "c1"), E("c1", "c2", "pass"), E("c2", "c1", "fail")])},
        ],
    }


def lawful_aider():
    """An aider node the law allows, for the verifier to prove the oracle does not
    refuse every aider node: a plan, the gate that grants it, work in between."""
    return spec([{"name": "a", "kind": "ask", "question": "plan"},
                 {"name": "g", "kind": "gate", "title": "open?", "grants": ["git_branch", AIDER_TOOL]},
                 {"name": "c", "kind": "run", "question": "open the line"},
                 {"name": "w", "kind": "aider", "question": "do it", "files": "a.py"}],
                [{"from": "a", "to": "g", "when": "always"},
                 {"from": "g", "to": "c", "when": "pass"},
                 {"from": "c", "to": "w", "when": "always"}])


def lawful_tool():
    """A tool node the law allows, for the verifier to prove the oracle does not
    refuse every tool node: work, then a call with its argument templated."""
    return spec([{"name": "a", "kind": "ask", "question": "plan"},
                 {"name": "t", "kind": "tool", "tool": "git_push", "args": {"project": "{{world}}"}}],
                [{"from": "a", "to": "t", "when": "always"}])


def lawful_return():
    """A return the law allows, for the verifier to prove the oracle does not
    refuse everything that loops: work, a check, and the check's way back."""
    return spec([{"name": "w", "kind": "run", "question": "do", "loops": 2},
                 {"name": "c", "kind": "eval", "node": "w",
                  "expected": "RAN:", "match": "contains"}],
                [{"from": "w", "to": "c", "when": "always"},
                 {"from": "c", "to": "w", "when": "fail"}])


def write_bytes(path, data):
    if isinstance(data, str):
        data = data.encode("utf-8")
    with open(path, "wb") as f:
        f.write(data)


def cut():
    doc = vectors()
    # refusals must actually refuse under our own topo()
    for r in doc["refusals"]:
        try:
            topo(r["spec"]["nodes"], r["spec"].get("edges", []))
            print("REFUSAL HOLE:", r["why"])
            sys.exit(1)
        except ValueError:
            pass
    write_bytes(FLOW, json.dumps(doc, indent=2, sort_keys=True) + "\n")
    print("flow -> %s" % FLOW)


def verify():
    doc = json.loads(open(FLOW, encoding="utf-8").read())
    want = vectors()
    print("\n  FLOW -- golden contract (N2, pinned oracle)")
    ok = True
    if doc != want:
        ok = False
        print("    [FAIL]  fixture drifted from the contract")
        for k in want:
            if doc.get(k) != want[k]:
                print("    drift:", k)
    else:
        print("    [PASS]  topo + branches + gate + receipt + budget + refusals")
    for r in want["refusals"]:
        try:
            topo(r["spec"]["nodes"], r["spec"].get("edges", []))
            ok = False
            print("    [FAIL]  refusal hole: %s" % r["why"])
        except ValueError:
            pass
    try:
        lawful = lawful_return()
        if topo(lawful["nodes"], lawful["edges"]) != ["w", "c"]:
            ok = False
            print("    [FAIL]  a lawful return is ordered wrongly")
    except ValueError as ex:
        ok = False
        print("    [FAIL]  a lawful return is refused: %s" % ex)
    try:
        lawful = lawful_aider()
        if topo(lawful["nodes"], lawful["edges"]) != ["a", "g", "c", "w"]:
            ok = False
            print("    [FAIL]  a lawful aider node is ordered wrongly")
    except ValueError as ex:
        ok = False
        print("    [FAIL]  a lawful aider node is refused: %s" % ex)
    try:
        lawful = lawful_tool()
        if topo(lawful["nodes"], lawful["edges"]) != ["a", "t"]:
            ok = False
            print("    [FAIL]  a lawful tool node is ordered wrongly")
    except ValueError as ex:
        ok = False
        print("    [FAIL]  a lawful tool node is refused: %s" % ex)
    for b in want["budget_cases"]:
        if over(b["elapsed_ms"], b["budget_s"]) != b["over"]:
            ok = False
            print("    [FAIL]  budget mispinned: %r" % b)
    ex = want["receipt_example"]
    if receipt(ex["run"], ex["node"], ex["output"], ex["ts"]) != ex["receipt"]:
        ok = False
        print("    [FAIL]  receipt formula does not reproduce")
    if ok:
        print()
        print("  PROVEN. The builder routes honestly.")
        return 0
    return 1


if __name__ == "__main__":
    # A WORD THIS DOES NOT KNOW WRITES NOTHING (2026-09-29). Anything that was
    # not `--verify` used to CUT, so `--check` -- prove.py's own word, typed by
    # a hand that meant "verify" -- rewrote the fixture from this file's
    # vectors and dropped the five the fixture had gained by hand.
    # The same refusal is every cutter's since that evening: tools/cut_words.py.
    from cut_words import word
    sys.exit(verify() if word(sys.argv[1:], __file__) == "--verify" else cut())
