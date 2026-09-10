# Kata k12 — Write-ahead log

**Track:** B — Storage & indexing · **Tier:** T1

## Problem
Design a durable append-only log that survives process crashes — the building block that
sits under a database's durability guarantee and a message broker's segment file.

## Requirements
**Functional**
- Append a record and return its offset.
- Force durability on demand (sync).
- Replay every durably-appended record in order after a restart.

**Non-functional**
- Crash-safe: a partial/torn write at the tail is *detected*, never replayed as a valid
  record.
- Appends are sequential; per-record overhead is small and bounded.

## Given
- Records up to a few KB; ~50k appends/second; a single writer.
- One or more on-disk segment files.

## Your core must
Expose a `Log` with:
- `Append(rec []byte) (offset uint64, err error)` — append one record, return its offset.
- `Sync() error` — force buffered records to stable storage.
- `Replay(fn func(offset uint64, rec []byte) error) error` — call fn for each record in
  append order.

Frame each record with a length prefix + CRC over the payload. **Tests must show:** a
torn/partial write at the tail is detected on replay and *not* surfaced as a valid record;
a full replay after a simulated crash returns exactly the records that were durably
appended; offsets are monotonic.

## Design should also address (in DESIGN.md)
- Record framing (length + CRC) and exactly how a torn tail is detected and truncated.
- The fsync policy — per-append vs batched/group commit — and the durability-vs-throughput
  tradeoff.
- Segment rotation and recovering (truncating) a corrupt tail on open.
- How this relates to a broker's segment log (cf. ferrumq / riffle) and a database WAL.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
