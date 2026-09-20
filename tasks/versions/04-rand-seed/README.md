# 04 — The change that says nothing

This program seeds the global `math/rand` source twice with the same value and
asks whether it got the same number back both times:

```go
rand.Seed(42)
first := rand.Int63()
rand.Seed(42)
fmt.Print(first == rand.Int63())
```

`rand.Seed` is deprecated. It still compiles. It still runs. There is no
warning, no diagnostic, and no error at any point.

It is run three times by one toolchain:

| slot | module declares | environment |
|---|---|---|
| `seed_honoured_at_go123` | `go 1.23` | — |
| `seed_honoured_at_go124` | `go 1.24` | — |
| `seed_honoured_at_go124_with_randseednop_0` | `go 1.24` | `randseednop=0` |

Predict `true` or `false` for each. Answering all three the same way scores one
of three.

```
undergo verify versions/04-rand-seed
```

The judge is the program's own output: the test writes a one-file module at each
go line, runs it with the toolchain in your PATH, and grades what it printed.

## Questions to answer in writing

1. `rand.Seed` is deprecated but still compiles and still runs. Say exactly what
   it does at go 1.24, and what it does not do.
2. Compare this with `versions/03-what-compiles`. Both are changes gated on the
   go line. Say how a developer finds out about each one, and which discovery
   mechanism you would rather have.
3. Give a realistic program that this change breaks, and say what its test suite
   would do — pass, fail, or flake.
4. Say why the standard library made the *global* functions ignore seeding
   rather than removing them, and what that buys.
5. `math/rand/v2` has no `Seed` at all. Say what a v2 package is allowed to do
   that a v1 package is not, and which promise makes that possible.
