# Feature Specification: AWS First Deploy

**Feature Branch**: `011-aws-deploy`

**Created**: 2026-08-18

**Status**: Draft

**Input**: First-time hosted deploy of the monorepo to AWS per STACK.md: CDK, EKS Fargate, ECR, RDS PostgreSQL 18, Amplify Hosting, Secrets Manager, GitHub Actions OIDC. Target **dev** as the live POC environment. Health/scaffold is enough (depends on `002-bootstrap`). Do **not** implement wake/sleep (`010`) or expenses domain.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Hosted backend is reachable (Priority: P1)

An operator (or CI acting for them) applies the AWS infrastructure and deploys the Nest backend to **dev**. The API is reachable over HTTPS. A caller can hit the **health** route successfully. The database for that environment is **Amazon RDS for PostgreSQL 18** (not Compose, not kind). Nest runs on **EKS Fargate**. Traffic reaches Nest through an **ALB** (STACK.md path). Backend images come from **ECR**.

**Why this priority**: Without a hosted API, there is no POC on AWS.

**Independent Test**: After deploy, health succeeds at the printed/dev API URL without any local Nest process.

**Acceptance Scenarios**:

1. **Given** bootstrap (`002`) exists and AWS credentials/OIDC are configured, **When** CI (or the documented first apply) deploys **dev**, **Then** Nest health succeeds on the hosted API URL.
2. **Given** that deploy, **When** an operator inspects compute, **Then** Nest is running on EKS Fargate (not EC2 node groups as the primary path).
3. **Given** that deploy, **When** they inspect the database, **Then** the app uses RDS PostgreSQL 18 for **dev** (small, cost-aware; single-AZ is acceptable).
4. **Given** a failed or partial apply, **When** the pipeline finishes, **Then** the failure is obvious (non-zero / failed job) and does not report success.

---

### User Story 2 - Hosted frontend uses the hosted API (Priority: P1)

The Next.js app for **dev** is hosted on **AWS Amplify Hosting**. It uses the generated API client against the **hosted** Nest URL (not localhost). A visitor can load the scaffold UI without running Next on their laptop.

**Why this priority**: The POC is a website plus API, not API-only.

**Independent Test**: Open the Amplify **dev** URL; scaffold page loads and health is called through the generated client against hosted Nest.

**Acceptance Scenarios**:

1. **Given** backend **dev** is healthy, **When** Amplify has built the connected **dev** branch, **Then** the hosted UI loads without a local frontend process.
2. **Given** that UI, **When** it calls health, **Then** it uses the generated client aimed at the hosted API, not an ad-hoc localhost URL.

---

### User Story 3 - CI deploys without long-lived AWS keys in git (Priority: P1)

GitHub Actions is the repeatable deploy path. It authenticates to AWS with **OIDC → IAM** (no access keys stored in the repo). On the agreed branch(es) it can: run **`pnpm verify`** on PRs (if not already from `002`), **build and push** the Nest image to ECR, **apply** CDK as needed, and **deploy** the Kubernetes overlay for **dev**. Image tags follow STACK.md (semantic-release version and/or git SHA).

**Why this priority**: Humans should not be the only way to ship, and keys must not live in git.

**Independent Test**: A PR run does not need AWS to pass verify. A deploy workflow on the **dev** path succeeds using OIDC and leaves a new image on EKS.

**Acceptance Scenarios**:

1. **Given** a pull request, **When** GitHub Actions runs verify, **Then** it can pass **without** deploying to AWS.
2. **Given** a commit on the **dev** deploy branch, **When** the deploy workflow runs, **Then** it uses GitHub OIDC (no committed AWS access keys) and updates **dev** Nest from ECR.
3. **Given** the repo, **When** someone searches git, **Then** they find no real AWS keys, RDS passwords, or `.env` with credentials.

---

### User Story 4 - Secrets and packaging match STACK.md (Priority: P1)

RDS credentials and other app secrets for **dev** live in **AWS Secrets Manager**. On EKS, **External Secrets Operator** syncs them into the Nest workload via **IRSA**. Kubernetes Secrets are not the AWS source of truth. Kustomize has a **dev** overlay (and staging/prod overlays as packaging). Helm is not introduced. While the cluster is **up**, Nest MAY scale toward zero with **KEDA**; tearing the **cluster** down is **`010`**, not this slice.

**Why this priority**: Wrong secret or packaging choices are expensive to undo later.

**Independent Test**: Nest on **dev** receives DB credentials from the Secrets Manager path; manifests are Kustomize overlays; no Helm charts.

**Acceptance Scenarios**:

1. **Given** **dev** is deployed, **When** Nest starts, **Then** it can connect to RDS using credentials that originated in Secrets Manager (not a password committed in git).
2. **Given** the repo, **When** a reviewer inspects deploy packaging, **Then** they see Kustomize base + **dev** (and staging/prod) overlays, not Helm as the primary tool.
3. **Given** this slice, **When** an operator wants to stop all hourly AWS cost, **Then** they still need **`010`** (this slice leaves **dev** **ready**, not sleeping).

---

### User Story 5 - Staging/prod shape exists without requiring a live prod (Priority: P2)

CDK stages and Kustomize overlays for **staging** and **prod** exist so promotion is the same idea as **dev**. This slice does **not** require a live **prod** (or even live **staging**) apply to be considered done. Default pipelines MUST NOT push every commit to **prod**.

