# 02 — The result nobody received

A worker goroutine computes something and sends it back. The caller has already
given up — it timed out, or its context was cancelled, and it moved on.

What happens to the worker?

`abandoned_test.go` starts five workers, each abandoned the moment it is
launched, and then asks the runtime's own goroutine profile whether each one is
still there. Predict the five answers.

| slot | the worker's body | the channel |
|---|---|---|
| `unbuffered_leaks` | `ch <- 1` | `make(chan int)` |
| `capacity_one_leaks` | `ch <- 1` | `make(chan int, 1)` |
| `select_with_default_leaks` | `select { case ch <- 1: default: }` | `make(chan int)` |
| `capacity_one_two_results_leaks` | `ch <- 1; ch <- 2` | `make(chan int, 1)` |
| `late_drain_leaks` | `ch <- 1`, and somebody else receives later | `make(chan int)` |

`true` means the goroutine is still parked when the profile is taken — it leaked.

Four of these five are designs you have written. One of them is the fix people
reach for, and one of them is that same fix applied to a case it does not cover.

```
undergo verify concurrency/02-abandoned-result
```

## Questions to answer in writing

1. For the first case, say exactly which runtime queue the worker is sitting in
   and what would have to happen for it to ever be scheduled again.
2. A buffer of one changes the answer. Say what the buffer is doing here in one
   sentence, and why calling it "a performance tweak" is wrong.
3. Case four has the same buffer and still leaks. State the general rule that
   case two and case four are both instances of.
4. Case three does not leak, but it is not free. Say what it costs, and name the
   situation where paying that cost is wrong.
5. `time.After(d)` returns a channel nobody may ever read from. Look up its
   capacity, and say why it is what it is.
6. Case five does not leak because a second goroutine drains the channel. Say
   why that is a fix you should be suspicious of, given that it is the one most
   likely to appear in a code review as "just drain it".
