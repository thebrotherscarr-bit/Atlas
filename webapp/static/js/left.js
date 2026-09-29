// WHAT'S LEFT -- everything still open, on one page.
//
// His words, 2026-09-29: "a page with any leftover tasks, checklists, seams ...
// any blockers or gates, or decisions ... problems or bugs, needs to be in one
// area where we can see what is actually left" -- and, when the hand answered
// with a file: "I said write a page on the webapp".
//
// THE LIST IS ONE DOCUMENT, WHATS_LEFT.md at the root of the ground, and this
// page is how it is read. The door's `records` tool serves it, the way it
// serves every document on Records; the page keeps no copy.
//
// EVERY NUMBER HERE IS COUNTED FROM THE LINES DRAWN. The document carries no
// tally of its own, so there is no second count to drift from the first.
//
// A LINE THE PAGE CANNOT READ IS SAID, NOT DROPPED. Every open line has a
// number (A1, C10) so the work can be ordered by saying it. A list line with
// no number, a number used twice, or a number under the wrong letter is named
// in red at the top: on a page whose job is to show what is left, a line that
// silently fell out is the one fault that matters.
//
// THE STYLES RIDE WITH THE PAGE. A dozen rules nothing else uses, written into
// the head once, rather than a block at the foot of app.css.
const Left = {
  DOC: 'WHATS_LEFT.md',

  // The letters that are alarms rather than tallies: what stops the release,
  // and what waits on him.
  TONE: { A: 'red', B: 'yellow' },

  async render(el) {
    this.style();
    el.innerHTML = '<div class="loading">Reading what is left...</div>';
    let d;
    try {
      d = JSON.parse(await App.tool('records', { name: this.DOC }));
    } catch (e) {
      el.innerHTML = `
        <div class="page-header"><div>
          <div class="page-title">What's left</div>
          <div class="page-subtitle">The list could not be read.</div>
        </div></div>
        <div class="card"><div class="eng-row eng-bad">
          <code>${escHtml(this.DOC)}</code> could not be read: ${escHtml(e.message || 'refused')}
          <span class="brief-src">records</span></div></div>`;
      return;
    }
    this.paint(el, this.parse(d.text || ''), d);
  },

  // The document, cut into what it holds: a title, the words before the first
  // section, and the sections -- each a run of numbered lines, checklist
  // lines, tables and plain words. A list line runs on through every line
  // indented beneath it, so a title may break across two.
  parse(text) {
    const lines = String(text).replace(/\r\n?/g, '\n').split('\n');
    const doc = { title: '', intro: [], sections: [], faults: [] };
    const seen = {};
    let sec = null;
    const into = () => (sec ? sec.blocks : doc.intro);
    let i = 0;
    while (i < lines.length) {
      const line = lines[i];
      let m;
      if (!line.trim() || /^-{3,}\s*$/.test(line)) { i++; continue; }
      if ((m = /^# (.+)$/.exec(line))) { doc.title = m[1].trim(); i++; continue; }
      if ((m = /^## (.+)$/.exec(line))) {
        const l = /^([A-Z])\. (.+)$/.exec(m[1].trim());
        sec = { letter: l ? l[1] : '', name: l ? l[2] : m[1].trim(), blocks: [], open: 0 };
        doc.sections.push(sec);
        i++; continue;
      }
      if ((m = /^### (.+)$/.exec(line))) { into().push({ kind: 'sub', text: m[1].trim() }); i++; continue; }
      if (line.startsWith('|')) {
        const rows = [];
        while (i < lines.length && lines[i].startsWith('|')) {
          const cells = lines[i].replace(/^\||\|\s*$/g, '').split('|').map(c => c.trim());
          if (!cells.every(c => /^:?-{3,}:?$/.test(c))) rows.push(cells);
          i++;
        }
        into().push({ kind: 'table', rows });
        continue;
      }
      if (line.startsWith('- ')) {
        let body = line.slice(2);
        i++;
        while (i < lines.length && /^ {2,}\S/.test(lines[i])) { body += ' ' + lines[i].trim(); i++; }
        into().push(this.entry(body, sec, seen, doc.faults));
        continue;
      }
      let para = line.trim();
      i++;
      while (i < lines.length && lines[i].trim() && !/^(#|\||- |-{3,}\s*$)/.test(lines[i])) {
        para += ' ' + lines[i].trim(); i++;
      }
      into().push({ kind: 'para', text: para });
    }
    return doc;
  },

  // One list line. `- [ ] words` is a checklist line; `- **A1. Title.** words
  // *(source)*` is a numbered one; anything else under a lettered section is a
  // fault, because it is a line nobody can name.
  entry(body, sec, seen, faults) {
    let m = /^\[( |x|X)\] (.*)$/.exec(body);
    if (m) {
      const done = m[1] !== ' ';
      if (sec && !done) sec.open++;
      return { kind: 'check', done, text: m[2].trim() };
    }
    let src = '';
    body = body.replace(/\s*\*\(([^)]*)\)\*\s*$/, (_, s) => { src = s.trim(); return ''; });
    m = /^\*\*([A-Z]\d+)\. (.+?)\*\*\s*(.*)$/.exec(body);
    if (!m) {
      if (sec && sec.letter) faults.push(`a line under ${sec.letter} has no number: "${body.slice(0, 80)}"`);
      return { kind: 'bullet', text: body, src };
    }
    const id = m[1];
    if (seen[id]) faults.push(`${id} is used twice`);
    seen[id] = true;
    if (sec && sec.letter && id[0] !== sec.letter) faults.push(`${id} sits under ${sec.letter}`);
    if (sec) sec.open++;
    return { kind: 'item', id, title: m[2].trim(), text: m[3].trim(), src };
  },

  // Escaped first, marked up after: the document's own words never reach the
  // page as markup, only the three marks this page knows.
  inline(s) {
    return escHtml(s)
      .replace(/\*\*(.+?)\*\*/g, '<b>$1</b>')
      .replace(/`([^`]+)`/g, '<code>$1</code>')
      .replace(/(^|[\s(])\*([^*\s][^*]*?)\*(?=[\s).,;:]|$)/g, '$1<i>$2</i>');
  },

  block(b) {
    switch (b.kind) {
      case 'sub':
        return `<div class="left-sub">${this.inline(b.text)}</div>`;
      case 'item':
        return `<div class="left-row" data-row>
          <span class="badge badge-${this.TONE[b.id[0]] || 'blue'} left-id">${escHtml(b.id)}</span>
          <div><div class="left-title">${this.inline(b.title)}</div>
            ${b.text ? `<div class="left-text">${this.inline(b.text)}</div>` : ''}
            ${b.src ? `<div class="left-src">from ${escHtml(b.src)}</div>` : ''}</div></div>`;
      case 'check':
        return `<div class="left-row${b.done ? ' left-done' : ''}" data-row>
          <span class="left-box${b.done ? ' on' : ''}"></span>
          <div class="left-text">${this.inline(b.text)}</div></div>`;
      case 'bullet':
        return `<div class="left-row" data-row><span></span>
          <div><div class="left-text">${this.inline(b.text)}</div>
            ${b.src ? `<div class="left-src">from ${escHtml(b.src)}</div>` : ''}</div></div>`;
      case 'table': {
        const [head, ...rows] = b.rows;
        return `<div class="table-wrap"><table>
          <thead><tr>${(head || []).map(c => `<th>${this.inline(c)}</th>`).join('')}</tr></thead>
          <tbody>${rows.map(r => `<tr>${r.map(c => `<td>${this.inline(c)}</td>`).join('')}</tr>`).join('')}</tbody>
        </table></div>`;
      }
      default:
        return `<div class="left-text left-para">${this.inline(b.text)}</div>`;
    }
  },

  paint(el, doc, d) {
    const lettered = doc.sections.filter(s => s.letter);
    const total = lettered.reduce((a, s) => a + s.open, 0);

    const tiles = lettered.map(s => `
      <div class="stat left-tile" data-jump="left-${s.letter}" title="Go to ${escHtml(s.name)}">
        <div class="stat-label">${escHtml(s.letter + ' · ' + s.name)}</div>
        <div class="stat-value ${s.open ? (this.TONE[s.letter] || 'blue') : 'green'}">${s.open}</div>
      </div>`).join('');

    const faults = doc.faults.length ? `<div class="card card-critical">
        <div class="card-title">The list has ${doc.faults.length} line${doc.faults.length === 1 ? '' : 's'} this page cannot name</div>
        ${doc.faults.map(f => `<div class="eng-row eng-bad">${escHtml(f)}</div>`).join('')}
        <div class="stat-note">Fix them in <code>${escHtml(this.DOC)}</code>. Every open line carries one number, used once.</div>
      </div>` : '';

    const sections = doc.sections.map(s => `
      <div class="card mt-16" data-sec id="left-${escHtml(s.letter || s.name)}">
        <div class="card-header">
          <span class="card-title">${s.letter ? escHtml(s.letter) + ' · ' : ''}${this.inline(s.name)}</span>
          ${s.letter ? `<span class="badge badge-${s.open ? (this.TONE[s.letter] || 'blue') : 'green'}">${s.open} open</span>` : ''}
        </div>
        ${s.blocks.map(b => this.block(b)).join('')}
      </div>`).join('');

    el.innerHTML = `<div id="left-page">
      <div class="page-header">
        <div>
          <div class="page-title">What's left</div>
          <div class="page-subtitle">${total} still open, in one place. Say a number to order the work.</div>
        </div>
        <div class="flex">
          <div class="search-bar"><input class="input" id="left-find" placeholder="Find in the list..."></div>
          <button class="btn" id="left-again">Read it again</button>
        </div>
      </div>
      ${faults}
      <div class="stats">${tiles}</div>
      ${doc.intro.map(b => this.block(b)).join('')}
      ${sections}
      <div class="stat-note mt-16">Read from <code>${escHtml(d.name || this.DOC)}</code> ·
        ${(Number(d.bytes || 0) / 1024).toFixed(1)} KB ·
        changed ${escHtml(when(Date.parse(d.modified) || 0))}
        <span class="brief-src">records · sha256 ${escHtml(String(d.sha256 || '').slice(0, 16))}</span></div>
    </div>`;

    document.getElementById('left-again').onclick = () => this.render(el);
    document.getElementById('left-find').oninput = (e) => this.find(e.target.value);
    el.querySelectorAll('[data-jump]').forEach(t => {
      t.onclick = () => {
        const to = document.getElementById(t.dataset.jump);
        if (to) to.scrollIntoView({ behavior: 'smooth', block: 'start' });
      };
    });
  },

  // A section with lines and none of them showing is put away with them; a
  // section that holds only words stays, because it was never filtered.
  find(q) {
    q = (q || '').trim().toLowerCase();
    document.querySelectorAll('#left-page [data-row]').forEach(r => {
      r.hidden = q !== '' && r.textContent.toLowerCase().indexOf(q) < 0;
    });
    document.querySelectorAll('#left-page [data-sec]').forEach(c => {
      const rows = [...c.querySelectorAll('[data-row]')];
      c.hidden = q !== '' && rows.length > 0 && rows.every(r => r.hidden);
    });
  },

  style() {
    if (document.getElementById('left-style')) return;
    const s = document.createElement('style');
    s.id = 'left-style';
    s.textContent = `
      .left-row { display: grid; grid-template-columns: 56px 1fr; gap: 12px;
        align-items: start; padding: 10px 0; border-bottom: 1px solid var(--border); }
      .left-row:last-child { border-bottom: 0; }
      .left-row[hidden] { display: none; }
      .left-id { justify-self: start; font-family: var(--font-mono); }
      .left-title { font-weight: 600; color: var(--text); font-size: var(--t-lg); line-height: 1.4; }
      .left-text { color: var(--text-2); font-size: var(--t-md); line-height: 1.6; margin-top: 2px; }
      .left-para { margin: 10px 0; }
      .left-src { font-size: var(--t-xs); color: var(--text-3); margin-top: 4px; }
      .left-sub { font-size: var(--t-xs); text-transform: uppercase; letter-spacing: 0.5px;
        color: var(--text-3); font-weight: 600; margin: 16px 0 2px; }
      .left-box { width: 14px; height: 14px; margin-top: 4px; border: 1px solid var(--border-2);
        border-radius: 3px; }
      .left-box.on { background: var(--green); border-color: var(--green); }
      .left-done .left-text { text-decoration: line-through; color: var(--text-3); }
      .left-tile { cursor: pointer; }`;
    document.head.appendChild(s);
  },
};
