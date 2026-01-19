package config

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// LoadInstructions reads .agent-layer/instructions/*.md in lexicographic order.
func LoadInstructions(dir string) ([]InstructionFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("missing instructions directory %s: %w", dir, err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".md") {
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("no instruction files found in %s", dir)
	}

	sort.Strings(names)

	files := make([]InstructionFile, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed reading instruction %s: %w", path, err)
		}
		data = bytes.TrimPrefix(data, utf8BOM)
		files = append(files, InstructionFile{
			Name:    name,
			Content: string(data),
		})
	}

	return files, nil
}

// WalkInstructionFiles is a helper to walk instruction files in a directory.
func WalkInstructionFiles(dir string, fn func(path string, entry fs.DirEntry) error) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if err := fn(path, entry); err != nil {
			return err
		}
	}
	return nil
}
