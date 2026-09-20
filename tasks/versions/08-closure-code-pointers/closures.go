// Package closures is about what a function value's pointer identifies.
//
// Go forbids comparing func values with ==, so the usual workaround is
// reflect.ValueOf(f).Pointer(). That returns the code pointer, which is not
// the same thing as the function value — and which compilers are free to
// share between closures.
package closures

import (
	"fmt"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/toolchain"
)

// Program prints two comparisons: two closures made from the SAME func
// literal over different captured values, and two SEPARATE literals with
// identical bodies.
const Program = `package main

import (
	"fmt"
	"reflect"
)

func mk(i int) func() int { return func() int { return i } }

func main() {
	a, b := mk(1), mk(2)
	f := func() int { return 42 }
	g := func() int { return 42 }
	fmt.Printf("%v|%v",
		reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer(),
		reflect.ValueOf(f).Pointer() == reflect.ValueOf(g).Pointer())
}
`

// Local is the toolchain running this test.
const Local = toolchain.Local

// Shares runs Program under the named toolchain and reports the two
// comparisons, same-site first.
func Shares(name string) (sameSite, distinctLiterals string, err error) {
	res, err := toolchain.Run(name, "1.21", Program)
	if err != nil {
		return "", "", err
	}
	a, b, ok := strings.Cut(res.Output, "|")
	if !ok {
		return "", "", fmt.Errorf("unexpected output %q: no separator", res.Output)
	}
	return a, b, nil
}
