# 04 — Every way out of a range body

A range-over-function loop runs its body inside a call to `yield`. So when the
body leaves early, something has to travel back to the sequence and tell it to
stop.

`yield.go` ranges over the same sequence six times — it yields 1 to 5 and
records what `yield` returned on each call — and differs only in **how the body
leaves**.

| slot | the exit |
|---|---|
| `complete_calls` | the loop runs to the end |
| `break_calls` | a plain `break` on the third value |
| `goto_calls` | a `goto` to a label **outside** the loop, on the second |
| `labelled_break_calls` | `break outer` from inside a nested loop |
| `return_calls` | `return` from the enclosing function, immediately |
| `panic_calls` | `panic` on the second value, recovered outside |
| `break_last_return` | what `yield` returned on its last call, for `break` |
| `panic_last_return` | the same, for `panic` |

`Calls` counts yield calls that **completed**. Five of these six exits work the
same way. One does not, and the two `*_last_return` slots are there to find it.

```
undergo verify iter/04-goto-and-labels
```

## Questions to answer in writing

1. There is exactly one channel of information from the loop back to the
   sequence. Name it, and say why one bit is enough for `break`, `goto`,
   labelled `break` and `return` — four different destinations.
2. `goto` out of a range body is legal Go. Say what the compiler has to
   generate for it, and why it cannot simply emit a jump.
3. `panic_calls` is 1 even though the body ran twice, and `panic_last_return`
   is `true`. Explain both, and say what that means about whether a panic uses
   the yield protocol at all.
4. Following from that: a sequence that holds a lock, or an open file, while
   yielding cannot rely on seeing `yield` return false. Say what it must do
   instead, and write the two lines.
5. A sequence that keeps calling `yield` after it returned false gets a runtime
   panic rather than being ignored. Say why that is the right design, and who
   the runtime is protecting.
