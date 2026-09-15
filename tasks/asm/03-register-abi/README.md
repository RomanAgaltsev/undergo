# 03 — Where the arguments actually go

Since Go 1.17, arguments travel in **registers** rather than on the stack. The
sequence is fixed, it is not very long, and integers and floats draw from two
independent pools.

`abi.go` declares four functions with deliberately awkward signatures, and a
helper that hands you the compiler's own listing. Nothing here needs running —
every answer is in that listing.

| slot | question |
|---|---|
| `nine_reads_stack` | does `Nine` (nine `int64`s) load any argument from the stack? |
| `ten_reads_stack` | does `Ten` (ten `int64`s) load any argument from the stack? |
| `pair_reads_stack` | does `SumPair` (one two-field struct) load from the stack? |
| `mixed_uses_floats` | does `Mixed` do its floating-point arithmetic in X registers? |

The first two differ by one argument. That is the whole point: they bracket a
number you should be able to name afterwards.

```
undergo verify asm/03-register-abi
```

You can read the listing yourself — but do it **after** committing your answers,
because the listing *is* the answer:

```
go build -gcflags=-S .
```

## Questions to answer in writing

1. `Nine` and `Ten` differ by one argument and one of them spills to the stack.
   State the number of integer argument registers, then name them in order.
2. `SumPair` takes **one** argument and does not touch the stack. Say what the
   ABI did with the struct, and describe a struct that would spill instead.
3. Integers and floats draw from independent sequences. Say what that means for
   a function taking eight `int64`s and eight `float64`s, and why the ABI was
   designed that way rather than with one shared sequence.
4. This ABI replaced one that passed everything on the stack. Say what the
   change bought, and why it needed a `//go:registerparams` transition and a
   whole release cycle rather than landing at once — what does the runtime need
   to know at every instruction for this to be safe?
5. `bodyOf` matches on `"." + name + " STEXT"` rather than on the package name.
   Find out what the listing actually calls a function, and say what the
   obvious matcher would have reported for all four slots.
