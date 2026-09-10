# Kata k2 — Unique ID generator

## Problem
Design a service that hands out unique, roughly time-ordered 64-bit IDs across many nodes with
no central coordinator on the hot path.

## Requirements
**Functional:** globally unique; k-sorted (roughly monotonic by time); no coordination per ID.
**Non-functional:** high throughput; no collisions; survives clock skew and regression without
emitting duplicates.

## Given
- Up to ~1024 nodes.
- ~10k IDs/second/node.
- Each ID fits in a signed 64-bit integer.

## Your core must
Expose a generator with `Next() (int64, error)` using a Snowflake-style bit layout
(timestamp + node + sequence), with a caller-supplied node ID. Tests must demonstrate:
- uniqueness across many IDs generated concurrently;
- monotonicity within a node;
- sequence rollover within a single millisecond (block until the next ms rather than collide);
- defined behavior on clock regression (return an error, or wait — your choice, tested).

## Design should also address
- The bit-allocation tradeoff: time range vs node count vs sequence/ms, and what each buys.
- Snowflake vs UUIDv4 vs a database ticket server — when each is right.
- The clock-regression strategy and why (error vs wait vs logical clock).
- What "roughly ordered" (k-sorted) means and why strict global order is expensive.

## Workflow
1. Write `DESIGN.md` against `../../../rubric/design-rubric.md`.
2. Get it graded; iterate until no dimension is "thin".
3. Build the core, test-first.
4. Write `REFLECTION.md`.
