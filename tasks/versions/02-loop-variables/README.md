# 02 — The bug that stopped existing

Three goroutines close over a loop variable. Once plainly, once with the `i := i`
shadow that a generation of Go programmers learned to write:

```go
for i := 0; i < 3; i++ {
	go func() { got = append(got, i) }()      // plain
}

for i := 0; i < 3; i++ {
	i := i                                     // shadowed
	go func() { got = append(got, i) }()
}
```

The goroutines wait on a start gate that closes after the loop, and the results
are sorted, so nothing here depends on the scheduler. Each version is compiled
and run by **one toolchain** — the one you have. Only the `go` line the module
declares changes.

| slot | module declares | which loop |
|---|---|---|
| `plain_at_go121` | `go 1.21` | plain |
| `plain_at_go122` | `go 1.22` | plain |
| `shadowed_at_go121` | `go 1.21` | shadowed |
| `shadowed_at_go122` | `go 1.22` | shadowed |

Predict each printed slice, exactly as `%v` renders it — `[0 1 2]`, say. Three
of the four are the same; the interesting content is which one is not.

```
undergo verify versions/02-loop-variables
```

The judge is the program's own output: the test writes a one-file module at each
go line, runs it with the toolchain in your PATH, and grades what it printed.

## Questions to answer in writing

1. Say what a `for` loop's variable *is*, in each of the two language versions —
   one sentence each, in terms of how many variables exist across three
   iterations.
2. `shadowed_at_go122` has the same value as `plain_at_go122`. Say what that
   means for the millions of `i := i` lines written before 2024, and whether
   deleting them is safe.
3. This change made existing programs behave differently, which Go 1 forbids.
   Say why it was allowed anyway, and name the one thing that made it
   acceptable.
4. Not every loop was affected. Give a loop whose behaviour is identical under
   both versions, and say what property of it makes the change invisible.
5. The change also applies to `for ... range`. Say whether a range loop's
   *value* variable was affected, and its *index*.
6. The program closes a start gate after the loop instead of letting the
   goroutines run as they are created. Say what the answer to `plain_at_go121`
   would be without that gate, and why that makes the pre-1.22 bug worse than a
   wrong answer.
