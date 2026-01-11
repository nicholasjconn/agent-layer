#!/usr/bin/env bash
set -euo pipefail

# Format and lint checker (shell + JS).
# Run from the agent-layer repo root.

say() { printf "%s\n" "$*"; }
die() {
  printf "ERROR: %s\n" "$*" >&2
  exit 1
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Require external tools.
require_cmd() {
  local cmd="$1" hint="$2"
  if ! command -v "$cmd" > /dev/null 2>&1; then
    die "$cmd not found. $hint"
  fi
}

require_cmd shfmt "Install shfmt (macOS: brew install shfmt; Ubuntu: apt-get install shfmt)."
require_cmd shellcheck "Install shellcheck (macOS: brew install shellcheck; Ubuntu: apt-get install shellcheck)."

# Resolve Prettier (local install preferred).
PRETTIER_BIN="$REPO_ROOT/node_modules/.bin/prettier"
if [[ -x "$PRETTIER_BIN" ]]; then
  PRETTIER="$PRETTIER_BIN"
elif command -v prettier > /dev/null 2>&1; then
  PRETTIER="$(command -v prettier)"
else
  die "prettier not found. Run: (cd .agent-layer && npm install) or install globally."
fi

# Collect shell sources for formatting and linting.
say "==> Shell format check (shfmt)"
shell_files=()
while IFS= read -r -d '' file; do
  shell_files+=("$file")
done < <(
  find "$REPO_ROOT" \
    -path "$REPO_ROOT/node_modules" -prune -o \
    -path "$REPO_ROOT/.git" -prune -o \
    -path "$REPO_ROOT/tmp" -prune -o \
    -type f \( -name "*.sh" -o -path "$REPO_ROOT/al" -o -path "$REPO_ROOT/.githooks/pre-commit" \) \
    -print0
)
if [[ "${#shell_files[@]}" -gt 0 ]]; then
  shfmt -d -i 2 -ci -sr "${shell_files[@]}"
fi

# Run shellcheck against the same shell sources.
say "==> Shell lint (shellcheck)"
if [[ "${#shell_files[@]}" -gt 0 ]]; then
  shellcheck "${shell_files[@]}"
fi

# Collect JS sources for formatting checks.
say "==> JS format check (prettier)"
js_files=()
while IFS= read -r -d '' file; do
  js_files+=("$file")
done < <(
  find "$REPO_ROOT" \
    -path "$REPO_ROOT/node_modules" -prune -o \
    -path "$REPO_ROOT/.git" -prune -o \
    -path "$REPO_ROOT/tmp" -prune -o \
    -type f \( -name "*.mjs" -o -name "*.js" \) \
    -print0
)
if [[ "${#js_files[@]}" -gt 0 ]]; then
  "$PRETTIER" --check "${js_files[@]}"
fi

say "==> Format check completed successfully"
