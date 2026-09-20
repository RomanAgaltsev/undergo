# 06 — The channel that was buffered

The first five tasks in this track change one line of `go.mod` and get two
answers out of one toolchain. `05-go-line-limits` found the edge of that: once a
GODEBUG has been **removed**, no go line reaches the behaviour it used to
restore, and its fifth written question asks what it would take to answer the
original question anyway.

This is that answer. It changes the **toolchain** instead.

The program is one line:

```go
func main() { fmt.Print(cap(time.After(time.Hour))) }
```

It is compiled and run three times:

| slot | toolchain | module declares |
|---|---|---|
| `cap_on_local` | the one you have | `go 1.27` |
| `cap_under_go1_22` | **go1.22.12** | `go 1.19` |
| `cap_on_local_at_go_line_119` | the one you have | `go 1.19` |

Predict each printed capacity, as a string — `"0"` or `"1"`.

The last two rows declare the same go line and differ only in which compiler and
runtime built them. That pair is the task.

```
undergo verify versions/06-timer-channel-buffer
```

The judge is the program's own output: the test writes a one-file module and
runs it under each toolchain, grading what it printed.

**This task fetches a Go toolchain the first time it runs**, which needs the
network once; afterwards it is in your module cache. If it cannot be fetched the
task skips and says so, and `undergo doctor` will tell you which toolchains
resolve here.

## Questions to answer in writing

1. `versions/05-go-line-limits` proved no go line reaches the buffered channel.
   Say what changed here that the go line could not change, in one sentence.
2. Two of these three runs use the same go line and disagree. Say which axis
   each of the two — the go line and the toolchain — actually selects.
3. `asynctimerchan` existed for four releases and was then removed. State what a
   GODEBUG is for, and what its removal costs somebody who had not finished
   migrating.
4. A buffered timer channel let the runtime deliver without a receiver. Name the
   bug that the buffer hid, and the one that removing it introduced for code
   written against the old behaviour.
5. This task downloads a toolchain. Say what that costs the repository's
   offline-by-default promise, and what the harness does about it.
