#!/usr/bin/env bats

# Tests for git hook behavior.
# Load shared helpers for temp roots and stub binaries.
load "helpers.bash"

# Test: pre-commit hook runs the test suite
@test "pre-commit hook runs the test suite" {
  local root bash_bin
  root="$(create_isolated_working_root)"
  bash_bin="$(command -v bash)"

  cat >"$root/.agent-layer/tests/run.sh" <<'EOF'
#!/usr/bin/env bash
echo "ran"
EOF
  chmod +x "$root/.agent-layer/tests/run.sh"

  run "$bash_bin" -c "cd '$root' && '$root/.agent-layer/.githooks/pre-commit'"
  [ "$status" -eq 0 ]
  [ "$output" = "ran" ]

  rm -rf "$root"
}
