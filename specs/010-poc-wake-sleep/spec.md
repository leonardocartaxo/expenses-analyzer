# Feature Specification: POC Wake and Sleep

**Feature Branch**: `010-poc-wake-sleep`

**Created**: 2026-08-18

**Status**: Draft

**Input**: POC is not 24/7. Sleep stops everything that *can* be stopped. Some AWS items still bill when “stopped” (especially storage). Wake recreates/starts the stack (created by `011-aws-deploy`) and waits until ready. Operator commands, not Org Admin. Follow STACK.md.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Wake until ready (Priority: P1)

An operator who wants to use the hosted POC runs wake. The command starts or **recreates** whatever sleep tore down (RDS start, EKS cluster, ALB/NAT, Nest), **waits until the stack is usable**, then reports **ready** with URLs. Wake may take a long time after a full sleep. The operator does not guess whether Postgres or the API is up.

**Why this priority**: Sleep is aggressive; wake must be a single command back to usable.

**Independent Test**: From sleeping, run wake; it exits 0 only after API health succeeds; stdout includes ready + URLs.

**Acceptance Scenarios**:

1. **Given** the POC is sleeping (cluster gone, RDS stopped), **When** the operator runs wake, **Then** the command does not exit successfully until the API is ready to serve.
2. **Given** wake is in progress, **When** the operator watches output, **Then** they see progress (cluster vs database vs API) and finally a clear **ready** message plus app/API URLs.
3. **Given** wake fails or times out, **When** the command ends, **Then** it exits non-zero and says what is still down in plain language.

---

### User Story 2 - Sleep stops everything stoppable (Priority: P1)

When the operator is done, they run sleep. **Every stoppable billable resource** in that POC env is stopped or torn down: RDS instance stopped, Fargate/Nest not running, **EKS cluster torn down** (no control-plane charge), ALB/NAT and similar hourly network gone.

Things that **cannot** be stopped without deleting data or the product (RDS storage while the instance is stopped, container images in ECR, secrets, etc.) MAY still bill. Sleep MUST list those leftovers. The next use requires wake.

**Why this priority**: This is a POC; idle EKS/ALB/NAT must not keep running 24/7.

**Independent Test**: After sleep, RDS is not available; EKS cluster is absent; status is sleeping and names remaining billable leftovers.

**Acceptance Scenarios**:

1. **Given** a ready POC, **When** the operator runs sleep, **Then** status is sleeping, RDS is not left available 24/7, and the EKS cluster is not left running.
2. **Given** sleep completed, **When** they run status, **Then** they see a list of resources that **still bill** while sleeping (cannot be stopped without deleting data).
3. **Given** already sleeping, **When** they run sleep again, **Then** the command is safe (idempotent) and still reports sleeping.

---

### User Story 3 - Status without starting (Priority: P2)

The operator can ask sleeping / waking / ready, what is still down, and what still bills while asleep, without starting wake.

**Why this priority**: Avoids waiting on a long wake just to check.

**Independent Test**: status while sleeping vs ready returns distinct outcomes; sleeping includes leftover-cost notes.

**Acceptance Scenarios**:

1. **Given** a sleeping stack, **When** they run status, **Then** it reports sleeping without starting RDS or creating a cluster.
2. **Given** a ready stack, **When** they run status, **Then** it reports ready.

---

### Edge Cases

- Hitting the website before wake finishes: UI/API show waking or retry; do not pretend bills loaded.
- Local host / Dev Container (Compose Postgres) and **kind** (`pnpm local:up`) do not use this slice (no RDS wake).
- Application admin / Org Admin MUST NOT be required to run wake.
- Amplify Hosting has **no stop**: it MAY remain as a leftover bill while sleeping unless the plan finds a supported pause/delete-and-recreate. Status MUST mention it if it still bills.
- Deleting the RDS instance (not stop) is **out of scope** for sleep — data must survive.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Operators MUST have **`pnpm wake`** that brings the POC stack back and **blocks until ready**, then prints ready + URLs (STACK.md).
- **FR-002**: Operators MUST have **`pnpm sleep`** that stops or tears down **all stoppable billable resources** in the POC env: at least stop RDS, stop Nest/Fargate, **tear down EKS** (no control-plane billing), and remove ALB/NAT (and equivalent hourly network) for that env.
- **FR-003**: Operators MUST have **`pnpm status`** reporting sleeping / waking / ready, what is still down, and **what still bills while sleeping**.
- **FR-004**: Ready MUST include Postgres available and API health succeeding. Wake MUST NOT report ready early.
- **FR-005**: Clients/UI MUST show waking or retry if requests arrive before ready.
- **FR-006**: Wake/sleep MUST target non-prod by default.
- **FR-007**: This slice MUST NOT grant Org Admin or Bill Manager cloud start/stop powers.
- **FR-008**: Sleep MUST NOT destroy durable data (RDS data volume / snapshots needed to wake). Images in ECR and secrets MAY remain.

### Key Entities

- **POC stack state**: sleeping | waking | ready.
- **Stoppable resource**: billed only while running or while the AWS resource exists in an hourly/on state (cluster, ALB, NAT, RDS *instance*, Fargate).
- **Leftover bill**: still charged after sleep (RDS storage, ECR, Secrets Manager, Amplify if not paused, log retention).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Operator goes from sleeping to using the app with wake + ready message only.
- **SC-002**: After sleep, RDS is not 24/7-available **and** EKS control plane is not left running for that POC env.
- **SC-003**: Status while sleeping lists leftover bills; operator can explain what they still pay.
- **SC-004**: Wake timeout or failure is obvious (non-zero exit, plain-language cause).

## Assumptions

- Exact CDK/Kustomize recreate steps are in this feature’s plan; this spec locks **maximize stoppable savings**.
- Wake after cluster teardown is slow; that is accepted for a POC.
- Bootstrap (`002`) stays first; **`011-aws-deploy`** creates the hosted stack; this slice is **after** that (do not block `002` or `011`).
