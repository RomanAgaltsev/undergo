package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"

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
func CIStubs(e Env, args []string) error {
	if err := noArgs("ci-stubs", args); err != nil {
		return err
	}
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}
	build, vet, skipped := partitionPackages(tasks, manifest.CurrentEnv(), func(pkg string) bool {
		return hasGoFiles(filepath.Join(e.Root, filepath.FromSlash(pkg)))
	})

	if err := runGoOver(e, "build", build); err != nil {
		return fmt.Errorf("gate 1: a task stub does not compile: %w", err)
	}
	if err := runGoOver(e, "vet", vet); err != nil {
		return fmt.Errorf("gate 1: a task stub does not vet clean: %w", err)
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

// buildBatch is how many package paths go to one invocation of the go tool.
//
// Windows caps a command line at 32,767 characters and a task path averages
// about 45, so a single invocation ceilings out near 700 packages. The
// catalogue is at 231 and grows every milestone — and the cliff would appear on
// the maintainer's own machine first, while CI stayed green, because Linux and
// macOS have far more headroom. Chunking removes the ceiling; the build cache
// makes the extra invocations nearly free.
const buildBatch = 200

// runGoOver runs one go subcommand over pkgs, in batches.
func runGoOver(e Env, verb string, pkgs []string) error {
	for chunk := range slices.Chunk(pkgs, buildBatch) {
		if err := runGo(e, append([]string{verb}, chunk...)); err != nil {
			return err
		}
	}
	return nil
}

func runGo(e Env, args []string) error {
	// Gate 1 hands the go tool every task package at once, so this is the
	// longest-running child process in the harness. Bounded for the same reason
	// as RunTests: a compiler that never returns should be a named failure
	// rather than a runner held until the job's own limit.
	ctx, cancel := context.WithTimeout(context.Background(), TaskTimeout*2)
	defer cancel()

	// The binary is the literal "go"; args are package paths this repository
	// built from its own catalogue, never anything a user typed.
	cmd := exec.CommandContext(ctx, "go", args...) //nolint:gosec // G204: fixed binary, catalogue-derived args
	cmd.Dir = e.Root
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err
	err := cmd.Run()
	if ctx.Err() != nil {
		return fmt.Errorf("timed out after %s", TaskTimeout*2)
	}
	return err
}
