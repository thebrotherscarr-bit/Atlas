//go:build windows

package engine

import "syscall"

// processGone reports whether pid names no running process, and every doubt
// answers NO -- the core's own rule (seatlog._alive): believing a live sitting
// dead would open a second engine under somebody's open sitting, the one thing
// the refusal in Open exists to stop. Two answers prove a death: the kernel has
// no such process, or it has one that has already exited.
//
// ASKED OF THE KERNEL, NEVER SIGNALLED. The same reason seatlog gives for
// ctypes over os.kill: on Windows a "liveness signal" is not a thing, and the
// portable-looking calls either do nothing or end the process they ask about.
func processGone(pid int) bool {
	if pid <= 0 {
		return false
	}
	const queryLimited = 0x1000 // PROCESS_QUERY_LIMITED_INFORMATION
	const stillActive = 259     // STILL_ACTIVE
	h, err := syscall.OpenProcess(queryLimited, false, uint32(pid))
	if err != nil {
		// ERROR_INVALID_PARAMETER is the kernel saying there is no such
		// process. Anything else -- access denied above all -- is a process
		// this door may not ask about, which is not a death.
		return err == syscall.Errno(87)
	}
	defer syscall.CloseHandle(h)
	var code uint32
	if err := syscall.GetExitCodeProcess(h, &code); err != nil {
		return false
	}
	return code != stillActive
}
