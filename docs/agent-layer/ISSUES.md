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

- Issue 2026-01-19 ceddb83: `.agent-layer/.env` overrides shell environment variables
    Priority: Medium. Area: environment handling.
    Description: When launching via `./al`, values from `.agent-layer/.env` override existing shell environment variables, and empty template keys can shadow valid tokens.
    Next step: Decide precedence and update environment merge logic or templates to avoid overriding with empty values; document the chosen behavior.

- Issue 2026-01-18 e5f6g7: Slash commands not output for antigravity
    Priority: Medium. Area: antigravity support.
    Description: Slash commands are not being output when antigravity mode is enabled.
    Next step: Investigate where slash commands are generated and ensure antigravity support is included.

- Issue 2026-01-18 h8i9j0: DECISIONS.md grows too large and consumes excessive tokens
    Priority: Medium. Area: project memory.
    Description: The decisions log grows unbounded as entries accumulate, eventually consuming too many tokens when agents read it for context.
    Next step: Consider archiving old decisions, summarizing completed phases, or splitting into a compact summary plus detailed archive.

- Issue 2026-01-18 b1c2d3: Memory file path convention investigation
    Priority: Low. Area: project memory.
    Description: Should memory files use full relative paths or just filenames in 01_memory.md and slash commands?
    Next step: Audit current usage and establish a single convention for referring to memory files.

- Issue 2026-01-18 e4f5g6: Memory file template structure investigation
    Priority: Medium. Area: templates.
    Description: Should templates in .agent-layer only contain headers, and how should generated content be handled when overwriting?
    Next step: Review existing template synchronization logic and define the intended behavior for content preservation.

- Issue 2026-01-18 a7b8c9: Boost coverage slash command too conservative
    Priority: High. Area: slash commands.
    Description: The boost-coverage command only picks one file at a time and stops too early. It should iterate until coverage targets are met, even if it requires many tests.
    Next step: Refactor the boost-coverage logic to support continuous iteration and multiple file targets in a single run.

- Issue 2026-01-18 f1g2h3: Codex initial reasoning effort should be high
    Priority: Medium. Area: configuration templates.
    Description: The initial example configuration for Codex sets reasoning_effort to "xhigh", which can be unnecessarily expensive or slow for a default. It should be "high".
    Next step: Update `internal/templates/config.toml` and `README.md` to use "high" instead of "xhigh".

- Issue 2026-01-18 i4j5k6: Remove MY_TOKEN from default configuration template
    Priority: Low. Area: configuration templates.
    Description: The `MY_TOKEN` placeholder in the default configuration template is unnecessary and should be removed to keep the default config clean. It can remain in the README as an example of environment variable usage.
    Next step: Remove `MY_TOKEN` from `internal/templates/config.toml`.

- Issue 2026-01-18 k7l8m9: Set MCP servers to enabled by default
    Priority: Medium. Area: configuration templates.
    Description: Current MCP server examples in the default configuration are disabled by default. They should be enabled by default to provide a better out-of-the-box experience when tokens are provided.
    Next step: Update `internal/templates/config.toml` to set `enabled = true` for default MCP servers.

- Issue 2026-01-18 l8m9n0: Limit exposed commands for GitHub MCP
    Priority: Medium. Area: mcp configuration.
    Description: The GitHub MCP server exposes many tools. We should explicitly list only the necessary commands in the configuration to reduce noise and potential security risks.
    Next step: Research useful GitHub MCP commands and configure `args` or `commands` whitelist in the default config template.
