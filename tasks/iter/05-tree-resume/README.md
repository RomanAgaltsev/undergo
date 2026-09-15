# 05 — A tree you can stop walking

Two functions over a binary search tree.

```go
func All(t *Tree) iter.Seq[int]
func Walk(t *Tree) (next func() (int, bool), stop func())
```

`All` is a **push** sequence: it hands values to a loop, in order. `Walk` is a
**pull** sequence: the caller asks for the next value when it wants one, and
can stop in between.

## How it is graded

| test | requirement |
|---|---|
| `TestAllIsInOrder` | every value, in order |
| `TestAllStopsDescending` | one `break` stops the whole walk — not just the frame that saw it |
| `TestAllEmptyTree` | a nil tree yields nothing |
| `TestWalkPullsInOrder` | the pull sequence agrees with the push one |
| `TestWalkResumes` | take three, then carry on from the fourth |
| `TestStopIsIdempotentAndFinal` | `stop` twice is safe; `next` reports false afterwards |
| `TestStopReleasesTheWalk` | twenty half-consumed walks, stopped, leave nothing behind |

The second and the last are the ones that fail if you take the obvious route.

```
undergo verify iter/05-tree-resume
```

## Questions to answer in writing

1. A recursive in-order walk has to propagate `yield`'s result up through every
   frame. Show the line people forget, and say exactly what the tree does
   without it — how many values does the sequence produce after a `break`?
2. `iter.Pull` turns a push sequence into a pull one. Something has to be
   suspended half-way through for that to work. Say what Go uses, and what it
   costs per `Pull`.
3. `TestStopReleasesTheWalk` starts twenty walks and stops them. Delete the
   `stop()` call from that loop, add `runtime.GC()`, and see what happens to
   the goroutine count. Say what that means for a `Pull` in a request handler.
4. Following from that: `next` returning `false` naturally also releases the
   walk. Say when you must call `stop` anyway, and write the `defer` you would
   use.
5. Compare this with `sched/05` in this repository. Say which of its leak
   shapes an abandoned `Pull` resembles, and whether the `goroutineleak`
   profile would name it.
