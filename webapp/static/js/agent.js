// ATLAS Agent -- the front page.
//
// THE OPERATOR, 2026-10-02: "this is what I am wanting, not some
// uber-complicated dashboard thing that I need a degree to figure out", and,
// shown a page out of his own AI Studio project, "that's what I am looking
// for". This is that page in this console's own plain JavaScript: a top bar
// with five doors, ONE terminal, and an Inspector that opens beside it when a
// door is pressed. Nothing else is on the page. The old Dashboard is still
// here, at /dashboard.
//
// WHAT IT TOOK FROM THAT PAGE, AND WHAT IT LEFT. The layout, the four mode
// tabs, the prompt line with its history, the five doors and the Inspector are
// his. The engine behind them is not what they said it was: the figure beside
// the wordmark was random, "WAL Active" was a label, the agent matched words
// and reported tool calls that never ran, the Aider tab appended a stub,
// committed it and said done, and the repair loop reported a repair that never
// happened. Here every figure is read from a tool or not shown (SPEC 3
// invariant 10), every failure is on the face of the answer (ESTATE LAW 5),
// and no tab claims what it cannot do: a tab that is not wired is marked so and
// says so when it is pressed (NOT_WIRED). None is left: all four are wired.
//
// BASH AND PYTHON ARE HIS TYPED SHELL (WHAT'S LEFT H15, 2026-10-03; his word of
// 2026-10-02: "typed by you, gated"). A line typed in either tab goes to the
// door's shell_run, which only his glass may call. A plain look runs at once;
// anything that writes shows a card in the thread and waits for his click; a
// secret or a path outside the ground is refused by name, and the card is not
// offered for a refusal. The page classifies nothing: it draws the one answer
// the door gives (state, class, why, hold, exit, output) and a failure is on the
// face of it. No shell tab needs an engine; the council's tab does.
//
// AIDER IS THE DOOR'S AIDER (WHAT'S LEFT H12 and B20, 2026-10-04; his word:
// "finish wiring in the Aider panel, it's on this system"). A line typed in
// the Aider tab is an instruction to the door's aider_run, which only his glass
// may call: Aider runs headless on a COPY of the files he names, behind a wall,
// against the coding seat's own model on this machine, and its edit is written
// back only while the world stands on a line of work and every Python file
// still parses. The page keeps the NAMES of the files in the chat (/add, /read,
// /drop) and judges nothing about them; the door refuses by name and the page
// draws the refusal. Nothing here commits, runs the suites or lands (RULE 6).
//
// THE COUNCIL IS THE AGENT. A line typed here is the REPL's own turn: the
// sealed law gate stamps it before any model reads a word, the one Router runs
// the tools, and the recompose puts every failure into the delivery. This page
// owns no council logic. It reads Run (council.js), the same object the
// Watchboard and the old Dashboard read, and it writes its turns into
// Chat.thread, so the three can never tell different stories.
//
// THE GATE IS HIS (RULE 6). A question from the council is answered in a field
// in the thread, never a pop-up, never a default. The Guardrails tab lists the
// writing calls parked at the door and approves or denies one only on his
// click; the page never answers for him.

const AG_ICON = {
  panel: '<svg class="ag-i" viewBox="0 0 24 24"><rect x="3" y="4" width="18" height="16" rx="2"/><path d="M15 4v16"/></svg>',
  out:   '<svg class="ag-i" viewBox="0 0 24 24"><path d="M14 8V6a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2-2v-2"/><path d="M9 12h12"/><path d="m18 9 3 3-3 3"/></svg>',
  trash: '<svg class="ag-i" viewBox="0 0 24 24"><path d="M4 7h16"/><path d="M10 11v6M14 11v6"/><path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2l1-12"/><path d="M9 7V4h6v3"/></svg>',
  copy:  '<svg class="ag-i" viewBox="0 0 24 24"><rect x="9" y="9" width="11" height="11" rx="2"/><path d="M5 15V6a2 2 0 0 1 2-2h9"/></svg>',
  help:  '<svg class="ag-i" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9"/><path d="M9.5 9.5a2.5 2.5 0 1 1 3.5 2.3c-.7.4-1 .9-1 1.7"/><path d="M12 17h.01"/></svg>',
  cpu:   '<svg class="ag-i" viewBox="0 0 24 24"><rect x="6" y="6" width="12" height="12" rx="2"/><rect x="9" y="9" width="6" height="6"/><path d="M9 3v3M15 3v3M9 18v3M15 18v3M3 9h3M3 15h3M18 9h3M18 15h3"/></svg>',
  mic:   '<svg class="ag-i" viewBox="0 0 24 24"><rect x="9" y="3" width="6" height="11" rx="3"/><path d="M5 11a7 7 0 0 0 14 0M12 18v3"/></svg>',
  x:     '<svg class="ag-i" viewBox="0 0 24 24"><path d="M6 6l12 12M18 6 6 18"/></svg>'
};

const $ag = (id) => document.getElementById(id);

