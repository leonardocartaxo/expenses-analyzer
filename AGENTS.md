# Expenses Analyzer — Agent Guide

Governing principles: [`.specify/memory/constitution.md`](.specify/memory/constitution.md) (process).  
Global product / business rules: [`docs/PRODUCT.md`](docs/PRODUCT.md) — do not invent product behavior that contradicts this file; do not duplicate it under `.specify/memory/`.

## Project shape

- Full rewrite (greenfield). Do not revive the deleted Go tree unless a current spec says so.
- Monorepo: **backend** + **frontend** (+ optional shared contracts). Stack is chosen once and recorded in plan / stack notes—do not invent a stack in code before that decision.
- Spec-driven workflow via Spec Kit skills under `.cursor/skills/` (`speckit-*`). Feature specs under `specs/` must align with `docs/PRODUCT.md`.
- Dev environment: **default local** is host Node + pnpm with **Compose Postgres only**. Optional **Dev Container** (`.devcontainer/`) for WebStorm / VS Code / Cursor (Node 24 bookworm-slim). Optional **kind** (`pnpm local:up`). See `STACK.md`. AWS first deploy is `011-aws-deploy`; wake/sleep is `010`.
- **Cursor CLI**: post-create installs `agent`. In a new Dev Container shell, run **`agent login`** so the CLI is authenticated (see `README.md` → Local development).

## Default loop

1. `speckit-constitution` — process/governance only  
2. Refine `docs/PRODUCT.md` when **global** product intent changes (then update the matching slice spec)  
3. `speckit-specify` → optional `speckit-clarify` — delivery slices that follow PRODUCT.md (do not copy PRODUCT.md)  
4. `speckit-plan` (stack + package layout)  
5. `speckit-tasks` → optional `speckit-analyze`  
6. `speckit-implement` — **test-first** (constitution VI): failing tests from the approved spec, then code; run `pnpm verify` after coherent batches  
7. `speckit-converge` — append remaining tasks; do not silently change specs  

Do **not** implement from PRODUCT.md alone. PRODUCT.md is the rulebook; the feature spec is the slice; tests lock the slice; code is last.

## Hard rules

- `docs/PRODUCT.md` owns global product rules; feature specs own slice behavior; plans own tech; tests lock the slice; harness owns pass/fail.
- **Test-first** for specified runtime behavior: tests MUST fail before the implementation that makes them pass. Tests MUST NOT invent behavior outside PRODUCT.md and the approved spec.
- No secrets in git. No cross-package boundary violations (UI in backend, DB in frontend).
- Prefer the smallest change that satisfies the approved spec.
