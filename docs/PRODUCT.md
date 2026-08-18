# Expenses Analyzer — Product

**Status**: Approved  
**Updated**: 2026-08-18  
**Approved**: 2026-08-18

This is the **only** source of truth for global product intent: what the app is, who it is for, how it works, and normative business rules. Feature specs under `specs/` must not contradict this document; they add delivery detail for one slice at a time.

Technical stack: [`STACK.md`](../STACK.md). Process / SDD: [`.specify/memory/constitution.md`](../.specify/memory/constitution.md). Delivery slices (draft): [`specs/README.md`](../specs/README.md).

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
| **User (account)** | Signs **up** and signs in; has **name** and **address** (free text). Has **no** personal bill list of their own. Accesses bills only through **organization membership**. |
| **Organization member** | Belongs to one or more **organizations**; works on that org’s bills according to their **roles**. |
| **Org Admin** | Organization role: manages membership/roles/invites in that org; full bill access in that org. |
| **Bill Manager** | Organization role: creates, imports, edits, deletes bills and related establishments; does not manage membership. |
| **Viewer** | Organization role: reads bills and summaries only. |
| **Application admin** | Platform role: manages the **application** (users, organizations, settings). Does **not** automatically see an org’s bills unless they are also a member. The first application admin is created by a **migration / seed script**, not by self-signup. |

**Rules of belonging**

- An **organization** has **many users** (a personal-type org may often have one member, but it is still an organization).  
- A **user** can belong to **many organizations** (e.g. one personal org and several company orgs).  
- In one organization, a user can hold **one or more** of: Org Admin, Bill Manager, Viewer.  
- Permissions in an org are the **union** of that user’s roles there.  
- A user **joins** an existing organization only when an **Org Admin** of that org **invites** them. There is no open join and no Bill Manager / Viewer invite.

**Organization types (v1)**

| Type | Intent |
| --- | --- |
| **personal** | Org created for an individual’s (or household-style) finances; still membership + roles + bills like any org. |
| **company** | Org for a business / formal organization. |

Organizations have a **name**, **address** (free text), **timezone**, **currency**, and **type**. The creator (Org Admin) **MUST choose the currency** when creating the org—there is no implied default. After create, only an **Org Admin** may change name, address, timezone, or currency.

Additional types may be added later without changing the rule that **all bills live on an organization**.

### Permission matrix (v1)

Capabilities below are **per organization**. Application admin is a platform role and is **not** a substitute for membership.

| Capability | Org Admin | Bill Manager | Viewer |
| --- | --- | --- | --- |
| View bills, transactions, establishments, dashboard | Yes | Yes | Yes |
| Create / edit / delete bills and transactions | Yes | Yes | No |
| CSV import (preview → accept or reject) | Yes | Yes | No |
| Create / edit / delete establishments | Yes | Yes | No |
| Assign / change establishment tags | Yes | Yes | No |
| Invite members, change roles, remove members | Yes | No | No |
| Edit org name, address, timezone, currency | Yes | No | No |
| Leave the organization | Yes, except last Org Admin | Yes | Yes |
| Delete the organization | Yes | No | No |

Any authenticated user may **create** a new organization (becomes Org Admin there). That is not an org-role permission.

---

## How it works (conceptual)

```text
Sign up (email, password, name, address)
   → sign in with email and password
      (forgot password → reset by email)
   → create an organization, or accept an Org Admin invite
   → choose an organization (personal-type, company-type, …)
        → manage that org’s bills only
             • add manually
             • import CSV → preview/edit → accept or reject
             • list / edit / delete (if role allows)
             • dashboard: month-by-month costs + breakdown by tags
```

1. **Organizations isolate data**  
   Org A’s bills never appear in Org B. Switching organization changes the entire working context. Users do not have a parallel personal bill store outside an org.

2. **Creating an organization**  
   Any authenticated user can create an organization of type **personal** or **company** and becomes **Org Admin** there. Personal use is a personal-type org—not a separate product mode with user-owned bills.  
   On create, the creator **MUST choose a currency** (no default). The org MUST have a **timezone** (default **UTC** if the creator does not pick one). Dashboard months use that timezone so all members see the same month boundaries. Only an **Org Admin** may later edit name, address, timezone, or currency. Changing currency MUST NOT convert existing amounts (numbers stay; the label changes).

