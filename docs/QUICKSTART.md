# QUICKSTART — Atlas in one hour

*Rewritten 2026-10-16 (the core's documentation pass). The version it replaced, of
2026-09-07, cloned an address that was not this repository and started nothing; it
stands in git.*

Atlas is the control plane of the core (Manjuel): the door (`atlas-mcp`, MCP and HTTP
on 127.0.0.1:8090) and the glass (`atlas-webapp`, the front page on 127.0.0.1:8091).
It does nothing on its own: it raises the core's engine for a world, and every write
is the operator's -- his click on the glass, or his card. Start with the core.

## You need

| Tool | Version | For |
|---|---|---|
| Windows | 10 1803 or newer | the store links Windows' own `winsqlite3` |
| Go | 1.26+ | the door and the glass |
| Rust stable and the MSVC build tools | -- | `atlas.exe`, which atlas-door and the door's battery shell |
| Python | 3.10+ | the core's engine and atlas's provers |
| Ollama | on 127.0.0.1:11434 | the models the core's seats name |

## 1. Clone both, the core first

    git clone https://github.com/thebrotherscarr-bit/Manjuel.git <YOUR-GROUND>
    cd <YOUR-GROUND>
    git clone https://github.com/thebrotherscarr-bit/Atlas.git atlas

Atlas lives INSIDE the core's folder, at exactly `atlas`; the core's `.gitignore`
expects it there.

## 2. Build

    cd atlas ; cargo build -p atlas ; cd ..
    cd atlas\line ; go build -o atlas-mcp.exe .\cmd\atlas-mcp ; cd ..\..
    cd atlas\webapp ; go build -o atlas-webapp.exe . ; cd ..\..

Not `go build ./...`: `line/cmd` holds five main packages, and Go discards the
binaries when it compiles more than one.

## 3. Start, and open

The door and the glass each need the service wire, `ATLAS_SERVICE`, from the core's
`.env` in their environment -- never on a command line. The exact lines, and what
each flag does, are in the core's `RUNBOOK.md`, "Starting the system". Then open
`http://127.0.0.1:8091`: the first time, it asks for your name and a PIN.

## 4. Prove

    cd atlas ; python tests\prove.py ; cd ..

THE BALL: every leg atlas has, one verdict -- PASS, FAIL, or ABSENT (a leg whose
dependency this machine does not hold, never counted as a pass). The core's own
suites are in its `TESTING.md`.

## What's next

- the core's `RUNBOOK.md` -- starting, running and stopping the whole system
- [../LAUNCH_PLAN.md](../LAUNCH_PLAN.md) -- where the work is going
- [ARCHITECTURE.md](ARCHITECTURE.md) -- how atlas is built
- [TOOL_PERMISSIONS.md](TOOL_PERMISSIONS.md) -- every tool the door serves, and who may call it
