# 05 — What a byte slice really costs

`make([]byte, n)` asks for `n` bytes. The allocator charges something else.

| slot | expression |
|---|---|
| `n1` | `make([]byte, 1)` |
| `n8` | `make([]byte, 8)` |
| `n9` | `make([]byte, 9)` |
| `n17` | `make([]byte, 17)` |
| `n33` | `make([]byte, 33)` |
| `n600` | `make([]byte, 600)` |
| `n1000` | `make([]byte, 1000)` |
| `n32768` | `make([]byte, 32768)` |
| `n32769` | `make([]byte, 32769)` |
| `one_pointer` | `make([]*int, 1)` |

Each is the `B/op` the benchmark reports for one allocation.

**These ten numbers are not produced by one rule.** There are three regimes and
two boundaries. Find the boundaries before you predict anything.

```
undergo verify alloc/05-size-classes
```

## Questions to answer in writing

1. One byte costs one byte, and nine bytes cost sixteen. Say which allocator
   handles the first case, and why charging the exact request is honest there
   rather than a rounding error.
2. `make([]*int, 1)` and `make([]byte, 8)` both cost eight bytes. Say why they
   arrive at the same number by different routes, and what would change if the
   element type were `*int` at length 2.
3. 32768 costs 32768 and 32769 costs 40960. Name the threshold, the unit it
   rounds to, and the fraction of that allocation that is wasted.
4. You now know the cost curve. Say what it implies for a buffer you grow by
   appending, and name the `append` behaviour that interacts with it.
5. The measurement imports `testing` from a non-test file. Say what
   `AllocedBytesPerOp` is actually reading, and why that is a fair way to
   answer this question.
