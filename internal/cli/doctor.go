package cli

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// Doctor reports what this machine can grade.
func Doctor(e Env, args []string) error {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(e.Err)
	if err := fs.Parse(args); err != nil {
		return err
	}

	env := manifest.CurrentEnv()
	fmt.Fprintf(e.Out, "go        %s\nplatform  %s/%s\ntoolchains %s\nrace      %s\n\n",
		env.GoVersion, env.GOOS, env.GOARCH, toolchainList(nil), raceSupport())

	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}
	var blocked int
	for _, t := range tasks {
		if ok, why := manifest.Gradeable(t, env); !ok {
			blocked++
			fmt.Fprintf(e.Out, "skip  %s: %s\n", t.ID, why)
		}
	}
	fmt.Fprintf(e.Out, "\n%d of %d tasks gradeable here\n", len(tasks)-blocked, len(tasks))
	return nil
}

// raceSupport reports whether `undergo ci-verify --race` can run here.
//
// The race detector is built on cgo, so it needs a C toolchain — which is the
// usual reason a Windows checkout cannot run the gate locally. Docker is the
// documented way round it, so this says so rather than leaving the reader to
// work out why -race refused.
func raceSupport() string {
	if _, err := exec.LookPath("gcc"); err == nil {
		return "available"
	}
	if _, err := exec.LookPath("clang"); err == nil {
		return "available"
	}
	if _, err := exec.LookPath("docker"); err == nil {
		return "no C toolchain; run `task test:docker` or `task race:docker` instead"
	}
	return "no C toolchain and no docker; the race gate cannot run here — CI is authoritative"
}

func toolchainList(t []string) string {
	if len(t) == 0 {
		return "(none besides the default)"
	}
	return strings.Join(t, ", ")
}
