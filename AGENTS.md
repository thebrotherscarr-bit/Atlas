# AGENTS — atlas sessions

**Preflight:** Read the core's `CLAUDE.md` and `law/` (the rules every hand
keeps) → its `WHATS_LEFT.md`, the Now line first (what is open, and where we
are working today) → `LAUNCH_PLAN.md` (the road) → work ONE piece → its
CHANGELOG entry. (Until 2026-10-16 this read THE_ROAD and SEAT_LOG's tail;
SEAT_LOG stopped on 2026-09-09.)

**2026-11-05: the Estate's design and doc standard.** Aurora is parked until Steward has visible computer control
(Steward's ADR-0004): no new Aurora work until then, and its look is never changed. The Estate's design now lives in
Steward's repository, and every repository follows its doc standard (the operator, 2026-11-05: "make this the
standard moving forward, this way nothing gets misaligned"):

- [docs/README.md](https://github.com/thebrotherscarr-bit/Steward/blob/main/docs/README.md) - the doc map: which
  document is for what
- [docs/ARCHITECTURE.md](https://github.com/thebrotherscarr-bit/Steward/blob/main/docs/ARCHITECTURE.md) - the Estate
  as it is now: parts, wires (ports, protocols), data, flows, rules, plan
- [docs/decisions/](https://github.com/thebrotherscarr-bit/Steward/tree/main/docs/decisions) - decision records in
  MAF's ADR format. A decision (a part, a wire, a port, a store of data, a rule) gets an ADR there before it is
  built; a feature gets a spec in Steward's `docs/specs/`; a multi-step build an ExecPlan in Steward's `docs/plans/`.

Where the Preflight above or LAUNCH_PLAN.md disagrees with those, they win. Every command below carries a hard time
limit when run by an agent (for example `timeout 1800 python tests/prove.py`). The pull-request checks are
`prove.yml`'s two jobs: the battery on Windows (`cargo build -p atlas --locked`, `python tests/prove.py`) and the Go
half (`go build ./...` and `go vet ./...` in `line/` and `webapp/`, then the suites).

---

## Ground truth

| Toolchain | Where | Notes |
|---|---|---|
| Rust stable | `~\.cargo\bin` (**off PATH** — prepend it) | spine: `core\`, `store\`, `apps\atlas` |
| Go 1.26+ | PATH | THE LINE: `line\` (atlas-mcp) |
| Python 3.10+ | `python` (no venv: `tests/prove.py` uses `sys.executable`) | stdlib-only tools in `tools\` |

**Zero external dependencies is standing law**: no crates.io, no npm, no
pip packages, no Go modules beyond stdlib. Hand-roll or refuse.

## Prove before LANDED — ONE command

```powershell
$env:Path = "$env:USERPROFILE\.cargo\bin;" + $env:Path
cargo build -p atlas                    # the door battery shells this binary
python tests\prove.py                   # THE BALL: all nine legs, one verdict
python tests\prove.py --live            # + the legs needing :8090 and Ollama
```

THE BALL gathers every prover atlas has — the Rust spine, both Go modules,
the two shipped batteries, **twenty-eight** golden verifiers (this file used
to name four), the six workflows, the E2E suite, and the GitHub workflows
that run it (GITHUB, 2026-10-09). Read `tests/PROVING.md` for the map.

It answers in **three** verdicts, not two. `ABSENT` means a leg named a
dependency this ground does not hold — the read-only source grounds
(`estate\`, `secondbrain\`) fourteen cutters were cut from, a binary not
built, a door not answering. ABSENT is never a pass and never a silent skip:
it prints the exact path that would answer it. Only `FAIL` exits red.

The legs also still run on their own:

```powershell
cargo test --workspace                                # each prints its own tally
cd line ; go build ./... ; go test ./... ; cd ..
cd line ; go run ./cmd/atlas-mcp --prove ; cd ..      # its battery, exit 0
python tools\<any_cutter>.py --verify
cargo run -q -p atlas -- --version                    # the number in VERSION
```

A piece is DONE when its check passes and its CHANGELOG entry is written.
Acceptance lives in the core's SPEC, section 4, and in the check
`LAUNCH_PLAN.md` names for each piece (ACCEPTANCE.md closed with the stones
on 2026-09-08; STATE_OF_BUILD.md stops at 2026-09-09).

## The laws (short form; CHARTER is authoritative)

1. Source grounds (`..\estate`, `..\secondbrain`) are READ-ONLY. Fold,
   never delete. Nothing invented — sources named.
2. Goldens BEFORE assertions: cut vectors from the read-only oracle first.
3. `can_approve` is false in every declaration ever. Approval lives in the
   operator's hand alone ("i am the approval").
4. Ports never move at cutover (SPEC_COMMANDS). Old verbs keep their verbs.
5. Proves are hermetic: temp grounds, never the live record.
6. Append-only witnesses: dated entries in CHANGELOG.md (SEAT_LOG.md and
   STATE_OF_BUILD.md, the stones' witnesses, stop at 2026-09-09).
7. VERSION pins the version (plain semver; the stone tag was struck
   2026-09-10); the version test catches drift.
8. Multi-tenancy: the MCP carries ALL projects; every call names its ground;
   strangers are refused by name.

## The household

40 enrolled seats (`agents\*.us`, folded via `tools\fold_agents.py`) —
Twelve Tribes under manjuel, council personas, lineage actors, gate,
seven doctrine seats. manjuel holds the charge as the operator's serf;
the operator is the approval. Undeclared actors get refused context by
name: enroll first (`atlas agent enroll <db> --dir agents`).

## Orientation door

```powershell
atlas orient --home .                                  # this repo
atlas orient --home ..\secondbrain\SecondBrain-collab\demo_vault\Manjuel
```
