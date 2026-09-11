# 02 — Who still shares the backing array?

One backing array, four views, six operations. `Run` does this:

```go
base       := []int{0, 1, 2, 3, 4, 5, 6, 7}
twoIndex   := base[2:4]
threeIndex := base[2:4:4]
twoIndex    = append(twoIndex, 99)
grown      := append(threeIndex, 77)
copy(base[1:], base[:3])
```

Without running anything, predict the state at the end:

| slot | question |
|---|---|
| `cap_two_index` | `cap(base[2:4])` |
| `cap_three_index` | `cap(threeIndex)` |
| `base_4` | `base[4]` |
| `base_1` | `base[1]` |
| `base_3` | `base[3]` |
| `grown_0` | `grown[0]` |
| `len_two_index` | `len(twoIndex)` |

```
undergo verify types/02-aliasing
```

`grown_0` is the one to think hardest about: `grown` came from `threeIndex`,
which is a view of `base`, and the last line rewrites `base`. Whether that
reaches `grown` is the whole question.

## Questions to answer in writing

1. Which single character in `base[2:4:4]` prevents the mutation that
   `base[2:4]` allows, and what is that form called?
2. `copy(base[1:], base[:3])` has overlapping source and destination. Is the
   result defined? Answer by citing the spec, not by reporting what you observed.
3. State a rule for when `append` may mutate a slice the caller still holds.
   Then give the API design that means you never have to state that rule.
