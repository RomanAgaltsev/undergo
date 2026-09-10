# Kata k24 — Write-back cache

**Track:** E — Caching & data access · **Tier:** T2

## Problem
Design a cache that absorbs writes in memory and flushes them to a slow backing store
lazily and in batches — trading a small window of durability risk for far higher write
throughput.

## Requirements
**Functional**
- Get and put by key.
- Track dirty (unflushed) entries.
- Flush dirty entries to the backing store, coalescing repeated writes to the same key.

**Non-functional**
- Reads see the latest written value (read-your-writes), flushed or not.
- Repeated writes to one key before a flush produce a single backing write.
- A bounded staleness window (max entries/time before a forced flush).

## Given
- Write-heavy; a backing store cheap in batches but expensive per write.
- Single process.

## Your core must
Expose a `Cache` with:
- `Get(key string) ([]byte, bool)` — latest value (including unflushed writes).
- `Put(key string, val []byte)` — write into the cache, marking the entry dirty.
- `Flush() error` — write all dirty entries to the backing store (coalesced), clearing dirt.

Track a dirty set and coalesce. **Tests must show:** after `Put` then `Flush`, the backing
store holds the latest value; multiple `Put`s to one key before a `Flush` cause a single
backing write; read-your-writes before flush; the staleness bound triggers/forces a flush.

## Design should also address (in DESIGN.md)
- Write-back vs write-through vs write-around — durability vs latency tradeoffs.
- The dirty-set and coalescing design.
- The crash-loss window and how a WAL (k12) would bound it.
- Evicting a dirty entry (it must be flushed first).

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
