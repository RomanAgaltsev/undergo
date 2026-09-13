# 03 — Merging two sequences, and stopping them

Two things to build:

**`(*Node).InOrder`** yields a binary search tree's values in sorted order, and
stops walking as soon as its consumer stops — including partway down the left
spine, where a plain recursive walk has no way to say so.

**`Merge(a, b)`** yields two sorted sequences in sorted order. This is the one
that needs `iter.Pull`: a merge has to hold one element from each side at the
same time to decide which comes next, and two push sequences cannot both be in
control.

```go
for v := range pull.Merge(left.InOrder(), right.InOrder()) {
	...
}
```

## The contract

- `InOrder` is sorted, handles a nil tree, and stops early when asked.
- `Merge` is sorted, keeps duplicates, and handles an empty side (either one, or
  both).
- When the consumer stops a `Merge` early, **both** sources are stopped.
- No goroutine outlives the merge.

```
undergo verify iter/03-pull
```

The last two are one requirement wearing two hats, and they are the point of the
task. A `Merge` that gets every value right and skips them leaves a goroutine
parked for the lifetime of the process, once per call.

## Questions to answer in writing

1. `iter.Pull` returns two functions. Say what the second one releases, and what
   is still running if you never call it.
2. Write `Merge` as two `for ... range` loops instead, and say exactly where it
   stops being writable. What is it about a push sequence that makes two of them
   impossible to drive together?
3. `next` returns `(v, ok)`. Distinguish `ok == false` from a sequence that
   genuinely yielded a zero value, and say how a caller tells them apart.
4. The leak test counts goroutines. Name the thing it would catch that
   `TestMergeStopsBothSources` would not, and the thing that test catches which
   counting goroutines would not.
