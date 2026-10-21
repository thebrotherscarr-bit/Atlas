//go:build !windows

package main

import "os/exec"

// hidden is the desktop app's console-less start; off Windows there is no
// console window to hide, so the command is started as it is.
func hidden(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}
