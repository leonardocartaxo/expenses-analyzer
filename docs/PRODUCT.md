# Expenses Analyzer — Product

**Status**: Draft  
**Updated**: 2026-08-10

This is the **only** source of truth for global product intent: what the app is, who it is for, how it works, and normative business rules. Feature specs under `specs/` must not contradict this document; they add delivery detail for one slice at a time.

Technical stack: [`STACK.md`](../STACK.md). Process / SDD: [`.specify/memory/constitution.md`](../.specify/memory/constitution.md).

---

## What this app is

**Expenses Analyzer** is a product for capturing and understanding **bills and spending** inside **organizations**.

There is **no** separate “user-owned bills” bucket. A user who wants personal finances **creates (or joins) an organization** of type **personal**. Companies and other kinds of groups use other organization types (for example **company**). All bills belong to an **organization**, never directly to a user account.

People record bills such as **credit cards, water, power**, and similar. **Manual bill entry is in v1** (required). **CSV import is in v1 for all supported bill kinds**—not only credit cards: normal bills (water, power, and similar) can be imported from CSV as well as credit card statements (PDF/Word import later). **Every bill is an association with an establishment**—including a **credit card bill** (the card/issuer is an establishment, just as a water company is). Credit card bills additionally contain **transactions**, each linked to further establishments (e.g. Luigi Restaurant, Max Cinemas)—searchable by address and tagged by type (food, transport, leisure, …).  

Members use a **dashboard** to see expenses **month by month**, analyze costs, and break spending down **by establishment tags** (food, transport, leisure, …) so they can tell whether they are spending too much in a category in a given month. The product helps the right members of each organization **see, correct, and summarize** that bill data—without mixing one organization’s data with another’s.

This document is **not** a technical plan and **not** a sprint backlog of user stories.

---

## Who it is for

| Audience | Relationship to the product |
| --- | --- |
| **User (account)** | Signs in; has **name** and **address**. Has **no** personal bill list of their own. Accesses bills only through **organization membership**. |
| **Organization member** | Belongs to one or more **organizations**; works on that org’s bills according to their **roles**. |
| **Org Admin** | Organization role: manages membership/roles in that org; full bill access in that org. |
| **Bill Manager** | Organization role: creates, imports, edits, deletes bills; does not manage membership. |
| **Viewer** | Organization role: reads bills and summaries only. |
| **Application admin** | Platform role: manages the **application** (users, organizations, settings). Does **not** automatically see an org’s bills unless they are also a member. |

**Rules of belonging**

- An **organization** has **many users** (a personal-type org may often have one member, but it is still an organization).  
- A **user** can belong to **many organizations** (e.g. one personal org and several company orgs).  
- In one organization, a user can hold **one or more** of: Org Admin, Bill Manager, Viewer.  
- Permissions in an org are the **union** of that user’s roles there.

**Organization types (v1)**

| Type | Intent |
| --- | --- |
| **personal** | Org created for an individual’s (or household-style) finances; still membership + roles + bills like any org. |
| **company** | Org for a business / formal organization. |

Organizations have a **name** and **address** (in addition to type).

Additional types may be added later without changing the rule that **all bills live on an organization**.

---

## How it works (conceptual)

```text
Sign in
   → choose an organization (personal-type, company-type, …)
        → manage that org’s bills only
             • add manually
             • import CSV → preview/edit → accept or reject
             • list / edit / delete (if role allows)
             • dashboard: month-by-month costs + breakdown by tags
```

1. **Organizations isolate data**  
   Org A’s bills never appear in Org B. Switching organization changes the entire working context. Users do not have a parallel personal bill store outside an org.

2. **Creating a personal org**  
   A user who wants “my bills” creates an organization with type **personal** (and typically becomes Org Admin there). That is how personal use works—not a separate product mode with user-owned bills.

3. **Bills are the core record**  
   A bill has at least: amount, date or billing period, type (e.g. credit card, water, power), optional note, and an **establishment**. Source is manual or CSV import. Every bill belongs to exactly one organization.

4. **Establishments on every bill; transactions add more establishments**  
   An **establishment** is a named party or place (utility provider, **credit card / issuer**, restaurant, cinema, etc.). It can have an **address** (searchable) and **category tags** such as food, transport, leisure, utilities, credit card (and similar). Establishments can be reused.

   - **Every bill**—including **credit card**, water, and power—MUST be associated with an establishment at bill level (the credit card itself / issuer is that establishment for a card bill; the water company for a water bill).  
   - **Credit card** (and similar) bills ALSO contain multiple **transactions** (charges on the statement). Each transaction is tied to an establishment (e.g. Luigi Restaurant).  
   - So a credit card bill has **one bill-level establishment** (the card) **plus** many transaction-level establishments (where the card was used).  
   - Transactions live **on the parent bill**; they are not a separate top-level bill list.  
   - Simple bills may have **zero transactions** but still have a bill-level establishment.

