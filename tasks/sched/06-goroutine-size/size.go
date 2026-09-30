// Package goroutinesize asks what a goroutine costs while it waits.
//
// Every measurement runs in a fresh child process, testdata/park, because the
// answer depends on the history of the process it is taken in: the runtime
// reuses freed stacks, and since Go 1.19 it sizes a new goroutine's stack from
// the goroutines it has already seen.
package goroutinesize

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// childEnv is the environment for every child process. GODEBUG is dropped
// because several of its settings change how stacks are sized. GOFLAGS is set
// to one of its own defaults rather than emptied: an empty GOFLAGS falls back
// to whatever `go env -w` stored, and a solver's -gcflags would reach the child.
func childEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		switch strings.ToUpper(name) {
		case "GODEBUG", "GOFLAGS":
			continue
		}
		env = append(env, kv)
	}
	return append(env, "GOFLAGS=-buildvcs=auto")
}

// Measure runs one phase of testdata/park in a fresh process and returns the
// number it prints.
func Measure(phase string) (int, error) {
	cmd := exec.Command("go", "run", "./testdata/park", phase)
	cmd.Env = childEnv()
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("phase %s: %w", phase, err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, fmt.Errorf("phase %s: output is not a number", phase)
	}
	return n, nil
}
