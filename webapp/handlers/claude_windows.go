package handlers

import (
	"os/exec"
	"syscall"
)

// hidden starts the Claude tab's harness with no console window
// (CREATE_NO_WINDOW). Windows only: the fields do not exist elsewhere.
func hidden(cmd *exec.Cmd) *exec.Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	return cmd
}
