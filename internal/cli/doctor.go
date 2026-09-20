package cli

import (
	"flag"
	"fmt"
	"maps"
	"os/exec"
	"slices"

	"github.com/RomanAgaltsev/undergo/internal/manifest"
)

// Doctor reports what this machine can grade.
func Doctor(e Env, args []string) error {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(e.Err)
	if err := fs.Parse(args); err != nil {
		return err
	}

	tasks, err := manifest.Walk(e.TasksDir())
	if err != nil {
		return err
	}

	env := manifest.CurrentEnv()
	fmt.Fprintf(e.Out, "go        %s\nplatform  %s/%s\nrace      %s\n",
		env.GoVersion, env.GOOS, env.GOARCH, raceSupport())
	reportToolchains(e, tasks)
	fmt.Fprintln(e.Out)

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

// reportToolchains says, for every toolchain the catalogue asks for, whether
// this machine can obtain it.
//
// A toolchain is no longer something you install and this command finds on
// PATH: it is named in GOTOOLCHAIN and fetched on demand, so the only useful
// question is whether it resolves here. Answering it may cost a download the
// first time, which also warms the cache for the run that follows.
func reportToolchains(e Env, tasks []*manifest.Task) {
	seen := map[string]bool{}
	for _, t := range tasks {
		for _, name := range t.Requires.Toolchains {
			seen[name] = true
		}
	}
	if len(seen) == 0 {
		fmt.Fprintln(e.Out, "toolchains no task requires one")
		return
	}
	wanted := slices.Sorted(maps.Keys(seen))
	for _, name := range wanted {
		status := "resolves"
		if !manifest.ToolchainAvailable(name) {
			status = "cannot be obtained here — it is fetched on demand and needs the network once"
		}
		fmt.Fprintf(e.Out, "toolchain %s  %s\n", name, status)
	}
}
