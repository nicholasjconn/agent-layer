import { test, describe, before, after } from "node:test";
import assert from "node:assert";
import fs from "node:fs";
import path from "node:path";
import os from "node:os";
import { runInstallConfig } from "../src/lib/install-config.mjs";

describe("src/lib/install-config.mjs", () => {
  let tmpDir;

  before(() => {
    tmpDir = fs.mkdtempSync(
      path.join(os.tmpdir(), "agent-layer-tests-install-config-"),
    );
  });

  after(() => {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  });

  test("runInstallConfig creates config files and enables agents", async () => {
    const parentRoot = path.join(tmpDir, "parent");
    const agentLayerRoot = path.join(parentRoot, ".agent-layer");
    fs.mkdirSync(agentLayerRoot, { recursive: true });

    fs.writeFileSync(path.join(agentLayerRoot, ".env.example"), "EXAMPLE=1\n");

    const templatesDir = path.join(
      agentLayerRoot,
      "config",
      "templates",
      "docs",
    );
    fs.mkdirSync(templatesDir, { recursive: true });
    const memoryFiles = [
      "ISSUES.md",
      "FEATURES.md",
      "ROADMAP.md",
      "DECISIONS.md",
      "COMMANDS.md",
    ];
    for (const name of memoryFiles) {
      fs.writeFileSync(path.join(templatesDir, name), `template-${name}\n`);
    }

    const agentsDir = path.join(agentLayerRoot, "config");
    fs.mkdirSync(agentsDir, { recursive: true });
    fs.writeFileSync(
      path.join(agentsDir, "agents.json"),
      `${JSON.stringify(
        {
          gemini: { enabled: false },
          claude: { enabled: false },
          codex: { enabled: false },
          vscode: { enabled: false },
        },
        null,
        2,
      )}\n`,
    );

    await runInstallConfig(
      { parentRoot, agentLayerRoot },
      { force: true, newInstall: true, nonInteractive: true },
    );

    assert.strictEqual(
      fs.readFileSync(path.join(agentLayerRoot, ".env"), "utf8"),
      "EXAMPLE=1\n",
    );

    for (const name of memoryFiles) {
      const content = fs.readFileSync(
        path.join(parentRoot, "docs", name),
        "utf8",
      );
      assert.strictEqual(content, `template-${name}\n`);
    }

    const launcher = fs.readFileSync(path.join(parentRoot, "al"), "utf8");
    assert.ok(launcher.includes(".agent-layer/agent-layer"));

    const gitignore = fs.readFileSync(
      path.join(parentRoot, ".gitignore"),
      "utf8",
    );
    assert.ok(gitignore.includes("# >>> agent-layer"));

    const agents = JSON.parse(
      fs.readFileSync(path.join(agentsDir, "agents.json"), "utf8"),
    );
    assert.strictEqual(agents.gemini.enabled, true);
    assert.strictEqual(agents.claude.enabled, true);
    assert.strictEqual(agents.codex.enabled, true);
    assert.strictEqual(agents.vscode.enabled, true);
  });
});
