# Issues

Purpose: Deferred defects, maintainability refactors, technical debt, risks, and engineering concerns.

Notes for updates:
- Add an entry only when you are not fixing it now.
- Keep each entry 3 to 5 lines (the first line plus 2 to 4 indented lines).
- Lines 2 to 5 must be indented by 4 spaces so they stay associated with the entry.
- Prevent duplicates by searching and merging.
- Remove entries when fixed.

Entry format:
- Issue YYYY-MM-DD abcdef: Short title
    Priority: Critical, High, Medium, or Low. Area: <area>
    Description: <observed problem or risk>
    Next step: <smallest concrete next action>
    Notes: <optional dependencies or constraints>

## Open issues

<!-- ENTRIES START -->

- Issue 2026-01-24 a1b2c3: VS Code slow first launch in agent-layer folder
    Priority: Low. Area: developer experience.
    Description: Launching VS Code in the agent-layer folder takes a very long time on first use, likely due to extension initialization, indexing, or MCP server startup.
    Next step: Profile VS Code startup to identify the bottleneck (extensions, language servers, MCP servers, or workspace indexing).

- Issue 2026-01-24 9a8b7c: Refactor global function patching to DI
    Priority: High. Area: architecture.
    Description: Systemic pattern of patching global variables (e.g. `var lookPath = exec.LookPath`) prevents parallel testing and invites race conditions.
    Next step: Refactor `internal/sync` to use interface-based dependency injection and enable `t.Parallel()`.
