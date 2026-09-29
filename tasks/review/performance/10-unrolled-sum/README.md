# Drill: 10-unrolled-sum

- **Category:** C7 — Performance (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** obvious

## PR description

> `SumUnrolled` replaces `Sum` on the metrics hot path: the loop is unrolled by
> four, and `BenchmarkSumUnrolled` runs about four times faster than `BenchmarkSum`.

## Files to review

- `drill.go`
- `drill_test.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/performance/10-unrolled-sum` until you submit.**
