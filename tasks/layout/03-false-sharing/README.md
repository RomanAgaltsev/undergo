# 03 — Which counters share a cache line

Eight workers, each incrementing its own counter. No two workers touch the same
counter, so there is no contention — and yet one of these two layouts is an order
of magnitude slower than the other.

Read `counters.go`. Without running anything, predict:

| slot | question |
|---|---|
| `sizeof_naive` | `unsafe.Sizeof(Naive{})` |
| `sizeof_padded` | `unsafe.Sizeof(Padded{})` |
| `naive_per_line` | how many `Naive` counters fit in one 64-byte line |
| `naive_neighbours_share` | do `NaiveCounters[0]` and `[1]` fall on the same line? `true`/`false` |
| `padded_neighbours_share` | the same question for `PaddedCounters` |

```
undergo verify layout/03-false-sharing
```

The throughput difference is **not** graded — a ratio depends on how many cores
you have, and you cannot be exactly right about it. It is the subject of the
first written question instead, and the benchmarks are here so you can measure
it:

```
go test -bench . -cpu 8
```

## Questions to answer in writing

1. Before running them, predict the ratio between `BenchmarkNaive` and
   `BenchmarkPadded` on your machine. Then run both and account for the gap
   between your guess and the measurement.
2. The benchmark increments through a plain `*uint64` with no lock and no
   atomic. Why does that not change the shape of the result — and why would you
   never ship it?
3. `Padded` spends 56 bytes of padding on every counter. At what number of
   workers does that stop being worth it, and what would you measure to find out
   rather than guess?
