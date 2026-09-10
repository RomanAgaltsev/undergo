# Drill: 09-poller

- **Category:** C5 — Resource leaks (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> `Monitor.Watch` starts a background goroutine that probes the database on a ticker and can be
> stopped via its context. Several resources aren't managed correctly — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-resource-leaks/09-poller` until you submit.**