3. **Joining an organization**  
   Only an **Org Admin** can invite someone (by email) into that org and assign roles. The invitee signs up or signs in, then accepts. An invite stays valid until the invitee **accepts** or an Org Admin **cancels** it. Bill Manager and Viewer cannot invite. There is no public join code and no open directory of orgs. Email verification is **not** required in v1. A user who does not remember their password MUST be able to **reset it by email**.

4. **Leaving, last admin, and deleting**  
   Members may leave an org. An Org Admin may remove other members and change roles. The **last remaining Org Admin** MUST NOT leave, be removed, or drop Org Admin until they assign Org Admin to another member **or** delete the organization. An Org Admin MAY delete the organization (that org’s memberships, bills, transactions, and establishments go with it).

5. **Bills are the core record**  
   A bill has at least: amount, a **single bill date**, type (from the closed v1 list), optional note, and an **establishment**. There is **no** separate billing-period start/end in v1. Source is manual or CSV import. Every bill belongs to exactly one organization.  
   The bill’s **amount is the parent total**: it is entered and stored on the bill itself. When transactions exist, the product MUST show the parent total, the sum of transaction amounts, and any difference—no silent mismatch. The parent total is **not** overwritten to equal the transaction sum.

6. **Establishments on every bill; transactions add more establishments**  
   An **establishment** is a named party or place (utility provider, **credit card / issuer**, restaurant, cinema, etc.). It has an **address** as **free text** (searchable) and **category tags** from the closed v1 list. Establishments belong to an organization and **MUST be reusable**: many bills and many transactions MAY point at the same establishment.

   - **Every bill**—including **credit card**, water, and power—MUST be associated with an establishment at bill level (the credit card itself / issuer is that establishment for a card bill; the water company for a water bill).  
   - **Credit card** bills ALSO contain multiple **transactions** (charges on the statement). Each transaction is tied to an establishment (e.g. Luigi Restaurant).  
   - So a credit card bill has **one bill-level establishment** (the card) **plus** many transaction-level establishments (where the card was used). Those transaction establishments MAY be shared with other bills and with the bill-level establishment list.  
   - Transactions live **on the parent bill**; they are not a separate top-level bill list.  
   - Simple bills may have **zero transactions** but still have a bill-level establishment.

7. **Manual entry**  
   Org Admin and Bill Manager type bill details (including picking or creating the bill’s establishment) and, when applicable, transactions + their establishments, and save them into the **current organization**. Members pick an existing establishment when they mean the same place. Deleting an establishment that any bill or transaction still references MUST be **blocked**.

8. **CSV import (v1) via templates — preview before save**  
   There is **no single industry-standard bill CSV**. Issuers and utilities use different columns, separators, dates, and amount conventions.  
   Authorized members upload a CSV and select (or the product matches) a **CSV template** (profile) for that source—e.g. a **credit card** issuer export **or** a **normal bill** export (water, power, and similar).  
   v1 ships **`expenses-card-statement-v1`** (card statement rows → one parent bill + transactions) and **`expenses-utility-bill-v1`** (**one row = one bill**, no transactions—water, power, gas, internet, phone). Additional templates may be added without changing the import gate.  
   CSV import is **not limited to credit cards**; templates cover both statement-style (often with transactions) and simple bill exports (bill + establishment, typically no transactions).  
   **Nothing is written to lasting bill data until the member accepts.** After upload, the app MUST show a **preview** of the proposed bills, establishments, and transactions. In preview, the member MUST be able to **correct** fields (and related proposed rows) before deciding.  
   The member then **accepts** (commit the previewed result) or **rejects** (discard the import; no lasting changes from that upload).  
   Unknown/unsupported templates or unreadable files fail clearly. **PDF and Word are not in v1.**  
   New templates can be added over time without changing the “CSV → preview/edit → accept or reject” product rule.

