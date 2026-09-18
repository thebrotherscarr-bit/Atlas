package db

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TraceWindow is how many of the newest traces the glass holds in memory.
//
// THE WHOLE STORE WAS REWRITTEN TO ADD ONE TRACE, UNTIL 2026-09-14. Every tool
// call the panel makes becomes a trace carrying the tool's full output, and
// AddTrace ended in Save(), which marshalled every trace ever recorded --
// indented -- into store.json and renamed it into place. Measured that
// afternoon: 29,156 traces in a 286 MB store.json, rewritten on every call;
// the process at 17.3 GB and ~7,300 CPU-seconds since a 13:21 start; tool
// calls taking 13-25 s each, slow enough that the Dashboard's Close the
// sitting did not reach the door for nine minutes.
//
// Now a trace is ONE LINE appended to traces.jsonl, and nothing is rewritten
// to add one. The ledger keeps every trace (fold, never delete). Memory keeps
// the newest TraceWindow, more than any reader asks for -- ListTraces 100, an
// agent's page 50, search 500 -- and a trace older than the window is still
// found by its id, from the ledger.
const TraceWindow = 1000

const (
	storeName  = "store.json"
	ledgerName = "traces.jsonl"
)

type DB struct {
	dir    string
	mu     sync.RWMutex
	traces []Trace // the newest TraceWindow, oldest first
	total  int     // every trace the ledger holds
	ragged bool    // the ledger's last line has no newline: a write died midway
	evals  []Eval
	agents []Agent
	msgs   []Message
	keys   map[string]string
	saveMu sync.Mutex // one Save at a time, from its snapshot to its rename
}

type Trace struct {
	ID         string    `json:"id"`
	AgentID    string    `json:"agent_id"`
	Tool       string    `json:"tool"`
	Input      string    `json:"input"`
	Output     string    `json:"output"`
	Hash       string    `json:"hash"`
	Status     string    `json:"status"`
	DurationMs int64     `json:"duration_ms"`
	CreatedAt  time.Time `json:"created_at"`
	Tenant     string    `json:"tenant,omitempty"`
}

