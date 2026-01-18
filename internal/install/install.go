package install

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicholasjconn/agent-layer/internal/templates"
)

const (
	gitignoreStart = "# >>> agent-layer"
	gitignoreEnd   = "# <<< agent-layer"
)

// Run initializes the repository with the required Agent Layer structure.
func Run(root string) error {
	if root == "" {
		return fmt.Errorf("root path is required")
	}

	steps := []func() error{
		func() error { return createDirs(root) },
		func() error { return writeTemplateFiles(root) },
		func() error { return writeTemplateDirs(root) },
		func() error { return updateGitignore(root) },
	}

	return runSteps(steps)
}

func runSteps(steps []func() error) error {
	for _, step := range steps {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
}

func createDirs(root string) error {
	dirs := []string{
		filepath.Join(root, ".agent-layer", "instructions"),
		filepath.Join(root, ".agent-layer", "slash-commands"),
		filepath.Join(root, ".agent-layer", "templates", "docs"),
		filepath.Join(root, "docs", "agent-layer"),
		filepath.Join(root, "tmp", "agent-layer", "runs"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func writeTemplateFiles(root string) error {
	files := []struct {
		path     string
		template string
		perm     fs.FileMode
	}{
		{filepath.Join(root, ".agent-layer", "config.toml"), "config.toml", 0o644},
		{filepath.Join(root, ".agent-layer", "commands.allow"), "commands.allow", 0o644},
		{filepath.Join(root, ".agent-layer", ".env"), "env", 0o600},
		{filepath.Join(root, ".agent-layer", "gitignore.block"), "gitignore.block", 0o644},
	}
	for _, file := range files {
		if err := writeTemplateIfMissing(file.path, file.template, file.perm); err != nil {
			return err
		}
	}
	return nil
}

func writeTemplateDirs(root string) error {
	dirs := []struct {
		templateRoot string
		destRoot     string
	}{
		{"instructions", filepath.Join(root, ".agent-layer", "instructions")},
		{"slash-commands", filepath.Join(root, ".agent-layer", "slash-commands")},
		{"docs/agent-layer", filepath.Join(root, "docs", "agent-layer")},
		{"docs/agent-layer", filepath.Join(root, ".agent-layer", "templates", "docs")},
	}
	for _, dir := range dirs {
		if err := writeTemplateDir(dir.templateRoot, dir.destRoot); err != nil {
			return err
		}
	}
	return nil
}

func updateGitignore(root string) error {
	blockPath := filepath.Join(root, ".agent-layer", "gitignore.block")
	blockBytes, err := os.ReadFile(blockPath)
	if err != nil {
		return fmt.Errorf("failed to read gitignore block %s: %w", blockPath, err)
	}
	return ensureGitignore(filepath.Join(root, ".gitignore"), string(blockBytes))
}

func writeTemplateIfMissing(path string, templatePath string, perm fs.FileMode) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}

	data, err := templates.Read(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

func writeTemplateDir(templateRoot string, destRoot string) error {
	return templates.Walk(templateRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel := strings.TrimPrefix(path, templateRoot+"/")
		if rel == path {
			return fmt.Errorf("unexpected template path %s", path)
		}
		destPath := filepath.Join(destRoot, rel)
		return writeTemplateIfMissing(destPath, path, 0o644)
	})
}

func ensureGitignore(path string, block string) error {
	block = strings.TrimRight(block, "\n") + "\n"
	contentBytes, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(path, []byte(block), 0o644)
	}

	content := string(contentBytes)
	start := strings.Index(content, gitignoreStart)
	end := strings.Index(content, gitignoreEnd)
	if start != -1 && end != -1 && end > start {
		end += len(gitignoreEnd)
		updated := content[:start] + block + content[end:]
		return os.WriteFile(path, []byte(updated), 0o644)
	}

	separator := "\n"
	if strings.HasSuffix(content, "\n") {
		separator = ""
	}
	updated := content + separator + block
	return os.WriteFile(path, []byte(updated), 0o644)
}
