# Drill: 09-refresh-cache

- **Category:** C1 — Concurrency & races (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> Adds a concurrent `Cache` backed by a map with an `RWMutex`, plus a background goroutine that
> stamps a refresh marker on an interval. `Get` and `Set` are the public accessors. Several
> concurrency mistakes are hiding here — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-concurrency/09-refresh-cache` until you submit.**
