# 04 — The loop that used to hang

`testdata/spin` starts a goroutine that asks for a garbage collection, then
runs a tight loop with **no function call, no allocation, and no channel
operation** in it. `GOMAXPROCS` is 1, so there is exactly one processor and the
two goroutines must take turns — if they can.

The program prints whether the collection finished *while the loop was still
running*. The frozen test runs it twice: once normally, once with
`GODEBUG=asyncpreemptoff=1`.

| slot | question |
|---|---|
| `gc_completes_normally` | did the collection finish during the spin, with the default settings? |
| `gc_completes_preemptoff` | did it finish during the spin, with `asyncpreemptoff=1`? |

One `GODEBUG` setting, two different worlds. Predict both.

```
undergo verify sched/04-async-preemption
```

## Questions to answer in writing

1. The loop contains no call, no allocation and no channel operation. Name the
   thing the compiler normally inserts that the runtime uses to stop a
   goroutine, and say why this loop has none of them.
2. With the default settings the runtime stops the loop anyway. Say what it
   actually *does* to interrupt it, and why that mechanism needs the goroutine's
   registers to be describable at an arbitrary instruction — not just at a call.
3. This was the behaviour of every Go program before 1.14. Say what symptom a
   production service showed, and why it looked like a garbage-collector problem
   rather than a scheduling one.
4. Read `testdata/spin`'s doc comment. It explains why the program samples an
   atomic from inside the loop instead of using `select` over a timeout. Say
   what would go wrong with the timeout version, and connect it to what you
   learned in `sched/02` about what a `select` does with two ready cases.
5. `GOMAXPROCS(1)` is deliberate. Say what would mask this effect at higher
   values, and whether that makes the problem gone or merely hidden.
