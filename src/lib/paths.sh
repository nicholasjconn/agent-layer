#!/usr/bin/env bash

# Working root discovery for agent-layer scripts.
# Finds the nearest ancestor directory that contains .agent-layer/.

find_working_root() {
  # Search up to 50 ancestors for .agent-layer and return the first match.
  local dir
  dir="$(cd "$1" && pwd)"
  for _ in {1..50}; do
    if [[ -d "$dir/.agent-layer" ]]; then
      printf "%s" "$dir"
      return 0
    fi
    local parent
    parent="$(cd "$dir/.." && pwd)"
    if [[ "$parent" == "$dir" ]]; then
      break
    fi
    dir="$parent"
  done
  return 1
}

resolve_working_root() {
  # Resolve the working root from a list of start directories.
  local start
  for start in "$@"; do
    if [[ -z "$start" ]]; then
      continue
    fi
    local root
    root="$(find_working_root "$start" || true)"
    if [[ -n "$root" ]]; then
      printf "%s" "$root"
      return 0
    fi
  done
  return 1
}
