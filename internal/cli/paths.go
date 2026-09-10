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
