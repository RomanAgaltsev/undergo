// Package preempt asks whether the runtime can stop a goroutine that never
// calls a function, allocates, or touches a channel.
package preempt

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunSpinner builds and runs testdata/spin under the given GODEBUG setting and
// reports whether the collection finished while the spin loop was still going.
func RunSpinner(asyncPreemptOff bool) (bool, error) {
	cmd := exec.Command("go", "run", "./testdata/spin")
	if asyncPreemptOff {
		cmd.Env = append(cmd.Environ(), "GODEBUG=asyncpreemptoff=1")
	}
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	switch text := strings.TrimSpace(string(out)); text {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("spin printed %q, expected true or false", text)
	}
}
