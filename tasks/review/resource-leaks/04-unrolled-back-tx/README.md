# Drill: 04-unrolled-back-tx

- **Category:** C5 — Resource leaks (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** obvious

## PR description

> `Transfer` debits one account and credits another inside a single DB transaction, committing
> at the end.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/resource-leaks/04-unrolled-back-tx` until you submit.**