5. **Manual entry**  
   Authorized members (per role) type bill details (including the bill’s establishment) and, when applicable, transactions + their establishments, and save them into the **current organization**.

6. **CSV import (v1) via templates — preview before save**  
   There is **no single industry-standard bill CSV**. Issuers and utilities use different columns, separators, dates, and amount conventions.  
   Authorized members upload a CSV and select (or the product matches) a **CSV template** (profile) for that source—e.g. a **credit card** issuer export **or** a **normal bill** export (water, power, and similar).  
   CSV import is **not limited to credit cards**; templates cover both statement-style (often with transactions) and simple bill exports (bill + establishment, typically no transactions).  
   **Nothing is written to lasting bill data until the member accepts.** After upload, the app MUST show a **preview** of the proposed bills, establishments, and transactions. In preview, the member MUST be able to **correct** fields (and related proposed rows) before deciding.  
   The member then **accepts** (commit the previewed result) or **rejects** (discard the import; no lasting changes from that upload).  
   Unknown/unsupported templates or unreadable files fail clearly. **PDF and Word are not in v1.**  
   New templates can be added over time without changing the “CSV → preview/edit → accept or reject” product rule.

7. **Dashboard and analysis**  
   In the current organization, members use a **dashboard** to view expenses **month by month**, inspect overall costs for each month, and break costs down **by establishment tags** (e.g. food, transport, leisure, utilities). That is how they see what kinds of expenses they have and whether a category is too high in a given month.  
   Tagging establishments (and keeping tags accurate) is part of enabling this analysis. Amounts that roll into tag totals come from the relevant spend (transactions when present; otherwise the bill tied to that establishment—product rules for attribution are refined in later specs if needed, but the intent is clear: **monthly view + tag breakdown**).

8. **Platform administration**  
   Application admins operate the product (users, orgs, settings). They are not a backdoor into every org’s bill list.

---

## Business rules (normative)

### Organizations and tenancy

- **BR-001**: All bills MUST belong to an **organization**. Users MUST NOT have a separate personal bill collection outside an organization.  
- **BR-002**: Organizations MUST have a **type**. v1 types include at least **personal** and **company**.  
- **BR-003**: Organizations MUST have a **name** and an **address**.  
- **BR-004**: Users MUST have a **name** and an **address**.  
- **BR-005**: Users MUST be able to **create** an organization (including type **personal**) and join organizations they are invited to or otherwise allowed to join.  
- **BR-006**: Bill and membership data MUST be isolated per organization (no cross-org leakage).  
- **BR-007**: A user MUST be able to belong to multiple organizations and switch among organizations they may access.

### Roles

- **BR-008**: v1 organization roles are exactly **Org Admin**, **Bill Manager**, and **Viewer**, with meanings in “Who it is for”.  
- **BR-009**: Multiple roles for one user in one org MUST combine as a **union** of permissions.  
- **BR-010**: **Application admin** manages platform concerns only; it MUST NOT grant automatic access to an organization’s bills without membership.

### Bills

- **BR-011**: Organizations MUST support many bills, including types such as credit card, water, and power.  
- **BR-012**: Authorized members MUST be able to create bills **manually** in the current organization. Manual entry is **in scope for v1** (not deferred).  
- **BR-013**: Authorized members MUST be able to create bills via **CSV import** (**also in v1**). CSV import MUST support **credit card bills and normal bills** (e.g. water, power)—not credit cards only.  
- **BR-014**: CSV import MUST use **CSV templates** (named profiles). The product MUST NOT assume one universal bill CSV layout; each template defines how a known export format maps into bills, establishments, and (when applicable) transactions. Templates MAY target credit-card exports or normal-bill exports.  
- **BR-015**: Members MUST be able to choose a CSV template when importing (or accept an automatically suggested template when detection is available).  
- **BR-016**: After a successful parse, the product MUST present an **import preview** of all proposed changes (bills, establishments, transactions as applicable) **before** any of those changes are committed as lasting organization data.  
- **BR-017**: In preview, members MUST be able to **edit/correct** proposed values (and adjust or drop proposed rows as the UI allows) so the preview matches what they intend to import.  
- **BR-018**: Members MUST be able to **accept** the preview (commit the current previewed result) or **reject** it (cancel the import). Rejecting MUST leave organization bill data unchanged by that import attempt.  
- **BR-019**: Until accept, the import MUST NOT permanently create or update bills/transactions/establishments from that file (draft/preview state only).  
- **BR-020**: Import from a CSV that matches no supported template, or fails to parse under the chosen template, MUST fail with a clear message—no silent partial commit of garbage data.  
- **BR-021**: In v1, non-CSV uploads (including PDF and Word) MUST be rejected as unsupported.  
- **BR-022**: Authorized members MUST be able to list, edit, and delete bills subject to role.  
- **BR-023**: Users MUST authenticate before accessing organization data.  
- **BR-024**: Validation and parse failures MUST be explained in plain language; the system MUST NOT silently invent amounts.  
- **BR-025**: Every bill MUST be associated with an **establishment**, including **credit card** bills (the card/issuer is the bill-level establishment) and “normal” bills such as water and power (the provider is the establishment).

