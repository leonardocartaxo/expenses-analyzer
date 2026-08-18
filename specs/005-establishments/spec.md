# Feature Specification: Establishments

**Feature Branch**: `005-establishments`

**Created**: 2026-08-18

**Status**: Draft

**Input**: Organization-scoped establishments with name, free-text address, closed tags, and reuse by selection per docs/PRODUCT.md.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create and list establishments (Priority: P1)

Org Admin or Bill Manager creates an establishment in the current org with a name, optional free-text address, and tags from the closed v1 list. Viewer can list/search but not create.

**Why this priority**: Every bill needs an establishment.

**Independent Test**: Create “Visa Rewards” with tag credit_card; it appears only in the current org.

**Acceptance Scenarios**:

1. **Given** a Bill Manager in org A, **When** they create an establishment with name and tag `utilities`, **Then** it is listed in org A.
2. **Given** a Viewer, **When** they try to create or edit an establishment, **Then** the product refuses.
3. **Given** an establishment in org A, **When** a member of org B lists establishments, **Then** they do not see org A’s establishment.

---

### User Story 2 - Search by address and change tags (Priority: P2)

Members with view access search establishments by address substring. Org Admin and Bill Manager change tags; Viewer cannot.

**Why this priority**: Dashboard categorization depends on tags; address search is in PRODUCT.md.

**Independent Test**: Establishment at “10 Pasta Lane” is found by query “Pasta”; Viewer cannot retag it.

**Acceptance Scenarios**:

1. **Given** an establishment with address “10 Pasta Lane”, **When** a member searches “pasta”, **Then** that establishment is in the results (case-insensitive substring).
2. **Given** a Bill Manager, **When** they set tags to `food`, **Then** those tags are stored and invalid tags are rejected.
3. **Given** a Viewer, **When** they try to change tags, **Then** the product refuses.

---

### User Story 3 - Share by selection (Priority: P1)

Members pick an existing establishment when they mean the same place. Duplicate names are allowed. Sharing is the same establishment record, not automatic merge by name.

**Why this priority**: Credit-card merchants and utilities are reused.

**Independent Test**: Two establishments named “Luigi” can exist; picking one id uses that one.

**Acceptance Scenarios**:

1. **Given** an existing establishment, **When** a member selects it again for another use, **Then** they are attaching to the same establishment, not creating a new one.
2. **Given** two establishments with the same name, **When** a member lists them, **Then** both appear as distinct records (address/tags can differ).

---

### Edge Cases

- Empty name rejected.
- Tags outside the closed list rejected in plain language.
- Deleting an establishment that any bill or transaction still references MUST be blocked with a plain-language message (BR-050). Unused establishments MAY be deleted by Org Admin or Bill Manager.
- Multiple tags on one establishment allowed.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Establishments MUST belong to an organization (no global shared directory across orgs).
- **FR-002**: Establishments MUST have a name; address is optional free text; tags MUST be from the closed v1 list: food, transport, leisure, utilities, credit_card, housing, health, shopping (BR-032).
- **FR-003**: Address search MUST support substring match (BR-032).
- **FR-004**: Org Admin and Bill Manager MUST create/edit/delete (when unused) and change tags; Viewer MUST read only (permission matrix, BR-029).
- **FR-005**: Many bills/transactions MUST be able to reference the same establishment (BR-033). This slice MUST expose select-existing as well as create-new.
- **FR-006**: Duplicate names MUST be allowed; the product MUST NOT auto-merge by name.
- **FR-007**: The product MUST block deleting an establishment that is still referenced by any bill or transaction (BR-050).

### Key Entities

- **Establishment**: org, name, free-text address, tags[] from closed list.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A Bill Manager can create an establishment and find it by a fragment of its address.
- **SC-002**: Viewer cannot change tags.
- **SC-003**: Closed-list tag values outside the list never persist.
- **SC-004**: Org isolation: 0 establishments leak across orgs in listing/search.

## Assumptions

- Bills/transactions that attach establishments are specified in 006; this slice can ship CRUD/search first. Delete-while-in-use MUST still be enforced once 006 exists (and MAY be a no-op until then because nothing references establishments yet).
