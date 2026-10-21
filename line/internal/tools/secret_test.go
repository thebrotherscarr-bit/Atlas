package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// THE SECRET CARD, held by strokes (secret.go): a key pasted into the masked field reaches the world's .env and
// nowhere else -- not the answer, not a refusal, not the holds record -- and only from his glass.

const secretProbe = "s3cr3t-VALUE-that-must-never-be-said-0123456789"

func secretCall(t *testing.T, name, value string, files map[string]string) (string, error, string) {
	t.Helper()
	reg, tr, tn := shellWorld(t, files)
	out, err := reg.Call(tr, "secret_set", map[string]any{"name": name, "value": value}, glass)
	return out, err, tn.Home
}

func readEnv(t *testing.T, home string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(home, ".env"))
	if err != nil {
		t.Fatalf("no .env: %v", err)
	}
	return string(b)
}

// nothingSaysTheValue fails if the value is in the answer, the error, or anything the tool wrote under state/.
func nothingSaysTheValue(t *testing.T, home, value, out string, err error) {
	t.Helper()
	if strings.Contains(out, value) {
		t.Fatalf("the answer carries the value: %q", out)
	}
	if err != nil && strings.Contains(err.Error(), value) {
		t.Fatalf("the refusal carries the value: %v", err)
	}
	_ = filepath.WalkDir(filepath.Join(home, "state"), func(p string, d os.DirEntry, werr error) error {
		if werr != nil || d.IsDir() {
			return nil
		}
		if b, rerr := os.ReadFile(p); rerr == nil && strings.Contains(string(b), value) {
			t.Fatalf("%s carries the value", p)
		}
		return nil
	})
}

// A KEY IS WRITTEN WHERE NONE WAS: .env is created, the answer says the name and the length and never the value, and
// the record line says the same.
func TestASecretIsWrittenIntoAFreshEnv(t *testing.T) {
	out, err, home := secretCall(t, "DISCORD_BOT_TOKEN", secretProbe, nil)
	if err != nil {
		t.Fatalf("refused: %v", err)
	}
	want := "secret DISCORD_BOT_TOKEN set in .env (47 chars)"
	if out != want {
		t.Fatalf("answered %q, want %q", out, want)
	}
	if got := readEnv(t, home); got != "DISCORD_BOT_TOKEN="+secretProbe+"\n" {
		t.Fatalf(".env is %q", got)
	}
	nothingSaysTheValue(t, home, secretProbe, out, err)
	rec, rerr := os.ReadFile(filepath.Join(home, "state", "holds.jsonl"))
	if rerr != nil || !strings.Contains(string(rec), `"what":"secret_set"`) || !strings.Contains(string(rec), want) {
		t.Fatalf("the record does not say the write: %q (%v)", rec, rerr)
	}
}

// AN EXISTING LINE IS REPLACED IN PLACE, every live one, and the file keeps its CRLF, its order, its comments and the
// other keys exactly as they were; a commented-out line and a longer name are not touched.
func TestASecretReplacesItsLineInPlaceKeepingCRLF(t *testing.T) {
	before := "# keys\r\nOLLAMA_HOST=127.0.0.1\r\nexport DISCORD_BOT_TOKEN=old\r\n# DISCORD_BOT_TOKEN=commented\r\nDISCORD_BOT_TOKEN_2=keep\r\n  DISCORD_BOT_TOKEN = older\r\nLAST=1\r\n"
	out, err, home := secretCall(t, "DISCORD_BOT_TOKEN", secretProbe, map[string]string{".env": before})
	if err != nil {
		t.Fatalf("refused: %v", err)
	}
	want := "# keys\r\nOLLAMA_HOST=127.0.0.1\r\nDISCORD_BOT_TOKEN=" + secretProbe + "\r\n# DISCORD_BOT_TOKEN=commented\r\nDISCORD_BOT_TOKEN_2=keep\r\nDISCORD_BOT_TOKEN=" + secretProbe + "\r\nLAST=1\r\n"
	if got := readEnv(t, home); got != want {
		t.Fatalf(".env is\n%q\nwant\n%q", got, want)
	}
	nothingSaysTheValue(t, home, secretProbe, out, err)
}

