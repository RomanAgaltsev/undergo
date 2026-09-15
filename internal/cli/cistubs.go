package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// CIStubs is gate 1: every task package compiles, and every machine-graded
// task also vets clean.
//
// Review and design tasks are built but never vetted. Their code is the thing
// under review, so a vet diagnostic is the answer — printing it in a public CI
// log spoils the drill for everyone. loupe made the same call for the same
// reason.
//
// A task pinned to another architecture or operating system is skipped rather
// than built, and the skip is reported. This is safe only because CI runs both
// architectures: ubuntu-latest is amd64 and builds every amd64-pinned task,
// macos-latest is arm64 and builds everything else. A task pinning a platform
// no runner has would be built nowhere, and nobody would be told.
func CIStubs(e Env, _ []string) error {
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}
	build, vet, skipped := partitionPackages(tasks, manifest.CurrentEnv(), func(pkg string) bool {
		return hasGoFiles(filepath.Join(e.Root, filepath.FromSlash(pkg)))
	})

	if len(build) > 0 {
		if err := runGo(e, append([]string{"build"}, build...)); err != nil {
			return fmt.Errorf("gate 1: a task stub does not compile: %w", err)
		}
	}
	if len(vet) > 0 {
		if err := runGo(e, append([]string{"vet"}, vet...)); err != nil {
			return fmt.Errorf("gate 1: a task stub does not vet clean: %w", err)
		}
	}
	for _, s := range skipped {
		fmt.Fprintf(e.Out, "skip  %s\n", s)
	}
	// Report what was skipped on the same line as what was done. A gate that
	// skips silently can report green over an empty run — gate 2 once claimed
	// 174 reference solutions proven having run three of them.
	fmt.Fprintf(e.Out, "gate 1: %d packages built, %d vetted, %d not buildable on %s/%s\n",
		len(build), len(vet), len(skipped), runtime.GOOS, runtime.GOARCH)
	return nil
}

// partitionPackages decides what to build, what to vet, and what this host
// cannot build at all. hasGo reports whether a package path holds any Go files;
// prose-only tasks have none.
func partitionPackages(tasks []*manifest.Task, env manifest.Env, hasGo func(string) bool) (build, vet, skipped []string) {
	for _, t := range tasks {
		pkg := "./tasks/" + t.ID
		if !hasGo(pkg) {
			continue
		}
		if ok, why := manifest.Buildable(t, env); !ok {
			skipped = append(skipped, fmt.Sprintf("%s (%s)", t.ID, why))
			continue
		}
		build = append(build, pkg)
		if t.Mode == manifest.ModeReview || t.Mode == manifest.ModeDesign {
			continue
		}
		vet = append(vet, pkg)
	}
	return build, vet, skipped
}

func hasGoFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".go" {
			return true
		}
	}
	return false
}

func runGo(e Env, args []string) error {
	cmd := exec.Command("go", args...)
	cmd.Dir = e.Root
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err
	return cmd.Run()
}
