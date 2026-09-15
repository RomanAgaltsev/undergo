# 05 — The alignment nobody guarantees

Three structs holding the same two things — a flag and a 64-bit counter —
declared three different ways:

```go
type Plain struct   { Flag bool; N int64 }
type Guarded struct { Flag bool; N atomic.Int64 }
type Hosted struct  { _ structs.HostLayout; Flag bool; N int64 }
```

| slot | question |
|---|---|
| `int64_align` | `unsafe.Alignof(int64(0))` |
| `atomic_size` | `unsafe.Sizeof` of an `atomic.Int64` |
| `plain_offset` | where `N` begins in `Plain` |
| `guarded_offset` | where `N` begins in `Guarded` |
| `hostlayout_same_size` | is `Hosted` the same size as `Plain`? |

The last one is the interesting slot. `structs.HostLayout` exists for a reason;
predict whether that reason shows up here.

```
undergo verify layout/05-alignment-64bit
```

## Questions to answer in writing

1. On a 32-bit platform `unsafe.Alignof(int64(0))` is **4**, not 8. Say what
   breaks when a 64-bit atomic operation lands on a 4-aligned address, and why
   `sync/atomic`'s documentation makes that the *caller's* problem rather than
   the package's.
2. `atomic.Int64` solves that without asking you to reorder anything. Say how,
   and what it costs — compare `atomic_size` against the size of a bare `int64`.
3. `structs.HostLayout` changed nothing here. Say what it is for, why it is a
   no-op on the platform you measured, and name a situation where it is not.
4. `layout/01` in this repository once asserted that "Go lays a struct out in
   declaration order". That sentence is wrong as stated, and
   `structs.HostLayout` exists *because* it is wrong. Say what the correct
   version is, and what the difference between the two claims lets a compiler
   do.
