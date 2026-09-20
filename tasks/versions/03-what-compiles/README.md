# 03 — Is this valid Go?

The question has no answer until you name a version.

Five snippets. Each is **compiled, never run**, by the one toolchain you have,
in a module declaring the `go` line shown. Predict `true` if it compiles and
`false` if it does not.

### `clear_min_max_at_go120` — at `go 1.20`

```go
m := map[string]int{"a": 1}
clear(m)
fmt.Println(min(3, 1, 2), max(3, 1, 2), len(m))
```

### `range_int_at_go121` — at `go 1.21`

```go
for range 3 {
}
```

### `range_func_at_go122` — at `go 1.22`

```go
func seq(yield func(int) bool) { /* ... */ }

for range seq {
}
```

### `new_expr_at_go125` — at `go 1.25`

```go
func f(x int) int { return x }

var _ = new(f(1))
```

### `generic_alias_at_go123` — at `go 1.23`

```go
type Pair[A any] struct{ V A }

type Alias[A any] = Pair[A]

_ = Alias[int]{V: 1}
```

And a sixth slot about the failures themselves:

| slot | question |
|---|---|
| `diagnostic_names_lang_flag` | does a refusal's message contain the text `-lang was set to`? |

```
undergo verify versions/03-what-compiles
```

The judge is the compiler: the test writes each snippet as a one-file module at
the stated go line and records whether the build succeeded.

Answering all five the same way is worth one point, and the one you get is not
the one you were aiming for.

## Questions to answer in writing

1. Four of these five are refused at the line shown. State the general rule you
   would infer from those four.
2. The fifth is not refused. Say what that does to your rule, and what it tells
   you about how the rule is implemented.
3. Read one diagnostic in full. It names two things: a required version and a
   flag. Say where the flag's value came from.
4. Your module declares `go 1.21` and your toolchain is 1.27. Say which of the
   two decides whether `for range 3` compiles, and which decides what
   optimisations the compiler applies.
5. Name the risk in lowering a module's go line to support older users, given
   what you now know about which features are gated and which are not.
