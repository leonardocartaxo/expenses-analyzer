# Feature Specification: Bills and Transactions

**Feature Branch**: `006-bills-transactions`

**Created**: 2026-08-18

**Status**: Draft

**Input**: Manual bills with required establishment; credit-card transactions; parent total is local; shared establishments per docs/PRODUCT.md.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Manual utility bill (Priority: P1)

Org Admin or Bill Manager creates a water/power/gas/internet/phone bill in the current org: amount, date or period, optional note, type from the closed list, required bill-level establishment (existing or new). Viewer can read, not create. No transactions required.

**Why this priority**: Manual entry is in v1 for normal bills.

**Independent Test**: Create a water bill linked to “City Water”; Viewer cannot delete it; other org does not see it.

**Acceptance Scenarios**:

1. **Given** a Bill Manager, **When** they save a `water` bill with amount, date, and an establishment, **Then** the bill is listed in the current org.
2. **Given** bill create, **When** establishment is missing, **Then** save is rejected.
3. **Given** a Viewer, **When** they try to create, edit, or delete a bill, **Then** the product refuses.
4. **Given** org A’s bill, **When** a member works in org B, **Then** they do not see that bill.

---

### User Story 2 - Credit card bill with transactions (Priority: P1)

Member creates a `credit_card` bill with bill-level establishment (the card). They add transactions, each with amount, date, optional note, and an establishment (merchant). Merchants MAY be existing establishments. Parent **bill amount is entered on the bill** and is not auto-set to the transaction sum.

**Why this priority**: Core credit-card model.

**Independent Test**: Card bill $200 with three transactions summing to $86.40; UI shows parent 200, sum 86.40, difference; Luigi selected on two txs is one establishment.

**Acceptance Scenarios**:

1. **Given** a `credit_card` bill, **When** the member adds transactions with merchant establishments, **Then** transactions appear only under that bill, not as a top-level bill list.
2. **Given** transactions whose amounts do not sum to the parent amount, **When** the member views the bill, **Then** they see parent total, transaction sum, and the difference; the parent amount is unchanged.
3. **Given** an existing merchant establishment, **When** two transactions select it, **Then** both reference that same establishment.
4. **Given** type `water`, **When** the member saves with zero transactions, **Then** the bill is valid.

---

### User Story 3 - Edit and delete (Priority: P2)

Authorized members edit bill fields, edit/remove transactions, and delete bills subject to role.

**Why this priority**: PRODUCT.md requires list/edit/delete.

**Independent Test**: Edit note; delete a transaction; delete the bill; Viewer cannot.

**Acceptance Scenarios**:

1. **Given** a Bill Manager, **When** they edit a bill note or amount, **Then** the stored parent total is the edited amount.
2. **Given** a Bill Manager, **When** they delete a bill, **Then** its transactions are gone and it no longer lists.
3. **Given** a Viewer, **When** they attempt delete, **Then** the product refuses.

---

### Edge Cases

- Closed type list: unknown type rejected.
- Transaction without establishment rejected.
- Parent amount missing/invalid rejected; do not invent amounts (BR-024).
- Adding transactions to a non-`credit_card` bill MUST be rejected with a clear message (BR-011, BR-039).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Bills MUST belong to the current organization and MUST have a bill-level establishment (BR-001, BR-025).
- **FR-002**: v1 types MUST be the closed list: credit_card, water, power, gas, internet, phone (BR-011).
- **FR-003**: Org Admin and Bill Manager MUST create/edit/delete bills and manage transactions; Viewer MUST read only (permission matrix, BR-012, BR-022, BR-037).
- **FR-004**: `credit_card` bills MUST support multiple transactions on the parent bill (BR-034, BR-036, BR-039).
- **FR-005**: Each transaction MUST have an establishment; that establishment MAY be shared (BR-035, BR-033).
- **FR-006**: Parent bill amount MUST be stored independently; UI MUST show parent total, transaction sum, and difference when transactions exist (BR-038).
- **FR-007**: Failures MUST be plain language; no silent invented amounts (BR-024).
- **FR-008**: Each bill MUST have a **single bill date**; no period start/end (BR-049).
- **FR-009**: Types other than `credit_card` MUST NOT have transactions (BR-011, BR-039).

### Key Entities

- **Bill**: org, type, parent amount, single date, note, source=manual, establishment.
- **Transaction**: parent bill, amount, date, note, establishment.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A Bill Manager can record a utility bill with a provider establishment and no transactions.
- **SC-002**: A Bill Manager can record a card bill with several merchant transactions, including two txs sharing one merchant.
- **SC-003**: When parent total ≠ transaction sum, both figures and the difference are visible; parent total does not change by itself.
- **SC-004**: Viewer cannot mutate bills; other orgs cannot see the bills.

## Assumptions

- CSV source and import preview are slice 007; this slice is manual only (`source = manual`).
- Amounts are in the org’s primary currency (no FX).
