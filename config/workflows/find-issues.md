---
description: Proactive, repo-adaptive code quality audit to surface over-engineering, bad practices, and “band-aids”, aligned with README and the project memory files. Report-first: no automatic issue logging unless the user explicitly opts in.
---

# Proactive code quality audit (report-first, memory-aware)

## Intent
Perform a structured audit to find:
- over-engineering / unnecessary complexity
- poor abstractions / architecture drift
- brittle fallbacks and defaults that mask failures
- “band-aids” / quick-and-dirty fixes that became permanent
- DRY violations and extensibility pain

This workflow is **report-first**:
- **Do not modify code**.
- **Do not write to** `docs/ISSUES.md`, `docs/FEATURES.md`, `docs/ROADMAP.md`, or `docs/DECISIONS.md` unless the user explicitly requests it.

---

## Project memory files (authoritative)
- `docs/ISSUES.md` — deferred defects, maintainability refactors, technical debt, risks.
- `docs/FEATURES.md` — backlog of deferred user feature requests (not yet scheduled into the roadmap).
- `docs/ROADMAP.md` — numbered phases; guides architecture and sequencing.
- `docs/DECISIONS.md` — rolling log of important decisions (brief).

If any are missing, create them from `config/templates/docs/<NAME>.md` (preserve headings and markers), but only if `write_mode=apply`.

---

## Inputs (optional)
If the user provides arguments after the command, interpret them as:

- `scope=uncommitted|last_commit|range|paths|repo` (default: `uncommitted`)
- `range=<git-range>` (only if `scope=range`, e.g. `HEAD~5..HEAD`)
- `paths=<comma-separated paths>` (only if `scope=paths`)
- `depth=narrow|deep|broad|full` (default: `full`)
  - `narrow`: only the narrow set (changed files)
  - `deep`: narrow + one layer outward (callers/callees/adjacent modules)
  - `broad`: repo-wide scan for smell patterns and hotspot spot-checking
  - `full`: narrow → deep → broad
- `risk=low|medium|high` (default: `medium`) — influences prioritization and recommended verification
- `report_path=quality_audit_report.md` (default: `quality_audit_report.md`)
- `max_findings=20` (default: `20`)
- `include_snippets=none|short` (default: `short`)
- `write_mode=report_only|propose|apply` (default: `report_only`)
  - `report_only`: report findings only
  - `propose`: include “ready-to-paste” entries for memory files, but do not write them
  - `apply`: update memory files using the project rules (requires explicit user request)
- `entry_id=auto|<token>` (default: `auto`) — used in proposed Issue/Feature/Decision entries (e.g., short SHA)

---

## Roles and handoffs (multi-agent)
1. **Standards Extractor**: read `docs/ROADMAP.md` then `docs/DECISIONS.md`, then scan `docs/FEATURES.md` and `docs/ISSUES.md`, then read `README.md` to derive the “audit lens”.
2. **Diff Auditor**: narrow audit on changed files per `scope`.
3. **Deep-Dive Investigator**: expand one layer outward from the highest-impact narrow findings.
4. **Broad Scanner**: lightweight smell scanning + hotspot spot-checking.
5. **Synthesizer**: deduplicate, prioritize, and map findings to roadmap/decisions/standards.
6. **Reporter**: write `report_path` and present a review summary to the user.

If only one agent is available, execute phases in this order with explicit headings.

---

## Guardrails
- Avoid speculation: each finding must include **evidence** (file + symbol/section + reasoning; optional short snippet).
- Avoid nitpicks: focus on issues that impact correctness, delivery speed, operability, or extensibility.
- No broad rewrites: this is an audit, not a refactor plan.
- Use the repo’s terminology. Avoid abbreviations in memory entries.

---

# Phase 0 — Preflight (Diff Auditor)

1. Confirm repo state:
   - `git status --porcelain`

2. Determine the narrow file set:
- `scope=uncommitted` (default):
  - staged: `git diff --name-only --staged`
  - unstaged: `git diff --name-only`
- `scope=last_commit`:
  - `git show --name-only --pretty="" HEAD`
- `scope=range`:
  - `git diff --name-only <range>`
- `scope=paths`:
  - use `paths=...`
- `scope=repo`:
  - narrow set is empty; proceed directly to broad scan

If `scope=uncommitted` yields no files, fall back to `scope=last_commit` and note it in the report.

3. Establish an entry identifier:
- if `entry_id=auto` and git is available: use `git rev-parse --short HEAD`
- otherwise: use `nogit`

**Deliverable**
- Narrow file list (or explicit “none”)
- Actual scope used (including fallback)
- Entry identifier token

---

# Phase 1 — Standards and roadmap lens (Standards Extractor)

Read in this order (when present):
1. `docs/ROADMAP.md`
2. `docs/DECISIONS.md`
3. `docs/FEATURES.md`
4. `docs/ISSUES.md`
5. `README.md`

Extract:
- architectural boundaries and layering rules
- dependency rules and patterns to avoid
- module ownership expectations
- planned future direction (roadmap phases)
- existing decisions that constrain changes

Produce an “audit lens” checklist:
- **Must align with** (roadmap/decisions)
- **Must avoid** (anti-patterns called out by repo standards)
- **Design priorities** (simplicity, explicitness, reliability, performance, etc.)

---

# Phase 2 — Narrow audit (Diff Auditor)

Trigger when `depth=narrow|deep|full`.

For each file in the narrow set:

