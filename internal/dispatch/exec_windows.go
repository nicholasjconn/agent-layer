//go:build windows

package dispatch

import (
	"os"
	"os/exec"
)

// execBinary runs the target binary on Windows and exits with its status code.
func execBinary(path string, args []string, env []string, exit func(int)) error {
	cmd := exec.Command(path, args[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exit(exitErr.ExitCode())
			return ErrDispatched
		}
		return err
	}
	exit(0)
	return ErrDispatched
}
