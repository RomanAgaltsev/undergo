// Package exitcode is about a test binary that stops in the middle, and the
// difference between a suite that passed and a suite that finished.
package exitcode

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

// Run runs the suite in testdata/suite and reports its combined output and
// whether `go test` reported success.
//
// A non-zero exit from `go test` is an answer, not an error: the error return
// is reserved for not being able to run it at all.
func Run() (output string, passed bool, err error) {
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = filepath.Join("testdata", "suite")
	out, runErr := cmd.CombinedOutput()

	var exit *exec.ExitError
	if runErr != nil && !errors.As(runErr, &exit) {
		return "", false, runErr
	}
	return strings.ReplaceAll(string(out), "\r\n", "\n"), runErr == nil, nil
}
