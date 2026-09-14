# 04 — Interning a stream of labels

A metrics pipeline receives hundreds of thousands of records a second. Each one
carries a label drawn from a small fixed set — sixteen service names, say — and
each arrives as a freshly allocated string off the wire. Stored as strings, every
record keeps its own copy alive.

Implement `Intern`, `Value` and `InternAll` on top of `unique.Handle`.

```go
tag := intern.Intern("service-name-label-07")
intern.Value(tag)   // -> "service-name-label-07"
```

## The contract

- Equal strings give handles that compare **equal with `==`**, whether or not the
  two strings share a backing array.
- Different strings give different handles.
- `Value` round-trips, including for the empty string.
- `InternAll` preserves order, and equal entries intern to the same handle.
- Interning 200,000 records drawn from 16 labels retains **less than half** what
  storing the strings does.

```
undergo verify weak/04-unique-interning
```

That last one is measured in the same run against its own baseline rather than
against a fixed number of bytes — written question 4 is about why it has to be.

## Questions to answer in writing

1. Comparing two handles is a pointer comparison. Say what that buys over `==` on
   the strings themselves, and name the case where it buys nothing at all.
2. What keeps an interned value alive, and what eventually frees it? Be specific
   about what has to become unreachable first.
3. Interning is not free. Describe the workload where it **loses** — and say what
   property of the key distribution decides it. Then say roughly where the
   break-even sits.
4. This task asserts a ratio between two measurements taken in the same run,
   never a byte count. Say what would go wrong with a byte count, and name two
   things that would move it without anything being wrong with the code.
5. `Tag` is declared as an alias (`type Tag = unique.Handle[string]`) rather than
   a defined type. Say what would break for callers if it were a defined type.
