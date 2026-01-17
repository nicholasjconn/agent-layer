# Changelog
All notable changes to this project will be documented in this file.

## v0.2.0 - 2026-01-17

Major architectural overhaul moving core logic from shell to Node.js.

### Added
- Per-agent opt-in configuration via `config/agents.json` with interactive setup prompt.
- HTTP transport support for MCP servers.
- Tavily MCP server for web search capabilities.
- `./al --version` flag with dirty suffix for non-tagged commits.
- User config preservation and backup during upgrades.

### Changed
- **Breaking:** CLI entrypoint is now `.agent-layer/agent-layer`; `./al` remains as the launcher wrapper in the parent root.
- Root resolution, environment loading, and cleanup moved from shell to Node.js (`src/lib/roots.mjs`, `src/lib/env.mjs`, `src/lib/cleanup.mjs`).
- Test framework migrated from Bats (shell) to Node.js native test runner.
- GitHub MCP server switched to hosted HTTP endpoint with PAT authentication.
- Architecture documentation updated to reflect new layer boundaries.

### Removed
- Shell scripts: `al`, `run.sh`, `setup.sh`, `clean.sh`, `check-updates.sh`, `open-vscode.command`.
- Shell-based root resolution: `src/lib/parent-root.sh`, `src/lib/temp-parent-root.sh`.

## v0.1.0 - 2026-01-12
Initial release.

### Added
- Installer for per-project setup that pins `.agent-layer/` to tagged releases, with upgrade, version, and dev-branch options.
- Repo-local `./al` launcher with sync and environment modes plus local update checks.
- Sync pipeline that generates client configs from `.agent-layer/config` sources.
- MCP prompt server that exposes workflows as prompts.
- Project memory templates and setup/bootstrap helpers.
