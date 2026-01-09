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
  [ "$status" -ne 0 ]
  run rg -n "run_shell_command\\(git status\\)" "$root/.gemini/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "run_shell_command\\(ls\\)" "$root/.gemini/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "\"some_tool\"" "$root/.gemini/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "\"extra\": true" "$root/.gemini/settings.json"
  [ "$status" -eq 0 ]

  run rg -n "Bash\\(bad:\\*\\)" "$root/.claude/settings.json"
  [ "$status" -ne 0 ]
  run rg -F "mcp__bad__*" "$root/.claude/settings.json"
  [ "$status" -ne 0 ]
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
  [ "$status" -ne 0 ]
  run rg -F 'git status(\\b.*)?$' "$root/.vscode/settings.json"
  [ "$status" -eq 0 ]
  run rg -F 'ls(\\b.*)?$' "$root/.vscode/settings.json"
  [ "$status" -eq 0 ]
  run rg -n "\"other\": 1" "$root/.vscode/settings.json"
  [ "$status" -eq 0 ]

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
