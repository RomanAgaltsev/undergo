# Drill: 04-ignored-close

- **Category:** C3 — Error handling (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** obvious

## PR description

> `WriteConfig` creates a file and JSON-encodes the config into it, using `defer f.Close()` to
> clean up.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/error-handling/04-ignored-close` until you submit.**
