#!/usr/bin/env bash
set -euo pipefail

# Hermetic CI test runner.
# Creates a temporary consumer workspace outside the repo, mounts .agent-layer,
# and runs the test suite in that isolated environment.

say() { printf "%s\n" "$*"; }
die() {
  printf "ERROR: %s\n" "$*" >&2
  exit 1
}

# Resolve the repo root (physical path to avoid symlink confusion)
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"

# Create workspace in system temp
WORKSPACE="$(mktemp -d "${TMPDIR:-/tmp}/agent-layer-ci.XXXXXX")"
cleanup() {
  if [[ -d "$WORKSPACE" ]]; then
    rm -rf "$WORKSPACE"
  fi
}
trap cleanup EXIT

say "==> Creating hermetic test workspace at $WORKSPACE"

# Create consumer roots (keep 2 to support "PWD in other repo" test scenarios)
ROOT_A="$WORKSPACE/root-a"
ROOT_B="$WORKSPACE/root-b"
mkdir -p "$ROOT_A" "$ROOT_B"

# Mount .agent-layer as a single symlink layer from temp -> checkout
# This creates a deterministic topology: /tmp/workspace/root-x/.agent-layer -> /repo
ln -s "$REPO_ROOT" "$ROOT_A/.agent-layer"
ln -s "$REPO_ROOT" "$ROOT_B/.agent-layer"

# Create the wrapper launcher (most realistic, matches what installer generates)
cat > "$ROOT_A/al" << 'EOF'
#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
exec "$SCRIPT_DIR/.agent-layer/al" "$@"
EOF
chmod +x "$ROOT_A/al"

say "==> Running sync in consumer root"
(cd "$ROOT_A" && node .agent-layer/src/sync/sync.mjs)

say "==> Running test suite"
"$REPO_ROOT/tests/run.sh" --work-root "$ROOT_A"

say "==> Tests completed successfully"
