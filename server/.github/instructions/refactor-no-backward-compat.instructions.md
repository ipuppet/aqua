---
description: "Use when performing refactoring, renaming, architecture cleanup, or repository-wide structural changes. Enforces no backward-compatibility layer unless explicitly requested."
name: "Refactor Without Backward Compatibility"
---
# Refactor Without Backward Compatibility

Apply these rules when the task is a refactor.

## 1) No Compatibility Layer by Default (Hard Rule)

- Do not add backward-compatibility code unless the user explicitly asks to preserve compatibility.
- Prefer one-pass, repo-wide cleanup over temporary adapters, aliases, bridges, or duplicate pathways.
- Remove obsolete paths during the same refactor when safe to do so.

## 2) Conflict Resolution

- If compatibility retention is required, it must be explicitly stated in the task.
- If an existing repository rule conflicts with a refactor request, surface the conflict clearly and ask for explicit override instructions.
