# 05 — Which of these comparisons panic

`a == b` on two `any` values compiles for **every** pair of types. Some of them
crash at run time.

`Panics(a, b)` reports whether the comparison panicked.

| slot | comparison |
|---|---|
| `slices_panic` | two `[]int` |
| `maps_panic` | two `map[int]int` |
| `funcs_panic` | two `func()` |
| `struct_slice_panics` | two `struct{ S []int }` |
| `array_of_slice_panics` | two `[1][]int` |
| `plain_struct_panics` | two `struct{ A, B int }` |
| `array_panics` | two `[2]int` |
| `different_types_panic` | a `[]int` against a `string` |
| `plain_struct_equal` | do the two plain structs compare **equal**? |

The eighth slot is the one to think hardest about: a slice is on one side of it.

```
undergo verify iface/05-interface-comparison
```

## Questions to answer in writing

1. Two interfaces holding slices panic rather than failing to compile. Say why
   the compiler cannot catch it, and what it would cost to make it able to.
2. A struct containing a slice panics, and so does an array of slices. Say how
   deep the rule goes, and state it in one sentence.
3. `different_types_panic` is the surprise. Explain it in terms of the **order**
   of the two steps a comparison performs, and say what that means for a
   heterogeneous `[]any` you are scanning for duplicates.
4. `iface/02` in this repository covers the typed-nil trap. Say whether a typed
   nil compares equal to a plain `nil`, and connect the answer to the two-word
   model of an interface value from `iface/01`.
5. You are writing a cache keyed on `any`. Say what can go wrong, and name two
   different ways to make it safe.
