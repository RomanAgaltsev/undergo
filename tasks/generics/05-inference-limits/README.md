# 05 — Which of these calls need a type argument?

Six call sites, each in its own complete program, none of them writing an
explicit type argument. Predict `true` if it compiles and `false` if it does not.

| slot | call site |
|---|---|
| `call1_compiles` | `First(nil)` |
| `call2_compiles` | `Zero()` |
| `call3_compiles` | `Sum(int32(1), int64(2))` |
| `call4_compiles` | `Apply(3, strconv.Itoa)` |
| `call5_compiles` | `Keys(map[string]int{"a": 1})` |
| `call6_compiles` | `var f func(int) int = Identity` |

The declarations they are calling:

```go
// call1
func First[T any](s []T) T { return s[0] }

// call2
func Zero[T any]() T { var z T; return z }

// call3
type Number interface{ ~int32 | ~int64 }
func Sum[T Number](a, b T) T { return a + b }

// call4
func Apply[T, R any](v T, f func(T) R) R { return f(v) }

// call5
func Keys[M ~map[K]V, K comparable, V any](m M) []K

// call6
func Identity[T any](v T) T { return v }
```

Answering all six the same way scores exactly half. One of the three that fails
does so for a different reason than the other two, and question 2 is about that.

```
undergo verify generics/05-inference-limits
```

The snippets are in `testdata/`, which the go tool ignores — three of them do not
compile, and a task directory that did not build would fail this repository's own
gates. You can build them yourself, but do it after you have committed your
answers.

## Questions to answer in writing

1. State the general rule for when a type parameter can be inferred. Your rule
   should predict all six of these without special cases.
2. Two of the failures report that the compiler *cannot infer* a type parameter.
   The third reports something else. Quote both messages and say what the
   difference is — absence of information versus what?
3. Which failures could be fixed by reordering the type parameter list or the
   function's arguments, and which cannot be fixed at any cost short of writing
   the type argument?
4. `call6` assigns a generic function to a variable of a concrete function type.
   Nothing is being *called*, so there is no argument list to infer from — say
   where `T` comes from instead, and what the compiler is unifying against.
