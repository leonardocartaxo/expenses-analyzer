# Feature Specification: Organizations and Membership

**Feature Branch**: `004-organizations`

**Created**: 2026-08-18

**Status**: Draft

**Input**: Create personal/company orgs, Org Admin invites, roles, leave/delete, last-admin rule per docs/PRODUCT.md.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create an organization (Priority: P1)

An authenticated user creates a personal or company organization with name, free-text address, a **chosen currency** (required, no default), optional timezone (default UTC), and becomes Org Admin. They can switch into that org as the working context.

**Why this priority**: Bills cannot exist without an org.

**Independent Test**: Sign in, create a personal org, appear as Org Admin, see empty bill context for that org only.

**Acceptance Scenarios**:

1. **Given** an authenticated user, **When** they create a personal or company org with name, address, and a chosen currency, **Then** the org exists and they are Org Admin.
2. **Given** create org, **When** they omit timezone, **Then** timezone is UTC.
3. **Given** create org, **When** they omit currency, **Then** creation is rejected.
4. **Given** two orgs, **When** they switch working organization, **Then** membership and later bill data are isolated to the selected org.
5. **Given** create company org, **When** they are not an application admin, **Then** creation still succeeds (no app-admin approval).

---

### User Story 2 - Org Admin invites a member (Priority: P1)

An Org Admin invites an email, assigning one or more of Org Admin, Bill Manager, Viewer. The invitee signs in (or signs up) and accepts. Bill Manager and Viewer cannot invite.

**Why this priority**: Only join path for existing orgs.

**Independent Test**: Org Admin invites a second user as Viewer; after accept, Viewer can switch into that org; Bill Manager cannot send an invite.

**Acceptance Scenarios**:

1. **Given** an Org Admin, **When** they invite an email with role Viewer, **Then** an invitation exists for that org.
2. **Given** a pending invite, **When** the invitee accepts while authenticated as that email, **Then** they gain membership with the invited roles.
3. **Given** a Bill Manager or Viewer, **When** they try to invite, **Then** the product refuses.
4. **Given** no invite, **When** a user tries to join an existing org, **Then** they cannot.
5. **Given** a pending invite, **When** an Org Admin cancels it, **Then** the invitee cannot accept it.

---

### User Story 3 - Roles, leave, last admin, delete (Priority: P2)

Org Admin changes roles and can remove members. Members can leave except the last Org Admin. Org Admin can delete the org.

**Why this priority**: Completes org lifecycle so later slices do not invent it.

**Independent Test**: Two Org Admins: one leaves successfully. One Org Admin: leave is blocked until they promote someone or delete.

**Acceptance Scenarios**:

1. **Given** a member with Viewer only, **When** Org Admin also grants Bill Manager, **Then** permissions are the union of both roles.
2. **Given** two Org Admins, **When** one leaves, **Then** they lose membership and the org remains.
3. **Given** a single Org Admin, **When** they try to leave or drop Org Admin, **Then** the product refuses until they assign Org Admin to someone else or delete the org.
4. **Given** an Org Admin, **When** they delete the org, **Then** that org’s memberships are gone and its data is not accessible.
5. **Given** an Org Admin, **When** they change name, address, timezone, or currency, **Then** the org profile updates and existing bill amounts are not converted.
6. **Given** a Bill Manager or Viewer, **When** they try to change org name, address, timezone, or currency, **Then** the product refuses.

---

### Edge Cases

- Invite to an email that already has membership in that org: reject or no-op with a clear message.
- User belongs to many orgs; switching must not mix data.
- Deleting org in this slice may have no bills yet; later slices MUST cascade bills/establishments (PRODUCT.md BR-042).
- Application admin without membership MUST NOT see the org as a working bill context.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Authenticated users MUST be able to create orgs of type personal or company (BR-002, BR-005).
- **FR-002**: Orgs MUST have name, free-text address, timezone (default UTC), and a currency **chosen at create with no default** (BR-003, BR-005, BR-031, BR-043).
- **FR-003**: Creator MUST become Org Admin (BR-005).
- **FR-004**: Users MUST join existing orgs only via Org Admin invitation; invites remain until accept or Org Admin cancel (BR-040).
- **FR-005**: Roles MUST be Org Admin, Bill Manager, Viewer; multiple roles union (BR-008, BR-009) and MUST match the PRODUCT.md permission matrix.
- **FR-006**: Users MUST belong to many orgs and switch working org (BR-007).
- **FR-007**: Last Org Admin MUST NOT abandon the org (BR-041).
- **FR-008**: Org Admin MUST be able to delete the org (BR-042).
- **FR-009**: Org data MUST be isolated (BR-006).
- **FR-010**: Application admin MUST NOT imply org membership (BR-010).
- **FR-011**: Only Org Admin MAY edit org name, address, timezone, or currency; changing currency MUST NOT convert amounts (BR-047).

### Key Entities

- **Organization**: type, name, address, timezone, currency.
- **Membership**: user, org, roles[].
- **Invitation**: org, email, roles[], invited by Org Admin.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can create a personal org and work in it as Org Admin without a second account.
- **SC-002**: A second person can enter that org only after an Org Admin invite + accept.
- **SC-003**: Bill Manager cannot invite; Viewer cannot invite.
- **SC-004**: Last Org Admin cannot leave; promoting another Org Admin unblocks leave.
- **SC-005**: Switching orgs never shows the other org’s membership list as the current context.

## Assumptions

- Currency code list is ISO 4217 (exact set in the plan, e.g. USD, EUR, BRL) without introducing FX.
- Email of invitee MUST match the signed-in account on accept.
