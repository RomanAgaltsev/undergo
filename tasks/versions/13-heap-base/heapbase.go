// Package heapbase is about an address that used to be the same every time.
//
// Go's heap has long started at a recognisable place, which made a printed
// pointer look like a stable fact about a program. Go 1.26 randomizes the base
// on 64-bit platforms. The task is which part of the address that changed — and
// which part was never stable to begin with.
package heapbase

import (
	"fmt"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/toolchain"
)

// Program prints the address of one heap-allocated object.
const Program = `package main

import "fmt"

type T struct{ a, b int }

func main() {
	x := new(T)
	fmt.Printf("%p", x)
}
`

// Local is the toolchain running this test.
const Local = toolchain.Local

// BasePrefix is the prefix a Go heap pointer carried for years.
const BasePrefix = "0xc000"

// Address runs Program under the named toolchain and returns the pointer it
// printed. The go line is 1.21 in every run and does not vary.
func Address(name string) (string, error) {
	res, err := toolchain.Run(name, "1.21", Program)
	if err != nil {
		return "", err
	}
	if !res.Exited0 {
		return "", fmt.Errorf("%s: program did not run: %s", name, res.Output)
	}
	addr := strings.TrimSpace(res.Output)
	if !strings.HasPrefix(addr, "0x") {
		return "", fmt.Errorf("%s: unexpected output %q: not a pointer", name, addr)
	}
	return addr, nil
}
