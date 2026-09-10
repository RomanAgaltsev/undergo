# Kata k31 — Top-K heavy hitters

**Track:** G — Probabilistic & analytics · **Tier:** T2

## Problem
Design a structure that reports the **K most frequent items** in a stream without
storing every distinct item — the "who are the top talkers" problem.

## Requirements
**Functional**
- Add an item.
- Query the current top-K items with their estimated counts.

**Non-functional**
- Bounded memory: O(K) (or O(1/ε)) counters, not O(distinct items).
- Error on the reported counts is bounded by the counter budget.

## Given
- Millions of distinct items; K ≈ 100.
- Single process; the stream does not fit in memory.

## Your core must
Expose a `TopK` with:
- `Add(item []byte)` — observe one occurrence.
- `Query(k int) []Entry` — the top k items, each an `Entry` of item + estimated count,
  in descending count order.

Implement a Space-Saving (or Misra-Gries) counter set with a fixed budget.
**Tests must show:** every true heavy hitter (frequency > N/K) is always reported;
the estimated counts are within the algorithm's guaranteed error; memory stays within
the fixed counter budget regardless of stream length.

## Design should also address (in DESIGN.md)
- Space-Saving vs Misra-Gries vs count-min-sketch + a heap — tradeoffs.
- The exact guarantee: which items are *guaranteed* to appear in the output, and which
  may be missed or over-counted.
- Merging summaries across shards (so per-partition top-Ks combine into a global one).

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
