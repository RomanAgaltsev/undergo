# 09 — The hash that is only a promise within one process

`hash/maphash` is the hash the runtime uses for maps, exposed as a package. It
takes a `Seed`, and `maphash.MakeSeed()` makes a new random one each time it is
called.

The same key, `"undergo"`, is hashed four ways:

| slot | question |
|---|---|
| `same_seed_same_key_agrees` | one seed, the same key hashed twice — same answer? |
| `different_seeds_agree` | two seeds from `MakeSeed()`, same key — same answer? |
| `string_and_bytes_agree` | `maphash.String(s, "undergo")` vs `maphash.Bytes(s, []byte("undergo"))` |
| `two_processes_agree` | the same program run twice, each making its own seed |

Predict `true` or `false` for each.

```
undergo verify edges/09-maphash-seeds
```

The judge is the package itself, plus one small program in `testdata/twice` that
the test runs twice and compares.

Note what is *not* asked: no slot grades a hash value. A hash value from this
package is not a fact about the key.

## Questions to answer in writing

1. State, in one sentence, what `hash/maphash` guarantees and what it refuses to
   guarantee. Your sentence should make `two_processes_agree` obvious.
2. `string_and_bytes_agree` is the one comparison that could have gone either
   way for a reason unrelated to seeds. Say what it is telling you, and why the
   package documents it explicitly.
3. Name three uses where this package is exactly right, and three where reaching
   for it is a bug that will not show up in testing. Be specific about when the
   bug surfaces.
4. Map iteration order is randomised in Go, and this hash is why. Say what
   attack the randomisation prevents, and what it costs.
5. You need a hash that agrees across processes and releases. Say which package
   you reach for instead, and what property you are now responsible for that
   `maphash` was handling for you.
