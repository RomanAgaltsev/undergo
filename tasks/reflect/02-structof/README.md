# 02 — Build a struct type at run time

Implement `Project`, which returns a **new** struct value holding only the named
fields of the one it was given:

```go
row := Row{ID: 7, Name: "ada", Email: "ada@example.com", Balance: 42}

got, err := structof.Project(row, "ID", "Balance")
json.Marshal(got)   // {"id":7,"balance":42}
```

The type does not exist until `Project` builds it. This is what a column
projection does when the columns are not known until the query arrives.

## The contract

- Exactly the named fields, **in the order they were asked for**, with their
  original types and struct tags — so `omitempty` still behaves as it did.
- A pointer to a struct counts as a struct.
- No fields named: an empty struct value, no error.
- An unknown field: an error naming the field.
- An unexported field: an **error**, not a panic.
- A non-struct argument: an error.

```
undergo verify reflect/02-structof
```

That "not a panic" is the sharp edge of this task, and it is not a style
preference — one of these cases fails in a way you cannot recover from after the
fact.

## Questions to answer in writing

1. `reflect.StructOf` does not return an error. List what it panics on, and say
   which of those your implementation has to rule out in advance.
2. The result comes back as `any`. Say what a caller can and cannot do with it,
   and why `json.Marshal` is happy while a type assertion is not.
3. Compare this with generating the same struct at compile time. Answer in terms
   of allocations per call and of what `reflect` caches between calls — does
   projecting the same two fields twice build the type twice?
4. `Type.FieldByName` returns a `StructField` whose `Index` is a slice, not an
   int. Say when it has more than one element, and what that means for a field
   you were asked to project.
