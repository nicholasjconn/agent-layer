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

- Issue 2026-01-22 71af41d: README missing overwrite prompt and force flag guidance
    Priority: Medium. Area: documentation.
    Description: README does not describe `--overwrite` prompting behavior or the `--force` flag, while other documentation does.
    Next step: Add a short install note in README covering prompts and non-interactive overwrites.
    Notes: Keep language aligned with Phase 7 install flow messaging.

- Issue 2026-01-22 71af41d: Installer script lacks `--force` passthrough
    Priority: Low. Area: installation.
    Description: `agent-layer-install.sh` does not accept or pass through `--force`, limiting non-interactive overwrites when re-running the installer.
    Next step: Add a `--force` flag to the script and forward it to `./al install`.
    Notes: Ensure usage text remains concise.

- Issue 2026-01-22 71af41d: Linux launcher always opens a terminal window
    Priority: Low. Area: launchers.
    Description: `open-vscode.desktop` sets `Terminal=true`, opening a terminal even on successful launches.
    Next step: Consider a launcher that only opens a terminal on failure or uses a non-terminal notification for errors.
    Notes: Keep `CODEX_HOME` behavior consistent across platforms.

- Issue 2026-01-22 71af41d: Install documentation fragmented across files
    Priority: Low. Area: documentation.
    Description: Install guidance is split between README, developer documentation, and the installer script, which increases the risk of inconsistent messaging.
    Next step: Consolidate or cross-link install flag guidance so changes remain synchronized.
    Notes: Focus on user-facing install flags and prompt behavior.

- Issue 2026-01-21 b1c2d3: Wizard Approval Mode should show descriptions
    Priority: Medium. Area: wizard.
    Description: Approval Mode selection in the wizard does not display the description of the mode, making it hard for users to remember details.
    Next step: Update the wizard's Approval Mode selection step to display the description for each option.

- Issue 2026-01-21 5af6278: Wizard run tests are oversized
    Priority: Low. Area: wizard tests.
    Description: `internal/wizard/run_test.go` is over 1200 lines, which makes review and maintenance harder.
    Next step: Split the test file by scenario into smaller files with shared helpers.

- Issue 2026-01-18 e4f5g6: Memory file template structure investigation
    Priority: Medium. Area: templates.
    Description: Should templates in .agent-layer only contain headers, and how should generated content be handled when overwriting?
    Next step: Review existing template synchronization logic and define the intended behavior for content preservation.

- Issue 2026-01-18 l8m9n0: Limit exposed commands for GitHub MCP
    Priority: Medium. Area: mcp configuration.
    Description: The GitHub MCP server exposes many tools. We should explicitly list only the necessary commands in the configuration to reduce noise and potential security risks.
    Next step: Research useful GitHub MCP commands and configure `args` or `commands` whitelist in the default config template.

- Issue 2026-01-21 f6g7h8: Centralized user-facing strings file
    Priority: Low. Area: code organization.
    Description: All user-facing language (messages, prompts, errors, help text) is scattered across the codebase. Centralizing it would make updates easier and ensure consistent tone.
    Next step: Audit existing user-facing strings and design a centralized strings module.

- Issue 2026-01-21 r5s6t7: Audit instructions and slash commands for duplicate content
    Priority: Low. Area: maintainability.
    Description: Instruction templates and slash command definitions contain duplicate instructions (e.g., memory formatting rules appear multiple times). This creates maintenance burden and inconsistency risk.
    Next step: Identify duplicate content across instruction and slash command files; consolidate into canonical locations.
