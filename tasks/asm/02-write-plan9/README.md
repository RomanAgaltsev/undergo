# 02 — Write one in Plan9

Implement `SumInt64` in `sum_amd64.s`. It adds a slice of `int64` and returns
the total.

```go
sum.SumInt64([]int64{1, 2, 3, 4, 5})   // -> 15
```

The Go side is already written for you, and it is worth reading before the
assembly:

- `sum.go` carries `//go:build amd64`, declares `SumInt64` **with no body** —
  which is how Go says "the definition is in assembly" — and marks it
  `//go:noescape`.
- `sum_fallback.go` carries `//go:build !amd64` and implements the same function
  in ordinary Go.

**Do not touch the fallback.** It is there because this repository builds every
task on every platform, including an arm64 macOS runner, where a package
containing only amd64 assembly fails with `missing function body`. Written
question 3 is about what that means for you.

## The contract

- Correct for the empty slice, one element, several, negatives, and odd lengths.
- Correct for a thousand elements.
- Allocates nothing — which is what `//go:noescape` is promising.

```
undergo verify asm/02-write-plan9
```

A slice argument arrives as **three** words. Work out what they are and where
the return value goes before writing an instruction; the `$0-N` in the `TEXT`
directive depends on it.

## A warning about your safety net

`go vet` reads Plan9 assembly and checks it against the Go declaration. It is
worth knowing exactly how far that goes, because it is not as far as you would
hope.

Get a **named offset** wrong and it tells you precisely:

```
sum_amd64.s:23:1: [amd64] SumInt64: invalid offset ret+16(FP); expected ret+24(FP)
```

Get the **declared argument size** in the `TEXT` directive wrong — `$0-24` where
it should be `$0-32` — and vet says nothing, the package builds, and the tests
pass. Try it.

## Questions to answer in writing

1. Derive the `$0-32` in the `TEXT` directive. Say what the `0` is, what the `32`
   is, and which words of the slice header each offset in the body refers to.
2. `//go:noescape` is a promise, not a request. Say exactly what you are
   promising the compiler, what it does with that promise, and what would happen
   if your assembly broke it.
3. On arm64 the fallback runs and every test passes — proving nothing about your
   assembly. Say how you would make that visible in CI rather than leaving a
   green tick that means the wrong thing.
4. `NOSPLIT` — say what it disables, what it buys, and when it is unsafe.
5. Vet catches a wrong offset and not a wrong argument size. Say why the second
   is harder to check, and what the consequence of getting it wrong actually is.
