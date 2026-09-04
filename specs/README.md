# Delivery slices

**Status**: Draft  
**Updated**: 2026-08-18

These folders are **feature specs** (what/why + acceptance). They MUST align with [`docs/PRODUCT.md`](../docs/PRODUCT.md). They are **not** a copy of PRODUCT.md. Plan a slice only after a human approves it.

`001-core-product` is a pointer only. Do not plan it.

| Order | Slice | Intent | Status |
| --- | --- | --- | --- |
| 1 | [`002-bootstrap`](002-bootstrap/spec.md) | Scaffold monorepo + `pnpm verify`. No expenses domain. | **Approved** (spec + plan) |
| 2 | [`011-aws-deploy`](011-aws-deploy/spec.md) | First hosted **dev**: CDK, EKS, ECR, RDS, Amplify, CI OIDC. No wake/sleep. | Draft |
| 3 | [`003-user-auth`](003-user-auth/spec.md) | Sign up, sign in, password reset by email. | Draft |
| 4 | [`004-organizations`](004-organizations/spec.md) | Create org, invite (Org Admin), membership, roles, leave/delete. | Draft |
| 5 | [`005-establishments`](005-establishments/spec.md) | Org establishments: name, free-text address, closed tags, reuse. | Draft |
| 6 | [`006-bills-transactions`](006-bills-transactions/spec.md) | Manual bills + credit-card transactions; parent total is local. | Draft |
| 7 | [`007-csv-import`](007-csv-import/spec.md) | Template CSV → editable preview → accept or reject. | Draft |
| 8 | [`008-dashboard`](008-dashboard/spec.md) | Month-by-month totals + tag breakdown (org timezone). | Draft |
| 9 | [`009-application-admin`](009-application-admin/spec.md) | First app admin via migration script; platform admin only. | Draft |
| 10 | [`010-poc-wake-sleep`](010-poc-wake-sleep/spec.md) | POC not 24/7: sleep stops all stoppable AWS; wake waits until ready. | Draft |

**`002-bootstrap` spec, plan, and tasks are ready.** Next is `/speckit-implement` (optional `/speckit-analyze` first). **`011-aws-deploy` is first-time AWS** (health-only is enough). **`010-poc-wake-sleep` is hosted-POC ops** after that stack exists — not a PRODUCT.md org feature. Later slices stay Draft until you approve them.

CSV templates (product artifacts, not a slice): [`docs/csv-templates/expenses-card-statement-v1.md`](../docs/csv-templates/expenses-card-statement-v1.md), [`docs/csv-templates/expenses-utility-bill-v1.md`](../docs/csv-templates/expenses-utility-bill-v1.md).
