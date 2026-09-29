# Drill: 11-fast-parse

- **Category:** C7 — Performance (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** subtle

## PR description

> Replace `strconv.ParseInt` on the ingest hot path with a hand-rolled `ParseID`.
> IDs always come from our own encoder as six ASCII digits, and the test pins
> `ParseID` to `strconv` over that domain. `BenchmarkParseID` is about five times
> faster than `BenchmarkParseStd`.

## Files to review

- `drill.go`
- `drill_test.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/performance/11-fast-parse` until you submit.**
