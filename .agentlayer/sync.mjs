#!/usr/bin/env node
/**
 * Agentlayer sync (Node-based generator)
 *
 * Generates per-client shim files from `.agentlayer/` sources:
 * - AGENTS.md
 * - CLAUDE.md
 * - GEMINI.md
 * - .github/copilot-instructions.md
 *
 * Generates per-client MCP configuration from `.agentlayer/mcp/servers.json`:
 * - .mcp.json              (Claude Code)
 * - .gemini/settings.json  (Gemini CLI)
 * - .vscode/mcp.json       (VS Code / Copilot Chat)
 *
 * Generates Codex Skills from `.agentlayer/workflows/*.md`:
 * - .codex/skills/<workflow>/SKILL.md
 *
 * Usage:
 *   node .agentlayer/sync.mjs
 *   node .agentlayer/sync.mjs --check
 *   node .agentlayer/sync.mjs --verbose
 */

import fs from "node:fs";
import path from "node:path";
import process from "node:process";

function usageAndExit(code) {
  console.error("Usage: node .agentlayer/sync.mjs [--check] [--verbose]");
  process.exit(code);
}

function parseArgs(argv) {
  const args = { check: false, verbose: false };
  for (const a of argv.slice(2)) {
    if (a === "--check") args.check = true;
    else if (a === "--verbose") args.verbose = true;
    else if (a === "-h" || a === "--help") usageAndExit(0);
    else usageAndExit(2);
  }
  return args;
}

function assert(cond, msg) {
  if (!cond) throw new Error(`agentlayer sync: ${msg}`);
}

function fileExists(p) {
  try {
    fs.accessSync(p, fs.constants.F_OK);
    return true;
  } catch {
    return false;
  }
}

function mkdirp(p) {
  fs.mkdirSync(p, { recursive: true });
}

function readUtf8(p) {
  return fs.readFileSync(p, "utf8");
}

function writeUtf8(p, content) {
  mkdirp(path.dirname(p));
  fs.writeFileSync(p, content, "utf8");
}

function rmrf(p) {
  fs.rmSync(p, { recursive: true, force: true });
}

function findRepoRoot(startDir) {
  let dir = path.resolve(startDir);
  for (let i = 0; i < 50; i++) {
    if (fileExists(path.join(dir, ".agentlayer"))) return dir;
    const parent = path.dirname(dir);
    if (parent === dir) break;
    dir = parent;
  }
  return null;
}

function listFiles(dir, suffix) {
  if (!fileExists(dir)) return [];
  return fs
    .readdirSync(dir)
    .filter((f) => f.endsWith(suffix))
    .sort()
    .map((f) => path.join(dir, f));
}

/**
 * Parse strict frontmatter (a small subset of YAML).
 *
 * Supported:
 * ---
 * name: find-issues
 * description: something
 * ---
 *
 * Rules:
 * - If the first line is not "---", frontmatter is treated as absent.
 * - If the first line is "---", a closing "---" is required.
 * - Only keys "name" and "description" are allowed.
 * - Duplicate keys are errors.
 * - If "name" exists, it must be non-empty.
 * - Empty bodies are errors (when used to generate a prompt/skill).
 */
