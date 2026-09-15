# 05 — The conversions that are not conversions

Converting a `[]byte` to a `string` normally **copies**. Strings are immutable
and slices are not, so sharing the bytes would let a later write to the slice
change a string that already exists.

Five expressions, each containing the identical text `string(b)`:

```go
lookup[string(b)]          // map index
string(b) == "hello"       // comparison
for range string(b)        // range
switch string(b) { ... }   // switch
s = string(b)              // assignment
```

| slot | allocations per run |
|---|---|
| `map_index` | |
| `comparison` | |
| `range_over` | |
| `type_switch` | |
| `assignment` | |

Four of these five write the same conversion and do not cost the same. Work out
which one differs and why before running it.

```
undergo verify types/05-compiler-free-conversions
```

## Questions to answer in writing

1. Say what the compiler must prove about the resulting string in order to skip
   the copy, and name the property that makes the odd one out fail that proof.
2. `types/03` in this repository buys the same thing with `unsafe`. Given this
   task's result, say when reaching for `unsafe` is actually justified — and
   what the compiler already covers for free.
3. The assignment allocates because the string outlives the expression. Say
   what would have to change about Go's strings for it not to.
4. Reverse the question: does `[]byte(s)` have an equivalent set of free cases?
   Name one if there is, and say why the situation is not symmetric.
5. Every sink in this package is typed rather than `any`. Say what would happen
   to the measurements if they were `any`, and why that would be measuring the
   wrong thing.
