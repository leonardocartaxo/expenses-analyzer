# Research: Monorepo Bootstrap

**Branch**: `002-bootstrap` | **Date**: 2026-08-18

STACK.md is the stack source. This file records bootstrap-only decisions, including STACK.md items that were still open and must be fixed for implement. AWS first deploy remains `011-aws-deploy`. Wake/sleep remains `010`.

## 1. Stack source and deferred AWS

**Decision**: Implement the locked STACK.md runtime (Node 24, TypeScript **6.0.3**, NestJS 11, Next.js App Router on React 19, PostgreSQL 18, TypeORM, pnpm, OpenAPI + Orval, Kustomize, `pnpm verify`). Do **not** add CDK, EKS, ECR, Amplify, Secrets Manager, External Secrets, KEDA, ALB controller, GitHub OIDC-to-AWS, or `pnpm wake` / `sleep` / `status`.

**Rationale**: The approved spec is health/scaffold only. `011-aws-deploy` owns first hosted apply. `010` owns POC idle control. Constitution V forbids extra layers.

**Alternatives considered**: Scaffolding empty CDK/Amplify stubs now — rejected (YAGNI, wrong slice). Using Compose for Nest/Next — rejected (STACK.md and FR-007: Compose is Postgres-only).

## 2. TypeScript version

**Decision**: One workspace TypeScript: **6.0.3**, used by `@expenses/backend`, `@expenses/frontend`, and `@expenses/api-client` (`tsc --noEmit` in `pnpm verify`). Prefer a root or catalog pin so packages do not drift.

**Rationale**: NestJS 11’s CLI still consumes the TypeScript **compiler API** (`nest build` / `nest start`). TypeScript **7.0.2** is newer but does not expose that API yet, so Nest CLI fails. Next.js 16 documents a **5.1** minimum and supports 6 without flags; 7 needs `experimental.useTypeScriptCli`. **6.0.3** is the latest 6.x both frameworks support as a normal `typescript` dependency.

**Alternatives considered**: TypeScript 7.0.2 — rejected for this slice (Nest CLI + dual-package workaround). TypeScript 5.9 — older than the 6 line both already support.

## 3. Package names and workspace layout

**Decision**: pnpm workspaces with `apps/*` and `packages/*`. Package names: `@expenses/backend`, `@expenses/frontend`, `@expenses/api-client`. Root package is private `expenses-analyzer`. Use `workspace:*` for the frontend → api-client dependency.

**Rationale**: Spec assumption and STACK.md default `@expenses/*`. Paths already named in STACK.md (`apps/backend`, `apps/frontend`, `packages/api-client`).

**Alternatives considered**: `@repo/*` — STACK.md still open; do not switch without a STACK.md change. npm/yarn — forbidden by STACK.md.

## 4. Verify harness composition

**Decision**: Root `pnpm verify` is **lint → typecheck → test**, fail-fast (`&&`). It MUST NOT start Compose or kind. It MUST NOT require package production builds or OpenAPI regeneration. A small root checker MUST fail if any workspace package is missing `lint`, `typecheck`, or `test`. Lint = ESLint flat config plus Prettier check. Typecheck = `tsc --noEmit` per package. Tests = Jest.

**Rationale**: Matches STACK.md and FR-002/FR-003/SC-006. STACK.md left “build / OpenAPI freshness in verify” open — keep verify fast and Docker-free for this slice.

**Alternatives considered**: Adding `pnpm build` or Orval freshness to verify — defer until a later plan (likely after generate is stable). Turborepo/Nx — extra orchestration, not required.

## 5. Health HTTP contract

**Decision**: Single route `GET /health`. Success **200** body `{ "status": "ok" }`. Failure (Postgres unreachable) **503** `{ "status": "error" }`. No other properties (no org/bill/user fields). Nest implements it with TypeORM `SELECT 1` (or equivalent ping). Do **not** adopt `@nestjs/terminus` default payload (`info` / `details`) for this slice.