// A NEW NAME IS APPENDED, in the file's own line ending, after a last line that had none; a byte-order mark stays.
func TestASecretIsAppendedInTheFilesOwnEnding(t *testing.T) {
	_, err, home := secretCall(t, "API_KEY", "abc", map[string]string{".env": "\xef\xbb\xbfA=1\r\nB=2"})
	if err != nil {
		t.Fatalf("refused: %v", err)
	}
	if got := readEnv(t, home); got != "\xef\xbb\xbfA=1\r\nB=2\r\nAPI_KEY=abc\r\n" {
		t.Fatalf(".env is %q", got)
	}
	_, err, home = secretCall(t, "API_KEY", "abc", map[string]string{".env": "A=1\n"})
	if err != nil {
		t.Fatalf("refused: %v", err)
	}
	if got := readEnv(t, home); got != "A=1\nAPI_KEY=abc\n" {
		t.Fatalf(".env is %q", got)
	}
}

// WHAT IT MUST REFUSE, IT REFUSES: a bad name, a value with a line break or a NUL, an empty or an over-long value, a
// .env that mixes line endings or is UTF-16. Each writes nothing and none says the value.
func TestASecretIsRefusedWhereItMustBe(t *testing.T) {
	cases := []struct{ why, name, value, env string }{
		{"lower-case name", "discord_token", secretProbe, ""},
		{"one-letter name", "A", secretProbe, ""},
		{"name with a dash", "MY-KEY", secretProbe, ""},
		{"name with an =", "KEY=X", secretProbe, ""},
		{"name starting with a digit", "1KEY", secretProbe, ""},
		{"name of 65", "A" + strings.Repeat("B", 64), secretProbe, ""},
		{"newline", "KEY", secretProbe + "\nEVIL=1", ""},
		{"carriage return", "KEY", secretProbe + "\r", ""},
		{"NUL", "KEY", secretProbe + "\x00", ""},
		{"empty", "KEY", "", ""},
		{"over 4096", "KEY", secretProbe + strings.Repeat("x", SecretMaxChars), ""},
		{"mixed endings", "KEY", secretProbe, "A=1\r\nB=2\n"},
		{"utf-16", "KEY", secretProbe, "A\x00=\x001\x00"},
	}
	for _, c := range cases {
		t.Run(c.why, func(t *testing.T) {
			files := map[string]string{}
			if c.env != "" {
				files[".env"] = c.env
			}
			out, err, home := secretCall(t, c.name, c.value, files)
			if err == nil || out != "" || !strings.Contains(err.Error(), "nothing was written") {
				t.Fatalf("not refused: %q / %v", out, err)
			}
			if c.value != "" {
				nothingSaysTheValue(t, home, secretProbe, out, err)
			}
			b, rerr := os.ReadFile(filepath.Join(home, ".env"))
			if c.env == "" && rerr == nil {
				t.Fatalf("a refusal created .env: %q", b)
			}
			if c.env != "" && string(b) != c.env {
				t.Fatalf("a refusal changed .env: %q", b)
			}
		})
	}
}

// ONLY HIS GLASS: a seat, an unnamed client and a replayed hold are refused by name, and nothing is written.
func TestASecretIsHisGlassAlone(t *testing.T) {
	reg, tr, tn := shellWorld(t, nil)
	for _, who := range []Caller{agent, {}, {Name: "k-agent"}} {
		out, err := reg.Call(tr, "secret_set", map[string]any{"name": "KEY", "value": secretProbe}, who)
		if err == nil || out != "" || !strings.Contains(err.Error(), "operator's own hand") {
			t.Fatalf("%q was answered: %q / %v", who.Name, out, err)
		}
		nothingSaysTheValue(t, tn.Home, secretProbe, out, err)
	}
	tool, ok := reg.Get("secret_set")
	if !ok || !tool.ServiceOnly || !tool.Writes || !hasArg(tool.Secrets, "value") {
		t.Fatalf("secret_set is not declared the operator's own writing hand with a secret value: %+v", tool)
	}
	for _, who := range []Caller{agent, {}, {Name: "operator via hold hold_1", Service: true, Hold: "hold_1"}} {
		out, err := tool.Fn(tn, map[string]any{"name": "KEY", "value": secretProbe, CallerKey: who})
		if err == nil || out != "" || !strings.Contains(err.Error(), "operator's own hand") {
			t.Fatalf("the tool ran for %q: %q / %v", who.Name, out, err)
		}
	}
	if exists(filepath.Join(tn.Home, ".env")) {
		t.Fatal("a caller that was not his glass wrote .env")
	}
}
