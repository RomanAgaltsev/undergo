package cli

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// withReadme gives repo(t)'s single task a README and, optionally, a mode other
// than the predict it ships with. Env's zero Out is nil, and Validate prints on
// success, so it also gets a writer.
func withReadme(t *testing.T, mode, readme string) Env {
	t.Helper()
	e := repo(t)
	e.Out = io.Discard
	dir := e.TaskDir("layout/01-struct-padding")

	if mode != "predict" {
		p := filepath.Join(dir, "task.yaml")
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		y := strings.Replace(string(b), "mode: predict", "mode: "+mode, 1)
		// A predict block on a non-predict task is itself a manifest error,
		// so the mode and its block move together.
		y = regexp.MustCompile(`(?m)^predict:.*
`).ReplaceAllString(y, "")
		if err := os.WriteFile(p, []byte(y), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	return e
}

// The policy is that a predict task's README names its instrument but does not
// print the command that runs it. A check nobody has watched go red is a check
// that proves nothing, so this test makes it go red on purpose.
func TestValidateRejectsAnInstrumentCommandInAPredictReadme(t *testing.T) {
	e := withReadme(t, "predict", "# t\n\nRun it:\n\n```\ngo build -gcflags=-m .\n```\n")
	if err := Validate(e, nil); err == nil {
		t.Fatal("expected validate to reject a predict README that prints an instrument command")
	}
}

func TestValidateAllowsAnInstrumentCommandInABuildReadme(t *testing.T) {
	e := withReadme(t, "build", "# t\n\nRun it:\n\n```\ngo build -gcflags=-m .\n```\n")
	if err := Validate(e, nil); err != nil {
		t.Fatalf("build mode permits the command: %v", err)
	}
}

// Naming the instrument in prose is exactly what the policy wants, so the rule
// must scan fenced blocks and not the whole document.
func TestValidateAllowsNamingTheInstrumentWithoutTheCommand(t *testing.T) {
	e := withReadme(t, "predict", "# t\n\nThe test compiles with -gcflags=-m and reads the diagnostics.\n")
	if err := Validate(e, nil); err != nil {
		t.Fatalf("naming the instrument is allowed: %v", err)
	}
}

// The written questions are post-prediction work, and several of them ask the
// solver to re-run the subject under a control. The rule governs the task body,
// where a command would be on screen while the prediction is still being made.
func TestValidateAllowsACommandInTheWrittenQuestions(t *testing.T) {
	e := withReadme(t, "predict", "# t\n\nPredict the sequence.\n\n"+
		QuestionsHeading+"\n\n1. Run it again under a control:\n\n   ```\n   go test -count=1 .\n   ```\n")
	if err := Validate(e, nil); err != nil {
		t.Fatalf("the written questions may print a command: %v", err)
	}
}

// A task with no written-questions section is scanned end to end: the narrower
// scope is a concession to a section that exists, not a way out of the rule.
func TestValidateScansTheWholeReadmeWhenThereAreNoQuestions(t *testing.T) {
	e := withReadme(t, "predict", "# t\n\nRun it:\n\n```\nGODEBUG=gctrace=1 go test .\n```\n")
	if err := Validate(e, nil); err == nil {
		t.Fatal("expected validate to reject a command with no questions section to sit under")
	}
}

// A `go` statement launches a goroutine, and a predict task about loop
// variables prints one. Only a real go subcommand is a command.
func TestValidateAllowsAGoStatementInAGoSnippet(t *testing.T) {
	e := withReadme(t, "predict", "# t\n\n```go\nfor i := range 3 {\n\tgo func() { _ = i }()\n}\n```\n")
	if err := Validate(e, nil); err != nil {
		t.Fatalf("a go statement is not a command: %v", err)
	}
}

// `go 1.21` is a go.mod directive. The versions track prints go.mod excerpts,
// so the rule must not read one as an invocation of the go tool.
func TestValidateAllowsAGoModDirective(t *testing.T) {
	e := withReadme(t, "predict", "# t\n\n```\nmodule m\n\ngo 1.21\n```\n")
	if err := Validate(e, nil); err != nil {
		t.Fatalf("a go.mod directive is not a command: %v", err)
	}
}

// A shell-tagged fence is still a command fence.
func TestValidateRejectsACommandInAShellFence(t *testing.T) {
	e := withReadme(t, "predict", "# t\n\n```sh\ngo test -run X .\n```\n")
	if err := Validate(e, nil); err == nil {
		t.Fatal("expected validate to reject a command in a ```sh fence")
	}
}

// Gate 1 skips what its host cannot build, which is only safe while some host
// builds everything. Gate 1 cannot check that — it sees one platform — so gate 3
// does. A task pinned where no runner runs would be built nowhere and skipped
// everywhere, and a skip nobody is told about reads exactly like a pass.
func TestValidateRejectsAPlatformNoRunnerHas(t *testing.T) {
	tests := []struct {
		name       string
		requires   string
		wantReject bool
	}{
		{name: "no pin at all", requires: "requires: {go: \"1.27\"}"},
		{name: "amd64, which ubuntu-latest is", requires: "requires: {go: \"1.27\", arch: [amd64]}"},
		{name: "arm64, which macos-latest is", requires: "requires: {go: \"1.27\", arch: [arm64]}"},
		{name: "both architectures", requires: "requires: {go: \"1.27\", arch: [amd64, arm64]}"},
		{
			name:       "an architecture no runner has",
			requires:   "requires: {go: \"1.27\", arch: [riscv64]}",
			wantReject: true,
		},
		{
			name:       "an operating system no runner has",
			requires:   "requires: {go: \"1.27\", os: [windows]}",
			wantReject: true,
		},
		{
			// Each half is covered by a different runner and neither covers
			// both, so the pair is reachable nowhere.
			name:       "an arch and an os that no single runner pairs",
			requires:   "requires: {go: \"1.27\", arch: [amd64], os: [darwin]}",
			wantReject: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := withReadme(t, "predict", "# t\n\nPredict the size.\n")
			p := filepath.Join(e.TaskDir("layout/01-struct-padding"), "task.yaml")
			b, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			y := strings.Replace(string(b), `requires: {go: "1.27"}`, tc.requires, 1)
			if err := os.WriteFile(p, []byte(y), 0o644); err != nil {
				t.Fatal(err)
			}

			err = Validate(e, nil)
			if tc.wantReject && err == nil {
				t.Fatalf("validate accepted %s", tc.requires)
			}
			if !tc.wantReject && err != nil {
				t.Fatalf("validate rejected a reachable platform: %v", err)
			}
		})
	}
}

// Spec §13 answers its top risk — solved work leaking into the public etalon —
// with "CI gate 3 rejects a committed plaintext solution/". It did not: the
// check named _solution/ alone, and so did .gitignore, while `solution/` is the
// name the overlay actually carries inside every seal. Both spellings now fail,
// and this is the test that says so.
func TestValidateRejectsPlaintextSolutionDirectories(t *testing.T) {
	for _, name := range []string{AuthoringDir, SolutionSubdir} {
		t.Run(name, func(t *testing.T) {
			e := withReadme(t, "predict", "# t\n\nPredict the size.\n")
			dir := filepath.Join(e.TaskDir("layout/01-struct-padding"), name)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "answer.go"),
				[]byte("package answer\n"), 0o644); err != nil {
				t.Fatal(err)
			}

			err := Validate(e, nil)
			if err == nil {
				t.Fatalf("validate accepted a committed plaintext %s/", name)
			}
			if !strings.Contains(err.Error(), name) {
				t.Errorf("error does not name %s/: %v", name, err)
			}
		})
	}
}
