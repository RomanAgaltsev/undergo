# Drill: 02-defer-in-loop

- **Category:** C5 — Resource leaks (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> `ConcatFiles` opens each file in a list, reads it, and concatenates the contents. It uses
> `defer f.Close()` so nothing is left open. `ReadOne` handles the single-file case.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-resource-leaks/02-defer-in-loop` until you submit.**
