#!/usr/bin/env bats

# Tests for the agent-config CLI entrypoint.
load "helpers.bash"

@test "agent-config CLI resolves agent-layer root from env" {
  local root script
  root="$(make_tmp_dir)"
  mkdir -p "$root/config"
  write_agent_config "$root/config/agents.json" false false false false

  script="$AGENT_LAYER_ROOT/src/lib/agent-config.mjs"

  run env AGENT_LAYER_ROOT="$root" node "$script" --print-shell codex
  [ "$status" -eq 0 ]
  [[ "$output" == *"enabled=false"* ]]

  rm -rf "$root"
}
