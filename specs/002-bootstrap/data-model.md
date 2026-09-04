# Data Model: Monorepo Bootstrap

**Branch**: `002-bootstrap` | **Date**: 2026-08-18

This slice has **no expenses-domain persistence**. There are no User, Organization, Bill, Establishment, or Transaction tables, routes, or screens (SC-005, FR-012). TypeORM is connected with **zero entities** and **`synchronize: false`**.

## Entities (scaffold)

### Workspace package

A pnpm workspace member that the root harness can lint, typecheck, and test.

| Field | Rule |
| --- | --- |
| Path | One of `apps/backend`, `apps/frontend`, `packages/api-client` |
| `name` | `@expenses/backend`, `@expenses/frontend`, `@expenses/api-client` |
| Scripts | MUST define `lint`, `typecheck`, `test` |
| Role | Backend owns HTTP + DB connection. Frontend owns UI. api-client owns generated OpenAPI client only (not a runtime pod) |

**Validation**: Root verify fails if any member is missing a required script (spec edge case). Frontend MUST depend on `@expenses/api-client` via `workspace:*`, not a copied URL helper.

### Health response

Non-domain liveness/readiness payload returned by `GET /health`.

| Field | Type | Values | Notes |
| --- | --- | --- | --- |
| `status` | string | `ok` \| `error` | Only property allowed |

**HTTP**:

- `200` + `status=ok` when Nest is serving and Postgres answers a ping (`SELECT 1`).
- `503` + `status=error` when the ping fails.

**Validation**: JSON MUST NOT include organization, bill, user, or other PRODUCT.md fields. OpenAPI `additionalProperties: false`. Frontend may render either outcome; the scaffold page still loads.

**State transitions**:

```text
Postgres reachable     → 200 { status: "ok" }
Postgres unreachable   → 503 { status: "error" }
```

No stored health rows. No history.

### Local database connection (config, not a table)

How Nest reaches PostgreSQL 18. Not an application entity.

| Field | Host / Compose | Dev Container | kind overlay |
| --- | --- | --- | --- |
| Engine | PostgreSQL 18.4 | same Compose | `postgres:18.4` pod |
| Host | `localhost` (published 5432) | `host.docker.internal` | Postgres Service DNS |
| Port | `5432` | `5432` | `5432` in-cluster |
| User / password / db | dummy `expenses` / `expenses` / `expenses` from env examples | same | dummy Secret in kind overlay |
| `synchronize` | `false` | `false` | `false` |

**Validation**: Committed `.env.example` / `.env.*.example` MUST NOT contain real credentials, AWS keys, or RDS passwords. Dummy local-only values are required so git stays secret-free (FR-013).

### Compose service

| Field | Rule |
| --- | --- |
| File | `compose.yaml` at repo root |
| Allowed services | `postgres` only |
| Image | `postgres:18.4` |
| App services | Nest and Next MUST NOT appear |

Used by host and Dev Container paths. **Not** used by `pnpm local:up`.

### kind workload

Runtime objects in overlay `deploy/kustomize/overlays/kind` (not domain entities).

| Workload | Image | Role |
| --- | --- | --- |
| Postgres | `postgres:18.4` | Stock engine; in-cluster only |
| Nest | `expenses-analyzer-backend:local` | Serves `/health` |
| Next | `expenses-analyzer-frontend:local` | Scaffold page; calls Nest via generated client using in-cluster `API_BASE_URL` |

`packages/api-client` is **not** a workload.

## Out of scope (must not appear)

- PRODUCT.md tables and enums (org types, roles, bill kinds, tags, CSV batches).
- AWS resources (RDS, EKS, Amplify, Secrets Manager).
- Wake/sleep state machines (`010`).
