# Kata k11 — Fair scheduler (deficit round-robin)

**Track:** A — Traffic & rate control · **Tier:** T2

## Problem
Design a scheduler that shares one constrained resource fairly across many competing
flows (tenants/clients) in proportion to configured weights, so no single flow can
starve the others — the "noisy neighbour" problem.

## Requirements
**Functional**
- Enqueue an item tagged with the flow it belongs to.
- Dequeue the next item, honouring each flow's fair share.
- Support per-flow weights (a heavier flow gets proportionally more throughput).

**Non-functional**
- Amortised O(1) work per dequeue (deficit round-robin), not O(flows).
- No starvation: any flow with queued items eventually makes progress.
- Work-conserving: never idle while any queue is non-empty.
- Bounded memory per active flow.

## Given
- Hundreds of flows with bursty arrivals; flows come and go.
- Single process; the scheduler feeds one worker/resource.

## Your core must
Expose a generic `Scheduler[T]` with:
- `Enqueue(flow string, item T)` — queue an item under its flow.
- `Dequeue() (T, bool)` — the next item to serve (zero value + false when empty).

Implement deficit round-robin (DRR): each active flow carries a deficit counter that
accrues a per-flow quantum (scaled by weight) each round. **Tests must show:** over a
run, per-flow throughput is proportional to weights within a tolerance; no flow starves
while it has queued items; a bursty/heavy flow cannot monopolise the output (others keep
making progress); the scheduler is work-conserving.

## Design should also address (in DESIGN.md)
- DRR vs weighted fair queueing (WFQ / virtual-time) vs weighted round-robin — the
  fairness-accuracy vs O(1) tradeoff.
- Quantum sizing and its effect on latency and fairness granularity.
- Handling flow churn: adding/removing active flows from the round without unfairness.
- Whether to key on `string` flows or a generic key type, and item ownership.

## Workflow
Write `DESIGN.md` → have it graded against `../../../rubric/design-rubric.md` → build the
core TDD (real, tested Go) → write `REFLECTION.md`.
