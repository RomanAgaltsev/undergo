# 09 — A key that is not a field name

Every other task in this track pins the go line and varies the compiler. This
one does the opposite: **one compiler, two language versions.** Nothing is
downloaded.

Go 1.27 changed what may appear on the left of a colon in a struct literal. It
used to be a field name of the struct type itself. Now it may be any valid
field selector for that type — so a literal can reach a field of an embedded
struct without naming the embedded struct.

These are the types. The name `X` exists at two different depths.

```go
type Deep struct {
	Z int
	X int
}

type Mid struct {
	Deep
	W int
	X int
}

type Top struct {
	Mid
	Y int
}
```

## Legality

Two literals are compiled — never run — at `go 1.26` and again at `go 1.27`:

```go
Top{W: 2}   // W is promoted through one level of embedding
Top{Z: 3}   // Z is promoted through two
```

| slot | literal | go line |
|---|---|---|
| `promoted_once_compiles_under_go1_26` | `Top{W: 2}` | `go 1.26` |
| `promoted_twice_compiles_under_go1_26` | `Top{Z: 3}` | `go 1.26` |
| `promoted_once_compiles_under_go1_27` | `Top{W: 2}` | `go 1.27` |

Predict `true` or `false` for each.

## Resolution

The second half is not about legality. This program is run at `go 1.27`:

```go
t := Top{X: 7}
fmt.Printf("%d|%d", t.Mid.X, t.Mid.Deep.X)
```

`X` is a valid selector at depth 1 and at depth 2. Exactly one of them ends up
holding `7`.

| slot | answer |
|---|---|
| `shadowed_key_sets_depth` | `1` or `2` — the depth of the `X` the key reached |

```
undergo verify versions/09-promoted-field-keys
```

**This task downloads nothing.** The go line is the only thing that varies, and
one installed toolchain answers at both settings.

## Questions to answer in writing

1. The compiler rejects one of these literals with a message naming a version.
   Which mechanism decides that — the toolchain you have installed, or
   something in the module? What is the difference, and which one can you
   change without downloading anything?
2. Before you look: if a name is a valid selector at two depths, which one
   should a struct literal pick, and why? Is that the same rule the selector
   expression `t.X` uses?
3. `Top{W: 2}` sets a field of an embedded struct without naming that struct.
   What does that cost a reader who is trying to find where `W` lives? Would
   you use this in code you maintain?
