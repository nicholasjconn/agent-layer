# Backlog

Note: This is an agent-layer memory file. It is primarily for agent use.

## Features and tasks (not scheduled)

<!-- ENTRIES START -->

- Backlog 2026-01-25 8b9c2d1: Define migration strategy for renamed/deleted template files
    Priority: High. Area: lifecycle management
    Description: Define how to handle template files that are renamed or deleted in future versions so they do not remain as stale orphans in user repos.
    Acceptance criteria: A clear design/decision is documented for how to detect and clean up obsolete template files during upgrades.
    Notes: Currently, al init only adds/updates; it does not remove files that vanished from the binary.

- Backlog 2026-01-25 9e3f1a2: Improve CLI output readability and formatting
    Priority: Medium. Area: developer experience
    Description: Enhance the human readability of CLI outputs (wizard, init, doctor, etc.) by adding colors, improved formatting, and better use of whitespace.
    Acceptance criteria: CLI commands consistently use semantic coloring and spacing to make reports and interactive prompts easier to interpret.
    Notes: Focus on making high-impact messages (errors, warnings, successes) visually distinct.

- Backlog 2026-01-25 a1b2c3d: Add interaction monitoring for prompt self-improvement
    Priority: Low. Area: agent intelligence
    Description: Add interaction monitoring to agent system instructions to self-improve all prompts, rules, and workflows based on usage patterns.
    Acceptance criteria: Monitoring captures interaction patterns and produces actionable suggestions for prompt improvements.
    Notes: Requires careful design to avoid privacy concerns and ensure suggestions are high-quality.

- Backlog 2026-01-25 b2c3d4e: Enable safe auto-approval for workflow slash commands
    Priority: Medium. Area: workflow automation
    Description: Enable safe auto-approval for slash-command workflows invoked through the workflow system.
    Acceptance criteria: Workflows can run with minimal human intervention where the operation is deemed safe.
    Notes: Requires clear safety criteria and audit trail for auto-approved actions.

- Backlog 2026-01-25 c3d4e5f: Auto-merge client-side edits back to agent-layer sources
    Priority: Medium. Area: config synchronization
    Description: Auto-merge client-side approvals or MCP server edits back into agent-layer sources.
    Acceptance criteria: Changes made in client configs are detected and merged back to source files.
    Notes: Needs conflict resolution strategy and user confirmation for ambiguous merges.

- Backlog 2026-01-25 d4e5f6a: Add task queueing system for chained operations
    Priority: Low. Area: workflow automation
    Description: Add a queueing system to chain tasks without interrupting the current task.
    Acceptance criteria: Users can queue multiple tasks that execute sequentially without manual intervention between them.
    Notes: Consider how to handle failures mid-queue and queue persistence across sessions.

- Backlog 2026-01-25 e5f6a7b: Add slash-command ordering guide
    Priority: Low. Area: documentation
    Description: Add a simple flowchart or rules-based guide for slash-command ordering.
    Acceptance criteria: Documentation clearly explains recommended slash-command sequences for common workflows.
    Notes: Should cover common scenarios like feature development, bug fixing, and code review.

- Backlog 2026-01-25 f6a7b8c: Build multi-agent chat tool
    Priority: Low. Area: agent collaboration
    Description: Build a Ralph Wiggum-like tool where different agents can chat with each other.
    Acceptance criteria: Agents can exchange messages and collaborate on tasks through a shared interface.
    Notes: Experimental feature; needs clear use cases and safety boundaries.

- Backlog 2026-01-25 a7b8c9d: Build unified documentation repository with MCP access
    Priority: Low. Area: knowledge management
    Description: Build a unified documentation repository with Model Context Protocol tool access for shared notes.
    Acceptance criteria: A central repository exists that agents can read/write through MCP tools.
    Notes: Consider access control, versioning, and conflict resolution.

- Backlog 2026-01-25 b8c9d0e: Add indexed chat history for searchable context
    Priority: Low. Area: knowledge management
    Description: Add indexed chat history in the unified documentation repository for searchable context.
    Acceptance criteria: Past conversations are indexed and searchable to provide relevant context to agents.
    Notes: Depends on unified documentation repository being implemented first.

- Backlog 2026-01-25 c9d0e1f: Persist conversation history in model-specific folders
    Priority: Low. Area: session management
    Description: Persist conversation history in model-specific local folders (e.g., `.agent-layer/gemini/`, `.agent-layer/openai/`).
    Acceptance criteria: Conversation history is saved locally per model and can be restored across sessions.
    Notes: Consider storage format, retention policy, and privacy implications.

- Backlog 2026-01-25 d0e1f2a: Implement full access mode for all agents
    Priority: Low. Area: agent permissions
    Description: Implement "full access" mode for all agents with security warnings (similar to Codex full-auto).
    Acceptance criteria: Users can opt-in to full access mode with clear warnings about security implications.
    Notes: Requires prominent warnings and easy way to revert to restricted mode.
