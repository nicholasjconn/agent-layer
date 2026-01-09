#!/usr/bin/env bats

load "helpers.bash"

@test "prompt MCP server exposes tools/list handler" {
  local server_file
  server_file="$AGENTLAYER_ROOT/mcp/agent-layer-prompts/server.mjs"

  [ -f "$server_file" ]
  grep -q "ListToolsRequestSchema" "$server_file"
  grep -q "setRequestHandler(ListToolsRequestSchema" "$server_file"
  grep -Eq "capabilities:.*tools" "$server_file"
}

@test "prompt MCP server fails fast when workflows are missing" {
  local server_file
  server_file="$AGENTLAYER_ROOT/mcp/agent-layer-prompts/server.mjs"

  [ -f "$server_file" ]
  grep -q "could not find .agent-layer/workflows" "$server_file"
  grep -q "no workflow files found" "$server_file"
}
