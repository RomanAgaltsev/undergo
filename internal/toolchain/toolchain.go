// Package toolchain builds and runs a small program under a named Go
// toolchain, fetching it on demand.
//
// Go 1.21 made GOTOOLCHAIN a version selector: naming a toolchain the machine
// does not have makes the go command download it as an ordinary module. That
// is how this repository reaches behaviour no go line can select — a GODEBUG
// that has been removed, a package that did not exist, a compiler that made a
// different choice.
//
// The go line and the toolchain are different axes. internal/goline varies the
// go line and pins the toolchain; this package varies the toolchain.
package toolchain

import (
	"context"
	"errors"
	"fmt"
	"go/version"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Local names the installed toolchain — no download, no network.
const Local = "local"

// minToolchain is the first release that can be selected this way, because
// GOTOOLCHAIN is itself a Go 1.21 feature.
const minToolchain = "go1.21"

// Result is what one toolchain invocation produced.
type Result struct {
	Output  string // combined stdout and stderr, trimmed of trailing newline
	Exited0 bool   // the command completed with status 0
}

// Run writes src as a single-file module declaring goLine and runs it under
// the named toolchain.
//
// A non-zero exit — including a compile failure — is reported in the Result,
// not as an error. The error return is reserved for the harness failing: a
// temp directory that cannot be made, a toolchain that cannot be obtained, or
// a go line the named toolchain cannot accept.
func Run(name, goLine, src string, env ...string) (Result, error) {
	return invoke(name, goLine, src, env, "run", ".")
}

// Build is Run without running the program.
func Build(name, goLine, src string) (Result, error) {
	return invoke(name, goLine, src, nil, "build", "-o", os.DevNull, ".")
}

var (
	availableMu sync.Mutex
	availableBy = map[string]bool{}
)

// Available reports whether the named toolchain can be obtained here, probing
// once per process and remembering the answer.
//
// The probe is `go version`, which has to resolve the toolchain to answer —
// so a successful probe also warms the module cache for the run that follows.
func Available(name string) bool {
	if name == Local {
		return true
	}
	availableMu.Lock()
	defer availableMu.Unlock()
	if got, ok := availableBy[name]; ok {
		return got
	}
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "version")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN="+name)
	ok := cmd.Run() == nil
	availableBy[name] = ok
	return ok
}

// Timeout bounds a toolchain invocation.
//
// These calls are download-bound rather than compute-bound: the first mention of
// a toolchain fetches it as an ordinary module, so the budget has to cover a
// cold module cache on a slow link. Without a deadline a stalled proxy hangs
// gate 2 indefinitely and reports nothing at all.
const Timeout = 10 * time.Minute

func invoke(name, goLine, src string, env []string, args ...string) (Result, error) {
	if err := checkToolchain(name, goLine); err != nil {
		return Result{}, err
	}
	if !Available(name) {
		return Result{}, fmt.Errorf("toolchain: %s cannot be obtained here; "+
			"it is fetched on demand and needs the network once", name)
	}

	dir, err := os.MkdirTemp("", "undergo-toolchain-")
	if err != nil {
		return Result{}, err
	}
	defer func() { _ = os.RemoveAll(dir) }()

	mod := fmt.Sprintf("module undergoprobe\n\ngo %s\n", goLine)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o644); err != nil {
		return Result{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(src), 0o644); err != nil {
		return Result{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	cmd.Env = append(append(os.Environ(), "GOTOOLCHAIN="+name), env...)

	out, err := cmd.CombinedOutput()
	res := Result{Output: strings.TrimRight(string(out), "\r\n"), Exited0: err == nil}

	// A killed process looks like a program that exited non-zero, which here
	// would be reported as the measured answer rather than as a harness failure.
	if ctx.Err() != nil {
		return Result{}, fmt.Errorf("running go %s under %s: timed out after %s",
			strings.Join(args, " "), name, Timeout)
	}
	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		return Result{}, fmt.Errorf("running go %s under %s: %w", strings.Join(args, " "), name, err)
	}
	return res, nil
}

// checkToolchain refuses a toolchain this mechanism cannot select, and a go
// line the named toolchain could not compile.
func checkToolchain(name, goLine string) error {
	want := "go" + goLine
	if !version.IsValid(want) {
		return fmt.Errorf("toolchain: %q is not a valid go line", goLine)
	}
	if name == Local {
		return nil
	}
	if !version.IsValid(name) {
		return fmt.Errorf("toolchain: %q is not a valid toolchain name", name)
	}
	if version.Compare(version.Lang(name), minToolchain) < 0 {
		return fmt.Errorf("toolchain: %s predates GOTOOLCHAIN, which arrived in %s",
			name, minToolchain)
	}
	if version.Compare(want, version.Lang(name)) > 0 {
		return fmt.Errorf("toolchain: go line %s is above %s; a module may not "+
			"declare a language version its compiler does not have", goLine, name)
	}
	return nil
}
