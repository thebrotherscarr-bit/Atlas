# Changelog

All notable changes to ATLAS will be documented in this file.

Format follows [Keep a Changelog](https://keepachangelog.com/).

Versions are plain semver from 0.1.2 on. Through 0.1.1 they carried a build
tag naming the stone that cut them — `0.1.0+a1` through `0.1.1+f1`. The
operator struck the moniker 2026-09-10: *"remove the moniker for the stones,
no letters in my versions."* Released headers below keep the tag they shipped
under, because they are the record of what happened.

## [Unreleased]

## [v0.2.5] — 2026-10-10 09:27 (tag on 144a425)

### A main on GitHub takes a pull request and nothing else: Land onto main and the Send of main refused by the door and greyed on the GitHub tab, where origin has the main line (2026-10-09; the core's WHAT'S LEFT I1, seventh step; his cards: "Mend Land onto main" and "Main is on GitHub")

Since the core's fifth step GitHub takes a change to main only through a pull request whose checks have passed, in both repositories and from him too, and the GitHub tab still offered Land onto main and then Send -- its own note said "Land onto main, then Send, still lands a line without one" -- a landing here that no send could carry, ending in GitHub's refusal. The marks' refusal offered the same road ("Send the main line first, then the mark"). On his card ("Main is on GitHub") the door decides it from what this machine already knows, with no network and whether or not the wall is open: origin has the world's main line.

`mainByPullRequest` (`gitctl.go`) is that judgement, once, in the door's words: "Refused: main is on GitHub, and GitHub takes a change to it only through a pull request whose checks have passed. Do the work on a line of work, send that line, and open a pull request from it (Pull requests, on the GitHub tab); once its checks pass, Merge on GitHub merges it there and brings the new main down here." `git_branch land` asks it before it lands, wherever he stands; `git_push` asks it when the line being sent is the main line; the lines list carries it as `why_not_land`, and the `git` state as `why_not_send` while he stands on main, the way the marks' list carries `why_not` (his word of 2026-09-17: "make the send button ask the door first"). A main origin has never had -- a world on this machine alone, or one never sent -- lands and sends as it always has, and a line of work is sent as before. The marks' refusal now says to bring the commit onto origin's main line first, "(a main origin has takes it only through a pull request)". On the GitHub tab, Land onto main is greyed with the door's sentence and the sentence is said once under the lines; Send to GitHub is greyed with it while he stands on such a main, and so is the council's Push on the page under Pages; and the Pull requests panel's note no longer says that Land onto main, then Send, lands a line. The head of `gitctl.go` says where a landing still goes.

Held by `TestAMainOriginHasTakesAPullRequestAndNothingElse` (new: before origin has main nothing is held back; then a landing refused by name and main unmoved, in the lines list's own words; a save made on main not sent and not on the remote, in the state's own words; and the line of work still sent, nothing held back on it), with the three mark strokes taught that the line reaches origin as a merged pull request brings it (a plain push in the stroke) and that the refusal names origin's main line; and on the glass by `TestAMainOnGitHubIsOfferedThePullRequestAndNothingElse` (new: Land, Send and the council's Push greyed with the door's two keys, the sentence said under the lines, the old road's words gone; and the door's land, `git_push`, lines list and `git` state each asking `mainByPullRequest`, which reads origin's main). Written first and run in his Bash tab: red for the right reasons -- the land went through ("Landed \"spur\" onto \"main\""), the marks still said "Send the main line first", and Land onto main was live -- with the stroke that now sends by a plain push green. Then the edits through `file_edit` on the line of work `land-by-pr` (opened from `provers-once`, whose pull request #8 waits on his press), each on his Approve; `node --check`, gofmt and vet clean, the door's tools package green (80.6 s), the glass's `server`, `handlers` and `db` green, and the battery in his PowerShell tab (`python atlas/tests/prove.py --check --quiet`, 97.9 s): 25 held, 14 absent, 0 broke.

Built in his Bash tab and placed on his card ("Place both, keep both"): the door `0c999f04...` on 8090 as pid 18656, started with the flags it ran with, and the glass `f6b98d3c...` on 8091 as pid 1840; the builds before them kept as `atlas-mcp.old-1009d.exe` (`eeeb07cd...`) and `atlas-webapp.old-1009b.exe` (`6705a091...`). Proved live: the placed door's lines list carries the sentence for both worlds and the `git` state none while he stands on a line of work, and the GitHub tab's Lines of work for atlas says it under the lines, with his session kept.

**What goes red if unplugged:** `TestAMainOriginHasTakesAPullRequestAndNothingElse`, if a landing onto or a send of a main origin has goes through, or the list or the state says other than the act; `TestAMainOnGitHubIsOfferedThePullRequestAndNothingElse`, if the page offers Land, Send or the council's Push there, greys them without the door's sentence, or the door's verbs stop asking the one judgement. **Not yet proved:** a press of Land onto main or Send on a main GitHub has, on the glass; he stands on a line of work, where neither is offered. **Named, not fixed:** the council's own push (the core's `git_cycle` and `git_push` skills) goes through the engine's `gitstate.py`, not this door, so on a main GitHub has it is refused by GitHub, not by name here; the page's Push is greyed for it.

### The provers themselves: release.yml's outside action named by its commit, and a line proved once, by its pull request (2026-10-09; the core's WHAT'S LEFT I1, sixth step; his word: "save and send, then go on to I1 step 6")

Two things about the workflows that prove atlas on GitHub. `release.yml`'s last step, the one that holds `contents: write` and publishes the draft, named `softprops/action-gh-release@v2`: a tag its owner moves, so a release ran whatever `v2` pointed at that day. And `prove.yml` ran on a bare `push:` beside `pull_request:`, so a line sent with its pull request open was proved twice, the same commit both times (four jobs for two), and a mark once more beside `release.yml`'s own battery.

Now `release.yml` names the commit `v2` pointed at when it was pinned, `3bb12739c298aeb8a4eeaf626c5b8d85266b0e65` (v2.6.2, read in his Bash tab with `gh api` on the action's tags), the tag beside it as a comment and a note saying it moves only by an edit there; GitHub's own `actions/checkout`, `setup-go` and `setup-python` stay on their major tags. And `prove.yml` pushes for `main` and the marks (`branches: [main]`, `tags: ["v*"]`) beside `pull_request:` and `workflow_dispatch:`: a line of work is proved once, by its pull request, which is the run main's ruleset waits for (`the battery (windows)`); main and a mark are proved by their own push, as before. The comments that said the battery runs "on every push" (in `prove.yml` twice, and in `release.yml`) say what runs it now.

The wire is a new leg of the battery, GITHUB (`tests/prove.py`, `leg_github`, after the goldens), one row per workflow: FAIL for an action from outside GitHub's own `actions/` named by anything but a 40-hex commit, and for a workflow that runs on pull requests and on a push to every line (a bare `push:`, `branches-ignore`, or `branches` of `*` alone; a push limited to the marks is not one); read as text by indentation, standard library only, and never ABSENT. The core's release gate holds the core's own workflows to the same two rules (the core's CHANGELOG, "core: a line is proved once"). `tests/PROVING.md` gives it the ninth row of the legs table, and `AGENTS.md` counts nine.

Proved on his glass, red first: the leg run alone in his Bash tab against the workflows as they stood, "FAIL prove.yml -- it runs on pull requests AND on a push to every line ..." and "FAIL release.yml -- softprops/action-gh-release@v2 names a tag its owner can move ..."; then, the edits made through `file_edit` on the line of work `provers-once` (opened from `protect-main`, whose pull request #7 waits on his press), each on his Approve, "PASS prove.yml -- no outside action; no line proved twice" and "PASS release.yml -- 1 outside action, each named by its commit; no line proved twice". And the whole battery in his PowerShell tab (`python atlas/tests/prove.py --check --quiet`, 190 s): 25 held, 14 absent, 0 broke. Nothing in the door or the glass moved: no build, no restart.

**What goes red if unplugged:** the GITHUB leg, on every pull request into main, where the battery is the check the ruleset waits for, if `release.yml`'s action is named by a tag again or `prove.yml` runs on every line's push beside its pull request. **Not yet proved:** one run per sent line on GitHub, and a release through the pinned action; the first is this line's own pull request, after his word to send it, the second the next mark (v0.2.5).

### The release flow's merge: `git_pr merge` by the line it opened, and waiting for green (2026-10-09; the core's WHAT'S LEFT I1, fifth step; his cards: "Everyone, me too" and "Inside its gates")

His word, 2026-10-09: "save and send, then go on to I1 step 5" -- main protected in both repositories, a pull request and green provers to change it, and the release flow taught to land its record through it. His cards the same hour: the rule binds everyone, him included ("Everyone, me too"), and the release flow merges inside its gates ("Inside its gates"). The rule itself is GitHub's (the core's CHANGELOG, "core: main takes a change only through a pull request"); the door's half is what a flow needs to merge what it opened.

`git_pr merge` takes `line` as well as `number`. His glass's button names the number; the release flow has the line it opened and no number to hand, so it names the line, gh is asked for the pull request from that line (`gh pr view <line>`), and the number gh answers is the one merged, exactly as the button's merge is. A name no line may carry is refused before gh is asked, a line with no pull request is said so, and a pull request gh answers for another line is refused. A line whose pull request GitHub has merged already is said so ("Pull request #N from \"<line>\" is merged on GitHub already") and main still comes down, so a release that stopped between its two merges can be fired again; a number already merged is still refused, as the button's press was. And `merge` takes `wait`, whole minutes from 0 to 30: a merge given minutes asks gh again every 30 s while the provers have not answered (none yet, or still running), merges the moment every check has passed, answers a failure at once, and is refused as before when the wait is spent; anything else as a wait is refused before gh is asked, and no wait is no waiting. Opening a pull request no longer tells him to land it by hand: "when every check on it has passed, Merge on GitHub (the GitHub tab, Pull requests) merges it there and brings main down here". The tool declares `line?` and `wait?` beside `number?`, and its description says both.

Held by `TestAFlowMergesThePullRequestFromItsLine` (new: a name no line may carry refused with gh never asked; a line with no pull request; a pull request answered for another line; then on a real line sent to an origin on this disk, merged by its line with gh asked to read `pr view spur` and to merge the number it answered, `pr merge 9 --merge --match-head-commit <head>`, and main down; and, GitHub having merged it already, said so with no second merge asked and main down) and `TestAMergeThatWaitsMergesOnceTheChecksPass` (new: five unlawful waits refused before gh; then none, then running, then passed merged after two waits; running then failed refused after one; running throughout refused once the wait is spent; no wait asks once). Each was run first in his Bash tab and was red for the right reason: the line's merge refused ("name the pull request to merge, by its number"), and the wait's stroke would not build (`undefined: prSleep`). And on the glass by `TestVersionControlListsAndOpensPullRequests`, which reads the door's `git_pr` arguments and now wants `line?` and `wait?` beside `number?`. Made on the line of work `protect-main` (opened from `merge-on-github`, whose pull request #6 waits on his press) through `file_edit`, each card on his Approve; gofmt and vet clean, the door's tools package green twice in his Bash tab (118.3 s, 115.4 s), and the glass's server green (0.8 s); nothing in the glass moved but its test.

Built in his Bash tab and placed on his cards, twice, the second over the first: `8537868f...` (the merge by line) as pid 26808, then `eeeb07cd...` (and the wait) as pid 7544 on 8090, started each time with the flags the door ran with; the builds before them kept as `atlas-mcp.old-1009b.exe` (`c45a5ef8...`, step 4's) and `atlas-mcp.old-1009c.exe` (`8537868f...`). Proved live: the door's tool list gives `git_pr` its `line` and its `wait` and the new description. The release flow that uses both is the core's (`flows/release.json` v5, folded on his card, "Fold it as v5").

**What goes red if unplugged:** the two new strokes, if a merge by line asks gh for anything but the pull request from that line and the number it answered, merges what another line opened, merges again what GitHub merged, or brings nothing down; or if a merge waits without a bound, waits on a failure, merges before every check has passed, or takes a wait it should refuse; and `TestVersionControlListsAndOpensPullRequests`, if the door stops declaring `line?` and `wait?`. At run time, the release flow's check after each merge (`is down on this machine`) stops the release if a merge did not land. **Not yet proved:** a merge by line, and a wait, on a real pull request; the first will be the next cut's (v0.2.5).

### Merge on GitHub: a pull request merged there as a merge commit, live only when every check on it has passed, and the new main brought down here (2026-10-09; the core's WHAT'S LEFT I1, fourth step; his card: "Merge on GitHub when green", and "Merge commit")

Until today a pull request opened from the GitHub tab landed nothing by itself: `git_pr` refused every action but `list` and `open`, saying "It never merges", and a line landed only as it always had, Land onto main (fast-forward) then Send, with GitHub marking the pull request merged once the main line held its saves. The door's genesis law kept `merge` out of its table ("absent BY CONSTRUCTION", `tools.go`), B16's `land` said the table "carries no `merge` by its founding law and never will" (`gitctl.go`), and the glass's strokes held the page to it. His word superseded that law for this: "we have to be able to merge, thats before the github functionality was complete, before hte remote repo was set up, and before we did proper CI and version controls. supersede that one". The rest of it stands: a forbidden verb is never free (`holds.go`, since 2026-09-25), and this one is his glass's alone.

`git_pr merge` (`number`, as JSON brings it or as typed digits) reads the pull request first (`gh pr view <n> --json ...`) and refuses by name one it is not told the number of, one gh cannot read or answers in a shape it cannot read, one that is not open, one that does not go into the main line, one with no checks, checks still running or checks that failed, one GitHub says conflicts, and one whose head gh did not say. Then it asks GitHub to merge exactly the head it read, as a merge commit (`gh pr merge <n> --merge --match-head-commit <sha>`): never `--admin`, never `--auto`, and the line is not deleted. Then it brings the new main line down, fast-forward only, as Take from GitHub does: taken where he stands on a clean main line (`pull --ff-only origin main`), moved where he stands on another line (`fetch origin main:main`), and fetched but left for him over unsaved work; it says which, and where main stands. The tool stays `ServiceOnly` and behind the same wall as every send, with `list` its one reading action; it declares `number?`, and any other action is refused naming the three.

On the GitHub tab each pull request's row carries Merge on GitHub, live only while that pull request's checks have all passed, and greyed otherwise with the reason (failed, still running, none yet). It is armed like a mark's Send: the first press says what will happen (a merge commit on GitHub, the new main brought down here, and that this page does not take it back), the second asks the door; a repaint disarms it, and the door's answer is kept under the list as an open's is. Land onto main, then Send, still lands a line without a pull request. Where the old law stood, the comments say it is superseded, dated, with his words: the genesis law in `tools.go`, B16's comment in `gitctl.go`, and the heads of `gitctl.go` and `gitpr.go`.

Held by `TestAPullRequestMergesOnlyWhenEveryCheckHasPassed` (new: a merge naming no number, an empty one, a word, zero, a negative and a fraction refused with gh never asked; then merged, closed, into another line, no checks, still running, failed, conflicting, no head and a shape it cannot read each refused by name, gh asked only to read; and gh failing said) and `TestAGreenPullRequestMergesOnGitHubAndTheNewMainComesDown` (new: on a real line sent to an origin on this disk, GitHub's own act stood in for by a merge commit made on that origin; gh asked to read and then to merge exactly `pr merge 9 --merge --match-head-commit <head>`, and nothing more; standing on the line, main here moves to GitHub's new main; on a clean main line it is taken and the merged work is on the disk; over unsaved work it is fetched and main here stays where it was), with `TestOpeningAPullRequestIsRefusedWhereItMustBe`'s merge now refused for naming no pull request. And on the glass by `TestVersionControlListsAndOpensPullRequests` (the button live only under `p.checks === 'passed'`, greyed with its reason, one merge call, from the panel's own `merge()`, armed by a first press; and the door's `case "merge":`, its `number?`, and `prMerge` judging the checks) and `TestVersionControlOffersLandingALineOntoMain` (the Lines-of-work box still asks the door for no merge). Written first and run in his Bash tab: red, the door's three for the right reason ("merge" is not something this does) and the glass's on the button. Then the door's and the page's edits through `file_edit` on the line of work `merge-on-github`, each on his Approve; `node --check`, gofmt and vet clean on both modules, the door's tools package green (73.5 s) and the glass's suite green (`db`, `handlers`, `server`).

Built in his Bash tab (door `c45a5ef8...`, glass `6705a091...`) and placed on his card ("Place both, keep both"): the door on 8090 as pid 32316 and the glass on 8091 as pid 32528, each started on the line it ran with, the builds before them kept as `atlas-mcp.old-1009.exe` (`a95a34db...`) and `atlas-webapp.old-1009.exe` (`e76db7e4...`). Proved live: the glass healthy (0.2.4) with his session kept; the GitHub tab's Pull requests for atlas drew the new words, with no pull request open; and the placed door's own tool list gives `git_pr` its `number` and its merge.

**What goes red if unplugged:** the two new strokes and the opening one, if a merge is made unnamed, on a pull request that is not open, not into the main line, not green, conflicting or of a head not read, if gh is asked for anything but that one merge commit, or if the new main does not come down, or is taken over unsaved work; `TestVersionControlListsAndOpensPullRequests`, if the button is live on a pull request whose checks have not all passed, merges on one press or from anywhere but the panel, or the door stops carrying `merge`. **Not yet proved:** a press of Merge on GitHub on a real pull request; none was open, and the first is this line's own, once it is sent and its checks pass. **Named, not fixed:** the door's published tool schema describes every optional argument as "the ground this call names" and types each as a string, `number` included; the glass sends a number, and the door reads either.

### Send to GitHub on a line of work that has never been sent, and the GitHub tab's words for it (2026-10-09; the core's WHAT'S LEFT I1, third step; his word: "save and send, then go on to I1 step 3")

A line of work has an upstream only once it has been sent, and ahead and behind are counted against it, so on a line that had never been sent the GitHub tab greyed Send to GitHub as "Nothing to send — GitHub already has every save here" and called the line "not linked to GitHub at all", though the world sends to GitHub (atlas's `github-tab` line read so on 2026-10-08, and was landed onto main by hand for it). The door was never the reason: `git_push` sends a line that has no upstream with `--set-upstream origin <line>` and says "This branch had never been sent before; it now follows origin/<line>." Only the page judged.

Now the page asks whether the world has an origin before it judges a line with no upstream (the `git` tool lists the remotes either way). With one, Send to GitHub is live, titled "Send <line> to GitHub for the first time, and link it there", and the GitHub row reads "this line of work has never been sent to GitHub — Send to GitHub sends it there and links it"; Take from GitHub says "Nothing to fetch — this line has never been sent, so GitHub holds nothing of it". With no origin, Send is greyed "Nowhere to send — this world has no remote called origin", and only then is a line "not linked to GitHub at all". A line that has been sent reads as before. Nothing in the door moved: `flows.js` alone.

Held by `TestALineThatHasNeverBeenSentCanBeSent` (new, `server/glass_test.go`): the page asks for an origin in both `controls` and `plain`; Send's test is the upstream-or-origin one and says "for the first time"; Take's words for a never-sent line; "not linked to GitHub at all" only for a world with no origin, and once in the file; and the door's side read off its own files: `git` still writes `upstream` and `remotes`, and `git_push` still sends a never-sent line with `--set-upstream origin` and says it was never sent before. Written first and run alone against the page as it stood: red on its first check (exit 1). Then the page's two edits through `file_edit` on the line of work `first-send`, each on his Approve (the lower one refused once, the upper having landed first, and was parked again); `node --check`, gofmt and vet clean and the glass's suite green in his Bash tab (`db`, `handlers`, `server`). Built there (`e76db7e4...`) and placed on his card ("Place it, keep both"): the glass on 8091 as pid 9812; the build it replaced (`6a0b90d5...`, the tray's) kept as `atlas-webapp.old.exe` and the one before it (`e32db082...`) as `atlas-webapp.old-1008.exe`; the door untouched. Proved live on his glass, his session kept through the restart: on the GitHub page, atlas on `first-send` read "this line of work has never been sent to GitHub — Send to GitHub sends it there and links it", its Send live with "Send first-send to GitHub for the first time, and link it there", and the core's `main` unchanged ("in step", Send greyed with "Nothing to send").

**What goes red if unplugged:** `TestALineThatHasNeverBeenSentCanBeSent`, if Send is greyed on a never-sent line again, the words call such a line not linked at all, or the door stops writing `upstream` and `remotes` or stops sending a never-sent line to origin. **Not yet proved:** a press of Send to GitHub on a line that has never been sent. `first-send` itself was sent from his Bash tab on his word ("save and send, then do the pull request"), so the press waits for the next line he opens. **Named, not fixed:** the full GitHub page's council path greys its Push the same old way (`g.ahead`), and it goes through the core's own push, whose first send of a line was not read here; and the number on the GitHub door counts a world's unsaved files and its saves ahead of an upstream, so the saves on a line that has never been sent are not in it.

### Three versions cut and never tagged are marked after the fact: v0.1.2, v0.1.3 and v0.1.4, each on the commit that cut its number (2026-10-09; his word: "post-script the versions that didnt get any tags or info ... tag them with what happened that day")

His word, 2026-10-09, and his card the same morning, "Both": a mark for each version cut and never tagged, each on its own card, the sending of them a separate word of his. In this repository three: v0.1.2 on `0586228`, v0.1.3 on `7688c04` and v0.1.4 on `467e94c`, each the commit that moved every pin to its number (every VERSION file and Cargo.toml say it there), annotated, in his name, its message saying what the version carried. Each was checked first as `git_tag` judges a cut (a lawful name, every stamp agreeing at the commit, the commit on the main line, no mark of that name already) and cut in his Bash tab on his card. Left unmarked: 0.1.0+f1 and 0.1.1+f1, which no commit of this repository declares. `[0.1.2]`'s heading names its mark, with a post-script under it, and `[0.1.5]`'s note that 0.1.3 and 0.1.4 were built and never tagged carries a post-script naming theirs; the words under both stand as they were. Nothing was sent: the marks are on this machine until his word. No code moved, and nothing was built or placed. The core's half (its v0.1.6 and v0.1.10, and BUILDPATH's list of the marks) is in the core's CHANGELOG, "core: the versions cut and never tagged are marked after the fact".

**What goes red if unplugged:** in the core, `test_the_plan_names_every_mark_where_it_sits`, which reads every mark this repository holds and wants each one named in the core's BUILDPATH by the commit it sits on, with every `## ` heading here that places a mark; the core's entry carries the suites' tally.

### The door calls the pages what the glass calls them: Aider's refusal and a parked call's answer name the GitHub tab, and `file_edit`'s held note the card at the foot of the terminal (2026-10-09; the core's WHAT'S LEFT I1, second step; his word: "alright, let's get to work.")

On 2026-10-08 the glass's Version control page became GitHub, Guardrails became Laws, and every call parked for his hand became a card pinned in the dock at the foot of the front page's terminal (the entries below). The door's own words went on naming the old pages, and both entries named them and left them for a door rebuild: Aider's refusal on the main line (`aider.go`, "Lines of work on Version control"), `heldAnswer`'s reply to a parked call (`holds.go`, answered "on the Version control page"), and `file_edit`'s held note ("waits for his Approve on the card in Guardrails") with the comment at the head of `edit.go`.

Now Aider's refusal sends him to "(`git_branch new <name>` through the door, or Lines of work on the GitHub tab) and write on it"; `heldAnswer` tells a parked caller that he answers it "on the card at the foot of the front page's terminal (or on the GitHub page), or with hold_answer"; and `file_edit`'s held note reads "Nothing has been written. It waits for his Approve on its card at the foot of the front page's terminal.", the comment at the head of `edit.go` saying the same. Nothing else in the door moved. The core's half (the tree doors' refusal, a comment and a report, in `skills.py` and `pipeline.py`) is in the core's CHANGELOG, "core: the code calls the pages what the glass calls them".

Held by `TestTheDoorCallsThePagesWhatTheGlassCallsThem` (new, `edit_test.go`): every file of the tools package but its tests is read whole, every run of spaces made one and the glue between two string literals and between two comment lines taken out, so a name broken over a line's end is found, and the test is red on "Version control" or "Guardrails", naming the file; it fails outright if `aider.go`, `edit.go` or `holds.go` was not read; and its own samples hold it both ways (a name split over two literals and a name in a comment's wrap found, the practice in lowercase not). The capital is the mark: version control in lowercase is the practice, not the page. The core's stroke holds the core's half the same way.

Made on the line of work `page-names` through `file_edit`, each card on his Approve; gofmt and vet clean and the door's tools package green in his Bash tab (`go test ./internal/tools/`, 50.3 s); and the core's suites green from his glass on the same tree (strokes 3739 of 3739, smoke 72 of 72), whose strokes read `aider.go`. Built and placed on his card ("Place it, keep both", 2026-10-09 07:36): the new `atlas-mcp.exe` (sha256 beginning `a95a34db020e`) answers on 8090 as pid 28732, started on the line the door ran with; the door before it (`7625c07f...`, 2026-10-08) is kept as `atlas-mcp.old.exe`, beside `atlas-mcp.old-1007.exe` (2026-10-07); the glass untouched. The first start line opened the start script in Notepad instead of running it (`Start-Process pwsh -File` reads `-File` as its own `-FilePath`); nothing ran, and the door was started again with the script handed to pwsh in `-ArgumentList` (the core's CHANGELOG names the slip). Proved live: every `file_edit` parked through the new door since, the core's record of this piece among them, came back with the new note, and he answered each from the dock.

**What goes red if unplugged:** `TestTheDoorCallsThePagesWhatTheGlassCallsThem`, if a file of the door's tools package names Version control or Guardrails again, the name broken over two literals or a comment's wrap included, or if `aider.go`, `edit.go` or `holds.go` stops being read. **Named, not fixed:** comments in the glass's own scripts and styles still name the old pages (`agent.js` 1062, `chat.js` 49 and 58, `workflows.js` 4, `agent.css` 424, and lines in `app.js` and `app.css`); they are comments, not words on a page, and nothing holds them.

### Approve cards stay in view: every card that waits for his hand is pinned at the foot of the terminal, in every tab, and Guardrails is Laws (2026-10-08; his word: "just pin them to the window so they actually stay at the botton if there is an approve card. should work on all tabs", and on his cards "A pinned tray", "First in v0.2.5" and "It becomes Laws")

Read fast, the front page's cards scrolled out of his view (his words that morning: "when we are reading through a lot of stuff quickly those approve cards are being moved out of my view"). A shell line that asked first drew its card inside its own entry in the scrolled output, the council's question drew at the output's end, and every other call the door parked (a `file_edit`, the council's writes) was listed only in the Inspector's Guardrails tab, which had to be opened and found. His word on the line itself: "Same as the cards in the chat, they are perfect ... Basically replacing the Guardrails tab, it's kind of not needed and is another place to have to search for".

Now one dock (`#ag-dock`, a sibling of the output, under it and the play frame and over the hints and the line he types in) holds every card that waits for his hand, whichever tab is open: the council's question with its field (`gateHtml`), each shell line of this thread that asks first (`shellCardHtml`: the card the thread carried, now with the line it asks about over it, since its entry may have scrolled away), and every other call the door has parked (`holdCardHtml`: the tool, who asked, the world and when, and the call exactly as it was parked; a `file_edit`'s change line by line as the door sets it out for his card, coloured as a diff, anything else its arguments whole). `paintDock` draws it from what was read and asks the door for nothing; it is drawn again with the thread, and only when what it shows has changed, so a button is never swapped under his pointer, and what he types into the council's question survives a redraw. A card leaves when it is answered. The queue is read quietly every 3 s while the page is in view and again when he comes back to it (`HOLD_POLL_MS`, `watch`), so a call parked by anything but this page reaches the dock while he looks; a read overtaken by a newer one is not drawn. A shell card's buttons go to `decide`, as before, and its run lands in its entry; any other card's go to the new `holdAnswer`, which says it is running until the door answers, then writes the door's own answer into the thread as a line of its own, chipped `door`, red when the answer is a refusal, a call that errored, or a `file_edit` replay that refused. A door started without `--auth`, or a queue that could not be read, is said in the dock in red: an empty dock would look exactly like a guarded one. The entry a held line made says where its card is. The dock keeps the thread's look, is capped at 45% of the window, scrolls inside itself, and gives way before the line he types in does.

The Guardrails door and tab are Laws: the laws it drew under its queue since 2026-10-06 are the whole tab, drawn by the Laws page's own object as before. `guardsHtml`, its two hold buttons and the Inspector button's amber dot are gone. The GitHub page under Pages keeps its own list of what waits, as it had.

Made on the line of work `approve-tray` as 31 `file_edit` cards he approved in Guardrails, its last use, each file's cards parked bottom-most first; every file then checked on disk. Held by `TestWhatWaitsForHisHandIsPinnedAtTheFootOfTheTerminal` (new): the dock between the output and the line he types in; no card drawn in what scrolls; the dock drawing the question, the held lines, the queue and the note when the door holds nothing, and calling nothing itself; drawn again by `paintOut` and `readHolds`; the queue read quietly every `HOLD_POLL_MS` while in view and no more once the page has gone; the dock's one handler sending a shell card's buttons to `decide`, any other card's to `holdAnswer` and the field to `answer`, with the output's own handler answering none; every key it reads off the queue one that `hold_list` writes (read off `holds.go`); `hold_answer`'s refusals and errors beginning with the words the page reads; a file's change shown as `file_edit` sets it out (`edit.go`); no Guardrails left; and the dock's cap and scroll. And by the strokes moved onto it: `TestTheLawsTabCarriesTheLaws` (in place of `TestGuardrailsCarriesTheLaws`), `TestSettingsIsADoorOfItsOwn` (beside Laws), `TestTheFrontPageFakesNothing` (`hold_answer` from `decide` and `holdAnswer` alone, each called once) and `TestTheShellTabsAreTheDoorsShellAndNothingElse` (the card in `shellCardHtml`, read for the door's keys). In his Bash tab: `node --check` on the three scripts, gofmt and vet clean, and the glass's suite green (`server` 0.391 s, `handlers` 3.181 s, `db`); built there (`6a0b90d5...`, from `3cbb809` and the unsaved edits) and placed on his card (pid 25368 from 15:37:00; the glass before it, `e32db082...`, kept as `atlas-webapp.old.exe`; the door untouched).

Proved live on his glass after the restart (health 0.2.4, unlocked): the bar's doors read Aider Pair, GitHub, Audit Ledger, Workflows, Registry & Docs, Laws and Settings, and the Inspector's Laws tab drew the laws (16 standing rules, 10 ledger entries, 6 of 7 laws sealed) with nothing waiting in it. Two shell lines that ask first (`git -C atlas log --oneline -3` and `ps`) and a `file_edit` (RUNBOOK.md's line on where a held call is decided) were parked; a 33 KB read scrolled the output past them and the terminal turned to the Agent tab, and all three stayed pinned over the line he types in, the file's change coloured. He answered all three from the dock: both runs landed in their entries, the edit was written (its answer a `door` line in the thread, `written`, the file's CRLF kept), and the dock emptied and hid itself.

**What goes red if unplugged:** `TestWhatWaitsForHisHandIsPinnedAtTheFootOfTheTerminal`, if a card is drawn in what scrolls again, the dock stops being read or drawn, its buttons reach anything but `decide`, `holdAnswer` and `answer`, a key it reads leaves `hold_list`, or Guardrails comes back; `TestTheLawsTabCarriesTheLaws`, if the tab stops drawing the Laws page's own object; `TestTheShellTabsAreTheDoorsShellAndNothingElse`, if a shell card's buttons are drawn outside its card. **Named, not fixed:** Deny was not pressed live (the `ps` card meant for it was approved), so the dock's Deny is held by its stroke until the next card he denies; the council's question was not raised live (it needs a turn that asks). The door's own words still send him to Guardrails: `edit.go`'s held note ("waits for his Approve on the card in Guardrails") and its header, beside `heldAnswer`'s Version control (v0.2.5's second step, a door rebuild). The hold log still writes `approved_ran` for a `file_edit` replay that refused; the dock's `door` line now says it as a failure.

### An approved hold runs in the world it was held for: a council's call held for atlas, approved from the glass, had run in the core (2026-10-08; his card: "Fix it first")

On 2026-10-08 he opened v0.2.5's line from the Agent tab rather than with a button (his word: "kind of want to use the system to work on this kind of thing"): `mcp_call atlas git_branch {"project": "atlas", "action": "new", "name": "approve-tray"}`, in sitting 361. The door held it for atlas, as it holds any writing call from the council's key (`k-ae2e9481`), and his Approve in Guardrails ran it: "Opened "approve-tray" and moved you onto it." In the core. `hold_answer` ran a held call with the tenant of the call that answered it (`tool.Fn(t, run)`), and the glass answers every hold as its default world, `research`; the hold kept its world's name and nothing else. So atlas stayed on main, the core moved to a new line, and the hold's story was split across two logs: `held` in atlas's own (`atlas/state/holds.jsonl`, whose first line it was) and `approved_ran` in the core's. No call had been held for another world before, which is why it had not shown; a save or a send held for atlas would have run in the core the same way.

A `Hold` now keeps the tenant it was judged in (`tn`, unexported: the queue still shows only its name), set by `park`, which every held call goes through (`Call`, `file_edit`, `shell_run`). `hold_answer` runs the call in that tenant and writes its answer (`denied`, `vanished`, `approved_ran`, `approved_errored`) to that world's log, beside the line that held it. A hold that carries no world is refused by name and nothing runs: a held call is never run in whatever world the answer came from. A call held and answered in one world runs and logs exactly as before, so the other strokes that answer holds (`edit_test.go`, `shell_test.go`, `tools_test.go`) did not move.

Held by `TestAHeldCallRunsInTheWorldItWasHeldFor` (`holds_test.go`): a call held for a second world, `atlas`, and approved from the default one runs in `atlas`, with its `held` and `approved_ran` lines both in that world's log and neither in the answering world's; and a hold with no world is refused with nothing run. Written on the line of work `approve-tray` through `file_edit`, each card on his Approve; gofmt and vet clean and the door's tools package green in his Bash tab (`ok atlas/line/internal/tools 49.295s`); built in his PowerShell tab (`7625c07f...`) and placed on his card (pid 33236 from 13:36:51; the door before it, `23adcaa3...`, kept as `atlas-mcp.old.exe`), restarted with the flags it ran with. Proved live in sitting 362: the council's `git_branch switch approve-tray`, held for atlas and approved from the glass, answered "Refused: there is unsaved work here", atlas's own answer (the core was clean and has no such line, where it would have been refused as an invalid reference), with both its lines in atlas's log and none in the core's.

**What goes red if unplugged:** `TestAHeldCallRunsInTheWorldItWasHeldFor`, if `hold_answer` runs or logs a held call in the answering world again, or runs one that carries no world. **Named, not fixed:** the hold log writes `approved_ran` for a `file_edit` whose replay refused ("the file has changed since the card was shown"): a refusal is the tool's answer, not an error, so the log cannot tell a written edit from a refused one. `heldAnswer` still names the old page (v0.2.5's second step).

## [v0.2.4] — 2026-10-08 10:51 (tag on af92401)

### GitHub is a tab of its own: Version control is renamed GitHub, with a door and an Inspector tab beside Aider Pair, drawn by the page's own object (2026-10-08; his word: "it should be in its own github tab on the inspector so its less confusing", "let's make sure we are using the existing infrastructure to underpin the new tab", and "let's rename it to Github. and then place in the new strokes")

From 2026-10-06 the Version control page's own object (`Flows.render(el, inPanel)`) was drawn under Aider in the Aider Pair tab, and from 2026-10-07 the unsent-work number sat on the Aider Pair door, so Aider Pair meant two things: Aider, and both worlds' git. On his word the page is renamed GitHub (`App.PAGES`, which the Pages menu, the palette and the crumb read, and the page's own title in `flows.js`) and has a door and an Inspector tab of its own beside Aider Pair (`agent.js`: `DOORS`, `TABS`). The GitHub tab draws the same object as before, told it stands in the panel (`githubHtml`, then `Flows.render($ag('ag-vc'), true)` in `paintPane`), with the whole-world diff under it: nothing about how it works changed, only where it stands, and it is still one object drawn in two places, so the tab and the page cannot drift. The unsent-work number moved to the GitHub door, its tooltip saying the GitHub tab is under it. Aider Pair draws Aider alone and ends with one line, on his card: "Aider's edits are unsaved work on your line: save and send them in the GitHub tab", whose GitHub opens the tab (`paneClick`'s `door`). The words that sent him to Version control by name (the Aider tab's title, its banner, a run's Undo line and the Workflows tab's button) say GitHub.

Made on the line of work `github-tab` in atlas as 20 `file_edit` cards he approved in Guardrails, each file's cards parked from the bottom up so each still matched when its turn came; every file on disk was then checked by hash against the exact text the cards showed.

Held by two strokes that replace the two that held the old arrangement. `TestGitHubIsATabOfItsOwn` (in place of `TestAiderPairCarriesVersionControl`): GitHub among the Inspector's tabs beside Aider Pair; the GitHub branch of `paintPane` draws `githubHtml` and the page's own object, and nothing else on the front page calls `Flows.render`; the whole-world diff under it; Aider Pair no longer draws the panel or the diff, and its one line opens the GitHub tab; no "Version control" left in the front page's words to him; in the panel the repositories alone; the page called GitHub at its head and in `App.PAGES`, and still drawn whole under Pages; and the panel's look. `TestUnsentWorkIsANumberOnTheGitHubDoor` (in place of `TestUnsentWorkIsANumberOnTheAiderPairDoor`): the number's place on the GitHub door and on no other, the door beside Aider Pair in the bar, and the reads that keep it. Both were run against the edited text before the cards were parked, and the old two failed on it. In his Bash tab: `node --check` on the three scripts, gofmt and vet clean, the two strokes run by name and passing, and the glass's suite green (`server`, `db`, `handlers`). Built there (`e32db082`) and placed on his card (the glass before it, `58f23d76`, kept as `atlas-webapp.old.exe`; the door untouched).

Proved live on his glass after the restart (pid 5584; health 0.2.2, the lock on, his session kept): the bar's doors read Aider Pair, GitHub, Audit Ledger, Workflows, Registry & Docs, Guardrails and Settings, with 4 on the GitHub door (the four files this piece changed, unsaved on `github-tab`) and its tooltip; the Inspector's GitHub tab drew both worlds in plain words with their buttons, and the whole-world diff under them; Aider Pair drew Aider and the one line, whose GitHub opened the GitHub tab; the Pages menu lists GitHub; and /flows opens titled GitHub, its crumb ATLAS / GitHub, drawn whole with the hold queue, the repositories, the council path and Recent.

**What goes red if unplugged:** `TestGitHubIsATabOfItsOwn`, if the tab, its pane, the page's own object, the one line, the name or the panel's shape goes; `TestUnsentWorkIsANumberOnTheGitHubDoor`, if the number leaves the GitHub door or the door leaves the bar. **Named, not fixed:** the door's own messages still say Version control in two places (`aider.go`: "Lines of work on Version control"; `holds.go`: a hold is answered "on the Version control page", which has been Guardrails' since 2026-10-06), and changing them is a door rebuild; comments in `chat.js` and `workflows.js` keep the old name as history; the page's own name inside the glass stays `flows` (`/flows`, `Flows`), as it did through the 2026-09-10 rename; and the core's WHATS_LEFT.md says, of H16, that unsent work is a number on the Aider Pair door.

### Pull requests through the glass: `git_pr` lists a world's open pull requests with their checks and opens one from a line of work on his click; it never merges (2026-10-07; his word: "set up the system to do a pull request and then get it pulled into main. full on version control", and on his card: "Open the wall for PRs")

Pull requests were kept out of the door on purpose: `gitctl.go` has said since it was written that what `gh` alone would add is GitHub-side objects (pull requests, issues, releases, CI runs), "a wall to open deliberately, not a dependency to acquire by accident", and the entry of 2026-09-18 below names them as deliberately not built. His card opened that wall for pull requests and for nothing else, and `gitctl.go` now says so under that passage. `git_pr` (`line/internal/tools/gitpr.go`) runs `gh`, the GitHub CLI already signed in on this machine, in the world's own ground through the door's `spawn` (stdin closed, every call bounded: 30 s to list, 60 s to open) with nothing that can stop for a human (`ghEnv`: no prompt, no pager, no colour, no update notice). `list`, its default and a read, gives the world's open pull requests, each with its line, where it goes, its link and its checks as one word (none, pending, failed or passed; a failure outranks anything still running). `open` takes the line he stands on into the main line, titled with the line's newest save and listing its saves, with no name but his on it; it is refused on the main line, over unsaved work, and for a line GitHub does not have exactly as it stands here (send it first), and a refused call never asks `gh`. It never merges, and any other action is refused saying so: a line still lands as it always has (move to the main line, Land onto main, which only fast-forwards, then Send), and GitHub marks the pull request merged once the main line holds its saves. It stands behind the same wall as every send (`GIT_REMOTE`) and is his glass's alone (`ServiceOnly`, and refused in the tool itself for any other caller). `tools.go` registers it with one line, `r.add(gitPRTool())`, and `docs/TOOL_PERMISSIONS.md` is regenerated (95 tools, 7 service-only).

On Version control each world has a Pull requests button beside its others (`webapp/static/js/flows.js`). It lists that world's open pull requests (number, title, which line into which, the checks word and the link), offers Open a pull request for the line the world stands on, keeps the door's answer under the list, and says there that a pull request lands nothing by itself. The page decides nothing of its own, and nothing on it merges.

Aider wrote `gitpr.go` and `gitpr_test.go` on the line of work `edit-tool` in eleven runs, each read back before the next, and gofmt finished their layout on his card. The registration in `tools.go`, the ruling in `gitctl.go`, the panel in `flows.js` and its stroke in `glass_test.go` went in through `file_edit`, each written on his Approve in Guardrails.

Held by `gitpr_test.go`, where a stand-in for `gh` records every call and no stroke reaches GitHub: `TestPullRequestsAreListedWithTheirChecks`, `TestAPullRequestsChecksReadAsOneWord`, `TestPullRequestsAreHisGlassAloneAndBehindTheWall`, `TestOpeningAPullRequestIsRefusedWhereItMustBe` and `TestASentLineOpensAPullRequestIntoTheMainLine` (a real line, sent to an origin on this disk); by the glass's `TestVersionControlListsAndOpensPullRequests` (the button, its list, one open call and no merge, and `gitpr.go` still declaring `git_pr` with `list` as its read); and by `TestEveryToolHasAPermissionRecord`, `TestThePermissionRecordHasOneRowPerTool` and `TestEveryServiceOnlyToolRefusesEveryOtherCaller`, which now reaches it. In his Bash tab: `node --check` clean, vet and both modules' suites green, gofmt clean; both built there and placed on his card, the builds before them kept as `.old.exe` (door 3a9ebde9, glass eb87f7a0), running as door 23adcaa3 and glass 58f23d76. Proved live: `edit-tool` was sent on his card and its push ran green in 3 m 0 s; pull request #2, `edit-tool` into `main`, was opened from Version control on his click; and Pull requests listed it with its checks passed.

**What goes red if unplugged:** the five strokes above, if a guard, the list or the open moves; the permission record's strokes, if the line comes out of `tools.go` or the declaration changes (it writes, `list` is its read, it is service-only); and `TestVersionControlListsAndOpensPullRequests`, if the page loses its button, its list or its one open call, or learns to merge. **Named, not fixed:** the glass's `TestEveryToolTheFrontPageAsksForIsCarriedByTheDoor` reads the door's tool names out of `tools.go` alone and the page's asks out of `agent.js` alone, so a tool declared in a file of its own (`file_edit`, `git_pr`) is one it cannot see; Version control's call to `git_pr` is held instead by the stroke that reads `gitpr.go`. RUNBOOK.md says the door serves 89 tools; it serves 95, and it does not yet say that `git_pr` needs `gh` installed and signed in on this machine.

### The operator's edit: `file_edit`, one exact edit of one text file, parked on his card and written only on his Approve (2026-10-07; his word: "go, build the edit tool in the door", and "use Aider, isn't that what the whole point of that tool is?")

Until today the hand changed a file through the Python tab with a re-spelling of its own that kept backslashes, `..` and line breaks out of the gate's sight, and the write card showed a statement, not the change. On 2026-10-07 that was stopped as not transparent, and he ordered a real edit tool in the door instead. `file_edit` (`line/internal/tools/edit.go`) takes three plain fields, `filepath`, `old` and `new`, never packed into a command. It refuses a path outside the world (the door's `jailed`), a secret, client material, a home or a reach outside the ground (`shellRefusals`, read on the path), what the seats and Aider never write (`aiderNeverWritten`), and any edit while an engine or a sitting is open (SITTING LAW 5); the old text must appear exactly once; the file keeps its own line endings, and a mixed one is refused; it takes UTF-8 text up to 4 MB. The first call writes nothing: it parks itself as `shell_run` does, carrying the change line by line, and the Inspector's Guardrails list ("Waiting for your hand") shows it with Approve and Deny. His Approve replays it through `hold_answer`, which reads the file again and writes only if the change is still exactly the one his card showed. It is his glass's alone (`ServiceOnly`). `tools.go` registers it with one line, `r.add(fileEditTool())`, put there through the Python tab on his card as one plain replacement; `docs/TOOL_PERMISSIONS.md` is regenerated (94 tools, 6 service-only).

Aider wrote both files on the line of work `edit-tool`, an instruction at a time (the door refuses one over 4,000 characters: "it is not an instruction, it is a file"), each read back before the next. Two runs that put a function inside another were taken back with Aider's own undo; runs that ran past 600 s, or whose edits matched nothing in the file, wrote nothing; short follow-ups fixed the rest, and gofmt finished the layout on his card.

Held by `edit_test.go`: `TestAnEditWaitsForHisApprovalAndWritesExactlyWhatWasShown`, `TestAnEditKeepsTheFilesOwnLineEndings` (line-feed text into a CRLF file stays CRLF), `TestAnEditIsRefusedWhereItMustBe`, `TestAnEditThatNoLongerMatchesWritesNothing`, `TestNoEditIsMadeWhileASittingIsOpen` and `TestFileEditIsHisGlassAlone`; with `TestEveryServiceOnlyToolRefusesEveryOtherCaller`, which now reaches it, and the permission record's two strokes. gofmt, vet and the door's whole suite green in his Bash tab; built there and placed on his card (the door before it kept as `atlas-mcp.old.exe`; the glass untouched). Proved live: this entry is the tool's first edit, parked with its change, shown in Guardrails, and written on his Approve.

**What goes red if unplugged:** the six strokes above; take the line out of `tools.go` and they go red with the permission record's two. **Named, not fixed:** the two refusals Aider's gate would not carry into an instruction (a secret file and a parent step) are held by the strokes of `shellRefusals` and `jailed`, not by `file_edit`'s own. A held edit shows in Guardrails as its arguments, not in the terminal's thread. RUNBOOK.md says the door serves 89 tools; it serves 94. Aider ran at about 1.5 tokens a second here, against 3.6 measured on 2026-10-04.

### A line from a script is a typed line: Agent.drive, kept in the front page (2026-10-07; his word: "write those helpers in so you don't need to rewrite them every time")

The hand drives the front page's terminal from the page's own script (a read in the Bash tab, a card it waits on, a turn), and the few lines it did that with lived only in the tab and were typed in again after every reload. `Agent.drive(tab, line, waitMs)` keeps them in the page: it turns to the tab as the tab's own button does, hands the line to `go()` exactly as if it had been typed there, and resolves with what the line made once it stops running (a shell's held, ran or refused, with its hold, class, reasons and output; a turn once it ends; a page command once it finishes), or null when nothing appears in time. Every gate, card, refusal and record is the typed line's own: it decides nothing, answers no card, and reaches the door by no road of its own.

Held by `TestALineFromAScriptIsATypedLine` (drive hands the line to `go()`, turns only to a wired tab, and calls no tool, API, fetch, hold answer, card decision, council send, shell or Aider run of its own), with `TestNoMethodOfTheFrontPageIsDefinedTwice` beside it (the page's own `send()` is the council's, so the helper is not called that). `node --check`, gofmt and vet clean and the glass's packages green in his Bash tab; built there and placed on his card (the glass before it kept as `.old.exe`; the door untouched). Proved live after a reload, with none of the old lines typed in again: `Agent.drive('bash', 'tail -n 1 sessions/sessions.jsonl')` turned to the Bash tab and came back ran, class read, exit 0, with the ledger's last line; and a read card sent through it from the Python tab came back held, with its hold id and the gate's reason, and ran on his click.

**What goes red if unplugged:** `TestALineFromAScriptIsATypedLine`.

### The front page is the window: it takes exactly the window's height, never scrolls, and the Inspector lies over the terminal on a narrow screen (2026-10-07; his word: "it should fit the browser window it sits in naturally", having sized his pane to the layout he wanted, and on his card "Fill, no scaling")

The front page took its height from a formula (the window less 125 px) with a floor of 560 px, so a window shorter than about 685 px cut the bottom off and the page scrolled; and below 1024 px wide the Inspector stacked under the terminal, which made the page two windows tall. Now the page is the window: `body.is-agent .main` and `.ag` are the window's height with nothing to scroll, the bar and the footer keep their size, and the terminal and the Inspector fill what is left between them (`.ag-main` with no floor, its two panes stretched), each scrolling inside itself as before; the text stays the size it is, and a wide window keeps the 1280 px middle. Below 1024 px the Inspector lies over the terminal's right side inside the window instead of stacking (below 768 px it covers the terminal, 12 px in from the edges), and closes as it always has. Every other page scrolls as a document, as before.

Held by `TestTheFrontPageIsTheWindow`, which replaces the stroke of 2026-10-03 that held the old fixed height and the stacked narrow layout down: the page and its frame are the window's height and hide their overflow, the terminal's row has no floor and fills, neither pane is given a height of its own, the narrow block does not stack and lays the Inspector over the terminal, and the output box scrolls. gofmt and vet clean and the glass's packages green in his Bash tab; built there and placed on his card (the glass before it kept as `.old.exe`; the door untouched). Proved live in the pane with the Inspector open: at his 1233 x 702 the page did not scroll and the terminal and the Inspector stood 584 px tall side by side with the footer on the window's edge, as in the screenshot he chose; at 900 x 600 the Inspector lay over the terminal's right side, inside the window; at 700 x 480 it covered the terminal with the input line and the footer in view; at 1233 x 480, which the old floor would have scrolled, both panes were 362 px tall and nothing scrolled; and Version control still scrolled as a document (1504 px in a 702 px window).

**What goes red if unplugged:** `TestTheFrontPageIsTheWindow`. **Named, not fixed:** the 6 px of sideways scroll measured in his pane that morning did not appear at any of the four sizes; the front page now hides any overflow rather than scrolling to it.

### The sidebar and the old Dashboard retire: every page is one list the front page reads, and what they carried is on the front page (2026-10-07; his word: "go, next: retire the sidebar and the old dashboard", and on his cards the covenant in the front page's footer, unsent work on the Aider Pair door, the idle line in the title, and every other page keeping just the crumb)

The last step of his design pass of 2026-10-06, once every tab it laid out was in the Inspector. The sidebar is gone from every page (its markup in `index.html`, its rules in `app.css` and `agent.css`, and `css/icons.css`, its glyphs), and the old Dashboard with it (`js/home.js`, its case in the page table and its line): `/dashboard` opens the front page, so a bookmark still lands. Every page is one list, `App.PAGES`, which the front page's Pages menu, the palette (Ctrl K) and each page's crumb read, so a page is named once and the crumb's ATLAS is every page's way home. What the sidebar carried is on the front page. Its bar already had the version, Boot and Close, the door's state, and who is signed in with Lock; on his cards, the covenant, read off the record as before (`App.paintCovenant`), is in its footer; unsent work is a number on the Aider Pair door, the door Version control stands under, read as the glass loads and after every act on Version control and hidden while nothing is owed (`App.paintOwed`, `App.showOwed`); and the idle line the Dashboard's hero said (the core's WHAT'S LEFT C17) is in the terminal's title, from five idle minutes, with the thirty-minute close and the minutes left (`Agent.idle`). What the Dashboard did for other pages moved to their owners: the thread kept at a turn's end and restored while its engine stands (`Chat.keep`, `Chat.showKept`), the line Version control's Recent stages into the terminal (`Chat.pending`), and the Projects card, which stands only in the Run tab now. The sidebar's own reads (its badges and its health check) went with it, and the words on other pages that sent him to the Dashboard send him to the front page.

Held by `TestTheFrontPageIsRoutedAndTheSidebarIsRetired` (no sidebar, glyphs or `home.js` in `index.html`; `/` and `/dashboard` to the front page and no Dashboard in the page table; every line of `App.PAGES` a page the router draws, at the address it reads it from; the Pages menu, the palette and the crumb reading the list; Boot and Close through the front page), `TestTheFooterReadsTheCovenantOffTheRecord`, `TestUnsentWorkIsANumberOnTheAiderPairDoor` and `TestTheFrontPageSaysWhatIdleCosts`; the watched-turn, toll, Audit Ledger, Run, Laws and What's-left strokes moved off what retired onto what stands. The seven scripts parse (`node --check`), gofmt and vet clean, and the glass's packages green in his Bash tab; built there and placed on his card (the glass before it kept as `.old.exe`, the new build 76 KB smaller; the door untouched). Proved live on his glass: no sidebar on the front page or on Version control, whose crumb led home; `/dashboard` opening the front page; the Pages menu listing the ten other pages and the palette all eleven; the covenant (`1512741580b7239b`) in the footer; 15 in amber on the Aider Pair door, the 15 files Version control showed changed; and, in sitting 356 on his card, one turn ("say hello in one short line", 53 s, tolled) kept at its end and back after a reload, staged from Version control's Recent into the terminal without running, and the idle line in the title at five idle minutes.

**What goes red if unplugged:** the four strokes named above, and the moved ones. **Named, not fixed:** the Dashboard's rules in `app.css` (the deck, the hero, the brief, `home-*`) are read by nothing now, and comments in `api.js`, `app.js`, `council.js`, `projects.js` and `handlers/trace_test.go` still name the retired `Home`. The front page scrolls sideways by 6 px in his pane (1238 x 702), measured during the live proof; fitting it to the window is planned on his word, and not built.

### Run carries the Watchboard's wire and the Dashboard's Projects, and a pick plays beside the terminal (2026-10-07; his word: "go, next: run", and on his cards the wire in place of the tab's own events, and a pick played "Beside the terminal")

The front page's Run tab drew the last turn's stages, its seats and a list of the turn's events of its own, which left the tokens out. That list is gone: the turn's events are drawn by the Watchboard's own wire (`Chat.paintWire`, which now takes the box it draws into and draws the Watchboard's own when handed none), with its one raw switch -- folded is a line per event with the tokens counted, raw is every frame whole -- so there is one wire and not two. The failures stay on the face of the turn, in its stages and NOT EVERYTHING RAN. Under the turn the tab draws the Dashboard's Projects by the card's own object (`Projects.render(el, inPanel)`): what the maker made, each one's versions and which is in hand, with Work on this and Put it down sent to the council through the terminal's Agent tab exactly as if typed there (`Agent.council`), and Open in its own tab. In the panel it draws no frame of its own: picking a project, or a version of one, plays it in the front page's one play frame beside the terminal (`Agent.pickProject`), which only a turn's delivery could fill until now, and brings it back if it was put away; a seat's words still never reach it. As a turn's events arrive only the turn is drawn again (`paintRun`), so the projects are not read again with every token; they read themselves again when the turn ends, as on the Dashboard. The Watchboard and the Dashboard under Pages still draw theirs whole, and the Watchboard's raw button now says which way the one switch stands. In the panel they take the Inspector's look (`.ag-wire`, `.ag-proj`), and the old list's rules went with it.

Held by `TestRunCarriesTheWireAndTheProjects`, beside `TestTheFrontPagePlaysWhatTheMakerMade`, whose words now name the second way onto the play frame. The three scripts parse (`node --check`) and the glass's packages are green in his Bash tab, gofmt and vet clean; built there and placed on his card (the glass before it kept as `.old.exe`; the door untouched). Proved live on his glass: the Run tab drew the thirteen projects with no frame of its own; picking countdown-timer played it beside the terminal (one frame on the page, no `sandbox` attribute), its version 1 played from the version list (`?v=1`), and after Put away a second pick brought it back; Work on this, pressed with the terminal on Bash, moved it to the Agent tab and was refused there as typed words are with no engine open; and one turn on his card (sitting 354, "good morning", delivered in 4.9 s, tolled) drew its stages, its seats and the wire: thirteen events, eight lines folded with the six token frames counted, thirteen raw, and the Watchboard's own button reading the same switch.

### Aider Pair carries Version control: the Version control page's own repositories, drawn under Aider (2026-10-06; his word: "go, next: aider pair")

An Aider edit is unsaved work on a line, and Version control is where it is saved; the tab drew its own small card per world (branch, changed, untracked, ahead, behind) and sent him to the Version control page to do anything. Under Aider it now draws that page's own object (`Flows.render(el, inPanel)`), told it stands in the panel, so there is one Version control and not two: each world in plain words -- the line it stands on, where it sends, the last save, what is unsaved, how it stands with GitHub -- with Save the work, Send to GitHub, Take from GitHub, its lines of work (move, land onto main, finish, start one) and its version marks (cut, send, remove, each armed), and its changed files each opening their diff. In the panel it draws the repositories alone: the hold queue is Guardrails', the way through the council is the terminal itself, and Recent is the terminal's own history. The tab's own world cards and its Version control button are gone; their Show the diff stays, on his card, as one Every change in a world button per world under the panel. The page under Pages still draws it whole. In the panel it takes the Inspector's look (`.ag-vc`).

Held by `TestAiderPairCarriesVersionControl`, beside the Version control page's Land stroke, unchanged and green. The glass's packages green in his Bash tab, gofmt and vet clean; built there and placed on his card (the glass before it kept as `.old.exe`; the door untouched). Proved live on his glass, reading only (nothing saved, sent or moved): under Aider's own cards the tab drew atlas and research in plain words with their five buttons each (Send and Take greyed, with nothing to send or take), atlas's lines of work (main, on GitHub), its four changed files, and every change in atlas at once from the kept button; no hold queue and no council card in the panel.

### Guardrails carries the Laws: the Laws page's own object, drawn under the hold queue (2026-10-06; his word: "go, next: guardrails")

The front page's Guardrails tab held the council's question and the writing calls parked at the door for his hand. Under them it now draws the Laws -- the Laws page's own object (`Laws.render(el, inPanel)`), told it stands in the panel, so there is one Laws and not two: the chain's verdict from `law_status`, why a law cannot be set or sealed while it cannot, the ledger's drafts with Seal behind its inline question, Set a law, the standing rules, the ledger's entries and every law on the chain, each opening in place, with the page's own find. In the panel it drops the page's title for a compact head (the counts, the find box, Read it again), and Read it again, Set and Seal draw it again where it stands rather than on the Laws page; the page under Pages still draws it whole. In the panel it takes the Inspector's look by pointing app.css's tokens at its palette (`.ag-laws`).

Held by `TestGuardrailsCarriesTheLaws`, beside the Laws page's own four strokes, unchanged and green. The glass's packages green in his Bash tab, gofmt and vet clean; built there and placed on his card (the glass before it kept as `.old.exe`; the door untouched). Proved live on his glass, reading only (no law set or sealed): the tab, under the hold queue ("Nothing is waiting"), showed 16 standing rules, 10 ledger entries and 6 of 7 laws sealed, the chain proving whole (10 links, head `237a5263b6f27eca`), and 33 rows; the find narrowed "publish" to RULE 12 and the ledger's entries 6 and 10; SITTING_LAWS.md opened in place; Read it again drew the laws again in the tab; and /laws still drew the page whole, title and all.

### Registry & Docs carries Tools with Call, the seats and Records' documents, and Settings is a door of its own (2026-10-06; his word: "go, next: registry & docs", and on the card: "Add another page along the top next to guardrails that is for settings. that'll include all the settings for messaging integration.")

The front page's Registry & Docs tab listed the tools and the documents. Each tool row now carries Call, which opens the Tools page's own form (`App.invokeTool`, in the modal: the tool's arguments over a JSON box, the answer with its hash and trace, and Call again). The fourteen seats are drawn by the Agents page's own card (`App.seatCard`), told it stands in the panel (`inPanel`): its name no longer leaves the front page, and a file button opens the declaration as it is on disk with its sha256 under the card, where the Agents page's seat page showed it; the prompt opens in place, and the Agents page's own filter narrows them. The documents carry what Records shows: when each changed, a sealed badge, every document of the kind, and the logs' newest-60 note.

Settings is the bar's sixth door, beside Guardrails, and an Inspector tab: the door's address (`mcp_url`) and the evals pass mark (`eval_threshold`), each shown as what is stored and what is used, saved only by the button beside it, and refused before saving unless it is an http(s) address or a number from 0 to 1 (or empty, for the default); and the messaging bridge as the door reports it (`team_status`), with how Discord, Slack and WhatsApp are connected -- each webhook URL and the hook secret placed by hand in the world's chat secrets file, never shown or taken by a page. The Settings page's three Provenance badges are not carried: that page wrote them in, and nothing reads them.

Held by `TestRegistryAndDocsCarriesTheToolsTheSeatsAndTheDocuments` and `TestSettingsIsADoorOfItsOwn`. The glass's packages green in his Bash tab, gofmt and vet clean; built there and placed on his card (the glass before it kept as `.old.exe`; the door untouched). Proved live on his glass: six doors in the bar with Settings beside Guardrails; Call on `muster` from the tab answered in the Tools page's form ("2 carried projects: atlas, research", with its hash and trace); the fourteen seat cards, none linking away, a prompt opened in place, Deep Researcher's file served with its sha256 (`e43fa6ada9ef7407`), and the filter narrowing "steward" to Neiro and Steward; doctrine's eight documents with their dates and the seven laws' sealed badges, and the logs' note; and Settings showing both settings stored as nothing (the address used as `http://127.0.0.1:8090`, the pass mark as 0.5) and the bridge's three platforms disconnected, with `foo` refused as an address and nothing stored.

### The Audit Ledger tab carries Records, Evals and the live check: what the build has proved, the estate and the live standups, the live check's button, and the evals scored by hand (2026-10-06; his word: "go, next: audit ledger")

The front page's Audit Ledger tab read the proof and drew only its sittings, in a card of its own. It now draws what the build has proved -- strokes, smoke, the live standup with its move against the last one, runs green, parity -- and the estate behind it -- sittings, tolls (with those written into SEAT_LOG), runs, the recent sittings and the newest five live standups -- by the console's one reader of the proof (`App.paintProof`), from the one `proofs` read the tab already made; its own Sittings card is gone. The rows were the Dashboard deck's own, so they are split out of it into `App.proofLedger` (and `paintProof`'s `ledger` placement), which the deck now calls too: the two cannot drift. The Dashboard's live check is in the tab (`standup_run`, the morning set), offered only while no engine is open, since the door refuses it then and the tab says why; its button stays grey through every redraw until it ends, the tally comes back in the card, and the proof is read again. The evals a person scored by hand are read where the Evals page reads them (none yet). The engine's report, the transcripts and the chain note stay. In the panel the proof takes the Inspector's look by pointing app.css's tokens at its palette (`.ag-proof`), and a tile's source line sits under its note, where a tile that narrow had cut it.

Held by `TestTheAuditLedgerCarriesTheProofAndTheLiveCheck`: one `proofs` read and the two places it is drawn; both drawn by `App.paintProof`; the deck and the tab drawing the rows from one function; the live check only from its button, as the morning set, and never offered while an engine is open; the evals read; the panel's look. The glass's packages green in his Bash tab, gofmt and vet clean; built there and placed on his card (the glass before it kept as `.old.exe`; the door untouched). Proved live on his glass: the tab showed strokes 3722/3722 and smoke 72/72 (10:53), the live standup 9/9 (2026-10-05 21:34), runs green 63/86 and parity 0.6808 -- the same five rows the Dashboard's deck drew from the same function; the estate (352 sittings, all closed; 280 tolled, 308 written into SEAT_LOG; 1,737 runs) and the last five of 86 live standups; "None scored yet"; and Run the live check offered, no engine being open. The live check itself was run from the tab on his own press: 9/9 in seven and a half minutes (sitting 353, tolled; `logs/standup_2026-10-06_114927.md`), its button grey until it ended, the tally in the card, and the rows read again (the live standup 9/9 at 11:49, runs green 64/87).

### The Workflows tab is the builder: the Inspector draws the Workflows page's own builder, Fire refuses an empty box and a second press, the door refuses a fire while a run is moving, and a refusal reaches the page in words (2026-10-06; his word: "go, start with workflows", the builder "Inside the Inspector")

The front page's Workflows tab drew the flows and the runs as text under a button that went to the Workflows page to do anything with them. It now draws that page's own builder -- the same object (`workflows.js`, `render(el, inPanel)`), without the page's header -- so there is one builder and not two that drift: the flows and their Open, a new flow, the steps and edges, Save, the boxes a flow needs and Fire, a gate's two buttons, and the runs waiting on him and the runs before. The page under Pages draws the same object, whole; the tab's "Open the builder" button is gone and "Version control" stays. The flow in hand outlives the redraw the tab gets every time the Inspector opens (an unsaved edit was lost to it), a run that comes back while the tab is shut waits on the record instead of throwing, and the builder takes the page's look by pointing app.css's tokens at the Inspector's palette (`.ag-wf`). The Inspector's tab is named "Registry & Docs", as the bar's door is.

The Fire hole (HANDOFF FOR 2026-10-06), closed both ways. The page refuses an empty box by name, and a second press while a run fired from it is moving, before the door is asked ("Nothing was fired."): a box left empty had gone to the engine as an empty value, which nothing refused, so the release flow's first step was handed no mark. And the door's `flow_run` refuses while a run is moving on the world, as `flow_resume` has since 2026-10-05 ("a run is moving on this world -- fire again once it stops"): it waited on the world's flow lock, so a second Fire queued behind the run in flight and started, unwatched, the moment that one stopped. A replay still waits its turn.

A refusal reaches the page in words. Every glass handler answers a failure as `{"error": words}`, and the page's client (`API.get`, `API.post`) threw `HTTP 502` and dropped them, so the door's refusals reached every page as a number -- the gate's "a run is moving on this world" among them since 2026-10-05. `API.fail` reads the words, and keeps the number only when there are none.

Held by `TestAFireIsRefusedWhileARunIsMoving` (the door: refused by name while the world's lock is held, nothing recorded, and the same fire runs once nothing moves; it fails at 5 s if the fire queues instead), `TestTheInspectorDrawsTheBuilderItself`, `TestFireRefusesAnEmptyBoxAndASecondPress` and `TestARefusalReachesThePageInWords` (the glass). The door's tools and flow packages (67 s, 7 s) and the glass's packages green in his Bash tab, gofmt and vet clean; both built there and placed on his card, last night's builds kept as `.old.exe`. Proved live on his glass, in the Inspector's tab: `release` fired with `mark` empty was refused on the page and `flows/runs.jsonl` did not move (339,736 bytes); a new flow saved with an eval that checks nothing came back as the door's own words (`refused: eval node "eval" names no node to check`, where it read `HTTP 502`) and nothing was folded; the unsaved flow was still in hand after the Inspector was shut and opened; and `smoke`, fired once on his word, ran COMPLETE in 5,074 ms with its waterfall in the panel, a second press while it moved was refused on the page, and the record grew by that one run.

## [v0.2.2] — 2026-10-05 20:51 (tag on acdceb6)

### A gate's decision answers the gate it was shown at: the door refuses one while a run is moving, and the page greys its gate buttons at the first click (2026-10-05, found by the first release run)

The first run of the release flow stopped at its save gate in the second it paused there: a "Stop it here" made while the first part was still running (the first gate's buttons stayed live for those fifteen minutes) had waited on the world's flow lock and answered the next gate, which nobody had seen. A "Carry it on" made the same way would have saved and sent past it. `flow_resume` now refuses while a run is moving on the world ("a run is moving on this world, so there is no gate to answer"), and nothing is decided; the Workflows page greys both gate buttons at the first click until the run stops and the panel is painted again, and gives them back if the door refuses the move. Held by `TestADecisionIsRefusedWhileARunIsMoving` and `TestAGatesButtonsGoGreyAtTheFirstClick`.

### The tool node: a flow calls the door's own tools by name, and the release's steps are one door tool (2026-10-05, the core's release workflow; his word: "make the release one workflow")

A ninth node kind, `tool`, calls one of the door's tools by name with the arguments the node carries, each rendered like a question (`{{mark}}`, `{{out_x}}`). No model is asked and no engine is opened, so a flow of tool nodes runs where a sitting must not be open: the live check opens its own. The flow law judges the node's shape -- a tool named in the door's shape (`ToolRe`), arguments named the same way, and no other kind carrying `tool` or `args` -- restated in the vector oracle (five refusals and a lawful tool node). The door judges the call when it runs (`councilEngine.Tool`): through `Call`, as the caller who fired or resumed the run; a tool that writes runs only when the nearest gate before the node grants it, whoever fired the run; a flow verb is refused (the run holds its world's flow lock); and the bare engine refuses a tool node by name.

`release_step` runs the core's `tests/cut.py` (`bump`, `index`, `check`, `ci`, `record`) in the world with the python the door runs the engine with, bounded at 45 minutes; a step that refused is an error, so its node fails and the release stops there. Held by `TestTheReleaseStepsRunFromTheCore`, `TestAToolNodesCallIsJudgedByTheDoor`, `TestAFlowsToolNodeCallsTheDoorItWasFiredThrough` (a reader answers the run; a writer past no granting gate fails it, fired from his glass) and the flow package's five tool strokes. The Workflows page fires a flow of tool nodes as it fires any other, and asks for every var the engine renders -- a tool node's `args` and a gate's `title` among them, which `openVars` did not read until the release flow's `mark` was found riding only there (held by `TestTheBuilderAsksForEveryVarTheEngineRenders`); its builder does not offer the kind yet.

## [v0.2.1] — 2026-10-05 15:16 (tag on 5c6beeb)

### The PowerShell tab: PowerShell 7 beside Bash, Python and Aider, every line on his card, and a shell line may run up to 300 s (2026-10-05, the core's WHAT'S LEFT H26; his word: "i should be able to run them in a tab alongside the agent, bash, python, and aider buttons or tabs. Let's get this done")

The door's `shell_run` takes `pwsh` (and `powershell`) beside `bash` and `python`. A PowerShell line goes through the bash gate for its refusals by name, and nothing the gate reads as a plain look lets it run unasked: the gate does not read PowerShell, so every line parks on his card. It runs `pwsh -NoProfile -NonInteractive`, the line carried in `ATLAS_SHELL_COMMAND` and read by `Invoke-Expression`, so the Windows command line never touches it, with the shell's built environment, output cap and secret scrub. The glass's fifth tab is PowerShell (`ps[world]>`), and every shell tab now asks the door for up to 300 s (`shellMaxWait`), where the page sent no time and every line got 30: the suites (about 175 s) and the `.ps1` scripts now fit.

Held by `TestAPowerShellLineAlwaysAsksFirst` (a cmdlet, a plain `ls` and a `git log` each park; `powershell` names the same shell) and by the glass's tab strokes (five tabs, each wired; `shellRun` sends `timeout: 300`). The refusal case is the bash gate's own, held there: a test line naming the secret file could not be typed through the glass, whose gate refuses it. Built in his Bash tab, placed on his card (the old builds kept as `.old.exe`), and run live on his glass: `7.6.6`, the ground's path, and a 35 s sleep that came back (35,347 ms). Not yet: a gate that reads PowerShell, so a plain look waits on the card where in Bash it runs at once.

## [v0.2.0] — 2026-10-05 13:39 (tag on e2db417)

### The Laws page: every law and how far the chain seals it, a law set and the ledger sealed from the glass (2026-10-05; his words, "set it as a law all the agents read ... THROUGH THE SYSTEM" and "a LAWS page ... Add it.")

Why: the laws live in the core's `law/` -- every `.md` there bound to `law/chain.jsonl` by `law/law.py`, and `LAW_LEDGER.md` the one law sealed as far as it is written,
which grows below its seal. Seeing how far each was sealed, setting a law or sealing the ledger meant typing at a terminal. On 2026-10-05 he set two laws (the core's
RULE 12, and the ledger's entries 6 and 7, SITTING LAWS 7 and 8, as drafts) and asked that it go through the system.

What, in the door (`line/internal/tools/law.go`, three tools):

- `law_status` reads: every law file in `law/` and how far the chain seals it (sealed, partly sealed, not on the chain, changed since sealed), read off `chain.jsonl` with
  law.py's own anchor; the bytes written below the ledger's seal; `law.py verify`'s own line as the verdict, where a clean exit that never says "proves whole" is not
  taken as whole; and why the two writers would refuse now, if they would.
- `law_add` writes one numbered entry to `law/LAW_LEDGER.md`, below its seal, in the ledger's own shape (an `entered:` and `from:` block, then the words). The file is
  opened for append only; the chain is walked before and after, and an entry the chain refuses after is taken back to the byte. Refused while a sitting or an engine is
  open on the world (SITTING LAW 5), over a chain that does not prove whole, for a line that would split the entry (`## ...` or a bare `---`), a control character, a
  title that is not one line of 3 to 120 characters, a law under 20 or over 8000 characters, and a ledger with a carriage return or no final line break.
- `law_seal` runs `law.py seal LAW_LEDGER.md --note ...` and nothing else: the chain is appended only through law.py. Refused when nothing is below the seal, a sitting is
  open, or the chain refuses.

Neither writer binds a law whole (`direct`), touches `CLAUDE.md`, or writes anywhere but the ledger and, through law.py, the chain; every act and every refusal is a line
in the holds record. Both are ServiceOnly: only his glass may call them, no role, seat or crossed gate's grant reaches them, and nothing parks them on a card -- his
click, behind the glass's own lock, is the decision. `docs/TOOL_PERMISSIONS.md` is regenerated (92 tools).

What, in the glass: a **Laws** line first under Record (a scales glyph) and `/laws` (`static/js/laws.js`). The chain card shows law.py's own verdict, green or red; a
card for the entries below the seal carries **Seal the ledger...**, which only asks (the question says how many bytes it binds and that a byte of it cannot change
after, with **Seal it** and **Not now**); a **Set a law** form takes its name, its words and where it came from; then `CLAUDE.md`'s standing rules in number order, the
ledger's entries each marked SEALED or DRAFT by where it stands against the seal's byte (counted in bytes, never characters), and every law file with its state. While
the door would refuse, both buttons are disabled with its reason and a card says it. Its reads are kept like What's left's, never quiet; it opens no pop-up and reaches
round the door for nothing.

**Proof.** In the door, 9 strokes (`line/internal/tools/law_test.go`) over a stand-in that speaks as law.py does; in the glass, 4 (`webapp/server/laws_test.go`), with
`TestNoPageOpensABrowserPopUp` and the permission record beside them. 22 reversals, each red and named, every file restored byte for byte. One stayed green at first,
the verdict taken from law.py's exit code alone, and `TestTheVerdictIsTheWalksOwnWordsNotItsExitCode` is the stroke written to turn it red. The door's reading of the
real law.py (its anchor, its verdict line, the chain's file and fields, the seal's words, an append below the seal and a seal over it on a copy of the real `law/`) is
held by the core's own suite: 15 checks and 12 reversals, in the core's CHANGELOG. A node run of the page with stubs, against the real `CLAUDE.md` and ledger: 25 checks.
And live, on a scratch pair in the app's browser pane (a door on 8190 and a glass on 8191 built from a synced copy, one throwaway ground holding copies of `CLAUDE.md`
and `law/`, a fresh wire): the page showed the chain proving whole at 6 links (head `07491469cd7d6d6c`), 14 standing rules, entries 1 to 5 sealed and 6 and 7 draft (2614
bytes), the ledger partly sealed and `SITTING_LAWS_2.md` not on the chain; a test law set from the form became entry 8 (344 bytes, LF, in the ledger's shape) with the
chain still whole; **Seal the ledger...** asked first, and **Seal it** sealed the copy to byte 18467 of 18467, the chain proving whole at 7 links; with a sitting
marked open the page said "Laws cannot be set or sealed now" in SITTING LAW 5's words, and a call made around the disabled button was refused by the door in the same
words. His ledger and chain did not move, and his door and glass were not touched.

**Placed on his card** (2026-10-05, 11:05, "Place and restart both"): the door and the glass were rebuilt from the tree, whose source was first checked byte for byte
against the copy that was proved (171 files), placed and restarted with the flags they had (door `e5713491...`, was `ba668ad1...`; glass `3db7cc11...`, was
`fa37925e...`; the old binaries kept outside the ground), and the page was opened on his own glass: the chain proving whole at 6 links, 14 standing rules, entries 6
and 7 drafts. **Not done.** The ledger's entries 6 and 7 stay drafts until he seals them, from this page or with `python law\law.py seal LAW_LEDGER.md`. The Records
page still marks the whole `law/` folder sealed (`records`); this page does not read that flag. **Restart:** done, on his card. **What goes red if unplugged:** the 9 door strokes, the 4 glass strokes, the permission
record, and the core's `test_the_door_reads_the_law_chain_as_law_py_writes_it`.

### The flow engine takes an `aider` node: its contract and its runner, and not yet the door's half (2026-10-05; the core's WHAT'S LEFT H17, first half)

Why: his ruling of B21 ("yes"): Aider makes the coder flow's code changes. H17 is that build; this is its engine half, and the piece is NOT finished.

What is built (`line/internal/flow`): `aider` is the eighth kind of the closed set. A node carries an instruction (`question`) and a file list (`files`, templated like a
question; the door judges every path). Refused at the save, never discovered when it runs: an aider node with no instruction or no files, `files` on any other kind, and
an aider node that can be reached without passing a gate whose `grants` name `aider_run` (`flow.AiderTool`) on every path; the nearest gate before it is the one that
counts, so a gate that grants nothing ends the hand's reach. The runner renders the instruction and the files and asks `Engine.Aider`, with the failed pass's words when
a check sent it back (a node that declares `loops`); an eval over an aider node scores the machine's own lines only (the evidence block), as over a `run` node; it is
retried like any node; a resume with no engine standing refuses at it; and the bare engine refuses it ("this flow has no door wired"). The Python cutter
(`tools/cut_flow_vectors.py`) learned the kind and its refusals, and the goldens were cut again: 8 kinds and 19 refusals, from 7 and 12, held by `TestFlowContract`.

**Not built: the door's half.** The door's council engine embeds the bare engine, so a flow fired through THE LINE refuses at an aider node with that sentence; nothing
yet asks the door's `aider_run` from a node, the coder tree is still its version 5, and the flow builder offers no `aider` kind. Those, and measuring the coder flow
with Aider on the six tasks of H12, are the rest of H17. Under RULE 11 the kind is declared and read by nothing on the door until then.

**Proof.** 8 strokes in `line/internal/flow/aider_test.go` and the regenerated contract; 16 reversals, each red by a failing test (four were first written so that they
only failed to compile, and were rewritten until they compiled and failed). The core's gate keeps its own copy of this law, and it did not know the kind until
the core's suite went red on the stroke that reconciles the two; the copy was restated the same day (the core's CHANGELOG). A door rebuild carries it, and nothing
on his door runs it yet. **What goes red if unplugged:** those 8 strokes, `TestFlowContract`, and in the core `test_the_flows_and_workflows_are_read_before_a_mark`.

### The front page plays what the maker made: the delivery's project in a sandboxed frame that shows the whole game (2026-10-05; the core's WHAT'S LEFT E1, first piece)

Why: the wife test. A person who is not the operator asked for a game and was handed a file path ("To try it, open projects\game-3\index.html in your browser"): the
new front page never read the delivery's `project` and had no way to play, so the game she asked for was on the old Dashboard's Projects card, another page
(finding 1 of `logs/pack_2026-10-03_181226.md`, run live on his glass on 2026-10-03).

What: a play frame between the terminal's output and its input line. When a turn ends and its delivery names a project (`project`, read off the maker and never
off a seat's words), the page asks the door's read-only `projects` tool which version that is (a quiet read) and frames the page the glass already serves at
`/api/projects/{name}/page?v=N`: sandboxed by the header it is served under and, as in the Dashboard, with no `sandbox` attribute (the app's own browser pane
refuses any frame that has one). A bar over it says the project and the version and offers "Open in its own tab" and Close.

Two things found by running it, each held now. (1) The terminal is rewritten whole on most events (`paintOut`), so a frame inside it would restart the game on every
line she types: the frame lives outside it, the same project and version are never loaded twice (a game she is playing keeps playing while she talks), a new
version or another project is loaded, and Close puts it away until a new version comes. (2) The game the maker made is an 800 by 600 canvas on a page with
`overflow: hidden` (`game-4`, version 4): in a small stage most of it was cut off, with no way to scroll to the rest. The frame is therefore given a window of 1000 by
720 and shrunk (never enlarged) to the room the stage has, centred in it and refitted when the room changes, so the whole game is on screen whatever the
size of the window; the stage is capped at 56 percent of the window and 640 px so the terminal keeps its place.

**Proof.** `TestTheFrontPagePlaysWhatTheMakerMade` pins the place (a sibling of the output, between it and the input line), that no frame is built in anything the
terminal rewrites, the frame's url, the missing `sandbox` attribute and referrer, that the project is read off the delivery and nowhere else, the quiet
`projects` read, that the same page is not loaded twice, Close, the fit and the height cap; 22 reversals, each red and named, every file restored byte for byte.
A node run of the logic with stubs (17 checks: one load per version, a repeat delivery leaves the frame alone, a new version reloads it, Close and staying closed and
coming back, the fit for a small, a big and an empty stage). And a scratch copy of the glass in the app's browser pane, against the real door: `game-4` (four
versions, made by the wife test) shown whole, shrunk to fit and centred, and Close putting it away and giving the terminal its area back. **Not shown here:** keys
and clicks through the scaled frame. The automation cannot deliver them into a cross-origin frame (it refuses to click a scaled one and sends keys to the top
page), so that rests on how browsers treat a scaled frame; the run on his own glass is below.

**Run live on his own glass** (2026-10-05, 08:04 to 08:07; sitting 342, opened and closed from the front page, three runs). The glass was placed on his card twice: the
first build (sha256 `8f8c44ae...`, replacing `0ff82efd...`) for the piece, the second (`fa37925e...`, glass pid 20416) for three fixes the first live run showed: when
the frame opened it took room from the terminal and the reply she had asked for scrolled out of view (the page now keeps the terminal's newest line in view when it
shows a page), the frame took 62 percent of the window and now takes 56, and two buttons on screen said Close (the frame's says Put away). Her three lines, typed as
she wrote them. (1) "If you can code, try making me a little game. something simple like a little sandbox miecrafty thing." The maker seated the Expert Coder alone and
made `game-5`, version 1, in 27.5 seconds; the delivery named it and the frame opened by itself on the front page with the game whole ("Arrow keys to move, Click to
start"), no file path to open (the reply's own words still tell her to open one: that is the engine's). (2) "i dont want it text based, i want a game i can play." was
answered "Nothing was saved: the page came back exactly as it was, so there is no new version to save." in 15.3 seconds, as on 10-03, and the frame stayed as it was: same
project, same version, so the game was not loaded again. (3) "thanks man, appreciate it" got a plain, friendly reply in 8.4 seconds ("I'm glad you enjoyed the game ...")
without the narration about "the operator" the same line got on 10-03 (one sample; the models vary). I did not play the game in the pane: the automation cannot deliver
keys or clicks into the frame. GitHub: core `53859d0` and atlas `8bff128` green for the piece.

**Not done, and said.** The engine's delivery text still tells her to open a file path, and the rest of what the live run found is the engine's (`manjuel/`), not
the page's: her second line's reply, small talk, "the game" meaning the oldest of that name, a blank canvas passing the maker's check. Nothing in `manjuel/` moved.
This is a glass change: it takes a rebuild and a restart of the glass, on his card. **What goes red if unplugged:** the stroke above; the Dashboard's own frame is untouched.

### The Linux leg went red on the Aider wall's probe test, and it was the test and not the wall (2026-10-05; follows "Aider behind the glass")

GitHub's `the go half, off Windows` (the portability probe: the first time these tests ran off Windows) failed on `cbec5eb`, and only on
`TestTheWallRefusesTheTwoKnownProbesQuietlyAndEverythingElseOutLoud`; the Windows battery, the gate, was green. The log needs a GitHub sign-in, so
the failing test was read from the run's ANNOTATIONS, which need none: the job's test step now writes the failing lines (`--- FAIL`, `FAIL`, `panic:`
and each `_test.go:line:` message) as annotations as well as in its log, the verdict still go test's own. The cause: the probe handed the wall
command STRINGS (`'git version'`, `'cmd /c "ver"'`), which Windows reads as command lines and every other system reads as the name of ONE program, so
on Linux the wall was shown the one-element list `['git version']`, did not match its two-word rule for the two known probes and said both out loud.
The wall was right and the probe was wrong. The probe now sends each command in the form its system uses: a string on Windows, as before, and the
same words as a list elsewhere, without the quote characters a list never carries. Proved here in the Windows form and, in a scratch copy only, with
the list form forced on, both through the real wall; GitHub's Linux leg is the third proof. I had vetted and compiled the Aider tests for Linux but
never run them there, and sent them: that is the fault. **What goes red if unplugged:** the test itself, and the annotation step (a red Linux leg
whose annotations name nothing).

### Aider behind the glass: the front page's Aider tab runs Aider on the files he names, behind a wall, on a line of work, and it was measured against the coder flow (2026-10-04; the core's WHAT'S LEFT H12 and B20)

His idea of 2026-10-02, as the list had it (WHAT'S LEFT H12): a page on the glass that drives Aider, open source and proven at search-and-replace edits. His answers of
2026-10-04, on the card: download it (`aider-chat`, from PyPI), put it in `aider/` at the ground's root, drive it with the coding seat's
own model on Ollama. And, in his words, "finish up h12 and b20, save and send, run it through the glass. also, finish wiring in the Aider
panel, it's on this system."

**What he gave (B20).** (1) The download: `aider-chat` 0.86.2 and the 110 packages it needs, about 640 MB, from PyPI, once, into a venv (a
Python 3.11: Aider does not run on the machine's own 3.14). (2) The folder: `aider/` at the ground's root, with `venv/` and `work/` (its home,
its temp folder and one folder per run). `.gitignore` names it, and was written BEFORE the folder was: `git_commit` stages everything, and one
save would otherwise have carried a venv. (3) The model: the coding seat's own, `qwen2.5-coder:14b`, read off `agents/expert_coder.md` at
every call (Model Target and Context, so a change he makes to the seat is the change Aider follows). No hosted route: RULE 4 as amended makes
that a measured, per-seat move of his, and nothing here needs one.

**Three door tools** (`line/internal/tools`: `aider.go` the tools, `aiderguard.go` the wall; the door serves 89 tools). `aider_run
{message, files, context?, timeout?}` and `aider_undo {run}` are `ServiceOnly` writers: no seat, agent or role is ever given Aider, and the
council cannot park a call of it for him to approve, exactly as with the shell. `aider_status` is a plain reader that the page asks quietly.
What a run is:

- **A copy, not the world.** He names the files Aider may change (at most 4) and optionally files it may only read (at most 4). The door copies
  exactly those into a scratch folder under `aider/work/runs/<id>/`, each in LF, and runs Aider there, headless, one message and out
  (`--message-file`), `--no-git`, the diff edit format, with every flag that stops it reaching for something it would be refused (no analytics,
  update check, URL scraping, shell suggestions, lint or test, a repository map; its model metadata is a file the door writes, so it asks no
  server). The scratch folder is its own empty git repository: Aider looks UP from where it stands for a repository even with `--no-git`, found
  the world's own `.git` on the first real run, and the wall rightly refused it (below).
- **The model is loaded first.** The door sends one empty request to Ollama at the seat's window (`num_ctx`) and `keep_alive` ten minutes before
  Aider is started. Measured on the first batch: a request that reached Ollama in the moment it was unloading the model hung four minutes inside
  Aider's own timeout until a request of that kind freed it. A rack that will not load the model is now a refusal that names it, and Aider is
  never started. The window is Aider's: a different one makes Ollama unload the model and load it again.
- **The window is the seat's.** The files and the instruction must fit the seat's declared Context (8,192 tokens: about 14 KB beside Aider's own
  prompt); a bigger set is refused up front, by size, rather than silently cut by the model's server. A single file over 64 KB is refused whatever
  the window.
- **Every gate the council's own tree doors keep.** A write lands only while the NEAREST repository the file belongs to (atlas/ carries its own)
  stands on a line of work: not main, not master, not detached, and a ground with no repository refuses. What a model may never write is the seats'
  own list restated (`aiderNeverWritten`: the law, the seats, the ledgers, the proof stamps, the governing files, a binary, a `.git`, `worlds/`,
  `projects/`, a secret, client material, a path outside the ground), and a new file lands in a folder that exists, never a new folder (RULE 8). The
  instruction is judged for a key, a secret file, client material and a path outside the ground before anything is touched. A file with MIXED line
  endings is refused; each file comes back in ITS OWN terminator whatever Aider wrote, and a new file takes its nearest sibling's.
- **What comes back is checked.** Only the files he named are looked at: anything else Aider made is reported and dropped, and a read-only file it
  changed is never written. A Python file it changed is parsed (never run) by the venv's own interpreter, and one that does not parse is NOT
  written. A file that changed while Aider worked is not overwritten. Every file lands or none does, with what was written put back if a write
  fails. The answer carries the diff (git's own, from the files as they were and as they stand), what Aider said, the tokens, and the run's id.
- **Undo.** The files as they were are kept beside the run; `aider_undo` puts them back exactly, removes what the run made, refuses a run twice,
  refuses a file that has changed since (it would lose what came after), and refuses on the main line. The last 30 runs are kept; an older
  one is git's to take back.
- **Recorded.** Every run, refusal and undo is a line in the world's `state/holds.jsonl`, the log `hold_answer` already writes, with the
  instruction (keys withheld).
- **One at a time**, because one model is in one memory; a second call while one runs is refused, not queued.

**The wall** (`aiderguard.go`). Measured on the first probe (Aider 0.86.2): `--yes-always` is not optional, since a run with nobody at the
keyboard would wait for a "y" for ever, and with it Aider will edit or create a file at ANY path the model names, outside its own root included.
So the door does not trust Aider's own confirmations. Aider runs inside a PEP 578 audit hook, the same wall the council's `run_python` child has
stood in since 2026-09-22, installed by source passed on the command line (`python -I -B -X utf8 -c <it>`) that nothing in the scratch folder can
edit: writes go only to the run's scratch folder, Aider's own home and temp, and the run's three history files; reads stay out of the ground and
the user's profile except `aider/` and the interpreter's own library; the network is loopback only (an address lookup or a connection to
anything else is refused by name, RULE 4); and no other process starts. Every refusal is said on stderr and lifted into the answer, so the card
shows what the wall stopped. And the wall ANNOUNCES ITSELF: once armed it writes a fresh random value the door put in the child's environment
(and the guard removed from it) to stderr, and a run whose stderr lacks that line is thrown away whole, however well it edited. The child's
environment is built from a list of names, never inherited: none of the door's keys, no service wire, no proxy, HOME and TEMP inside `aider/work`.
**Honest limits, the same as the precedent's:** an audit hook is not a kernel sandbox, and CPython's own documentation says so; native code
loaded through `ctypes` (which Aider's own terminal libraries import, so it cannot be refused) steps around it. The wall bounds an accident (a
model naming the wrong path, a library phoning home), not a hostile program. The scratch copy, the line of work, the suites and his click are what
stand behind it.

**What the first real runs taught** (Aider 0.86.2, the real model, through the door's own code, in a throwaway repository; each is a stroke now).
(1) The scratch folder, inside the world, let Aider find the world's own `.git` and its `.aider.conf.yml`, the wall refused the read, and the
first run died on an uncaught PermissionError: the scratch folder is its own repository now. (2) Every healthy run asks for two things the wall
refuses, `git version` (GitPython at import) and the Windows shell's `ver` (the platform module, three ways); a card that warned of the wall on
every run would teach him to read past the one where it matters, so those two are refused WITHOUT a word and everything else is refused
out loud. (3) Aider asks for its history files' folder, the run's own, on every start: asking for a folder that already exists is no change, and
is the operating system's to answer. (4) The Ollama hang above. (5) Aider on Windows writes CRLF; with `--line-endings lf` and the door putting
each file's own terminator back, a CRLF file comes back CRLF and an LF file LF. Aider is deterministic at temperature 0: two attempts at one task
came back byte for byte the same (4.1k tokens in, 417 out), so repeats teach nothing and the measurement below uses more tasks instead.

**The tab** (`js/agent.js`, `css/agent.css`). Wired: it was marked unwired and said so. A line typed in the Aider tab, or prefixed from any tab
(`aider /files`), goes to the door; it needs no engine. `/add path` names a file Aider may change, `/read path` one it may only read, `/drop`,
`/files`, `/status`, `/undo [run]` (the last by default), `/help`; anything else is an instruction. The page keeps only the NAMES of the files in
the chat and judges nothing about them: a path that may not be written, a main line, a model not on the rack are refused at the door, by name, and
the page draws the refusal ("REFUSED BY NAME", no button). A run's card shows: a chip for what happened (written, no change, withheld, timed out,
exit N), the model, the line, the tokens; the files with +/-; the diff, coloured; "THE WALL STOPPED" for anything the wall refused; "NOT WRITTEN"
with the reason when an edit was withheld; what Aider said (a toggle); what it made that was not kept; and an Undo button for a run that wrote and
for nothing else. While a run waits the card says so and its clock moves: at this machine's speed a run is minutes. The title bar says the
model, the line of work (or that it is on main and Aider cannot write), the files in the chat and whether a run is going. The Inspector's Aider
Pair tab now shows the same status, the files in the chat (each can be dropped), the recent runs (each can be taken back) and, as before, where each
world's work stands with its diff. Nothing on the page commits, runs the suites or lands (RULE 6): the edit is unsaved work on the line, and
Version control is where it is saved.

**Measured against the coder flow (H12).** The line asked whether Aider edits better than the council with the models this machine has, "measured, not
argued". Six small changes to one real file, `manjuel/mathkit.py` (226 lines, 8 KB: it fits the seat's window): add `correlation` after `covariance`;
give `mean` an optional `trim`; make `variance` refuse a negative `ddof`; add `zscores` after `stdev`; make `cosine` answer `None` for a zero vector;
add `clamp` after `cosine`. The scorer was written before any attempt was run: an attempt counts only if it changed the file, the file parses, the
behaviour asked for is right (a battery per task) and nothing else moved (every definition the task does not name is identical as a syntax tree, and a
differential battery over the file's public functions answers exactly as the original's). What a model says of its own work counts for nothing.
**Aider, through `aider_run` exactly as the glass calls it, one attempt per task: 6 of 6**, in 70, 100, 123, 140, 262 and 270 seconds (16 minutes in
all; at temperature 0 it answers the same twice, so a second attempt teaches nothing). **The council, the same words handed to a council turn (the coder
flow's `attempt` is one of these), in a throwaway copy of the ground on a line of work with the real engine, one attempt per task: 3 of 6**
(`correlation`, `mean`'s `trim` and `clamp`, in 417 and about 600 and about 600 seconds). The other three changed nothing: `variance` and `zscores`
never reached an edit (the answer is the file read back), and `cosine` asked `ground_edit` for an edit of the wrong shape ("needs exactly one `@@ OLD`
line and one `@@ NEW` line ... Found 0 and 0"), after which the closing words claimed the change in their first lines and described the refused edit
further down, and the engine's own "NOT EVERYTHING RAN" note under them says which is true. The council's turns took 466 to 648 seconds, 56 minutes in all.
**What the numbers do not say.** The council had one sample a task and one pass, where the flow gives `attempt` two. The scorer first counted those
three unchanged files as written, because the throwaway copy's file was CRLF (its git made it so) and the ground's is LF: found on reading the table,
corrected before anything was concluded, and no verdict moved. And Aider needs the whole file in a window of 8,192 tokens, where the council's edit
works on a named definition of a file of any size: of the 31 modules in the core's `manjuel/`, 12 fit (the largest, `skills.py`, is 249 KB), and 14
of the 26 files in this package. So on small whole files Aider was faster and right more often; on most of the ground it cannot be used with this
seat's window. Whether it becomes the flow's `attempt` is his (WHAT'S LEFT B21). The measuring scripts and the twelve results are in this session's
scratch folder, not in the ground.

**What goes red if unplugged** (RULE 11). The door: `TestAiderIsHisAloneByEveryRoute` and `TestEveryServiceOnlyToolRefusesEveryOtherCaller` (no
other hand is given Aider), `TestAiderWritesOnlyOnALineOfWork`, `TestAiderRefusesWhatTheSeatsAreRefused` and the core's
`test_the_aider_door_keeps_what_the_seats_keep` (the seats' lists and the door's copy are equal from both sides), `TestAiderRefusesAKeyAndRecordsEverythingWithoutOne`,
`TestAiderWritesEachFileInItsOwnTerminator`, `TestAiderOnlyWritesBackTheFilesHeNamed`, `TestAiderIsBelievedOnlyBehindAWallThatSaidSo`,
`TestAiderWritesNothingFromARunThatDidNotFinish`, `TestAiderWithholdsAnEditToAFileThatMovedWhileItWorked`,
`TestAiderWithholdsAPythonFileThatDoesNotParse`, `TestAiderUndoPutsBackExactlyAndRefusesWhatMovedSince`, `TestAiderRunsOneAtATime`,
`TestAiderFollowsTheCodingSeatAndRefusesWhatCannotFit`, `TestAiderLoadsTheModelAtTheSeatsWindowBeforeItStartsAider`,
`TestAiderScratchFolderIsItsOwnRepository`, `TestAiderIsStartedWithTheFlagsThatKeepItOnTheLeash`, `TestAiderGetsABuiltEnvironmentWithNoKeyInIt`,
`TestAiderWaitIsBoundedAndHasADefault`, `TestAiderKeepsTheLastThirtyRunsAndNothingElse`, `TestAiderSaysWhyItCannotRun`; the wall, against a real
Python behind the real guard: `TestTheWallShutsWhatItNamesAndLeavesOpenWhatAiderNeeds` (26 attempts, 28 where the machine has an address of its own,
to write, make a folder, read, list, remove, rename, copy, connect, look up, bind and start, each judged),
`TestTheWallRefusesTheTwoKnownProbesQuietlyAndEverythingElseOutLoud`, `TestTheWallDoesNotStartWithoutItsEnvironment`,
`TestTheWallsAnnouncementNamesTheRunsOwnNonce`, `TestTheGuardSourceIsSafeToCarryOnAWindowsCommandLine` (a double quote in the guard is read by two
sets of rules on a Windows command line: a convention doing a type's job, held now); the glass: `TestTheAiderTabIsTheDoorsAiderAndNothingElse` (the
keys the page reads off the door's documents are keys the door has, the Undo button is drawn for a run that wrote and nothing else, the page asks for
three tools and never one that saves, lands or answers a hold, holds no list of what may be written, and `/drop` says what it dropped), with `TestTheFrontPageFakesNothing`,
`TestNoMethodOfTheFrontPageIsDefinedTwice` and `TestEveryToolTheFrontPageAsksForIsCarriedByTheDoor` carrying the tab. **Reversals, each red and named,
every file restored byte for byte: 111 (the door and the wall 67, the page 36, the core's stroke 8).** The scorer for the measurement was proved on six correct solutions and fourteen plausible wrong ones
before it was trusted with a model's work.

**Not done, and said.** Aider is not the coder flow's `attempt` engine: that is his decision, with the measurement above in front of him (WHAT'S
LEFT B21). No seat, agent or hosted route is given Aider. The suites and the Land click are not touched. The door and the glass are placed (below).

**Placed.** On his card ("Place and restart both") a build of the tree replaced the live door and glass, and both were started again with the start
lines they had: door pid 21608, 12,801,536 bytes, sha256 `ba668ad1...` (the one it replaced: `e4392bb3...`); glass pid 25236, 11,177,472 bytes, sha256
`0ff82efd...` (replacing `8ddafc33...`). Two builds of the tree came out byte for byte the same, and a third, taken just before the copy, matched
them. The old binaries are kept outside the ground. No sitting was open (341 closed 2026-10-03 18:11). The glass came back on its PIN lock, and the PIN
is his.

**Run live on his own glass** (2026-10-04, 21:11 to 21:15, after the first save of both repositories, atlas `406e5e3` and core `f1fd61e`: the door's
`git_branch switch` refuses a tree with unsaved work and Aider needs a line of work, so a live run could only follow a save). On his page the Aider tab read
"Aider 0.86.2 is ready: qwen2.5-coder:14b, a window of 8192 tokens (about 14 KB of files at a time)". An instruction on `main` was refused by name. A line of
work, `aider-proof`, was opened through the call the Version control page's button makes. The `clamp` task of the measurement, handed to the real
`manjuel/mathkit.py`, was written in 64 seconds (4.1k tokens sent, 339 received, 6 lines added), and the measurement's own scorer, run on it against the
committed file, says wrote, parses, correct and clean; it is byte for byte the file the measurement's Aider run made, so the earlier OK stands for it too.
The nested `atlas/` repository, standing on `main` while the core was on its line, refused an edit by name. A new file, `aider_proof.py`, was made in 9.8
seconds and parses. Both were taken back with the Undo button: `manjuel/mathkit.py` hashed to the committed blob again (`9534e60a...`), the new file was
gone and the tree was clean. The world's `state/holds.jsonl` holds the six lines (refused, wrote, undone, refused, wrote, undone), and the Inspector's Aider
Pair read the status, the two runs "taken back", atlas on `main` with nothing changed and the core on its line with nothing changed. The core went back to
`main`, the line was closed, and no sitting was opened.

### The front page's window keeps the size it starts at, and the terminal scrolls inside it (2026-10-03; the core's WHAT'S LEFT C37)

His word, 2026-10-03, watching the app's browser pane while a session ran: "the window just keep growing, that doesnt make sense. the
window should be locked to that scale it starts at and just scroll the terminal instead of the whole page."

**The cause.** `agent.css` gives the terminal window a fixed height (`height: calc(100vh - 125px)`) and its output box
`overflow-y: auto`, which works on a wide screen. Below 1024 px wide -- and the app's pane is 888 -- the same file turns the main
area into a column, and in a column the window's `flex: 1` sizes its HEIGHT: by its content, overriding the fixed height. The window
grew with every line, the page with it, and the output box was never short enough to scroll. Measured on his live page: the document
3696 px tall at 888 px wide (4042 by the end of a session), the window 3578 against the 592 it is meant to be.

**The fix** is one rule in the narrow-screen block, with the reason beside it: `.ag-win { flex: none; }`.

**Proved in a browser, not only in the file**, on a scratch glass at the pane's own 888 x 717, with the same 150-line flood before and
after. Before: the document 15,186 px, the window 15,068, the output box 14,960 px tall with nothing to scroll. After: the document
717 px (the viewport: the page does not scroll), the window 592, the output box 485 px showing the end of 14,960 px of content and
scrolling inside itself. And at 1280 x 800 (the side-by-side layout, unchanged: the window 675 = the viewport less 125, the output box
scrolling) and at 600 x 800 (stacked, no horizontal scroll). The glass's whole Go suite is green and `go vet` is quiet.

**What goes red if unplugged:** `TestTheWindowIsNotSizedByItsContentOnANarrowScreen` -- the window and the panel keep their given
height, the narrow block tells the window not to flex, the output box scrolls -- with three reversals, each red and restored byte for
byte. A rule read as text can be in the file and still lose, so the test is the wire and the browser measurement is the proof.

**Not done, and said.** With the Inspector open on a screen narrower than 1024 px the panel still STACKS below the window (a second
full-height box), so the page scrolls to reach it; that is the layout's old design and was not asked about. **Placed** on his card
("Place and restart the glass"): the proved build over `atlas-webapp.exe` (a fresh build of the tree is byte for byte the same; the old one is
kept outside the ground), only the glass restarted (pid 18916; the door untouched), and measured on his own page at his own 888 x 717 after he
unlocked it: the served stylesheet has the rule; the page 717 px, the window 592, `flex: 0 0 auto`; and after a 60-line flood the page still 717
and the window still 592, with the output box holding 4,925 px and scrolling inside its 485.

### The operator's typed shell: Bash and Python on the front page, gated, and the rack's portless host mended (2026-10-03; the core's WHAT'S LEFT H15 and C36)

The operator's word of 2026-10-02, on the card that ruled it, was that Bash and Python are "typed by you, gated". Asked
today what the Bash tab should run and how the Python tab should behave, he chose Git Bash, and a Python session that
remembers its names (the "python loop that works" he said he struggled to get before). And, mid-turn, about the rack:
"the ollama host is always the same served from the desktop app, so whatever we need to do to rectify that, get it
while you are at it."

**The rack (C36).** This machine sets `OLLAMA_HOST=127.0.0.1`. `rack.Host()` turned the bare host into
`http://127.0.0.1`, which dials port 80, so `rack_list` called a standing Ollama silent and the front page said "rack
silent". A host with no port now means Ollama's own, 11434 (`rack.DefaultPort`), in the one function every caller goes
through (the rack tools, `chat`, `play`); a port that is named is never changed, and an outward host is still refused
once it has gained one. The core's engine never read the variable (it binds the rack on loopback itself), so the door
was the only place that disagreed. Live on the placed door: "rack 11 voices".

**`shell_run`, a tool of the door that only his glass may call** (`line/internal/tools`: `shell.go` the tool,
`shellgate.go` the classifier, `shellpy.go` the Python session). `shell_run {shell: bash|python, command, timeout?,
reset?}`. What you type runs as you, in the world the page is on, and is judged by reading it, never by asking a model:

- A plain look runs at once: list, cat, head, tail, wc, grep, rg, find without `-delete` or `-exec`, echo, `git status`,
  `git log`, `git diff` and the other read-only git verbs, a pipe of such, and, in Python, a calculation over numbers,
  text and containers of them.
- Anything else shows a card in the thread and waits for his click: it writes, deletes, installs, reaches the network,
  commits, pushes, runs a program the gate cannot read, or uses shell syntax it does not model (a substitution, a
  variable, a brace, a background job). The default is the card: the read list is short and exact, and a command the
  gate has never heard of is asked about, not run. The card is the hold queue the door already has (`hold_list`,
  `hold_answer`, the Guardrails tab), so approving runs exactly what was parked, and re-judges it first.
- A secret file (`.env` and the rest of `is_secret`'s names), client material, a key typed into a command, a home
  directory or a path outside the ground is refused by name, and no click lifts a refusal. A name dressed in quotes or
  backslashes is the name it spells, and a glob that can reach a dotfile is refused.

**What holds it shut.** No agent is ever given the shell: `Tool.ServiceOnly` is a new declaration, enforced in
`Registry.Call` before RBAC or the holds look at it, so a role cannot grant the tool and an agent cannot park a call of
it for him to approve. The attempt is written down; the tool refuses for itself too if it is reached any other way; and
the permission record (`docs/TOOL_PERMISSIONS.md`, now 86 tools) prints the declaration. Approval is a typed field of
the caller (`Caller.Hold`, set by `hold_answer` alone), not a word in a name. The law is walked before every run
(`law/law.py verify`, cached against what `law/` holds, as the engine's own gate does). The child gets a BUILT
environment with none of the door's keys in it (the door holds the service wire), and what it prints is scrubbed of
every secret value in the ground's `.env` and of every variable of the door's own whose name says it is one. A run ends
at its limit (30 s, 300 at most), with its children (taskkill's tree on Windows), keeps the first 64 KB of what it
printed, and an entry is at most 8,000 characters. Every run, hold and refusal is a line in the world's
`state/holds.jsonl` with the command (keys withheld): the log `hold_answer` already writes, not a new one. `spawn`
gained four optional fields for this (`CleanEnv`, `MaxBytes`, `Grace`, `KillTree`); every other caller is as it was.

**Bash** is Git Bash, found by the `git` on the PATH (the bare word `bash` on this machine is the WSL stub in
WindowsApps, which the door never uses), and an entry reaches it in an environment variable that `eval` reads, never
as an argument of `bash -c` (below). **Python** is one process per world, started from the ground with the python
the door runs the engine with, speaking JSON lines over a private copy of its pipe (a child process's output lands on a
separate channel and cannot pose as an answer; `input()` reads nothing). The driver judges the entry with `ast`, never
`exec`: names are trusted for what they hold now, read off the live session, so a name an approved cell bound to a file
does not make `x.read()` look like arithmetic. An entry that runs past its limit ends the session and its names go, and
the next answer says so; `Reset` or `/reset` ends it by hand.

**The page** (`js/agent.js`): the Bash and Python tabs are wired; Aider stays marked unwired. A line typed in a shell
tab, or prefixed from any tab (`$ ls`, `! ls`, `>>> 1 + 1`), goes to `shell_run`; a shell needs no engine and does not
wait behind the council. The page decides nothing about a command: it draws the one document the door answers in
(`state`, `class`, `why`, `hold`, `exit`, `output`...), a card for a held answer and for nothing else, "REFUSED BY NAME"
with no button for a refusal, a bad chip and box for a non-zero exit, "ENDED" for a run cut short, and "answered
elsewhere" for a card whose hold was answered in Guardrails. Each tab has its own prompt, placeholder and things to try.

**Proved.** 309 commands, each said 11 ways (padded, tabbed, behind a comment, chained behind a harmless look, piped,
stderr folded in), 3,399 judgements, none allowed to differ from the plain phrasing; 87 Python entries judged by the
real session; the tool run for real against Git Bash and Python in a temp ground (no agent is handed it, a name is not
an approval, a write waits and approving runs exactly what was parked, a refusal is refused again at the replay, the
child inherits no key, a secret value reaches no page, a broken law refuses the run, every outcome is recorded, a
runaway ends); a differential test of the gate's lexer against the real Git Bash through the real route (33
quotings); and 45 reversals (4 on the rack, 28 on the door, 13 on the page), each red and named, every file
restored byte for byte. The whole door and glass suites green, `go vet` clean. Then the built pair was run for real: a
scratch door and glass built from this tree (their own wire and lock, a throwaway ground), driven in the browser with
real clicks and typing. Bash: a listing, a log, the three refusals, a recursive grep whose secret came back withheld, a
write that showed a card and ran on a click, a delete that was denied and left the file. Python: `x = 41`, `x + 1` =
42, an import that asked and ran, `/reset` and the NameError after it, the `$ ` prefix routed from the Python tab, the
held call listed in Guardrails. The ground's holds log held all of it.

**Found by that run, and mended.** I named a new method `shell`, which silently REPLACED the page's own `shell()` (the
template it is drawn from; the later key of an object literal wins), and the page drew "[object Promise]". Every stroke
that reads the source as text had stayed green. The method is `shellRun` now, and
`TestNoMethodOfTheFrontPageIsDefinedTwice` says so for every method of the page. The run also showed a refusal said
twice, a denied card that wore its waiting time as its run time, and a card that kept its buttons after Guardrails
answered it; all three are mended and held.

**A hole in the door as first placed, found an hour later and mended.** To prove the gate and bash see the same words I
wrote the dull test: split a list of awkward quotings both ways and compare. Handed to `bash -c` on the Windows command
line, a doubled backslash outside double quotes was HALVED by Git Bash's runtime: `cat .\\env` was lexed as the file
`.\env` (not a secret's name) and run as `cat .\env`, which bash reads as `.env`. The name refusal for `.env` could be
walked round for as long as that door stood (about an hour); the secret VALUES were still scrubbed from what it printed,
so nothing showed. The entry travels in `ATLAS_SHELL_COMMAND` now and runs under `eval`, which is what the lexer
models, and `TestTheLexerSeesTheWordsBashSees` (33 quotings through the real route) and
`TestADoubledBackslashIsNotASecretName` hold it. The mistake was mine: I had assumed a command line is a faithful
carrier for a command, and no test of the gate had ever compared its words with bash's.

**Placed.** On the operator's card of 2026-10-03 ("Place and restart, hold the save and send") the door (12,340,736 to
12,595,200 bytes, sha256 `7f36ab1e414c...`) and the glass (11,140,096 to 11,152,384 bytes, `53c89db99d1d...`) were
placed; a fresh build of this tree comes out byte for byte the same. Both were restarted (door pid 17064, glass pid
22816; the old binaries are kept outside the ground). The door was down about two minutes: PowerShell split the quoted
`--manjuel "python ..."` value of its start line into two arguments and the door refused it by name; it was started
again with the value quoted as the runbook has it. Checked on his own page, read-only: the status line says "rack 11
voices", `ls` and `git log` read his ground, `env` showed a card that was denied (nothing ran), and `1 + 1`,
`sum(range(101))` and `/reset` ran in Python.

**Placed again.** On a second card the same hour ("Place the fixed door and restart it") the fixed door (12,595,200 to
12,596,224 bytes, sha256 `e4392bb39683...`; a fresh build of the tree is byte for byte the same) was placed and
restarted (door pid 9248; the glass is unchanged, pid 22816). The fixed route was run for real on the scratch pair:
quotings, pipes, `cat .\env` refused by name, `cat .\\env` finding no such file. The last look at the fixed door on
his own page did not happen: his glass had gone to its lock screen, and the PIN is his.

**Not done, and said.** Python's `def`, `class` and a call to a function you defined are always a card (nothing
remembers which of your functions are pure yet). `/reset` is not written to the holds log. The shell tabs keep no
scrollback across a page reload (the record is the holds log). The Aider tab (B20, H12) and retiring the sidebar (H16)
are the next pieces. Saved and sent on his word ("go"), after the ground's suites ran green from his glass.

### The front page: one terminal, five doors and an Inspector, in the glass's own JavaScript (2026-10-02; the core's WHAT'S LEFT H14)

The operator, shown a page out of his own AI Studio project: "that's what I am looking for", and, asked where it
should go, "replace the Dashboard now". `/` is now that page, written in this console's plain JavaScript and CSS
(`js/agent.js`, `css/agent.css`): a top bar with five doors (Aider Pair, Audit Ledger, Workflows, Registry & Docs,
Guardrails), a status line, a Pages menu, an Inspector toggle and who is signed in; ONE terminal with four mode tabs
(Agent, Bash, Python, Aider), a prompt line with history, the microphone and suggestions; and an Inspector that opens
beside it with seven tabs (Run, Aider Pair, Audit Ledger, Workflows, Guardrails, Docs & Registry, Rack). The sidebar
is hidden on this page and on no other. The launchpad it replaced is kept whole at `/dashboard`, with its own
sidebar link, and every other page is where it was.

The Agent tab is the council: a line typed there is the REPL's own turn through `Run`, so the law gate stamps it before
any model reads a word, and the page writes its turns into `Chat.thread`, which the Watchboard and the old Dashboard
read too. Failures are on the face of the answer ("NOT EVERYTHING RAN"); a question from the council is a card in the
thread with a field under it (a blank answer is an answer, as the door already treats it, and while a question is open
the prompt line takes nothing else); a turn can be cancelled; a typed `/command` prints into the terminal; and
`/boot`, `/close` and `/clear` are the page's own. Boot opens the engine, loads the pipeline's models and keeps the
engine's own `/status` report. The sidebar's button and the palette boot through the same functions.

What it did not take from the page it came from. That page showed a tokens-per-second figure drawn from a random
number, a "WAL Active" that was a label, an agent that matched words and reported tool calls that never ran, a pair
programmer that appended a stub, committed it and said done, a repair loop that reported a repair that never
happened, a seeded approval queue that gated nothing, and a login with a token the page invented. None of it is here.
Every figure on this page is read from a tool or not shown. Bash, Python and Aider are in the bar and, pressed, say
they are not wired yet (each is its own piece, behind the door's gate; the operator ruled "typed by you, gated" for
Bash and Python). The Audit Ledger shows the engine's own `/status` report, the sittings off `proofs` and the
transcripts off `records`, and does not call `verify_chain`, which cannot yet read the pen's links (the core's H7).
Guardrails is the door's real hold queue, with an Approve and a Deny that run only on his click. The prototype is
unpacked and installed, untracked, at `worlds/agentos` in the core's ground, for reference only; nothing of its code
is here.

The door's `hold_list` joins the glass's quiet reads (`backgroundReads`): the Inspector's dot asks how many writing
calls wait when the page opens and after each turn, and a "New trace recorded" toast at every one of those would be
the noise the 2026-09-16 ruling ended. The door declares it `Writes: false`.

**Proved.** Three new tests in `server/glass_test.go`. The front page is routed at `/` with the launchpad kept at
`/dashboard`, and boot and close are its own. It fakes nothing: no random figure, no sign-in of its own, the unwired
tabs marked and saying so, the failures drawn on the answer, the gate answered only from the field and the hold queue
only from a button. And every tool it asks the door for is in the door's tool table, every quiet read is on the
glass's list, and every route it fetches is served. Fifteen reversals, each red and named, restored byte for byte;
one of them stayed green and showed a real gap in the first stroke, which looked for the failure face anywhere in the
file, and it now looks in the terminal's own. The whole glass suite green (`db`, `handlers`, `server`), `go vet`
clean. Then the built page was run for real: a copy of the glass built from this tree, on port 8191 with its own data
and lock, against the live door. From the page it booted an engine (sitting 334; `/warm`, then `/status`, whose memory
line read "chain whole"), ran one real council turn (310 s, two seats, four tools, delivered, the engine's own "READ IN
PART, NOT WHOLE" line on the answer), ran a typed `/command`, made the council stop on its own question and answered
it in the thread (a `/remember` declined at its last question; `memory.md` was not touched), cancelled a turn, read
every Inspector tab off live data, and closed the sitting (tolled). That run found, and this code fixes, a top bar
that made the page scroll sideways under about 1,200 px, a Seats stage that counted sittings and not seats, a gate
that refused a blank answer, and an answered question that printed nothing.

**Placed.** On the operator's word of 2026-10-03 ("go", with his sitting 335 open) the rebuilt glass was placed
(11,140,096 bytes, sha256 `d7e29c67d1d0...`, the one the tests and the live run proved; the old binary is kept
outside the ground) and the live glass on 8091 was restarted. It serves exactly the files in this tree, byte for
byte, and his first tab shows the front page signed in, with the open sitting reflected ("engine open", "sitting
335"). The glass embeds its pages, so any later change to them needs the same: a rebuild, a placement and a restart.

**Not done, and said.** Bash and Python on the gate, the Aider tab (which waits on B20) and retiring the sidebar are
the next pieces (the core's H15, H12 and H16). The status line quotes `rack_list`, which reads a portless
`OLLAMA_HOST` as port 80 and so reports the rack silent while Ollama stands (the core's C36).

### The door's tools have permission records: read off the table and the policy, held current by a test (2026-10-02; the core's WHAT'S LEFT D7)

The tool server's tools declared what they are (writes or reads, which actions of a writing tool only
read, which arguments are secrets, what they need beyond a directory) and the shipped role policy decided
who may call what, but nothing set the two side by side where the operator could read them in one
sitting: "the door's 85 tools have no permission records" was the line (HANDOFF 09-28).
`docs/TOOL_PERMISSIONS.md` is the record, and it is READ OFF the tool table and the policy by the
functions that judge a real call (`rbac.Can`, `Tool.Writes`, `Tool.Reads`, `HeldExempt`), not written
beside them, so it cannot say what the door does not do: per tool, its tier, whether it writes, its
reading actions, its secret arguments, whether a call from anything but the glass is held, and what each
shipped role (operator, steward, agent, guest) is told -- yes, reads only, or refused. 85 tools: 48 read,
37 write (two of the writers, `git_branch` and `git_tag`, have reading actions the agent role may call);
tiers core 76, spine 1, engine 8; 4 carry a secret argument. A tool added, a flag turned or a role
changed makes the record stale, and a stale copy fails `TestEveryToolHasAPermissionRecord`: the change
shows up in review, before it reaches a door. Regenerate with
`go test ./internal/tools -run TestEveryToolHasAPermissionRecord -update`.

What it does not claim, said in the file itself: a tool's `Writes` is its author's declaration, and
nothing here proves a tool that says it only reads never writes -- that is each tool's own test. The
record makes a change of declaration visible; it does not make a false one true.

**Proved.** Four tests in `internal/tools/permissions_test.go`: the record is current; it has a column
for every role the policy ships and for none it does not; it has exactly one row per tool; and the shipped
roles answer as the record says (the operator is refused nothing, the guest is allowed nothing, the
agent is refused every tool that only writes). Six reversals on a scratch copy, each red and named,
restored byte for byte: a tool's declaration flipped, the agent role allowed to edit, a role added to
the policy and not to the record, the exempt list grown, a row deleted, a cell edited by hand. The
package green on the ground (46 s). Go's test cache does not see a file outside the module, so a hand
edit of the record can be answered from the cache by a repeat of the same command; run it with
`-count=1`, which `tests/prove.py` does and GitHub, which has no cache, does by construction (said in
the test's header).

No binary moved: a test file and a document. **What goes red if unplugged:** those four tests, and
`tests/prove.py`'s `go test ./... -count=1` on every push.

### The chain verdict in Python, and `verify_chain` rewired to it: the first half of the fold (2026-10-01; his ruling of the same morning, the core's WHAT'S LEFT B14)

B14 read "retire the unused Rust". Measured before anything moved, the Rust was not unused: it was
the door's spine for one tool, `verify_chain` (the chain verdict EMPTY | INTACT | FLIP | TAMPER),
and the canon that `cargo test` proves over the 21 golden-master chains -- and that spine is
itself a byte-faithful mirror of a Python oracle outside this ground. His ruling: the verdict in
Python first, proved against the goldens, then the Rust folded to an attic. This is the first half.

`tools/chain_verify.py` (stdlib only; reads a chain, writes nothing) is the verdict written from
the Rust in the ground, line for line: an entry's hashed view (ts, kind, n, payload, prev, actor,
body_v, where present); its hash sha256(prev + canon(view)) under the entry's own `body_v` (1 the
board shape, 2 the links shape, 3 JCS -- UTF-16 key order, no whitespace, floats and out-of-range
integers refused; 3 when unstamped); the weld following each entry's STORED hash; a broken weld
TAMPER at its index; a hash that does not recompute while the weld holds a named FLIP with the
chain still appendable; an empty file EMPTY. One place the Rust and the oracle part, and the
oracle is followed: a chain whose rows carry no `hash` at all is SKIPPED with the reason, not
TAMPER at entry 1 -- the goldens say skipped. Its words: `chain verify <path>` (the one line the
Rust printed, so no caller reads anything new; exit 0 while appendable), `--json <path>` (the
oracle's dict), `--goldens` (every chain and injection in `tests/fixtures/canon/chain_verdicts.json`
re-verdicted here and compared).

The door's `verify_chain` now spawns that script through the one spawn contract -- the engine's
own interpreter (the first word of the core command when it is a python, else `python`), bounded
at 30 s, the words kept whole on error -- and finds it the way `findAtlas` walked for the binary:
up from the door's own executable, from the working directory, and from the tenant's home (whose
`atlas/` holds the repository when the tenant is the core ground), never on faith from a flag. No
tool in the door's table shells the Rust binary any longer (`atlas-door`, a separate binary, still
does, for `atlas trade`). `TestAToolsOwnWordsSurviveItsError` no longer skips when the spine is
not built: there is nothing to build.

**Proved, on the ground.** `--goldens`: 29 of 29 verdicts match (21 chains, 8 injections). The
canon over `tests/fixtures/canon/vectors.json`: 356 of 356 (250 outputs byte-equal, 106 refused
where the vector refuses). Against the Rust CLI's own line over the 21 goldens: 18 identical, the
other 3 SKIPPED here where the oracle recorded them skipped, 0 differ.
`TestVerifyChainReadsAGoldenThroughThePythonSpine`: a golden verifies INTACT in the spine's shape;
a file that is not a chain errs with `verdict=TAMPER ... broke_at=Some(0)` in its own words; the
tool's description no longer says Rust. `tests/prove.py` gains the leg `chain_verify (python
spine)`, never ABSENT (the goldens are committed): 23 held, 14 absent, 0 broke. The verifier is not
a cutter -- it takes no cutter's word and writes nothing -- so `cut_words` leaves it alone.

**What goes red if unplugged:** the Go test, the proof's new leg, and `verify_chain` itself,
which errs by name when the script is not beside the door.

**Left behind by the rewire, named and not fixed (2026-10-02, read before the fold was begun):** the
door's `--atlas-bin` flag, `Options.AtlasBin` and `findAtlas` are read by nothing now -- declared,
correct, and loose; and `tests/e2e/_start_mcp.ps1` still refuses to start the door unless the Rust is
built, for a tool that no longer needs it. The Rust stays (ruled 2026-10-02, B14 closed), so both are
harmless dead code. The comment on `TierSpine` and the CLI reference's line for `verify_chain`, which
said "Rust", are corrected.

**The second half is a decision, not a build.** Reading what the fold touches found what the first
measurement missed: `line/cmd/atlas-door` -- the D2 trade-ops page, the Go cutover of the Python
`door.py` (docs/rollback/D2_DOOR.md) -- reads and writes through `atlas trade`, the Rust `store`
crate, and every atlas release ships `atlas.exe` and `atlas-door.exe`. With the verdict in Python
nothing in the tool server needs the Rust; those two shipped products do. Put to him 2026-10-02.
**Ruled 2026-10-02, his word: "keep the Rust, close B14".** The Rust stays; nothing is folded.

### `land`: a line of work landed onto the main line from Version control (2026-10-01; his ruling of 2026-09-30, the core's WHAT'S LEFT B16)

The table carries no `merge` by its founding law and never will; landing a line of work was his
terminal's. `git_branch` gains the action `land`, narrow on purpose: fast-forward only (the main
line has not moved since the line began), from the main line (you stand on it), over saved work
(the tree is clean); a line that does not exist, the main line onto itself, and a join with two
parents are each refused by name, and the last names whose act it is ("your terminal's, not this
button's"). Landing does not close the line -- `close` is its own act and refuses unmerged work
by itself. The Lines-of-work box on Version control offers "Land onto main" for a line that is
not the main line, only when the main line is the one you stand on.

**Proved.** `TestLandingALineIsAFastForwardOntoTheMainLine` (refused from the line itself, the
main line onto itself, an absent line, a dirty tree; then landed, main at the line's head, the
line still there) and `TestLandingRefusesWhenTheMainLineHasMovedOn` (refused, main unmoved);
`TestVersionControlOffersLandingALineOntoMain` reads the page's source for the button, its gate
and the absence of any forbidden verb. `go test ./internal/tools/` and the glass's packages green
on a scratch copy; the door and the glass rebuilt from it and PLACED on his allowance at 08:58 --
the door pid 108440, the glass pid 35444, both hashing as built, the live check 9/9 through them.

**What goes red if unplugged:** the three tests; the door's refusal of an unknown action, which
names `land` now.

### `seed_catalog.py --verify` checks only (2026-09-30, the core's WHAT'S LEFT C35)

It said "check only" and created `data/master.db` when the file was absent, seeding it before it
looked -- so a verify on a fresh clone wrote, and prove.py runs it on every proof. A verify with
no database says so and exits 1 with nothing written; `--verify --reset` is refused; the counts a
verify reports are read, not seeded (`tally`). The seed itself is unchanged.

**Proved.** On a scratch copy: verify with the database present changes no byte; without it,
"absent", exit 1, no `data/`; `--verify --reset` refused, exit 2; the plain seed still seeds. THE
WIRE: `tools/cut_words.py --verify` grew a leg that runs the script in a bare copy and wants
exactly that -- with the old script back it goes red ("exit 0; WROTE data/"). On the ground the
leg passes and `seed_catalog.py --verify` answers VERIFY OK with `data/master.db` untouched.

**What goes red if unplugged:** that leg of `cut_words --verify`, here and on GitHub's push.

### The glass's batch: nine lines of the list, one rebuild (operator, 2026-09-29: "do the list top to bottom")

WHAT'S LEFT C17, C23, C25, C26, C27, C28, C29, D4 and D5 -- all in the glass but for the TUI's
half of C29. Nothing in the door moved.

- **The idle line says what idle costs (C17).** The hero's clock said "idle 14m" and not that
  the core closes a sitting nobody has used for thirty minutes between turns (serve.py's
  IDLE_CLOSE); it says so now, with the minutes left. The thirty is `Home.IDLE_CLOSE_MIN`, a
  copy of the core's constant named as one -- the wire does not carry it.
- **`Run.check` tells a silent door from a world with no engine (C23).** The glass answers 502
  with a JSON body when the door does not answer; the body parsed, `open` was absent, and the
  hero said "No engine" -- boot one -- over a door that was down. Anything but 200 is the door
  silent. **And the store hands back a copy (C23):** `db.GetAgent` returned a pointer into its
  own slice, read after the lock was gone while `UpsertAgent` may overwrite that element or
  grow the slice under it; `GetTrace` already copied for the same reason.
- **A paused run is found after a reload (C25).** It lived only in the tab's memory; the pause
  sits in `flows/runs.jsonl` for as long as he likes (the spec's §4.9). The runs list is read
  from the record on arrival -- every flow's, with none open -- and each PAUSED run is one click
  from its waterfall (`flow_status`) and its two buttons again.
- **The Tools page asks in the page (C26).** Call collected its arguments with a native
  `prompt()`, which the desktop app's pane dismisses unseen and which the spec's §4.7 forbids
  anywhere. The modal the page already owns takes the arguments as a form -- the tool's declared
  arguments listed over the box -- and paints the answer with a way to call again. The eval
  scorer's three `prompt()`s went the same way. No page opens a browser pop-up now.
- **A watched turn ends the way every turn ends (C27).** `endWatched` said 'done'; Home and Chat
  listen for 'end', so a turn started in another window kept its bubble marked live after its
  delivery had landed.
- **Closing says when no toll is owed (C28).** Three places said closing "pays its toll"; the
  core pays one only if a turn ran. The door counts the runs; the glass reads the count.
- **The covenant is read off the record (C29).** The sidebar and the TUI's banner carried the
  house covenant as a literal -- the last two of the 59 copies the door stopped minting from on
  2026-09-25. `vc.CovenantOf(home)` reads it off the operator's own declaration
  (`agents/operator.us`); the TUI's `banner` prints it or says "none declared", in lines of one
  width; the sidebar reads it through the door (`us_to_vc` on that declaration) and says when the
  door did not answer.
- **The builder has a box for loops (D4):** on the five kinds that do work, never on an eval or
  a gate, which the engine refuses by name; typed as a number so the door does not refuse the
  whole save. A looping flow no longer has to be written as JSON by hand.
- **The Dashboard runs the live check (D5):** a button beside Boot, offered only with no engine
  open because the standup opens a sitting of its own; the tally lands in the boot report and
  the deck reads its line back from `tests/run_history.jsonl`.

**Proved.** Ten strokes: `db.TestGetAgentHandsBackACopyAndNotTheStoresOwnMemory` (a caller's
write must not reach the store; an upsert must not move under a reader's copy);
`vc.TestTheCovenantIsReadOffTheOperatorsOwnDeclaration`;
`atlas-tui.TestTheBannerReadsTheCovenantOffTheRecord` (one width, no literal in main.go); and
in `server/glass_test.go` the pages' own source held to each promise the way the What's-left
stroke does -- `TestNoPageOpensABrowserPopUp` over every script, and one stroke each for C23's
read, C27, C17 with D5, C28, C29's sidebar, and D4 with C25. Fourteen reversals on a scratch
copy, red by name: the store's own memory handed out; the 502 read as no engine; 'done'; a
prompt() back; the literal back in the sidebar and in the banner; the sidebar's read never
made; the covenant read off the wrong declaration; the idle line silent on the close; the
button gone; the toll words back; the loops box gone from `run`; loops offered on a gate; the
runs not read on arrival. `go vet` and `gofmt` clean, `node --check` on every script touched;
on the ground `python tests/prove.py --check` **22 held - 14 absent - 0 broke**.

**RESTART REQUIRED, AND DONE:** the glass rebuilt from the ground and placed on his standing
word (pid 27136, the 0.1.9 build of 2026-09-29, stopped by pid and path; the build placed
hashing as built; pid 111424, 2026-09-30 07:29:38, "service wire held"). Live after it, in the pane: the sidebar's covenant read through the door; "Run the live check" beside Boot; the hero saying when the toll is owed; the Tools page's Call for `muster` asked in the modal and answered inside it; the Workflows page reading 4,045 characters of runs on arrival -- no run was paused in the record, so the card had nothing to list -- and the loops box on `coder`'s work steps and not on its check or its gate.

**What goes red if unplugged:** the ten strokes, here and on GitHub's push.

**Not built, and why (D6, a merge button):** `merge` is a forbidden verb on this door by
construction (the first lines of `internal/tools/tools.go`), so a merge button needs his
ruling before it needs code; it is put to him on the list (B16).

### The door's batch: a flow's turn is sent unattended, `run_start` takes a head, and the two streams that had kept an old fault (operator, 2026-09-29: "do the list top to bottom")

Three lines of the core's WHAT'S LEFT, all in THE LINE: C14's other half, C18 and C24.

- **A flow's turn is sent unattended (C14).** `engine.RunUnattended` is a second entry point
  beside `Run`, and it alone puts `"unattended": true` on the wire (serve.py, protocol 1). The
  council's `Turn` -- a flow's `run` node -- uses it; `run_start` and the glass's stream, his
  own turns, say nothing. On that word the core no longer asks "retry / skip / abort?" when a
  seat marked `On Fail: prompt` fails: it skips the seat, the prompt's own default, and writes
  that it did (the core's CHANGELOG, 2026-09-29). Until this the door sent every turn as though
  his hand were on the keyboard, so a flow stopped at that question, `Turn` refused it -- a
  flow cannot answer -- and the run died on a seat's first failure. Any OTHER question still
  stops the turn and is refused in words: that gate is his (RULE 6).
- **`run_start` takes a head (C18).** `voice?` runs every seat on one model for the turn and
  `voices?` names a model per seat over it -- flow_run's two words, read by flow_run's reader
  (`seatVoices`, which `flowVoices` and the new `runHead` both call), so a head that flow_run
  would refuse is refused here in run_start's own name, and before the engine is looked for.
  No head named sends nothing: the ground's declared targets, the wire as it was.
- **`/run/listen` does not stall behind a browser that left (C24).** Its sink was the bare
  send `/run/stream` gave up on 2026-09-15: sixty-four events of room and nobody reading once
  the tab had closed, so the sixty-fifth stopped the engine's reader mid-capture and the world's
  run lock was held for good. An event that would wait on a reader who has gone is dropped.
- **`/chat/stream` has one writer, and none once the browser has gone (C24).** The token
  callback ran on the send's goroutine and wrote each token to the socket itself: two
  goroutines on one ResponseWriter, and -- once the browser had gone and the handler had
  returned -- writes to a ResponseWriter net/http had taken back. Tokens are handed to the
  handler's goroutine, which owns the socket, as `/run/stream` hands over the engine's events.
  The send still finishes, whole and witnessed, with nobody watching.
- **`tests/PROVING.md`, `docs/SPEC_CONTROL_CENTER.md`:** what of `runstream.go` is covered now
  and what is not; the reconciliation table's `run_start` row carries the head and the word.

**Proved.** Strokes on the far end of the pipe: this test binary run again as a child that
speaks serve.py's wire and hands every row back (`tools/council_test.go`; the engine's and
httpserver's stubs extended). `TestAnUnattendedTurnSaysSoOnTheWireAndAnAttendedOneSaysNothing`
(engine); `TestAFlowsTurnIsSentUnattendedAndHisOwnIsNot` and
`TestRunStartTakesAHeadTheWayFlowRunDoes` (tools); `TestACaptureOutlivesTheBrowserThatStartedIt`
and `TestTheChatStreamHasOneWriterAndNoneOnceTheBrowserHasGone` (httpserver; the second against a
loopback stand-in for the rack and a socket that counts the writes made after the handler
returned). Eight reversals on a scratch copy, red by name: the flow's turn sent as though he
were at the prompt; the word never put on the wire, seen from the engine and from the door;
run_start reading its head and dropping it; its refusal speaking in flow_run's name; the head
gone from its arguments; `/run/listen`'s sink waiting on a browser that has gone (red after the
twenty seconds the stroke allows); `/chat/stream`'s tokens written by the send's own goroutine.
On the scratch copy `go test ./...` green but for the six `git_tag` strokes this machine's
scratch folder cannot run, and the battery PROVEN, 142 legs, 85 tools; on the ground `python
tests/prove.py --check` **22 held - 14 absent - 0 broke**, `gofmt` and `go vet` clean.

**RESTART REQUIRED, AND DONE:** the door rebuilt from the ground and placed on his standing word
of the evening (pid 276, the 08:55 build, stopped by pid and path; the new build hashing as
built; pid 106660, 18:30:28). Live after it: sitting 306 opened from the glass, one
`run_start` with `voices: {Steward: llama3.2:latest}` delivered in 6.4 s, closed with its toll.

**What goes red if unplugged:** the five strokes above, here and on GitHub's push; and a flow
whose seat fails, which stops at the question again.

**Not built, and where it went:** C23 names the glass's code (`webapp/db`'s `GetAgent`,
`council.js`'s `Run.check`), not the door's; it is the glass's batch.

### Every cutter asks its word before it acts: a word it does not know writes nothing (operator, 2026-09-29: "do the list top to bottom")

The flow cutter learned this at v0.1.9, on the day the hand cut its fixture with `--check`. THE
COUNT IN THAT ENTRY WAS SHORT. It said eighteen others ended on the same line. There were
TWENTY-FOUR others with the fault: eighteen on the one line, `cut_mesh_vectors` on two branches,
`cut_canon_vectors`, `cut_jesster_vectors` and `cut_us_vectors` through a `main()` that read the
word itself -- and two that are not named `cut_*_vectors.py` at all, `cut_chain_verdicts.py` and
`fold_agents.py`, found by reading every script in `tools/`.

- **`tools/cut_words.py` (new):** the words a cutter knows, in one place. `word(argv, script)`
  answers `--verify` or `--cut` (said, or nothing said). For anything else it prints the usage
  and exits 2, having written nothing.
- **Twenty-five scripts ask it,** as the first act of their `__main__` block. The flow cutter's
  own guard of v0.1.9 is the shared one now; what it refuses, and how it says so, did not move.
- **Run with `--verify`, `cut_words.py` is the proof, and a leg of THE BALL** by the rule
  prove.py already had (every file in `tools/` that answers `--verify`). WHICH SCRIPTS IT ASKS
  IS A RULE, NOT A LIST AND NOT A NAME: every script beside it that takes its word off
  `sys.argv` and knows `--verify`. Each one's `__main__` block is READ first, and a script that
  does not ask is failed and NOT RUN -- running it with a strange word is the fault itself. The
  rest are run with a word nobody knows and must exit 2 with every file under `tests/fixtures`
  and `agents` byte for byte as it was.
- **Not touched, and why:** `seed_catalog.py` parses its words with argparse, which refuses an
  unknown one itself; `check_trade_parity.py`, `cut_fixtures.py` and `lint_us.py` have no second
  mode to mistake.
- **`AGENTS.md`, `CONTRIBUTING.md`, `tests/PROVING.md`:** twenty-eight verifiers where they said
  twenty-seven; PROVING says what a cutter's words are, and its standing line carries the date
  and the command it was measured with.

**Proved.** On a scratch copy, then on the ground, byte for byte the same files:
`python tools/cut_words.py --verify`: PROVEN, 25 cutters. Reversals on a scratch copy, red by
name: `fold_agents` no longer asking; `cut_chain_verdicts` importing the word and never asking
it; a vectors cutter reading its word the old way (each failed "does not ask `word` before it
acts; not run"); `word` letting a strange word through (25 failed); a cutter that refuses by
the book and writes into `agents` first, and one that writes into `tests/fixtures` first (each
failed "exit 2; MOVED:" with the file named). AND THE RULE ITSELF: with the scripts picked by
name again, a `fold_agents` that does not ask went unseen and the proof passed at 23 -- which
is why they are picked by what they do. `python tests/prove.py --check`, on the ground:
**22 held - 14 absent - 0 broke** (21 before: the new leg). No caller hands a cutter any other
word: `prove.ps1` and prove.py say `--verify`, and nothing in `.github/` calls one.

**What goes red if unplugged:** prove's `cut_words` leg, here and on GitHub's push.

**Found, not built:** `seed_catalog.py --verify` is described as "check only" and is not. It
never reads the flag: it creates `data/master.db` when that is absent and seeds it before it
checks (the seed adds nothing twice). prove.py runs it on every proof (the core's WHAT'S LEFT,
C35).

## [v0.1.9] — 2026-09-29 15:15 (tag on b1059a1)

### The flow verifier knows the bounded return, and a word it does not know writes nothing (operator, 2026-09-29: "do the list top to bottom")

**Found on the way to the mark.** `python tests/prove.py --check` said "20 held - 14 absent -
1 broke": `cut_flow_vectors`, "fixture drifted from the contract / drift: refusals". It had been
red since 2026-09-28 (`325888c`), when the bounded return landed: `tests/fixtures/flow_vectors.json`
gained five refusal vectors by hand, the Go strokes that read it passed, and the script that is
the fixture's ORACLE was never taught them -- its `topo()` knew no `loops` at all.

- **`tools/cut_flow_vectors.py`:** the oracle learns the return, in the runner's own shape
  (`flow.go`'s `loopsOf` and `lawfulReturns`). `loops` is 0 to 5 and never on an eval or a gate;
  an eval's fail-edge into a node that declares `loops` is a RETURN and leaves Kahn's count; a
  check returns once; a return goes back; the body it goes around holds no gate and does work;
  a declared ceiling is reached by something. `vectors()` carries the five the fixture holds.
- **A lawful return is proved NOT refused.** Every refusal vector only asserts that a spec is
  refused; an oracle that refused everything would pass them all. `--verify` now also walks
  work, a check and the check's way back, and wants the order `w, c`.
- **A word it does not know writes nothing.** Anything that was not `--verify` used to CUT.
  THE HAND DID EXACTLY THAT TODAY: it ran the script with `--check`, prove.py's own word, and
  the fixture was rewritten from the script's seven vectors, the five gone. Nothing was saved in
  between; the file was put back byte for byte from the last save (`git show HEAD:`), and the
  diff against it reads nothing. `--cut` and no word at all still cut; anything else is refused
  with the usage and exit 2.

**Proved.** `--verify`: PROVEN. A cut now reproduces the saved fixture byte for byte (7,185
bytes), asked of the module without writing. Reversals on a scratch copy, red: a ceiling
nothing returns to let through; a gate inside the return let through; the return not known
(the lawful return refused: "want exactly one start, got 0"); a vector renamed (drift). NOT
RED, and why: `loops` on a check let through -- that spec is also refused for re-doing no work,
so the rule is a better sentence and not a separate wall. `python tests/prove.py --check`:
**21 held - 14 absent - 0 broke**.

**What goes red if unplugged:** prove's `cut_flow_vectors` leg.

**Found, not built:** eighteen other `tools/cut_*_vectors.py` end on the same line and cut on
any word that is not `--verify` (the core's WHAT'S LEFT, C32).

### What's left: everything still open, on one page of the glass (operator, 2026-09-29: "write a page on the webapp")

His order, in his words: "a page with any leftover tasks, checklists, seams ... any blockers or
gates, or decisions ... problems or bugs, needs to be in one area where we can see what is
actually left".

- **`webapp/static/js/left.js` (new):** the page. It reads ONE document, `WHATS_LEFT.md` at the
  root of the ground, through the door's `records` tool -- the same read Records makes -- and
  keeps no copy. A tile for each lettered section with the count of its open lines (A, what
  stops the release, in red; B, what waits on him, in yellow), a card for each section, one row
  for each line: its number, what it is, the words, and where it came from. A find box over
  every row; a section with no row showing is put away with them. "Read it again" reads it again.
- **Every number is counted from the lines drawn.** The document carries no tally of its own.
- **A line the page cannot name is said in red at the top:** a list line under a lettered
  section with no number, a number used twice, a number under the wrong letter.
- **`webapp/static/index.html`:** the menu line, second under Work, and the script, loaded before
  `app.js`. **`webapp/static/js/app.js`:** the router's case. **`webapp/static/css/icons.css`:**
  the glyph, on the same grid and stroke as the rest.
- **Nothing in the door moved.** `records` already serves every root `.md` of a ground.

**Proved.** `server.TestTheWhatsLeftPageIsWired`, and four reversals red (the router's case
removed; the script gone from index.html; the page reading another document; the menu line
removed). `go test ./...` green on a scratch copy of `webapp`, `gofmt -l` and `go vet` clean.
The page's reader run on the real list under node: ten sections, 80 open lines, no fault. Drawn
in the app's browser pane against his running glass with its own styles: 80 rows, 7 tiles, the
find box answering `router` with A1, B10, C1, C2, C5, C9, E3.

**Restart required** for the webapp: the page is `go:embed`ed.

**What goes red if unplugged:** the stroke above; and on the page, the red card at its top.

### Changed — a run node's turn is handed the failed pass as its feed (operator, 2026-09-29: "carry the failed pass without the door's name")

- `flow/run.go`, `execNode` "run": `eng.Turn(ctx, objective, feed, method)` had carried "" for the
  feed; it now carries `fail_<node>` -- why the last pass was sent back and what the checked node
  answered, bounded (loopVars, `carried`) -- which `councilEngine.Turn` forwards to the council as
  the turn's source material. Seeded only for a node that declares `loops` and empty until a pass
  was sent back, so a first pass and a node that never returns hand over nothing; a question may
  still name `{{fail_attempt}}`. WHY: the fifth firing of the core's coder-tree measured that a
  sent-back pass quotes a door's own reply, and the council reads a tool's name in the OBJECTIVE
  as a request for that tool, so a question carrying the pass shut the Coder's window on every
  retry and v4 carried nothing. The feed is the council's second channel: shown to every seat,
  routed on by none.
- Stroked: `TestTheFailedPassRidesAsTheFeedNeverInTheObjective` (a first pass fed nothing; the
  retry's feed carrying the failure and the first answer; the objective the words alone on both);
  the loop stroke unchanged. Reversal R13 (the feed handed as "" again): red. The flow package
  green, `gofmt` and `go vet` clean. Battery: PROVEN, 142 legs, 85 tools, from the built binary.

### Added — the suites run from the glass, and the head is read off their stamp (operator, 2026-09-28: "if its on the glass, and the record mathes, id call it proof"; by card, suite_run)

The third piece of the code safety pass, and the check the coder's loop will steer on. Until now
the strokes and the smoke ran only on his terminal, and a hand proved on a mirror; his ruling makes
a run from the glass, whose record matches, proof.

- `suite_run` (`set?` = strokes | smoke | both, the default; TierEngine): the door runs the world's
  own `tests/test_manjuel.py` and `tests/smoke_cli.py` with the python it runs the engine with, ONE
  AFTER THE OTHER in the release gate's order, bounded at thirty minutes together. The suites stamp
  their own proof as they always have -- `tests/last_run.json` says `running` at the start so a
  crash cannot leave a green stamp, the tally at the end, and a line in `run_history` -- and the
  head is read BACK OFF THAT STAMP, never parsed from the script's prose: `SUITES both on
  "<world>" · strokes: N/N green · exit 0 · smoke: N/N green · exit 0`, with a suite that never
  finished said as `DID NOT FINISH (the stamp still says running)`. The body carries each suite's
  `[FAIL]` lines first (bounded at sixty) and its last lines with the tally; a pass is a count, not
  a line. ONE SUITE AT A TIME, MACHINE-WIDE (`suiteLock`, TryLock): a second call while one runs is
  refused, not queued -- a stroke counts the browsers alive on this computer and two suites at once
  red it for each other. Refuses by name: no core command; a world that carries no suite; a set not
  one of the three. A READER in the door's eyes (`Writes: false`): it writes the suites' own stamps
  and nothing of the work, so the coder's loop can ask it without a hand parking the call (LAW_003
  §3: the gate stands at the end, not inside the loop).
- Proved: `TestTheSuitesRunFromTheGlassAndTheHeadIsReadOffTheStamp` (the three refusals, the tool
  declared a reader; then against stand-ins that stamp the way the real suites do: the head read
  off the stamp suite by suite with each exit, the body carrying the failures once and the tallies
  and no pass line, the strokes before the smoke, the smoke set leaving the strokes untouched, a
  crashed suite said off its own stamp, a second call while one runs refused); two battery legs.
  Reversals: R1 the set running one suite where it should run two -- 1 red; R2 the verdict left off the head -- 1; R3 the lock released before the run -- 1; R4 the failures dropped from the body -- 1.

**RESTART REQUIRED:** the door carries this once rebuilt and restarted; **85 tools** by its own
battery (84 before). **WHAT GOES RED IF THIS COMES UNPLUGGED:** the stroke, if a set stops running
its suites in order (R1), the head stops reading the stamp (R2), the lock is struck (R3) or the
failures leave the body (R4); the battery, if the door stops naming its missing python or refusing
a world with no suites. NOT BUILT: the coder-on-the-tree flow that asks it (next by the card's
order), and a Dashboard button.

### Added — a check's fail-edge may return the run to the work, bounded (operator, 2026-09-28: "create the bounded back-edge looping"; by card, the ceiling lives on the node returned to)

LAW_003 was sealed 2026-09-17 and the flow engine still refused every cycle by name: the coder
flow's retry was one repair pass unrolled into a line, and a FAIL could never be re-attempted. A
node may now be RETURNED TO, and only with the law's three bounds, each refused by name at the save
(section 5) and none inferred at run time.

- `flow.Node.Loops` (0 to `MaxLoops`, 5): how many times THIS node may be returned to -- the
  ceiling, on the node the work starts at, where a reader meets it. A return is a check's `fail`
  edge into a node that declares it (`loopsOf`, ONE reading for `Validate` and the runner); every
  other back-edge is still the cycle it always was. `Validate` refuses: a ceiling nothing returns
  to ("a ceiling read by nothing"); a return that does not go back; a gate inside the return (a
  gate stands at the end of a loop, section 3); a return that re-does no work (section 4); `loops`
  on an eval or a gate (the work is returned to, never the verdict); a check that returns twice.
- The runner: when the check fails and the ceiling has room, a `loop` line is written (`node`,
  `check`, `checked`, `return` k, `ceiling`, `why`, `status: returned`), the body -- the node, the
  check and every node on a path between -- is unfired, and the walk resumes AT the node, so every
  pass writes fresh node lines with receipts (section 2c). At the ceiling the `loop` line says
  `spent` and the fail is a fail: a forward fail-edge if the check has one, else FAIL. The budget
  binds every pass; a cancel still stops. `pass_<node>` ("1" the first time) and `fail_<node>`
  (why the last pass was sent back and what the checked node answered, bounded to 4000 runes,
  empty until a return) are seeded on every pass, so the retry is told what failed instead of
  rolling the same dice; `Resume` rebuilds both off the run's own `loop` lines by the same
  function. `flow_status` renders each return and a spent ceiling with its why.
- Proved: `TestACheckThatFailsReturnsTheRunToTheWorkAndTheCeilingHolds` (three passes under a
  ceiling of two, the work re-done and the retry told the failure, the record `attempt verify
  verdict loop attempt ... land:paused` exactly, the waterfall, the last pass standing as the
  answer; the never-passing case ending FAIL after exactly three passes with `spent` on the record;
  the same spec without `loops` refused as a cycle), `TestASpentCeilingFollowsAForwardFailEdge`,
  `TestAResumeRebuildsTheLoopCount`, `TestAReturnIsRefusedWhereItIsNotALoop` (ten refusals by
  name; two lawful shapes, a return to the start node among them); five refusal vectors added to
  `tests/fixtures/flow_vectors.json`, held by `TestFlowContract`; two battery legs. Reversals: R1
  the ceiling not read -- 1 red; R2 no loop line -- 1 red; R3 the unreached ceiling accepted -- 1
  red and the fixture's golden; R4 the resume not rebuilding -- 1 red; R5 the body not unfired -- 1
  red.

- FOLDED ON THE PLACED DOOR, as the glass's wire, on his card: the core's `coder` flow as v13 --
  `brief -> attempt(loops 2) -> verify -> verdict`, `verdict --fail--> attempt`, `--pass--> land`;
  `repair`, `recheck` and `proof` struck, v1-v12 kept on disk.

**RESTART REQUIRED:** the door carries this once rebuilt and restarted (done 2026-09-28, pid 86180).
**WHAT GOES RED IF THIS COMES UNPLUGGED:** the runner stroke, if the loop stops reading the ceiling (R1), stops writing the return
down (R2) or stops re-firing the work (R5); the refusal stroke AND the fixture's golden, if the save
stops refusing a ceiling read by nothing (R3); the resume stroke, if a carried-on run forgets its
count (R4). NOT BUILT: a `loops` box in the glass's builder -- a looped spec is folded with
`flow_save` and the JSON.

### Added — the standup runs from the glass (operator, 2026-09-28: "fire the standup through the glass"; by card, a door tool that runs it)

- `standup_run` (`set?` = morning | court | all; TierEngine; writes): the door runs the world's own
  `tests/standup.py` with the python it runs the engine with (`--manjuel`'s first word), bounded at
  thirty minutes, in the world, and hands back one head -- `STANDUP <set> on "<world>" · exit N ·
  <tally> · report <path>` -- over the whole of what the script printed (its words outrank the exit
  code, ADR-006 item 2; exit 1 is a missed expectation, 2 could not run; `\r\n` read as one shape).
  Refuses by name: no core command; a world that carries no `tests/standup.py`; a set not one of the
  three; an engine open on the world, or a sitting open in its ledger -- the standup opens and tolls
  a sitting of its own, and the sitting line is the lock. It takes no rack lock: the standup is a
  whole process on the rack, as a council turn is. A script that prints no tally is said so, never
  invented.
- Proved: `TestTheStandupRunsInTheWorldAndRefusesByName` (the four refusals; then, where a python
  is on the PATH, against a stand-in script: the court's flag reaches the script, the head carries
  set, exit, tally and report, the morning set runs bare, and a script with no tally is said so);
  two battery legs. Reversals: R6 the set's flag dropped -- 1 red; R7 the open sitting not refused
  -- 1 red.

- FIRED FROM THE GLASS on the placed door, his order: the morning set, sitting 284, 9/9 LIVE, report
  `logs/standup_2026-09-28_104238.md`, the tool's head `STANDUP morning on "research" · exit 0 ·
  9/9 cases met their expectations. · report ...`; then the court on his ruling ("measure first"),
  sitting 285, 0/1, `exit status 1` carried in the head with the script's whole account under it.
  The glass's Tools page asks for arguments with `prompt()`, which the desktop's browser pane
  dismisses unseen; the page's own `App.tool` made the same request.

**RESTART REQUIRED:** the door carries this once rebuilt and restarted (done 2026-09-28, pid 86180);
**84 tools** by its own battery and by `tools/list` (83 before). **WHAT GOES RED IF THIS COMES UNPLUGGED:** the tools stroke, if the set stops
reaching the script (R6) or an open sitting stops refusing (R7); the battery, if the door stops
naming its missing python or refusing a world with no standup. NOT BUILT: a Dashboard button for it
-- the glass's own piece; it is fired from the glass's Tools page.

### Added — a gate declares what its crossing grants (operator, 2026-09-28: "then the crossed gate"; by card, "the gate declares what it grants")

Under `--auth` the council's writes parked even after the operator had resumed a gate that asked
exactly that question, so `version-tag` could read the marks and never cut or send one through the
council: his `continue` at the gate and his approval of the hold were one decision asked twice, and
the flow recorded the hold's words as the node's answer.

- `flow.Node.Grants`: a GATE's declaration of the writing tools its `continue` authorises, from that
  gate until the next gate or the end of the run. Only a gate carries grants (`Validate` refuses the
  rest by name, and an empty name); the pause and the resume lines carry them, and `flow_status`
  says "continue grants the council: ..." beside the question. `runFrom` binds a crossed gate's
  grants onto the engine that fires the nodes after it (`onHand`, a copy like `onHead`).
- THE LINE: `councilEngine.WithHand` takes the crossing (a copy; the registered engine stays
  unbound) and `underHand` raises the hand over the world for the length of one turn and lowers it
  after. `tools.Call`, about to park a write from a key, asks whether the hand stands over that tool
  on that world: if it does the call runs and the holds record gets a `crossed` line naming the
  tool, the caller and "run <id> gate <name>"; if not it parks as before. RBAC still judges first.
  `flow_save` refuses a grant naming a tool the door does not carry, or one that does not write, at
  the save.
- Proved: flow -- `TestACrossedGateCarriesItsGrantsToTheNodesAfterItAndNoFurther` (the pause and
  the resume carry the grants, the waterfall says them, the node after the crossed gate runs under
  them and the node after a gate that grants nothing runs under none; an engine that takes no hand
  is left alone, a gate with no grants binds nothing, a bound copy leaves the engine unbound) and
  `TestOnlyAGateGrants`. tools -- `TestACrossedGateCarriesTheHandToTheToolsItGrantsAndNoOther` on a
  real repository with holds armed (the cut parks with no hand; under the hand it reaches the tool,
  which refuses on its own law, while an ungranted writer still parks; lowered after the turn the
  same cut parks again; exactly one `crossed` line naming the run and the gate; a handless council
  raises nothing) and `TestFlowSaveRefusesAGrantOnAStrangerOrAReader`. Two battery legs. Reversals:
  R1 the door never asks for the hand -- 1 red; R2 the crossing never bound -- 1 red; R3 the save
  accepting any grant -- 1 red and the battery; R4 the crossing leaving no record line -- 1 red.
- `version-tag`'s two gates, `judge` and `send_gate`, declare `git_tag`: folded as v3 through
  `flow_save` on his word ("fold it, as the glass"), everything else in the spec unchanged. The
  core's RUNBOOK says what a gate grants.
- PROVED LIVE on the placed door, on his word: `version-tag` v3 fired on research for a mark it
  already carries; the read step listed the marks through the council (not held); the judge gate
  paused saying "continue grants the council: git_tag"; the crossing ran the cut under the hand --
  the council's `git_tag cut` reached the tool, which refused on its own law ("already exists
  here ... A mark is never moved"), the proof step failed honestly, no mark moved, and the holds
  record carries the `crossed` line naming the run and the gate.

**RESTART REQUIRED:** the door carries this once rebuilt and restarted (done 2026-09-28). **WHAT GOES RED IF THIS
COMES UNPLUGGED:** the tools stroke, if the door stops asking for the hand (R1) or stops writing the
crossing down (R4); the flow stroke, if the crossing stops being bound (R2); the save stroke and the
battery, if a grant is no longer judged at the save (R3).

### Fixed — a secret parked in the holds is withheld where the queue is shown (operator, 2026-09-28: "take the hold queue's key argument next"; by card, "withhold it where it is shown")

Found on 2026-09-28 while widening the council's key: a parked call keeps its arguments whole, and
`hold_list` handed them to the glass whole, so the re-proof `key` a keyed caller sends to
`auth_key_scope` (or create, or revoke) would have sat in plain view on Version control (RULE 7).
The widening rode the glass's wire that day to avoid it.

- `Tool.Secrets` names the arguments whose values are secrets; the four auth verbs declare `key`.
  `hold_list` shows a parked call with those arguments as `[withheld: a secret; the parked call
  keeps it]`, and any value in the door's own key shape (`auth.KeyRe`) is withheld even where nobody
  declared it, on a tool the door no longer carries included. The parked call keeps the real values
  and approving runs it whole. The holds record never carried arguments.
- Proved: `TestASecretArgumentIsDeclaredWhereverAKeyIsTaken` -- every tool that takes `key` declares
  it, every declaration names an argument the tool takes, and `shown` withholds by the declaration
  (any shape), by the shape (any name) and on an uncarried tool, while a plain word stays.
  `TestAParkedSecretIsWithheldWhereItIsShownAndRunsWhole` -- a scope move parked from a key: the
  held answer, the queue and the record all without the plaintext, the rest of the call shown whole,
  and the approval running it with the real key so the scope moves. One battery leg. Reversals: R1
  the verb's declaration struck -- 1 red; R2 the shape belt struck -- 1 red; R3 `hold_list` handing
  the arguments over whole -- 1 red and the battery; R4 a reading verb's declaration struck -- 1 red.

**WHAT GOES RED IF THIS COMES UNPLUGGED:** the declaration stroke, if a tool takes a key and does
not say so (R1, R4) or the shape stops being read (R2); the parked stroke and the battery leg, if
the queue hands the arguments over whole again (R3).

### Added — `auth_key_scope`: a live key's tenants move without the secret (operator, 2026-09-28: "then the key's scope"; by card, "a scope verb, then widen it")

The council's key (`k-ae2e9481`, `manjuel-council`) was minted for `research` alone while RUNBOOK
said research and atlas, so on the atlas tenant the door refused it `403` before RBAC was asked --
and the only verb that could have changed that was `auth_key_create`: another key, and the
plaintext rotated through `.env`. A scope is a fact about the key's record, not about the secret,
so it moves without the plaintext ever existing again.

- `auth.Scope(home, id, tenants, by)`: the list replaced whole on a LIVE key -- widen and narrow are
  the same act -- never empty (a key that carries nothing is a revocation by another name; revoke
  it), a revoked or unknown id refused by id, and an audit line `scope` with the list it set.
- `auth_key_scope` (`id`, `tenants`, `key?`): possession re-proved like create and revoke (a live key
  named); every tenant named must be one this door carries, refused by name otherwise -- a scope
  onto a stranger is a permission read by nothing; it writes, so from any hand but the glass it
  parks in the holds for the operator's approval, which runs exactly the parked call.
- Proved: `TestAKeysScopeMovesWithoutTheSecret` (auth) -- the same plaintext verifies wider after the
  move, `ScopeOK` carries what was given and nothing else, narrowing replaces the list whole, the
  three refusals, and the audit holding create, scope, scope, revoke in order with what each set.
  `TestAKeysScopeIsMovedByAProvedHandOntoCarriedGroundOrHeld` (tools) through `Call` with holds
  armed -- no re-proof, a stranger tenant, an unknown id and no tenants each refused with nothing
  moved; from a key the move is HELD and moves nothing until the glass approves the hold, which runs
  it, and the same plaintext carries both grounds; the glass with a live key moves and narrows at
  once. Three battery legs. Reversals on a scratch copy: R1 the re-proof struck -- 1 red; R2 the
  audit struck -- 1 red; R3 a stranger accepted -- 1 red and the battery; R4 the verb declared
  non-writing -- 1 red. Restored: green but the scratchpad's seven.
- The core's `RUNBOOK.md` says the key was minted for research alone and how a scope moves.
- DONE ON THE PLACED DOOR, on his word: `k-ae2e9481` widened to `research,atlas` (the store and the
  audit line say so); before it, `muster` on atlas with the council's bearer was `403`; after it,
  the roster, `tenant_rbac_check` ALLOWED as `steward`, and `git_tag list` answered -- the list,
  not a hold. FOUND ON THE WAY, NOT BUILT: a call parked in the holds has its ARGUMENTS shown on
  Version control (`hold_list`), so a re-proof `key` parked there would be displayed -- RULE 7. The
  widening therefore rode the glass's own wire (the operator's hand, no hold) with the key in the
  body over loopback. What the queue shows of a `key` argument is a ruling.

**RESTART REQUIRED:** the door carries the verb once rebuilt and restarted (done 2026-09-28). **WHAT GOES RED IF THIS COMES
UNPLUGGED:** the tools stroke, if the re-proof, the carried-tenant refusal or the writing declaration
goes (R1, R3, R4); the auth stroke, if the move stops being audited (R2).

### Fixed — a reading action of a writing tool is not held (operator, 2026-09-28: "take the held git_tag list next")

`version-tag`'s read step was parked on 2026-09-26: with `--auth` on, the council's key asked
`git_tag` for the LIST of marks and the hold queue parked it as `hold_1790443435789_1_git_tag`,
because `Writes` is one flag for the whole tool and `git_tag` bundles `list` with `cut`, `send` and
`remove` under one `action` argument. `git_branch` has the same shape (`list` beside `new`, `switch`,
`close`). So under `--auth` the flow could not even read the marks through the council, and RBAC
judged the list an `edit`.

A writing tool now declares the actions that only read -- `Tool.Reads`; `git_tag` and `git_branch`
declare `list`, which is also their default -- and `Tool.WritesFor(args)` is what this CALL does:
the declaration narrowed by the action, the empty action being the tool's own (reading) default.
`tools.Call` asks it once and both the hold queue and RBAC read that one answer: `git_tag list` is
answered for any bearer and is a `read` to a role; every other action is parked for the operator's
hand and is an `edit`, exactly as before. `actionOf` is the one reader of the word -- lower-cased,
trimmed -- and both handlers read through it, so the door and the tool cannot judge two different
words.

- `tools.Tool.Reads`, `Tool.WritesFor`, `actionOf`; `Call` narrows RBAC's kind and the hold by the
  call; `git_tag` and `git_branch` registered with `Reads: []string{"list"}`; `toolGitTag` and
  `toolGitBranch` read the action through `actionOf`.
- Proved: `TestAReadingActionIsDeclaredOnlyWhereItCanBeRead` -- a reading action is declared only on
  a tool that writes and takes an action, as one lower-case word; the two that earned it keep
  `list`; and `WritesFor` on the declaration alone (the default and the word read; `cut`, `send`,
  `remove`, `new` and nonsense write; a tool without reading actions is what it says whatever the
  action). `TestAReadingActionOfAWritingToolIsNotHeld` on a real repository with holds armed -- the
  list answered in three spellings (none, `list`, ` LIST `) and `git_branch list` too, the cut and
  the open parked with nothing landed, the glass never held, and a role that may not write (the
  shipped `agent`) listing the marks and refused the cut by kind before any hold. One battery leg
  (holds armed, no repository: the list reaches the tool and is refused as "not a repository" --
  the tool speaking, not the queue -- while the cut is parked). Reversals on a scratch copy: R1
  `git_tag` declares no reading action -- 2 red and the battery; R2 `WritesFor` ignores the action
  -- 2 red; R3 the door reads the word raw -- 2 red; R4 the handler reads it its own way -- 1 red.
  Restored: green but the scratchpad's seven "Filename too long" strokes, unchanged.
- `docs/ARCHITECTURE.md` 4d says it; the core's `RUNBOOK.md` names which actions wait.

NOT CHANGED, and named: under `--auth` the flow's `cut` and `send` steps are still parked -- that is
RULE 6 working -- so `version-tag` cannot cut through the council; its read and check steps can now
list. Whether a gate the operator crossed carries his hand to the next node's writes is a ruling,
not a patch.

**RESTART REQUIRED:** the door carries it once rebuilt and restarted. **WHAT GOES RED IF THIS COMES
UNPLUGGED:** the two strokes and the battery leg, if `git_tag` stops declaring `list` (R1) or `Call`
stops asking `WritesFor` (R2); the repository stroke, if the door and the handler stop reading the
one word (R3, R4).

### Fixed — the RBAC role model: a role speaks in the kinds the door declares (operator, 2026-09-26: "kinds from the tool's own declaration")

Found closing P0-13 and left as a design ruling: the shipped `DefaultPolicy` roles carried permissions
by KIND (`read`, `edit`, `bash`, `net`, `tools`) while `rbac.Can` looked up TOOL NAMES and a `*`
wildcard, so assigning ANY shipped role to a key denied it every tool -- and no stroke had ever assigned
one to find out. His ruling by card: the kinds come from the tool's own declaration. `Can` now asks a
role three questions in order -- the tool BY NAME, the WILDCARD, then the KINDS the call carries -- and
the kinds are not words the policy invents: `edit` is a tool that declares `Writes: true`, `read` one
that does not, and `tools` (calling at all) rides on every call. A name beats the wildcard and the
wildcard beats a kind, so a role may deny `edit` and still allow `git_tag`. A deny on any carried kind
denies; an allow needs every carried kind; a role that speaks to neither has not answered, and the roles
it implies are asked the same three questions, breadth-first. `bash` and `net` named nothing the door
declares and are gone from the shipped roles; a policy on disk that still carries them (this
repository's own `rbac.json` does) is read as before, and those keys are simply never asked.

- `rbac.Kind`, `KindsOf(writes)`, `IsKind`; `rbac.Can(p, agent, tool, writes)` -- every caller says
  whether the tool writes, so a kind is never guessed from a name; `rbac.RoleOf` for the hold record.
- `tenant.CheckRBAC(agent, tool, writes)`, `tenant.RoleOf`; `tools.Call` hands `Tool.Writes` down with
  the name; `tenant_rbac_check` looks the tool up at this door for its declaration and refuses a name
  the door does not carry rather than judging a word.
- Proved: `internal/rbac` has its first prover, `rbac_test.go` (8) -- every shipped permission is a kind
  the door declares or `*`, every shipped role speaks to all three kinds and the policy ships open; the
  kinds a call carries come from the declaration; each shipped role over a reader and a writer (every
  row was DENIED before); name over wildcard over kind, both ways; an allow needs every carried kind
  and a deny needs one; implied roles asked after the assigned one, a cycle ending, a missing role
  skipped; no role, no entry; a policy on disk carrying `bash`/`net` still deciding by the live kinds.
  Through the door, `TestAShippedRoleAssignedToAKeyDecidesByWhatTheToolDeclares` assigns all four
  shipped roles to keys on a real registry: operator and steward read and write, agent reads and is
  refused a writer by kind with nothing landed, guest is refused calling at all, and `tenant_rbac_check`
  answers the same with the same declaration. One battery leg. Reversals on a scratch copy: R1 `Call`
  passes `false` for the declaration -- the door stroke and the battery red; R2 `bash` back in a
  shipped role -- 1 red; R3 the kinds never asked (the original fault) -- 5 red in the package and 1
  through the door; R4 the tool's name never asked -- 1 red. Restored: green but the scratchpad's
  seven "Filename too long" strokes, identical at HEAD.
- `docs/ARCHITECTURE.md` 4d says the three questions; `tests/PROVING.md` takes `internal/rbac` off the
  unproven list (and `internal/tenant`, off it since 2026-09-11 unrecorded), with the count measured:
  two packages carry no prover, `cmd/atlas-vc` and `internal/orient`.
- **`rbac.json`** (this repository's own -- the atlas tenant's policy), on his word. Found on the
  way: since 2026-09-09 it assigned `agent` to an actor named `5 carried projects` -- muster's first
  line, pasted into the glass's Assign Role box -- which put the atlas tenant in closed mode (the
  2026-09-25 boot line named only `research` as open) and refused every keyed caller there
  `no assigned role`. His ruling by card: "Steward is the default actor." The door names a bearer by
  its key id (`httpserver`: `Caller{Name: cred.keyID}`), and the council's key is `k-ae2e9481`
  (`manjuel-council`), so that id now holds `steward` -- reads and writes, the writes held for his
  hand -- and the accidental actor is struck. The roles in the file keep their dead `bash`/`net`
  keys; nothing reads them.

**RESTART REQUIRED:** the door carries the role model once rebuilt and restarted; the glass is
untouched. **WHAT GOES RED IF THIS COMES UNPLUGGED:** the compiler first -- `CheckRBAC` and `Can` take
the declaration as a parameter, and no caller can leave it out; then the door stroke and the battery
leg, if `Call` stops passing `Tool.Writes` (R1); and `TestEveryShippedPermissionIsAKindTheDoorDeclares`,
if a shipped role names a word no tool declares (R2).

### Fixed — P0-13: RBAC judges the transport, and open mode is said (SPEC_CONTROL_CENTER §12.5; operator, 2026-09-25: "identity from the transport; open mode said out loud")

RBAC failed open twice. The dispatch read `actor` off the caller's OWN args and ran no check when it
was absent -- a gate any caller stepped around by saying nothing -- while `holds.go`, two lines below,
had already refused to read the args for exactly that reason. And a tenant with no roles assigned
allowed everything in silence. Since the door was armed this afternoon it knows who is calling from
the transport, so the check runs on EVERY call with that identity: the glass (the service wire) is
the operator's hand and passes; anything else is judged by the name the door gave it -- a verified
key's id, or what a stdio client called itself -- against the tenant's policy. `actor` in the args is
a label for the record and decides nothing. An empty policy is open mode still, and is SAID: the
boot line names each open tenant (`rbac open (no roles assigned) on: research, atlas`), and every
hold record carries `rbac: open (no roles assigned)` or `rbac: role <name>`.

- `tools.Call`: the check on `caller`, never `args["actor"]`. `tenant.RBACOpen`; `tools.RBACLine`;
  `Hold.RBAC` and the holds log line.
- Proved: `TestRBACJudgesTheTransportNotTheArgs` -- a denying role refused by the caller's own name,
  an allowing role passed, an unassigned caller refused, the args' `actor` decides nothing (the
  forgery path), the glass never judged, open mode passing everyone and SAID by name on the boot
  line and on the parked hold; a battery leg reads the boot line's word. Reversals: the args-actor
  check restored, the hold's line blanked -- each reds its stroke.
- FOUND ON THE WAY, NOT BUILT: the shipped `DefaultPolicy` roles carry permissions by KIND
  (`read`, `edit`, `bash`, `net`, `tools`) while `rbac.Can` looks up TOOL NAMES and a `*` wildcard --
  so assigning a shipped role to a key today denies it every tool. The policy model and its roles
  disagree; that is a design ruling, not a patch.

### Fixed — P0-14: a forbidden verb is never free (B1-02; operator, 2026-09-25: "the invariant that is true")

The absence test compared WHOLE tool names to bare verbs -- no tool is named `commit` -- so it had
never fired. Measured: five tools carry a forbidden verb as a word (`git_commit`, `git_push`,
`mesh_post`, `chat_send`, `team_send` -- the last three by D4's verbs, the first two by B1-02's), and
every one declares `Writes: true`, so every one is held for any caller but the glass. That is the
invariant worth keeping, and it is what the test measures now.

- `tools.ForbiddenVerbs`, one source; `tools.ForbiddenWord` matches the verb as a WORD (`git_commit`
  carries commit; `commitment` carries nothing); `tools.HeldExempt`.
- Proved: `TestAForbiddenVerbIsNeverFree` over the whole registry (and refuses to be vacuous -- at
  least one tool must carry a verb), and the battery leg that replaced the one that never fired.
  Reversals: `git_commit` declared non-writing, the verb matched as a substring -- each reds.

**RESTART REQUIRED:** the door carries P0-12, P0-13 and P0-14 once rebuilt and restarted.

### Fixed — P0-12: an optional argument is optional (SPEC_CONTROL_CENTER §12.5; operator, 2026-09-25: "then the harm list")

`endsWithOptional` demanded TWO trailing `?` -- a shape no tool has ever declared -- so every
argument on every tool, `project?` included, was published as REQUIRED, and the description's own
promise ("omit it and the ground the session was opened in answers") was a law no client could read.
`httpserver` carried a second copy of the same wrong rule, so both schema builders (`GET /tools` and
`tools/list`) agreed with each other and with nothing. Listed as P0 on 2026-09-10; verified live
2026-09-24; closed today.

- `protocol.Optional` / `protocol.TrimOptional`, exported: one `?`, one source. `httpserver`'s copies
  read them; `protocol`'s private names remain as thin wrappers for the one caller inside the package.
- Proved: `protocol_test.go` (the package's first test -- `project?` optional, `path` required, the
  marker stripped and nothing else) and `httpserver_test.go`'s schema walk now refuses any tool that
  publishes `project` as required. Reversal: the two-`?` rule restored reds both.

**RESTART REQUIRED:** the door carries this once rebuilt and restarted.

**WHAT GOES RED IF THIS COMES UNPLUGGED:** the two tests above; a client reading the schema would be
told, again, that `project` is mandatory on 82 tools.

## [v0.1.8] — 2026-09-25 10:36 (tag on 56a3078)

### Fixed — a credential's issuer is minted in the record's own covenant, not in a literal (operator, 2026-09-25: "one source: the manifest; two readers")

The `covenant` on a `.us` record is the DID namespace every credential from that manifest is minted
in (`did:atlas:<covenant>:<id>`, and the reporting line the same way). The issuer carried the same id
by hand: once in the door's `us_to_vc` (`tools.go`), once as `atlas-vc --issuer`'s default — 59
copies of one fact across two repositories, and nothing to say if one moved. The core's reconciler
found the field LOOSE on its side the day LOOSE existed; the core now checks that no record drifts
from the rest of the manifest, and this is the door's half.

- `vc.IssuerFor(block)` mints `did:atlas:<covenant>:operator` from the record, and refuses by name a
  record that declares no covenant. `vc.FromFile(path, "")` mints; a named issuer is a hand's explicit
  choice and is used as given.
- `us_to_vc` passes no issuer; `atlas-vc --issuer` defaults to none, with the help text saying so.
- Proved: `vc_test.go` (the package's first test — subject, reporting line and issuer share one
  namespace; no covenant, no credential; a named issuer obeyed), and one leg in the door's battery: a
  record planted with covenant `feedfacecafebeef` comes back issued in that namespace.
- Still a literal, flagged and not built: `atlas-tui`'s banner prints `covenant: 1512741580b7239b` as
  text (`cmd/atlas-tui/main.go:368`), and the glass's sidebar shows the same string.

**RESTART REQUIRED:** the door carries this once `atlas-mcp.exe` is rebuilt and restarted.

**WHAT GOES RED IF THIS COMES UNPLUGGED:** the `vc` test and the battery leg — the tool hard-coding
the issuer again reds the leg; `IssuerFor` returning a literal reds the test.

### Fixed — the engine names a corrupt flow instead of hiding it (operator, 2026-09-25: "the engine shouldn't hide a corrupt spec, either")

`flow.List` read each spec back and, when one would not parse, skipped it without a word
(`line/internal/flow/flow.go`, `if err != nil { continue }`) — so a corrupt flow was invisible to
`flow_list` until somebody fired it, and a `.json` named outside the name law vanished the same way.
The core's release gate found the skip the first day it read `flows/` (its `flows` check, 2026-09-25);
the door itself said nothing.

- `flow.List` now returns `([]Spec, []Unread, error)`: every flow it could read, and every `.json`
  under `flows/` it could not, each with why — corrupt, or a name no tool can reach. Folded versions
  (`<name>.v<k>.json`) are history and are neither; `Get(name, k)` reaches them.
- `flow_list` prints them under `UNREADABLE — n (not hidden; fix or remove the file)`, beside the
  flows that can be read; a folder holding only unreadable files says `FLOWS — none that can be
  read.` rather than `no flows yet`. Its description says so.
- Proved: `TestListHidesNothing` (flow — a corrupt file, a misnamed file, a folded version that is
  neither), and one leg in the door's own battery (`atlas-mcp --prove`): a corrupt spec planted
  beside a good one is named, and the good one is still listed.

**RESTART REQUIRED:** the door carries this only once `atlas-mcp.exe` is rebuilt and restarted; the
running door is the old build until then.

**WHAT GOES RED IF THIS COMES UNPLUGGED:** `TestListHidesNothing` and the battery leg — and, from
the other side of the seam, the core gate's `flows` check, which refuses the same file by name.

### Fixed — the glass carries a service wire, so the door can be armed (operator, 2026-09-24: "build out the auth")

`webapp/main.go` read `h.ConfigureAuth(true, "", "data/sessions.json")`, with a comment saying *"the
service wire to the door stays empty -- the door is unchanged."* True when it was written, false from
the day `--auth` existed. Every one of the six send-sites guards on `h.service != ""`, so **the glass
sent no `Authorization` header at all** — and arming the door would have answered every page 401:
records, rack, worlds, Version control. The PIN would not have helped; that is a different gate.

It now reads **`ATLAS_SERVICE`** — the same variable the door already falls back to
(`cmd/atlas-mcp/main.go`), so the two agree by reading one place rather than by someone setting two.
An environment variable and not a flag, by RULE 7: a key never rides a command line where `ps` can read
it. Unset is the ordinary case and is exactly the behaviour the glass had before.

The boot line now says **whether** a wire is held, never its value (`team_status`'s shape). A start
with an armed door and no key is the one failure this piece exists to prevent, and without that line it
would look exactly like an ordinary start.

**WHAT GOES RED IF THIS COMES UNPLUGGED** (core RULE 11): `TestTheLockIsSwitchedOnWhereTheGlassStarts`
now asks whether the call CARRIES anything, not only whether it is there — it passed for weeks over a
hard-coded `""`. Proved by reversal: restore the empty wire and it reds naming the 401. A second stroke,
`TestTheGlassSendsItsServiceWireAndOnlyWhenItHasOne`, counts sends against guards across all four
handler files, so an empty wire can never become a bearer of `""`.

One thing found in the building and recorded rather than quietly fixed: the RULE 7 assertion first
grepped the whole file for `auth-service` and went red on the COMMENT explaining the change, which names
the door's flag to say the two read one variable. Prose quotes the thing it is explaining — the same
fault `contains` was narrowed for on 2026-09-12. It asks the IMPORT BLOCK for `"flag"` instead:
structure rather than wording.

Not in this piece, and named so it is not mistaken for done: **P0-13**, RBAC failing open twice
(`tools.go` runs no check when a call names no `actor`; `tenant.go` allows all on an empty policy, and
`DefaultPolicy` ships it empty). It does not block arming — the service wire short-circuits `gateCall`
before any scope check, and no tenant key has ever been minted — but the holds are only as good as it
once one is.

**Arming still needs his hand:** a rebuilt `atlas-webapp.exe` placed, the door restarted with `--auth`,
and `ATLAS_SERVICE` set in both processes' environment.

### Added — and a head per seat (operator, 2026-09-23: "then B underneath it")

`flow_run` takes `voices` beside `voice`: seat → model, applied OVER it.

**THE TWO REACH DIFFERENT THINGS, and the asymmetry is the shape of what is underneath.** `voice` is ONE
model, which is what an `ask`, `prompt`, `seat` or `memory` node measures against — so it defaults
every node that pins none, and is the whole roster to a `run` node. `voices` only means anything to a
`run` node: that is the one kind reaching the estate's ROSTER rather than a single voice, so there is
nothing for a per-seat map to say to the others, and handing one to them would quietly pick an entry.

**Why it exists:** a run-level voice answers "is this flow better on that model". It cannot answer "does
the STEWARD raise a flag where it used to announce", because it moves the whole roster and the answer
becomes a fact about two changes at once.

### Changed — the run's head is a type, not a string

`flow.Head{Voice, Voices}` and `engine.Head{Model, Voices}` replace the bare `voice string` that landed
in `1091ae1` an hour earlier, before anything depended on it. **They must move together**: a caller
passing one and forgetting the other would fire a parity that varied the whole roster while its record
said it varied one seat, and no signature with two loose strings beside each other stays honest for
long. The two packages keep SEPARATE types on purpose — `Voice` means "one model" to a flow and "every
seat" to the council — and `councilEngine.Turn` is the one place that knows both, so neither package
has to learn the other's vocabulary.

- **`internal/flow/run.go`:** `Head` + `Named()`; `RunOn` takes it; `tidy` drops blanks so a head of
  whitespace cannot make a run's record say it was headed; `onHead`/`startHead` replace
  `onVoice`/`startVoice`; `Status` renders a `seats:` line and `Compare` names both sides' seats, in
  SEAT order rather than Go's map order — two runs' lines have to be readable against each other.
- **`internal/tools/tools.go`:** `flowVoices` reads the map off a call as an object or as an object in
  a string, exactly as `inputs` is read, and REFUSES every other shape rather than dropping it; the
  council's `WithHead` binds a copy; `flow_run` declares `voices?`.
- **`internal/engine/engine.go`:** `Head` on the objective row — `model`, `voices`, or neither.
- **`docs/SPEC_CONTROL_CENTER.md`:** the `flow_run` row carries `voices?`.

Measured on a mirror: every package green but the seven long-path git strokes, whose red set is
identical at HEAD on the same mirror. Reversed five ways — the council dropping the map, a seat map
reaching nodes that measure one voice, a replay dropping it, a bad shape accepted rather than refused,
and the seats rendered in map order — each reds its own stroke and no other.

### Added — the head belongs to the run, not the spec (operator, 2026-09-23: "let's do C first, then B underneath it")

`flow_run` takes an optional `voice`: the model THIS RUN answers on.

**WHY IT CANNOT LIVE IN THE SPEC.** A model could only be named per NODE, so measuring one flow on two
models meant folding two specs -- and two specs stop being one experiment the moment either is edited,
which is precisely the parity the pair was folded to measure. The spec stays model-agnostic and the
head is named when the run is FIRED, so **the same questions on two heads is a comparison by
construction** rather than by a hand's promise that it kept the two in step.

**IT REACHES EVERY NODE THAT PINS NONE, AND ONLY THOSE.** A node carrying its own `voice` keeps it --
the spec said that one aloud. Everything else answers on the run's head, the council included. Without
that last part a "parity run" would have varied only its `run` nodes while its `ask`, `prompt`, `seat`
and `memory` nodes went on reaching the declared targets, and reported itself as a whole-flow
comparison anyway.

**AND IT GOES IN THE START LINE, which is what makes the number mean anything.** `flow_status` names
the head that answered. `flow_resume` and `flow_replay` read it back off the run's own record instead
of taking one from the caller -- a gate is FOR walking away, and coming back to finish a parity run's
second half on the declared targets would leave one run measuring two models with nothing saying so.
`flow_compare` names both heads when they differ, and stays quiet when they do not, because two runs
on one head is the model's own variance and is a different reading.

Per-seat heads, so a parity can vary ONE seat, is the narrower ruling and is his (B).

- **`internal/flow/run.go`:** `RunOn(home, eng, s, inputs, voice)`; `Run` is now the declared-targets
  door and delegates with `""`, so the thirty-odd existing call sites are unchanged and the two
  cannot drift. `onVoice` binds the head onto the engine through an optional `WithVoice` interface --
  the same shape `Resume` already used for `Ready()` -- and `startVoice` reads it back off a start
  line. `runFrom` carries it and defaults each node's empty `Voice` to it; `Status` and `Compare`
  render it.
- **`internal/tools/tools.go`:** `councilEngine` carries `voice` and hands it to the engine's
  objective row; `WithVoice` binds a COPY, so the registered engine is untouched and two flows in
  flight cannot cross heads. `flow_run` declares `voice?` on the wire; `flow_resume` and `flow_replay`
  deliberately do not.
- **`internal/engine/engine.go`:** `Run(objective, feed, method, model, sink)` -- the tag rides on the
  objective row that Manjuel's headless door now reads. `runstream.go` and `run_start` pass `""`.
- **`docs/SPEC_CONTROL_CENTER.md`:** the `flow_run` row carries `voice?` and what it does.

Measured on a mirror: every package green except the seven git strokes that are red at HEAD on the
same mirror for the same reason (`Filename too long` under a long scratch path) -- byte-identical red
sets, baseline and working. Reversed seven ways; each switch-off reds its own stroke and no other.

---

## [v0.1.7] — 2026-09-23 07:30 (tag on 063a152)

**THE MARK IS CUT**, 2026-09-23 07:30, annotated `EVERY PIN IN STEP`, on `063a152`, through the
door's `git_tag`. The heading read `[Unreleased]` until the mark existed; the words under it are
as they were written.

**AND IT IS THE FIRST ATLAS MARK CUT OVER BINARIES THAT AGREE WITH IT.** `v0.1.6` was cut on
2026-09-18 from the Cut button with only its root `VERSION` bumped: its binaries answered 0.1.5,
the spine's version strokes had been red unseen since `0c65afc` because the push check skipped
cargo, and its release workflow never built. All of that was fixed on 2026-09-22 in `3a07dfd`,
and this mark is what that fix was for -- the door read `VERSION` AT `063a152`, found 0.1.7, and
`staleStamps` found no stamp disagreeing.

**WHAT THE GATE SAID** (2026-09-23, on a mirror of this tree after the bump): `.\version.ps1
sync` — All 11 pins in sync at 0.1.7 · gofmt clean in both modules · `go build ./...` and `go
test ./...` green in `line` and in `webapp` · `python tests/prove.py --check` PROVEN. The
operator's terminal remains the proof; a mirror is the hand's own check.

**WHAT IS IN IT.** The version control fix itself, the door that opens a world whose sitting's
process is gone, flows that report truthfully, and the two reads at the door that handed out this
estate's keys.

### 0.1.7: every pin and every claim moved together

His word, 2026-09-23: *"bump the versions by 1 on both core and atlas"*. No code moved; the pins
and the docs only, so nothing needs rebuilding or restarting.

- **All eleven pins to 0.1.7** by `.\version.ps1 set 0.1.7` -- the eight `VERSION` files, plus
  `Cargo.toml`, `core/src/version.rs` and `Cargo.lock`. `.\version.ps1 sync` answers "All 11 pins
  in sync: 0.1.7". `Cargo.lock` moves with them because of the fix landed 2026-09-22 in
  `3a07dfd`; before that a lock left at the old number made `cargo build --locked` fail with
  exit 101, which is how `v0.1.6` came to be cut over binaries that answered 0.1.5.
- **24 version CLAIMS moved** in `AGENTS.md`, `README.md`, `docs/ACCEPTANCE.md`,
  `docs/OLLAMA_PROVER.md`, `docs/PIPELINES.md`, `docs/WORKFLOWS.md` and
  `tests/e2e/E2E_SCENARIOS.md` -- the lines that assert what a binary PRINTS, which `version.ps1`
  deliberately leaves to a hand. **Not swept, and each for a reason:** this CHANGELOG (history --
  moving it would rewrite the record), `LAUNCH_PLAN.md` (a world's, closed until he points at
  it), `docs/SPEC_CONTROL_CENTER.md` (a roadmap OF numbers, not a claim about now) and
  `DELIVERABLE.md`'s "0.1.6 is already cut", which is true until the new mark exists.
- **`docs/ACCEPTANCE.md` said the pins were 8; they are 11.** `Cargo.toml` and
  `core/src/version.rs` were pins the count never learned, and `Cargo.lock` became one on
  2026-09-22. The line now names all eleven and says `.\version.ps1 sync` is what proves them,
  which is what CI asks.

**The mark is NOT cut here.** Its gate is `python tests/prove.py --check`, both Go modules green,
gofmt clean and the door's battery -- and the core's own gate waits on a live standup. This
CHANGELOG stays under `[Unreleased]` until the mark exists.

### The door hands out what changed, not whatever is on the disk: `git_diff` and `read_plan` both served `.env`

His word, 2026-09-22: *"4. Code safety"* -- the handoff's fourth piece, from the review the same
day. **RESTART REQUIRED, THE DOOR:** `line/internal/tools/gitstate.go` and
`line/internal/tenant/tenant.go` moved. No sitting was open.

**WHAT IT WAS.** Two read tools, one fault, and in both cases the path jail was working: the file
asked for is INSIDE the world, so a jail about where a path LANDS could never have caught it.
Any program on this computer can reach :8090.

    git_diff             `git diff` says nothing about a file git has never seen, so the tool
                         fell through to serving an untracked file WHOLE. That is right for a
                         new file on its way to a save and wrong for everything a world keeps
                         out of its history on purpose. `git_diff file=.env` came back as the
                         estate's keys, in full.

    read_plan            `Tenant.PlanPath` ended by trying the caller's own `which` as a path
                         inside Home -- and `resolve` hands an ABSOLUTE path straight back
                         unchanged. So `read_plan which=.env` was the estate's keys again, an
                         absolute path was any file on this machine, and `../` walked out of
                         the world entirely.

**WHAT IT IS NOW.**

- **`line/internal/tools/gitstate.go`:** before the untracked fallback, `.env` and `.env.*` are
  refused BY NAME whether or not a world remembered to ignore them (RULE 7: keys are never
  printed), and anything `git check-ignore` reports as ignored is refused with it. git's own
  ignore rules are the world's own statement of what is not part of the work, which is exactly
  the question `git_diff` asks -- so `worlds/`, `vault/`, data and logs come with it for free,
  and no second list has to be kept in step.
- **`line/internal/tenant/tenant.go`:** a name that arrives on a CALL now resolves only under
  `plans/`, which is what the tool has always said it serves. The manifest's own `plans` map and
  the five built-in names stay trusted -- they are configuration the operator wrote, not input.
  An absolute path, a drive letter, `.`, `..`, an escape from Home and a directory all answer
  `""`, which is the honest denial the function already had for an unknown plan. Measured first:
  no caller in this repository depended on the bare in-Home fallback.

Strokes: `TestGitDiffWillNotServeWhatTheWorldKeepsOutOfItsHistory` (five files refused, none
handing out its contents, an unignored `.env` still refused by name) and
`TestAPlanNameCannotBecomeAPath` (ten names refused, four ways that must not fire). Both ways
proven: the whole Go battery green, and each fix switched off in turn on a mirror turns its own
stroke red -- the tenant one printing the hole verbatim, `read_plan ".env"` resolving to the
ground's key file.

The ways that must not fire are stroked too: an ordinary untracked file is still served whole, a
tracked file's real change still diffs, and a plan under `plans/`, the built-in `road` and a
manifest-mapped name all still resolve.

### Flows report truthfully: a turn that did not deliver, a check that stays failed, every gate resumable, and a lock per world

His word, 2026-09-22: *"2. Flows report truthfully"* -- the handoff's second piece, from the review
the same day. **RESTART REQUIRED, THE DOOR:** `line/internal/flow/` and `line/internal/tools/`
moved. No sitting was open.

**WHAT IT WAS.** Six findings, each read in the code before it was touched:

    a turn that did      councilEngine.Turn returned the final event's text whatever kind it
    not deliver          was. `refused` (the law gate), `aborted`, `cancelled`, `unreachable`
                         and `command` -- the last being a turn that ran no pipeline at all,
                         which is also what a runtime error leaves, carrying the objective's
                         own words back -- were all handed to the flow as the node's OUTPUT.
                         The run log wrote them "ok", so a flow could walk its whole happy
                         path having done nothing, and the node's recorded answer could be
                         its own question
    a check that         a node line carried `status`, which is "ok" for an eval that
    came back a pass     ANSWERED, pass or fail alike. Resume rebuilt `pass` as true for
                         every ok line, so a check that FAILED came back from a gate as one
                         that passed: the fail branch went silent and the pass branch fired
    the second gate      any `resumed` line made the whole run "not paused", so a flow with
                         two gates could never pass the second -- while flow_status and
                         flow_runs went on showing PAUSED and the tool went on telling the
                         operator to resume it
    an answered gate     a resumed gate fires in memory and writes no node line, so the NEXT
    forgotten            resume walked back to it and paused on it again. Found by the
                         two-gate stroke below, not by reading
    a cancel called      flow_cancel ends the run's context, and both a cancel and a spent
    out of time          budget arrived as the same dead context: the verdict said
                         OUT_OF_TIME over a run the operator had stopped himself. The cancel
                         could not reach a turn in flight at all -- `Turn` ignored the
                         context, so a cancelled flow drove the council to the end of the
                         turn and only the NEXT node saw it
    a resume with no     `continue` walked into the next `run` node, which cannot open an
    engine burned the    engine; the node errored, the run took a terminal line, and a run
    run                  with one never resumes again. So a gate answered after the engine's
                         thirty-minute idle close -- which is exactly the walk-away a gate is
                         FOR -- destroyed the work it guarded
    one flow froze       flow_run, flow_resume and flow_replay held `askLock` -- one mutex
    every world          across every tenant -- for a WHOLE run. The glass's chat, every
                         prompt run, every key mint and every other world's flow waited on
                         it. SPEC_CONTROL_CENTER 4.6 asks for exactly this, in its own words

**THE CHANGE.**

    deliveryOf        only `delivery` is a turn's answer; every other terminal event is a
                      refusal that names the kind and quotes what it said. The kind IS the
                      verdict, so the flow's run node fails instead of recording nothing
                      as something
    the node line     carries `pass`, the node's own outcome, beside `status`
    Resume            reads that outcome back (a run logged before today is read from the
                      eval's own answer, and nothing in the record is rewritten); reads a
                      `resumed` line as the gate it answered standing fired; and judges
                      only the pause that is STILL OPEN, so every gate resumes once
    Ready             flow.Resume asks the engine, through an optional interface, whether
                      one is standing for the `run` nodes still to fire -- BEFORE it
                      appends anything, so a refusal leaves the run at its gate
    verdictFor        a cancelled context is STOPPED, a spent budget is OUT_OF_TIME
    Turn              watches the run's context and sends the engine a cancel when it
                      dies: Ctrl-C is what the engine understands
    flowLock          one mutex per world for flows. The rack stays one queue -- the
                      council engine's Ask, RunPrompt and SeatAsk take askLock per CALL --
                      so a flow's model calls still queue behind rack_ask and let go
                      between nodes. Turn is not among them: run_start does not take
                      askLock either, and one engine per world is already the invariant

**STROKES, +6 in `internal/flow` and +2 in `internal/tools`.** A check that failed comes back
failed and its fail branch fires; one that passed comes back passed; a run logged before the
outcome was written down is read from the eval's own answer; a flow with two gates pauses,
resumes, pauses and completes, and refuses a third resume; a cancelled run is STOPPED and a spent
budget is still OUT_OF_TIME; a resume with no engine is refused, writes NOTHING into the run, and
the same gate resumes once an engine stands. In `tools`: only a delivery is an answer (every other
terminal kind refuses, naming itself, and returns no text), and the flow lock is per world, is not
askLock, and is the one the three flow tools actually take -- read out of the source, the way
`internal/flow`'s own no-finish stroke reads it.

**PROVEN BY REVERSAL, seven undos on a scratch copy**, each turning its own stroke red and put
back green: the delivery check off; the outcome not written down; any prior resume closing the
run; an answered gate leaving no record; a cancel called out of time; the engine check skipped;
and flow_run back on the global lock. `gofmt`, `go vet` (Windows and GOOS=linux) clean; in `line`,
every package green but the six `git_tag` push strokes the scratch path always breaks.

**THE NEW DOOR IS NOT IN PLACE.** Built in scratch from sources byte-identical to these: sha256
4430227a4544eb3e, 12,166,144 bytes. Placing it over `line/atlas-mcp.exe` and restarting the door
(pid 5712) is his allowance. Until then flows run as they did.

**Named, not fixed.**

    the other tools   askLock is still one mutex for `remember`, `rack_ask`, the chat and the
                      key store. SPEC_CONTROL_CENTER 4.6 wants it keyed by tenant home for
                      all of them; this piece took the flows, which were the ones holding it
                      for minutes at a time
    a cancelled node  is logged `fail` with its error, and the run is STOPPED. The node did
                      not fail at its work, and the line says what it was doing when the hand
                      stopped it

### A world a crash used to lock: the door reads the pid on the open line

His word, 2026-09-22: *"1. Survives crashes"* -- the handoff's first piece. The engine's half is in
the core's CHANGELOG. **RESTART REQUIRED, THE DOOR:** `line/internal/engine/` moved, so none of this
is live until the door is rebuilt, placed and restarted. No sitting was open.

**WHAT IT WAS.** `Open` refused any world whose ledger's last line had no `ended`, and never read
the `pid` every line has carried since 2026-09-17. An engine that crashed or was killed left a line
that refused its world to every Boot, until someone opened a REPL there to reap it -- and the
headless engine this door spawns never reaped.

**THE CHANGE.**

    SittingOrphaned   new: the open line's pid, and whether that process is PROVABLY gone. No
                      pid, an unreadable ledger, a live process, or one this door may not ask
                      about all answer "not orphaned" -- the core's own rule (seatlog._alive)
    Open              refuses as before unless the open line is orphaned. Then it opens the
                      world, says so in the door's notes, and the engine it spawns reaps the
                      line as it starts (the core's serve.main, with the REPL's own reaper).
                      The door still never writes the ledger
    processgone_*     the kernel asked, never signalled: OpenProcess and the exit code on
                      Windows (ERROR_INVALID_PARAMETER is "no such process"; access denied is
                      not a death), signal 0 and ESRCH elsewhere. The first platform split in
                      THE LINE; GOOS=linux builds and vets
    lastSitting       the ledger read, shared by both; SittingOpen answers exactly as before

**STROKES, +2 in `internal/engine`.** `TestAnOrphanedSittingIsOnlyOneWhoseProcessIsProvablyGone`: an
open line whose process has come and gone is orphaned and names its pid; a live pid, no pid, a
closed line, an unreadable last line and no ledger are not. `TestOpenOpensAWorldWhoseSittingsProcessIsGone`:
the door opens such a world on a real (stub) engine, and the same line with a live process behind
it is still refused by name.

**PROVEN BY REVERSAL, on a scratch copy.** With the liveness answer forced to "alive", both go red,
the second with the exact old refusal; put back, green. `gofmt` and `go vet` clean, on Windows and
for GOOS=linux; in `line`, every package green but the six `git_tag` push strokes the scratch path
always breaks.

**THE NEW DOOR IS NOT IN PLACE.** Built in scratch from sources byte-identical to these: sha256
025ed38d99fdda4f, 12,151,296 bytes, `--version` 0.1.6. Placing it over `line/atlas-mcp.exe` and
restarting the door (pid 8116) is his allowance. Until then a crashed engine's world is refused as
before.

**AND PLACED, THEN SAVED AND SENT, 2026-09-22, on his word** ("Place the door"; "Both
repositories"). The door (pid 8116) was checked by pid and path, stopped by pid and replaced with
this build, which hashes as built; the build it ran is kept in the hand's scratch. It runs as **pid
5712** on 127.0.0.1:8090 alone, on the command line it had: `/health` 0.1.6, 82 tools. The
`atlas-mcp.exe` from `Desktop\Archive` (pid 23164) was left alone, and the glass was not touched,
so his session held. Saved through the council in sitting 263 (12:28-12:31, four runs, closed with
its toll): atlas `12c8574` (`3a07dfd..12c8574`), these five files, `main` alone, GitHub level.
Written after the save; it rides with the next.

**Named, not fixed.**

    env_list          still says "sat in elsewhere" for an orphaned line, though Boot now opens
                      it. One line in tools.go, left for a piece that reads that file whole
    the spec          SPEC_CONTROL_CENTER's "on its own restart it must reap orphans" is met from
                      the other side: an engine whose door dies closes itself on the hang-up --
                      now even with a question pending -- and a dead one's line no longer refuses

### Every version stamp in step: the stale 0.1.6, the lock the bump tool forgot, and a push check that runs the whole battery

His word, 2026-09-22: *"the version control for atlas sounds like a fairly simple fix, let's get
that knocked out"* -- after the review found `v0.1.6` cut with only its root `VERSION` bumped.
**RESTART REQUIRED, THE DOOR:** `line/internal/tools/gitctl.go` moved, so the new mark check is
not live until the door is rebuilt, placed and restarted. The glass's Go did not move, but its
`VERSION` stamp did, so its `/api/health` says 0.1.5 until it is rebuilt too. No sitting was open.

**WHAT IT WAS.** Measured on scratch copies of this tree, never on the tree itself:

    the drift         v0.1.6 (0c65afc, 2026-09-18) was cut from the panel's Cut button with the
                      root VERSION at 0.1.6 and the other nine pins -- seven VERSION files,
                      Cargo.toml, core/src/version.rs -- still at 0.1.5. Every binary built
                      from that mark answers 0.1.5
    red, unseen       at that commit, and on main until today, `cargo test --workspace` is red
                      twice: the spine's own `version-cross` stroke (`atlas --prove`) and
                      core's `no_stone_tag_and_no_whitespace`. prove.yml ran THE BALL as
                      `--check`, "the fast legs only (no cargo)", so every push stayed green
    no release        release.yml's pin step (`version.ps1 sync`) exits 1 on that tree with
                      nine MISMATCH lines, before anything is built. v0.1.5's run had died at
                      its Go tests (the workflow's own note), and the fix landed six minutes
                      after that tag. So no release has been built from either mark. GitHub's
                      run page was not opened: that needs his word
    the door let it   tagCut judged a mark against the root VERSION alone
    the tool forgot   `version.ps1 set` moved Cargo.toml and not Cargo.lock, and every
    the lock          `cargo build --locked` after it refuses ("cannot update the lock file ...
                      because --locked was passed") -- the first step of prove.yml and of
                      release.yml. Measured both ways on a scratch copy. And `sync` asked only
                      whether the wanted text appeared ANYWHERE in a pin's file, which a lock
                      holding three crates satisfies with one
    two harnesses     tests/e2e/ollama_prover.py and tests/workflows/wf_operator.py said
                      VERSION = "0.1.3", so the live prover's S1 version checks and S9-1 failed
                      on the number alone

**THE CHANGE.**

    every stamp at    `version.ps1 set 0.1.6`, run on the ground with the fixed tool: seven
    0.1.6             VERSION files, Cargo.toml, core/src/version.rs, and Cargo.lock's three
                      crates. `version.ps1 sync`: "All 11 pins in sync: 0.1.6"
    version.ps1       Cargo.lock is a pin (its three crates, one pattern), and `sync` finds
                      every match of every pin and holds each one to the version
    the door          tagCut refuses a mark while any file called VERSION anywhere in the tree
                      at that commit, or the root Cargo.toml's [package] or
                      [workspace.package] version, says another number, and names each stamp
                      and what it says. Read out of git at that commit, never off the disk; a
                      stamp git cannot read is named, not passed. The core has no VERSION file
                      and no Cargo.toml, so nothing changes for its marks
    prove.yml         `version.ps1 sync` first, on every push; then THE BALL whole -- cargo test
                      and the door's own battery with it -- where it was `--check`. The comment
                      above that line had promised cargo all along
    release.yml       its pin step is "Every version pin agrees", named without a count
    release.ps1       its last lines said nothing builds a tag. release.yml has since
                      2026-09-12, and it stops at a DRAFT; that is what it says now
    the harnesses     both read atlas's root VERSION instead of carrying a number
    the docs          the version each binary prints, 0.1.5 -> 0.1.6, in the places 0.1.5's own
                      bump moved (AGENTS, ACCEPTANCE, E2E_SCENARIOS, OLLAMA_PROVER, PIPELINES,
                      WORKFLOWS) and in README; dated measurements left as they were.
                      THE_ROAD's version line names v0.1.6 and what happened to it; PIPELINES,
                      DELIVERABLE and PROVING say what the two workflows run now; DELIVERABLE's
                      release line says `patch`, because 0.1.6 is taken

**STROKES, +2 in `internal/tools` (`gitctl_test.go`).** `TestAMarkWaitsForEveryVersionStamp`: a
world whose root VERSION says 0.1.6 while `line/cmd/atlas-mcp/VERSION` and Cargo.toml say 0.1.5
is refused v0.1.6, naming both stamps, and no mark exists after; once both move, the same name
lands. `TestCargoVersionReadsOnlyTheCratesOwnNumber`: a workspace's or a crate's own number,
double- or single-quoted, and nothing for a member that inherits it or for a dependency.

**PROVEN BY REVERSAL, on a scratch copy.** With the new check switched off, the first stroke goes
red; put back, green. Its name is short on purpose: under its first, longer name the scratch
folder's path pushed git past MAX_PATH, and a stamp came back "could not be read" -- the guard
failing shut, as written, over a fault of the machine.

**ON SCRATCH COPIES, BEFORE AND AFTER** (no .git, no target/; Go's cache and temp in scratch):

    version.ps1 sync       before: nine MISMATCH, exit 1    after: "All 11 pins in sync: 0.1.6",
                                                                   under Windows PowerShell and pwsh
    cargo build --locked   Cargo.toml at 0.1.6 with the lock at 0.1.5: refused. Both moved: builds
    cargo test             before: red, the two version     after: ABSENT, 55 strokes held --
                           strokes                                 stopped only by the oracle legs
    go test, line          the same six git_tag push strokes red before and after ("Filename too
                           long", the scratch path); every other package ok
    THE BALL, whole,       21 held, 15 absent, 1 broke -- the one broke being those six strokes
    spine built first
    gofmt -l, go vet       clean

**THE NEW BUILDS ARE NOT IN PLACE.** Both wait in the hand's scratch, built from sources
byte-identical to these:

- the door: sha256 993170cb8efd0d4e, 12,147,712 bytes; `--version` answers 0.1.6;
- the glass: sha256 7616272baea14ee2, 11,034,624 bytes; nothing moved in it but its VERSION.

Placing them and restarting is his allowance. Until then the running door judges a mark by the
root VERSION alone, and both answer 0.1.5. None of this reaches GitHub, including the two
workflows, until atlas is saved and sent.

**SAVED AND SENT, THEN PLACED, 2026-09-22, on his word** (asked, and answered: "Save and send";
"Door and glass"). Saved through the council in sitting 262 (09:41-09:43, two runs, closed with
its toll): atlas `3a07dfd` (`d94c1e9..3a07dfd`), these 29 files and nothing else, `main` alone,
and GitHub's `main` matches. Then, with no sitting open and no engine standing, the glass (pid
27600) and the door (pid 26876) were each checked by pid and path, stopped by pid, and replaced
with these builds, which hash as built. The builds they replaced (door 3c0ce9c7, glass 91a6875d)
are kept in the hand's scratch.

    the door    pid 8116, on the command line it had, from `line\`: 127.0.0.1:8090 alone,
                82 tools, `/health` says 0.1.6
    the glass   pid 24548, from `webapp\`: 127.0.0.1:8091 alone, gate on, his lock found,
                `/api/health` says 0.1.6

The `atlas-mcp.exe` running from `Desktop\Archive` (pid 23164) was left alone. His Dashboard
opens on the lock screen; the PIN is his. This paragraph is written after the save and rides with
the next one.

**Named, not fixed.**

    the dead mark     v0.1.6 stays where it is: a mark is never moved, and GitHub may hold it.
                      The next release is the next number, cut and sent by him
    the bump          moving the stamps is still `version.ps1` on a terminal; the glass has no
                      button for it
    the core          its marks are unchanged by this. Its own pins are pyproject.toml and
                      manjuel/__init__.py, and the door does not read the second

### The maker's projects on the glass: the list, the page, and the words to pick one up

His word, 2026-09-21: *"go on piece 2"* -- the maker's piece 2 in the core's BUILDPATH, "THE
PAGE ON THE GLASS -- a preview pane beside the run, and a project list to pick one from, which
is also where a project is put DOWN". The engine's half is in the core's CHANGELOG. **RESTART
REQUIRED, BOTH: the door and the glass.** `line/internal/tools/` and `webapp/` moved, and the
glass embeds its pages, so none of this is live until both binaries are rebuilt, placed and
restarted. No sitting was open.

**WHAT IT WAS.** The maker saves each thing it makes as its own git repository under the core's
`projects/`. Nothing on the glass could see one: no list, no page, and no way to put a project
down short of closing the sitting.

**THE CHANGE.**

    projects          the door's 82nd tool, internal/tools/projects.go, new, Writes: false.
                      `list`: every folder under projects/ that holds a .git of its own,
                      each with its versions read off its own log (the maker's "Version N:"
                      taken off the note), newest first. `page`: one project's index.html
                      as it stands, or as version N was, whole; refused past a megabyte.
                      Every git call names the project's own .git and work tree outright,
                      with a ten-second deadline (ESTATE LAW 7)
    a name, not       `^[a-z0-9][a-z0-9-]{0,63}$`, the maker's own shape, checked before any
    a path            path is built. And every step is asked, with Lstat, whether it is a
                      PLAIN folder, so a project that is a link is refused by name. Measured
                      first: Go 1.26's EvalSymlinks does not follow a Windows junction, so
                      the first cut's inside-check let a junction's .git through. A
                      junction needs no privilege to make
    the page route    GET /api/projects/{name}/page?v=N, handlers/projects.go, new. The name
                      and version are checked before the door is asked. A door refusal is
                      said, 404, and no page is served. The page goes out as text/html
                      under Content-Security-Policy: sandbox allow-scripts allow-modals --
                      NEVER allow-same-origin -- with default-src 'none', connect-src
                      'none', frame-ancestors 'self', nosniff, no-store and no referrer
    the header is     the frame carries no sandbox ATTRIBUTE: the app's own browser pane
    the wall          refuses any frame that does (net::ERR_BLOCKED_BY_CLIENT, measured). The
                      header alone, measured in that pane on a scratch server serving the
                      glass's exact policy: the page's script ran and its canvas drew; its
                      origin was "null"; reading the parent, the parent's DOM, the cookie
                      and localStorage threw SecurityError; fetch threw
    the card          static/js/projects.js, new: Projects, on the Dashboard under the run.
                      It shows every project with its versions and when it last moved, the
                      page in a frame, a version picker and "Open in its own tab". "Work on
                      this" and "Put it down" put the engine's own words in the box and run
                      them ("work on the <name> project", "put the project down"), so the
                      turn is in the record like any other. The delivery's `project` is how
                      the card learns what is in hand. It is kept per world in the settings
                      store and keyed to the session, so a closed engine's project is never
                      shown as held
    not traces        `projects` joins backgroundReads. The card reads the list when the
                      Dashboard opens and after every turn, and those reads are answered and
                      not kept. That is this hand's call, made in the shape of his D1 ruling,
                      and it is one line to undo. A page view goes straight to the door and
                      is not kept either

**STROKES, 66 -> 73 in `internal/tools` and 23 -> 27 in `webapp/handlers`.** `projects_test.go`
(7) is hermetic. Each stroke runs on a temp home that is itself a git repository, so a project
with no .git of its own would read as the HOME's history if anything let git walk up. It proves:

- the list is read off each project's own history and sorted newest first, dated so that
  alphabetical order is NOT the answer; a world with none lists none;
- a page is served as it stands and as version 1 was; a version past the last is refused,
  naming the range;
- the page is never HTML-escaped on the way;
- a name that is a path is refused; a project that is a link is refused, as a symlink where the
  machine allows one and otherwise a junction by `mklink /J`;
- an unknown action is refused by name.

`projects` is in the registry's non-writing list, and `TestEveryCoreToolStandsAlone` calls it.
`handlers/projects_test.go` (4): the page is served under the sandbox header; a bad name or
version never reaches the door; a door refusal is said and serves no page; and the list is a
background read.

**PROVEN BY REVERSAL on scratch copies: twelve undos, each turning red.**

- The door, seven: a link walked into (the first cut's missing plain-folder test); a folder with
  no history of its own read as a project; a name not held to the maker's shape; the list
  unsorted; "Version N:" left on the note; the page escaped on its way out; a version past the
  last served as the last.
- The glass, five: the page given the glass's origin; the header dropped; a bad name, and a bad
  version, still asking the door; the list kept as a trace.

One undo was wrong at first and is corrected. The bad-version undo replaced the check with
`if false`, which left `err` unused, so the package did not compile, and the script read "no
failing stroke" as green. Rewritten to compile, it goes red. With every undo put back:

- `gofmt`, `go vet` and `node --check` clean;
- every package in `line` green except six `git_tag` push strokes. They cannot push under the
  scratch folder's long path: git fails writing an object with "Filename too long", past
  Windows' path limit. The same six fail the same way at HEAD, so they measure the scratch,
  not this piece;
- `webapp`: db 13, handlers 27, server 6.

**AND ON REAL BINARIES, IN SCRATCH: 20 of 20.** Built from these sources: a scratch door on :8098
carrying `proof`, a mirror of the core with two projects in the maker's own format, and a scratch
glass on :8097 with a data folder of its own, aimed at that door and nowhere else.

- The list is the two projects with the maker's notes, and left no trace.
- The page as it stands and as version 1 was: each 200, under the sandbox header, as HTML, not
  cached. Version 9 is refused by the door in its own words, a name that is a path is refused,
  and with no session there is no page (401).
- An engine opened on `proof`. "work on the snake-game project" was answered by the engine with
  no seat, and the delivery named `snake-game`. "put the project down" was answered, and the
  delivery named none.
- The sitting was closed and tolled, and both projects were on disk exactly as they were.

His door and glass kept running throughout and were not touched. The scratch pair was stopped by
pid. The sittings opened were the mirror's, and his ledger gained nothing.

**THE NEW BUILDS ARE NOT IN PLACE.** Both wait in the hand's scratch. They were built after the
last source line moved, from sources byte-identical to these:

- the door: sha256 3c0ce9c774622738, 12,138,496 bytes;
- the glass: sha256 91a6875d9acb8ed0, 11,034,112 bytes.

Placing them over `line/atlas-mcp.exe` and `webapp/atlas-webapp.exe` and restarting both is his
allowance. Until then the door serves 81 tools and the Dashboard has no Projects card. And the
card's buttons speak to an engine: an engine started before the core's half of this piece does
not know the words.

**AND PLACED, 2026-09-22, on his word: "place them and restart the door and the glass".** No
sitting was open and no engine stood (`/run/state` open false, nothing running under the door).
The glass (pid 7160) and the door (pid 22844) were each checked by pid and path, then stopped by
pid alone. A second `atlas-mcp.exe` was running from `Desktop\Archive`, outside the ground and
not this hand's; it was left alone. Both builds were copied into place and hash as built. The
builds they replaced are kept in the hand's scratch (door ea98a7a7, glass 62b7dbb1).

    the door    pid 26876, started on the command line it had, from `line\`: 127.0.0.1:8090
                alone, carrying research and atlas, 82 tools with `projects` among them.
                Asked on the ground, `projects` lists snake-game with its three versions and
                serves version 1's page
    the glass   pid 27600, from `webapp\`: 127.0.0.1:8091 alone, gate on, and his lock found
                (a name set, no fresh setup offered). `projects.js` is served byte-identical
                to disk, and a page asked with no session is refused, 401
    the page    his Dashboard opened in the app's pane on the lock screen, with the Projects
                card behind it. The PIN was his to type, and he typed it. Signed in, the card
                listed snake-game with its three versions and their notes, none in hand, and
                the frame drew the game's board (the page answered 200). Nothing was clicked

The door's own notes went to `line/mcp.err.log` this start: stdout and stderr are kept apart,
the arrangement this repository's `.gitignore` names, and `mcp.log` is the empty stdout.

**Named, not fixed.**

    the storage      a page in the frame has an origin of its own, so its localStorage
                     throws, and a game there forgets its high score. Opened from
                     projects\<name>\ in a browser, it keeps it. Giving the frame the
                     glass's origin would give the page the glass's cookie
    README's line    the opening sentence still says ATLAS treats every tool call as a
                     record; the exception paragraph under it now names the Projects card's
                     two reads

### The lock: one user, one PIN, this computer only

His words, 2026-09-21, in order: *"Simple login system for now, user/pin to start"*; asked
who may open the glass, *"This PC only"*; *"I like this idea of multi-roles, all working
under a single user"*; and *"Let's make the thing at least semi-secure"*. **RESTART
REQUIRED** -- the glass's Go and its embedded pages moved, so none of this is live until the
new binary is placed and the glass restarted. The door was not touched. No sitting was open.

**WHAT IT WAS.** The glass listened on `":" + port` -- every address the machine has -- and
its gate had never once been closed: `ConfigureAuth` had no caller (the finding in "The glass
had one test function in 2,800 lines", below). Every page, and every tool call the glass
passes to the door, answered whoever reached it.

**THE CHANGE.**

    the lock          handlers/lock.go, new. The first time the glass opens it asks the
                      person at this computer for a name and a PIN of 4 to 8 digits; after
                      that it is a lock screen. The PIN is never stored: PBKDF2-SHA256 over
                      a random 16-byte salt, 600,000 rounds, in data/user.json (0600,
                      written beside and renamed over). Five wrong PINs in a row close it
                      for a minute, and the right PIN waits the minute out too. Setup is
                      taken only from this computer -- judged by the connection, never by
                      a header -- and only while nobody is set up. A user file that cannot
                      be read is said, and never taken for "nobody is set up": that would
                      hand the lock to whoever opened the page next
    three faces       GET /api/lock, POST /api/setup, POST /api/unlock -- open through the
                      gate, as /api/login is, because the lock screen must ask and answer
                      before anyone is signed in
    the operator's    a PIN session names no tenant, which is how every face here already
    session           reads "see and name everything": the glass as it was with the gate
                      off, behind a PIN. rpcCallAs, scopeProject, CallTool and StreamChat
                      pin a project only when a session names one, so the page's project
                      switch works as before and a key's session is still pinned to its
                      tenant. /api/me says who is signed in
    the cookie        SameSite Strict, was Lax; HttpOnly as before. Lax still sends the
                      cookie when another site links the browser here, and a GET such as
                      /api/council/stream?objective=... runs a turn
    this PC only      the glass listens on 127.0.0.1 and nowhere else (server.go, addr)
    switched on       main.go closes the gate at every start and names the lock. The
                      service wire to the door stays empty
    the page          static/js/lock.js, new. The lock screen draws over everything:
                      "Welcome" with a name, a PIN and the PIN again; "Hi, <name>" with a
                      PIN; a reset screen for a damaged file; and the minute counted down
                      with the button held. Opening it reloads the page, so every face
                      starts over signed in. The sidebar names who is signed in beside a
                      Lock button, and api.js raises the lock screen on any 401
    the prover        tests/e2e/ollama_prover.py holds no PIN, nor should it: when health
                      says auth is on, S9-2 to S9-5 are proved LOCKED -- each must be
                      refused 401 -- and E2E_SCENARIOS.md says so under S9

**FORGOT THE PIN:** delete `webapp\data\user.json` and reload; the glass asks for a new one.
Nothing else is lost -- the record lives elsewhere.

**STROKES, 15 -> 23 in `webapp/handlers` and 4 -> 6 in `webapp/server`.** `handlers/lock_test.go`
is new and hermetic -- temp files, a fake clock, a fake door. A PIN is 4 to 8 ASCII digits and
nothing else (Arabic-Indic digits are digits to unicode, not to a keypad); a name is cleaned and
bounded; setup is refused from another machine and a second time, writes nothing when refused,
and never writes the PIN; the right PIN opens, a wrong one is counted and says how many tries are
left, the cookie is HttpOnly and Strict, and the session is the operator's; five wrong close the
lock for sixty seconds on the fake clock, the right PIN waits too, and the count starts over
after; the lock state tells the page what to draw, a damaged file included; the operator names
any world through rpcCallAs, scopeProject, CallTool and a chat stream while a key's session stays
pinned in all four, and no session reaches the door at all; an unconfigured lock says so. In
`server`: the glass listens on 127.0.0.1, and main.go still closes the gate and names the lock --
because the gate stood unused for weeks for want of ONE line, with every stroke above it green.
The open-path stroke holds the three faces.

**PROVEN BY REVERSAL, sixteen undos on a scratch copy**, each turning its own stroke red and no
other: setup from any machine; setup taken twice; a damaged file read as nobody; no lockout; the
right PIN skipping the wait; the tries left not counted; any unicode digit taken as a PIN; the
cookie back to Lax; the operator's empty tenant pinned in rpcCallAs, in scopeProject, in CallTool
and in StreamChat; a key's session no longer pinned; the glass on every address; the three faces
left gated; and main.go starting the glass unlocked. Put back, green. `gofmt` and `go vet` clean;
`db` 13, `handlers` 23, `server` 6.

**AND ON THE REAL BINARY, IN SCRATCH** -- built from these sources and run on its own port, :8097,
with a data folder of its own. His glass and his door were neither stopped nor changed, and
every check below reads the glass's own store, not the door. It bound 127.0.0.1 alone. Fourteen of fourteen over HTTP: health says
the gate is on; a data face and a setting refuse a caller with no session, 401 "login required";
the lock greets its one user; a second setup is refused; a wrong PIN is counted; the right one
opens, with the cookie HttpOnly and SameSite=Strict; the session reads the store and `/api/me`
names the user with no tenant; Lock ends the session on the server; five wrong PINs close the
lock, the right one waits with Retry-After 60, and the lock state carries the wait. The user file
it opened was written by a separate Python script with `hashlib`'s own PBKDF2, so the file is the
standard algorithm and not only this code's reading of it. The browser pane drew the welcome, the
lock with the user's name, and the countdown -- and opened by the name `localhost`, as his
Dashboard tab is, it reached the glass: the first connection about 0.3 s later than 127.0.0.1
would, because the browser tries IPv6 first, and 2 ms a request after that. The prover's S9
against it: S9-2 to S9-5 refused, as the lock requires. Stopped by pid.

**THE NEW BUILD IS NOT IN PLACE.** It waits in the hand's scratch, sha256 62b7dbb1a5a7b9a0,
10,676,736 bytes. Placing it over `webapp/atlas-webapp.exe` and restarting the glass is his
allowance; until then the glass runs the build from before this piece, gate open.

**AND PLACED THE SAME DAY**, on his word: *"place it and restart the glass"*. With no sitting
open and no engine running, the glass -- pid 24656, the only process listening on :8091, and
running from `webapp/atlas-webapp.exe`, checked by pid and path before anything was stopped --
was stopped; the build was copied over that file and hashes as built, the sources it was built
from unchanged; and the glass was started on its own command line. It listens on
127.0.0.1:8091 alone (pid 7160). `/api/health` says the gate is on, `/api/lock` says nobody is
set up, the data faces refuse a caller with no session, and his Dashboard tab reloaded onto the
Welcome screen. The name and the PIN are his to type. The build it replaced is kept in the
hand's scratch; the door was not touched.

**Named, not fixed.**

    the door         atlas-mcp still answers any program on this computer with no key
                     (`--auth` off), so the lock is the glass's alone. The door's own gate
                     and its holds are the other half of THE REACH
    the key login    /api/login (N6) still stands beside the PIN and opens a session,
                     pinned to its tenant, for any key the door verifies -- and while the
                     door is open, anything that reaches it can mint the first key, because
                     key creation is open on an empty store
    the roles        "multi-roles, all working under a single user" is the next piece, not
                     this one
    S9-1 and S9-6    the prover declares VERSION = "0.1.3" while the glass reports 0.1.5,
                     so S9-1 fails on that alone, before this piece and after it; and S9-6
                     passes on any failure to connect, as it always has, so under the lock
                     it passes without proving a stream
    a short PIN      anyone who can READ user.json could guess the PIN offline. The file
                     sits on his own disk beside the glass's database, and a person holding
                     that disk holds the machine already

## [0.1.6] — 2026-09-18 (tag on 0c65afc)

### 0.1.6 — THE DOOR'S OWN QUARTER, and the number is his: "atlas needs its own number too"

**THE MARK IS CUT**, 2026-09-18 17:52, annotated `THE DOORS OWN QUARTER`, on `0c65afc` --
from the panel's own Cut button, which filled `v0.1.6` itself from what VERSION declares at
that commit. The heading read `Unreleased` until the mark existed; the words under it are
unchanged.

**CORRECTED 2026-09-22: the mark holds 0.1.6 in its root `VERSION` only.** At `0c65afc` the other
nine pins say 0.1.5, so every binary built from this mark answers 0.1.5, `cargo test` is red there
(the spine's version strokes), and release.yml refuses it at its pin step -- no release was built
from it. The mark stays, because a mark is never moved. See [Unreleased], "Every version stamp in
step".

**WHAT THE NUMBER HOLDS** -- twenty entries below this one and above `## [0.1.5]`, eight
saves since `v0.1.5` was cut on 3dacdbc:

    the gate           RULE 6 stops being a convention: a writing call that is not the
                       operator's parks at the door and waits for his hand
    the loop           a node may be retried and a verdict may not
    the marks          a mark stands on the main line and leaves only after its history
                       has; the Send button and the version-tag flow both ask the door
                       BEFORE they offer; a mark that never left can be taken back; and
                       the card names which GitHub each world sends to
    the traces         the trace ledger, and a quarter of work making it cheap -- the
                       Dashboard's own reads are not traces, one event stream per tab,
                       the record read once per change, the health check resting while
                       nobody looks, and the store no longer rewritten whole to log a call
    the honesty        the glass stops waiting forever on a door that is gone, the door
                       stops waiting on a browser that left and sees an engine that died,
                       and the boot is no longer counted as a run

**PROVED ON HIS GROUND, 2026-09-18**: every Go package green in `line` (12 packages, the
tools battery included) and `webapp` (db, handlers, server), `go vet` clean, both binaries
built and running on this source.

### A mark can be taken back, and the panel says which GitHub it sends to

On his word, 2026-09-18: *"implement any missing features for github repo management that
arent on our version control panel yet."* **RESTART REQUIRED** -- `line/internal/tools/`
and `webapp/static/js/flows.js` moved, and the webapp EMBEDS its JavaScript, so neither
half is live until both binaries are rebuilt and restarted. No sitting was open.

**WHAT WAS ACTUALLY MISSING**, measured by reading the door's whole git surface against
the page: `git_tag` listed marks, cut them and sent them, and could not REMOVE one -- and
`git_remote` has answered where a world sends since 2026-09-10 while nothing in the glass
has ever called it. Everything else the panel needs was already there.

**THE DAY EARNED THE FIRST ONE TWICE.** `v0.1.12` was cut at a terminal on `7e64f20`, a
save that declares 0.1.11 and carried none of the work the mark names; the glass's own Cut
would have refused it on two separate guards. Taking it back then had to happen OUTSIDE the
glass, at the same terminal that made the mess -- `git tag -d`, no door, no record.

**`git_tag remove`, AND WHAT IT WILL NOT DO IS THE POINT.** A mark GitHub already has is
never withdrawn from here: somebody may have fetched it, and a name that vanishes leaves
them holding a version this ground no longer knows -- the same harm as moving one, which
`tagCut` has always refused in the same breath. Three guards, in this order:

    the name       mechanics only (`badMarkName`): a leading dash git would read as a
                   flag, the characters no ref may carry. NOT the version law -- four
                   of the six marks this estate has actually had to remove were named
                   `0.1.4`, `0.1.5`, with no `v`, and a remove that demanded a lawful
                   name would refuse exactly the marks that most need removing
    it exists      an uncut mark is named, not shrugged at
    GitHub         asked, and a question that CANNOT be asked fails shut: with the wall
                   down "I could not ask" is not "it is not there", because the wall may
                   have been open when the mark was cut and sent

**THE BUTTON ASKS THE DOOR FIRST, for this verb too**, which is the doctrine the Send
button earned the day before. The list now carries `removable` and, when it is false,
`why_not_remove` in the door's own words -- one judgement (`removeRefusal`) called by the
act before it deletes and by the list for every mark, so the greyed button and the refusal
behind it cannot drift apart. `remoteTags` came out of `tagList` for the same reason: the
list asks GitHub once for every mark, the act asks about one, and both must count an
unreachable remote as NOT KNOWN rather than as an empty set.

**ON THE GLASS.** A `Remove` button per mark, armed twice like Send -- it is the only
button on the page that destroys anything -- and greyed with the door's sentence where the
door would refuse. The row already says *on GitHub*, which IS the reason, so nothing is
repeated underneath it.

**AND THE ROW THAT NAMES THE REPOSITORY.** Every line on that card spoke of "GitHub"
without ever saying WHICH: this estate carries more than one world and they point
different places. `Sends to` now reads it from `git_remote` -- and prints the HOST the
door lifted out plus the path after it, never the raw URL, because a remote URL can carry
a token in its userinfo (RULE 7). A door that does not answer leaves the row saying so and
the rest of the card standing.

**PROVED.** Every Go package in `line` and `webapp` green, `go vet` clean. Four new
strokes: the sent mark refused and the local one taken back with the commit it stood on
named and that commit still readable afterwards; an unlawfully named mark removed anyway;
the wall-shut refusal; and the list's answer byte-identical to the act's. **By reversal,
four ways** -- the wall-shut guard, the GitHub-has guard, the name law and the list's
question each redden exactly the strokes that own them and nothing else. The stroke that
used `delete` as its example of a verb this door does NOT have now uses `move`: a verb list
that grows takes its own stroke with it.

**WHAT WAS DELIBERATELY NOT BUILT.** Pull requests, issues, releases and CI runs are
GitHub-side objects and need someone else's server -- RULE 4, and this file's own header
has said since it was written that they are "a separate wall to open deliberately, not a
dependency to acquire by accident". Pointing a world at a different remote is still not a
button (`git_remote` stays read-only): that is a decision, not a verb.

### The version-tag flow asks the door before it asks the hand

On his word, 2026-09-17: *"fix the version-tag flow's send node too"* -- and, when the hand
went looking in the wrong place, *"we are working on 127.0.0.1:8091/workflows ... there is
all the version control, EVERYTHING you need is on the webapp already."* He was right: the
builder came back on 2026-09-11 and the flow was edited THROUGH IT, not by hand on the
spec file. **No binary moved; no restart is owed.** No sitting was open.

**WHAT IT WAS.** `version-tag` v1 cut a mark, proved the cut, then paused at `send_gate`
with a static warning about irreversibility and carried on to `send`. Nothing in it asked
whether the mark COULD be sent. With the door's new refusal in place -- a mark leaves only
after origin's main line carries its commit -- the flow would have failed at its LAST node,
on `sent`, after the hand had already answered the only question it was asked.

**v2, EIGHT NODES.** One new step, `check` (kind `run`), between `proof` and `send_gate`:
it asks `git_tag list` again now the mark exists and reports THAT MARK ONLY -- whether the
door says it is sendable and, if not, the door's `why_not` word for word. The wiring moved
with it: `proof -pass-> check -always-> send_gate`, so everything after the cut still hangs
off the cut's own proof.

    read       run   the marks, and what the ground declares
    judge      gate  stop if the number is wrong
    cut        run   cut {{mark}} with {{what}}
    proof      eval  contains `Cut {{mark}} at`, scored on the EVIDENCE
    check      run   NEW -- would this door send {{mark}}, and if not, why not
    send_gate  gate  now renders {{out_check}}: the door's answer is IN the question
    send       run   one mark, by name
    sent       eval  contains `Sent {{mark}} to origin.`

**THE GATE CARRIES THE DOOR'S OWN SENTENCE.** A gate title is rendered against the run's
vars, so `{{out_check}}` puts what the door actually said into the pause the operator walks
back to -- hours later, on the waterfall, whole. The title also names the remedy in his own
terms: send the main line first from Version control, then fire this flow again. The
`send` step's objective says the same thing from the other side: report the refusal word
for word, and never push anything else -- the main line is sent by a hand, never by a flow.

**PROVED ON THE SAVED SPEC, against a stub engine in a temp home** -- no network, no
engine, no sitting, and the estate's record untouched. Three ways: with the door answering
`"sendable": true`, the send gate pauses carrying that; with the door refusing, the gate
carries the refusal's own words; and a refused CUT still ends the run FAIL with `check`,
`send_gate` and `send` never firing, which is the older rule this rewiring had to keep.

**AND BY REVERSAL.** The same harness against the folded `version-tag.v1.json` reddens both
door-answer strokes -- the old gate carried no answer because nothing had asked -- while
the refused-cut stroke stays green, because that half did not change.

**FOLDED, NOT REWRITTEN.** Saving through the page wrote v2 and kept v1 whole beside it.
The flow has still NEVER BEEN FIRED; firing it opens a sitting and spends the council, and
that is his.

### The Send button asks the door before it offers itself

On his word, 2026-09-17: *"make the send button ask the door first."* **RESTART REQUIRED,
BOTH: the door and the glass.** No sitting was open.

**WHAT IT WAS.** The Version marks panel offered Send on every mark GitHub lacked, and
the door's answer arrived after the click -- in fact after two, because sending arms
first. On a mark that is the worst moment to learn it: the answer is never "try again"
but "that mark may not go at all", and the operator has already decided by then. The
guards placed the same morning made this certain rather than theoretical: three of the
four shapes a mark can have are refused, and the panel advertised all of them.

**THE DOOR ANSWERS THE QUESTION BEFORE IT IS ASKED** (`internal/tools/gitctl.go`).
`sendRefusal` is now the one place that judges whether a mark may leave -- the name, the
mark's existence, the commit, and whether origin's main line carries it. `tagSend` calls
it before the push; `tagList` calls it for every mark GitHub does not already have, and
the list carries the verdict per mark: `sendable`, and `why_not` in the door's own
sentence when it is false. The wall is answered per mark too, in the wall's own words, so
a shut wall offers nothing and says why. A mark GitHub already has is not asked about --
there is no button on it.

**ONE JUDGEMENT, ONE WORDING.** Two judgements would drift, and then a button would say
one thing while the act behind it said another. The stroke holds them together: it reads
the list's `why_not` for a mark and the refusal the send itself returns, and fails unless
they are the same string.

**THE GLASS RENDERS IT AND DOES NOT SECOND-GUESS IT** (`webapp/static/js/flows.js`). A
mark the door will not send is a dead button carrying the refusal, and the reason is
printed under the table as well -- a refusal that lives only in a tooltip is a refusal
nobody reads. A `sendable:false` with no reason keeps the button live and lets the door
speak on the click, because a glass that invents its own refusal is the fault this piece
closes, from the other side.

**STROKES 61 -> 62** in `internal/tools`. `TestTheListSaysWhetherAMarkCouldBeSentAndWhyNot`
runs against a bare repository on this disk -- no stroke reaches a network -- and proves
all five answers: a mark origin's main line carries is offered; one ahead of its line is
not, and carries the reason; the list and the send refuse in the same words; sending the
line makes the same mark offerable; a name that is not `vMAJOR.MINOR.PATCH` is refused
before any of that; and a shut wall offers nothing and names the dial.

**PROVEN BY REVERSAL** on scratch copies, six ways: the piece as written is green; the
door as it stood this morning with the new stroke kept reddens it; the list not asking at
all, not judging the wall, and not judging the name each redden it alone; and making the
send speak its own words instead of the shared one reddens the new stroke AND the two
older send strokes -- which is the drift the design prevents, shown as it would look.
`gofmt` and `go vet` clean, every package in `line/` green, the door's battery 125/125,
the surface at 81 tools.

**AND FIRED, NOT ONLY READ.** A scratch door and a scratch glass, on their own ports,
against a scratch world carrying four marks -- one origin's main line carries, one ahead
of the line, one named `0.1.7`, and one on a side history. The panel greyed three with
their reasons under the table and left one live; clicking a greyed one did nothing and
armed nothing; the live one armed on the first click and sent on the second. The scratch
remote ended up holding `v0.1.5` and NOT the side history's commit, and that row then read
*on GitHub* with no button at all. His own door and glass were not touched: they kept
running throughout, and the scratch pair was stopped by pid.

### A mark stands on the main line, and leaves only after its history has

On his word, 2026-09-17: *"fix the tags ... make sure the version tags are being used
properly."* **RESTART REQUIRED -- and done: the door was rebuilt and restarted the same
hour (pid 14512).**

**WHAT IT WAS FOR.** Six marks on the core's ground -- `v0.1.0`, `v0.1.1`, `v0.1.3`,
`v0.1.4` and the lightweight `0.1.4` and `0.1.5` -- pointed into `pre-strip-master`, the
history kept from before `worlds/` was stripped, which still carries 383 paths under it,
268 of them under a `vault/`. Four of the six had names `git_tag` calls lawful, and the
Version marks panel offers a Send button for every mark GitHub lacks. Two clicks would
have published that history. `tagSend` checked the name, the wall and that the mark
existed; nothing asked WHERE its commit stood. The marks themselves were removed on the
core's side (see its CHANGELOG); this is the door learning to refuse the shape.

**THE CHANGE, both halves in `internal/tools/gitctl.go`.**

    cut    a mark is cut only on a commit the main line carries -- `main`, or
           `master` where a world calls it that. A world with neither is told so.
           The check sits after "a mark is never moved" and before the unsaved-work
           one, so the older refusals keep their order and their words
    send    a mark is sent only when origin's main line already carries its commit,
           as this machine last saw it. Then the push publishes the mark and nothing
           else. A mark ahead of its line is told to send the line first
    carries `merge-base --is-ancestor`, and any doubt -- a ref that is not there, a
           commit git cannot read -- is NO. A guard that fails open is not a guard

**STROKES, 59 -> 61 in `internal/tools`.** `TestAMarkIsCutOnlyOnTheMainLine`: a cut on a
side line is refused and creates nothing, and the main line's own commit still takes the
mark. `TestAMarkLeavesOnlyAfterItsHistoryHas` runs against a bare repository on this disk
-- no stroke reaches a network -- and proves the whole shape: a mark whose commit origin
carries is sent and the remote has it; a newer mark whose line has not been sent is
refused and NEITHER the mark NOR its commit reaches the remote; sending the line lets the
same mark through; and a mark made by hand around this door, on a line that was never
sent, is refused with its history left at home. That last case is the six, in miniature.

**PROVEN BY REVERSAL** on scratch copies: the door as it was (with the new strokes) turns
both red -- which is the measurement of what the old door would have done; each guard
undone alone turns only its own stroke red; judging the LOCAL main line instead of
origin's turns the send stroke red; and `carries` failing open turns both red. `gofmt`
and `go vet` clean, every package in `line/` green, and the door's own battery 125/125
with the surface at 81 tools.

**AND ON HIS DOOR.** Stopped by pid, the build placed, started on RUNBOOK's own line: it
lists the three marks the core carries, all three already on GitHub, and a cut asked at
`baa4f32` -- a commit `main` does not carry -- comes back *"Refused: baa4f32 is not on the
main line (main)."* Nothing was created, and nothing was sent.

**LEFT AS IT IS.** The Version marks panel still offers Send on any mark GitHub lacks and
learns the refusal after the click; nothing off the main line exists to press it on today.
*(CLOSED THE SAME DAY on his word -- the entry above.)*
The `version-tag` flow's `send` node now needs the main line sent first -- its own gate
does not say so yet, and it has never been fired.

### The record caught up with the week

On his word, 2026-09-17: *"tick the finished work, update the records to reflect the current
system."* The core's half, and the list of what is still owed, are in the core's CHANGELOG
("The record caught up with the week") and TASKS. Here, docs only; no code moved.

    README.md      "Version: 0.1.0+f1" -> 0.1.5, which VERSION says and the binaries
                   answer; the door's "78 tools" -> 81, counted off /tools; and, under
                   the sentence that says every tool call is a record, one paragraph
                   naming the exception D1 made: the Dashboard's own refresh reads are
                   answered and not kept. The sentence itself stands
    THE_ROAD.md    the 2026-09-08 snapshot kept as written, and a dated section under
                   it: the door's 81 tools, 31 with no page in the glass; the versions
                   since the cutover; what is unreleased; and that SEAT_LOG.md and
                   STATE_OF_BUILD.md, the road's own witnesses, were last written
                   2026-09-09

**Measured for it.** Of the door's 81 tools, 31 are named nowhere in `webapp/` outside its
tests -- a word match over the pages and the Go, the same one the 2026-09-11 count used. The
one mark this repository holds is `v0.1.5`, on 3dacdbc.

### The Dashboard's own reads are no longer traces

On his ruling, 2026-09-16: *"D1 b D2 30 minutes  D3 no"* -- the three decisions the
optimization pass left for him. This is D1, as he chose it: the Dashboard's own reads are
answered and not kept. D2, an idle engine closing itself after thirty minutes, is the next
piece, on his word. D3, a slower refresh while no engine is open, is ruled out.

**WHAT WAS THERE.** `Home.read` asks the door for `muster`, `rack_list` and `proofs` when the
Dashboard opens, when it is shown again and every fifteen seconds it is on screen.
`CallTool` kept every call as a trace carrying the tool's whole output and sent it to every
open tab, and every tab toasted "New trace recorded". One refresh was three traces and about
22 KB, so a Dashboard on screen for an eight-hour day wrote about 5,760 traces and 42 MB, and
the Traces page, which lists the newest 100, showed only its last eight minutes. Those three
tools were 99% of the 29,312 traces the ledger held when it was folded on 2026-09-14.

**THE CHANGE.**

    Home.read        asks its three as background reads: App.tool(n, a, true)
    App.tool         passes `background` on; nothing else sets it
    API.callTool     sends "background": true only when asked, so every other call's body
                     is exactly what it was
    CallTool         a call marked background, for a tool in backgroundReads, is answered
                     -- the door is still asked, and the reply still carries its output,
                     hash and duration -- and not kept: no trace, no broadcast, no toast,
                     and no trace_id in the reply, because there is none to name
    backgroundReads  muster, rack_list, proofs, each of which the door declares Writes:
                     false. Anything else marked background is kept as before, so no call
                     that writes can leave the record by asking to

Kept as before: the same three asked by anything else -- Records reading `proofs`, the end of
a turn repainting it, a call from the Tools page -- the sidebar's badges, every write and
every turn. A background read that fails is shown on the page as before, and is not kept
either. The 310 MB already in the ledger stays where it is.

**STROKES, 10 -> 15 in `webapp/handlers`**, in one new file, `handlers/trace_test.go`, against
a loopback door that answers every call: the three asked in the background leave no trace and
name none; no open tab hears of them; the door is still asked, and its words come back with
their hash and duration; the same three asked plainly, or with `background` false, are kept
and announced; and `git_commit`, `env_close`, `records` and `seats` marked background are kept
and announced. On the glass as it was, the first two are red and the other three green. The
module green on the mirror: `db` 13, `handlers` 15, `server` 4.

**AND THE PAGE, OUTSIDE A BROWSER.** A harness in scratch, not shipped, loads the real api.js,
council.js, home.js and app.js with the browser stubbed and reads every body the page posts.
Opening the Dashboard, and one refresh, each send the three and only the three, all marked
background; the page still paints from their answers, and a refused one still reaches the
brief; Records, the end of a turn, a bare App.tool, the sidebar's badges and a bare
API.callTool send no `background` key at all. Before: 9 of 11. After: 11 of 11.

**PROVEN BY REVERSAL, on a scratch copy.** Six undos of the glass -- the glass as it was, any
tool marked background left out, the list naming a call that writes, a background read still
announced, one still kept, a reply naming a trace it did not keep -- and seven of the page --
the scripts as they were, Home.read without the flag, App.tool dropping it, API.callTool
dropping it, API.callTool marking every call, API.callTool always sending the key, App.tool
marking every call -- each turned its own strokes red and no other. One stroke has no undo of
its own: that a background read is still answered, which nothing in this change would break.
`gofmt` and `go vet` clean, `node --check` on the three scripts.

**AND ON REAL BINARIES, IN A REAL PAGE.** His glass and his door were both down when this was
proved, and neither was started. The binary from before this piece and the new build each ran
from scratch, on a port and with a store of their own, and the browser pane opened each one's
Dashboard. With the door down every call came back "mcp error:"; what is kept does not depend
on the answer.

    opening the Dashboard   before: 6 kept -- seats, records, muster 2, rack_list, proofs
                            after:  3 kept -- seats, records, muster (the sidebar's)
    one refresh             before: 3 more, and the toast
    two refreshes           after:  none, no toast, and the brief still showed the rack refused
    a plain muster call     after:  kept, and the toast

The new build served api.js, app.js and home.js byte-identical to disk. Both were stopped by
pid.

**THE NEW BUILD IS NOT IN PLACE.** Copying it over `webapp/atlas-webapp.exe` was refused by the
session's permission check, and nothing was forced: the copy never ran, so the binary there is
still the build from before this piece, and a glass started from it asks as it always did. The
new build waits in the hand's scratch.

**AND ON 2026-09-17 IT WAS PLACED.** On his word, *"finish d1, then proceed to d2"*. Tried on
that order, the copy was refused by the permission check a second time; it ran after he answered
*"Yes, copy it now"* to a question naming that one copy, with no atlas-webapp running. It is the
build this entry describes -- the forty files under `webapp/` were compared with the tree it was
built from, unchanged -- and `webapp/atlas-webapp.exe` now hashes as that build does: sha256
e01791bd86729738, 10,613,248 bytes. Its file time still reads the build's, 2026-09-16 15:19,
because the copy keeps it. A copy of the binary it replaced, the piece-6 build of 2026-09-15,
is in the hand's scratch. The glass was down and was not started: its next start serves the new
build, and a Dashboard tab left open from before keeps the old scripts, which mark nothing as
background, until it is reloaded.

**AND ON HIS GLASS, THE SAME DAY.** On his word, *"start the door and glass, test it live"*,
both were started on RUNBOOK's own lines, and the glass served api.js, app.js and home.js
byte-identical to disk. The browser pane opened the Dashboard with the page's `fetch`, `toast`
and `App.onEvent` wrapped to log every call, and its Boot was pressed once, for D2's own test.

    opening the Dashboard   5 kept -- seats, records and muster, the sidebar's, and git
                            twice. The same page kept 8 on 2026-09-15
    the Boot                env_open kept, as a call that writes is
    one plain muster call   kept; its trace_added reached the page and the toast rose,
                            which is the proof that the wrapped functions see what they
                            are there to see
    the refreshes           408 background reads from 08:19:58 to 08:53:43, three every
                            15 seconds, 366 of them while the engine stood open, and not
                            one kept: the ledger went from 30,561 to 30,563 in that time,
                            and the two are the Boot and the plain call. In the 324 read
                            after the plain call was checked, no trace_added reached the
                            page and no toast rose
    afterwards              the wrappers put back and the tab closed; the glass and the
                            door left running

The first toast check was blind and is not counted. It watched for new elements, and `toast`
rewrites the one `#toast` element in place, so it could never have seen one. The wrapped
`toast` replaced it before any of the numbers above were taken.

**Named, not fixed.** The opening sentence of `README.md` still says ATLAS treats every tool
call as a record. It speaks for the whole harness rather than the glass, and was already wider
than the glass's ledger, which has only ever held the calls that passed through the glass.

### The glass stops waiting forever on the door, and its store's saves stop racing

On his order, 2026-09-15: *"continue to piece 6"* -- piece 6 of the optimization pass's plan,
"the glass stops waiting forever on the door, and its key store's save race is closed".

**WHAT WAS THERE.**

    no timeout       every call the glass made to the door -- /run/state, /tools, every
                     tool call, every proxied flow, prompt and chat call, and every
                     stream -- went through http.DefaultClient or http.Get, and neither
                     ever times out. A door that stopped answering held each caller, a
                     goroutine and a connection, for as long as it stayed silent.
                     ESTATE LAW 7: bounded everything, timeouts on calls
    a racing Save    Save took the store's read lock only to copy four fields and
                     marshalled after letting go -- but a map and three slices copy as
                     references to the same data. A SetKey landing mid-marshal ended in
                     Go's fatal "concurrent map iteration and map write", which takes
                     the whole glass down, or in a panic inside the JSON encoder; an
                     UpsertAgent was a torn read. And every Save wrote the one
                     store.json.tmp: sixteen Saves at once failed 265 to 277 times in
                     320, on that file's open or its rename, and every caller of Save
                     discards its error

**THE CHANGE.**

    pollWait 15 s    /run/state, which every page polls, and /tools. Both answer in
                     milliseconds
    callWait 30 min  every tool call, proxied call and stream. The longest fixed bound
                     the door sets on any tool of its own is rack_pull's download,
                     thirty minutes, so nothing the door promises to finish is cut
                     short. A flow's budget_s and a prompt_eval's cases are bounded by
                     what the operator wrote, not by the door: past thirty minutes the
                     glass stops waiting, and the door still finishes them and writes
                     their runs
    streams          also end when the browser that opened them leaves, even while the
                     door has yet to send a byte. Tool calls do not: their answer is
                     still written to the trace ledger and announced to every open
                     window after the one that asked has gone
    Save             holds a new saveMu from its snapshot to its rename, and marshals
                     under the read lock, so the Save that renames last wrote the
                     newest state

**STROKES, 11 -> 13 in `webapp/db` and 5 -> 10 in `webapp/handlers`**, the second five in one
new file, `handlers/door_test.go`. A door that never answers -- a real loopback server that
accepts every request and says nothing -- holds each call to its bound: /run/state and
/tools come back 502 saying the deadline passed; a tool call says so too, and its trace
records it; a proxied call; the council and chat streams end when their browser leaves
before the door has spoken; a WebSocket turn tells its socket. For the store: eight writers
set 800 keys, each write saving, and all 800 are on disk; sixteen Saves at once, twenty times
each, all land.

**PROVEN BY REVERSAL, on a scratch copy.** The whole piece undone: every one of those strokes
red, the store's taking the test process down with the encoder's panic. The old Save alone:
crashed 3 runs in 3, twice with the encoder's index out of range and once with the fatal
concurrent map write. saveMu removed and the marshal kept under the lock: the simultaneous
Saves red 3 in 3. The marshal moved back after the lock with saveMu kept: crashed 3 in 3.
Each of the seven door-call undos -- /run/state without its bound, /tools back on http.Get,
a tool call, a proxied call, the council stream and the chat stream untied from their
browser, the WebSocket turn -- turned its own stroke red and no other. The race detector
could not be used, because this machine has no C compiler for cgo, so the store's strokes
prove the race by its crashes and its failed Saves instead. `gofmt` and `go vet` clean over
the module, and the module green on the mirror.

**AND ON THE GLASS ITSELF.** The door was not touched. The glass -- pid 18408, the only
atlas-webapp running, checked by pid and path before anything was stopped -- was stopped,
the new binary put in place, and started on the same command line it had. /api/health
answers 200 with 29,415 traces; /run/state and all 81 tools reach the door through it; the
store opened with its settings.

**NOT FIRED LIVE.** No silent door was stood up against the live glass; the strokes' silent
door is a real server of its own.

**Named, not fixed.** `GetAgent` hands back a pointer into the agents slice, so a reader
holding it races an `UpsertAgent` the same way the Save did. And `Run.check` reads a 502 from
/council/state as "no engine open", not as a door that did not answer -- true of a door that
is down before this change, and now of one that is silent for fifteen seconds. The plan's
other pieces are his to call; none was started.

### The door stops waiting on a browser that left, and sees an engine that died

On his order, 2026-09-15: *"keep going"* -- piece 5 of the optimization pass's plan, "the
door drains a turn's events when a browser disconnects, and notices an engine that
crashed".

**TWO FAULTS IN THE DOOR, BOTH SILENT.**

    a closed tab     /run/stream hands each engine event to the socket through a queue
    stopped the      of 256. When the browser left, the handler returned and nothing
    turn             read the queue again -- and the sink was a bare send, so the first
                     event past the room left stopped the pump, the pump stopped
                     reading the engine, and the turn never ended. Run holds the
                     world's run lock for the whole turn, so every later turn on that
                     world waited on it forever. The glass's proxy closes its door
                     connection when a browser leaves, so closing a tab mid-turn, or
                     restarting the glass, was enough. The comment on that case read
                     "The turn keeps running inside the engine and its transcript
                     still lands"
    a dead engine    Alive read `exec.Cmd.ProcessState == nil`, and exec.Cmd fills that
    read as alive    only in Wait -- or Run, which calls Wait -- and the door calls
                     neither. It was nil for the life of every process: a crashed
                     engine stayed in the registry, /run/state said open, and the glass
                     offered a send box onto a process that was gone. That is the fault
                     I5 was written to fix, and the check it was fixed with could never
                     see it

**THE CHANGE.**

    the sink         sends; once the request is gone, an event that would wait is
                     dropped, because there is no one left to show it to. While the
                     browser is there it still waits, so a slow reader still slows the
                     turn, as before
    one waiter       Open starts one goroutine that waits on the process and closes
                     `exited` when it ends, asked to or not. Alive reads it. shutdown
                     reads it where it used to start a Wait of its own, and still kills
                     after the grace. os.Process.Wait, not exec.Cmd.Wait: that one
                     closes the stdout pipe the moment the process ends, while pump may
                     still be reading a dying engine's last lines

Both comments that said otherwise are corrected in place, dated, their sentences kept.

**MEASURED ON MIRRORS OF THE DOOR BEFORE AND AFTER, with a real child process.** The strokes
spawn the test binary itself as the engine, speaking serve.py's wire: `opened`, a delivery
per objective, `closed` on close, one objective that ends the process mid-turn with no
terminal event, and one that speaks as many tokens as it is told. It proves the door's
side of the wire and nothing about the engine.

    the crash    before: ten seconds after the process ended, Alive still said standing
                 after:  gone, and not handed back by the registry, inside 0.14 s
    the tab      before: a client that read the first line and hung up left a
                 50,000-token turn unfinished at 20 s; the stroke took 80 s, the rest of
                 it a close sitting out the 60 s grace
                 after:  the turn ended and was counted inside 0.37 s
    a close      0.01 s before and after, its toll reported paid

**STROKES, 10 -> 12 in `internal/engine` and 12 -> 13 in `internal/httpserver`**, and one
stub-engine function in each that skips unless it is the spawned child. PROVEN BY
REVERSAL, three undos on a scratch copy -- Alive back on ProcessState, shutdown deaf to
the waiter, the sink back to a bare send -- each turns its own stroke red and no other.
One line has no stroke: os.Process.Wait over exec.Cmd.Wait. `gofmt` and `go vet` clean,
the module green on the mirror, and the door's own battery 125 of 125 on the new binary.

**AND ON THE DOOR ITSELF.** With no sitting open and no engine standing, the new binary was
put in place and started on its own command line. It serves 81 tools; `env_list` names
research and atlas, both closed; the proofs reply's sha256 is unchanged; the glass reaches
it.

**AND A PROCESS THAT WAS NOT THIS HAND'S WAS STOPPED WITH IT, which is the lesson.** The
restart found the running door by its NAME rather than by the pid it had just looked at,
and two processes answered: this door, and an `atlas-mcp.exe` running from
`Desktop\Archive\atlas\line` with `--atlas-bin` and no `--http`, started by something else
during this piece. Both were stopped. That process was outside the ground and not this
hand's to touch, and it was not restarted: whatever launched it owns it. A door is stopped
by the pid that was looked at, never by its name.

**NOT FIRED LIVE.** Both faults need an open engine, and opening one opens a sitting in his
record.

**Named, not fixed.** `/run/listen` has the same shape with a queue of 64; how many lines a
capture sends through it is voice.py's to say, and that file was not read in full.
`/chat/stream` has a different one: when the browser leaves, its handler returns while the
send's goroutine goes on writing tokens to the response. The plan's other pieces are his
to call; none was started.

### The door reads the record once per change, and looks for old code only where the engine keeps it

On his order, 2026-09-15: *"continue the work"* -- piece 4 of the optimization pass's plan,
"the door caches the record files until they change, and limits its running-old-code
check to `manjuel/`".

**WHAT THE PASS MEASURED, on the ground, from a mirror.**

    proofs     15-16 ms and 7.6 MB of garbage a call, 11 ms and 6.0 MB of it parsing
               sessions.jsonl -- 846 KB, append-only, changed once a sitting. The
               Dashboard asks on every refresh
    the walk   every /run/state an open engine answered walked the whole ground past
               a skip list: 289 folders, 5,891 files, 27 ms and 3.5 MB. 30 of the 64
               .py it counted were atlas/'s tools and tests, which no engine imports,
               so an edit to one said "restart" over code the engine never ran

**THE PLAN WAS WRONG ABOUT ONE FOLDER, AND THE DISK SAID SO.** `manjuel/` is not all the
code a running engine holds. The law gate loads `law/law.py` into the process, that
loads the pen's `links.py`, and `links.py` imports `jesster` by name -- which Python
keeps until the process ends. A walk of `manjuel/` alone would miss an edit to the pen
while the engine runs the old copy of it, so the walk enters `law/` as well. `law.py`
and `links.py` themselves are loaded fresh each time the gate walks the chain, so an
edit to either raises the row without the engine being stale. The old walk counted them
too, and law.py's own header calls the pen never edited.

**THE CHANGE.**

    engine.CodeChanged   reads the top of the ground and enters manjuel/ and law/;
                         every other folder at the top is skipped whole, and
                         __pycache__ anywhere. Same question, same answer: the
                         newest .py the engine holds, and when
    tools.fromRecord     run_history.jsonl, parity_history.jsonl, sessions.jsonl and
                         SEAT_LOG.md are read once, and what each was made into is
                         kept -- the standups, the parity rows, the counts and the
                         twelve recent sittings, never the bytes -- until the file's
                         size or modification time moves. One stat a file a call.
                         last_run.json, a few hundred bytes, is read every call as
                         before
    three rules          the stat is taken before the read, so a write landing
                         between them costs one extra read and never an old answer;
                         a read that fails is not kept; nothing kept is changed
                         after, and every call builds its own reply around it

**LEFT AS IT WAS, AND WHY.** `SittingOpen` reads sessions.jsonl whole, and the plan named
it. It runs on `env_open` and `env_list` alone -- nothing polls it -- and it is the guard
between a second engine and a forked ledger, so it still reads the disk every time.
The listing of `logs/` for standup reports, 1.3 ms, is still read every call: a folder's
time is a weaker key than a file's, and a report missing from the page would be the page
lying about the record.

**MEASURED ON THE GROUND, from mirrors of the door before and after, reading the same
files.**

    proofs      16.0 ms, 7,775 KB a call  ->  the first call 15.8 ms, then 1.6 ms, 352 KB
    the walk    27.4 ms, 3,552 KB a call  ->  0.40 ms, 29 KB
    the reply   14,828 bytes, sha256 05fd341190cdf6b1: before, after, first call and
                thirty-first, all byte-identical, and identical to the live door's
    the answer  manjuel/pipeline.py at 2026-09-14 12:41:20, before and after

**STROKES, 7 -> 10 in `internal/engine` and 55 -> 59 in `internal/tools`.** The walk
ignores atlas/'s tools, a skill's .py, agent_workspace, worlds and __pycache__ when each
is newest; it finds manjuel.py, the package, law.py and the pen's jesster.py when each
is; and an engine is stale when the pen changes after it started, and not when an atlas
tool does. For proofs: an unchanged record is read once over four calls and answered
byte for byte the same; a ledger that grew is read again with its time held still; a
file rewritten to the same size is read again when its time moves; a read that failed is
reported and not kept. PROVEN BY REVERSAL, twelve undos on a scratch copy -- the old skip
list, law/ dropped, the folder rule applied at every depth, __pycache__ entered, nothing
kept, the size dropped from the key, the time dropped, a failed read kept, and each of
the four call sites reading its file directly again -- each turns its own strokes red and
no other. Two lines have no stroke: the stat-before-read order, whose race a test cannot
time, and the copy that keeps twelve rows rather than the array behind them. `gofmt` and
`go vet` clean, the module green on the mirror before and after, and the door's own
battery 125 of 125 on the new binary.

**AND ON THE DOOR ITSELF.** With no sitting open and no engine standing, the door was
stopped, the new binary put in place, and started on its own command line, unchanged. It
serves 81 tools; `env_list` names research and atlas, both closed; the glass reaches it.
`proofs` over the door's own /rpc, twelve calls each: median 13.8 ms before, 3.1 ms after,
the reply's sha256 unchanged.

**NOT FIRED LIVE.** The walk runs only for an open engine, and opening one opens a sitting
in his record, so the new walk has run against the ground from a mirror and not inside the
door. The next Boot is the proof.

**WHAT IT SAVES:** per Dashboard tab on screen between sittings, 11 to 14 ms and 7.2 MB of
garbage off every refresh -- 60 to 80 CPU-seconds in the door for a day of refreshes -- and
with an engine standing, 27 ms and 3.4 MB off every /run/state. The plan's other pieces are
his to call; none was started.

### The health check rests while nobody is looking

On his order, 2026-09-15: *"keep going"* -- piece 3 of the optimization pass's plan,
"the health check pauses while a tab is hidden, the same rule the refresh already
follows".

**WHAT THE PASS FOUND.** `App.loadHealth` asked `/api/health` every five seconds, on
every page, for the life of the tab, and rescheduled itself whether or not anyone
could see the page -- 17,280 requests a day per open tab, to paint one dot and a
version in the sidebar. The Dashboard's own refresh already stops while a tab is
hidden and catches up the moment it is shown (`Home.watch`); this never learned
that rule. Measured the day before, in a hidden pane: 12 health checks in 69 seconds.

**THE CHANGE, in `loadHealth` alone.**

    hidden          asks nothing and leaves nothing scheduled
    shown again     checks at once -- that is the moment the dot is read -- and
                    resumes the five-second cycle. The listener is installed by
                    the first check, so the rule lives in one function
    one cycle       a check started by the return to the tab can overlap a timed
                    one still waiting on its answer; the timer is cleared again
                    before the next is set, so two cycles can never form
    hidden mid-check  the answer still paints; nothing is scheduled after it

**MEASURED.** A harness in scratch, not shipped, drives the real app.js with timers
captured and visibility switchable:

    before   5 of 9: a hidden tab asked four times across three due timers, kept
             a check scheduled, and nothing checked on the way back
    after    9 of 9: a visible tab checks and keeps checking; a hidden tab asks
             nothing and schedules nothing; coming back checks at once and
             resumes one cycle; racing checks leave one cycle; a tab hidden
             mid-check schedules nothing

PROVEN BY REVERSAL: each of four parts undone on a copy turns its own check red.
The reversal script itself was wrong once and is corrected: it expected removing the
first hidden-test to also leave a check scheduled, and the second test, at the
reschedule, rightly still prevents that. The harnesses for pieces 1 and 2 and the
watched-turn fix still hold 9 of 9, 8 of 8 and 16 of 16. `node --check` clean,
`go vet` clean, the webapp's tests green, built.

**AND ON THE GLASS ITSELF**, rebuilt and restarted. The browser pane's page reports
itself hidden, and it made no health check at all in 23 seconds after load. With the
page marked visible and the visibility event fired in the real browser, it checked
at 0 ms, then at 5.0 s and 5.0 s, and painted the dot green with version 0.1.5;
marked hidden again, the next due check asked nothing and left no timer.

**WHAT IT SAVES:** every health request from a tab nobody is looking at -- up to
17,280 a day per tab. The plan's other pieces are his to call; none was started.

### A tab holds one event stream, not two

On his order, 2026-09-15: *"keep going"* -- piece 2 of the optimization pass's plan,
"one event stream per tab, shared by both listeners".

**WHAT THE PASS FOUND.** Two readers follow the broadcast bus: `App.onEvent`, for the
toasts, and `Run.mirror`, which shows a turn another window started. Each called
`API.sse`, and `API.sse` opened a new `EventSource` every time it was called -- so
every tab held two `/api/events` streams, and the webapp sent that tab every
broadcast twice: each trace with its tool's whole output, each token of each turn.
Measured on the glass before this: one Dashboard tab, two `GET /api/events`.

**AND ITS RECONNECT MULTIPLIED.** An error fires on every failed retry, and each one
queued its own reconnect three seconds out. In the harness, two streams erroring
twice opened six before they settled.

**THE CHANGE, in `API.sse` alone.** The first caller opens the stream; every caller
is a listener on it, and each message is parsed once and handed to each. A listener
that throws is caught on its own, so it cannot stop the next one -- the isolation
the two separate streams used to give. One reconnect is pending at a time, or none.
Neither caller changed.

**MEASURED.** A harness in scratch, not shipped, loads the real api.js and council.js
with the browser stubbed:

    before   4 of 8: a tab with both readers opened two streams; repeated errors
             opened six and left two open
    after    8 of 8: one stream; a broadcast reaches the toast reader and the
             mirror, which opens its watched turn from it; a listener that
             throws does not stop the next; a malformed line reaches no one;
             repeated errors reconnect once, leave one stream open, and the new
             stream still serves both readers

PROVEN BY REVERSAL: each of four parts undone on a copy turns its own checks red.
Piece 1's harness still holds 9 of 9 and the watched-turn harness 16 of 16.
`node --check` clean, `go vet` clean, the webapp's tests green, built.

**AND ON THE GLASS ITSELF**, rebuilt and restarted: a Dashboard tab opens one
`GET /api/events`, the stream stands open with both readers on it, and one refresh's
three trace broadcasts arrived over it with the toast raised.

**WHAT IT SAVES:** half of every broadcast byte the webapp sends each open tab --
every trace and every token, once instead of twice. The plan's other pieces are his
to call; none was started.

### The Dashboard reads the record once per refresh, not twice

On his order, 2026-09-15: *"continue"* -- piece 1 of the optimization pass's plan,
"one `proofs` fetch per refresh".

**WHAT THE PASS MEASURED.** One open Dashboard tab made exactly 960 tool calls an
hour, four per 15-second refresh, and two of the four were the same `proofs` call.
`Home.read()` asked for it in its own `Promise.all` and kept the answer, then called
`App.paintProof`, which asked again. The comment above that code read "proofs is
asked for ONCE and kept"; the code beneath it did not. Over 2026-09-12..14 the trace
ledger held 14,440 `proofs` calls at ~20 KB each -- half of every call the glass
made, and 290 MB of its 300. Opening the page read it three times: `render()`
painted the scores on a fetch of its own before `read()` fetched twice more.

**THE CHANGE, SMALL ON PURPOSE.**

    App.paintProof(box, only, read)   takes an answer the caller already holds --
                                      the tool's text, or the {err} Home.read's
                                      ask() hands back -- and fetches only when
                                      it is handed nothing
    Home.read                         hands its own read over, refusal included
    Home.render                       paints no scores of its own; read() does,
                                      and the box says "Reading the record..."
                                      until it lands
    unchanged                         Records and the end of a turn still ask for
                                      themselves: neither holds an answer, and a
                                      turn can change the record

The comment that said otherwise is corrected in place, dated, its sentence kept.

**MEASURED.** A harness in scratch, not shipped, loads the real council.js, home.js
and app.js with the browser stubbed and counts every tool call:

    before   opening the page reads proofs 3 times; one refresh is 4 calls with
             proofs twice; a refused proofs read is asked for twice
    after    9 of 9: opening reads proofs once; a refresh is 3 calls; a refusal is
             painted as refused without a second ask; Records and a turn's end
             still read for themselves

PROVEN BY REVERSAL: each of four parts undone on a copy turns its own checks red.
The harness for the watched-turn fix still holds 16 of 16. `node --check` on both
scripts, `go vet` clean, the webapp's tests green, built.

**AND ON THE GLASS ITSELF**, rebuilt and restarted: opening the Dashboard wrote 8
traces, `proofs` once, where it wrote 10 with `proofs` three times the day before;
one refresh, with the page's own timer paused so nothing else could land, wrote 3
traces and 22 KB -- `muster`, `rack_list`, `proofs` -- where the tick had written 4
and about 42 KB. The scores and the sitting card painted from that one read.

**PER OPEN DASHBOARD TAB, PER DAY:** 5,760 fewer `proofs` calls and about 116 MB less
trace ledger -- a quarter of the calls and nearly half the bytes. The plan's other
pieces are his to call; none was started.

### The glass stops rewriting its whole trace store to record one call

On his order, 2026-09-14: *"address the found issues"* -- the first piece of the
day's second diagnostics pass.

**WHAT THE PASS MEASURED.** `atlas-webapp` held 17.3 GB and had used ~7,300
CPU-seconds since a 13:21 start. Tool calls through it took 3-25 s (`muster` 2.9,
`git` 13.2, `env_list` 16.0, `records` 25.3), and the Dashboard's Close the
sitting did not reach the door for 9 min 20 s: sitting 224 stood open with its
engine idle, and the close went through only after a second tab was shut.

**ONE CAUSE.** `CallTool` keeps every call as a trace carrying the tool's full
output, and `db.AddTrace` ended in `Save()`, which marshalled every trace ever
recorded -- indented -- into `store.json` and renamed it into place. At the stop:
29,312 traces, 303 MB, rewritten on every call, and every `SetKey`, `AddEval`
and `AddMessage` rewrote them too. 99% of them are the panel's own polls:
`proofs` 49%, `muster` 25%, `rack_list` 25%.

**THE LEDGER** (`webapp/db/db.go`):

    traces.jsonl   one line appended per trace, and nothing rewritten to add one.
                   Every trace is kept
    the window     memory holds the newest 1,000 -- more than any reader asks for
                   (ListTraces 100, an agent's page 50, search 500)
    by id          a trace older than the window is read from the ledger, so an
                   eval that names an old trace still opens it
    the count      /api/health reads the ledger's tally; it copied the whole list
                   to measure it, on a route every open tab polls
    store.json     evals, agents, messages and keys -- no traces
    the fold       a store.json that carries traces moves them into the ledger at
                   start: written beside it, renamed in one step, and only then is
                   store.json rewritten without them. Run twice, it writes no
                   trace twice. A store whose traces cannot be moved refuses to
                   start, rather than lose them at the next save
    a torn line    half a line left by a killed process is not counted, and the
                   next trace starts on a line of its own

**MEASURED BEFORE AND AFTER, on synthetic stores in scratch** -- no record content
left the ground; each binary on its own port, the door unreachable, twenty traces
posted one after another:

    5,000 traces, 57 MB      392 ms a trace, 1.1 GB peak   ->  0.9 ms, 171 MB
    15,000 traces, 173 MB    1,176 ms, 2.6 GB              ->  0.8 ms, 480 MB (the fold)

**AND ON THE GLASS ITSELF.** Stopped, rebuilt in place, restarted on RUNBOOK's own
line; the door untouched. The fold moved 29,312 traces into a 301 MB ledger of
29,312 lines, `/api/health` reads 29,312, and the oldest trace (`muster`,
2026-09-12) and the newest (`proofs`) open by id with their hashes unchanged.
`store.json` is 2,821 bytes. The same six tool calls: 0.01-0.37 s. The process:
80 MB. A fresh Dashboard load: 14 calls, median 9 ms, no script error.

**STROKES, 6 -> 11 in `webapp/db`.** A trace is one line and store.json is not
touched (its mtime stamped first, so a same-bytes rewrite still shows); the window
is bounded while the count and an out-of-window lookup survive a reopen; a
pre-ledger store folds in order with nothing lost, and folds nothing the next
time; a fold that died between its rename and its rewrite writes no trace twice;
a torn last line is skipped and the trace after it kept.
`TestTheSaveIsAtomicAndLeavesNoScratch` drove its save through `AddTrace`, which no
longer saves; it is driven through `SetKey`, and now also refuses a store.json
that carries traces -- the guard unchanged. PROVEN BY REVERSAL, seven undos on a
scratch copy: the save put back into AddTrace, the traces put back into
store.json, the window unbounded, the ledger lookup removed, the fold's
duplicate check removed, the torn-line newline removed, the torn-line test
loosened -- each turns its own stroke red and no other. `gofmt` and `go vet`
clean; the webapp's tests green.

**NOT CHANGED, AND NAMED FOR HIM.** The ledger still grows with every call the
panel makes, and its own polls are nearly all of them; whether a panel's reads
should be traces at all is his call. And `Save()` still marshals the keys map
after releasing its lock, so a `SetKey` landing mid-marshal can crash the process
on a concurrent map write -- true before this change, and not touched by it.

### The boot no longer shows up as a run, and a watched turn's clock stops

On his order, 2026-09-14: *"keep going"* -- the third piece of the day's
diagnostics findings.

**WHAT THE DIAGNOSTICS PASS SAW.** After a Boot, the Dashboard's run card held a
turn nobody had typed, its rows drawn as `unnamed` and raw JSON, under a badge
reading "done" whose seconds never stopped climbing.

**THREE FAULTS, ONE PATH: `Run.mirror()`.**

    the tab watched       the mirror follows the broadcast bus so a turn
    its own boot          started in another window shows here, and ignores
                          its own echo -- but `running` covers only turns
                          begun through `Run.start`. `Home.bootStep` opens
                          /council/stream itself, so /warm came back as a
                          WATCHED turn, and /status's lines landed in it too,
                          because that turn never ended
    a watched turn        it ended on `delivery` or `refused` alone. A
    never ended           /command ends with `command`, and `aborted`,
                          `cancelled` and `unreachable` end turns as well;
                          none of them closed it, so the clock ran on
    the door's frames     `pipeSSE` re-broadcasts each `data:` line WITHOUT its
    became `unnamed`      frame name, so stream_open, stream_end and
                          stream_error arrived with no `event` and were kept as
                          `unnamed` rows the run card drew as raw JSON. The
                          runner keeps no row for any of them

**FIXED IN THE GLASS ALONE; the door and the engine are untouched.**

    Run.own          counts council streams a tab opens outside Run.start.
                     bootStep raises it for its stream and lowers it
                     ECHO_GRACE_MS (a second) after the stream ends, because
                     the bus copy of a line can trail the direct one
    Run.TERMINAL     serve.py's TERMINAL, spelled the same: delivery, refused,
                     aborted, cancelled, unreachable, command. A watched turn
                     ends on any of them. `error` is not one, and stays not one
    door frames      told apart by the one field each carries that the others
                     do not -- runstream.go gives stream_open `objective`,
                     stream_error `error`, stream_end `dropped` -- and kept as
                     no row. stream_error ends a watched turn as refused; a
                     turn waiting on its gate stays open, so the answer's
                     stream still lands in it
    App.runRow       `command` has a row of its own instead of raw JSON

council.js's header said PROTOCOL 1 carries seventeen events. It has carried
nineteen since `heard` and `command` joined on 2026-09-09; corrected in place,
and dated.

**PROVED OUTSIDE A BROWSER, BOTH WAYS.** This console has no JS suite, so the
proof is a node harness in scratch, not shipped: council.js and home.js loaded
with the browser stubbed, and the exact bus sequences replayed -- the booting
tab's own /warm and /status, a watched /warm, a delivery, a door refusal, a gate
and its answer, the three other terminals, and a non-terminal `error`.

    before the fix   5 of 16 held. The booting tab's run card got a "/warm"
                     turn holding `unnamed, command, unnamed, unnamed, command,
                     unnamed` -- the sighting, reproduced
    after            16 of 16
    reversal         each of four parts undone on a copy turns its own checks
                     red: the echo test in council.js, the count in bootStep,
                     the terminal list, and the frames

`node --check` on the three scripts, `go vet` clean, the webapp's tests green,
built. The glass was stopped, rebuilt in place and restarted on RUNBOOK's own
command line; it serves all three scripts byte-identical to disk, and a fresh
page load logs no script error.

**NOT FIRED LIVE.** No engine was booted to watch it end to end, because a boot
opens a sitting in his record. The next Boot is the proof -- and an open
dashboard tab needs a reload before it carries the new scripts.

**Named, not fixed:** the mirror ends a watched turn with `done`, and neither the
Dashboard's thread nor the Watchboard listens for it -- both close a bubble on
`end` -- so a watched turn's bubble stays marked live after it lands.

### The record caught up with its own tag

On his order, 2026-09-14: *"bring the record up to date first."*

**v0.1.5 WAS CUT ON 2026-09-12 AND THIS FILE NEVER SAID SO.** Everything from
0.1.3 on still sat under `[Unreleased]`, so the one heading a reader trusts to
mean "not shipped" covered two days of shipped work. `[0.1.5] — 2026-09-12
(tag on 3dacdbc)` heads it now. 0.1.3 and 0.1.4 were built and never tagged;
they ride inside that tag, and the heading says so.

**AND TWO BLOCKS I WROTE INTO THE RELEASE HAD LANDED AFTER IT.** The release.yml
fix and the version-tag flow were appended to the git_tag entry as it grew, but
`git log v0.1.5..HEAD` puts both commits after the tag. Each is its own entry
under `[Unreleased]` now, words unchanged, carrying the commit it landed in.
Moving them is not rewriting them: a release's section should describe what the
tag holds.

**THREE DOCUMENTS PINNED TO 0.1.5 WERE WRONG ABOUT 0.1.5 -- MEASURED AT THE
TAG, NOT AGAINST TODAY.** The tagged tree was exported with `git archive v0.1.5`,
and its own door was built and asked:

    the door's surface   DELIVERABLE, ACCEPTANCE and PIPELINES said 78; the
                         tag's door reports 79 -- git_tag landed inside the
                         release, after the count was taken
    line test fns        DELIVERABLE said 229; `go test -list` on the tag
                         lists 167, and every func in the test files, helpers
                         included, is 239 -- 229 matches neither
    webapp test fns      DELIVERABLE said 1 fn / 1 pkg; the tag has 15 / 3
    the glass's pass     DELIVERABLE said it had not happened; 572186d, the
                         commit that wrote that sentence, is the one that did it
    workflows            DELIVERABLE's inventory said ONE; its own CI section,
                         and the tag, hold two
    release.yml          PIPELINES' last line had it running all six pipelines
                         and building cross-platform; the section above that
                         line already says it does neither

The 125 prove strokes that ACCEPTANCE and PIPELINES state were built from the
same tag and held.

**THE ORIGINALS STAND.** Every correction is dated and sits beside the number it
corrects; no count was overwritten. These are snapshots of a release, and the
0.1.3 entry below already holds the ruling on that shape: renaming a snapshot
"would claim it was proven at 0.1.3, which it was not."

**NOT CARRIED FORWARD, deliberately.** The door carries 81 tools today, not 79:
`hold_list` and `hold_answer` came after the tag. That number belongs to the next
deliverable, written when the next number is cut -- not to a document about
0.1.5. ADR-006 keeps its 78 and its 71 core, because a decision record's
measurements are part of what was decided.

### The hold queue reaches the glass, and the whole loop was FIRED

**AND THE FIRST THING FIRING IT FOUND WAS THAT THE WARNING COULD NOT BE SEEN.**
`hold_list` refused any caller it could not verify as the service wire -- and
with `--auth` off the door verifies nobody, so the glass asked, was refused,
and rendered NOTHING. The one state where a reader most needs telling that
RULE 6 is not being enforced was the one state he could not see. The disarmed
answer is now given to ANYONE: nothing parks while holds are off, so there is
no queue to leak, and the warning is the entire value of the reply. Armed, the
queue itself stays his -- it names callers, projects and arguments.

**"Waiting for your hand"** sits ABOVE the repositories on Version control,
because a call parked on his decision outranks the state of a tree nobody is
asking about. Each row carries the tool, who asked, the project, and the
arguments whole, with Approve and Deny. Approving runs EXACTLY the parked call.

**THE WHOLE LOOP WAS THEN FIRED AGAINST A LIVE ARMED DOOR** -- a second door on
:8099 with `--auth` and a service wire, over a throwaway ground, because eight
green strokes do not prove the transport actually produces the Caller:

    no key                401 before the tool is ever reached. With the gate
                          armed an unkeyed stranger never meets the hold.
    a KEYED SEAT commits  HELD, named by its own key id, parked with an id --
                          and `git log` confirmed NOTHING WAS COMMITTED.
    the seat approves it  "answering a hold is the operator's act alone".
    the seat reads it     refused; the queue is not a stranger's to read.
    the service wire      sees the parked call with its arguments whole.
    it approves           the commit lands, with the message the SEAT wrote,
                          from the args it was parked with.
    state/holds.jsonl     `held` then `approved_ran`, both naming k-768ca119
                          rather than the operator who released it.

So the hold's real subject is a caller that HAS a key and is not the glass --
a seat given its own bearer. That is exactly the hand RULE 6 was written about.

The probe door and its ground were destroyed afterwards; the estate's own door
was never touched and still reports `auth=false`.

### RULE 6 stops being a convention and becomes a gate

The second of the three things worth taking from the sovereign-microkernel
read: its `high_impact` + `pending_approvals` queue. Everything else that
MCPServer offered, this door already had in a real form rather than a
simulated one.

**WHAT WAS ACTUALLY TRUE HERE BEFORE THIS.** RULE 6 says no agent commits,
pushes, lands, approves or authorises a spend -- and this door enforced that by
NOBODY HAVING WIRED AN AGENT TO THE WRITING VERBS. That is a convention. A seat
that reached `git_commit` would have committed.

**AND THE DOOR COULD NOT TELL WHO WAS ASKING.** `Registry.Call` took
`(tenants, name, args)` and no principal at all: the operator's glass and a
seat were the same request. The `actor` argument RBAC already reads comes
straight out of the caller's own args, so anyone may declare themselves anyone.

    Caller{Name, Service}   built by the DOOR from the transport. `Service` is
                            the only bit that decides anything: the request
                            carried the service wire, so it is his own glass.
    two call sites          httpserver (from the credential the gate already
                            builds) and protocol (stdio is never the glass, by
                            construction rather than by check).
    the hold                a Writes:true tool called by a non-service caller
                            PARKS. Nothing runs. The caller is told HELD, told
                            nothing was written, told which rule holds it, and
                            told not to work around it.
    hold_list/hold_answer   the operator's alone -- both refuse a non-service
                            caller. If an agent could answer its own hold the
                            queue would be a formality it walks through itself.
    the record              every park, approval and denial appends to the
                            ground's own `state/holds.jsonl`, written BEFORE
                            the answer: a hold he cannot find later is a call
                            that vanished.

**IT IS INERT WITHOUT `--auth`, AND IT SAYS SO IN THREE PLACES.** With the gate
off the door reads no credential, so `Service` is false for everyone; holding
on that would stop his own panel while stopping no agent that thought to send a
header. So holds arm off the same dial as the gate, `hold_list` answers
`"armed": false` with the reason and the fix, and the boot line now reads
`holds off -- any caller may write` or `holds ARMED`. A guard that quietly does
nothing is worse than no guard, because it is believed.

**A CORRECTION I MADE TO MY OWN COMMENT, BECAUSE THE MUTATION POINTED AT THE
WRONG LINE.** The header first claimed the load-bearing guard was `Call`
deleting the reserved key out of inbound args. Mutating that line away changed
NO verdict -- the decision never consults the args, and the door overwrites
them on the way out. The line that actually holds is `!caller.Service`: point
it at the args instead and a forged call writes straight through, which the
stroke reports as `WROTE smuggled`. The comment now says which is which, and
the strip is named as belt and braces rather than as the reason.

**AND THE STROKES FOUND A REAL BUG BEFORE THE FIELD DID.** A hold's id was
`hold_<millis>_<tool>`, so TWO PARKS IN THE SAME MILLISECOND COLLIDED: the
second silently replaced the first in the map, one call vanished, and approving
the surviving id would have run the other one's arguments. Two agent writes
inside one millisecond is a Tuesday. A sequence under the same lock fixes it.

Eight strokes, each both ways -- held and not held, read and write, glass and
seat, armed and disarmed, approved and denied, and a hold answered twice.
Nineteen packages green, gofmt clean, the door's own battery PROVEN.

**WHAT THIS CHANGES FOR HIM TODAY: NOTHING, DELIBERATELY.** The running door
was started without `--auth`, so holds are off and every path behaves exactly
as it did. Arming it is his call and it is a relaunch, not an edit -- and the
glass has no panel for the queue yet, so until that lands the answer path is
`hold_answer` through the door. Both are named here rather than discovered.


### A node may be retried, and a verdict may not

From the read of the three `sovereign-agent-harness-&-os-microkernel` folders
he built with Google, on his order to diff them against this estate and merge
what is worth merging.

**THE DIFF BETWEEN HIS THREE VERSIONS COLLAPSED FIRST.** `microkernel.py`,
`cli.py`, `tui.py` and `gateway.py` are BYTE-IDENTICAL in all three -- the
Python backend never changed across v1, v2 and v3. v2 added a prover, docs,
deploy and CI; v3 is v2 with two frontend files moved and six spent patch
scripts deleted. So there was one system to read, not three.

**AND MOST OF IT, THIS GROUND ALREADY HAD -- REALLY, WHERE THAT ONE PRETENDS.**
Measured, not asserted:

    its PPMIRouter        keyword scoring in an information-theory coat; the
                          slot regexes it declares are never called, and three
                          of its five intents are scored and never acted on.
                          Against manjuel's intent.py, which is stroked.
    its MCPServer         four tools over a two-entry DICT. `execute_bash`
                          RETURNS "[SIMULATED_EXECUTION_STDOUT] ... exit code
                          0" WITHOUT RUNNING ANYTHING, always success. Against
                          the door's real tools, ground.Barred and the jail.
    its AgentDaemon       never calls a model. The "neural loop" is a hardcoded
                          if/elif named `simulated_decision`.
    its telemetry         VRAM 4.8/24.0, 42.6 tok/s, "ledger_verified": true --
                          asserted, never measured. Its CLI prints a DIFFERENT
                          hardcoded set (VRAM=14.2GB, CPU=22%) two files away.
    its ledger            hashes prev|timestamp|payload and LEAVES event_type
                          OUT of the signature, so any row's type can be
                          rewritten and verify_ledger_integrity still says
                          VERIFIED. Against law/law.py, which signs the link.
    its checkpointing     the one thing I expected to be a genuine gap here --
                          and it is not. flow.Resume already rebuilds outputs,
                          outText and pass from the append-only run log and
                          runFrom skips every fired node. Ours is append-only
                          (ESTATE LAW 8); its UPSERTs one row.

Two of its files cannot run on this machine at all: `tui.py` dies on a bare
`import curses` (Windows Python has no `_curses`), and `LedgerDB.query()` does
not exist -- nor does the table `ledger_events`, the column `previous_hash`, or
the column `payload`. Four call sites in the CLI and TUI use them. Proven, not
read: `python cli.py ledger tail` raises AttributeError.

**WHAT WAS ACTUALLY WORTH TAKING WAS ONE IDEA: PER-STEP RETRY.** This engine
had none -- a node whose engine failed to answer ended the run, and a flow that
lost a race with a model still loading was simply dead.

    Node.Retries      0 to MaxRetries(5), and 0 is the default, so every flow
                      folded before this behaves EXACTLY as it did.
    retryWait         1s, 2s, 4s, then held at 8s. Backing off because the
                      usual cause needs a moment; capped because doubling
                      forever is the unbounded ticker ESTATE LAW 7 refuses.

**AND THE LINE THAT MATTERS MORE THAN THE FEATURE: RETRY ANSWERS AN ERROR AND
NEVER A VERDICT.** `rerr` is the engine failing to answer at all -- no door, a
dead socket, a template that will not render. A node that ANSWERED and was then
judged FAIL never reaches the retry loop, and `Validate` refuses `retries` on
`eval` and `gate` outright, by name, so the other reading is not available even
to someone looking for it. Re-rolling a check until it says PASS is the
laundering path this engine spent 2026-09-12 closing, and it would have arrived
dressed as a reliability feature.

**EVERY ATTEMPT IS WRITTEN DOWN, AND THE CLOCK KEEPS RUNNING.** Each abandoned
attempt gets its own line in `runs.jsonl` with its error and its number, and
`Status` renders it on the waterfall like any other node -- a silent retry
hides the flakiness it papers over. The attempt AND the pause after it are both
counted into `elapsed`, so a node that keeps failing runs the flow OUT_OF_TIME
rather than past it; every millisecond is accounted exactly once.

**THE GLASS OFFERS IT ONLY WHERE THE ENGINE ALLOWS IT.** A `retries` box on
ask, run, seat, prompt and memory -- and NOT on eval or gate, because a field
the engine will refuse is a field the glass must not offer. It is typed as a
NUMBER on the way out: Go unmarshals `Retries` into an int and rejects a string
outright, so a box left as text would not degrade, it would fail the whole
save. `version` had been special-cased inline for that exact reason; the set is
named now so the next number field cannot forget.

Six strokes, each both ways, and proved non-vacuous: with the retry loop
mutated to `false &&`, the two that assert retrying go red and the one that
asserts UNCHANGED behaviour without retries stays green, which is the pair
doing its job. Both Go modules green, gofmt clean.

### A version-control flow, and the council's only road to the new verb measured before it was built on

*Landed after the `v0.1.5` tag, in 6e8248a. First written inside the release's git_tag entry, and moved here on 2026-09-14 when 0.1.5 was marked released, so that section describes only what the tag holds. The words are unchanged.*

**AND A VERSION-CONTROL FLOW, BECAUSE THERE WAS NONE.** His words: *"add the
workflow if there is no version control workflow existing."* The engine held
`coder`, `smoke` and `version-bump` -- and `version-bump` only READS
pyproject twice around a gate. `version-tag` v1, seven nodes, 900s:

    read       run   list the marks in {{world}} and say what the ground
                     declares, which file said so, and what GitHub has
    judge      gate  stop if the number is wrong -- the version file is
                     bumped and SAVED before a mark is cut, never after
    cut        run   cut {{mark}} with the message {{what}}
    proof      eval  contains `Cut {{mark}} at`, scored on the EVIDENCE
    send_gate  gate  the irreversible half, named before it is taken
    send       run   send that one mark by name
    sent       eval  contains `Sent {{mark}} to origin.`

THE EDGE CONDITIONS ARE THE DESIGN. `proof -> send_gate` is `pass` and there
is NO fail edge anywhere, so a cut that did not happen ends the run FAIL and
the send gate is never offered -- the rule this estate learned the hard way
today, that no gate is offered for work that failed. Both irreversible steps
sit behind a `gate`, whose verdict is PAUSED: the engine prepares, the hand
decides (RULE 6).

**THE COUNCIL HAS NO TAG SKILL, AND THE ROUTE IT DOES HAVE WAS MEASURED, NOT
ASSUMED.** A `run` node goes to the COUNCIL -- manjuel's Router and its 42
skills -- and none of them is a tag. Its only road to this new verb is
`mcp_call`. That road was tested before a single node was written on it, with
the read-only action and no model in the loop:

    atlas git_tag {"project": "atlas", "action": "list"}   -> v0.1.5, 3dacdbc
    atlas git_tag                                          -> v0.1.11, c766ce7

So three of the four links are proven: skill -> door -> tool -> answer. THE
FOURTH IS NOT. Whether the Router reliably dispatches `mcp_call` and writes
the right JSON into `<content>` needs a model and a real mark, and there is no
mark left to cut today. **THE FLOW IS SAVED AND VALIDATED, AND IT HAS NEVER
BEEN FIRED** -- said here rather than discovered later, because a flow that
has only been read is exactly what this version is named after. The honest
first firing is the next number; a first-class `git_tag` skill in the core
would remove the uncertain link altogether, and that is a decision, not a fix.

### The release workflow built the spine in release and the door looked in debug

*Landed after the `v0.1.5` tag, in 94c085d. First written inside the release's git_tag entry, and moved here on 2026-09-14 when 0.1.5 was marked released, so that section describes only what the tag holds. The words are unchanged.*

**AND THEN THE TAG WENT, AND THE RELEASE WORKFLOW FAILED ON ITS OWN FAULT.**
`release.yml` was written this morning and had never been FIRED. The first
real tag through it, v0.1.5, died at `Go tests, both modules`:

    prove refused: no working atlas binary behind "atlas"
    build it first (cargo build -p atlas) or set ATLAS_BIN.

`Build the spine` runs `--release`, so the binary is in `target/release`.
`findAtlas` walks only `target/DEBUG` -- while `prove_test.go`'s skip check
deliberately looks in BOTH profiles, its own comment reading "a tree built
with --release would skip here and refuse there". So the stroke declined to
skip, the finder found nothing, and the leg went red. `prove.yml` never met
this because it builds the DEBUG profile, with a comment saying it builds
first precisely so that leg is an answer rather than a shrug.

The step now names `ATLAS_BIN` at the release binary and refuses outright if
it is missing. That is a better answer than also building debug: a release
run should prove THE BYTES IT IS ABOUT TO PUBLISH, not a second copy built a
different way. The disagreement between finder and skip-check is left alone
and named -- which profiles findAtlas walks is how the door behaves at
RUNTIME, and that is a decision rather than a workflow fix.

**WHAT THIS MEANS FOR v0.1.5, said plainly.** The mark is sound and the code
under it is proven: `prove.yml` passed on 3dacdbc, both jobs, and the core's
passed on c766ce7. What did NOT happen is the draft release -- the packaging
run failed, and it cannot be re-run green, because a dispatch uses the
workflow file AS IT STANDS AT THAT REF and the fix is not there. THE MARK IS
NOT MOVED TO FETCH IT: that is the rule this door enforces on everyone else,
and it is not waived for its author. The fix rides the next number.

## [0.1.5] — 2026-09-12 (tag on 3dacdbc)

0.1.3 and 0.1.4 were built and never tagged, so both ship inside this tag; their
headers below mark where each began.

*Post-script, 2026-10-09: both have marks of their own since, cut after the fact on his word on the commits that cut their numbers, v0.1.3 on `7688c04` and v0.1.4 on `467e94c` (every VERSION file and Cargo.toml agree at each), each with its message. The headers below stand as they were.*

### 0.1.5 — THE FLOW CONFIRMATION

His word, 2026-09-12: *"0.1.5 the flow confirmation."*

AND THE NAME IS THE FINDING. Every fault under this heading was a flow saying
a thing had happened when it had not, and every fix was the same move: make the
CONFIRMATION mean something. A check that asked whether code RAN and called
that correct. A marker that travelled from one node's objective into another's
prose. A marker a seat could simply write. A marker quoted out of the very
requirement it was meant to verify, and found inside the clause that said it
did NOT match. A repair nobody judged. And under all of it, a release script
that could not be parsed at all, reporting at the end that work was happening
somewhere it was not.

None of them was found by reading. Every one was found by FIRING the thing and
then refusing to believe the green.

### There was no way to cut a tag, and the panel is where he cuts them

His order, 2026-09-12: *"cut the tags through the dashboard"*, then *"add the
workflow if there is no version control workflow existing, use the system to
run the tag updates."*

**THE ANSWER TO THE FIRST HALF WAS THAT IT COULD NOT BE DONE.** Every path was
checked before anything was clicked, and none of them could cut a mark:

    the door        seven git tools -- git, git_diff, git_commit, git_push,
                    git_pull, git_branch, git_remote. No tag.
    git_branch      `action` is list | new | switch | close. Lines of work
                    only; it had never known what a tag was.
    the glass       Save the work · Send to GitHub · Take from GitHub ·
                    Lines of work.
    the council     the core's own skills are git_status, git_commit,
                    git_init, git_pull, git_push, git_cycle. No tag there
                    either.

So `release.yml` had been waiting since this morning for an event nothing in
this estate could produce.

**`git_tag`: list, cut, send.** One more verb in gitctl.go, under the header
that file already carries -- stdin closed, jailed to the tenant's Home, and
sending walled. A mark is the one artifact a stranger takes on faith: he
fetches v0.1.5 and believes it is 0.1.5 because the name says so, and unlike a
line of work a mark is not expected to move under him. Every refusal exists to
keep that sentence true.

    plain semver only    vMAJOR.MINOR.PATCH and nothing else. EARNED TODAY:
                         release.ps1 would have turned a plain 0.1.5 into
                         `v0.1.5+f1`. release.yml refuses that on arrival;
                         this refuses it before the mark EXISTS, which is the
                         half that matters once a mark has been fetched.
    the ground agrees    the name must equal the version the world declares --
                         VERSION, or pyproject.toml for the core -- and the
                         refusal quotes what the file actually says.
    AT THE COMMIT        and this is the whole design. The version is read out
                         of git at the commit being marked, never off the
                         disk. The disk is the TIP; a mark is often cut at an
                         older commit, and checking a tag against a VERSION
                         that moved AFTER it is exactly how a green check
                         passes a wrong tag. Same arithmetic release.yml does
                         by checking the tag out first.
    never moved          a name that already exists is refused, and told where
                         it points. Re-pointing a fetched tag is a force-push
                         whose victim never finds out: his clone keeps the old
                         object and agrees with nobody.
    unsaved work         a mark at HEAD over a dirty tree is refused -- AND
                         THAT CHECK COMES BEFORE THE VERSION ONE, deliberately.
                         The ordinary way to arrive here is to bump the version
                         file and forget to save it; checked the other way
                         round that person is told the mark and the ground
                         disagree, which is true and no help at all.
    one mark, named      send pushes `refs/tags/<name>`, never `--tags`. Same
                         class as `git push --all`, which is in CLAUDE.md
                         RULE 1 because it has already cost this estate
                         something: a verb that looks like it acts on the thing
                         you named quietly acts on all of them.

**THE STROKES WERE PROVED NON-VACUOUS, not just green.** A guard and its stroke
landing together have no red to point at, so the central guard was MUTATED --
`if "v"+declared != name` replaced with `if false` -- and both strokes went red
on exactly the right thing: a mark named v0.2.0 cut against a ground saying
0.1.5. Guard restored, green again. Ten strokes, each both ways.

**AND THE PANEL GREW THE BUTTON.** "Version marks" opens beside "Lines of
work": every mark with what it points at, when, and whether GitHub has it --
asked OF GitHub, and only while the wall is open, because a local repository
genuinely does not know and `sent_known` says which answer you are reading. The
name box is filled with the ground's own number rather than left blank, since
the one lawful answer is already known and retyping it is only a chance to typo
it. Sending asks first, and says what it is starting: where a release workflow
is set up, the tag arriving is what fires it.

**A FINDING THAT CAME FREE, AND IS NOT MINE TO FIX.** `prove.ps1` parses now
(it could not, until this morning) -- and it exits at STEP 1 OF 6. The Rust
`us_parity` test reads `estate/Agents/.us/Agents.us`, and `estate/` has never
existed in this ground: `git log --all -- estate` returns nothing, ever. So
steps 2 through 6 -- the Go tests, the Python verifiers, the 58 MCP strokes,
enrollment, the webapp build -- have never run from that script on this
machine. They were run by hand here instead, and all pass. Whether that fixture
should skip when its oracle is absent or the oracle should be brought in is a
decision, and the oracle is private material, so it is named and left.

**AND THEN THE PANEL SHOWED THE WRONG SHA, WHICH IS WHY YOU FIRE IT.** The
first mark was cut with every stroke in this file green, and the glass
printed `v0.1.5 | 1beae6d` -- while the mark was on `ae31e7c`. An ANNOTATED
tag is its own object with its own sha, so `%(objectname)` is THAT, not the
commit; the column means "what it points at", and the number under it was one
the reader would not find anywhere in the log. `%(*objectname)` dereferences,
and the plain field stays as the fallback for a lightweight tag, which is
already its own commit.

The stroke came after the sighting and was proved against it: mutated back to
`f[1]`, `TestTheColumnNamesTheCommitAndNotTheMarksOwnObject` goes red naming
both shas. It also pins the OTHER half -- that the mark really is annotated --
because on a git that quietly made a lightweight tag there would be no second
sha to get wrong, and the message the operator wrote would be stored nowhere.

The mark was cut, seen to be wrong, and REMOVED BEFORE IT WAS SENT, so nothing
outside this machine ever saw it. Removing it took a git command, because the
panel deliberately has no un-cut button: that is a real gap and it is named
rather than filled, since a verb that deletes marks is a different decision
from a verb that cuts them.

**AND THE SEND BUTTON DID NOTHING AT ALL, WHICH WAS THE THIRD SIGHTING.**
Sending a mark is the one act on this page that starts something on the far
side, so the first cut asked with a `confirm()`. It was DISMISSED WITHOUT
EVER BEING SHOWN -- the button was pressed, the panel still said "only on
this machine", and no refusal appeared anywhere, because a dismissed confirm
is not an error. A browser dialog is the one piece of interface the glass
does not control, and nothing else on this page opens a modal.

It ASKS IN THE PANEL now: the first press ARMS the button -- it says "Click
again to send", and the answer line says what sending starts and that a mark
someone has fetched cannot be pulled back. One mark is armed at a time by
construction, and any repaint of the card disarms, so a flag cannot outlive
the row it was set on and arm a button he never pressed.

That is three faults in one afternoon found by USING the thing, with every
stroke green throughout: a sha that pointed at nothing a reader could find, a
button that did nothing, and before either of them a panel that had no such
button at all. It is the same lesson this version is named for.

Both Go modules green and gofmt-clean; the door's own battery PROVEN.

### The glass had one test function in 2,800 lines, and a release workflow that was only ever claimed

Two of the review's open items, closed.

**THE GLASS: 1 test function -> 15, over three packages.** `webapp` is ~2,800
lines across eight packages and had exactly one stroke in all of it
(`handlers/ws_test.go`). ADR-006 made the same measurement about the door --
the protocol layer and the tenant model were its two least-tested things, and
both were given first strokes on 2026-09-11. The glass never had that pass,
and the operator named it the one thing to fix before field testing.

    db/db_test.go        the round trip survives reopening (every write calls
                         Save, and this is where one that does not shows); the
                         save is tmp-then-rename with no scratch left behind;
                         newest-first and the limit; upsert replaces rather
                         than duplicates; GetAgents hands back a COPY so a
                         caller's mutation cannot reach the store
    handlers/wall_test   `visible` in all four cases, including the widest
                         thing it does -- AN UNTENANTED ROW IS VISIBLE TO
                         EVERY TENANT, deliberate so the glass does not blank
                         on the day auth goes on, and pinned rather than
                         assumed; the three filters drop another tenant's rows
                         and return [] not nil; /health hides global counts
                         under auth and still answers
    server/server_test   every embedded file gets a QUOTED validator, same
                         bytes same tag; the session gate's open paths, its
                         JSON 401 for /api/*, its redirect for a page, and
                         /metrics and /ws both behind it; a refused request is
                         still COUNTED, because a counter that only saw
                         successes would hide the burst worth seeing

**AND A FINDING THE STROKES TURNED UP: THE AUTH GATE CANNOT BE TURNED ON.**
`ConfigureAuth` has NO CALLER anywhere in the module. `main.go` never calls
it, so `authOn` is false for the life of every process, `h.sessions` is never
initialised, and `ATLAS_AUTH=1` -- documented in two separate comments -- is
read nowhere. The gate in `server.gated` is correct and inert. Both states are
pinned: `TestWithAuthOffTheGateIsInert` is what production actually does, and
`TestWithAuthOnTheGateRefusesTheRightThings` configures it by hand to show the
mechanism was already right. **The only missing piece is the wiring, and
wiring it is a decision** -- the env var, the session path, the service wire --
so it is named here rather than invented.

**Not reachable, said plainly rather than worked around:** the route table is
built inside `ListenAndServe`, which then binds a port, so the mux itself
cannot be tested without refactoring that function. `.woff2` MIME registration
lives in `main.go`'s `func main`, which has nothing to call.

**THE RELEASE WORKFLOW: WRITTEN, NOT DELETED.** `DELIVERABLE.md` named
`.github/workflows/release.yml` for weeks and it had never existed;
`release.ps1` signed off with "GitHub Actions will build binaries and create
the release" while nothing did. The claim is now true.

`release.yml` fires on a `v*` tag and **proves before it publishes** -- a tag
is the one artifact a stranger takes on faith, so the battery runs on the
TAGGED commit rather than trusting whatever was green when it was cut:

    the tag equals VERSION   three things must agree -- the ref, the file, and
                             plain semver. A tag carrying a build moniker, or
                             pointing at a commit nobody bumped, dies before a
                             single artifact exists. `release.ps1` would have
                             produced `v0.1.5+f1` from a plain 0.1.5 two days
                             ago; this catches that class by arithmetic.
    ten pins in sync         version.ps1 sync, which held six until today
    the battery + both Go    a red battery publishes nothing
    every binary ASKED       six of them, one at a time, measured not asserted
                             -- atlas-tui prefixes its own name, and until
                             today it and atlas-vc printed hardcoded literals,
                             which are exactly the two a released artifact
                             would have carried
    a DRAFT release          RULE 6. Everything above the line is proof;
                             releasing is his click. `workflow_dispatch` runs
                             every proof and stops there, so this file can be
                             exercised without cutting anything.

It does not cut the tag. No agent tags.

### The release path could not run, and had not been able to for two days

**Three faults, stacked, each hiding the one under it.**

**1. It would not PARSE.** `release.ps1` is UTF-8 with NO BOM and carried seven
em-dashes. PowerShell 5.1 reads a BOM-less file as ANSI, so each dash became
three bytes of nonsense and broke the string on the `<ver>` help line — *"The
'<' operator is reserved for future use."* Proven by parsing HEAD's own bytes:
re-encoded with a BOM it parses clean; as it actually sits on disk it does not.
`prove.ps1`, which step 1 of 4 calls, had the same single dash and the same
fault. **Neither local script has ever run on this machine.** Both are pure
ASCII now, so no BOM has to survive a future edit.

**2. Under that, the moniker.** `version.ps1` was cured of the stone tag on
2026-09-10 — *"remove the moniker for the stones, no letters in my versions"* —
and **its only caller was not**. `release.ps1` computed `0.1.5+f1` and handed
it to `version.ps1 set`, which now REFUSES a tag. So it did not merely cut a
bad tag: it would have died at step 2 of 4, *after* running the whole prove
suite. One half of a fix landed and the other half was left holding the bug.

**3. And it lied on the way out.** Its closing line said *"GitHub Actions will
build binaries and create the release."* Nothing does — `.github/workflows/`
holds one file, `prove.yml`, and it runs on push/PR, not on tags. It now says
so plainly. A script that tells the operator work is happening elsewhere, when
it is not, is worse than one that says nothing: he stops looking.
`--dry-run` was never a thing either; PowerShell takes `-DryRun`.

**Fixed — and `version.ps1` now moves every pin, not six of ten.** Its list
held six VERSION files. The two it was missing were exactly the two that
printed a hardcoded literal until this morning (`atlas-vc` had no VERSION file
at all; the webapp said `0.1.3` by hand in `/health` and in the Prometheus
gauge). Cargo.toml and `core/src/version.rs`'s assertion were a WARNING telling
the reader to run a stroke **that does not exist** — there is no `version-cross`
leg in `tests/prove.py`. They are moved now, not mentioned: a bump that left
them behind would ship a Cargo.toml disagreeing with every binary, and
`release.ps1` would commit and tag it.

`sync` checks all ten and says so: `All 10 pins in sync: 0.1.5`.

**A REGRESSION I CAUSED AND REVERTED, recorded because it is the lesson.** The
pin rewrite used `Get-Content`/`Set-Content` and CORRUPTED
`core/src/version.rs` on its first run: PowerShell 5.1's `Get-Content` read the
BOM-less file as ANSI, so every multi-byte sequence became separate Latin-1
characters — `§` to `Â§`, `—` to `â€"` — and `Set-Content -Encoding utf8` then
re-encoded that mojibake *and* added a BOM. Reverted from HEAD and rewritten
with `[IO.File]::ReadAllText` plus `UTF8Encoding($false)`. The diff is now two
version lines and nothing else, no BOM, `§` and `—` intact. **A version bump
must move a version and touch nothing else** — and the fault was the same
encoding trap as fault 1, one layer down.

**Measured.**

    release.ps1 0.1.5 -DryRun   Target 0.1.5   Tag v0.1.5
    release.ps1 patch  -DryRun  Target 0.1.5   Tag v0.1.5
    release.ps1 0.1.5+f1        Refused: carries a build tag.  exit 1
    all three .ps1              parse clean on PowerShell 5.1
    five binaries + the glass   answer 0.1.5, asked one by one

**And DELIVERABLE.md was a 0.1.2 / 2026-09-08 snapshot end to end** — header,
Version field, and a binaries table claiming 0.1.2 for six binaries while
omitting `atlas-vc` entirely. Its own usage block still taught
`.\release.ps1 0.1.2+f2`, the banned moniker, in the delivery instructions. The
inventory called `.github/workflows/` "2 workflows: 6 parallel pipelines" and
`line/` "92 tests" (it is 229 fns over 19 pkgs). Every version claim is
measured now; the 112/112 results and the git-history block are LABELLED as
0.1.2-era fossils rather than swept, because they are the record of what that
release proved.

### The repair path is judged too, and the verdict now means something

`recheck -> land always`. A run that failed its requirement, repaired and
rechecked arrived at the hand with **no judgement of the repaired work** — the
same fault as the original green-on-wrong-code, moved one edge down. It was the
last of the three holes that firing this flow found, and the only one left open
when the other two landed.

**Fixed — `proof`, a second eval on `recheck`, holding the repair to the SAME
expectation.** `coder` is at v12, eight nodes.

**IT HAS NO FAIL EDGE, AND THAT IS THE DESIGN.** An eval that fails with no
fail edge stops the run (`run.go`: `if !hasFailEdge(...) { return VerdictFail }`),
and there is nothing to steer to anyway — the retry is UNROLLED, so there is no
second repair. Work that still does not meet the requirement must not be
OFFERED for landing. So the verdict carries information it did not before:

    PAUSED   it passed, and your hand decides
    FAIL     it did not, and no gate is offered for it

Nothing is thrown away either way: every attempt stays in the workspace with
what each run said.

Strokes: `TestARepairThatWorksReachesTheGate` — judge fails, repair and recheck
fire, proof passes, the run PAUSES at the gate with all four named;
`TestARepairThatDidNotWorkNeverReachesTheGate` — the same path with a repair
that did not fix it FAILS, never fires `land`, and never sets a paused node.
The flow lives in gitignored runtime state, so the strokes are how this shape
travels at all.

**Also measured live, and it closes a gap in the previous landing's own
judgement.** That landing noted the pass side had never been proven on real
code — every live run had ended verdict-fail. It has now: objective "print the
6th Fibonacci number", `FIB6: 8` expected, and the coder got it first try —
artifact prints `FIB6: 8`, the evidence block carries it, verdict `pass`, gate.

### An eval scores evidence, and prose is not scored at all

**The run that forced it.** With the correctness check in, the coder was asked
for a script printing `FIB6: 8`. It wrote a fibonacci that never sets
`fib_sequence[1]`, so it printed `FIB6: 0`. **The verdict passed.** The seat
had reported the failure perfectly — *"printed `FIB6: 0`, which is not the
expected output of `FIB6: 8`"* — and a `contains "FIB6: 8"` found the marker
inside the clause saying it did not match.

Requiring the verdict block (landed an hour earlier) did not stop it: the block
was present, because `run_python` really had run. **Evidence that something ran
is not evidence that the marker came from what ran.** And the second road is
worse — naming the marker in the objective puts it in the brief, the brief puts
it in the node's objective, and the seat quotes it. The guidance to name the
marker was manufacturing the false pass.

**Fixed — for a `run` node, the check reads only the machine's lines.**

    evidenceOf()     the text after the verdict marker, and whether there is
                     any. The LAST marker wins -- belt to the braces below.
    flow/run.go      an eval judging a `run` node scores that and nothing
                     else. The answer stays whole as what a person reads at
                     the gate; it is simply no longer scoreable.

**And the marker itself is now unforgeable.** `appendVerdicts` used to SKIP
when a marker was already present — "the news reaches the gate exactly once".
It does, but a SEAT can write that string, and a seat that did would have
suppressed the machine's block and left its own words sitting exactly where an
eval now reads evidence from. It is strip-then-append: every seat-written
marker line is removed, then the real block is written. A turn with no verdicts
leaves no marker at all, because a suppressed block is indistinguishable from
an invented one.

`TestAppendVerdicts`' once-only assertion was REWRITTEN the same day it was
written: it held the IMPLEMENTATION (byte-identical on a second call) where the
property is what matters (exactly one marker, the real lines under it). The
guard is the same guard.

Strokes: `TestTheVerdictScoresEvidenceAndNotProse` carries the real answer —
requirement quoted, failure honestly reported, evidence saying `FIB6: 0` — and
must fail, with the right-output twin passing so the rule did not become a
refusal of everything. `TestAnObjectiveCannotSatisfyItself` holds the other
road. `TestARunNodeThatCalledNoToolCannotBeJudged` and
`TestAVoiceIsJudgedWithoutToolEvidence` still stand: an `ask`, `prompt` or
`memory` node holds no tools by definition and is judged on its answer, which
is the only thing it has.

Measured live, re-firing the exact case that had passed:

    prose     contains "FIB6: 8"   True
    evidence  contains "FIB6: 8"   False
    evidence  run_python: RAN: fibonacci.py / --- stdout --- / FIB6: 5
    verdict   fail -> repair -> recheck -> land (gate)

**Still open, and it is the last of them:** `recheck -> land` is UNJUDGED. Only
`verify` is scored, so a run that fails the verdict, repairs and rechecks
reaches the gate with no judgement of the repaired work.

### No evidence is not a verdict

**How it was found: by the run immediately after the correctness check
landed.** A `verify` node came back in **5.6 seconds** with no verdict block at
all — it had called no tool — and the eval failed it. Taking the fail edge was
right. Recording it as `fail: expected contains "X"` was not: failing because
the work was WRONG and failing because NOBODY WATCHED THE WORK are different
facts, and a reader at the gate could not tell them apart.

**And the pass side was the real hole.** The marker a check hunts is a string,
and a seat can WRITE the string without anything having run. That is the same
laundering that made a pasted `RAN:` pass earlier the same day, arriving by a
new road: not a marker that travelled between nodes, but one a seat simply
asserted. A check that accepts it is scoring testimony again (LAW 5).

**Fixed — an eval judging a `run` node requires the machine's own block before
it judges at all.**

    play.ToolVerdictHead   the marker moved to `play`, which imports nothing
                           of ours. `tools` WRITES the block and `flow` now
                           READS it; `tools` imports `flow`, so it could live
                           in neither, and two copies of the string would be
                           one rule with two spellings.
    flow/run.go            execNode takes the spec's nodes, so an eval can ask
                           what KIND the node it judges is. A `run` node whose
                           answer carries no block returns
                           `fail: NO EVIDENCE -- <node> is a run node that
                           called no tool`, and says nothing was judged.

**IT CANNOT PASS, and that is the point rather than a nicety.** Requiring the
block means the evidence was machine-emitted from the tool results, not typed
by the seat being judged.

**ONLY FOR `run` NODES.** An `ask`, `prompt` or `memory` node holds no tools by
definition, so demanding tool evidence there would refuse every honest eval
over a voice — `branchSpec`'s does exactly that, and a stroke holds it.

Strokes: `TestARunNodeThatCalledNoToolCannotBeJudged` — the 5.6-second turn,
whose prose literally contains `RAN:` and must still fail, with the record
naming why; `TestEvidenceLetsTheJudgementStand` — the same answer WITH the
block passes, and a witnessed failure still fails, so the rule did not become
a rubber stamp; `TestAVoiceIsJudgedWithoutToolEvidence` — the ask-node eval is
untouched. The stub engine gained a canned `Turn` answer so the rule could be
struck both ways; with none it returns the old string and every existing
stroke reads as it did.

**Also — the builder now states the discipline it had left implicit.** The
expectation box said only "what expect is, for this run", and a run whose code
was CORRECT failed because the hand wrote `Refused` where the program printed
`Refusing`. The block now says the test is exact and case-sensitive, and that
the marker belongs in the objective and repeated in the box. A looser test
would go green on work that only sounded right, which is the thing all of this
exists to refuse.

**Still open from the same run:** `recheck -> land` remains UNJUDGED. Only
`verify` is scored, so a run that fails the verdict, repairs and rechecks
reaches the gate with no judgement of the repaired work. That is the last of
the three holes firing this flow found, and it is not closed here.

### The check scored liveness and called it correctness

**How it was found: by firing it.** The coder flow was given a real task — "turn
a version string into a release target, and it must REFUSE any string carrying
a plus build tag INSTEAD OF STRIPPING OR DEFAULTING IT". The coder wrote
`version_str.split('+')[0]`: it stripped, the one behaviour the objective named
and forbade. **The flow went green.** `check` had asked whether `run_python`
said `RAN:`, and it had. The delivery even wrote both halves of the
contradiction in one sentence: *"It refuses to process any string carrying a
plus build tag, effectively ignoring it."*

A flow that goes green on wrong code is worse than one that goes red. The green
is the thing a reader trusts, and it launders wrong work to a gate.

**Fixed — an eval's `expected` is rendered, so a check can hold a node to an
expectation supplied at FIRE time.** It was a literal, so a check could only
ever test something written when the flow was FOLDED. That is enough for
liveness and cannot express correctness, because what a correct run prints is a
fact about THIS request.

    flow/run.go      play.Render on nd.Expected before scoring; the fail
                     message carries the RENDERED want, because "expected
                     contains {{expect}}" tells a reader nothing
    scoreNode        takes `want` already rendered -- the caller owns the
                     templating so this stays one question
    workflows.js     openVars scans `expected` too, and the comment that said
                     it deliberately did not is corrected in the same stroke.
                     A gate's title is still excluded: gates never reach
                     execNode, so a box for it would fill nothing

**THE EXPECTATION COMES FROM THE HAND, NOT A MODEL.** A model that both states
what correct output looks like and writes the code can agree with itself, and
agreeing with itself is the disease. `play.Render` refuses a missing var, so a
flow that templates an expectation nobody supplied fails at that node instead
of scoring against an empty string.

Strokes: `TestExpectationComesFromTheFiring` (met -> pass; **ran and WRONG ->
fail**, which is the run above), `TestAnExpectationNobodySuppliedIsRefused`,
`TestTheFailMessageNamesTheRenderedWant`, and
`TestAFoldedLiteralExpectationIsUnchanged` -- because every spec written before
this carries a plain literal and must score exactly as it did.

**Measured live.** `coder` v11: the same objective that went green now fails.

    verify    run_python: RAN: version_to_release_target.py   (exit 0)
    verdict   fail: expected contains "Refused"
              -> repair -> recheck -> land (gate)

**THREE WIRING ATTEMPTS, EACH CORRECTED BY A RUN, AND THE LESSONS ARE THE
VALUE:**

    v9   check(liveness) gated verdict(correctness). Wrong: a task whose
         correct behaviour is a NON-ZERO EXIT fails the liveness gate. A
         correct refusal exits 1.
    v10  check and verdict both hung off verify as recorders. Wrong, and the
         engine says so plainly: `if !hasFailEdge(...) { return VerdictFail }`
         -- AN EVAL IS A GATE, NEVER A PASSIVE RECORDER. Removing check's fail
         edge made the whole run FAIL before verdict could fire.
    v11  one eval, on the requirement. Liveness is not a second gate; it is
         evidence, and it is already in the verdict block a reader sees.

**And the expectation is only as good as the observable the objective names.**
A run where the code was CORRECT still failed the verdict, because the
expectation said `Refused` and the program printed `Refusing`. `contains` is
exact and case-sensitive and did what it was told. The cure is a spec -- the
objective naming the marker, the expectation matching it -- not fuzzy matching,
which would reintroduce the laundering this exists to stop. The brief now also
states that the file is run with NO arguments, after a coder wrote an
`sys.argv`-reading script that `run_python` structurally cannot invoke.

**Open, named, not fixed — two holes this found and did not close:**

    the retry is unjudged   `recheck -> land always`. Only `verify` is judged,
                            so a run that fails the verdict, repairs and
                            rechecks reaches the gate with NO correctness
                            judgement of the repaired work. Same shape as the
                            original bug, one branch over. Wants a second eval
                            after recheck.
    verify can run nothing  a verify node came back in 5.6s with no verdict
                            block at all -- it called no tool, so the verdict
                            had nothing to judge and failed for want of
                            evidence rather than for wrong work. The two are
                            not the same and the record should not conflate
                            them.

### 0.1.4 — THE DELIVERY PACKAGING

His word, 2026-09-12: *"atlas 0.1.4 - the delivery packaging."*

`VERSION` is the single authority and it is now true in all EIGHT files,
not six: the root, `line/`, one beside each of the five `line/cmd/*`
mains, and a new `webapp/handlers/VERSION`. Two of those did not exist
this morning, and their absence is the whole entry:

    atlas-vc    printed the literal "0.1.3" unconditionally and had no
                VERSION file at all -- the fifth command, outside the
                scheme the other four were in, while ACCEPTANCE.md
                counted "all 6 files" as though it were not
    atlas-tui   had a VERSION file beside it and did not read it: it
                asked the Rust spine for --version and fell back to a
                LITERAL when the spine was absent. So it printed a stale
                number on exactly the machines where the spine is not
                built -- the fresh clone `prove.py` keeps an ABSENT
                branch for, and the one case ACCEPTANCE.md is written to
                catch. The one place the staleness could not be noticed
                was the one place it lived.
    webapp      said "0.1.3" by hand in TWO unrelated functions --
                /health's `version` field and the Prometheus gauge
                `atlas_server_info{server="atlas-webapp",version=...}`.
                A version a monitoring system scrapes is the last one
                anybody re-reads.

All five commands and the glass now answer from the file: measured, not
asserted -- `atlas-mcp 0.1.4`, `atlas-tui atlas-tui 0.1.4` (with the
spine absent, which is the fixed path), `atlas-vc 0.1.4`, `atlas-door
0.1.4`, `atlas-town 0.1.4`. `Cargo.toml` and `core/src/version.rs`'s own
assertion moved with them, and every doc line that ASSERTS what the build
prints -- ACCEPTANCE, PIPELINES, OLLAMA_PROVER, WORKFLOWS, E2E_SCENARIOS,
AGENTS -- moved too. The released headers below did not: they are the
record of what shipped.

### The gofmt gate repeated the fault it was named after

Its own note, written this morning: *"three files had never been through gofmt
and nobody knew, because the check that found the first one was scoped to a
single directory."* It was then written with `working-directory: line`.

SIX FILES IN `webapp/` HAD NEVER BEEN FORMATTED — `db/db.go` and
`handlers/{flows,prompts,session,team,ws}.go`, untouched since they landed in
`77162d9` on 2026-09-09 — and the gate could not see one of them, because
webapp is a separate Go module. Found by running gofmt by hand during a
function check, not by the gate.

Five were struct-tag alignment and nothing else; their token streams with all
whitespace removed hash identical before and after. `ws.go` was NOT: gofmt
1.19 and later read the hanging indent under `Frames out:` as a CODE BLOCK and
would have reflowed it to a tab and split the label from its own list. That
comment was reflowed flat by hand instead, so it says the same thing, reads the
same way, and gofmt now has nothing left to do to it — checked, zero lines.

THE GATE NOW NAMES BOTH MODULES, and deliberately does not just run `gofmt -l .`
from the root: that walks `target/` and every vendored tree, which is how a
gate gets slow and then gets deleted. It also reports both before exiting, so
one run names every offender instead of one per push. A third Go module has to
be added to that list on the day it is created, and the note in the workflow
says so.

Proven to bite: an unformatted function was put into `webapp/handlers/team.go`
and the gate's own script exited 1 naming the file, then the file was restored
and it exited 0. `go build` and `go vet` clean across both modules.

### 0.1.3

His word: atlas is 0.1.3. `VERSION` is the single authority — `core/src/version.rs`
reads it with `include_str!` — but it is not the only pin, and `version.ps1`
says so itself in yellow every time it runs: it moves six VERSION files and
then names `python tests/prove.py` as the thing that finds the rest.

Fourteen real pins moved: the six VERSION files, the workspace package in
`Cargo.toml` (and `Cargo.lock`'s three crates, regenerated by cargo rather than
typed), the spine's own `assert_eq!`, `atlas-tui` and `atlas-vc`'s `--version`,
the glass's `/api/health` and its metric label, and the two test harnesses.
Then the live documents that ASSERT a version: `docs/ACCEPTANCE.md`,
`docs/PIPELINES.md`, `docs/OLLAMA_PROVER.md`, `tests/e2e/E2E_SCENARIOS.md`,
`AGENTS.md` and two `**Version:**` headers.

WHAT WAS NOT TOUCHED, AND WHY EACH ONE STAYS:

    core/src/version.rs:3    a doc comment reading "e.g. `0.1.2`" -- an example
    gitctl_test.go:146       `v0.1.2` is a BRANCH NAME in a name-law fixture
    docs/WORKFLOWS.md:408    sample message text inside an example payload
    version.ps1, release.ps1 help text and the record of the struck moniker
    DELIVERABLE.md           dated 2026-09-08, with its own "verified" table
                             and a frozen git-history block -- a snapshot OF
                             0.1.2. Renaming it would claim it was proven at
                             0.1.3, which it was not.

Proven by the spine's own battery rather than by reading: `version-pin
file=0.1.3 bin=0.1.3` and `version-cross all agree: 0.1.3`. Both binaries were
rebuilt and restarted on their own command lines and now answer 0.1.3 — the
door carrying exactly the two worlds it carried before, atlas and research.

### The workflow builder gets its face back

The DAG builder came off the panel 2026-09-10 — "not used, wipe it" — and
`/flows` became Version control. What was wiped was THE PAGE. The engine was
never touched, and this is what was sitting behind the missing page the whole
time: `internal/flow` at 1,140 lines with 455 lines of strokes, ten `flow_*`
tools on the door, nine handlers, all nine routes wired in `server.go`, and a
method in `api.js` for every one of them. Proven live before a line of the page
was written: save -> list -> get -> run -> status, a real llama3.2 call,
receipts, budget bar, verdict COMPLETE.

So nothing here rebuilds a workflow system. `static/js/workflows.js` DRAWS the
one that was already there, and every button is a call the LINE already answers.

THE PATTERNS, in this engine's own terms (his reference,
workflowbuilder.io/blog/agentic-workflow-patterns):

    prompt chaining     nodes joined by `always` edges; topo order is the chain
    routing             an `eval` node, then `pass` / `fail` edges off it
    parallelization     branches declared, run sequentially -- one card, one
                        rack queue, and interleaved model output cannot be read
    reflection          unrolled into fixed passes, because Validate REFUSES
                        cycles. That is what the article recommends anyway.
    human-in-the-loop   a `gate` node: the run stops at PAUSED and waits

All seven node kinds are offered and no eighth is, because `Kinds` is a closed
set and anything else is refused by name at save. The page does not pre-judge a
spec: Validate lives in the LINE, its refusals name the node and the reason, and
a second opinion drawn here is how two validators drift apart. A gate offers
exactly the two moves the engine accepts, and neither is the default.

ITS OWN PAGE, NOT THE OLD ROUTE. `/flows` stays Version control — that is the
git overwatch he uses every day, and taking the route back would cost him the
page he actually stands on. The builder is `/workflows`, its own line on the
panel.

Proven from the browser, not from curl: a two-step flow built in the UI, folded
to v1, fired, `COMPLETE · fired 2 · 57290ms` with a receipt per node.

**Restart required** — the webapp embeds its static files (`go:embed`), so
`:8091` does not carry this until it is rebuilt and restarted. Tested on a
scratch binary on `:8099`; his running server was not touched.

- **`flows/` is ignored.** The engine writes `<home>/flows/` the moment a flow
  is saved or fired — specs, their folded history, `runs.jsonl`. Runtime state,
  the same shape as `state/` beside it, and the provers build their own homes in
  temp dirs rather than reading it.

### The seat log travels, so a clone can prove itself

`atlas/SEAT_LOG.md` was gitignored and never reached a clone, so `cargo test
--workspace` failed there on the `orient-pack` stroke (`LOG=false`) while
passing on the ground. CI never saw it: the gate runs `prove.py --check`, which
skips the cargo leg. Untracked from `.gitignore` on his word, 2026-09-11.

### The criteria documents, measured rather than remembered

`docs/ACCEPTANCE.md` and `docs/PIPELINES.md` are what a stranger reads to learn
what PASS looks like here, so a stale number in them is not cosmetic — it is a
gate reporting the wrong verdict. Every count in both was run rather than read:

    atlas-mcp --prove         said 58 strokes        is 125        CORRECTED
    cutters answering verify  said "all 17 pass"     is 24 legs:   CORRECTED
                                                     12 byte-identical,
                                                     12 ABSENT (oracle)
    agents/docs/*.md          said 40 doc files      is 41         CORRECTED
    atlas-town --prove        said 11 strokes        is 11         held
    atlas-door --prove        said 13 strokes        is 13         held
    agents/*.us               said 40                is 40         held
    agents/modules/*.us       said 4                 is 4          held
    skills/*/SKILL.md         said 3                 is 3          held
    go test ./...             said 92+               is 138 across 19 pkgs, held

The cutter line was wrong twice over, which is why it did not become "all 24
pass". Half of them cannot verify on any machine but the one they were cut from:
twelve re-cut byte-identical and twelve report ABSENT naming `ATLAS_ORACLE_ROOT`,
because the private oracle ground does not ship and never will. A gate that says
"all pass" over that is unpassable by construction, the same defect as the
`--describe` row fixed earlier today. It now states both verdicts, which is the
three-verdict doctrine this repo already holds everywhere else.

Untouched on purpose: `DELIVERABLE.md` is dated 2026-09-08 and carries its own
"verified 2026-09-08" table and a frozen git-history block — a release snapshot,
not a living description, and its numbers are the record of that day.

### Two checks stop naming a command that cannot answer them

`docs/ACCEPTANCE.md` and `docs/PIPELINES.md` both gated on `atlas-mcp
--describe` listing 72 tools. It does not list tools. `--describe` prints one
line — the door's name and what it is — and returns (`cmd/atlas-mcp/main.go`),
and no flag the door declares lists them at all. So the gate could never pass at
any number, which is why renumbering it to 78 with the rest of the sweep was
refused: a fresh coat of paint on a check that does nothing is worse than a
stale one, because the stale one still reads as suspect.

Both now name `atlas-mcp --prove`, which really counts. `prove.go` builds the
surface in-process, calls `tools/list` against it, and reports `surface carries
78 tools (>=20)` — no server, no model, which is what every neighbouring row in
those tables already assumes. PIPELINES' own failure policy for this pipeline,
"Tool count < 20 → fail", is that stroke's gate written out in prose; it had
simply never been pointed at the stroke. The LINE table keeps its eleven
stages, so the summary that counts them stays true.

Measured off the three batteries rather than adjusted:

    atlas-mcp --prove    125 strokes, surface carries 78 tools
    atlas-town --prove    11 strokes
    atlas-door --prove    13 strokes

Which names the next stale number and leaves it standing: both documents still
call `atlas-mcp --prove` 58 strokes. It is 125. Town's 11 and the door's 13 are
right as written.

### The last leg that called an absence a failure

The door's battery learned this doctrine at 07:00 on 2026-09-11. One layer up,
`tests/prove.py` had never learned it, and it was found the same way: by a
fresh clone.

`absent_dependency` recognises an absence only when a refusal NAMES A PATH — it
scans for one, checks it does not exist, and checks it lies outside the atlas
tree. `check_trade_parity` shells the Rust spine, so on an unbuilt tree it
refuses with a COMMAND instead ("REFUSED: build the binary first (cargo build
-p atlas)"), which that scan cannot see. It fell through to FAIL.

So the whole battery went red on any machine that had not yet run cargo build —
which is every fresh clone, and the first thing a second machine does. Measured
by parking the binary and re-running:

    before                       20 held - 14 absent - 1 broke - exit 1
    after                        20 held - 15 absent - 0 broke - exit 0
    after, with the spine built  21 held - 14 absent - 0 broke - exit 0

The new branch is gated on the binary being GENUINELY ABSENT, so it can never
turn a real break into an absence: with the binary on disk it is unreachable,
which the third line proves — that leg still reports the ORACLE absence there,
not this one. `leg_go` has applied the same rule two functions up since it was
written. This is that rule reaching the last leg without it.

- **The tool count was stale in two places.** `README.md` and `DELIVERABLE.md`
  said the door serves 72 tools. It serves 78, counted off the wire — `tools/list`
  and `GET /tools` agree, and every one of them carries a description and an
  inputSchema. The table under DELIVERABLE's heading lists 24 of them and always
  did; only the count moved.

### The first CI run found two, on its first try

atlas's CI went up and the Linux probe went red immediately — which is what
that job is for. The gate passed (the battery, green on Windows, 1m47s). The
probe found one workflow error of mine and one real dependency on the
operator's platform.

- **The door's battery reported ABSENT as FAIL.** Every stroke in
  `cmd/atlas-door` goes through the Rust spine, so on a machine where it has
  not been built there is nothing to prove and nothing has gone wrong. It said
  FAIL. That made a fresh clone look broken: the first thing `go test ./...`
  said on a clean checkout was a red door, with the actual cause — one
  undocumented `cargo build` — buried in a refusal nobody read. Found twice in
  one day: on a clean clone of both repos, where it was the ONLY red in sixteen
  packages, and then by CI, whose Linux job cannot build the spine at all
  (`store/src/ffi.rs` links Windows' `winsqlite3`) and so could never have
  passed it. It now SKIPS with the command that would answer it, which is the
  shape `tests/prove.py` has held since it was written: ABSENT is not a pass
  and it is not a failure.
- **The stub engine was `cmd /c echo`, hardcoded, in two places.** `cmd` does
  not exist off Windows. Nobody could know: until today these suites had never
  run anywhere but one machine. `echoEngine()` picks by `runtime.GOOS`, and
  `splitCommand` takes both spellings to the same argv.

Proven both ways before it was pushed: **19 packages green with the spine
built, and 19 green with `ATLAS_BIN` pointed at a path that cannot exist** —
the closest proxy this machine can run for the Linux job. With the spine the
door PASSES; without it the door SKIPS and says why.


### atlas gets a CI, and it needs no network

Until today every claim in this repo rested on one machine — 78 tools, 19 green
packages, a 35-leg battery — while the core beside it has proved itself on four
matrix legs per push since 2026-09-03. A stranger had to take all of it on
faith. `.github/workflows/prove.yml` is the answer, in the shape the core's
workflow already established.

**It fetches nothing.** `Cargo.lock` holds three packages and all three are
this workspace's own; both `go.mod` files are stdlib only; `tests/prove.py`
imports nothing outside the standard library. Nothing after `checkout` touches
the network. Stated in the file as a PROPERTY, so a change that needs a
dependency is the regression rather than a surprise.

**Two jobs, and they are separate on purpose.**

- **The gate, on Windows.** Toolchains named in the log, `cargo build -p atlas
  --locked` first (the battery reports the spine ABSENT when unbuilt — honest,
  but then the door leg proves nothing), a `gofmt -l` that must come back
  empty, then THE BALL. It does not re-run cargo test or go test separately:
  the battery runs them itself, and two batteries that can disagree are worse
  than one that cannot.

  Windows is not a preference. `store/src/ffi.rs` does
  `#[link(name = "winsqlite3")]` against the Windows SDK with no `cfg(windows)`
  guard and no alternative backend — the store links the OS's own SQLite rather
  than vendoring one (ADR-004). On any other platform cargo fails at link time.

- **The portability probe, on Linux.** Its own job so a red there can never
  mask the gate. The Go half has no platform gating at all — no `_windows.go`,
  no `//go:build windows`, and `GOOS=linux go build ./...` was verified clean
  before the file was written. But the Go SUITES have never once run off
  Windows. They look portable (the Windows-shaped strings in them are test
  INPUTS, paths the wall must refuse) and looking portable is not being
  portable. **A red there is a finding, not a broken workflow** — it would mean
  the suites picked up a dependency on the operator's platform, which is
  exactly what the core's own matrix exists to catch.

Verified locally before it was committed: the YAML parses to two jobs,
`gofmt -l .` is empty, `cargo build -p atlas --locked` succeeds, and
`tests/prove.py --check` exits 0 with 21 held, 14 absent, 0 broke.


### gofmt comes back empty

Three files had never been through it: `internal/tools/gitctl_test.go` (a map
literal whose keys were padded to the wrong width), `cmd/atlas-vc/main.go` (a
doc comment in the pre-1.19 spelling, leading spaces where Go now wants a tab
block), and `internal/engine/engine.go` (one struct field padded for an
alignment group a comment had broken). Eleven lines between them, every one
cosmetic — no semantic change, 19 packages still green.

Worth naming because of how it was nearly missed: the first check ran
`gofmt -l internal/tools/` and found ONE file, so the report said one file. The
module-wide run found three. A scoped check answers the question it was
scoped to, not the question that was asked.

`gofmt -l .` now returns nothing, which is the precondition for making it a CI
gate rather than something a hand has to remember.


### ADR-006 accepted, and items 2 through 6 built

**Item 2 — the tier is a field, and the stroke corrected the ADR.** `Tool`
carries a `Tier`: CORE (nothing but a directory), SPINE (the Rust binary),
ENGINE (`--manjuel`). The zero value is CORE, so the common case is free and
the declaration stays honest rather than ceremonial.
`TestEveryCoreToolStandsAlone` calls every tool claiming CORE against a bare
temp tenant with the engine unwired AND the spine denied — and **found four
the ADR's grep had missed**. The table said 2 engine tools; there are **6**
(`env_open`, `env_close`, `run_start`, `run_answer`, `run_cancel`,
`ask_steward`), three of which refuse with *"no engine is open"*, a phrase the
grep was not looking for. The real split is **71 core / 6 engine / 1 spine**,
and the ADR was corrected to the measured number. Denying the spine mattered
too: `verify_chain` passed as CORE until `ATLAS_BIN` was pointed at a path that
cannot exist, because on a machine where cargo has run the walk simply finds
the binary.

**Item 4 — first strokes on `tenant` (288 lines, zero) and `ground` (200
lines, zero).** `ground` is the package that carried a world from outside the
estate onto the dashboard, and nothing was broken in it — Detect and Siblings
both did exactly what they were written to do. What was missing was any stroke
stating what that IS, so the blast radius of a launch directory was
discoverable only by suffering it. Now pinned: nearest ground wins; a `.us`
module outranks a folder name; Siblings scans **exactly one level** and a
grandchild is never carried; the attic and the vendored trees are never
grounds. And for `tenant`: an unknown project refuses BY NAME and hands a
stranger no context; case and space cannot fork a tenant; Home is absolute,
because every wall check downstream is a prefix test against it.

**Item 5 — one spawn contract.** Four seams, four different answers to the
same four questions. Two that actually bit: `gitstate` DISCARDED stderr, so a
git failure there was a shrug, and the spine call had **no timeout at all**, so
a wedged Rust binary hung the tool call and through it the door, forever
(ESTATE LAW 7 is bounded everything). `spawn()` closes stdin at one place,
bounds every child, keeps both streams whether it succeeded or failed, and
names a timeout as a timeout instead of "signal: killed".

Two of the six seams are deliberately left, and say why in the source:
`internal/engine` spawns over `StdinPipe` because that pipe IS the wire, and a
helper that closes stdin would break it by doing its job; and `cmd/atlas-door`
is a different binary whose sharing would need a new leaf package, because
`internal/tools` already imports `internal/engine`. A new package is a new
folder and folders are the operator's to place (RULE 8) — so it was named
here rather than invented.

**Item 6 — the oracle ground is nameable.** Every cutter computed its source
as `ATLAS.parent`. Before the 2026-09-10 split that was true; after it, the
parent is the manjuel core, so 14 legs reported ABSENT naming paths that **have
never existed on any machine**. Right verdict, phantom reason, nothing to act
on. They now take `ATLAS_ORACLE_ROOT`, falling back to the historical location
so nothing that worked stops working, and the ABSENT line names the dial. The
distinction that makes atlas standalone is written into `tests/PROVING.md`:
**the goldens travel, the cutters do not** — every golden is committed and the
Rust implementation is proved against them anywhere; only the RE-CUT needs the
private ground, and that ground never ships.

19 Go packages green, 0 failing. THE BALL: 21 held, 14 absent, 0 broke.


### The door names what it carries

It printed `carried 2`. A count cannot be checked against intent — two is two
whether the two are the ones you meant or not. It now prints:

    ground: atlas (...\Research\atlas via AGENTS.md); carrying 2: atlas, research

**Earned the same day, by the incident it would have prevented.** The door was
started from the ground root, so `ground.Detect` resolved the ground to
`research` and `ground.Siblings()` read the parent of THAT — the desktop —
carrying every neighbour holding an `AGENTS.md`. A world outside the estate
rode in, the dashboard read its git state, and the first anyone knew was a
sidebar badge reading 83,303. Every line needed to catch it at boot was already
there except the names.

**And the launch point is the whole control, which was got wrong out loud
first.** `Detect` walks UP and stops at the FIRST marker it finds, and
`Siblings` scans the parent of what it found — not the parent of the working
directory. From `atlas\line` the first marker is atlas's own `AGENTS.md`, so
the scan is of `Research` and the door carries exactly `atlas, research`. From
the ground root the first marker is the core's, and the scan is of the desktop.
No flag, no env var and no code change was needed for the behaviour the
operator asked for; a different `-WorkingDirectory` was the entire answer, and
this hand told him otherwise before checking. RUNBOOK's own start line already
does `cd atlas\line` first.

### The first strokes on the MCP, and two bugs they caught immediately

ADR-006 measured the door and found the thing that explains a two-day failure
streak: `internal/httpserver` is 800 lines, it is what makes atlas-mcp an MCP
server rather than a pile of functions, and it had **zero tests**. The tool
bodies behind it carry 755 strokes; the wire in front of them carried none.
Twelve strokes now stand there, hermetic — a real registry over a real temp
tenant that deliberately does NOT look like manjuel, asserting only what the
wire says.

- **A tool's own words outrank the Go error** (ADR-006 item 1). The tools/call
  error branch discarded `out` one line before it would have been sent and
  returned `err.Error()` instead — which for anything shelling a subprocess is
  the bare string `exit status 1`. `verify_chain` exposed it: it captures the
  Rust spine's `CombinedOutput` INTO `out`, so the diagnosis was in hand and
  thrown away. **Protocol-level**: every tool that errored was losing whatever
  it had written. `isError` still tells the caller it failed; it now also says
  why.
- **THE REGRESSION STROKE WAS WORTHLESS UNTIL IT WASN'T.** It passed with the
  fix reverted, which means it proved nothing. The reason was a second bug, in
  the resolver written earlier the same day: `findAtlas` walked from `"."`, and
  `filepath.Dir(".")` is `"."` — so that root broke out of the walk on its
  first step and never climbed. The door never noticed because its own
  executable path walks fine; a TEST binary lives in a build temp dir, where
  that root was the only one that could reach the tree and it was the one doing
  nothing. So the spine was never found, `out` was empty, and both branches
  produced identical text. Fixed to walk from an absolute cwd. **The stroke now
  fails with the fix reverted** (`the caller got only "exit status 1"`) and
  passes with it applied, which is the only thing that makes it a test.

**What the twelve pin:** the protocol version is a promise to every client
(`2025-06-18`); a notification is answered with silence; every advertised tool
carries a description and an inputSchema and never requires a property it does
not declare; `-32601` for an unknown method and `-32602` for an unknown tool,
each naming what it did not know; `-32700` for a malformed line; the id comes
back; a tool's refusal is `isError` and NOT a transport error, because a client
that retries transport errors would retry a refusal forever.

**And one the review found that the ADR missed:** `GET /tools` and `tools/list`
are two copies of the same 28-line rendering. A stroke now holds them
byte-identical, because a REST caller and an MCP caller disagreeing about what
a tool takes is the worst kind of quiet.

17 Go packages green, 0 failing.


### The spine is found, not assumed

`verify_chain` was dead on every machine where the Rust spine had been built
but not installed on PATH — which is every fresh clone. `--atlas-bin` defaults
to the bare word `"atlas"`, `atlas-door` walked the built tree to find it, and
`atlas-mcp` never did: the same estate answered differently depending on which
door you came through. `tests/PROVING.md` has named this a trap since it was
written, and RUNBOOK's start line still did not carry the flag.

- **The walk moved to `internal/tools`**, the one place that actually shells
  the binary, so every caller gets it and there is no second copy to drift.
  Order: an explicit `--atlas-bin` that is not the placeholder, then
  `ATLAS_BIN`, then PATH, then the built tree.
- **It starts from the running binary's own location.** A first cut walked up
  from the tenant home and from cwd; for `atlas-mcp` the tenant home is the
  CORE ground and the spine lives DOWN from there in `atlas/target/`, so the
  walk climbed past Desktop and found nothing. `atlas-mcp.exe` sits at
  `<repo>/line/` and the spine is built at `<repo>/target/` — two up and back
  down, which the walk now covers. Both `debug` and `release` profiles.
- **It still names what it could not run.** Nothing found returns `"atlas"`
  unchanged, so the refusal a caller sees is the same honest one as before.

Proved live, with no `--atlas-bin` passed at all: `verify_chain` on
`law/chain.jsonl` went from `exec: "atlas": not found in %PATH%` to
`verdict=FLIP entries=4`. 16 Go packages green after the change.

### The README stopped asking for what it does not need

Found by cloning this repo onto a clean tree and following it literally.

- **The MSVC toolchain is named.** `store/src/ffi.rs` links Windows' own
  `winsqlite3`, so `cargo build` needs `link.exe`. A fresh PC with only rustup
  dies on `linker 'link.exe' not found`, which says nothing about this
  project. Several GB of prerequisite, named nowhere until today.
- **Python 3.14 was never needed.** The core's `pyproject.toml` says `>=3.10`
  and CI proves 3.10 and 3.13. The `.venv` the same line named is created by
  nothing and is gitignored; `tests/prove.py` uses `sys.executable`.
- **`cd line ; go build ./...` produced no binaries and no error.** `line/cmd`
  holds five main packages and Go discards every result when it compiles more
  than one. It is a compile check that reads like a build.


### The rule that kept winning arguments it was not in

`.covenant` being declared twice was one instance of something; this is the
rest of it, found by measuring instead of by eye. Two audits were written
against the stylesheet: one for the same selector declared twice, one for the
shape that actually did the damage — two DIFFERENT classes worn on one element
where the later one silently takes a property.

**The culprit is `.muted`, and it has now ambushed three things.** It is worn
as a colour utility, 62 times across this console, but it also sets a
`font-size` and a `margin-top`, and at one class of specificity the later rule
wins. It took the hero verdict from 40px down to 12px this morning. It was
also cutting Version control's footnote — *"Every button here is your hand,
not the machine's"* — from 16px of top margin to 6px, which is why that line
has been sitting closer to the buttons than it was written to.

- **The spacing utilities moved to the foot of the file.** A utility exists to
  set one property and `.mt-16` was losing that property to a colour class.
  Utilities come last; that is the whole reason they are a category. Measured
  after: the footnote gets its 16px.
- **The hazard is written down at `.muted`'s own definition**, in the terms a
  hand needs: wearing it beside any class that sets `font-size` or
  `margin-top` means `.muted` wins unless that class is declared after it.
  Stripping the two extra properties would resize all 62 uses, which is a
  change nobody asked for, so the rule stands and the trap is documented where
  it is stepped in.
- **`.seat-open` declared `color: inherit` and never got it** — `.card-title`
  is later and sets the colour. What the link actually wanted was to look like
  every other card title, which is what it was already doing, so the dead
  declaration is gone and the hover stays (a pseudo-class outranks a bare
  class). Measured: the seat name renders at `--accent-2`, identical to a
  title that is not a link.
- **`.sidebar-footer`'s `font-size: 11px`** was re-stated verbatim as
  `var(--t-xs)` further down and never read.

**What the audits say now.** Three pairs still collide and all three are the
INTENDED rule winning: `.btn-sm` over `.btn`, `.home-box` over `.card`,
`.chat-foot` over `.muted`. Three more — `.empty-text`, `.textarea`, `th` — hold
an earlier value a later rule deliberately refines, which is a cascade doing
its job rather than a defect. No ambushes remain.

Swept afterwards across all eight pages: none scrolls sideways, and the crumb
is hidden on the Dashboard and present everywhere else.

### Two the sweep turned up on the way

- **"Through the council" was pushing its own input out of its card.** The
  header's right-hand group was `flex-shrink: 0`, which is right for a row of
  buttons and wrong the moment a FIELD is in it — the base field rule is
  `width: 100%`, and a full-width input in a group that will not give anything
  back goes straight through the card's edge. A field in a header now sizes to
  the room it is given. Measured at 1400px: one row, 48px, inside the card.
- **A changed file's label ran into its filename with no gap.** The label was
  an inline-block with `min-width: 150px`, which is a floor and not a ceiling,
  and *"changed, ready to save"* is wider than that in the mono face. It is
  two real columns now: the label takes what it needs up to a cap, the path
  takes the rest and wraps rather than pushing the card.

### The Aurora scheme

The operator pointed at a console he built for an earlier version of this
estate and said what he wanted from it: *"i like the current dashboard layout,
and the colors/flow of the one i sent."* So the layout is untouched — the
deck, the hero, the ledger, the crumb all stand — and the palette and the
flow devices are his.

**The palette now means something.** This console was Catppuccin-adjacent: a
neutral near-black under a cornflower blue. Aurora's ground is blue-GREEN at
the root (`#05080a`), its text is teal-tinted rather than grey, and it runs a
phosphor cyan with an amber second. The colours are not decoration — the
**amber is what the estate calls SEALED**, the **phosphor is what it calls
PROVEN**, and the orange marks a section of the record. Fifty-five hardcoded
`rgba()` literals of the old palette were still scattered through the file
behind the tokens; every one of them moved.

- **Mono-first, as Aurora is.** Its whole console is one fixed pitch and that
  is most of why it reads the way it does: every label, value, field and
  control on one grid. The prose face is kept for the handful of places this
  console carries real SENTENCES — a hero line, a page subtitle, the brief,
  a delivery — where a fixed pitch costs more in reading than it earns.
- **The glowing edge.** A 2px bar down the left of the panel the page is
  about, with a 12px bloom off it, in the verdict's own colour — so the edge
  says the same thing as the word beside it, or it says nothing.
- **The section rule.** Aurora's headings are a tracked micro-label followed
  by a hairline running to the panel's edge. That single device is what makes
  its panels read as parts of one instrument rather than a stack of boxes;
  every card header in this console now closes the same way.
- **The washed ground**, a cyan bleed at top centre and an amber at top right,
  fixed behind everything for one paint.
- **The pinned action wears the seal.** Booting opens a sitting and closing
  pays its toll: both move the record, so that one button is amber while cyan
  stays the colour of reading.

### Four that were wrong, found by looking at every page

- **The Watchboard's badge read "delivered — 54.3s" over a panel reading
  "nothing has run in this tab yet."** `Run.check()` is where a turn kept from
  a previous visit comes back, and the only handler it fires is `state`, which
  repaints the badge alone. The board now paints from what that read found.
- **A dead badge on Evals.** The run card moved to the Dashboard and took its
  log element with it; `paintRun` returns at its first line when that element
  is absent, so `#ev-run-state` was never painted once and sat in the header
  of every visit showing a hardcoded em dash — a control that looks like a
  reading and is a literal.
- **`.covenant` was declared twice** in one stylesheet, and the later rule won
  silently. Two rules for one class is exactly how `.muted` came to steal the
  hero verdict's font-size this morning. One rule now, and the hairline above
  it does the lifting rather than the colour.
- **Three tables were written bare**, and a bare table pushes its card, which
  pushes the grid, which scrolls the page sideways: Settings was 601px of
  content in a 595px main because of a badge reading `can_approve:false`.
  Wrapping the three by hand fixes the three and not the fourth, so a table
  in a card scrolls in its own box whether or not anyone remembered the
  wrapper. Button labels stopped breaking mid-word in the same pass — the
  mono face is wider, and a row of buttons wraps between buttons or it stops
  reading as a row of controls.

### The watchboard: the council from the inside

Chat was a second conversation — the same box as the launchpad, the same
`Chat.thread` array, the same bubbles, plus scrollback. Two of its jobs were
real and neither needed a whole page: answering a gate, and showing each
seat's words as they stream. Both belong in a watchboard. His words:
*"basically like a multi-panel watchboard to see the backend of all the system
and core work so literally every tool call and everything is being landed on a
page. this would be like the internal chat of the models themselves."*

**The feed was already there and nothing rendered it.** council.js has kept
EVERY event of every turn since it was written — *"the record of a run is the
events"* — and the only reader was a single-column card that drops, on
purpose, the one kind that matters most here:

    case 'token': return '';   // App.runRow

The seats' own words. The step-by-step shows what the council DID; this shows
what the models SAID while doing it. Same wire, nothing new asked of the
engine, nothing new stored, no second definition of anything.

- **Four panels**, because they answer four different questions and reading
  them interleaved is what made the single column unreadable. THE TURN: what
  was asked, under which pipeline, and the delivery. THE FLOOR: every seat that
  took it, its model, and its raw output whole. THE TOOLS: every call, the
  arguments in, the result out, `failed` read off the engine's own field. THE
  WIRE: every event in order, unreduced — with a `raw` toggle that prints
  each whole payload, because *every event* has to mean every event or the
  panel is only another summary. An event kind this build has never met is
  printed as its own JSON rather than dropped.
- **Tokens are counted, not listed, when the wire is folded** — a turn carries
  hundreds and they are shown whole on the floor. `raw` lists them too.
- **The past is on the page.** The live turn is one turn; `logs/` holds every
  run this ground has made, served whole with a sha256 by `records`. That
  answers the other half of what he asked: *"should we move the session
  tracking over there? more of the in-depth view."*
- **The gate stays** — a council question is answered in a field on this page,
  never a `prompt()`, never a default, never a guess (RULE 6). The label in the
  panel is Watchboard; the route and `data-page` stay `chat`, exactly as flows
  did when it became Version control.

**Proved against a real turn, not a stub.** Engine booted, sitting 196,
`git status` run from the watchboard's own box: 2 seats, 1 tool call, **148
events**, 54.3s, delivered. The floor carried Router (`qwen3.5:4b`, 280 chars)
and Steward (`llama3.2:latest`, 409 chars) streaming their own words under
their own names; the wire carried the arithmetic guard, the skill decision, the
call, the result, both token runs and the delivery. Sitting closed and tolled.

- **Two clocks on one page stopped disagreeing.** The facts row is painted per
  EVENT, and a seat can think for a minute without putting one on the wire — so
  it read 2.6s beside a badge reading 49s. The turn's clock now ticks into its
  own span, the same shape as the engine's age on the Dashboard.
- **`when()` takes an epoch, not an ISO string**, and hands anything else
  straight back — so the earlier-runs table printed raw `2026-09-11T02:59:20Z`
  in UTC. Records already knew this and wrapped it in `Date.parse`.

### The Dashboard answers the question it is actually asked

Rebuilt against what the console is FOR, at the operator's word: *"think about
what the dashboard is used for and what would make sense to put where
functionally"*, and against a read of what dashboards do now.

**The order was wrong, and it was wrong in the one way that costs something.**
He sits down and asks, in this order: *can I work at all — what do I want done
— what is happening — is the ground sound — is anything waiting on me.* The
page answered them almost backwards. The ENGINE CARD, which gates every other
thing here — nothing typed into that box runs without an engine — sat BELOW
the box it gates, three scrolls down. Meanwhile the five score cards, which
move perhaps twice a day, held the top-left quadrant.

- **The hero is the sitting now.** Open or not, on which world, since when,
  with the one button that changes it. The engine card is folded into it and
  gone as a card; its whole content was one sentence, which is the exact
  complaint its own comment made about the card it had replaced.
- **With one override: a red build outranks an unopened engine.** Booting onto
  a broken build without being told is worse than not knowing the engine is
  shut, so red anywhere raises an alarm in the hero that NAMES what fell and
  the first thing to fall. A measure that was never run raises a yellow one:
  *not-tried is not the same as passed*, and only one of those is good news.
- **The proof drops to a ledger beside it** — rows, not cards. Five cards of
  identical weight made the eye do the ranking; a ruled list puts every value
  in one column where they can be compared in a single sweep, which is the
  only reason to show them together. Three elements above the fold now
  (the sitting, the ledger, the box) where there were nine.
- **The brief rose above the box, and green went actually silent.** It is the
  interrupt channel — a gate waiting, a refusal, an engine running code older
  than the ground — and it arrived UNDER the box it should have changed what
  he typed into. When nothing needs him it now renders nothing at all; the
  page subtitle still carries the all-clear, which is where a quiet statement
  belongs.

**What the market read gave us, and what it did not.** The 2026 consensus
across the tools this one sits beside is decision-first rather than KPI-grid:
one north-star answer in the top-left, four to six supporting measures, and
nothing above the fold that does not change the next move — with Nielsen
Norman's finding as the hard edge, that a reader scans in a Z and gives up
past about seven competing elements. That is the shape taken. What was NOT
taken: sparklines, donuts, and an "AI summary" of numbers the record already
states plainly. Every line still names the file it was read from, which is
this console's own law and worth more than any chart.

### The seat page stops being a lie

- **`/agents/:id` was an orphan AND wrong.** Nothing on the seats page linked
  to it, so the only way in was to type the URL — and when you did it read
  `API.getAgent`, which queries the webapp's own SQLite `agents` table. Nothing
  writes that table. The list beside it reads `agents/*.md` through the `seats`
  tool, so a ground holding fourteen declared seats answered **404 for every
  one of them**, and the fields it was built to show (office, reports_to, mode,
  permissions) do not exist in a declaration at all. Fourth instance of a page
  counting the webapp's store instead of asking the record, and the last one
  standing. A route nobody can reach is a route nobody notices is broken.
- **It reads the record now**, keyed on the file stem — stable, unique,
  url-safe, and already what the record calls the document. It carries what
  the card cannot: the declaration with nothing folded, the system prompt open
  rather than behind a toggle, every pipeline placement with its step and
  condition, and **the file itself with its sha256**, served by `records`.
  `/agents/steward` went from HTTP 404 to seven fields, five placements and a
  3422-character prompt.
- **An absent name is denied honestly** and the denial lists the fourteen that
  are there, each one a link — the same doctrine `read_doctrine` already held
  itself to.
- **A seat's name is the way in.** It reads as a title until you go near it:
  on a page of fourteen cards, a row of blue links is a row of noise.

### Three that were silently wrong

- **`.hero-verdict muted` came out at 12px.** `.muted` further down the file
  sets a font-size and a margin, has the same specificity, and is later — so
  it won. A utility class used as a state name is a collision waiting for
  whichever rule happens to be written last. The state has its own name now.
- **`window.Home` is not a thing.** `Home` is declared `const` at the top of
  home.js, and a top-level `const` in a classic script binds in the global
  LEXICAL scope, never as a property of `window`. The guard was always false,
  so the hero came up with its kicker and an empty body — which looks exactly
  like a failed read.
- **A system prompt is prose, so it wraps.** The base `pre` rule now scrolls
  anything too wide, and a paragraph you have to drag sideways is not served
  whole in any sense that matters. Link-buttons stopped underlining in the
  same stroke: several buttons here are `<a>`, and they were carrying the
  browser's underline through the button's own chrome.

### The console says where you are, and what moved
- **A crumb, painted once by the router.** Every page but the Dashboard now
  opens with `ATLAS / <page>`, and a detail view carries its third step
  (`ATLAS / Agents / steward`) because there the depth is a fact. It lives
  above the routed content rather than inside each page's header: the router
  is the only thing that knows where you are, and fourteen hand-written copies
  of one line is fourteen chances for one to be wrong. The page's NAME is read
  off its own nav link, so the relabel to *Version control* reached the crumb
  without a second edit.
- **The standup card carries a delta.** `+1 vs last`, read from
  `run_history.jsonl`, which holds every prior live run. It is the ONLY card
  that gets one: nothing else on that row has a second measurement to compare
  against, and an unchanged reading and a never-measured-twice reading are
  different claims. A card with nothing to compare shows no delta rather than
  a zero.
- **Fields are styled by the element, not only by the class.** Ten inputs in
  this console were written without `class="input"` and came out as the
  browser's white box with black text — the save-message field on Version
  control was a bare white slab sitting on a dark card. Classing the ten by
  hand fixes the ten and not the eleventh, so the rule now binds to `input`,
  `select` and `textarea` themselves, with `:not()` guards keeping it off the
  controls that are not fields. The placeholder took the ground's own muted
  colour in the same stroke; it had been inheriting a grey that read nearly as
  loud as a typed value.
- **A git act on Version control re-proves the sidebar badge.** The badge
  repainted only on a run event, so a save or a send made from that very page
  left the panel showing the count from before the act — the one moment the
  number is most obviously being watched. Seen live: research pushed, the page
  said *in step*, the badge still read 7. It repaints that badge ALONE
  (`paintOwed`, the same shape as `paintProof(box, only)`): it is the only one
  an act on that page can move, and doing all four costs six tool calls in a
  row — long enough that the first cut still showed the stale number when the
  eye went looking for it.
- **Settings stopped scrolling sideways.** A `1fr` grid column is
  `minmax(auto, 1fr)`, and `auto` there means MIN-CONTENT, so one long
  unbreakable line inside a card — the team-bridge readout — made the column
  refuse to narrow and pushed the whole page. Measured: 737px of content in a
  595px main. `min-width: 0` on the grid children, and `pre` scrolls in its
  own box rather than shoving the page.

### The readers are proven too
- **15 more strokes on `internal/tools`**, on the three tools that READ the
  estate's own record — where a quiet wrong answer does the most damage. Every
  one of those three exists because an earlier page counted the webapp's own
  SQLite store instead of asking the record and showed four zeros on a full
  estate; nothing had ever held them to it. **34 strokes in the package now.**
- **`records`**: sorted by what a document *is*, with the law above the seats
  (alphabetical would invert that), `law/` marked sealed, only markdown served,
  and the sha256 receipt matching the bytes actually handed over. The path
  stroke is the one worth reading — **a name is matched against the listing and
  never joined onto Home**, so `../.env`, an absolute path and `law/../.env` are
  not *defended against*, they simply are not in the list and cannot resolve.
- **`proofs`**: a closing sitting line folds over its opening one (get that
  wrong and every sitting counts twice), a half-written last line — what a
  killed process leaves — is dropped without losing the lines before it, one
  absent file never blanks the other two, and the counts the core owns are
  *named* rather than recounted.
- **`seats`**: fields are whatever a declaration carries, including one invented
  tomorrow, because the core's own parser names no field either; the colon comes
  off a key; a `System Prompt` with an empty value takes the body beneath it and
  stops being counted twice.
- **The registry contract is pinned.** A tool that writes must declare it —
  that flag is what the review-only table refuses by, so a mislabelled tool
  would hand counsel a hand instead of eyes.
- **`runstream.go` stays uncovered, deliberately**, and is named as such in
  `tests/PROVING.md`. `RunStream`, `AnswerStream` and `ListenStream` all need a
  live engine on an open sitting; a mock there would prove the mock.


### The console takes a keyboard
- **Ctrl+K / Cmd+K opens a palette.** Every page and the engine's verbs, three
  letters and Enter. Checked on the event rather than the platform, because a
  Mac keyboard on a Windows box happens and neither should have to learn the
  other's key. **It navigates and it boots, and nothing else** — a palette that
  can run anything is a second command surface to learn and a second place for a
  destructive verb to hide. Close is in it because it *pays the toll*, the one
  routine act the record depends on and the easiest to forget.
- **What it offers is what is true right now.** The engine rows read `Run`'s live
  state, so a palette opened after a crash does not offer to close a sitting that
  already closed.
- **Tables read down a column.** A stripe barely there — enough to hold a line,
  not enough to become a pattern — drawn as a background so a hovered row still
  wins without a specificity fight. Figures are right-aligned, which is what
  makes two numbers comparable at a glance now that they are tabular.
- **Waiting looks like the thing that is coming.** "Reading both grounds..." with
  a spinner is a sentence about the machine; a skeleton says how much is arriving
  and where it will sit, so the layout does not jump when it lands. It is the one
  piece of non-user-triggered motion here and it earns it — a still grey block
  reads as a broken image — and it is switched off under
  `prefers-reduced-motion`, along with the mic pulse.
- **An empty screen is an invitation**, so it gets room and a plain sentence
  rather than a shrug in the corner of a card.


### The console is rearranged around what you actually look at
- **Recent moved with it**, and the nav item is **Version control** now. The
  route and `data-page` stay `flows`: the route is what every link, bookmark and
  history entry points at, and `data-page` is the key `css/icons.css` draws the
  glyph from. Clicking a recent objective still fills the box — the box is on
  the Dashboard, so it stages on `Home.pending` and the router goes there.
  Nothing is run by a click.
- **`escHtml` did not escape quotes, and seventeen attributes were built with
  it.** `textContent -> innerHTML` escapes `&`, `<` and `>` and nothing else, so
  a double quote passed through untouched — fine in a text position, wrong in an
  attribute. Found by moving Recent: the objective `git commit: "a test of the
  recent card"` rendered as `data-say="git commit: "`, the attribute ending at
  the operator's own quote with the rest of his sentence loose as markup.
  Escaped at the source rather than at seventeen call sites, because a rule that
  must be remembered at every use is forgotten at the eighteenth; `&quot;`
  renders as `"` in a text position, so the other 209 uses are unaffected.
- **The Dashboard opens with the scores.** strokes, smoke, standup, standups run
  and parity moved off Records, where they were a page you had to go to, and now
  head the launchpad. They answer *is the build sound* — the question worth
  answering before you type into the box beneath them.
- **The run sits under the box that started it.** The step-by-step was a card on
  Evals, a page away from the thing that produced it. Same renderer, same
  `ev-run` id — `App.paintRun` draws it unchanged, because a second copy of that
  renderer is exactly the drift this estate keeps writing docstrings about.
- **The engine card moved down.** A card whose whole content is one sentence was
  opening the page and crowding it. Still above the brief, because nothing below
  it runs until an engine is open.
- **Every git surface is on one page.** The repository card left the Dashboard
  for Flows — "this needs to go with the other github stuff" — so the door path
  (the per-world buttons, no engine needed) and the council path (law gate,
  recorded run) now sit one above the other. Two ways in, on purpose: *"there is
  a series of redundancies.. its called safety, bud."*
- **The flow builder came off.** Registry, editor, SVG graph, fire, resume, the
  runs table, compare, replay and the town board — all rendered, none ever used.
  **Nothing behind them was deleted:** the `/api/flows` routes, `flow.go`,
  `run.go` and their strokes all stand, so a builder is a render away.
- **Records leads with the documents.** What this ground *carries* is what the
  page is opened for; the sittings and standups are the history behind it.
- **Live standups show the last 5 of 33, newest first.** It listed every run ever
  recorded, oldest first, so the one that mattered — the last — was at the bottom
  of a table that grew a row every morning. The rest are named in
  `tests/run_history.jsonl`.


### Added
- **The Flows page controls git, in your own words.** The overwatch card was
  read-only; it now carries the verbs. Per world: a message box and **Save the
  work**, **Send to GitHub**, **Take from GitHub**, and **Lines of work** —
  which lists every branch as *you are here · the main line · on GitHub*, and
  opens, moves to, or finishes with one. A button that cannot work is greyed
  out **with the reason in its tooltip**, so it says why before it is pressed
  rather than after.

  Five new door tools behind it (`internal/tools/gitctl.go`, split from
  `gitstate.go` because that file's own header promises "REMOTE OPERATIONS ARE
  REPORTED, NEVER PERFORMED" and a file must not quietly stop meaning what it
  says at the top): `git_commit`, `git_push`, `git_pull`, `git_branch`,
  `git_remote`. Every one closes its child's stdin, jails paths to the
  tenant's Home, and reads the remote wall without ever opening it.

  **Nothing here fires on its own.** RULE 6 is untouched: these are the
  buttons on the operator's glass and the hand on the button is his.

  **Plain git, not `gh`.** The GitHub CLI is free, open source, installed and
  authenticated on this machine, and it still does not go in: RULE 4 walls
  anything needing someone else's server "even when the remote thing is
  better, free, or open source", and atlas law 6 is hand-roll or refuse.
  Everything named — branches, mains, open and closed, sending — is plain git.
  What `gh` alone would add is GitHub-side objects (pull requests, issues,
  releases, CI runs): a wall to open deliberately, not a dependency to acquire
  by accident.

- **`internal/tools` has strokes for the first time.** 19 of them
  (`gitctl_test.go`), on the package that carries every MCP tool handler and
  was named the estate's biggest hole in `tests/PROVING.md` that morning.
  Hermetic by law 5: each builds its own repository in `t.TempDir()`, none
  touches the record, none reaches a network — the two about sending prove the
  *refusal*, which is the only half provable without one.

### Fixed
- **The wall is the estate's, not one repository's.** The panel told the
  operator two different stories about one ruling — `research` "Sending
  allowed" and `atlas` "Sending OFF" side by side — because `dial()` read only
  `<Home>/.env`, and a carried tenant has no `.env` of its own. He had not
  shut a wall for atlas; atlas was looking in the wrong place. `dial()` now
  reads one level up, and that bound is not invented: it is the door's own
  ground law from `main.go` ("one level up, one level across"). It stops
  there, because walking to the filesystem root would leave the ground
  (RULE 1) and a stray `.env` in `Desktop\` must never open this estate's wall.

- **`readGit()` threw on every navigation away from the dashboard.** It
  captured `home-git-card`, `home-git` and `home-git-controls`, *then* awaited
  the door. A route away during that round trip left all three pointing at
  detached nodes — and writing `innerHTML` into a detached node SUCCEEDS,
  which is what hid it; the throw landed one line later on
  `document.getElementById('git-commit')` returning null. Both the 15-second
  poll and every turn-end call it, so it fired constantly and killed the rest
  of the handler each time. The elements are re-acquired after the await and
  the buttons are found through the bar, not the document.

- **The door was started without `--atlas-bin`, so nothing could reach the
  spine.** `verify_chain` answered `exec: "atlas": executable file not found
  in %PATH%`, which is what put three of the six workflows in the red.
  `atlas-mcp` defaults the flag to the bare string `"atlas"` and does none of
  the built-tree lookup `atlas-door` does. Relaunched with it;
  `verify_chain` now answers `verdict=INTACT`. The code default is still a
  trap and is written up as a ruling in `tests/PROVING.md`.

- **The door's `--manjuel` had been mangled into an error.** `mcp.err` held
  `refused: "C:/.../manjuel.py" is not a landed command` — the `python `
  prefix was lost by a PowerShell `-ArgumentList` earlier the same day, so the
  council engine was never wired. Restarted correctly; `mcp.err` now carries
  only its startup notes.

- **`tests/e2e/_start_mcp.ps1` pointed at the wrong tree.** All three of its
  lines named `C:\Users\novad\Desktop\Archive\atlas` — the pre-split copy,
  outside the ground RULE 1 fences, which nobody edits. A suite started
  against it would have proven a repository no one was changing. Every path
  now derives from the script's own location.

## [0.1.2] — 2026-09-10 (tag on 0586228)

*Post-script, 2026-10-09: marked after the fact on his word, on `0586228`, the commit that cut 0.1.2 (every VERSION file and Cargo.toml say 0.1.2 there), its message what it carries.*

### Changed
- **The stone moniker is struck from the version.** `0.1.1+f1` → `0.1.2`, in
  all fourteen places it was pinned: the six `VERSION` files, `Cargo.toml`,
  `core/src/version.rs`, both Go `main.go`s, both webapp handlers, and the two
  test harnesses. The stones are still named where they belong — THE_ROAD,
  STATE_OF_BUILD, the released CHANGELOG headers — but the version is now only
  what the thing IS, not which sitting cut it.

  **The drift catcher was inverted, not deleted.** `core/src/version.rs`'s
  tests used to REQUIRE a `+` and a stone; they now REFUSE one, so the
  moniker cannot creep back unnoticed. And `version.ps1` — the canonical
  bumper — was force-appending a stone in `Format-Version` and defaulting a
  stoneless version to `"f1"` in `Parse-Version`, so the very next
  `bump patch` would have quietly rewritten `0.1.2` as `0.1.2+f1` and
  re-broken every pin. The stone is gone from it entirely and it now refuses a
  version carrying one. Same trap as the flow cutter, one file over.

  **The prove battery caught what grep missed.** Five `VERSION` files under
  `line/` carry no extension, so the first sweep walked past all of them; the
  Rust spine's `version-cross` stroke named every one by path on the next run.
  `version.ps1` now says out loud that its six files are not the only pins and
  points at `tests/prove.py` for the rest.

### Added
- **The engine card says how long it has been standing, and when it is idle.**
  The card already showed `started 10:58:58 AM` — a clock time you have to
  subtract from to learn anything. It now ticks a live elapsed beside it, and
  turns amber with `— idle Nm` once nothing has run for five minutes.

  The measurement behind it, taken across every sitting the core has ever
  recorded: a standup gets 69 seconds of engine time per run; sittings of two
  runs or fewer get 208, and there are 63 of them — **5.3 engine-hours for 91
  runs**. Twenty-two sittings were never closed at all. Sitting 74 held an
  engine thirty minutes for 2 runs, 82 held one fifty-four minutes for 5, and
  166 held one **sixteen minutes for zero**. The record had known for weeks;
  nothing on the glass said a word.

  **Idle is measured from the last turn, not from boot.** A first cut went
  amber only when nothing had EVER run, which misses the shape the waste
  actually takes — 74 and 82 both did work, then sat. A running turn is never
  idle however slow the model is: this must not scold a slow rack, only an
  engine nobody is using.

  Ticks at one second, not on the 15s poll, which would read as a broken clock.
  The interval clears itself the moment the span leaves the DOM — an
  idle-engine warning that leaked timers would be its own joke.

  `static/js/home.js` (`Home.age`), `static/css/app.css` (`.eng-idle`).
  Operator: "we can just add that to the dashboard to view as tasks are
  running in real time, right?"

- **The card counts the turns, from the door.** `/run/state` now carries `runs`
  and `last_run`; the engine counts every turn it pumps, in `Engine.tick`.
  Before this the glass used `Run.turn` — only what THIS TAB had seen, empty
  after a reload — so an engine that had run ten turns read as untouched.

  **A slash command is housekeeping, not work.** Booting sends `/warm` and
  `/status` through the same path, so a freshly opened engine reported "2 runs"
  before anyone asked it anything, and **"0 runs" — the state most worth
  shouting about — was unreachable.** Counted in `Run` after `pump` returns
  rather than inside `pump`, because pump cannot see the objective and only the
  caller knows what it was. An `Answer` always counts: that is his hand.

  **Idleness can never predate the engine.** With no runs, the client fell back
  to its own `turn.ended` — which survives a reboot — and a
  THIRTY-EIGHT-SECOND-OLD engine reported `idle 14m`, counting from a turn a
  previous engine had run. The floor is now this engine's own start.

  `internal/engine/engine.go`, `internal/tools/runstream.go`,
  `static/js/council.js`, `static/js/home.js`.

- **A turn is visible in every browser, not only the one that started it.**
  Operator: "i want to be able to see it runing in sync on my chromium browser
  with your internal browser."

  `/council/stream` belongs to whoever opened it, so a second window could not
  know a turn was running and would not learn until its own 15s poll — by then
  the turn was usually over. **Storage could never have fixed this**:
  sessionStorage is per-tab, localStorage is per-browser. The server is the
  only place two browsers can agree, so `pipeSSE` now mirrors every council
  line onto the broadcast bus `/api/events` and `/ws` already serve, and
  `Run.mirror()` follows it. The runner ignores its own echo (`running` is the
  whole test) or every event would be doubled in its trace.

  **A watched turn needed somewhere to land.** The runner pushes its own pair
  into the thread in `ask()`; a mirroring window never did, so every event
  found no live bubble and the page sat on "waiting for the engine" while the
  answer streamed past it.

- **THE BUS HAS NEVER SENT THE TYPE IT IS KEYED ON.** `SSE` marshalled
  `e.Data` alone, dropping `e.Type` — so `App.onEvent`, which switches on
  `e.type` across eleven events, **had never fired once**. The struct already
  carried `json:"type"`; only that one line disagreed. Fixed, and the toasts
  work for the first time.

- **The boot report is kept.** It lived only in a `<pre>` that `render()`
  rebuilds empty and hidden, so a refresh discarded the RECORD, GATE, RACK and
  VOICE the boot had reported and the only way back was rebooting a healthy
  engine. Now written to the settings store, debounced, **keyed to the
  session** so a dead engine's report is never painted over a live one.

- **The conversation is kept, so a reloaded browser is not an empty box.**
  Mirroring live events was only half of "keep the boot report and the run
  turns": `Chat.thread` is in-memory per tab and `Run.turn` is sessionStorage,
  also per tab, so a browser that reloaded AFTER a turn showed the engine card
  and nothing else. That is exactly how he was checking — "i have been
  reloading my external browser tab every once in a while to watch and see if
  there is parity" — and there never could be. The thread now goes to the
  settings store, keyed to the session, and any browser restores it on load.

- **`onRun` HAS THROWN AT THE END OF EVERY TURN SINCE a014c77.** That pass took
  the sittings strip off the launchpad at his markup and deleted
  `Home.readSittings` — but left `this.readSittings()` standing in the turn-end
  handler. A `TypeError` on every completed turn, silent, killing the rest of
  the handler. It cost nothing while it was the last statement, and the moment
  anything was added after it that thing simply never ran: the kept
  conversation looked correct in every unit and stored nothing at all. Found
  by asking the live page which of the three calls threw, rather than reasoning
  about it.

- **A council bubble's words are in `turn.answer`, not `text`.** `line()`
  renders from `m.turn` and leaves `m.text` empty, so a first cut of the keeper
  filtered on `m.text` and discarded every answer, keeping only what he typed.

- **Only a real start opens a watched turn.** The mirror opened one on ANY
  event, so the trailing lines of a turn the tab had just run — arriving after
  `running` went false — opened a second, empty turn and pushed "(a turn
  started in another window)" into the runner's own conversation.

- **The broadcast bus was 32 deep and dropping.** `broadcast` drops rather than
  blocks, which is right — one slow watcher must never hold up a turn — but at
  32 a council turn's token events made dropping the normal case. 512 deep, and
  the writer now drains in batches and flushes once instead of flushing per
  event, which is what let it fall behind in the first place.

- **THE BALL — one command proves the whole of atlas.** `tests/prove.py`
  gathers the eight places atlas proves itself: the Rust spine, both Go
  modules, the two batteries shipped inside the binaries, **twenty-seven**
  golden verifiers, the six workflows and the E2E suite. AGENTS.md named
  four verifiers; nobody had run the other twenty-two, which is how two legs
  stayed red for weeks without anyone seeing it.

  **It answers in three verdicts, not two.** `ABSENT` means a leg named a
  dependency this ground does not hold — a read-only source ground never
  copied in, a binary not built, a door not answering. ABSENT is never
  counted as a pass and never silently skipped: it prints the exact path or
  command that would answer it, and it does not exit red, because nothing is
  broken — something is missing, and the difference is the whole point. Only
  FAIL exits red. Every child runs with `stdin=DEVNULL`, so the ball is safe
  to run from inside a live engine turn.

  Standing today: **22 held · 15 absent · 0 broke**. The fifteen absences are
  fourteen cutters plus one Rust stroke whose read-only oracles
  (`estate\`, `secondbrain\`) are not in this ground. Their goldens are all
  here and green; what is gone is the ability to re-cut them and to notice
  the source drifting. That is a ruling for the operator — restore the
  grounds read-only, or retire those cutters with the goldens frozen as the
  authority. Substituting anything would break law 2 outright.

  `tests/PROVING.md` is the map: every leg, what it costs, what it needs,
  and the seventeen packages that carry no prover at all — `internal/tools`
  first among them, where every MCP tool handler lives.

### Fixed
- **`internal/rack` had been red since the repo split.** Four strokes wanted
  `tests/fixtures/rack_open_ground/state/rack_ledger.jsonl`. Git does not
  track empty directories, so the fixture ground came over empty in the H0
  pull and nothing said a word. The right fix was not a hand-written file
  but running the cutter that owns it — `tools/cut_rack_open_vectors.py`
  builds that ground by definition — after which the five tracked goldens
  judged it byte-exact.

- **The door battery had no binary to shell.** `TestDoorProveStrokesGreen`
  drives `cmd/atlas-door` against the Rust spine over a real loopback
  socket, and the Rust binary was never built in this ground. `cargo build
  -p atlas`. The prover had been behaving correctly the whole time —
  refusing by name rather than fabricating, exactly as law says.

- **`cut_flow_vectors.py` was behind its own fixture, and would have eaten
  it.** `ce9c390` added the `run` node kind to the flow contract and updated
  `flow.go` *and* `flow_vectors.json` — but not the cutter that owns the
  fixture. So `--verify` had been red since; worse, running the cutter
  *without* `--verify` would have rewritten the fixture, stripped `run` back
  out, and turned `internal/flow` red with no visible cause. The cutter now
  carries the seven-kind contract and refuses a `run` node with no
  objective, exactly as `flow.go` does. A full re-cut is now a no-op. (The
  seventh refusal's keys also sort like every other entry now — the giveaway
  that it had been hand-added rather than cut.)

### Note
- Static assets are compiled into the binary (`//go:embed static/...`), so a
  change under `static/` needs `go build -o atlas-webapp.exe .` and a restart
  before it reaches the glass. Editing the file alone does nothing.

## [0.1.1+f1] — 2026-09-08

### Added
- **E2E test suite**: 84 scenarios across 12 categories
  - S1: Binary Smoke (5), S2: MCP Tool Surface (8), S3: Write Tools (5)
  - S4: Mesh B2 (8), S5: Rack F1 (10), S6: Guard (5), S7: Tenant (4)
  - S8: Ollama Integration (10), S9: Webapp (6), S10: Cross-Impl (4)
  - S11: Agent Lifecycle (6), S12: Prove Chains (3)
- **6 CI/CD pipelines** (parallel jobs in ci.yml)
  - SPINE: Rust core (12 stages), LINE: Go MCP (11 stages)
  - GOLDEN: Python verifiers (11 stages), AGENT: 40-seat validation (11 stages)
  - WEBAPP: GUI + API (18 stages), OLLAMA: Live integration (14 stages)
- **6 agent workflows** with Python scripts
  - SCOUT: read-only survey (5 steps)
  - STEWARD: memory + testimony via Ollama (5 steps)
  - MESH: encrypted communication (7 steps)
  - GATE: injection blocking, PII stripping (5 steps)
  - TOWN: task scheduling (4 steps)
  - OPERATOR: full lifecycle (12 steps)
- **Ollama prover**: stdlib-only Python harness
  - Tests 11 models (qwen3.5, llama3.2, phi4-mini, gemma4, deepseek-r1, etc.)
  - Direct Ollama API + MCP integration + webapp API
  - JSON output, exit codes, timeout handling
- **Documentation**: PIPELINES.md, WORKFLOWS.md, OLLAMA_PROVER.md, ACCEPTANCE.md, E2E_SCENARIOS.md

### Changed
- `ci.yml` expanded from 1 pipeline to 6 parallel pipelines
- `generate` calls use chat API for thinking-model compatibility (qwen3.5)

## [0.1.0+f1] — 2026-09-08

### Added
- **Rust core**: prov-hash, ground-prove, atlas-api, store crates
- **Go MCP server**: 25 tools + HTTP transport + embedded GUI (port 8090)
- **Go TUI**: Dashboard, traces, agents, chain view, mesh, rack, format output
- **Go webapp**: Observability GUI with real-time SSE (port 8091)
  - Dashboard with stats, recent traces, agent roster
  - Agent registry with full .us declaration viewer
  - Trace log with SHA-256 hash chain
  - Tool surface with JSON argument invocation
  - Evaluations engine with pass/fail scoring
  - Message bus with Discord/Slack/WhatsApp adapters
  - Settings for MCP connection, eval thresholds, provenance
  - Search across agents and traces
- **Agent system**: 40 .us declarations validated, enrolled, documented
- **Skills**: prove, orient, enroll
- **Python verifiers**: cut_canon_vectors, cut_chain_verdicts, cut_us_vectors, fold_agents
- **TUI extractField**: Auto-find array field when `*` hits a map
- **YAML format**: Fixed `[]map[string]any` to `[]any` for type switch compatibility

### Governance
- `can_approve: false` structural enforcement across all 40 agents
- Covenant hash `1512741580b7239b` verified in all declarations
- Append-only witnesses: SEAT_LOG.md + STATE_OF_BUILD.md
- Forbidden verbs absent by construction
- Zero external dependencies enforced

### Documentation
- 54+ specification and design documents
- 40 agent documentation files with hierarchy index
- 3 skill documentation files
- GUI gap analysis (Atlas vs Langfuse/Langsmith)
- Build plan with phased acceptance criteria
- MIT License
- Contributing guide
- Full changelog

### Tests
- 92 Go tests (line workspace)
- 85 Rust tests (cargo test --workspace)
- 4 Python verifiers
- MCP 58-stroke prove
- TUI end-to-end verification against live MCP
- Agent enrollment validation (40/40)

### Fixed
- extractField: `agents.*.name` transform path now works correctly
- extractField: Auto-discovers first array-valued field when `*` encounters a map
- YAML format: `[]map[string]any` changed to `[]any` in rbac/agents list handlers
- TUI help: Changed `*.name` example to `name`
