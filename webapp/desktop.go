package main

// THE DESKTOP APP (2026-10-21, his word: "compile this thing and turn it into a
// desktop app like ollama desktop or codex or hermes desktop or the new chatgpt
// app"). The glass already is the app -- one page, one terminal, the Inspector,
// the dock -- so this is not a second front end: it is the glass started as a
// program of its own. Run as `atlas-desktop.exe` (or with `--desktop`) it
//
//   - reads the service wire (ATLAS_SERVICE) from the ground's .env into its own
//     environment, the way the start lines in RUNBOOK do by hand -- never printed,
//     never on a command line (RULE 7);
//   - starts Ollama and the door when they are not answering, exactly as
//     after-a-reboot is done by hand: the door from atlas/line with the RUNBOOK
//     line, its output to mcp.log and mcp.err.log beside it, no console window;
//   - serves the glass itself, or uses the one already running on its port;
//   - opens it in a window of its own -- Microsoft Edge's app mode, which every
//     Windows 11 machine carries, so nothing is downloaded and nothing new is
//     installed (RULE 4): no address bar, no tabs, its own taskbar entry -- and
//     ends when that window is closed. The door and Ollama keep running, as
//     Ollama's own desktop app leaves its server running.
//
// The lock is unchanged: the window opens on the same PIN page the glass shows a
// browser.

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// desktopMode says whether this start is the desktop app's.
func desktopMode() bool {
	exe, _ := os.Executable()
	if strings.Contains(strings.ToLower(filepath.Base(exe)), "desktop") {
		return true
	}
	for _, a := range os.Args[1:] {
		if a == "--desktop" || a == "-desktop" {
			return true
		}
	}
	return false
}

// answering says whether something listens on a loopback port.
func answering(port string) bool {
	c, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 400*time.Millisecond)
	if err != nil {
		return false
	}
	c.Close()
	return true
}

// waitFor waits until a loopback port answers, or the time is up.
func waitFor(port string, d time.Duration) bool {
	end := time.Now().Add(d)
	for time.Now().Before(end) {
		if answering(port) {
			return true
		}
		time.Sleep(250 * time.Millisecond)
	}
	return false
}

// desktopLog is where the app says what it did, beside the glass's own logs:
// a GUI program has no console to say it on.
func desktopLog(format string, a ...any) {
	f, err := os.OpenFile(filepath.Join("data", "desktop.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, time.Now().Format("2006-01-02 15:04:05")+"  "+format+"\n", a...)
}

// serviceFromDotenv puts ATLAS_SERVICE in this process's environment from the
// ground's .env when it is not there already. Only that one name is read; the
// value is never printed or logged.
func serviceFromDotenv(ground string) {
	if os.Getenv("ATLAS_SERVICE") != "" {
		return
	}
	f, err := os.Open(filepath.Join(ground, ".env"))
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if !strings.HasPrefix(line, "ATLAS_SERVICE=") {
			continue
		}
		v := strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "ATLAS_SERVICE=")), `"'`)
		if v != "" {
			os.Setenv("ATLAS_SERVICE", v)
			desktopLog("the service wire is held (read from .env)")
		}
		return
	}
}

// startOllama starts Ollama's own app when the rack does not answer.
func startOllama() {
	if answering("11434") {
		return
	}
	app := filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Ollama", "ollama app.exe")
	if _, err := os.Stat(app); err != nil {
		desktopLog("Ollama is not answering and its app is not at %s", app)
		return
	}
	if err := hidden(exec.Command(app)).Start(); err != nil {
		desktopLog("Ollama could not be started: %v", err)
		return
	}
	desktopLog("Ollama started: %v", waitFor("11434", 20*time.Second))
}

// startDoor starts the door with RUNBOOK's line when it does not answer.
func startDoor(ground string) {
	if answering("8090") {
		return
	}
	line := filepath.Join(ground, "atlas", "line")
	exe := filepath.Join(line, "atlas-mcp.exe")
	if _, err := os.Stat(exe); err != nil {
		desktopLog("the door is not answering and is not at %s", exe)
		return
	}
	g := filepath.ToSlash(ground)
	cmd := hidden(exec.Command(exe, "--http", "127.0.0.1:8090", "--auth",
		"--tenant", "research="+g, "--tenant", "atlas="+g+"/atlas",
		"--default-project", "research", "--manjuel", "python "+g+"/manjuel.py"))
	cmd.Dir = line
	if out, err := os.OpenFile(filepath.Join(line, "mcp.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644); err == nil {
		cmd.Stdout = out
	}
	if errf, err := os.OpenFile(filepath.Join(line, "mcp.err.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644); err == nil {
		cmd.Stderr = errf
	}
	if err := cmd.Start(); err != nil {
		desktopLog("the door could not be started: %v", err)
		return
	}
	desktopLog("the door started (pid %d): %v", cmd.Process.Pid, waitFor("8090", 20*time.Second))
}

// edge finds Microsoft Edge.
func edge() string {
	for _, root := range []string{os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"),
		filepath.Join(os.Getenv("LOCALAPPDATA"))} {
		if root == "" {
			continue
		}
		p := filepath.Join(root, "Microsoft", "Edge", "Application", "msedge.exe")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// openWindow opens the glass in a window of its own and waits until it is
// closed. Its own profile folder, under the glass's data/, is what makes the
// window its own process (and keeps its cookies apart from his browser's).
func openWindow(port string) error {
	browser := edge()
	if browser == "" {
		return fmt.Errorf("Microsoft Edge was not found; open http://127.0.0.1:%s/ in a browser", port)
	}
	profile, _ := filepath.Abs(filepath.Join("data", "desktop-profile"))
	cmd := exec.Command(browser, "--app=http://127.0.0.1:"+port+"/", "--user-data-dir="+profile,
		"--window-size=1440,920", "--no-first-run", "--no-default-browser-check",
		"--disable-features=Translate")
	if err := cmd.Start(); err != nil {
		return err
	}
	desktopLog("the window opened (pid %d)", cmd.Process.Pid)
	return cmd.Wait()
}

// desktopPrepare is everything before the glass serves: its own folder as the
// working directory (the glass reads data/ relative to it), the wire, Ollama and
// the door. It returns true when a glass already answers on the port, so this
// start only opens a window onto it.
func desktopPrepare(port string) bool {
	exe, err := os.Executable()
	if err == nil {
		os.Chdir(filepath.Dir(exe))
	}
	ground, _ := filepath.Abs(filepath.Join("..", ".."))
	desktopLog("starting: the ground is %s", ground)
	serviceFromDotenv(ground)
	startOllama()
	startDoor(ground)
	return answering(port)
}
