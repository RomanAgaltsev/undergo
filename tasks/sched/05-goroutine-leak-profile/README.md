# 05 — Which leaks does the runtime name

`leaks.go` plants one goroutine leak of each classic shape:

| function | the leak |
|---|---|
| `BlockedSend` | send on an unbuffered channel nobody receives from |
| `NilRecv` | receive from a nil channel |
| `WaitGroupNeverDone` | `Wait` on a counter nothing decrements |
| `MutexNeverUnlocked` | `Lock` a mutex the caller is still holding |
| `ContextNeverDone` | wait on a context that has no cancellation |

Then it asks Go 1.27's `goroutineleak` profile which of them it can name, and
reads `Count()` **both before and after** `WriteTo`.

| slot | question |
|---|---|
| `blocked_send_named` … `context_named` | did the profile name that shape? |
| `count_before_writeto` | what `Count()` returned before `WriteTo` ran |
| `count_after_writeto` | what `Count()` returned after |

Two of those slots are the same profile, queried twice, with nothing in between
but a call that writes to a buffer. **If you predict the same number for both,
say why you expected that before you run it.**

```
undergo verify sched/05-goroutine-leak-profile
```

## Questions to answer in writing

1. One of the five shapes is not really its own shape. Say which, and what it
   reduces to.
2. `count_before_writeto` and `count_after_writeto` differ. Say where the work
   happens, and why an API would be built that way. Then say what a monitoring
   endpoint that calls only `Count()` would report, and how long that bug would
   survive in production.
3. A goroutine blocked on a mutex the caller still holds is different from one
   that is merely slow. Say whether the runtime can distinguish "leaked" from
   "waiting a long time", and what property it actually decides on.
4. `Survey` pins `GOMAXPROCS` to 1 and calls `runtime.Gosched` ten times before
   profiling. Both are deliberate and Go's own test corpus does the same. Say
   what each is for.
5. Read how `Survey` decides whether a shape was named. It matches on something
   other than the package name. Find out what a profile actually calls a
   function, and say what the obvious matcher would have reported instead —
   and why that wrong answer would have looked entirely plausible.