### Dashboard and tag analysis

- **BR-026**: Members with view access MUST have a **dashboard** for the current organization that shows expenses **month by month**.  
- **BR-027**: For a selected month, the dashboard MUST show overall costs so members can analyze what was spent.  
- **BR-028**: Members MUST be able to see costs **categorized by establishment tags** (e.g. food, transport, leisure, utilities) for a selected month, so they can see what kinds of expenses they have and whether a category is high.  
- **BR-029**: Members MUST be able to assign and change **establishment tags** (subject to role) so dashboard categorization stays accurate.  
- **BR-030**: Dashboard figures MUST respect organization isolation (only the current org’s data).

### Transactions and establishments

- **BR-031**: An **establishment** MUST have a name; it MUST support an **address** that members can search; it MUST support **category tags** (v1 examples: food, transport, leisure, utilities; more tags may be added later).  
- **BR-032**: Establishments MAY be reused across many bills and transactions; an establishment is not itself a bill.  
- **BR-033**: Bills of types that carry line detail (at least **credit card** in v1) MUST support **transactions**: discrete charges **on that parent bill**.  
- **BR-034**: Each transaction MUST be linkable to an **establishment** (e.g. Luigi Restaurant, Max Cinemas).  
- **BR-035**: Transactions MUST NOT appear as a detached top-level bill list; they belong to their parent bill.  
- **BR-036**: Authorized members MUST be able to add, edit, and remove transactions on a parent bill, and create/update establishment details, subject to the same bill permissions as editing that bill.  
- **BR-037**: When transactions exist, the parent bill’s meaningful total MUST stay consistent with its transactions (sum of transaction amounts matches the bill total, or the product clearly shows how parent total and transactions relate—no silent mismatch).  
- **BR-038**: Bill types that do not need line detail (e.g. many utility bills) MAY have zero transactions, but MUST still have a bill-level establishment; credit card bills in v1 MUST allow multiple transactions with establishments.

---

## Usage overview *(lightweight — not delivery user stories)*

Illustrative flows for shared understanding. Detailed user stories belong in later feature specs.

1. **Personal use** — User creates a **personal**-type organization, becomes Org Admin, adds bills there.  
2. **Company use** — User works in a **company**-type org with other members and roles.  
3. **Viewer is restricted** — Can read that org’s bills; cannot delete or import.  
4. **Same person, two orgs** — Switches from personal org to company org; data set changes completely.  
5. **App admin** — Manages users/orgs on the platform; still cannot open Org D’s bills without membership.  
6. **CSV import** — Chooses a **CSV template** (card or normal bill), uploads file → **previews** proposed changes → **edits corrections in preview** → **accepts** (commit) or **rejects** (discard; no lasting import).  
7. **Dashboard** — Views month-by-month totals and breakdown by tags (food, transport, leisure, …) to judge category spend.  
8. **Credit card bill** — Bill-level **establishment** = the credit card/issuer; plus **transactions**, each with their own establishment (restaurant, cinema, …).  
9. **Utility bill** — Bill-level **establishment** = the provider; typically no transactions.

---

## Success criteria (product-level)

- **SC-001**: Stakeholders can explain that bills live only on organizations (including personal-type orgs), not on user accounts.  
- **SC-002**: Acceptance tests for later features can map cleanly to BR-001–BR-038 (no contradictory tenancy, role, bill, import-preview, dashboard, transaction, or establishment rules).  
- **SC-003**: v1 import scope is unambiguous: **CSV via templates** for **credit card and normal bills** in; PDF/Word out; no assumption of one universal CSV layout.  
- **SC-004**: App admin vs org bill access is unambiguous: platform only, no break-glass bills in v1.  
- **SC-005**: v1 org types **personal** and **company** are unambiguous.  
- **SC-006**: Stakeholders understand every bill—including credit card—has a bill-level **establishment**, and credit card bills also have **transactions → establishments**.  
- **SC-007**: Stakeholders understand import is **upload → preview (editable) → accept or reject**, with no lasting bill data from that file until accept.  
- **SC-008**: Stakeholders understand the **dashboard** shows month-by-month costs and **tag-based** category breakdown (food, transport, leisure, …).