9. **Dashboard and analysis**  
   In the current organization, members use a **dashboard** to view expenses **month by month** in the **organization timezone**, inspect overall costs for each month, and break costs down **by establishment tags**. That is how they see what kinds of expenses they have and whether a category is too high in a given month.  
   Tagging establishments (and keeping tags accurate) is part of enabling this analysis.

   **Attribution (v1, locked):** overall month cost uses **parent bill amounts**. If a bill has transactions, tag totals use **transaction** amounts and each transaction establishment’s tags. If a bill has zero transactions, tag totals use the **parent bill amount** and the bill-level establishment’s tags. Overall may differ from the sum of tag totals when parent total ≠ transaction sum; both views stay visible. An amount is counted **once** in overall cost. If an establishment has multiple tags, that amount appears in **each** tagged category. Untagged spend goes in an **untagged** bucket (not omitted).

10. **Platform administration**  
    Application admins operate the product (users, orgs, settings). They are not a backdoor into every org’s bill list. The first application admin account is created by a **migration / seed script**. An existing application admin MAY grant application admin to another user. Platform “settings” in v1 MAY be a stub.

---

## Business rules (normative)

### Organizations and tenancy

- **BR-001**: All bills MUST belong to an **organization**. Users MUST NOT have a separate personal bill collection outside an organization.  
- **BR-002**: Organizations MUST have a **type**. v1 types include at least **personal** and **company**.  
- **BR-003**: Organizations MUST have a **name** and an **address**. Address is **free text**.  
- **BR-004**: Users MUST have a **name** and an **address**. Address is **free text**.  
- **BR-005**: Any authenticated user MUST be able to **create** an organization of either v1 type (**personal** or **company**). The creator MUST become **Org Admin** of that org. Creating a company org MUST NOT require application-admin approval. On create, the creator MUST **choose a currency**; the product MUST NOT assign a default currency.  
- **BR-006**: Bill and membership data MUST be isolated per organization (no cross-org leakage).  
- **BR-007**: A user MUST be able to belong to multiple organizations and switch among organizations they may access.  
- **BR-040**: A user MUST join an existing organization **only** via an **Org Admin invitation** (invitee email). There MUST be no open join, join codes, or invites from Bill Manager / Viewer. An invitation MUST remain valid until the invitee **accepts** or an Org Admin **cancels** it.  
- **BR-047**: Only an **Org Admin** MAY edit the organization’s name, address, timezone, or currency. Changing currency MUST NOT convert existing amounts.  
- **BR-041**: Members MAY leave an organization. An Org Admin MAY remove other members and change roles. The last remaining Org Admin MUST NOT leave, be removed, or drop the Org Admin role until they assign Org Admin to another member or delete the organization.  
- **BR-042**: An Org Admin MAY **delete** the organization. Deletion MUST remove that org’s memberships, bills, transactions, and establishments (no leftover access to that org’s bills).  
- **BR-043**: Each organization MUST have a **timezone**. Dashboard “month by month” MUST use that timezone. If none is chosen at create, the timezone MUST default to **UTC**.

### Roles

- **BR-008**: v1 organization roles are exactly **Org Admin**, **Bill Manager**, and **Viewer**, with meanings in “Who it is for” and the **permission matrix**.  
- **BR-009**: Multiple roles for one user in one org MUST combine as a **union** of permissions.  
- **BR-010**: **Application admin** manages platform concerns only; it MUST NOT grant automatic access to an organization’s bills without membership.  
- **BR-044**: The first application admin MUST be created by a **migration / seed script**, not by ordinary sign-up. An existing application admin MAY grant application admin to another user. Platform settings in v1 MAY be a stub.

### Accounts

- **BR-023**: Users MUST authenticate with **email and password** before accessing organization data. v1 MUST NOT require SSO, SAML, OAuth social login, or magic links.  
- **BR-045**: Anyone MUST be able to **sign up** for an account with email, password, name, and address. Sign-up does not by itself grant application-admin or membership in an existing org.  
- **BR-048**: v1 MUST NOT require **email verification** to sign up or sign in.  
- **BR-052**: A user who does not remember their password MUST be able to **reset it by email**: request a reset (from sign-in), receive an email, set a new password, then sign in with the new password. This is not magic-link sign-in.

