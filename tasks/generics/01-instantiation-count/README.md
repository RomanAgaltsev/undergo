# 01 — How many bodies does the compiler actually emit?

`shapes.go` instantiates one generic **method** and one generic **function** over
the same five type pairs:

    (int, string)   (int, int64)   (string, int)   (MyInt, string)   (*int, bool)

Go does not stencil a fresh copy per type pair, and it does not hand everything to
one shared body either. Read `shapes.go`, then predict what the compiler put in
the archive.

| slot | question |
|---|---|
| `method_shape_bodies` | GC-shape bodies emitted for `MapM` |
| `method_concrete_wrappers` | fully concrete (non-shape) `MapM` symbols |
| `method_dicts` | `MapM` dictionaries |
| `func_shape_bodies` | GC-shape bodies emitted for `MapF` |
| `func_concrete_wrappers` | fully concrete (non-shape) `MapF` symbols |
| `func_dicts` | `MapF` dictionaries |

Five call sites each. If all six of your answers are 5, you have made the same
mistake three times. If your two rows are identical, you have made a different
one.

```
undergo verify generics/01-instantiation-count
```

The test builds this package into an **archive** and reads its symbol table.

Note that a symbol table holds more than function bodies. Deciding what counts as
an instantiation is part of the question.

## Questions to answer in writing

1. There are fewer bodies than dictionaries. Which two call sites share a body,
   and what rule decides that two types may share one? Name the rule, not just
   the pair.
2. If a body is shared, how does the shared code know which concrete type it is
   working on? Say what a dictionary holds, and who passes it.
3. The method and the function differ in exactly one column. Which one, and why
   does binding `T` on the receiver rather than at the call change what the
   compiler must emit?
4. Rebuild with `-gcflags=all=-l` and count again. Then link a `main` that calls
   `Callsites` and count in the **binary** instead. One of those changes nothing
   and the other changes a great deal — explain which, and why this task grades
   an archive.
