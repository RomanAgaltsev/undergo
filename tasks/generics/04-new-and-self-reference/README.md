# 04 — Two things Go 1.26 let you write

Implement `SumAll`, which folds a slice using the element type's own `Add`:

```go
total, err := fold.SumAll([]Money{{100}, {250}, {5}})   // -> *Money{355}
```

The signature is already written for you, and it leans on **both** of Go 1.26's
changes to the language specification. One of them is in the constraint:

```go
type Adder[A Adder[A]] interface {
	Add(A) A
}
```

Before 1.26 that was rejected outright — a generic type could not refer to itself
in its own type parameter list.

## The contract

- The fold runs left to right, starting from the first element.
- An empty slice (or a nil one) returns a nil pointer and `ErrEmpty`.
- A one-element slice returns a pointer to that element.
- The result does not point into the caller's slice.
- Exactly **one** allocation happens, and it is the result.

```
undergo verify generics/04-new-and-self-reference
```

Allocate that result with `new` taking an **expression** — the second 1.26
change. No test can tell the two spellings apart, so this one is on your honour;
question 2 is where you show you know the difference.

## Questions to answer in writing

1. Write `Adder` as it had to be written before Go 1.26. What did the extra type
   parameter buy, and what could a caller get wrong with it that they cannot get
   wrong now?
2. Compare `return new(total), nil` with `p := new(A); *p = total; return p, nil`.
   What does `new` accept now that it did not before, and does the difference
   show up in the generated code, the allocation count, or neither?
3. `SumAll` needs an error for the empty case because the constraint promises
   only `Add`. Name the constraint you would add to remove the error from the
   signature, and say what it would cost the types that want to satisfy `Adder`.
4. `Money.Add` has a value receiver, so `Money` satisfies `Adder[Money]`. What
   would change for a caller if `Add` were declared on `*Money` instead?
