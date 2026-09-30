# 10 — Once, and the second caller

`once.go` has the `Once` everybody writes first:

```go
func (o *Naive) Do(f func()) {
	if o.done.CompareAndSwap(false, true) {
		f()
	}
}
```

It can never run `f` twice. `Race` puts it — and `sync.Once` — in the one
situation where that is not enough: a first caller is inside `f`, blocked, and a
second caller arrives.

| slot | question |
|---|---|
| `naive_runs_f_once` | with two callers, does `Naive` call a function exactly once? |
| `naive_second_returns_early` | does `Naive`'s second caller return while the first is still inside `f`? |
| `once_second_returns_early` | does `sync.Once`'s? |

The judge is `Race`: it blocks the first caller's `f` on a channel, starts the
second caller, and watches for it to return for up to a second before releasing
`f`.

```
undergo verify concurrency/10-once-early-return
```

## Questions to answer in writing

1. `Naive` runs `f` once. Quote the sentence in `sync.Once`'s documentation that
   it breaks anyway.
2. Say what a caller of `Naive.Do` can observe when it returns early, and why the
   race detector would catch that only on some runs.
3. Read `sync/once.go`. What does `done` mean there that it does not mean in
   `Naive`, and why can the fast path still be a single atomic load?
4. Fix `Naive` without a mutex. Say what you need besides the one flag, and what
   a caller who finds `f` still running has to do.
5. This task waits up to a second for a caller *not* to return. Say why it
   cannot use `testing/synctest` to avoid the wait, and why a loaded machine can
   slow it down but never flip an answer.
