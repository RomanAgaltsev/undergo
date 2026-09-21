// Package promoted asks which struct-literal keys are legal, and which field a
// legal one reaches.
//
// Go 1.27 allows a key to be any valid field selector for the struct type, so a
// literal may set a promoted field of an embedded struct directly. Legality is
// resolved from the main module's go line and not from the toolchain, so one
// compiler answers both halves of this task.
package promoted

import (
	"fmt"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/goline"
)

// types is the shared declaration block: two levels of embedding, and a field
// name that exists at both depths.
const types = `type Deep struct {
	Z int
	X int
}

type Mid struct {
	Deep
	W int
	X int
}

type Top struct {
	Mid
	Y int
}
`

// ProgramOnce sets a field promoted through one level of embedding.
const ProgramOnce = `package main

import "fmt"

` + types + `
func main() { fmt.Print(Top{W: 2}.W) }
`

// ProgramTwice sets a field promoted through two levels.
const ProgramTwice = `package main

import "fmt"

` + types + `
func main() { fmt.Print(Top{Z: 3}.Z) }
`

// ProgramShadow names a field that exists at depth 1 and at depth 2, and prints
// both so the caller can see which one the key reached.
const ProgramShadow = `package main

import "fmt"

` + types + `
func main() {
	t := Top{X: 7}
	fmt.Printf("%d|%d", t.Mid.X, t.Mid.Deep.X)
}
`

// Compiles reports whether src compiles at the given go line, and returns the
// compiler's diagnostic when it does not.
func Compiles(goLine, src string) (ok bool, diagnostic string, err error) {
	res, err := goline.Build(goLine, src)
	if err != nil {
		return false, "", err
	}
	return res.Exited0, res.Output, nil
}

// Shadowed runs ProgramShadow and returns the two X values it printed, the
// depth-1 field first.
func Shadowed(goLine string) (mid, deep string, err error) {
	res, err := goline.Run(goLine, ProgramShadow)
	if err != nil {
		return "", "", err
	}
	if !res.Exited0 {
		return "", "", fmt.Errorf("program did not run: %s", res.Output)
	}
	a, b, ok := strings.Cut(res.Output, "|")
	if !ok {
		return "", "", fmt.Errorf("unexpected output %q: no separator", res.Output)
	}
	return a, b, nil
}
