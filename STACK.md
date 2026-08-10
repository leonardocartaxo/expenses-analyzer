# Stack decision

**Status:** core runtime locked (2026-08-10); remaining tooling still open

## Locked

| Layer | Choice | Notes |
| --- | --- | --- |
| Runtime | **Node.js 24** | Matches Dev Container base (`node:24-bookworm-slim`) |
| Language | **TypeScript** | Implied by NestJS + React tooling; exact TS version pinned at scaffold |
| Backend | **NestJS 11** (latest `@nestjs/*`) | Pin current latest at scaffold (today: `@nestjs/core` **11.1.29**) |
| Frontend | **React 19** (latest `react` / `react-dom`) | Pin current latest at scaffold (today: **19.2.8**) |
| Database | **PostgreSQL 18** (wire-compatible) | App MUST run against stock **PostgreSQL 18** locally/Compose, and be deployable to both **Amazon RDS for PostgreSQL** and **Amazon Aurora PostgreSQL**. Prefer portable SQL/features; avoid engine-specific extensions or Aurora-only APIs unless gated and documented. Pin current latest minor for local/RDS (today: **18.4**); Aurora uses the closest supported Aurora PostgreSQL-compatible version. Do not use PG 19 beta. |
| Default DB host | **Amazon RDS for PostgreSQL** | Default managed DB for POC/dev (and typical envs): **RDS PostgreSQL 18**, prefer small single-AZ for cost. **Aurora PostgreSQL** remains a supported alternate target; do not assume Aurora-only features. |
| RDS idle control | **Wake/sleep + Lambda starter** | **Primary:** repo `wake` / `sleep` scripts (or button) start/stop the POC stack — at least **start/stop RDS** and coordinate Nest/KEDA as needed; wait until Postgres is available before treating the API as ready. **Secondary:** a **Lambda in front of the API path** detects traffic while RDS is stopped and calls `StartDBInstance` (and related wake steps). **Accept that early requests can fail for minutes** while RDS starts; clients/UI should retry or show “waking.” Not per-request instant DB; storage still bills when stopped. |
| IaC | **AWS CDK** | AWS infrastructure (EKS/Fargate wiring, ECR, ALB controller deps as applicable, RDS, wake/sleep Lambda, IAM, Amplify-related resources as needed) is defined with **AWS CDK** (TypeScript preferred to match the monorepo). Not Terraform unless a later plan changes this. |
| CI/CD | **GitHub Actions + Amplify build** | **GitHub Actions** is the primary pipeline: PR/`verify`, **semantic-release** on the release branch, Nest image build/push to **ECR**, Kustomize/EKS deploy, and CDK diff/deploy as applicable. Auth to AWS via **GitHub OIDC → IAM** (no long-lived access keys). **Amplify Hosting** builds/deploys the Next.js app from the connected branch; Actions still runs frontend checks in `pnpm verify`. Not CodePipeline-first unless a later plan changes this. |
| Secrets | **Secrets Manager + External Secrets Operator** | App secrets (especially **RDS credentials**) live in **AWS Secrets Manager**. On EKS, **External Secrets Operator** syncs them into the Nest workload (env/files) using **IRSA**. Optional **SSM Parameter Store** for non-secret config only. GitHub Actions uses OIDC (not long-lived AWS keys) and MUST NOT be the runtime source of DB passwords. Local/dev uses gitignored `.env` / Compose; never commit real credentials. Kubernetes Secrets are not the AWS source of truth. |
| Environments | **local + dev + staging + prod** | **local:** Dev Container, Compose Postgres 18, kind. **dev / staging / prod:** separate AWS env wiring (CDK stages/stacks, Kustomize overlays, Amplify branches, Secrets Manager paths, RDS instances as applicable). Promote via CI (verify → release → deploy). Wake/sleep RDS automation targets non-prod by default unless explicitly enabled for prod. |
| Cloud | **AWS** | Production deployment target for the monorepo stack |
| Backend deploy | **Kubernetes on Amazon EKS** | NestJS backend ships as container image(s) and runs on **EKS** (managed Kubernetes). Manifests/charts MUST be plain Kubernetes so the same workload can run on any conformant cluster locally or elsewhere; EKS is the AWS production target. |
| EKS compute | **Fargate + KEDA scale-to-zero** | Backend pods run on **EKS Fargate**. **KEDA** autoscales the Nest Deployment, including **min replicas 0** when idle. Scale-from-zero for HTTP uses the **KEDA HTTP Add-on** (interceptor) behind the ALB so requests can wake Nest; accept cold starts. EKS control-plane cost still applies while the cluster exists; tear down or stop the cluster when the POC is unused if near-zero cloud spend is required. Not EC2 node groups unless a later plan changes this. |
| Image registry | **Amazon ECR** | Nest backend images are built in CI, tagged (semantic-release version and/or git SHA), pushed to **ECR**, and pulled by EKS Fargate. Not Docker Hub/GHCR for the primary deploy path. |
| Ingress | **ALB via AWS Load Balancer Controller** | Public/private HTTP(S) to Nest uses Kubernetes **Ingress** with **ingressClassName: alb**. Controller provisions an **Application Load Balancer**. Fargate requires **IP target type** (`alb.ingress.kubernetes.io/target-type: ip`). Traffic path: **ALB → KEDA HTTP Add-on → Nest Service/pods**. Not NLB-as-primary for the API unless a later plan changes this. |
| Package manager | **pnpm** workspaces | Root workspace monorepo; `pnpm-workspace.yaml` + `pnpm-lock.yaml`. Prefer pnpm over npm/yarn/Lerna for installs and scripts. |
| Release | **semantic-release** | Conventional Commits → version bump, changelog, git tag / GitHub Release. Repo-level (or backend-focused) release; stamp that version on the backend container image for EKS. Not Lerna/Changesets unless we later publish multiple independent libraries. |
| Frontend shell | **Next.js** (latest App Router on React 19) | React UI is a Next.js app. Deployed to **AWS Amplify Hosting**. Pin current latest Next at scaffold. |
| Frontend deploy | **AWS Amplify Hosting** | Next.js app builds and hosts via Amplify (SSR/ISR as supported by Amplify). Not OpenNext or EKS for the frontend unless a later plan changes this. |
| ORM | **TypeORM** | NestJS ↔ PostgreSQL via TypeORM (`@nestjs/typeorm`). Must stay portable across stock Postgres 18, RDS, and Aurora (no engine-specific features). Use migrations; do not rely on `synchronize` in shared/non-local environments. |
| API contracts | **OpenAPI + Orval** | Nest publishes OpenAPI (e.g. `@nestjs/swagger`). Spec is exported to `packages/api-client/openapi.json`. **Orval** generates a typed fetch client under `packages/api-client` (`@expenses/api-client`). Next.js imports that package; do not hand-edit generated output. Spec is source of truth for HTTP routes and payloads. |
| Verify harness | **`pnpm verify`** | Root gate: **lint → typecheck → test** (fail fast). Lint = **ESLint** (flat config) + **Prettier**; typecheck = **`tsc --noEmit`** per package; tests = **Jest** (one runner for backend and frontend unit/integration). Packages expose `lint` / `typecheck` / `test`; root orchestrates via pnpm. Deploy-to-AWS checks stay in CI/CD, not local verify. |
| Local runtime | **Compose + local Kubernetes (kind)** | **Docker Compose** for stock **PostgreSQL 18** (and similar sidecars). **kind** runs the Nest backend with the same Kubernetes packaging used for EKS. Dev Container integrates via Docker socket (or equivalent) so `kind` / `kubectl` work from the container. Next.js still develops locally (Amplify for hosted deploys). |
| K8s packaging | **Kustomize** | Backend manifests as Kustomize **base + overlays** (e.g. kind, staging, prod). Same base applied locally and on EKS; env differences via overlays. Not Helm unless a later plan changes this. |
| Repo shape | Monorepo: **`apps/backend`**, **`apps/frontend`**, **`packages/api-client`** | Full rewrite; no revive of the prior Go tree. Workspace via pnpm. |

## Still open (decide in first `/speckit-plan`)

- AWS shape: Amplify app/branch wiring per env for the Next.js package; CDK app/package layout (stages for dev/staging/prod); Actions workflow layout (verify vs release vs deploy per env)
- RDS wake details: how Lambda sits in front of ALB (e.g. Lambda@Edge / CloudFront / separate wake endpoint), idle detection for auto-sleep, and IAM permissions for StartDBInstance/StopDBInstance (dev/staging first; prod policy explicit)
- External Secrets Operator install/layout alongside Kustomize overlays (`kind`, `dev`, `staging`, `prod`)
- Whether `pnpm verify` also includes package builds and/or OpenAPI client freshness checks (add once scaffolded)
- kind cluster config details (node count, local registry / image load workflow) in the Dev Container
- Exact apps/packages naming if we prefer `@repo/*` scopes over `@expenses/*`
- Kustomize directory layout under the repo (e.g. `deploy/backend/` with overlays per env)

## Environment

Dev Container is in place (`.devcontainer/`). Local Compose mirrors the AWS-bound services once those choices are locked in the first plan.

## Pinning policy

Record major lines here; pin exact package versions in lockfiles when packages are scaffolded. Prefer latest stable on the locked major lines unless a plan documents an exception.
