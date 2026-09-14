# 03 — What a frame costs

Five functions in `frames.go`, three of them holding 1 KiB arrays. Predict what
the compiler gave each one for a stack frame.

| slot | question |
|---|---|
| `helper_frame_is_zero` | does a leaf that returns `n * 2` get a frame at all? |
| `unused_array_costs_frame` | does an array the compiler can see through cost anything? |
| `one_array_costs_frame` | does one genuinely live array cost a frame? |
| `two_live_is_about_double_one` | are two simultaneously-live arrays about twice one? |
| `disjoint_scopes_reduce_frame` | does putting the second array in its own scope shrink the frame? |

```
undergo verify compiler/03-stack-slots
```

The last one is the question worth having an opinion about, and the common
opinion is wrong.

Every slot is a **comparison**, never a byte count, because a frame size is not
portable — the same function here measures eight bytes larger on arm64 than on
amd64. Written question 4 is about what that means for testing.

## Questions to answer in writing

1. "Disjoint live ranges" — define it precisely, and say what has to be true of
   two locals for a compiler to be *allowed* to share one slot between them.
2. `disjoint_scopes_reduce_frame` comes out the way it does even though the two
   arrays provably never coexist. Say what these arrays have in common that
   blocks the sharing — and note that it is the same property that made them
   survive at all.
3. `unused_array_costs_frame` and `one_array_costs_frame` differ by one thing.
   Say what, and what that tells you about reading frame sizes as a proxy for
   "how much memory does this function use".
4. Frame size differs between architectures. Say what you would and would not
   assert about it in a test you expect to run on more than one machine.
5. A deeply recursive function multiplies its frame by its depth. Say what that
   means for the two-array version, and what you would change first.
