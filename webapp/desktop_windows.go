package main

import (
	"os/exec"
	"syscall"
)

// hidden starts a process with no console window (CREATE_NO_WINDOW), the desktop
// app's way of starting the door and Ollama behind its own window. Windows only:
// the fields do not exist elsewhere, which is why this lives in its own file.
func hidden(cmd *exec.Cmd) *exec.Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	return cmd
}
