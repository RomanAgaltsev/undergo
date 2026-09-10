# Kata k27 — Notification fan-out

**Track:** F — Real-time & fan-out · **Tier:** T2

## Problem
Design a notification fan-out engine: one logical notification for a user must be
delivered to each of that user's registered devices exactly once, retrying transient
delivery failures with backoff — without ever delivering a duplicate to a device.

## Requirements
**Functional**
- Register / unregister a device for a user.
- Push a notification for a user; it targets all of that user's devices.
- Deliver queued notifications, retrying transient failures.

**Non-functional**
- Per-device de-duplication: a device receives a given notification at most once.
- Transient failures are retried with backoff; permanent failures give up after a bound.
- No duplicate delivery when a retry follows a partial success.

## Given
- Users with several devices each; a transport that can fail transiently (timeouts).
- Bursty pushes (a broadcast to many users at once).

## Your core must
Expose a `Fanout` with:
- `Register(user, device string)` — add a device to a user (and an unregister counterpart).
- `Push(user string, notif Notif)` — enqueue a notification (with a stable id) for all the
  user's devices.
- `Deliver()` — attempt delivery of pending items, applying retry/backoff.

Track per-device delivery state and dedup by `(device, notif-id)`. **Tests must show:**
a notification reaches each of the user's devices exactly once; a device that failed
transiently is retried (with increasing backoff) and eventually delivered; a retry after
a partial success does **not** re-deliver to already-delivered devices.

## Design should also address (in DESIGN.md)
- The dedup key & where per-device delivery state is stored (and its TTL / cleanup).
- Retry/backoff policy (exponential + jitter) and the give-up bound; dead-letter handling.
- At-least-once transport + idempotent delivery ⇒ effective exactly-once (cf. k19).
- Fan-out to a huge audience: batching and avoiding a thundering herd on the transport.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
