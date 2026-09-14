# 01 — Which goroutine runs first

`runnext.go` starts four goroutines from one goroutine, numbered 1 to 4, under
`GOMAXPROCS=1`. Nothing yields between the `go` statements. Each goroutine
appends its own number to a slice as soon as it runs.

There is exactly one P, so there is exactly one answer. Predict it.

| slot | question |
|---|---|
| `first` | the number the first goroutine to run appended |
| `second` | the second |
| `third` | the third |
| `fourth` | the fourth |

If your four answers are `1 2 3 4`, you have assumed a queue. If they are
`4 3 2 1`, you have assumed a stack. It is neither.

```
undergo verify sched/01-runnext-ordering
```

## Questions to answer in writing

1. Exactly one of the four is out of position. Name the structure it sits in,
   and say why the runtime has it at all — what does it cost to put a newly
   spawned goroutine at the back of the queue instead?
2. The spawner does not yield. When, precisely, does it stop running and let
   the first spawned goroutine start? Name the call.
3. Raise `GOMAXPROCS` to 4 and the answer stops being stable. Say what becomes
   true that was not true before — and why that makes this task pin
   `GOMAXPROCS` rather than ask about the general case.
4. A goroutine in that structure is not stealable by another P for a short
   window. Say why that matters for a two-goroutine ping-pong, and what it
   would cost if it were stealable immediately.
5. Run this task's test under `go test -race`. Two of the four slots change.
   Say which two, why *those* two, and what that tells you about using the race
   detector to investigate a scheduling question. This program contains no data
   race at all — so what is the detector changing?

