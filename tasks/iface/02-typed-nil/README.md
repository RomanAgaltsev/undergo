# 02 — The nil that is not nil

An interface value holding a nil pointer is the most consequential confusion in
Go. Six comparisons over the same value.

| slot | the comparison |
|---|---|
| `true_nil_is_nil` | `var err error; err == nil` |
| `nil_pointer_in_interface_is_nil` | `var p *MyErr; var err error = p; err == nil` |
| `returned_nil_pointer_is_nil` | the same, arriving from a function returning `*MyErr` |
| `comparing_to_typed_nil_works` | `err == (*MyErr)(nil)` |
| `reflect_says_nil` | `reflect.ValueOf(err).IsNil()` |
| `errors_is_nil_agrees_with_equality` | does `errors.Is(err, nil)` give the same answer as `err == nil`? |

```
undergo verify iface/02-typed-nil
```

The `review/typed-nil` track drills spotting this in code. This one asks what the
machine actually does, which is the half that makes the rule stick.

## Questions to answer in writing

1. State the rule for when an interface value equals nil, in terms of both of its
   words.
2. This bites hardest in a function that returns a concrete error type. Write the
   signature that causes it and the one that avoids it, and say which line is
   actually at fault.
3. `go vet` has a check in this area. Name it, and say precisely what it catches
   and what it does not.
4. `reflect` and `==` disagree about this value. Say which question each one is
   answering — neither is wrong.
5. `errors.Is(err, nil)` agrees with `==` here. Say why, and whether you would
   ever write it.
