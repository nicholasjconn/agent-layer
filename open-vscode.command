#!/usr/bin/env bash
set -euo pipefail

# Launch VS Code with CODEX_HOME set to the repo-local .codex directory.
# Requires the VS Code `code` CLI on PATH.
# Set OPEN_VSCODE_NO_CLOSE=1 to keep the Terminal window open on macOS.

# Resolve the repo root and expected CODEX_HOME path.
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd -L)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd -L)"
CODEX_HOME="$ROOT/.codex"

# Fail fast if the repo-local CODEX_HOME has not been generated yet.
if [[ ! -d "$CODEX_HOME" ]]; then
  echo "error: CODEX_HOME directory not found at $CODEX_HOME" >&2
  echo "Run: node \"$ROOT/.agent-layer/src/sync/sync.mjs\"" >&2
  exit 1
fi

# Require the VS Code CLI so we can launch the workspace.
if ! command -v code >/dev/null 2>&1; then
  echo "error: VS Code 'code' CLI not found in PATH." >&2
  echo "In VS Code, run: Shell Command: Install 'code' command in PATH" >&2
  exit 1
fi

# Launch VS Code with the repo-local CODEX_HOME.
export CODEX_HOME
code_status=0
code "$ROOT" || code_status=$?
if [[ $code_status -ne 0 ]]; then
  echo "error: failed to launch VS Code (exit $code_status)" >&2
  exit "$code_status"
fi

# Optionally close the macOS Terminal window after launch.
if [[ -z "${OPEN_VSCODE_NO_CLOSE:-}" ]] && [[ "${TERM_PROGRAM:-}" == "Apple_Terminal" ]]; then
  if command -v osascript >/dev/null 2>&1; then
    osascript <<'EOF' >/dev/null 2>&1 || true
tell application "Terminal"
  if (count of windows) > 0 then
    try
      close front window
    end try
  end if
end tell
EOF
  fi
fi
