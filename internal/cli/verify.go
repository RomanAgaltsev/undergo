package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

	passed, err := RunTests(e, t, work)
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

// RunTests runs `go test` in dir with the task's environment, streaming output.
//
// The task is taken but unused in v1: every mode is verified the same way, by
// running the frozen tests in the work directory. It stays in the signature
// because task.yaml carries an overridable verify command.
func RunTests(e Env, _ *manifest.Task, dir string) (bool, error) {
	cmd := exec.Command("go", "test", "-count=1", "-v", ".")
	cmd.Dir = dir
	cmd.Stdout = e.Out
	cmd.Stderr = e.Err
	cmd.Env = append(os.Environ(),
		"UNDERGO_PREDICTION="+filepath.Join(dir, predict.DefaultFile))

	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	// A non-zero exit is a failed task, not a broken harness.
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return false, nil
	}
	return false, err
}
