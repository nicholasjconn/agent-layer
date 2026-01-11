#!/usr/bin/env bash
set -euo pipefail

# .agent-layer/check-updates.sh
# Warn when .agent-layer is pinned to an older tag (local tags only).

# Resolve the repo root relative to this script.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Exit quietly when git is unavailable or .agent-layer is not a repo.
if ! command -v git > /dev/null 2>&1; then
  exit 0
fi
if ! git -C "$ROOT/.agent-layer" rev-parse --is-inside-work-tree > /dev/null 2>&1; then
  exit 0
fi

# Skip checks if the current commit is not tagged.
current_tag="$(git -C "$ROOT/.agent-layer" describe --tags --exact-match 2> /dev/null || true)"
if [[ -z "$current_tag" ]]; then
  exit 0
fi

# Compare the current tag against the latest local tag.
latest_tag="$(git -C "$ROOT/.agent-layer" tag --list --sort=-v:refname | head -n 1)"
if [[ -z "$latest_tag" || "$current_tag" == "$latest_tag" ]]; then
  exit 0
fi

# Emit an upgrade reminder when a newer tag exists locally.
echo "warning: .agent-layer is on tag $current_tag; latest local tag is $latest_tag." >&2
echo "warning: Run '.agent-layer/agent-layer-install.sh --upgrade' from the working repo root to upgrade." >&2
