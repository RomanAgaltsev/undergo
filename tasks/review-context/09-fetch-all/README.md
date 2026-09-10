# Drill: 09-fetch-all

- **Category:** C4 — Context & cancellation (see `../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> `FetchAll` retrieves every URL in sequence, giving each a one-second timeout, and returns the
> collected results. Several context mistakes are hiding here — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../rubric/review-rubric.md`; write findings in the
`../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review-context/09-fetch-all` until you submit.**
