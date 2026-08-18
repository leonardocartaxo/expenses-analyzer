# Feature Specification: Dashboard

**Feature Branch**: `008-dashboard`

**Created**: 2026-08-18

**Status**: Draft

**Input**: Month-by-month costs and establishment-tag breakdown in the organization timezone per docs/PRODUCT.md.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Month totals for the current org (Priority: P1)

A member with view access opens the dashboard, selects a month, and sees overall costs for that month in the org’s primary currency. Months are bounded using the **organization timezone**. Another org’s bills never appear.

**Why this priority**: Primary analysis surface.

**Independent Test**: Bills in January (org TZ) total X; dashboard January shows X; switching org changes the figures.

**Acceptance Scenarios**:

1. **Given** bills in the current org, **When** the member opens the dashboard for a month, **Then** they see that month’s overall cost.
2. **Given** an org timezone other than UTC, **When** a bill’s timestamp falls on a UTC month boundary, **Then** the bill is counted in the month implied by the org timezone.
3. **Given** org A data, **When** the member switches to org B, **Then** dashboard figures are only org B.

---

### User Story 2 - Breakdown by establishment tags (Priority: P1)

For the selected month, the member sees costs grouped by establishment tags (food, transport, leisure, …) so they can judge whether a category is high.

**Why this priority**: PRODUCT.md success criterion for analysis.

**Independent Test**: Tagged merchant transactions in the month appear under those tags; untagged spend is visible as untagged/other, not silently dropped.

**Acceptance Scenarios**:

1. **Given** transactions whose establishments have tag `food`, **When** the member views that month, **Then** the food category includes those amounts.
2. **Given** a utility bill with no transactions and establishment tag `utilities`, **When** they view that month, **Then** the utilities category includes that bill amount.
3. **Given** a Viewer, **When** they open the dashboard, **Then** they can see totals and tag breakdown but cannot change tags (tags are changed in establishments/bills slices).

---

### Edge Cases

- Month with no bills: show zero, not an error.
- Parent total ≠ transaction sum: show overall (parent-based) and tag breakdown without inventing a third silent total (BR-051).
- Establishment with multiple tags: amount appears in **each** tag in the breakdown; overall month cost counts it once.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Members with view access MUST have a dashboard for the current org showing expenses month by month (BR-026, BR-027, BR-030).
- **FR-002**: Month boundaries MUST use the organization timezone (BR-043).
- **FR-003**: Members MUST see costs categorized by establishment tags for the selected month (BR-028).
- **FR-004**: Dashboard MUST use the org’s single primary currency (BR-031).
- **FR-005**: Viewer MUST be able to read the dashboard; they MUST NOT need Bill Manager rights.
- **FR-006**: Attribution MUST follow BR-051: overall = parent bill amounts; tags from transactions when present, otherwise the bill; multi-tag in each category; untagged bucket; do not hide parent vs tag-sum differences.

### Key Entities

- **Month summary**: org, year-month (org TZ), overall cost, per-tag costs.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A member can select a month and see an overall cost that only includes the current org.
- **SC-002**: Tag breakdown is visible for the same month so a category can be compared to the overall.
- **SC-003**: Changing org timezone (if allowed in orgs slice) changes which month a boundary bill falls into.
- **SC-004**: Viewer can complete the dashboard view without edit permissions.

## Assumptions

- Untagged spend appears in an **untagged** bucket, not omitted (also required by BR-051).
