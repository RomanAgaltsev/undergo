package manifest

import (
	"fmt"
	"go/version"
	"io/fs"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
)

// Walk loads and validates every task.yaml under root, sorted by ID.
func Walk(root string) ([]*Task, error) {
	var out []*Task
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "task.yaml" {
			return nil
		}
		t, err := Load(p)
		if err != nil {
			return err
		}
		if err := Validate(t); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		out = append(out, t)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// Env is the machine a task would be graded on.
type Env struct {
	GoVersion  string
	GOOS       string
	GOARCH     string
	Toolchains []string
}

// CurrentEnv describes this machine.
func CurrentEnv() Env {
	return Env{
		GoVersion:  runtime.Version(),
		GOOS:       runtime.GOOS,
		GOARCH:     runtime.GOARCH,
		Toolchains: installedToolchains(),
	}
}

// installedToolchains reports alternate Go toolchains found on PATH, e.g. go1.21.13.
func installedToolchains() []string {
	var out []string
	for _, name := range []string{"go1.21.13", "go1.22.12", "go1.23.12", "go1.24.6", "go1.25.1", "go1.26.0"} {
		if _, err := exec.LookPath(name); err == nil {
			out = append(out, name)
		}
	}
	return out
}

// Gradeable reports whether a task can be graded in env, and why not if it cannot.
func Gradeable(t *Task, e Env) (bool, string) {
	if t.Requires.Go != "" {
		if version.Compare(e.GoVersion, "go"+t.Requires.Go) < 0 {
			return false, fmt.Sprintf("needs Go %s or newer, this is %s", t.Requires.Go, e.GoVersion)
		}
	}
	if t.Requires.MaxGo != "" {
		if version.Compare(e.GoVersion, "go"+t.Requires.MaxGo) > 0 {
			return false, fmt.Sprintf("needs Go %s or older, this is %s", t.Requires.MaxGo, e.GoVersion)
		}
	}
	if len(t.Requires.Arch) > 0 && !slices.Contains(t.Requires.Arch, e.GOARCH) {
		return false, fmt.Sprintf("answers are %s-specific, this is %s",
			strings.Join(t.Requires.Arch, "/"), e.GOARCH)
	}
	if len(t.Requires.OS) > 0 && !slices.Contains(t.Requires.OS, e.GOOS) {
		return false, fmt.Sprintf("needs %s, this is %s", strings.Join(t.Requires.OS, "/"), e.GOOS)
	}
	for _, want := range t.Requires.Toolchains {
		if !slices.Contains(e.Toolchains, want) {
			return false, fmt.Sprintf("needs the %s toolchain: go install golang.org/dl/%s@latest", want, want)
		}
	}
	return true, ""
}
