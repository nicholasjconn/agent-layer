package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckStructure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "doctor-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Test missing directories
	results := CheckStructure(tmpDir)
	failCount := 0
	for _, r := range results {
		if r.Status == StatusFail {
			failCount++
		}
	}
	if failCount != 2 {
		t.Errorf("Expected 2 failures for empty directory, got %d", failCount)
	}

	// Test existing directories
	if err := os.Mkdir(filepath.Join(tmpDir, ".agent-layer"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "docs/agent-layer"), 0755); err != nil {
		t.Fatal(err)
	}
	results = CheckStructure(tmpDir)
	for _, r := range results {
		if r.Status != StatusOK {
			t.Errorf("Expected OK status for existing directory %s, got %s", r.CheckName, r.Status)
		}
	}
}

func TestFindSecrets(t *testing.T) {
	secrets := make(map[string]bool)
	findSecrets("Bearer ${MY_TOKEN} and ${ANOTHER_ONE}", secrets)

	if !secrets["MY_TOKEN"] {
		t.Error("Expected MY_TOKEN to be found")
	}
	if !secrets["ANOTHER_ONE"] {
		t.Error("Expected ANOTHER_ONE to be found")
	}
	if len(secrets) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(secrets))
	}
}
