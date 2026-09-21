package cli

import (
	"bytes"
	"context"
	"errors"
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
	requireToolchains := fs.Bool("require-toolchains", false,
		"fail rather than skip when a task's required toolchain cannot be obtained")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("usage: undergo ci-verify [--race] [--require-toolchains]: %w", ErrUsage)
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
			// A skip nobody is told about is indistinguishable from success.
			// In CI these tasks are proven or the gate is red — but only when
			// the reason is a toolchain. Declaring requires.toolchains is not
			// the same as being skipped for one: a task may also be pinned to
			// an architecture, and failing that here would send whoever reads
			// the log hunting for a fetch that never failed.
			if *requireToolchains && unobtainableToolchain(t) != "" {
				failed = append(failed, fmt.Sprintf("%s: %s", t.ID, why))
				continue
			}
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

	// Vet the reference solution before running it.
	//
	// Gate 1 vets the stub; nothing vetted the answer. A sealed blob is opaque
	// to every tool by design, and .golangci.yml excludes the work tree, so no
	// reference solution in this repository had ever been read by a linter or by
	// vet — 96 solutions over sixteen milestones, held up as exemplary Go.
	//
	// Output goes to the same sink for the same reason the tests do: a vet
	// diagnostic on a solution can quote the solution.
	var sink bytes.Buffer
	if err := vetSolution(work, &sink); err != nil {
		if path, werr := writeFailureLog(e, t.ID, sink.Bytes()); werr == nil {
			return fmt.Errorf("%w; output written to %s", err, path)
		}
		return err
	}

	// Capture everything: a reference solution must never reach a CI log.
	quiet := Env{Root: e.Root, Out: &sink, Err: &sink, Race: e.Race, CI: true}
	passed, err := RunTests(quiet, work)
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

// vetSolution runs go vet over the unsealed reference solution, writing any
// diagnostic to sink rather than to the log.
func vetSolution(work string, sink *bytes.Buffer) error {
	ctx, cancel := context.WithTimeout(context.Background(), TaskTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "vet", ".")
	cmd.Dir = work
	cmd.Stdout, cmd.Stderr = sink, sink
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("go vet timed out after %s", TaskTimeout)
		}
		return errors.New("the reference solution does not vet clean")
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
	for line := range strings.SplitSeq(output, "\n") {
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

// unobtainableToolchain returns the first toolchain the task requires that
// cannot be obtained here, or "" when every one of them can — including when
// the task requires none at all.
func unobtainableToolchain(t *manifest.Task) string {
	for _, want := range t.Requires.Toolchains {
		if !manifest.ToolchainAvailable(want) {
			return want
		}
	}
	return ""
}
