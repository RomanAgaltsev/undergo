# Drill: 10-benchmark-wrong-type

- **Category:** C10 — Testing gaps (see `../../../../rubric/bug-taxonomy.md`)
- **Tier:** obvious

## PR description

> Swaps `SliceQueue` for the new `RingQueue` on the dispatch path, to put a
> ceiling on its memory. The benchmarks show the swap costs nothing:
> `BenchmarkRingQueue` and `BenchmarkSliceQueue` both report about 27 ns/op and
> one allocation per op, so this is a pure memory win.

## Files to review

- `drill.go`
- `queue_test.go` — **this is what you're reviewing.**

## Your task

Review the **tests** against `../../../../rubric/review-rubric.md`; write findings in
the `../../../../rubric/submission-format.md` format. **Do not run `undergo reveal review/testing/10-benchmark-wrong-type` until you submit.**
