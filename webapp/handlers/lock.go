// THE LOCK (2026-09-21). His words, in order: "Simple login system for now,
// user/pin to start"; "multi-roles, all working under a single user"; "make the
// thing at least semi-secure"; and, asked who may open the glass, "This PC only".
//
// ONE USER, ONE PIN. The first time the glass opens it asks the person at this
// computer for a name and a PIN of 4 to 8 digits; after that it is a lock
// screen. A PIN session is the OPERATOR's: it names no tenant, so it sees and
// names every project exactly as the glass did with the gate off (session.go,
// createUser). The roles that will work under that one user are the next
// piece, not this one.
//
// SEMI-SECURE, AND WHAT THAT MEANS, said plainly rather than implied:
//   - the PIN is never stored: PBKDF2-SHA256 over a random salt, 600,000
//     rounds, in the glass's own data folder (user.json, 0600);
//   - five wrong PINs in a row close the lock for a minute;
//   - setup is taken only from this computer, and only while nobody is set up;
//   - the glass listens on this computer only (server.go, addr).
//
// What it is NOT: a PIN is short, so anyone able to READ user.json could guess
// it offline in minutes. That file sits on the operator's own disk beside the
// glass's database, and a person holding that disk holds the machine already.
// FORGOT THE PIN: delete webapp\data\user.json and reload -- the glass asks
// for a new one. Nothing else is lost; the record lives elsewhere.
package handlers

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

// The numbers the lock keeps. Named, so the page and the strokes read the same
// values the handlers enforce.
const (
	PINMin    = 4
	PINMax    = 8
	pinRounds = 600_000
	pinTries  = 5
	pinWait   = time.Minute
	nameMax   = 40
)

// lockFile is what user.json holds. No PIN, only what checks one.
type lockFile struct {
	Name    string `json:"name"`
	Salt    string `json:"salt"`
	Hash    string `json:"hash"`
	Rounds  int    `json:"rounds"`
	Created string `json:"created"`
}

// lockBox is the lock's state: where the user lives and how many wrong PINs
// have come in a row. The count lives in memory on purpose -- a restart forgets
// it, and restarting the glass needs this computer already.
type lockBox struct {
	mu    sync.Mutex
	path  string
	fails int
	until time.Time
	now   func() time.Time
}

// ConfigureLock names the file the one user lives in. main.go calls it beside
// ConfigureAuth; with no lock configured the three lock faces say so.
func (h *Handlers) ConfigureLock(path string) {
	if path == "" {
		path = "data/user.json"
	}
	h.lock = &lockBox{path: path, now: time.Now}
}

// validPIN is the one rule: 4 to 8 digits, nothing else.
func validPIN(pin string) error {
	if len(pin) < PINMin || len(pin) > PINMax {
		return fmt.Errorf("a PIN is %d to %d digits", PINMin, PINMax)
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return fmt.Errorf("a PIN is digits only")
		}
	}
	return nil
}

// cleanName keeps a name a person typed as a person would read it: spaces
// collapsed, control characters refused, and short enough for a lock screen.
func cleanName(s string) (string, error) {
	name := strings.Join(strings.Fields(s), " ")
	if name == "" {
		return "", fmt.Errorf("a name is needed")
	}
	if len([]rune(name)) > nameMax {
		return "", fmt.Errorf("a name is at most %d characters", nameMax)
	}
	for _, c := range name {
		if unicode.IsControl(c) {
			return "", fmt.Errorf("a name is letters, numbers and spaces")
		}
	}
	return name, nil
}

func hashPIN(pin string, salt []byte, rounds int) ([]byte, error) {
	return pbkdf2.Key(sha256.New, pin, salt, rounds, 32)
}

// read answers (nil, nil) when nobody is set up -- the one state setup is
// taken in. A file that is there and cannot be read is an ERROR, never "nobody
// is set up": treating a damaged file as absent would hand the lock to the
// next person who opened the page.
func (b *lockBox) read() (*lockFile, error) {
	raw, err := os.ReadFile(b.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var f lockFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("the user file is damaged (%v)", err)
	}
	if f.Name == "" || f.Salt == "" || f.Hash == "" || f.Rounds <= 0 {
		return nil, fmt.Errorf("the user file is incomplete")
	}
	return &f, nil
}

