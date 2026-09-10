package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// CIStubs is gate 1: every task package compiles, and every machine-graded
// task also vets clean.
//
// Review and design tasks are built but never vetted. Their code is the thing
// under review, so a vet diagnostic is the answer — printing it in a public CI
// log spoils the drill for everyone. loupe made the same call for the same
// reason.
func CIStubs(e Env, _ []string) error {
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}
	build, vet := partitionPackages(tasks, func(pkg string) bool {
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
	fmt.Fprintf(e.Out, "gate 1: %d packages built, %d vetted\n", len(build), len(vet))
	return nil
}

// partitionPackages decides what to build and what to vet. hasGo reports
// whether a package path holds any Go files; prose-only tasks have none.
func partitionPackages(tasks []*manifest.Task, hasGo func(string) bool) (build, vet []string) {
	for _, t := range tasks {
		pkg := "./tasks/" + t.ID
		if !hasGo(pkg) {
			continue
		}
		build = append(build, pkg)
		if t.Mode == manifest.ModeReview || t.Mode == manifest.ModeDesign {
			continue
		}
		vet = append(vet, pkg)
	}
	return build, vet
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
