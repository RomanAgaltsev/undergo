# Kata k09 — Circuit breaker

**Track:** A — Traffic & rate control · **Tier:** T1

## Problem
Design a circuit breaker that protects a caller from a failing downstream dependency:
fail fast while the dependency is unhealthy, then probe for recovery — the "stop
hammering a service that's already down" problem.

## Requirements
**Functional**
- Gate each call with an allow/deny decision.
- Record the outcome of a call (success or failure).
- Move through the states closed → open → half-open → closed as outcomes and time dictate.

**Non-functional**
- The gate decision is O(1) on the hot path.
- Concurrency-safe (many callers share one breaker).
- Bounded memory for the failure window regardless of traffic volume.

## Given
- A dependency called ~10k times/second; a burst of failures should open the breaker
  quickly.
- Single process; one breaker instance per dependency.

## Your core must
Expose a `Breaker` with:
- `Allow() bool` — may this call proceed right now?
- `Record(ok bool)` — report the outcome of a call that was allowed.

Model the closed/open/half-open states over a rolling failure window (count- or
time-bucketed), with a configurable failure threshold, open duration, and half-open
probe count. **Tests must show:** the breaker trips to *open* after the failure
threshold is crossed; while open, `Allow()` denies until the open duration elapses;
in *half-open* it admits only a limited number of probes and closes on success (or
re-opens on failure); correctness under concurrent callers.

## Design should also address (in DESIGN.md)
- Rolling-window design: a count-based window vs sliding time buckets, and the memory
  tradeoff.
- Failure threshold as a raw count vs a failure ratio (and the minimum-volume guard).
- The half-open probe policy — how many probes, and what closes vs re-opens the breaker.
- How the breaker composes with timeouts and retries, and why it isn't a substitute for
  either. Distributed/shared breaker state is out of scope for the core.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
