// Package negzero is about a value that is zero by one test and not by another.
//
// IEEE-754 has two zeros. The negative one compares equal to the positive one
// with ==, and differs from it in every bit but none of its value. Go 1.22
// changed what reflect.Value.IsZero says about it — and the interesting part is
// what Go 1.21 said, which was several different things at once, depending on
// how the value was reached and how big it was.
package negzero

import (
	"fmt"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/toolchain"
)

// Program prints five IsZero answers about one negative zero, reached five
// ways: bare, inside a small struct, inside a struct too large for reflect's
// comparable fast path, inside an array, and as the small struct's field.
const Program = `package main

import (
	"fmt"
	"math"
	"reflect"
)

type Small struct{ F float64 }

type Big struct {
	F   float64
	Pad [1024]byte
}

func main() {
	nz := math.Copysign(0, -1)
	small := Small{F: nz}
	big := Big{F: nz}
	fmt.Printf("%v|%v|%v|%v|%v",
		reflect.ValueOf(nz).IsZero(),
		reflect.ValueOf(small).IsZero(),
		reflect.ValueOf(big).IsZero(),
		reflect.ValueOf([2]float64{nz, nz}).IsZero(),
		reflect.ValueOf(small).Field(0).IsZero())
}
`

// Local is the toolchain running this test.
const Local = toolchain.Local

// Reading is what one toolchain reported, in the order Program prints.
type Reading struct {
	Scalar      string
	SmallStruct string
	BigStruct   string
	Array       string
	Field       string
}

// Read runs Program under the named toolchain. The go line is 1.21 in every run
// and does not vary; only the compiler does.
func Read(name string) (Reading, error) {
	res, err := toolchain.Run(name, "1.21", Program)
	if err != nil {
		return Reading{}, err
	}
	if !res.Exited0 {
		return Reading{}, fmt.Errorf("%s: program did not run: %s", name, res.Output)
	}
	parts := strings.Split(res.Output, "|")
	if len(parts) != 5 {
		return Reading{}, fmt.Errorf("%s: unexpected output %q: want 5 fields", name, res.Output)
	}
	return Reading{
		Scalar:      parts[0],
		SmallStruct: parts[1],
		BigStruct:   parts[2],
		Array:       parts[3],
		Field:       parts[4],
	}, nil
}
