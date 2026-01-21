package sync

import (
	"strings"
	"testing"

	"github.com/nicholasjconn/agent-layer/internal/config"
)

func ptr(i int) *int {
	return &i
}

func ptrBool(b bool) *bool {
	return &b
}

func TestCheckInstructionSize_BelowThreshold(t *testing.T) {
	instructions := []config.InstructionFile{
		{Name: "test.md", Content: strings.Repeat("a", 100)}, // 100 chars = 25 tokens at 4 chars/token
	}
	cfg := config.WarningsConfig{
		InstructionTokenThreshold: ptr(50),
	}

	warning := CheckInstructionSize(instructions, cfg)
	if warning != nil {
		t.Fatalf("expected no warning, got: %s", warning.Message)
	}
}

func TestCheckInstructionSize_AboveThreshold(t *testing.T) {
	instructions := []config.InstructionFile{
		{Name: "test.md", Content: strings.Repeat("a", 300)}, // 300 chars = 75 tokens at 4 chars/token
	}
	cfg := config.WarningsConfig{
		InstructionTokenThreshold: ptr(50),
	}

	warning := CheckInstructionSize(instructions, cfg)
	if warning == nil {
		t.Fatalf("expected warning, got nil")
	}
	if !strings.Contains(warning.Message, "exceed") {
		t.Fatalf("expected warning message to contain 'exceed', got: %s", warning.Message)
	}
}

func TestCheckInstructionSize_AtBoundary(t *testing.T) {
	// Exactly at threshold should NOT warn (only > threshold warns)
	instructions := []config.InstructionFile{
		{Name: "test.md", Content: strings.Repeat("a", 200)}, // 200 chars = 50 tokens at 4 chars/token
	}
	cfg := config.WarningsConfig{
		InstructionTokenThreshold: ptr(50),
	}

	warning := CheckInstructionSize(instructions, cfg)
	if warning != nil {
		t.Fatalf("expected no warning at boundary, got: %s", warning.Message)
	}
}

func TestCheckInstructionSize_DisabledThreshold(t *testing.T) {
	instructions := []config.InstructionFile{
		{Name: "test.md", Content: strings.Repeat("a", 1000000)}, // Large content
	}
	cfg := config.WarningsConfig{
		InstructionTokenThreshold: nil, // Disabled
	}

	warning := CheckInstructionSize(instructions, cfg)
	if warning != nil {
		t.Fatalf("expected no warning when threshold disabled, got: %s", warning.Message)
	}
}

func TestEstimatedCharsPerToken_Constant(t *testing.T) {
	// Verify the constant is set to the expected value
	if EstimatedCharsPerToken != 4 {
		t.Fatalf("expected EstimatedCharsPerToken to be 4, got %d", EstimatedCharsPerToken)
	}
}

func TestCheckMCPServerCount_BelowThreshold(t *testing.T) {
	servers := []config.MCPServer{
		{ID: "s1", Enabled: ptrBool(true)},
		{ID: "s2", Enabled: ptrBool(true)},
	}
	cfg := config.WarningsConfig{
		MCPServerThreshold: ptr(5),
	}

	warnings := CheckMCPServerCount(servers, []string{"claude"}, cfg)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

func TestCheckMCPServerCount_AboveThreshold(t *testing.T) {
	servers := []config.MCPServer{
		{ID: "s1", Enabled: ptrBool(true)},
		{ID: "s2", Enabled: ptrBool(true)},
		{ID: "s3", Enabled: ptrBool(true)},
		{ID: "s4", Enabled: ptrBool(true)},
		{ID: "s5", Enabled: ptrBool(true)},
		{ID: "s6", Enabled: ptrBool(true)}, // 6 servers > threshold of 5
	}
	cfg := config.WarningsConfig{
		MCPServerThreshold: ptr(5),
	}

	warnings := CheckMCPServerCount(servers, []string{"claude"}, cfg)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if !strings.Contains(warnings[0].Message, "claude") {
		t.Fatalf("expected warning to mention client name, got: %s", warnings[0].Message)
	}
}

func TestCheckMCPServerCount_PerClient(t *testing.T) {
	// Server s1-s6 apply to claude only, s7-s8 apply to gemini only
	servers := []config.MCPServer{
		{ID: "s1", Enabled: ptrBool(true), Clients: []string{"claude"}},
		{ID: "s2", Enabled: ptrBool(true), Clients: []string{"claude"}},
		{ID: "s3", Enabled: ptrBool(true), Clients: []string{"claude"}},
		{ID: "s4", Enabled: ptrBool(true), Clients: []string{"claude"}},
		{ID: "s5", Enabled: ptrBool(true), Clients: []string{"claude"}},
		{ID: "s6", Enabled: ptrBool(true), Clients: []string{"claude"}}, // 6 for claude
		{ID: "s7", Enabled: ptrBool(true), Clients: []string{"gemini"}},
		{ID: "s8", Enabled: ptrBool(true), Clients: []string{"gemini"}}, // 2 for gemini
	}
	cfg := config.WarningsConfig{
		MCPServerThreshold: ptr(5),
	}

	warnings := CheckMCPServerCount(servers, []string{"claude", "gemini"}, cfg)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning (for claude only), got %d", len(warnings))
	}
	if !strings.Contains(warnings[0].Message, "claude") {
		t.Fatalf("expected warning for claude, got: %s", warnings[0].Message)
	}
}

