package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ModulePath identifies an undergo checkout.
const ModulePath = "github.com/RomanAgaltsev/undergo"

// Env is everything a command needs that is not its own flags.
type Env struct {
	Root string
	Out  io.Writer
	Err  io.Writer

	// Race runs a task's frozen tests under the race detector.
	//
	// It lives here rather than in a parameter because it describes how this
	// invocation runs rather than what it runs, and because RunTests is
	// reached from two commands that both need to pass it through unchanged.
	//
	// The detector needs cgo, so a machine with no C toolchain cannot honour
	// this. `undergo doctor` reports whether this one can.
	Race bool
}

// TasksDir is the task catalogue.
func (e Env) TasksDir() string { return filepath.Join(e.Root, "tasks") }

// WorkDir is where the solver's copy of a task lives. Gitignored.
func (e Env) WorkDir(id string) string {
	return filepath.Join(e.Root, "work", filepath.FromSlash(id))
}

// TaskDir is the catalogue directory for a task.
func (e Env) TaskDir(id string) string {
	return filepath.Join(e.TasksDir(), filepath.FromSlash(id))
}

// ProgressPath is the solver's record. Gitignored.
func (e Env) ProgressPath() string { return filepath.Join(e.Root, ".undergo", "progress.yaml") }

// RevealDir is where a revealed solution is written. Gitignored.
func (e Env) RevealDir(id string) string {
	return filepath.Join(e.Root, ".undergo", "reveals", filepath.FromSlash(id))
}

// FailureLogPath is where gate 2 writes the test output of a reference solution
// that did not pass. Gitignored.
//
// Gate 2 captures a failing task's output rather than printing it, because a
// reference solution must never reach a public CI log. That makes a failure
// undiagnosable from the log alone: "frozen tests did not pass" says nothing
// about which slot was wrong or whether the package even compiled.
//
// Writing the output to a file under .undergo/ resolves both needs. Only the
// path is printed, so CI logs stay clean, while whoever is standing at the
// checkout can read what actually happened. .undergo/ is gitignored, so the
// file cannot be committed by accident.
func (e Env) FailureLogPath(id string) string {
	return filepath.Join(e.Root, ".undergo", "ci-failures", filepath.FromSlash(id)+".log")
}

// FindRoot walks up from start to the undergo module root.
func FindRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		b, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil && strings.Contains(string(b), ModulePath) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not inside an undergo checkout (looked for go.mod declaring %s)", ModulePath)
		}
		dir = parent
	}
}

// splitID separates the single positional task id from flag arguments, so that
// `undergo reveal <id> --stuck` works as well as `undergo reveal --stuck <id>`.
// Go's flag package stops parsing at the first non-flag argument, so without
// this the documented argument order is rejected. Only safe for commands whose
// flags are all boolean, which is every command that takes an id.
func splitID(args []string) (id string, flags []string) {
	for _, a := range args {
		if id == "" && !strings.HasPrefix(a, "-") {
			id = a
			continue
		}
		flags = append(flags, a)
	}
	return id, flags
}
