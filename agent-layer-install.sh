#!/usr/bin/env bash
set -euo pipefail

# Agent Layer installer/upgrader.
# Run this from the parent repo root (the parent of .agent-layer/).

say() { printf "%s\n" "$*"; }
die() {
  printf "ERROR: %s\n" "$*" >&2
  exit 1
}

usage() {
  cat << 'EOF'
Usage: agent-layer-install.sh [--force] [--upgrade] [--version <tag>] [--latest-branch <branch>] [--repo-url <url>]

Installs/updates agent-layer in the current working repo and sets up a local launcher.
Defaults to the latest tagged release for new installs (detached HEAD).

Options:
  --force, -f       Overwrite ./al and allow user config to be replaced during upgrades
  --upgrade, -u     Upgrade .agent-layer to the latest tagged release (detached)
  --version <tag>   Install a specific tagged release (detached)
  --latest-branch   Update .agent-layer to the latest commit of a branch (detached; dev)
  --repo-url <url>  Override the agent-layer repo URL
  --help, -h        Show this help
EOF
}

# Default option values and repo URL configuration.
FORCE="0"
UPGRADE="0"
VERSION=""
VERSION_SET="0"
LATEST_BRANCH=""
LATEST_BRANCH_SET="0"
REPO_URL_DEFAULT="https://github.com/nicholasjconn/agent-layer.git"
REPO_URL_OVERRIDE="0"
REPO_URL="${AGENT_LAYER_REPO_URL:-$REPO_URL_DEFAULT}"
if [[ -n "${AGENT_LAYER_REPO_URL:-}" ]]; then
  REPO_URL_OVERRIDE="1"
fi

# Parse CLI flags and reject unknown options.
while [[ $# -gt 0 ]]; do
  case "$1" in
    --force | -f)
      FORCE="1"
      ;;
    --upgrade | -u)
      UPGRADE="1"
      ;;
    --version)
      [[ $# -ge 2 ]] || die "--version requires a value"
      VERSION="$2"
      VERSION_SET="1"
      shift
      ;;
    --version=*)
      VERSION="${1#*=}"
      VERSION_SET="1"
      ;;
    --latest-branch)
      [[ $# -ge 2 ]] || die "--latest-branch requires a value"
      LATEST_BRANCH="$2"
      LATEST_BRANCH_SET="1"
      shift
      ;;
    --latest-branch=*)
      LATEST_BRANCH="${1#*=}"
      LATEST_BRANCH_SET="1"
      ;;
    --repo-url)
      [[ $# -ge 2 ]] || die "--repo-url requires a value"
      REPO_URL="$2"
      REPO_URL_OVERRIDE="1"
      shift
      ;;
    --repo-url=*)
      REPO_URL="${1#*=}"
      REPO_URL_OVERRIDE="1"
      ;;
    --help | -h)
      usage
      exit 0
      ;;
    *)
      die "Unknown argument: $1"
      ;;
  esac
  shift
done

# Validate incompatible or missing flag combinations.
if [[ "$UPGRADE" == "1" && -n "$LATEST_BRANCH" ]]; then
  die "Choose only one: --upgrade or --latest-branch <branch>"
fi
if [[ "$UPGRADE" == "1" && -n "$VERSION" ]]; then
  die "Choose only one: --upgrade or --version <tag>"
fi
if [[ -n "$VERSION" && -n "$LATEST_BRANCH" ]]; then
  die "Choose only one: --version <tag> or --latest-branch <branch>"
fi
if [[ "$LATEST_BRANCH_SET" == "1" && -z "$LATEST_BRANCH" ]]; then
  die "--latest-branch requires a value"
fi
if [[ "$VERSION_SET" == "1" && -z "$VERSION" ]]; then
  die "--version requires a value"
fi

# Require git up front because install/upgrade uses it heavily.
command -v git > /dev/null 2>&1 || die "git is required (not found)."

# Confirm we are running from the repo root (not inside .agent-layer).
PARENT_ROOT="$(pwd -P)"
if [[ "$(basename "$PARENT_ROOT")" == ".agent-layer" ]]; then
  die "Run this from the parent repo root (parent of .agent-layer/), not inside .agent-layer/."
