# 05 — The carry Go cannot reach

Add two 128-bit numbers, each held as a pair of `uint64`:

```go
func Add128(loA, hiA, loB, hiB uint64) (lo, hi uint64)
```

`add128_amd64.s` already has a body. It adds the low words, throws the carry
away, and always reports a high word of zero. Fix it.

The point is what you need to reach for. Adding the low words can set a carry,
and the high words must add that carry in. **Go has no expression for the carry
flag** — there is no operator, no builtin, and no way to ask whether the last
addition overflowed except by comparing afterwards.

## How it is graded

- `TestCarryPropagates` — the one case the stub gets wrong, named so it fails
  first
- `TestEdges` — every combination of 0, 1, all-ones, all-ones−1 and 1<<63
- `TestAgainstReference` — ten thousand random inputs from a fixed seed,
  compared **bit for bit** against `math/bits.Add64`

No timing, no benchmark. The grade is whether your answer is identical to the
reference, on every input.

```
undergo verify asm/05-carry-chain
```

## Questions to answer in writing

1. Name the two instructions you used. Say what the second one reads that the
   first one wrote, and why no Go expression can name it.
2. `math/bits.Add64` exists and compiles to exactly this pair. Say what writing
   it yourself taught you that reading its documentation would not — and say
   how you would have written a correct `Add128` in pure Go without it.
3. The `TEXT` directive says `$0-48`. Show the arithmetic, and say where the two
   results live.
4. Extend this to 256 bits in your head. Say how many carry-adds you need, and
   why the chain cannot be parallelised across the words — then say what that
   implies for big-integer performance generally.