## 2A) Over-engineering / complexity signals
Examples (language-agnostic):
- multiple abstraction layers for simple operations
- indirection that obscures behavior (factories/wrappers without clear payoff)
- configuration/flags controlling many unrelated behaviors
- “framework inside the codebase” patterns (custom DSL/IOC) without strong justification

## 2B) Bad practices / brittleness
- swallowed errors / overly broad exception handling
- side effects hidden inside helpers with unclear ownership
- hidden global state, implicit environment assumptions
- confusing naming, unclear invariants, missing validation at boundaries
- unclear separation of concerns vs the README’s layering model

## 2C) Fallbacks, defaults, and band-aids
- magic defaults that mask missing configuration
- “temporary” workarounds that became permanent
- retries/timeouts/backoffs missing or inconsistent
- suppression markers (lint/type ignores) concentrated in one area

## 2D) Record evidence-backed findings
Each finding should include:
- **Title**
- **Location**: file + symbol/section
- **Evidence**: what you observed (optional short snippet)
- **Impact**: correctness | operability | maintainability | extensibility | performance | security
- **Severity**: High | Medium | Low
- **Recommendation**: concrete next step and smallest viable improvement
- **Alignment**: how it matches/conflicts with ROADMAP/DECISIONS/README

---

# Phase 3 — Deep dive (Deep-Dive Investigator)

Trigger when `depth=deep|full`.

1. Select up to 5 highest-impact narrow findings.
2. For each, go one layer outward:
   - callers/call-sites
   - direct dependencies
   - adjacent modules in the same feature area
   - tests that cover the area (if present)
3. Determine whether the issue is localized or systemic:
   - repeated patterns across modules
   - inconsistent abstractions
   - duplicated logic or inconsistent defaults

Refine recommendations into “next-step sized” actions.

---

# Phase 4 — Broad audit (Broad Scanner)

Trigger when `depth=broad|full`.

## 4A) Identify hotspots
Use one or more of:
- largest modules/files by line count
- high-churn areas (recent history) when practical
- core infrastructure directories named in README
- areas scheduled soon in ROADMAP (avoid spending effort on code slated for replacement; flag instead)

## 4B) Lightweight repo-wide smell scanning
Use tooling only if available; otherwise use simple searches (`rg` preferred, else `grep -R`).

Suggested patterns (adapt as appropriate):
- markers: `TODO|FIXME|HACK|WORKAROUND|TEMP|XXX`
- suppressed checks: `eslint-disable|@ts-ignore|nolint|noqa|type: ignore`
- suspicious defaults: `default|fallback|if.*None|null|undefined`
- broad catches: `except:|except Exception|catch \(.*\)`
- “temporary” flags: `feature flag|kill switch|temporary flag`

Summarize into themes, not raw match dumps.

## 4C) Spot-check for architectural drift
Pick a few representative areas and check:
- is the layering consistent with README?
- are boundaries stable with ROADMAP direction?
- do decisions in DECISIONS still hold, or are they being violated?

---

# Phase 5 — Synthesis and prioritization (Synthesizer)

1. Deduplicate findings across phases.
2. Prioritize using a consistent rubric:

**Severity rubric**
- **High**: likely correctness/security issue, production incident risk, or blocks ROADMAP.
- **Medium**: significant maintainability cost, likely bug source, extensibility pain.
- **Low**: clarity/polish; not urgent.

3. Select up to `max_findings` for the report.

---

# Phase 6 — Report and review (Reporter)

Create `report_path` with:

## Required sections
1. **Executive summary** (5–10 bullets)
2. **Standards & roadmap lens** (from Phase 1)
3. **Top findings** (prioritized list)
4. **Themes and systemic risks**
5. **Suggested next steps**
   - propose small coherent batches (≤ 3) rather than a long backlog

## Optional: ready-to-paste memory entries
If `write_mode=propose|apply`, include an “Entry candidates” section.

### Candidate entry rules (mandatory)
- Entries must be **3 to 5 lines**.
- No abbreviations. Prefer short sentences.
- Search existing memory files first and avoid near-duplicates.

Formats:
- Issues: `- Issue YYYY-MM-DD <entry_id>: <short title>`
- Features: `- Feature YYYY-MM-DD <entry_id>: <short title>` (user-visible only)
- Decisions: `- Decision YYYY-MM-DD <entry_id>: <short title>`

**Important:** If `write_mode=propose`, do not modify files—only provide candidates.

---

# Phase 7 — Optional apply step (only if explicitly requested)

Proceed only if `write_mode=apply`.

1. Ensure all memory files exist (create from templates if missing).
2. For each approved finding:
   - Add to `docs/ISSUES.md` if it is a defect, refactor, technical debt, reliability/security/performance risk, or test gap.
   - Add to `docs/FEATURES.md` only if it is a user-visible capability request.
   - Add to `docs/DECISIONS.md` only if it is a significant decision and the user wants it logged.
   - Consider `docs/ROADMAP.md` updates only if the user explicitly wants roadmap edits.
3. Deduplicate by merging into existing entries; keep entries compact.

---

## In-chat review step (mandatory)
Present the top findings (High + Medium) to the user and ask which items should be:
- logged to ISSUES, FEATURES, DECISIONS
- scheduled into ROADMAP
- turned into a fix plan
- deprioritized/ignored

---

## Definition of done
- `report_path` exists and is readable.
- Findings are evidence-backed, deduplicated, and aligned to README/ROADMAP/DECISIONS.
- No code changes made unless separately requested.
- Memory files are not modified unless the user explicitly requests `write_mode=apply`.
