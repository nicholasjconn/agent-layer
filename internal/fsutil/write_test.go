package fsutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFileAtomicCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	err := WriteFileAtomic(path, []byte("hello"), 0644)
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
}

func TestWriteFileAtomicOverwritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	err := os.WriteFile(path, []byte("old"), 0600)
	require.NoError(t, err)

	err = WriteFileAtomic(path, []byte("new"), 0600)
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "new", string(data))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestWriteFileAtomicFailures(t *testing.T) {
	t.Run("invalid dir", func(t *testing.T) {
		err := WriteFileAtomic("/invalid/path/config.toml", []byte("data"), 0644)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create temp file")
	})

	t.Run("target is directory", func(t *testing.T) {
		dir := t.TempDir()
		target := filepath.Join(dir, "target")
		require.NoError(t, os.Mkdir(target, 0755))

		err := WriteFileAtomic(target, []byte("data"), 0644)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rename temp file")
	})
}

func TestSyncDir(t *testing.T) {
	err := syncDir("/invalid/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "open dir")
}

func TestSyncDir_Success(t *testing.T) {
	dir := t.TempDir()
	err := syncDir(dir)
	assert.NoError(t, err)
}

func TestWriteFileAtomic_NoPermission(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "file.txt")
	// Subdir doesn't exist, and we won't be able to create temp file
	err := WriteFileAtomic(path, []byte("data"), 0644)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create temp file")
}

func TestWriteFileAtomic_TempFileCleanup(t *testing.T) {
	// Test that temp file is cleaned up when rename fails
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	// Create target as directory so rename fails
	require.NoError(t, os.Mkdir(target, 0755))

	err := WriteFileAtomic(target, []byte("data"), 0644)
	assert.Error(t, err)

	// Verify no temp files remain
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	for _, entry := range entries {
		if entry.Name() != "target" {
			t.Errorf("unexpected file found: %s", entry.Name())
		}
	}
}

func TestSyncDir_File(t *testing.T) {
	// Try to sync a file instead of a directory
	dir := t.TempDir()
	filePath := filepath.Join(dir, "file.txt")
	require.NoError(t, os.WriteFile(filePath, []byte("test"), 0644))

	// This may or may not error depending on OS
	// The point is to exercise the code path
	_ = syncDir(filePath)
}
