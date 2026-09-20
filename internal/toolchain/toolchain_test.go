package toolchain

import (
	"strings"
	"testing"
)

const helloSrc = `package main

import "fmt"

func main() { fmt.Print("hello") }
`

// The whole point of the package: a named toolchain runs the program, and it
// is not the one running the tests.
func TestRunsUnderTheNamedToolchain(t *testing.T) {
	if !Available("go1.22.12") {
		t.Skip("go1.22.12 cannot be fetched here")
	}
	const src = `package main

import (
	"fmt"
	"runtime"
)

func main() { fmt.Print(runtime.Version()) }
`
	got, err := Run("go1.22.12", "1.19", src)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got.Output != "go1.22.12" {
		t.Errorf("ran under %q, want go1.22.12", got.Output)
	}
}

func TestLocalRunsUnderTheInstalledToolchain(t *testing.T) {
	got, err := Run(Local, "1.21", helloSrc)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got.Output != "hello" {
		t.Errorf("Output = %q; want hello", got.Output)
	}
}

// A module may not declare a language version its compiler does not have.
func TestRefusesAGoLineAboveTheNamedToolchain(t *testing.T) {
	if _, err := Run("go1.22.12", "1.24", helloSrc); err == nil {
		t.Fatal("accepted go 1.24 under go1.22.12")
	}
}

func TestRefusesAToolchainOlderThanGOTOOLCHAINItself(t *testing.T) {
	if _, err := Run("go1.13", "1.13", helloSrc); err == nil {
		t.Fatal("accepted a toolchain older than go1.21")
	}
}

// A compile failure is an answer, not a harness failure.
func TestCompileFailureIsAResult(t *testing.T) {
	const src = `package main

func main() { var x int }
`
	got, err := Build(Local, "1.21", src)
	if err != nil {
		t.Fatalf("Build returned an error for a compile failure: %v", err)
	}
	if got.Exited0 {
		t.Error("Exited0 = true for a program that does not compile")
	}
	if !strings.Contains(got.Output, "declared and not used") {
		t.Errorf("diagnostic missing:\n%s", got.Output)
	}
}

// An unobtainable toolchain is a typed error, never a panic or a hang.
func TestUnobtainableToolchainIsAnError(t *testing.T) {
	if Available("go1.99.99") {
		t.Fatal("Available said yes to a version that does not exist")
	}
	if _, err := Run("go1.99.99", "1.21", helloSrc); err == nil {
		t.Fatal("Run accepted a toolchain that cannot be fetched")
	}
}
