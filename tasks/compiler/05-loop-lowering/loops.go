// Package loops asks which hand-written loops the compiler recognises and
// replaces with a single runtime call, and which it leaves as loops.
//
// Nothing here is measured at run time. Every answer is in the listing.
package loops

import (
	"os/exec"
	"strings"
)

// ZeroLoop zeroes a byte slice the long way.
func ZeroLoop(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// ZeroClear zeroes it with the builtin.
func ZeroClear(b []byte) { clear(b) }

// CopyLoop copies element by element.
//
// staticcheck's S1001 says to use copy here, and it is right — that is the
// point of the task, and the suppression is why the loop survives to be
// compared against CopyBuiltin.
//
//nolint:staticcheck // the loop is the subject of the exercise
func CopyLoop(dst, src []byte) {
	for i := range src {
		dst[i] = src[i]
	}
}

// CopyBuiltin copies with the builtin.
func CopyBuiltin(dst, src []byte) { copy(dst, src) }

// ZeroPointers nils out a slice of pointers the long way. The elements hold
// pointers, which changes what the compiler is allowed to do.
func ZeroPointers(s []*int) {
	for i := range s {
		s[i] = nil
	}
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
// The listing names a symbol by its full import path, so the match is on
// "."+name+" STEXT".
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

// Calls reports whether fn's body mentions the named runtime symbol.
func Calls(listing, fn, symbol string) bool {
	return strings.Contains(bodyOf(listing, fn), symbol)
}
