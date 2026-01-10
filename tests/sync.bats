#!/usr/bin/env bats

load "helpers.bash"

@test "sync generates Codex config and instructions" {
  local root
  root="$(create_working_root)"

  run bash -c "cd \"$root\" && node .agent-layer/sync/sync.mjs"
  [ "$status" -eq 0 ]

  [ -f "$root/.codex/config.toml" ]
  [ -f "$root/.codex/AGENTS.md" ]
  grep -q '^\[mcp_servers\.' "$root/.codex/config.toml"
  grep -q 'GENERATED FILE' "$root/.codex/AGENTS.md"

  rm -rf "$root"
}

@test "sync overwrites command allowlists from policy" {
  local root
  root="$(create_working_root)"

  mkdir -p "$root/.gemini" "$root/.claude" "$root/.vscode"
  cat >"$root/.gemini/settings.json" <<'EOF'
{
  "tools": { "allowed": ["run_shell_command(bad)", "some_tool"], "extra": true }
}
EOF
  cat >"$root/.claude/settings.json" <<'EOF'
{
  "permissions": { "allow": ["Bash(bad:*)", "mcp__bad__*", "Edit"], "extra": true }
}
EOF
  cat >"$root/.vscode/settings.json" <<'EOF'
{
  "chat.tools.terminal.autoApprove": { "/^bad(\\b.*)?$/": true },
  "other": 1
}
EOF

  run bash -c "cd \"$root\" && node .agent-layer/sync/sync.mjs"
  [ "$status" -eq 0 ]

  run rg -n "run_shell_command\\(bad\\)" "$root/.gemini/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "run_shell_command\\(git status\\)" "$root/.gemini/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "run_shell_command\\(ls\\)" "$root/.gemini/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "\"some_tool\"" "$root/.gemini/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "\"extra\": true" "$root/.gemini/settings.json"
  [ "$status" -eq 0 ]

  run rg -n "Bash\\(bad:\\*\\)" "$root/.claude/settings.json"
  [ "$status" -eq 0 ]
  run rg -F "mcp__bad__*" "$root/.claude/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "Bash\\(git status:\\*\\)" "$root/.claude/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "Bash\\(ls:\\*\\)" "$root/.claude/settings.json"
  [ "$status" -eq 0 ]
  run rg -F "mcp__context7__*" "$root/.claude/settings.json"
  [ "$status" -eq 0 ]
  run rg -F "\"Edit\"" "$root/.claude/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "\"extra\": true" "$root/.claude/settings.json"
  [ "$status" -eq 0 ]

  run rg -F 'bad(\\b.*)?$' "$root/.vscode/settings.json"
  [ "$status" -eq 0 ]
  run rg -F 'git status(\\b.*)?$' "$root/.vscode/settings.json"
  [ "$status" -eq 0 ]
  run rg -F 'ls(\\b.*)?$' "$root/.vscode/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "\"other\": 1" "$root/.vscode/settings.json"
  [ "$status" -eq 0 ]

  rm -rf "$root"
}

@test "sync fails when policy contains unsafe argv token" {
  local root
  root="$(create_sync_working_root)"

  cat >"$root/.agent-layer/policy/commands.json" <<'EOF'
{
  "version": 1,
  "allowed": [
    { "argv": ["git", "status:all"] }
  ]
}
EOF

  run bash -c "cd \"$root\" && node .agent-layer/sync/sync.mjs"
  [ "$status" -ne 0 ]
  [[ "$output" == *"unsupported characters"* ]]

  rm -rf "$root"
}

@test "sync --overwrite removes divergent allowlists and mcp entries" {
  local root
  root="$(create_working_root)"

  mkdir -p "$root/.gemini" "$root/.claude" "$root/.vscode" "$root/.codex/rules"
  cat >"$root/.gemini/settings.json" <<'EOF'
{
  "tools": { "allowed": ["run_shell_command(bad)", "some_tool"], "extra": true },
  "mcpServers": { "extra": { "command": "node", "args": [] } }
}
EOF
  cat >"$root/.claude/settings.json" <<'EOF'
{
  "permissions": { "allow": ["Bash(bad:*)", "mcp__bad__*", "Edit"], "extra": true }
}
EOF
  cat >"$root/.vscode/settings.json" <<'EOF'
{
  "chat.tools.terminal.autoApprove": { "/^bad(\\b.*)?$/": true },
  "other": 1
}
EOF
  cat >"$root/.mcp.json" <<'EOF'
{
  "mcpServers": { "extra": { "command": "node", "args": [] } }
}
EOF
  cat >"$root/.vscode/mcp.json" <<'EOF'
{
  "servers": { "extra": { "type": "stdio", "command": "node", "args": [] } }
}
EOF
  mkdir -p "$root/.codex"
  cat >"$root/.codex/config.toml" <<'EOF'
# GENERATED FILE - DO NOT EDIT DIRECTLY
[mcp_servers.extra]
command = "node"
EOF
  cat >"$root/.codex/rules/agent-layer.rules" <<'EOF'
prefix_rule(pattern=["bad"], decision="allow", justification="legacy")
EOF

  run bash -c "cd \"$root\" && node .agent-layer/sync/sync.mjs --overwrite"
  [ "$status" -eq 0 ]

  run rg -n "run_shell_command\\(bad\\)" "$root/.gemini/settings.json"
  [ "$status" -ne 0 ]
  run rg -n "Bash\\(bad:\\*\\)" "$root/.claude/settings.json"
  [ "$status" -ne 0 ]
  run rg -F 'bad(\\b.*)?$' "$root/.vscode/settings.json"
  [ "$status" -ne 0 ]

  run rg -n "\"extra\": \\{" "$root/.gemini/settings.json"
  [ "$status" -ne 0 ]
  run rg -n "\"extra\": \\{" "$root/.mcp.json"
  [ "$status" -ne 0 ]
  run rg -n "\"extra\": \\{" "$root/.vscode/mcp.json"
  [ "$status" -ne 0 ]
  run rg -n "mcp_servers\\.extra" "$root/.codex/config.toml"
  [ "$status" -ne 0 ]
  run rg -n "\\[\"bad\"\\]" "$root/.codex/rules/agent-layer.rules"
  [ "$status" -ne 0 ]

  rm -rf "$root"
}

