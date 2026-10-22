#!/usr/bin/env pwsh
# update.ps1 -- Atlas updates itself from GitHub's main, as SecondBrain does.
#
# His word, 2026-10-22: "atlas should auto-update just like the secondbrain
# does." The glass's Update button runs this (handlers/update.go); it can also
# be run by hand. In order:
#
#   1. refuse while a sitting is open (RULE 9: nothing moves under his sitting),
#      while this checkout is not on main, or while it holds unsaved work;
#   2. take main from GitHub, fast-forward only;
#   3. build the door (line/cmd/atlas-mcp) and the glass (webapp) beside the
#      running ones, as *.new.exe;
#   4. stop the door (8090) and the glass (8091), keep each old program as
#      *.old-<stamp>.exe, put the new ones in place (atlas-desktop.exe is the
#      glass's own program under its desktop name, so it is renewed too);
#   5. start the door with RUNBOOK's line and the glass, each with the service
#      wire read from the ground's .env (never printed, never on a command line).
#
# Anything that fails before step 4 leaves both running as they were.
# Pure ASCII, so PowerShell 5.1 reads it whatever the code page.

param(
    [string]$Log = ""
)

$ErrorActionPreference = 'Stop'
$atlas = $PSScriptRoot
$ground = Split-Path $atlas -Parent
if (-not $Log) { $Log = Join-Path $atlas 'update.log' }

function Say([string]$t) {
    $line = "{0}  {1}" -f (Get-Date -Format 'yyyy-MM-dd HH:mm:ss'), $t
    Add-Content -Path $Log -Value $line
    Write-Output $line
}

# Named Invoke-AtlasGit and calling git.exe: PowerShell's names ignore case, so a helper called
# Git that runs `git` calls itself until the call depth overflows (the first two presses, 2026-10-22).
function Invoke-AtlasGit([string[]]$a) {
    $env:GIT_TERMINAL_PROMPT = '0'
    $out = & git.exe -C $atlas @a 2>&1
    if ($LASTEXITCODE -ne 0) { throw "git $($a -join ' ') failed: $out" }
    return ($out | Out-String).Trim()
}

function PidOn([int]$port) {
    $c = Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($c) { return $c.OwningProcess }
    return $null
}

function Wire {
    $f = Join-Path $ground '.env'
    if (-not (Test-Path $f)) { return '' }
    $line = Get-Content $f | Where-Object { $_ -match '^\s*ATLAS_SERVICE\s*=' } | Select-Object -First 1
    if (-not $line) { return '' }
    return ($line -split '=', 2)[1].Trim().Trim('"', "'")
}

try {
    Say "update: starting in $atlas"

    # 1. the refusals
    $ledger = Join-Path $ground 'sessions\sessions.jsonl'
    if (Test-Path $ledger) {
        $last = Get-Content $ledger -Tail 1
        if ($last) {
            $row = $last | ConvertFrom-Json
            if (-not $row.ended) { throw "sitting $($row.n) is open; close it first (RULE 9: nothing moves under a sitting)" }
        }
    }
    $branch = Invoke-AtlasGit @('rev-parse', '--abbrev-ref', 'HEAD')
    if ($branch -ne 'main') { throw "this checkout is on '$branch', not main; the update only moves main" }
    & git.exe -C $atlas diff-index --quiet HEAD -- 2>$null
    if ($LASTEXITCODE -ne 0) { throw "this checkout has unsaved work; save or put it aside first" }

    # 2. main from GitHub, forward only
    Invoke-AtlasGit @('fetch', '-q', 'origin', 'main') | Out-Null
    $before = Invoke-AtlasGit @('rev-parse', '--short', 'HEAD')
    Invoke-AtlasGit @('merge', '-q', '--ff-only', 'origin/main') | Out-Null
    $after = Invoke-AtlasGit @('rev-parse', '--short', 'HEAD')
    Say "update: main $before -> $after"

    # 3. build beside the running programs
    $line = Join-Path $atlas 'line'
    $web = Join-Path $atlas 'webapp'
    Push-Location $line
    & go build -o atlas-mcp.new.exe ./cmd/atlas-mcp
    $ok = $LASTEXITCODE -eq 0
    Pop-Location
    if (-not $ok) { throw 'the door did not build' }
    Push-Location $web
    & go build -o atlas-webapp.new.exe .
    $ok = $LASTEXITCODE -eq 0
    Pop-Location
    if (-not $ok) { throw 'the glass did not build' }
    Say 'update: built the door and the glass'

    # 4. stop, keep the old, place the new
    Start-Sleep -Seconds 2   # the glass answers the button before it is stopped
    foreach ($port in 8090, 8091) {
        $p = PidOn $port
        if ($p) { Stop-Process -Id $p -Force; Say "update: stopped pid $p on $port" }
    }
    Start-Sleep -Milliseconds 800
    $stamp = Get-Date -Format 'MMdd-HHmmss'
    $swaps = @(
        @{ Dir = $line; Name = 'atlas-mcp.exe'; New = 'atlas-mcp.new.exe' },
        @{ Dir = $web; Name = 'atlas-webapp.exe'; New = 'atlas-webapp.new.exe' }
    )
    foreach ($s in $swaps) {
        $cur = Join-Path $s.Dir $s.Name
        if (Test-Path $cur) { Rename-Item $cur ($s.Name -replace '\.exe$', ".old-$stamp.exe") }
        Move-Item (Join-Path $s.Dir $s.New) $cur
    }
    $desk = Join-Path $web 'atlas-desktop.exe'
    if (Test-Path $desk) { Rename-Item $desk "atlas-desktop.old-$stamp.exe" }
    Copy-Item (Join-Path $web 'atlas-webapp.exe') $desk
    Say "update: placed; the old programs are kept as *.old-$stamp.exe"

    # 5. start them again
    $wire = Wire
    if ($wire) { $env:ATLAS_SERVICE = $wire }
    $g = $ground -replace '\\', '/'
    $doorArgs = @('--http', '127.0.0.1:8090', '--auth', '--tenant', "research=$g", '--tenant', "atlas=$g/atlas",
        '--default-project', 'research', '--manjuel', "`"python $g/manjuel.py`"")
    Start-Process -FilePath (Join-Path $line 'atlas-mcp.exe') -ArgumentList $doorArgs -WorkingDirectory $line -WindowStyle Hidden `
        -RedirectStandardOutput (Join-Path $line 'mcp.update.log') -RedirectStandardError (Join-Path $line 'mcp.update.err.log')
    Start-Process -FilePath (Join-Path $web 'atlas-webapp.exe') -WorkingDirectory $web -WindowStyle Hidden
    foreach ($port in 8090, 8091) {
        $up = $false
        foreach ($i in 1..40) { if (PidOn $port) { $up = $true; break }; Start-Sleep -Milliseconds 500 }
        Say "update: port $port answering: $up"
    }
    Say "update: done, Atlas is on $after"
}
catch {
    Say "update: REFUSED or FAILED -- $($_.Exception.Message)"
    foreach ($n in 'line\atlas-mcp.new.exe', 'webapp\atlas-webapp.new.exe') {
        $f = Join-Path $atlas $n
        if (Test-Path $f) { Remove-Item $f -Force }
    }
    exit 1
}
