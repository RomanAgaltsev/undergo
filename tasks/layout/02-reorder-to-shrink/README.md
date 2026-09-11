# 02 — What field order costs across a million records

`Event` and `Packed` hold the same six fields in different orders. The service
keeps a million of them in a slice.

Read `record.go`. Without running anything, predict:

| slot | question |
|---|---|
| `sizeof_event` | `unsafe.Sizeof(Event{})` |
| `sizeof_packed` | `unsafe.Sizeof(Packed{})` |
| `offset_event_kind` | `unsafe.Offsetof(Event{}.Kind)` |
| `waste_per_record` | the difference between the two sizes, in bytes |
| `waste_mib` | that difference across `Million` records, in whole MiB |

The last two are arithmetic on the first two: if the sizes are right and the
waste is wrong, you have found a different kind of mistake.

```
undergo verify layout/02-reorder-to-shrink
```

## Questions to answer in writing

1. `Packed` orders its fields largest-alignment-first. Is that always optimal,
   or does it only happen to be optimal here? Construct a counter-example, or
   argue that none exists.
2. `go vet` has a `fieldalignment` check that finds this mechanically. Why is it
   not on by default?
3. The service could also shrink `Retries` from a `uint16` to a `uint8`. How
   many bytes would that save per record — and why is the answer probably zero
   for one of these two structs but not the other?
