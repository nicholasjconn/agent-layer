//go:build windows

package dispatch

import (
	"os"

	"golang.org/x/sys/windows"
)

// lockFile acquires an exclusive lock on the file handle.
func lockFile(file *os.File) error {
	var overlapped windows.Overlapped
	return windows.LockFileEx(windows.Handle(file.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, &overlapped)
}

// unlockFile releases the lock on the file handle.
func unlockFile(file *os.File) error {
	var overlapped windows.Overlapped
	return windows.UnlockFileEx(windows.Handle(file.Fd()), 0, 1, 0, &overlapped)
}
