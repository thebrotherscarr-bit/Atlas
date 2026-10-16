# ATLAS Documentation Index

*Every document in the Atlas documentation suite, organized by purpose.*

---

## Getting Started

| Document | Purpose |
|---|---|
| [QUICKSTART.md](QUICKSTART.md) | Get running in 5 minutes |
| [TUTORIAL.md](TUTORIAL.md) | 5 progressive exercises with answer key |
| [WALKTHROUGH.md](WALKTHROUGH.md) | End-to-end system tour |
| [GUIDED_DEMO.md](GUIDED_DEMO.md) | Live demo walkthrough |

## Architecture & Design

| Document | Purpose |
|---|---|
| [ARCHITECTURE.md](ARCHITECTURE.md) | System design, layers, trust model |
| [PARITY_PLAN.md](PARITY_PLAN.md) | Cross-language parity matrix |
| [PLAN_GUI_IDE.md](PLAN_GUI_IDE.md) | GUI/UX/IDE design plan |
| [BUILD_PLAN_TUI_AGENTS.md](BUILD_PLAN_TUI_AGENTS.md) | TUI test, agent docs, skills build plan |

## Reference

| Document | Purpose |
|---|---|
| [CLI_REFERENCE.md](CLI_REFERENCE.md) | Every command, every flag |
| [COMPLIANCE.md](COMPLIANCE.md) | EU AI Act, SOC 2, ISO 42001 mapping |
| [TOOL_PERMISSIONS.md](TOOL_PERMISSIONS.md) | Every tool the door offers: what it declares, who may call it (generated, held current by a test) |

## Also in this folder

| Document | Purpose |
|---|---|
| [SPEC_CONTROL_CENTER.md](SPEC_CONTROL_CENTER.md) | The glass's spec (its PRD) of 2026-09-08 |
| [WORKFLOWS.md](WORKFLOWS.md) | Six agent workflows (stamped 0.1.8) |
| [PIPELINES.md](PIPELINES.md) | Six CI/CD pipelines, a design (stamped 0.1.8; CI runs `prove.yml`) |
| [OLLAMA_PROVER.md](OLLAMA_PROVER.md) | The Ollama prover (stamped 0.1.8) |
| [ACCEPTANCE.md](ACCEPTANCE.md) | Pipeline acceptance criteria (stamped 0.1.8) |
| [OPENAPI.yaml](OPENAPI.yaml) | The door's HTTP API (stamped 0.1.0+f1) |
| [GUI_GAP_ANALYSIS.md](GUI_GAP_ANALYSIS.md) | The GUI against the market: what it has, what is expected, what is missing |

## Specifications

| Document | Purpose |
|---|---|
| [specs/SPEC_US.md](../specs/SPEC_US.md) | .us declaration format |
| [specs/SPEC_SEAM.md](../specs/SPEC_SEAM.md) | Cross-language seam contract |
| [specs/SPEC_CANON.md](../specs/SPEC_CANON.md) | Canonical JSON forms |
| [specs/SPEC_CHAINS.md](../specs/SPEC_CHAINS.md) | Hash chain format & verdicts |
| [specs/SPEC_SQLITE.md](../specs/SPEC_SQLITE.md) | SQLite store design |
| [specs/SPEC_COMMANDS.md](../specs/SPEC_COMMANDS.md) | Old ↔ new command map |
| [specs/SPEC_COVENANT.md](../specs/SPEC_COVENANT.md) | Covenant derivation |
| [specs/SPEC_US_MESH.md](../specs/SPEC_US_MESH.md) | Mesh protocol spec |
| [specs/SPEC_VC.md](../specs/SPEC_VC.md) | W3C VC adapter spec |
| [specs/SPEC_US_SCHEMA.json](../specs/SPEC_US_SCHEMA.json) | JSON Schema for .us |
| [specs/SPEC_MCP.json](../specs/SPEC_MCP.json) | The door's MCP spec, as JSON |

## Agent Documentation

