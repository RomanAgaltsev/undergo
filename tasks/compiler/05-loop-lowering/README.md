# 05 — Which loops become one call

Five functions. Two write a loop the compiler might recognise, two use the
builtin that does the same job, and one writes the loop over a slice of
**pointers**.

```go
func ZeroLoop(b []byte)          { for i := range b { b[i] = 0 } }
func ZeroClear(b []byte)         { clear(b) }
func CopyLoop(dst, src []byte)   { for i := range src { dst[i] = src[i] } }
func CopyBuiltin(dst, src []byte){ copy(dst, src) }
func ZeroPointers(s []*int)      { for i := range s { s[i] = nil } }
```

The test reads the compiler's own listing and asks, for each, whether the body
contains a call to a runtime routine.

| slot | question |
|---|---|
| `zero_loop_memclr` | does `ZeroLoop` call `memclrNoHeapPointers`? |
| `zero_clear_memclr` | does `ZeroClear`? |
| `copy_loop_memmove` | does `CopyLoop` call `memmove`? |
| `copy_builtin_memmove` | does `CopyBuiltin`? |
| `zero_pointers_memclr` | does `ZeroPointers` call `memclrNoHeapPointers`? |

Two of the five pairs look symmetric. They are not.

```
undergo verify compiler/05-loop-lowering
```

## Questions to answer in writing

1. One hand-written loop was replaced by a runtime call and the other was not.
   Say what the compiler must prove before it may make that substitution, and
   why the second loop fails that test where the first passes.
2. Given your answer to 1, say what `clear` buys you over the loop — and what
   `copy` buys you over *its* loop, which is a different answer.
3. `ZeroPointers` writes `nil`, which needs no write barrier at all. Say why it
   is still not lowered, and what the runtime routine's name is telling you.
4. Name the `-d=ssa/...` flag that would show you which pass is responsible, and
   say what you would look for in its output.
