# Kata k14 — LSM-tree KV engine

**Track:** B — Storage & indexing · **Tier:** T2

## Problem
Design a write-optimized embedded key-value store — the log-structured merge-tree engine
that sits under LevelDB, RocksDB, and Cassandra.

## Requirements
**Functional**
- Get / Put / Delete by key.
- Ordered Scan over a key range.
- Deletes via tombstones (a delete shadows older values).

**Non-functional**
- Writes land in an in-memory memtable, then flush to immutable sorted files.
- Reads merge the memtable and the sorted files, newest wins.
- Compaction reclaims space and drops shadowed keys and obsolete tombstones.

## Given
- Millions of keys; write-heavy; values up to a few KB.
- Single process; on-disk sorted files.

## Your core must
Expose an `Engine` with:
- `Get(key []byte) ([]byte, bool)` — newest value, or not-found.
- `Put(key, val []byte)` and `Delete(key []byte)` (tombstone).
- `Scan(from, to []byte) Iterator` — ordered, deduplicated key range.

Build a memtable (skiplist or ordered map) + immutable sorted SSTables + a compaction
step. **Tests must show:** read-your-writes across the memtable and flushed SSTables
(newest wins); a deleted key is not returned even when an older SSTable still holds it;
after compaction the shadowed/older versions and dropped tombstones are gone and `Scan`
returns the merged, ordered, deduplicated view.

## Design should also address (in DESIGN.md)
- Memtable choice (skiplist vs balanced tree) and the flush trigger (size/count).
- SSTable layout (sorted blocks + sparse index) and how a bloom filter (k13) lets `Get`
  skip SSTables that can't contain the key.
- Compaction strategy — size-tiered vs leveled — and its write / read / space
  amplification tradeoffs.
- The WAL (k12) needed to make the memtable durable (name it even if outside the core).

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
