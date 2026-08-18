# CSV template: `expenses-card-statement-v1`

**Status**: Draft (product template — not an issuer-specific bank export)  
**Source kind**: `credit_card`  
**Display name**: Expenses Analyzer card statement (v1)

Canonical v1 credit-card import profile. One file can contain one or more card bills. Each data row is a **transaction**. Rows that share the same bill identity group into one **parent bill**. The parent bill amount is **local** (repeated on each row of the group); it is not required to equal the sum of `tx_amount`.

Normative import gate (preview → accept or reject) is in [`docs/PRODUCT.md`](../PRODUCT.md).

## File rules

| Rule | Value |
| --- | --- |
| Encoding | UTF-8 (no BOM required; BOM allowed) |
| Delimiter | comma `,` |
| Header | **Required**; column names must match exactly (case-sensitive) |
| Record | One transaction per data row |
| Dates | `YYYY-MM-DD` |
| Amounts | Decimal with `.` as decimal separator; no thousands separator; no currency symbol; **positive = spend** |
| Empty optional fields | Allowed (`bill_note`, `tx_note`, `card_address`, `merchant_address`, `merchant_tags`) |
| Tags | Pipe-separated (`food\|leisure`); each tag MUST be from the closed v1 list in PRODUCT.md |
| Quotes | RFC 4180-style quoting when a field contains comma, quote, or newline |

Currency is the **organization’s** primary currency, not a column.

## Columns

| Column | Maps to | Required |
| --- | --- | --- |
| `bill_date` | Parent bill date (statement / period end) | Yes |
| `bill_amount` | Parent bill **local** total | Yes |
| `bill_note` | Parent bill note | No |
| `card_name` | Bill-level establishment **name** (card / issuer) | Yes |
| `card_address` | Bill-level establishment **address** (free text) | No |
| `tx_date` | Transaction date | Yes |
| `tx_amount` | Transaction amount | Yes |
| `tx_note` | Transaction note | No |
| `merchant_name` | Transaction establishment **name** | Yes |
| `merchant_address` | Transaction establishment **address** (free text) | No |
| `merchant_tags` | Transaction establishment tags (pipe-separated) | No |

The bill-level establishment (the card) is tagged `credit_card` by this template. Members can change tags in preview or after accept.

## Grouping (rows → bills)

Rows belong to the **same parent bill** when all of these match:

- `bill_date`
- `bill_amount`
- `bill_note`
- `card_name`
- `card_address`

Each group becomes:

1. **One bill** — type `credit_card`, amount = `bill_amount`, date = `bill_date`, note = `bill_note`, establishment = card (`card_name` / `card_address`).
2. **One transaction per row** — date/amount/note from `tx_*`, establishment = merchant (`merchant_name` / `merchant_address` / `merchant_tags`).

Establishments **may be shared**: the same `merchant_name` + `merchant_address` in several rows (including across bills in the file) SHOULD preview as one reusable establishment, and the member MUST be able to point those rows at an establishment that already exists in the organization.

If `bill_amount` ≠ sum of `tx_amount` in the group, preview MUST show parent total, transaction sum, and the difference (PRODUCT.md BR-038).

## Failures (no lasting commit)

Reject the file before preview commit when:

- Not CSV / not UTF-8 text
- Missing header or unknown required column
- A required field is empty
- A date is not `YYYY-MM-DD`
- An amount is not a finite decimal number
- A `merchant_tags` value is not in the closed v1 tag list
- Zero data rows

Do not invent amounts or drop invalid rows silently.

## Example

See [`examples/expenses-card-statement-v1.example.csv`](examples/expenses-card-statement-v1.example.csv).

That sample is two Visa bills. **Luigi Restaurant** appears on both (shared establishment). January parent total `200.00` is **not** equal to the transaction sum (`45.90 + 28.00 + 12.50 = 86.40`); that mismatch must stay visible in preview.
