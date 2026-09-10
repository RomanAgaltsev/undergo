# Kata k22 — Gossip membership (SWIM)

**Track:** D — Coordination & consensus · **Tier:** T3

## Problem
Design decentralized failure detection and membership: with no central coordinator, nodes
gossip to detect failures and converge on a shared member list — the SWIM protocol behind
Serf / HashiCorp memberlist and Cassandra-style clusters.

## Requirements
**Functional**
- Maintain a per-node view of cluster membership.
- Detect a failed member and propagate that across the cluster.
- Let a node join and leave the cluster.

**Non-functional**
- Failure detected within a bounded number of protocol rounds.
- Membership views converge across nodes (eventual agreement).
- Bounded false-positive rate under a given message-drop rate.

## Given
- Hundreds of nodes; a lossy network; no central registry.
- Each node can only talk to a few peers per round (bounded fan-out).

## Your core must
Expose a `Cluster` with:
- `Members() []Member` — the node's current view of live members.

Drive it with SWIM: direct **ping**, **ack**, **indirect probe** via k random peers, and a
**suspicion** timeout before declaring a member dead, with membership updates piggybacked on
messages. **Tests must show:** a stopped member is detected within bounded rounds; membership
converges across nodes after a change; the false-positive (wrongly-suspected) rate stays
bounded under the modeled drop rate.

## Design should also address (in DESIGN.md)
- Direct vs indirect probing and why indirect probes cut false positives.
- The suspicion mechanism (suspect → confirm/refute) and its timeout sizing.
- Dissemination: piggybacked gossip vs a separate anti-entropy round; convergence time.
- Incarnation numbers to refute stale "dead" rumors; scalability of the per-round fan-out.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
