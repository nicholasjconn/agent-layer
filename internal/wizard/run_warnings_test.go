package wizard

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nicholasjconn/agent-layer/internal/warnings"
)

func TestRun_WarningsConfirmError(t *testing.T) {
	root := t.TempDir()
	setupRepo(t, root)
	configDir := filepath.Join(root, ".agent-layer")
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(basicAgentConfig()), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, ".env"), []byte(""), 0600))

	confirmCalls := 0
	ui := &MockUI{
		NoteFunc:        func(title, body string) error { return nil },
		SelectFunc:      func(title string, options []string, current *string) error { return nil },
		MultiSelectFunc: func(title string, options []string, selected *[]string) error { return nil },
		ConfirmFunc: func(title string, value *bool) error {
			confirmCalls++
			// Restore missing is first, warnings is second
			if strings.Contains(title, "Enable warnings") {
				return errors.New("warnings confirm error")
			}
			*value = true
			return nil
		},
	}

	err := Run(root, ui, func(r string) ([]warnings.Warning, error) { return nil, nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "warnings confirm error")
}

func TestRun_WarningsEnabled_HappyPath(t *testing.T) {
	root := t.TempDir()
	setupRepo(t, root)
	configDir := filepath.Join(root, ".agent-layer")
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(basicAgentConfig()), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, ".env"), []byte(""), 0600))

	ui := &MockUI{
		NoteFunc:        func(title, body string) error { return nil },
		SelectFunc:      func(title string, options []string, current *string) error { return nil },
		MultiSelectFunc: func(title string, options []string, selected *[]string) error { return nil },
		InputFunc: func(title string, value *string) error {
			// Accept the defaults for all warning thresholds
			return nil
		},
		ConfirmFunc: func(title string, value *bool) error {
			if strings.Contains(title, "Enable warnings") {
				*value = true // Enable warnings
			} else {
				*value = true
			}
			return nil
		},
	}

	err := Run(root, ui, func(r string) ([]warnings.Warning, error) { return nil, nil })
	require.NoError(t, err)

	// Verify warnings thresholds are in config
	data, _ := os.ReadFile(filepath.Join(configDir, "config.toml"))
	assert.Contains(t, string(data), "[warnings]")
}

func TestRun_WarningsEnabled_InputError(t *testing.T) {
	root := t.TempDir()
	setupRepo(t, root)
	configDir := filepath.Join(root, ".agent-layer")
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(basicAgentConfig()), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, ".env"), []byte(""), 0600))

	ui := &MockUI{
		NoteFunc:        func(title, body string) error { return nil },
		SelectFunc:      func(title string, options []string, current *string) error { return nil },
		MultiSelectFunc: func(title string, options []string, selected *[]string) error { return nil },
		InputFunc: func(title string, value *string) error {
			// Error on first input
			return errors.New("input error")
		},
		ConfirmFunc: func(title string, value *bool) error {
			if strings.Contains(title, "Enable warnings") {
				*value = true // Enable warnings
			} else {
				*value = true
			}
			return nil
		},
	}

	err := Run(root, ui, func(r string) ([]warnings.Warning, error) { return nil, nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "input error")
}

func TestRun_WarningsEnabled_SecondInputError(t *testing.T) {
	root := t.TempDir()
	setupRepo(t, root)
	configDir := filepath.Join(root, ".agent-layer")
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(basicAgentConfig()), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, ".env"), []byte(""), 0600))

	inputCalls := 0
	ui := &MockUI{
		NoteFunc:        func(title, body string) error { return nil },
		SelectFunc:      func(title string, options []string, current *string) error { return nil },
		MultiSelectFunc: func(title string, options []string, selected *[]string) error { return nil },
		InputFunc: func(title string, value *string) error {
			inputCalls++
			if inputCalls == 2 {
				return errors.New("second input error")
			}
			return nil
		},
		ConfirmFunc: func(title string, value *bool) error {
			if strings.Contains(title, "Enable warnings") {
				*value = true
			} else {
				*value = true
			}
			return nil
		},
	}

	err := Run(root, ui, func(r string) ([]warnings.Warning, error) { return nil, nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "second input error")
}

