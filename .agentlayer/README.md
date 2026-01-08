# Agentlayer (repo-local agent standardization)

This repository includes a repo-local standardization layer under `.agentlayer/`.

## What this is for

**Goal:** make agentic tooling consistent across Claude Code, Gemini CLI, VS Code/Copilot, and Codex by keeping a **single source of truth** in the repo, then generating the small set of per-client shim/config files those tools require.

**Primary uses**
- A unified instruction set (system/developer-style guidance) usable across tools.
- Repeatable “workflows” exposed as:
  - MCP prompts (commands) in clients that support MCP prompts.
  - Codex Skills (repo-local) for Codex.
- A repo-committed MCP server catalog, projected into each client’s config format.
- A lightweight setup flow that works in any project repo.

## Prerequisites

Required:
- **Node.js + npm** (recommended: manage a pinned version via `mise`, `asdf`, `volta`, or `nvm`)
- **git** (recommended; required for enabling hooks)

Optional (depending on which clients you use):
- VS Code (Copilot Chat)
- Gemini CLI
- Claude Code
- Codex CLI / Codex VS Code extension

## Quickstart

From repo root:

1) **Run setup (installs deps, enables hooks, verifies everything)**

```bash
chmod +x .agentlayer/setup.sh
./.agentlayer/setup.sh
```

2) **Create your Agentlayer env file (recommended: agent-only secrets)**

```bash
# Recommended: keep agent-only secrets separate from project env vars
cp .env.example .agentlayer/.env
# edit .agentlayer/.env; do not commit it
```

If you also use a project/dev `.env` for your application, keep it in `./.env` and do not mix agent-only tokens into it.

3) **Create the repo-local launcher `./al` (recommended)**

```bash
chmod +x ./al
```

Default behavior: sync every run via `node .agentlayer/sync.mjs`, then load `.agentlayer/.env` and exec the command (via `.agentlayer/with-env.sh`).
If you want a different default, open `al` and uncomment exactly one option block.

Examples:

```bash
./al gemini
./al claude
./al codex
```

4) **Edit sources of truth**
- Unified instructions: `.agentlayer/instructions/*.md`
- Workflows: `.agentlayer/workflows/*.md`
- MCP server catalog: `.agentlayer/mcp/servers.json`

5) **Regenerate after changes**

```bash
node .agentlayer/sync.mjs
```

## How to use (day-to-day)

### Prefer `./al` for running CLIs

`./al` is intentionally minimal. By default it:

1) Runs `node .agentlayer/sync.mjs`
2) Loads `.agentlayer/.env` via `.agentlayer/with-env.sh`
3) Execs the command

Examples:

```bash
./al gemini
./al claude
./al codex
```

To change defaults, edit `al` and uncomment exactly one option block (env-only, sync-only, sync-check + regen, or include project env).
For a one-off project env run, use:

```bash
./.agentlayer/with-env.sh --project-env gemini
```

`with-env.sh` resolves the repo root for env file paths and does not change your working directory.

### Debugging trick (verify env vars are being loaded)

```bash
./al env | grep -E 'GITHUB_TOKEN|CONTEXT7_API_KEY'
```

### What files you should and should not edit

**Edit these (sources of truth):**
- `.agentlayer/instructions/*.md`
- `.agentlayer/workflows/*.md`
- `.agentlayer/mcp/servers.json`

**Do not edit these directly (generated):**
- `AGENTS.md`
- `CLAUDE.md`
- `GEMINI.md`
- `.github/copilot-instructions.md`
- `.mcp.json`
- `.gemini/settings.json`
- `.vscode/mcp.json`
- `.codex/skills/*/SKILL.md`

If you accidentally edited a generated file, revert it (example):

```bash
git checkout -- .mcp.json
```

### Instruction file ordering (why the numbers)

`sync.mjs` concatenates `.agentlayer/instructions/*.md` in **lexicographic order**.

Numeric prefixes (e.g. `00_`, `10_`, `20_`) ensure a **stable, predictable ordering** without requiring a separate manifest/config file. If you add new instruction fragments, follow the same pattern.

## Refresh / restart guidance (failure modes)

General rule:
- After changing `.agentlayer/**` → run `node .agentlayer/sync.mjs` (or run your CLI via `./al ...`) → then refresh/restart the client as needed.

### MCP prompt server (workflows as “commands”)

Workflows are exposed as MCP prompts by:
- `.agentlayer/mcp/agentlayer-prompts/server.mjs`

If you changed `.agentlayer/workflows/*.md`:
- run `node .agentlayer/sync.mjs` (or `./al <cmd>`)
- then refresh MCP discovery in your client (or restart the client/session)

### Client-specific notes

**Gemini CLI**
- If instructions don’t update: refresh memory (often `/memory refresh`).
- If prompts/tools don’t update: refresh MCP (often `/mcp refresh`).

**VS Code / Copilot Chat**
- MCP prompts can be invoked as `/mcp.<server>.<prompt>`.
- If prompts/tools look stale: restart the MCP server and/or use “Chat: Reset Cached Tools”.

**Claude Code**
- If MCP config changed: restart Claude Code.
- If prompts/tools look stale: reconnect or restart the session.

**Codex**
- Codex does not reliably expose MCP prompts as slash commands in the client UI.
- Workflows are available via generated **repo-local Skills** under `.codex/skills/`.
- After updating workflows and running sync, restart Codex if it doesn’t pick up skills immediately.

## What’s inside this repository

### Source-of-truth directories
- `.agentlayer/instructions/`  
  Unified instruction fragments (concatenated into shims).
- `.agentlayer/workflows/`  
  Workflow definitions (exposed as MCP prompts; also used to generate Codex skills).
- `.agentlayer/mcp/servers.json`  
  Canonical MCP server list (no secrets inside).

### Generated outputs
- Instruction shims:
  - `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `.github/copilot-instructions.md`
- MCP configs projected per client:
  - `.mcp.json`, `.gemini/settings.json`, `.vscode/mcp.json`
- Codex skills:
  - `.codex/skills/*/SKILL.md`

### Scripts
- `.agentlayer/setup.sh`  
  One-shot setup (install MCP deps, enable hooks, validate).
- `.agentlayer/sync.mjs`  
  Generator (“build”) for all shims/configs/skills.
- `./al`  
  Repo-local launcher; default is sync every run, then load `.agentlayer/.env`. Edit `al` to choose a different behavior.

## FAQ / Troubleshooting

### “I edited generated JSON and now things are broken.”
Generated JSON files (`.mcp.json`, `.vscode/mcp.json`, `.gemini/settings.json`) do not allow comments and may be strict-parsed by clients.

Fix:
1) revert the generated file(s)
2) edit the source-of-truth (`.agentlayer/mcp/servers.json`)
3) run `node .agentlayer/sync.mjs`

### “I edited instructions but the agent didn’t follow them.”
- Did you run `node .agentlayer/sync.mjs` (or run via `./al ...`)?
- Did you restart the session/client (many tools read system instructions at session start)?
- For Gemini CLI, run memory refresh; for others, start a new session.

### “I edited workflows but the prompt/command list didn’t update.”
- Run `node .agentlayer/sync.mjs`
- Refresh MCP discovery (Gemini: `/mcp refresh`; VS Code: restart server/reset cached tools; Claude: reconnect/restart).

### “Commits are failing after enabling hooks.”
The hook runs:

```bash
node .agentlayer/sync.mjs --check
```

If it fails, run:

```bash
node .agentlayer/sync.mjs
```

Then commit again.

### “Can I rename the instruction files?”
Yes. Keep numeric prefixes if you want stable ordering without changing `sync.mjs`.
