# ATLAS — Sovereign Agent Harness

A polyglot agent provenance system with cryptographic chain integrity, structural governance, and a modern observability webapp.

**Version:** in `VERSION` (v0.2.5 at the last cut, 2026-10-10)  
**Covenant:** 1512741580b7239b  
**License:** MIT

## What This Is

ATLAS is an operator-first agent harness that treats every tool call, trace, and agent declaration as a cryptographically chained, append-only record. It is not a framework — it is a governance layer that sits on top of any agent runtime and refuses to let the system lie about what happened.

One exception, since 2026-09-16 and on the operator's ruling: the Dashboard's own refresh reads (`muster`, `rack_list`, `proofs`, asked every 15 seconds while the page is on screen) are answered and not kept. Since 2026-09-21 the Projects card's two reads are answered the same way: its list (`projects`, asked when the page opens and after every turn) and the page its frame shows, served at `/api/projects/<name>/page`. Every other tool call through the glass is kept as a trace.

### Core Properties

| Property | Enforcement |
|---|---|
| `can_approve` is always `false` | Structural grammar — approval lives in the operator's hand alone |
| Append-only | Fold, never delete — every change leaves a witness |
| Chain integrity | SHA-256 hashing links traces to their inputs/outputs |
| Forbidden verbs | No agent approves, and none commits, pushes or merges alone: the door's git tools run from the operator's glass or wait on his card (a merge on GitHub allowed on his word, 2026-10-09) |
| Multi-tenancy | MCP carries ALL projects; every call names its ground |
| Zero external deps | No crates.io, no npm, no pip, no Go modules beyond stdlib |

## Architecture

```
atlas/
├── core/           Rust crate: store, prov-hash, ground-prove, atlas-api
├── store/          Rust crate: sqlite persistence
├── apps/atlas/     Rust CLI binary
├── line/           Go workspace: atlas-mcp, atlas-tui, atlas-town, atlas-door, atlas-vc
├── webapp/         Go + vanilla JS: the glass, the front page (port 8091)
├── tools/          Python: verifier scripts (stdlib-only)
├── specs/          Specification documents
├── agents/         .us agent declarations (40 seats)
├── skills/         Prove, orient, enroll skills
├── docs/           Full documentation
└── tests/          Integration tests
```

## Binaries

| Binary | Language | Description |
|---|---|---|
| `atlas` | Rust | CLI: version, ground-init, ground-prove |
| `atlas-mcp` | Go | The door: MCP and HTTP on :8090; its tools in `docs/TOOL_PERMISSIONS.md` |
| `atlas-tui` | Go | Terminal UI: dashboard, traces, agents, chain view |
| `atlas-town` | Go | Town square: multi-agent coordination |
| `atlas-door` | Go | Door: per-module MCP gateway |
| `atlas-vc` | Go | `.us` declarations as W3C verifiable credentials (`specs/SPEC_VC.md`) |
| `atlas-webapp` | Go | The glass: the front page on :8091 -- one terminal, its Inspector, every card |

## Quick Start

### Prerequisites

- Rust stable (`~\.cargo\bin`) **plus the MSVC toolchain** — Visual Studio
  Build Tools with "Desktop development with C++". `store/src/ffi.rs` links
  Windows' own `winsqlite3`, so `cargo build` needs `link.exe`. Without it
  the build dies on `linker 'link.exe' not found`, which says nothing about
  this project. Several GB of prerequisite, named nowhere until 2026-09-11.
- Go 1.26+ — both `line/go.mod` and `webapp/go.mod` say `go 1.26`.
- Python 3.10+ for the provers. This line said **3.14**, and nothing here
  needs it: the core's `pyproject.toml` says `>=3.10` and CI proves 3.10 and
  3.13. The `.venv` it also named is created by nothing and is gitignored —
  `python tests/prove.py` uses `sys.executable` and is the portable path.
- Windows 10 1803 or newer, for that same `winsqlite3` link.

### Build

```powershell
# Rust
$env:Path = "$env:USERPROFILE\.cargo\bin;" + $env:Path
cargo build --workspace

# Go
cd line ; go build -o atlas-mcp.exe .\cmd\atlas-mcp ; cd ..
# NOT `go build ./...` — line/cmd holds FIVE main packages, and Go
# discards the results when it compiles more than one. That command is a
# compile check that leaves you with no binary and no error.

# Webapp
cd webapp ; go build -o atlas-webapp.exe . ; cd ..
```

### Run

```powershell
# MCP server
cd line ; go run ./cmd/atlas-mcp

# TUI
cd line ; go run ./cmd/atlas-tui

# Webapp
cd webapp ; go run .
# Open http://localhost:8091
```

### Prove

```powershell
python tests\prove.py            # THE BALL: every leg atlas has, one verdict (AGENTS.md)
python tests\prove.py --live     # + the legs needing :8090 and Ollama
```

The legs still run on their own (`cargo test --workspace`, `go test ./...` in
`line` and `webapp`, `go run ./cmd/atlas-mcp --prove` in `line`, and any cutter
in `tools\` with `--verify`); `tests/PROVING.md` is the map.

## The glass

- **The front page** — one terminal with five tabs (Agent, Bash, Python,
  PowerShell, Aider) and an Inspector beside it: Run, What's left, Aider Pair,
  GitHub, Audit Ledger, Workflows, Registry & Docs, Laws, Settings
- **Cards** — every write that waits for the operator is a card pinned at the
  foot of the terminal, in every tab
- **Live** — each seat's tokens and its thought as they arrive (SSE)
- **Locked** — a PIN, this computer only (127.0.0.1)

The Trace Log, Messages and the Playground still draw at their own addresses;
nothing on the front page reaches them (the core's WHAT'S LEFT D19).

## Agent System

40 `.us` agent declaration files, validated by `atlas agent enroll --dry`:

- **manjuel** — Operator's serf, record keeper
- **Twelve Tribes** — Core household (aria, rex, mason, luna, quinn, alex, scout, courier, scribe, steward, neiro, builder)
- **Council personas** — Advisory seats
- **Lineage actors** — Historical chain
- **Gate** — Access control
- **Doctrine seats** — Law and foundation

Every agent carries:
- `can_approve: false` (structural, not configurable)
- `reports_to` chain (resolvable to root)
- Covenant hash `1512741580b7239b`
- Permissions declaration

## Specifications

| Spec | Location | Status |
|---|---|---|
| SPEC_US.md | `specs/SPEC_US.md` | Current |
| SPEC_COMMANDS.md | `specs/SPEC_COMMANDS.md` | Current |
| LAUNCH_PLAN.md | `LAUNCH_PLAN.md` | The road (2026-10-16) |
| THE_ROAD.md | `THE_ROAD.md` | The stone road, finished |

`specs/SPEC.md`, `SPEC_STORE.md` and `SPEC_TUI.md`, named here until 2026-10-16,
are not in `specs/`; `docs/README.md` lists what is.

## The Laws

1. Source grounds are READ-ONLY. Fold, never delete.
2. Goldens BEFORE assertions.
3. `can_approve` is false in every declaration ever.
4. Ports never move at cutover.
5. Proves are hermetic: temp grounds, never the live record.
6. Append-only witnesses.
7. VERSION pins the stone.
8. Multi-tenancy: strangers refused by name.

## Documentation

See `docs/README.md` for the full index.

## Contributing

See `CONTRIBUTING.md`.

## License

MIT — see `LICENSE`.
