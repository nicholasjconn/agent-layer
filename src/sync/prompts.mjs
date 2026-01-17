import fs from "node:fs";
import path from "node:path";
import { REGEN_COMMAND } from "./constants.mjs";
import {
  banner,
  parseFrontMatter,
  resolveWorkflowName,
  slugify,
} from "./instructions.mjs";
import { failOutOfDate } from "./outdated.mjs";
import {
  fileExists,
  listFiles,
  mkdirp,
  readUtf8,
  rmrf,
  writeUtf8,
} from "./utils.mjs";

/**
 * @typedef {{ check: boolean, verbose: boolean }} SyncArgs
 */

/**
 * Format a YAML folded block scalar with stable wrapping.
 * @param {string} value
 * @param {number} maxWidth
 * @returns {string}
 */
function formatYamlFoldedBlock(value, maxWidth = 78) {
  const normalized = String(value ?? "").trim();
  if (!normalized) return "  ";

  const words = normalized.split(/\s+/u).filter(Boolean);
  /** @type {string[]} */
  const lines = [];
  let line = "";

  for (const word of words) {
    if (!line) {
      line = word;
      continue;
    }
    if (line.length + 1 + word.length <= maxWidth) {
      line = `${line} ${word}`;
    } else {
      lines.push(line);
      line = word;
    }
  }

  if (line) lines.push(line);
  return lines.map((entry) => `  ${entry}`).join("\n");
}

/**
 * Render VS Code prompt frontmatter.
 * @param {string} name
 * @param {string} description
 * @returns {string}
 */
function renderPromptFrontmatter(name, description) {
  let out = `---\nname: ${name}\n`;
  if (description.trim()) {
    out += `description: >-\n${formatYamlFoldedBlock(description)}\n`;
  }
  out += `---\n\n`;
  return out;
}

/**
 * Check whether a VS Code prompt file is generated.
 * @param {string} promptFile
 * @returns {boolean}
 */
function isGeneratedVscodePrompt(promptFile) {
  if (!fileExists(promptFile)) return false;
  const txt = readUtf8(promptFile);
  return (
    txt.includes("GENERATED FILE") &&
    txt.includes(`Regenerate: ${REGEN_COMMAND}`)
  );
}

/**
 * Generate VS Code prompt files from workflow definitions.
 * @param {string} repoRoot
 * @param {string} workflowsDir
 * @param {SyncArgs} args
 * @returns {void}
 */
export function generateVscodePrompts(repoRoot, workflowsDir, args) {
  if (!fileExists(workflowsDir)) {
    throw new Error(
      `agent-layer sync: missing workflows directory at ${workflowsDir}. ` +
        "Restore .agent-layer/config/workflows before running sync.",
    );
  }
  const promptsRoot = path.join(repoRoot, ".vscode", "prompts");
  mkdirp(promptsRoot);

  const workflowFiles = listFiles(workflowsDir, ".md");
  if (workflowFiles.length === 0) {
    throw new Error(
      `agent-layer sync: no workflow files found in ${workflowsDir}. ` +
        "Add at least one .md file to .agent-layer/config/workflows.",
    );
  }

  const expectedFiles = new Set();
  const slugToName = new Map();

  for (const wfPath of workflowFiles) {
    const md = readUtf8(wfPath);
    const { meta, body } = parseFrontMatter(md, wfPath);

    const name = resolveWorkflowName(meta, wfPath);
    const slug = slugify(name);

    if (slugToName.has(slug)) {
      throw new Error(
        `agent-layer sync: workflow slug collision: "${name}" and "${slugToName.get(slug)}" both map to "${slug}". ` +
          "Rename one workflow name to avoid collisions.",
      );
    }
    slugToName.set(slug, name);

    if (!body.trim().length) {
      throw new Error(`agent-layer sync: workflow body is empty for ${wfPath}`);
    }

    const description = meta.description ? String(meta.description) : "";
    const promptFile = path.join(promptsRoot, `${slug}.prompt.md`);
    expectedFiles.add(promptFile);

    const content =
      renderPromptFrontmatter(name, description) +
      banner(
        `.agent-layer/config/workflows/${path.basename(wfPath)}`,
        REGEN_COMMAND,
      ) +
      `${body.trimEnd()}\n`;

    if (args.check) {
      const old = fileExists(promptFile) ? readUtf8(promptFile) : null;
      if (old !== content) {
        failOutOfDate(
          repoRoot,
          [promptFile],
          "VS Code prompt files are generated from .agent-layer/config/workflows/*.md.",
        );
      }
    } else {
      writeUtf8(promptFile, content);
      if (args.verbose) console.log(`wrote: ${promptFile}`);
    }
  }

  if (fileExists(promptsRoot)) {
    const entries = fs.readdirSync(promptsRoot, { withFileTypes: true });
    for (const entry of entries) {
      if (!entry.isFile()) continue;
      if (!entry.name.endsWith(".prompt.md")) continue;

      const promptFile = path.join(promptsRoot, entry.name);
      if (expectedFiles.has(promptFile)) continue;
      if (!isGeneratedVscodePrompt(promptFile)) continue;

      if (args.check) {
        failOutOfDate(
          repoRoot,
          [promptFile],
          "Stale generated VS Code prompt file found (no matching workflow).",
        );
      }

      rmrf(promptFile);
      if (args.verbose) console.log(`removed stale prompt: ${promptFile}`);
    }
  }
}
