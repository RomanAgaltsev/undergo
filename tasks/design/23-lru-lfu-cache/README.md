# Kata k23 — LRU/LFU cache core

**Track:** E — Caching & data access · **Tier:** T1

## Problem
Design the in-memory cache primitive that sits under every larger cache: O(1) get/put with
capacity-bounded eviction by a replacement policy (LRU and/or LFU).

## Requirements
**Functional**
- Get and put by key.
- Evict when over capacity, by the chosen policy (least-recently-used and/or
  least-frequently-used).

**Non-functional**
- O(1) get and put — no scan of the entries.
- Fixed capacity.
- Accurate hit/miss accounting.

## Given
- Capacity up to ~1M entries; read-heavy.
- Single process.

## Your core must
Expose a generic `Cache[K comparable, V any]` with:
- `Get(key K) (V, bool)` — value and hit, or the zero value and miss.
- `Put(key K, val V)` — insert/update, evicting per policy when over capacity.

Implement O(1) operations via a map + an intrusive doubly-linked list (LRU) or frequency
buckets (LFU). **Tests must show:** eviction order is correct (LRU evicts the
least-recently-used entry; LFU the least-frequently-used); hit/miss counters are accurate;
operations are O(1) (no linear scan on the hot path).

## Design should also address (in DESIGN.md)
- LRU (map + intrusive list) vs LFU (O(1) frequency-bucket design) vs adaptive (ARC / 2Q).
- The intrusive-list layout that keeps get/put O(1).
- Thread-safety approach (single mutex vs sharded locks) and its throughput cost.
- How a TTL layer sits on top of this core (tie-in to k4).

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
