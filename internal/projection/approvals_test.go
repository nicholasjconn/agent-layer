package projection

import (
	"testing"

	"github.com/conn-castle/agent-layer/internal/config"
)

func TestBuildApprovals(t *testing.T) {
	cfg := config.Config{
		Approvals: config.ApprovalsConfig{Mode: "commands"},
	}
	result := BuildApprovals(cfg, []string{"git status"})
	if !result.AllowCommands || result.AllowMCP {
		t.Fatalf("unexpected approvals flags: %+v", result)
	}
	if len(result.Commands) != 1 || result.Commands[0] != "git status" {
		t.Fatalf("unexpected commands: %+v", result.Commands)
	}
}
