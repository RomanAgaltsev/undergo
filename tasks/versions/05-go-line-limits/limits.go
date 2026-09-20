// Package limits finds the edge of what the go.mod line can reach backwards to.
package limits

import "github.com/RomanAgaltsev/undergo/internal/goline"

// CapProgram prints the capacity of the channel time.After returns. Before Go
// 1.23 that channel was buffered; from 1.23 it is not.
const CapProgram = `package main

import (
	"fmt"
	"time"
)

func main() { fmt.Print(cap(time.After(time.Hour))) }
`

// StartupProgram announces itself from init and from main, so that a failure
// before either can be told apart from a failure inside them.
const StartupProgram = `package main

import "fmt"

func init() { fmt.Print("init ") }

func main() { fmt.Print("main") }
`

// AfterCap runs CapProgram at the given go line and returns the printed
// capacity as a string.
func AfterCap(line string) (string, error) {
	res, err := goline.Run(line, CapProgram)
	if err != nil {
		return "", err
	}
	return res.Output, nil
}

// Startup runs StartupProgram at the given go line with the given extra
// environment, and returns its combined output.
func Startup(line string, env ...string) (string, error) {
	res, err := goline.Run(line, StartupProgram, env...)
	if err != nil {
		return "", err
	}
	return res.Output, nil
}
