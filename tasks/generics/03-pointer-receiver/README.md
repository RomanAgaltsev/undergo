# 03 — Calling a pointer method on a slice of values

`SetStamp` is declared on `*Event`. The caller holds a `[]Event`.

```go
events := []Event{{Name: "a"}, {Name: "b"}}
StampAll[Event](events, time.Now())
```

`Event` does not implement `Stampable` — its method set does not include a
pointer method — so the obvious constraint does not compile:

```go
func StampAll[T Stampable](items []T, at time.Time) []T   // []Event will not satisfy this
```

Implement `StampAll` with the constraint already written in `stamp.go`, which
says that `PT` is *the pointer to* `T` **and** that `PT` is `Stampable`.

## The contract

- Every element of the **original** slice is stamped, not a copy of it.
- The returned slice shares its backing array with the argument.
- A nil or empty slice returns without panicking.
- The walk allocates nothing.

```
undergo verify generics/03-pointer-receiver
```

Note that `T` must still be written explicitly at the call site — `StampAll[Event](...)`
— while `PT` is not. Written question 2 is about why.

## Questions to answer in writing

1. State, in terms of method sets, why `[]Event` cannot satisfy `[]T` where
   `T Stampable`. Which method set does `Event` have, and which does `*Event`?
2. `StampAll[Event](events, when)` names one type argument and leaves the other
   to inference. Why can `PT` be inferred once `T` is known, and why can `T` not
   be inferred from `events` alone?
3. Write the version of this function that existed before type parameters —
   the one operating on `[]*Event` — and say what it forces on every caller who
   has a `[]Event` in hand.
4. Converting each element to a `Stampable` interface value and calling through
   that also compiles, and also allocates nothing. Explain why boxing a *pointer*
   costs no allocation when boxing a struct value does, and say what the
   interface version costs instead.
