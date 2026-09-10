# Drill: 08-pipeline

- **Category:** C4 — Context & cancellation (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> `Pipeline` tags the context with a request id, fetches data, and notifies a downstream service
> under a one-second timeout. Several context mistakes are hiding here — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-context/08-pipeline` until you submit.**
