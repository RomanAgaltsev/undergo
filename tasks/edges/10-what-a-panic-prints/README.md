# 10 — What a panic message knows about its value

The `panic:` line is printed by the **runtime**, while the stack is unwinding,
without calling into `fmt`. So what it can say about your value is limited by
what the runtime knows how to format — not by what `fmt` could have managed.

Four programs, each panicking with one value. The types:

```go
type Code int

type Str string

type MyErr struct{ S string }

func (e MyErr) Error() string { return e.S }
```

And the four panics:

| slot | the value | predict |
|---|---|---|
| `named_int_line` | `Code(7)` | the whole `panic: ...` line, exactly |
| `named_string_line` | `Str("boom")` | the whole line, exactly |
| `error_value_line` | `MyErr{"broke"}` | the whole line, exactly |
| `plain_struct_line_has_address` | `struct{ A int }{3}` | does its line contain `0x`? |

The package is `main`, so a type named in the output is spelled `main.Something`.
Quote your strings exactly as the runtime would.

```
undergo verify edges/10-what-a-panic-prints
```

The judge is the runtime: the test runs each program and grades the line
beginning `panic: `.

Three of these four lines are worth one lesson each, and they are three
different lessons.

## Questions to answer in writing

1. Two of the values are named types with an underlying basic type, and both
   print their type. Say why that is more useful than printing the value alone,
   with an example of a panic you would otherwise misdiagnose.
2. `MyErr` prints differently from the other three, and the difference is not
   about its shape. Say what the runtime checked, and what it called.
3. The fourth value prints an address. Say why the runtime stops there for that
   type but not for the others — and say what it would have cost to do better.
4. From your answer to 3, state the rule for what to panic *with*, if you want
   the message to be readable when it reaches a log at 3am.
5. `recover()` hands you the value itself, not this string. Say when the printed
   line is all you will ever have, and what that means for choosing panic values
   in a library other people call.
