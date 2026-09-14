# 01 — The door that was closed

`//go:linkname` lets a package refer to a symbol by its linker name, including
unexported symbols in other packages. It is the mechanism behind a generation of
libraries reaching into the runtime for things the standard library did not
export.

Go 1.23 closed most of it. Four programs in `testdata/`, each linknaming a
different target. Predict which still link.

| slot | the target |
|---|---|
| `gcstart_links` | `runtime.gcStart` |
| `nanotime_links` | `runtime.nanotime` |
| `mallocgc_links` | `runtime.mallocgc` |
| `own_package_links` | a function in the program's own package |
| `gcstart_links_with_optout` | `runtime.gcStart`, built with `-ldflags=-checklinkname=0` |

```
undergo verify edges/01-linkname
```

Three of the four runtime targets are internal, unexported functions, and they do
not all behave the same way. Working out what separates them is the task.

## Questions to answer in writing

1. Say what `//go:linkname` does, and why it requires importing `unsafe` even
   when no unsafe operation appears in the file.
2. Go 1.23 refused a *subset* of linknames. From the measurements, say what the
   criterion is — and note that it is not "is the symbol exported".
3. `own_package_links` — say why this case was never in scope for the
   restriction.
4. The opt-out flag exists and is documented as temporary. Say what you would do
   in a library that depends on a now-refused linkname, and what you would not do.
5. This change broke real code on upgrade. Describe, in general terms, the kind of
   library that was affected and why it had reached for `linkname` in the first
   place.
