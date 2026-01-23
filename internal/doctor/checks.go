package doctor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/conn-castle/agent-layer/internal/config"
)

// CheckStructure verifies that the required project directories exist.
func CheckStructure(root string) []Result {
	var results []Result
	paths := []string{".agent-layer", "docs/agent-layer"}

	for _, p := range paths {
		fullPath := filepath.Join(root, p)
		info, err := os.Stat(fullPath)
		if err != nil {
			results = append(results, Result{
				Status:         StatusFail,
				CheckName:      "Structure",
				Message:        fmt.Sprintf("Missing required directory: %s", p),
				Recommendation: "Run `al init` to initialize the repository.",
			})
			continue
		}
		if !info.IsDir() {
			results = append(results, Result{
				Status:         StatusFail,
				CheckName:      "Structure",
				Message:        fmt.Sprintf("%s exists but is not a directory", p),
				Recommendation: "Ensure the path is a directory or re-run `al init`.",
			})
			continue
		}
		results = append(results, Result{
			Status:    StatusOK,
			CheckName: "Structure",
			Message:   fmt.Sprintf("Directory exists: %s", p),
		})
	}
	return results
}

// CheckConfig validates that the configuration file can be loaded and parsed.
func CheckConfig(root string) ([]Result, *config.ProjectConfig) {
	var results []Result
	cfg, err := config.LoadProjectConfig(root)
	if err != nil {
		results = append(results, Result{
			Status:         StatusFail,
			CheckName:      "Config",
			Message:        fmt.Sprintf("Failed to load configuration: %v", err),
			Recommendation: "Check .agent-layer/config.toml for syntax errors.",
		})
		return results, nil
	}

	results = append(results, Result{
		Status:    StatusOK,
		CheckName: "Config",
		Message:   "Configuration loaded successfully",
	})
	return results, cfg
}

// CheckSecrets scans the configuration for missing environment variables.
func CheckSecrets(cfg *config.ProjectConfig) []Result {
	var results []Result
	required := config.RequiredEnvVarsForMCPServers(cfg.Config.MCP.Servers)

	// Scan .env for missing values
	for _, secret := range required {
		val, ok := cfg.Env[secret]
		if !ok || val == "" {
			// Check if it's in the actual environment
			if os.Getenv(secret) == "" {
				results = append(results, Result{
					Status:         StatusFail,
					CheckName:      "Secrets",
					Message:        fmt.Sprintf("Missing secret: %s", secret),
					Recommendation: fmt.Sprintf("Add %s to .agent-layer/.env or your environment.", secret),
				})
			} else {
				results = append(results, Result{
					Status:    StatusOK,
					CheckName: "Secrets",
					Message:   fmt.Sprintf("Secret found in environment: %s", secret),
				})
			}
		} else {
			results = append(results, Result{
				Status:    StatusOK,
				CheckName: "Secrets",
				Message:   fmt.Sprintf("Secret found in .env: %s", secret),
			})
		}
	}

	if len(required) == 0 {
		results = append(results, Result{
			Status:    StatusOK,
			CheckName: "Secrets",
			Message:   "No required secrets detected in configuration.",
		})
	}

	return results
}

// CheckAgents reports which agents are enabled or disabled.
func CheckAgents(cfg *config.ProjectConfig) []Result {
	var results []Result
	agents := []struct {
		Name    string
		Enabled *bool
	}{
		{"Gemini", cfg.Config.Agents.Gemini.Enabled},
		{"Claude", cfg.Config.Agents.Claude.Enabled},
		{"Codex", cfg.Config.Agents.Codex.Enabled},
		{"VSCode", cfg.Config.Agents.VSCode.Enabled},
		{"Antigravity", cfg.Config.Agents.Antigravity.Enabled},
	}

	for _, a := range agents {
		if a.Enabled != nil && *a.Enabled {
			results = append(results, Result{
				Status:    StatusOK,
				CheckName: "Agents",
				Message:   fmt.Sprintf("Agent enabled: %s", a.Name),
			})
		} else {
			results = append(results, Result{
				Status:    StatusWarn,
				CheckName: "Agents",
				Message:   fmt.Sprintf("Agent disabled: %s", a.Name),
			})
		}
	}
	return results
}
