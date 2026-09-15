# 03 — What a timer channel holds

`timers.go` asks two separate questions about timers, and two different Go
releases answer them.

The first three slots ask for the **capacity** of the channel you get from
`time.NewTimer`, `time.NewTicker`, and `time.After`.

The fourth builds a `Ticker`, never calls `Stop`, drops the last reference to
it, and asks whether the garbage collector took it.

| slot | question |
|---|---|
| `timer_cap` | `cap(t.C)` for a `*time.Timer` |
| `ticker_cap` | `cap(t.C)` for a `*time.Ticker` |
| `after_cap` | `cap` of the channel `time.After` returns |
| `dropped_ticker_collected` | was an un-`Stop`ped, unreferenced `Ticker` reclaimed? |

If all three capacity answers come out the same, that is either a real finding
or a sign you answered once and copied it. Be able to say which.

```
undergo verify sched/03-timer-channels
```

## Questions to answer in writing

1. A buffered timer channel lets a fire be delivered to a receiver that has
   already stopped caring. Say what bug that caused in practice, and why
   changing the buffer is a **behaviour** change rather than only a performance
   one.
2. `dropped_ticker_collected` is gated on the **`go.mod` go line**, not on the
   toolchain. Create a scratch module with `go 1.22`, run the same code under
   this toolchain, and predict what changes. Then say why Go gated it on the
   module rather than on the compiler.
3. Does calling `Stop` change the answer to question 4's slot? Answer in terms
   of *who holds a reference*, not in terms of politeness.
4. `time.After` inside a `select` in a loop was the classic leak. Say whether
   these changes fix it, and what it still costs while the timer is pending.
