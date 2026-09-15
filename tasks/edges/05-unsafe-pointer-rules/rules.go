// Package rules asks which unsafe.Pointer misuses go vet can actually find.
//
// unsafe.Pointer's documentation lists six legal conversion patterns and says
// plainly that anything else is invalid. testdata/patterns holds eight
// snippets: five of them follow a listed pattern and three do not.
//
// The snippets carry //go:build ignore so they are not part of this package.
// They have to be, because a task whose files deliberately fail vet cannot live
// in a package that gate 1 vets.
package rules

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// Snippets are the eight files, in the order the task asks about them.
var Snippets = []string{"p1", "p2", "p3", "p4", "p5", "v1", "v2", "v3"}

// VetRejects reports whether `go vet -unsafeptr` complains about one snippet.
func VetRejects(name string) (bool, error) {
	path := filepath.Join("testdata", "patterns", name+".go")
	out, err := exec.Command("go", "vet", "-unsafeptr", path).CombinedOutput()
	if err == nil {
		return false, nil
	}
	text := string(out)
	// Distinguish a vet diagnostic from a snippet that does not compile: the
	// second would make the answer meaningless.
	if strings.Contains(text, "possible misuse") {
		return true, nil
	}
	return false, &vetError{name: name, out: text}
}

type vetError struct {
	name string
	out  string
}

func (e *vetError) Error() string {
	return "vetting " + e.name + " failed for a reason other than unsafeptr: " + e.out
}
