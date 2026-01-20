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

func TestLoadDefaultMCPServers(t *testing.T) {
	defaults, err := loadDefaultMCPServers()
	require.NoError(t, err)
	assert.NotEmpty(t, defaults)

	// Check for expected defaults
	ids := make(map[string]bool)
	for _, s := range defaults {
		ids[s.ID] = true
	}
	assert.True(t, ids["github"])
}

func TestAppendMissingDefaultMCPServers_Error(t *testing.T) {
	content := "[mcp]\n"
	tree, err := toml.LoadBytes([]byte(content))
	require.NoError(t, err)

	err = appendMissingDefaultMCPServers(tree, []string{"non-existent"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing default MCP server template")
}

func TestAppendMissingDefaultMCPServers_Empty(t *testing.T) {
	content := "[mcp]\n"
	tree, err := toml.LoadBytes([]byte(content))
	require.NoError(t, err)

	// Empty missing list should return nil immediately
	err = appendMissingDefaultMCPServers(tree, []string{})
	assert.NoError(t, err)
}

func TestMissingDefaultMCPServers_EmptyID(t *testing.T) {
	defaults := []DefaultMCPServer{
		{ID: "github"},
	}
	// Server with empty ID should be skipped
	servers := []config.MCPServer{
		{ID: ""},
		{ID: "github"},
	}

	missing := missingDefaultMCPServers(defaults, servers)
	assert.Empty(t, missing)
}

func TestMcpServerTrees_UnexpectedType(t *testing.T) {
	// Create a tree where mcp.servers is not []*toml.Tree
	content := `[mcp]
servers = "not-an-array"
`
	tree, err := toml.LoadBytes([]byte(content))
	require.NoError(t, err)

	_, err = mcpServerTrees(tree)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected type")
}

func TestMcpServerTrees_Nil(t *testing.T) {
	content := `[mcp]
`
	tree, err := toml.LoadBytes([]byte(content))
	require.NoError(t, err)

	servers, err := mcpServerTrees(tree)
	assert.NoError(t, err)
	assert.Nil(t, servers)
}

func TestAppendMissingDefaultMCPServers_McpServerTreesError(t *testing.T) {
	// Create a tree where mcp.servers has unexpected type
	content := `[mcp]
servers = "bad"
`
	tree, err := toml.LoadBytes([]byte(content))
	require.NoError(t, err)

	err = appendMissingDefaultMCPServers(tree, []string{"github"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected type")
}
