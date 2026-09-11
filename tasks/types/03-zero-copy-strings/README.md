# 03 — Convert between string and []byte without copying

Implement `StringToBytes` and `BytesToString` in `zerocopy.go`.
`zerocopy_test.go` is the specification; do not weaken it.

`string(b)` and `[]byte(s)` already do this — by copying. Copying is why they are
safe. Your versions must share the backing array instead, which is why the doc
comments carry warnings that the standard conversions do not need.

Three things are tested, and the third is the one that matters:

- the content round-trips, including the empty and nil cases;
- **neither conversion allocates**, measured with `testing.AllocsPerRun` — a
  correct-but-copying implementation passes everything else and fails here;
- the sharing is observable: mutating the slice changes the string.

```
undergo start types/03-zero-copy-strings
undergo verify types/03-zero-copy-strings
```

## Questions to answer in writing

1. Which of the six legal `unsafe.Pointer` conversion patterns does your
   implementation rely on — or does it avoid `unsafe.Pointer` entirely?
2. What do `unsafe.StringData` and `unsafe.SliceData` promise for an empty
   string and a zero-length slice? Quote the documentation, then say what your
   code does about it and why the tests would not have caught it.
3. The standard library performs this conversion internally in several places.
   Find one, and state the invariant that makes it safe there but not in general.
