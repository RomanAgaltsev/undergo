# Drill: 07-body-on-error-status

- **Category:** C5 — Resource leaks (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> `Get` performs an HTTP GET, returns an error for any non-200 status, and otherwise reads and
> returns the body (closing it with `defer`).

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-resource-leaks/07-body-on-error-status` until you submit.**
