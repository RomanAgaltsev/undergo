# Kata k30 — Count-min sketch

**Track:** G — Probabilistic & analytics · **Tier:** T1

## Problem
Design a structure that estimates the **frequency** of individual items in a very
large stream using sublinear memory — trading a bounded overestimate for space.

## Requirements
**Functional**
- Add (increment) an item's count by n.
- Estimate an item's total count.

**Non-functional**
- Never *under*-estimate the true count.
- Overestimate is bounded by ε·N with probability at least 1 − δ (N = total added).
- Memory is fixed by the chosen (ε, δ); independent of the number of distinct items.

## Given
- Millions of distinct items; billions of updates.
- Memory budget of a few hundred KB.
- Single process.

## Your core must
Expose a `CMS` with:
- `Add(item []byte, n uint64)` — increment the item's counters by n.
- `Estimate(item []byte) uint64` — the estimated count (the min across rows).

Implement a d×w counter grid with d independent hashes, sized from (ε, δ)
(w ≈ ⌈e/ε⌉, d ≈ ⌈ln 1/δ⌉). **Tests must show:** `Estimate` is never below the true
count; across a workload the overestimate stays within the ε·N bound; memory is the
constant d×w regardless of distinct-item count.

## Design should also address (in DESIGN.md)
- The width/depth (ε, δ) sizing tradeoff — accuracy vs memory.
- Count-min sketch vs count-min-mean vs an exact frequency map.
- How this extends to finding heavy hitters (the tie-in to k31 top-K).
- Conservative update, and the collision/overestimate source.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
