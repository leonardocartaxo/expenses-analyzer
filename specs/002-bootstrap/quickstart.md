# Quickstart: Monorepo Bootstrap

**Branch**: `002-bootstrap` | **Date**: 2026-08-18

Validate the scaffold against [spec.md](./spec.md). Details: [data-model.md](./data-model.md), [contracts/health.openapi.yaml](./contracts/health.openapi.yaml), [contracts/operator-commands.md](./contracts/operator-commands.md).

Do **not** run the three local paths at the same time. Preferred order: **host → Dev Container → kind**.

## Prerequisites

- Node.js 24 and pnpm (Corepack).
- Docker (Compose + kind).
- Host kind path also needs `kind` and `kubectl`.
- Dummy env from committed examples (copy to gitignored `.env`); no real secrets.

## 0. Install and verify (no Docker)

```bash
pnpm install
pnpm verify
```

**Expect**: lint, then typecheck, then Jest, all pass. This must work with Compose and kind stopped.

## 1. Host + Compose Postgres (default)

```bash
docker compose up -d
# copy env examples; host DATABASE_HOST=localhost; API_BASE_URL=http://localhost:3001
pnpm --filter @expenses/backend start:dev
pnpm --filter @expenses/frontend dev
```

**Expect**:

- `docker compose ps` shows **only** PostgreSQL 18 (no Nest/Next services).
- `GET http://localhost:3001/health` → `{"status":"ok"}`.
- `http://localhost:3000` scaffold page uses `@expenses/api-client` (not a hardcoded fetch URL in the UI).

Stop Nest/Next and `docker compose down` before path 3.

## 2. Dev Container + Compose Postgres

Open the repo in the Dev Container. Postgres is still Compose, not the Dev Container itself.

```bash
docker compose up -d
# env: DATABASE_HOST=host.docker.internal
pnpm --filter @expenses/backend start:dev
pnpm --filter @expenses/frontend dev
```

**Expect**: same health + scaffold behavior as host. `docker compose ps` still shows Postgres as its **own** container.

## 3. kind (`pnpm local:up`)

```bash
pnpm local:up
pnpm local:status
```

**Expect**: command returns only after Nest health is 200; prints URLs (frontend **8080**, backend **8081** per plan). `kubectl` shows Postgres, Nest, and frontend pods Ready. Open the printed frontend URL; scaffold still uses the generated client.

```bash
pnpm local:down
pnpm local:down   # idempotent
```

**Expect**: kind workloads gone; second down still succeeds.

## Negative checks

- `pnpm verify` with Compose/kind down still passes.
- Health JSON has no organization/bill/user fields.
- Git has `.env.example` (and similar) only — no committed `.env` with credentials.
- `compose.yaml` has no Nest/Next services.
- No CDK/EKS/Amplify apply in this slice (`011-aws-deploy`).
