# 01 — Reading the compiler's assembly

Four small functions in `subject.go`. Predict what the compiler made of them,
then read the listing and find out.

| slot | question |
|---|---|
| `add_frame_bytes` | the local frame size of `AddInts`, in bytes |
| `add_is_nosplit` | is `AddInts` marked `nosplit`? |
| `add_uses_ax_and_bx` | does its body read the two arguments from `AX` and `BX`? |
| `lenof_moves_bx_to_ax` | does `LenOf` compile to a single register move? |
| `firstof_can_panic` | does `FirstOf` retain a bounds-check panic? |
| `divmod_frame_is_zero` | does `DivMod` get by with no frame? |

```
undergo verify asm/01-read-plan9
```

Two of these are about the calling convention and four are about what a function
still needs the stack for. The judge is the compiler's own `-S` listing: the test
compiles this package and reads what it emitted.

## Questions to answer in writing

1. Name the calling convention Go uses on amd64, when it arrived, and what it
   replaced. Say where the first three integer arguments go and where the result
   comes back.
2. A `string` is two words and a slice is three. Say which registers each occupies
   and what `LenOf` therefore has to do.
3. `FirstOf` indexes a slice and keeps no bounds check. Say what let the compiler
   drop it — and note that it is the same mechanism a whole task in the
   `compiler` track is about.
4. Two of these functions have a zero frame and one does not. Say what the one
   with a frame needs it for.
5. `nosplit` — say what it means, why a small leaf function gets it, and what
   would go wrong if a large function were marked with it.
