# Kata k17 — Pub/sub broker

**Track:** C — Messaging & streaming · **Tier:** T2

## Problem
Design a publish/subscribe broker that fans each message on a topic out to many
independent subscribers, each consuming at its own pace, with at-least-once delivery.

## Requirements
**Functional**
- Publish a message to a topic.
- Subscribe to a topic and receive its messages.
- Ack a delivered message; each subscriber has an independent cursor.

**Non-functional**
- At-least-once delivery.
- A slow subscriber does not stall or block the others.
- Acks advance that subscriber's cursor so acked messages aren't redelivered.

## Given
- Thousands of subscribers across topics; bursty publishing.
- Single process.

## Your core must
Expose a `Broker` with:
- `Publish(topic string, msg []byte)` — fan a message out to the topic's subscribers.
- `Subscribe(topic string) *Sub` — a new subscription with an independent cursor.

and a `Sub` with a receive method and `Ack(id uint64)`. **Tests must show:** a published
message reaches all current subscribers (at-least-once); a slow/blocked subscriber does
not prevent others from receiving; `Ack` advances the subscriber's cursor so acked
messages are not redelivered.

## Design should also address (in DESIGN.md)
- Per-subscriber queues vs a shared log + per-subscriber cursors (the fan-out memory
  tradeoff).
- Handling a slow consumer: buffer, drop, or backpressure — and the consequences.
- At-least-once vs at-most-once, and redelivery of unacked messages.
- Durable vs in-memory subscriptions; wildcard / hierarchical topics as an extension.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
