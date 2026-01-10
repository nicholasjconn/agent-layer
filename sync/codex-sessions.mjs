import fs from "node:fs";
import path from "node:path";
import { fileExists, isPlainObject, readUtf8 } from "./utils.mjs";

/**
 * Load Codex session JSONL files under a sessions directory.
 * @param {string} sessionsDir
 * @returns {string[]}
 */
function listCodexSessionFiles(sessionsDir) {
  if (!fileExists(sessionsDir)) return [];

  /** @type {string[]} */
  const files = [];
  const stack = [sessionsDir];
  while (stack.length) {
    const dir = stack.pop();
    if (!dir) continue;
    const entries = fs.readdirSync(dir, { withFileTypes: true });
    for (const entry of entries) {
      const full = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        stack.push(full);
      } else if (entry.isFile() && entry.name.endsWith(".jsonl")) {
        files.push(full);
      }
    }
  }
  return files.sort((a, b) => fs.statSync(b).mtimeMs - fs.statSync(a).mtimeMs);
}

/**
 * Check whether a value is a string array.
 * @param {unknown} value
 * @returns {value is string[]}
 */
function isStringArray(value) {
  return Array.isArray(value) && value.every((v) => typeof v === "string");
}

/**
 * Try to extract an argv array from a JSONL record.
 * @param {unknown} record
 * @returns {string[] | null}
 */
function extractArgv(record) {
  if (!isPlainObject(record)) return null;
  const directCandidates = [record.command, record.argv];
  for (const cand of directCandidates) {
    if (isStringArray(cand)) return cand;
  }

  const candidates = [
    record.msg,
    record.tool,
    record.params,
    record.request,
    record.input,
    record.command,
  ];
  for (const cand of candidates) {
    if (!isPlainObject(cand)) continue;
    const nestedCandidates = [cand.command, cand.argv];
    for (const nested of nestedCandidates) {
      if (isStringArray(nested)) return nested;
    }
    const deeperCandidates = [
      cand.params,
      cand.request,
      cand.input,
      cand.command,
    ];
    for (const deeper of deeperCandidates) {
      if (!isPlainObject(deeper)) continue;
      const deepArgv = deeper.command ?? deeper.argv;
      if (isStringArray(deepArgv)) return deepArgv;
    }
  }
  return null;
}

/**
 * Extract an event type from a record.
 * @param {unknown} record
 * @returns {string}
 */
function extractEventType(record) {
  if (!isPlainObject(record)) return "";
  if (typeof record.type === "string") return record.type;
  const msg = record.msg;
  if (isPlainObject(msg) && typeof msg.type === "string") return msg.type;
  return "";
}

/**
 * Check whether an event type is an approval request.
 * @param {string} eventType
 * @returns {boolean}
 */
function isApprovalRequestType(eventType) {
  return eventType.endsWith("_approval_request");
}

/**
 * Check whether a record looks like an approval event.
 * @param {unknown} record
 * @returns {boolean}
 */
function isApprovalRecord(record) {
  if (!isPlainObject(record)) return false;
  const eventType = extractEventType(record);
  if (eventType && isApprovalRequestType(eventType)) return true;
  const text = [
    record.type,
    record.event,
    record.name,
    record.action,
    record.kind,
  ]
    .filter((v) => typeof v === "string")
    .join(" ")
    .toLowerCase();

  const hasHint =
    text.includes("approval") ||
    text.includes("approve") ||
    text.includes("permission");
  const decision = record.decision ?? record.status ?? record.result;
  const approved = record.approved ?? record.allowed ?? record.allow;

  if (approved === true) return true;
  if (typeof decision === "string") {
    const normalized = decision.toLowerCase();
    if (
      normalized === "allow" ||
      normalized === "approved" ||
      normalized === "yes"
    )
      return true;
  }
  return hasHint;
}

/**
 * Scan Codex session logs for approvals not in policy.
 * @param {string} workingRoot
 * @param {Set<string>} policySet
 * @param {number} maxFiles
 * @returns {{ items: import("./divergence.mjs").ApprovalItem[], notes: string[] }}
 */
export function scanCodexSessionApprovals(
  workingRoot,
  policySet,
  maxFiles = 1,
) {
  const notes = [];
  if (!workingRoot) return { items: [], notes };
  const sessionsDir = path.join(workingRoot, ".codex", "sessions");

  let files = [];
  try {
    files = listCodexSessionFiles(sessionsDir);
  } catch {
    notes.push(`Could not read Codex sessions under ${sessionsDir}`);
    return { items: [], notes };
  }

  if (!files.length) return { items: [], notes };
  if (maxFiles > 0) files = files.slice(0, maxFiles);

  /** @type {import("./divergence.mjs").ApprovalItem[]} */
  const items = [];
  const seen = new Set();
  for (const filePath of files) {
    let lines = [];
    try {
      lines = readUtf8(filePath).split(/\r?\n/);
    } catch {
      notes.push(`Could not read ${filePath}`);
      continue;
    }
    for (let idx = 0; idx < lines.length; idx++) {
      const line = lines[idx].trim();
      if (!line) continue;
      let record;
      try {
        record = JSON.parse(line);
      } catch {
        continue;
      }
      const eventType = extractEventType(record);
      const isRequest = eventType && isApprovalRequestType(eventType);
      if (!isRequest && !isApprovalRecord(record)) continue;

      const argv = extractArgv(record);
      if (!argv) continue;
      const prefix = argv.join(" ");
      if (policySet.has(prefix)) continue;
      if (seen.has(prefix)) continue;
      seen.add(prefix);
      items.push({
        kind: "approval",
        source: "codex-session",
        filePath: `${filePath}:${idx + 1}`,
        raw: line,
        prefix,
        argv,
        parseable: true,
        ...(isRequest ? { reason: `eventType=${eventType}` } : {}),
      });
    }
  }

  return { items, notes };
}
