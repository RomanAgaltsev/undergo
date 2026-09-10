# Kata k10 — Adaptive concurrency limiter

**Track:** A — Traffic & rate control · **Tier:** T2

## Problem
Design a limiter that bounds the number of concurrent in-flight requests based on
*observed latency* rather than a hand-tuned fixed number — so the system sheds load
before it collapses, and grows the limit when there's headroom.

## Requirements
**Functional**
- Acquire a permit before starting work; if the system is at its limit, shed (deny).
- Release the permit when work finishes, reporting the observed latency.
- Raise the concurrency limit when latency is low, lower it when latency rises.

**Non-functional**
- Concurrency-safe.
- The limit never exceeds a configured hard ceiling.
- The controller converges without wild oscillation.

## Given
- A service whose safe concurrency varies with downstream latency; ~10k req/s.
- Single process; one limiter guarding the resource.

## Your core must
Expose a `Limiter` with:
- `Acquire() (Permit, bool)` — returns a permit and true, or the zero permit and false
  when at the limit (the caller sheds).
- and a `Permit` with `Release(latency time.Duration)` — return the permit and feed the
  observed latency back to the controller.

Implement an AIMD or gradient controller (Little's-law style): the limit increases while
latency stays near its baseline and decreases as latency climbs. **Tests must show:**
under low/steady latency the limit grows toward the ceiling; under rising latency the
limit shrinks and `Acquire` begins shedding; the limit never exceeds the hard ceiling;
correctness under concurrent callers.

## Design should also address (in DESIGN.md)
- AIMD vs a gradient controller (à la Netflix concurrency-limits) vs a direct Little's-law
  estimate — tradeoffs.
- How "rising latency" is detected: current RTT against a baseline / minimum RTT.
- Whether to queue or shed when at the limit, and the latency cost of each.
- Why a fixed concurrency limit is fragile as downstream conditions change.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
