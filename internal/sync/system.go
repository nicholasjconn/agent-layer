package sync

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/conn-castle/agent-layer/internal/fsutil"
)

// System abstracts system-level operations to enable dependency injection in sync logic.
type System interface {
	LookPath(file string) (string, error)
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	WriteFileAtomic(filename string, data []byte, perm os.FileMode) error
	MarshalIndent(v any, prefix, indent string) ([]byte, error)
	ReadFile(name string) ([]byte, error)
}

// RealSystem implements System using actual system calls.
type RealSystem struct{}

// LookPath searches for an executable named file in the directories named by the PATH environment variable.
func (RealSystem) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// Stat returns a FileInfo describing the named file.
func (RealSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// MkdirAll creates a directory named path, along with any necessary parents.
func (RealSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// WriteFileAtomic writes data to a file atomically by writing to a temp file and renaming.
func (RealSystem) WriteFileAtomic(filename string, data []byte, perm os.FileMode) error {
	return fsutil.WriteFileAtomic(filename, data, perm)
}

// MarshalIndent returns the JSON encoding of v with indentation.
func (RealSystem) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// ReadFile reads the named file and returns the contents.
func (RealSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
