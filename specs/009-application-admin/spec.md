# Feature Specification: Application Admin

**Feature Branch**: `009-application-admin`

**Created**: 2026-08-18

**Status**: Draft

**Input**: First application admin via migration/seed script; platform management of users, organizations, settings; no break-glass into org bills per docs/PRODUCT.md.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Seed the first application admin (Priority: P1)

Operators run a migration or seed script that creates the first application admin account. Ordinary sign-up does not produce an application admin.

**Why this priority**: Platform operations cannot start from a random self-signup.

**Independent Test**: Fresh database + seed script yields one app admin; a normal sign-up user is not app admin.

**Acceptance Scenarios**:

1. **Given** an empty user table, **When** the migration/seed script runs with the configured admin email, **Then** that user exists with application admin.
2. **Given** a person who only signed up through the app, **When** they inspect their platform role, **Then** they are not application admin.
3. **Given** the seed script run twice, **When** the admin already exists, **Then** it does not create a second admin for the same email (idempotent).

---

### User Story 2 - Manage users and organizations on the platform (Priority: P2)

An application admin can list users and organizations, and operate platform settings. They cannot open an organization’s bills unless they are also a member of that org.

**Why this priority**: PRODUCT.md platform role.

**Independent Test**: App admin lists orgs; opening org D bills is denied without membership; after being invited as Viewer they can read bills via normal membership.

**Acceptance Scenarios**:

1. **Given** an application admin, **When** they open platform user/org lists, **Then** they see platform records (users, org names/types) as specified in the plan.
2. **Given** an application admin who is not a member of org D, **When** they request org D’s bills, dashboard, or establishments, **Then** access is denied.
3. **Given** the same admin later accepts an Org Admin invite to org D, **When** they work in org D, **Then** they have org permissions from membership, not from the platform role.

---

### Edge Cases

- Seed credentials MUST NOT be committed; email/password come from local/env configuration.
- Removing the last application admin MUST be refused until another application admin exists.
- Application admin cannot bypass Org Admin invite to become a member.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The first application admin MUST be created by a migration/seed script (BR-044).
- **FR-002**: Sign-up MUST NOT grant application admin (BR-045, BR-010).
- **FR-003**: Application admin MUST be able to manage platform concerns: users, organizations, and settings. Settings in v1 MAY be a stub (e.g. read-only app metadata) (BR-044).
- **FR-004**: Application admin MUST NOT grant automatic access to an organization’s bills, transactions, establishments, or dashboard (BR-010).
- **FR-005**: Seed/migration MUST be idempotent for the first admin email.
- **FR-006**: An existing application admin MAY grant application admin to another user (BR-044).
- **FR-007**: The product MUST refuse removing the last application admin.

### Key Entities

- **Application admin**: platform flag/role on a user, seeded not self-served.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: After seed, operators can sign in as the application admin without using the public sign-up path to obtain that role.
- **SC-002**: App admin without membership cannot read any org’s bills in tests.
- **SC-003**: A normal sign-up never becomes application admin.

## Assumptions

- Seed admin email/password are provided via environment for local/dev; production secret source is STACK.md (Secrets Manager) in a later deploy slice.
