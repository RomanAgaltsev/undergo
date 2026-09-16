# 08 — Fan-out, fan-in, and closing the cascade

Implement four functions in `pipeline.go` that compose into a cancellable
pipeline:

```
Gen ──► FanOut ──┬─► Square ──┐
                 ├─► Square ──┼──► FanIn ──► you
                 ├─► Square ──┤
                 └─► Square ──┘
```

| function | contract |
|---|---|
| `Gen(ctx, vs...)` | emit each of `vs`, then close. Stop early on cancel, **close in every case** |
| `Square(ctx, in)` | read until `in` closes or `ctx` is cancelled, emitting `v*v`; close on the way out |
| `FanOut(ctx, in, n)` | `n` channels fed from `in`; each value goes to exactly one; all closed when `in` closes or `ctx` is cancelled |
| `FanIn(ctx, ins...)` | merge; close once every input is drained or `ctx` is cancelled. Order across inputs is unspecified |

## What the tests check

Nine tests. The two that will decide whether you have understood this:

- **`TestCancelWithoutDrainingLeaksNothing`** — the consumer cancels and walks
  away. Nobody drains. Every goroutine in the graph must still unwind.
- **`TestStageSendUnblocksOnCancelAlone`** — one `Square` holding a value that
  nobody will ever take. Only `ctx` can free it.

The second exists because the first is not enough. Work out why before you read
the explanation: it is the most useful thing in this task.

```
undergo verify concurrency/08-close-cascade
```

## Questions to answer in writing

1. "The sender closes" is the rule. Give the reason — not the rule restated, but
   what closing *asserts*, and why only the sender is in a position to assert it.
2. `FanIn` has several senders feeding one channel. Say who closes it, and why
   that answer needs a `sync.WaitGroup` rather than a counter.
3. What does a receiver cause if it closes the channel it is reading? Name the
   specific failure and say why it is worse than the alternatives.
4. A stage must handle cancellation on **both** its receive and its send. Say
   what goes wrong with only the receive guarded, and construct the smallest
   graph that exposes it.
5. `TestCancelMidStreamLeaksNothing` drains after cancelling.
   `TestCancelWithoutDrainingLeaksNothing` does not. State what draining fixes,
   and what it hides.
6. Every channel here is unbuffered. Say what changes about cancellation
   behaviour if `Square`'s output has a buffer of 100, and whether it makes the
   leak more or less likely to be caught by a test.
