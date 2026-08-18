# Stack decision

**Status:** core runtime locked (2026-08-10); remaining tooling still open

## Locked

| Layer | Choice | Notes |
| --- | --- | --- |
| Runtime | **Node.js 24** | Matches Dev Container base (`node:24-bookworm-slim`) |
| Language | **TypeScript** | Implied by NestJS + React tooling; exact TS version pinned at scaffold |
| Backend | **NestJS 11** (latest `@nestjs/*`) | Pin current latest at scaffold (today: `@nestjs/core` **11.1.29**) |
| Frontend | **React 19** (latest `react` / `react-dom`) | Pin current latest at scaffold (today: **19.2.8**) |
| Database | **PostgreSQL 18** (wire-compatible) | App MUST run against stock **PostgreSQL 18** locally (**Compose Postgres-only** for host and Dev Container, **and** a Postgres pod on **kind**) and be deployable to **Amazon RDS for PostgreSQL** and **Amazon Aurora PostgreSQL**. Prefer portable SQL; avoid engine-specific extensions or Aurora-only APIs unless gated and documented. Pin current latest minor for local/RDS (today: **18.4**); Aurora uses the closest supported Aurora PostgreSQL-compatible version. Do not use PG 19 beta. |
| Default DB host | **Amazon RDS for PostgreSQL** | Default managed DB for POC/dev (and typical envs): **RDS PostgreSQL 18**, prefer small single-AZ for cost. **Aurora PostgreSQL** remains a supported alternate target; do not assume Aurora-only features. |
| RDS idle control | **Wake/sleep + Lambda starter** | This is a **POC, not 24/7**. **Primary operator path:** **`pnpm wake`** / **`pnpm sleep`** / **`pnpm status`**. **Sleep MUST stop or tear down every stoppable billable resource** in the POC env: at least **stop RDS**, **Fargate/Nest to zero**, **tear down the EKS cluster** (control plane MUST NOT keep billing), and **remove ALB/NAT** (and similar hourly network) that exist only for that env. **Wake** recreates/starts that stack, **blocks until usable**, then prints **ready** + URLs (exit 0). Accept a long wake (tens of minutes) after full teardown. **Secondary:** Lambda may start RDS on stray traffic **only while** enough network/API path still exists; after cluster teardown, **`pnpm wake` is the way back**. If someone hits the app before ready, UI shows “waking.” Wake/sleep is an **operator** tool, not Org Admin. |
| IaC | **AWS CDK** | AWS infrastructure (EKS/Fargate wiring, ECR, ALB controller deps as applicable, RDS, wake/sleep Lambda, IAM, Amplify-related resources as needed) is defined with **AWS CDK** (TypeScript preferred to match the monorepo). Not Terraform unless a later plan changes this. |
| CI/CD | **GitHub Actions + Amplify build** | **GitHub Actions** is the primary pipeline: PR/`verify`, **semantic-release** on the release branch, Nest image build/push to **ECR**, Kustomize/EKS deploy, and CDK diff/deploy as applicable. Auth to AWS via **GitHub OIDC → IAM** (no long-lived access keys). **Amplify Hosting** builds/deploys the Next.js app from the connected branch; Actions still runs frontend checks in `pnpm verify`. Not CodePipeline-first unless a later plan changes this. |
| Secrets | **Secrets Manager + External Secrets Operator** | App secrets (especially **RDS credentials**) live in **AWS Secrets Manager**. On EKS, **External Secrets Operator** syncs them into the Nest workload (env/files) using **IRSA**. Optional **SSM Parameter Store** for non-secret config only. GitHub Actions uses OIDC (not long-lived AWS keys) and MUST NOT be the runtime source of DB passwords. Local/dev uses gitignored `.env` / cluster secrets examples; never commit real credentials. Kubernetes Secrets are not the AWS source of truth. |
| Environments | **local + dev + staging + prod** | **local (three paths, this order):** (1) **Host:** Nest + Next on the machine; Compose **Postgres only**. (2) **Dev Container:** Nest + Next in the container; same Compose Postgres. (3) **kind** (`pnpm local:up`): Postgres, Nest, Next as **pods**. **dev / staging / prod:** AWS. Wake/sleep targets non-prod by default. |
| Cloud | **AWS** | Production deployment target for the monorepo stack |
| Backend deploy | **Kubernetes on Amazon EKS** | NestJS backend ships as container image(s) and runs on **EKS** (managed Kubernetes). Manifests/charts MUST be plain Kubernetes so the same workload can run on any conformant cluster locally or elsewhere; EKS is the AWS production target. |
| EKS compute | **Fargate + KEDA scale-to-zero; cluster down when sleeping** | While **ready**, backend pods run on **EKS Fargate**; **KEDA** may scale Nest to **min 0** when idle. While **sleeping**, the **EKS cluster MUST be torn down** so the control plane does not keep billing. Wake recreates the cluster from the same Kustomize/CDK packaging. Not EC2 node groups unless a later plan changes this. |
| Image registry | **Amazon ECR** | Nest backend images are built in CI, tagged (semantic-release version and/or git SHA), pushed to **ECR**, and pulled by EKS Fargate. Not Docker Hub/GHCR for the primary deploy path. |
| Ingress | **ALB via AWS Load Balancer Controller** | Public/private HTTP(S) to Nest uses Kubernetes **Ingress** with **ingressClassName: alb**. Controller provisions an **Application Load Balancer**. Fargate requires **IP target type** (`alb.ingress.kubernetes.io/target-type: ip`). Traffic path: **ALB → KEDA HTTP Add-on → Nest Service/pods**. Not NLB-as-primary for the API unless a later plan changes this. |
| Package manager | **pnpm** workspaces | Root workspace monorepo; `pnpm-workspace.yaml` + `pnpm-lock.yaml`. Prefer pnpm over npm/yarn/Lerna for installs and scripts. |
| Release | **semantic-release** | Conventional Commits → version bump, changelog, git tag / GitHub Release. Repo-level (or backend-focused) release; stamp that version on the backend container image for EKS. Not Lerna/Changesets unless we later publish multiple independent libraries. |
| Frontend shell | **Next.js** (latest App Router on React 19) | React UI is a Next.js app. Deployed to **AWS Amplify Hosting**. Pin current latest Next at scaffold. |
| Frontend deploy | **AWS Amplify Hosting** | Next.js app builds and hosts via Amplify (SSR/ISR as supported by Amplify). Not OpenNext or EKS for the frontend unless a later plan changes this. |
| ORM | **TypeORM** | NestJS ↔ PostgreSQL via TypeORM (`@nestjs/typeorm`). Must stay portable across stock Postgres 18, RDS, and Aurora (no engine-specific features). Use migrations; do not rely on `synchronize` in shared/non-local environments. |
| API contracts | **OpenAPI + Orval** | Nest publishes OpenAPI (e.g. `@nestjs/swagger`). Spec is exported to `packages/api-client/openapi.json`. **Orval** generates a typed fetch client under `packages/api-client` (`@expenses/api-client`). Next.js imports that package; do not hand-edit generated output. Spec is source of truth for HTTP routes and payloads. |
| Verify harness | **`pnpm verify`** | Root gate: **lint → typecheck → test** (fail fast). Lint = **ESLint** (flat config) + **Prettier**; typecheck = **`tsc --noEmit`** per package; tests = **Jest** (one runner for backend and frontend unit/integration). Packages expose `lint` / `typecheck` / `test`; root orchestrates via pnpm. Deploy-to-AWS checks stay in CI/CD, not local verify. |
| Local runtime | **Host (default) → Dev Container → kind** | **Three local paths:** (1) **Default:** Nest and Next on the **host** (Node + pnpm); **Docker Compose runs PostgreSQL 18 only** (no app Compose services). (2) **Dev Container:** Nest and Next **inside** the container; **same Compose Postgres-only**. (3) **kind:** `pnpm local:up` deploys Postgres + Nest + Next as **pods** (wait until health + URLs; `local:down` / `local:status`). Do not require these to share ports at the same time. Hosted UI remains Amplify. First AWS apply is `011-aws-deploy`; wake/sleep is `010`. |
| K8s packaging | **Kustomize** | Manifests as Kustomize **base + overlays** (kind, staging, prod). **kind overlay MUST run local Postgres 18 + Nest + frontend** so the developer can see those pieces as pods. Same Nest packaging toward EKS; env differences via overlays. Not Helm unless a later plan changes this. |
| Wake/sleep operator UX | **`pnpm wake` / `pnpm sleep` / `pnpm status`** | **sleep:** stop/tear down **all stoppable** billable POC resources (see RDS idle control). **wake:** bring them back and **wait until ready**. **status:** sleeping / waking / ready, plus a list of **what still bills while sleeping** (cannot be stopped without deleting data). Optional later: a button that runs the same flow. Auto-sleep on idle is still open. |
| Repo shape | Monorepo: **`apps/backend`**, **`apps/frontend`**, **`packages/api-client`** | Full rewrite; no revive of the prior Go tree. Workspace via pnpm. |