---

## Domain model (conceptual)

Readable sketch only — exact schema is chosen in a later plan.

```text
USER
----
id
name                 # required
address              # required
email
isAppAdmin

ORGANIZATION
------------
id
name                 # required
address              # required
type                 # personal | company

MEMBERSHIP
----------
id
userId
organizationId
roles[]              # org_admin | bill_manager | viewer

ESTABLISHMENT
-------------
id
organizationId
name
address              # searchable
tags[]               # food | transport | leisure | utilities | ...

BILL
----
id
organizationId
establishmentId      # required (card/issuer OR utility provider, etc.)
type                 # credit_card | water | power | ...
amount
date / period
note?
source               # manual | csv

TRANSACTION
-----------
id
billId
establishmentId      # required
amount
date
note?
```

```text
User ──< Membership >── Organization
Organization ──< Bill ──────────> Establishment
                  └──< Transaction ──> Establishment
Organization ──< Establishment
```

---

### CSV template (conceptual)

```text
CSV TEMPLATE (profile)
----------------------
id / key                 # e.g. "issuer_x_card_export_v1"
displayName              # shown when user picks a template
sourceKind               # credit_card | water | power | ...  (card OR normal bills)
delimiter / encoding / dateFormat / amountRules
columnMapping            # how file columns map to bill / establishment / transaction fields
                             # (transactions optional — often empty for normal bills)
```

Which concrete issuer/utility templates ship in v1 is decided in a later feature spec; the product rule is **template-based CSV for credit card and normal bills**, not one global format and not cards-only.

---

## Assumptions

- Auth method (email/password, etc.) is chosen in a later auth spec; this doc only requires authentication.  
- One primary currency per organization unless a later spec adds multi-currency.  
- A personal-type org may start with a single member; multi-member personal orgs are allowed if the product invites them.  
- Exact limits on who may create company vs personal orgs can be refined in a later feature spec.  
- There is **no industry-standard bill CSV**; supported layouts are whatever **CSV templates** the product defines.  
- LLM-assisted mapping is optional later; v1 parsers are expected to be **deterministic per template**, with **editable preview** and explicit **accept / reject**.  
- Exact rules for how bill-level vs transaction-level amounts roll into a tag on the dashboard can be refined in a later analytics feature spec; the product intent is month view + tag breakdown so members can spot overspending by category.  
- Monorepo and AWS are delivery context, not product rules (`STACK.md`).

---

## Out of scope (v1)

- User-owned bills outside an organization  
- PDF / Word bill parsing  
- Bank feeds / open banking  
- App-admin break-glass into org bills  
- AI advice, auto-pay, investments, full budgeting suites  
- Native mobile apps  
- SSO/SAML (unless a later spec adds them)

---

## Decisions log

| Topic | Decision |
| --- | --- |
| User profile | **name** and **address** required |
| Organization profile | **name** and **address** required (plus type) |
| User-owned bills | **No** — bills only on organizations |
| Personal use | Via organization type **personal** |
| Org types (v1) | At least **personal** and **company** |
| Multi-org membership | Yes |
| v1 org roles | Org Admin, Bill Manager, Viewer |
| App admin | Platform only; no automatic bill access |
| Bills | **Manual entry in v1** + CSV import in v1; **every bill** (incl. credit card) has a bill-level **establishment** |
| Dashboard | Month-by-month expenses + analysis by **establishment tags** (food, transport, leisure, …) |
| CSV | **Template/profile based** for **credit card and normal bills** (water, power, …)—not cards-only; no single universal CSV standard |
| CSV import UX | **Upload → editable preview → accept or reject**; no lasting DB commit until accept |
| CSV parsing approach (v1) | Deterministic parsers per template + preview/edit/confirm (LLM optional later) |
| Credit card | Bill establishment = card/issuer; transactions each have their own establishment |
| Establishment | Reusable; name, address, tags; used by bills and by transactions |
| Establishment tags (v1 examples) | food, transport, leisure, utilities |
| PDF/Word | Later |
| Bill type examples | Credit card (with transactions), water, power, etc. |
| User stories in this doc | Deferred — definition only; stories in later feature specs |
| Canonical location | `docs/PRODUCT.md` (do not duplicate in `.specify/memory/`) |
