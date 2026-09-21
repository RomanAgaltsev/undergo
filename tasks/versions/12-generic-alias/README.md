# 12 — The alias that adds nothing

`type B = A` declares no new type. It declares another name for one, and
nothing at run time can tell `B` from `A`.

This task asks whether that still holds when the alias takes a type parameter.

## The program

```go
type Box[T any] struct{ V T }

type Alias[T any] = Box[T]

func main() {
	var a Alias[int]
	var b Box[int]
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	fmt.Printf("%v|%s", ta == tb, ta.String())
}
```

It is compiled under two toolchains, and run under one. **The go line is
`go 1.23` in every run and does not vary** — only the compiler does. That is not
incidental here: read the slot table and ask why the go line could not have been
the axis for this one.

| slot | answer |
|---|---|
| `compiles_under_go1_23` | does the declaration compile on the old toolchain? |
| `compiles_on_local` | does it compile on the one you have? |
| `diagnostic_names_an_experiment` | does the old compiler's message mention `GOEXPERIMENT`? |
| `instantiations_are_identical` | is `reflect.TypeOf(a) == reflect.TypeOf(b)`? |
| `alias_type_name` | the exact string `ta.String()` prints |

The last one is a string, not a boolean. Write it exactly as the program would.

```
undergo verify versions/12-generic-alias
```

**This task fetches a Go toolchain the first time it runs.** If it cannot be
fetched the task skips and says so.

## Questions to answer in writing

1. Answer for the non-generic case first: after `type B = A`, how many types
   exist? How many type descriptors? Then say what you expect for `Alias[int]`
   and `Box[int]`, and why.
2. `alias_type_name` has two plausible answers and only one is right. What would
   the other one imply about what an alias is?
3. The go line is fixed at `1.23` in every run of this task, and the behaviour
   still differs. What does that tell you about where this feature was gated?
   Contrast it with task 09, where the go line is the *only* thing that varies.
4. If you upgraded a dependency that started using a generic alias in its public
   API, what would break for a caller on an older toolchain — and would the
   error mention the dependency or your own code?
