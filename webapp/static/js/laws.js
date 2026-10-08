// THE LAWS -- what every hand and every agent reads first, and the one place he sets one.
//
// His words, 2026-10-05: "set it as a law all the agents read ... THROUGH THE SYSTEM", then "a LAWS page! great idea, Add it."
//
// WHAT IT SHOWS, AND WHERE IT READS IT. The laws are documents -- CLAUDE.md (the standing rules) and law/*.md (the sitting laws, the estate laws, the
// founding laws and the ledger) -- and the door's `records` tool serves them, as it serves WHATS_LEFT.md to the What's-left page. HOW FAR EACH ONE IS
// SEALED is NOT read from `records`: it marks the whole law/ folder "sealed", which stops being true the day the ledger grows a draft below its seal.
// It is the door's `law_status`: the chain's own links beside law.py's own walk, whose verdict is shown whole and is never this page's opinion.
//
// WHAT IT DOES, which is two things and both are his hand's: set a law (one numbered entry appended to the ledger, below its seal, as a DRAFT every hand
// already reads) and seal the ledger (what is below the seal becomes law and is frozen). Both are ServiceOnly on the door and neither runs while a sitting
// is open (SITTING LAW 5); the page says so and asks nothing it knows would be refused. Sealing is a second, separate click, behind an inline question
// that says what it does -- this console opens no browser pop-up.
//
// THE STYLES RIDE WITH THE PAGE, as left.js's do.
const Laws = {
  RULES: 'CLAUDE.md',
  LEDGER: 'law/LAW_LEDGER.md',

  // THE SAME LAWS IN TWO PLACES (2026-10-06, his word: "go, next: guardrails"; the tab is Laws since 2026-10-08, his card:
  // "It becomes Laws"). The front page's Laws tab draws THIS object -- `inPanel`, without the page's title -- and the page
  // under Pages draws it whole; Read it again, Set and Seal draw it again where it stands.
  async render(el, inPanel) {
    this.inPanel = !!inPanel;
    this.style();
    el.innerHTML = '<div class="loading">Reading the laws...</div>';
    let st;
    try {
      st = JSON.parse(await App.tool('law_status', {}));
    } catch (e) {
      el.innerHTML = `
        <div class="page-header"><div>
          <div class="page-title">Laws</div>
          <div class="page-subtitle">The chain could not be read.</div>
        </div></div>
        <div class="card"><div class="eng-row eng-bad">
          the laws could not be read: ${escHtml(e.message || 'refused')}
          <span class="brief-src">law_status</span></div></div>`;
      return;
    }
    const read = async (name) => {
      try { return JSON.parse(await App.tool('records', { name })); }
      catch (e) { return { name, error: e.message || 'refused' }; }
    };
    const [rules, ledger] = await Promise.all([read(this.RULES), read(this.LEDGER)]);
    this.paint(el, st, rules, ledger);
  },

  // The standing rules: `## RULE 12 — TITLE` and everything under it, to the next `## ` heading or the `---` that ends a section. Shown in number order,
  // not the order the file happens to hold them in (RULE 9 stands before RULE 1 there, for a reason of its own).
  parseRules(text) {
    const out = [];
    const lines = String(text || '').replace(/\r\n?/g, '\n').split('\n');
    let cur = null;
    for (const line of lines) {
      const m = /^## (RULE \d+[a-z]?) [—-] (.+)$/.exec(line);
      if (m) { cur = { id: m[1], title: m[2].trim(), body: [] }; out.push(cur); continue; }
      if (/^## /.test(line) || /^---\s*$/.test(line)) { cur = null; continue; }
      if (cur) cur.body.push(line);
    }
    const num = (r) => { const m = /RULE (\d+)([a-z]?)/.exec(r.id); return Number(m[1]) + (m[2] ? 0.5 : 0); };
    return out.map(r => ({ ...r, body: r.body.join('\n').replace(/^\n+|\n+$/g, '') })).sort((x, y) => num(x) - num(y));
  },

  // The ledger: `## 6. Title` entries, each SEALED or DRAFT by where its heading stands against the seal's byte (the seal is always laid at the end of what
  // was written, so an entry is wholly on one side of it).
  parseLedger(text, sealedBytes) {
    const src = String(text || '');
    const enc = new TextEncoder();
    const out = [];
    const re = /^## (\d+)\. (.+)$/gm;
    let m;
    const heads = [];
    while ((m = re.exec(src))) heads.push({ n: Number(m[1]), title: m[2].trim(), at: m.index, end: m.index + m[0].length });
    heads.forEach((h, i) => {
      const stop = i + 1 < heads.length ? heads[i + 1].at : src.length;
      let body = src.slice(h.end, stop).replace(/\n---\s*$/, '').replace(/^\n+|\n+$/g, '');
      let meta = '';
      const f = /^```\n([\s\S]*?)\n```\n?/.exec(body);
      if (f) { meta = f[1]; body = body.slice(f[0].length).replace(/^\n+/, ''); }
      out.push({ n: h.n, title: h.title, meta, body, sealed: enc.encode(src.slice(0, h.at)).length < sealedBytes });
    });
    return out;
  },

  // Escaped first, marked up after: the documents' own words never reach the page as markup, only the two marks this page knows.
  inline(s) {
    return escHtml(s).replace(/\*\*(.+?)\*\*/g, '<b>$1</b>').replace(/`([^`]+)`/g, '<code>$1</code>');
  },

  badge(state) {
    const tone = { 'sealed': 'green', 'partly sealed': 'yellow', 'not on the chain': 'yellow', 'changed since sealed': 'red' }[state] || 'muted';
    return `<span class="badge badge-${tone}">${escHtml(state.toUpperCase())}</span>`;
  },

  row(id, title, badgeHtml, bodyHtml, extra) {
    return `<details class="laws-row" data-row ${extra || ''}>
      <summary><span class="badge badge-blue laws-id">${escHtml(id)}</span>
        <span class="laws-title">${this.inline(title)}</span>${badgeHtml || ''}</summary>
      <div class="laws-body">${bodyHtml}</div></details>`;
  },

  paint(el, st, rules, ledger) {
    const files = st.files || [];
    const led = files.find(f => f.name === 'LAW_LEDGER.md') || {};
    const sealedFiles = files.filter(f => f.state === 'sealed').length;
    const entries = ledger.error ? [] : this.parseLedger(ledger.text, led.sealed_bytes || 0);
    const ruleList = rules.error ? [] : this.parseRules(rules.text);
    const drafts = entries.filter(e => !e.sealed);
    const blocked = st.why || '';

    const chain = st.ok
      ? `<div class="card"><div class="card-header"><span class="card-title">The chain</span>
          <span class="badge badge-green">PROVES WHOLE</span></div>
          <div class="laws-verdict">${escHtml(st.verdict)}</div>
          <div class="stat-note">${st.links} links · head <code>${escHtml(st.head || '')}</code> · the walk is law.py's own</div></div>`
      : `<div class="card card-critical"><div class="card-header"><span class="card-title">The chain REFUSES</span>
          <span class="badge badge-red">RED</span></div>
          <div class="laws-verdict">${escHtml(st.verdict)}</div>
          <div class="stat-note">While it refuses the law gate refuses every run and no law can be set or sealed. Put the byte back (or
            the file) and read again.</div></div>`;

    const wait = blocked
      ? `<div class="card card-wait"><div class="card-title">Laws cannot be set or sealed now</div>
          <div class="eng-row">${escHtml(blocked)}</div></div>` : '';

    const draft = st.draft_bytes > 0
      ? `<div class="card"><div class="card-header"><span class="card-title">${drafts.length || 'Some'} ${drafts.length === 1 ? 'entry' : 'entries'} below the ledger's seal</span>
          <span class="badge badge-yellow">DRAFT</span></div>
          <div class="left-text">${escHtml(String(st.draft_bytes))} bytes (${drafts.map(e => 'entry ' + e.n).join(', ') || 'unread'}) are written and are read by every
            hand already; they are not frozen. Sealing binds the ledger as it stands.</div>
          <div class="laws-actions"><button class="btn btn-primary" id="laws-seal" ${blocked ? 'disabled title="' + escHtml(blocked) + '"' : ''}>Seal the ledger...</button>
            <span id="laws-seal-ask" hidden>
              <span class="left-text">This binds all ${escHtml(String(led.bytes || ''))} bytes of the ledger on the chain. After it a byte of it cannot
                change: a sealed entry is changed only by a new entry that amends it.</span>
              <button class="btn btn-danger" id="laws-seal-go">Seal it</button>
              <button class="btn" id="laws-seal-no">Not now</button></span>
            <span id="laws-seal-msg" class="laws-msg"></span></div></div>` : '';

    const rulesHtml = ruleList.length
      ? ruleList.map(r => this.row(r.id, r.title, '', `<div class="laws-text">${this.inline(r.body)}</div>`)).join('')
      : `<div class="eng-row eng-bad">${escHtml(this.RULES)} could not be read: ${escHtml(rules.error || 'no rule in it')}</div>`;

    const ledgerHtml = entries.length
      ? entries.map(e => this.row('Entry ' + e.n, e.title, e.sealed ? this.badge('sealed') : '<span class="badge badge-yellow">DRAFT</span>',
          `${e.meta ? `<div class="laws-meta">${escHtml(e.meta)}</div>` : ''}<div class="laws-text">${this.inline(e.body)}</div>`)).join('')
      : `<div class="eng-row eng-bad">${escHtml(this.LEDGER)} could not be read: ${escHtml(ledger.error || 'no entry in it')}</div>`;

    const filesHtml = files.map(f => this.row(f.kind || 'LAW', f.name, this.badge(f.state),
      `<div class="laws-meta">${escHtml(f.name)} · ${(f.bytes / 1024).toFixed(1)} KB · sha256 ${escHtml(f.fingerprint || '')}…${
        f.sealed_bytes ? ' · sealed to byte ' + f.sealed_bytes + (f.draft_bytes ? ', ' + f.draft_bytes + ' bytes draft below it' : '') : ''}${
        f.by ? ' · laid by ' + escHtml(f.by) : ''}${f.links && f.links.length ? ' · link' + (f.links.length > 1 ? 's' : '') + ' ' + f.links.join(', ') : ''}</div>
       <div class="laws-text" data-file="${escHtml(f.name)}"><span class="stat-note">Open to read it.</span></div>`, `data-lazy="${escHtml(f.name)}"`)).join('');

    const counts = `${ruleList.length} standing rules, ${entries.length} ledger entries, ${sealedFiles} of ${files.length} laws sealed`;
    el.innerHTML = `<div id="laws-page">
      ${this.inPanel ? `<div class="laws-head">
        <div class="stat-note">What every hand and every agent reads first: ${counts}.</div>
        <div class="laws-actions"><input class="input" id="laws-find" placeholder="Find in the laws..."><button class="btn" id="laws-again">Read it again</button></div>
      </div>` : `<div class="page-header">
        <div>
          <div class="page-title">Laws</div>
          <div class="page-subtitle">What every hand and every agent reads first: ${counts}.</div>
        </div>
        <div class="flex">
          <div class="search-bar"><input class="input" id="laws-find" placeholder="Find in the laws..."></div>
          <button class="btn" id="laws-again">Read it again</button>
        </div>
      </div>`}
      ${chain}${wait}${draft}
      <div class="card mt-16" data-sec>
        <div class="card-header"><span class="card-title">Set a law</span></div>
        <div class="left-text">It is appended to the ledger below the seal, as a draft: <code>law/</code> is read first by every hand (RULE 0, SITTING LAW 6), so it
          is read from the moment it is set. Sealing is a separate click, above.</div>
        <div class="laws-form">
          <input class="input" id="laws-title" maxlength="120" placeholder="Its name, one line (e.g. SITTING LAW 9 -- nothing is sent without his word)">
          <textarea class="textarea" id="laws-text" rows="7" placeholder="The law, in the words it is to be read in."></textarea>
          <textarea class="textarea" id="laws-from" rows="2" placeholder="Where it came from: his words, the date (optional)"></textarea>
          <div class="laws-actions"><button class="btn btn-primary" id="laws-add" ${blocked ? 'disabled title="' + escHtml(blocked) + '"' : ''}>Set it as a draft</button>
            <span id="laws-msg" class="laws-msg"></span></div>
        </div>
      </div>
      <div class="card mt-16" data-sec>
        <div class="card-header"><span class="card-title">The standing rules · ${escHtml(this.RULES)}</span>
          <span class="badge badge-blue">${ruleList.length}</span></div>
        <div class="stat-note">Read first by every hand, every turn. A rule copied from the ledger binds here, as it always has.</div>
        ${rulesHtml}
      </div>
      <div class="card mt-16" data-sec>
        <div class="card-header"><span class="card-title">The ledger · ${escHtml(this.LEDGER)}</span>
          <span class="badge badge-${drafts.length ? 'yellow' : 'green'}">${drafts.length ? drafts.length + ' draft' : 'all sealed'}</span></div>
        ${ledgerHtml}
      </div>
      <div class="card mt-16" data-sec>
        <div class="card-header"><span class="card-title">The laws on the chain · law/</span>
          <span class="badge badge-${sealedFiles === files.length ? 'green' : 'yellow'}">${sealedFiles} of ${files.length} sealed</span></div>
        ${filesHtml}
      </div>
      <div class="stat-note mt-16">Read from the door: <span class="brief-src">law_status</span> <span class="brief-src">records</span>${
        ledger.sha256 ? ' · ledger sha256 ' + escHtml(String(ledger.sha256).slice(0, 16)) : ''}</div>
    </div>`;

    this.bind(el, st);
  },

  bind(el, st) {
    const $ = (id) => document.getElementById(id);
    $('laws-again').onclick = () => this.render(el, this.inPanel);
    $('laws-find').oninput = (e) => this.find(e.target.value);
    el.querySelectorAll('details[data-lazy]').forEach(d => {
      d.addEventListener('toggle', async () => {
        if (!d.open || d.dataset.loaded) return;
        d.dataset.loaded = '1';
        const box = d.querySelector('[data-file]');
        try {
          const doc = JSON.parse(await App.tool('records', { name: 'law/' + d.dataset.lazy }));
          box.innerHTML = this.inline(doc.text || '');
        } catch (e) {
          delete d.dataset.loaded;
          box.innerHTML = `<span class="eng-bad">could not be read: ${escHtml(e.message || 'refused')}</span>`;
        }
      });
    });

    const say = (id, text, bad) => { const m = $(id); if (m) { m.textContent = text; m.className = 'laws-msg' + (bad ? ' laws-bad' : ''); } };

    if ($('laws-add')) $('laws-add').onclick = async () => {
      const title = $('laws-title').value.trim(), text = $('laws-text').value.trim(), from = $('laws-from').value.trim();
      if (!title || !text) { say('laws-msg', 'A law needs a name and its words.', true); return; }
      $('laws-add').disabled = true;
      say('laws-msg', 'Setting it...');
      try {
        const r = JSON.parse(await App.tool('law_add', { title, text, from }));
        if (r.state !== 'added') { say('laws-msg', 'Not set: ' + (r.why || 'refused'), true); $('laws-add').disabled = false; return; }
        toast('Law set as entry ' + r.entry + ' -- a draft until it is sealed');
        this.render(el, this.inPanel);
      } catch (e) {
        say('laws-msg', 'Not set: ' + (e.message || 'refused'), true);
        $('laws-add').disabled = false;
      }
    };

    if ($('laws-seal')) {
      $('laws-seal').onclick = () => { $('laws-seal').hidden = true; $('laws-seal-ask').hidden = false; };
      $('laws-seal-no').onclick = () => { $('laws-seal-ask').hidden = true; $('laws-seal').hidden = false; };
      $('laws-seal-go').onclick = async () => {
        $('laws-seal-go').disabled = true;
        say('laws-seal-msg', 'Sealing...');
        try {
          const r = JSON.parse(await App.tool('law_seal', {}));
          if (r.state !== 'sealed') { say('laws-seal-msg', 'Not sealed: ' + (r.why || 'refused'), true); $('laws-seal-go').disabled = false; return; }
          toast('The ledger is sealed to byte ' + r.sealed_to);
          this.render(el, this.inPanel);
        } catch (e) {
          say('laws-seal-msg', 'Not sealed: ' + (e.message || 'refused'), true);
          $('laws-seal-go').disabled = false;
        }
      };
    }
  },

  find(q) {
    q = (q || '').trim().toLowerCase();
    document.querySelectorAll('#laws-page [data-row]').forEach(r => {
      r.hidden = q !== '' && r.textContent.toLowerCase().indexOf(q) < 0;
    });
  },

  style() {
    if (document.getElementById('laws-style')) return;
    const s = document.createElement('style');
    s.id = 'laws-style';
    s.textContent = `
      .laws-row { border-bottom: 1px solid var(--border); padding: 6px 0; }
      .laws-row:last-child { border-bottom: 0; }
      .laws-row[hidden] { display: none; }
      .laws-row > summary { cursor: pointer; display: flex; gap: 12px; align-items: center; list-style: none; padding: 6px 0; }
      .laws-row > summary::-webkit-details-marker { display: none; }
      .laws-id { font-family: var(--font-mono); white-space: nowrap; }
      .laws-title { font-weight: 600; color: var(--text); font-size: var(--t-md); flex: 1; }
      .laws-body { padding: 4px 0 12px 8px; }
      .laws-text { color: var(--text-2); font-size: var(--t-md); line-height: 1.6; white-space: pre-wrap; }
      .laws-meta { font-family: var(--font-mono); font-size: var(--t-xs); color: var(--text-3); white-space: pre-wrap; margin-bottom: 8px; }
      .laws-verdict { font-family: var(--font-mono); font-size: var(--t-sm); white-space: pre-wrap; margin: 6px 0; }
      .laws-form { display: grid; gap: 8px; margin-top: 10px; }
      .laws-actions { display: flex; gap: 10px; align-items: center; flex-wrap: wrap; margin-top: 10px; }
      .laws-msg { color: var(--text-2); font-size: var(--t-sm); }
      .laws-msg.laws-bad, .laws-text .eng-bad { color: var(--red); }
      .laws-head { margin-bottom: 12px; }
      .laws-head .laws-actions { flex-wrap: nowrap; margin-top: 6px; }
      .laws-head .input { flex: 1; width: auto; min-width: 0; }`;
    document.head.appendChild(s);
  },
};
