# 02 — Time that costs nothing

`bubble.go` holds three ordinary functions. They sleep, they run a ticker, they
ask what year it is. Nothing in them knows anything about testing.

The frozen test calls all three inside `synctest.Test`, which runs them in a
*bubble*. Predict what they report from in there.

| slot | question |
|---|---|
| `slept_minutes` | what `SleepTwice(1h, 30m)` measured, in whole minutes |
| `ticks` | what `TickCount(1s, 10.5s)` counted |
| `year` | what `Year()` returned |
| `fake_exceeds_real` | whether the time the bubble reported sleeping exceeds the real time the test took |

Two of these four are the answers you would give without the bubble. Two are
not. Decide which before you run it.

Note the ticker runs for **10.5 seconds**, not 10. That half-second is not
decoration — question 5 is about why it is there.

```
undergo verify sched/02-synctest-bubble
```

## Questions to answer in writing

1. The test sleeps ninety minutes of bubble time and finishes in milliseconds.
   Say what has to be true of **every** goroutine in the bubble before the clock
   is allowed to jump, and name the term for that state.
2. Name one operation that will hang a bubble forever instead of advancing its
   clock, and say why the runtime cannot tell it apart from ordinary work.
3. What does `synctest.Wait()` do that `time.Sleep(0)` does not? Say when you
   would need it.
4. The year is 2000. Say what breaks in code that uses `time.Now()` to build an
   ID or seed a generator, and argue whether that is a feature of the bubble or
   a hazard of it.
5. Change the ticker's total from `10500*time.Millisecond` to exactly
   `10*time.Second` and run the test twenty times. The count is no longer
   stable, inside a bubble whose whole purpose is determinism. Say exactly what
   is undecided, and say which kind of non-determinism a fake clock removes and
   which kind it does not.
