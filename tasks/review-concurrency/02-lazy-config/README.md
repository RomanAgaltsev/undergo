# Drill: 02-lazy-config

- **Category:** C1 — Concurrency & races (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> `Config` now loads its settings lazily on first `Get` instead of eagerly at startup, so
> processes that never read config don't pay the load cost. A `initialized` flag makes sure
> we only load once.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-concurrency/02-lazy-config` until you submit.**
