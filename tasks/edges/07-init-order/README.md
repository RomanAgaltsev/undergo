# 07 — Before main, in what order

Everything a Go program does before `main` is specified. Not conventional, not
implementation-defined: specified, in the language spec, as an algorithm.

`testdata/prog` is three packages that announce themselves. `main` imports
`beta`, `beta` imports `alpha`, and each package prints when its variables are
initialised and when its `init` runs.

```go
// package alpha
var A = announce("alpha.A")
func init() { fmt.Println("alpha.init") }

// package beta — imports alpha
var B = announce("beta.B (sees " + alpha.A + ")")
func init() { fmt.Println("beta.init") }

// package main — imports beta
var Second = announce("main.Second") + First
var First = announce("main.First")
func init() { fmt.Println("main.init") }
```

Read all three files, then predict:

| slot | question |
|---|---|
| `first_line` | the first line the program prints |
| `main_first_var` | which of `main`'s two variable announcements comes first |
| `init_runs_after_its_package_vars` | does a package's `init` run after that package's variables? |
| `full_order` | every line before `main.main`, comma-separated, in order |

Note that `Second` is declared **above** `First`. That is not a typo.

```
undergo verify edges/07-init-order
```

The judge is the program itself: the test runs `testdata/prog` and reads the
lines it printed, in the order it printed them.

## Questions to answer in writing

1. State the rule that decides the order of package-level variables *within* one
   package. It is not the order they are written in; say what it is instead.
2. `main.Second` is declared first and initialised second. Say what would happen
   if you removed its dependency on `First` — and what that tells you about how
   fragile a program that cares about this order would be.
3. State the rule that decides the order of *packages*. Then say what happens
   with two packages that do not import each other, and whether you may rely on
   the answer.
4. A package has two files, each with an `init`. Say what decides which runs
   first, and why relying on it is a bad idea even though it is deterministic.
5. Initialisation can fail — a `var` whose initialiser panics. Say where that
   panic surfaces, what the exit status is, and why no `recover` you write can
   help.