### Bills

- **BR-011**: Organizations MUST support many bills. v1 bill types are a **closed list**: `credit_card`, `water`, `power`, `gas`, `internet`, `phone`. Members MUST NOT invent new types in v1. Only `credit_card` MUST support transactions in v1. Types other than `credit_card` MUST have zero transactions; adding transactions to them MUST be rejected.  
- **BR-049**: Every bill MUST have a **single bill date**. v1 MUST NOT store a separate billing-period start/end.  
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
- **BR-024**: Validation and parse failures MUST be explained in plain language; the system MUST NOT silently invent amounts.  
- **BR-025**: Every bill MUST be associated with an **establishment**, including **credit card** bills (the card/issuer is the bill-level establishment) and “normal” bills such as water and power (the provider is the establishment).  
- **BR-046**: v1 MUST ship **`expenses-card-statement-v1`** and **`expenses-utility-bill-v1`** as defined under [`docs/csv-templates/`](csv-templates/). The utility template is **one row = one bill** (no transactions). Additional templates MAY be added later without changing BR-014–BR-021.

### Dashboard and tag analysis

- **BR-026**: Members with view access MUST have a **dashboard** for the current organization that shows expenses **month by month**.  
- **BR-027**: For a selected month, the dashboard MUST show overall costs so members can analyze what was spent.  
- **BR-028**: Members MUST be able to see costs **categorized by establishment tags** (e.g. food, transport, leisure, utilities) for a selected month, so they can see what kinds of expenses they have and whether a category is high.  
- **BR-029**: **Org Admin** and **Bill Manager** MUST be able to assign and change **establishment tags** so dashboard categorization stays accurate. Viewer MUST NOT change tags.  
- **BR-030**: Dashboard figures MUST respect organization isolation (only the current org’s data).  
- **BR-031**: Each organization MUST use a **single primary currency** for bill and dashboard amounts in v1. That currency MUST be chosen at organization creation (BR-005).

### Transactions and establishments

- **BR-032**: An **establishment** MUST have a name; it MUST support an **address** as **free text** that members can search (substring is enough in v1); it MUST support **category tags** from the **closed v1 list**: `food`, `transport`, `leisure`, `utilities`, `credit_card`, `housing`, `health`, `shopping`. Members MUST NOT invent new tags in v1. An establishment MAY have more than one tag from that list.  
- **BR-033**: Establishments MUST be **sharable**: many bills and many transactions in the same organization MAY reference the **same** establishment. An establishment is not itself a bill. Members MUST be able to pick an existing establishment instead of creating a duplicate when they mean the same place. Duplicate names are allowed; sharing is by selecting the same establishment, not by matching names automatically.  
- **BR-034**: Bills of types that carry line detail (at least **credit card** in v1) MUST support **transactions**: discrete charges **on that parent bill**.  
- **BR-035**: Each transaction MUST be linkable to an **establishment** (e.g. Luigi Restaurant, Max Cinemas). That establishment MAY already be used by other transactions or bills in the org.  
- **BR-036**: Transactions MUST NOT appear as a detached top-level bill list; they belong to their parent bill.  
- **BR-037**: Authorized members MUST be able to add, edit, and remove transactions on a parent bill, and create/update establishment details, subject to the same bill permissions as editing that bill.  
- **BR-038**: The parent bill has its **own** amount (parent total). When transactions exist, the product MUST display the parent total, the **sum of transaction amounts**, and any **difference**. The product MUST NOT silently change the parent total to match the sum, and MUST NOT hide a mismatch.  
- **BR-039**: Bill types that do not need line detail (water, power, gas, internet, phone) MUST have zero transactions, but MUST still have a bill-level establishment; credit card bills in v1 MUST allow multiple transactions with establishments.  
- **BR-050**: The product MUST **block** deleting an establishment that is still referenced by any bill or transaction in that organization. Unused establishments MAY be deleted by Org Admin or Bill Manager.  
- **BR-051**: Dashboard attribution MUST follow the locked rules in “Dashboard and analysis” (overall = parent bill amounts; tags from transactions when present, otherwise the bill; multi-tag appears in each tag; untagged bucket; do not hide parent vs tag-sum differences).

