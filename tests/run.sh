#!/usr/bin/env bash
set -euo pipefail

# Run formatting checks and the Bats test suite.

say() { printf "%s\n" "$*"; }
die() {
  printf "ERROR: %s\n" "$*" >&2
  exit 1
}

# Parse --work-root so the runner can be invoked from the agent-layer repo.
usage() {
  cat << 'EOF'
Usage: tests/run.sh [--work-root <path>]

Run formatting checks and the Bats suite. When running from inside the
agent-layer repo itself, pass --work-root to a consumer root that contains
a .agent-layer/ directory.
EOF
}

work_root=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --work-root)
      shift
      if [[ $# -eq 0 || -z "${1:-}" ]]; then
        die "--work-root requires a path."
      fi
      work_root="$1"
      shift
      ;;
    --work-root=*)
      work_root="${1#*=}"
      if [[ -z "$work_root" ]]; then
        die "--work-root requires a path."
      fi
      shift
      ;;
    -h | --help)
      usage
      exit 0
      ;;
    *)
      die "Unknown option: $1 (run --help for usage)."
      ;;
  esac
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PATHS_SH="$REPO_ROOT/src/lib/paths.sh"
if [[ ! -f "$PATHS_SH" ]]; then
  die "Missing src/lib/paths.sh (expected in the agent-layer repo)."
fi
# shellcheck disable=SC1090
source "$PATHS_SH"
if [[ -n "$work_root" ]]; then
  if [[ ! -d "$work_root" ]]; then
    die "--work-root does not exist: $work_root"
  fi
  work_root="$(cd "$work_root" && pwd)"
  if [[ ! -d "$work_root/.agent-layer" ]]; then
    die "--work-root must contain a .agent-layer directory: $work_root"
  fi
  cd "$work_root"
  WORKING_ROOT="$work_root"
  AGENTLAYER_ROOT="$work_root/.agent-layer"
  export WORKING_ROOT AGENTLAYER_ROOT
else
  if ! resolve_working_root "$SCRIPT_DIR" "$PWD" > /dev/null; then
    die "Missing .agent-layer/ directory in this path or any parent. Re-run with --work-root <path>."
  fi
  # Resolve entrypoint helpers so the runner works from any directory.
  ENTRYPOINT_SH="$SCRIPT_DIR/.agent-layer/src/lib/entrypoint.sh"
  if [[ ! -f "$ENTRYPOINT_SH" ]]; then
    ENTRYPOINT_SH="$SCRIPT_DIR/src/lib/entrypoint.sh"
  fi
  if [[ ! -f "$ENTRYPOINT_SH" ]]; then
    ENTRYPOINT_SH="$SCRIPT_DIR/../src/lib/entrypoint.sh"
  fi
  if [[ ! -f "$ENTRYPOINT_SH" ]]; then
    die "Missing src/lib/entrypoint.sh (expected near .agent-layer/)."
  fi
  # shellcheck disable=SC1090
  source "$ENTRYPOINT_SH"
  resolve_entrypoint_root || exit $?
fi

# Require external tools used by tests.
require_cmd() {
  local cmd="$1" hint="$2"
  if ! command -v "$cmd" > /dev/null 2>&1; then
    die "$cmd not found. $hint"
  fi
}

require_cmd git "Install git (dev-only)."
require_cmd node "Install Node.js (dev-only)."
require_cmd rg "Install ripgrep (macOS: brew install ripgrep; Ubuntu: apt-get install ripgrep)."

# Resolve the Bats binary (allow override via BATS_BIN).
BATS_BIN="${BATS_BIN:-bats}"
if ! command -v "$BATS_BIN" > /dev/null 2>&1; then
  die "bats not found. Install bats-core (macOS: brew install bats-core; Ubuntu: apt-get install bats)."
fi

# Run the Bats test suite.
say "==> Tests (bats)"
"$BATS_BIN" "$AGENTLAYER_ROOT/tests"
