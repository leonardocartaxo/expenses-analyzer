# Feature Specification: Monorepo Bootstrap

**Feature Branch**: `002-bootstrap`

**Created**: 2026-08-18

**Status**: Draft

**Input**: Scaffold the expenses-analyzer monorepo per STACK.md with health/scaffold only; no expenses domain; defer full AWS deploy.

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

1. **Given** the scaffold, **When** a caller hits the backend health route, **Then** they receive a successful health response with no organization or bill data.
2. **Given** the scaffold, **When** a developer starts the frontend, **Then** a scaffold page loads and uses the generated API client for health (not an ad-hoc backend URL typed in the UI).
3. **Given** the OpenAPI spec for the health route, **When** the API client is generated, **Then** the frontend imports that package and generated output is not hand-edited.

---

### User Story 3 - Local Postgres and Kubernetes packaging stubs (Priority: P2)

Local development can start stock PostgreSQL 18 via Compose. Kubernetes packaging for the backend exists as Kustomize stubs suitable for kind. Full AWS (CDK, EKS, Amplify, RDS wake) is out of this slice.

**Why this priority**: Local runtime should match STACK.md without locking AWS deploy details.

**Independent Test**: Compose starts Postgres 18; Kustomize base/overlay files exist and apply conceptually to kind (cluster create may be documented rather than required in verify).

**Acceptance Scenarios**:

1. **Given** Docker is available, **When** the developer starts Compose for Postgres 18, **Then** a stock PostgreSQL 18 instance is reachable with credentials from a gitignored local env example (no real secrets in git).
2. **Given** the repo, **When** a reviewer inspects Kustomize packaging, **Then** there is a backend base plus at least a kind overlay stub.

---

### Edge Cases

- Verify must fail if a workspace package is missing `lint`, `typecheck`, or `test`.
- Committed env examples MUST NOT contain real credentials.
- Health route MUST NOT return organization, bill, or user domain data.
- Stack choices in [`STACK.md`](../../STACK.md) MUST be followed; do not introduce a second package manager, ORM, or frontend framework.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST be a pnpm workspace monorepo with `apps/backend`, `apps/frontend`, and `packages/api-client` as defined in STACK.md.
- **FR-002**: Root `pnpm verify` MUST run lint → typecheck → test and fail fast.
- **FR-003**: Lint MUST be ESLint (flat config) plus Prettier; typecheck MUST be `tsc --noEmit` per package; tests MUST be Jest.
- **FR-004**: Backend MUST be NestJS 11 with a health/scaffold HTTP route only (no expenses domain entities).
- **FR-005**: Frontend MUST be Next.js (App Router) on React 19 with a scaffold page only.
- **FR-006**: Backend MUST publish OpenAPI for the health route; Orval MUST generate the client into `packages/api-client`; frontend MUST consume that package.
- **FR-007**: Docker Compose MUST define stock PostgreSQL 18. TypeORM MAY connect for a smoke/health check; MUST NOT `synchronize` as the documented non-local default.
- **FR-008**: Kustomize stubs MUST exist for the Nest backend (base + kind overlay). Helm MUST NOT be introduced.
- **FR-009**: This slice MUST NOT implement sign-up, organizations, bills, CSV, dashboard, or AWS deploy (CDK/EKS/Amplify/RDS wake).
- **FR-010**: Secrets MUST NOT be committed; provide `.env` examples only.

### Key Entities

- **Workspace package**: backend, frontend, or api-client member of the pnpm workspace.
- **Health response**: non-domain liveness/readiness payload.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new developer can install and pass `pnpm verify` without implementing any PRODUCT.md domain feature.
- **SC-002**: Health can be called through the generated client from the frontend scaffold in one local runbook path.
- **SC-003**: Compose Postgres 18 starts from the documented local command.
- **SC-004**: No expenses-domain tables, routes, or screens ship in this slice.

## Assumptions

- STACK.md core runtime is locked; this spec implements that scaffold, not a stack debate.
- Full AWS deploy is a later slice; local Compose + Kustomize stubs are enough here.
- kind cluster create may be documented; `pnpm verify` does not require a live cluster.
- Package scope names default to `@expenses/*` unless a later plan changes STACK.md.
