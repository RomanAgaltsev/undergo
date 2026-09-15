# 04 — What the built-in types actually are

`Header` nests one of each of the types whose size people guess at:

```go
type Header struct {
	S   string
	Sl  []int
	M   map[int]int
	Ch  chan int
	F   func()
	I   any
	Ptr *int
}
```

Every answer below is a **compile-time constant**. Nothing is measured, and
nothing depends on how much data the values hold — which the last slot is there
to make concrete.

| slot | question |
|---|---|
| `string_size` | `unsafe.Sizeof` of the `string` field |
| `slice_size` | of the `[]int` field |
| `map_size` | of the `map[int]int` field |
| `iface_size` | of the `any` field |
| `total_size` | of the whole `Header` |
| `map_offset` | `unsafe.Offsetof(h.M)` |
| `big_slice_size` | `Sizeof` of a slice holding **one million** ints |

Four of the seven fields have the same size. Decide whether that means they are
the same shape.

```
undergo verify layout/04-headers
```

## Questions to answer in writing

1. A map, a channel and a func value are all the same size. Say what those
   bytes actually are in each case, and where a map's buckets live.
2. `big_slice_size` is the same as `slice_size`. Say why, and say what that
   means for passing a large slice to a function — then say what it means for
   passing a large **array**.
3. Add up the seven field sizes and compare with `total_size`. Account for any
   difference, and say what property of these particular fields makes the answer
   come out as it does.
4. Every answer here would change on a 32-bit platform. Say which, and by how
   much. Then say whether this task should have pinned `requires.arch`, and
   defend your answer.