fi

# Enforce or warn about git repo state to support hooks.
if git rev-parse --show-toplevel > /dev/null 2>&1; then
  GIT_ROOT="$(git rev-parse --show-toplevel)"
  if [[ "$GIT_ROOT" != "$PARENT_ROOT" ]]; then
    die "Run this from the repo root: $GIT_ROOT"
  fi
else
  say "WARNING: This directory does not appear to be a git repo."
  say "If you meant a different folder, stop now."
  say "You can continue, but hooks won't be enabled until you init git."
  if [[ -t 0 ]]; then
    read -r -p "Continue anyway? [y/N] " reply
    case "$reply" in
      y | Y | yes | YES) ;;
      *)
        die "Aborted."
        ;;
    esac
  else
    die "Not a git repo and no TTY available to confirm. Re-run from a TTY or after init."
  fi
fi

AGENT_LAYER_DIR="$PARENT_ROOT/.agent-layer"
NEW_INSTALL="0"

# Resolve the fetch target for upgrades (explicit URL or origin remote).
resolve_fetch_target() {
  if [[ "$REPO_URL_OVERRIDE" == "1" ]]; then
    printf "%s" "$REPO_URL"
    return 0
  fi
  if git -C "$AGENT_LAYER_DIR" remote get-url origin > /dev/null 2>&1; then
    printf "%s" "origin"
    return 0
  fi
  die "No origin remote found. Use --repo-url <url> or set AGENT_LAYER_REPO_URL."
}

# Ensure a specific tag exists in the remote before cloning.
ensure_version_exists_remote() {
  local version="$1"
  local remote="$2"
  local refs

  say "==> Checking for tag '$version' in $remote"
  if ! refs="$(git ls-remote --tags "$remote" "refs/tags/$version" "refs/tags/$version^{}" 2> /dev/null)"; then
    die "Failed to query tags from $remote; cannot verify '$version'."
  fi
  if [[ -z "$refs" ]]; then
    die "Tag '$version' not found; cannot install requested version."
  fi
}

# User-managed config paths that should be preserved during upgrades.
USER_CONFIG_PATHS=(
  "config/agents.json"
  "config/mcp-servers.json"
  "config/policy/commands.json"
)
USER_CONFIG_DIRS=(
  "config/instructions"
  "config/workflows"
)
USER_CONFIG_BACKUP_DIR=""