func TestCheckMCPServerCount_DisabledServers(t *testing.T) {
	servers := []config.MCPServer{
		{ID: "s1", Enabled: ptrBool(true)},
		{ID: "s2", Enabled: ptrBool(true)},
		{ID: "s3", Enabled: ptrBool(false)}, // Disabled
		{ID: "s4", Enabled: ptrBool(false)}, // Disabled
		{ID: "s5", Enabled: ptrBool(false)}, // Disabled
		{ID: "s6", Enabled: ptrBool(false)}, // Disabled
	}
	cfg := config.WarningsConfig{
		MCPServerThreshold: ptr(2), // Only 2 are actually enabled
	}

	warnings := CheckMCPServerCount(servers, []string{"claude"}, cfg)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings (disabled servers don't count), got %d", len(warnings))
	}
}

func TestCheckMCPServerCount_Disabled(t *testing.T) {
	servers := []config.MCPServer{
		{ID: "s1", Enabled: ptrBool(true)},
		{ID: "s2", Enabled: ptrBool(true)},
		{ID: "s3", Enabled: ptrBool(true)},
		{ID: "s4", Enabled: ptrBool(true)},
		{ID: "s5", Enabled: ptrBool(true)},
		{ID: "s6", Enabled: ptrBool(true)},
	}
	cfg := config.WarningsConfig{
		MCPServerThreshold: nil, // Disabled
	}

	warnings := CheckMCPServerCount(servers, []string{"claude"}, cfg)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings when threshold disabled, got %d", len(warnings))
	}
}

func TestCheckMCPServerCount_NilEnabled(t *testing.T) {
	servers := []config.MCPServer{
		{ID: "s1", Enabled: nil}, // nil Enabled should be treated as disabled
		{ID: "s2", Enabled: nil},
		{ID: "s3", Enabled: nil},
	}
	cfg := config.WarningsConfig{
		MCPServerThreshold: ptr(1),
	}

	warnings := CheckMCPServerCount(servers, []string{"claude"}, cfg)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings (nil Enabled = disabled), got %d", len(warnings))
	}
}

func TestEnabledClientNames(t *testing.T) {
	agents := config.AgentsConfig{
		Gemini:      config.AgentConfig{Enabled: ptrBool(true)},
		Claude:      config.AgentConfig{Enabled: ptrBool(true)},
		Codex:       config.CodexConfig{Enabled: ptrBool(false)},
		VSCode:      config.AgentConfig{Enabled: ptrBool(true)},
		Antigravity: config.AgentConfig{Enabled: nil},
	}

	names := EnabledClientNames(agents)
	expected := []string{"claude", "gemini", "vscode"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d clients, got %d: %v", len(expected), len(names), names)
	}
	for i, name := range names {
		if name != expected[i] {
			t.Fatalf("expected %s at position %d, got %s", expected[i], i, name)
		}
	}
}

func TestEnabledClientNames_AllEnabled(t *testing.T) {
	agents := config.AgentsConfig{
		Gemini:      config.AgentConfig{Enabled: ptrBool(true)},
		Claude:      config.AgentConfig{Enabled: ptrBool(true)},
		Codex:       config.CodexConfig{Enabled: ptrBool(true)},
		VSCode:      config.AgentConfig{Enabled: ptrBool(true)},
		Antigravity: config.AgentConfig{Enabled: ptrBool(true)},
	}

	names := EnabledClientNames(agents)
	expected := []string{"antigravity", "claude", "codex", "gemini", "vscode"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d clients, got %d: %v", len(expected), len(names), names)
	}
	for i, name := range names {
		if name != expected[i] {
			t.Fatalf("expected %s at position %d, got %s", expected[i], i, name)
		}
	}
}

func TestEnabledClientNames_NoneEnabled(t *testing.T) {
	agents := config.AgentsConfig{}

	names := EnabledClientNames(agents)
	if len(names) != 0 {
		t.Fatalf("expected 0 clients, got %d: %v", len(names), names)
	}
}
