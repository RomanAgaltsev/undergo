# 03 — When the compiler knows the type

A call through an interface is a load and an indirect jump — unless the compiler
can work out which concrete type is on the other side, in which case it replaces
the call with a direct one. That is **devirtualization**, and it reports itself.

Six call sites in `dispatch.go`. Predict which get resolved.

| slot | the call site |
|---|---|
| `one_type_devirtualized` | one concrete type assigned, then called |
| `reassigned_devirtualized` | two types assigned on different branches |
| `via_parameter_devirtualized` | the interface arrives as a parameter |
| `via_global_devirtualized` | a package-level interface variable |
| `from_constructor_devirtualized` | assigned from a function returning a concrete type |
| `in_slice_devirtualized` | an element of a `[]Shape` |

```
undergo verify iface/03-devirtualization
```

Two of the six resolve. One pair differs only by a line that may never execute.

## Questions to answer in writing

1. State what the compiler must be able to prove. Then say why
   `reassigned_devirtualized` comes out as it does even when `pick` is always
   false at run time.
2. `via_global_devirtualized` involves a variable assigned exactly once, in this
   package, to one concrete type. Say why that is still not enough.
3. Devirtualization enables something else, and that something is usually the
   larger win. Name it and explain the order of events.
4. `in_slice_devirtualized` — say what would have to be true of the slice for the
   answer to change, and whether any compiler could establish it.
5. PGO extends devirtualization. Say what a profile provides that static analysis
   cannot, and what the compiler must still emit as a result.
