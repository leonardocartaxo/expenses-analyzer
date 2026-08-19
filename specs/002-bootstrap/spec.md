# Feature Specification: Monorepo Bootstrap

**Feature Branch**: `002-bootstrap`

**Created**: 2026-08-18

**Status**: Approved

**Approved**: 2026-08-19

**Input**: Scaffold the expenses-analyzer monorepo per STACK.md with health/scaffold only. Local in this order: (1) Nest+Next on the host + Compose Postgres only, (2) Dev Container + Compose Postgres, (3) kind (`pnpm local:up`) as pods. No expenses domain; defer AWS to `011-aws-deploy`.

## Clarifications

### Session 2026-08-19

- Q: When PostgreSQL is not reachable, should the health route still report success? → A: No — health succeeds only if Postgres can be pinged; otherwise it fails (unhealthy)
- Q: If health is unhealthy (for example Postgres is down), should the scaffold page still load? → A: Yes — the page loads and shows unhealthy/error from the generated client

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Empty workspace is installable and verifiable (Priority: P1)

A developer clones the repo, installs workspace dependencies, and runs the root verify command. Lint, typecheck, and tests all run and pass on the scaffold (no product features required).

**Why this priority**: Nothing else can land until the harness exists.

**Independent Test**: Fresh install + root verify succeeds with no expenses/auth/org code.

**Acceptance Scenarios**:

1. **Given** a clean checkout of this feature, **When** the developer installs workspace dependencies and runs the root verify command, **Then** lint, typecheck, and tests complete successfully (fail-fast order: lint → typecheck → test).
2. **Given** verify is running, **When** any of lint, typecheck, or tests fail, **Then** the command fails and later steps in that run are not treated as success.

---

### User Story 2 - Backend and frontend packages exist as a shell (Priority: P1)

The monorepo contains a backend app, a frontend app, and a shared API-client package. Each package is a real workspace member. The backend exposes a health/scaffold endpoint only. The frontend renders a scaffold page. The API client is generated from the backend OpenAPI document (or an equivalent committed spec for the health route).

**Why this priority**: Package boundaries must exist before any domain feature.

**Independent Test**: A reviewer can open backend, frontend, and api-client packages and run each package’s lint/typecheck/test scripts.

**Acceptance Scenarios**:

1. **Given** the scaffold and reachable PostgreSQL, **When** a caller hits the backend health route, **Then** they receive a successful health response with no organization or bill data.
2. **Given** the backend is running but PostgreSQL is unreachable, **When** a caller hits the health route, **Then** they receive an unsuccessful (unhealthy) response with no organization or bill data.
3. **Given** the scaffold, **When** a developer starts the frontend, **Then** a scaffold page loads and uses the generated API client for health (not an ad-hoc backend URL typed in the UI).
4. **Given** the frontend is running and health is unhealthy, **When** a developer opens the scaffold page, **Then** the page still loads and shows the unhealthy/error outcome from the generated client (it MUST NOT go blank or only render when health succeeds).
5. **Given** the OpenAPI spec for the health route, **When** the API client is generated, **Then** the frontend imports that package and generated output is not hand-edited.

---

### User Story 3 - Host machine + Compose Postgres (Priority: P1)

**First / default local option.** A developer does **not** use the Dev Container. They run **Nest and Next on their machine** (host Node + pnpm). **Docker Compose runs PostgreSQL 18 only** — no Nest/Next Compose services. Apps connect to that Postgres on localhost. Full AWS is out of this slice.

**Why this priority**: Lowest-friction inner loop for people who already have Node on the laptop.

**Independent Test**: Compose Postgres up; Nest and Next start on the host; health works against that Postgres.

**Acceptance Scenarios**:

1. **Given** Docker and Node/pnpm on the host, **When** the developer starts Compose, **Then** only PostgreSQL 18 runs as a container (Nest and Next are not Compose services).
2. **Given** that Postgres, **When** they start Nest and Next **on the host**, **Then** the scaffold frontend uses the generated client against Nest health.
3. **Given** committed env examples, **When** they point the host apps at Compose Postgres, **Then** no real secrets are required in git.

---

### User Story 4 - Dev Container + Compose Postgres (Priority: P1)

**Second local option.** A developer works **inside the Dev Container**. They start **Nest and Next in that container**. **PostgreSQL 18** uses the **same Compose-Postgres-only** idea (its own container), not the Dev Container as the DB server.

**Why this priority**: Same inner loop for IDEs that prefer the Dev Container (WebStorm Remote Dev, Cursor/VS Code).

**Independent Test**: Dev Container + Compose Postgres; Nest and Next start in the container; health works.

**Acceptance Scenarios**:

1. **Given** the Dev Container and Compose, **When** the developer starts Postgres via Compose, **Then** PostgreSQL 18 is in its **own** container.
2. **Given** that Postgres, **When** they start Nest and Next **inside the Dev Container**, **Then** the scaffold frontend uses the generated client against Nest health.

---

### User Story 5 - Local kind deploy (Priority: P1)

**Third local option.** A developer can run **`pnpm local:up`**. That starts **kind**, applies Kustomize, and runs **PostgreSQL 18, Nest, and Next as pods**. The command **waits until health is ready** and prints URLs. `local:down` / `local:status` exist. Optional; not required for daily host or Dev Container work.

**Why this priority**: Practice the k8s packaging locally before EKS.

