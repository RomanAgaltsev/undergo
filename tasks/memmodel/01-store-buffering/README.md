# 01 — The outcome that should be impossible

Two goroutines, four statements:

```
goroutine A: x = 1; r1 = y
goroutine B: y = 1; r2 = x
```

Write out every interleaving of those four statements. You will find `(1,1)`,
`(0,1)` and `(1,0)`. You will not find `(0,0)` — that would need both goroutines
to read a value from before the other's write, and there is no ordering of four
statements in which that happens.

`sb.go` runs this half a million times, in two variants: once with ordinary
variables, once with `sync/atomic`.

| slot | question |
|---|---|
| `plain_00_observed` | did `(0,0)` ever happen with ordinary variables? |
| `atomic_00_observed` | did `(0,0)` ever happen with atomics? |

Before running it, write down all four outcomes and mark each **legal** or
**illegal** for both variants. The slots ask about *observation*; the questions
below ask about *legality*, and they are not the same question.

```
undergo verify memmodel/01-store-buffering
```

## Questions to answer in writing

1. List the four outcomes and say which are legal under the Go memory model for
   the plain variant. Then say whether "legal" and "observed" are the same
   question, and which one this task grades.
2. The atomic variant makes one outcome illegal. Name the property Go's atomics
   have that C++'s `memory_order_relaxed` does not, and say what it costs.
3. This program contains a data race on purpose. Say what `-race` does to it,
   and why the task sets `requires.default_build`.
4. `(0,0)` is observable on amd64, which is the *strongly* ordered
   architecture. Name the hardware structure responsible, and say why a machine
   that does not reorder loads and stores still needs a fence.
