# Go 1.27

Source: https://go.dev/doc/go1.27

## Size-specialized allocation routines make small allocations cheaper
→ candidate | alloc
"The compiler now generates calls to size-specialized memory allocation
routines, reducing the cost of some small (<80 byte) memory allocations by up to
30%", at a cost of "about 60 KB" of binary size. The opt-out
`GOEXPERIMENT=nosizespecializedmalloc` makes this measurable **within one
toolchain**: build the same benchmark both ways and predict the delta, and the
size threshold at which it disappears. No second toolchain needed, so this is an
ordinary optimize or predict task rather than a version-diff one.

## Closures may now share a code pointer
→ candidate | edges
"The compiler now generates simpler names for function literals" and "may also
combine multiple instances of the same function literal ... to share the same
code in the compiled binary". The consequence is the task: "function literals
with different captured closure data may have equal code pointers in more
cases". Predict whether two closures over different variables compare equal via
`reflect.Value.Pointer`, and explain why comparing function code pointers was
always wrong — the release notes link to the documented warning.

## Generic methods
→ candidate | generics
Methods may now declare their own type parameters. The task is not "use it": it
is to predict how many instantiations GC-shape stenciling emits for a generic
method, verified against the symbol table, and how that differs from a generic
function with the same shapes.

## The goroutineleak profile is generally available
→ candidate | sched
`runtime/pprof` gains a GA `goroutineleak` profile, and the
`goroutineleakprofile` GOEXPERIMENT is deleted. A leaked goroutine is one
"blocked on some concurrency primitive ... that cannot possibly become
unblocked". A build task: plant a leak of each shape — channel, `sync.Mutex`,
`sync.Cond` — and make the profile name all of them.

## Timer channels are always unbuffered
→ candidate | sched
The `asynctimerchan` GODEBUG "has been removed permanently", so channels from
package `time` are now always synchronous regardless of GODEBUG. A predict task
on the interleaving of a timer fire and a receive is now stable enough to grade.
The 1.23–1.26 behaviour becomes a version-diff task once E2 is swept.

## The simd package is experimental behind GOEXPERIMENT
→ candidate | asm
A portable SIMD surface enabled with `GOEXPERIMENT=simd`. Beating a pure-Go
baseline with it, and explaining where the win comes from, is an optimize task —
but it is gated on the experiment, so the task must say so in `requires`.

## encoding/json is backed by v2
→ no action
Stricter defaults and faster unmarshal, reversible with `GOEXPERIMENT=nojsonv2`.
The observable surface is library behaviour rather than machine behaviour.

## Tracebacks include pprof goroutine labels
→ no action
Useful in production, but not something a solver could be instructively wrong
about.

## Removed GODEBUG settings are now accepted at their final default
→ no action
A go-command compatibility rule, not runtime behaviour.

## crypto/mldsa, ML-DSA in crypto/tls, MLKEM1024
→ no action
Post-quantum cryptography is outside what this repo is about.

## macOS 13 minimum, ppc64 ELFv2, bzr support removed
→ no action
Platform support, not observable Go behaviour.
