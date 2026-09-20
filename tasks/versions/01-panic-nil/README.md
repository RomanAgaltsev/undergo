# 01 — What recover returns for panic(nil)

This program panics with `nil` and prints the dynamic type of whatever
`recover()` hands back:

```go
func main() {
	defer func() { fmt.Printf("%T", recover()) }()
	panic(nil)
}
```

It is compiled and run **four times, by one toolchain** — the one you already
have. Nothing about the source changes. What changes is the `go` line the module
declares, and the environment it runs in.

| slot | module declares | environment |
|---|---|---|
| `type_at_go120` | `go 1.20` | — |
| `type_at_go121` | `go 1.21` | — |
| `type_at_go121_with_panicnil_1` | `go 1.21` | `panicnil=1` |
| `type_at_go120_with_panicnil_0` | `go 1.20` | `panicnil=0` |

Predict the exact string `%T` prints in each. Two of the four are the same
string; which two is the question.

```
undergo verify versions/01-panic-nil
```

The judge is the program's own output: the test writes a one-file module at each
go line, runs it with the toolchain in your PATH, and grades what it printed.

## Questions to answer in writing

1. `panic(nil)` was legal and useless for twelve years. Say what a caller could
   not distinguish before 1.21, and why that made `recover()` unreliable.
2. Two of these four runs use the same toolchain and the same source and
   disagree. Name the file that decides, and the line in it.
3. The last two rows set `GODEBUG` explicitly. State the precedence rule between
   an explicit setting and the go line, in one sentence, and say which way it
   goes.
4. Go 1 promises that a program that works today keeps working. This change
   alters what a working program does. Explain how both are true at once.
5. Find `panicnil` in `$GOROOT/src/internal/godebugs/table.go` and say what its
   `Changed` and `Old` fields mean. Then find a setting in that table with no
   `Changed` value and say what its absence implies.
