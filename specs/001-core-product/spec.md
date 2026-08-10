# Feature Specification: Core Product Definition

**Feature Branch**: `001-core-product`

**Created**: 2026-08-10

**Status**: Pointer — canonical product rules live elsewhere

## Canonical document

**Global business rules (what / who / how) are maintained only in:**

→ **[`docs/PRODUCT.md`](../../docs/PRODUCT.md)**

Do **not** copy product rules into this file. Edit `docs/PRODUCT.md` when product intent changes.

## Spec Kit note

- This folder exists so Spec Kit / agents can record that a product-definition pass happened (`001-core-product`).
- Delivery work uses **new** feature specs (`002-…`, `003-…`) that **must align with** `docs/PRODUCT.md`.
- Do **not** run `/speckit-plan` against this pointer as if it were an implementable feature backlog; plan concrete slices (auth, orgs, bills, CSV, bootstrap) instead.

## Related

- Stack: [`STACK.md`](../../STACK.md)
- Process: [`.specify/memory/constitution.md`](../../.specify/memory/constitution.md)
- Agent guide: [`AGENTS.md`](../../AGENTS.md)