## Still open (decide in the relevant plan)

- AWS shape (decide in **`011-aws-deploy` plan**): Amplify app/branch wiring per env; CDK app/package layout (stages for dev/staging/prod); Actions workflow layout (verify vs release vs deploy per env); External Secrets Operator layout with Kustomize overlays (`kind`, `dev`, `staging`, `prod`)
- RDS wake wiring (decide in **`010` plan**): how Lambda sits in front of ALB when the API path still exists; idle **auto-sleep**; IAM for Start/Stop RDS and cluster teardown/recreate (dev/staging first; prod policy explicit). Manual wake/sleep that **maximizes stoppable savings** is locked above.
- Whether `pnpm verify` also includes package builds and/or OpenAPI client freshness checks (add once scaffolded)
- kind cluster config details (node count, local registry / image load workflow) in the Dev Container
- Exact apps/packages naming if we prefer `@repo/*` scopes over `@expenses/*`
- Kustomize directory layout under the repo (e.g. `deploy/backend/` with overlays per env)

## Environment

Dev Container remains available (`.devcontainer/`). Local default is **host + Compose Postgres**. Then Dev Container. Then **kind**. First AWS deploy is **`011-aws-deploy`**. Wake/sleep is **`010`**.

## Pinning policy

Record major lines here; pin exact package versions in lockfiles when packages are scaffolded. Prefer latest stable on the locked major lines unless a plan documents an exception.
