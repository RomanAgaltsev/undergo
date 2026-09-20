# 08 — Two closures, one code pointer

Go will not let you write `f == g` on two func values. The usual workaround is
to compare what `reflect` calls their pointers:

```go
func mk(i int) func() int { return func() int { return i } }

a, b := mk(1), mk(2)          // same literal, different captured values
f := func() int { return 42 } // two separate literals,
g := func() int { return 42 } // identical bodies

reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
reflect.ValueOf(f).Pointer() == reflect.ValueOf(g).Pointer()
```

Both comparisons are run under two toolchains. The go line is `go 1.21` in every
run and does not vary — only the compiler does.

| slot | comparison | toolchain |
|---|---|---|
| `same_site_shares_on_local` | `a` vs `b` | the one you have |
| `same_site_shares_under_go1_26` | `a` vs `b` | **go1.26.0** |
| `distinct_literals_share_on_local` | `f` vs `g` | the one you have |
| `distinct_literals_share_under_go1_26` | `f` vs `g` | **go1.26.0** |

Predict `true` or `false` for each. Exactly one of the four differs between the
two toolchains, and two of them are a control.

```
undergo verify versions/08-closure-code-pointers
```

The judge is the program's own output: the test runs it under each toolchain and
grades what it printed.

**This task fetches a Go toolchain the first time it runs.** If it cannot be
fetched the task skips and says so.

## Questions to answer in writing

1. Go forbids `f == g` on func values. Say why, and say what
   `reflect.ValueOf(f).Pointer()` returns instead.
2. `a` and `b` came from the same literal and capture different values. Say what
   is actually different between them at run time, and what is the same.
3. One slot changed between two adjacent releases. Say what the compiler started
   doing, and why it is allowed to.
4. Somebody registers callbacks in a map keyed by `reflect.ValueOf(fn).Pointer()`
   to deduplicate them. Describe what breaks on the upgrade, and say whether
   their tests would have caught it.
5. Give a correct way to give a function value an identity you can compare.
   There is more than one; say what each costs.
