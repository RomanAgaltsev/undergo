# Kata k21 — Leader election

**Track:** D — Coordination & consensus · **Tier:** T2

## Problem
Design leader election over a small cluster: the nodes agree on exactly one leader,
re-elect promptly when the leader fails, and never allow two leaders at once (no
split-brain) — the coordination primitive under Raft, Kubernetes leases, and friends.

## Requirements
**Functional**
- Each node advances on a tick and participates in elections.
- Report a node's current role (follower / candidate / leader).
- Elect a new leader after the current one stops.

**Non-functional**
- At most one leader at steady state (safety).
- Prompt re-election after leader loss (liveness within bounded ticks).
- No split-brain under the modeled message loss.

## Given
- A small cluster (3 or 5 nodes) with an abstract, lossy message channel.
- Randomized timeouts are available to break symmetry.

## Your core must
Expose a `Node` with:
- `Tick()` — advance the node's clock/state one step.
- `State() Role` — the node's current role.

Run a lease/term-based election (Raft-style terms + randomized election timeouts) over a
small cluster abstraction. **Tests must show:** at steady state exactly one node is leader;
after the leader stops ticking, a new leader is elected within bounded rounds; under the
modeled message loss, two nodes are never leaders in the same term (no split-brain).

## Design should also address (in DESIGN.md)
- Terms/epochs and the "one leader per term" invariant; why randomized timeouts avoid livelock.
- Quorum/majority votes and how they prevent split-brain during a partition.
- How this is the election half of **k7 Raft**; where it stops short of full consensus.
- Detecting leader failure (heartbeat/lease) and the availability gap during re-election.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
