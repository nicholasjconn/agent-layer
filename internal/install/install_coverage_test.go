package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteVersionFile_EmptyExisting(t *testing.T) {
	root := t.TempDir()
	inst := &installer{root: root, pinVersion: "1.0.0"}
	alDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(alDir, 0o755); err != nil {
		t.Fatal(err)
	}
	pinFile := filepath.Join(alDir, "al.version")
	if err := os.WriteFile(pinFile, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	err := inst.writeVersionFile()
	if err == nil || !strings.Contains(err.Error(), "existing pin file") {
		t.Fatalf("expected error for empty pin file, got %v", err)
	}
}

func TestWriteVersionFile_ReadError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping read error test as root")
	}
	root := t.TempDir()
	inst := &installer{root: root, pinVersion: "1.0.0"}
	alDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(alDir, 0o755); err != nil {
		t.Fatal(err)
	}
	pinFile := filepath.Join(alDir, "al.version")
	// Make a directory to force read error
	if err := os.Mkdir(pinFile, 0o755); err != nil {
		t.Fatal(err)
	}

	err := inst.writeVersionFile()
	if err == nil || !strings.Contains(err.Error(), "failed to read") {
		t.Fatalf("expected error for read failure, got %v", err)
	}
}

func TestWriteVersionFile_MkdirError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping mkdir error test as root")
	}
	root := t.TempDir()
	inst := &installer{root: root, pinVersion: "1.0.0"}

	// Create parent dir as read-only to force MkdirAll to fail when creating subdir
	parent := filepath.Join(root, "parent")
	if err := os.Mkdir(parent, 0o555); err != nil {
		t.Fatal(err)
	}
	// Use a read-only directory as root to force MkdirAll failure when creating .agent-layer.
	inst.root = parent

	err := inst.writeVersionFile()
	if err == nil || !strings.Contains(err.Error(), "failed to create directory") {
		t.Fatalf("expected error for mkdir failure, got %v", err)
	}
}

