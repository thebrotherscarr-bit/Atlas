package main

import (
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// The desktop app's pieces that do not need a screen: how it knows it is the
// desktop app, how it tells a port answers, and that it takes exactly one name
// out of .env -- the service wire -- and nothing else.

func TestDesktopModeIsTheFlagOrTheName(t *testing.T) {
	kept := os.Args
	defer func() { os.Args = kept }()
	os.Args = []string{"atlas-webapp.exe"}
	if desktopMode() {
		t.Fatal("a plain start of the glass is not the desktop app")
	}
	os.Args = []string{"atlas-webapp.exe", "--desktop"}
	if !desktopMode() {
		t.Fatal("--desktop starts the desktop app")
	}
}

func TestAnsweringAndWaitFor(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skip("no loopback listener here")
	}
	_, port, _ := net.SplitHostPort(l.Addr().String())
	if !answering(port) || !waitFor(port, time.Second) {
		t.Fatal("a listening port answers")
	}
	l.Close()
	if answering(port) {
		t.Fatal("a closed port does not answer")
	}
}

func TestServiceFromDotenvTakesOnlyTheWire(t *testing.T) {
	dir := t.TempDir()
	body := "OTHER_KEY=not-taken\n# ATLAS_SERVICE=commented\nATLAS_SERVICE=\"wire-123\"\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	kept, had := os.LookupEnv("ATLAS_SERVICE")
	os.Unsetenv("ATLAS_SERVICE")
	os.Unsetenv("OTHER_KEY")
	defer func() {
		if had {
			os.Setenv("ATLAS_SERVICE", kept)
		} else {
			os.Unsetenv("ATLAS_SERVICE")
		}
	}()
	serviceFromDotenv(dir)
	if os.Getenv("ATLAS_SERVICE") != "wire-123" {
		t.Fatalf("the wire is read, quotes stripped: got %q", os.Getenv("ATLAS_SERVICE"))
	}
	if os.Getenv("OTHER_KEY") != "" {
		t.Fatal("no other name in .env is taken")
	}
	os.Setenv("ATLAS_SERVICE", "already")
	serviceFromDotenv(dir)
	if os.Getenv("ATLAS_SERVICE") != "already" {
		t.Fatal("a wire already set is never replaced")
	}
}
