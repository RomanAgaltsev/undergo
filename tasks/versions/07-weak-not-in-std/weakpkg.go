// Package weakpkg asks a question the go line cannot answer: does this import
// exist?
//
// The weak package arrived in Go 1.24. A module may declare any go line it
// likes; the standard library it compiles against is the one inside the
// toolchain.
package weakpkg

import "github.com/RomanAgaltsev/undergo/internal/toolchain"

// Program uses weak.Make, which needs the weak package.
const Program = `package main

import (
	"fmt"
	"weak"
)

func main() {
	v := new(int)
	p := weak.Make(v)
	fmt.Print(p.Value() != nil)
}
`

// Local is the toolchain running this test.
const Local = toolchain.Local

// Compiles reports whether Program compiles under the named toolchain, and
// returns the compiler's diagnostic when it does not.
func Compiles(name, goLine string) (ok bool, diagnostic string, err error) {
	res, err := toolchain.Build(name, goLine, Program)
	if err != nil {
		return false, "", err
	}
	return res.Exited0, res.Output, nil
}
