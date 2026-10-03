// Rack prove: the golden ladder byte-exact from folded tags, a stub door
// on loopback, refusal of outward hosts, and honest silence when down.
// Hermetic: fixture files, httptest loopback, closed-port silence.
package rack

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "tests", "fixtures", name))
	if err != nil {
		t.Fatalf("fixture missing: %s", err)
	}
	return string(b)
}

func TestLadderMatchesGolden(t *testing.T) {
	tags := fixture(t, "rack_tags.json")
	want := fixture(t, "rack_ladder.txt")
	voices, err := ParseTags([]byte(tags))
	if err != nil {
		t.Fatal(err)
	}
	if got := Ladder(voices); got != want {
		t.Fatalf("ladder mismatch:\n go %q\n py %q", got, want)
	}
}

func TestStubDoorListsTiers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(fixture(t, "rack_tags.json")))
	}))
	defer srv.Close()
	t.Setenv("OLLAMA_HOST", srv.URL)
	host, err := Host()
	if err != nil {
		t.Fatal(err)
	}
	voices, err := List(host)
	if err != nil {
		t.Fatal(err)
	}
	if len(voices) != 9 {
		t.Fatalf("want 9 stub voices, got %d", len(voices))
	}
	ladder := Ladder(voices)
	for _, want := range []string{"scout (", "voice (", "mind (", "qwen2.5-coder:14b", "llama3.2:latest"} {
		if !strings.Contains(ladder, want) {
			t.Fatalf("ladder missing %q", want)
		}
	}
}

func TestOutwardHostRefused(t *testing.T) {
	for _, h := range []string{"http://example.com:11434", "http://192.168.1.10:11434", "https://10.0.0.5"} {
		t.Setenv("OLLAMA_HOST", h)
		if _, err := Host(); err == nil {
			t.Fatalf("outward host %q admitted", h)
		}
	}
	// This-host spellings dial loopback, never outward: admitted.
	for _, h := range []string{"0.0.0.0:11434", "http://0.0.0.0:11434", "http://127.0.0.1:11434"} {
		t.Setenv("OLLAMA_HOST", h)
		if _, err := Host(); err != nil {
			t.Fatalf("this-host %q refused: %s", h, err)
		}
	}
}

// A host with no port is Ollama's port (WHAT'S LEFT C36). This machine sets
// OLLAMA_HOST=127.0.0.1; the rack read it as http://127.0.0.1, dialled port 80
// and called a standing Ollama silent. A port that is named is never changed.
func TestPortlessHostMeansOllamasPort(t *testing.T) {
	for in, want := range map[string]string{
		"":                    "http://127.0.0.1:11434",
		"127.0.0.1":           "http://127.0.0.1:11434",
		"localhost":           "http://localhost:11434",
		"0.0.0.0":             "http://0.0.0.0:11434",
		"http://127.0.0.1":    "http://127.0.0.1:11434",
		"http://127.0.0.1/":   "http://127.0.0.1:11434",
		"http://127.0.0.1:":   "http://127.0.0.1:11434",
		"[::1]":               "http://[::1]:11434",
		"127.0.0.1:11500":     "http://127.0.0.1:11500",
		"http://127.0.0.1:80": "http://127.0.0.1:80",
	} {
		t.Setenv("OLLAMA_HOST", in)
		got, err := Host()
		if err != nil {
			t.Fatalf("OLLAMA_HOST=%q refused: %s", in, err)
		}
		if got != want {
			t.Fatalf("OLLAMA_HOST=%q read as %q, want %q", in, got, want)
		}
	}
	// Adding the port must not make an outward host lawful.
	for _, h := range []string{"example.com", "http://192.168.1.10", "https://10.0.0.5"} {
		t.Setenv("OLLAMA_HOST", h)
		if _, err := Host(); err == nil {
			t.Fatalf("outward host %q admitted once it gained a port", h)
		}
	}
}

func TestSilenceHonest(t *testing.T) {
	// Port 1 on loopback refuses: deterministic silence, honest reason.
	t.Setenv("OLLAMA_HOST", "http://127.0.0.1:1")
	host, err := Host()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := List(host); err == nil {
		t.Fatal("silence reported as voices")
	} else if !strings.Contains(err.Error(), "silent") {
		t.Fatalf("silence without the honest word: %s", err)
	}
}
