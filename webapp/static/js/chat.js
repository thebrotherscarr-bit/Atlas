// ATLAS Chat -- the conversation with the council, and the Watchboard's panels.
//
// The operator, 2026-09-10: "review that chat panel, do we really need that?
// ... more of the in-depth view", then "have that string all of the CoT,
// thinking, tool calls, etc. basically like a multi-panel watchboard to see
// the backend of all the system and core work so literally every tool call
// and everything is being landed on a page. this would be like the internal
// chat of the models themselves."
//
// THE FEED WAS ALREADY THERE. council.js has kept EVERY event of every turn
// since it was written -- "the record of a run is the events" -- so the
// Watchboard asked nothing new of the engine and stored nothing new: the
// seats' own words, every tool call and every event are read off Run.turn.
//
// THE WATCHBOARD IS THE RUN TAB'S NOW (2026-10-10, WHAT'S LEFT D17's last
// part; his word: "Watchboard can go into runs", and on his card, "Yes, floor
// and tools too"). It was a page of its own from 2026-09-10; its address opens
// the front page with the Inspector's Run tab open, and the tab draws the turn
// with the three panels below, each handed the box it draws into:
//
//   THE FLOOR  every seat that took it, its model, and its raw output whole
//   THE TOOLS  every call: who made it, the arguments in, the result out
//   THE WIRE   every event in order, unreduced, including kinds this build
//              has never heard of -- an unknown event is still something that
//              happened, and an omission here looks like nothing happening
//
// What the page carried besides was already the front page's: the turn's
// facts and failures are on the Run tab's own card and on the terminal's entry
// for the turn, the box he typed in is the terminal, the council's question is
// a card in the dock at its foot, and the earlier runs are the Audit Ledger's
// Transcripts. Nothing here listens to a turn or draws a page: the front page
// does both, and asks these panels to draw again as the turn moves.
//
// The object is still called Chat because it still owns the conversation with
// the council: the front page reads and writes `Chat.thread`, and since the
// Dashboard retired (2026-10-07) Chat also keeps it (keep, showKept).
const Chat = {
  thread: [],          // [{who:'him'|'council', text, turn?}] -- the front page's turns
  shut: {},            // seats folded shut by hand, by index; open by default, and again with every turn
  raw: false,          // the wire: every frame, or tokens folded
  TAIL: 6,             // how many exchanges are kept for another browser

  // ---- THE CONVERSATION IS KEPT WHERE BOTH BROWSERS CAN READ IT ----------
  //
  // Moved here from the Dashboard (home.js) when it retired, 2026-10-07: the
  // thread is Chat's, so keeping it is too. `Chat.thread` is in-memory per tab,
  // so a browser that reloaded after a turn showed an empty conversation; what
  // renders -- who spoke and what was said -- goes to the settings store, which
  // both browsers read. The events, seats and tokens behind a turn stay in the
  // transcript on disk, the real record. KEYED TO THE SESSION, like the boot
  // report, so one engine's conversation is never painted under another's.
  keep() {
    clearTimeout(this._keep);
    this._keep = setTimeout(() => {
      // A COUNCIL BUBBLE'S WORDS ARE IN `turn.answer`, NOT IN `text`.
      const said = (this.thread || [])
        .slice(-this.TAIL)
        .map(m => ({
          who: m.who,
          text: m.who === 'him' ? (m.text || '')
            : ((m.turn && (m.turn.answer || m.turn.refusal)) || m.text || '')
        }))
        .filter(m => m.text.trim());
      if (!said.length) return;
      API.setSetting('thread.' + (Run.world || 'research'),
        JSON.stringify({ session: Run.session || '', said })).catch(() => {});
    }, 600);
  },

  // Restored only while the engine that wrote it stands, and only into a tab
  // that holds no conversation of its own.
  async showKept() {
    if ((this.thread || []).length) return;
    try {
      const r = await API.getSetting('thread.' + (Run.world || 'research'));
      const kept = JSON.parse((r && r.value) || '{}');
      if (!kept.said || !kept.said.length) return;
      if (!Run.engineOpen || (kept.session && kept.session !== Run.session)) return;
      this.thread = kept.said.map(m => ({ who: m.who, text: m.text }));
    } catch { /* nothing kept yet */ }
  },

  // ---- the Watchboard's panels, drawn into the Run tab's boxes -----------

  // THE FLOOR: the thing the Watchboard was made for. Every seat that took the
  // turn, its model, and its OWN WORDS whole -- not the recompose, not a
  // summary. A seat's name folds its words shut, and opens them again.
  paintFloor(box, n) {
    if (!box) return;
    const t = Run.turn;
    const seats = (t && t.seats) || [];
    if (n) n.textContent = seats.length ? seats.length + ' took the floor' : '';
    if (!seats.length) {
      box.innerHTML = `<div class="empty-text">No seat has spoken yet.</div>`;
      return;
    }
    box.innerHTML = seats.map((s, i) => {
      const shut = this.shut[i] === true;
      const words = (s.text || '').trim();
      return `<div class="wb-seat${shut ? ' shut' : ''}">
        <button class="wb-seat-head" type="button" data-seat="${i}">
          <span class="wb-seat-name">${escHtml(s.seat || 'seat')}</span>
          <span class="wb-seat-model">${escHtml(s.model || '')}</span>
          <span class="wb-seat-n">${words.length ? words.length + ' chars' : 'silent so far'}</span>
        </button>
        ${shut ? '' : `<div class="wb-words">${words
          ? escHtml(words)
          : '<span class="muted">nothing on the wire from this seat yet</span>'}</div>`}
      </div>`;
    }).join('');
    box.querySelectorAll('[data-seat]').forEach(b => {
      b.onclick = () => {
        const i = +b.dataset.seat;
        this.shut[i] = !this.shut[i];
        this.paintFloor(box, n);
      };
    });
    if (Run.running) box.scrollTop = box.scrollHeight;
  },

  // THE TOOLS: who called what, with which arguments, and what came back.
  // `failed` is the ENGINE'S own field — the pipeline's test, not a reading of
  // the words that came back.
  paintTools(box, n) {
    if (!box) return;
    const t = Run.turn;
    const calls = ((t && t.events) || []).filter(e => e._kind === 'tool' || e._kind === 'tool_result');
    const tools = (t && t.tools) || [];
    if (n) n.textContent = tools.length ? tools.length + ' called' : '';
    if (!calls.length) {
      box.innerHTML = `<div class="empty-text">No skill has been called in this turn.</div>`;
      return;
    }
    box.innerHTML = calls.map(d => {
      const name = d.action || d.tool || d.name || '?';
      if (d._kind === 'tool') {
        return `<div class="wb-call">
          <div class="wb-call-head"><b>${escHtml(name)}</b>${d.seat
            ? ` <span class="muted">called by ${escHtml(d.seat)}</span>` : ''}</div>
          ${d.args ? `<pre class="wb-args">${escHtml(JSON.stringify(d.args, null, 1))}</pre>` : ''}
        </div>`;
      }
      return `<div class="wb-call ${d.failed ? 'bad' : 'ok'}">
        <div class="wb-call-head"><b>${d.failed ? 'FAILED' : 'ok'}</b> ${escHtml(name)}</div>
        ${d.error ? `<div class="wb-err">${escHtml(d.error)}</div>` : ''}
        ${d.text ? `<pre class="wb-args">${escHtml(String(d.text).slice(0, 4000))}</pre>` : ''}
      </div>`;
    }).join('');
    if (Run.running) box.scrollTop = box.scrollHeight;
  },

  // THE WIRE: every event, in order, unreduced. `raw` prints each whole
  // payload; folded prints one line each. Tokens are COUNTED rather than
  // listed when folded — a turn carries hundreds and they are shown whole on
  // the floor — but `raw` lists them too, because "every event" has to mean
  // every event or this panel is only another summary.
  //
  // THE RUN TAB'S (2026-10-07, and the one place it is drawn since 2026-10-10).
  // The front page's Run tab hands over the box it draws into, and its raw
  // switch is this one (this.raw), so a switch and its wire cannot drift.
  paintWire(box, n) {
    if (!box) return;
    const t = Run.turn;
    const evs = (t && t.events) || [];
    if (n) n.textContent = evs.length ? evs.length + ' events' : '';
    if (!evs.length) {
      box.innerHTML = `<div class="empty-text">The wire is quiet.</div>`;
      return;
    }
    if (this.raw) {
      box.innerHTML = evs.map(d =>
        `<div class="wb-ev"><span class="wb-kind">${escHtml(String(d._kind))}</span>` +
        `<code>${escHtml(JSON.stringify(d))}</code></div>`).join('');
    } else {
      let frames = 0, chars = 0;
      const out = [];
      const flush = () => {
        if (!frames) return;
        out.push(`<div class="wb-ev wb-dim"><span class="wb-kind">token</span>` +
          `<code>${frames} frames · ${chars} chars · shown whole on the floor</code></div>`);
        frames = 0; chars = 0;
      };
      for (const d of evs) {
        if (d._kind === 'token') { frames++; chars += (d.text || '').length; continue; }
        flush();
        out.push(`<div class="wb-ev"><span class="wb-kind">${escHtml(String(d._kind))}</span>` +
          `<code>${escHtml(this.oneLine(d))}</code></div>`);
      }
      flush();
      box.innerHTML = out.join('');
    }
    if (Run.running) box.scrollTop = box.scrollHeight;
  },

  // One line for an event, without inventing a shape for a kind this build has
  // never met: anything unrecognised is printed as its own JSON.
  oneLine(d) {
    const k = d._kind;
    if (k === 'seat') return (d.seat || '') + ' · ' + (d.model || '');
    if (k === 'run') return 'pipeline ' + (d.pipeline || '') + (d.transcript ? ' · ' + d.transcript : '');
    if (k === 'tool') return (d.action || d.tool || d.name || '?') +
      (d.args ? ' ' + JSON.stringify(d.args).slice(0, 160) : '');
    if (k === 'tool_result') return (d.failed ? 'FAILED ' : 'ok ') +
      (d.action || d.tool || d.name || '?') + (d.error ? ' · ' + d.error : '');
    if (k === 'report' || k === 'note') return (d.text || '').trim().slice(0, 200);
    if (k === 'needs_answer') return (d.prompt || '').slice(0, 200);
    if (k === 'delivery') return (d.pipeline || '') + ' · ' + String(d.text || '').length + ' chars';
    const copy = Object.assign({}, d);
    delete copy._kind;
    return JSON.stringify(copy).slice(0, 300);
  },

  // `logs/2026-09-10_200628_git_commit_a_skill_is_offered....md` is a stamp
  // and an objective welded together. The objective is the part a person
  // recognises, so it is what the row shows (the Audit Ledger's Transcripts).
  runName(name) {
    const base = String(name).replace(/^logs\//, '').replace(/\.md$/, '');
    const m = base.match(/^(\d{4}-\d\d-\d\d)_(\d{6})_(.*)$/);
    return m ? m[3].replace(/_/g, ' ') : base;
  }
};
