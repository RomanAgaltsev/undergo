# 01 — What an interface value is

An interface value is not a pointer and not a struct. Measure it.

| slot | question |
|---|---|
| `empty_interface_words` | how many machine words is an `any`? |
| `method_interface_words` | how many is a one-method interface? |
| `slice_words` | how many is a `[]int`, for comparison? |
| `boxing_small_int_allocs` | does `BoxSink = Small` (7) allocate? |
| `boxing_large_int_allocs` | does `BoxSink = Runtime` (1 << 20) allocate? |
| `boxing_pointer_allocs` | does boxing a `*Dog` allocate? |
| `boxing_struct_allocs` | does boxing a `Dog` value allocate? |

```
undergo verify iface/01-eface-vs-iface
```

The four boxing slots are not all the same, and the reasons differ: one is about
a cache, one is about what fits, one is about what a pointer already is.

The sizes are in **words** rather than bytes so the answer does not depend on the
machine — but it still depends on the machine being 64-bit. Written question 4 is
about that.

## Questions to answer in writing

1. Name the two words in an empty interface value and the two in a one-method
   interface. Say what differs between them, and what the difference means at a
   call site.
2. Boxing a pointer costs nothing while boxing a struct costs an allocation.
   Explain in terms of what the data word holds.
3. A `[]any` holding a million pointers against a `[]*T` holding the same
   million: say how much more the first costs, and why.
4. This task declares `requires.arch: [amd64, arm64]`. Say what the answers would
   be on a 32-bit target and which slots would change.
5. `any` and `interface{ Speak() string }` are the same size. Say what converting
   between them costs, and when it happens.