---

## Usage overview *(lightweight — not delivery user stories)*

Illustrative flows for shared understanding. Detailed user stories belong in later feature specs.

1. **Personal use** — User signs up, creates a **personal**-type organization, becomes Org Admin, adds bills there.  
2. **Company use** — User works in a **company**-type org with other members and roles after an Org Admin invite.  
3. **Viewer is restricted** — Can read that org’s bills; cannot delete, import, or change tags.  
4. **Same person, two orgs** — Switches from personal org to company org; data set changes completely.  
5. **App admin** — Seeded by migration script; manages users/orgs on the platform; still cannot open Org D’s bills without membership.  
6. **CSV import** — Chooses **`expenses-card-statement-v1`** or **`expenses-utility-bill-v1`**, uploads file → **previews** proposed changes → **edits corrections in preview** → **accepts** (commit) or **rejects** (discard; no lasting import).  
7. **Dashboard** — Views month-by-month totals (org timezone) and breakdown by tags (food, transport, leisure, …) to judge category spend.  
8. **Credit card bill** — Bill-level **establishment** = the credit card/issuer; plus **transactions**, each with their own establishment (restaurant, cinema, …). Parent total may differ from the transaction sum; both are shown.  
9. **Utility bill** — Bill-level **establishment** = the provider; **no** transactions; a **single bill date**.  
10. **Shared establishment** — Luigi Restaurant is selected on several transactions (and possibly other bills) instead of creating a new establishment each time.  
11. **Last Org Admin** — Cannot leave until they promote someone else or delete the org.  
12. **In-use establishment** — Delete is blocked while any bill or transaction still points at it.  
13. **Forgot password** — From sign-in, requests a reset email, sets a new password, signs in.

---

## Success criteria (product-level)

- **SC-001**: Stakeholders can explain that bills live only on organizations (including personal-type orgs), not on user accounts.  
- **SC-002**: Acceptance tests for later features can map cleanly to BR-001–BR-052 (no contradictory tenancy, role, bill, import-preview, dashboard, transaction, or establishment rules).  
- **SC-003**: v1 import scope is unambiguous: **CSV via templates** for **credit card and normal bills** in; PDF/Word out; shipped templates **`expenses-card-statement-v1`** and **`expenses-utility-bill-v1`** (one row = one bill).  
- **SC-004**: App admin vs org bill access is unambiguous: platform only, no break-glass bills in v1; first app admin comes from a migration/seed script.  
- **SC-005**: v1 org types **personal** and **company** are unambiguous.  
- **SC-006**: Stakeholders understand every bill—including credit card—has a bill-level **establishment**, and credit card bills also have **transactions → establishments**.  
- **SC-007**: Stakeholders understand import is **upload → preview (editable) → accept or reject**, with no lasting bill data from that file until accept.  
- **SC-008**: Stakeholders understand the **dashboard** shows month-by-month costs (org timezone) and **tag-based** category breakdown (food, transport, leisure, …).  
- **SC-009**: Stakeholders can explain join = Org Admin invite only; sign-up is self-serve; last Org Admin cannot abandon the org.  
- **SC-010**: Stakeholders understand the parent bill total is independent of the transaction sum, and a mismatch is visible.  
- **SC-011**: Stakeholders understand establishments are shared by selection, and v1 bill types and tags are closed lists.  
- **SC-012**: Currency is chosen at org create; only Org Admin edits org profile; bills have a single date; in-use establishments cannot be deleted.

---

## Domain model (conceptual)

Readable sketch only — exact schema is chosen in a later plan.

