# Drill: 03-uncancellable-worker

- **Category:** C4 — Context & cancellation (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> Adds a background `Run` worker that ticks once a second, saves a running count under a
> per-tick timeout, and saves a heartbeat. It's meant to stop when its context is cancelled.
> Several context mistakes are hiding here — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-context/03-uncancellable-worker` until you submit.**
