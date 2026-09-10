# Kata k3 — Consistent-hash load balancer

## Problem
Design a load balancer / partition router that maps keys to a changing set of backend members
with minimal remapping when membership changes.

## Requirements
**Functional:** add/remove members at runtime; map any key to a member; even load distribution.
**Non-functional:** **minimal key movement** (~1/N of keys) when a member joins or leaves;
low, predictable lookup latency on the hot path.

## Given
- Up to ~100 members.
- Millions of keys.
- `Pick` is called on the request hot path.

## Your core must
Expose a ring with `Add(member string)`, `Remove(member string)`, and `Pick(key string) string`,
using virtual nodes per member. Tests must demonstrate:
- load is spread across members within a tolerance (e.g. no member gets wildly more than 1/N);
- adding or removing one member remaps only a small fraction of keys (close to 1/N), i.e. keys
  not near the change keep their member.

## Design should also address
- The virtual-node count tradeoff: distribution variance vs memory/lookup cost.
- Consistent hashing vs **rendezvous (HRW) hashing** — when each wins.
- Consistent hashing with bounded loads (capping any member's share).
- The hash function choice and why (distribution, speed, not cryptographic).

## Workflow
1. Write `DESIGN.md` against `../../../rubric/design-rubric.md`.
2. Get it graded; iterate until no dimension is "thin".
3. Build the core, test-first.
4. Write `REFLECTION.md`.
