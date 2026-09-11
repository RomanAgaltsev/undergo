# 01 — The exact capacity sequence, and where doubling stops

`Sequence(600)` appends 600 ints to a nil slice one at a time and records every
distinct capacity `append` chose along the way.

Read `growth.go`. Without running anything, predict:

| slot | question |
|---|---|
| `first_cap` | the first capacity in the sequence |
| `third_cap` | the third capacity in the sequence |
| `cap_after_256` | the capacity that immediately follows 256 |
| `cap_after_512` | the capacity that immediately follows 512 |
| `distinct_caps` | how many distinct capacities appear in total |

Two of these are where the naive model breaks. If you answered `cap_after_512`
with 1024, you have made two separate mistakes and it is worth finding both.

```
undergo verify types/01-cap-growth
```

Answers are for the **default build** of the pinned toolchain. That caveat is
load-bearing here, and question 3 is about why.

## Questions to answer in writing

1. Above a few hundred elements the growth factor is no longer 2. What is it,
   roughly, and why was it changed?
2. The capacities are not exactly `oldcap × factor` either. What rounds them,
   and to what? Work out what `cap_after_512` would have been with no rounding.
3. Run the sequence twice more:

   ```
   go test -gcflags=all=-d=variablemakehash=n -count=1 .
   ```

   and then again with `Sink` deleted from `growth.go`. Each changes the answer.
   Explain which part of the sequence each one changes, and say which of the
   three results you would call "how append grows".
