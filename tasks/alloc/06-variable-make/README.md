# 06 — A make whose length is not a constant

`varmake.go` has two functions. Each makes a slice whose length is a parameter,
touches it, and returns its length; neither slice leaves its frame. One makes
bytes, the other makes `int64`s.

A slice that never leaves its frame sounds like the easy case for escape
analysis. For a *constant* length it is. For a length known only at run time,
the answer changed in Go 1.25, and it is not a yes or a no.

| slot | question |
|---|---|
| `m_says_escapes` | does the compiler's escape report say `make([]byte, n)` in `Bytes` escapes to the heap? |
| `bytes_32_allocs` | heap allocations per call of `Bytes(32)` |
| `bytes_33_allocs` | heap allocations per call of `Bytes(33)` |
| `words_4_allocs` | heap allocations per call of `Words(4)` |
| `words_5_allocs` | heap allocations per call of `Words(5)` |
| `bytes_1_allocs_hash_off` | heap allocations per call of `Bytes(1)`, compiled with the compiler's `variablemakehash` debug switch set to `n` |
| `m_says_escapes_hash_off` | the escape report's verdict on `make([]byte, n)` with that switch set |

The judges are the compiler and the allocator, not a rule of thumb. The test
compiles this package with escape diagnostics on and reads the compiler's own
verdict on the `make` in `Bytes`; then it recompiles the package's test in a
child process and counts heap allocations with `testing.AllocsPerRun`, once as
normal and once with the debug switch set. It never prints what it measured.

```
undergo verify alloc/06-variable-make
```

## Questions to answer in writing

1. The escape report's verdict and the allocation count disagree for one of these
   calls. Say what the escape report is a statement *about*, and what you may
   therefore never conclude from reading it.
2. `Bytes` and `Words` flip at different lengths. Say what unit the limit is
   measured in, and predict the largest `n` for which `make([]Point, n)` stays on
   the stack when `Point` is `struct{ X, Y, Z float64 }`.
3. The debug switch changes the allocation count and leaves the escape verdict
   alone. Say what that tells you about *where* in the compiler this optimisation
   lives, relative to escape analysis.
4. Before Go 1.25, how would you have kept a small scratch buffer of variable
   length off the heap by hand? Is that trick still worth writing today?
5. `hash` in the switch's name is not about maps. Find out what a hash-debug flag is
   for, and say why a compiler would ship an off-switch for an optimisation.
