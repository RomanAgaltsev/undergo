# 05 — The method the linker cannot drop

`Widget` has four methods: `Alpha`, `Bravo`, `Charlie`, `Delta`. Two functions
call one of them:

```go
func CallDirect(w Widget) int              { return w.Alpha() }
func CallByName(w Widget, name string) int { ... MethodByName(name) ... }
```

The test builds **three** artifacts and counts how many of the four methods
survive as symbols in each:

| slot | artifact |
|---|---|
| `archive_retained` | the package compiled to an archive (`go build -o lib.a .`) |
| `direct_binary_retained` | a binary whose `main` calls only `CallDirect` |
| `byname_binary_retained` | a binary whose `main` calls `CallByName` |

Out of four. **One of these three numbers is the one you expected, one is
higher, and one is lower.**

```
undergo verify reflect/05-methodbyname-linker
```

## Questions to answer in writing

1. The direct binary retains **zero** methods, not one, even though
   `CallDirect` calls `Alpha`. Explain in two steps what happened to `Alpha`.
2. The archive retains all four regardless of what anything calls. Say why an
   archive cannot answer a question about reachability at all — then name the
   task in this repository that counts symbols **in an archive** precisely
   because a binary would lie, and say what makes the two questions different.
3. `MethodByName` takes a string that could come from anywhere. Say what the
   linker would have to prove in order to drop a method anyway, and why
   `-ldflags=-w` does not help.
4. This is the mechanism behind a class of production bug where a method works
   in development and is missing from a trimmed build. Name a framework style
   that triggers it, and say what such a framework must do to stay safe.
5. `countMethods` anchors its pattern with `$`. Find out what else `go tool nm`
   emits per symbol, and say what the count would have been without the anchor.
