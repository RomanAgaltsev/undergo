package cli

import (
	"flag"
	"fmt"
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
	fmt.Fprintf(e.Out, "go        %s\nplatform  %s/%s\ntoolchains %s\n\n",
		env.GoVersion, env.GOOS, env.GOARCH, toolchainList(env.Toolchains))

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

func toolchainList(t []string) string {
	if len(t) == 0 {
		return "(none besides the default)"
	}
	return strings.Join(t, ", ")
}
