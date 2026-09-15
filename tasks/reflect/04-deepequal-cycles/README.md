# 04 — Equality that terminates

Write a structural equality check:

```go
func Equal(a, b any) bool
```

It must follow pointers, slices, maps and structs, compare **contents** rather
than addresses, and terminate even when either value contains a cycle.

Match `reflect.DeepEqual`'s semantics, because that is what the tests compare
you against: a nil slice and an empty slice are **not** equal, and neither are a
nil map and an empty map.

## How it is graded

| test | requirement |
|---|---|
| `TestAgainstDeepEqual` | eighteen acyclic cases must agree with `reflect.DeepEqual` exactly |
| `TestTerminatesOnSelfReference` | two identical one-element rings are equal, and the call returns |
| `TestCyclicAgainstAcyclic` | a ring is **not** equal to a chain that looks the same for one step |
| `TestLongerCycle` | two identical two-element rings are equal |

The cycle tests run `Equal` in a goroutine with a five-second budget, so a
non-terminating implementation fails rather than hanging.

```
undergo verify reflect/04-deepequal-cycles
```

## Questions to answer in writing

1. The naive recursive version does not terminate on a cyclic value. Say what
   state you added to fix it, and where in the recursion you consult it.
2. What is that state keyed on? Try keying it on a **single** pointer instead of
   a pair and run the tests. Say which test catches it, and explain why keying
   on one side alone gives wrong *answers* rather than merely being slower.
3. When your check finds a pair it is already comparing, it has to return
   something. Say what it returns and why that is correct rather than merely
   convenient — what is the assumption, and what eventually discharges it?
4. `reflect.DeepEqual` says a nil slice and an empty slice differ. Argue for and
   against that choice, then say which you would want in a test-assertion
   library and why.
5. Say when you would reach for this in production code, and what you would
   reach for instead.
