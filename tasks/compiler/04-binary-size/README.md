# 04 — What actually costs bytes

Six programs in `testdata/`, all built the same way. One is a baseline; the other
five each add exactly one thing to it. Predict which of them made the binary
bigger, and which added the most.

| slot | what the program adds |
|---|---|
| `fmt_grows_materially` | one `fmt.Println` |
| `reflect_grows_materially` | one `reflect.TypeOf` |
| `generic_grows_materially` | a generic function instantiated over four types |
| `iface_method_grows_materially` | a method reachable only through an interface |
| `unused_table_grows_materially` | a 512 KiB package-level array nothing references |
| `largest_addition` | the name of the snippet that added the most |

"Materially" means by more than 8 KiB. Written question 5 is about why the task
cannot ask a sharper question than that.

```
undergo verify compiler/04-binary-size
```

Every slot is a **comparison**, never a byte count: an absolute size depends on
the platform, the toolchain and the linker, and would differ on every machine
that ran this.

## Questions to answer in writing

1. The unreferenced 512 KiB table. Say whether it costs anything, what decides
   that, and what single change would flip the answer.
2. The method reachable only through an interface. Say why the linker cannot drop
   it, and what that means for a package that defines many methods on a type used
   behind an interface.
3. `reflect` has a reputation for bloating binaries. Say what the actual
   mechanism is, and why this measurement is smaller than you might expect.
4. `fmt` is the largest by a wide margin. Say what it drags in, and why a program
   that only ever prints a string still pays for it.
5. Four of these slots use an 8 KiB threshold and one asks about a single byte.
   Given that the same snippet measures +24 bytes on Linux and exactly 0 on
   Windows, say why a strict comparison against the baseline is the wrong
   question — and which of the five genuinely is exactly zero everywhere.

6. Name two build flags that reduce binary size, and say precisely what each one
   gives up.
