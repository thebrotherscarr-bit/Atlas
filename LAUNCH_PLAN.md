# LAUNCH_PLAN — the harness, built on the recorder

*Rewritten 2026-10-16 on the operator's word: "Rewrite the system design and launch_plan. Supersede spec 8.2, on my word";
and of what it is for: "This is the plan, we are doing the harness. we built the recorder, its perfect, now we build the
harness on top of it." The plan it replaces (2026-09-08) stands in git at atlas `5648290`. What the harness is: the core's
SYSTEM_DESIGN, in his notes (tag `system-design`; the file retired 2026-10-23). What is open, item by item: the core's WHATS_LEFT.md -- H31 is the item, and its Now line says where we
are working today.*

---

## 0. Where it stands, 2026-10-16 — read off the disk

| thing | state | proof |
|---|---|---|
| the recorder | whole and green: strokes 3942 of 3942 and smoke 72 of 72 (2026-10-14, 19:00 and 18:54) | the core's `tests/last_run.json`; core commit `49dd3e6` |
| the repositories | core main at `7f69ab1`, atlas main at `5648290`; main protected in both since 2026-10-09 (a line of work, a pull request, green checks, his merge) | git log; GitHub |
| the stack | the door on :8090 (build `b4b5d8ba`), the glass on :8091 (`fe726e95`), Ollama on :11434 | the running processes |
| the last sitting | 382, closed 2026-10-14 19:06 | the core's `sessions/sessions.jsonl` |
| the rack | 11 models on one 16 GB RX 6800 XT; the driver's model, `qwen3.5:9b`, is 6.6 GB | the core's `rack.md` |
| the loop's law | LAW_003 THE LOOP (his word, 2026-09-17; sealed, link 5 of the chain): a loop is lawful with a declared ceiling, a stop a machine checks and every pass in the record, and the gate stands at the end | the core's `law/chain.jsonl` |
| the driver | not built; H31 is written and its step (2) planned on his answers of 2026-10-16 | WHATS_LEFT H31 |
| in its way | C71: Steward's call to the door gives up at 295 s while a strokes run takes about 313 s; and the door's `aider_run` answers his glass alone | WHATS_LEFT C71; H31 step (2)(e) |
| the documents | behind what was built since early October; the documentation pass (I4) is pulled forward to now | the pass's findings, under F |

---

## 1. LAUNCHED, this time

One sentence, provable: **he types a job and its check in the Agent tab, and Steward carries it to a green check through
the engine's skills and the door's tools -- every write on his card, every pass in the record -- and the report is written
from the record.**

Proved by three real jobs, live on his glass (H31 step (3)):

1. **C71** -- a strokes run asked through Steward comes back with its tally. *(R1: the set-up mends C71 first.)*
2. **A save and send of a finished line** -- its check: the line's pull request stands on GitHub with its checks green; the
   merge stays his press.
3. **HANDOFF's finding (bb)** -- the boot's gate finds the day's HANDOFF block; its check: a stroke shown red first, green
   after.

Then v0.2.6 is cut by his release flow (I2), and SPEC's first line gains what H31 says it never had: it finishes jobs.

---

## 2. Week one, 2026-10-16 to 10-22: the documents, the set-up, the loop

Each piece on his word, planned on What's left's Now line before it runs; each step in the glass's tabs with every write on
his card (RULE 14); each proof on his glass (RULE 13).

| # | piece | its check | waits on |
|---|---|---|---|
| 1 | **The documentation pass** (I4, pulled forward): SYSTEM_DESIGN.md, this file and SPEC 8.2 first; then every root document of both repositories read whole, its findings written under F before any mend, each mend on his card | every count and claim a document makes matches the disk; What's left's page stroke green | -- |
| 2 | **The set-up: C71** (H31 (2)(a)): a strokes run asked through the door comes back with its tally, because the suites are the check for code | the strokes asked through `suite_run` in the Agent tab return their tally inside the turn | R1, R2 |
| 3 | **The door change** (H31 (2)(e)): `aider_run` callable by the driver, every run parked for his card | a door test shown red first: the driver's call parks; his glass's call runs as it does today | his word, given: "the door change, i siad it, we do it" |
| 4 | **The loop**, in the engine (H31 (2)(b) to (d) and (f)): a job is a task that names its check; the tools are read, write, measure and locate, with Aider for code; LAW_003's three bounds; the report from the record | strokes, each shown red first: a job with no check refused by name; a job with no ceiling refused; the stop reads the check's own evidence, never the seat's words; every pass in the transcript, the failed ones included | R3, R4 |

A change to `manjuel/` is delivered with "restart required"; a change to the door or the glass is a built binary he places
on his card. Each piece names what would go red if it came unplugged (RULE 11).

---

## 3. Week two, 2026-10-23 to 10-29: the three jobs, then the cut

