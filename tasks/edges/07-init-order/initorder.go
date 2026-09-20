// Package initorder is about everything that happens before main.
//
// The program in testdata/prog is three packages deep and announces every
// step of its own initialisation. Nothing in it is concurrent and nothing in
// it is timed: the order is specified by the language, not decided by the
// implementation.
package initorder

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// Order runs testdata/prog and returns the lines it printed before main ran,
// in order.
func Order() ([]string, error) {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = filepath.Join("testdata", "prog")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var lines []string
	for line := range strings.SplitSeq(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "main.main") {
			break
		}
		lines = append(lines, line)
	}
	return lines, nil
}
