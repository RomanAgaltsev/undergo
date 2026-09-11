# 04 — Why boxing 255 is free and boxing 256 is not

`Box` does one thing: it puts an `int` into an `any`. Whatever that costs is the
cost of the conversion.

Without running anything, predict:

| slot | question |
|---|---|
| `allocs_0` | allocations per call of `Box(0)` |
| `allocs_255` | the same for `Box(255)` |
| `allocs_256` | the same for `Box(256)` |
| `allocs_100000` | the same for `Box(100000)` |
| `first_allocating` | the **smallest** non-negative int whose boxing allocates |

The last slot is the point of the task. "Small integers are cheap" is not an
answer; a number is.

```
undergo verify alloc/04-boxing-cache
```

## Questions to answer in writing

1. Two `any` values boxing the same `int` — are they `==`? Does your answer
   change above the threshold? Is what you observe a language guarantee or an
   artefact of the implementation?
2. Does the same saving apply to a `byte`, an `int8`, a `rune`, a `uint64` of the
   same value? Predict before you measure.
3. The threshold is exactly what it is because of the element type of a table in
   the runtime. Name the table, and explain why *that* size is the boundary.
