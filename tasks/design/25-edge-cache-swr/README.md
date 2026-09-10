# Kata k25 — Edge cache (stale-while-revalidate)

**Track:** E — Caching & data access · **Tier:** T2

## Problem
Design a read-through edge cache that stays fast even when entries expire: return a stale
value immediately and refresh it in the background (stale-while-revalidate), while
protecting the origin from stampedes.

## Requirements
**Functional**
- Get a key with a loader that fetches from the origin.
- On a fresh hit, return the cached value without calling the loader.
- On a stale hit, return the stale value immediately and refresh asynchronously.
- On a miss, load; coalesce concurrent loads for the same key.

**Non-functional**
- A stale-but-usable entry is served without blocking on the loader.
- Exactly one in-flight refresh/load per key under a stampede.
- Bounded staleness (a max-stale beyond which Get blocks on a fresh load).

## Given
- An origin that is slow and occasionally down; high read rate with hot keys.
- Single process (one edge node).

## Your core must
Expose a `Cache` with:
- `Get(key string, loader func() ([]byte, error)) ([]byte, error)` — return fresh from
  cache, or stale-now + refresh-in-background, or load on miss; coalescing concurrent
  loads for the same key.

Implement TTL + stale-while-revalidate + request coalescing. **Tests must show:** a fresh
entry is served from cache (loader not called); a stale entry is returned immediately while
a single background refresh runs; under a stampede of concurrent `Get`s for one missing key
the loader runs exactly once; (optional) serve-stale-on-error when the loader fails.

## Design should also address (in DESIGN.md)
- TTL vs stale-while-revalidate vs stale-if-error (RFC 5861) semantics.
- Request coalescing / singleflight to prevent origin stampedes.
- Background-refresh scheduling and how max-staleness is bounded.
- Negative caching and cache penetration for keys that don't exist.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
