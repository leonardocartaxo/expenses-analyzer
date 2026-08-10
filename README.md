# Expenses Analyzer

Personal and organization **bills / expenses** analysis: all bills live on **organizations** (including type **personal** or **company**). See [`docs/PRODUCT.md`](docs/PRODUCT.md).

**Global product rules (what / who / how):** [`docs/PRODUCT.md`](docs/PRODUCT.md) — edit that file when product intent changes.  
Feature specs under [`specs/`](specs/) deliver slices that must align with PRODUCT.md.  
This README is the human entrypoint for the reboot.

This repo is a **full rewrite** (greenfield monorepo). Do not revive the old Go tree unless a current spec says so.

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

Application packages are not scaffolded yet; stack and Spec Kit governance are in place.

## Spec-driven development (how we build)

**Specs own behavior. Plans own tech. Tasks own work. The harness owns pass/fail.**

Governance: [`.specify/memory/constitution.md`](.specify/memory/constitution.md) (process).  
Product: [`docs/PRODUCT.md`](docs/PRODUCT.md).  
Agent entrypoint: [`AGENTS.md`](AGENTS.md).  
Skills live under [`.cursor/skills/`](.cursor/skills/) (`speckit-*`).

Use **slash commands** in Cursor chat to start each phase. Put the feature details in the same message after the command.

### Day-to-day loop (one feature)

```text
specify → clarify? → plan → tasks → analyze? → implement → converge?
```

| Step | Command | Output | Your job |
| --- | --- | --- | --- |
| 1 | `/speckit-specify` | `specs/NNN-name/spec.md` | Describe **what/why**; approve the spec |
| 2 | `/speckit-clarify` | Updated spec | Answer up to a few focused questions |
| 3 | `/speckit-plan` | `plan.md` (+ design notes) | Approve tech approach (uses `STACK.md`) |
| 4 | `/speckit-tasks` | `tasks.md` | Skim the checklist |
| 5 | `/speckit-analyze` | Consistency report | Optional before implement |
| 6 | `/speckit-implement` | Code | Review PRs; require **`pnpm verify`** |
| 7 | `/speckit-converge` | New tasks only | Use when something from the spec is still missing |

**Do not** jump to implement without an approved spec (and plan for non-trivial work).  
**Do not** change product behavior only in code—amend the spec first.

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

Execute tasks.md for bootstrap. After each coherent batch run pnpm verify
and fix failures before continuing.
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

Execute tasks for create-expense. Run pnpm verify after each coherent batch.
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
| Bug already covered by an approved spec | Fix in code + `pnpm verify` |
| New behavior or changed intent | `/speckit-specify` (or amend the spec), then plan/tasks/implement |

Example (covered by an existing spec):

```text
Fix the expense form so empty amount shows the validation message from the
create-expense spec. Run pnpm verify when done.
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
/speckit-implement        # build + pnpm verify
/speckit-converge         # append missed tasks; don’t rewrite specs
/speckit-checklist        # optional custom checklist
/speckit-taskstoissues    # optional: tasks → GitHub issues
```

## Local development

1. Open the repo in the **Dev Container** ([`.devcontainer/`](.devcontainer/)) — Cursor, VS Code, or WebStorm Remote Dev.  
2. After bootstrap is implemented:
   - Compose → Postgres 18  
   - kind + Kustomize → Nest API  
   - Next.js locally (Amplify for hosted UI)  
3. Never commit secrets or `.env` files with credentials.  
4. Definition of done: **`pnpm verify`** (and CI on the PR).

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

1. Refine and approve [`docs/PRODUCT.md`](docs/PRODUCT.md) (global business rules).  
2. `/speckit-specify` for **bootstrap** or a delivery slice that follows PRODUCT.md.  
3. Implement only from approved feature specs.
