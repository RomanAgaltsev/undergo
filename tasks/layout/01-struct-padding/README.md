# 01 — Padding, alignment, and the trailing zero-size field

Three structs. Two of them hold identical fields in a different order. One of
them ends in `struct{}`.

Read `padding.go`. Without running anything, predict:

| slot | question |
|---|---|
| `sizeof_header` | `unsafe.Sizeof(Header{})` |
| `sizeof_reordered` | `unsafe.Sizeof(Reordered{})` |
| `sizeof_trailing` | `unsafe.Sizeof(Trailing{})` |
| `offset_kind` | `unsafe.Offsetof(Header{}.Kind)` |
| `align_header` | `unsafe.Alignof(Header{})` |

Write your answers into `prediction.yaml`, then:

```
undergo verify layout/01-struct-padding
```

Guessing after running the test is the one way to waste this task.

## Questions to answer in writing

1. Which two fields in `Header` are separated by padding, and how much?
2. `Trailing` ends with a zero-size field. Why does that cost anything?
3. On a 32-bit target, which of your five answers change, and why?
