# 09 — What an atomic costs

`cost.go` adds one to an `int64` three ways — a plain `x++`, `atomic.Int64.Add`,
and `x++` under a `sync.Mutex` — and times each: first from one goroutine, then
with every P adding to the same value at once.

| slot | question |
|---|---|
| `atomic_costs_more_than_plain` | from one goroutine, does `Add` cost more per increment than `x++`? |
| `mutex_costs_more_than_atomic` | from one goroutine, does Lock, `x++`, Unlock cost more than `Add`? |
| `sharing_slows_atomic` | does each `Add` cost more when every P is adding to the same value than when one goroutine has it to itself? |

Answer `true` or `false`.

The judge is the clock. Each arm runs five times and its fastest run counts; the
test compares the per-increment costs and never prints them. A slot is graded
only when the two costs it compares are at least 1.5× apart — if a busy machine
squeezes a pair closer than that, the test says the run was inconclusive rather
than grade you against noise. It needs `GOMAXPROCS` of at least 2.

```
undergo verify concurrency/09-atomic-cost
```

## Questions to answer in writing

1. Before measuring: how many times more does an uncontended `Add` cost than
   `x++` on your machine? Then measure it (the hint says how) and account for the
   gap between your guess and the number.
2. There is no shared plain `x++` arm. Say why it would be cheaper than the
   shared `Add`, and what the race detector would say about it.
3. Read `sync.Mutex`'s `Lock` and `Unlock` fast paths. How many atomic operations
   does an uncontended Lock/Unlock pair cost, and does that match your measured
   mutex/atomic ratio?
4. In the shared arms no goroutine ever waits for another, yet each `Add` costs
   far more. Say what the cores are doing instead of waiting.
5. No slot compares the shared mutex with the shared atomic, because no single
   answer held: on a 4-vCPU CI runner they came out about equal, on a 12-thread
   desktop the mutex cost 6–13×. Measure yours. Then name two things a
   contended `Lock` does that an `Add` never does, and one thing that lets it
   keep up anyway.
