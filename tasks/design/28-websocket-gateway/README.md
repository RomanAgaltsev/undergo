# Kata k28 — Websocket gateway registry

**Track:** F — Real-time & fan-out · **Tier:** T2

## Problem
Design the in-memory registry at the heart of a websocket gateway: it holds every live
connection, maps users to their connections, and routes a message to one user or
broadcasts to all — while a slow connection must never stall the hub or other clients.

## Requirements
**Functional**
- Register / unregister a live connection (a user may have several).
- Send a message to a specific user (all their connections).
- Broadcast a message to every connection.

**Non-functional**
- A targeted send reaches only the target user's connections.
- A slow / blocked connection is shed (dropped or its buffer bounded), not allowed to
  block the hub or other clients.
- Concurrency-safe under registrations, sends, and disconnects happening at once.

## Given
- Hundreds of thousands of concurrent connections; a skewed distribution across users.
- Each connection has a bounded send buffer; consumers read at varying speeds.

## Your core must
Expose a `Hub` with:
- `Register(conn *Conn)` / `Unregister(conn *Conn)` — manage the connection set.
- `Send(user string, msg []byte)` — deliver to a user's live connections.
- `Broadcast(msg []byte)` — deliver to all connections.

Shard the registry by user and apply **per-connection backpressure** (bounded buffer).
**Tests must show:** a targeted `Send` reaches only that user's connections and no others;
a slow connection (full buffer) is dropped/shed without blocking `Send`/`Broadcast` to the
rest; `Broadcast` reaches every registered connection.

## Design should also address (in DESIGN.md)
- The registry data model (user → set of conns) and why it's sharded (lock contention).
- Backpressure policy on a full send buffer: drop-oldest, drop-newest, or disconnect.
- The per-connection writer goroutine and clean unregister (no goroutine/connection leak).
- Cross-node routing: how a multi-gateway deployment finds the node holding a user's conn
  (cf. presence k26) — sketch, not build.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
