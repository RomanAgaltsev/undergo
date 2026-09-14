# 01 — A cache that lets go

An ordinary map is a strong reference: whatever you put in it lives as long as
the map does. Build one that does not.

```go
c := cache.New[string, Payload]()
c.Put("a", payload)

// ... while something else still holds payload:
v, ok := c.Get("a")   // -> payload, true

// ... once nothing does, after a collection:
v, ok = c.Get("a")    // -> nil, false
```

## The contract

- `Put` stores a **weak** reference. The cache never keeps a value alive.
- `Get` returns the value while it is still reachable elsewhere, and
  `(nil, false)` once it is not — the same answer it gives for a key that was
  never stored.
- `Len` counts **entries, not live values**. An entry whose value has been
  collected is still an entry.
- `Reap` removes exactly those entries, and leaves the live ones alone.
- `Put` on an existing key replaces it.

```
undergo verify weak/01-weak-cache
```

`Len` disagreeing with reality is deliberate, not a bug to design out. Written
question 2 is about why it has to.

## Questions to answer in writing

1. A cache is supposed to keep things. Say what a weak cache is *for* — name a
   workload where it is right, and one where an ordinary LRU is right and this
   would be actively worse.
2. `Len` counts entries whose values are gone. Why can the map not drop an entry
   the moment its value is collected? Say what the runtime would have to offer
   for that to be possible, and what it would cost.
3. `Get` returns `(nil, false)` both for a collected value and for a key that was
   never stored. Should a caller be able to tell those apart? Argue it either
   way, then say which you would ship and why.
4. This cache is not safe for concurrent use. Say what you would add — and
   whether `weak.Pointer` itself needs protecting, or only the map.
5. The test that proves the cache lets go wraps its allocation in a function
   literal, and another test calls `runtime.KeepAlive`. Explain what each one is
   defending against. Both tests pass for the wrong reason without them.
