//go:build !windows

package engine

import (
	"errors"
	"syscall"
)

// processGone reports whether pid names no running process, and every doubt
// answers NO (see processgone_windows.go). Signal 0 sends nothing and asks only
// whether the process exists; ESRCH is the one answer that proves it does not.
func processGone(pid int) bool {
	if pid <= 0 {
		return false
	}
	return errors.Is(syscall.Kill(pid, 0), syscall.ESRCH)
}
