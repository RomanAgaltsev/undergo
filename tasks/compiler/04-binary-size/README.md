# 04 — What actually costs bytes

Six programs in `testdata/`, all built the same way. One is a baseline; the other
five each add exactly one thing to it. Predict which of them made the binary
bigger, and which added the most.

| slot | what the program adds |
|---|---|
| `fmt_grows_binary` | one `fmt.Println` |
| `reflect_grows_binary` | one `reflect.TypeOf` |
| `generic_grows_binary` | a generic function instantiated over four types |
| `iface_method_grows_binary` | a method reachable only through an interface |
| `unused_table_grows_binary` | a 512 KiB package-level array nothing references |
| `largest_addition` | the name of the snippet that added the most |

Four of the five grow it, or fewer. Two of these have loud reputations that the
measurement does not support.

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
5. Name two build flags that reduce binary size, and say precisely what each one
   gives up.
