# 04 — A reader that never blocks the writer

Implement a **sequence lock**: one writer, many readers, readers that never
block the writer and never take a lock.

The protocol is short. The writer bumps a counter, writes the payload, and bumps
it again. A reader reads the counter, reads the payload, reads the counter
again, and retries if it changed or was odd.

```go
func (s *Seqlock) Write(p Point)
func (s *Seqlock) Read() Point
```

`Point` is three `int64`s, so a write to it is not one instruction and a torn
read is both possible and detectable.

## What you have to satisfy

| test | requirement |
|---|---|
| `TestNoTornReads` | with one writer and four readers, no reader ever sees a `Point` assembled from two different writes |
| `TestReadSeesLatest` | a read after a completed write returns it |
| `TestReadsAreConcurrent` | readers do not exclude one another — a mutex-based "seqlock" fails here |
| the race gate | `go test -race` is **silent** |

That last one is half the exercise. A design that passes the first three and
trips the detector is not finished.

```
undergo verify memmodel/04-seqlock
```

## Questions to answer in writing

1. The sequence number is incremented twice per write rather than once. Say what
   its parity means to a reader, and what would break if you incremented once.
2. A C seqlock reads the payload with **plain** loads and separates them from
   the sequence reads with explicit memory fences. Go has no fence. Say what you
   used instead, what it cost you, and why the plain-load version is not merely
   unidiomatic here but actually a bug.
3. A reader can be starved: if the writer is fast enough, `Read` may retry
   forever. Say whether that is a livelock, and what you would change if readers
   had to be guaranteed progress.
4. The payload is three `int64`s. Say what changes if it is one, and whether a
   seqlock is still worth having. Then say what changes if it contains a
   **pointer**, and why that case is much worse than the other two.
