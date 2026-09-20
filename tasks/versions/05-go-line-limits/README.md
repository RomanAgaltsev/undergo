# 05 — What the go line cannot reach

The first four tasks in this track all worked the same way: change one line of
`go.mod`, get a different answer out of the same toolchain. This one tries the
same trick and does not get it.

**Program one** prints the capacity of the channel `time.After` returns. Before
Go 1.23 that channel was buffered; from 1.23 it is not, and `asynctimerchan=1`
was the GODEBUG that restored the old behaviour.

```go
func main() { fmt.Print(cap(time.After(time.Hour))) }
```

**Program two** announces itself twice, so that a failure *before* the program
runs can be told apart from a failure inside it.

```go
func init() { fmt.Print("init ") }

func main() { fmt.Print("main") }
```

Program two is run once, at `go 1.27`, with `GODEBUG=asynctimerchan=1` in the
environment — reaching for the switch that used to restore the old timer.

| slot | question |
|---|---|
| `after_cap_at_go119` | the capacity printed by program one at `go 1.19`, as a string |
| `after_cap_at_go127` | the same at `go 1.27` |
| `init_ran_with_removed_godebug` | did program two's output contain `init`? |
| `output_says_fatal_error` | did program two's output contain `fatal error`? |

```
undergo verify versions/05-go-line-limits
```

The judge is the programs' own output: the test writes a one-file module at each
go line, runs it with the toolchain in your PATH, and grades what it printed.

## Questions to answer in writing

1. `time.After`'s channel was buffered before Go 1.23 and is not now. Say what
   the go line does to that fact, and why.
2. `init_ran_with_removed_godebug` is `false`. Say how much of the program ran,
   and where in the process's life the failure happened.
3. That failure is a `fatal error`, not a panic. Say what that rules out, and
   name the other task in this catalogue that grades the difference.
4. Go could have ignored a removed setting silently. Say why it does not, and
   what a silent ignore would have cost somebody migrating a large codebase.
5. Given this task's result, state what you would now need in order to answer
   "what was `cap(time.After(d))` in Go 1.22?" — and say why that is a different
   kind of exercise from the rest of this track.
