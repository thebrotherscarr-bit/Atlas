// THE CLAUDE TAB -- Claude in the Inspector (2026-10-22, his word: "set the
// model ... a little box for ... instructions for this session ... talk to you
// directly and then watch you work. In that little tab.").
//
// The glass runs the Claude Code harness in the ground (handlers/claude.go) and
// streams every event back: its words, each tool it calls, each result. Its
// hands are the door's, so anything it writes parks as a card in the dock.
//
// The model, the instructions and the conversation it continues are kept in
// the glass's own settings (/api/settings), so they are the same in every
// window. New conversation starts a fresh one.
const ClaudeTab = {
  info: null,
  busy: false,
  log: [],          // what this window has drawn: {kind, text, name?}

  async setting(key, value) {
    if (value === undefined) {
      try { return (await (await fetch('/api/settings/' + key)).json()).value || ''; } catch { return ''; }
    }
    await fetch('/api/settings/' + key, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ value }) });
  },

  async render(el) {
    this.style();
    el.innerHTML = '<div class="ag-empty">Reading...</div>';
    try { this.info = await (await fetch('/api/claude/info')).json(); } catch (e) { this.info = { error: e.message }; }
    const [model, instructions] = await Promise.all([this.setting('claude.model'), this.setting('claude.instructions')]);
    const i = this.info || {};
    const pick = model || (i.claude || [])[0] || '';
    const opt = (m) => `<option value="${escHtml(m)}"${m === pick ? ' selected' : ''}>${escHtml(m)}</option>`;
    el.innerHTML = `<div class="cl-tab">
      ${i.found === false ? '<div class="ag-card"><h4>Claude Code was not found</h4><p>Install the Claude desktop app, or set ATLAS_CLAUDE_BIN to claude.exe.</p></div>' : ''}
      <div class="cl-row">
        <label class="cl-lbl">Model</label>
        <select id="cl-model" class="input">
          <optgroup label="Claude${i.key ? ' (your key in SecondBrain)' : ' (your Claude login)'}">${(i.claude || []).map(opt).join('')}</optgroup>
          <optgroup label="Local (Ollama)">${(i.local || []).map(opt).join('')}</optgroup>
        </select>
      </div>
      <div class="cl-row cl-col">
        <label class="cl-lbl">Instructions for this session</label>
        <textarea id="cl-instr" class="input" rows="3" placeholder="Anything Claude should keep in mind while it works here...">${escHtml(instructions)}</textarea>
      </div>
      <div class="cl-bar">
        <button class="btn" id="cl-new">New conversation</button>
        <span class="cl-status" id="cl-status"></span>
      </div>
      <div class="cl-log" id="cl-log"></div>
      <div class="cl-send">
        <textarea id="cl-msg" class="input" rows="2" placeholder="Talk to Claude. Enter sends, Shift+Enter is a new line."></textarea>
        <button class="btn btn-primary" id="cl-go">Send</button>
      </div>
    </div>`;
    const $ = (id) => document.getElementById(id);
    $('cl-model').onchange = (e) => this.setting('claude.model', e.target.value);
    $('cl-instr').onchange = (e) => this.setting('claude.instructions', e.target.value);
    $('cl-new').onclick = async () => { await this.setting('claude.session', ''); this.log = []; this.paintLog(); this.status('A new conversation starts with the next message.'); };
    $('cl-go').onclick = () => this.send();
    $('cl-msg').onkeydown = (e) => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); this.send(); } };
    this.paintLog();
    this.status(this.busy ? 'Working...' : (i.running ? 'Claude is working on a message from another window.' : ''));
  },

  status(t) { const s = document.getElementById('cl-status'); if (s) s.textContent = t || ''; },

  add(entry) { this.log.push(entry); this.paintLog(); },

  paintLog() {
    const box = document.getElementById('cl-log');
    if (!box) return;
    box.innerHTML = this.log.map(e => {
      if (e.kind === 'you') return `<div class="cl-you">${escHtml(e.text)}</div>`;
      if (e.kind === 'claude') return `<div class="cl-claude">${escHtml(e.text)}</div>`;
      if (e.kind === 'tool') return `<div class="cl-tool"><span class="cl-tname">${escHtml(e.name)}</span> ${escHtml(e.text)}</div>`;
      if (e.kind === 'result') return `<details class="cl-result"><summary>${escHtml(e.text.split('\n')[0].slice(0, 140))}</summary><pre>${escHtml(e.text)}</pre></details>`;
      if (e.kind === 'done') return `<div class="cl-done">${escHtml(e.text)}</div>`;
      return `<div class="cl-err">${escHtml(e.text)}</div>`;
    }).join('');
    box.scrollTop = box.scrollHeight;
  },

  // One stream-json event, as the harness emits it, drawn as what it is.
  event(ev) {
    if (ev.type === 'system' && ev.subtype === 'init') {
      if (ev.session_id) this.setting('claude.session', ev.session_id);
      this.status(`Working on ${ev.model || 'the model'}...`);
      return;
    }
    const parts = (ev.message && Array.isArray(ev.message.content)) ? ev.message.content : [];
    if (ev.type === 'assistant') {
      for (const p of parts) {
        if (p.type === 'text' && p.text.trim()) this.add({ kind: 'claude', text: p.text });
        if (p.type === 'tool_use') this.add({ kind: 'tool', name: p.name.replace(/^mcp__/, ''), text: JSON.stringify(p.input || {}).slice(0, 300) });
      }
    } else if (ev.type === 'user') {
      for (const p of parts) {
        if (p.type !== 'tool_result') continue;
        const c = Array.isArray(p.content) ? p.content.map(x => x.text || '').join('\n') : String(p.content || '');
        this.add({ kind: 'result', text: c || '(no output)' });
      }
    } else if (ev.type === 'result') {
      if (ev.session_id) this.setting('claude.session', ev.session_id);
      const cost = ev.total_cost_usd ? ` · $${Number(ev.total_cost_usd).toFixed(3)}` : '';
      this.add({ kind: ev.is_error ? 'err' : 'done', text: `${ev.is_error ? 'Stopped' : 'Done'} in ${((ev.duration_ms || 0) / 1000).toFixed(1)} s${cost}` });
    }
  },

  async send() {
    const msg = document.getElementById('cl-msg');
    const text = (msg && msg.value || '').trim();
    if (!text || this.busy) return;
    msg.value = '';
    this.busy = true;
    this.add({ kind: 'you', text });
    this.status('Starting...');
    const go = document.getElementById('cl-go');
    if (go) { go.textContent = 'Stop'; go.onclick = () => fetch('/api/claude/stop', { method: 'POST' }); }
    try {
      const resp = await fetch('/api/claude/turn', {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          message: text,
          model: document.getElementById('cl-model')?.value || '',
          instructions: document.getElementById('cl-instr')?.value || '',
          session: await this.setting('claude.session'),
        }),
      });
      if (!resp.ok) {
        const e = await resp.json().catch(() => ({}));
        throw new Error(e.error || ('HTTP ' + resp.status));
      }
      const reader = resp.body.getReader();
      const dec = new TextDecoder();
      let buf = '';
      for (;;) {
        const { value, done } = await reader.read();
        if (done) break;
        buf += dec.decode(value, { stream: true });
        let cut;
        while ((cut = buf.indexOf('\n\n')) >= 0) {
          const block = buf.slice(0, cut);
          buf = buf.slice(cut + 2);
          const isDone = /^event: done/m.test(block);
          const data = block.split('\n').filter(l => l.startsWith('data: ')).map(l => l.slice(6)).join('\n');
          if (!data) continue;
          let ev;
          try { ev = JSON.parse(data); } catch { continue; }
          if (isDone) { if (!ev.ok && ev.error) this.add({ kind: 'err', text: ev.error }); }
          else this.event(ev);
        }
      }
    } catch (e) {
      this.add({ kind: 'err', text: e.message || 'refused' });
    } finally {
      this.busy = false;
      this.status('');
      const g = document.getElementById('cl-go');
      if (g) { g.textContent = 'Send'; g.onclick = () => this.send(); }
    }
  },

  style() {
    if (document.getElementById('cl-style')) return;
    const s = document.createElement('style');
    s.id = 'cl-style';
    s.textContent = `
      .cl-tab { display: flex; flex-direction: column; gap: 10px; height: 100%; min-height: 0; }
      .cl-row { display: flex; gap: 10px; align-items: center; }
      .cl-col { flex-direction: column; align-items: stretch; gap: 4px; }
      .cl-lbl { font-size: var(--t-xs); color: var(--text-3); text-transform: uppercase; letter-spacing: 0.5px; }
      .cl-row select { flex: 1; }
      .cl-bar { display: flex; gap: 10px; align-items: center; }
      .cl-status { font-size: var(--t-xs); color: var(--text-3); }
      .cl-log { flex: 1; min-height: 180px; overflow-y: auto; border: 1px solid var(--border); border-radius: 8px;
        padding: 10px; display: flex; flex-direction: column; gap: 8px; font-size: var(--t-md); }
      .cl-you { align-self: flex-end; max-width: 85%; background: var(--surface-2, #1c1c22); padding: 6px 10px;
        border-radius: 8px; white-space: pre-wrap; }
      .cl-claude { white-space: pre-wrap; line-height: 1.55; color: var(--text); }
      .cl-tool { font-family: var(--font-mono); font-size: var(--t-xs); color: var(--text-2); word-break: break-all; }
      .cl-tname { color: var(--blue, #60a5fa); font-weight: 600; }
      .cl-result { font-size: var(--t-xs); color: var(--text-3); }
      .cl-result pre { white-space: pre-wrap; max-height: 240px; overflow: auto; margin: 4px 0 0; }
      .cl-done { font-size: var(--t-xs); color: var(--green, #34d399); }
      .cl-err { font-size: var(--t-xs); color: var(--red, #f87171); white-space: pre-wrap; }
      .cl-send { display: flex; gap: 8px; align-items: flex-end; }
      .cl-send textarea { flex: 1; resize: vertical; }`;
    document.head.appendChild(s);
  },
};
