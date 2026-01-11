#!/usr/bin/env bash
set -euo pipefail

# Combined format check + tests for local development.
# Runs format-check.sh then ci.sh.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "==> Running format checks"
"$SCRIPT_DIR/format-check.sh"

echo ""
echo "==> Running tests"
"$REPO_ROOT/tests/ci.sh"

echo ""
echo "==> All checks passed!"
