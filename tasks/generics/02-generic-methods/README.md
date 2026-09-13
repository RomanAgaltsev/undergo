# 02 — A typed Map, as a method

`Set[T]` is an ordered collection. Implement `Map`, which applies a function to
every item and returns the results **as a different type**:

```go
func (s *Set[T]) Map[R any](
	ctx context.Context,
	fn func(context.Context, T) (R, error),
) ([]R, error)
```

`R` belongs to the method, not to the receiver. Before Go 1.27 a method could not
declare its own type parameters, so this had to be a package-level function
taking the `Set` as an argument.

Callers write no type arguments:

```go
s := set.New(1, 2, 3)

labels, err := s.Map(ctx, func(ctx context.Context, v int) (string, error) {
	return strconv.Itoa(v), nil
})
```

## The contract

- Results come back **in input order**, one per item.
- An empty `Set` returns a nil slice and a nil error, and never calls `fn`.
- `ctx` is checked **before** each item, so a context already cancelled on entry
  returns `ctx.Err()` without calling `fn` at all.
- The first error stops the walk — later items are not passed to `fn` — and comes
  back as a `*MapError` carrying the **index** of the item that failed, wrapping
  the original error so `errors.Is` still reaches it.
- Exactly **one** slice is allocated.

```
undergo verify generics/02-generic-methods
```

That last rule is a real test, not advice. An implementation that appends to a nil
slice passes everything else and fails it.

## Questions to answer in writing

1. Write the pre-1.27 version of this API — the one that had to be a function.
   What did callers lose, beyond the spelling?
2. At the call site above, `T` and `R` are both inferred. Say where each one is
   inferred *from*, and what you would have to write if `fn` were `nil`.
3. This returns after the first failure. A batch processor usually wants every
   error instead. What would the signature become, and what does that cost the
   caller who only wanted the happy path?
4. `Map` is a method on `*Set[T]`, and interfaces may not declare type parameters.
   What does that mean for a `Mapper` interface you might want to write, and how
   would you work around it?
