# 04 — The pragma that promises

`noescape_amd64.s` holds two functions with **byte-identical bodies**. Both sum
four `int64`s behind a pointer. Compare them if you like — they are the same
instructions.

The only difference is in `noescape.go`:

```go
//go:noescape
func SumFree(p *[4]int64) int64

func SumKept(p *[4]int64) int64
```

`CallFree` and `CallKept` each build a four-element array and pass its address
to one of them.

| slot | question |
|---|---|
| `call_free_allocs` | allocations per call for the one with the pragma |
| `call_kept_allocs` | allocations per call for the one without |
| `same_result` | do the two return the same number? |

```
undergo verify asm/04-noescape
```

## Questions to answer in writing

1. The two functions have identical bodies and different allocation behaviour.
   Say what the pragma tells the compiler, and why the compiler cannot work it
   out for itself.
2. `//go:noescape` is a **promise, not a check**. Write (on paper) an assembly
   function for which the promise is false, and say precisely what goes wrong —
   name the assumption the garbage collector makes that you have broken, and
   say why the resulting bug would appear far from its cause.
3. Both functions are declared `NOSPLIT`. Say what that promises and what
   happens if *that* promise is false.
4. Does `//go:noescape` do anything for a function whose arguments contain no
   pointers? Say why the answer is what it is.
5. The `TEXT` directives say `$0-16`. Show the arithmetic. Then say which part
   of that line `go vet` actually checks and which part it does not — and what
   would happen if you wrote `$0-8`.
