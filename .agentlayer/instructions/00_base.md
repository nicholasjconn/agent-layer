# Instructions (Ultimate Source of Truth)

## Guiding Principles
1. **Fail Loudly & Quickly:** Never make assumptions to keep the code running. If an error occurs, the system must fail immediately and audibly (console error, exception, etc.) rather than degrade silently. Silent failures are worse than no tests.
2. **Single Source of Truth:** Every piece of data (environment variables, database state, configuration, derived metrics) must have a single, canonical source. Never maintain separate mutable state; always derive it from the source.
3. **No Logic in Frontend:** The frontend is strictly a presentation layer. It must **never** compute metrics, logic, or business rules. It only displays pre-computed data from the backend.
4. **No Fallbacks:** Do not implement fallback logic or default values that mask missing data, environment variables, or constants. If required input is missing or malformed, the system must fail rather than guessing.
5. **Strict UTC Only:** All internal time representations, storage, and API transport must be in UTC. No timezone info leaks into the system.

---

## Critical Protocol
1. **Mandatory Q&A:** If the user asks a direct question, you MUST answer it EXPLICITLY in your response text. If the prompt also requires code generation, answer the questions FIRST before discussing the code.
2. **Clarify Ambiguity:** If a decision is unclear or the prompt is ambiguous, pause and ask the user for clarification before generating code.

## Code Quality & Philosophy
1. **Adhere to Best Practices:** Follow industry standards for the language/framework in use. If a request violates these, specifically warn the user and ask for confirmation.
2. **Prioritize Clarity:** Write well-abstracted code where appropriate, but prioritize clarity, readability, and extensibility over clever abstraction.
3. **No Band-Aids:** Never apply quick fixes. Address the root cause, even if it is a major undertaking.
    * *Constraint:* If the fix requires significantly refactoring many files, ask for confirmation first.
4. **Search for Excellence:** Always search for the best solution, not just the easiest one.
5. **Strict Scope:** The agent should only make changes that are directly requested. Keep solutions simple and focused.
6. Never reduce the code coverage percentage in order to pass the tests.
7. **Packages:** Always use the command line to find the latest versions of packages to use. Do not rely on your training or memory to pick a package version.
8. **Strict Typing:** All Python code must use type hints (`typing`), and all TypeScript/JS must use strict types. Function signatures must have docstrings describing arguments and return values. Do not rely on implicit typing.

## Workflow & Safety
1. **Code Verification:** ALWAYS read and understand relevant files before proposing edits. Do not speculate about code you have not inspected.
2. **Context Economy:** When searching for files or context, strictly limit your search to the specific service or directory relevant to the request. Do not scan the entire monorepo unless explicitly necessary.
3. **Git Safety:** Never commit changes to git. Explicitly ask the user to commit changes for you.
4. **Temporary Files:** Generate all temporary/debug files in the `./.agentlayer/tmp` directory.
5. **Schema Safety:** Never modify the database schema via raw SQL or direct tool access. Always generate a proper migration file (e.g., Alembic, Prisma, Django) and ask the user to apply it.
6. **Debugging Strategy:** Debugging follows a strict Red-Green loop: reproduce the bug with a persistent test case, then fix it. **Avoid one-off scripts** unless identifying the issue with a test case is impossible. If a one-off script is absolutely required, it must be generated in `./.agentlayer/tmp` and **deleted** immediately after use.
7. **Definition of Done:** A task is not complete until relevant Markdown files are updated and the code is fully documented.
8. **Environment Variables:** Never modify the .env file. Only ever modify the .env.example file. If a change needs to be made to the .env file, ask the user to make it.

