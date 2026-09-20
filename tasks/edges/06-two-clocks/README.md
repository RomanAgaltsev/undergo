# 06 — The time that is not equal to itself

`time.Now()` does not return one reading. It returns two: a **wall clock**,
which NTP or an administrator can move in either direction, and a **monotonic**
reading, which cannot go backwards. Both live in the same `time.Time`.

Some operations carry the monotonic reading forward. Some drop it. And two
`Time` values naming the very same instant are not necessarily equal.

Read `clocks.go`. One instant is taken, `ref`, and `stripped` is `ref.Round(0)`.
Predict `true` or `false` for each:

| slot | question |
|---|---|
| `now_carries_monotonic` | does `time.Now()` carry a monotonic reading? |
| `add_keeps_monotonic` | does `ref.Add(time.Second)` still carry one? |
| `utc_keeps_monotonic` | does `ref.UTC()`? |
| `round_zero_keeps_monotonic` | does `ref.Round(0)`? |
| `json_round_trip_keeps_monotonic` | does a `Time` survive JSON with one? |
| `equal_method_says_same` | `ref.Equal(stripped)` |
| `double_equals_says_same` | `ref == stripped` |
| `sub_of_stripped_pair_matches` | is `ref.Sub(stripped)` exactly zero? |

Two of those last three disagree with each other about the same pair of values.
That disagreement is the task.

```
undergo verify edges/06-two-clocks
```

The judge is the standard library: the test asks each question directly and
grades the answer. A monotonic reading is visible only through `Time.String`,
which appends ` m=±<seconds>` when one is present.

## Questions to answer in writing

1. Say what the monotonic reading is *for* — name the concrete failure it
   prevents, in a program that measures how long something took.
2. `Equal` and `==` disagree about `ref` and `stripped`. State what each of them
   compares, and say which one you should use on a `time.Time`, always.
3. Give the rule for which operations strip the monotonic reading. Your rule
   should predict `UTC`, `Round(0)` and `AddDate` without special cases.
4. A `time.Time` goes into a map key, or into a struct compared with `==`, or
   into a `cmp.Diff` in a test. Say what goes wrong, and say when it goes wrong
   — always, or only sometimes.
5. `Round(0)` looks like a no-op rounding. Say why that is the documented way to
   strip the monotonic reading, and what it says about the API that this is the
   spelling.
