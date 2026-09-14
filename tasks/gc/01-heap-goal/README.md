# 01 — What GOGC actually sets

GOGC does not set a collection frequency. It sets the **heap goal**: the size the
heap is allowed to reach before the next collection starts. The documented rule
is

    goal = live x (1 + GOGC/100)

so at GOGC=100 the collector aims to let the heap reach twice the live set.

`goal.go` measures the real goal against that rule, over a 4 MiB live set and
again over a 64 MiB one.

| slot | question |
|---|---|
| `formula_gogc_100` | what the formula predicts at GOGC=100, to one decimal |
| `large_heap_ratio_gogc_100` | the measured goal ÷ live over 64 MiB, at GOGC=100 |
| `large_heap_ratio_gogc_400` | the same at GOGC=400 |
| `small_heap_exceeds_formula` | over 4 MiB, is the real ratio *above* the formula at all three settings? |
| `large_heap_matches_formula` | over 64 MiB, is it within 2% of the formula at all three? |
| `gap_shrinks_as_heap_grows` | is the distance from the formula smaller on the large heap? |

The first three have one answer between them if the formula is the whole story.
The last three are where you find out whether it is.

```
undergo verify gc/01-heap-goal
```

Each measurement forces **two** collections before reading the goal. Written
question 3 is about why one is not enough.

## Questions to answer in writing

1. The formula holds on a large heap and not on a small one. Name at least two
   things counted toward the goal that are not your live heap, and say why they
   matter proportionally more when the heap is small.
2. From your answer to 1, predict what happens to the ratio on a heap of 100 KiB.
   Then say what that means for a small, short-lived process — a CLI, a Lambda —
   that never allocates much.
3. The measurement forces two collections and reads the goal after the second.
   Say what the goal read after the *first* one is based on, and what a task that
   read it then would be measuring.
4. `GOGC=off` sets no goal at all. Say what still triggers a collection, and what
   `GOMEMLIMIT` has to do with it.
5. Suppose a service's live set doubles in a single cycle. Using the formula, say
   what the goal is for the cycle *after* that, and what that means for peak RSS.
