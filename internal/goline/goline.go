// Package goline builds and runs a small program at a chosen go.mod language
// version, using the toolchain that is already installed.
//
// Go 1.21 made the go line the compatibility mechanism: GODEBUG defaults,
// several compiler semantics and language legality are all resolved from the
// main module's go directive. So one toolchain can produce two different
// behaviours from one source file, which is what the versions track grades.
//
// This package pins the toolchain and varies the go line. internal/toolchain
// is the other axis, and does the running for both.
package goline

import (
	"fmt"
	"go/version"
	"runtime"

	"github.com/RomanAgaltsev/undergo/internal/toolchain"
)

// Result is what one toolchain invocation produced.
type Result = toolchain.Result

// Run writes src as a single-file module declaring the given go line, runs it
// with the installed toolchain, and returns its combined output.
//
// Extra env entries are appended to the child's environment as "K=V".
//
// A non-zero exit — including a compile failure — is reported in the Result,
// not as an error. The error return is reserved for a failure of the harness
// itself: a temp directory that cannot be made, or a missing go command.
func Run(line, src string, env ...string) (Result, error) {
	if err := checkCeiling(line); err != nil {
		return Result{}, err
	}
	return toolchain.Run(toolchain.Local, line, src, env...)
}

// Build writes src as a single-file module declaring the given go line and
// compiles it without running it. On failure Output holds the compiler's
// diagnostic.
func Build(line, src string) (Result, error) {
	if err := checkCeiling(line); err != nil {
		return Result{}, err
	}
	return toolchain.Build(toolchain.Local, line, src)
}

// checkCeiling refuses a go line newer than the toolchain running this code.
// Such a line is not a version diff; it is a toolchain download waiting to
// happen, and this package's whole promise is that it never fetches one.
func checkCeiling(line string) error {
	want := "go" + line
	if !version.IsValid(want) {
		return fmt.Errorf("goline: %q is not a valid go line", line)
	}
	have := version.Lang(runtime.Version())
	if version.Compare(want, have) > 0 {
		return fmt.Errorf("goline: go line %s is newer than this toolchain (%s); "+
			"a version-diff task may only look backwards", line, have)
	}
	return nil
}
