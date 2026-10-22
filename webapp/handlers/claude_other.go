//go:build !windows

package handlers

import "os/exec"

// hidden: off Windows there is no console window to hide.
func hidden(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}
