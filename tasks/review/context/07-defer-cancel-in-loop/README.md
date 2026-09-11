# Drill: 07-defer-cancel-in-loop

- **Category:** C4 — Context & cancellation (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> `ProcessAll` runs `fn` for every id, giving each call its own one-second timeout context and
> deferring the cancel.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/context/07-defer-cancel-in-loop` until you submit.**
