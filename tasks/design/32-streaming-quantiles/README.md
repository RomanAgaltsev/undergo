# Kata k32 — Streaming quantiles

**Track:** G — Probabilistic & analytics · **Tier:** T2

## Problem
Design a structure that estimates **quantiles** (p50, p90, p99, p999) over an unbounded
stream of numbers using memory that does not grow with the stream — the "what's my p99
latency" problem.

## Requirements
**Functional**
- Add a value.
- Query the estimated value at quantile q ∈ [0, 1].
- Merge two summaries (for distributed aggregation).

**Non-functional**
- Rank error is bounded by a configured ε.
- Memory is independent of the number of values added.
- Mergeable without re-reading the original values.

## Given
- Billions of values; need p50/p90/p99/p999 (tail accuracy matters).
- Single process; summaries are also merged across shards.

## Your core must
Expose a `Q` with:
- `Add(v float64)` — observe one value.
- `Quantile(q float64) float64` — the estimated value at quantile q.
- `Merge(other *Q) error` — combine two summaries.

Implement a t-digest (or Greenwald-Khanna) summary. **Tests must show:** the estimate at
p50/p90/p99 is within the rank error ε versus an exact sort; `Merge` gives a summary
close to one built from both streams together; memory (centroids/tuples) stays bounded
as the stream grows.

## Design should also address (in DESIGN.md)
- t-digest vs Greenwald-Khanna vs KLL — tail accuracy vs uniform accuracy vs simplicity.
- Rank error vs value error, and why tail quantiles are the hard case.
- Why mergeability matters for aggregating per-node summaries into a global one.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