```text
USER
----
id
name                 # required
address              # required, free text
email
password             # sign-up / sign-in
isAppAdmin           # first one seeded by migration script

ORGANIZATION
------------
id
name                 # required
address              # required, free text
type                 # personal | company
currency             # required at create; chosen by creator; no default
timezone             # required; default UTC; dashboard months use this

MEMBERSHIP
----------
id
userId
organizationId
roles[]              # org_admin | bill_manager | viewer

INVITATION
----------
id
organizationId
email                # invitee
roles[]              # assigned on accept
invitedByUserId      # must be Org Admin

ESTABLISHMENT
-------------
id
organizationId
name
address              # free text, searchable
tags[]               # closed v1: food | transport | leisure | utilities |
                     #             credit_card | housing | health | shopping

BILL
----
id
organizationId
establishmentId      # required (card/issuer OR utility provider, etc.)
type                 # closed v1: credit_card | water | power | gas | internet | phone
amount               # parent total (local to the bill; not auto-synced to tx sum)
date                 # single bill date (no period start/end in v1)
note?
source               # manual | csv

TRANSACTION
-----------
id
billId
establishmentId      # required; MAY be shared with other txs/bills
amount
date
note?
```

```text
User ──< Membership >── Organization
User ──< Invitation  ── Organization
Organization ──< Bill ──────────> Establishment
                  └──< Transaction ──> Establishment
Organization ──< Establishment
```

---

### CSV template (conceptual)

```text
CSV TEMPLATE (profile)
----------------------
id / key                 # e.g. "expenses-card-statement-v1"
displayName              # shown when user picks a template
sourceKind               # credit_card | water | power | ...  (card OR normal bills)
delimiter / encoding / dateFormat / amountRules
columnMapping            # how file columns map to bill / establishment / transaction fields
                             # (transactions optional — often empty for normal bills)
```

v1 ships **`expenses-card-statement-v1`** and **`expenses-utility-bill-v1`**: see [`docs/csv-templates/`](csv-templates/). The product rule remains **template-based CSV for credit card and normal bills**, not one global format and not cards-only.

---

## Confirmed assumptions

These are **accepted product decisions** (not open questions). Later feature specs may refine mechanics but must not contradict them.

- **Auth**: Users **sign up** and sign in with **email and password** in v1. **No email verification** required. **Password reset by email** is in v1. Session/token mechanics belong in a later auth feature spec; the product method is not open.  
- **Join**: Existing orgs are joinable only by **Org Admin invite**. Invite lasts until accept or Org Admin cancel.  
- **First app admin**: Created by **migration / seed script**. Existing app admin MAY grant the role to another user. Settings MAY be a stub.  
- **Addresses**: **Free text** for users, organizations, and establishments. Establishment search is substring in v1.  
- **Bill types (v1)**: Closed list `credit_card`, `water`, `power`, `gas`, `internet`, `phone`. Non-card types MUST NOT have transactions.  
- **Tags (v1)**: Closed list `food`, `transport`, `leisure`, `utilities`, `credit_card`, `housing`, `health`, `shopping`.  
- **Timezone**: Per organization; dashboard months use it; default UTC. Only Org Admin may change it.  
- **Parent total**: Stored on the bill; not auto-equalized to transaction sum; mismatch is shown.  
- **Bill date**: Single date; no period range in v1.  
- **Shared establishments**: Many bills/transactions may select the same establishment. In-use establishments cannot be deleted.  
- **Currency**: Chosen **at org create** (required, no default). One primary currency per organization in v1; multi-currency / FX is out of scope. Changing currency does not convert amounts. Only Org Admin may change it.  
- **Org profile edits**: Org Admin only (name, address, timezone, currency).  
- **Personal orgs**: A personal-type org may start with a single member; multi-member personal orgs are allowed.  
- **Who may create which org type**: Any authenticated user may create **personal** or **company** organizations. No app-admin gate in v1.  
- **CSV reality**: There is **no industry-standard bill CSV**; supported layouts are only what **CSV templates** define.  
- **CSV parsing (v1)**: Deterministic parsers **per template**, plus **editable preview** and explicit **accept / reject**. LLM-assisted mapping is optional later, never a substitute for preview/validation.  
- **v1 CSV templates**: `expenses-card-statement-v1` and `expenses-utility-bill-v1` (one row = one bill).  
- **PDF/Word (later)**: When added, prefer hybrid extraction with the **same preview → accept/reject** gate; not blind auto-commit.  
- **Dashboard tag totals**: Locked attribution in “Dashboard and analysis” (not deferred).  
- **Delivery context**: Monorepo and AWS are not product rules (`STACK.md`). Delivery slices live under `specs/` and stay **Draft** until approved; do not treat this file as a backlog.

