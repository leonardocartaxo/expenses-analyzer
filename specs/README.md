# Delivery slices

**Status**: Draft  
**Updated**: 2026-08-18

These folders are **feature specs** (what/why + acceptance). They MUST align with [`docs/PRODUCT.md`](../docs/PRODUCT.md). They are **not** a copy of PRODUCT.md and are **not** approved for `/speckit-plan` until a human approves the slice.

`001-core-product` is a pointer only. Do not plan it.

| Order | Slice | Intent |
| --- | --- | --- |
| 1 | [`002-bootstrap`](002-bootstrap/spec.md) | Scaffold monorepo + `pnpm verify`. No expenses domain. |
| 2 | [`003-user-auth`](003-user-auth/spec.md) | Sign up, sign in, password reset by email. |
| 3 | [`004-organizations`](004-organizations/spec.md) | Create org, invite (Org Admin), membership, roles, leave/delete. |
| 4 | [`005-establishments`](005-establishments/spec.md) | Org establishments: name, free-text address, closed tags, reuse. |
| 5 | [`006-bills-transactions`](006-bills-transactions/spec.md) | Manual bills + credit-card transactions; parent total is local. |
| 6 | [`007-csv-import`](007-csv-import/spec.md) | Template CSV → editable preview → accept or reject. |
| 7 | [`008-dashboard`](008-dashboard/spec.md) | Month-by-month totals + tag breakdown (org timezone). |
| 8 | [`009-application-admin`](009-application-admin/spec.md) | First app admin via migration script; platform admin only. |

Implement **002-bootstrap** first after PRODUCT.md and that spec are approved. Later slices depend on earlier ones.

CSV templates (product artifacts, not a slice): [`docs/csv-templates/expenses-card-statement-v1.md`](../docs/csv-templates/expenses-card-statement-v1.md), [`docs/csv-templates/expenses-utility-bill-v1.md`](../docs/csv-templates/expenses-utility-bill-v1.md).