**Rationale**: Spec entity is a non-domain liveness/readiness payload. A two-value `status` field is enough for Next, Orval, and kind wait-on-health. Terminus JSON is Nest-specific noise in the shared contract.

**Alternatives considered**: Terminus `TypeOrmHealthIndicator` as the public body — rejected for contract simplicity (internal ping is still TypeORM). Process-only health that ignores Postgres — weaker match for “health works against that Postgres.”

## 6. TypeORM without domain entities

**Decision**: Wire `@nestjs/typeorm` + `typeorm` with `synchronize: false` in all documented configs. No entities, no migrations in this slice. Connection settings from env (see data-model). Nest **must boot for tests and OpenAPI export without a live Postgres**; tests mock the DataSource. Local host/Dev Container/kind run against real PostgreSQL 18.

**Rationale**: FR-007 allows TypeORM to connect and forbids `synchronize` as the non-local default. Constitution bootstrap exception: verify cannot require Compose. Empty schema is correct because SC-005 forbids expenses-domain tables.

**Alternatives considered**: Skip TypeORM until a domain slice — rejected; STACK.md already locked TypeORM and local Postgres is in this spec. SQLite for tests — rejected; not the locked engine.

## 7. OpenAPI + Orval client

**Decision**: Nest `@nestjs/swagger` is the runtime source. Export committed `packages/api-client/openapi.json`. Orval (`client: 'fetch'`) generates committed TypeScript under `packages/api-client` (e.g. `src/generated/`). Frontend imports `@expenses/api-client` only. Do not hand-edit generated files. Runtime base URL via Orval `baseUrl.runtime` using `process.env.API_BASE_URL` (Next server-side scaffold page). Root script `pnpm generate:api` = export OpenAPI (no live DB) then Orval. Export/generate is **not** part of `pnpm verify`.

**Rationale**: STACK.md contract-first rule. Committing generated output lets verify and the frontend typecheck without Nest running. Server Components avoid CORS for the scaffold page.

**Alternatives considered**: Hand-written fetch in the UI — forbidden (FR-006). Regenerating on every verify — slower, needs Nest; deferred. React Query/SWR Orval clients — extra frontend stack for a scaffold page.

## 8. Default ports and three local paths

**Decision**: Paths in this order, not concurrent.

| Path | Nest + Next | PostgreSQL 18 | Nest | Next | Notes |
| --- | --- | --- | --- | --- | --- |
| 1. Host (default) | Host Node + pnpm | Compose only | `localhost:3001` | `localhost:3000` | Dev Container not required |
| 2. Dev Container | Processes inside the container | Same Compose file, **separate** container | `localhost:3001` in the container | `localhost:3000` | DB host `host.docker.internal` |
| 3. kind | Pods | Postgres **pod** | printed URL (host **8081**) | printed URL (host **8080**) | `pnpm local:up` |

Compose publishes Postgres `5432` for paths 1–2. kind does **not** publish 5432 on the host (avoids clashing with Compose). kind uses extraPortMappings 8080/8081 so it is less likely to collide with the host inner loop; still document “do not run host apps and kind at once.”

**Rationale**: Spec user stories 3–5 and STACK.md local order. Port split (Next 3000, Nest 3001) avoids the default Nest/Next clash.

**Alternatives considered**: Nest/Next as Compose services — forbidden. Dev Container as the Postgres process — forbidden (FR-009). Shared ports 3000/3001 for kind — more clash with host.

## 9. Compose Postgres-only

**Decision**: Root `compose.yaml` with a single service `postgres`, image `postgres:18.4` (STACK.md pin). Dummy local credentials only (`expenses` / `expenses` / db `expenses`) referenced from `.env.example`, not real secrets. No Nest, Next, or admin-UI services.

**Rationale**: FR-007/FR-013. Same file for host and Dev Container.

**Alternatives considered**: `postgres:18` floating tag — worse reproducibility. pgAdmin sidecar — out of spec.