// write lands the whole file at once: written beside, then renamed over, so a
// crash mid-write leaves the old file or none -- never half of one.
func (b *lockBox) write(f lockFile) error {
	if err := os.MkdirAll(filepath.Dir(b.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	tmp := b.path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, b.path)
}

// wait is the whole seconds left on a closed lock, rounded up, 0 when it is
// open. Caller holds mu.
func (b *lockBox) wait() int {
	left := b.until.Sub(b.now())
	if left <= 0 {
		return 0
	}
	secs := int(left / time.Second)
	if left%time.Second != 0 {
		secs++
	}
	return secs
}

// fromThisComputer is judged from the connection, never from a header a
// caller could write.
func fromThisComputer(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(strings.Trim(host, "[]"))
	return ip != nil && ip.IsLoopback()
}

func jsonCode(w http.ResponseWriter, code int, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

// LockState answers the page's first question before anyone is signed in:
// is someone set up here, and what name does the lock screen greet. OPEN, and
// it says nothing a stranger could use -- the only person who can reach it is
// at this computer (server.go).
func (h *Handlers) LockState(w http.ResponseWriter, r *http.Request) {
	if h.lock == nil {
		jsonErr(w, 503, "the lock is not configured")
		return
	}
	h.lock.mu.Lock()
	f, err := h.lock.read()
	wait := h.lock.wait()
	h.lock.mu.Unlock()
	sess, ok := h.sessionOf(r)
	out := map[string]any{
		"setup":     f == nil && err == nil,
		"name":      "",
		"wait":      wait,
		"signed_in": ok && sess != nil,
		"pin_min":   PINMin,
		"pin_max":   PINMax,
	}
	if f != nil {
		out["name"] = f.Name
	}
	if err != nil {
		out["damaged"] = err.Error()
	}
	jsonResp(w, out)
}

// LockSetup takes the one user, once, from this computer, and signs them in.
func (h *Handlers) LockSetup(w http.ResponseWriter, r *http.Request) {
	if h.lock == nil || h.sessions == nil {
		jsonErr(w, 503, "sign-in is not switched on")
		return
	}
	if !fromThisComputer(r) {
		jsonErr(w, 403, "setup is only taken from this computer")
		return
	}
	var body struct {
		Name string `json:"name"`
		PIN  string `json:"pin"`
	}
	if err := json.NewDecoder(lockBody(r)).Decode(&body); err != nil {
		jsonErr(w, 400, "invalid json")
		return
	}
	name, err := cleanName(body.Name)
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	if err := validPIN(body.PIN); err != nil {
		jsonErr(w, 400, err.Error())
		return
	}

	h.lock.mu.Lock()
	defer h.lock.mu.Unlock()
	// ONCE. Held under the lock's own mutex, so two setups racing cannot both
	// find nobody set up.
	if f, err := h.lock.read(); f != nil || err != nil {
		jsonErr(w, 409, "someone is already set up here -- sign in instead")
		return
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		jsonErr(w, 500, "no randomness for the salt")
		return
	}
	key, err := hashPIN(body.PIN, salt, pinRounds)
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	if err := h.lock.write(lockFile{
		Name: name, Salt: hex.EncodeToString(salt), Hash: hex.EncodeToString(key),
		Rounds: pinRounds, Created: time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		jsonErr(w, 500, "the user could not be saved: "+err.Error())
		return
	}
	sess := h.sessions.createUser(name)
	setSessionCookie(w, sess.Token)
	jsonResp(w, map[string]any{"status": "ok", "name": name})
}

// Unlock checks a PIN against the one user. Wrong ones are counted; the fifth
// in a row closes the lock for a minute, and the right PIN waits it out too.
func (h *Handlers) Unlock(w http.ResponseWriter, r *http.Request) {
	if h.lock == nil || h.sessions == nil {
		jsonErr(w, 503, "sign-in is not switched on")
		return
	}
	var body struct {
		PIN string `json:"pin"`
	}
	if err := json.NewDecoder(lockBody(r)).Decode(&body); err != nil {
		jsonErr(w, 400, "invalid json")
		return
	}

	h.lock.mu.Lock()
	defer h.lock.mu.Unlock()
	if wait := h.lock.wait(); wait > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(wait))
		jsonCode(w, 429, map[string]any{
			"error": fmt.Sprintf("too many tries -- wait %d seconds", wait), "wait": wait})
		return
	}
	f, err := h.lock.read()
	if err != nil {
		jsonErr(w, 500, err.Error()+" -- delete webapp\\data\\user.json to set up again")
		return
	}
	if f == nil {
		jsonErr(w, 409, "nobody is set up yet")
		return
	}
	matched := false
	if validPIN(body.PIN) == nil {
		salt, err1 := hex.DecodeString(f.Salt)
		want, err2 := hex.DecodeString(f.Hash)
		if err1 == nil && err2 == nil {
			got, err := hashPIN(body.PIN, salt, f.Rounds)
			matched = err == nil && subtle.ConstantTimeCompare(got, want) == 1
		}
	}
	if !matched {
		h.lock.fails++
		if h.lock.fails >= pinTries {
			h.lock.fails = 0
			h.lock.until = h.lock.now().Add(pinWait)
			wait := h.lock.wait()
			w.Header().Set("Retry-After", strconv.Itoa(wait))
			jsonCode(w, 429, map[string]any{
				"error": fmt.Sprintf("too many tries -- wait %d seconds", wait), "wait": wait})
			return
		}
		left := pinTries - h.lock.fails
		jsonCode(w, 401, map[string]any{
			"error": fmt.Sprintf("that PIN is not it -- %d tries left", left), "tries_left": left})
		return
	}
	h.lock.fails = 0
	sess := h.sessions.createUser(f.Name)
	setSessionCookie(w, sess.Token)
	jsonResp(w, map[string]any{"status": "ok", "name": f.Name})
}

// lockBody bounds what the lock reads: a name and a PIN are a few dozen bytes,
// and nothing sent here should be read whole.
func lockBody(r *http.Request) io.Reader {
	return io.LimitReader(r.Body, 4<<10)
}
