import { isPlainObject } from "./utils.mjs";

/**
 * Parse Codex config toml sections into a map of name -> block lines.
 * @param {string} content
 * @returns {{ header: string[], sections: Map<string, string[]> }}
 */
export function parseCodexConfigSections(content) {
  const lines = content.split(/\r?\n/);
  const sections = new Map();
  const header = [];
  let currentName = null;
  let currentLines = [];

  const flush = () => {
    if (currentName) {
      sections.set(currentName, currentLines.slice());
    } else if (currentLines.length) {
      header.push(...currentLines);
    }
    currentLines = [];
  };

  for (const line of lines) {
    const match = line.match(/^\[mcp_servers\.([^\]]+)\]$/);
    if (match) {
      flush();
      const rawName = match[1];
      let name = rawName;
      if (rawName.startsWith('"') && rawName.endsWith('"')) {
        try {
          name = JSON.parse(rawName);
        } catch {
          name = rawName;
        }
      } else if (rawName.startsWith("'") && rawName.endsWith("'")) {
        name = rawName.slice(1, -1);
      }
      currentName = name;
      currentLines.push(line);
    } else {
      currentLines.push(line);
    }
  }
  flush();

  return { header, sections };
}

/**
 * Parse a Codex config section into a server definition.
 * @param {string[]} lines
 * @returns {{ server: { transport: "stdio", command: string, args: string[], envVars: string[], envVarsKnown: boolean } | { transport: "http", url: string, bearerTokenEnvVar: string|null }, reason?: string }}
 */
export function parseCodexServerSection(lines) {
  let command = null;
  let args = [];
  let envVars = [];
  let envVarsKnown = false;
  let url = null;
  let bearerTokenEnvVar = null;

  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith("#")) continue;
    if (trimmed.startsWith("[") && trimmed.endsWith("]")) {
      if (!trimmed.startsWith("[mcp_servers.")) break;
      continue;
    }
    const cmdMatch = trimmed.match(/^command\s*=\s*(.+)$/);
    if (cmdMatch) {
      const raw = cmdMatch[1].trim();
      try {
        command = JSON.parse(raw);
      } catch {
        return { reason: "command is not a JSON string" };
      }
      continue;
    }
    const urlMatch = trimmed.match(/^url\s*=\s*(.+)$/);
    if (urlMatch) {
      const raw = urlMatch[1].trim();
      try {
        url = JSON.parse(raw);
      } catch {
        return { reason: "url is not a JSON string" };
      }
      continue;
    }
    const bearerMatch = trimmed.match(/^bearer_token_env_var\s*=\s*(.+)$/);
    if (bearerMatch) {
      const raw = bearerMatch[1].trim();
      try {
        bearerTokenEnvVar = JSON.parse(raw);
      } catch {
        return { reason: "bearer_token_env_var is not a JSON string" };
      }
      continue;
    }
    const argsMatch = trimmed.match(/^args\s*=\s*(.+)$/);
    if (argsMatch) {
      const raw = argsMatch[1].trim();
      try {
        const parsed = JSON.parse(raw);
        if (
          !Array.isArray(parsed) ||
          parsed.some((v) => typeof v !== "string")
        ) {
          return { reason: "args is not a string array" };
        }
        args = parsed;
      } catch {
        return { reason: "args is not valid JSON" };
      }
      continue;
    }
    const envMatch = trimmed.match(/^env\s*=\s*\{(.+)\}$/);
    if (envMatch) {
      envVarsKnown = true;
      const body = envMatch[1];
      const matches = body.matchAll(/([A-Za-z_][A-Za-z0-9_]*)\s*=/g);
      envVars = Array.from(matches).map((m) => m[1]);
    }
  }

  if (url) {
    if (command) {
      return { reason: "command is not allowed for HTTP servers" };
    }
    return {
      server: {
        transport: "http",
        url,
        bearerTokenEnvVar,
      },
    };
  }

  if (bearerTokenEnvVar) {
    return { reason: "bearer_token_env_var requires url" };
  }

  if (!command) {
    return { reason: "missing command" };
  }
  return {
    server: { transport: "stdio", command, args, envVars, envVarsKnown },
  };
}

/**
 * Determine whether an MCP server entry is enabled.
 * @param {unknown} entry
 * @returns {boolean}
 */
export function isEnabledServer(entry) {
  if (!isPlainObject(entry)) return false;
  if (entry.enabled === undefined) return true;
  return entry.enabled === true;
}
