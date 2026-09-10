# Kata k29 — HyperLogLog cardinality estimator

**Track:** G — Probabilistic & analytics · **Tier:** T1

## Problem
Design a data structure that estimates the number of **distinct** items seen in a
high-volume stream, using a small, **fixed** amount of memory regardless of how many
distinct items there are.

## Requirements
**Functional**
- Add an item to the estimator.
- Report the estimated distinct-count (cardinality) at any time.
- Merge two estimators into one whose estimate is the cardinality of the *union*.

**Non-functional**
- Memory is fixed and configured up front (independent of cardinality).
- The estimate stays within a known standard error for a chosen precision.
- Mergeable without re-reading the original streams.

## Given
- Cardinalities up to ~10^9 distinct items.
- Memory budget of a few KB (e.g. precision p = 14 → 2^14 = 16384 registers).
- Single process; concurrent `Add` is a design consideration, not a hard requirement.

## Your core must
Expose an `HLL` with:
- `Add(item []byte)` — hash the item and update one register.
- `Count() uint64` — the current cardinality estimate.
- `Merge(other *HLL) error` — union into the receiver (same precision required).

Implement registers + the harmonic-mean estimate with small-range and large-range
corrections. **Tests must show:** the estimate is within the standard error
(≈ 1.04 / √m) across several true cardinalities; `Merge` of two estimators equals the
estimate you'd get by adding both streams to one; memory (register count) is constant
regardless of how many items were added.

## Design should also address (in DESIGN.md)
- The precision (register count) vs accuracy vs memory tradeoff.
- HyperLogLog vs linear counting vs an exact set — when each wins.
- The hash-function quality requirement (uniformity) and why it matters.
- The small-range (linear counting) and large-range corrections, and where the plain
  estimator is biased.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
