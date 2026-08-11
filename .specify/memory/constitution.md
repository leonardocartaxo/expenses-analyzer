<!--
Sync Impact Report
- Version change: (none) → 1.0.0
- Modified principles: N/A (initial ratification)
- Added sections: Core Principles (I–V), Monorepo & Product Constraints, Development Workflow, Governance
- Removed sections: N/A
- Deferred: concrete language/framework versions (locked in bootstrap 2 / first `/speckit-plan`)
-->

# Expenses Analyzer Constitution

## Core Principles

### I. Spec-Driven Development (NON-NEGOTIABLE)
Every feature MUST begin with a Spec Kit feature spec (what/why and acceptance criteria)
before implementation. Technical choices belong in the plan; work belongs in tasks.
Agents MUST NOT invent product behavior outside an approved spec. When intent changes,
update the spec first, then re-plan or converge—do not silently diverge in code.

### II. Monorepo With Clear Package Boundaries
This repository is a single monorepo containing at least a **backend** package and a
**frontend** package, plus optional **shared** packages for cross-cutting contracts.
Backend owns APIs, persistence, and domain rules. Frontend owns UI and client UX.
Shared code MUST be limited to stable contracts (types, schemas, OpenAPI clients, constants)—
never UI components in backend or database access in frontend. Features that span packages
MUST be specified once and implemented with explicit tasks per package.

### III. Contract-First Across Backend and Frontend
Cross-package interfaces MUST be defined before parallel implementation (API schemas,
events, shared types, or equivalent). Frontend MUST consume the published contract, not
ad-hoc backend internals. Breaking contract changes require a versioned migration plan
and updates to both packages in the same feature workstream when feasible.

### IV. Harness as Definition of Done
A change is done only when the project verification harness passes (lint, typecheck,
unit/integration tests, and any package-specific checks wired into the root verify
command). Agents MUST run the harness after meaningful task batches and fix failures
before claiming completion. Specs MUST include testable acceptance criteria the harness
can enforce.

### V. Simplicity and Full Rewrite Discipline
This reboot is a **full rewrite**: do not port legacy structure, patterns, or debt from
the previous Go codebase unless a current spec explicitly requires equivalent behavior.
Prefer the smallest design that satisfies the spec (YAGNI). New abstractions require
justification in the plan. Avoid premature microservices, dual write paths, or unused
framework layers.

## Monorepo & Product Constraints

- **Product**: personal/household expenses analysis—capture, categorize, and understand
  spending (exact scope per feature specs).
- **Layout (target)**: package roots such as `apps/` or `packages/` for backend and
  frontend; exact names are fixed when the stack is chosen. Specs and Spec Kit assets
  remain at the repo root (`.specify/`, `specs/`).
- **Stack**: language and frameworks are **not** ratified in this constitution. They MUST
  be chosen once (bootstrap 2) and recorded in the first technical plan and a short
  stack decision note. Until then, do not scaffold application code that locks a stack.
- **Secrets**: never commit credentials, `.env` files with secrets, or private keys.
  Use env examples and local-only overrides.
- **Quality bar**: typed boundaries where the stack supports them; no unchecked `any`
  leakage across package contracts; structured logging on the backend once scaffolded.

## Development Workflow

1. **Constitution** — governance lives here; amend before changing process rules.
2. **Specify** — `/speckit-specify` (and `/speckit-clarify` when ambiguous).
3. **Plan** — `/speckit-plan` with stack and monorepo package design.
4. **Tasks** — `/speckit-tasks`; prefer `/speckit-analyze` before implement for
   multi-package features.
5. **Implement** — `/speckit-implement` task-by-task; run harness after each coherent
   batch; human review for contract and UX-sensitive changes.
6. **Converge** — `/speckit-converge` after implement to append missed work—never
   silently rewrite specs during converge.

Human approval gates: after specify and after plan for non-trivial features. Small
bugfixes that do not change product intent may skip a new feature branch of specs only
when an existing spec already covers the behavior; otherwise clarify or amend the spec.

## Governance

This constitution supersedes informal habits and ad-hoc agent prompts. Amendments MUST
update this file, bump **Version** (MAJOR for incompatible principle changes, MINOR for
new principles/sections, PATCH for clarifications), set **Last Amended** to the change
date, and note impact in the Sync Impact Report comment. PRs and agent sessions MUST
respect Spec Kit artifacts and the harness. Complexity beyond the constitution requires
explicit plan justification. Runtime agent entrypoint: `AGENTS.md` (points here).

**Version**: 1.0.0 | **Ratified**: 2026-08-10 | **Last Amended**: 2026-08-10
