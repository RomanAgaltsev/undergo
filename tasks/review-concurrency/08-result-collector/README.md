# Drill: 08-result-collector

- **Category:** C1 — Concurrency & races (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> `Collect` fans `fn` out across all inputs concurrently, gathers the results, and closes a
> signal channel the first time a result is zero. Several concurrency mistakes are hiding here
> — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-concurrency/08-result-collector` until you submit.**
