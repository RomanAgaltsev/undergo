# Kata k07 — Replicated KV (Raft)

**Track:** D — Coordination & consensus · **Tier:** T3

## Problem
Design a consistent, fault-tolerant key-value store replicated across a small cluster:
writes are agreed by a consensus protocol (Raft) so the store stays linearizable and
survives a minority of node failures — including a leader crash. The hard core of a
replicated state machine.

## Requirements
**Functional**
- Put and Get key-value pairs against the cluster.
- Replicate each write through a leader to a majority before it's committed and applied.
- Elect a new leader and keep serving after the old leader fails.

**Non-functional**
- Linearizable reads and writes on the happy path.
- Tolerates a minority of node failures (a majority quorum stays available).
- A committed write is never lost across a leader change.

## Given
- A small cluster (3 or 5 nodes); an unreliable network (message loss/delay is modeled).
- Far more reads than writes, but writes must be durable and agreed.

## Your core must
Expose a replicated `KV` over a small `Cluster` of `Node`s:
- `Put(key, value string) error` — replicate and apply a write via the leader.
- `Get(key string) (value string, ok bool, err error)` — a linearizable read.

Back it with a **Raft** log (leader election + log replication + apply to the state
machine). **Tests must show:** a value written before a leader failover is still readable
after the new leader is elected (agreement under leader change); happy-path reads/writes
are linearizable; a write is committed only once replicated to a majority.

## Design should also address (in DESIGN.md)
- The Raft pieces: leader election (cf. k21), log replication, commit index, and apply.
- Why a majority quorum, and the safety vs availability behaviour under a partition.
- Linearizable reads: read-index / lease reads vs routing reads through the log.
- Log compaction / snapshotting to bound the log; membership changes as an extension.
- You may reuse patterns from **ferrumq** (log replication) and pair this with **MIT 6.824**
  — reference them, don't depend on them.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
