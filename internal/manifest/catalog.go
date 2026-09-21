package manifest

import (
	"fmt"
	"go/version"
	"io/fs"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/toolchain"
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
	slices.SortFunc(out, func(a, b *Task) int { return strings.Compare(a.ID, b.ID) })
	return out, nil
}

// Env is the machine a task would be graded on.
type Env struct {
	GoVersion string
	GOOS      string
	GOARCH    string
}

// CurrentEnv describes this machine.
func CurrentEnv() Env {
	return Env{
		GoVersion: runtime.Version(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
	}
}

// Buildable reports whether a task's package can be compiled on env, and why
// not if it cannot.
//
// It is deliberately narrower than Gradeable. Only the architecture and the
// operating system decide whether a package physically builds: an amd64-only
// .s file leaves its Go declaration without a body on arm64, and the build
// fails with "missing function body".
//
// A task needing a newer Go or an extra toolchain still compiles with what is
// here. The first should fail loudly rather than be skipped, and the second
// affects only running, so neither belongs in this check.
func Buildable(t *Task, e Env) (bool, string) {
	if len(t.Requires.Arch) > 0 && !slices.Contains(t.Requires.Arch, e.GOARCH) {
		return false, fmt.Sprintf("needs %s, this is %s",
			strings.Join(t.Requires.Arch, "/"), e.GOARCH)
	}
	if len(t.Requires.OS) > 0 && !slices.Contains(t.Requires.OS, e.GOOS) {
		return false, fmt.Sprintf("needs %s, this is %s",
			strings.Join(t.Requires.OS, "/"), e.GOOS)
	}
	return true, ""
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
		if !ToolchainAvailable(want) {
			return false, fmt.Sprintf("needs the %s toolchain, which could not be "+
				"obtained here (it is fetched on demand and needs the network once)", want)
		}
	}
	return true, ""
}

// ToolchainAvailable is how Gradeable decides whether a required toolchain can
// be obtained. It is a variable so tests can answer without the network.
//
// Availability used to mean "found on PATH", scanned from a hardcoded list of
// six names. No task ever declared requires.toolchains, so that function only
// ever returned nothing — and it could not express go1.19, which is exactly the
// kind of version a toolchain-pair task wants. A toolchain is now obtained by
// naming it in GOTOOLCHAIN and letting the go command fetch it.
var ToolchainAvailable = toolchain.Available
