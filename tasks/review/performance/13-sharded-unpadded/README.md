# Drill: 13-sharded-unpadded

- **Category:** C7 — Performance (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> Moves the per-user event counter to sixteen shards with a lock each, so
> goroutines counting different users never wait for each other.
> `BenchmarkInc` gives every goroutine a user of its own. On my machine
> `-cpu 1,2,4,8` reports 11, 24, 34 and 26 ns/op — the rise is the price of
> running in parallel, and it levels off.

## Files to review

- `drill.go`
- `drill_test.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/performance/13-sharded-unpadded` until you submit.**
