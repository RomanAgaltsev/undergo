package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// A task's answer must come from the task, not from the solver's shell. The
// frozen tests of every predict task run through RunTests, and whatever
// environment they get is also what every go command they start inherits — so
// a GOFLAGS=-gcflags=... or a GODEBUG in the caller's environment would move
// answers measured in-process and in children alike.
func TestRunTestsGivesTasksACleanEnvironment(t *testing.T) {
	t.Setenv("GOFLAGS", "-gcflags=all=-N")
	t.Setenv("GODEBUG", "gctrace=1")

	dir := t.TempDir()
	files := map[string]string{
		"go.mod":   "module probe\n\ngo 1.27\n",
		"probe.go": "package probe\n",
		"probe_test.go": `package probe

import (
	"os"
	"testing"
)

func TestEnvironment(t *testing.T) {
	if v := os.Getenv("GOFLAGS"); v != "-buildvcs=auto" {
		t.Errorf("GOFLAGS = %q, want -buildvcs=auto", v)
	}
	if v, ok := os.LookupEnv("GODEBUG"); ok {
		t.Errorf("GODEBUG = %q reached the task", v)
	}
}
`,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	var out bytes.Buffer
	passed, err := RunTests(Env{Out: &out, Err: &out}, dir)
	if err != nil {
		t.Fatal(err)
	}
	if !passed {
		t.Errorf("the task saw the caller's environment:\n%s", out.String())
	}
}
