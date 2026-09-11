# Drill: 01-ignored-ctx

- **Category:** C4 — Context & cancellation (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** obvious

## PR description

> Adds two HTTP fetch helpers that both take a `context.Context` so callers can cancel or
> time-limit the request: `Fetch` and `FetchWithCtx`.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/context/01-ignored-ctx` until you submit.**
