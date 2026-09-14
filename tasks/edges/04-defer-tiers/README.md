# 04 — Open-coded or not

`defer` has more than one implementation. In the good case the compiler
**open-codes** it — the deferred call is emitted inline at each return, costing
almost nothing. In the bad case it falls back to the runtime, allocating a record
and pushing it onto a chain.

Six functions in `defers.go`. Predict which fall back.

| slot | the function |
|---|---|
| `one_uses_heap_defer` | one `defer` |
| `several_uses_heap_defer` | three, in a straight line |
| `in_loop_uses_heap_defer` | one, inside a `for` |
| `in_branch_uses_heap_defer` | one, inside an `if` |
| `many_uses_heap_defer` | nine, in a straight line |
| `with_recover_uses_heap_defer` | one, whose body calls `recover` |

```
undergo verify edges/04-defer-tiers
```

Two of the six fall back, for two different reasons. One of those reasons is a
number you can find from the other slots.

## Questions to answer in writing

1. Say what open-coding a `defer` means — where the deferred call ends up, and
   what the compiler must know to put it there.
2. `in_branch_uses_heap_defer` and `in_loop_uses_heap_defer` differ. Say what the
   compiler knows in one case that it cannot know in the other.
3. From the slots, deduce the limit on how many defers can be open-coded in one
   function. Say why a limit has to exist at all.
4. `with_recover_uses_heap_defer` — say whether `recover` changes the tier, and
   what it does change.
5. "`defer` in a hot loop" is a known smell. Say what is actually wrong with it,
   what you would write instead, and what you give up.
