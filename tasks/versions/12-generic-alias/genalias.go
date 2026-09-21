// Package genalias asks what a parameterized type alias costs.
//
// An alias is meant to be another spelling of a type rather than a new type.
// Generic aliases arrived in Go 1.24; before that the same declaration needed
// an experiment. The question is whether the alias produces a second
// instantiation or shares the one it names.
package genalias

import (
	"fmt"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/toolchain"
)

// Program declares a generic type and a parameterized alias of it, then reports
// whether the two instantiated types are identical and what the alias's type
// descriptor calls itself.
const Program = `package main

import (
	"fmt"
	"reflect"
)

type Box[T any] struct{ V T }

type Alias[T any] = Box[T]

func main() {
	var a Alias[int]
	var b Box[int]
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	fmt.Printf("%v|%s", ta == tb, ta.String())
}
`

// Local is the toolchain running this test.
const Local = toolchain.Local

// goLine is fixed for every run in this task: the compiler is the only axis.
const goLine = "1.23"

// Compiles reports whether Program compiles under the named toolchain, and
// returns the compiler's diagnostic when it does not.
func Compiles(name string) (ok bool, diagnostic string, err error) {
	res, err := toolchain.Build(name, goLine, Program)
	if err != nil {
		return false, "", err
	}
	return res.Exited0, res.Output, nil
}

// Identity runs Program under the named toolchain and reports whether the two
// instantiations are the same type, and what the alias's type is called.
func Identity(name string) (identical, typeName string, err error) {
	res, err := toolchain.Run(name, goLine, Program)
	if err != nil {
		return "", "", err
	}
	if !res.Exited0 {
		return "", "", fmt.Errorf("%s: program did not run: %s", name, res.Output)
	}
	a, b, ok := strings.Cut(res.Output, "|")
	if !ok {
		return "", "", fmt.Errorf("%s: unexpected output %q: no separator", name, res.Output)
	}
	return a, b, nil
}
