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
