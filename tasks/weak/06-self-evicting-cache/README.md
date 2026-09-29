# 06 — A cache that forgets on its own

`weak/01`'s cache kept dead entries until somebody called `Reap`. Build one that
needs no `Reap`: every entry removes itself once its value has been collected.

```go
c := cache.New[string, Payload]()
c.Put("a", v)          // weak: c does not keep v alive
v2 := &Payload{}
c.Put("a", v2)         // replaces the entry; v may now die

// ... after v is collected and its cleanup has run:
got, ok := c.Get("a")  // -> v2, true — v's cleanup must not remove v2's entry
```

## The contract

- `Put` stores a **weak** reference and replaces any existing entry for the key.
  `Put(k, nil)` removes `k`.
- An entry **removes itself** some time after its value is collected — no `Reap`.
- A removal only ever removes the entry that was made for the value that died.
  **A later `Put` under the same key survives the earlier value's removal.**
- `Get` returns the value while it is reachable elsewhere, `(nil, false)` otherwise.
- `Len` counts entries; a dead entry may be counted until it has removed itself.
- The cache is **safe for concurrent use**.

```
undergo verify weak/06-self-evicting-cache
```

The tests poll for up to five seconds for entries to disappear; a correct cache
takes milliseconds. Concurrency safety is checked two ways: the concurrent test
usually trips the runtime's own map check, but only the race detector is
reliable, and it needs a C toolchain. Without one, run the race gate in Docker:

```
task race:docker
```

## Questions to answer in writing

1. Describe exactly how an old value's cleanup could remove a newer entry, and
   why checking the *value* at cleanup time cannot prevent it.
2. Your cleanup almost certainly captures the cache. Say what that keeps alive,
   for how long, and whether it matters.
3. The documentation says cleanups may run concurrently with one another, unlike
   finalizers. Say what that forced on your design that a single-threaded program
   would not otherwise need.
4. `Len` can still count a dead entry. Say what that window is and why no design
   can close it.
5. Evicting lazily in `Get` needs no cleanups at all. Compare the two designs, and
   describe a workload where the lazy one grows without bound.