func TestWriteVersionFile_PromptError(t *testing.T) {
	root := t.TempDir()
	inst := &installer{
		root:       root,
		overwrite:  true, // Must be true to trigger prompt
		pinVersion: "1.0.0",
		prompt: func(path string) (bool, error) {
			return false, fmt.Errorf("prompt failed")
		},
	}
	alDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(alDir, 0o755); err != nil {
		t.Fatal(err)
	}
	pinFile := filepath.Join(alDir, "al.version")
	if err := os.WriteFile(pinFile, []byte("0.9.0"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := inst.writeVersionFile()
	if err == nil || !strings.Contains(err.Error(), "prompt failed") {
		t.Fatalf("expected error from prompt, got %v", err)
	}
}

func TestWriteVersionFile_NoOverwrite(t *testing.T) {
	root := t.TempDir()
	inst := &installer{
		root:       root,
		overwrite:  true, // Must be true to trigger prompt
		pinVersion: "1.0.0",
		prompt: func(path string) (bool, error) {
			return false, nil
		},
	}
	alDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(alDir, 0o755); err != nil {
		t.Fatal(err)
	}
	pinFile := filepath.Join(alDir, "al.version")
	if err := os.WriteFile(pinFile, []byte("0.9.0"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := inst.writeVersionFile()
	if err != nil {
		t.Fatalf("expected nil error for declined overwrite, got %v", err)
	}
	if len(inst.diffs) != 1 {
		t.Errorf("expected 1 diff, got %d", len(inst.diffs))
	}
}

func TestShouldOverwrite_PromptError(t *testing.T) {
	inst := &installer{
		overwrite: true,
		prompt: func(path string) (bool, error) {
			return false, fmt.Errorf("prompt error")
		},
	}
	_, err := inst.shouldOverwrite("path")
	if err == nil {
		t.Fatal("expected error from prompt")
	}
}

func TestShouldOverwrite_Force(t *testing.T) {
	inst := &installer{overwrite: true, force: true}
	overwrite, err := inst.shouldOverwrite("path")
	if err != nil {
		t.Fatal(err)
	}
	if !overwrite {
		t.Error("expected true when forced")
	}
}

func TestShouldOverwrite_MissingPrompt(t *testing.T) {
	inst := &installer{overwrite: true, prompt: nil}
	_, err := inst.shouldOverwrite("path")
	if err == nil {
		t.Fatal("expected error when prompt handler is missing")
	}
}

func alwaysOverwrite(_ string) (bool, error) {
	return true, nil
}

func TestFileMatchesTemplate_ReadError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "unreadable")
	if err := os.WriteFile(path, []byte("content"), 0o200); err != nil { // Write only
		t.Fatalf("write file: %v", err)
	}
	// Try to read it - might succeed as root/owner on some systems, but let's try.
	// A better way is to pass a directory path as a file.
	dirPath := filepath.Join(root, "dir")
	if err := os.Mkdir(dirPath, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Reading a directory as a file fails on many systems, or use a non-existent file.
	// But fileMatchesTemplate does os.ReadFile(path). If path doesn't exist, it returns error.
	_, err := fileMatchesTemplate("non-existent-file", "config.toml")
	if err == nil {
		t.Fatalf("expected error for non-existent file")
	}
}

func TestFileMatchesTemplate_TemplateReadError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "file")
	if err := os.WriteFile(path, []byte("content"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Pass an invalid template path
	_, err := fileMatchesTemplate(path, "non-existent-template")
	if err == nil {
		t.Fatalf("expected error for non-existent template")
	}
}

func TestWriteGitignoreBlock_TemplateReadError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")

	// Pass an invalid template path
	err := writeGitignoreBlock(path, "non-existent-template", 0o644, nil, nil)
	if err == nil {
		t.Fatalf("expected error for non-existent template")
	}
}

func TestWriteGitignoreBlock_FileReadError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")

	// Create a directory at the path to force ReadFile to fail (is a directory)
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	err := writeGitignoreBlock(path, "gitignore.block", 0o644, nil, nil)
	if err == nil {
		t.Fatalf("expected error when reading directory as file")
	}
}

func TestWriteTemplateFile_TemplateReadError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "file")

	// Pass an invalid template path
	err := writeTemplateFile(path, "non-existent-template", 0o644, nil, nil)
	if err == nil {
		t.Fatalf("expected error for non-existent template")
	}
}

func TestWriteTemplateFile_StatError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping stat error test as root")
	}
	root := t.TempDir()
	sub := filepath.Join(root, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	target := filepath.Join(sub, "target")
	// Make sub unreadable to force stat error on target
	if err := os.Chmod(sub, 0o000); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	defer func() { _ = os.Chmod(sub, 0o755) }()

	err := writeTemplateFile(target, "config.toml", 0o644, nil, nil)
	if err == nil {
		t.Fatalf("expected error for stat failure")
	}
}

func TestWriteTemplateFile_MkdirError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping mkdir error test as root")
	}
	root := t.TempDir()
	// Create a file where the directory should be
	path := filepath.Join(root, "file", "target")
	if err := os.WriteFile(filepath.Join(root, "file"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	err := writeTemplateFile(path, "config.toml", 0o644, nil, nil)
	if err == nil {
		t.Fatalf("expected error for mkdir failure")
	}
}

func TestWriteTemplateFile_WriteError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "file")
	// Create a directory where the file should be to cause WriteFile to fail
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	err := writeTemplateFile(path, "config.toml", 0o644, alwaysOverwrite, nil)
	if err == nil {
		t.Fatalf("expected error for write failure")
	}
}

func TestWriteTemplateFile_FileMatchesError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping file read error test as root")
	}
	root := t.TempDir()
	path := filepath.Join(root, "file")
	if err := os.WriteFile(path, []byte("x"), 0o000); err != nil {
		t.Fatalf("write file: %v", err)
	}
	defer func() { _ = os.Chmod(path, 0o644) }()

	// fileMatchesTemplate should fail to read 'file'
	err := writeTemplateFile(path, "config.toml", 0o644, nil, nil)
	if err == nil {
		t.Fatalf("expected error for fileMatchesTemplate failure")
	}
}

