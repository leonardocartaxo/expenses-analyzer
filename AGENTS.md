# Expenses Analyzer — Agent Guide

Governing principles: [`.specify/memory/constitution.md`](.specify/memory/constitution.md).

## Project shape

- Full rewrite (greenfield). Do not revive the deleted Go tree unless a current spec says so.
- Monorepo: **backend** + **frontend** (+ optional shared contracts). Stack is chosen once and recorded in plan / stack notes—do not invent a stack in code before that decision.
- Spec-driven workflow via Spec Kit skills under `.cursor/skills/` (`speckit-*`).
- Dev environment: open the repo in the **Dev Container** (`.devcontainer/`). Works with **WebStorm** (`customizations.jetbrains`) and **VS Code / Cursor** (`customizations.vscode`). Runtime Compose (DB/app services) comes after the stack is locked.

## Default loop

1. `speckit-constitution` — process/governance only  
2. `speckit-specify` → optional `speckit-clarify`  
3. `speckit-plan` (stack + package layout)  
4. `speckit-tasks` → optional `speckit-analyze`  
5. `speckit-implement` — run the root verify harness after coherent batches  
6. `speckit-converge` — append remaining tasks; do not silently change specs  

## Hard rules

- Specs own behavior; plans own tech; harness owns pass/fail.
- No secrets in git. No cross-package boundary violations (UI in backend, DB in frontend).
- Prefer the smallest change that satisfies the approved spec.
