# Kata k26 — Presence tracker

**Track:** F — Real-time & fan-out · **Tier:** T2

## Problem
Design a presence system that tracks who is online and, per room/channel, who is
currently present — driven by client heartbeats with a TTL, so a client that goes away
(closed tab, dropped connection) is marked offline without an explicit logout.

## Requirements
**Functional**
- Record a heartbeat for a user.
- Report whether a user is currently online.
- Join / leave a room, and list a room's currently-online members.

**Non-functional**
- Offline detection is timely: a user with no heartbeat within the TTL is offline.
- Bounded memory: expired presence state is reclaimed.
- Concurrency-safe under many concurrent heartbeats.

## Given
- Millions of connected users; heartbeats every few seconds; a TTL of tens of seconds.
- Users are members of many rooms; rooms range from tiny to very large.

## Your core must
Expose a `Presence` with:
- `Beat(user string)` — record/refresh a user's heartbeat.
- `Online(user string) bool` — true if the last heartbeat is within the TTL.
- `Join(room, user string)` / `Leave(room, user string)` — room membership.
- `Room(id string) []string` — the room's currently-**online** members.

Use heartbeat + TTL expiry. **Tests must show:** a user is online right after a beat and
offline once the TTL elapses without one; `Room` returns exactly the online members
(joined, not expired, not left); join/leave update membership correctly.

## Design should also address (in DESIGN.md)
- Where the TTL clock lives: lazy expiry on read vs a sweeper, and the memory/latency tradeoff.
- Sharding presence state by user for concurrency at scale.
- Fan-out of presence *changes* to interested subscribers (cf. k27 / k28) — in scope vs out.
- Cross-node presence (a user connected to one gateway, queried from another): sketch the approach.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