---

## Out of scope (v1)

- User-owned bills outside an organization  
- PDF / Word bill parsing  
- Bank feeds / open banking  
- App-admin break-glass into org bills  
- AI advice, auto-pay, investments, full budgeting suites  
- Native mobile apps  
- SSO/SAML, OAuth social login, and magic-link sign-in (unless a later spec adds them)  
- Open join / join codes for organizations  
- Custom bill types or custom tags beyond the closed v1 lists  
- Auto-syncing parent bill total to the sum of transactions  
- Email verification  
- Billing-period start/end (single bill date only)  
- Deleting an establishment that is still in use  
- Defaulting currency on org create

---

## Decisions log

| Topic | Decision |
| --- | --- |
| User profile | **name** and **address** required |
| Organization profile | **name** and **address** required (plus type); **Org Admin only** may edit name, address, timezone, currency |
| Address format | **Free text** (user, org, establishment); establishment search = substring |
| User-owned bills | **No** — bills only on organizations |
| Personal use | Via organization type **personal** |
| Org types (v1) | At least **personal** and **company** |
| Who may create orgs | Any authenticated user may create **personal** or **company**; no app-admin approval |
| Who may invite | **Org Admin** only |
| Invite lifetime | Until accept or Org Admin cancel |
| Multi-org membership | Yes |
| v1 org roles | Org Admin, Bill Manager, Viewer — see permission matrix |
| Org lifecycle | Members may leave; last Org Admin cannot abandon; Org Admin may delete org |
| App admin | Platform only; no automatic bill access; **first admin via migration/seed script**; existing app admin MAY grant the role; settings MAY be a stub |
| Sign-up | **Self-serve** email + password + name + address; **no** email verification |
| Password reset | **Yes** — forgot password → reset by email; then sign in with the new password |
| Bills | **Manual entry in v1** + CSV import in v1; **every bill** (incl. credit card) has a bill-level **establishment** |
| Bill date | **Single date**; no period start/end |
| Bill types (v1) | Closed: credit_card, water, power, gas, internet, phone; non-card = **no transactions** |
| Parent bill total | **Local to the parent**; show vs transaction sum; do not auto-overwrite |
| Dashboard | Month-by-month in **org timezone** (default UTC) + analysis by **establishment tags**; attribution **locked** |
| CSV | **Template/profile based** for **credit card and normal bills**—not cards-only; no single universal CSV standard |
| CSV import UX | **Upload → editable preview → accept or reject**; no lasting DB commit until accept |
| CSV parsing approach (v1) | Deterministic parsers per template + preview/edit/confirm (LLM optional later) |
| v1 shipped CSV templates | **`expenses-card-statement-v1`** and **`expenses-utility-bill-v1`** (one row = one bill) |
| Credit card | Bill establishment = card/issuer; transactions each have their own establishment |
| Establishment | Reusable / **shared by selection**; name, free-text address, tags; used by bills and by transactions |
| Delete establishment | **Blocked** while any bill or transaction references it |
| Establishment tags (v1) | Closed: food, transport, leisure, utilities, credit_card, housing, health, shopping |
| Currency (v1) | **Chosen at org create** (required, no default); one primary currency; no FX; change does not convert amounts |
| Auth method | **Email and password** in v1; SSO/SAML/OAuth social/magic link out |
| Personal org membership | May start with one member; multi-member allowed |
| PDF/Word parsing | Later; same preview → accept/reject gate when added |
| PDF/Word | Later |
| User stories in this doc | Deferred — definition only; stories in later feature specs |
| Canonical location | `docs/PRODUCT.md` (do not duplicate in `.specify/memory/`) |
