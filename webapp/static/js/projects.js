// Projects -- what the maker made, on the Dashboard beside the run (the maker's
// piece 2, 2026-09-21; the operator: "go on piece 2").
//
// THE LIST IS READ OFF EACH PROJECT'S OWN HISTORY, by the door's `projects`
// tool: its versions and the maker's own notes for them. Nothing here keeps a
// second account of what a project is.
//
// THE PAGE IS SHOWN SANDBOXED, BY THE HEADER IT IS SERVED UNDER
// (handlers/projects.go): an origin of its own, never this glass's, so a
// model's code cannot reach the glass with his session behind it. Its scripts
// run -- a game is a script -- and reach nothing else.
//
// THE FRAME CARRIES NO `sandbox` ATTRIBUTE, ON PURPOSE, AND MEASURED
// (2026-09-21). The first cut set one as a second wall, and the Claude app's
// own browser pane -- where he watches -- refuses any frame that carries it
// (net::ERR_BLOCKED_BY_CLIENT, even for a plain page with no header), so the
// preview came up blank there. The header alone was then measured in that
// same browser: the page framed, its script and canvas ran, its origin was
// "null", and its reach for the parent, the parent's DOM, cookies, storage and
// fetch were each refused. The header is the wall, and a stroke holds it.
//
// PICKING ONE UP AND PUTTING IT DOWN ARE WORDS. "Work on this" and "Put it
// down" put a sentence in the box and run it, exactly as if he had typed it:
// the engine answers (the maker's route), the turn is in the record like any
// other, and the conversation shows what was said. Nothing on this card writes.
//
// WHICH ONE IS IN HAND comes from the engine's own delivery (`project`), the
// only place it lives. It is kept where both browsers can read it, keyed to
// the sitting, the way the boot report and the thread are kept.
const Projects = {
  list: [],
  error: '',
  selected: '',
  version: 0,          // 0 is the page as it stands; N is version N
  inHand: '',
  session: null,
  bound: false,

  async render() {
    if (!document.getElementById('home-projects')) return;
    if (!this.bound) { Run.on((w) => this.onRun(w)); this.bound = true; }
    this.paint();
    await this.readHeld();
    await this.read();
  },

  // One read of the list. A background read: the glass answers it and keeps
  // no trace of it, as it does for the Dashboard's other reads (D1 b).
  async read() {
    try {
      const d = JSON.parse(await App.tool('projects', { action: 'list' }, true));
      this.list = Array.isArray(d.projects) ? d.projects : [];
      this.error = '';
    } catch (e) {
      this.error = e.message || 'unreadable';
    }
    this.paint();
  },

  key() { return 'maker.' + (Run.world || 'research'); },

  async readHeld() {
    this.session = Run.session || '';
    this.inHand = '';
    if (!Run.engineOpen) return;
    try {
      const r = await API.getSetting(this.key());
      const kept = JSON.parse((r && r.value) || '{}');
      // A project held by an engine that is no longer standing is not held.
      if (kept.session && kept.session === Run.session) this.inHand = kept.project || '';
    } catch { /* nothing kept yet */ }
  },

  keepHeld(project) {
    this.inHand = project;
    API.setSetting(this.key(), JSON.stringify({ session: Run.session || '', project }))
      .catch(() => {});
  },

  onRun(what) {
    if (!document.getElementById('home-projects')) return;
    if (what === 'state') {
      // A sitting opened or closed: what was in hand went with the old one.
      if ((Run.session || '') !== this.session) this.readHeld().then(() => this.paint());
      return;
    }
    if (what !== 'end' && what !== 'done') return;
    const d = Run.turn && Run.turn.delivery;
    if (d && typeof d.project === 'string') {
      this.keepHeld(d.project);
      // The project a turn made, changed or picked up is the one to show.
      if (d.project) { this.selected = d.project; this.version = 0; }
    }
    this.read();
  },

  // "Work on this" and "Put it down": words in the box, run as his own.
  say(words) {
    const input = document.getElementById('home-input');
    if (!input) return;
    input.value = words;
    Home.go();
  },

  pageUrl(name, version) {
    return API.base + '/projects/' + encodeURIComponent(name) + '/page' +
      (version ? '?v=' + version : '');
  },

  paint() {
    const box = document.getElementById('home-projects');
    if (!box) return;
    const held = this.inHand
      ? '<span class="badge badge-green">in hand: ' + escHtml(this.inHand) + '</span>'
      : '<span class="badge">none in hand</span>';
    const head = '<div class="card-header"><span class="card-title">Projects</span>' +
      '<span class="flex">' + held + ' <span class="brief-src">projects</span></span></div>';

    if (this.error) {
      box.innerHTML = head + '<div class="muted" style="padding:var(--s3) 0">' +
        'The projects could not be read: ' + escHtml(this.error) + '</div>';
      return;
    }
    if (!this.list.length) {
      box.innerHTML = head + '<div class="muted" style="padding:var(--s3) 0">' +
        'Nothing has been made yet. Ask for something in the box above — ' +
        '"make me a snake game" — and it will be here.</div>';
      return;
    }
    if (!this.list.some(p => p.name === this.selected)) {
      this.selected = this.list.some(p => p.name === this.inHand) ? this.inHand : this.list[0].name;
      this.version = 0;
    }
    const p = this.list.find(x => x.name === this.selected);
    const vs = p.versions || [];

    const rows = this.list.map(x => {
      const n = (x.versions || []).length;
      const on = x.name === this.selected;
      return '<div class="proj-row" data-name="' + escHtml(x.name) + '" style="padding:var(--s2) var(--s3);' +
        'border-radius:var(--radius);cursor:pointer;margin-bottom:var(--s1);' +
        'border:1px solid ' + (on ? 'var(--accent)' : 'var(--border)') + ';' +
        'background:' + (on ? 'var(--accent-dim)' : 'transparent') + '">' +
        '<div><b>' + escHtml(x.name) + '</b>' +
        (x.name === this.inHand ? ' <span class="badge badge-green">in hand</span>' : '') + '</div>' +
        '<div class="muted" style="font-size:var(--t-sm)">' +
        (n === 1 ? '1 version' : n + ' versions') +
        (x.updated ? ' · ' + escHtml(timeAgo(x.updated)) : '') +
        (x.error ? ' · <span style="color:var(--red)">' + escHtml(x.error) + '</span>' : '') +
        '</div></div>';
    }).join('');

    const shown = this.version || vs.length;
    const options = ['<option value="0"' + (this.version ? '' : ' selected') + '>as it stands' +
      (vs.length ? ' (version ' + vs.length + ')' : '') + '</option>']
      .concat(vs.slice().reverse().map(v => '<option value="' + v.n + '"' +
        (this.version === v.n ? ' selected' : '') + '>version ' + v.n + ' — ' +
        escHtml(v.note || '') + '</option>'));
    const act = p.name === this.inHand
      ? '<button class="btn btn-sm" type="button" id="proj-down">Put it down</button>'
      : '<button class="btn btn-sm btn-primary" type="button" id="proj-up">Work on this</button>';
    const url = this.pageUrl(p.name, this.version);

    box.innerHTML = head +
      '<div style="display:flex;gap:var(--s4);flex-wrap:wrap;align-items:flex-start">' +
      '<div style="flex:1 1 200px;min-width:0">' + rows + '</div>' +
      '<div style="flex:3 1 420px;min-width:0">' +
      '<div style="display:flex;gap:var(--s2);align-items:center;flex-wrap:wrap;margin-bottom:var(--s2)">' +
      '<select class="select" id="proj-version" style="flex:1 1 220px;min-width:0">' + options.join('') + '</select>' +
      act +
      '<a class="btn btn-sm" href="' + escHtml(url) + '" target="_blank" rel="noopener noreferrer">Open in its own tab</a>' +
      '</div>' +
      '<iframe id="proj-frame" title="' + escHtml(p.name) + ', version ' + shown + '" ' +
      'referrerpolicy="no-referrer" src="' + escHtml(url) + '" ' +
      'style="width:100%;height:440px;border:1px solid var(--border-2);border-radius:var(--radius);background:#fff"></iframe>' +
      '<div class="muted" style="font-size:var(--t-sm);margin-top:var(--s1)">' +
      'Shown sandboxed: it runs here and cannot reach this glass. It is also ' +
      '<code>projects\\' + escHtml(p.name) + '\\index.html</code>, to open in any browser.</div>' +
      '</div></div>';

    box.querySelectorAll('.proj-row').forEach(el => {
      el.onclick = () => { this.selected = el.dataset.name; this.version = 0; this.paint(); };
    });
    document.getElementById('proj-version').onchange = (e) => {
      this.version = Number(e.target.value) || 0;
      this.paint();
    };
    const up = document.getElementById('proj-up');
    if (up) up.onclick = () => this.say('work on the ' + p.name + ' project');
    const down = document.getElementById('proj-down');
    if (down) down.onclick = () => this.say('put the project down');
  }
};
