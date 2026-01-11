#!/usr/bin/env bash
set -euo pipefail

# ./al - repo-local launcher
#
# Purpose:
#   Select a single default behavior for the Agent Layer runner.
#
# Instructions:
#   - Uncomment exactly one option in the "Options" section below.
#   - Keep all other options commented.
#   - Do not edit the "path glue" section below.
#   - Do not leave multiple options uncommented.

# +------------------------------------------------------------+
# | al path glue (do not edit)                                 |
# | Keeps ./al working as a symlink or as .agent-layer/al.      |
# | Resolves RUNNER to the correct .agent-layer/run.sh.         |
# +------------------------------------------------------------+
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Probe for the runner in plausible locations:
# 1. Alongside this script (when invoked as .agent-layer/al)
# 2. In .agent-layer relative to invocation dir (when invoked as ./al wrapper)
if [[ -f "$SCRIPT_DIR/run.sh" ]]; then
  RUNNER="$SCRIPT_DIR/run.sh"
elif [[ -f "$SCRIPT_DIR/.agent-layer/run.sh" ]]; then
  RUNNER="$SCRIPT_DIR/.agent-layer/run.sh"
else
  echo "ERROR: Cannot locate run.sh (expected in $SCRIPT_DIR or $SCRIPT_DIR/.agent-layer)" >&2
  exit 2
fi

RUNNER_DIR="$(cd "$(dirname "$RUNNER")" && pwd)"

# Optional update check. Comment this out to skip update checks.
"$RUNNER_DIR/check-updates.sh" || true

# Options (choose one). Keep exactly one "exec" line uncommented.
#
# Option A (default): sync every run, load only .agent-layer/.env, then exec.
exec "$RUNNER" "$@"

# Option B: load .agent-layer/.env only (no sync).
# exec "$RUNNER" --env-only "$@"

# Option C: sync only (do not load any .env files).
# exec "$RUNNER" --sync-only "$@"

# Option D: check sync state and regenerate if stale, then env-only.
# exec "$RUNNER" --check-env "$@"

# Option E: sync every run, load .agent-layer/.env and .env, then exec.
# exec "$RUNNER" --project-env "$@"
