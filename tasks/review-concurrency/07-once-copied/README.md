# Drill: 07-once-copied

- **Category:** C1 — Concurrency & races (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> `Resource.Get` uses a `sync.Once` so the connection is initialized exactly once, no matter how
> many times (or from how many goroutines) `Get` is called.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-concurrency/07-once-copied` until you submit.**
