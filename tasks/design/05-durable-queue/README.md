# Kata k05 — Durable message queue

**Track:** C — Messaging & streaming · **Tier:** T2

## Problem
Design a durable, replayable message queue: messages are persisted in order and consumed
by offset, so a consumer can restart and resume, or replay from an earlier point — the
log at the heart of Kafka and similar brokers.

## Requirements
**Functional**
- Append (produce) a message and return its offset.
- Read messages starting from a given offset.
- Commit a consumer's offset; a consumer resumes from its last committed offset.

**Non-functional**
- Durable — messages and committed offsets survive a restart.
- Ordered within the log.
- Append and sequential read are efficient (append-only).

## Given
- Millions of messages of a few KB each; a single log/partition.
- Multiple consumers reading at independent offsets.

## Your core must
Expose a `Queue` with:
- `Append(msg []byte) (offset uint64, err error)` — persist a message, return its offset.
- `Read(offset uint64, max int) ([]Message, error)` — up to max messages from offset.
- `Commit(consumer string, offset uint64) error` — record a consumer's progress.
- `Committed(consumer string) uint64` — a consumer's resume offset.

Persist with an append-only log + per-consumer offsets. **Tests must show:** messages
survive a restart and replay in order; a consumer resumes from its committed offset and
uncommitted messages are redelivered; committing advances the offset so committed messages
are not redelivered.

## Design should also address (in DESIGN.md)
- The append-only segment log + index (cf. the k12 WAL, and ferrumq / riffle).
- Where per-consumer offsets are stored and how they're made durable.
- At-least-once vs exactly-once delivery, and where a dedup filter (k19) fits.
- Partitioning for parallel consumers; retention / compaction of old segments.
- You may reuse patterns from **ferrumq / riffle** — reference them, don't depend on them.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