function parseFrontMatter(markdown, sourcePath = "<workflow>") {
  // Strip UTF-8 BOM if present.
  if (markdown.charCodeAt(0) === 0xfeff) {
    markdown = markdown.slice(1);
  }

  const lines = markdown.split(/\r?\n/);

  if (lines.length === 0 || lines[0] !== "---") {
    return { meta: {}, body: markdown };
  }

  const meta = {};
  let i = 1;
  let closed = false;

  for (; i < lines.length; i++) {
    const line = lines[i];

    if (line === "---") {
      closed = true;
      break;
    }

    if (!line.trim()) continue; // allow blank lines in header
    if (line.trim().startsWith("#")) continue; // allow comment lines

    const idx = line.indexOf(":");
    if (idx === -1) {
      throw new Error(
        `agentlayer sync: invalid workflow frontmatter in ${sourcePath}: expected "key: value" but got: ${line}`
      );
    }

    const k = line.slice(0, idx).trim();
    const v = line.slice(idx + 1).trim();

    if (!k) {
      throw new Error(
        `agentlayer sync: invalid workflow frontmatter in ${sourcePath}: empty key in line: ${line}`
      );
    }

    if (!["name", "description"].includes(k)) {
      throw new Error(
        `agentlayer sync: unsupported workflow frontmatter key "${k}" in ${sourcePath}. Allowed keys: name, description`
      );
    }

    if (Object.prototype.hasOwnProperty.call(meta, k)) {
      throw new Error(
        `agentlayer sync: duplicate workflow frontmatter key "${k}" in ${sourcePath}`
      );
    }

    const vv = v.replace(/^["']|["']$/g, "");
    if (k === "name" && !vv.trim()) {
      throw new Error(
        `agentlayer sync: workflow frontmatter "name" must be non-empty in ${sourcePath}`
      );
    }
    meta[k] = vv;
  }

  if (!closed) {
    throw new Error(
      `agentlayer sync: invalid workflow frontmatter in ${sourcePath}: missing closing "---"`
    );
  }

  const body = lines.slice(i + 1).join("\n").replace(/^\n+/, "");
  return { meta, body };
}

function slugify(s) {
  return (
    String(s)
      .trim()
      .toLowerCase()
      .replace(/[^a-z0-9._-]+/g, "-")
      .replace(/-+/g, "-")
      .replace(/^-|-$/g, "") || "workflow"
  );
}

function banner(sourceHint, regenHint) {
  return (
    `<!--\n` +
    `  GENERATED FILE - DO NOT EDIT DIRECTLY\n` +
    `  Source: ${sourceHint}\n` +
    `  Regenerate: ${regenHint}\n` +
    `-->\n\n`
  );
}

function concatInstructions(instructionsDir) {
  const files = listFiles(instructionsDir, ".md");
  const chunks = [];
  for (const f of files) {
    const name = path.basename(f);
    chunks.push(
      `<!-- BEGIN: ${name} -->\n` +
        readUtf8(f).trimEnd() +
        `\n<!-- END: ${name} -->\n`
    );
  }
  return chunks.join("\n").trimEnd() + "\n";
}

function relPath(repoRoot, absPath) {
  const r = path.relative(repoRoot, absPath);
  return r.split(path.sep).join("/");
}

function failOutOfDate(repoRoot, changedAbsPaths, extraMessage = "") {
  const rels = changedAbsPaths.map((p) => relPath(repoRoot, p));

  const instructionShims = [];
  const mcpConfigs = [];
  const codexSkills = [];
  const other = [];

  for (const rp of rels) {
    if (
      rp === "AGENTS.md" ||
      rp === "CLAUDE.md" ||
      rp === "GEMINI.md" ||
      rp === ".github/copilot-instructions.md"
    ) {
      instructionShims.push(rp);
    } else if (
      rp === ".mcp.json" ||
      rp === ".gemini/settings.json" ||
      rp === ".vscode/mcp.json"
    ) {
      mcpConfigs.push(rp);
    } else if (rp.startsWith(".codex/skills/")) {
      codexSkills.push(rp);
    } else {
      other.push(rp);
    }
  }

  console.error("agentlayer sync: generated files are out of date.");
  if (extraMessage) console.error(extraMessage);
  console.error("");
  console.error("Do NOT edit generated files directly.");
  console.error("");

  if (instructionShims.length) {
    console.error("Instruction shims (edit: .agentlayer/instructions/*.md):");
    for (const p of instructionShims) console.error(`  - ${p}`);
    console.error("");
  }

  if (mcpConfigs.length) {
    console.error("MCP config files (edit: .agentlayer/mcp/servers.json):");
    for (const p of mcpConfigs) console.error(`  - ${p}`);
    console.error("");
  }

  if (codexSkills.length) {
    console.error("Codex skills (edit: .agentlayer/workflows/*.md):");
    for (const p of codexSkills) console.error(`  - ${p}`);
    console.error("");
  }

  if (other.length) {
    console.error("Other generated files:");
    for (const p of other) console.error(`  - ${p}`);
    console.error("");
  }

  console.error("Fix:");
  console.error("  1) Edit the source-of-truth file(s) listed above");
  console.error("  2) Run: node .agentlayer/sync.mjs");
  console.error("");
  console.error("If you accidentally edited a generated file, revert it (example):");
  console.error("  git checkout -- .mcp.json");
  console.error("");
  console.error("Files that would change:");
  for (const p of rels.sort()) console.error(`  - ${p}`);

  process.exit(1);
}

function diffOrWrite(outputs, args, repoRoot) {
  const changed = [];
  for (const [outPath, content] of outputs) {
    const old = fileExists(outPath) ? readUtf8(outPath) : null;
    if (old !== content) {
      changed.push(outPath);
      if (!args.check) writeUtf8(outPath, content);
    }
    if (args.verbose) {
      console.log(
        `${old === content ? "ok" : args.check ? "needs-update" : "wrote"}: ${outPath}`
      );
    }
  }

  if (args.check && changed.length) {
    failOutOfDate(repoRoot, changed);
  }
  return changed.length > 0;
}

function validateServerCatalog(parsed, filePath) {
  assert(parsed && typeof parsed === "object", `${filePath} must contain a JSON object`);

  if (parsed.defaults !== undefined) {
    assert(
      parsed.defaults && typeof parsed.defaults === "object" && !Array.isArray(parsed.defaults),
      `${filePath}: defaults must be an object`
    );
    if (parsed.defaults.vscodeEnvFile !== undefined) {
      assert(
        typeof parsed.defaults.vscodeEnvFile === "string",
        `${filePath}: defaults.vscodeEnvFile must be a string`
      );
    }
    // Back-compat: allow defaults.geminiTrust but prefer defaults.trust.
    if (parsed.defaults.trust !== undefined) {
      assert(
        typeof parsed.defaults.trust === "boolean",
        `${filePath}: defaults.trust must be boolean`
      );
    }
    if (parsed.defaults.geminiTrust !== undefined) {
      assert(
        typeof parsed.defaults.geminiTrust === "boolean",
        `${filePath}: defaults.geminiTrust must be boolean`
      );
    }
  }

  assert(Array.isArray(parsed.servers), `${filePath}: servers must be an array`);

  const seen = new Set();
  for (const s of parsed.servers) {
    assert(s && typeof s === "object" && !Array.isArray(s), `${filePath}: each server must be an object`);
    assert(typeof s.name === "string" && s.name.trim(), `${filePath}: server.name must be a non-empty string`);
    assert(!seen.has(s.name), `${filePath}: duplicate server name "${s.name}"`);
    seen.add(s.name);

    if (s.enabled !== undefined) {
      assert(typeof s.enabled === "boolean", `${filePath}: ${s.name}.enabled must be boolean`);
    }
    if (s.trust !== undefined) {
      assert(typeof s.trust === "boolean", `${filePath}: ${s.name}.trust must be boolean`);
    }
    // Back-compat: per-server geminiTrust (prefer trust)
    if (s.geminiTrust !== undefined) {
      assert(typeof s.geminiTrust === "boolean", `${filePath}: ${s.name}.geminiTrust must be boolean`);
    }

    if (s.transport !== undefined) {
      assert(typeof s.transport === "string", `${filePath}: ${s.name}.transport must be a string`);
      assert(
        s.transport === "stdio",
        `${filePath}: ${s.name}.transport must be "stdio" (this generator supports only stdio currently)`
      );
    }

    assert(typeof s.command === "string" && s.command.trim(), `${filePath}: ${s.name}.command must be a non-empty string`);

    if (s.args !== undefined) {
      assert(Array.isArray(s.args), `${filePath}: ${s.name}.args must be an array`);
      assert(s.args.every((x) => typeof x === "string"), `${filePath}: ${s.name}.args must be string[]`);
    }

    if (s.envVars !== undefined) {
      assert(Array.isArray(s.envVars), `${filePath}: ${s.name}.envVars must be an array`);
      assert(s.envVars.every((x) => typeof x === "string"), `${filePath}: ${s.name}.envVars must be string[]`);
    }

    // Optional Gemini allow/deny lists.
    if (s.includeTools !== undefined) {
      assert(Array.isArray(s.includeTools), `${filePath}: ${s.name}.includeTools must be an array`);
      assert(s.includeTools.every((x) => typeof x === "string"), `${filePath}: ${s.name}.includeTools must be string[]`);
    }
    if (s.excludeTools !== undefined) {
      assert(Array.isArray(s.excludeTools), `${filePath}: ${s.name}.excludeTools must be an array`);
      assert(s.excludeTools.every((x) => typeof x === "string"), `${filePath}: ${s.name}.excludeTools must be string[]`);
    }
    if (s.includeTools !== undefined && s.excludeTools !== undefined) {
      throw new Error(
        `agentlayer sync: ${filePath}: ${s.name} cannot set both includeTools and excludeTools`
      );
    }
  }
}

function loadServerCatalog(repoRoot) {
  const filePath = path.join(repoRoot, ".agentlayer", "mcp", "servers.json");
  if (!fileExists(filePath)) return { defaults: {}, servers: [] };
  const parsed = JSON.parse(readUtf8(filePath));
  validateServerCatalog(parsed, filePath);
  const servers = Array.isArray(parsed.servers) ? parsed.servers : [];
  const defaults = parsed.defaults ?? {};
  return { defaults, servers };
}

function enabledServers(servers) {
  const enabled = servers.filter(
    (s) => s && s.name && (s.enabled === undefined || s.enabled === true)
  );

  // Validate schema to avoid silently generating broken configs.
  for (const s of enabled) {
    const transport = s.transport ?? "stdio";
    if (transport !== "stdio") {
      throw new Error(
        `agentlayer sync: unsupported transport '${transport}' for server '${s.name}'. ` +
          `This generator currently supports only stdio servers.`
      );
    }
    if (!s.command || typeof s.command !== "string") {
      throw new Error(`agentlayer sync: server '${s.name}' missing valid 'command'.`);
    }
    if (s.args !== undefined && !Array.isArray(s.args)) {
      throw new Error(`agentlayer sync: server '${s.name}' has non-array 'args'.`);
    }
    if (s.envVars !== undefined && !Array.isArray(s.envVars)) {
      throw new Error(`agentlayer sync: server '${s.name}' has non-array 'envVars'.`);
    }
  }

  return enabled;
}

function renderMcpConfigs(repoRoot, catalog) {
  const defaults = catalog.defaults ?? {};
  const servers = enabledServers(catalog.servers ?? []);

  // NOTE: VS Code can load env from an envFile. Default remains project root .env
  // unless you set defaults.vscodeEnvFile to "${workspaceFolder}/.agentlayer/.env".
  const vscodeEnvFile = defaults.vscodeEnvFile ?? "${workspaceFolder}/.env";

  // Single generic trust field (applied to Gemini today; ignored elsewhere).
  // Back-compat: accept defaults.geminiTrust if defaults.trust is not present.
  const defaultTrust =
    defaults.trust === undefined
      ? (defaults.geminiTrust === undefined ? false : Boolean(defaults.geminiTrust))
      : Boolean(defaults.trust);

  // VS Code
  const vscode = { servers: {} };
  for (const s of servers) {
    vscode.servers[s.name] = {
      type: "stdio",
      command: s.command,
      args: s.args ?? [],
      envFile: vscodeEnvFile,
    };
  }

  // Claude Code
  const claude = { mcpServers: {} };
  for (const s of servers) {
    const env = {};
    for (const v of s.envVars ?? []) env[v] = `\${${v}}`;
    claude.mcpServers[s.name] = {
      command: s.command,
      args: s.args ?? [],
      ...(Object.keys(env).length ? { env } : {}),
    };
  }

  // Gemini CLI
  const gemini = { mcpServers: {} };
  for (const s of servers) {
    const env = {};
    for (const v of s.envVars ?? []) env[v] = `\${${v}}`;

    // Back-compat: per-server geminiTrust if trust is not present.
    const trust =
      s.trust === undefined
        ? (s.geminiTrust === undefined ? defaultTrust : Boolean(s.geminiTrust))
        : Boolean(s.trust);

    const entry = {
      command: s.command,
      args: s.args ?? [],
      ...(Object.keys(env).length ? { env } : {}),
      trust,
    };

    if (Array.isArray(s.includeTools)) entry.includeTools = s.includeTools;
    if (Array.isArray(s.excludeTools)) entry.excludeTools = s.excludeTools;

    gemini.mcpServers[s.name] = entry;
  }

  return [
    [
      path.join(repoRoot, ".vscode", "mcp.json"),
      JSON.stringify(vscode, null, 2) + "\n",
    ],
    [path.join(repoRoot, ".mcp.json"), JSON.stringify(claude, null, 2) + "\n"],
    [
      path.join(repoRoot, ".gemini", "settings.json"),
      JSON.stringify(gemini, null, 2) + "\n",
    ],
  ];
}

function isGeneratedCodexSkill(skillDir) {
  const p = path.join(skillDir, "SKILL.md");
  if (!fileExists(p)) return false;
  const txt = readUtf8(p);
  return (
    txt.includes("GENERATED FILE - DO NOT EDIT DIRECTLY") &&
    txt.includes("Regenerate: node .agentlayer/sync.mjs")
  );
}

function generateCodexSkills(repoRoot, workflowsDir, args) {
  const codexSkillsRoot = path.join(repoRoot, ".codex", "skills");
  mkdirp(codexSkillsRoot);

  const workflowFiles = listFiles(workflowsDir, ".md");
  const expectedFolders = new Set();
  const slugToName = new Map();

  for (const wfPath of workflowFiles) {
    const md = readUtf8(wfPath);
    const { meta, body } = parseFrontMatter(md, wfPath);

    const fallbackName = path.basename(wfPath, ".md");
    const name = (meta.name && String(meta.name).trim()) ? String(meta.name).trim() : fallbackName;
    const description = meta.description ? String(meta.description) : "";
    const folder = slugify(name);

    if (slugToName.has(folder)) {
      throw new Error(
        `agentlayer sync: workflow slug collision: "${name}" and "${slugToName.get(folder)}" both map to "${folder}". ` +
          `Rename one workflow name to avoid collisions.`
      );
    }
    slugToName.set(folder, name);

    expectedFolders.add(folder);

    const skillDir = path.join(codexSkillsRoot, folder);
    const skillFile = path.join(skillDir, "SKILL.md");

    assert(name.trim().length > 0, `workflow name resolved to empty for ${wfPath}`);
    assert(body.trim().length > 0, `workflow body is empty for ${wfPath}`);

    // YAML frontmatter must be first line for compatibility.
    const content =
      `---\n` +
      `name: ${name}\n` +
      `description: ${description}\n` +
      `---\n\n` +
      banner(`.agentlayer/workflows/${path.basename(wfPath)}`, "node .agentlayer/sync.mjs") +
      `# ${name}\n\n` +
      (description ? `${description}\n\n` : "") +
      `${body.trimEnd()}\n`;

    mkdirp(skillDir);

    if (args.check) {
      const old = fileExists(skillFile) ? readUtf8(skillFile) : null;
      if (old !== content) {
        failOutOfDate(
          repoRoot,
          [skillFile],
          "Codex skills are generated from .agentlayer/workflows/*.md."
        );
      }
    } else {
      writeUtf8(skillFile, content);
      if (args.verbose) console.log(`wrote: ${skillFile}`);
    }
  }

  // Remove stale generated skills (no matching workflow).
  if (fileExists(codexSkillsRoot)) {
    const entries = fs.readdirSync(codexSkillsRoot, { withFileTypes: true });
    for (const e of entries) {
      if (!e.isDirectory()) continue;
      const dirName = e.name;
      if (expectedFolders.has(dirName)) continue;

      const skillDir = path.join(codexSkillsRoot, dirName);
      if (!isGeneratedCodexSkill(skillDir)) continue;

      if (args.check) {
        const maybeSkill = path.join(skillDir, "SKILL.md");
        failOutOfDate(
          repoRoot,
          [fileExists(maybeSkill) ? maybeSkill : skillDir],
          "Stale generated Codex skill found (no matching workflow)."
        );
      }

      rmrf(skillDir);
      if (args.verbose) console.log(`removed stale generated Codex skill: ${skillDir}`);
    }
  }
}

function main() {
  const args = parseArgs(process.argv);
  const repoRoot = findRepoRoot(process.cwd());
  if (!repoRoot) {
    console.error("agentlayer sync: could not find repo root containing .agentlayer/");
    process.exit(2);
  }

  const instructionsDir = path.join(repoRoot, ".agentlayer", "instructions");
  const workflowsDir = path.join(repoRoot, ".agentlayer", "workflows");

  const unified =
    banner(".agentlayer/instructions/*.md", "node .agentlayer/sync.mjs") +
    concatInstructions(instructionsDir);

  const outputs = [
    [path.join(repoRoot, "AGENTS.md"), unified],
    [path.join(repoRoot, "CLAUDE.md"), unified],
    [path.join(repoRoot, "GEMINI.md"), unified],
    [path.join(repoRoot, ".github", "copilot-instructions.md"), unified],
  ];

  const catalog = loadServerCatalog(repoRoot);
  outputs.push(...renderMcpConfigs(repoRoot, catalog));

  diffOrWrite(outputs, args, repoRoot);
  generateCodexSkills(repoRoot, workflowsDir, args);

  if (!args.check) {
    console.log("agentlayer sync: updated shims + MCP configs + Codex skills");
  }
}

main();
