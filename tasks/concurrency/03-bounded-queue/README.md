# 03 — A bounded queue on sync.Cond

Implement `Queue[T]` in `queue.go`: a bounded FIFO where `Push` blocks while the
queue is full, `Pop` blocks while it is empty, and `Close` unblocks everyone.

Use `sync.Cond`. A channel would be shorter and is the right answer in almost
every real program — that is question 4.

## The contract

| method | blocks while | returns `ErrClosed` when |
|---|---|---|
| `Push(v)` | the queue is full | closed on entry, or closed while waiting for room |
| `Pop()` | the queue is empty | closed **and drained** |
| `Close()` | never | — |

Read the second row of the `Pop` column twice. **A closed queue still yields what
it already holds.** A `Pop` that checks for closure before checking for items is
wrong, and `TestCloseDrainsRemaining` will say so.

`Close` may be called more than once, and concurrently with everything else.

## What the tests are actually checking

Ten tests. Eight of them confirm the obvious. Two do not, and they are the
reason this task exists:

- **`TestCloseWakesAllWaiters`** — four goroutines are parked in `Pop`, then
  `Close` runs once. All four must come back.
- **`TestWokenLosersKeepWaiting`** — three goroutines are parked in `Pop`, then a
  single `Push` arrives. One of them gets the value. The other two must still be
  waiting, and neither may report the queue closed, because it is not.

If your implementation passes the other eight and fails these two, you have
written the two classic `sync.Cond` bugs, one each.

```
undergo verify concurrency/03-bounded-queue
```

## Questions to answer in writing

1. `sync.Cond.Wait` must be called with the lock held, and it releases the lock
   before parking. Say what is true of the shared state at the instant `Wait`
   returns — and therefore what `Wait` returning does and does not tell you.
2. `TestWokenLosersKeepWaiting` fails for one specific one-word change to a
   correct solution. Name the change and explain the failure.
3. `TestCloseWakesAllWaiters` fails for a different one-word change. Name it,
   and say why no later event rescues the goroutines it strands.
4. Rewrite the queue with a buffered channel instead — actually write it, it is
   about fifteen lines. Then state the one thing your `Cond` version does that
   the channel version cannot, and decide whether you would ship it.
5. This solution can use one `Cond` for both waiting sides, or one for each.
   Describe the cost of the single-`Cond` version under load, and name the term
   for what happens to the waiters that were woken for nothing.
