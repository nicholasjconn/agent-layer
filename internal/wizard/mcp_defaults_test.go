package wizard

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nicholasjconn/agent-layer/internal/config"
)

func TestMissingDefaultMCPServers(t *testing.T) {
	defaults := []DefaultMCPServer{
		{ID: "github"},
		{ID: "context7"},
		{ID: "tavily"},
	}
	servers := []config.MCPServer{
		{ID: "github"},
		{ID: "tavily"},
	}

	missing := missingDefaultMCPServers(defaults, servers)

	assert.Equal(t, []string{"context7"}, missing)
}

func TestAppendMissingDefaultMCPServers(t *testing.T) {
	content := "[mcp]\n"
	missing := []string{"github"}

	updated, err := appendMissingDefaultMCPServers(content, missing)
	require.NoError(t, err)

	assert.Contains(t, updated, "[[mcp.servers]]")
	assert.True(t, strings.Contains(updated, `id = "github"`))
}
