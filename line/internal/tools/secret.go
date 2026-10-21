package tools

// THE SECRET CARD (2026-10-21). Steward asks the operator for a key -- a bot token, an API key -- with a needs_answer
// that names the .env variable (`"secret": "DISCORD_BOT_TOKEN"`); the glass paints a MASKED field for it, and what he
// pastes there goes from the glass straight to this tool, which writes NAME=value into the world's own .env. The glass
// then answers the engine's question with the word `set`, never the value. So the value never crosses the engine, a
// model, the transcript, the wire's events or a log: RULE 7, keys are silent.
//
// THE OPERATOR'S HAND ALONE. ServiceOnly: Registry.Call refuses every caller but his glass before RBAC or the holds look
// at it, and the tool refuses for itself too, so no seat, agent or replayed hold can write a key. It is not parked for a
// card: the paste into the masked field on his card IS his act.
//
// WHAT IT SAYS AND WHAT IT KEEPS. Its answer, and the one line it writes to the world's holds record, say only
// "secret NAME set in .env (N chars)". The value is never in a refusal, an error, a record or an answer; the file's
// other lines are read only to be written back as they were, and never returned. The argument is declared secret
// (Tool.Secrets), so anywhere a call of it were SHOWN it would read as Withheld.
//
// ITS RULES. The name is ^[A-Z][A-Z0-9_]{1,63}$. The value is one line: no newline, carriage return or NUL, not empty,
// at most 4096 characters. An existing `NAME=` line (also `export NAME=`) is replaced in place, every one of them; with
// none, `NAME=value` is appended. The file keeps its own line endings and its byte-order mark, and is created (LF) when
// it is absent. A .env that mixes two line endings, is not UTF-8 text, or is not a plain file is refused, and nothing is
// written. It is not refused while a sitting is open: the question it answers is asked from inside one, and .env is
// read by the engine and indexed by nothing (RULE 7).

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"atlas/line/internal/tenant"
)

