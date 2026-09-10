# Kata k20 — Distributed lock

**Track:** D — Coordination & consensus · **Tier:** T2

## Problem
Design a distributed lock with crash safety: at most one holder at a time, the lock is
freed automatically if the holder crashes (lease expiry), and every acquisition hands out
a strictly increasing **fencing token** so a stalled old holder can't corrupt shared state
after its lease has expired.

## Requirements
**Functional**
- Acquire the lock with a TTL; return a fencing token on success.
- Renew (extend) the lease while still held.
- Release the lock.

**Non-functional**
- Mutual exclusion: never two live holders at once.
- A crashed holder's lease expires and frees the lock without manual intervention.
- Fencing tokens strictly increase, so a stale holder is fenced out at the resource.

## Given
- Many contending clients; clients can pause (GC/STW) or crash at any time.
- Clocks are imperfect; a holder may not renew in time.

## Your core must
Expose a `Lock` with:
- `Acquire(ttl time.Duration) (token uint64, ok bool)` — grant the lock + a fencing token.
- `Renew(token uint64) bool` — extend the lease if still the holder.
- `Release(token uint64)` — release the lock.

Use a lease + a monotonic fencing token. **Tests must show:** two acquires can't both hold
the lock (mutual exclusion); a holder that stops renewing loses the lock after the TTL and
another client can acquire; each successful acquire returns a strictly larger fencing token.

## Design should also address (in DESIGN.md)
- Why a lease + fencing token, and the classic "GC pause after acquire" failure that fencing fixes.
- Where the lock state lives (single authority vs a quorum — cf. k7 Raft) and its availability.
- Clock assumptions: monotonic vs wall-clock, and the safety margin on the TTL.
- Reentrancy, fairness, and thundering-herd on release as extensions.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
