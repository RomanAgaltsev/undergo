# 04 — What a million-entry map actually costs

`Build` fills a `map[int64]int64` with 2²⁰ entries. The test measures the live
heap either side of building it, with a collection on both sides, and divides by
the number of entries.

Without running anything, predict:

| slot | question |
|---|---|
| `payload_bytes_per_entry` | the bytes of key and value you actually asked for, per entry |
| `bytes_per_entry_rounded` | what the map really costs per entry, rounded to the nearest 8 |

The first is arithmetic. The second is the task: predict the whole cost, not the
overhead ratio, and round it to the nearest multiple of 8.

```
undergo verify types/04-map-memory
```

## Questions to answer in writing

1. Where does the difference between the two numbers go? Account for it in terms
   of what a Swiss table stores per slot — not as "map overhead".
2. `BuildWithHint` builds the same map with the size known up front. Predict
   whether that changes the final memory, then measure it. Were you right, and
   what does the hint actually change?
3. 2²⁰ entries is close to a worst case for this measurement. Build maps of
   900,000 and 1,100,000 entries and compare bytes per entry. Explain the shape
   of what you see — and say what that implies about quoting "bytes per entry"
   for a map at all.
