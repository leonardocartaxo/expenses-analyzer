# Implementation Plan: Monorepo Bootstrap

**Branch**: `002-bootstrap` | **Date**: 2026-08-18 | **Spec**: [spec.md](./spec.md)

**Status**: Approved

**Approved**: 2026-08-19

**Input**: Feature specification from `/specs/002-bootstrap/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Scaffold the expenses-analyzer **pnpm monorepo** per [STACK.md](../../STACK.md): NestJS 11 health API, Next.js (React 19) scaffold page, Orval client in `packages/api-client`, and root **`pnpm verify`** (lint → typecheck → Jest). No PRODUCT.md domain. Local in this order: **(1) host Nest+Next + Compose PostgreSQL 18 only**, **(2) Dev Container + same Compose Postgres**, **(3) kind via `pnpm local:up`** (Postgres + Nest + Next as pods). **AWS first deploy is out of scope** (`011-aws-deploy`). Wake/sleep is `010`.

## Technical Context

**Language/Version**: **TypeScript 6.0.3** on **Node.js 24**, shared by backend, frontend, and `packages/api-client`. Package manager **pnpm** workspaces (`packageManager` field + Corepack).

**Primary Dependencies**: **NestJS 11** (`@nestjs/common` / `core` / `platform-express` / `swagger` / `typeorm`, pin current 11.x at scaffold; STACK.md notes `@nestjs/core` **11.1.29** as of 2026-08-10). **Next.js** latest App Router on **React 19** (pin at scaffold; STACK.md notes React **19.2.8**). **Orval** fetch client. **TypeORM**. ESLint flat config + Prettier. **Kustomize**. **kind**. Docker Compose.

**Storage**: **PostgreSQL 18.4** — Compose-only for host and Dev Container; stock Postgres pod on kind. No domain tables. TypeORM `synchronize: false`. RDS/Aurora are **not** this slice.

**Testing**: **Jest** per package, orchestrated by root `pnpm verify`. Offline unit/contract tests only in the harness. Compose/kind are quickstart, not verify.

**Target Platform**: Developer laptop (host or Dev Container) and local Kubernetes (**kind**). Production AWS is `011-aws-deploy`.

**Project Type**: Monorepo web application (`apps/backend` + `apps/frontend` + `packages/api-client`).

**Performance Goals**: N/A for scaffold. `pnpm verify` and `local:up` must be usable on a single developer machine; `local:up` may take minutes (image build + kind).

**Constraints**: Follow STACK.md; no second package manager, ORM, or UI framework. No secrets in git. Host / Dev Container / kind need not run concurrently. `pnpm verify` MUST pass without Compose, kind, or AWS. Constitution VI: test-first after the runner exists (bootstrap exception).

**Scale/Scope**: Three workspace packages, one health route, one scaffold page, one Compose service, one kind overlay. No expenses features, no AWS apply.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Status |
| --- | --- | --- |
| I. Spec-driven | Approved `specs/002-bootstrap/spec.md`; plan does not add product behavior | **PASS** |
| II. Monorepo boundaries | `apps/backend` (API+DB), `apps/frontend` (UI), `packages/api-client` (contracts only) | **PASS** |
| III. Contract-first | `GET /health` OpenAPI before parallel UI/API work; frontend consumes `@expenses/api-client` | **PASS** |
| IV. Harness as DoD | Root `pnpm verify` is created in this slice; done = harness green | **PASS** |
| V. Simplicity / rewrite | Greenfield Nest/Next; no Go revival; no Helm/CDK/microservices | **PASS** |
| VI. Test-first | Harness bootstrap exception: Jest first, then failing tests for remaining ACs, then code | **PASS** |
| Secrets | `.env` gitignored; examples only; kind dummy Secret not real RDS | **PASS** |
| No stack invention | Choices from STACK.md only | **PASS** |

No unjustified complexity. **Complexity Tracking** left empty.

### Post-design re-check

Design artifacts (`research.md`, `data-model.md`, `contracts/`, `quickstart.md`) still: health-only payload; TypeORM with zero entities; Compose Postgres-only; kind overlay separate from AWS overlays; verify offline. **PASS** — proceed to `/speckit-tasks`.

## Project Structure

### Documentation (this feature)

```text
specs/002-bootstrap/
├── plan.md              # This file
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/
│   ├── health.openapi.yaml
│   └── operator-commands.md
└── tasks.md             # Phase 2 — NOT created by /speckit-plan
```

### Source Code (repository root)

```text
apps/backend/
├── src/
│   ├── main.ts
│   ├── app.module.ts
│   └── health/                 # GET /health only
├── test/
├── Dockerfile                  # kind (and later EKS in 011)
├── package.json                # @expenses/backend
└── tsconfig.json

apps/frontend/
├── app/                        # Next.js App Router scaffold page
├── test/
├── Dockerfile                  # kind only (Amplify is 011)
├── package.json                # @expenses/frontend
└── tsconfig.json

packages/api-client/
├── openapi.json                # exported from Nest (committed)
├── orval.config.ts
├── src/                        # Orval output (committed, do not hand-edit)
└── package.json                # @expenses/api-client

compose.yaml                    # postgres:18.4 only
.env.example
.env.backend.example
.env.frontend.example

deploy/
├── kind/
│   └── cluster.yaml            # 1 control-plane, ports 8080/8081
└── kustomize/
    ├── base/                   # Nest Deployment + Service
    └── overlays/
        └── kind/               # + Postgres + frontend + local images

scripts/
├── verify.mjs                  # required-scripts check + fail-fast chain
└── local/
    ├── up.sh
    ├── down.sh
    └── status.sh

.github/workflows/verify.yml    # pnpm verify only; no AWS
.eslintrc / eslint.config.mjs
.prettierrc
pnpm-workspace.yaml
package.json                    # verify, local:*, generate:api
```

Existing `.devcontainer/` is **extended** (pnpm, ports, docker.sock / Docker-outside-of-Docker, kind/kubectl, `host.docker.internal`). Do not make it mandatory for path 1.

**Structure Decision**: STACK.md repo shape (`apps/backend`, `apps/frontend`, `packages/api-client`) plus `deploy/kustomize` for kind. Frontend lives in the kind overlay, not in Kustomize base, so `011` can reuse Nest base toward EKS while Amplify hosts Next. No `dev`/`staging`/`prod` overlays in this slice.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

None.

## Phase 0 & Phase 1 notes (for implement)

- Pin versions at scaffold per STACK.md; TypeScript **6.0.3** for all workspace packages. Do not change locked majors.
- Nest listens on **3001** by default (Next on **3000**) so the two do not clash.
- Test-first sequence: (1) workspace + Jest so tests execute, (2) failing tests for health shape, package scripts, Compose postgres-only, env examples, api-client import, local:* script presence, (3) implementation until `pnpm verify` is green.
- Live host/Dev Container/kind flows are proven via [quickstart.md](./quickstart.md), not the harness.
- `pnpm generate:api` must work without Postgres (see research.md).
- Do not implement `011` or `010` artifacts (CDK, Amplify, wake/sleep, ESO, KEDA, ALB).
