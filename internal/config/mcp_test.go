package config

import "testing"

func TestAppliesToClient(t *testing.T) {
	server := MCPServer{Clients: []string{"gemini", "codex"}}
	if !server.AppliesToClient("gemini") {
		t.Fatalf("expected gemini to apply")
	}
	if server.AppliesToClient("vscode") {
		t.Fatalf("expected vscode not to apply")
	}

	empty := MCPServer{}
	if !empty.AppliesToClient("any") {
		t.Fatalf("expected empty clients to apply")
	}
}
