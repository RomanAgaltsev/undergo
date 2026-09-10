# Kata k19 — Exactly-once dedup

**Track:** C — Messaging & streaming · **Tier:** T1

## Problem
Design a deduplication filter that suppresses duplicate message deliveries so downstream
processing is effectively exactly-once — keyed by message id, with a TTL to bound memory.

## Requirements
**Functional**
- Given a message id, report whether it has been seen before (and record it on first sight).
- Expire old ids after a TTL so memory stays bounded.

**Non-functional**
- Bounded memory under high id churn (TTL eviction).
- O(1) check.
- Concurrency-safe.

## Given
- Millions of ids/second; a dedup window measured in minutes.
- Single process.

## Your core must
Expose a `Dedup` with:
- `Seen(id string) bool` — true if the id was already seen; on first sight, record it and
  return false.

Include TTL-based eviction. **Tests must show:** first sight of an id returns false and
records it; a repeat within the TTL returns true (duplicate); after the TTL elapses the id
is evicted and treated as new; memory stays bounded under high churn.

## Design should also address (in DESIGN.md)
- An exact set + TTL vs a probabilistic filter (bloom / count-min — cf. k13 / k30) and the
  false-positive tradeoff.
- Sizing the TTL / window to cover the delivery system's redelivery window.
- How this composes with at-least-once delivery (k5 / k17) to give effective exactly-once.
- Sharding for concurrency under high id rates.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
