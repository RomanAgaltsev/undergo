// Package panicprint asks what the runtime is able to say about the value you
// panicked with.
//
// The message is printed by the runtime, after the stack has begun unwinding
// and without calling into the fmt package — so what it can say is limited by
// what the runtime knows how to format, not by what fmt could have managed.
package panicprint

import (
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/goline"
)

// Kinds are the four values panicked with, each in its own program.
const (
	NamedInt    = "Code(7)"
	NamedString = `Str("boom")`
	ErrorValue  = `MyErr{"broke"}`
	PlainStruct = "struct{ A int }{3}"
)

const program = `package main

type Code int

type Str string

type MyErr struct{ S string }

func (e MyErr) Error() string { return e.S }

func main() { panic(VALUE) }
`

// PanicLine runs a program that panics with the given expression and returns
// the runtime's "panic: ..." line.
func PanicLine(value string) (string, error) {
	res, err := goline.Run("1.27", strings.Replace(program, "VALUE", value, 1))
	if err != nil {
		return "", err
	}
	for line := range strings.SplitSeq(strings.ReplaceAll(res.Output, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(line, "panic: ") {
			return strings.TrimSpace(line), nil
		}
	}
	return "", nil
}
