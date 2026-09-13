# 02 — What stopping does to an iterator

`Counted.Seq` yields `1..10` and counts its own calls to `yield`. Five consumers
range over it and leave the loop in five different ways. Read `contract.go` and
predict how many times `yield` was called in each case.

| slot | how the consumer leaves |
|---|---|
| `plain_yields` | runs to the end |
| `break_yields` | `break` at element 3 |
| `return_yields` | `return` from the enclosing function at element 3 |
| `labelled_yields` | `break outer` on the **second** pass, at element 3 |
| `panicking_yields` | panics at element 5, recovered outside the loop |
| `misbehaving_panics` | `true` or `false`: does ranging over `Misbehaving` panic? |

Two of these are equal. One is larger than 10. Read each consumer carefully
enough to say why before you look.

```
undergo verify iter/02-yield-contract
```

## Questions to answer in writing

1. `yield` is called by the iterator and returns a `bool`. Say who writes the
   function that returns `false`, and what the `for ... range` statement compiles
   into for a `break`.
2. `break` and `return` give the same count. Explain why the iterator cannot tell
   them apart, and name something the *consumer* can observe that differs.
3. `Misbehaving` keeps calling `yield` after it returned `false`. Quote the
   runtime's message in full, and say why this is a panic rather than a silently
   ignored call — what would break if it were ignored?
4. Put a `defer` inside the loop body, and another inside `Counted.Seq`. For the
   `break` case, say when each one runs relative to the other.