## 10. Dev Container updates

**Decision**: Keep the existing Node 24 bookworm-slim image. Add: Corepack/pnpm in post-create; port forwards 3000/3001/5432; `host.docker.internal:host-gateway`; Docker-outside-of-Docker (docker.sock) so `docker compose up` and kind use **host** Docker; kubectl + kind CLIs for the optional kind path. Postgres remains a Compose container, not PID 1 in the Dev Container.

**Rationale**: US4 + FR-009. Docker-outside-of-Docker matches “same Compose Postgres-only” and lets kind run from the container without DinD complexity.

**Alternatives considered**: Adding a Compose `dev` app service for Nest/Next — forbidden. Docker-in-Docker — heavier, breaks “Postgres is its own container on the host engine.”

## 11. kind + Kustomize layout

**Decision**:

- Cluster: one control-plane node, name `expenses-analyzer`, config `deploy/kind/cluster.yaml`.
- Images: build Nest/Next locally, `kind load docker-image` with tags `expenses-analyzer-backend:local` and `expenses-analyzer-frontend:local`. `imagePullPolicy: Never` for those. Postgres image pulled (`postgres:18.4`), `IfNotPresent`. **No local registry** in this slice.
- Packaging: `deploy/kustomize/base` = Nest Deployment + Service only (portable toward EKS in `011`). `deploy/kustomize/overlays/kind` adds Postgres, frontend Deployment/Service, local images, dummy Secret/ConfigMap. **No Helm. No ingress-nginx. No ALB annotations.** No `dev`/`staging`/`prod` overlays here (`011` adds them).
- Frontend image uses Next `output: 'standalone'`. In-cluster `API_BASE_URL` points at the Nest Service DNS. Browser talks only to the frontend NodePort.
- Scripts: `pnpm local:up` / `local:down` / `local:status` → `scripts/local/*.sh`. `local:up` waits until `GET /health` is 200, then prints URLs. `local:down` is `kind delete cluster` (idempotent). `local:status` does not touch AWS.

**Rationale**: Closes STACK.md open items for kind node count, image load, and Kustomize directory layout, within this spec. Base-without-Postgres keeps RDS for `011`.

**Alternatives considered**: Local registry / Tilt / Skaffold — extra moving parts. Helm — forbidden. Putting frontend in base — wrong for Amplify-hosted UI on AWS. Three-node kind — unnecessary for a scaffold.

## 12. CI in this slice

**Decision**: Add a **verify-only** GitHub Actions workflow (checkout, Corepack/pnpm, `pnpm verify`) with **no AWS credentials**. Do **not** add semantic-release, ECR push, or CDK.

**Rationale**: README already treats PR verify as part of done. `011` FR-006 says PR verify may already exist from `002`. A workflow with no AWS is not “AWS first deploy.”

**Alternatives considered**: Defer all Actions to `011` — possible, but then PRs have no harness until that slice. Full release pipeline now — belongs with image tagging in `011`.

## 13. Test-first (constitution VI)

**Decision**: Use the harness bootstrap exception: stand up Jest so tests can run, then add **failing** tests for remaining acceptance criteria, then implement. `pnpm verify` tests stay offline (mocks + file/contract assertions). Live Compose/kind checks are **quickstart**, not verify.

**Rationale**: Constitution VI plus SC-006. Spec edge cases (missing package scripts, env examples, Compose Postgres-only, health shape) are assertable without Docker.

**Alternatives considered**: Implementing the scaffold then adding tests — forbidden. Requiring kind in CI — contradicts SC-006.

## 14. Logging and other quality bars

**Decision**: Nest default logger is enough for this slice. Do not add a log pipeline. Structured JSON logging can wait for a later slice if needed; do not block bootstrap on pino.

**Rationale**: Health/scaffold has no PII. Constitution quality bar should not expand this slice into observability product work.

**Alternatives considered**: nestjs-pino now — extra dependency with no spec acceptance criterion.
