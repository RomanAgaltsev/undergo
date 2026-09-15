//go:build amd64

// Package abi asks where the register ABI puts arguments and results.
//
// Since Go 1.17 arguments travel in registers rather than on the stack. There
// is a fixed sequence of them, it is not very long, and integers and floats
// draw from two independent sequences.
//
// Nothing here needs running. Every answer is in the listing the compiler will
// print for you.
package abi

import (
	"os/exec"
	"strings"
)

// Nine takes nine integer arguments.
func Nine(a, b, c, d, e, f, g, h, i int64) int64 {
	return a + b + c + d + e + f + g + h + i
}

// Ten takes one more than Nine.
func Ten(a, b, c, d, e, f, g, h, i, j int64) int64 {
	return a + b + c + d + e + f + g + h + i + j
}

// Pair is a two-field struct passed by value.
type Pair struct{ X, Y int64 }

// SumPair takes the struct rather than its fields.
func SumPair(p Pair) int64 { return p.X + p.Y }

// Mixed interleaves integer and floating-point arguments.
func Mixed(a int64, x float64, b int64, y float64) float64 {
	return float64(a+b) + x + y
}

// Listing returns the compiler's assembly for this package.
func Listing() (string, error) {
	out, err := exec.Command("go", "build", "-gcflags=-S", ".").CombinedOutput()
	if err != nil && len(out) == 0 {
		return "", err
	}
	return string(out), nil
}

// bodyOf returns the listing lines belonging to one function.
//
// The listing names a symbol by its FULL IMPORT PATH, so the match is on
// "."+name+" STEXT" rather than on the package clause. Matching the short name
// finds nothing and every question then answers itself incorrectly.
func bodyOf(listing, name string) string {
	var b strings.Builder
	marker := "." + name + " STEXT"
	on := false
	for _, line := range strings.Split(listing, "\n") {
		switch {
		case strings.Contains(line, marker):
			on = true
		case on && strings.Contains(line, " STEXT"):
			on = false
		}
		if on {
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	return b.String()
}

// ReadsArgumentFromStack reports whether fn loads any argument from the stack
// rather than taking all of them in registers.
func ReadsArgumentFromStack(listing, fn string) bool {
	return strings.Contains(bodyOf(listing, fn), "(SP)")
}

// UsesFloatRegisters reports whether fn does floating-point arithmetic in the
// X register file.
func UsesFloatRegisters(listing, fn string) bool {
	return strings.Contains(bodyOf(listing, fn), "ADDSD")
}
