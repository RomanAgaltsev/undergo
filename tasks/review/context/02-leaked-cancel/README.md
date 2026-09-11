# Drill: 02-leaked-cancel

- **Category:** C4 — Context & cancellation (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> `Query` runs a SQL query under a 2-second timeout and collects the first column of each row.
> It cancels the context when it's done.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/context/02-leaked-cancel` until you submit.**