// SecretNameRe is the shape a secret's name must have: the .env variable the engine's needs_answer names.
var SecretNameRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]{1,63}$`)

// SecretMaxChars is the longest value secret_set writes.
const SecretMaxChars = 4096

// secretLock serialises the writes to a .env: two keys pasted at once must not each write back the file the other
// had not yet changed.
var secretLock sync.Mutex

// secretRefuse is a refusal that names why and never the value.
func secretRefuse(format string, a ...any) (string, error) {
	return "", fmt.Errorf("refused: "+format+": nothing was written", a...)
}

// toolSecretSet is secret_set: NAME=value into <world>/.env, from his glass alone.
func toolSecretSet(t tenant.Tenant, args map[string]any) (string, error) {
	c := callerOf(args)
	if !c.Service || c.Hold != "" {
		return "", errors.New("secret_set is the operator's own hand: his glass alone may call it, and never by a replayed hold")
	}
	name := strings.TrimSpace(str(args, "name"))
	value, isStr := args["value"].(string)
	if !SecretNameRe.MatchString(name) {
		return secretRefuse("the name must be an upper-case .env variable (A-Z, 0-9 and _, 2 to 64 characters, a letter first)")
	}
	if !isStr || value == "" {
		return secretRefuse("the value is empty")
	}
	if strings.ContainsAny(value, "\n\r\x00") {
		return secretRefuse("the value carries a line break or a NUL; a key is one line")
	}
	if !utf8.ValidString(value) {
		return secretRefuse("the value is not UTF-8 text")
	}
	n := utf8.RuneCountInString(value)
	if n > SecretMaxChars {
		return secretRefuse("the value is %d characters, over the %d a key may be", n, SecretMaxChars)
	}

	secretLock.Lock()
	defer secretLock.Unlock()
	if err := secretWrite(filepath.Join(t.Home, ".env"), name, value); err != nil {
		return "", err
	}
	said := fmt.Sprintf("secret %s set in .env (%d chars)", name, n)
	if reg := regOf(args); reg != nil {
		reg.record(t.Home, "secret_set", Hold{
			ID: fmt.Sprintf("secret_%d_%s", time.Now().UnixMilli(), name), Tool: "secret_set",
			Caller: callerLabel(c), Project: t.Name, RBAC: rbacState(t, c),
		}, said)
	}
	return said, nil
}

// secretWrite puts name=value into the .env at abs, replacing every live `name=` line in place or appending one, and
// keeping the file's own line endings and byte-order mark. A refusal never carries the file's text or the value.
func secretWrite(abs, name, value string) error {
	st, err := os.Lstat(abs)
	switch {
	case errors.Is(err, os.ErrNotExist):
		if werr := os.WriteFile(abs, []byte(name+"="+value+"\n"), 0o600); werr != nil {
			return fmt.Errorf("refused: .env could not be created (%s): nothing was written", secretErrKind(werr))
		}
		return nil
	case err != nil:
		return fmt.Errorf("refused: .env could not be read (%s): nothing was written", secretErrKind(err))
	case !st.Mode().IsRegular():
		return errors.New("refused: .env is not a plain file (a folder or a link): nothing was written")
	}
	b, err := os.ReadFile(abs)
	if err != nil {
		return fmt.Errorf("refused: .env could not be read (%s): nothing was written", secretErrKind(err))
	}
	if !utf8.Valid(b) || strings.IndexByte(string(b), 0) >= 0 {
		return errors.New("refused: .env is not UTF-8 text (PowerShell's > writes UTF-16): nothing was written")
	}
	eol, ok := aiderEol(b)
	if !ok {
		return errors.New("refused: .env mixes two line endings, so a write cannot keep what it has: nothing was written")
	}
	flat := aiderLF(b)
	const bomMark = "\xef\xbb\xbf" // UTF-8's byte-order mark, which PowerShell's -Encoding UTF8 writes
	bom := ""
	if strings.HasPrefix(flat, bomMark) {
		bom, flat = bomMark, strings.TrimPrefix(flat, bomMark)
	}
	lines := strings.Split(flat, "\n")
	set := name + "=" + value
	found := false
	for i, l := range lines {
		if secretLineNames(l, name) {
			lines[i] = set
			found = true
		}
	}
	out := strings.Join(lines, "\n")
	if !found {
		if out != "" && !strings.HasSuffix(out, "\n") {
			out += "\n"
		}
		out += set + "\n"
	}
	if err := os.WriteFile(abs, []byte(bom+aiderApplyEol(out, eol)), st.Mode().Perm()); err != nil {
		return fmt.Errorf("refused: .env could not be written (%s): nothing was written", secretErrKind(err))
	}
	return nil
}

// secretLineNames is whether one .env line sets name, as the engine's dotenv.py reads it: leading space and an
// `export ` ignored, the key trimmed, a comment never.
func secretLineNames(line, name string) bool {
	l := strings.TrimSpace(line)
	if l == "" || strings.HasPrefix(l, "#") {
		return false
	}
	if len(l) > 7 && strings.EqualFold(l[:7], "export ") {
		l = strings.TrimLeft(l[7:], " \t")
	}
	k, _, has := strings.Cut(l, "=")
	return has && strings.TrimSpace(k) == name
}

// secretErrKind is an OS error said without its path or text: the kind of fault only.
func secretErrKind(err error) string {
	switch {
	case errors.Is(err, os.ErrPermission):
		return "permission denied"
	case errors.Is(err, os.ErrNotExist):
		return "not found"
	}
	return "an OS error"
}

// secretSetTool is secret_set's declaration; tools.go registers it with one line, r.add(secretSetTool()).
func secretSetTool() Tool {
	return Tool{
		Name:        "secret_set",
		Writes:      true,
		ServiceOnly: true,
		Secrets:     []string{"value"},
		Args:        []string{"name", "value", "project?"},
		Fn:          toolSecretSet,
		Description: "the operator's own hand: write one key into this world's .env as NAME=value, from the masked field on the card a needs_answer naming `secret` paints; an existing NAME= line is replaced in place, else one is appended, and the file keeps its own line endings. The name is an upper-case variable (^[A-Z][A-Z0-9_]{1,63}$), the value one line of at most 4096 characters. The value is never answered, recorded or logged: the answer says only the name and its length. His glass alone may call it",
	}
}
