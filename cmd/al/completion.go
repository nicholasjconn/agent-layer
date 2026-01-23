//go:build !windows

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/conn-castle/agent-layer/internal/fsutil"
)

var (
	userHomeDir = os.UserHomeDir
	lookPath    = exec.LookPath
	execCommand = exec.Command
)

// newCompletionCmd builds the completion subcommand with optional install behavior.
func newCompletionCmd() *cobra.Command {
	var install bool
	cmd := &cobra.Command{
		Use:       "completion [bash|zsh|fish]",
		Short:     "Generate shell completion scripts",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish"},
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]
			script, err := generateCompletion(cmd.Root(), shell)
			if err != nil {
				return err
			}
			if !install {
				_, err := fmt.Fprint(cmd.OutOrStdout(), script)
				return err
			}
			return installCompletion(shell, script, cmd.OutOrStdout())
		},
	}
	cmd.Flags().BoolVar(&install, "install", false, "Install the completion script for the specified shell")
	return cmd
}

// generateCompletion renders the completion script for the selected shell.
func generateCompletion(root *cobra.Command, shell string) (string, error) {
	var buf bytes.Buffer
	switch shell {
	case "bash":
		if err := root.GenBashCompletion(&buf); err != nil {
			return "", err
		}
	case "zsh":
		if err := root.GenZshCompletion(&buf); err != nil {
			return "", err
		}
	case "fish":
		if err := root.GenFishCompletion(&buf, true); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported shell %q (supported: bash, zsh, fish)", shell)
	}
	return buf.String(), nil
}

// installCompletion writes the completion script to the user-level install location.
func installCompletion(shell string, script string, out io.Writer) error {
	path, note, err := completionInstallPath(shell)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create completion dir: %w", err)
	}
	if err := fsutil.WriteFileAtomic(path, []byte(script), 0o644); err != nil {
		return fmt.Errorf("write completion file: %w", err)
	}

	if _, err := fmt.Fprintf(out, "Installed %s completion to %s\n", shell, path); err != nil {
		return err
	}
	if note != "" {
		if _, err := fmt.Fprintln(out, note); err != nil {
			return err
		}
	}
	return nil
}

// completionInstallPath returns the destination path and any follow-up note to display.
func completionInstallPath(shell string) (string, string, error) {
	switch shell {
	case "bash":
		xdgData, err := xdgDataHome()
		if err != nil {
			return "", "", err
		}
		path := filepath.Join(xdgData, "bash-completion", "completions", "al")
		note := "Bash completion requires bash-completion to be enabled in your shell."
		return path, note, nil
	case "fish":
		xdgConfig, err := xdgConfigHome()
		if err != nil {
			return "", "", err
		}
		path := filepath.Join(xdgConfig, "fish", "completions", "al.fish")
		note := "Restart fish or open a new terminal to enable completions."
		return path, note, nil
	case "zsh":
		dir, ok := firstWritableFpath()
		if ok {
			return filepath.Join(dir, "_al"), "", nil
		}
		xdgData, err := xdgDataHome()
		if err != nil {
			return "", "", err
		}
		fallbackDir := filepath.Join(xdgData, "zsh", "site-functions")
		note := fmt.Sprintf("Add this to your .zshrc before compinit:\n  fpath=(%s $fpath)", fallbackDir)
		return filepath.Join(fallbackDir, "_al"), note, nil
	default:
		return "", "", fmt.Errorf("unsupported shell %q (supported: bash, zsh, fish)", shell)
	}
}

// xdgDataHome resolves the XDG data home directory.
func xdgDataHome() (string, error) {
	if value := strings.TrimSpace(os.Getenv("XDG_DATA_HOME")); value != "" {
		return value, nil
	}
	home, err := userHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".local", "share"), nil
}

// xdgConfigHome resolves the XDG config home directory.
func xdgConfigHome() (string, error) {
	if value := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME")); value != "" {
		return value, nil
	}
	home, err := userHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".config"), nil
}

// firstWritableFpath returns the first writable directory in $fpath, if any.
func firstWritableFpath() (string, bool) {
	for _, dir := range parseFpathEnv() {
		if dir == "" {
			continue
		}
		if writableDir(dir) {
			return dir, true
		}
	}

	zshPath, err := lookPath("zsh")
	if err != nil {
		return "", false
	}
	cmd := execCommand(zshPath, "-c", "print -l $fpath")
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	for _, dir := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}
		if writableDir(dir) {
			return dir, true
		}
	}
	return "", false
}

// parseFpathEnv splits $FPATH into a list of directories.
func parseFpathEnv() []string {
	raw := strings.TrimSpace(os.Getenv("FPATH"))
	if raw == "" {
		return nil
	}
	return strings.Split(raw, string(os.PathListSeparator))
}

// writableDir reports whether dir exists and is writable.
func writableDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}
	tmp, err := os.CreateTemp(dir, "al-fpath-*")
	if err != nil {
		return false
	}
	name := tmp.Name()
	_ = tmp.Close()
	_ = os.Remove(name)
	return true
}
