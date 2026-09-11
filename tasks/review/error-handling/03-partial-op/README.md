# Drill: 03-partial-op

- **Category:** C3 — Error handling (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> `CreateOrder` runs a three-step operation — write, charge, notify — and rolls back on
> failure. It uses a named return and a deferred cleanup. A few of the error paths don't do
> what the description claims; find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review/error-handling/03-partial-op` until you submit.**
