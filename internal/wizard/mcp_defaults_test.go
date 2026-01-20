package wizard

import (
	"testing"

	toml "github.com/pelletier/go-toml"
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

	tree, err := toml.LoadBytes([]byte(content))
	require.NoError(t, err)

	err = appendMissingDefaultMCPServers(tree, missing)
	require.NoError(t, err)

	updated, err := tree.ToTomlString()
	require.NoError(t, err)

	assert.Contains(t, updated, "[[mcp.servers]]")
	assert.Contains(t, updated, `id = "github"`)
}
