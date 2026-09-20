// Package panicnil asks what recover() hands back when a program panics with
// nil, under two declared language versions and two GODEBUG settings.
package panicnil

import "github.com/RomanAgaltsev/undergo/internal/goline"

// Program panics with nil and prints the dynamic type of whatever recover
// returns. It is compiled by the toolchain you already have; only the go.mod
// line changes between runs.
const Program = `package main

import "fmt"

func main() {
	defer func() { fmt.Printf("%T", recover()) }()
	panic(nil)
}
`

// RecoveredType runs Program at the given go line with the given extra
// environment and reports the printed type.
func RecoveredType(line string, env ...string) (string, error) {
	res, err := goline.Run(line, Program, env...)
	if err != nil {
		return "", err
	}
	return res.Output, nil
}