type Eval struct {
	ID        string    `json:"id"`
	TraceID   string    `json:"trace_id"`
	Name      string    `json:"name"`
	Score     float64   `json:"score"`
	Passed    bool      `json:"passed"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
	Tenant    string    `json:"tenant,omitempty"`
}

type Agent struct {
	ID          string `json:"id"`
	Office      string `json:"office"`
	ReportsTo   string `json:"reports_to"`
	Role        string `json:"role"`
	Mode        string `json:"mode"`
	Permissions string `json:"permissions"`
	USContent   string `json:"us_content"`
	Tenant      string `json:"tenant,omitempty"`
}

type Message struct {
	ID        string    `json:"id"`
	Channel   string    `json:"channel"`
	Platform  string    `json:"platform"`
	AgentID   string    `json:"agent_id"`
	Content   string    `json:"content"`
	Direction string    `json:"direction"`
	CreatedAt time.Time `json:"created_at"`
	Tenant    string    `json:"tenant,omitempty"`
}

// fileData is store.json. Traces is READ, never written: a store saved before
// the ledger carries them inline, and load() folds them out.
type fileData struct {
	Traces []Trace           `json:"traces,omitempty"`
	Evals  []Eval            `json:"evals"`
	Agents []Agent           `json:"agents"`
	Msgs   []Message         `json:"messages"`
	Keys   map[string]string `json:"keys"`
}

// Open refuses only when a store's inline traces cannot be moved into the
// ledger: starting anyway would mean the next Save() wrote store.json without
// them.
func Open(dir string) (*DB, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	d := &DB{
		dir:  dir,
		keys: make(map[string]string),
	}
	if err := d.load(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DB) load() error {
	data, err := os.ReadFile(filepath.Join(d.dir, storeName))
	if err == nil {
		var fd fileData
		if json.Unmarshal(data, &fd) == nil {
			d.evals = fd.Evals
			d.agents = fd.Agents
			d.msgs = fd.Msgs
			if fd.Keys != nil {
				d.keys = fd.Keys
			}
			if len(fd.Traces) > 0 {
				if err := d.fold(fd.Traces); err != nil {
					return fmt.Errorf("the store's %d traces could not be moved into %s, "+
						"so store.json was left as it is: %w", len(fd.Traces), ledgerName, err)
				}
				if err := d.Save(); err != nil {
					return fmt.Errorf("the traces are in %s, and store.json could not be "+
						"rewritten without them: %w", ledgerName, err)
				}
			}
		}
	}
	return d.readLedger()
}

// fold moves the traces a store.json carried inline into the ledger, oldest
// first, and returns only once every one of them is there.
//
// IT IS SAFE TO RUN TWICE. The new ledger is written beside the old one -- the
// store's traces, then every line the ledger already held that is not one of
// them -- and renamed over it in one step. A fold that dies midway leaves the
// old ledger and store.json as they were; a start that dies after the rename
// and before store.json is rewritten folds again next time, and no trace is
// written twice.
func (d *DB) fold(inline []Trace) error {
	ledger := filepath.Join(d.dir, ledgerName)
	tmp := ledger + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	fail := func(err error) error {
		f.Close()
		os.Remove(tmp)
		return err
	}
	w := bufio.NewWriter(f)
	ids := make(map[string]bool, len(inline))
	for _, t := range inline {
		line, err := json.Marshal(t)
		if err != nil {
			return fail(err)
		}
		w.Write(line)
		w.WriteByte('\n')
		if t.ID != "" {
			ids[t.ID] = true
		}
	}
	_, err = eachLine(ledger, func(line []byte) {
		var head struct {
			ID string `json:"id"`
		}
		if json.Unmarshal(line, &head) == nil && head.ID != "" && ids[head.ID] {
			return
		}
		w.Write(line)
		w.WriteByte('\n')
	})
	if err != nil && !os.IsNotExist(err) {
		return fail(err)
	}
	if err := w.Flush(); err != nil {
		return fail(err)
	}
	if err := f.Sync(); err != nil {
		return fail(err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, ledger); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// readLedger counts every trace in the ledger and holds the newest
// TraceWindow. Only the lines kept are decoded, so a start reads a long ledger
// without holding it.
func (d *DB) readLedger() error {
	ring := make([][]byte, TraceWindow)
	next, total := 0, 0
	ragged, err := eachLine(filepath.Join(d.dir, ledgerName), func(line []byte) {
		if !whole(line) {
			return
		}
		ring[next%TraceWindow] = line
		next++
		total++
	})
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	kept := next
	if kept > TraceWindow {
		kept = TraceWindow
	}
	d.traces = make([]Trace, 0, kept)
	for i := next - kept; i < next; i++ {
		var t Trace
		if json.Unmarshal(ring[i%TraceWindow], &t) == nil {
			d.traces = append(d.traces, t)
		}
	}
	d.total = total
	d.ragged = ragged
	return nil
}

// eachLine hands fn every complete line of a file, newline removed. A last
// line with no newline is a write that died midway: it is not handed on, and
// the first result says it was there.
func eachLine(path string, fn func(line []byte)) (ragged bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	r := bufio.NewReaderSize(f, 1<<20)
	for {
		line, err := r.ReadBytes('\n')
		if err == io.EOF {
			return len(line) > 0, nil
		}
		if err != nil {
			return false, err
		}
		line = bytes.TrimRight(line, "\r\n")
		if len(line) > 0 {
			fn(line)
		}
	}
}

// whole is the cheap test for a line that is a trace and not the fragment a
// torn write left behind: one JSON object, opened and closed.
func whole(line []byte) bool {
	return len(line) > 1 && line[0] == '{' && line[len(line)-1] == '}'
}

// Save writes store.json: evals, agents, messages and keys.
//
// ONE AT A TIME, AND MARSHALLED UNDER THE LOCK (2026-09-15). It took the read
// lock only to copy the four fields and marshalled after releasing it -- but a
// map and three slices copy as references to the same data. A SetKey landing
// mid-marshal was either Go's fatal "concurrent map iteration and map write",
// which ends the process, or a panic in the JSON encoder when the map grew under
// it; an UpsertAgent was a torn read. And every Save wrote the one
// store.json.tmp: sixteen at once failed 265 to 277 times in 320, on that
// file's open or its rename, while every caller of Save discards the error.
// Now each Save holds saveMu from its snapshot to its rename, and takes the
// snapshot inside it, so the Save that renames last wrote the newest state.
func (d *DB) Save() error {
	d.saveMu.Lock()
	defer d.saveMu.Unlock()
	d.mu.RLock()
	data, err := json.MarshalIndent(fileData{
		Evals:  d.evals,
		Agents: d.agents,
		Msgs:   d.msgs,
		Keys:   d.keys,
	}, "", "  ")
	d.mu.RUnlock()
	if err != nil {
		return err
	}
	tmp := filepath.Join(d.dir, storeName+".tmp")
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(d.dir, storeName))
}

// AddTrace appends one line to the ledger and nothing else is written. A line
// the disk refused is said on stderr: the trace is still shown, and it will not
// be there after a restart.
func (d *DB) AddTrace(t Trace) {
	line, err := json.Marshal(t)
	d.mu.Lock()
	defer d.mu.Unlock()
	if err == nil {
		err = d.appendLine(line)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "trace %s not written to %s: %v\n", t.ID, ledgerName, err)
	} else {
		d.total++
	}
	if len(d.traces) < TraceWindow {
		d.traces = append(d.traces, t)
		return
	}
	copy(d.traces, d.traces[1:])
	d.traces[len(d.traces)-1] = t
}

// appendLine writes one line at the end of the ledger. Held under d.mu, so two
// traces never interleave. After a torn last line the new line starts on a
// line of its own, rather than being glued to the fragment and lost with it.
func (d *DB) appendLine(line []byte) error {
	f, err := os.OpenFile(filepath.Join(d.dir, ledgerName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	buf := make([]byte, 0, len(line)+2)
	if d.ragged {
		buf = append(buf, '\n')
	}
	buf = append(buf, line...)
	buf = append(buf, '\n')
	_, werr := f.Write(buf)
	cerr := f.Close()
	if werr != nil {
		return werr
	}
	if cerr != nil {
		return cerr
	}
	d.ragged = false
	return nil
}

// TraceCount is every trace the ledger holds, not the window.
func (d *DB) TraceCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.total
}

// GetTraces reads the window, newest first. No reader asks for more than it
// holds.
func (d *DB) GetTraces(agentID string, limit int) []Trace {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var result []Trace
	for i := len(d.traces) - 1; i >= 0; i-- {
		if agentID == "" || d.traces[i].AgentID == agentID {
			result = append(result, d.traces[i])
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result
}

// GetTrace hands back a copy: the window moves under a pointer into it. A trace
// older than the window is read from the ledger, so an eval that names an old
// trace still opens it.
func (d *DB) GetTrace(id string) *Trace {
	d.mu.RLock()
	for i := len(d.traces) - 1; i >= 0; i-- {
		if d.traces[i].ID == id {
			t := d.traces[i]
			d.mu.RUnlock()
			return &t
		}
	}
	d.mu.RUnlock()
	if id == "" {
		return nil
	}
	quoted, err := json.Marshal(id)
	if err != nil {
		return nil
	}
	needle := append([]byte(`"id":`), quoted...)
	var found *Trace
	eachLine(filepath.Join(d.dir, ledgerName), func(line []byte) {
		if found != nil || !bytes.Contains(line, needle) {
			return
		}
		var t Trace
		if json.Unmarshal(line, &t) == nil && t.ID == id {
			found = &t
		}
	})
	return found
}

func (d *DB) AddEval(e Eval) {
	d.mu.Lock()
	d.evals = append(d.evals, e)
	d.mu.Unlock()
	d.Save()
}

func (d *DB) GetEvals(traceID string) []Eval {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var result []Eval
	for _, e := range d.evals {
		if traceID == "" || e.TraceID == traceID {
			result = append(result, e)
		}
	}
	return result
}

func (d *DB) UpsertAgent(a Agent) {
	d.mu.Lock()
	for i := range d.agents {
		if d.agents[i].ID == a.ID {
			d.agents[i] = a
			d.mu.Unlock()
			d.Save()
			return
		}
	}
	d.agents = append(d.agents, a)
	d.mu.Unlock()
	d.Save()
}

func (d *DB) GetAgents() []Agent {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]Agent, len(d.agents))
	copy(out, d.agents)
	return out
}

func (d *DB) GetAgent(id string) *Agent {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for i := range d.agents {
		if d.agents[i].ID == id {
			return &d.agents[i]
		}
	}
	return nil
}

func (d *DB) AddMessage(m Message) {
	d.mu.Lock()
	d.msgs = append(d.msgs, m)
	d.mu.Unlock()
	d.Save()
}

func (d *DB) GetMessages(channel string, limit int) []Message {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var result []Message
	for i := len(d.msgs) - 1; i >= 0; i-- {
		if channel == "" || d.msgs[i].Channel == channel {
			result = append(result, d.msgs[i])
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result
}

func (d *DB) SetKey(key, value string) {
	d.mu.Lock()
	d.keys[key] = value
	d.mu.Unlock()
	d.Save()
}

func (d *DB) GetKey(key string) string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.keys[key]
}

func (d *DB) Close() error {
	return d.Save()
}
