#!/usr/bin/env bash
set -euo pipefail

# Repo-local helper to load agentlayer env into the environment, then exec the command.
# Usage:
#   ./.agentlayer/with-env.sh claude
#   ./.agentlayer/with-env.sh codex
#   ./.agentlayer/with-env.sh --project-env gemini

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# Keep the caller's working directory; use ROOT only for env file paths.

INCLUDE_PROJECT_ENV=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --project-env) INCLUDE_PROJECT_ENV=1; shift ;;
    --help|-h)
      echo "Usage: ./.agentlayer/with-env.sh [--project-env] <command> [args...]"
      exit 0
      ;;
    --) shift; break ;;
    *) break ;;
  esac
done

if [[ $# -lt 1 ]]; then
  echo "Usage: ./.agentlayer/with-env.sh [--project-env] <command> [args...]"
  exit 2
fi

AGENT_ENV="$ROOT/.agentlayer/.env"

if [[ -f "$AGENT_ENV" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "$AGENT_ENV"
  set +a
fi

PROJECT_ENV="$ROOT/.env"
if [[ "$INCLUDE_PROJECT_ENV" -eq 1 && -f "$PROJECT_ENV" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "$PROJECT_ENV"
  set +a
fi

exec "$@"
