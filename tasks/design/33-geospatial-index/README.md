# Kata k33 — Geospatial index

**Track:** G — Probabilistic & analytics · **Tier:** T2

## Problem
Design an in-memory index over 2-D points (latitude/longitude) that answers
**k-nearest-neighbour** and **bounding-box** queries efficiently — the "which drivers
are near this rider" problem.

## Requirements
**Functional**
- Insert / update a point by id.
- Delete a point by id.
- k-nearest query: the k closest points to a location.
- Range query: all points inside a bounding box.

**Non-functional**
- Query cost sublinear in the number of points.
- Correct across cell/quadrant boundaries (no missed neighbours at edges).

## Given
- Millions of points (e.g. live vehicle locations), churning frequently.
- High query rate; single process, in-memory.

## Your core must
Expose a `Geo` with:
- `Insert(id string, lat, lng float64)` — add or move a point.
- `Delete(id string)` — remove a point.
- `Nearest(lat, lng float64, k int) []string` — the k nearest ids, closest first.
- `Range(minLat, minLng, maxLat, maxLng float64) []string` — ids inside the box.

Implement a geohash bucketing (or quadtree). **Tests must show:** `Nearest` is correct
even when the true neighbours fall in an *adjacent* cell (the boundary problem — you must
search neighbouring cells); `Range` returns exactly the points inside the box; `Delete`
and re-`Insert` (move) are reflected in later queries.

## Design should also address (in DESIGN.md)
- Geohash vs quadtree vs k-d tree vs S2/H3 — tradeoffs for this workload.
- The boundary problem: why a single-cell lookup misses near neighbours, and how the
  neighbour-cell search fixes it.
- Cell precision vs cell count, and how it interacts with point density.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
