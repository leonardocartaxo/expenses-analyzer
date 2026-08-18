# CSV template: `expenses-utility-bill-v1`

**Status**: Draft  
**Source kind**: normal bills (`water` | `power` | `gas` | `internet` | `phone`)  
**Display name**: Expenses Analyzer utility bill (v1)

Canonical v1 **one row = one bill** import. No transactions. Use this for water, power, gas, internet, and phone—not for credit-card statements (`expenses-card-statement-v1`).

Normative import gate (preview → accept or reject) is in [`docs/PRODUCT.md`](../PRODUCT.md).

## File rules

| Rule | Value |
| --- | --- |
| Encoding | UTF-8 (no BOM required; BOM allowed) |
| Delimiter | comma `,` |
| Header | **Required**; column names must match exactly (case-sensitive) |
| Record | **One bill per data row** (zero transactions) |
| Dates | `YYYY-MM-DD` (single **bill date**; no period start/end) |
| Amounts | Decimal with `.` as decimal separator; no thousands separator; no currency symbol; **positive = spend** |
| Empty optional fields | Allowed (`note`, `provider_address`, `provider_tags`) |
| Tags | Pipe-separated (`utilities\|housing`); each tag MUST be from the closed v1 list in PRODUCT.md |
| Quotes | RFC 4180-style quoting when a field contains comma, quote, or newline |

Currency is the **organization’s** primary currency, not a column.

## Columns

| Column | Maps to | Required |
| --- | --- | --- |
| `bill_date` | Bill date | Yes |
| `amount` | Bill amount (parent total; there are no transactions) | Yes |
| `note` | Bill note | No |
| `bill_type` | Bill type: `water`, `power`, `gas`, `internet`, or `phone` | Yes |
| `provider_name` | Bill-level establishment **name** | Yes |
| `provider_address` | Bill-level establishment **address** (free text) | No |
| `provider_tags` | Establishment tags (pipe-separated) | No |

`bill_type` MUST NOT be `credit_card`. Card statements use `expenses-card-statement-v1`.

## Grouping

There is **no grouping**. Each data row becomes exactly one bill of `bill_type` with establishment `provider_name` / `provider_address` / `provider_tags`, and **zero** transactions.

The same provider on several rows SHOULD preview as a **shared** establishment when name and address match, and the member MUST be able to point a row at an establishment that already exists in the organization.

## Failures (no lasting commit)

Reject the file before preview commit when:

- Not CSV / not UTF-8 text
- Missing header or unknown required column
- A required field is empty
- A date is not `YYYY-MM-DD`
- An amount is not a finite decimal number
- `bill_type` is not one of `water`, `power`, `gas`, `internet`, `phone`
- A `provider_tags` value is not in the closed v1 tag list
- Zero data rows

Do not invent amounts or drop invalid rows silently.

## Example

See [`examples/expenses-utility-bill-v1.example.csv`](examples/expenses-utility-bill-v1.example.csv).

That sample is three bills. **City Water** is used on two rows (shared establishment).
