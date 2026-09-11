# Go 1.24

Source: https://go.dev/doc/go1.24

## The builtin map is a Swiss table
→ candidate | types
The new implementation is "based on Swiss Tables", disabled with
`GOEXPERIMENT=noswissmap` — so the before-and-after is measurable in one
toolchain. The task the `types` track already plans: predict a map's real memory
for a given key count and load factor under both, and explain where the
difference comes from. Beware the shape of the question: "is it smaller" is a
coin flip, "how many bytes and why" is a task.

## The weak package
→ candidate | weak
"A low-level primitive provided to enable the creation of memory-efficient
structures, such as weak maps ... canonicalization maps ... and various kinds of
caches." This is the floor for the entire `weak` track: every task in it needs
`requires.go` of at least 1.24.

## runtime.AddCleanup replaces SetFinalizer
→ candidate | weak
The notes enumerate the differences precisely, which is what makes it gradeable:
"multiple cleanups may be attached to a single object, cleanups may be attached
to interior pointers, cleanups do not generally cause leaks when objects form a
cycle, and cleanups do not delay the freeing of an object or objects it points
to." Four distinct claims, four predictions — give a solver a cyclic structure
with both mechanisms attached and ask which runs, and when.

## sync.Map is a hash-trie
→ candidate | types
"Modifications of disjoint sets of keys are much less likely to contend on
larger maps, and there is no longer any ramp-up time required to achieve
low-contention loads." `GOEXPERIMENT=nosynchashtriemap` restores the old one, so
the contention curve can be measured both ways. The "ramp-up time" claim is the
interesting half: predict the shape of throughput over time under the old
implementation and explain what was warming up.

## A new runtime-internal mutex
→ candidate | sched
Disabled with `GOEXPERIMENT=nospinbitmutex`. Only observable through contention
behaviour, so the task must be a measured one — and it risks being a benchmark
without a mechanism, which is the failure mode to avoid. Write it only if the
spin/park boundary can be made visible.

## Runtime CPU overhead down 2–3% overall
→ no action
"Several performance improvements to the runtime have decreased CPU overheads by
2–3% on average across a suite of representative benchmarks." An aggregate of the
three entries above; the components are the tasks, not the total.

## testing/synctest arrives as an experiment
→ no action
Behind `GOEXPERIMENT=synctest` here and generally available in 1.25, where it is
a candidate. Record the API change: 1.24 has `synctest.Run`, 1.25 has
`synctest.Test`, so a task written against the GA form does not compile here.

## crypto/mlkem, crypto/hkdf, crypto/pbkdf2, crypto/sha3
→ no action
Cryptographic primitives, outside what this repo asks about.

## The linker emits a GNU build ID and a Mach-O UUID
→ no action
Build metadata, not observable Go behaviour.
