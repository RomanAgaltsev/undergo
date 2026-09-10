# Drill: 06-lock-ordering

- **Category:** C1 — Concurrency & races (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> Adds `Transfer` to move money between two `Account`s, locking both accounts' mutexes so the
> two balance updates are atomic. `Balance` reads a single account safely.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format. **Do not run `undergo reveal review-concurrency/06-lock-ordering` until you submit.**