func TestWriteGitignoreBlock_WriteError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	// Create a directory where the file should be to cause WriteFile to fail
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	err := writeGitignoreBlock(path, "gitignore.block", 0o644, alwaysOverwrite, nil)
	if err == nil {
		t.Fatalf("expected error for write failure")
	}
}

func TestWriteTemplateFiles_Error(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping write error test as root")
	}
	root := t.TempDir()
	inst := &installer{root: root}

	// Create .agent-layer directory and make it read-only to force write failure
	alDir := filepath.Join(root, ".agent-layer")
	if err := os.Mkdir(alDir, 0o555); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defer func() { _ = os.Chmod(alDir, 0o755) }()

	// writeTemplateFiles tries to write .agent-layer/config.toml, which should fail
	if err := inst.writeTemplateFiles(); err == nil {
		t.Fatalf("expected error from writeTemplateFiles")
	}
}

func TestWriteTemplateFiles_GitignoreBlockError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping write error test as root")
	}
	root := t.TempDir()
	inst := &installer{root: root}

	// Create .agent-layer directory
	alDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(alDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create gitignore.block as a directory to force writeGitignoreBlock to fail
	// (writeGitignoreBlock tries to ReadFile or WriteFile, failing on dir)
	blockPath := filepath.Join(alDir, "gitignore.block")
	if err := os.Mkdir(blockPath, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// writeTemplateFiles should succeed for other files but fail on gitignore.block
	if err := inst.writeTemplateFiles(); err == nil {
		t.Fatalf("expected error from writeTemplateFiles (gitignore block)")
	}
}

func TestWriteTemplateDirs_Error(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping write error test as root")
	}
	root := t.TempDir()
	inst := &installer{root: root}

	// Create .agent-layer as read-only. writeTemplateDirs tries to create .agent-layer/instructions inside it.
	alDir := filepath.Join(root, ".agent-layer")
	if err := os.Mkdir(alDir, 0o555); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defer func() { _ = os.Chmod(alDir, 0o755) }()

	// writeTemplateDirs tries to write to .agent-layer/instructions
	if err := inst.writeTemplateDirs(); err == nil {
		t.Fatalf("expected error from writeTemplateDirs")
	}
}

func TestRun_Error(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping write error test as root")
	}
	root := t.TempDir()

	// Make docs read-only so creating docs/agent-layer fails.
	docsDir := filepath.Join(root, "docs")
	if err := os.Mkdir(docsDir, 0o555); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defer func() { _ = os.Chmod(docsDir, 0o755) }()

	if err := Run(root, Options{}); err == nil {
		t.Fatalf("expected error from Run")
	}
}

func TestUpdateGitignoreContent_Empty(t *testing.T) {
	content := ""
	block := "# >>> agent-layer\nblock\n# <<< agent-layer\n"
	updated := updateGitignoreContent(content, block)
	expected := "# >>> agent-layer\nblock\n# <<< agent-layer\n"
	if updated != expected {
		t.Fatalf("expected %q, got %q", expected, updated)
	}
}

func TestRenderGitignoreBlock_Empty(t *testing.T) {
	// When block is empty, TrimRight("", "\n") returns "", and Split("", "\n") returns [""]
	// So len(lines) is 1, not 0. The len(lines) == 0 branch is unreachable with normal strings.
	result := renderGitignoreBlock("")
	if result == "" {
		t.Fatalf("expected non-empty result")
	}
}

func TestRenderGitignoreBlock_Normal(t *testing.T) {
	block := "# >>> agent-layer\nentry1\nentry2\n# <<< agent-layer\n"
	result := renderGitignoreBlock(block)

	// Should contain the hash prefix
	if !strings.Contains(result, gitignoreHashPrefix) {
		t.Fatalf("expected hash line in result")
	}

	// Should preserve original content
	if !strings.Contains(result, "entry1") || !strings.Contains(result, "entry2") {
		t.Fatalf("expected original content preserved")
	}
}

func TestUpdateGitignoreContent_EmptyExisting(t *testing.T) {
	content := ""
	block := "# >>> agent-layer\nblock\n# <<< agent-layer\n"
	updated := updateGitignoreContent(content, block)

	if !strings.Contains(updated, "block") {
		t.Fatalf("expected block to be added")
	}
}

