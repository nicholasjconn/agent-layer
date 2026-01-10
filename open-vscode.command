#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd -L)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd -L)"
CODEX_HOME="$ROOT/.codex"

if [[ ! -d "$CODEX_HOME" ]]; then
  echo "error: CODEX_HOME directory not found at $CODEX_HOME" >&2
  echo "Run: node \"$ROOT/.agent-layer/sync/sync.mjs\"" >&2
  exit 1
fi

if ! command -v code >/dev/null 2>&1; then
  echo "error: VS Code 'code' CLI not found in PATH." >&2
  echo "In VS Code, run: Shell Command: Install 'code' command in PATH" >&2
  exit 1
fi

export CODEX_HOME
code_status=0
code "$ROOT" || code_status=$?
if [[ $code_status -ne 0 ]]; then
  echo "error: failed to launch VS Code (exit $code_status)" >&2
  exit "$code_status"
fi

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
