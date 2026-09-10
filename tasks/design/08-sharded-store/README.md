# Kata k08 — Sharded + replicated store (capstone)

**Track:** H — Application systems · **Tier:** T3

## Problem
The capstone. Design a horizontally-scalable, fault-tolerant key-value store: keys are
**sharded** across many replica groups, each group is **replicated** for durability and
availability, and a client's Put/Get is **routed** to the right shard. This is the
composition of three earlier katas — the k3 consistent-hash ring (routing), the k5 durable
log (per-replica persistence), and the k7 Raft replication (per-shard agreement) — into one
system, the way real stores (Cassandra, TiKV, CockroachDB) are built.

## Requirements
**Functional**
- Put and Get key-value pairs against the cluster.
- Route each key to its shard (a replica group) deterministically.
- Replicate each shard's writes so a shard survives a minority replica failure.

**Non-functional**
- Scales horizontally: adding shards adds capacity, with minimal key remapping.
- A write is durable and linearizable within its shard.
- A minority replica loss in a shard is tolerated; a whole-shard loss degrades only that
  key range, not the whole store.

## Given
- Many shards, each a small (3- or 5-node) replica group; more keys than one node holds.
- An unreliable network (message loss/delay modeled); reads outnumber writes.

## Your core must
Expose a `Store` that composes the earlier cores:
- `Put(key, value string) error` — route to the key's shard, then replicate + apply there.
- `Get(key string) (value string, ok bool, err error)` — route + linearizable read.

Compose: a **k3-style consistent-hash ring** to map a key → shard; a **k7-style replicated
KV** (Raft) per shard for agreement; a **k5-style durable log** per replica for persistence.
**Tests must show:** a key always routes to the same shard, and keys spread across shards
(routing); a write is durable across a minority replica failure within its shard (replicated
durability); losing a shard's replicas affects only that shard's key range, and the rest of
the store keeps serving (behavior on shard/replica loss).

## Design should also address (in DESIGN.md)
- Composition boundaries: what the ring (k3), the per-shard Raft (k7), and the per-replica
  log (k5) each own, and the interfaces between them.
- Rebalancing when shards are added/removed: how much data moves (cf. k3's ~1/N remapping)
  and how it moves without downtime.
- Cross-shard operations: why single-key ops stay simple and what a multi-key/transaction
  would cost (2PC / no cross-shard atomicity in scope).
- Routing state: how a client/coordinator learns the shard map and handles a stale map.
- The first bottleneck as the cluster grows, and which axis scales vs which doesn't.
- You may build on your own **k3 / k5 / k7** cores (or reference them) and echo **ferrumq** —
  reference, don't hard-depend; the capstone is the *integration*, not re-deriving Raft.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
