# 06 — Where the collector looks for pointers

Five functions in `scan.go` each allocate a `Payload`, store the **only**
reference to it somewhere, force collections, and report whether the Payload
survived.

| slot | where the only reference is kept |
|---|---|
| `in_bytes_survives` | its bits, written into a `[]byte` |
| `through_view_survives` | a `*Payload` field of a struct laid over that same `[]byte` |
| `in_pointers_survives` | a `[]unsafe.Pointer` |
| `in_words_survives` | its address, in a `[]uintptr` |
| `in_anys_survives` | a `[]any` |
| `vet_flags_a_store` | does `go vet` object to anything in this package? |

Predict `true` or `false` for each.

The judge is the collector itself: each Payload is watched through a
`weak.Pointer`, which does not keep it alive and reports nil once it has been
reclaimed. The last slot is judged by running `go vet` over this package.

```
undergo verify gc/06-what-the-collector-scans
```

Two of the first five answers are `true`. The second slot is the one this task
exists for.

## Questions to answer in writing

1. The collector found the Payload in one arena and not in another that held the
   same bits. Say what the collector consults to decide whether a word is a
   pointer, and when that is decided for a heap allocation.
2. `ThroughView` stores through a type that *has* a pointer field. Say why that did
   not help, and name the one thing the view's type does change about the store.
3. Neither `go vet` nor the race detector's `checkptr` objects to any store here.
   Say what each one checks, and why neither can see this.
4. A generic arena hands out `*T` carved from one big `[]byte`. For which `T` is
   that safe? Say how you would enforce it — and whether a type constraint can.
5. Go 1.25's release notes say its new stack allocation "amplifies the effects of
   incorrect unsafe.Pointer usage". Connect that sentence to `alloc/06`.
