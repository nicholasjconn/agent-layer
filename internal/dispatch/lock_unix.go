//go:build !windows

package dispatch

import (
	"os"

	"golang.org/x/sys/unix"
)

// lockFile acquires an exclusive advisory lock on the file.
func lockFile(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_EX)
}

// unlockFile releases the advisory lock on the file.
func unlockFile(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_UN)
}
