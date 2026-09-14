package bce

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// funcRange is the line span of one function in bce.go.
type funcRange struct {
	name       string
	start, end int
}

// functionRanges reads bce.go and returns the line span of every function.
//
// The spans come from the source rather than being written down, because a
// bounds-check diagnostic lands on the indexing expression, which is rarely the
// line you would have guessed.
func functionRanges(t *testing.T) []funcRange {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "bce.go", nil, 0)
	if err != nil {
		t.Fatalf("parsing bce.go: %v", err)
	}

	var out []funcRange
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		out = append(out, funcRange{
			name:  fn.Name.Name,
			start: fset.Position(fn.Pos()).Line,
			end:   fset.Position(fn.End()).Line,
		})
	}
	if len(out) == 0 {
		t.Fatal("no functions found in bce.go")
	}
	return out
}

var checkRE = regexp.MustCompile(`bce\.go:(\d+):\d+: Found IsInBounds`)

// checked compiles this package with the bounds-check report on and returns the
// set of functions that kept at least one check.
func checked(t *testing.T) map[string]bool {
	t.Helper()

	out, err := exec.Command("go", "build", "-gcflags=-d=ssa/check_bce", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -d=ssa/check_bce: %v\n%s", err, out)
	}

	ranges := functionRanges(t)
	kept := map[string]bool{}
	for line := range strings.SplitSeq(string(out), "\n") {
		m := checkRE.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		at, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		for _, r := range ranges {
			if at >= r.start && at <= r.end {
				kept[r.name] = true
			}
		}
	}
	return kept
}

// TestPredictions asks the compiler which checks survived and compares with
// your answers. It never prints a line number.
func TestPredictions(t *testing.T) {
	kept := checked(t)

	predict.Check(t, map[string]any{
		"ranged_checked":          kept["Ranged"],
		"counted_checked":         kept["Counted"],
		"two_slices_checked":      kept["TwoSlices"],
		"unchecked_checked":       kept["Unchecked"],
		"after_one_check_checked": kept["AfterOneCheck"],
		"next_element_checked":    kept["NextElement"],
	})
}
