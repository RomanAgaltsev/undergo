# 04 — What the profile lets the compiler assume

`Shape` has three implementations. Two functions call `Area` through the
interface:

```go
func Hot(s Shape) float64  { return s.Area() }
func Cold(s Shape) float64 { return s.Area() + 1 }
```

A CPU profile is checked in as `default.pgo`. It was collected from a workload
that calls `Hot` with a `Rect` overwhelmingly, and barely touches `Cold`.

The test builds the package **twice** — once with the profile, once with
`-pgo=off` — and counts what the compiler says it did.

| slot | question |
|---|---|
| `devirtualized_with_pgo` | how many interface call sites were devirtualized, with the profile |
| `devirtualized_without_pgo` | the same, without it |
| `target_is_rect` | does a note name `Rect.Area`? |
| `target_is_circle` | does a note name `Circle.Area`? |

There are **two** interface call sites. Predict the first number carefully.

```
undergo verify compiler/04-pgo-devirtualization
```

## Questions to answer in writing

1. Devirtualizing turns an interface call into a direct one — but the profile
   is only a guess about the future. Say what the compiler must still emit for
   the cases the guess gets wrong, and why the transformation is a win anyway.
2. Only one of the two call sites was devirtualized. Say why, and what that
   tells you about how PGO decides.
3. Devirtualizing is mostly valuable because of what it *enables*. Say what
   that is, and connect it to the inlining budget you measured in
   `compiler/01`.
4. The `Area` implementations are deliberately written too large to inline.
   Remove that — make them one-liners — regenerate the profile, and see what
   happens to the counts. Explain the result. (Read the package doc comment
   first; it tells you what to expect and you should be able to say why.)
5. A profile describes the past. Say what happens when production traffic
   changes shape so that `Circle` dominates, and whether the compiled program
   becomes **wrong** or merely **slower**.
