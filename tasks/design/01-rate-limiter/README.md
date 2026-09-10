# Kata k1 — Rate limiter

## Problem
Design an in-process API rate limiter that caps how fast each client (identified by a key) may
make requests.

## Requirements
**Functional:** per-key limits; an allow/deny decision per request; configurable rate and burst.
**Non-functional:** goroutine-safe; low, predictable latency; memory bounded as the number of
distinct keys grows.

## Given
- ~100k active client keys.
- ~10k requests/second aggregate.
- Single process (distributed limiting is a design-discussion extension, not part of the core).

## Your core must
Expose a limiter with `Allow(key string) bool` (optionally `AllowN(key string, n int) bool`),
implement **both** a token-bucket and a sliding-window strategy behind a common interface, and
be safe under concurrent callers. Tests must demonstrate:
- a burst up to capacity is allowed, then further requests are denied;
- tokens/allowance refill over time;
- correctness under concurrent goroutines hitting the same key (run with `-race`).

## Design should also address
- Token-bucket vs sliding-window vs fixed-window: accuracy, burst behavior, memory.
- How you'd extend this to **distributed** rate limiting (shared counter, cell-based, etc.).
- The clock source — why monotonic time, and how testability drives injecting it.
- How memory stays bounded when keys are unbounded (eviction / TTL of idle buckets).

## Workflow
1. Write `DESIGN.md` against `../../../rubric/design-rubric.md`.
2. Get it graded; iterate until no dimension is "thin".
3. Build the core, test-first.
4. Write `REFLECTION.md`.
