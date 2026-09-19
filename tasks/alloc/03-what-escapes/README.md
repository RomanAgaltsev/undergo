# 03 — Which of these values escape to the heap?

Four functions, each building one `Holder`. Some of those Holders are stack
allocated and some are moved to the heap.

Read `escape.go`. Without running anything, predict `true` or `false` for each:

| slot | function |
|---|---|
| `local_escapes` | `Local` |
| `returned_escapes` | `Returned` |
| `stored_escapes` | `Stored` |
| `interfaced_escapes` | `Interfaced` |

The judge is not a rule of thumb: the test compiles this package with
`-gcflags=-m` and reads the compiler's own escape diagnostics, attributing each
one to the function that contains it.

```
undergo verify alloc/03-what-escapes
```

## Questions to answer in writing

1. `Interfaced` returns `any(h)` by value. Nothing anywhere takes the address of
   `h`. So why is its answer what it is, and what exactly is on the heap?
2. Escape analysis is usually summarised as "does this value outlive its
   frame?". Describe a value that outlives its frame and still does not escape.
3. Run `go build -gcflags=-m=2 .` and explain one line of the extra output that
   `-m` alone does not print.
