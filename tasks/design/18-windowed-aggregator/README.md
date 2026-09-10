# Kata k18 — Windowed stream aggregator

**Track:** C — Messaging & streaming · **Tier:** T2

## Problem
Design a stream aggregator that groups an unbounded event stream into time windows and
emits an aggregate as each window closes — handling out-of-order (late) events via a
watermark. The core of Flink / Kafka Streams windowing.

## Requirements
**Functional**
- Push events, each carrying an event-time and a value.
- Emit the aggregate for each window as it closes.
- Support tumbling, sliding, and session windows.

**Non-functional**
- Bounded memory: windows past the watermark are emitted and evicted.
- Late events within the allowed lateness still update their window.
- Deterministic window boundaries.

## Given
- A high event rate with out-of-order events and bounded lateness.
- Single process.

## Your core must
Expose an `Agg` with:
- `Push(event Event)` — where `Event` carries an event-time and a value.
- `Emit() []Window` — the windows that have closed, each with its aggregate.

Support tumbling / sliding / session windowing with a watermark. **Tests must show:**
events fall into the correct window boundaries; a late event before the watermark updates
its window; a session window merges events within the gap and closes after it; windows
past the watermark are emitted and their state evicted.

## Design should also address (in DESIGN.md)
- Tumbling vs sliding vs session window semantics.
- Event-time vs processing-time, and the watermark that bounds lateness.
- State size and eviction after windows close.
- Incremental aggregation (running sum/count) vs buffering raw events; exactly-once emission.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
