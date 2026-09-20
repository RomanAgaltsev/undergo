package goline

import (
	"strings"
	"testing"
)

const helloSrc = `package main

import "fmt"

func main() { fmt.Print("hello") }
`

func TestRunExecutesAtAGoLine(t *testing.T) {
	got, err := Run("1.21", helloSrc)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !got.Exited0 {
		t.Fatalf("Exited0 = false; output:\n%s", got.Output)
	}
	if got.Output != "hello" {
		t.Errorf("Output = %q; want %q", got.Output, "hello")
	}
}

// The go line must actually reach the compiler: `for range 3` is a compile
// error before 1.22 and legal from 1.22.
func TestGoLineReachesTheCompiler(t *testing.T) {
	const src = `package main

func main() {
	for range 3 {
	}
}
`
	old, err := Build("1.21", src)
	if err != nil {
		t.Fatalf("Build(1.21): %v", err)
	}
	if old.Exited0 {
		t.Fatal("expected `for range 3` to be rejected at go 1.21")
	}
	if !strings.Contains(old.Output, "requires go1.22 or later") {
		t.Errorf("diagnostic did not name the requirement:\n%s", old.Output)
	}

	now, err := Build("1.22", src)
	if err != nil {
		t.Fatalf("Build(1.22): %v", err)
	}
	if !now.Exited0 {
		t.Errorf("expected `for range 3` to compile at go 1.22:\n%s", now.Output)
	}
}

// A failing program is an answer, not an error.
func TestNonZeroExitIsAResultNotAnError(t *testing.T) {
	const src = `package main

import "os"

func main() { os.Exit(3) }
`
	got, err := Run("1.21", src)
	if err != nil {
		t.Fatalf("Run returned an error for a non-zero exit: %v", err)
	}
	if got.Exited0 {
		t.Error("Exited0 = true for a program that exited 3")
	}
}

func TestEnvIsPassedThrough(t *testing.T) {
	const src = `package main

import (
	"fmt"
	"os"
)

func main() { fmt.Print(os.Getenv("UNDERGO_PROBE")) }
`
	got, err := Run("1.21", src, "UNDERGO_PROBE=set")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got.Output != "set" {
		t.Errorf("Output = %q; want %q", got.Output, "set")
	}
}

// The ceiling guard. A go line above the installed toolchain would make the go
// command download one, which breaks the offline promise silently.
func TestRefusesAGoLineAboveTheToolchain(t *testing.T) {
	if _, err := Run("99.0", helloSrc); err == nil {
		t.Fatal("Run accepted a go line above the toolchain")
	}
	if _, err := Build("99.0", helloSrc); err == nil {
		t.Fatal("Build accepted a go line above the toolchain")
	}
}
