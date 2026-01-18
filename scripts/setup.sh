#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root_dir"

tool_bin="${AL_TOOL_BIN:-$root_dir/.tools/bin}"
gocache="${GOCACHE:-$root_dir/.cache/go-build}"
gomodcache="${GOMODCACHE:-$root_dir/.cache/go-mod}"

mkdir -p "$tool_bin" "$gocache" "$gomodcache"

if ! command -v go >/dev/null 2>&1; then
  echo "Go is required but was not found in PATH. Install Go 1.25.6+ and retry." >&2
  exit 1
fi
if ! command -v git >/dev/null 2>&1; then
  echo "git is required but was not found in PATH. Install Git and retry." >&2
  exit 1
fi
if ! command -v make >/dev/null 2>&1; then
  echo "make is required but was not found in PATH. Install Xcode Command Line Tools and retry." >&2
  exit 1
fi

echo "Downloading Go module dependencies..."
GOCACHE="$gocache" GOMODCACHE="$gomodcache" go mod download

echo "Installing pinned Go tools into $tool_bin..."
make tools TOOL_BIN="$tool_bin" GO_CACHE="$gocache" GO_MOD_CACHE="$gomodcache"

if command -v pre-commit >/dev/null 2>&1; then
  if hooks_path="$(git config --get core.hooksPath)"; then
    if [[ -n "$hooks_path" ]]; then
      cat >&2 <<EOF
pre-commit refuses to install when git config core.hooksPath is set (currently: $hooks_path).

Find where it is set:
  git config --show-origin --get core.hooksPath

Fix (recommended):
  - Repo-local:   git config --unset-all core.hooksPath
  - Global:      git config --global --unset-all core.hooksPath

Then rerun:
  ./scripts/setup.sh
EOF
      exit 1
    fi
  fi

  echo "Installing pre-commit hook..."
  pre-commit install --install-hooks
else
  echo "pre-commit not found; skipping git hook installation." >&2
  echo "Install pre-commit to enable local hooks, then run: pre-commit install --install-hooks" >&2
fi
