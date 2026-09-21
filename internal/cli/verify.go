package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
	"github.com/RomanAgaltsev/undergo/internal/predict"
	"github.com/RomanAgaltsev/undergo/internal/progress"
)

// Verify grades the solver's work directory and records the outcome.
func Verify(e Env, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: undergo verify <id>")
	}
	t, err := find(e, args[0])
	if err != nil {
		return err
	}
	if ok, why := manifest.Gradeable(t, manifest.CurrentEnv()); !ok {
		fmt.Fprintf(e.Out, "SKIP %s: %s\n", t.ID, why)
		return nil
	}
	if t.Mode == manifest.ModeReview || t.Mode == manifest.ModeDesign {
		fmt.Fprintf(e.Out,
			"%s is a %s task: it is graded by you, against the sealed key.\n"+
				"Write your answer, then: undergo reveal %s\n", t.ID, t.Mode, t.ID)
		return nil
	}

	work := e.WorkDir(t.ID)
	if _, err := os.Stat(work); err != nil {
		return fmt.Errorf("no work directory — run: undergo start %s", t.ID)
	}

	passed, err := RunTests(e, work)
	if err != nil {
		return err
	}

	rec, err := progress.Load(e.ProgressPath())
	if err != nil {
		return err
	}
	rec.RecordRun(t.ID, passed)
	if err := rec.Save(e.ProgressPath()); err != nil {
		return err
	}
	if !passed {
		return fmt.Errorf("%s: not yet", t.ID)
	}
	fmt.Fprintf(e.Out, "PASS %s\n", t.ID)
	return nil
}

// TaskTimeout bounds one task's frozen tests.
//
// A frozen test can hang rather than fail. A synctest bubble waiting on a mutex
// never returns, which is a trap the concurrency track found by falling into it,
// and gate 2 runs 96 packages one after another. go test's own -timeout is per
// package and defaults to ten minutes, so one task that never returns could hold
// a runner for hours before anybody saw a diagnostic.
//
// The slowest task measured is well under a minute. This is a wide margin, not a
// budget: it exists to turn a hang into a named failure, not to police speed.
const TaskTimeout = 5 * time.Minute

// RunTests runs `go test` in dir with the task's environment, streaming output.
//
// Every mode is verified the same way: by running the frozen tests in the work
// directory. The parameter here used to be an unused *manifest.Task, kept in the
// signature for task.yaml's `verify` field — a field no code ever read, on all
// 267 manifests, for sixteen milestones. The field is gone and so is the
// parameter; dir is what this needs.
func RunTests(e Env, dir string) (bool, error) {
	args := []string{"test", "-count=1", "-v"}
	if e.Race {
		args = append(args, "-race")
	}
	args = append(args, ".")

	ctx, cancel := context.WithTimeout(context.Background(), TaskTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err
	cmd.Env = append(os.Environ(),
		"UNDERGO_PREDICTION="+filepath.Join(dir, predict.DefaultFile))

	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	// A killed process exits non-zero like a failing one, so the deadline is
	// checked first: "did not pass" and "never returned" are different answers
	// and only one of them is about the solution.
	if ctx.Err() != nil {
		return false, fmt.Errorf("timed out after %s — the frozen tests are not returning", TaskTimeout)
	}
	// A non-zero exit is a failed task, not a broken harness.
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return false, nil
	}
	return false, err
}
