# 02 — Ranking four collectors

`Cycles(gogc, 48)` runs the same workload — a 2 MiB live set, 48 MiB of
immediately-dead garbage — at four values of GOGC, and counts the collections
each one performed.

**This task does not ask you how many.** It asks for the order.

| slot | question |
|---|---|
| `ordering` | the four settings ranked by collections performed, most first, joined by `>` — e.g. `gogc100>gogc50>gogc800>gogc400` |
| `gogc50_beats_gogc800` | `true` or `false`: does GOGC=50 collect more often than GOGC=800? |

```
undergo verify gc/02-cycle-ordering
```

## Why it does not ask for a count

Because there is not one. Running this workload ten times on one idle machine
gave, at GOGC=100:

    39, 41, 46, 53, 59, 60, 60, 65, 67, 68

Holding the live set constant — which `Cycles` does — narrowed it to `16–20` and
no further. There is no tightening of the workload that makes it exact.

The **ordering**, measured the same way, came out identical every time.

So the count is not a property of your program. The ordering is. Getting that
distinction right is worth more than any number this task could have asked you
for, and it is the same distinction you will need for anything you ever measure
in production.

## Questions to answer in writing

1. Name three things that can move the collection count between two runs of an
   identical program on an idle machine.
2. The count is not reproducible and the ordering is. Say what that tells you
   about which of the two is a property of the program, and give one other
   measurement from your own experience that has the same shape.
3. You are asked to halve a service's GC CPU and you have exactly this data.
   Say which number you would put in the ticket, and what you would measure
   instead of cycle count.
4. `GOGC=off` is a fifth setting. Say where it ranks, and whether that is a
   useful answer or a degenerate one.
5. GOGC=800 collects far less often than GOGC=50. Name what you pay for that, and
   the deployment where the trade is wrong.
