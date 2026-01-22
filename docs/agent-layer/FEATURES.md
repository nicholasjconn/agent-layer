# Features

Purpose: Backlog of deferred user-visible feature requests (not yet scheduled into the roadmap).

Notes for updates:
- Add an entry only when it is a new user-visible capability.
- Keep each entry 3 to 5 lines (the first line plus 2 to 4 indented lines).
- Lines 2 to 5 must be indented by 4 spaces so they stay associated with the entry.
- Prevent duplicates by searching and merging.
- When a feature is scheduled into a roadmap phase, move it into `ROADMAP.md` and remove it from this file.
- When a feature is implemented, ensure it is removed from this file.

Entry format:
- Feature YYYY-MM-DD abcdef: Short title
    Priority: Critical, High, Medium, or Low. Area: <area>
    Description: <what the user should be able to do>
    Acceptance criteria: <clear condition to consider it done>
    Notes: <optional dependencies or constraints>

## Backlog (not scheduled)

<!-- ENTRIES START -->

- Feature 2026-01-21 7b2c9d: Explicit `al install --upgrade` command
    Priority: Medium. Area: installation
    Description: Add an `--upgrade` flag to the `install` command to explicitly download and replace the existing `al` binary with the latest version.
    Acceptance criteria: Running `./al install --upgrade` downloads the latest binary and refreshes the installation.

- Feature 2026-01-21 4e8f1a: `al install --version vX.X.X` command
    Priority: Medium. Area: installation
    Description: Allow users to specify a specific version of Agent Layer to install or switch to using a `--version` flag.
    Acceptance criteria: Running `./al install --version v1.2.3` downloads and installs the specified version of the binary.
