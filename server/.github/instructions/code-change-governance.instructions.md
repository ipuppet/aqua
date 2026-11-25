---
description: "Use when modifying code or editing project files in this repository. Enforces auditable file edits and mandatory pre-read of project standards."
name: "Code Change Governance"
applyTo: "**"
---
# Code Change Governance

Apply these rules to all code modification tasks in this repository.

## 1) Traceable Editing Only (Hard Rule)

- Use built-in file tools for reading and writing project files.
- Do not modify files through terminal-based rewrite/edit mechanisms (for example, `sed`/`awk`/`perl` replacements, shell redirection overwrite, generated patch scripts).
- Terminal usage is allowed for non-editing operations such as build, test, run, lint, and diagnostics.
- Keep changes reviewable and reversible at file level.

## 2) Read Standards Before Editing (Hard Rule)

Before making any code change, first check project documentation in this order:

1. `README.md`
2. `docs/standards/ENGINEERING.md`
3. `docs/standards/ARCHITECTURE.md`
4. Related service docs under `docs/services/`

If no applicable rule is found, follow widely accepted best practices for the language and framework in use.

## 3) Conflict Resolution

- If an existing repository rule conflicts with a task request, surface the conflict clearly and ask for explicit override instructions.
- For refactoring compatibility behavior, follow the separate on-demand refactor instruction.
