# 03 — Which of these reflect calls panic?

Eight operations in `settable.go`, each doing one thing with `reflect`. Predict
`true` if it panics and `false` if it does not.

| slot | what it does |
|---|---|
| `set_through_value` | `reflect.ValueOf(c)` then set a field |
| `set_through_pointer` | `reflect.ValueOf(&c).Elem()` then set a field |
| `set_unexported` | set an unexported field, through a pointer |
| `set_wrong_type` | `SetString` on an `int` field |
| `set_slice_element` | `reflect.ValueOf(s).Index(0).SetInt(...)` on a `[]int` |
| `set_map_value_in_place` | `MapIndex(...)` then set a field on the result |
| `elem_of_nil_pointer` | `Elem()` on a nil `*Config`, then set |
| `append_and_assign` | `reflect.Append` and assign the result back |

Five of the eight panic. Two of the pairs look alike and differ for reasons worth
being able to state.

```
undergo verify reflect/03-settability
```

## Questions to answer in writing

1. State the settability rule in one sentence, covering both of its conditions.
2. `set_through_value` and `set_slice_element` both start from a value, not a
   pointer. One panics and the other does not. Explain the difference in terms of
   what was copied.
3. Among the five that panic, there are **four** distinct reasons. Group them and
   name each reason. Which one is not about settability at all?
4. `CanSet` exists. Rewrite `set_through_value` so it reports rather than
   panicking, and say why a library that takes `any` from its caller should
   prefer that — and what it should do about the case in question 3 that `CanSet`
   would not have caught.
