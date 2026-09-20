// Package compiles asks a question with no answer until you name a version:
// is this valid Go?
package compiles

import "github.com/RomanAgaltsev/undergo/internal/goline"

// The five snippets. Each is compiled, never run.
const (
	ClearMinMax = `package main

import "fmt"

func main() {
	m := map[string]int{"a": 1}
	clear(m)
	fmt.Println(min(3, 1, 2), max(3, 1, 2), len(m))
}
`

	RangeInt = `package main

func main() {
	for range 3 {
	}
}
`

	RangeFunc = `package main

func seq(yield func(int) bool) {
	for i := 0; i < 3; i++ {
		if !yield(i) {
			return
		}
	}
}

func main() {
	for range seq {
	}
}
`

	NewExpr = `package main

func f(x int) int { return x }

var _ = new(f(1))

func main() {}
`

	GenericAlias = `package main

type Pair[A any] struct{ V A }

type Alias[A any] = Pair[A]

func main() { _ = Alias[int]{V: 1} }
`
)

// Compiles reports whether src compiles at the given go line, and returns the
// compiler's diagnostic when it does not.
func Compiles(line, src string) (ok bool, diagnostic string, err error) {
	res, err := goline.Build(line, src)
	if err != nil {
		return false, "", err
	}
	return res.Exited0, res.Output, nil
}