func TestUpdateGitignoreContent_NoTrailingNewline(t *testing.T) {
	content := "existing"
	block := "# >>> agent-layer\nblock\n# <<< agent-layer\n"
	updated := updateGitignoreContent(content, block)

	if !strings.Contains(updated, "existing") {
		t.Fatalf("expected existing content preserved")
	}
	if !strings.Contains(updated, "block") {
		t.Fatalf("expected block to be appended")
	}
}

func TestEnsureGitignore_WriteErrorAfterRead(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping write error test as root")
	}
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")

	// Write initial content so read succeeds
	if err := os.WriteFile(path, []byte("existing\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Make directory read-only so WriteFileAtomic fails
	if err := os.Chmod(root, 0o555); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	defer func() { _ = os.Chmod(root, 0o755) }()

	block := "# >>> agent-layer\nblock\n# <<< agent-layer\n"
	err := ensureGitignore(path, block)
	if err == nil {
		t.Fatalf("expected error for write failure")
	}
}

func TestEnsureGitignore_WriteErrorNewFile(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping write error test as root")
	}
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")

	// Make directory read-only so WriteFileAtomic fails for new file
	if err := os.Chmod(root, 0o555); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	defer func() { _ = os.Chmod(root, 0o755) }()

	block := "# >>> agent-layer\nblock\n# <<< agent-layer\n"
	err := ensureGitignore(path, block)
	if err == nil {
		t.Fatalf("expected error for write failure")
	}
}

func TestWriteGitignoreBlock_MkdirError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping mkdir error test as root")
	}
	root := t.TempDir()
	// Create a file where the directory should be
	parentFile := filepath.Join(root, ".agent-layer")
	if err := os.WriteFile(parentFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	path := filepath.Join(parentFile, "gitignore.block")
	err := writeGitignoreBlock(path, "gitignore.block", 0o644, nil, nil)
	if err == nil {
		t.Fatalf("expected error for mkdir failure")
	}
}

func TestWriteGitignoreBlock_NewFileWriteError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping write error test as root")
	}
	root := t.TempDir()
	agentLayerDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(agentLayerDir, 0o555); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defer func() { _ = os.Chmod(agentLayerDir, 0o755) }()

	path := filepath.Join(agentLayerDir, "gitignore.block")
	err := writeGitignoreBlock(path, "gitignore.block", 0o644, nil, nil)
	if err == nil {
		t.Fatalf("expected error for write failure")
	}
}

func TestWriteGitignoreBlock_OverwriteWriteError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("skipping write error test as root")
	}
	root := t.TempDir()
	agentLayerDir := filepath.Join(root, ".agent-layer")
	if err := os.MkdirAll(agentLayerDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	path := filepath.Join(agentLayerDir, "gitignore.block")
	// Write some content so file exists
	if err := os.WriteFile(path, []byte("# >>> agent-layer\nold\n# <<< agent-layer\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Make directory read-only so WriteFileAtomic fails
	if err := os.Chmod(agentLayerDir, 0o555); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	defer func() { _ = os.Chmod(agentLayerDir, 0o755) }()

	// Use overwrite=true to trigger the overwrite path
	err := writeGitignoreBlock(path, "gitignore.block", 0o644, alwaysOverwrite, nil)
	if err == nil {
		t.Fatalf("expected error for write failure during overwrite")
	}
}

func TestWarnDifferences_RelPathError(t *testing.T) {
	// This test verifies the fallback path when filepath.Rel fails.
	// filepath.Rel can fail if paths are on different volumes (Windows)
	// or if there's some other OS-specific issue.
	// We use a path that cannot be made relative to root to trigger fallback.
	inst := &installer{
		root:      "/some/root",
		overwrite: false,
		diffs:     []string{"/completely/different/path/file.txt"},
	}

	// warnDifferences should not panic and should use absolute path as fallback
	// when Rel fails (on some systems) or succeeds with a long relative path.
	inst.warnDifferences()
	// Success = no panic, the function handled the path
}
