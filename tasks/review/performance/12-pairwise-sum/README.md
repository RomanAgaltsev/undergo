# Drill: 12-pairwise-sum

- **Category:** C7 — Performance (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** multi-bug

## PR description

> Two hot-path improvements for the nightly report. `Total` now sums pairwise —
> fewer bounds checks, half the iterations — and `Report` renders with a
> `strings.Builder` instead of concatenating. `BenchmarkReport` is several times
> faster than `BenchmarkReportBaseline`.
> Several problems are hiding here — find as many as you can.

## Files to review

- `drill.go`
- `drill_test.go`

## Your task

Review against `../../../../rubric/review-rubric.md`; write findings in the
`../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/performance/12-pairwise-sum` until you submit.**
