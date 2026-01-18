package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCreatesStructure(t *testing.T) {
	root := t.TempDir()
	if err := Run(root); err != nil {
		t.Fatalf("Run error: %v", err)
	}

	expectFiles := []string{
		filepath.Join(root, ".agent-layer", "config.toml"),
		filepath.Join(root, ".agent-layer", "commands.allow"),
		filepath.Join(root, ".agent-layer", ".env"),
		filepath.Join(root, ".agent-layer", "gitignore.block"),
		filepath.Join(root, "docs", "agent-layer", "ISSUES.md"),
	}
	for _, path := range expectFiles {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	gitignorePath := filepath.Join(root, ".gitignore")
	data, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("read gitignore: %v", err)
	}
	if !strings.Contains(string(data), gitignoreStart) {
		t.Fatalf("expected gitignore block to be present")
	}
}

func TestEnsureGitignoreCreatesFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	block := "# >>> agent-layer\nal\n# <<< agent-layer\n"

	if err := ensureGitignore(path, block); err != nil {
		t.Fatalf("ensureGitignore error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read gitignore: %v", err)
	}
	if string(data) != block {
		t.Fatalf("unexpected gitignore content: %q", string(data))
	}
}

func TestEnsureGitignoreReplacesBlock(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	original := "keep\n# >>> agent-layer\nold\n# <<< agent-layer\nend\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatalf("write gitignore: %v", err)
	}

	block := "# >>> agent-layer\nnew\n# <<< agent-layer\n"
	if err := ensureGitignore(path, block); err != nil {
		t.Fatalf("ensureGitignore error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read gitignore: %v", err)
	}
	if !strings.Contains(string(data), "new") || strings.Contains(string(data), "old") {
		t.Fatalf("expected block to be replaced, got %q", string(data))
	}
}

func TestEnsureGitignoreAppendsBlock(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	original := "keep\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatalf("write gitignore: %v", err)
	}

	block := "# >>> agent-layer\nnew\n# <<< agent-layer\n"
	if err := ensureGitignore(path, block); err != nil {
		t.Fatalf("ensureGitignore error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read gitignore: %v", err)
	}
	if !strings.Contains(string(data), "new") {
		t.Fatalf("expected appended block, got %q", string(data))
	}
}

func TestEnsureGitignorePartialBlock(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	original := "keep\n# >>> agent-layer\nold\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatalf("write gitignore: %v", err)
	}

	block := "# >>> agent-layer\nnew\n# <<< agent-layer\n"
	if err := ensureGitignore(path, block); err != nil {
		t.Fatalf("ensureGitignore error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read gitignore: %v", err)
	}
	if !strings.Contains(string(data), "new") {
		t.Fatalf("expected block to be appended")
	}
}

func TestWriteTemplateIfMissingExisting(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "config.toml")
	if err := os.WriteFile(path, []byte("custom"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if err := writeTemplateIfMissing(path, "config.toml", 0o644); err != nil {
		t.Fatalf("writeTemplateIfMissing error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "custom" {
		t.Fatalf("expected existing file to remain")
	}
}

func TestWriteTemplateDirMissing(t *testing.T) {
	root := t.TempDir()
	err := writeTemplateDir("missing-root", root)
	if err == nil {
		t.Fatalf("expected error for missing template root")
	}
}

func TestWriteTemplateIfMissingInvalidTemplate(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "config.toml")
	err := writeTemplateIfMissing(path, "missing-template", 0o644)
	if err == nil {
		t.Fatalf("expected error for missing template")
	}
}

func TestWriteTemplateIfMissingStatError(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	path := filepath.Join(file, "config.toml")
	if err := writeTemplateIfMissing(path, "config.toml", 0o644); err == nil {
		t.Fatalf("expected error for stat failure")
	}
}

func TestEnsureGitignoreReadError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	err := ensureGitignore(path, "# >>> agent-layer\n# <<< agent-layer\n")
	if err == nil {
		t.Fatalf("expected error for directory path")
	}
}

func TestCreateDirsError(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := createDirs(file); err == nil {
		t.Fatalf("expected error")
	}
}

func TestUpdateGitignoreMissingBlock(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".agent-layer"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := updateGitignore(root); err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunStepsError(t *testing.T) {
	err := runSteps([]func() error{
		func() error { return fmt.Errorf("boom") },
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunMissingRoot(t *testing.T) {
	if err := Run(""); err == nil {
		t.Fatalf("expected error for missing root")
	}
}