| # | piece | its check |
|---|---|---|
| 5 | **The three jobs**, live on his glass, one at a time (H31 step (3)) | each job's own check, as section 1 names it |
| 6 | **The numbers**: each job's passes, minutes, tool calls, cards and model loads, read off the record | written on H31 and in the CHANGELOG |
| 7 | **Another model, only if the numbers ask** (R5): Nemotron 3.5 Bolt on the job the 9B failed | the same job's check, in the same record |
| 8 | **Folding on the numbers**: D18, H23 and H25 close with H31; any other layer the driver's runs show is not needed (the routing, the relay of seats) is put to him, then folded, never deleted (LAW 1) | the suites green after each fold |
| 9 | **The cut**: v0.2.6, by his release flow | the release gate; GitHub's checks |

If week one runs long, week two starts later; the order does not change.

---

## 4. Rulings owed -- each put to him on a card when its piece comes up

| # | ruling | gates |
|---|---|---|
| R1 | **The first proof job.** H31 makes C71 both the set-up and the first of the three jobs; once the set-up mends it, it cannot also prove the driver. Either the driver takes a different first job (the hand's suggestion: C73, a real fault in the engine's check of numbers, which a stroke can hold), or C71 is left for the driver and, until then, its check for code is the strokes' own stamp (`tests/last_run.json`) | piece 2 |
| R2 | **C71's lever**: the skill's bound (`MANJUEL_SKILL_TIMEOUT`, 300 s) raised past the strokes' ~313 s, or a strokes run that answers at once and is read off its stamp | piece 2 |
| R3 | **How a job is typed** in the Agent tab: the words for the goal and for its check | piece 4 |
| R4 | **The first ceiling**: passes, minutes and tool calls per job; generous, tuned once it works | piece 4 |
| R5 | **Bolt, only on a measured failure.** Its official build (24 GB) is made for llama.cpp and is larger than the card; no Ollama build was found on 2026-10-16; a 4 GB build is reported for 2026-10-23. A second runtime beside Ollama is his ruling | piece 7 |
| R6 | **The cut**: his release flow and its three gates | piece 9 |

---

## 5. Risks, and the move the day each one shows

| risk | signal | move |
|---|---|---|
| the 9B loses the thread on a long chain | the record shows the same call again and again, a call with made-up arguments, or the ceiling reached | read it off the three jobs first; another model only on those numbers (SITTING LAW 3, R5) |
| cards pile up | a job waits on him longer than it works | count the cards per job in the record and put the count to him; LAW_003 asks for the gate only at the landing, so whether each edit on a line of work keeps its own card is his ruling, never a standing yes (RULE 2) |
| the driver's model and Aider's coder share one card | model loads take a job's minutes | measure first; one model at a time, `rack_unload` between |
| a check weaker than its job | the check passes and the thing asked about is still there | each check names the thing itself and is shown red before the job runs (Nigel's rule: a change is kept only when the check passes and the thing asked about is gone) |
| a long tool outruns its bound | another C71 | each bound read against the tool's measured time when the job is written (LAW 7) |
| the hand keeps going | a piece lands with a second one attached | RULE 10: one piece, then stop |
| the documents go stale again | a count or a claim no longer matches the disk | each piece writes its document lines in the same pass |

---

## 6. The road after the launch

Each on his word, in What's left's order:

- **H30**, the honesty gate on every answer: the report and the record never part.
- **H20**, a context cache: the window a long chain needs (with C70's step (2)).
- **H2 to H4**, memory that comes back, notices and consolidates: the driver starts from what was decided before.
- **I3 to I6**: v0.2.7 (both repositories keep what earlier versions promised), v0.2.8 (what is left of the documentation
  pass), v0.2.9 (packaging and the deliverable), v0.3.0 (a full release, packaged and ready to deploy).
- **The cage**: Microsoft's MXC or NVIDIA's OpenShell under the door's shell and Aider -- his call.
- **H8**, a seat on a hosted route, only on a measured need; **H10**, a job that runs without a sitting he opened -- later,
  and his.

---

## 7. The daily rhythm

**Morning**: the What's left tab first -- its Now line says where we are working today, and the day's plan is written there,
on his card, before a step runs. **Each step**: one line before it, in the glass's tabs, every write on his card, the result
checked where he can see it. **Each piece**: its CHANGELOG entry and the document lines it changed; "restart required" when
`manjuel/` moved; what would go red if it came unplugged. **Each save**: his word in chat, for that save -- a line of work, a
pull request, green checks, and the merge his press on GitHub. **Never**: `git status` or `git diff` from the hand's own
shell; an edit while a sitting is open; a folder unasked; the hand's name on his work; a send without his word; the
engine's `do the git cycle` on a line that has an upstream (it would send with no card).

---

## 8. Read first

`CLAUDE.md`'s list, every time (RULE 0); then What's left's Now line; then the core's SYSTEM_DESIGN (his notes, tag `system-design`) and this file.
