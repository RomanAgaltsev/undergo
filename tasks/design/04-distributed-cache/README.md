# Kata k04 — Distributed cache

**Track:** E — Caching & data access · **Tier:** T2

## Problem
Design a cache that fronts a slow backing store for many application nodes: it evicts
under memory pressure, expires entries by TTL, and protects the backing store from
stampedes by loading each missing key only once.

## Requirements
**Functional**
- Get a value by key.
- Set a value with a time-to-live.
- On a miss, load from the backing store exactly once even under concurrent requests
  for the same key (singleflight), and cache the result.

**Non-functional**
- Concurrency-safe.
- Bounded memory via capacity eviction.
- A hot missing key triggers exactly one backing-store load, not one per caller.

## Given
- ~100k hot keys; ~50k gets/second; a backing store that is expensive to hit.
- Multiple application instances (how the cache is distributed is a design concern).

## Your core must
Expose a `Cache` with:
- `Get(key string) ([]byte, bool)` — cached value, or miss.
- `Set(key string, val []byte, ttl time.Duration)` — store with expiry.
- `GetOrLoad(key string, loader func() ([]byte, error)) ([]byte, error)` — return the
  cached value or run the loader exactly once for concurrent misses (singleflight).

Include capacity eviction and TTL expiry. **Tests must show:** eviction under capacity;
TTL expiry removes entries; under N concurrent `GetOrLoad` calls for the same missing key
the loader runs exactly once and every caller receives the value; concurrency-safe.

## Design should also address (in DESIGN.md)
- How the cache is distributed across nodes — sharding by consistent hashing (cf. k3) vs a
  replicated cache — and the tradeoffs.
- Invalidation / coherence across nodes when a value changes.
- Thundering-herd protection (singleflight) and negative caching for missing keys.
- You may reuse patterns from **streamcache** — reference it in the design, don't depend on it.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
