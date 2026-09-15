package cli

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
	"github.com/RomanAgaltsev/undergo/internal/seal"
)

// SolutionSubdir is where a sealed blob keeps the files that overlay a task.
const SolutionSubdir = "solution"

// preflightRace refuses the race gate when this machine cannot build with it.
//
// Without this, every single task fails with "frozen tests did not pass" —
// which is true of the exit code and false about the cause. A machine with no
// C toolchain would be told it has 22 broken tasks rather than one missing
// compiler, and the gate would be measuring the environment instead of the
// solutions.
func preflightRace() error {
	if out, err := exec.Command("go", "env", "CGO_ENABLED").Output(); err == nil {
		if strings.TrimSpace(string(out)) == "0" {
			return fmt.Errorf("ci-verify --race: CGO_ENABLED=0, and the race detector is built on cgo; " +
				"set CGO_ENABLED=1, or run `task race:docker`")
		}
	}

	for _, cc := range []string{"gcc", "clang"} {
		if _, err := exec.LookPath(cc); err == nil {
			return nil
		}
	}
	return fmt.Errorf("ci-verify --race: no C toolchain (gcc or clang) on PATH, and the race detector needs one; " +
		"run `task race:docker` instead, or `go run ./cmd/undergo doctor` to see what this machine can do")
}

// CIVerify is gate 2: every sealed reference solution must unseal, compile and
// pass its own task's frozen tests. Test output is captured and only a verdict
// is printed, so a solution never reaches a public CI log.
func CIVerify(e Env, args []string) error {
	fs := flag.NewFlagSet("ci-verify", flag.ContinueOnError)
	fs.SetOutput(e.Err)
	race := fs.Bool("race", false, "run every reference solution under the race detector")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("usage: undergo ci-verify [--race]")
	}
	e.Race = *race
	if e.Race {
		if err := preflightRace(); err != nil {
			return err
		}
	}

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
		// A task pinned to the default build is skipped by the race gate
		// rather than failed by it: -race changes the compiler's escape and
		// stack-allocation decisions, so racing such a task measures the
		// instrumentation. The plain gate still proves it.
		if e.Race && t.Requires.DefaultBuild {
			fmt.Fprintf(e.Out, "skip  %s (answers are pinned to the default build; -race changes it)\n", t.ID)
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
	quiet := Env{Root: e.Root, Out: &sink, Err: &sink, Race: e.Race}
	passed, err := RunTests(quiet, t, work)
	if err != nil {
		return err
	}
	if !passed {
		where := ""
		if path, werr := writeFailureLog(e, t.ID, sink.Bytes()); werr == nil {
			where = "; output written to " + path
		}
		if verdicts := slotVerdicts(t, sink.String()); verdicts != "" {
			return fmt.Errorf("frozen tests did not pass (%s)%s", verdicts, where)
		}
		return fmt.Errorf("frozen tests did not pass%s", where)
	}
	return nil
}

// slotVerdicts names the slots a failing predict task got wrong, so a CI log
// says which question failed rather than only that something did.
//
// This is safe to print, and only for a predict task. predict.Check emits one
// line per slot of the form `slot "name": correct` — the slot names are already
// public in task.yaml, and the package deliberately never prints a measured
// value, because doing so would hand over the answer. Nothing else from the
// captured output is included: a build or optimize task's output can contain
// the overlaid reference source, which must never reach a public log.
func slotVerdicts(t *manifest.Task, output string) string {
	if t.Mode != manifest.ModePredict {
		return ""
	}
	var wrong []string
	for _, line := range strings.Split(output, "\n") {
		text := strings.TrimSpace(line)
		i := strings.Index(text, `slot "`)
		if i < 0 || strings.HasSuffix(text, ": correct") {
			continue
		}
		wrong = append(wrong, text[i:])
	}
	if len(wrong) == 0 {
		return ""
	}
	return strings.Join(wrong, "; ")
}

// writeFailureLog saves a failing task's captured output where a maintainer can
// read it, and returns the path.
//
// The path is safe to print; the contents are not. See Env.FailureLogPath.
func writeFailureLog(e Env, id string, output []byte) (string, error) {
	path := e.FailureLogPath(id)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, output, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// writeReferencePredictions installs the sealed answer sheet for a predict task.
func writeReferencePredictions(staging, work string) error {
	b, err := os.ReadFile(filepath.Join(staging, "prediction.yaml"))
	if err != nil {
		return fmt.Errorf("a predict task's blob must carry prediction.yaml: %w", err)
	}
	return os.WriteFile(filepath.Join(work, "prediction.yaml"), b, 0o644)
}
