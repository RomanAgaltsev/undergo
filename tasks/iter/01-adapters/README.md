# 01 — Five adapters over a sequence

Implement `Map`, `Filter`, `Take`, `Zip` and `Chunk` over `iter.Seq` and
`iter.Seq2`. Each is four or five lines. The difficulty is not the
transformation — it is that an adapter is a **consumer** of the sequence it
wraps and a **producer** for the sequence consuming it, and it has to honour the
yield contract in both directions.

```go
for v := range adapters.Take(adapters.Map(src, double), 3) {
	...
}
```

## The contract

- Every adapter stops pulling from its source as soon as its own consumer stops.
- `Take(seq, n)` **stops** the source after `n` — it does not filter afterwards,
  so `Take` over an infinite sequence terminates.
- `Take(seq, n)` for `n <= 0` yields nothing and never touches the source.
- `Zip` stops at the shorter side and does not drain the longer one.
- `Chunk` emits a short final chunk, and `size <= 0` yields nothing.
- Chunks do not share a backing array with one another.

```
undergo verify iter/01-adapters
```

If the test run hangs instead of failing, that is a result too: something is
draining a sequence that never ends.

## Questions to answer in writing

1. What does `yield` returning `false` mean, and whose responsibility is it to
   act on it? Trace one `break` in the example above through all three layers.
2. `Chunk` cannot reuse a single buffer, even though it would be faster and the
   sizes test would still pass. Say precisely what breaks, and name the standard
   library function whose documentation makes the same promise `Chunk` has to.
3. Four of these adapters are written with `for ... range seq`. `Zip` cannot be.
   Explain what `Zip` needs that a `range` loop cannot give it.
4. `Take` stops its source by returning from inside the range body. What would
   happen instead if it counted down and simply stopped calling `yield`, leaving
   the loop running?
