# Drill: 11-model-swallows-errors

- **Category:** C10 — Testing gaps (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> Adds a model-based test for `Store`: ten thousand random puts and deletes,
> mirrored into a plain map and checked after every step. It found nothing,
> which is the point.

## Files to review

- `drill.go`
- `store_test.go` — **this is what you're reviewing.**

## Your task

Review the **test** against `../../../../rubric/review-rubric.md`. It's green and it
runs ten thousand steps — but could it ever fail? Write findings in the
`../../../../rubric/submission-format.md` format.
**Do not run `undergo reveal review/testing/11-model-swallows-errors` until you submit.**