is_user_config_path() {
  local path="$1"
  case "$path" in
    config/instructions/*) return 0 ;;
    config/workflows/*) return 0 ;;
  esac
  for entry in "${USER_CONFIG_PATHS[@]}"; do
    if [[ "$path" == "$entry" ]]; then
      return 0
    fi
  done
  return 1
}

assert_only_user_config_changes() {
  local dirty line path non_config
  dirty="$(git -C "$AGENT_LAYER_DIR" status --porcelain)"
  [[ -z "$dirty" ]] && return 0

  non_config=()
  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    path="${line:3}"
    if [[ "$path" == *" -> "* ]]; then
      path="${path##* -> }"
    fi
    if ! is_user_config_path "$path"; then
      non_config+=("$path")
    fi
  done <<< "$dirty"

  if [[ "${#non_config[@]}" -gt 0 ]]; then
    die ".agent-layer has uncommitted changes outside user config. Commit or stash before upgrading."
  fi
}

should_preserve_user_config() {
  [[ "$FORCE" != "1" && "$NEW_INSTALL" != "1" ]]
}

backup_user_config() {
  local entry dir backup_dir
  mkdir -p "$AGENT_LAYER_DIR/tmp"
  backup_dir="$(mktemp -d "$AGENT_LAYER_DIR/tmp/user-config-backup.XXXXXX")"
  USER_CONFIG_BACKUP_DIR="$backup_dir"

  for entry in "${USER_CONFIG_PATHS[@]}"; do
    if [[ -e "$AGENT_LAYER_DIR/$entry" ]]; then
      mkdir -p "$backup_dir/$(dirname "$entry")"
      cp -a "$AGENT_LAYER_DIR/$entry" "$backup_dir/$entry"
    fi
  done

  for dir in "${USER_CONFIG_DIRS[@]}"; do
    if [[ -d "$AGENT_LAYER_DIR/$dir" ]]; then
      mkdir -p "$backup_dir/$dir"
      cp -a "$AGENT_LAYER_DIR/$dir/." "$backup_dir/$dir/"
    fi
  done
}

reset_user_config_worktree() {
  local entry dir
  for entry in "${USER_CONFIG_PATHS[@]}"; do
    git -C "$AGENT_LAYER_DIR" checkout -q -- "$entry" > /dev/null 2>&1 || true
    git -C "$AGENT_LAYER_DIR" clean -fd -- "$entry" > /dev/null 2>&1 || true
  done

  for dir in "${USER_CONFIG_DIRS[@]}"; do
    git -C "$AGENT_LAYER_DIR" checkout -q -- "$dir" > /dev/null 2>&1 || true
    git -C "$AGENT_LAYER_DIR" clean -fd -- "$dir" > /dev/null 2>&1 || true
  done
}

restore_user_config() {
  local entry dir backup_dir
  backup_dir="$USER_CONFIG_BACKUP_DIR"
  [[ -z "$backup_dir" || ! -d "$backup_dir" ]] && return 0

  for entry in "${USER_CONFIG_PATHS[@]}"; do
    if [[ -e "$backup_dir/$entry" ]]; then
      mkdir -p "$AGENT_LAYER_DIR/$(dirname "$entry")"
      cp -a "$backup_dir/$entry" "$AGENT_LAYER_DIR/$entry"
    fi
  done

  for dir in "${USER_CONFIG_DIRS[@]}"; do
    if [[ -d "$backup_dir/$dir" ]]; then
      mkdir -p "$AGENT_LAYER_DIR/$dir"
      cp -a "$backup_dir/$dir/." "$AGENT_LAYER_DIR/$dir/"
    fi
  done

  rm -rf "$backup_dir"
  USER_CONFIG_BACKUP_DIR=""
}

# Upgrade .agent-layer to the latest local tag.
upgrade_agent_layer() {
  local fetch_target latest_tag current_commit current_tag changes

  if should_preserve_user_config; then
    assert_only_user_config_changes
    backup_user_config
    reset_user_config_worktree
  elif [[ -n "$(git -C "$AGENT_LAYER_DIR" status --porcelain)" ]]; then
    die ".agent-layer has uncommitted changes. Commit or stash before upgrading."
  fi

  fetch_target="$(resolve_fetch_target)"

  say "==> Fetching tags for .agent-layer"
  git -C "$AGENT_LAYER_DIR" fetch --tags "$fetch_target"

  latest_tag="$(git -C "$AGENT_LAYER_DIR" tag --list --sort=-v:refname | head -n 1)"
  [[ -n "$latest_tag" ]] || die "No tags found after fetching; cannot install latest release. Use --latest-branch <branch> for dev builds."

  current_commit="$(git -C "$AGENT_LAYER_DIR" rev-parse --short HEAD)"
  current_tag="$(git -C "$AGENT_LAYER_DIR" describe --tags --exact-match 2> /dev/null || true)"

  say "==> Current version: ${current_tag:-$current_commit}"
  say "==> Latest tag: $latest_tag"

  if [[ "$current_tag" == "$latest_tag" ]]; then
    say "==> .agent-layer is already up to date."
    if should_preserve_user_config; then
      restore_user_config
    fi
    return 0
  fi

  say "==> Checking out $latest_tag"
  git -C "$AGENT_LAYER_DIR" checkout -q "$latest_tag"
  if should_preserve_user_config; then
    restore_user_config
  fi

  say "==> Changes since ${current_tag:-$current_commit}:"
  changes="$(git -C "$AGENT_LAYER_DIR" --no-pager log --oneline "$current_commit..$latest_tag" || true)"
  if [[ -n "$changes" ]]; then
    printf "%s\n" "$changes"
  else
    say "  (no commits listed)"
  fi

  say "==> Note: .agent-layer is now on a detached HEAD at $latest_tag."
}

# Update .agent-layer to a specific tag.
install_agent_layer_version() {
  local version="$1"
  local fetch_target current_commit current_tag changes

  if should_preserve_user_config; then
    assert_only_user_config_changes
    backup_user_config
    reset_user_config_worktree
  elif [[ -n "$(git -C "$AGENT_LAYER_DIR" status --porcelain)" ]]; then
    die ".agent-layer has uncommitted changes. Commit or stash before updating."
  fi

  fetch_target="$(resolve_fetch_target)"

  say "==> Fetching tags for .agent-layer"
  git -C "$AGENT_LAYER_DIR" fetch --tags "$fetch_target"

  if ! git -C "$AGENT_LAYER_DIR" rev-parse -q --verify "refs/tags/$version" > /dev/null; then
    die "Tag '$version' not found after fetching; cannot install requested version."
  fi

  current_commit="$(git -C "$AGENT_LAYER_DIR" rev-parse --short HEAD)"
  current_tag="$(git -C "$AGENT_LAYER_DIR" describe --tags --exact-match 2> /dev/null || true)"

  say "==> Current version: ${current_tag:-$current_commit}"
  say "==> Requested tag: $version"

  if [[ "$current_tag" == "$version" ]]; then
    say "==> .agent-layer is already at $version."
    if should_preserve_user_config; then
      restore_user_config
    fi
    return 0
  fi

  say "==> Checking out $version"
  git -C "$AGENT_LAYER_DIR" checkout -q "$version"
  if should_preserve_user_config; then
    restore_user_config
  fi

  say "==> Changes since ${current_tag:-$current_commit}:"
  changes="$(git -C "$AGENT_LAYER_DIR" --no-pager log --oneline "$current_commit..$version" || true)"
  if [[ -n "$changes" ]]; then
    printf "%s\n" "$changes"
  else
    say "  (no commits listed)"
  fi

  say "==> Note: .agent-layer is now on a detached HEAD at $version."
}

# Update .agent-layer to the latest commit on a specific branch.
latest_branch_agent_layer() {
  local branch="$1"
  local fetch_target current_commit latest_commit changes

  if should_preserve_user_config; then
    assert_only_user_config_changes
    backup_user_config
    reset_user_config_worktree
  elif [[ -n "$(git -C "$AGENT_LAYER_DIR" status --porcelain)" ]]; then
    die ".agent-layer has uncommitted changes. Commit or stash before updating."
  fi

  fetch_target="$(resolve_fetch_target)"

  say "==> Fetching latest commit for branch '$branch'"
  git -C "$AGENT_LAYER_DIR" fetch "$fetch_target" "$branch"

  latest_commit="$(git -C "$AGENT_LAYER_DIR" rev-parse --short FETCH_HEAD)"
  current_commit="$(git -C "$AGENT_LAYER_DIR" rev-parse --short HEAD)"

  if [[ "$current_commit" == "$latest_commit" ]]; then
    say "==> .agent-layer is already at latest $branch ($latest_commit)."
  else
    say "==> Current commit: $current_commit"
    say "==> Latest $branch commit: $latest_commit"
  fi
  say "==> Checking out latest $branch commit"
  git -C "$AGENT_LAYER_DIR" checkout -q --detach FETCH_HEAD
  if should_preserve_user_config; then
    restore_user_config
  fi

  if [[ "$current_commit" != "$latest_commit" ]]; then
    changes="$(git -C "$AGENT_LAYER_DIR" --no-pager log --oneline -n 5 FETCH_HEAD || true)"
    if [[ -n "$changes" ]]; then
      say "==> Recent commits:"
      printf "%s\n" "$changes"
    fi
  fi

  say "==> Note: .agent-layer is now on a detached HEAD at $latest_commit."
  say "==> To update again, re-run the installer with: --latest-branch $branch"
}

# Ensure .agent-layer exists, then apply the requested upgrade behavior.
if [[ ! -e "$AGENT_LAYER_DIR" ]]; then
  NEW_INSTALL="1"
  [[ -n "$REPO_URL" ]] || die "Missing repo URL (set AGENT_LAYER_REPO_URL or use --repo-url)."
  if [[ -n "$VERSION" ]]; then
    ensure_version_exists_remote "$VERSION" "$REPO_URL"
  fi
  say "==> Cloning agent-layer into .agent-layer/"
  git clone "$REPO_URL" "$AGENT_LAYER_DIR"
  if [[ -n "$VERSION" ]]; then
    install_agent_layer_version "$VERSION"
  elif [[ -n "$LATEST_BRANCH" ]]; then
    latest_branch_agent_layer "$LATEST_BRANCH"
  else
    upgrade_agent_layer
  fi
else
  if [[ -d "$AGENT_LAYER_DIR" ]]; then
    if git -C "$AGENT_LAYER_DIR" rev-parse --is-inside-work-tree > /dev/null 2>&1; then
      if [[ -n "$VERSION" ]]; then
        install_agent_layer_version "$VERSION"
      elif [[ "$UPGRADE" == "1" ]]; then
        upgrade_agent_layer
      elif [[ -n "$LATEST_BRANCH" ]]; then
        latest_branch_agent_layer "$LATEST_BRANCH"
      else
        say "==> .agent-layer exists and is a git repo; leaving as-is"
      fi
    else
      die ".agent-layer exists but is not a git repo. Move it aside or remove it, then re-run."
    fi
  else
    die ".agent-layer exists but is not a directory. Move it aside or remove it, then re-run."
  fi
fi

# Run setup to generate configs and install MCP prompt server dependencies.
if [[ -f "$AGENT_LAYER_DIR/agent-layer" ]]; then
  install_config_args=(--install-config --parent-root "$PARENT_ROOT" --agent-layer-root "$AGENT_LAYER_DIR")
  if [[ "$FORCE" == "1" ]]; then
    install_config_args+=(--force)
  fi
  if [[ "$NEW_INSTALL" == "1" ]]; then
    install_config_args+=(--new-install)
  fi
  "$AGENT_LAYER_DIR/agent-layer" "${install_config_args[@]}"
  say "==> Running setup"
  "$AGENT_LAYER_DIR/agent-layer" \
    --setup \
    --parent-root "$PARENT_ROOT" \
    --agent-layer-root "$AGENT_LAYER_DIR"
else
  die "Missing .agent-layer/agent-layer"
fi

# Print next steps for configuring and running the tools.
say ""
say "Next steps (required):"
say "  1) Copy .agent-layer/.env.example to .agent-layer/.env and fill in tokens."
say "  2) Review and edit config files:"
say "     - .agent-layer/config/agents.json"
say "     - .agent-layer/config/mcp-servers.json"
say "     - .agent-layer/config/policy/commands.json"
say "     - .agent-layer/config/instructions/*.md"
say "     - .agent-layer/config/workflows/*.md"
say ""
say "After config changes, re-run setup:"
say "  ./al --setup"
say ""
say "Launch an agent:"
say "  ./al gemini"
say "  ./al claude"
say "  ./al codex"

cleanup_installer() {
  local script_source script_dir script_path repo_root rel_path

  script_source="${BASH_SOURCE[0]}"
  if [[ -z "$script_source" ]]; then
    return 0
  fi
  if [[ "$script_source" == "-" || "$script_source" == "bash" || "$script_source" == */bash ]]; then
    return 0
  fi
  if [[ ! -f "$script_source" ]]; then
    return 0
  fi

  script_dir="$(cd "$(dirname "$script_source")" && pwd -P)"
  script_path="$script_dir/$(basename "$script_source")"

  if [[ "$script_path" == "$AGENT_LAYER_DIR/"* ]]; then
    say "==> Installer lives in .agent-layer; keeping it in place."
    return 0
  fi

  if repo_root="$(git -C "$script_dir" rev-parse --show-toplevel 2> /dev/null)"; then
    if [[ "$script_path" == "$repo_root/"* ]]; then
      rel_path="${script_path#"$repo_root/"}"
      if git -C "$repo_root" ls-files --error-unmatch "$rel_path" > /dev/null 2>&1; then
        say "==> Installer is tracked in git; keeping it in place."
        return 0
      fi
    fi
  fi

  if [[ -w "$script_path" ]]; then
    rm -f "$script_path"
    say "==> Removed downloaded installer: $script_path"
  else
    say "==> Installer still present at $script_path (remove it if you no longer need it)."
  fi
}

cleanup_installer