| Document | Purpose |
|---|---|
| [../agents/docs/README.md](../agents/docs/README.md) | Agent household index (40 seats) |
| [../agents/docs/operator.md](../agents/docs/operator.md) | The hand; holds the gate |
| [../agents/docs/manjuel.md](../agents/docs/manjuel.md) | The living House heart |
| [../agents/docs/steward.md](../agents/docs/steward.md) | Plans, specs, keeps THE_ROAD |
| [../agents/docs/archivist.md](../agents/docs/archivist.md) | Keeps the record findable |
| [../agents/docs/neiro.md](../agents/docs/neiro.md) | Proposes, aligns, files |
| [../agents/docs/jesster.md](../agents/docs/jesster.md) | Verifies, refutes, returns |
| + 34 more | See agents/docs/README.md |

## Skills

| Document | Purpose |
|---|---|
| [../skills/prove/SKILL.md](../skills/prove/SKILL.md) | Run the full prove battery |
| [../skills/orient/SKILL.md](../skills/orient/SKILL.md) | Assemble the orientation pack |
| [../skills/enroll/SKILL.md](../skills/enroll/SKILL.md) | Enroll .us declarations into master.db |

## Architecture Decision Records

| Document | Purpose |
|---|---|
| [ADR/001-polyglot-seam.md](ADR/001-polyglot-seam.md) | Why polyglot with exec seams |
| [ADR/002-fold-never-delete.md](ADR/002-fold-never-delete.md) | Why append-only |
| [ADR/003-can-approve-false.md](ADR/003-can-approve-false.md) | Why structural refusal |
| [ADR/004-socket-free-kernels.md](ADR/004-socket-free-kernels.md) | Why zero network in C++ |
| [ADR/005-golden-master-parity.md](ADR/005-golden-master-parity.md) | Why golden-master testing |
| [ADR/006-the-door-is-a-product.md](ADR/006-the-door-is-a-product.md) | Why the door is a product: `atlas-mcp` as a standalone MCP server |

## Rollback Notes

| Document | Purpose |
|---|---|
| [rollback/A1_SPINE.md](rollback/A1_SPINE.md) | A1 rollback procedure |
| [rollback/B1_LINE.md](rollback/B1_LINE.md) | B1 rollback procedure |
| [rollback/C1_FACES.md](rollback/C1_FACES.md) | C1 rollback procedure |
| [rollback/D1_TOWN.md](rollback/D1_TOWN.md) | D1 rollback procedure |
| [rollback/D2_DOOR.md](rollback/D2_DOOR.md) | D2 rollback procedure |
| [rollback/E1_KERNELS.md](rollback/E1_KERNELS.md) | E1 rollback procedure |
| [rollback/F1_HARVEST.md](rollback/F1_HARVEST.md) | F1 rollback procedure |
| [rollback/GX02_PORT_COMMAND_REPORT.md](rollback/GX02_PORT_COMMAND_REPORT.md) | Port command verification |

## Project Management

| Document | Purpose |
|---|---|
| [../LAUNCH_PLAN.md](../LAUNCH_PLAN.md) | The road: the harness built on the recorder, its order and its two weeks (2026-10-16) |
| [../STATE_OF_BUILD.md](../STATE_OF_BUILD.md) | Append-only progress log (the stones, to 2026-09-09) |
| [../SEAT_LOG.md](../SEAT_LOG.md) | Append-only activity log |
| [../ACCEPTANCE.md](../ACCEPTANCE.md) | Testable acceptance criteria (the stones, closed 2026-09-08) |
| [../THE_ROAD.md](../THE_ROAD.md) | The stone road, finished; it points at LAUNCH_PLAN |
| [../CHANGELOG.md](../CHANGELOG.md) | Version history |

## Root Documents

| Document | Purpose |
|---|---|
| [../README.md](../README.md) | Project overview |
| [../CHARTER.md](../CHARTER.md) | Founding rulings & standing laws |
| [../THE_CATALOG.md](../THE_CATALOG.md) | Module inventory |
| [../CONTRIBUTING.md](../CONTRIBUTING.md) | Contribution guide |
| [../SECURITY.md](../SECURITY.md) | Security policy |
| the core's `CLAUDE.md` | The operator's rules for any hand -- in the core repository, not this one |

## SDKs

None in this repository. The two links that stood here (`../sdk/python`, `../sdk/npm`)
named folders atlas does not hold (found 2026-10-16).

## Testing

| Document | Purpose |
|---|---|
| [../tests/e2e/SPEC_SHEET.md](../tests/e2e/SPEC_SHEET.md) | E2E test specification |

---

*Listed against the folder on 2026-10-16 (the core's documentation pass); the 2026-09-08 count it replaced is in git.*
