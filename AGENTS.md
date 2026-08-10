# Expenses Analyzer — Agent Guide

Governing principles: [`.specify/memory/constitution.md`](.specify/memory/constitution.md) (process).  
Global product / business rules: [`docs/PRODUCT.md`](docs/PRODUCT.md) — do not invent product behavior that contradicts this file; do not duplicate it under `.specify/memory/`.

## Project shape

- Full rewrite (greenfield). Do not revive the deleted Go tree unless a current spec says so.
- Monorepo: **backend** + **frontend** (+ optional shared contracts). Stack is chosen once and recorded in plan / stack notes—do not invent a stack in code before that decision.
- Spec-driven workflow via Spec Kit skills under `.cursor/skills/` (`speckit-*`). Feature specs under `specs/` must align with `docs/PRODUCT.md`.
- Dev environment: open the repo in the **Dev Container** (`.devcontainer/`). Works with **WebStorm** (`customizations.jetbrains`) and **VS Code / Cursor** (`customizations.vscode`). Base image is **Node 24 bookworm-slim** (glibc) because JetBrains Remote Dev does not run on Alpine/musl. Runtime Compose (DB/app services) comes after the stack is locked.

## Default loop

1. `speckit-constitution` — process/governance only  
2. Refine `docs/PRODUCT.md` when product intent changes (global rules)  
3. `speckit-specify` → optional `speckit-clarify` — delivery slices that follow PRODUCT.md  
4. `speckit-plan` (stack + package layout)  
5. `speckit-tasks` → optional `speckit-analyze`  
6. `speckit-implement` — run the root verify harness after coherent batches  
7. `speckit-converge` — append remaining tasks; do not silently change specs  

## Hard rules

- `docs/PRODUCT.md` owns global product rules; feature specs own slice behavior; plans own tech; harness owns pass/fail.
- No secrets in git. No cross-package boundary violations (UI in backend, DB in frontend).
- Prefer the smallest change that satisfies the approved spec.
