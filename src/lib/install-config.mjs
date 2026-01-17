import fs from "node:fs";
import path from "node:path";
import readline from "node:readline";
import { fileExists, readUtf8 } from "../sync/utils.mjs";
import { loadAgentConfig } from "./agent-config.mjs";

export const INSTALL_CONFIG_USAGE = [
  "Usage:",
  "  ./al --install-config [--force] [--new-install] [--non-interactive] [--parent-root <path>] [--temp-parent-root] [--agent-layer-root <path>]",
].join("\n");

const GITIGNORE_BLOCK = [
  "# >>> agent-layer",
  ".agent-layer/",
  "",
  "# Agent Layer launcher",
  "al",
  "",
  "# Agent Layer-generated instruction shims",
  "AGENTS.md",
  "CLAUDE.md",
  "GEMINI.md",
  ".github/copilot-instructions.md",
  "",
  "# Agent Layer-generated client configs + artifacts",
  ".mcp.json",
  ".codex/",
  ".gemini/",
  ".claude/",
  ".vscode/mcp.json",
  ".vscode/prompts/",
  "# <<< agent-layer",
].join("\n");

const MEMORY_FILES = [
  "ISSUES.md",
  "FEATURES.md",
  "ROADMAP.md",
  "DECISIONS.md",
  "COMMANDS.md",
];

/**
 * @typedef {{ force: boolean, newInstall: boolean, nonInteractive: boolean }} InstallConfigArgs
 */

/**
 * Parse arguments for the install-config subcommand.
 * @param {string[]} argv
 * @returns {InstallConfigArgs}
 */
export function parseInstallConfigArgs(argv) {
  const args = {
    force: false,
    newInstall: false,
    nonInteractive: false,
  };
  for (const arg of argv) {
    if (arg === "--force") {
      args.force = true;
      continue;
    }
    if (arg === "--new-install") {
      args.newInstall = true;
      continue;
    }
    if (arg === "--non-interactive") {
      args.nonInteractive = true;
      continue;
    }
    throw new Error(`agent-layer cli: install-config unknown argument: ${arg}`);
  }
  return args;
}

/**
 * Print a message.
 * @param {string} message
 * @returns {void}
 */
function say(message) {
  process.stdout.write(`${message}\n`);
}

/**
 * Ensure the .env exists under the agent-layer root.
 * @param {string} agentLayerRoot
 * @returns {void}
 */
function ensureEnvFile(agentLayerRoot) {
  const envPath = path.join(agentLayerRoot, ".env");
  if (fileExists(envPath)) {
    say("==> .agent-layer/.env already exists; leaving as-is");
    return;
  }
  const examplePath = path.join(agentLayerRoot, ".env.example");
  if (!fileExists(examplePath)) {
    throw new Error("Missing .agent-layer/.env.example; cannot create .env");
  }
  fs.copyFileSync(examplePath, envPath);
  say("==> Created .agent-layer/.env from .env.example");
}

/**
 * Create a readline interface for prompts.
 * @returns {import(\"node:readline\").Interface}
 */
function createPromptInterface() {
  return readline.createInterface({
    input: process.stdin,
    output: process.stdout,
  });
}

/**
 * Prompt for a yes/no response.
 * @param {import(\"node:readline\").Interface} rl
 * @param {string} prompt
 * @param {boolean} defaultYes
 * @returns {Promise<boolean>}
 */
async function promptYesNo(rl, prompt, defaultYes) {
  while (true) {
    const answer = await new Promise((resolve) => {
      rl.question(prompt, resolve);
    });
    const normalized = String(answer ?? "").trim();
    if (!normalized) return defaultYes;
    if (/^(y|yes)$/i.test(normalized)) return true;
    if (/^(n|no)$/i.test(normalized)) return false;
    say("Please answer y or n.");
  }
}

/**
 * Ensure a memory file exists, optionally prompting for replacement.
 * @param {string} parentRoot
 * @param {string} agentLayerRoot
 * @param {string} name
 * @param {boolean} interactive
 * @param {import(\"node:readline\").Interface | null} rl
 * @returns {Promise<void>}
 */
async function ensureMemoryFile(
  parentRoot,
  agentLayerRoot,
  name,
  interactive,
  rl,
) {
  const filePath = path.join(parentRoot, "docs", name);
  const templatePath = path.join(
    agentLayerRoot,
    "config",
    "templates",
    "docs",
    name,
  );
  const relPath = path.relative(parentRoot, filePath);

  if (!fileExists(templatePath)) {
    throw new Error(
      `Missing template: ${path.relative(agentLayerRoot, templatePath)}`,
    );
  }

  if (fileExists(filePath)) {
    if (!interactive || !rl) {
      say(`==> ${relPath} exists; leaving as-is (no TTY to confirm)`);
      return;
    }
    const keep = await promptYesNo(
      rl,
      `${relPath} exists. Keep it? [Y/n] `,
      true,
    );
    if (keep) {
      say(`==> Keeping existing ${relPath}`);
      return;
    }
    fs.mkdirSync(path.dirname(filePath), { recursive: true });
    fs.copyFileSync(templatePath, filePath);
    say(`==> Replaced ${relPath} with template`);
    return;
  }

  fs.mkdirSync(path.dirname(filePath), { recursive: true });
  fs.copyFileSync(templatePath, filePath);
  say(`==> Created ${relPath} from template`);
}

/**
 * Ensure all memory files exist using the templates.
 * @param {string} parentRoot
 * @param {string} agentLayerRoot
 * @param {boolean} interactive
 * @param {import(\"node:readline\").Interface | null} rl
 * @returns {Promise<void>}
 */
