package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
	"github.com/RomanAgaltsev/undergo/internal/seal"
)

// SolutionSubdir is where a sealed blob keeps the files that overlay a task.
const SolutionSubdir = "solution"

// CIVerify is gate 2: every sealed reference solution must unseal, compile and
// pass its own task's frozen tests. Test output is captured and only a verdict
// is printed, so a solution never reaches a public CI log.
func CIVerify(e Env, _ []string) error {
	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}

	var failed []string
	var proven, skipped int
	for _, t := range tasks {
		if t.Mode == manifest.ModeReview || t.Mode == manifest.ModeDesign {
			fmt.Fprintf(e.Out, "skip  %s (%s is graded by a person)\n", t.ID, t.Mode)
			skipped++
			continue
		}
		if ok, why := manifest.Gradeable(t, manifest.CurrentEnv()); !ok {
			fmt.Fprintf(e.Out, "skip  %s: %s\n", t.ID, why)
			skipped++
			continue
		}
		if err := proveOne(e, t); err != nil {
			failed = append(failed, t.ID)
			fmt.Fprintf(e.Out, "FAIL  %s: %v\n", t.ID, err)
			continue
		}
		proven++
		fmt.Fprintf(e.Out, "ok    %s\n", t.ID)
	}
	if len(failed) > 0 {
		return fmt.Errorf("%d reference solution(s) did not pass: %v", len(failed), failed)
	}
	// Count what was actually proven, not what was looked at: a summary that
	// reports the whole catalogue when it verified three of it is the kind of
	// green that hides an empty run.
	fmt.Fprintf(e.Out, "\n%d reference solutions proven, %d skipped\n", proven, skipped)
	return nil
}

// proveOne unseals a task's solution over a fresh copy and runs its tests.
func proveOne(e Env, t *manifest.Task) error {
	work := e.WorkDir(filepath.Join("_ci", t.ID))
	if err := os.RemoveAll(work); err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(work) }()

	if err := copyTask(t.Dir, work); err != nil {
		return err
	}
	blob, err := os.ReadFile(filepath.Join(t.Dir, SealName))
	if err != nil {
		return fmt.Errorf("no sealed solution")
	}
	staging, err := os.MkdirTemp("", "undergo-seal-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(staging) }()
	if err := seal.Extract(string(blob), staging); err != nil {
		return err
	}

	// A predict task has nothing to overlay: the stub is the code, and the
	// blob carries only the answer sheet and the prose. Every other mode
	// must ship a solution/ that replaces the stub.
	src := filepath.Join(staging, SolutionSubdir)
	switch _, statErr := os.Stat(src); {
	case statErr == nil:
		if err := copyTask(src, work); err != nil {
			return err
		}
	case t.Mode != manifest.ModePredict:
		return fmt.Errorf("blob has no %s/ to overlay", SolutionSubdir)
	}
	if t.Mode == manifest.ModePredict {
		if err := writeReferencePredictions(staging, work); err != nil {
			return err
		}
	}

	// Capture everything: a reference solution must never reach a CI log.
	var sink bytes.Buffer
	quiet := Env{Root: e.Root, Out: &sink, Err: &sink}
	passed, err := RunTests(quiet, t, work)
	if err != nil {
		return err
	}
	if !passed {
		return fmt.Errorf("frozen tests did not pass")
	}
	return nil
}

// writeReferencePredictions installs the sealed answer sheet for a predict task.
func writeReferencePredictions(staging, work string) error {
	b, err := os.ReadFile(filepath.Join(staging, "prediction.yaml"))
	if err != nil {
		return fmt.Errorf("a predict task's blob must carry prediction.yaml: %w", err)
	}
	return os.WriteFile(filepath.Join(work, "prediction.yaml"), b, 0o644)
}
