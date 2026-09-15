package shapes

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var (
	profileOnce sync.Once
	profilePath string
	profileErr  error
)

// profile collects a CPU profile of this package's own benchmark and returns
// the file's path.
//
// It is generated rather than checked in, and that is not a convenience. A PGO
// profile records symbols by their FULL IMPORT PATH, so a profile collected at
// one path is silently ignored when the same code is compiled at another —
// which is exactly what happens when a task directory is copied elsewhere to be
// verified. Generating it here makes the task work wherever it is run from.
func profile() (string, error) {
	profileOnce.Do(func() {
		dir, err := os.MkdirTemp("", "undergo-pgo-")
		if err != nil {
			profileErr = err
			return
		}
		path := filepath.Join(dir, "cpu.pgo")

		cmd := exec.Command("go", "test",
			"-run", "^$",
			"-bench", "^BenchmarkProfileWorkload$",
			"-benchtime", "2s",
			"-cpuprofile", path,
			"-o", filepath.Join(dir, "bench.test"),
			".")
		if out, err := cmd.CombinedOutput(); err != nil {
			profileErr = &profileError{out: string(out), err: err}
			return
		}
		profilePath = path
	})
	return profilePath, profileErr
}

type profileError struct {
	out string
	err error
}

func (e *profileError) Error() string { return "collecting a profile: " + e.out }
func (e *profileError) Unwrap() error { return e.err }
