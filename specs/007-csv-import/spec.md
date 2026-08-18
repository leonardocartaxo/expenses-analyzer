# Feature Specification: CSV Import

**Feature Branch**: `007-csv-import`

**Created**: 2026-08-18

**Status**: Draft

**Input**: Template-based CSV import with editable preview and accept/reject; ship expenses-card-statement-v1 and expenses-utility-bill-v1 (one row = one bill) per docs/PRODUCT.md and docs/csv-templates/.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Import a card statement through preview (Priority: P1)

Org Admin or Bill Manager chooses template `expenses-card-statement-v1`, uploads a valid CSV, sees a preview of proposed bills, card establishments, merchants, and transactions, edits if needed, then accepts. Lasting org data appears only after accept.

**Why this priority**: CSV is in v1; preview gate is non-negotiable.

**Independent Test**: Upload the example file; preview two Visa bills; Luigi shared; January mismatch visible; accept writes data; a second run that rejects leaves data unchanged.

**Acceptance Scenarios**:

1. **Given** a Bill Manager, **When** they upload a valid `expenses-card-statement-v1` file, **Then** they see a preview and the org’s lasting bills are unchanged until accept.
2. **Given** preview, **When** they accept, **Then** proposed bills/transactions/establishments are committed to the current org (source = csv) according to the current preview.
3. **Given** preview, **When** they reject or cancel, **Then** no lasting changes from that upload remain.
4. **Given** the example file’s January group, **When** they view preview, **Then** parent total 200.00, transaction sum 86.40, and the difference are visible.

---

### User Story 2 - Import utility bills, one row each (Priority: P1)

Org Admin or Bill Manager chooses `expenses-utility-bill-v1`, uploads a valid CSV, and previews one bill per row with no transactions. Shared providers reuse one establishment. Accept commits; reject does not.

**Why this priority**: PRODUCT.md requires CSV for normal bills in v1, not cards only.

**Independent Test**: Upload the utility example; three bills; City Water shared; no transactions.

**Acceptance Scenarios**:

1. **Given** a Bill Manager, **When** they upload a valid `expenses-utility-bill-v1` file, **Then** preview shows one bill per row and zero transactions.
2. **Given** that preview, **When** they accept, **Then** the org has those bills with bill-level establishments and no transactions.
3. **Given** a row with `bill_type=credit_card`, **When** they import with this template, **Then** parse fails with a clear message.

---

### User Story 3 - Correct rows in preview (Priority: P1)

In preview, the member can correct fields and drop proposed rows so the committed result matches intent.

**Why this priority**: PRODUCT.md requires editable preview.

**Independent Test**: Change a merchant name in preview; accept; stored establishment/transaction uses the corrected name.

**Acceptance Scenarios**:

1. **Given** a preview, **When** the member edits a proposed amount or name, **Then** accept commits the edited values, not the original file values for that field.
2. **Given** a preview, **When** the member drops a proposed transaction row, **Then** accept does not create that transaction.

---

### User Story 4 - Unsupported files fail clearly (Priority: P2)

PDF/Word and unreadable or wrong-template CSVs fail with a plain-language message and no partial lasting commit.

**Why this priority**: BR-020, BR-021.

**Independent Test**: Upload a `.pdf`; rejected. Upload CSV with missing columns; rejected. Org bills unchanged.

**Acceptance Scenarios**:

1. **Given** a non-CSV upload, **When** the member imports, **Then** the product rejects it as unsupported.
2. **Given** a CSV that fails `expenses-card-statement-v1` parse rules, **When** they import with that template, **Then** they see a clear error and no lasting bill changes.
3. **Given** a Viewer, **When** they try to import, **Then** the product refuses.

---

### Edge Cases

- Zero data rows: fail, no preview commit.
- Invalid tag in `merchant_tags`: fail with the offending value named.
- Shared merchant across bills in one file: preview one establishment reused.
- Member can map merchants to an establishment that already exists in the org during preview.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Import MUST use named CSV templates; v1 MUST include `expenses-card-statement-v1` and `expenses-utility-bill-v1` as specified under [`docs/csv-templates/`](../../docs/csv-templates/) (BR-014, BR-046).
- **FR-002**: Member MUST choose a template (suggestion optional) (BR-015).
- **FR-003**: After successful parse, product MUST show editable preview of bills, establishments, and transactions (BR-016, BR-017).
- **FR-004**: Accept MUST commit current preview; reject MUST leave lasting bill data unchanged (BR-018, BR-019).
- **FR-005**: Parse/template failure MUST be plain language with no silent partial lasting commit (BR-020, BR-024).
- **FR-006**: Non-CSV including PDF/Word MUST be rejected in v1 (BR-021).
- **FR-007**: Only Org Admin and Bill Manager MAY import (permission matrix).
- **FR-008**: Import MUST be scoped to the current organization (BR-006).
- **FR-009**: Parent totals in preview MUST follow BR-038 (do not auto-overwrite to transaction sum).
- **FR-010**: `expenses-utility-bill-v1` MUST treat each data row as one bill with zero transactions (BR-046, BR-049).

### Key Entities

- **CSV template**: key, display name, parse rules, column mapping.
- **Import preview**: draft proposed bills, establishments, transactions; not lasting until accept.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The shipped card example CSV can be previewed and accepted into a test org with two card bills and a shared Luigi establishment.
- **SC-002**: Rejecting preview results in zero bills created from that attempt.
- **SC-003**: Invalid files never create lasting bills.
- **SC-004**: Viewer cannot import.
- **SC-005**: The shipped utility example CSV can be previewed and accepted as three bills, zero transactions, with City Water shared.

## Assumptions

- Where the uploaded file and preview live (memory, temporary table, or object storage such as S3) is **deferred** until this slice is planned. Product rule stays: no lasting domain commit until accept.
- LLM mapping is out of scope.