const Agent = {
  // ---- what the page is made of ------------------------------------------
  MODES: [
    { id: 'agent',  label: 'Agent',  badge: 'default', wired: true,  tip: 'Ask the council' },
    { id: 'bash',   label: 'Bash',   badge: '$',       wired: true,  tip: 'Git Bash in this world. A plain look runs at once; anything that writes asks first' },
    { id: 'python', label: 'Python', badge: '>>>',     wired: true,  tip: 'A Python session that keeps its names. A plain calculation runs at once; the rest asks first' },
    { id: 'aider',  label: 'Aider',  badge: 'pair',    wired: true,  tip: 'Aider edits the files you name, on a line of work, with the coding seat\'s model on this machine. /add names a file, /undo takes a run back' }
  ],
  // What each tab offers to try, one click away.
  HINTS: {
    agent: ['/status', '/help', '/memory', '/sittings'],
    bash: ['ls', 'pwd', 'git status', 'git log --oneline -5'],
    python: ['1 + 1', 'x = 41', 'x + 1', '/reset'],
    aider: ['/files', '/status', '/undo', '/help']
  },
  PLACEHOLDER: {
    agent: 'Ask the council, run a /command (e.g. /status), or say what you want done...',
    bash: 'A command, run as you. A plain look runs at once; anything that writes asks first...',
    python: 'Python, in a session that keeps its names. Plain calculations run at once; the rest asks first...',
    aider: 'Say what Aider should change. Name its files first: /add path/to/file.py (/help lists the rest). It writes only on a line of work...'
  },
  // The five doors in the bar, each of which opens the Inspector on its tab.
  DOORS: [
    { id: 'aider',  label: 'Aider Pair' },
    { id: 'ledger', label: 'Audit Ledger' },
    { id: 'flows',  label: 'Workflows' },
    { id: 'docs',   label: 'Registry & Docs' },
    { id: 'guards', label: 'Guardrails' }
  ],
  // The Inspector's tabs: the five doors, and the two the bar does not carry.
  TABS: [
    { id: 'run',    label: 'Run' },
    { id: 'aider',  label: 'Aider Pair' },
    { id: 'ledger', label: 'Audit Ledger' },
    { id: 'flows',  label: 'Workflows' },
    { id: 'guards', label: 'Guardrails' },
    { id: 'docs',   label: 'Docs & Registry' },
    { id: 'rack',   label: 'Rack' }
  ],
  // WHAT A TAB THAT IS NOT WIRED SAYS WHEN IT IS PRESSED. It says so. A tab
  // that looked wired and answered with a stub is the fault the page it came
  // from was full of.
  NOT_WIRED: {
    // Every tab is wired (Bash and Python on 2026-10-03, Aider on 2026-10-04). The next tab that is not is marked
    // `wired: false` in MODES and says so here, in words, under its own id.
  },
  HIST_KEY: 'atlas.agent.hist',
  // The window a game is drawn for. A game is made for a window of its own size (an 800 by 600 canvas is common), and the page it is served as cannot be
  // scrolled to, so the play frame is given this much and shrunk to fit the room the stage has (fitPlay).
  PLAY_W: 1000,
  PLAY_H: 720,

  // ---- state --------------------------------------------------------------
  mode: 'agent',
  entries: [],            // the terminal, in order; a turn's words live in Chat.thread
  hist: [],               // what was typed, newest first
  hi: -1,
  open: false,            // the Inspector
  tab: 'run',
  rack: null,             // read off rack_list, never remembered
  holdState: null,        // read off hold_list
  bound: false,
  booting: false,
  aiderFiles: [],         // the names of the files Aider may change; the door judges them when a run is asked
  aiderRead: [],          // the names of the files it may only read
  aiderSt: null,          // read off aider_status, never remembered
  aiderLast: '',          // the last run that wrote, for /undo
  play: null,             // the page the maker made, played beside the terminal: { name, version, closed }; never remembered
  _seq: 0,
  _pane: 0,

  // ---- the page -----------------------------------------------------------

  async render(el) {
    this.el = el;
    document.body.classList.add('is-agent');
    if (!this.bound) { Run.on((w) => this.onRun(w)); this.bound = true; this.loadHist(); }
    el.innerHTML = this.shell();
    this.wire();
    this.paintMode();
    this.paintTop();
    this.paintTitle();
    this.paintOut();
    this.paintPanel();
    const input = $ag('ag-input');
    // An objective staged from Version control's Recent lands in the line,
    // unrun, for him to read and press.
    if (Home.pending) { input.value = Home.pending; Home.pending = ''; }
    input.focus();
    this.who();
    this.version();
    // STATE FIRST, then what needs it: the kept conversation is restored only
    // while an engine stands, so the check must come back before it is asked.
    await Run.check();
    try { await Home.showKeptThread(); } catch { /* nothing kept yet */ }
    this.adopt();
    this.paintTitle();
    this.paintTop();
    this.paintOut();
    this.readStatus();
    this.readAider();
    this.watch(true);
  },

  shell() {
    const doors = this.DOORS.map(d =>
      `<button type="button" data-door="${d.id}">${escHtml(d.label)}</button>`).join('');
    const tabs = this.TABS.map(t =>
      `<button type="button" data-tab="${t.id}">${escHtml(t.label)}</button>`).join('');
    return `<div class="ag" id="ag">
      <header class="ag-top"><div class="ag-top-in">
        <div class="ag-brand"><span class="ag-dot" id="ag-dot"></span><span class="ag-word">ATLAS</span><span class="ag-ver" id="ag-ver"></span></div>
        <nav class="ag-nav" id="ag-nav">${doors}</nav>
        <div class="ag-right">
          <div class="ag-stat" id="ag-stat"><span class="ag-sdot" id="ag-sdot"></span><span id="ag-stat-a">reading</span><span class="ag-slash ag-rack">/</span><span id="ag-stat-b" class="ag-rack"></span></div>
          <div class="ag-pages-wrap">
            <button type="button" class="ag-chipbtn" id="ag-pages" title="Every other page">Pages</button>
            <div class="ag-menu" id="ag-menu" hidden></div>
          </div>
          <button type="button" class="ag-chipbtn" id="ag-insp" title="Toggle the inspector">${AG_ICON.panel}<span>Inspector</span><span class="ag-adot" id="ag-adot" hidden></span></button>
          <div class="ag-user"><span id="ag-who"></span><button type="button" class="ag-icon-btn" id="ag-lock" title="Lock the console" aria-label="Lock the console">${AG_ICON.out}</button></div>
        </div>
      </div></header>
      <main class="ag-main">
        <section class="ag-win" id="ag-win">
          <div class="ag-bar">
            <div class="ag-bar-l"><div class="ag-lights"><i></i><i></i><i></i></div><div class="ag-title" id="ag-title"></div></div>
            <div class="ag-bar-r"><div class="ag-modes" id="ag-modes"></div>
              <button type="button" class="ag-tool" id="ag-clear" title="Clear the screen (Ctrl+L)" aria-label="Clear the screen">${AG_ICON.trash}</button></div>
          </div>
          <div class="ag-out" id="ag-out" aria-live="polite"></div>
          <div class="ag-play" id="ag-play" hidden></div>
          <div class="ag-hints" id="ag-hints" hidden></div>
          <div class="ag-in">
            <span class="ag-p agent" id="ag-prompt"></span>
            <input id="ag-input" type="text" autocomplete="off" spellcheck="false"
                   placeholder="Ask the council, run a /command (e.g. /status), or say what you want done..." />
            <div class="ag-in-r">
              <button type="button" class="ag-tool" id="ag-mic" title="Speak (local whisper, nothing leaves this machine)" aria-label="Speak">${AG_ICON.mic}</button>
              <button type="button" class="ag-tool" id="ag-help" title="Show suggestions" aria-label="Show suggestions">${AG_ICON.help}</button>
              <button type="button" class="ag-run" id="ag-go">Run</button>
            </div>
          </div>
        </section>
        <aside class="ag-panel" id="ag-panel" hidden>
          <div class="ag-ph"><span>inspector <i style="color:var(--ag-n700);font-style:normal">/</i> <b id="ag-ph-tab"></b></span>
            <button type="button" class="ag-tool" id="ag-pclose" title="Close the inspector" aria-label="Close the inspector">${AG_ICON.x}</button></div>
          <div class="ag-tabs" id="ag-tabs">${tabs}</div>
          <div class="ag-pane" id="ag-pane"></div>
        </aside>
      </main>
      <footer class="ag-foot"><div class="ag-foot-in"><span>ATLAS - sovereign agent harness</span><span>Offline-first &middot; Local runtime</span></div></footer>
    </div>`;
  },

  wire() {
    const on = (id, ev, fn) => { const e = $ag(id); if (e) e['on' + ev] = fn; };
    document.querySelectorAll('#ag-nav [data-door]').forEach(b => { b.onclick = () => this.door(b.dataset.door); });
    document.querySelectorAll('#ag-tabs [data-tab]').forEach(b => { b.onclick = () => { this.tab = b.dataset.tab; this.paintPanel(); }; });
    on('ag-insp', 'click', () => { this.open = !this.open; this.paintPanel(); });
    on('ag-pclose', 'click', () => { this.open = false; this.paintPanel(); });
    on('ag-lock', 'click', () => Lock.lock());
    on('ag-clear', 'click', () => { this.entries = []; this.paintOut(); });
    on('ag-go', 'click', () => this.go());
    on('ag-mic', 'click', () => this.mic());
    on('ag-help', 'click', () => { const h = $ag('ag-hints'); h.hidden = !h.hidden; $ag('ag-help').classList.toggle('on', !h.hidden); });
    on('ag-pages', 'click', (e) => { e.stopPropagation(); this.menu(); });
    this.paintMode();

    const input = $ag('ag-input');
    input.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') { e.preventDefault(); this.go(); }
      else if (e.key === 'ArrowUp') {
        e.preventDefault();
        if (this.hi + 1 < this.hist.length) { this.hi++; input.value = this.hist[this.hi]; }
      } else if (e.key === 'ArrowDown') {
        e.preventDefault();
        if (this.hi > 0) { this.hi--; input.value = this.hist[this.hi]; }
        else { this.hi = -1; input.value = ''; }
      } else if (e.ctrlKey && (e.key === 'l' || e.key === 'L')) {
        e.preventDefault(); this.entries = []; this.paintOut();
      }
    });
    input.addEventListener('input', () => { $ag('ag-go').disabled = !input.value.trim(); });
    $ag('ag-go').disabled = true;

    // The terminal is re-drawn whole on most events, so its buttons are
    // listened for here, once, and not on each redraw.
    const out = $ag('ag-out');
    out.addEventListener('click', (e) => {
      const b = e.target.closest('[data-act]');
      if (!b) return;
      const act = b.dataset.act;
      if (act === 'boot') this.boot();
      else if (act === 'close') this.closeSitting();
      else if (act === 'cancel') Run.cancel();
      else if (act === 'copy') {
        const en = this.entries.find(x => String(x.id) === b.dataset.id);
        if (en && navigator.clipboard) navigator.clipboard.writeText(this.textOf(en)).then(() => toast('Copied'), () => {});
      } else if (act === 'run') { this.open = true; this.tab = 'run'; this.paintPanel(); }
      else if (act === 'aider-undo') this.aiderUndo(b.dataset.run, '/undo ' + b.dataset.run);
      else if (act === 'aider-said') {
        const en = this.entries.find(x => String(x.id) === b.dataset.id);
        if (en) { en.showSaid = !en.showSaid; this.paintOut(false); }
      }
      else if (act === 'shell-yes' || act === 'shell-no') {
        // THE ONLY PLACE A SHELL CARD IS ANSWERED FROM: a button he presses (RULE 6).
        const en = this.entries.find(x => String(x.id) === b.dataset.id);
        if (en) this.decide(en, act === 'shell-yes' ? 'approve' : 'deny');
      }
    });
    out.addEventListener('submit', (e) => {
      if (e.target.id !== 'ag-gate-form') return;
      e.preventDefault();
      this.answer(e.target.querySelector('input').value);
    });
    $ag('ag-title').addEventListener('click', (e) => {
      const b = e.target.closest('[data-act]');
      if (!b) return;
      if (b.dataset.act === 'boot') this.boot(); else if (b.dataset.act === 'close') this.closeSitting();
      else if (b.dataset.act === 'reset') this.shellReset();
    });
    $ag('ag-pane').addEventListener('click', (e) => this.paneClick(e));
    $ag('ag-pane').addEventListener('input', (e) => { if (e.target.id === 'ag-tools-q') this.filterTools(e.target.value); });
    if (!this._esc) {
      this._esc = true;
      document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && $ag('ag-menu') && !$ag('ag-menu').hidden) $ag('ag-menu').hidden = true;
      });
      document.addEventListener('click', () => { const m = $ag('ag-menu'); if (m) m.hidden = true; });
    }
  },

  loadHist() { try { const r = localStorage.getItem(this.HIST_KEY); if (r) this.hist = JSON.parse(r).slice(0, 100); } catch { /* per-viewer convenience only */ } },
  saveHist() { try { localStorage.setItem(this.HIST_KEY, JSON.stringify(this.hist.slice(0, 100))); } catch { /* per-viewer convenience only */ } },

  // ---- the bar ------------------------------------------------------------

  // Who is signed in, from the lock itself (/api/me), beside the button that
  // locks the console again.
  async who() {
    try {
      const r = await fetch('/api/me');
      if (!r.ok) return;
      const me = await r.json();
      const el = $ag('ag-who');
      if (el && me && me.user) el.textContent = me.user;
    } catch { /* the lock screen says the rest */ }
  },

  version() {
    const v = document.getElementById('version');
    const el = $ag('ag-ver');
    if (el && v && v.textContent && v.textContent !== 'loading...') el.textContent = v.textContent;
    else setTimeout(() => { const e2 = $ag('ag-ver'); const v2 = document.getElementById('version'); if (e2 && v2 && v2.textContent !== 'loading...') e2.textContent = v2.textContent; }, 1500);
  },

  // THE STATUS LINE SAYS ONLY WHAT A READ ANSWERED: the door and the engine
  // from Run's last check, the rack from rack_list. Where the prototype this
  // page came from showed a tokens-per-second figure, this shows what the rack
  // reported; a figure the record cannot prove is not shown.
  paintTop() {
    const dot = $ag('ag-dot');
    if (!dot) return;
    let cls, word;
    if (Run.unreachable) { cls = 'bad'; word = 'door silent'; }
    else if (!Run.engineOpen) { cls = 'warn'; word = 'no engine'; }
    else { cls = 'ok'; word = 'engine open'; }
    dot.className = 'ag-dot ' + cls;
    $ag('ag-sdot').className = 'ag-sdot ' + cls;
    $ag('ag-stat-a').textContent = word;
    const b = $ag('ag-stat-b');
    if (!this.rack) b.textContent = 'rack: reading';
    else b.textContent = this.rack.ok ? 'rack ' + this.rack.n + ' voice' + (this.rack.n === 1 ? '' : 's') : 'rack silent';
    $ag('ag-stat').title = this.rack && !this.rack.ok ? this.rack.says : 'rack_list, read just now';
    const waiting = !!this.asking();
    const held = this.holdState ? this.holdState.n : 0;
    $ag('ag-adot').hidden = !(waiting || held > 0);
    $ag('ag-insp').title = waiting ? 'The council is asking you something'
      : held ? held + ' writing call' + (held === 1 ? '' : 's') + ' waiting for your hand'
      : 'Toggle the inspector';
  },

  menu() {
    const m = $ag('ag-menu');
    if (!m.hidden) { m.hidden = true; return; }
    const links = [...document.querySelectorAll('.nav-link')].filter(a => a.dataset.page !== 'agent');
    m.innerHTML = links.map(a => `<a href="${escHtml(a.getAttribute('href'))}">${escHtml((a.childNodes[0].textContent || '').trim())}</a>`).join('');
    m.querySelectorAll('a').forEach(a => {
      a.onclick = (e) => { e.preventDefault(); m.hidden = true; history.pushState(null, '', a.getAttribute('href')); App.router(); };
    });
    m.hidden = false;
  },

  // ---- the terminal window --------------------------------------------------

  paintMode() {
    const box = $ag('ag-modes');
    if (!box) return;
    box.innerHTML = this.MODES.map(m =>
      `<button type="button" data-mode="${m.id}" class="${m.id === this.mode ? 'on' : ''}${m.wired ? '' : ' off'}"
        title="${escHtml(m.tip)}"><span>${escHtml(m.label)}</span><small>${escHtml(m.badge)}</small></button>`).join('');
    box.querySelectorAll('[data-mode]').forEach(b => {
      b.onclick = () => {
        const m = this.MODES.find(x => x.id === b.dataset.mode);
        if (!m.wired) { this.say('(' + m.label + ' tab)', this.NOT_WIRED[m.id], 'ERROR'); return; }
        this.mode = m.id; this.paintMode(); this.paintTitle(); this.paintOut(); $ag('ag-input').focus();
        if (m.id === 'aider') this.readAider();
      };
    });
    this.paintInput();
  },

  // The line under the terminal belongs to the tab: its prompt, its placeholder and what it
  // offers to try.
  paintInput() {
    const pr = $ag('ag-prompt');
    if (pr) { pr.className = 'ag-p ' + this.mode; pr.innerHTML = this.promptText(); }
    const input = $ag('ag-input');
    if (input) input.placeholder = this.PLACEHOLDER[this.mode] || this.PLACEHOLDER.agent;
    const box = $ag('ag-hints');
    if (box) {
      box.innerHTML = '<span>Try:</span>' + (this.HINTS[this.mode] || []).map(h =>
        `<button type="button" data-say="${escHtml(h)}">${escHtml(h)}</button>`).join('');
      box.querySelectorAll('[data-say]').forEach(b => { b.onclick = () => this.go(b.dataset.say); });
    }
  },

  // The prompt a line is typed after. An entry keeps the prompt it was typed under, so the
  // scrollback still says which tab each line went to.
  promptText(kind) {
    const k = kind || this.mode;
    const w = escHtml(Run.world || 'research');
    if (k === 'bash') return 'bash[' + w + ']$';
    if (k === 'python') return 'py[' + w + ']&gt;&gt;&gt;';
    if (k === 'aider') return 'aider[' + w + ']&gt;';
    return 'agent[' + w + ']&gt;';
  },

  // Bash and Python are the operator's typed shell; they run on the door and need no engine.
  isShell(k) { return (k || this.mode) === 'bash' || (k || this.mode) === 'python'; },

  paintTitle() {
    const el = $ag('ag-title');
    if (!el) return;
    const w = escHtml(Run.world || 'research');
    let state;
    if (Run.unreachable) {
      state = '<span class="ag-bad">door silent</span>';
    } else if (this.isShell()) {
      // A shell needs the door, not an engine: it says what it is, and the Python session can be ended.
      state = this.mode === 'bash'
        ? '<span title="Git Bash, run in this world as you; the law gate is walked before every run">git bash &middot; gated</span>'
        : '<span title="One Python process for this world; it keeps its names until it is reset, runs past its limit, or the door restarts">keeps its names &middot; gated</span>' +
          '<button type="button" class="ag-link" data-act="reset">Reset</button>';
    } else if (this.isAider()) {
      state = this.aiderTitleState();
    } else if (!Run.engineOpen) {
      state = '<span class="ag-warn">no engine</span><button type="button" class="ag-link" data-act="boot">Boot</button>';
    } else {
      state = '<span title="' + (Run.runs || 0) + ' turn' + (Run.runs === 1 ? '' : 's') + ' run in this sitting">sitting ' +
        escHtml(String(Run.sitting || '?')) + ' &middot; ' + escHtml(this.age()) + '</span>' +
        (Run.stale ? '<span class="ag-warn" title="' + escHtml((Run.staleFile || 'manjuel') + ' changed after this engine started; a reboot picks it up (seats, skills and pipelines reload without one)') + '">old code</span>' : '') +
        '<button type="button" class="ag-link" data-act="close">Close</button>';
    }
    el.innerHTML = AG_ICON.cpu + '<b>' + escHtml(this.mode) + '</b><span class="ag-sl">/</span><span>project:' + w + '</span><span class="ag-sl">&middot;</span>' + state;
    const pr = $ag('ag-prompt');
    if (pr) { pr.className = 'ag-p ' + this.mode; pr.innerHTML = this.promptText(); }
  },

  // How long this engine has stood, said in words a person reads, from the
  // time the door says it started.
  age() {
    const ms = Date.now() - new Date(Run.started).getTime();
    if (!isFinite(ms) || ms < 0) return '';
    const s = Math.floor(ms / 1000);
    return s < 60 ? s + 's' : s < 3600 ? Math.floor(s / 60) + 'm' : Math.floor(s / 3600) + 'h ' + Math.floor((s % 3600) / 60) + 'm';
  },

  bannerHtml() {
    const w = escHtml(Run.world || 'research');
    if (this.isShell()) {
      const world = Run.world || 'research';
      const lines = this.mode === 'bash'
        ? ['Git Bash, run in ' + world + ' as you. A plain look (ls, cat, grep, git status, git log, git diff) runs at once.',
           'Anything that writes, deletes, installs, reaches the network, or saves or sends work shows a card here and waits for your click. A secret file, client material, a key typed into a command or a path outside the ground is refused by name, and no click lifts a refusal.',
           'The law is walked before every run, the command runs with none of the door\'s keys in its environment, what it prints is scrubbed of every secret the ground holds, and each run is written to the ground\'s holds log.']
        : ['Python, in one session for ' + world + ' that keeps its names between entries: x = 41, then x + 1.',
           'A plain calculation runs at once. An import, a file, a call to anything but a plain-data function, a def or a class shows a card and waits for your click. An entry that runs past 30 s ends the session and its names go; Reset (or /reset) ends it yourself.',
           'The same gate as Bash: secrets and paths outside the ground are refused by name, the session holds none of the door\'s keys, and what it prints is scrubbed.'];
      return '<div class="ag-e"><pre class="ag-pre dim">' + escHtml(Run.unreachable
        ? 'The door did not answer. A shell runs on the door, so nothing below can run until it does.' : lines.join('\n')) + '</pre></div>';
    }
    if (this.isAider()) return this.aiderBannerHtml();
    const l1 = Run.unreachable
      ? 'The door did not answer. Nothing below can run until it does.'
      : Run.engineOpen
        ? 'The council stands on ' + (Run.world || 'research') + ' (sitting ' + (Run.sitting || '?') + '). What you type goes to it: the law gate stamps it before any model reads a word, the router runs the tools, and every failure is shown with the answer.'
        : 'No engine is open on ' + (Run.world || 'research') + '. Press Boot, or type /boot: booting opens a sitting, and closing it writes its end and pays the toll if a turn ran.';
    const l2 = 'Your words go to the council as they are, and so does any /command the engine knows (try /help or /status). The page itself answers /boot, /close and /clear.';
    const l3 = 'Bash, Python and Aider are wired in the tabs above: Bash and Python run what you type as you, and anything that writes asks first; Aider edits the files you name, on a line of work.';
    void w;
    return '<div class="ag-e"><pre class="ag-pre dim">' + escHtml(l1 + '\n' + l2 + '\n' + l3) + '</pre></div>';
  },

  // What an entry says, in one place: a live or finished turn, a command's
  // printed text, a refusal, or the words kept from an earlier visit.
  // The console prints a rule and "Objective: ..." before every command, and the
  // prompt line above the output already says what was asked.
  bare(s) {
    return String(s || '').replace(/^\s*={10,}[ \t]*\n[ \t]*Objective:[^\n]*\n/, '').replace(/^(\s*\n)+/, '').replace(/\s+$/, '');
  },

  textOf(en) {
    if (en.kind === 'local') return en.text || '';
    if (en.kind === 'shell') return (en.ans && en.ans.output) || '';
    if (en.kind === 'aider') return en.note || (en.ans && (en.ans.diff || en.ans.said)) || '';
    const t = en.cn && en.cn.turn;
    if ((en.isCommand || en.isAnswer) && en.cmdText != null && this.bare(en.cmdText)) return this.bare(en.cmdText);
    if (t) {
      if (t.answer) return t.answer;
      if (t.refusal) return t.refusal;
      return '';
    }
    return en.cn ? (en.cn.text || '') : '';
  },

  entryHtml(en) {
    const id = en.id;
    if (en.kind === 'shell') return this.shellEntryHtml(en);
    if (en.kind === 'aider') return this.aiderEntryHtml(en);
    if (en.kind === 'local') {
      const bad = en.status === 'ERROR';
      return `<div class="ag-e" id="ag-e-${id}"><div class="ag-line"><div class="ag-line-l"><span class="ag-p agent">${this.promptText('agent')}</span><span class="ag-cmd">${escHtml(en.command)}</span></div>
        <div class="ag-meta"><span class="ag-chip${bad ? ' bad' : ''}">${en.status === 'RUNNING' ? 'running' : 'page'}</span>${en.ms != null ? '<span>' + en.ms + 'ms</span>' : ''}</div></div>
        <pre class="ag-pre${bad ? ' bad' : ''}" id="ag-o-${id}">${escHtml(en.text || '')}</pre></div>`;
    }
    const t = en.cn && en.cn.turn;
    const live = !!(en.cn && en.cn.live);
    let route = '';
    if (t) {
      if (t.delivery) route = t.pipeline || 'delivered';
      else if (t.refusal) route = t.verdict || 'refused';
      else if (live) route = 'running';
      else if (t.waiting) route = 'waiting';
      else if (en.isAnswer) route = 'answer';
      else if (en.isCommand) route = 'command';
      else route = t.verdict || '';
    } else if (en.isAnswer) route = 'answer';
    else if (en.isCommand) route = 'command';
    const failed = /refused|aborted|cancelled|unreachable|broke/.test(route);
    const chips = [];
    if (route) chips.push(`<span class="ag-chip${failed ? ' bad' : ''}">${escHtml(route)}</span>`);
    if (t) {
      const nt = (t.tools || []).length;
      if (nt) chips.push(`<span class="ag-note">${nt} tool${nt === 1 ? '' : 's'}</span>`);
      chips.push(`<span id="ag-m-${id}">${escHtml(Run.elapsed(t))}</span>`);
    }
    const fails = t ? Run.failures(t) : [];
    const oot = t ? Run.outOfTime(t) : [];
    const lines = [];
    if (fails.length) lines.push(`<div class="ag-notrun"><b>NOT EVERYTHING RAN</b><br>${fails.map(escHtml).join('<br>')}<br><span class="ag-src">machine-emitted from what happened, not a seat's account of it</span></div>`);
    if (oot.length) lines.push(`<div class="ag-notrun"><b>OUT OF TIME</b><br>${oot.map(escHtml).join('<br>')}</div>`);
    if (t && t.dropped) lines.push(`<div class="ag-notrun"><b>${t.dropped} events were dropped</b> - this page did not see everything that ran</div>`);
    let foot = '';
    if (t && t.delivery && !live) {
      const ns = Run.steps(t).length || (t.seats || []).length;
      foot = `<div class="ag-foot-line"><span>${ns} seat${ns === 1 ? '' : 's'}</span>${t.transcript ? '<span>transcript ' + escHtml(t.transcript) + '</span>' : ''}<button type="button" data-act="run">open the run</button></div>`;
    }
    const deep = !!(t && t.delivery);
    const body = this.textOf(en);
    return `<div class="ag-e" id="ag-e-${id}"><div class="ag-line"><div class="ag-line-l"><span class="ag-p agent">${this.promptText('agent')}</span><span class="ag-cmd">${en.command === '' && en.isAnswer ? '(blank answer)' : escHtml(en.command)}</span></div>
      <div class="ag-meta">${chips.join('')}<button type="button" class="ag-copy" data-act="copy" data-id="${id}" title="Copy the output" aria-label="Copy the output">${AG_ICON.copy}</button></div></div>
      ${body ? `<pre class="ag-pre${failed ? ' bad' : deep ? ' deep' : ''}" id="ag-o-${id}">${escHtml(body)}</pre>` : ''}${lines.join('')}${foot}</div>`;
  },

  // A line typed in a shell tab, as the door answered it. The page decides nothing: it draws the
  // answer's state, class, reasons, hold and output, and a failure (a non-zero exit, a refusal, a run
  // that was ended) is on the face of it. A card is drawn only for what the door said was held.
  shellEntryHtml(en) {
    const id = en.id;
    const a = en.ans || null;
    const kind = en.shell === 'python' ? 'python' : 'bash';
    const why = (a && a.why) || [];
    const chips = [];
    let bad = false;
    if (en.state === 'running') chips.push('<span class="ag-chip">running</span>');
    else if (en.state === 'held') chips.push('<span class="ag-chip">asks first</span>');
    else if (en.state === 'denied') chips.push('<span class="ag-chip bad">denied</span>');
    else if (en.state === 'gone') chips.push('<span class="ag-chip">answered elsewhere</span>');
    else if (en.state === 'refused') { chips.push('<span class="ag-chip bad">refused</span>'); bad = true; }
    else if (en.state === 'error') { chips.push('<span class="ag-chip bad">could not run</span>'); bad = true; }
    else if (en.state === 'ran' && a) {
      chips.push('<span class="ag-chip">' + (a.approved ? 'approved by you' : escHtml(a.class || 'ran')) + '</span>');
      if (a.exit != null) {
        const nz = a.exit !== 0;
        if (nz) bad = true;
        chips.push('<span class="ag-chip' + (nz ? ' bad' : '') + '">exit ' + escHtml(String(a.exit)) + '</span>');
      }
    }
    // Only a run or a refusal has a time: a card that waited for his hand has not taken that long to do anything.
    if (en.state === 'ran' || en.state === 'refused') {
      const ms = a && a.ms != null ? a.ms : en.ms;
      if (ms != null) chips.push('<span>' + ms + 'ms</span>');
    }
    const out = this.textOf(en);
    const lines = [];
    if (en.state === 'held') {
      lines.push(`<div class="ag-gate"><b>This asks first</b>
        <div class="ag-gate-q">${escHtml(why.map(w => '- ' + w).join('\n'))}</div>
        <div class="ag-btns"><button type="button" class="ag-btn go" data-act="shell-yes" data-id="${id}"${en.deciding ? ' disabled' : ''}>${en.deciding ? 'Running...' : 'Approve - run it'}</button>
          <button type="button" class="ag-btn no" data-act="shell-no" data-id="${id}"${en.deciding ? ' disabled' : ''}>Deny</button></div>
        <span class="ag-src">parked at the door as ${escHtml((a && a.hold) || '')}; approving runs exactly this entry and nothing else</span></div>`);
    }
    if (en.state === 'refused') lines.push(`<div class="ag-notrun"><b>REFUSED BY NAME</b><br>${why.map(escHtml).join('<br>')}<br><span class="ag-src">no click lifts a refusal; nothing ran</span></div>`);
    if (en.state === 'denied') lines.push('<div class="ag-notrun"><b>DENIED</b> by you. Nothing ran.</div>');
    if (en.state === 'gone') lines.push('<div class="ag-src">This was answered somewhere else (the Guardrails tab shows what it said), or the door restarted and dropped it. Nothing more happens from here.</div>');
    if (en.state === 'error') lines.push(`<div class="ag-notrun"><b>COULD NOT RUN</b><br>${escHtml(en.err || '')}</div>`);
    if (a && a.timed_out) lines.push(`<div class="ag-notrun"><b>ENDED</b> ${escHtml(a.note || 'it ran past its limit')}</div>`);
    else if (a && a.note) lines.push(`<div class="ag-src">${escHtml(a.note)}</div>`);
    if (a && a.truncated) lines.push('<div class="ag-src">the output was cut at 64 KB</div>');
    return `<div class="ag-e" id="ag-e-${id}"><div class="ag-line"><div class="ag-line-l"><span class="ag-p ${kind}">${en.prompt || this.promptText(kind)}</span><span class="ag-cmd">${escHtml(en.command)}</span></div>
      <div class="ag-meta">${chips.join('')}<button type="button" class="ag-copy" data-act="copy" data-id="${id}" title="Copy the output" aria-label="Copy the output">${AG_ICON.copy}</button></div></div>
      ${out ? `<pre class="ag-pre${bad ? ' bad' : ''}" id="ag-o-${id}">${escHtml(out)}</pre>` : ''}${lines.join('')}</div>`;
  },

  // THE GATE, IN THE THREAD. The council stopped to ask, so the question is
  // here with a field under it. Nothing is assumed on his behalf (RULE 6).
  gateHtml() {
    const asking = this.asking();
    if (!asking || Run.running) return '';
    return `<div class="ag-gate"><b>The council is asking</b><div class="ag-gate-q">${escHtml(asking)}</div>
      <form id="ag-gate-form"><input type="text" autocomplete="off" placeholder="Your answer - nothing is assumed on your behalf" /><button class="ag-btn go" type="submit">Answer</button></form></div>`;
  },

  liveHtml() {
    if (!Run.running) return '';
    return `<div class="ag-live"><span class="pulse">&#9679;</span><span>${escHtml(Run.nowLine())} &middot; ${escHtml(Run.elapsed())}</span><button type="button" data-act="cancel">cancel</button></div>`;
  },

  // WHAT HE IS TYPING INTO THE GATE SURVIVES A REDRAW. The screen is redrawn
  // on every state check, and an answer half written when one landed would
  // otherwise vanish; the field takes focus once, when the question first
  // appears, and is left alone after that.
  paintOut(stick = true) {
    const box = $ag('ag-out');
    if (!box) return;
    const near = box.scrollHeight - box.scrollTop - box.clientHeight < 90;
    const old = box.querySelector('#ag-gate-form input');
    const typed = old ? old.value : '';
    const hadFocus = !!old && document.activeElement === old;
    box.innerHTML = this.bannerHtml() + this.entries.map(e => this.entryHtml(e)).join('') +
      this.gateHtml() + '<div id="ag-live">' + this.liveHtml() + '</div>';
    if (stick || near) this.stick();
    const gi = box.querySelector('#ag-gate-form input');
    if (!gi) { this._gateSeen = false; return; }
    if (typed) gi.value = typed;
    if (!this._gateSeen || hadFocus) { gi.focus(); this._gateSeen = true; }
  },

  // The council has stopped to ask. Until he answers, the engine takes nothing
  // else: an objective sent now would queue behind the question and hang.
  asking() { return (Run.turn && Run.turn.waiting) || Run.pending || ''; },

  stick() { const box = $ag('ag-out'); if (box) box.scrollTop = box.scrollHeight; },

  // One token must not redraw the screen: only the live entry's words and the
  // live line move.
  paintLive() {
    const en = this.live();
    if (!en) return;
    const pre = $ag('ag-o-' + en.id);
    if (!pre) { this.paintOut(); return; }
    pre.textContent = this.textOf(en);
    const live = $ag('ag-live');
    if (live) live.innerHTML = this.liveHtml();
    const box = $ag('ag-out');
    if (box && box.scrollHeight - box.scrollTop - box.clientHeight < 140) this.stick();
  },

  live() {
    for (let i = this.entries.length - 1; i >= 0; i--) {
      const e = this.entries[i];
      if (e.cn && e.cn.live) return e;
    }
    return null;
  },

  // A seat thinking says nothing on the wire, sometimes for a minute; the
  // clock keeps moving or a working engine reads as a hung one.
  tick(on) {
    clearInterval(this._tick);
    this._tick = null;
    if (!on) return;
    this._tick = setInterval(() => {
      if (!Run.running || !$ag('ag-out')) { this.tick(false); return; }
      const live = $ag('ag-live');
      if (live) live.innerHTML = this.liveHtml();
    }, 1000);
  },

  // Turns started on another page or in another window are this conversation
  // too: the shared thread is read, and a pair this screen has not drawn is
  // drawn.
  adopt() {
    const known = new Set(this.entries.filter(e => e.cn).map(e => e.cn));
    const th = Chat.thread || [];
    for (let i = 0; i < th.length; i++) {
      const m = th[i];
      if (m.who !== 'him') continue;
      const cn = th[i + 1] && th[i + 1].who === 'council' ? th[i + 1] : null;
      if (!cn || known.has(cn)) continue;
      this.entries.push({ id: ++this._seq, kind: 'turn', command: m.text, him: m, cn, isCommand: /^\//.test(m.text || '') });
      known.add(cn);
    }
  },

  // ---- what Run tells us ------------------------------------------------------

  onRun(what) {
    if (!$ag('ag-out')) return;                       // not the page on screen
    if (what === 'state') { this.paintTitle(); this.paintTop(); this.paintOut(false); return; }
    if (what === 'start') {
      if (Run.turn && Run.turn.watching && !(Chat.thread || []).some(m => m.live)) {
        const him = { who: 'him', text: Run.turn.objective || '' };
        const cn = { who: 'council', text: '', live: true };
        Chat.thread.push(him, cn);
        this.entries.push({ id: ++this._seq, kind: 'turn', command: him.text, him, cn, isCommand: /^\//.test(him.text) });
      }
      this.tick(true);
      this.paintOut();
      this.paintTop();
      return;
    }
    if (what === 'event') {
      const en = this.live();
      if (en) en.cn.turn = Run.turn;
      this.paintLive();
      if (this.open && this.tab === 'run') this.schedulePane();
      return;
    }
    if (what === 'end') {
      const en = this.live();
      if (en) { en.cn.turn = Run.turn; en.cn.live = false; }
      this.tick(false);
      this.paintOut();
      this.paintTop();
      if (this.open && this.tab === 'run') this.paintPane();
      Home.keepThread();
      const made = Run.turn && Run.turn.delivery;
      if (made && typeof made.project === 'string' && made.project) this.playProject(made.project);
      Run.check().then(() => { this.paintTitle(); this.paintTop(); this.paintOut(false); });
      this.readHolds();
      const input = $ag('ag-input');
      if (input && !document.querySelector('#ag-gate-form')) input.focus();
    }
  },

  // THE PLAY FRAME (WHAT'S LEFT E1, 2026-10-05). A person asked for a game and was handed a file path to open: this page never read the delivery's
  // `project`, so the game she asked for was on another page. The delivery names the project in hand (read off the maker, never off a seat's words);
  // the page asks the door's read-only `projects` tool which version that is and frames the page the glass already serves, sandboxed by the header it
  // is served under (handlers/projects.go) and, as there, with no `sandbox` attribute (the app's own browser pane refuses any frame that has one).
  // The frame lives OUTSIDE #ag-out, which is rewritten whole on most events: a frame inside it would restart the game on every line she types. The
  // same page is never loaded twice, so a game she is playing keeps playing while she talks; a new version of it, or another project, is loaded.
  async playProject(name) {
    let version = 0;
    try {
      const d = JSON.parse(await App.tool('projects', { action: 'list' }, true));
      const p = (d.projects || []).find(x => x.name === name);
      version = p && Array.isArray(p.versions) ? p.versions.length : 0;
    } catch { /* the page as it stands */ }
    const s = this.play;
    if (s && s.name === name && s.version === version) return;
    this.play = { name, version, closed: false };
    this.paintPlay();
  },

  paintPlay() {
    const box = $ag('ag-play');
    if (!box) return;
    const s = this.play;
    if (!s || s.closed) { box.hidden = true; box.innerHTML = ''; return; }
    const url = Projects.pageUrl(s.name, s.version);
    box.hidden = false;
    box.innerHTML = '<div class="ag-play-bar"><b>' + escHtml(s.name) + '</b><span>' + (s.version ? 'version ' + s.version : 'as it stands') + '</span>' +
      '<span class="ag-play-sp"></span>' +
      '<a href="' + escHtml(url) + '" target="_blank" rel="noopener noreferrer" title="Opens the page alone, in its own tab">Open in its own tab</a>' +
      '<button type="button" id="ag-play-x" title="Put it away; it comes back when a new version is made" aria-label="Put it away">Close</button></div>' +
      '<div class="ag-play-body"><iframe title="' + escHtml(s.name) + '" referrerpolicy="no-referrer" src="' + escHtml(url) + '"></iframe></div>';
    const x = $ag('ag-play-x');
    if (x) x.onclick = () => { s.closed = true; this.paintPlay(); };
    this.fitPlay();
    if (typeof ResizeObserver === 'function') {
      if (this._po) this._po.disconnect();
      const body = box.querySelector('.ag-play-body');
      if (body) { this._po = new ResizeObserver(() => this.fitPlay()); this._po.observe(body); }
    }
  },

  // The whole game on screen: the frame is a window of PLAY_W by PLAY_H, shrunk (never enlarged) to the room the stage has and centred in it. Clicks and keys
  // reach the page through the shrinking, as they do through any scaled frame.
  fitPlay() {
    const box = $ag('ag-play');
    const body = box && box.querySelector && box.querySelector('.ag-play-body');
    const fr = body && body.querySelector('iframe');
    if (!fr) return;
    const w = body.clientWidth, h = body.clientHeight;
    if (!w || !h) return;
    const k = Math.min(w / this.PLAY_W, h / this.PLAY_H, 1);
    fr.style.width = this.PLAY_W + 'px';
    fr.style.height = this.PLAY_H + 'px';
    fr.style.transformOrigin = '0 0';
    fr.style.transform = 'scale(' + k + ')';
    fr.style.left = Math.max(0, (w - this.PLAY_W * k) / 2) + 'px';
    fr.style.top = Math.max(0, (h - this.PLAY_H * k) / 2) + 'px';
  },

  schedulePane() {
    if (this._sp) return;
    this._sp = setTimeout(() => { this._sp = null; if (this.open && this.tab === 'run') this.paintPane(); }, 250);
  },

  // ---- sending ------------------------------------------------------------

  // A line may name its tab: `$ ls` or `! ls` is Bash, `>>> 1 + 1` is Python, from any tab.
  prefixMode(raw) {
    let m;
    if ((m = /^(?:\$|!)\s+([\s\S]*)$/.exec(raw))) return { mode: 'bash', text: m[1] };
    if ((m = /^>>>\s+([\s\S]*)$/.exec(raw))) return { mode: 'python', text: m[1] };
    if ((m = /^aider\s+([\s\S]*)$/i.exec(raw))) return { mode: 'aider', text: m[1] };
    return null;
  },

  async go(rawIn) {
    const input = $ag('ag-input');
    const raw = (rawIn != null ? rawIn : (input ? input.value : '')).trim();
    if (!raw) return;
    if (input) { input.value = ''; $ag('ag-go').disabled = true; }
    this.hist = [raw].concat(this.hist.filter(h => h !== raw));
    this.hi = -1;
    this.saveHist();

    if (raw === '/clear') { this.entries = []; this.paintOut(); return; }
    if (raw === '/boot') return this.boot();
    if (raw === '/close') return this.closeSitting();
    // A shell line: typed in a shell tab, or prefixed from any tab. A shell needs the door, not an
    // engine, and it does not wait behind the council.
    const via = this.prefixMode(raw);
    // The Aider tab's line goes to the door's Aider, and needs no engine either.
    if (via ? via.mode === 'aider' : this.isAider()) return this.aiderGo(via ? via.text : raw);
    const shell = via ? via.mode : (this.isShell() ? this.mode : '');
    if (shell) {
      if (raw === '/reset' && shell === 'python') return this.shellReset();
      return this.shellRun(shell, via ? via.text : raw);
    }
    if (Run.running) { toast('A turn is already running', 'error'); return; }
    if (this.asking()) {
      return this.say(raw, 'The council is waiting on your answer above, and it takes nothing else until you give it. Answer it first. Nothing was sent.', 'ERROR');
    }
    if (!Run.engineOpen) {
      return this.say(raw, Run.unreachable
        ? 'The door did not answer, so nothing can run. Nothing was sent.'
        : 'No engine is open on ' + (Run.world || 'this world') + ', so nothing was sent. Press Boot (or type /boot) first: that opens a sitting, and the sitting line is the lock.', 'ERROR');
    }
    this.send(raw);
  },

  // ---- the shell tabs -------------------------------------------------------------

  async shellRun(kind, text) {
    const en = { id: ++this._seq, kind: 'shell', shell: kind, prompt: this.promptText(kind), command: text, state: 'running', t0: Date.now() };
    this.entries.push(en);
    this.paintOut();
    await this.shellAsk(en, { shell: kind, command: text });
  },

  // /reset in the Python tab, or its Reset button: the session ends and its names go.
  async shellReset() {
    const en = { id: ++this._seq, kind: 'shell', shell: 'python', prompt: this.promptText('python'), command: '/reset', state: 'running', t0: Date.now() };
    this.entries.push(en);
    this.paintOut();
    await this.shellAsk(en, { shell: 'python', reset: true });
  },

  async shellAsk(en, args) {
    let said;
    try { said = await App.tool('shell_run', args); }
    catch (e) { en.state = 'error'; en.err = e.message || 'refused'; en.ms = Date.now() - en.t0; this.paintOut(); return; }
    this.shellTake(en, said);
  },

  // What the door answered, in the one shape it answers in (state, class, why, hold, exit, output).
  shellTake(en, said) {
    let a = null;
    try { a = JSON.parse(said); } catch { /* reported below */ }
    en.ms = Date.now() - en.t0;
    if (!a || typeof a !== 'object' || !a.state) {
      en.state = 'error';
      en.err = 'The door answered something that is not the shell\'s shape: ' + String(said).slice(0, 300);
    } else { en.ans = a; en.state = a.state; }
    this.readHolds();
    this.paintOut();
  },

  // A CARD IS ANSWERED BY A BUTTON HE PRESSES AND NOWHERE ELSE (RULE 6). Approving runs, at the door,
  // exactly the entry that was parked and re-judges it first: a refusal is refused again. The replay's
  // answer is the shell's own document, after the words hold_answer puts before it.
  async decide(en, decision) {
    const a = en.ans;
    if (!a || !a.hold || en.deciding) return;
    en.deciding = true;
    this.paintOut();
    let said;
    try { said = await App.tool('hold_answer', { id: a.hold, decision }); }
    catch (e) { said = 'Refused: ' + (e.message || 'refused'); }
    en.deciding = false;
    en.ms = Date.now() - en.t0;
    if (decision === 'deny') {
      en.state = 'denied';
    } else {
      const m = /ran:\s*(\{[\s\S]*\})\s*$/.exec(said);
      let ran = null;
      if (m) { try { ran = JSON.parse(m[1]); } catch { /* shown as said, below */ } }
      if (ran && ran.state) { en.ans = ran; en.state = ran.state; }
      else { en.state = 'error'; en.err = said; }
    }
    this.readHolds();
    this.paintOut();
  },

  // ---- the Aider tab ----------------------------------------------------------------
  //
  // AIDER IS THE DOOR'S AIDER (WHAT'S LEFT H12 and B20, 2026-10-04; his word: "finish wiring in the Aider panel, it's on
  // this system"). A line typed here is an instruction to the door's aider_run, which only his glass may call: Aider runs
  // headless on a COPY of the files he named, behind a wall, against the coding seat's own model on this machine, and the
  // door writes its edit back only while the world stands on a line of work and every Python file still parses. The page
  // keeps the NAMES of the files in the chat and judges nothing about them: a path that may not be written, a main line,
  // a model that is not on the rack are refused at the door, by name, and the page draws the refusal. The Undo button
  // and /undo are the door's aider_undo; nothing here commits, runs the suites or lands (RULE 6): the edit is unsaved
  // work on the line, and Version control is where it is saved.

  isAider(k) { return (k || this.mode) === 'aider'; },

  // What the door says of Aider right now, read quietly. The title bar and the Inspector draw from it; the page never
  // answers for it.
  async readAider() {
    try { this.aiderSt = JSON.parse(await App.tool('aider_status', {}, true)); }
    catch (e) { this.aiderSt = { error: e.message || 'unreadable' }; }
    if (this.mode === 'aider') this.paintTitle();
  },

  aiderTitleState() {
    const s = this.aiderSt;
    if (!s) return '<span>reading...</span>';
    if (s.error) return '<span class="ag-bad" title="' + escHtml(s.error) + '">status unreadable</span>';
    if (!s.installed) return '<span class="ag-warn" title="' + escHtml((s.why || [])[0] || '') + '">not installed</span>';
    const parts = ['<span title="Aider ' + escHtml(s.version || '') + ', headless, behind a wall; the model is the coding seat\'s, on this machine">' + escHtml(s.model || 'no model') + '</span>'];
    parts.push(s.on_line
      ? '<span title="Aider writes only on a line of work">line <b>' + escHtml(s.line) + '</b></span>'
      : '<span class="ag-warn" title="The main line is yours: open a line of work (Version control) and Aider can write on it">' + (s.line ? 'on ' + escHtml(s.line) : 'no line of work') + '</span>');
    if (s.rack !== 'up' || !s.on_rack) parts.push('<span class="ag-bad" title="' + escHtml((s.why || []).join('; ')) + '">model not reachable</span>');
    const n = this.aiderFiles.length;
    const all = this.aiderFiles.concat(this.aiderRead.map(f => f + ' (read-only)'));
    parts.push('<span title="' + escHtml(all.join('\n') || 'none yet: /add path') + '">' + n + ' file' + (n === 1 ? '' : 's') + (this.aiderRead.length ? ' + ' + this.aiderRead.length + ' read' : '') + '</span>');
    if (s.busy) parts.push('<span class="ag-warn">busy</span>');
    return parts.join('<span class="ag-sl">&middot;</span>');
  },

  aiderBannerHtml() {
    const world = Run.world || 'research';
    const s = this.aiderSt;
    const lines = [
      'Aider, run on the files you name in ' + world + ', with the coding seat\'s model on this machine. /add path names a file Aider may change, /read path one it may only read, /drop path takes one out, /files lists them, /undo takes the last run back, /help says all of it. Then say what to change.',
      'It works on a copy, behind a wall: it writes only to its own scratch folder, reaches nothing but this machine, and starts no other process. The changes come back only while ' + world + ' stands on a line of work (never the main line), only if every Python file still parses, and each file keeps its own line endings.',
      'Nothing is saved, run or sent from here. The edit is unsaved work on the line: run the suites, save it in Version control, and the Land click is yours.'
    ];
    let state = '';
    if (Run.unreachable) state = 'The door did not answer, so Aider cannot run.';
    else if (s && s.error) state = 'Aider\'s state could not be read: ' + s.error;
    else if (s && !s.installed) state = 'Aider is not installed in ' + world + ': ' + ((s.why || [])[0] || '');
    else if (s && (s.why || []).length) state = 'Not ready: ' + s.why.join('; ') + '.';
    else if (s) state = 'Aider ' + (s.version || '') + ' is ready: ' + (s.model || '') + ', a window of ' + (s.context || '?') + ' tokens (about ' + (s.budget_kb || '?') + ' KB of files at a time).';
    return '<div class="ag-e"><pre class="ag-pre dim">' + escHtml(lines.join('\n') + (state ? '\n' + state : '')) + '</pre></div>';
  },

  aiderHelp() {
    return [
      '/add path [path...]   files Aider may change (up to the door\'s limit)',
      '/read path [path...]  files Aider may only read, for context',
      '/drop path | all      take a file out of the chat',
      '/files                 what is in the chat now',
      '/status               what the door says of Aider now',
      '/undo [run]           take a run back (the last one by default)',
      'anything else         an instruction: what Aider should change in those files',
      'A line may also name this tab from any other: `aider /files`.'
    ].join('\n');
  },

  aiderFilesText() {
    if (!this.aiderFiles.length && !this.aiderRead.length) return 'No files in the chat yet. /add path/to/file.py names one Aider may change; /read path/to/file.md one it may only read.';
    return 'Aider may change: ' + (this.aiderFiles.join(', ') || '(nothing yet)') + (this.aiderRead.length ? '\nRead only: ' + this.aiderRead.join(', ') : '');
  },

  aiderStatusText() {
    const s = this.aiderSt;
    if (!s) return 'Not read yet.';
    if (s.error) return 'Aider\'s state could not be read: ' + s.error;
    const out = ['installed: ' + (s.installed ? 'yes' + (s.version ? ', ' + s.version : '') : 'no'),
      'model: ' + (s.model || 'none declared') + (s.context ? ', a window of ' + s.context + ' tokens (about ' + (s.budget_kb || '?') + ' KB of files)' : ''),
      'rack: ' + (s.rack === 'up' ? 'up' : 'silent') + (s.model ? (s.on_rack ? ', has the model' : ', does not have the model') : ''),
      'line of work: ' + (s.on_line ? s.line : (s.line ? s.line + ' (the main line is yours)' : 'none')),
      'busy: ' + (s.busy ? 'yes' : 'no')];
    if ((s.why || []).length) out.push('NOT READY: ' + s.why.join('; '));
    return out.join('\n');
  },

  // A line typed in the Aider tab: a /command of the tab's own, or an instruction for the door.
  async aiderGo(text) {
    const t = String(text || '').trim();
    const m = /^\/([A-Za-z]+)\s*([\s\S]*)$/.exec(t);
    if (!m) return this.aiderRun(t);
    const cmd = m[1].toLowerCase();
    const rest = m[2].trim();
    if (cmd === 'add') return this.aiderNames('files', rest, t);
    if (cmd === 'read') return this.aiderNames('read', rest, t);
    if (cmd === 'drop') return this.aiderDrop(rest, t);
    if (cmd === 'files') return this.aiderSay(t, this.aiderFilesText());
    if (cmd === 'status') { await this.readAider(); return this.aiderSay(t, this.aiderStatusText()); }
    if (cmd === 'undo') return this.aiderUndo(rest || this.aiderLast, t);
    if (cmd === 'help') return this.aiderSay(t, this.aiderHelp());
    return this.aiderSay(t, 'The Aider tab has no /' + cmd + '. /help lists what it answers.', 'ERROR');
  },

  aiderSplit(rest) {
    return String(rest || '').split(/[\s,;]+/).map(n => n.replace(/\\/g, '/').replace(/^(\.\/)+/, '')).filter(Boolean);
  },

  aiderNames(kind, rest, typed) {
    const names = this.aiderSplit(rest);
    if (!names.length) return this.aiderSay(typed, 'Name a file: ' + (kind === 'files' ? '/add path/to/file.py' : '/read path/to/file.md') + '. More than one is fine, separated by spaces.', 'ERROR');
    const mine = kind === 'files' ? this.aiderFiles : this.aiderRead;
    const other = kind === 'files' ? this.aiderRead : this.aiderFiles;
    const said = [];
    for (const n of names) {
      const i = other.indexOf(n);
      if (i >= 0) other.splice(i, 1);
      if (mine.includes(n)) { said.push(n + ' is already in the chat'); continue; }
      mine.push(n);
      said.push('added ' + n + (kind === 'files' ? '' : ' (read-only)'));
    }
    this.paintTitle();
    if (this.open && this.tab === 'aider') this.paintPane();
    return this.aiderSay(typed, said.join('\n') + '\n' + this.aiderFilesText());
  },

  aiderDrop(rest, typed) {
    const names = this.aiderSplit(rest);
    if (!names.length) return this.aiderSay(typed, 'Name the file to take out: /drop path/to/file.py, or /drop all.', 'ERROR');
    const had = this.aiderFiles.concat(this.aiderRead);
    const all = names.length === 1 && names[0].toLowerCase() === 'all';
    const gone = all ? had : had.filter(f => names.includes(f));
    this.aiderFiles = this.aiderFiles.filter(f => !gone.includes(f));
    this.aiderRead = this.aiderRead.filter(f => !gone.includes(f));
    const said = [];
    if (gone.length) said.push('dropped ' + gone.join(', '));
    else if (all) said.push('nothing was in the chat');
    const missed = all ? [] : names.filter(n => !had.includes(n));
    if (missed.length) said.push(missed.join(', ') + (missed.length === 1 ? ' was' : ' were') + ' not in the chat');
    this.paintTitle();
    if (this.open && this.tab === 'aider') this.paintPane();
    return this.aiderSay(typed, said.join('\n') + '\n' + this.aiderFilesText());
  },

  // The page answered by itself, not the door; the chip says so.
  aiderSay(command, text, status) {
    this.entries.push({ id: ++this._seq, kind: 'aider', prompt: this.promptText('aider'), command, state: 'note', note: text, bad: status === 'ERROR', t0: Date.now(), ms: 0 });
    this.paintOut();
  },

  // One instruction, sent as typed with the names of the files in the chat. The run takes a minute or more, and the
  // entry says so while it waits; nothing in the world has changed until the door answers.
  async aiderRun(text) {
    const files = this.aiderFiles.slice();
    const read = this.aiderRead.slice();
    const en = { id: ++this._seq, kind: 'aider', prompt: this.promptText('aider'), command: text, state: 'running', t0: Date.now(), files, read };
    this.entries.push(en);
    this.aiderTick(true);
    this.paintOut();
    let said;
    try { said = await App.tool('aider_run', { message: text, files: files.join('\n'), context: read.join('\n') }); }
    catch (e) { en.state = 'error'; en.err = e.message || 'refused'; en.ms = Date.now() - en.t0; this.aiderTick(this.entries.some(x => x.kind === 'aider' && x.state === 'running')); this.paintOut(); return; }
    this.aiderTake(en, said);
  },

  // Taking a run back is a button he presses or a line he types, and nothing else.
  async aiderUndo(run, typed) {
    const id = String(run || '').trim();
    if (!id) return this.aiderSay(typed, 'There is no run in this tab to take back yet. /undo <run id> takes back one from the Recent runs in the Aider Pair inspector.', 'ERROR');
    const en = { id: ++this._seq, kind: 'aider', prompt: this.promptText('aider'), command: typed || '/undo ' + id, state: 'running', t0: Date.now(), undoOf: id };
    this.entries.push(en);
    this.paintOut();
    let said;
    try { said = await App.tool('aider_undo', { run: id }); }
    catch (e) { en.state = 'error'; en.err = e.message || 'refused'; en.ms = Date.now() - en.t0; this.paintOut(); return; }
    this.aiderTake(en, said);
  },

  // What the door answered, in the one shape it answers in (state, why, changed, diff, said, note, guarded, ...).
  aiderTake(en, said) {
    let a = null;
    try { a = JSON.parse(said); } catch { /* reported below */ }
    en.ms = Date.now() - en.t0;
    if (!a || typeof a !== 'object' || !a.state) {
      en.state = 'error';
      en.err = 'The door answered something that is not Aider\'s shape: ' + String(said).slice(0, 300);
    } else {
      en.ans = a;
      en.state = a.state;
      if (en.undoOf && a.state === 'undone') {
        for (const x of this.entries) { if (x.ans && x.ans.run === en.undoOf) x.undone = true; }
        if (this.aiderLast === en.undoOf) this.aiderLast = '';
      } else if (!en.undoOf && a.state === 'ran' && a.run && (a.changed || []).length && !a.withheld) this.aiderLast = a.run;
    }
    this.aiderTick(this.entries.some(x => x.kind === 'aider' && x.state === 'running'));
    this.readAider().then(() => { if (this.open && this.tab === 'aider') this.paintPane(); });
    this.paintOut();
  },

  since(t0) { return Math.floor((Date.now() - t0) / 1000) + 's'; },

  // A run waits for minutes; the clock on its chip keeps moving so a working Aider does not read as a hung one.
  aiderTick(on) {
    clearInterval(this._at);
    this._at = null;
    if (!on) return;
    this._at = setInterval(() => {
      const live = this.entries.filter(e => e.kind === 'aider' && e.state === 'running');
      if (!live.length || !$ag('ag-out')) { this.aiderTick(false); return; }
      for (const e of live) { const el = $ag('ag-at-' + e.id); if (el) el.textContent = this.since(e.t0); }
    }, 1000);
  },

  diffHtml(text) {
    return String(text || '').split('\n').map(l => {
      const e = escHtml(l);
      if (l.startsWith('+++') || l.startsWith('---') || l.startsWith('diff ')) return '<span class="ag-d-file">' + e + '</span>';
      if (l.startsWith('@@')) return '<span class="ag-d-hunk">' + e + '</span>';
      if (l.startsWith('+')) return '<span class="ag-d-add">' + e + '</span>';
      if (l.startsWith('-')) return '<span class="ag-d-del">' + e + '</span>';
      return e;
    }).join('\n');
  },

  // A line typed in the Aider tab, as the door answered it. The page decides nothing: it draws the answer's state,
  // reasons, files, diff, words and notes, and a failure (a refusal, a wall that stopped something, an edit that was
  // withheld, a run that was ended) is on the face of it. The Undo button is drawn for a run that wrote and for nothing else.
  aiderEntryHtml(en) {
    const id = en.id;
    const a = en.ans || null;
    const why = (a && a.why) || [];
    const chips = [];
    let bad = !!en.bad;
    if (en.state === 'running') chips.push('<span class="ag-chip">running <span id="ag-at-' + id + '">' + this.since(en.t0) + '</span></span>');
    else if (en.state === 'note') chips.push('<span class="ag-chip' + (bad ? ' bad' : '') + '">page</span>');
    else if (en.state === 'refused') { chips.push('<span class="ag-chip bad">refused</span>'); bad = true; }
    else if (en.state === 'error') { chips.push('<span class="ag-chip bad">could not run</span>'); bad = true; }
    else if (en.state === 'undone') chips.push('<span class="ag-chip">taken back</span>');
    else if (en.state === 'ran' && a) {
      if (a.withheld) { chips.push('<span class="ag-chip bad">withheld</span>'); bad = true; }
      else if (a.timed_out) { chips.push('<span class="ag-chip bad">timed out</span>'); bad = true; }
      else if (a.exit != null && a.exit !== 0) { chips.push('<span class="ag-chip bad">exit ' + escHtml(String(a.exit)) + '</span>'); bad = true; }
      else if ((a.changed || []).length) chips.push('<span class="ag-chip">written</span>');
      else chips.push('<span class="ag-chip">no change</span>');
      if (a.model) chips.push('<span>' + escHtml(a.model) + '</span>');
      if (a.line) chips.push('<span>on ' + escHtml(a.line) + '</span>');
      if (a.tokens) chips.push('<span>' + escHtml(String(a.tokens).replace(/\.$/, '')) + '</span>');
    }
    // Only a finished answer has a time: a card still waiting has not taken that long to do anything.
    if (en.state === 'ran' || en.state === 'refused' || en.state === 'undone') {
      const ms = a && a.ms != null ? a.ms : en.ms;
      if (ms != null && ms > 0) chips.push('<span>' + ms + 'ms</span>');
    }
    const out = en.state === 'note' ? (en.note || '') : '';
    const lines = [];
    if (en.state === 'running') lines.push('<div class="ag-src">' + (en.undoOf ? 'putting the files back' : 'Aider is working on a copy of ' + escHtml((en.files || []).join(', ') || 'no files') + '. A run takes minutes (the model reads the whole file and writes the edit at the speed of this machine), and nothing in the world changes until it has finished and its edit has been checked.') + '</div>');
    if (en.state === 'refused') lines.push(`<div class="ag-notrun"><b>REFUSED BY NAME</b><br>${why.map(escHtml).join('<br>')}<br><span class="ag-src">no click lifts a refusal; nothing ran and nothing was written</span></div>`);
    if (en.state === 'error') lines.push(`<div class="ag-notrun"><b>COULD NOT RUN</b><br>${escHtml(en.err || '')}</div>`);
    if (a && (a.guarded || []).length) lines.push(`<div class="ag-notrun"><b>THE WALL STOPPED</b><br>${a.guarded.map(escHtml).join('<br>')}</div>`);
    if (a && a.withheld) lines.push(`<div class="ag-notrun"><b>NOT WRITTEN</b><br>${escHtml(a.note || '')}</div>`);
    else if (a && a.timed_out) lines.push(`<div class="ag-notrun"><b>ENDED</b> ${escHtml(a.note || 'it ran past its limit')}</div>`);
    else if (a && a.exit != null && a.exit !== 0) lines.push(`<div class="ag-notrun"><b>FAILED</b> ${escHtml(a.note || '')}</div>`);
    else if (a && a.note) lines.push(`<div class="ag-src">${escHtml(a.note)}</div>`);
    if (a && (a.changed || []).length) lines.push('<div class="ag-src">' + (a.withheld ? 'the edit that was NOT written: ' : '') + a.changed.map(f => escHtml((f.new ? 'new ' : '') + f.file + (f.added || f.removed ? ' +' + f.added + ' -' + f.removed : ''))).join(' &middot; ') + '</div>');
    if (a && a.diff) lines.push('<pre class="ag-pre diff' + (a.withheld ? ' bad' : '') + '">' + this.diffHtml(a.diff) + '</pre>');
    if (en.state === 'ran' && a && a.run && (a.changed || []).length && !a.withheld && !en.undone) {
      lines.push(`<div class="ag-btns"><button type="button" class="ag-btn no" data-act="aider-undo" data-run="${escHtml(a.run)}">Undo this run</button><span class="ag-src">run ${escHtml(a.run)}; unsaved on ${escHtml(a.line || 'the line')} until Version control saves it</span></div>`);
    }
    if (en.undone) lines.push('<div class="ag-src">taken back</div>');
    if (a && a.said) lines.push(`<div><button type="button" class="ag-link" data-act="aider-said" data-id="${id}">${en.showSaid ? 'hide' : 'show'} what Aider said</button></div>` + (en.showSaid ? `<pre class="ag-pre dim">${escHtml(a.said)}</pre>` : ''));
    if (a && (a.ignored || []).length) lines.push(`<div class="ag-src">Aider also touched, and none of it was kept: ${a.ignored.map(escHtml).join('; ')}</div>`);
    if (a && a.truncated) lines.push('<div class="ag-src">what Aider said was cut at 16 KB</div>');
    return `<div class="ag-e" id="ag-e-${id}"><div class="ag-line"><div class="ag-line-l"><span class="ag-p aider">${en.prompt || this.promptText('aider')}</span><span class="ag-cmd">${escHtml(en.command)}</span></div>
      <div class="ag-meta">${chips.join('')}<button type="button" class="ag-copy" data-act="copy" data-id="${id}" title="Copy the output" aria-label="Copy the output">${AG_ICON.copy}</button></div></div>
      ${out ? `<pre class="ag-pre${bad && en.state === 'note' ? ' bad' : ''}" id="ag-o-${id}">${escHtml(out)}</pre>` : ''}${lines.join('')}</div>`;
  },

  // The page answered by itself, not the engine; the chip says so.
  say(command, text, status) {
    this.entries.push({ id: ++this._seq, kind: 'local', command, text, status: status || 'SUCCESS', ms: 0 });
    this.paintOut();
  },

  send(raw) {
    const him = { who: 'him', text: raw };
    const cn = { who: 'council', text: '', live: true };
    Chat.thread.push(him, cn);
    const en = { id: ++this._seq, kind: 'turn', command: raw, him, cn, isCommand: /^\//.test(raw) };
    this.entries.push(en);
    Run.start({ objective: raw });
    if (en.isCommand) this.tapText(en);
    this.paintOut();
  },

  answer(text) {
    if (Run.running) return;
    const him = { who: 'him', text };
    const cn = { who: 'council', text: '', live: true };
    Chat.thread.push(him, cn);
    const en = { id: ++this._seq, kind: 'turn', command: text, him, cn, isCommand: false, isAnswer: true };
    this.entries.push(en);
    Run.start({ answer: text });
    this.tapText(en);
    this.paintOut();
  },

  // A /command prints; it does not deliver. What it printed is on the wire as
  // `text`, which Run keeps out of its events (the console copy of everything
  // else would double a run), so this listens to the same stream for it.
  tapText(en) {
    en.cmdText = '';
    const es = Run.es;
    if (!es) return;
    es.addEventListener('engine', (e) => {
      let d; try { d = JSON.parse(e.data); } catch { return; }
      if (d.event === 'text') { en.cmdText += d.text || ''; this.paintLive(); }
    });
  },

  // Speak, and the words land IN THE LINE for him to read and fix. Nothing is
  // sent on his behalf.
  mic() {
    const btn = $ag('ag-mic');
    const input = $ag('ag-input');
    if (Run.listening) { Run.stopHearing(); btn.classList.remove('hot'); return; }
    if (!Run.engineOpen) { toast('No engine is open on this world', 'error'); return; }
    btn.classList.add('hot');
    Run.hear(
      () => {},
      (text) => {
        btn.classList.remove('hot');
        input.value = (input.value ? input.value + ' ' : '') + text;
        $ag('ag-go').disabled = !input.value.trim();
        input.focus();
        input.setSelectionRange(input.value.length, input.value.length);
      },
      (why) => { btn.classList.remove('hot'); toast(why, 'error'); });
  },

  // ---- the sitting ------------------------------------------------------------

  // One /command through the council stream, its printed text handed to
  // onText as it arrives and returned whole. Counted in Run.own so this tab's
  // own echo on the bus is not drawn as a turn somebody else started.
  stream(command, onText) {
    return new Promise((done) => {
      let got = '';
      Run.own++;
      let over = false;
      let es = null;
      const end = () => {
        if (over) return;
        over = true;
        try { if (es) es.close(); } catch { /* already closed */ }
        setTimeout(() => { Run.own = Math.max(0, Run.own - 1); }, Run.ECHO_GRACE_MS);
        done(got);
      };
      const add = (s) => { got += s; onText(s); };
      es = new EventSource(API.base + '/council/stream?' + new URLSearchParams({ objective: command }));
      es.addEventListener('engine', (e) => {
        let d; try { d = JSON.parse(e.data); } catch { return; }
        if (d.event === 'text') add(d.text || '');
        if (d.event === 'error') add('\n' + (d.text || '') + '\n');
      });
      es.addEventListener('stream_end', end);
      es.addEventListener('stream_error', (e) => {
        let d = {}; try { d = JSON.parse(e.data); } catch { /* no body */ }
        add('\nREFUSED: ' + (d.error || 'the step was refused') + '\n');
        end();
      });
      es.onerror = end;
    });
  },

  local(command) {
    const en = { id: ++this._seq, kind: 'local', command, text: '', status: 'RUNNING', t0: Date.now(), ms: null };
    this.entries.push(en);
    this.paintOut();
    return en;
  },

  paintLocal(en) {
    const pre = $ag('ag-o-' + en.id);
    if (!pre) { this.paintOut(); return; }
    pre.textContent = en.text;
    this.stick();
  },

  // The boot, in the order the Dashboard does it: close an open sitting, open
  // a fresh engine, load the pipeline's models, then the boot report. A step
  // that fails says so while the rest still runs.
  async boot() {
    if (this.booting) return;
    this.booting = true;
    const en = this.local('/boot');
    const add = (s) => { en.text += s; this.paintLocal(en); };
    try {
      if (Run.engineOpen) {
        add('closing the open sitting first...\n');
        add((await App.tool('env_close', {})) + '\n');
      }
      add('opening a fresh engine...\n');
      add((await App.tool('env_open', {})) + '\n');
    } catch (e) {
      add('REFUSED: ' + (e.message || 'refused') + '\n');
      en.status = 'ERROR';
      en.ms = Date.now() - en.t0;
      await Run.check();
      this.booting = false;
      this.paintTitle(); this.paintTop(); this.paintOut();
      return;
    }
    await Run.check();
    this.paintTitle(); this.paintTop();
    add('\nloading this pipeline\'s models (/warm)...\n');
    await this.stream('/warm', add);
    add('\nthe boot report (/status)...\n');
    const report = await this.stream('/status', add);
    this.keepBoot(report);
    en.status = 'SUCCESS';
    en.ms = Date.now() - en.t0;
    await Run.check();
    this.booting = false;
    this.paintTitle(); this.paintTop(); this.paintOut();
    this.readStatus();
  },

  // Closing is not a formality: it writes `ended`, and pays the toll if a turn
  // ran. A sitting left open is what makes the next open refuse.
  async closeSitting() {
    const en = this.local('/close');
    const add = (s) => { en.text += s; this.paintLocal(en); };
    try {
      add((await App.tool('env_close', {})) + '\n');
      en.status = 'SUCCESS';
    } catch (e) {
      add('REFUSED: ' + (e.message || 'refused') + '\n');
      en.status = 'ERROR';
    }
    en.ms = Date.now() - en.t0;
    await Run.check();
    this.paintTitle(); this.paintTop(); this.paintOut();
  },

  // The boot report is kept where both browsers can read it, under the key the
  // Dashboard already uses, keyed to the session so a dead engine's report is
  // never shown over a live one.
  keepBoot(text) {
    if (!text || !text.trim()) return;
    this._report = text;
    API.setSetting('boot.' + (Run.world || 'research'), JSON.stringify({ session: Run.session || '', text })).catch(() => {});
  },

  // ---- reads --------------------------------------------------------------

  async readStatus() {
    try { this.rack = this.parseRack(await App.tool('rack_list', {}, true)); }
    catch (e) { this.rack = { ok: false, says: e.message || 'unreadable', tiers: [], n: 0 }; }
    this.paintTop();
    await this.readHolds();
  },

  async readHolds() {
    try {
      const d = JSON.parse(await App.tool('hold_list', {}, true));
      this.holdState = { armed: !!d.armed, n: (d.held || []).length, why: d.why_not || '' };
      if (d.armed) this.settleCards(new Set((d.held || []).map(h => h.id)));
    } catch { this.holdState = null; }
    this.paintTop();
  },

  // A shell card whose hold has left the door's queue was answered somewhere else (the Guardrails
  // tab) or dropped when the door restarted. It stops offering buttons that can only be refused.
  settleCards(waiting) {
    let moved = false;
    for (const en of this.entries) {
      if (en.kind === 'shell' && en.state === 'held' && !en.deciding && en.ans && !waiting.has(en.ans.hold)) { en.state = 'gone'; moved = true; }
    }
    if (moved) this.paintOut(false);
  },

  // The rack's own ladder, parsed the way the old Dashboard parsed it. A voice
  // is counted only if the tool listed it; silence is reported as silence.
  parseRack(text) {
    const t = String(text || '');
    if (/nothing fabricated/.test(t)) return { ok: false, says: t.split('\n')[0], tiers: [], n: 0 };
    const tiers = [];
    let cur = null;
    for (const line of t.split('\n')) {
      let m = /^\s{2}(\S+) \(([^)]+)\):/.exec(line);
      if (m) { cur = { name: m[1], cap: m[2], voices: [] }; tiers.push(cur); continue; }
      m = /^\s{4}- (\S+) · ([0-9.]+)GB · (\S+)/.exec(line);
      if (m && cur) cur.voices.push({ tag: m[1], gb: parseFloat(m[2]), family: m[3] });
    }
    const n = tiers.reduce((a, x) => a + x.voices.length, 0);
    return { ok: true, says: '', tiers, n };
  },

  // The check is repeated while the tab is on screen and not busy, and again
  // the moment he returns to it: an engine closes itself after thirty idle
  // minutes and the title must not go on saying it is open.
  watch(on) {
    clearInterval(this._w);
    this._w = null;
    if (this._wv) { document.removeEventListener('visibilitychange', this._wv); this._wv = null; }
    if (!on) return;
    const again = () => {
      if (!$ag('ag-out')) { this.watch(false); return; }
      if (document.hidden || Run.running || this.booting) return;
      Run.check().then(() => { this.paintTitle(); this.paintTop(); });
      if (this.isAider()) this.readAider();
    };
    this._wv = () => { if (!document.hidden) again(); };
    document.addEventListener('visibilitychange', this._wv);
    this._w = setInterval(again, 30000);
  },

  // ---- the Inspector --------------------------------------------------------

  door(id) {
    this.open = true;
    this.tab = id;
    this.paintPanel();
  },

  paintPanel() {
    const p = $ag('ag-panel');
    if (!p) return;
    p.hidden = !this.open;
    $ag('ag-insp').classList.toggle('on', this.open);
    document.querySelectorAll('#ag-nav [data-door]').forEach(b => b.classList.toggle('on', this.open && b.dataset.door === this.tab));
    if (!this.open) return;
    const tab = this.TABS.find(t => t.id === this.tab) || this.TABS[0];
    $ag('ag-ph-tab').textContent = tab.label;
    document.querySelectorAll('#ag-tabs [data-tab]').forEach(b => b.classList.toggle('on', b.dataset.tab === this.tab));
    this.paintPane();
  },

  // Each tab reads its own tool and draws what came back. A stamp keeps a slow
  // read from landing on the tab that replaced it.
  async paintPane() {
    const box = $ag('ag-pane');
    if (!box || !this.open) return;
    const id = this.tab;
    const stamp = ++this._pane;
    const put = (html) => { if (this._pane === stamp && this.tab === id && $ag('ag-pane')) $ag('ag-pane').innerHTML = html; };
    if (id === 'run') { put(this.runHtml()); return; }
    put('<div class="ag-empty">Reading...</div>');
    try {
      if (id === 'aider') put(await this.aiderHtml());
      else if (id === 'ledger') put(await this.ledgerHtml());
      else if (id === 'flows') put(await this.flowsHtml());
      else if (id === 'guards') put(await this.guardsHtml());
      else if (id === 'docs') put(await this.docsHtml());
      else if (id === 'rack') put(await this.rackHtml());
    } catch (e) {
      put(`<div class="ag-card"><h4>Could not be read</h4><p>${escHtml(e.message || 'refused')}</p></div>`);
    }
  },

  card(title, body, src) {
    return `<div class="ag-card"><h4>${title}</h4>${body}${src ? `<div class="ag-src">${escHtml(src)}</div>` : ''}</div>`;
  },

  // ---- Run: the last turn, read off the events the engine sent ---------------

  stages(t) {
    const ev = t.events || [];
    const kinds = new Set(ev.map(e => e._kind));
    const refused = kinds.has('refused');
    const seats = t.seats || [];
    const tools = t.tools || [];
    const failed = tools.filter(x => x.failed).length;
    const delivered = !!t.delivery;
    const idle = !Run.running;
    const ran = seats.length > 0 || tools.length > 0 || delivered;
    // A seat sits again for every tool round, so the sittings are not the seats:
    // the delivery's own steps are, and until it arrives the distinct names.
    const nSeats = ((t.delivery && t.delivery.steps) || []).length || new Set(seats.map(s => s.seat)).size;
    return [
      { n: 'Asked', s: 'done', v: 'in' },
      { n: 'Law gate', s: refused ? 'bad' : ran ? 'done' : (Run.running ? 'act' : ''), v: refused ? 'refused' : ran ? 'passed' : '' },
      { n: 'Seats', s: seats.length ? (Run.running && !delivered ? 'act' : 'done') : '', v: nSeats ? String(nSeats) : '' },
      { n: 'Tools', s: tools.length ? (failed ? 'bad' : tools.some(x => !x.done) ? 'act' : 'done') : (idle && ran ? 'skip' : ''),
        v: tools.length ? tools.length + (tools.length === 1 ? ' call' : ' calls') + (failed ? ', ' + failed + ' failed' : '') : (idle && ran ? 'none' : '') },
      { n: 'Needs you', s: t.waiting ? 'wait' : (delivered ? 'skip' : ''), v: t.waiting ? 'asking' : (delivered ? 'no' : '') },
      { n: 'Delivered', s: delivered ? 'done' : (t.refusal ? 'bad' : ''), v: delivered ? Run.elapsed(t) : (idle && t.verdict ? t.verdict : '') }
    ];
  },

  evRow(d) {
    const k = d._kind;
    const e = escHtml;
    switch (k) {
      case 'opened': return `<div class="ag-ev"><b>opened</b> sitting ${e(String(d.sitting ?? ''))} &middot; session ${e(d.session || '')}</div>`;
      case 'run': return `<div class="ag-ev"><b>run</b> pipeline ${e(d.pipeline || '')}${d.review_only ? ' &middot; review only' : ''}${d.transcript ? '<br><span class="ag-src">transcript ' + e(d.transcript) + '</span>' : ''}</div>`;
      case 'seat': return `<div class="ag-ev"><b>${e(d.seat || 'seat')}</b> ${e(d.model || '')}${d.timeout ? ' &middot; timeout ' + e(String(d.timeout)) + 's' : ''}</div>`;
      case 'token': return '';
      case 'text': return '';
      case 'tool': return `<div class="ag-ev"><b>skill</b> ${e(d.action || d.tool || d.name || '?')}${d.seat ? ' <span class="ag-src">called by ' + e(d.seat) + '</span>' : ''}</div>`;
      case 'tool_result': return `<div class="ag-ev ${d.failed ? 'bad' : 'ok'}"><b>${d.failed ? 'FAILED' : 'ok'}</b> ${e(d.action || d.tool || d.name || '?')}${d.error ? ' - ' + e(d.error) : ''}${d.failed && d.text ? ' - ' + e(String(d.text).slice(0, 240)) : ''}</div>`;
      case 'report': case 'note': return `<div class="ag-ev"><b>${k}</b> ${e((d.text || '').trim().slice(0, 300))}</div>`;
      case 'needs_answer': return `<div class="ag-ev gate"><b>THE COUNCIL IS ASKING</b> ${e(d.prompt || '')}</div>`;
      case 'delivery': return `<div class="ag-ev ok"><b>DELIVERY</b> ${e(d.pipeline || '')}${d.elapsed != null ? ' &middot; ' + e(String(d.elapsed)) + 's' : ''}</div>`;
      case 'refused': case 'aborted': case 'cancelled': case 'unreachable': case 'error':
        return `<div class="ag-ev bad"><b>${e(String(k).toUpperCase())}</b> ${e(d.text || d.error || '')}</div>`;
      case 'closed': return `<div class="ag-ev"><b>closed</b> ${e(d.text || 'the sitting is tolled')}</div>`;
      case 'command': return `<div class="ag-ev"><b>command</b> ${e(d.text || '')} <span class="ag-src">finished without running a pipeline</span></div>`;
      default: return `<div class="ag-ev"><b>${e(String(k))}</b> ${e(JSON.stringify(d).slice(0, 300))}</div>`;
    }
  },

  runHtml() {
    const t = Run.turn;
    if (!t) return '<div class="ag-empty">No turn in this tab yet. Say something to the council and the whole turn lands here.</div>';
    const st = this.stages(t).map(s => `<div class="ag-st ${s.s}"><b>${s.n}</b><span>${escHtml(s.v)}</span></div>`).join('');
    const fails = Run.failures(t);
    const steps = (t.delivery && t.delivery.steps) || [];
    const stepRows = steps.map(s => `<tr><td>${escHtml(String(s.seat || ''))}</td><td>${escHtml(String(s.model || ''))}</td>
      <td>${escHtml(String(s.elapsed ?? ''))}s</td><td>${s.skipped ? 'skipped' : s.error ? 'error' : 'ran'}</td></tr>`).join('');
    return this.card(escHtml(t.objective || 'the last turn') + (Run.running ? '' : ' <span class="ag-chip">' + escHtml(t.verdict || 'done') + '</span>'),
        `<div class="ag-stages">${st}</div>` +
        (fails.length ? `<div class="ag-notrun"><b>NOT EVERYTHING RAN</b><br>${fails.map(escHtml).join('<br>')}</div>` : ''),
        'read off the events the engine sent; nothing here is inferred from a seat\'s own words') +
      (stepRows ? this.card('Seats', `<table><thead><tr><th>seat</th><th>model</th><th>elapsed</th><th></th></tr></thead><tbody>${stepRows}</tbody></table>`, 'the delivery\'s own per-seat facts') : '') +
      this.card('Events', (t.events || []).map(d => this.evRow(d)).join('') || '<div class="ag-empty">No events yet.</div>',
        (t.events || []).length + ' events kept' + (t.thinned ? ' - this turn came back from storage without its record' : ''));
  },

  // ---- Audit Ledger -----------------------------------------------------------

  async ledgerHtml() {
    const w = Run.world || 'research';
    let kept = this._report || '';
    if (!kept) {
      try {
        const r = await API.getSetting('boot.' + w);
        const k = JSON.parse((r && r.value) || '{}');
        if (k.text && (!k.session || k.session === Run.session) && Run.engineOpen) kept = k.text;
      } catch { /* none kept */ }
    }
    let p = null;
    try { p = JSON.parse(await App.tool('proofs', {}, true)); } catch { /* shown as unread below */ }
    const rc = (p && p.record) || null;
    const rows = rc && !rc.error ? (rc.recent || []).slice().reverse().map(r => `<tr>
      <td>${escHtml(String(r.n))}</td><td>${escHtml(String(r.started || '').replace('T', ' ').slice(0, 16))}</td>
      <td>${r.ended ? escHtml(String(r.ended).slice(11, 16)) : '<span style="color:var(--ag-rose)">open</span>'}</td>
      <td>${escHtml(String(r.runs))}</td><td>${r.toll_paid ? 'tolled' : 'no toll'}</td></tr>`).join('') : '';
    const report = this.card('The engine\'s own report',
      kept ? `<pre>${escHtml(kept)}</pre>` : '<p>No report is kept for this sitting.</p>',
      kept ? 'what /status printed at boot, kept' : 'press the button to read it now') +
      `<div class="ag-btns"><button type="button" class="ag-btn" data-act="status"${Run.engineOpen ? '' : ' disabled'}>Read /status now</button></div>`;
    const sit = this.card('Sittings',
      rc && !rc.error
        ? `<div class="ag-row"><span class="k">sittings</span><span class="v">${rc.sittings}</span></div>
           <div class="ag-row"><span class="k">still open</span><span class="v">${rc.still_open || 0}</span></div>
           <div class="ag-row"><span class="k">tolled</span><span class="v">${rc.tolled}</span></div>
           <div class="ag-row"><span class="k">runs recorded</span><span class="v">${rc.runs}</span></div>
           ${rows ? `<table><thead><tr><th>#</th><th>opened</th><th>closed</th><th>runs</th><th>toll</th></tr></thead><tbody>${rows}</tbody></table>` : ''}`
        : `<p>The record could not be read${rc && rc.error ? ': ' + escHtml(rc.error) : ''}.</p>`,
      'proofs, read from sessions/sessions.jsonl and never counted here');
    const logs = this.card('Transcripts', '<div id="ag-logs"><p>Every run is written down in logs/ with a sha256 receipt.</p></div>' +
      '<div class="ag-btns"><button type="button" class="ag-btn" data-act="logs">List the newest</button></div>', 'records');
    const note = this.card('Chain verdicts',
      '<p>The door\'s own chain checker cannot read the pen\'s links yet, so it would call a whole chain FLIP (WHAT\'S LEFT H7). Until it can, the words to read are the engine\'s own, in the report above: its record, gate and memory lines say whether the law and the memory are whole.</p>',
      'verify_chain is not asked here on purpose');
    return report + sit + logs + note;
  },

  // ---- Aider Pair ---------------------------------------------------------------

  // The Inspector's side of the tab: what the door says of Aider now, the files in the chat (each can be dropped), the
  // runs this world kept (each can be taken back), and where each world's work stands -- its branch and what changed --
  // because an edit is unsaved work on a line, and Version control is where it is saved.
  async aiderHtml() {
    await this.readAider();
    const s = this.aiderSt || {};
    const parts = [];
    if (s.error) {
      parts.push(this.card('Aider\'s state could not be read', `<p>${escHtml(s.error)}</p>`, 'aider_status'));
    } else {
      const rows = [
        ['installed', s.installed ? 'yes' + (s.version ? ', ' + s.version : '') : 'no'],
        ['folder', s.folder || ''],
        ['model', s.model ? s.model + ' (the coding seat\'s)' : 'none declared'],
        ['window', s.context ? s.context + ' tokens, about ' + (s.budget_kb || '?') + ' KB of files' : 'unknown'],
        ['rack', (s.rack === 'up' ? 'up' : 'silent') + (s.model ? (s.on_rack ? ', has the model' : ', does not have the model') : '')],
        ['line of work', s.on_line ? s.line : (s.line ? s.line + ' (the main line is yours)' : 'none')],
        ['busy', s.busy ? 'yes: a run is going' : 'no']
      ];
      parts.push(this.card('Aider ' + (s.installed ? 'is installed' : 'is not installed here'),
        rows.map(([k, v]) => `<div class="ag-row"><span class="k">${escHtml(k)}</span><span class="v">${escHtml(String(v))}</span></div>`).join('') +
        ((s.why || []).length ? `<div class="ag-notrun"><b>NOT READY</b><br>${s.why.map(escHtml).join('<br>')}</div>` : ''),
        'aider_status, read just now'));
    }
    const fileRows = (list, ro) => list.map(f => `<div class="ag-row"><span class="k">${escHtml(f)}${ro ? ' <span class="ag-src">read-only</span>' : ''}</span><span class="v"><button type="button" class="ag-link" data-act="aider-drop" data-file="${escHtml(f)}">drop</button></span></div>`).join('');
    parts.push(this.card('In the chat',
      (fileRows(this.aiderFiles, false) + fileRows(this.aiderRead, true)) ||
        '<p>No files yet. In the Aider tab, /add path/to/file.py names a file Aider may change, and /read path one it may only read.</p>',
      'kept in this page; the door judges every name when a run is asked'));
    const runs = (s.runs || []).map(r => `<div class="ag-row"><span class="k">${escHtml(r.run)}<br><span class="ag-src">${escHtml((r.files || []).join(', '))}${r.line ? ' &middot; ' + escHtml(r.line) : ''}</span></span>
      <span class="v">${r.undone ? 'taken back' : `<button type="button" class="ag-btn no" data-act="aider-undo" data-run="${escHtml(r.run)}">Undo</button>`}</span></div>`).join('');
    parts.push(this.card('Recent runs', runs || '<p>No run yet.</p>', 'the last runs this world kept; an older one is git\'s to take back'));
    let worlds = [];
    let said = '';
    try {
      const m = await App.tool('muster', {}, true);
      worlds = String(m).split('\n').map(x => x.trim()).filter(x => x && !x.endsWith(':'));
    } catch (e) { said = e.message || 'unreadable'; }
    for (const wd of worlds) {
      let d = null;
      try { d = JSON.parse(await App.tool('git', { project: wd })); } catch (e) { d = { error: e.message || 'unreadable' }; }
      const keys = ['branch', 'changed', 'untracked', 'ahead', 'behind'].filter(k => d && d[k] != null && d[k] !== '');
      parts.push(this.card(escHtml(wd),
        d && d.error ? `<p>${escHtml(d.error)}</p>`
          : d && d.is_repo === false ? '<p>Not a repository.</p>'
          : keys.map(k => `<div class="ag-row"><span class="k">${escHtml(k)}</span><span class="v">${escHtml(String(d[k]))}</span></div>`).join('') || '<p>Nothing to report.</p>',
        'git') + `<div class="ag-btns"><button type="button" class="ag-btn" data-act="diff" data-world="${escHtml(wd)}">Show the diff</button></div>`);
    }
    if (said) parts.push(this.card('The worlds could not be read', `<p>${escHtml(said)}</p>`, 'muster'));
    parts.push('<div id="ag-diff"></div>');
    parts.push(`<div class="ag-btns"><button type="button" class="ag-btn" data-act="goto" data-path="/flows">Version control</button></div>`);
    return parts.join('');
  },

  // ---- Workflows ------------------------------------------------------------------

  async flowsHtml() {
    let flows = '', runs = '';
    try { flows = (await API.listFlows()).flows || ''; } catch (e) { flows = 'The flows could not be read: ' + (e.message || 'refused'); }
    try { runs = (await API.listRuns('')).runs || ''; } catch (e) { runs = 'The runs could not be read: ' + (e.message || 'refused'); }
    const cut = (s, n) => s.length > n ? '...' + s.slice(-n) : s;
    return this.card('Flows', `<pre>${escHtml(flows.trim() || 'none declared')}</pre>`, 'flow_list') +
      this.card('Recent runs', `<pre>${escHtml(cut(runs.trim(), 3000) || 'none on record')}</pre>`, 'flow_runs') +
      `<div class="ag-btns"><button type="button" class="ag-btn" data-act="goto" data-path="/workflows">Open the builder</button>
        <button type="button" class="ag-btn" data-act="goto" data-path="/flows">Version control</button></div>`;
  },

  // ---- Guardrails -------------------------------------------------------------------

  async guardsHtml() {
    let d = null, err = '';
    try { d = JSON.parse(await App.tool('hold_list', {}, true)); } catch (e) { err = e.message || 'unreadable'; }
    const parts = [];
    const asking = this.asking();
    if (asking) parts.push(this.card('The council is asking',
      `<p>${escHtml(asking)}</p><p>Answer it in the terminal; the question is waiting there.</p>`, 'the last run'));
    if (err) {
      parts.push(this.card('Waiting for your hand', `<p>The hold queue could not be read: ${escHtml(err)}</p>`, 'hold_list'));
    } else if (!d.armed) {
      parts.push(this.card('Waiting for your hand',
        `<p>Nothing is being held, and nothing can be. ${escHtml(d.why_not || 'The door was started without --auth, so it cannot tell a seat from this page.')} RULE 6 is a convention again until the door is restarted with --auth and a service wire.</p>`, 'hold_list'));
    } else if (!(d.held || []).length) {
      parts.push(this.card('Waiting for your hand', '<p>Nothing is waiting. A writing call from anything but this page parks here for your decision.</p>', 'hold_list'));
    } else {
      parts.push(this.card('Waiting for your hand',
        `<p>${d.held.length} writing call${d.held.length === 1 ? '' : 's'} parked at the door. Approving runs exactly the call shown, with the arguments it was parked with.</p>` +
        d.held.map(h => `<div class="ag-card"><h4>${escHtml(h.tool)}</h4>
          <p>asked by ${escHtml(h.caller)} &middot; ${escHtml(h.project || '')} &middot; ${escHtml(h.when || '')}</p>
          <pre>${escHtml(JSON.stringify(h.args || {}, null, 1))}</pre>
          <div class="ag-btns"><button type="button" class="ag-btn go" data-act="hold" data-decision="approve" data-id="${escHtml(h.id)}">Approve - run it</button>
          <button type="button" class="ag-btn no" data-act="hold" data-decision="deny" data-id="${escHtml(h.id)}">Deny</button></div></div>`).join('') +
        '<div id="ag-holdout"></div>', 'hold_list'));
    }
    return parts.join('');
  },

  // ---- Registry & Docs ----------------------------------------------------------------

  async docsHtml() {
    let tools = [];
    let said = '';
    try { tools = (await API.tools()).tools || []; } catch (e) { said = e.message || 'unreadable'; }
    this._tools = tools;
    return this.card('The tool surface', said ? `<p>${escHtml(said)}</p>` :
        `<input type="text" id="ag-tools-q" placeholder="Filter ${tools.length} tools..." autocomplete="off" /><div id="ag-tools">${this.toolRows(tools.slice(0, 80), tools.length)}</div>`,
        'tools/list - the door\'s own registry, read just now') +
      this.card('Documents', '<div id="ag-docs"><p>The estate\'s own documents, sorted by what they are, each served whole with a sha256 receipt.</p></div>' +
        '<div class="ag-btns"><button type="button" class="ag-btn" data-act="kinds">List the kinds</button></div>', 'records');
  },

  toolRows(list, total) {
    const rows = list.map(t => `<div class="ag-row"><span class="k">${escHtml(t.name)}</span><span class="v" title="${escHtml(t.description || '')}">${escHtml(Object.keys((t.inputSchema && t.inputSchema.properties) || {}).join(', ') || 'no arguments')}</span></div>`).join('');
    return (rows || '<div class="ag-empty">Nothing matches.</div>') + (list.length < total ? `<div class="ag-src">showing ${list.length} of ${total}; type to narrow</div>` : '');
  },

  filterTools(q) {
    const s = String(q || '').trim().toLowerCase();
    const all = this._tools || [];
    const hit = s ? all.filter(t => (t.name + ' ' + (t.description || '')).toLowerCase().includes(s)) : all;
    const box = $ag('ag-tools');
    if (box) box.innerHTML = this.toolRows(hit.slice(0, 80), hit.length);
  },

  // ---- Rack ----------------------------------------------------------------------------

  async rackHtml() {
    await this.readStatus();
    const r = this.rack;
    if (!r || !r.ok) return this.card('The rack is silent', `<p>${escHtml((r && r.says) || 'no answer')}</p>`, 'rack_list - nothing is fabricated');
    const biggest = Math.max(1, ...r.tiers.flatMap(t => t.voices.map(v => v.gb)));
    return this.card(r.n + ' local voices across ' + r.tiers.length + ' tier' + (r.tiers.length === 1 ? '' : 's'),
      r.tiers.map(t => `<div class="ag-src" style="margin-top:6px">${escHtml(t.name)} ${escHtml(t.cap)}</div>` +
        t.voices.map(v => `<div class="ag-voice"><span>${escHtml(v.tag)}</span><i style="width:${Math.max(4, (v.gb / biggest) * 90)}px"></i><em>${v.gb.toFixed(1)}GB</em></div>`).join('')).join(''),
      'rack_list - loopback only');
  },

  // ---- the Inspector's buttons ------------------------------------------------------------

  async paneClick(e) {
    const b = e.target.closest('[data-act]');
    if (!b) return;
    const act = b.dataset.act;
    if (act === 'goto') { history.pushState(null, '', b.dataset.path); App.router(); return; }
    if (act === 'aider-undo') { this.aiderUndo(b.dataset.run, '/undo ' + b.dataset.run); return; }
    if (act === 'aider-drop') { this.aiderDrop(b.dataset.file, '/drop ' + b.dataset.file); return; }
    if (act === 'status') {
      b.disabled = true;
      const t = await this.stream('/status', () => {});
      this.keepBoot(t);
      this.paintPane();
      return;
    }
    if (act === 'logs') {
      const box = $ag('ag-logs');
      box.innerHTML = '<p>Reading...</p>';
      try {
        const d = JSON.parse(await App.tool('records', { kind: 'logs' }));
        const docs = d.documents || [];
        box.innerHTML = docs.length
          ? `<table><thead><tr><th>run</th><th>when</th></tr></thead><tbody>${docs.slice(0, 40).map(x =>
              `<tr class="click" data-act="doc" data-name="${escHtml(x.name)}"><td>${escHtml(Chat.runName(x.name))}</td><td>${escHtml(when(Date.parse(x.modified) || 0))}</td></tr>`).join('')}</tbody></table><div id="ag-doc"></div>`
          : '<p>No run has been written down yet.</p>';
      } catch (er) { box.innerHTML = `<p>${escHtml(er.message || 'refused')}</p>`; }
      return;
    }
    if (act === 'kinds') {
      const box = $ag('ag-docs');
      box.innerHTML = '<p>Reading...</p>';
      try {
        const d = JSON.parse(await App.tool('records', {}));
        box.innerHTML = (d.kinds || []).map(k => `<div class="ag-row click" style="cursor:pointer" data-act="kind" data-kind="${escHtml(k.kind)}"><span class="k">${escHtml(k.kind)}</span><span class="v">${k.count}</span></div>`).join('') + '<div id="ag-doc"></div>';
        this._recs = d;
      } catch (er) { box.innerHTML = `<p>${escHtml(er.message || 'refused')}</p>`; }
      return;
    }
    if (act === 'kind') {
      const k = ((this._recs && this._recs.kinds) || []).find(x => x.kind === b.dataset.kind);
      const box = $ag('ag-doc');
      if (!k || !box) return;
      box.innerHTML = `<table><thead><tr><th>${escHtml(k.kind)}</th><th>size</th></tr></thead><tbody>${k.documents.slice(0, 60).map(x =>
        `<tr class="click" data-act="doc" data-name="${escHtml(x.name)}"><td>${escHtml(x.name)}</td><td>${(x.bytes / 1024).toFixed(1)} KB</td></tr>`).join('')}</tbody></table><div id="ag-doc2"></div>`;
      return;
    }
    if (act === 'doc') {
      const box = $ag('ag-doc2') || $ag('ag-doc');
      if (!box) return;
      box.innerHTML = '<p>Reading...</p>';
      try {
        const d = JSON.parse(await App.tool('records', { name: b.dataset.name }));
        box.innerHTML = `<div class="ag-src">${escHtml(d.name)} &middot; ${d.bytes} bytes &middot; sha256 ${escHtml(String(d.sha256 || '').slice(0, 16))}${d.sealed ? ' &middot; sealed' : ''}</div><pre>${escHtml(d.text || '')}</pre>`;
      } catch (er) { box.innerHTML = `<p>${escHtml(er.message || 'refused')}</p>`; }
      return;
    }
    if (act === 'diff') {
      const box = $ag('ag-diff');
      box.innerHTML = this.card('Diff - ' + escHtml(b.dataset.world), '<p>Reading...</p>', 'git_diff');
      try {
        const out = await App.tool('git_diff', { project: b.dataset.world });
        box.innerHTML = this.card('Diff - ' + escHtml(b.dataset.world), `<pre>${escHtml(String(out).slice(0, 12000) || 'No uncommitted changes.')}</pre>`, 'git_diff');
      } catch (er) { box.innerHTML = this.card('Diff - ' + escHtml(b.dataset.world), `<p>${escHtml(er.message || 'refused')}</p>`, 'git_diff'); }
      return;
    }
    if (act === 'hold') {
      const out = $ag('ag-holdout');
      b.disabled = true;
      if (out) out.innerHTML = `<pre>${b.dataset.decision === 'approve' ? 'Running it...' : 'Denying...'}</pre>`;
      let said;
      try { said = await App.tool('hold_answer', { id: b.dataset.id, decision: b.dataset.decision }); }
      catch (er) { said = 'Refused: ' + (er.message || 'refused'); }
      await this.readHolds();
      await this.paintPane();
      const o2 = $ag('ag-holdout');
      if (o2) o2.innerHTML = `<pre>${escHtml(said)}</pre>`;
      else if (this.tab === 'guards') { const pane = $ag('ag-pane'); if (pane) pane.insertAdjacentHTML('beforeend', `<pre>${escHtml(said)}</pre>`); }
    }
  }
};
