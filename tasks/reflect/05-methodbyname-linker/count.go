package widget

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
)

// methodSymbol matches a Widget method's own symbol and nothing else.
//
// The trailing anchor matters. `go tool nm` also emits .arginfo1,
// .argliveinfo and go:info.*$abstract entries per symbol, so an unanchored
// match turns four into twelve.
var methodSymbol = regexp.MustCompile(`\.Widget\.(Alpha|Bravo|Charlie|Delta)$`)

// countMethods reports how many of Widget's four methods appear as symbols in
// the object at path.
func countMethods(path string) (int, error) {
	out, err := exec.Command("go", "tool", "nm", path).Output()
	if err != nil {
		return 0, err
	}
	n := 0
	for _, line := range splitLines(string(out)) {
		if methodSymbol.MatchString(line) {
			n++
		}
	}
	return n, nil
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := range len(s) {
		if s[i] == '\n' {
			out = append(out, trimCR(s[start:i]))
			start = i + 1
		}
	}
	return append(out, trimCR(s[start:]))
}

func trimCR(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\r' {
		return s[:len(s)-1]
	}
	return s
}

// Counts builds three artifacts and reports how many Widget methods survive in
// each: the package archive, a binary that calls one method directly, and a
// binary that reaches a method by name.
func Counts() (archive, direct, byname int, err error) {
	dir, err := os.MkdirTemp("", "undergo-linker-")
	if err != nil {
		return 0, 0, 0, err
	}
	defer func() { _ = os.RemoveAll(dir) }()

	exeSuffix := ""
	if runtime.GOOS == "windows" {
		exeSuffix = ".exe"
	}

	builds := []struct {
		out string
		pkg string
	}{
		{filepath.Join(dir, "lib.a"), "."},
		{filepath.Join(dir, "direct"+exeSuffix), "./testdata/direct"},
		{filepath.Join(dir, "byname"+exeSuffix), "./testdata/byname"},
	}
	counts := make([]int, len(builds))
	for i, b := range builds {
		if out, berr := exec.Command("go", "build", "-o", b.out, b.pkg).CombinedOutput(); berr != nil {
			return 0, 0, 0, &buildError{pkg: b.pkg, out: string(out), err: berr}
		}
		n, cerr := countMethods(b.out)
		if cerr != nil {
			return 0, 0, 0, cerr
		}
		counts[i] = n
	}
	return counts[0], counts[1], counts[2], nil
}

type buildError struct {
	pkg string
	out string
	err error
}

func (e *buildError) Error() string { return "building " + e.pkg + ": " + e.out }
func (e *buildError) Unwrap() error { return e.err }
