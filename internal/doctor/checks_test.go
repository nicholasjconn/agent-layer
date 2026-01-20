package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasjconn/agent-layer/internal/config"
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

func TestCheckSecretsUsesRequiredEnvVars(t *testing.T) {
	t.Setenv("HEADER_TOKEN", "present")

	cfg := &config.ProjectConfig{
		Config: config.Config{
			MCP: config.MCPConfig{
				Servers: []config.MCPServer{
					{
						ID:      "demo",
						URL:     "https://example.com/${URL_TOKEN}",
						Command: "run-${CMD_TOKEN}",
						Args:    []string{"--token", "${ARG_TOKEN}"},
						Headers: map[string]string{"Authorization": "Bearer ${HEADER_TOKEN}"},
						Env:     map[string]string{"API_KEY": "${ENV_TOKEN}"},
					},
				},
			},
		},
		Env: map[string]string{
			"ARG_TOKEN": "set",
		},
	}

	results := CheckSecrets(cfg)
	expected := map[string]Status{
		"Missing secret: URL_TOKEN":                 StatusFail,
		"Missing secret: CMD_TOKEN":                 StatusFail,
		"Missing secret: ENV_TOKEN":                 StatusFail,
		"Secret found in .env: ARG_TOKEN":           StatusOK,
		"Secret found in environment: HEADER_TOKEN": StatusOK,
	}

	for msg, status := range expected {
		found := false
		for _, result := range results {
			if result.Message != msg {
				continue
			}
			if result.Status != status {
				t.Fatalf("expected %q status %s, got %s", msg, status, result.Status)
			}
			found = true
			break
		}
		if !found {
			t.Fatalf("expected result message %q", msg)
		}
	}
}
