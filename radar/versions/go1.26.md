# Go 1.26

Source: https://go.dev/doc/go1.26

## The Green Tea garbage collector is on by default
→ candidate | gc
"Between a 10–40% reduction in garbage collection overhead in real-world
programs that heavily use the garbage collector", from "better locality and CPU
scalability" when marking and scanning small objects. `GOEXPERIMENT=nogreenteagc`
turns it off, so the comparison is an A/B within one toolchain — and that opt-out
"is expected to be removed in Go 1.27", which dates any task built on it.

## Green Tea leans on vector instructions on newer amd64
→ candidate | gc
A further "10% in garbage collection overhead" on "Intel Ice Lake or AMD Zen 4
and newer", because the collector "leverages vector instructions for scanning
small objects when possible". A predict task here must declare `requires.arch`
and say what it depends on — the same answer is wrong on older amd64 and on
arm64, which is exactly the kind of dependency `requires` exists to record.

## The compiler stack-allocates slice backing stores in more cases
→ candidate | alloc
"The compiler can now allocate the backing store for slices on the stack in more
situations." This is the sharpest allocation task in the release because it comes
with its own A/B: `-gcflags=all=-d=variablemakehash=n` disables the new stack
allocations, and `bisect -compile=variablemake` finds the specific one. Predict
which of several `make` calls reach the heap, then prove it.

This is not hypothetical: while building `alloc/01-zero-alloc-join` a benchmark
fixture that wrote `make([]byte, 64)` to a blank identifier measured **zero**
allocations, because the slice never escaped. A task on this mechanism would
have made that obvious rather than surprising.

## cgo call overhead is down about 30%
→ candidate | edges
"The baseline runtime overhead of cgo calls has been reduced by ~30%." A measured
cgo-boundary task needs a C toolchain, so it must declare that and skip cleanly
where one is absent — worth writing only once the harness can express that
requirement.

## The heap base address is randomized on 64-bit
→ candidate | edges
"The runtime now randomizes the heap base address at startup", disabled with
`GOEXPERIMENT=norandomizedheapbase64`. The task is the consequence, not the
feature: which observations about pointer values stop being reproducible across
runs, and which were never reproducible in the first place.

## The goroutineleak profile arrives as an experiment
→ no action
Behind `GOEXPERIMENT=goroutineleakprofile` here and generally available in 1.27,
where it is already a candidate. Its 1.26 limitation is worth remembering when
that task is written: leaks reachable "through global variables or the local
variables of runnable goroutines" are not detected.

## runtime/secret
→ no action
Experimental, Linux amd64/arm64 only, and about forward secrecy rather than
machine behaviour a solver could predict.

## simd/archsimd
→ no action
Superseded as a candidate by the portable `simd` package in 1.27, which is
already tracked there. The architecture-specific API is "not yet considered
stable".

## WebAssembly heaps grow in smaller increments
→ no action
Real, but there is no wasm track and no plan for one.
