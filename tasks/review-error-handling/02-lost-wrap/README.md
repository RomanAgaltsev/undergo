# Drill: 02-lost-wrap

- **Category:** C3 — Error handling (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> `LoadUser` now returns a domain `ErrNotFound` when the row is missing, and `StatusFor`
> maps that to a `404`. Everything else stays a `500`. Wiring looks straightforward.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-error-handling/02-lost-wrap` until you submit.**
