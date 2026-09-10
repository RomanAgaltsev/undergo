# Kata k36 — Durable scheduler / delay queue

**Track:** H — Application systems · **Tier:** T2

## Problem
Design a durable scheduler: register a task to fire at a future time, and have it fire at
or after that time even across a process restart — a delay queue / timer wheel with leases
on fired tasks. The engine behind delayed jobs, retries, and reminders.

## Requirements
**Functional**
- Schedule a task to fire at a given time; return a handle.
- Poll for tasks that are now due.
- Cancel a scheduled task before it fires.

**Non-functional**
- Durable: scheduled tasks survive a restart.
- Fires at or after the due time — never early.
- A due task handed to a poller is leased so it isn't double-fired within the lease.

## Given
- Millions of scheduled tasks; due times from seconds to days out.
- Multiple pollers (workers) draining due tasks concurrently.

## Your core must
Expose a `Sched` with:
- `Schedule(at time.Time, task Task) (id string)` — register a future task.
- `Poll() []Task` — the tasks now due, each handed out under a lease.
- `Cancel(id string) bool` — cancel a task before it fires.

Use a timer wheel / delay queue + a durable store + a lease on fired tasks. **Tests must
show:** a task fires at/after its due time and not before; scheduled tasks survive a
simulated restart (rebuilt from the durable store); a leased due task is not handed to a
second poller within its lease; cancelling before the due time prevents the fire.

## Design should also address (in DESIGN.md)
- Timer wheel vs a delay-queue heap: the tradeoff at scale and for far-future due times.
- Durability of the schedule (append-only log — cf. k12 WAL, ferrumq / riffle).
- The lease / visibility timeout so a crashed worker's task is redelivered (cf. k5 queue).
- At-least-once vs exactly-once firing and idempotent tasks (cf. k19 dedup).

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
