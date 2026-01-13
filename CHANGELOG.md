# Changelog
All notable changes to this project will be documented in this file.

## v0.1.0 - 2026-01-12
Initial release.

### Added
- Installer for per-project setup that pins `.agent-layer/` to tagged releases, with upgrade, version, and dev-branch options.
- Repo-local `./al` launcher with sync and environment modes plus local update checks.
- Sync pipeline that generates client configs from `.agent-layer/config` sources.
- MCP prompt server that exposes workflows as prompts.
- Project memory templates and setup/bootstrap helpers.
