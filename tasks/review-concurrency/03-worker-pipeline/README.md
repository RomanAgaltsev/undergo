# Drill: 03-worker-pipeline

- **Category:** C1 — Concurrency & races (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> Adds a fan-out worker pool: `Process` spins up `n` workers that each pull jobs off a
> channel, double them, and append to a shared results slice. `Warm` starts a background
> drainer to pre-warm a stream. Several things are subtly off — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-concurrency/03-worker-pipeline` until you submit.**
