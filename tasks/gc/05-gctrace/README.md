# 05 — Reading what the collector says

`GODEBUG=gctrace=1` makes the runtime write a line per collection to stderr:

```
gc 1 @0.001s 4%: 0+0.52+0 ms clock, 0+0.52/1.0/0+0 ms cpu, 3->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 12 P
```

Parse it.

```go
func Parse(trace string) ([]Cycle, error)
func TotalForced(cs []Cycle) int
```

`testdata/trace.txt` is a **recorded** trace from a real program — fifteen
collections, two of them forced by `runtime.GC`. It is a fixture rather than a
live run on purpose: GC timing is not reproducible, and `gc/02` and `gc/03` in
this repository exist because of it. A recording is.

## How it is graded

| test | requirement |
|---|---|
| `TestParsesEveryCollection` | fifteen cycles, numbered 1 to 15 |
| `TestIgnoresNonCollectionLines` | the scavenger line and the `GC forced` banner are not cycles |
| `TestFirstCycleFields` | every field of the first line, including `WallMS` |
| `TestForcedCollections` | two forced, and the right two |
| `TestMalformedLineIsAnError` | a broken collection line is an error, not a zero `Cycle` |
| `TestEmptyTrace` | an empty trace is not an error |

The second test is worth reading before you start. One of the lines it is about
contains the word **forced**.

```
undergo verify gc/05-gctrace
```

## Questions to answer in writing

1. The heap section reads `3->4->0 MB`. Say what each of the three numbers is,
   and which one `GOGC` is calculated from.
2. One line is marked `(forced)`. Say what forces a collection, and why the
   pacer treats a forced one differently from one it scheduled.
3. The percentage is cumulative since program start, not per cycle. Say what
   that makes it useless for, and what you would use instead to answer "is the
   collector costing me too much *right now*".
4. You have just parsed a debug format with no compatibility promise. Say what
   you would do differently for a tool you intended to keep, and name the
   package that offers a supported alternative — then say what that package
   cannot tell you that this line can.
