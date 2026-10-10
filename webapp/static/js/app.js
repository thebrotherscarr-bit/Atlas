// ATLAS App — SPA router and pages
const App = {
  currentPage: 'agent',
  data: {},

  // EVERY PAGE, IN ONE LIST (2026-10-07: the sidebar retired at his word, "go, next: retire the sidebar and the old
  // dashboard"). The sidebar's links were the list until then; now the front page's Pages menu, the palette and the
  // crumb all read this one. The old Dashboard is not on it: /dashboard opens the front page.
  PAGES: [
    { page: 'agent', href: '/', label: 'Home' },
    { page: 'left', href: '/left', label: "What's left" },
    { page: 'chat', href: '/chat', label: 'Watchboard' },
    { page: 'agents', href: '/agents', label: 'Agents' },
    { page: 'workflows', href: '/workflows', label: 'Workflows' },
    { page: 'laws', href: '/laws', label: 'Laws' },
    { page: 'evals', href: '/evals', label: 'Evals' },
    { page: 'records', href: '/records', label: 'Records' },
    { page: 'flows', href: '/flows', label: 'GitHub' },
    { page: 'tools', href: '/tools', label: 'Tools' },
    { page: 'settings', href: '/settings', label: 'Settings' }
  ],

  // A PAGE WITH A TAB OPENS ITS TAB (2026-10-10, his word: "everything stays within the new interface, as previously
  // judged"; his cards: "Every page with a tab", "Keep them, opening the tab"). Workflows, Laws, GitHub and Settings stay
  // on the list above, so the Pages menu, the palette and an old bookmark still reach them, and each address opens the
  // front page with that Inspector tab open, as /dashboard opens the front page: the tab draws what the page drew, beside
  // the cards that wait for his hand. The page's address, and the Inspector tab it opens. AND THE SIX THAT HAD NO TAB, put
  // to him page by page the same day (D17's last part, his cards): the Watchboard opens Run, which draws its floor, tool
  // calls and wire; Agents, Tools and Records open Registry & Docs, which carries their seats, tools and documents; Evals
  // opens the Audit Ledger, which carries the hand-scored evals; and What's left has a tab of its own.
  TABBED: {
    workflows: 'flows', laws: 'laws', flows: 'github', settings: 'settings',
    chat: 'run', agents: 'docs', tools: 'docs', records: 'docs', evals: 'ledger', left: 'left'
  },

  init() {
    this.router();
    window.addEventListener('popstate', () => this.router());
    // Ctrl+K / Cmd+K, anywhere. The one global key this console binds.
    Palette.bind();
    // Unsent work, read once as the glass loads (it was the sidebar's badge); the front page's GitHub door shows it.
    this.paintOwed();
    API.sse((e) => this.onEvent(e));
    // Every page follows a turn started in another browser, not just the
    // one that asked for it. Idempotent: mirror() returns at once if it is
    // already listening.
    Run.mirror();
  },

  router() {
    // THE OLD DASHBOARD IS RETIRED (2026-10-07): a bookmark to it opens the front page, which carries what it had.
    if (location.pathname === '/dashboard') history.replaceState(null, '', '/');
    // A PAGE WITH A TAB OPENS ITS TAB (2026-10-10, TABBED): the front page, with the Inspector open on it. A front page
    // already drawn opens the tab where it stands, and nothing on it is drawn again.
    const first = location.pathname.slice(1).split('/')[0];
    const tab = Object.hasOwn(this.TABBED, first) ? this.TABBED[first] : '';
    if (tab) {
      history.replaceState(null, '', '/');
      if (this.currentPage === 'agent' && document.getElementById('ag')) { Agent.door(tab); return; }
      Agent.open = true;
      Agent.tab = tab;
    }
    const path = location.pathname.slice(1) || 'agent';
    const parts = path.split('/');
    this.currentPage = parts[0];
    this.pageParam = parts[1] || null;
    // THE FRONT PAGE IS DRAWN ON ITS OWN GROUND (agent.css, under this class); every other page is the console's.
    document.body.classList.toggle('is-agent', this.currentPage === 'agent');
    this.paintCrumb();
    this.render();
  },

  // THE CRUMB SAYS THE ESTATE AND THE PAGE, AND STOPS -- unless the path
  // really is deeper, which is the only case where a third step is a fact
  // rather than furniture. The page's NAME comes off the one page list
  // (PAGES), so a label renamed there (Flows -> Version control, 2026-09-10)
  // is renamed here in the same stroke and cannot drift. With the sidebar
  // retired (2026-10-07) the crumb's ATLAS is every page's way home.
  paintCrumb() {
    const el = document.getElementById('crumb');
    if (!el) return;
    const link = this.PAGES.find(p => p.page === this.currentPage);
    // A page off the list (traces, messages, playground) still routes, so it
    // still gets a crumb -- titled from the path when the list has no line.
    const name = link
      ? link.label
      : this.currentPage.charAt(0).toUpperCase() + this.currentPage.slice(1);
    const home = this.currentPage === 'agent';
    const parts = [`<a href="/" onclick="event.preventDefault();history.pushState(null,'','/');App.router();">ATLAS</a>`];
    if (!home) {
      parts.push('<span class="sep">/</span>');
      if (this.pageParam) {
        const back = '/' + this.currentPage;
        parts.push(`<a href="${back}" onclick="event.preventDefault();history.pushState(null,'','${back}');App.router();">${escHtml(name)}</a>`);
        parts.push('<span class="sep">/</span>');
        parts.push(`<span class="here">${escHtml(decodeURIComponent(this.pageParam))}</span>`);
      } else {
        parts.push(`<span class="here">${escHtml(name)}</span>`);
      }
    }
    el.innerHTML = parts.join('');
    el.hidden = home;
  },

  // UNSENT WORK, ON THE GITHUB DOOR (2026-10-08, on his card: "A GitHub door in the top bar"; it was on the Aider Pair
  // door from 2026-10-07, on his card: "On the Aider Pair door"). The sidebar's badge on Version control moved with the
  // page, now GitHub: a number on the front page's GitHub door, shown only while some world holds work GitHub does not
  // have -- a WARNING rather than a tally, so the sidebar's three plain counts did not come with it. A NUMBER THE RECORD
  // CAN PROVE: read off the tools that own it, across every carried world, and hidden when the read fails rather than
  // guessed. Read when the glass loads and after every act on the GitHub page (flows.js), the one place the number
  // moves -- the same reads the sidebar made -- and kept, so the front page, drawn again on every visit, shows it
  // without asking again (showOwed).
  async paintOwed() {
    let owed = null;
    try {
      const m = await this.tool('muster', {});
      const worlds = (m || '').split(String.fromCharCode(10))
        .map(s => s.trim()).filter(s => s && !s.endsWith(':'));
      let n = 0;
      for (const w of worlds) {
        try {
          const g = JSON.parse(await this.tool('git', { project: w }));
          if (g.is_repo) n += (g.changed || 0) + (g.untracked || 0) + (g.ahead || 0);
        } catch { /* one unreadable world must not blank the others */ }
      }
      owed = n;
    } catch { owed = null; }
    this._owed = owed;
    this.showOwed();
  },

  showOwed() {
    const el = document.getElementById('ag-owed');
    if (!el) return;
    if (!this._owed) { el.hidden = true; return; }
    el.textContent = String(this._owed);
    el.title = this._owed + ' unsaved or unsent across the worlds: the GitHub tab is under this door';
    el.hidden = false;
  },

  // THE COVENANT IS READ OFF THE RECORD, NEVER TYPED HERE (2026-09-29, WHAT'S
  // LEFT C29). The sidebar carried the house covenant as a literal, one of the
  // 59 copies the door stopped minting from on 2026-09-25; the operator's own
  // declaration says it, and us_to_vc mints a credential from that declaration
  // in exactly that namespace. Read once, as a background read, and kept; a
  // door that cannot serve it leaves the line saying so rather than a number
  // from memory, and the next drawing asks again. Since the sidebar retired
  // (2026-10-07, on his card: "Front page footer") the line is the front
  // page's footer, drawn again on every visit and painted from the keeping.
  async paintCovenant() {
    const el = document.getElementById('covenant');
    if (!el) return;
    if (this._covenant) { el.textContent = this._covenant; return; }
    if (this._covenantRead) return;
    this._covenantRead = true;
    try {
      const vc = JSON.parse(await this.tool('us_to_vc', { path: 'agents/operator.us', project: 'atlas' }, true));
      const c = (vc.credentialSubject || {}).covenant || '';
      this._covenant = c ? 'covenant: ' + c : 'covenant: none declared in the record';
    } catch (e) {
      this._covenantRead = false;
      const now = document.getElementById('covenant');
      if (now) now.textContent = 'covenant: not read (the door did not answer)';
      return;
    }
    const now = document.getElementById('covenant');
    if (now) now.textContent = this._covenant;
  },

  onEvent(e) {
    if (e.type === 'trace_added') toast('New trace recorded');
    if (e.type === 'eval_added') toast('New eval recorded');
    if (e.type === 'message_sent') toast('Message sent');
    if (e.type === 'message_received') {
      toast('Message received');
      if (this.currentPage === 'messages') this.render();
    }
    if (e.type === 'chat_opened') toast('Chat session opened');
    if (e.type === 'chat_done') toast('Chat turn witnessed');
    if (e.type === 'prompt_saved') toast('Prompt version folded');
    if (e.type === 'prompt_ran') toast('Prompt measured');
    if (e.type === 'prompt_evaled') toast('Eval scored');
    if (e.type === 'flow_saved') toast('Flow version folded');
    if (e.type === 'flow_ran') toast('Flow fired');
    if (e.type === 'flow_resumed') toast('Flow moved by the hand');
  },

  async render() {
    const el = document.getElementById('content');
    switch (this.currentPage) {
      // THE FRONT PAGE (agent.js): one terminal, five doors, an Inspector.
      case 'agent': await Agent.render(el); break;
      // The launchpad (home.js) retired on 2026-10-07 at his word; the router sends /dashboard to the front page.
      case 'traces': this.pageParam ? await this.renderTraceDetail(el) : await this.renderTraces(el); break;
      case 'playground': await Play.render(el); break;
      case 'messages': await this.renderMessages(el); break;
      // Every other page on the list is one of the front page's Inspector tabs since 2026-10-10 -- Workflows, Laws,
      // GitHub and Settings first, then (D17's last part, his cards) What's left, the Watchboard, Agents, Tools, Evals
      // and Records: the router opens them (TABBED) and draws none of them here.
      default: el.innerHTML = '<div class="empty"><div class="empty-icon">?</div><div class="empty-text">Page not found</div></div>';
    }
  },

  // === THE DOOR'S TOOLS, AS THE PAGES ASK THEM ===
  // A number the record cannot prove is not shown -- SPEC 3 invariant 10.
  //
  // `background` IS FOR THE POLLING READS AND NOTHING ELSE (2026-09-16). The
  // front page's quiet reads pass it (the Dashboard's did, until it retired on
  // 2026-10-07), and the glass answers such a read without keeping it as a
  // trace (handlers.go, backgroundReads). Every other caller leaves it off, so
  // whatever he asks for himself is kept as before.
  async tool(name, args, background) {
    const r = await API.callTool(name, args || {}, background);
    try {
      const env = JSON.parse(r.output);
      if (env.error) throw new Error(env.error.message || 'refused');
      const c = env.result && env.result.content;
      const text = c && c[0] ? c[0].text : '';
      if (env.result && env.result.isError) throw new Error(text || 'refused');
      return text;
    } catch (e) { throw new Error(e.message || 'unreadable answer'); }
  },

  // === THE SEATS ===
  // A seat, as its card: what Registry & Docs draws for each seat declared in agents/, read through the `seats` tool from
  // agents/*.md and pipelines.md -- the source of truth. The Agents page drew the same cards until its address opened that
  // tab (2026-10-10, D17's last part).
  seatCard(s, i) {
    if (s.error) {
      return `<div class="card seat-card"><div class="eng-row eng-bad">
        <b>${escHtml(s.file)}</b> could not be read: ${escHtml(s.error)}</div></div>`;
    }
    const f = s.fields || {};
    const order = s.field_order || Object.keys(f);
    // The model and the gate lead, because they are what he tunes.
    const lead = ['Model Target', 'Stage', 'When', 'Wakes On', 'Wakes'];
    const rest = order.filter(k => lead.indexOf(k) < 0 && (f[k] || '').trim());

    const stands = (s.stands_in || []).map(st =>
      `<span class="badge badge-blue" title="${escHtml(st.note || '')}">${escHtml(st.pipeline)}
       <span class="muted">#${st.step}</span>${st.when ? ' · ' + escHtml(st.when) : ''}</span>`).join(' ');

    // THE NAME DOES NOT LEAVE THE FRONT PAGE (2026-10-06): the card stands in the Inspector's Registry & Docs tab, and its
    // file, with its receipt, opens under it from a button of its own (agent.js, `seat-file`) -- what the seat's own page
    // showed until the Agents page's address opened the tab (2026-10-10, D17's last part).
    const title = `<span class="card-title">${escHtml(s.name)}</span>`;
    return `<div class="card seat-card" data-seat="${escHtml((s.name || '').toLowerCase())}">
      <div class="card-header">
        ${title}
        <span class="flex">
          ${f['Model Target'] ? `<code class="seat-model">${escHtml(f['Model Target'])}</code>` : ''}
          ${s.prompt ? `<button class="btn btn-sm" data-prompt="${i}">prompt</button>` : ''}
          <button class="btn btn-sm" data-act="seat-file" data-i="${i}" data-file="agents/${escHtml(s.file)}">file</button>
        </span>
      </div>
      <div class="seat-rows">
        ${lead.filter(k => (f[k] || '').trim()).map(k =>
          `<div class="seat-row"><span class="seat-k">${escHtml(k)}</span>
           <span class="seat-v">${escHtml(f[k])}</span></div>`).join('')}
        ${rest.map(k =>
          `<div class="seat-row"><span class="seat-k">${escHtml(k)}</span>
           <span class="seat-v muted">${escHtml(f[k])}</span></div>`).join('')}
      </div>
      <div class="seat-stands">
        ${stands || '<span class="muted">stands in no pipeline — racked, summoned when its flag is raised</span>'}
        <span class="brief-src">${escHtml(s.file)}</span>
      </div>
      ${s.prompt ? `<pre class="seat-prompt" id="prompt-${i}" hidden>${escHtml(s.prompt)}</pre>` : ''}
      <div id="seat-file-${i}"></div>
    </div>`;
  },

  filterAgents(q) {
    q = (q || '').trim().toLowerCase();
    document.querySelectorAll('#agent-grid .seat-card').forEach(c => {
      const hay = (c.dataset.seat || '') + ' ' + c.textContent.toLowerCase();
      c.hidden = q !== '' && hay.indexOf(q) < 0;
    });
  },

  // === TRACES ===
  async renderTraces(el) {
    el.innerHTML = '<div class="loading">Loading traces...</div>';
    try {
      const data = await API.listTraces();
      const traces = data.traces || [];
      el.innerHTML = `
        <div class="page-header">
          <div>
            <div class="page-title">Trace Log</div>
            <div class="page-subtitle">${traces.length} traces recorded</div>
          </div>
          <div class="search-bar">
            <input class="input" placeholder="Filter traces..." oninput="App.filterTraces(this.value)">
          </div>
        </div>
        <div class="card">
          ${traces.length === 0
            ? '<div class="empty"><div class="empty-icon">&#128269;</div><div class="empty-text">No traces yet. Call a tool from Registry & Docs to start recording.</div></div>'
            : '<div class="table-wrap"><table><thead><tr><th>Tool</th><th>Agent</th><th>Status</th><th>Duration</th><th>Hash</th><th>Time</th></tr></thead><tbody>' +
              traces.map(t => `
                <tr onclick="location.href='/traces/${t.id}'" style="cursor:pointer">
                  <td><code>${escHtml(t.tool)}</code></td>
                  <td>${escHtml(t.agent_id || '—')}</td>
                  <td><span class="badge ${t.status === 'ok' ? 'badge-green' : 'badge-red'}">${escHtml(t.status)}</span></td>
                  <td>${t.duration_ms}ms</td>
                  <td><span class="hash">${escHtml((t.hash || '').slice(0, 16))}...</span></td>
                  <td>${timeAgo(t.created_at)}</td>
                </tr>
              `).join('') +
              '</tbody></table></div>'}
        </div>
      `;
    } catch (e) {
      el.innerHTML = `<div class="empty"><div class="empty-icon">!</div><div class="empty-text">${escHtml(e.message)}</div></div>`;
    }
  },

  filterTraces(q) {
    q = q.toLowerCase();
    document.querySelectorAll('tbody tr').forEach(row => {
      row.style.display = row.textContent.toLowerCase().includes(q) ? '' : 'none';
    });
  },

  async renderTraceDetail(el) {
    el.innerHTML = '<div class="loading">Loading trace...</div>';
    try {
      const data = await API.getTrace(this.pageParam);
      const t = data.trace;
      const evals = data.evals || [];
      el.innerHTML = `
        <div class="page-header">
          <div>
            <div class="page-title">Trace: ${escHtml(t.tool)}</div>
            <div class="page-subtitle">${escHtml(t.id)}</div>
          </div>
          <a href="/traces" class="btn" onclick="event.preventDefault();history.pushState(null,'','/traces');App.router();">Back</a>
        </div>
        <div class="grid-2">
          <div class="card">
            <div class="card-header"><span class="card-title">Trace Details</span></div>
            <table>
              <tr><td>Tool</td><td><code>${escHtml(t.tool)}</code></td></tr>
              <tr><td>Agent</td><td>${escHtml(t.agent_id || '—')}</td></tr>
              <tr><td>Status</td><td><span class="badge ${t.status === 'ok' ? 'badge-green' : 'badge-red'}">${escHtml(t.status)}</span></td></tr>
              <tr><td>Duration</td><td>${t.duration_ms}ms</td></tr>
              <tr><td>Time</td><td>${new Date(t.created_at).toLocaleString()}</td></tr>
              <tr><td>Hash</td><td><span class="hash">${escHtml(t.hash || '—')}</span></td></tr>
            </table>
          </div>
          <div class="card">
            <div class="card-header"><span class="card-title">Input / Output</span></div>
            <div class="form-group"><div class="form-label">Input</div><pre style="font-size:12px;color:var(--text-2);background:var(--bg);padding:12px;border-radius:var(--radius);overflow-x:auto;max-height:200px">${escHtml(t.input || '—')}</pre></div>
            <div class="form-group"><div class="form-label">Output</div><pre style="font-size:12px;color:var(--text-2);background:var(--bg);padding:12px;border-radius:var(--radius);overflow-x:auto;max-height:200px">${escHtml(t.output || '—')}</pre></div>
          </div>
        </div>
        <div class="card mt-16">
          <div class="card-header">
            <span class="card-title">Evals</span>
            <button class="btn btn-sm" onclick="App.addEval('${t.id}')">+ Add Eval</button>
          </div>
          ${evals.length === 0
            ? '<div class="empty"><div class="empty-text">No evals for this trace.</div></div>'
            : '<div class="table-wrap"><table><thead><tr><th>Name</th><th>Score</th><th>Passed</th><th>Detail</th></tr></thead><tbody>' +
              evals.map(e => `
                <tr>
                  <td>${escHtml(e.name)}</td>
                  <td><span class="eval-score ${e.passed ? 'pass' : 'fail'}">${e.score}</span></td>
                  <td><span class="badge ${e.passed ? 'badge-green' : 'badge-red'}">${e.passed ? 'PASS' : 'FAIL'}</span></td>
                  <td>${escHtml(e.detail || '—')}</td>
                </tr>
              `).join('') + '</tbody></table></div>'}
        </div>
      `;
    } catch (e) {
      el.innerHTML = `<div class="empty"><div class="empty-icon">!</div><div class="empty-text">${escHtml(e.message)}</div></div>`;
    }
  },

  // THE EVAL SCORER ASKS IN THE PAGE (2026-09-29): its three prompt()s became
  // one form in the modal, for the reason invokeTool gives below.
  addEval(traceId) {
    const modal = document.getElementById('modal');
    const content = document.getElementById('modal-content');
    content.innerHTML = `
      <h3 style="margin-bottom:12px">Score trace <code>${escHtml(traceId)}</code></h3>
      <form id="eval-form">
        <div class="form-group"><div class="form-label">Eval name</div>
          <input class="input" id="eval-name" autocomplete="off"></div>
        <div class="form-group"><div class="form-label">Score (0-1)</div>
          <input class="input" id="eval-score" type="number" min="0" max="1" step="0.01" value="0"></div>
        <div class="form-group"><div class="form-label">Detail (optional)</div>
          <input class="input" id="eval-detail" autocomplete="off"></div>
        <div class="flex-between mt-16">
          <span class="muted" id="eval-note"></span>
          <span class="flex" style="gap:8px">
            <button class="btn" type="button" id="eval-close">Close</button>
            <button class="btn btn-primary" type="submit">Add the eval</button>
          </span>
        </div>
      </form>`;
    modal.style.display = 'flex';
    document.getElementById('eval-close').onclick = () => this.closeModal();
    document.getElementById('eval-name').focus();
    document.getElementById('eval-form').onsubmit = async (e) => {
      e.preventDefault();
      const name = document.getElementById('eval-name').value.trim();
      if (!name) { document.getElementById('eval-note').textContent = 'Name it first.'; return; }
      const score = parseFloat(document.getElementById('eval-score').value || '0');
      const detail = document.getElementById('eval-detail').value || '';
      let threshold = 0.5;
      try {
        const s = await API.getSetting('eval_threshold');
        if (s.value) threshold = parseFloat(s.value);
      } catch {}
      const passed = score >= threshold;
      try {
        await API.addEval({ id: 'e-' + Date.now(), trace_id: traceId, name, score, passed, detail });
        this.closeModal();
        toast('Eval added');
        this.router();
      } catch (err) { toast(err.message, 'error'); }
    };
  },

  // THE TOOLS PAGE ASKS IN THE PAGE, NOT IN A POP-UP (2026-09-29, WHAT'S LEFT
  // C26). This collected the arguments with a native prompt(), which the
  // desktop app's browser pane dismisses unseen -- so Call did nothing there
  // -- and which this glass's own rule forbids everywhere else (the spec's
  // §4.7: no prompt() anywhere; prompts.js says the same). The modal the
  // console already owns takes the arguments as a form: the tool's declared
  // arguments listed over the box so nobody has to remember them, JSON in,
  // and the answer painted into the same modal with a way to call again.
  // Since 2026-10-10 the Tools page's address opens Registry & Docs, whose
  // Call hands this form the list it read (App._tools).
  invokeTool(name, raw) {
    const tool = (this._tools || []).find(t => t.name === name) || {};
    const props = Object.keys((tool.inputSchema && tool.inputSchema.properties) || {});
    const req = (tool.inputSchema && tool.inputSchema.required) || [];
    const modal = document.getElementById('modal');
    const content = document.getElementById('modal-content');
    content.innerHTML = `
      <h3 style="margin-bottom:12px">Call <code>${escHtml(name)}</code></h3>
      <div class="muted" style="margin-bottom:8px">${escHtml(tool.description || '')}</div>
      <form id="tool-call-form">
        <div class="form-group">
          <div class="form-label">Arguments (JSON)${props.length
            ? ' — ' + props.map(p => escHtml(p) + (req.indexOf(p) >= 0 ? '' : '?')).join(', ') : ''}</div>
          <textarea class="textarea" id="tool-call-args" rows="4" spellcheck="false">${escHtml(raw || '{}')}</textarea>
        </div>
        <div class="flex-between mt-16">
          <span class="muted" id="tool-call-note"></span>
          <span class="flex" style="gap:8px">
            <button class="btn" type="button" id="tool-call-close">Close</button>
            <button class="btn btn-primary" type="submit">Call it</button>
          </span>
        </div>
      </form>`;
    modal.style.display = 'flex';
    document.getElementById('tool-call-close').onclick = () => this.closeModal();
    const box = document.getElementById('tool-call-args');
    box.focus();
    document.getElementById('tool-call-form').onsubmit = (e) => {
      e.preventDefault();
      let args;
      try { args = JSON.parse(box.value || '{}'); }
      catch { document.getElementById('tool-call-note').textContent = 'That is not JSON.'; return; }
      this.callTool(name, args, box.value);
    };
  },

  async callTool(name, args, raw) {
    const content = document.getElementById('modal-content');
    content.innerHTML = `<div class="loading">Calling ${escHtml(name)}...</div>`;
    try {
      const result = await API.callTool(name, args);
      content.innerHTML = `
        <h3 style="margin-bottom:12px">Tool Result: <code>${escHtml(name)}</code></h3>
        <div class="form-group"><div class="form-label">Output</div><pre style="font-size:12px;background:var(--bg);padding:12px;border-radius:var(--radius);max-height:300px;overflow:auto">${escHtml(result.output || '')}</pre></div>
        <div class="form-group"><div class="form-label">Hash</div><span class="hash">${escHtml(result.hash || '')}</span></div>
        <div class="flex-between mt-16">
          <span style="font-size:12px;color:var(--text-3)">Trace: ${escHtml(result.trace_id || '')} — ${result.duration_ms || 0}ms</span>
          <span class="flex" style="gap:8px">
            <button class="btn" id="tool-call-again">Call again</button>
            <button class="btn" id="tool-call-done">Close</button>
          </span>
        </div>`;
      document.getElementById('tool-call-again').onclick = () => this.invokeTool(name, raw);
      document.getElementById('tool-call-done').onclick = () => this.closeModal();
    } catch (e) {
      content.innerHTML = `<div class="empty"><div class="empty-icon">!</div><div class="empty-text">${escHtml(e.message)}</div></div>
        <button class="btn mt-16" id="tool-call-done">Close</button>`;
      document.getElementById('tool-call-done').onclick = () => this.closeModal();
    }
  },

  closeModal() {
    document.getElementById('modal').style.display = 'none';
  },

  // WHAT THIS WORLD HAS PROVED, read from its own record and nowhere else.
  // Every card names the file it came from; a world that never ran a suite
  // says so rather than rendering as a zero, because "zero passed" and "never
  // run" are opposite claims.
  // Paints into whichever box it is given -- it lived on Evals, then Records,
  // then the Dashboard, and now the front page's Audit Ledger tab. The id is
  // the caller's business.
  //
  // `only` SPLITS THE TWO HALVES, because they answer different questions and
  // they belong on different pages now (the operator, 2026-09-10): the SCORES
  // are "is the build sound", which is the first thing the Dashboard should
  // say; the ESTATE is "what has sat here", which is the record, and is drawn
  // under the rows in the Audit Ledger tab since Records' address opened the
  // tabs (2026-10-10). One reader, one `proofs` call, two placements.
  //   'scores' -- strokes, smoke, standup, standups run, parity
  //   'estate' -- sittings, tolls, runs, and the live standups table
  //   'ledger' -- the proof's rows (proofLedger), for the front page's Audit Ledger tab
  //   omitted  -- both, as before
  //
  // `read` IS AN ANSWER THE CALLER ALREADY HOLDS (2026-09-15): the tool's text,
  // or the {err} Home.read's ask() hands back for a refusal. The Dashboard's
  // 15-second read asked for `proofs` in its own Promise.all and this function
  // then asked again, so every refresh read the record files twice and wrote
  // two 20 KB traces -- 14,440 `proofs` calls over 2026-09-12..14, half of
  // everything the glass called. A caller holding the answer hands it over; a
  // caller that holds none still asks.
  async paintProof(boxId, only, read) {
    const box = document.getElementById(boxId);
    if (!box) return;
    let p;
    try {
      let text = read;
      if (text === undefined) {
        box.innerHTML = '<div class="loading">Reading the record...</div>';
        text = await this.tool('proofs', {});
      }
      if (text && typeof text === 'object') throw new Error(text.err || 'unreadable');
      p = JSON.parse(text);
    } catch (e) {
      box.innerHTML = `<div class="card"><div class="eng-row eng-bad">
        The record could not be read: ${escHtml(e.message || 'refused')}
        <span class="brief-src">proofs</span></div></div>`;
      return;
    }

    const cards = [];
    // `delta` is OPTIONAL and is only ever passed where the record holds a
    // previous value to compare against. A card with nothing to compare shows
    // no delta rather than a zero -- "unchanged" and "never measured twice"
    // are different claims, and the second is the true one here.
    const card = (label, value, tone, note, src, delta) => cards.push(
      `<div class="stat">
         <div class="stat-head">
           <div class="stat-label">${escHtml(label)}</div>
           ${delta ? `<span class="stat-delta ${delta.dir}">${escHtml(delta.text)}</span>` : ''}
         </div>
         <div class="stat-value ${tone || ''}">${value}</div>
         ${note ? `<div class="stat-note">${note}</div>` : ''}
         <div class="brief-src">${escHtml(src)}</div></div>`);

    // ---- the suites, as they stamped themselves ----------------------
    const su = p.suites || {};
    for (const name of ['strokes', 'smoke']) {
      const r = su[name];
      if (!r) {
        card(name, '<span class="muted">never run</span>', '', p.suites_error || 'no stamp in this world', 'tests/last_run.json');
        continue;
      }
      card(name, `${r.passed}<span class="muted">/${r.total}</span>`,
           r.green ? 'green' : 'red',
           r.green ? 'green · ' + when(r.at) : 'RED · ' + ((r.failures || []).join(', ') || 'see the run'),
           'tests/last_run.json');
    }

    // ---- the live standups -------------------------------------------
    const runs = p.standups || [];
    if (!runs.length) {
      card('standup', '<span class="muted">never run live</span>', '',
           p.standups_error || 'a dry run does not count', 'tests/run_history.jsonl');
    } else {
      const last = runs[runs.length - 1];
      const greens = runs.filter(r => r.green).length;
      // THE ONE HONEST DELTA IN THIS CONSOLE. run_history.jsonl holds every
      // prior live run, so the move from the previous one is read, not
      // guessed. Nothing else here has a second measurement to compare
      // against, so nothing else gets a delta.
      let move = null;
      if (runs.length > 1) {
        const prev = runs[runs.length - 2];
        const d = (last.passed || 0) - (prev.passed || 0);
        if (d !== 0) move = { dir: d > 0 ? 'up' : 'down',
                              text: (d > 0 ? '+' : '') + d + ' vs last' };
      }
      card('standup', `${last.passed}<span class="muted">/${last.total}</span>`,
           last.green ? 'green' : 'red',
           (last.green ? 'green' : 'RED — ' + ((last.failed || []).join(', ') || '?')) +
           ' · ' + when(last.at), 'tests/run_history.jsonl', move);
      card('standups run', `${greens}<span class="muted">/${runs.length}</span>`,
           greens === runs.length ? 'green' : 'yellow',
           'green of all live runs on record', 'tests/run_history.jsonl');
    }

    // ---- parity -------------------------------------------------------
    const par = p.parity || [];
    if (!par.length) {
      card('parity', '<span class="muted">never run</span>', '',
           p.parity_error || 'chain vs bare calls, per case', 'sessions/parity_history.jsonl');
    } else {
      const lp = par[par.length - 1];
      card('parity', String(lp.mean ?? '—'), 'blue',
           `${lp.scored ?? '?'} of ${lp.cases ?? '?'} cases scored · ${when(lp.at)}`,
           'sessions/parity_history.jsonl');
    }

    // ---- the estate's own state, before its scores -------------------
    // The sittings ARE the prior record. Read exactly from sessions.jsonl,
    // where a closing line supersedes its opening one.
    const rc = p.record || {};
    let estate = '';
    if (rc.error) {
      estate = `<div class="card"><div class="eng-row eng-bad">The record could not be read:
        ${escHtml(rc.error)}<span class="brief-src">sessions/sessions.jsonl</span></div></div>`;
    } else if (rc.sittings != null) {
      const rows = (rc.recent || []).slice().reverse().map(r => `<tr>
        <td class="num">${escHtml(String(r.n))}</td>
        <td>${escHtml(String(r.started || '').replace('T', ' '))}</td>
        <td>${r.ended ? escHtml(String(r.ended).slice(11)) : '<span class="tool-bad">still open</span>'}</td>
        <td class="num">${escHtml(String(r.runs))}</td>
        <td>${r.toll_paid ? '<span class="badge badge-green">tolled</span>'
                          : '<span class="badge badge-yellow">no toll</span>'}</td></tr>`).join('');
      estate = `<div class="card">
        <div class="card-title">The estate <span class="muted">— read from the record, not counted here</span></div>
        <div class="stats">
          <div class="stat"><div class="stat-label">sittings</div>
            <div class="stat-value">${rc.sittings}</div>
            <div class="stat-note">${rc.still_open ? '<span class="tool-bad">' + rc.still_open + ' still open</span>' : 'all closed'}</div>
            <div class="brief-src">sessions/sessions.jsonl</div></div>
          <div class="stat"><div class="stat-label">tolled</div>
            <div class="stat-value ${rc.tolled === rc.sittings ? 'green' : 'yellow'}">${rc.tolled}<span class="muted">/${rc.sittings}</span></div>
            <div class="stat-note">${rc.seat_log_tolls != null ? rc.seat_log_tolls + ' written into SEAT_LOG' : 'SEAT_LOG unreadable'}</div>
            <div class="brief-src">sessions/sessions.jsonl · SEAT_LOG.md</div></div>
          <div class="stat"><div class="stat-label">runs recorded</div>
            <div class="stat-value blue">${rc.runs}</div>
            <div class="stat-note">objectives the council actually ran</div>
            <div class="brief-src">sessions/sessions.jsonl</div></div>
        </div>
        ${rows ? `<div class="table-wrap"><table><thead><tr><th class="num">sitting</th><th>opened</th><th>closed</th><th class="num">runs</th><th>toll</th></tr></thead><tbody>${rows}</tbody></table></div>` : ''}
        <div class="stat-note" style="margin-top:10px">
          ${escHtml((rc.counted_by_the_engine || []).join(' and '))} are counted by the core's own rules
          (memory.py's entry pattern; a SELECT against index/vectors.db) and are shown whole in the
          boot report in the front page's Audit Ledger tab. They are not recounted here: a second definition of "an entry"
          would drift from the core's the first time it changed.
        </div></div>`;
    }

    // THE NEWEST FIVE, not all thirty-three. This ran every live standup ever
    // recorded, oldest first, so the run that matters — the last one — sat at
    // the bottom of a table that grew a row every morning. Five, newest first,
    // matching the sittings table above it.
    const RECENT = 5;
    const shown = runs.slice(-RECENT).reverse();
    const older = runs.length - shown.length;
    const standups = runs.length ? `<div class="card"><div class="card-title">Live standups
        <span class="muted">— the last ${shown.length} of ${runs.length}, newest first</span></div>
        <div class="table-wrap"><table><thead><tr><th>when</th><th>score</th><th>failed</th><th>report</th></tr></thead><tbody>` +
        shown.map(r => `<tr>
          <td>${escHtml(when(r.at))}</td>
          <td><span class="badge ${r.green ? 'badge-green' : 'badge-red'}">${r.passed}/${r.total}</span></td>
          <td>${(r.failed || []).length ? escHtml((r.failed || []).join(', ')) : '<span class="muted">—</span>'}</td>
          <td><code>${escHtml(r.report || '')}</code></td></tr>`).join('') +
        `</tbody></table>` +
        (older ? `<div class="stat-note" style="margin-top:10px">${older} earlier
           run${older === 1 ? '' : 's'} are in <code>tests/run_history.jsonl</code>.</div>` : '') +
        `</div>` : '';

    const scores = `<div class="stats">${cards.join('')}</div>`;
    box.innerHTML = only === 'ledger' ? this.proofLedger(p)
                  : only === 'scores' ? scores
                  : only === 'estate' ? estate + standups
                  : estate + scores + standups;
  },

  // WHAT THIS BUILD HAS PROVED, AS ROWS (2026-10-06: split out of the Dashboard's deck, which retired with the
  // Dashboard on 2026-10-07). The front page's Audit Ledger tab draws them -- one reader of the proof, the way
  // paintProof's `only` places the scores and the estate.
  proofLedger(p) {
    const su = p.suites || {};
    const runs = p.standups || [];
    const last = runs.length ? runs[runs.length - 1] : null;
    const par = p.parity || [];
    const rows = [];
    const row = (k, note, v, tone, delta) => rows.push(
      `<div class="led-row">
         <div class="led-k">${escHtml(k)}<span class="led-note">${escHtml(note)}</span></div>
         <div class="led-v ${tone || ''}">${v}${delta
           ? `<span class="led-delta ${delta.dir}">${escHtml(delta.text)}</span>` : ''}</div>
       </div>`);
    const frac = (r) => `${r.passed}<span class="muted">/${r.total}</span>`;
    const none = '<span class="muted">—</span>';

    for (const name of ['strokes', 'smoke']) {
      const r = su[name];
      row(name, r ? when(r.at) : (p.suites_error || 'never run'),
          r ? frac(r) : none, r ? (r.green ? 'green' : 'red') : '');
    }

    if (last) {
      // THE ONE HONEST DELTA. run_history.jsonl holds every prior live run, so
      // the move from the previous one is read, not guessed. Nothing else here
      // has a second measurement to compare against, so nothing else gets one.
      let move = null;
      if (runs.length > 1) {
        const d = (last.passed || 0) - (runs[runs.length - 2].passed || 0);
        if (d !== 0) move = { dir: d > 0 ? 'up' : 'down', text: (d > 0 ? '+' : '') + d };
      }
      row('live standup', when(last.at), frac(last), last.green ? 'green' : 'red', move);
      const greens = runs.filter(r => r.green).length;
      row('runs green', 'of every live run on record',
          `${greens}<span class="muted">/${runs.length}</span>`,
          greens === runs.length ? 'green' : 'yellow');
    } else {
      row('live standup', p.standups_error || 'a dry run does not count', none, '');
    }

    if (par.length) {
      const lp = par[par.length - 1];
      row('parity', `${lp.scored ?? '?'} of ${lp.cases ?? '?'} cases · ${when(lp.at)}`,
          String(lp.mean ?? '—'), 'blue');
    } else {
      row('parity', p.parity_error || 'chain against bare calls', none, '');
    }

    return `<div class="ledger">
        <div class="led-head">What this build has proved</div>
        ${rows.join('')}
        <div class="led-foot"><span class="brief-src">proofs · read, never counted here</span></div>
      </div>`;
  },

  // === MESSAGES (team bridge) ===
  parseTeamHistory(text) {
    const rows = [];
    let cur = null;
    for (const line of (text || '').split('\n')) {
      let m = /^(inbound|outbound) (\S+)\/(\S+) · (\S+) · receipt ([0-9a-f]+)/.exec(line);
      if (m) { cur = { dir: m[1], platform: m[2], channel: m[3], ts: m[4], receipt: m[5], content: '' }; rows.push(cur); continue; }
      if (cur && line.startsWith('  ')) cur.content += (cur.content ? '\n' : '') + line.slice(2);
    }
    return rows;
  },

  async renderMessages(el) {
    el.innerHTML = '<div class="loading">Loading messages...</div>';
    try {
      const [hist, presence] = await Promise.all([
        API.teamHistory('', '', 100).catch(() => ({ history: '' })),
        API.teamStatus().catch(() => ({ status: '' }))
      ]);
      const rows = this.parseTeamHistory(hist.history);
      el.innerHTML = `
        <div class="page-header">
          <div>
            <div class="page-title">Messages</div>
            <div class="page-subtitle">${rows.length} crossings · receipts on every row</div>
          </div>
        </div>
        <div class="card"><div class="card-title">Bridge</div><pre id="team-presence">${escHtml(presence.status || 'bridge silent')}</pre></div>
        <div class="card mt-16"><div class="card-title">Send (guarded, receipted)</div>
          <div class="flex">
            <select id="team-platform"><option value="discord">discord</option><option value="slack">slack</option><option value="whatsapp">whatsapp</option></select>
            <input id="team-channel" type="text" placeholder="channel" value="general" />
            <input id="team-content" type="text" placeholder="Message..." style="flex:1" />
            <button class="btn btn-primary" id="team-send">Send</button>
          </div>
          <div id="team-send-status" class="muted"></div>
        </div>
        <div class="card mt-16">
          ${rows.length === 0
            ? '<div class="empty"><div class="empty-icon">&#9993;</div><div class="empty-text">The bridge holds nothing yet. Sends and deliveries appear here with receipts.</div></div>'
            : '<div class="table-wrap"><table><thead><tr><th>Dir</th><th>Platform</th><th>Channel</th><th>Content</th><th>Receipt</th><th>Time</th></tr></thead><tbody>' +
              rows.map(m => `
                <tr>
                  <td><span class="badge ${m.dir === 'outbound' ? 'badge-blue' : 'badge-muted'}">${escHtml(m.dir)}</span></td>
                  <td><span class="badge badge-${m.platform === 'discord' ? 'purple' : m.platform === 'slack' ? 'blue' : 'green'}">${escHtml(m.platform)}</span></td>
                  <td><code>${escHtml(m.channel)}</code></td>
                  <td>${escHtml((m.content || '').slice(0, 80))}${(m.content || '').length > 80 ? '...' : ''}</td>
                  <td><span class="hash">${escHtml(m.receipt.slice(0, 16))}</span></td>
                  <td>${timeAgo(m.ts)}</td>
                </tr>
              `).join('') + '</tbody></table></div>'}
        </div>
      `;
      document.getElementById('team-send').onclick = () => this.teamSend();
    } catch (e) {
      el.innerHTML = `<div class="empty"><div class="empty-icon">!</div><div class="empty-text">${escHtml(e.message)}</div></div>`;
    }
  },

  async teamSend() {
    const platform = document.getElementById('team-platform').value;
    const channel = document.getElementById('team-channel').value.trim();
    const content = document.getElementById('team-content').value.trim();
    if (!content) return;
    try {
      const r = await API.teamSend(platform, channel, content);
      document.getElementById('team-send-status').textContent = r.text;
      document.getElementById('team-content').value = '';
      this.router();
    } catch (e) { document.getElementById('team-send-status').textContent = 'Refused: ' + e.message; }
  },

  async runProve() {
    toast('Running prove...');
    try {
      const result = await API.callTool('verify_chain', { path: '.', project: 'atlas' });
      toast('Prove complete: ' + (result.output || '').slice(0, 80));
    } catch (e) { toast('Prove failed: ' + e.message, 'error'); }
  }
};

document.addEventListener('DOMContentLoaded', () => App.init());
