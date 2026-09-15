# 05 — Six legal patterns, and eight snippets

`unsafe.Pointer`'s documentation lists **six** legal conversion patterns and
says that any other use is invalid. `testdata/patterns` holds eight snippets.
Five follow a listed pattern; three do not.

| snippet | what it does |
|---|---|
| `p1` | `*float64` → `*uint64` through `unsafe.Pointer` |
| `p2` | pointer arithmetic folded into a single conversion expression |
| `p3` | `unsafe.Pointer` → `uintptr`, printed, never converted back |
| `p4` | `unsafe.Slice` over a pointer the caller owns |
| `p5` | `unsafe.String` over a slice's data pointer |
| `v1` | address stored in a `uintptr` **variable**, converted back later |
| `v2` | `unsafe.Add` walking well past the end of the allocation |
| `v3` | address round-tripped through a **struct field** |

For each one, predict whether `go vet -unsafeptr` rejects it.

| slot | snippet |
|---|---|
| `vet_rejects_p1` … `vet_rejects_p5` | the five patterns |
| `vet_rejects_v1`, `vet_rejects_v2`, `vet_rejects_v3` | the three violations |

Three snippets are invalid. Predict carefully whether vet finds all three.

```
undergo verify edges/05-unsafe-pointer-rules
```

## Questions to answer in writing

1. List the six legal patterns in your own words.
2. One violation is **not** rejected by vet. Name it, say what is wrong with it,
   and say what a moving garbage collector would do to it.
3. Say why vet is conservative here — what would it have to compute in order to
   catch the one it missed, and what would it cost in false positives?
4. `checkptr` is enabled by `-race`. Say what it catches that vet cannot, and
   why it needs the program to **run** rather than to be read.
5. Go's collector does not currently move heap objects. Say what would break
   first if it started, and whether "it works today" is a defence.
6. The snippets carry `//go:build ignore`. Say why a task about code that fails
   vet cannot keep that code in its own package.
