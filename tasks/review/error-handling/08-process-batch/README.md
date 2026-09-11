# Drill: 08-process-batch

- **Category:** C3 — Error handling (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> `ProcessBatch` runs `handle` on each item, logging and skipping failures, and returns an error
> if the batch as a whole failed (including the empty-batch case). Several error-handling
> mistakes are hiding here — find as many as you can.

## Files to review

- `drill.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format (one line per finding). **Do not run `undergo reveal review/error-handling/08-process-batch` until you submit.**