@test "sync fails when instructions directory is missing" {
  local root
  root="$(create_sync_working_root)"

  rm -rf "$root/.agent-layer/instructions"

  run bash -c "cd \"$root\" && node .agent-layer/sync/sync.mjs"
  [ "$status" -ne 0 ]
  [[ "$output" == *"missing instructions directory"* ]]

  rm -rf "$root"
}

@test "sync fails when instructions directory has no markdown files" {
  local root
  root="$(create_sync_working_root)"

  rm -f "$root/.agent-layer/instructions/"*.md

  run bash -c "cd \"$root\" && node .agent-layer/sync/sync.mjs"
  [ "$status" -ne 0 ]
  [[ "$output" == *"no instruction files found"* ]]

  rm -rf "$root"
}

@test "sync fails when workflows directory is missing" {
  local root
  root="$(create_sync_working_root)"

  rm -rf "$root/.agent-layer/workflows"

  run bash -c "cd \"$root\" && node .agent-layer/sync/sync.mjs"
  [ "$status" -ne 0 ]
  [[ "$output" == *"missing workflows directory"* ]]

  rm -rf "$root"
}

@test "sync fails when workflows directory has no markdown files" {
  local root
  root="$(create_sync_working_root)"

  rm -f "$root/.agent-layer/workflows/"*.md

  run bash -c "cd \"$root\" && node .agent-layer/sync/sync.mjs"
  [ "$status" -ne 0 ]
  [[ "$output" == *"no workflow files found"* ]]

  rm -rf "$root"
}

@test "sync fails when MCP server catalog is missing" {
  local root
  root="$(create_sync_working_root)"

  rm -f "$root/.agent-layer/mcp/servers.json"

  run bash -c "cd \"$root\" && node .agent-layer/sync/sync.mjs"
  [ "$status" -ne 0 ]
  [[ "$output" == *"servers.json not found"* ]]

  rm -rf "$root"
}

@test "inspect scans only working-root Codex sessions" {
  local root external
  root="$(create_working_root)"

  mkdir -p "$root/.codex/sessions/2025/01/01"
  cat >"$root/.codex/sessions/2025/01/01/rollout-local.jsonl" <<'EOF'
{"msg":{"type":"exec_approval_request","command":["echo","hi"],"cwd":"/tmp"}}
EOF

  external="$(make_tmp_dir)"
  mkdir -p "$external/sessions/2025/01/01"
  cat >"$external/sessions/2025/01/01/rollout-external.jsonl" <<'EOF'
{"msg":{"type":"exec_approval_request","command":["echo","external"],"cwd":"/tmp"}}
EOF

  run bash -c "cd \"$root\" && CODEX_HOME=\"$external\" node .agent-layer/sync/inspect.mjs > \"$root/out.json\""
  [ "$status" -eq 0 ]

  run node -e "const data=require(process.argv[1]); if (data.summary.approvals !== 1) process.exit(1); if (data.divergences.approvals[0].prefix !== 'echo hi') process.exit(1);" "$root/out.json"
  [ "$status" -eq 0 ]
  run node -e "const data=require(process.argv[1]); if (data.divergences.approvals.some((a)=>a.prefix==='echo external')) process.exit(1);" "$root/out.json"
  [ "$status" -eq 0 ]

  rm -rf "$root" "$external"
}

@test "inspect handles Codex config with empty args" {
  local root
  root="$(create_sync_working_root)"

  cat >"$root/.agent-layer/mcp/servers.json" <<'EOF'
{
  "version": 1,
  "servers": [
    {
      "name": "empty-args",
      "enabled": true,
      "transport": "stdio",
      "command": "node",
      "args": [],
      "envVars": []
    }
  ]
}
EOF

  mkdir -p "$root/.codex"
  cat >"$root/.codex/config.toml" <<'EOF'
# GENERATED FILE - DO NOT EDIT DIRECTLY
[mcp_servers.empty-args]
command = "node"
EOF

  run bash -c "cd \"$root\" && node .agent-layer/sync/inspect.mjs > \"$root/out.json\""
  [ "$status" -eq 0 ]

  run node -e "const data=require(process.argv[1]); if (data.summary.mcp !== 0) process.exit(1);" "$root/out.json"
  [ "$status" -eq 0 ]

  rm -rf "$root"
}

@test "sync warning counts repo-local Codex approvals" {
  local root external
  root="$(create_working_root)"

  mkdir -p "$root/.codex/sessions/2025/01/01"
  cat >"$root/.codex/sessions/2025/01/01/rollout-local.jsonl" <<'EOF'
{"type":"exec_approval_request","command":["rm","-f","README.md"]}
{"msg":{"type":"exec_approval_request","command":["whoami"]}}
EOF

  external="$(make_tmp_dir)"
  mkdir -p "$external/sessions/2025/01/01"
  cat >"$external/sessions/2025/01/01/rollout-external.jsonl" <<'EOF'
{"msg":{"type":"exec_approval_request","command":["echo","external"]}}
EOF

  run bash -c "cd \"$root\" && CODEX_HOME=\"$external\" node .agent-layer/sync/sync.mjs"
  [ "$status" -eq 0 ]
  [[ "$output" == *"approvals: 2"* ]]

  rm -rf "$root" "$external"
}
