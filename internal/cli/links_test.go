package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBrokenRelPathsFindsBothLinkAndBacktickForms(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "review-concurrency", "01-request-counter")
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "rubric"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "rubric", "review-rubric.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	body := "Review against `../../../rubric/review-rubric.md`, see " +
		"[taxonomy](../../../rubric/bug-taxonomy.md), and `../../rubric/design-rubric.md`.\n"

	got := brokenRelPaths(taskDir, body)
	if len(got) != 2 {
		t.Fatalf("broken = %v, want the two that do not resolve", got)
	}
	for _, want := range []string{"../../../rubric/bug-taxonomy.md", "../../rubric/design-rubric.md"} {
		var found bool
		for _, g := range got {
			if g == want {
				found = true
			}
		}
		if !found {
			t.Errorf("%q not reported broken", want)
		}
	}
}
