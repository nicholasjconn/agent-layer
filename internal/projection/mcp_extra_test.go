package projection

import (
	"testing"

	"github.com/nicholasjconn/agent-layer/internal/config"
)

func TestResolveMCPServersUnsupportedTransport(t *testing.T) {
	enabled := true
	servers := []config.MCPServer{
		{ID: "bad", Enabled: &enabled, Transport: "grpc"},
	}
	_, err := ResolveMCPServers(servers, map[string]string{}, "gemini", nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolveMCPServersHeaderMissingEnv(t *testing.T) {
	enabled := true
	servers := []config.MCPServer{
		{
			ID:        "http",
			Enabled:   &enabled,
			Clients:   []string{"gemini"},
			Transport: "http",
			URL:       "https://example.com?token=${TOKEN}",
			Headers:   map[string]string{"Authorization": "Bearer ${TOKEN}"},
		},
	}
	_, err := ResolveMCPServers(servers, map[string]string{}, "gemini", nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolveMCPServersCommandMissingEnv(t *testing.T) {
	enabled := true
	servers := []config.MCPServer{
		{
			ID:        "stdio",
			Enabled:   &enabled,
			Clients:   []string{"gemini"},
			Transport: "stdio",
			Command:   "${CMD}",
		},
	}
	_, err := ResolveMCPServers(servers, map[string]string{}, "gemini", nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}
