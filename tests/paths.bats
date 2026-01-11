#!/usr/bin/env bats

# Tests for working root discovery helpers.
# Load shared helpers for temp roots and stub binaries.
load "helpers.bash"

# Test: resolve_working_root finds .agent-layer from working root
@test "resolve_working_root finds .agent-layer from working root" {
  local root
  root="$(create_working_root)"

  ROOT="$root" AGENTLAYER_ROOT="$AGENTLAYER_ROOT" run bash -c \
    'source "$AGENTLAYER_ROOT/src/lib/paths.sh"; resolve_working_root "$ROOT"'

  [ "$status" -eq 0 ]
  [ "$output" = "$root" ]
  rm -rf "$root"
}

# Test: resolve_working_root fails when .agent-layer is missing from ancestors
@test "resolve_working_root fails when .agent-layer is missing from ancestors" {
  local root
  root="/"

  ROOT="$root" AGENTLAYER_ROOT="$AGENTLAYER_ROOT" run bash -c \
    'source "$AGENTLAYER_ROOT/src/lib/paths.sh"; resolve_working_root "$ROOT"'

  [ "$status" -ne 0 ]
}
