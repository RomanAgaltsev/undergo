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
