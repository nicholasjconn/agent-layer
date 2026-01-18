# Roadmap (vNext / Go rewrite)

Purpose: Phased plan of work; guides architecture decisions and sequencing.

Update rules:
- The roadmap is a single list of numbered phases.
- Do not renumber completed phases (phases marked with ✅).
- You may renumber incomplete phases when updating the roadmap (for example, to insert a new phase).
- Incomplete phases include Goal, Tasks (checkbox list), and Exit criteria.

Phase template (completed):
## Phase N ✅ — <phase name>
- <Accomplishment summary bullet>
- <Accomplishment summary bullet>

Phase template (incomplete):
## Phase N — <phase name>

### Goal
- <What success looks like for this phase, in 1 to 3 bullet points.>

### Tasks
- [ ] <Concrete deliverable-oriented task>
- [ ] <Concrete deliverable-oriented task>

### Exit criteria
- <Objective condition that must be true to call the phase complete.>

---

## Phases

<!-- PHASES START -->

## Phase 1 ✅ — Define the vNext contract (docs-first)
- Defined the vNext product contract: repo-local `./al`, config-only `.agent-layer/`, required `docs/agent-layer/` memory, always-sync-on-run.
- Created simplified `README.md`, `DECISIONS.md`, and `ROADMAP.md` as the foundation for the Go rewrite.

## Phase 2 ✅ — Repository installer + skeleton (single command install)
- Implemented repo initialization (`al install`), gitignore management, and template seeding.
- Added release workflow + installer script for repo-local binary installation.

## Phase 3 ✅ — Core sync engine (parity with current generators)
- Implemented config parsing/validation, instruction + workflow parsing, and deterministic generators for all clients.
- Wired the internal MCP prompt server into Gemini/Claude configs and added golden-file tests.

## Phase 4 ✅ — Agent launchers (Gemini/Claude/Codex/VS Code/Antigravity)
- Added shared launch pipeline and client launchers with per-agent model/effort wiring.
- Ensured Antigravity runs with generated instructions and slash commands only.

## Phase 5 — MCP servers + approvals (v2 simplified model)

### Goal
- MCP server config and the 4-mode approvals policy are easy to understand and consistent.

### Tasks
- [x] Implement `[[mcp.servers]]` projection for both `transport = "http"` and `transport = "stdio"` (including env wiring).
- [x] Implement `${ENV_VAR}` substitution from `.agent-layer/.env` where needed for config generation.
- [ ] Wire `.agent-layer/.env` tokens into generated client configs (client-specific best practice).
- [x] Implement approvals modes: `all`, `mcp`, `commands`, `none` and generate per-client projections.
- [ ] Implement `al doctor` to report missing tokens, disabled servers, and common misconfigurations.

### Exit criteria
- MCP + approvals behavior matches the v2 contract and is documented.

## Phase 6 — UX polish for adoption

### Goal
- High-quality day-to-day UX without increasing conceptual complexity.

### Tasks
- [ ] Implement interactive `al wizard` for agent enablement + model selection + Codex reasoning.
- [ ] Implement shell completions (`al completion bash|zsh|fish|powershell`).
- [x] Implement `al vscode` launcher behavior for CODEX_HOME.
- [ ] Add optional OS-native launcher artifacts (macOS .app, Windows shortcut, Linux .desktop).
- [x] Implement concurrency-safe run dirs: `tmp/agent-layer/runs/<run-id>/` and export `AL_RUN_DIR`.

### Exit criteria
- Tab completion works; wizard works; VS Code Codex extension launch path is reliable; concurrent runs do not collide.

## Phase 7 — Migration + compatibility

### Goal
- Existing users can adopt v2 without rebuilding their mental model.

### Tasks
- [ ] Provide a one-time migration command to translate v1 config layout into v2 layout (where possible).
- [ ] Document “v1 → v2” mapping and breaking changes.
- [ ] Ensure install/upgrade path is stable and reversible.

### Exit criteria
- A v1 repo can migrate to v2 with minimal manual edits and maintain equivalent behavior.