async function ensureMemoryFiles(parentRoot, agentLayerRoot, interactive, rl) {
  say("==> Ensuring project memory files exist");
  for (const name of MEMORY_FILES) {
    await ensureMemoryFile(parentRoot, agentLayerRoot, name, interactive, rl);
  }
}

/**
 * Write the repo-local launcher script.
 * @param {string} parentRoot
 * @returns {void}
 */
function writeLauncher(parentRoot) {
  const launcher = [
    "#!/usr/bin/env bash",
    "set -euo pipefail",
    "",
    "# Repo-local launcher.",
    "# This script delegates to the managed Agent Layer entrypoint in .agent-layer/.",
    "# If you prefer, replace this file with a symlink to .agent-layer/agent-layer.",
    'SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"',
    'exec "$SCRIPT_DIR/.agent-layer/agent-layer" "$@"',
    "",
  ].join("\n");
  const launcherPath = path.join(parentRoot, "al");
  fs.writeFileSync(launcherPath, launcher);
  fs.chmodSync(launcherPath, 0o755);
}

/**
 * Create or preserve the repo-local launcher based on force.
 * @param {string} parentRoot
 * @param {boolean} force
 * @returns {void}
 */
function ensureLauncher(parentRoot, force) {
  const launcherPath = path.join(parentRoot, "al");
  if (fileExists(launcherPath)) {
    if (force) {
      say("==> Overwriting ./al");
      writeLauncher(parentRoot);
    } else {
      say("==> NOTE: ./al already exists; not overwriting.");
      say("==> Re-run the installer with --force to replace ./al.");
    }
    return;
  }
  say("==> Creating ./al");
  writeLauncher(parentRoot);
}

/**
 * Update the agent-layer block in .gitignore.
 * @param {string} parentRoot
 * @returns {void}
 */
function updateGitignore(parentRoot) {
  const gitignorePath = path.join(parentRoot, ".gitignore");
  const blockLines = GITIGNORE_BLOCK.split("\n");
  let lines = [];
  if (fileExists(gitignorePath)) {
    const raw = readUtf8(gitignorePath);
    lines = raw === "" ? [] : raw.split(/\r?\n/);
  }

  const output = [];
  let found = false;
  let inBlock = false;
  for (const line of lines) {
    if (line === "# >>> agent-layer") {
      if (!found) {
        output.push(...blockLines);
        found = true;
      }
      inBlock = true;
      continue;
    }
    if (inBlock) {
      if (line === "# <<< agent-layer") {
        inBlock = false;
      }
      continue;
    }
    output.push(line);
  }

  if (!found) {
    if (output.length > 0 && output[output.length - 1] !== "") {
      output.push("");
    }
    output.push(...blockLines);
  }

  const content = `${output.join("\n")}\n`;
  fs.writeFileSync(gitignorePath, content);
}

/**
 * Configure enabled agents during new installs.
 * @param {string} agentLayerRoot
 * @param {boolean} interactive
 * @param {import(\"node:readline\").Interface | null} rl
 * @returns {Promise<void>}
 */
async function configureAgents(agentLayerRoot, interactive, rl) {
  const configPath = path.join(agentLayerRoot, "config", "agents.json");
  if (!fileExists(configPath)) {
    throw new Error(
      "Missing .agent-layer/config/agents.json; cannot configure enabled agents.",
    );
  }

  let enableAll = true;
  let enableGemini = true;
  let enableClaude = true;
  let enableCodex = true;
  let enableVscode = true;

  if (interactive && rl) {
    say("==> Choose which agents to enable (press Enter for yes).");
    enableGemini = await promptYesNo(rl, "Enable Gemini CLI? [Y/n] ", true);
    enableClaude = await promptYesNo(
      rl,
      "Enable Claude Code CLI? [Y/n] ",
      true,
    );
    enableCodex = await promptYesNo(rl, "Enable Codex CLI? [Y/n] ", true);
    enableVscode = await promptYesNo(
      rl,
      "Enable VS Code / Copilot Chat? [Y/n] ",
      true,
    );
    enableAll = false;
  }

  if (!interactive || !rl || enableAll) {
    say("==> Non-interactive install: enabling all agents");
    enableGemini = true;
    enableClaude = true;
    enableCodex = true;
    enableVscode = true;
  }

  const config = loadAgentConfig(agentLayerRoot);
  config.gemini.enabled = enableGemini;
  config.claude.enabled = enableClaude;
  config.codex.enabled = enableCodex;
  config.vscode.enabled = enableVscode;

  fs.writeFileSync(configPath, `${JSON.stringify(config, null, 2)}\n`);
  say("==> Updated .agent-layer/config/agents.json");
}

/**
 * Run installer configuration tasks.
 * @param {{ parentRoot: string, agentLayerRoot: string }} roots
 * @param {InstallConfigArgs} args
 * @returns {Promise<void>}
 */
export async function runInstallConfig(roots, args) {
  const parentRoot = roots.parentRoot;
  const agentLayerRoot = roots.agentLayerRoot;
  const interactive = !args.nonInteractive && Boolean(process.stdin.isTTY);

  let rl = null;
  if (interactive) {
    rl = createPromptInterface();
  }

  try {
    ensureEnvFile(agentLayerRoot);
    await ensureMemoryFiles(parentRoot, agentLayerRoot, interactive, rl);
    ensureLauncher(parentRoot, args.force);
    say("==> Updating .gitignore (agent-layer block)");
    updateGitignore(parentRoot);
    if (args.newInstall) {
      await configureAgents(agentLayerRoot, interactive, rl);
    }
  } finally {
    if (rl) rl.close();
  }
}