func TestRun_WarningsEnabled_ThirdInputError(t *testing.T) {
	root := t.TempDir()
	setupRepo(t, root)
	configDir := filepath.Join(root, ".agent-layer")
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(basicAgentConfig()), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, ".env"), []byte(""), 0600))

	inputCalls := 0
	ui := &MockUI{
		NoteFunc:        func(title, body string) error { return nil },
		SelectFunc:      func(title string, options []string, current *string) error { return nil },
		MultiSelectFunc: func(title string, options []string, selected *[]string) error { return nil },
		InputFunc: func(title string, value *string) error {
			inputCalls++
			if inputCalls == 3 {
				return errors.New("third input error")
			}
			return nil
		},
		ConfirmFunc: func(title string, value *bool) error {
			if strings.Contains(title, "Enable warnings") {
				*value = true
			} else {
				*value = true
			}
			return nil
		},
	}

	err := Run(root, ui, func(r string) ([]warnings.Warning, error) { return nil, nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "third input error")
}

func TestRun_WarningsEnabled_FourthInputError(t *testing.T) {
	root := t.TempDir()
	setupRepo(t, root)
	configDir := filepath.Join(root, ".agent-layer")
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(basicAgentConfig()), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, ".env"), []byte(""), 0600))

	inputCalls := 0
	ui := &MockUI{
		NoteFunc:        func(title, body string) error { return nil },
		SelectFunc:      func(title string, options []string, current *string) error { return nil },
		MultiSelectFunc: func(title string, options []string, selected *[]string) error { return nil },
		InputFunc: func(title string, value *string) error {
			inputCalls++
			if inputCalls == 4 {
				return errors.New("fourth input error")
			}
			return nil
		},
		ConfirmFunc: func(title string, value *bool) error {
			if strings.Contains(title, "Enable warnings") {
				*value = true
			} else {
				*value = true
			}
			return nil
		},
	}

	err := Run(root, ui, func(r string) ([]warnings.Warning, error) { return nil, nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fourth input error")
}

func TestRun_WarningsEnabled_FifthInputError(t *testing.T) {
	root := t.TempDir()
	setupRepo(t, root)
	configDir := filepath.Join(root, ".agent-layer")
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(basicAgentConfig()), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, ".env"), []byte(""), 0600))

	inputCalls := 0
	ui := &MockUI{
		NoteFunc:        func(title, body string) error { return nil },
		SelectFunc:      func(title string, options []string, current *string) error { return nil },
		MultiSelectFunc: func(title string, options []string, selected *[]string) error { return nil },
		InputFunc: func(title string, value *string) error {
			inputCalls++
			if inputCalls == 5 {
				return errors.New("fifth input error")
			}
			return nil
		},
		ConfirmFunc: func(title string, value *bool) error {
			if strings.Contains(title, "Enable warnings") {
				*value = true
			} else {
				*value = true
			}
			return nil
		},
	}

	err := Run(root, ui, func(r string) ([]warnings.Warning, error) { return nil, nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fifth input error")
}

func TestRun_WarningsEnabled_SixthInputError(t *testing.T) {
	root := t.TempDir()
	setupRepo(t, root)
	configDir := filepath.Join(root, ".agent-layer")
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(basicAgentConfig()), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, ".env"), []byte(""), 0600))

	inputCalls := 0
	ui := &MockUI{
		NoteFunc:        func(title, body string) error { return nil },
		SelectFunc:      func(title string, options []string, current *string) error { return nil },
		MultiSelectFunc: func(title string, options []string, selected *[]string) error { return nil },
		InputFunc: func(title string, value *string) error {
			inputCalls++
			if inputCalls == 6 {
				return errors.New("sixth input error")
			}
			return nil
		},
		ConfirmFunc: func(title string, value *bool) error {
			if strings.Contains(title, "Enable warnings") {
				*value = true
			} else {
				*value = true
			}
			return nil
		},
	}

	err := Run(root, ui, func(r string) ([]warnings.Warning, error) { return nil, nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sixth input error")
}
