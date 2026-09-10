# Kata k16 — Time-series store

**Track:** B — Storage & indexing · **Tier:** T2

## Problem
Design a store for huge volumes of timestamped numeric points with aggressive
compression and fast range queries — the engine behind Prometheus / Facebook Gorilla.

## Requirements
**Functional**
- Append a (series, timestamp, value) point.
- Range-query a series over [from, to].

**Non-functional**
- High compression on typical data (regular timestamps, slowly-changing values).
- Append is amortized O(1); timestamps within a series are non-decreasing.

## Given
- Millions of points/second; mostly-regular scrape intervals.
- Single process.

## Your core must
Expose a `Store` with:
- `Append(series string, t int64, v float64) error` — append one point to a series.
- `Range(series string, from, to int64) []Point` — the points in the time range.

Compress with delta-of-delta timestamp encoding + Gorilla XOR float encoding. **Tests must
show:** lossless roundtrip (decoded points equal what was appended); an out-of-order append
(t ≤ last for the series) is rejected or has defined behavior; for a regular series the
compressed size is materially smaller than raw 16-byte points (assert a ratio).

## Design should also address (in DESIGN.md)
- Delta-of-delta for timestamps and XOR for floats — why they shrink scrape data so well.
- Chunk/block boundaries and when a chunk is closed and sealed.
- Handling counter resets and irregular intervals gracefully.
- The read path: decompress a chunk, then filter to [from, to].

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
