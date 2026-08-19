# Operator command contracts (002-bootstrap)

These are the developer/operator interfaces for this slice. HTTP is specified in [health.openapi.yaml](./health.openapi.yaml).

## `pnpm verify`

- **Order (fail-fast)**: lint → typecheck → test.
- **Must pass** with Compose and kind **not** running.
- **Must fail** (non-zero) if any workspace package lacks `lint`, `typecheck`, or `test`.
- **Must fail** if lint, typecheck, or tests fail; later steps in that run are not success.
- **Must not** require AWS credentials or start local clusters.

## `pnpm generate:api`

- Export Nest OpenAPI to `packages/api-client/openapi.json` **without** a live Postgres.
- Run Orval into `packages/api-client`.
- Generated TypeScript is committed; do not hand-edit it.
- Not part of `pnpm verify`.

## Host / Dev Container Postgres

- Start/stop with Docker Compose (`docker compose up` / `down`) using root `compose.yaml`.
- **Not** `pnpm local:up`.
- Compose MUST run PostgreSQL 18 only.

## `pnpm local:up` (kind only)

- Create/reuse kind cluster `expenses-analyzer`.
- Build and load Nest + Next images; apply Kustomize kind overlay (Postgres + Nest + Next pods).
- **Block** until `GET /health` returns 200.
- Print frontend and backend URLs, then exit 0.
- Exit non-zero if health never becomes ready.

## `pnpm local:status`

- Report kind stack ready vs not-ready (pods + health).
- MUST NOT start AWS or call `pnpm wake` / `sleep` / `status`.

## `pnpm local:down`

- Tear down the kind cluster/workloads.
- Idempotent if already down (exit 0).
