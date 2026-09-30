# 06 — What a waiting goroutine costs

Goroutines are called cheap. This task asks how cheap, and what they cost once
they have done some work.

`testdata/park` parks goroutines on a channel and measures the runtime's stack
memory before and after, in four phases:

| phase | the goroutines |
|---|---|
| `parked` | 10 000, each waiting on a channel straight away |
| `deep` | 1 000, each waiting at the bottom of a 32-deep recursion through 1 KiB frames |
| `returned` | 1 000 that made the same recursion, came back up, and then wait |
| `unreachable` | 10 000 waiting on a channel that nothing but they can reach |

| slot | question |
|---|---|
| `stack_kib_per_goroutine` | stack per `parked` goroutine, in KiB: `<1`, `1-10`, `10-100` or `>100` |
| `grows_with_stack_depth` | does a `deep` goroutine hold more than twice the stack of a `parked` one? |
| `shrinks_after_return` | a few collections later, does a `returned` goroutine hold less than half the stack of a `deep` one? |
| `parked_reclaimed_by_gc` | after two collections, has the runtime reclaimed any of the `unreachable` goroutines? |

The judge is the runtime's own accounting (`runtime.MemStats.StackInuse`) and
goroutine count. Every phase runs in a fresh child process, so one phase cannot
leave anything behind for the next. The test never prints a size.

```
undergo verify sched/06-goroutine-size
```

## Questions to answer in writing

1. Your bracket holds on every OS; the exact figure does not. Find the starting
   stack size for linux and for Windows in `runtime/stack.go`, and say why
   Windows reserves more.
2. An OS thread's stack is typically 1–8 MiB. Using your measured figure, how
   many parked goroutines fit in one thread's stack?
3. `shrinks_after_return` needed more collections than the other phases. Say
   when the runtime shrinks a stack, by how much at a time, and why that is the
   collector's job.
4. Since Go 1.19 a new goroutine's first stack is not always the minimum. Find
   `adaptivestackstart` in the runtime and say why this task runs each phase in
   its own process.
5. The `unreachable` goroutines can never run again. Say why the collector keeps
   them anyway, and what `sched/05` offers instead.
