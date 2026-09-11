# 02 — How many allocations does each of these do?

Five small functions in `snippets.go`. Each one's allocation count differs from
its neighbours' for a *different reason*.

Read them. Without running anything, predict `allocs/op` for each:

| slot | function |
|---|---|
| `concat_two` | `ConcatTwo("alpha", "beta")` |
| `sprintf_int` | `SprintfInt(42)` |
| `itoa_int` | `ItoaInt(42)` |
| `builder_join` | `BuilderJoin("alpha", "beta")` |
| `make_local` | `MakeLocal()` |

Counts are whole numbers, as reported by `testing.AllocsPerRun`. The measurement
assigns each result to a typed package-level variable, so what you are counting
is the function's own allocations and not the cost of boxing its result.

```
undergo verify alloc/02-allocs-per-op
```

Two of these five are 0 and they are 0 for unrelated reasons. One of them is
higher than almost everyone guesses. If you get four right and one wrong, the
wrong one is the one worth writing up.

## Questions to answer in writing

1. `ItoaInt(42)` and `SprintfInt(42)` produce the same string. Account for every
   allocation the second does that the first does not — and note that the
   difference is smaller than the folklore claims. Why?
2. `BuilderJoin` uses the type that exists specifically to avoid allocations, and
   is not the cheapest here. What is it doing, and what one-line change would
   make it the cheapest?
3. `MakeLocal` asks for 64 bytes and gets its count anyway. What single edit
   would change that count, and what is the compiler analysis that decides it?
