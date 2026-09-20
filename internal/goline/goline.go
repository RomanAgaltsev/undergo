// Package goline builds and runs a small program at a chosen go.mod language
// version, using the toolchain that is already installed.
//
// Go 1.21 made the go line the compatibility mechanism: GODEBUG defaults,
// several compiler semantics and language legality are all resolved from the
// main module's go directive. So one toolchain can produce two different
// behaviours from one source file, which is what the versions track grades.
package goline

import (
	"errors"
	"fmt"
	"go/version"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Result is what one toolchain invocation produced.
type Result struct {
	Output  string // combined stdout and stderr, trimmed of trailing newline
	Exited0 bool   // the command completed with status 0
}

// Run writes src as a single-file module declaring the given go line, runs it,
// and returns its combined output.
//
// Extra env entries are appended to the child's environment as "K=V".
//
// A non-zero exit — including a compile failure — is reported in the Result,
// not as an error. The error return is reserved for a failure of the harness
// itself: a temp directory that cannot be made, or a missing go command.
func Run(line, src string, env ...string) (Result, error) {
	return invoke(line, src, env, "run", ".")
}

// Build writes src as a single-file module declaring the given go line and
// compiles it without running it. On failure Output holds the compiler's
// diagnostic.
func Build(line, src string) (Result, error) {
	return invoke(line, src, nil, "build", "-o", os.DevNull, ".")
}

func invoke(line, src string, env []string, args ...string) (Result, error) {
	if err := checkCeiling(line); err != nil {
		return Result{}, err
	}

	dir, err := os.MkdirTemp("", "undergo-goline-")
	if err != nil {
		return Result{}, err
	}
	defer func() { _ = os.RemoveAll(dir) }()

	mod := fmt.Sprintf("module undergoprobe\n\ngo %s\n", line)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o644); err != nil {
		return Result{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(src), 0o644); err != nil {
		return Result{}, err
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	// GOTOOLCHAIN=local is load-bearing: without it a go line above the
	// installed toolchain makes the go command fetch one, which would break
	// the offline-by-default promise on a solver's machine.
	cmd.Env = append(append(os.Environ(), "GOTOOLCHAIN=local"), env...)

	out, err := cmd.CombinedOutput()
	res := Result{Output: strings.TrimRight(string(out), "\r\n"), Exited0: err == nil}

	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		return Result{}, fmt.Errorf("running go %s: %w", strings.Join(args, " "), err)
	}
	return res, nil
}

// checkCeiling refuses a go line newer than the toolchain running this code.
// Such a line is not a version diff; it is a toolchain download waiting to
// happen, and GOTOOLCHAIN=local would turn it into a confusing failure inside
// a task rather than a clear one here.
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
