# 01 — Join strings without allocating

`JoinInto` is correct and allocates on every call. Make it allocate **nothing**
per operation when the caller supplies a destination with enough capacity, while
still producing exactly what `strings.Join` would.

The graded assertion is `allocs/op == 0` from `TestOptimizeTarget`. The
correctness tests are frozen; do not weaken them.

```
undergo start alloc/01-zero-alloc-join
undergo verify alloc/01-zero-alloc-join
```

## Questions to answer in writing

1. The baseline builds a string with `+=` in a loop. How many allocations does
   that do for six parts, and why is it not six?
2. Why does `append(dst, s...)` with a string operand not allocate, when
   `[]byte(s)` does?
3. Your version writes into `buf[:0]` each iteration. What would go wrong if the
   caller kept the *previous* returned slice as well?
