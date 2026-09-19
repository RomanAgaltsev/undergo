# 02 — Which bounds checks survive

Every slice index in Go is checked at run time — unless the compiler can prove it
does not need to be. Six functions in `bce.go`. Predict which keep at least one
check.

| slot | the function |
|---|---|
| `ranged_checked` | `for _, x := range xs` |
| `counted_checked` | `for i := 0; i < len(xs); i++ { xs[i] }` |
| `two_slices_checked` | `for i := range xs { ys[i] }` |
| `unchecked_checked` | `xs[i]`, with `i` a parameter |
| `after_one_check_checked` | four indexes after one `len(xs) < 4` test |
| `next_element_checked` | `xs[i+1]` in a loop to `len(xs)-1` |

```
undergo verify compiler/02-bounds-checks
```

Two are checked and four are not. One of the four will probably surprise you; it
is the one where the compiler has to do arithmetic to reach its conclusion.

The judge is the compiler's own report: the test compiles this package with
`-d=ssa/check_bce` and reads which indexes it still names.

## Questions to answer in writing

1. State what the compiler must prove to drop a check, and where it gets the
   facts from.
2. `ranged_checked` and `two_slices_checked` are the same loop shape, and differ.
   Say exactly which fact the compiler has in one case and not the other — and
   say what one line would restore it.
3. `next_element_checked` requires reasoning about `i+1`. Say what chain of facts
   gets the compiler there.
4. The old idiom `_ = xs[n-1]` before a run of indexes: say what it was for, and
   whether `after_one_check_checked` suggests it is still needed.
5. A dropped check is a proof, not a guess. Say what would happen if the proof
   were wrong, and why that is not a risk you are taking.
