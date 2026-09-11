package escape

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

// funcRange is the line span of one function in escape.go.
type funcRange struct {
	name       string
	start, end int
}

// functionRanges reads escape.go and returns the line span of every function.
//
// The spans are derived from the source rather than written down, because a
// diagnostic does not always land on the line you expect: the compiler reports
// Interfaced's escape on the conversion, not on the composite literal. Matching
// by enclosing function is both simpler and immune to edits above.
func functionRanges(t *testing.T) []funcRange {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "escape.go", nil, 0)
	if err != nil {
		t.Fatalf("parsing escape.go: %v", err)
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
		t.Fatal("no functions found in escape.go")
	}
	return out
}

var diagnosticRE = regexp.MustCompile(`^\./escape\.go:(\d+):\d+: (.*)$`)

// escapedFunctions compiles this package with escape diagnostics on and returns
// the set of functions the compiler reported something escaping from.
func escapedFunctions(t *testing.T) map[string]bool {
	t.Helper()

	cmd := exec.Command("go", "build", "-gcflags=-m", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=-m: %v\n%s", err, out)
	}

	ranges := functionRanges(t)
	escaped := map[string]bool{}
	for line := range strings.SplitSeq(string(out), "\n") {
		m := diagnosticRE.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		message := m[2]
		if !strings.Contains(message, "escapes to heap") && !strings.Contains(message, "moved to heap") {
			continue
		}
		at, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		for _, r := range ranges {
			if at >= r.start && at <= r.end {
				escaped[r.name] = true
			}
		}
	}
	return escaped
}

// TestPredictions asks the compiler which Holders escape and compares its answer
// with yours. It never prints the compiler's output.
func TestPredictions(t *testing.T) {
	escaped := escapedFunctions(t)

	predict.Check(t, map[string]any{
		"local_escapes":      escaped["Local"],
		"returned_escapes":   escaped["Returned"],
		"stored_escapes":     escaped["Stored"],
		"interfaced_escapes": escaped["Interfaced"],
	})
}
