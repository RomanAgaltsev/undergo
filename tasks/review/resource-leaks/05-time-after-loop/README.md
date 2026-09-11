# Drill: 05-time-after-loop

- **Category:** C5 — Resource leaks (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** obvious

## PR description

> `Poll` runs `check` once a second until it succeeds or the context is cancelled, using
> `time.After` for the interval.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/resource-leaks/05-time-after-loop` until you submit.**
