---
description: "Task list for 002-bootstrap monorepo scaffold"
---

# Tasks: Monorepo Bootstrap

**Input**: Design documents from `/specs/002-bootstrap/`

**Prerequisites**: [plan.md](./plan.md) (Approved 2026-08-19), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Tests**: Included. Constitution VI bootstrap exception: stand up Jest so tests can run, then add **failing** tests for remaining acceptance criteria, then implement. `pnpm verify` stays offline (no Compose/kind). Live paths are [quickstart.md](./quickstart.md).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Backend: `apps/backend/`
- Frontend: `apps/frontend/`
- API client: `packages/api-client/`
- Root harness: `package.json`, `scripts/`, `compose.yaml`
- kind: `deploy/kind/`, `deploy/kustomize/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: pnpm workspace and TypeScript 6.0.3 pin so packages can be added

- [ ] T001 Create `pnpm-workspace.yaml` with `apps/*` and `packages/*`
- [ ] T002 Create root `package.json` (`private`, name `expenses-analyzer`, `engines.node` >= 24, `packageManager` pnpm via Corepack, shared `typescript` **6.0.3**)
- [ ] T003 [P] Create `apps/backend/package.json` named `@expenses/backend` with placeholder `lint` / `typecheck` / `test` scripts
- [ ] T004 [P] Create `apps/frontend/package.json` named `@expenses/frontend` with placeholder `lint` / `typecheck` / `test` scripts and `workspace:*` dep on `@expenses/api-client`
- [ ] T005 [P] Create `packages/api-client/package.json` named `@expenses/api-client` with placeholder `lint` / `typecheck` / `test` scripts
- [ ] T006 Add `tsconfig.base.json` at repo root targeting TypeScript 6.0.3 (`strict` on)
- [ ] T007 [P] Add `apps/backend/tsconfig.json`, `apps/frontend/tsconfig.json`, and `packages/api-client/tsconfig.json` extending `tsconfig.base.json`
- [ ] T008 [P] Add root `eslint.config.mjs` (flat) and `.prettierrc`
- [ ] T009 Run `pnpm install` so `pnpm-lock.yaml` is produced (no secrets; do not commit unless asked)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Jest and package scripts so tests can execute (constitution VI harness bootstrap exception). No user-story acceptance behavior yet beyond “the runner runs.”

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T010 Add Jest to `@expenses/backend`, `@expenses/frontend`, and `@expenses/api-client` with `apps/backend/jest.config.ts`, `apps/frontend/jest.config.ts`, and `packages/api-client/jest.config.ts`
- [ ] T011 Wire each package `lint` (ESLint + Prettier check), `typecheck` (`tsc --noEmit`), and `test` (Jest) in `apps/backend/package.json`, `apps/frontend/package.json`, and `packages/api-client/package.json`
- [ ] T012 Add one trivial passing smoke test per package in `apps/backend/test/smoke.spec.ts`, `apps/frontend/test/smoke.spec.ts`, and `packages/api-client/test/smoke.spec.ts` so Jest executes (not health/domain behavior)
- [ ] T013 Create `scripts/verify.mjs` that orchestrates lint → typecheck → test and root `package.json` scripts `lint`, `typecheck`, `test`, `verify` calling it
- [ ] T014 Confirm `pnpm test` runs via root `package.json` and the three Jest configs without Compose or kind

**Checkpoint**: Foundation ready — `pnpm test` executes. Remaining acceptance criteria MUST be failing tests before the code that satisfies them.

---

## Phase 3: User Story 1 - Empty workspace is installable and verifiable (Priority: P1) 🎯 MVP

**Goal**: Root `pnpm verify` is fail-fast lint → typecheck → test and fails if any workspace package is missing those scripts.

**Independent Test**: Fresh `pnpm install` + `pnpm verify` succeeds with no expenses/auth/org code; Compose/kind not required.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T015 [US1] Add failing tests in `scripts/verify.test.mjs` that `scripts/verify.mjs` runs lint then typecheck then test and exits non-zero on first failure (later steps not success)
- [ ] T016 [US1] Add a failing test in `scripts/verify.test.mjs` that verify fails when a workspace package lacks `lint`, `typecheck`, or `test`

### Implementation for User Story 1

- [ ] T017 [US1] Implement required-script check and fail-fast chain in `scripts/verify.mjs` until T015 and T016 pass
- [ ] T018 [US1] Confirm `pnpm verify` (root `package.json` → `scripts/verify.mjs`) passes on a clean tree without Compose or kind (SC-001, SC-006)

**Checkpoint**: User Story 1 is the MVP harness. Stop and validate before health/UI.

---

## Phase 4: User Story 2 - Backend and frontend packages exist as a shell (Priority: P1)

**Goal**: Nest health (Postgres ping), OpenAPI + Orval client, Next scaffold page that still renders when health is unhealthy.

**Independent Test**: Reviewer can run each package’s lint/typecheck/test; mocked health 200/503; frontend imports `@expenses/api-client`.

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T019 [P] [US2] Add failing contract tests in `apps/backend/test/health.spec.ts` for `GET /health`: 200 `{ status: "ok" }` when ping succeeds; 503 `{ status: "error" }` when Postgres is unreachable; no org/bill/user fields (`additionalProperties` false per `specs/002-bootstrap/contracts/health.openapi.yaml`)
- [ ] T020 [P] [US2] Add a failing test in `packages/api-client/test/openapi.spec.ts` that committed `packages/api-client/openapi.json` documents `/health` and `HealthResponse` without domain fields
- [ ] T021 [P] [US2] Add failing tests in `apps/frontend/test/scaffold.spec.tsx` that the scaffold page imports `@expenses/api-client` (not an ad-hoc fetch URL) and still renders when health is unhealthy

### Implementation for User Story 2

- [ ] T022 [US2] Scaffold NestJS 11 in `apps/backend/src/main.ts` and `apps/backend/src/app.module.ts` with TypeORM `synchronize: false`, zero entities, env-based Postgres config (Nest boots for tests without live DB)
- [ ] T023 [US2] Implement `GET /health` in `apps/backend/src/health/` (TypeORM `SELECT 1` / ping; 200 `ok` / 503 `error` only)
- [ ] T024 [US2] Add `@nestjs/swagger` and export script `apps/backend/src/export-openapi.ts` writing `packages/api-client/openapi.json` without a live Postgres
- [ ] T025 [US2] Add `packages/api-client/orval.config.ts` (`client: 'fetch'`, `baseUrl.runtime` `process.env.API_BASE_URL`) and generate committed client under `packages/api-client/src/` (do not hand-edit generated files)
- [ ] T026 [US2] Add root `package.json` script `generate:api` (export then Orval)
- [ ] T027 [US2] Scaffold Next.js App Router in `apps/frontend/app/page.tsx` (and `apps/frontend/app/layout.tsx`) using `@expenses/api-client`; show ok and unhealthy outcomes; `API_BASE_URL` default `http://localhost:3001`
- [ ] T028 [US2] Re-run `pnpm verify` via `scripts/verify.mjs` and keep it green without Compose

**Checkpoint**: Health contract + scaffold page + generated client. Live Postgres is US3.

---

## Phase 5: User Story 3 - Host machine + Compose Postgres (Priority: P1)

**Goal**: Default local path — Nest/Next on the host; Compose runs PostgreSQL 18 only; env examples with dummy credentials.

**Independent Test**: Parse/assert Compose is Postgres-only; env examples exist and contain no real secrets. Live `docker compose up` is quickstart, not `pnpm verify`.

### Tests for User Story 3 ⚠️

- [ ] T029 [P] [US3] Add a failing test in `scripts/compose.test.mjs` that root `compose.yaml` has a single `postgres` service using image `postgres:18.4` and no Nest/Next services
- [ ] T030 [P] [US3] Add a failing test in `scripts/env-examples.test.mjs` that `.env.example`, `.env.backend.example`, and `.env.frontend.example` exist, are not gitignored exceptions for real `.env`, and contain only dummy local credentials (no AWS keys / RDS passwords)

### Implementation for User Story 3

- [ ] T031 [US3] Add root `compose.yaml` with PostgreSQL 18.4 only (publish `5432`; dummy user/password/db `expenses`)
- [ ] T032 [P] [US3] Add `.env.example`, `.env.backend.example` (`DATABASE_HOST=localhost`, Nest `3001`), and `.env.frontend.example` (`API_BASE_URL=http://localhost:3001`) — dummy values only
- [ ] T033 [US3] Confirm `.gitignore` keeps `.env` / `.env.*` ignored except `!.env.example` and `!.env.*.example`
- [ ] T034 [US3] Document host start commands in `README.md` (Compose Postgres, then `pnpm --filter @expenses/backend start:dev` and `pnpm --filter @expenses/frontend dev`)

**Checkpoint**: Host path is specified and file-tested. Manual Compose is quickstart §1.

---

## Phase 6: User Story 4 - Dev Container + Compose Postgres (Priority: P1)

**Goal**: Nest/Next inside the Dev Container; same Compose Postgres in its own container; not the Dev Container as the DB.

**Independent Test**: Dev Container config reaches Compose Postgres via `host.docker.internal`; Postgres is not PID 1 of the workspace container.

### Tests for User Story 4 ⚠️

- [ ] T035 [US4] Add a failing test in `scripts/devcontainer.test.mjs` that `.devcontainer/devcontainer.json` does not run Postgres as the workspace service and documents `host.docker.internal` / port forwards `3000`, `3001`, `5432`

### Implementation for User Story 4

- [ ] T036 [US4] Extend `.devcontainer/devcontainer.json` with port forwards, `host.docker.internal:host-gateway`, and Docker-outside-of-Docker (docker.sock) so Compose/kind use host Docker
- [ ] T037 [US4] Update `.devcontainer/post-create.sh` to enable Corepack/pnpm, and install `kubectl` + `kind` for the optional kind path
- [ ] T038 [US4] Add Dev Container notes to `README.md` (same `compose.yaml`; `DATABASE_HOST=host.docker.internal`; Dev Container not required for host path)

**Checkpoint**: US3 Compose file remains the only Postgres for this path.

---

## Phase 7: User Story 5 - Local kind deploy (Priority: P1)

**Goal**: `pnpm local:up` / `down` / `status` for kind: Postgres + Nest + Next pods; wait until health 200; print URLs on 8080/8081. Verify still passes with kind down.

**Independent Test**: Scripts and Kustomize overlay exist and describe Postgres, Nest, and frontend. Live kind is quickstart §3, not the harness.

### Tests for User Story 5 ⚠️

- [ ] T039 [P] [US5] Add failing tests in `scripts/local.test.mjs` that `scripts/local/up.sh`, `scripts/local/down.sh`, and `scripts/local/status.sh` exist and root `package.json` exposes `local:up`, `local:down`, `local:status`
- [ ] T040 [P] [US5] Add a failing test in `scripts/kustomize-kind.test.mjs` that `deploy/kustomize/overlays/kind` includes Postgres, Nest, and frontend (and does not introduce Helm)

### Implementation for User Story 5

- [ ] T041 [US5] Add `deploy/kind/cluster.yaml` (one control-plane named `expenses-analyzer`, extraPortMappings host `8080`/`8081`)
- [ ] T042 [US5] Add Kustomize Nest base in `deploy/kustomize/base/` (Deployment + Service; no Postgres; no ALB)
- [ ] T043 [US5] Add `deploy/kustomize/overlays/kind/` with Postgres 18.4, frontend Deployment/Service, local image tags `expenses-analyzer-backend:local` / `expenses-analyzer-frontend:local`, `imagePullPolicy: Never` for app images, dummy Secret (not AWS)
- [ ] T044 [P] [US5] Add `apps/backend/Dockerfile` for the Nest image
- [ ] T045 [P] [US5] Add `apps/frontend/Dockerfile` with Next `output: 'standalone'`; in-cluster `API_BASE_URL` to the Nest Service
- [ ] T046 [US5] Implement `scripts/local/up.sh` (kind create/reuse, build, `kind load docker-image`, apply overlay, wait until `GET /health` is 200, print URLs, non-zero if never ready)
- [ ] T047 [P] [US5] Implement `scripts/local/down.sh` (idempotent `kind delete cluster --name expenses-analyzer`)
- [ ] T048 [P] [US5] Implement `scripts/local/status.sh` (pods + health; must not call AWS or `pnpm wake`/`sleep`)
- [ ] T049 [US5] Wire `local:up` / `local:down` / `local:status` in root `package.json`
- [ ] T050 [US5] Confirm `pnpm verify` via `scripts/verify.mjs` still passes with Compose and kind not running (SC-006)

**Checkpoint**: kind packaging is in-repo; live `pnpm local:up` is quickstart §3.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: CI verify, docs, no domain leak, no AWS in this slice

- [ ] T051 [P] Add `.github/workflows/verify.yml` (checkout, Corepack/pnpm, `pnpm verify`; **no AWS credentials**)
- [ ] T052 [P] Align `README.md` local-development section with host → Dev Container → kind and `pnpm verify` as definition of done
- [ ] T053 Add a test in `scripts/no-domain.test.mjs` that `apps/backend/src` and `apps/frontend/app` contain no organization/bill/user domain routes or screens (SC-005)
- [ ] T054 Confirm no CDK/EKS/Amplify/Helm/`pnpm wake` files under `deploy/` or `.github/`; `pnpm verify` via `scripts/verify.mjs` green
- [ ] T055 Follow `specs/002-bootstrap/quickstart.md` host path once (Compose + health `ok`) as a manual check; do not require it in CI

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational — **MVP harness**
- **User Story 2 (Phase 4)**: Depends on US1 (`pnpm verify` exists) so new tests fail then pass under the harness
- **User Story 3 (Phase 5)**: Depends on US2 (apps exist to document host start); Compose not required for US2 unit tests
- **User Story 4 (Phase 6)**: Depends on US3 (`compose.yaml` is the shared Postgres)
- **User Story 5 (Phase 7)**: Depends on US2 (Nest/Next to containerize)
- **Polish (Phase 8)**: Depends on US1–US5 tasks in repo (live kind optional)

### User Story Dependencies

- **User Story 1 (P1)**: After Phase 2 — no other stories
- **User Story 2 (P1)**: After US1 — health + client + page
- **User Story 3 (P1)**: After US2 — host Compose + env examples
- **User Story 4 (P1)**: After US3 — Dev Container uses the same Compose file
- **User Story 5 (P1)**: After US2 — kind overlay; independent of US4

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Contract/OpenAPI before Nest export and Orval
- Health before frontend client usage
- Story complete before the next phase unless marked independent (US5 vs US4)

### Parallel Opportunities

- T003, T004, T005 (package.json stubs)
- T007, T008 (tsconfigs vs eslint/prettier)
- T019, T020, T021 (US2 failing tests)
- T029, T030 (US3 failing tests)
- T039, T040 (US5 failing tests)
- T044, T045 (Dockerfiles)
- T047, T048 (down/status scripts)
- T051, T052 (CI vs README)

---

## Parallel Example: User Story 2

```bash
# Failing tests together:
Task: "Contract tests in apps/backend/test/health.spec.ts"
Task: "OpenAPI test in packages/api-client/test/openapi.spec.ts"
Task: "Scaffold tests in apps/frontend/test/scaffold.spec.tsx"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (TypeScript **6.0.3**)
2. Complete Phase 2: Foundational (Jest can run)
3. Complete Phase 3: User Story 1 (`pnpm verify`)
4. **STOP and VALIDATE**: `pnpm verify` without Docker
5. Then US2 (health/scaffold) — this slice is not demoable as a product UI until US2

### Incremental Delivery

1. Setup + Foundational → tests can run
2. US1 → `pnpm verify` (MVP harness)
3. US2 → health + Next + Orval
4. US3 → Compose Postgres-only + env examples
5. US4 → Dev Container wiring
6. US5 → kind scripts + Kustomize
7. Polish → Actions + README + no-domain check

### Parallel Team Strategy

1. Team completes Setup + Foundational + US1 together
2. Then: Developer A US2; Developer B can draft US3 compose tests in parallel after US1
3. US4 after US3; US5 after US2

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to spec user stories US1–US5
- TypeScript **6.0.3** for backend, frontend, and api-client (not 7)
- Do not implement `011-aws-deploy` or `010` wake/sleep
- Live Compose/kind are quickstart, not `pnpm verify`
- Verify tests fail before implementing each story
- Prefer Conventional Commits; do not commit secrets
