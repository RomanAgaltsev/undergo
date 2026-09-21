# 13 — The heap that moved

For years, an early heap pointer in a Go program printed as something like
`0xc00000a0c0`. The prefix was recognisable enough that people treated it as
meaningful on sight.

Go 1.26 randomizes the heap base on 64-bit platforms.

## The program

```go
type T struct{ a, b int }

func main() {
	x := new(T)
	fmt.Printf("%p", x)
}
```

It runs under two toolchains. The go line is `go 1.21` in every run and does not
vary — only the compiler does.

| slot | question |
|---|---|
| `base_prefix_is_c000_under_go1_25` | does the printed address start with `0xc000`? |
| `base_prefix_is_c000_on_local` | does it on the toolchain you have? |

Two slots. Predict `true` or `false` for each.

```
undergo verify versions/13-heap-base
```

**This task fetches a Go toolchain the first time it runs.** If it cannot be
fetched the task skips and says so.

## The exercise that is not graded

Do this before you read the explanation, because it is the actual point of the
task and no slot will check it for you.

Run the program **ten times** under `go1.25.7`, then ten times under your own
toolchain:

```
GOTOOLCHAIN=go1.25.7 go run .
```

Split each address into two parts: the high end, and the last three hex digits.
Ask, separately for each part and each toolchain, whether it is the same from
run to run. You should get four answers, and they are not all what the version
note would lead you to expect.

**This is deliberately not a graded slot.** Work out why before you look it up:
what would a test have to assume in order to grade "these two runs differ", and
how often would that assumption hold?

## Questions to answer in writing

1. Which part of the address did Go 1.26 change? Which part was never stable,
   and what decides it?
2. You find a bug report that quotes a pointer value and argues from it. What is
   wrong with the argument, and was it wrong before 1.26 too?
3. Why would a language runtime randomize this on purpose? Name what it costs
   and who it protects against.
4. Answer the question in the section above: why is the run-to-run observation
   taught here in prose instead of graded as a slot? What would the test have to
   do to grade it honestly, and what would that cost?