**Why this priority**: Avoids surprise prod bills and irreversible first apply; still avoids a one-off snowflake **dev**.

**Independent Test**: Overlays/stages for staging and prod are in repo; only **dev** must be proven live.

**Acceptance Scenarios**:

1. **Given** the repo after this slice, **When** a reviewer looks at env packaging, **Then** **dev**, **staging**, and **prod** are distinct overlays/stages.
2. **Given** CI, **When** a normal feature branch is pushed, **Then** **prod** is not deployed.

---

### Edge Cases

- `pnpm verify` and local host / Dev Container / kind MUST still work with AWS down or missing credentials.
- This slice MUST NOT implement sign-up, organizations, bills, CSV, dashboard, or **`pnpm wake` / `sleep` / `status`**.
- First-time AWS account setup (CDK bootstrap, GitHub OIDC provider, Amplify app linked to the repo) MAY be a documented operator step; it MUST NOT require committing secrets.
- A failed image push or Kustomize apply MUST NOT be reported as a successful **dev** deploy.
- Aurora remains a **supported alternate** per STACK.md; the default created for **dev** is **RDS PostgreSQL 18**, not Aurora-only features.
- kind (`002`) stays local; this slice MUST NOT replace kind with EKS for laptop development.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: This slice MUST provision AWS infrastructure with **AWS CDK** (TypeScript preferred per STACK.md): at least **dev** EKS (Fargate), ECR, RDS PostgreSQL 18, IAM/IRSA, Secrets Manager for app secrets, ALB controller dependencies as required, and Amplify-related resources as needed. Not Terraform unless a later plan changes STACK.md.
- **FR-002**: The **first live hosted environment** MUST be **dev** (non-prod POC). A caller MUST reach Nest **health** on that environment over HTTPS.
- **FR-003**: Nest on **dev** MUST run on **EKS Fargate**. Ingress MUST use **ALB** with **IP target type** as required for Fargate (STACK.md). **KEDA** MAY scale Nest toward min 0 while the cluster exists.
- **FR-004**: Nest images MUST be built in CI, pushed to **ECR**, and pulled by EKS. Tags MUST include semantic-release version and/or git SHA per STACK.md.
- **FR-005**: The Next.js **dev** app MUST be hosted on **AWS Amplify Hosting** and MUST call the hosted Nest API through the generated client.
- **FR-006**: GitHub Actions MUST authenticate to AWS via **OIDC → IAM**. Long-lived AWS access keys MUST NOT be committed. PR **verify** MUST NOT require a successful AWS deploy.
- **FR-007**: App secrets (especially RDS credentials) MUST live in **AWS Secrets Manager**. **External Secrets Operator** MUST sync them into Nest on EKS using **IRSA**. GitHub Actions MUST NOT be the runtime source of DB passwords.
- **FR-008**: Kubernetes packaging MUST be **Kustomize** base + overlays including **dev**, **staging**, and **prod** (kind overlay remains `002`). Helm MUST NOT be introduced as the primary packager.
- **FR-009**: Staging and prod overlays/stages MUST exist. This slice MUST NOT require a live **prod** deploy. Default CI MUST NOT deploy **prod** on every branch push.
- **FR-010**: This slice MUST NOT implement PRODUCT.md domain features, and MUST NOT implement **`pnpm wake` / `sleep` / `status`** (those remain `010`). After this slice, **dev** is **ready** (running), not sleeping.
- **FR-011**: Secrets MUST NOT be committed. Document operator bootstrap (account, CDK bootstrap, OIDC, Amplify repo link) without real credentials in git.
- **FR-012**: Local paths from `002` (host + Compose Postgres, Dev Container, kind) MUST remain valid and MUST NOT depend on this AWS stack.

### Key Entities

- **Hosted environment**: named AWS env (`dev`, `staging`, `prod`); this slice proves **dev** live.
- **Deploy pipeline**: GitHub Actions jobs that verify, publish the Nest image, apply infra/workload for an environment.
- **App secret**: credential or config that must not live in git (especially RDS).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Someone with no local Nest/Next running can open the **dev** website and see the scaffold UI talk to hosted health.
- **SC-002**: Hosted **dev** API health succeeds after a CI (or documented first) deploy, using RDS as the database.
- **SC-003**: A reviewer can confirm git contains no AWS keys or database passwords.
- **SC-004**: A pull request can pass verify without any AWS account.
- **SC-005**: Staging/prod packaging exists in the repo; **prod** is not updated by ordinary feature-branch pushes.
- **SC-006**: Wake/sleep operator commands are still absent (owned by `010`).

## Assumptions

- `002-bootstrap` is implemented first (packages, health, kind overlay, verify). This slice may deploy that health-only app; it does not wait for `003`–`009`.
- STACK.md is the stack source; this spec does not reopen Nest vs Next vs pnpm.
- **dev** is the POC hosted env and may incur continuous cost until `010` exists; that is accepted for this slice.
- Exact CDK app layout, Amplify branch-per-env wiring, AWS region, and Actions workflow file layout are decided in this feature’s **plan**.
- CDK bootstrap and GitHub OIDC trust are one-time in the target AWS account; the repo documents them.
- `010-poc-wake-sleep` runs **after** this slice so sleep can tear down a stack this slice knows how to create.
- `packages/api-client` is not an AWS workload.
