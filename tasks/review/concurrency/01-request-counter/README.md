# Drill: 01-request-counter

- **Category:** C1 — Concurrency & races (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** obvious

## PR description

> Adds a tiny `Server` that tracks how many times each URL path is hit. To keep the
> request path fast, we do the counter bookkeeping in a background goroutine so the
> handler returns immediately with `202 Accepted`.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/concurrency/01-request-counter` until you submit.**
