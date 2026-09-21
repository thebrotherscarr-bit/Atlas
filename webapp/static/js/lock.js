// THE LOCK SCREEN (2026-09-21). One user, one PIN, this PC only -- the glass's
// front door for a person rather than a key (handlers/lock.go has the rest).
//
// IT DRAWS OVER EVERYTHING. The page beneath may already have asked the server
// for things and been refused; none of that matters until the lock opens. And
// opening it RELOADS the page, so every face starts over signed in instead of
// staying half-drawn from the refusals it got while locked.
const Lock = {
  shown: false,
  timer: null,

  // On load: signed in, or locked? /api/me answers 401 when the lock is shut.
  async boot() {
    try {
      const r = await fetch('/api/me');
      if (r.status === 401) { this.show(); return; }
      this.paintWho(await r.json());
    } catch { /* the glass itself is unreachable; app.js says so */ }
  },

  // The sidebar says who is signed in, beside the button that locks it again.
  paintWho(me) {
    const line = document.getElementById('lock-line');
    const who = document.getElementById('lock-who');
    if (!line || !who || !me || !me.user) return;
    who.textContent = me.user;
    line.hidden = false;
  },

  async show() {
    if (this.shown) return;
    this.shown = true;
    let st = {};
    try { st = await (await fetch('/api/lock')).json(); } catch {}
    const box = document.createElement('div');
    box.id = 'lock';
    box.setAttribute('role', 'dialog');
    box.setAttribute('aria-modal', 'true');
    box.style.cssText = 'position:fixed;inset:0;z-index:10000;display:flex;' +
      'align-items:center;justify-content:center;background:var(--bg);' +
      'font-family:var(--font-prose);color:var(--text)';
    box.innerHTML = this.card(st);
    document.body.appendChild(box);
    this.wire(st);
  },

  card(st) {
    const min = st.pin_min || 4, max = st.pin_max || 8;
    const label = (id, text) =>
      `<label for="${id}" style="display:block;margin:14px 0 6px;color:var(--text-2);font-size:var(--t-sm)">${text}</label>`;
    const field = 'width:100%;box-sizing:border-box;padding:10px 12px;font-size:var(--t-lg);' +
      'background:var(--bg-3);color:var(--text);border:1px solid var(--border-2);border-radius:var(--radius)';
    const pin = (id, text) => label(id, text) +
      `<input id="${id}" type="password" inputmode="numeric" autocomplete="off" pattern="[0-9]*"
        minlength="${min}" maxlength="${max}" style="${field};letter-spacing:6px">`;
    const shell = (title, words, body, button) => `
      <form id="lock-form" autocomplete="off" style="width:340px;max-width:calc(100vw - 32px);padding:28px;
        background:var(--bg-2);border:1px solid var(--border);border-radius:var(--radius-lg);box-shadow:var(--shadow)">
        <div style="font-family:var(--font);color:var(--accent);letter-spacing:3px;font-size:var(--t-sm)">ATLAS</div>
        <div style="font-size:var(--t-2xl);margin:6px 0 8px">${title}</div>
        <div style="color:var(--text-2);font-size:var(--t-base);line-height:1.5">${words}</div>
        ${body}
        <div id="lock-say" role="status" style="min-height:20px;margin-top:12px;color:var(--red);font-size:var(--t-sm)"></div>
        ${button ? `<button id="lock-go" type="submit" class="btn btn-primary" style="width:100%;margin-top:6px;padding:10px">${button}</button>` : ''}
      </form>`;
    if (st.damaged) {
      return shell('The sign-in file needs a reset',
        `It could not be read (${escHtml(st.damaged)}). Delete <code>atlas\\webapp\\data\\user.json</code>
         and reload this page to pick a name and PIN again. Nothing else is lost.`, '', '');
    }
    if (st.setup) {
      return shell('Welcome',
        `This dashboard opens with a PIN. Pick your name and a PIN of ${min} to ${max} numbers --
         you type the PIN each time you sit down.`,
        label('lock-name', 'Your name') +
        `<input id="lock-name" type="text" maxlength="40" autocomplete="off" style="${field}">` +
        pin('lock-pin', 'PIN') + pin('lock-pin2', 'PIN again'), 'Start');
    }
    return shell(`Hi, ${escHtml(st.name || 'there')}`,
      'Type your PIN to open the dashboard.', pin('lock-pin', 'PIN'), 'Unlock');
  },

  say(text) {
    const el = document.getElementById('lock-say');
    if (el) el.textContent = text || '';
  },

  // Disabled AND seen to be: the console's button style has no disabled look,
  // and a person waiting out the minute should not have to guess.
  button(enabled) {
    const btn = document.getElementById('lock-go');
    if (!btn) return;
    btn.disabled = !enabled;
    btn.style.opacity = enabled ? '' : '0.45';
    btn.style.cursor = enabled ? '' : 'not-allowed';
  },

  // A closed lock counts down where the person is looking, and the button
  // waits with it, so nobody spends the minute pressing it.
  countdown(seconds) {
    clearInterval(this.timer);
    let left = seconds;
    const tick = () => {
      if (left <= 0) {
        clearInterval(this.timer);
        this.timer = null;
        this.say('');
        this.button(true);
        return;
      }
      this.say(`Too many tries. Try again in ${left} seconds.`);
      this.button(false);
      left -= 1;
    };
    tick();
    this.timer = setInterval(tick, 1000);
  },

  wire(st) {
    const form = document.getElementById('lock-form');
    if (!form) return;
    const first = document.getElementById(st.setup ? 'lock-name' : 'lock-pin');
    if (first) first.focus();
    if (st.wait) this.countdown(st.wait);
    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      if (this.timer) return;
      const value = (id) => (document.getElementById(id) || {}).value || '';
      const pin = value('lock-pin');
      let url = '/api/unlock', body = { pin };
      if (st.setup) {
        if (!value('lock-name').trim()) { this.say('Type a name first.'); return; }
        if (pin !== value('lock-pin2')) { this.say('The two PINs are not the same.'); return; }
        url = '/api/setup';
        body = { name: value('lock-name'), pin };
      }
      this.button(false);
      try {
        const r = await fetch(url, {
          method: 'POST', headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body)
        });
        const d = await r.json().catch(() => ({}));
        if (r.ok) { location.reload(); return; }
        if (r.status === 429 && d.wait) { this.countdown(d.wait); return; }
        this.say(d.error || 'That did not work.');
        const box = document.getElementById('lock-pin');
        if (box && !st.setup) { box.value = ''; box.focus(); }
      } catch {
        this.say('The dashboard did not answer. Is it still running?');
      }
      this.button(true);
    });
  },

  // The sidebar's Lock button: the session ends here and on the server.
  async lock() {
    try { await fetch('/api/logout', { method: 'POST' }); } catch {}
    location.reload();
  }
};

document.addEventListener('DOMContentLoaded', () => Lock.boot());