**Independent Test**: local-up; Postgres, Nest, and frontend pods Ready; health works; local-down removes the kind stack.

**Acceptance Scenarios**:

1. **Given** Docker is available, **When** the developer runs **`pnpm local:up`**, **Then** kind runs Postgres 18 + Nest + frontend as pods, and the command succeeds only after Nest health is ready (then it prints URLs).
2. **Given** that stack is up, **When** they list pods, **Then** they see Postgres, Nest, and the frontend as running containers.
3. **Given** the stack is up, **When** they run **`pnpm local:status`**, **Then** they see ready vs not-ready without starting AWS.
4. **Given** the stack is up, **When** they run **`pnpm local:down`**, **Then** those kind workloads are torn down (idempotent if already down).
5. **Given** `pnpm verify`, **When** Compose and kind are not running, **Then** verify can still pass.

---

### Edge Cases

- Verify must fail if a workspace package is missing `lint`, `typecheck`, or `test`.
- Committed env examples MUST NOT contain real credentials.
- Health route MUST NOT return organization, bill, or user domain data.
- If PostgreSQL is unreachable, health MUST NOT report success (process-up alone is not healthy).
- If health is unhealthy, the scaffold page MUST still load and show that outcome via the generated client (MUST NOT blank out or require successful health to render).
- Host, Dev Container, and kind MUST NOT be required to run at the same time (port clashes are OK to avoid; document if they conflict).
- Stack choices in [`STACK.md`](../../STACK.md) MUST be followed; do not introduce a second package manager, ORM, or frontend framework.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST be a pnpm workspace monorepo with `apps/backend`, `apps/frontend`, and `packages/api-client` as defined in STACK.md.
- **FR-002**: Root `pnpm verify` MUST run lint → typecheck → test and fail fast.
- **FR-003**: Lint MUST be ESLint (flat config) plus Prettier; typecheck MUST be `tsc --noEmit` per package; tests MUST be Jest.
- **FR-004**: Backend MUST be NestJS 11 with a health/scaffold HTTP route only (no expenses domain entities). Health MUST succeed only when PostgreSQL is reachable (ping); if Postgres is unreachable, health MUST fail (unhealthy). Health MUST NOT succeed merely because the process is up.
- **FR-005**: Frontend MUST be Next.js (App Router) on React 19 with a scaffold page only. The page MUST load whether health succeeds or fails, MUST use the generated API client, and MUST show the unhealthy/error outcome when health is not successful.
- **FR-006**: Backend MUST publish OpenAPI for the health route; Orval MUST generate the client into `packages/api-client`; frontend MUST consume that package.
- **FR-007**: Docker Compose MUST run **PostgreSQL 18 only** (no Nest/Next Compose services). That Compose file is the DB for **host** and **Dev Container** inner loops. TypeORM MAY connect; MUST NOT `synchronize` as the documented non-local default.
- **FR-008**: **Host path (default):** Nest and Next MUST run on the developer’s machine (Node + pnpm). They MUST connect to Compose Postgres. The Dev Container MUST NOT be required for this path.
- **FR-009**: **Dev Container path:** Nest and Next MAY run inside the Dev Container against the same Compose-Postgres-only setup. Postgres MUST NOT run as the DB server inside the Dev Container PID namespace.
- **FR-010**: The kind overlay MUST run stock **PostgreSQL 18** plus Nest and frontend as cluster workloads.
- **FR-011**: Kustomize MUST include a Nest base plus a **kind** overlay that deploys Postgres, Nest, and the frontend. Helm MUST NOT be introduced.
- **FR-012**: This slice MUST NOT implement sign-up, organizations, bills, CSV, dashboard, or AWS deploy (CDK/EKS/Amplify). AWS first deploy is `011-aws-deploy`; wake/sleep is `010`.
- **FR-013**: Secrets MUST NOT be committed; provide `.env` examples only.
- **FR-014**: Operators MUST have **`pnpm local:up`**, **`pnpm local:down`**, and **`pnpm local:status`** for the **kind** path. **local:up** MUST wait until Nest health is ready, then print URLs. Host/Dev Container Postgres start is Compose, not `local:up`.

### Key Entities

- **Workspace package**: backend, frontend, or api-client member of the pnpm workspace.
- **Health response**: non-domain liveness/readiness payload. Success only when PostgreSQL answers a ping; failure when it does not.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new developer can install and pass `pnpm verify` without implementing any PRODUCT.md domain feature.
- **SC-002**: On the **host** (no Dev Container), health can be called through the generated client from Next against Nest, with Postgres only in Compose.
- **SC-003**: In the **Dev Container**, the same health path works with Compose Postgres.
- **SC-004**: `pnpm local:up` brings up **kind** with Postgres, Nest, and frontend as Ready pods and prints URLs after health succeeds.
- **SC-005**: No expenses-domain tables, routes, or screens ship in this slice.
- **SC-006**: `pnpm verify` passes without Compose or kind running.

## Assumptions

- STACK.md core runtime is locked; this spec implements that scaffold, not a stack debate.
- Full AWS deploy is `011-aws-deploy`; wake/sleep is `010`.
- **Three** local paths, in this order: host + Compose Postgres, Dev Container + Compose Postgres, kind.
- `pnpm verify` does not require Compose or kind.
- The three paths need not run concurrently.
- Package scope names default to `@expenses/*` unless a later plan changes STACK.md.
- `packages/api-client` is a library, not a pod.
