# 04 — What an assertion costs

A type assertion is not one thing. Some compile to a couple of instructions and
some call into the runtime, and the difference is not where you would guess.

Five assertions in `assert.go`. Predict which need the runtime's help.

| slot | the assertion |
|---|---|
| `concrete_comma_ok_calls_runtime` | `v.(Dog)`, two-value form |
| `concrete_panicking_calls_runtime` | `v.(Dog)`, one-value form |
| `empty_to_interface_calls_runtime` | `any` asserted to a one-method interface |
| `interface_to_interface_calls_runtime` | one non-empty interface to another |
| `type_switch_calls_runtime` | a three-case type switch over concrete types |

```
undergo verify iface/04-assertion-cost
```

The first two differ only in whether they take the `ok`. That is the pair worth
thinking hardest about — and the answer is not that the safe one is slower.

## Questions to answer in writing

1. Asserting to a **concrete** type and asserting to an **interface** are
   different jobs. Say what each one has to check, and why only one of them can
   be done with a comparison.
2. The comma-ok and panicking forms give different answers. Say what the extra
   call in one of them is for, and whether it is on the hot path.
3. A type switch over three concrete types does not call the runtime. Say what it
   compiles to instead, and what would change if one case were an interface.
4. None of these assertions **allocates**, including the interface ones. Say why
   — in terms of what an interface value's two words hold.
5. `itab` lookups are cached. Say what is cached, what the key is, and what that
   means for the first call against the millionth.
