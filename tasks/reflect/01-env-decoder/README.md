# 01 — Decode environment variables into a struct by tag

Implement `Decode` in `decoder.go`. `decoder_test.go` is the specification; do
not weaken it.

The interesting part is not the string parsing. It is settability: `reflect`
will hand you `Value`s you are not allowed to write to, and the rules for which
ones are which are the actual lesson.

```
undergo start reflect/01-env-decoder
undergo verify reflect/01-env-decoder
```

## Questions to answer in writing

1. `reflect.ValueOf(c).Field(0).SetString("x")` panics even when `c` is a
   struct with an exported first field. Why?
2. How do you detect an unexported field, and why can't you set one even via
   `unsafe`-free reflection?
3. What does this cost per call compared with a hand-written decoder, and where
   does the time go?
