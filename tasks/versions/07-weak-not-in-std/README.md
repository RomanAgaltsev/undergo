# 07 — The package that did not exist

`versions/03-what-compiles` asked whether a snippet is valid Go, and the answer
turned on the `go` line. From four of its five snippets you could reasonably
conclude that the go line governs what compiles.

This program holds the go line **fixed** at `go 1.23` and compiles it twice
anyway:

```go
import (
	"fmt"
	"weak"
)

func main() {
	v := new(int)
	p := weak.Make(v)
	fmt.Print(p.Value() != nil)
}
```

| slot | toolchain | module declares |
|---|---|---|
| `compiles_on_local` | the one you have | `go 1.23` |
| `compiles_under_go1_23` | **go1.23.12** | `go 1.23` |

And one slot about the failure, if there is one:

| slot | question |
|---|---|
| `diagnostic_says_not_in_std` | does the compiler's message contain `is not in std`? |

Predict `true` or `false` for each.

```
undergo verify versions/07-weak-not-in-std
```

The judge is the compiler: the test writes the same source as a one-file module
twice and records whether each build succeeded.

**This task fetches a Go toolchain the first time it runs.** If it cannot be
fetched the task skips and says so; `undergo doctor` reports which resolve here.

## Questions to answer in writing

1. `versions/03-what-compiles` showed legality gated by the go line. This one
   holds the go line fixed and legality still changes. State the rule that
   covers both.
2. Say where the standard library physically lives, and what that means for a
   module that lowers its go line to support older users.
3. The diagnostic names a path. Say what that path tells you about how the
   toolchain was obtained.
4. You maintain a library that must build on the last three Go releases. Say
   what your go line should be, and what it does *not* protect you from.
5. `weak` has no replacement on an older toolchain. Say what a library does when
   it needs a package that may be absent, and name the mechanism.
