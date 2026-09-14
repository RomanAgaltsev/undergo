# 02 — Where nothing lives

A `struct{}` occupies zero bytes. It still has to have an address, because Go
lets you take one. Six questions about what that address is.

| slot | the question |
|---|---|
| `two_empty_structs_share_an_address` | do two separate `new(struct{})` compare equal? |
| `two_zero_arrays_share_an_address` | the same for `new([0]int)` |
| `trailing_field_adds_padding` | does a zero-size field at the **end** of a struct make it bigger? |
| `leading_field_adds_padding` | does one at the **start**? |
| `zero_length_slices_share_data` | do `a[0:0]` and `a[4:4]` have the same data pointer? |
| `empty_slice_equals_nil` | does `[]int{}` compare equal to `nil`? |

```
undergo verify edges/02-zero-size
```

Two of these differ only in where a field sits. That pair is the one with a
practical consequence.

## A warning about the first two

The specification does **not** guarantee the answer. It says two distinct
zero-size variables *may* have the same address — may, not must. This task grades
what this toolchain does on the default build, and written question 4 is about
what you should therefore never assert in your own tests.

## Questions to answer in writing

1. Quote what the spec actually says about the address of a zero-size variable.
   Then say what `runtime.zerobase` is and how it produces the measured answer.
2. A zero-size field adds padding at the end of a struct and not at the start.
   Say why — what would go wrong without the padding?
3. From your answer to 2, say what `structs.HostLayout` has to do with any of
   this, and why it exists.
4. You measured a specific answer to a question the spec leaves open. Say what you
   would and would not write a test against, and what you would do if you needed
   the behaviour to be stable.
5. `empty_slice_equals_nil` — say what distinguishes the two values, and name a
   place where the distinction has bitten someone (JSON is the usual one).
