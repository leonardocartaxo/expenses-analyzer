# Expenses Analyzer

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-rewrite-blue)](#status)
[![Spec Kit](https://img.shields.io/badge/process-spec--driven-informational)](#spec-driven-development-how-we-build)

Capture, categorize, and understand **bills and spending** inside **organizations** (personal or company). All bills live on an organization—not on a user account.

| | |
| --- | --- |
| **Name** | Expenses Analyzer |
| **Repository** | https://github.com/leonardocartaxo/expenses-analyzer |
| **Product** | [`docs/PRODUCT.md`](docs/PRODUCT.md) (**Approved**) |
| **Stack** | [`STACK.md`](STACK.md) |
| **License** | [MIT](LICENSE) |
| **Maintainer** | [Leonardo Cartaxo](https://github.com/leonardocartaxo) |

This README is the human entrypoint for the reboot. Global product rules: [`docs/PRODUCT.md`](docs/PRODUCT.md). Feature specs under [`specs/`](specs/) deliver slices that must align with PRODUCT.md. Edit PRODUCT.md when product intent changes.

## Status

This repo is a **full rewrite** (greenfield monorepo). Application packages (`apps/`, `pnpm verify`) are **not scaffolded yet**. Do not revive the old Go tree unless a current spec says so.

## Table of contents

- [What’s locked (stack)](#whats-locked-stack)
- [Spec-driven development (how we build)](#spec-driven-development-how-we-build)
- [Local development](#local-development)
- [Repo map (target)](#repo-map-target)
- [Commits and releases](#commits-and-releases)
- [Where to go next](#where-to-go-next)
- [Contributing](#contributing)
- [License](#license)

---

## What’s locked (stack)

See [`STACK.md`](STACK.md) for the full decision table. Short version:

| Area | Choice |
| --- | --- |
| Runtime | Node.js 24, TypeScript |
| Backend | NestJS 11 + TypeORM → PostgreSQL 18 |
| Frontend | Next.js (React 19) on AWS Amplify |
| API contract | OpenAPI from Nest → Orval client (`packages/api-client`) |
| Monorepo | pnpm workspaces: `apps/backend`, `apps/frontend`, `packages/api-client` |
| Local | Dev Container + Compose (Postgres) + kind (Nest) + Kustomize |
| Cloud | AWS: EKS Fargate + KEDA scale-to-zero, ALB, ECR, RDS (Aurora-compatible), CDK, GitHub Actions |
| Done when | `pnpm verify` (lint → typecheck → Jest) passes |

## Spec-driven development (how we build)

**PRODUCT.md owns global rules. Specs own a delivery slice. Plans own tech. Tests lock the slice. Code makes tests pass. The harness owns pass/fail.**

Governance (including **test-first**): [`.specify/memory/constitution.md`](.specify/memory/constitution.md).  
Product: [`docs/PRODUCT.md`](docs/PRODUCT.md) (**Approved**).  
Agent entrypoint: [`AGENTS.md`](AGENTS.md).  
Skills live under [`.cursor/skills/`](.cursor/skills/) (`speckit-*`).

### What you edit when

| If you need to… | Do this |
| --- | --- |
| Change a **global** product rule (who owns bills, roles, CSV gate, …) | Edit **`docs/PRODUCT.md` first**. Then amend or add a feature spec that follows it. |
| Deliver a **slice** of behavior already in PRODUCT.md | Write/approve a spec under `specs/` (`002-…`). Do **not** copy PRODUCT.md into the spec. |
| Choose **how** to build it | `/speckit-plan` after the spec is approved (`STACK.md`). |
| Make it real | **Tests first** (from the spec), then code until they pass. |

You do **not** add a paragraph to PRODUCT.md and then implement. PRODUCT.md is the rulebook; a slice spec is the ticket; tests prove the ticket; code is last.

```text
PRODUCT.md (global rules, already approved)
   → feature spec (this slice) → approve
      → plan → tasks
         → failing tests from the spec
            → implement until tests + pnpm verify pass
```

Use **slash commands** in Cursor chat to start each phase. Put the feature details in the same message after the command.

### Day-to-day loop (one feature)

```text
specify → clarify? → plan → tasks → analyze? → implement (test-first) → converge?
```

| Step | Command | Output | Your job |
| --- | --- | --- | --- |
| 1 | `/speckit-specify` | `specs/NNN-name/spec.md` | Slice **what/why** that follows PRODUCT.md; approve the spec |
| 2 | `/speckit-clarify` | Updated spec | Answer up to a few focused questions |
| 3 | `/speckit-plan` | `plan.md` (+ design notes) | Approve tech approach (uses `STACK.md`) |
| 4 | `/speckit-tasks` | `tasks.md` | Skim the checklist |
| 5 | `/speckit-analyze` | Consistency report | Optional before implement |
| 6 | `/speckit-implement` | Tests, then code | **Test-first** (constitution VI): write failing tests from the spec, then implement; require **`pnpm verify`** |
| 7 | `/speckit-converge` | New tasks only | Use when something from the spec is still missing |

**Do not** jump to implement without an approved spec (and plan for non-trivial work).  
**Do not** implement specified behavior before the failing tests that encode it.  
**Do not** change product behavior only in code—amend PRODUCT.md and/or the spec first.  
**Do not** invent tests for behavior that is not in PRODUCT.md or the approved spec.

---

### Full example A — Bootstrap (scaffold the monorepo)

Use this when creating the empty Nest/Next/api-client tree and harness (no expenses domain yet).

**1. Specify**

```text
/speckit-specify

Bootstrap the expenses-analyzer monorepo per STACK.md:
- pnpm workspaces with apps/backend (NestJS), apps/frontend (Next.js),
  packages/api-client (Orval)
- root pnpm verify = ESLint + Prettier + tsc + Jest
- Docker Compose for Postgres 18
- kind + Kustomize stubs for the Nest backend
- no real expenses domain yet — health/scaffold only
- defer full CDK/EKS/Amplify deploy to a later feature
```

Then open `specs/00x-…/spec.md`, read acceptance criteria, and **approve** (or ask for edits).

**2. Clarify** (only if the agent asks)

```text
/speckit-clarify
```

**3. Plan** (after you approve the spec)

```text
/speckit-plan

Plan the approved bootstrap spec using STACK.md.
Defer full AWS deploy; stub local Compose/kind/Kustomize only.
```

**Approve** `plan.md`.

**4. Tasks**

```text
/speckit-tasks
```

Optional:

```text
/speckit-analyze
```

**5. Implement**

```text
/speckit-implement

Execute tasks.md for bootstrap. Test-first: after the Jest harness can run,
add failing tests for health/scaffold acceptance, then implement.
After each coherent batch run pnpm verify and fix failures before continuing.
```

**6. Converge** (if something from the spec is still missing)

```text
/speckit-converge

Compare the repo to the bootstrap spec/plan/tasks and append remaining
unbuilt work as new tasks. Do not silently change the spec.
```

Then run `/speckit-implement` again for any new tasks.

---

### Full example B — Product feature (create expense)

After bootstrap exists and product rules are specified.

**1. Specify**

```text
/speckit-specify

As a user I can create an expense with amount (positive integer cents) and
currency (USD or BRL). The system persists it and returns an id plus fields.
The UI has a simple form that uses the generated API client.
Invalid amounts show a clear error. Include acceptance criteria for the harness.
```

**2–7. Same sequence**

```text
/speckit-plan
```

```text
/speckit-tasks
```

```text
/speckit-analyze
```

```text
/speckit-implement

Execute tasks test-first: failing tests from the spec, then code.
Run pnpm verify after each coherent batch.
```

```text
/speckit-converge
```

---

### Specifying product / business rules

Edit the canonical file **[`docs/PRODUCT.md`](docs/PRODUCT.md)** (what / who / how / BR-*).  
Do not create `.specify/memory/product.md`. Do not duplicate rules into feature specs.

Delivery slices still use Spec Kit:

```text
/speckit-specify

Implement … following docs/PRODUCT.md (BR-…).
…
```

Example shape for a later slice (not for redefining the whole product):

```text
/speckit-specify

Org membership and roles per docs/PRODUCT.md:
Org Admin, Bill Manager, Viewer; multi-org membership; space isolation.
Include acceptance criteria for the harness.
```

---

### Small bugs

| Situation | What to do |
| --- | --- |
| Bug already covered by an approved spec, **no test catches it** | Add a **failing** test from that spec, then fix until `pnpm verify` passes |
| Bug already caught by the harness | Fix the code; keep the test |
| New behavior or changed intent | Amend PRODUCT.md if it is a global rule; amend or add a spec; then plan/tasks; **test-first** implement |

Example (covered by an existing spec):

```text
Add a failing test for the empty-amount validation from the spec, then
fix the form so that test passes. Run pnpm verify when done.
```

---

### Cheat sheet

```text
/speckit-constitution     # rare — process/governance only
/speckit-specify          # always first for a feature
/speckit-clarify          # optional
/speckit-plan             # after you approve the spec
/speckit-tasks            # after you approve the plan
/speckit-analyze          # optional
/speckit-implement        # test-first, then pnpm verify
/speckit-converge         # append missed tasks; don’t rewrite specs
/speckit-checklist        # optional custom checklist
/speckit-taskstoissues    # optional: tasks → GitHub issues
```

## Local development

1. Open the repo in the **Dev Container** ([`.devcontainer/`](.devcontainer/)) — Cursor, VS Code, or WebStorm Remote Dev.  
2. **Cursor CLI (`agent`) login** (required inside the Dev Container if you use the Cursor agent CLI):
   - Post-create installs the Cursor CLI (`agent`) into `~/.local/bin`.
   - Open a **new terminal** in the container (so `PATH` picks up `~/.local/bin`), then run:

```bash
agent login
```

   - Complete the browser/device login flow when prompted. Without this, `agent` commands from the Dev Container will not be authenticated.
   - Check: `agent --version` (and that `agent` is on your `PATH`).

3. After bootstrap is implemented:
   - Compose → Postgres 18  
   - kind + Kustomize → Nest API  
   - Next.js locally (Amplify for hosted UI)  
4. Never commit secrets or `.env` files with credentials.  
5. Definition of done: **`pnpm verify`** (and CI on the PR).

Until bootstrap lands, there is no `pnpm verify` or `apps/` tree yet—use Spec Kit example A above to create them.

## Repo map (target)

```text
apps/backend          NestJS API
apps/frontend         Next.js UI
packages/api-client   Orval-generated OpenAPI client
docs/PRODUCT.md       Global product / business rules (canonical)
deploy/               Kustomize overlays (kind, dev, staging, prod) — when added
specs/                Feature specs (Spec Kit); align with PRODUCT.md
.specify/             Spec Kit templates, constitution, scripts
STACK.md              Locked stack decisions
AGENTS.md             Agent / SDD entrypoint
```

## Commits and releases

- Prefer [Conventional Commits](https://www.conventionalcommits.org/) (feeds **semantic-release**).  
- CI: GitHub Actions (verify, release, ECR/EKS); frontend host builds via Amplify.  
- Do not invent stack or product behavior outside approved specs and `STACK.md`.

## Where to go next

1. [`docs/PRODUCT.md`](docs/PRODUCT.md) is **Approved**. Delivery slices: [`specs/README.md`](specs/README.md).  
2. Approve and `/speckit-plan` **bootstrap** (`002-bootstrap`) next.  
3. Implement only from approved feature specs.

## Contributing

1. Open the repo in the **Dev Container**.  
2. Follow [Spec-driven development](#spec-driven-development-how-we-build): PRODUCT.md → slice spec → plan → **test-first** → `pnpm verify`.  
3. Prefer [Conventional Commits](https://www.conventionalcommits.org/).  
4. Do not commit secrets or `.env` files with credentials.

Issues and pull requests: https://github.com/leonardocartaxo/expenses-analyzer

## License

[MIT](LICENSE) © 2026 Leonardo Cartaxo
