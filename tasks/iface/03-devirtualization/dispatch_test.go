package dispatch

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

type funcRange struct {
	name       string
	start, end int
}

// functionRanges reads dispatch.go and returns each function's line span.
//
// The spans come from the source rather than being written down: a
// devirtualization diagnostic lands on the call, not on the declaration.
func functionRanges(t *testing.T) []funcRange {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "dispatch.go", nil, 0)
	if err != nil {
		t.Fatalf("parsing dispatch.go: %v", err)
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
		t.Fatal("no functions found in dispatch.go")
	}
	return out
}

var devirtRE = regexp.MustCompile(`dispatch\.go:(\d+):\d+: devirtualizing`)

// devirtualized compiles this package with the compiler's reasoning on and
// returns the set of functions containing a devirtualized call.
func devirtualized(t *testing.T) map[string]bool {
	t.Helper()

	out, err := exec.Command("go", "build", "-gcflags=-m=2", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=-m=2: %v\n%s", err, out)
	}

	ranges := functionRanges(t)
	found := map[string]bool{}
	for line := range strings.SplitSeq(string(out), "\n") {
		m := devirtRE.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		at, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		for _, r := range ranges {
			if at >= r.start && at <= r.end {
				found[r.name] = true
			}
		}
	}
	return found
}

// TestPredictions asks the compiler which calls it resolved and compares with
// your answers. It never prints a diagnostic.
func TestPredictions(t *testing.T) {
	d := devirtualized(t)

	predict.Check(t, map[string]any{
		"one_type_devirtualized":         d["OneType"],
		"reassigned_devirtualized":       d["Reassigned"],
		"via_parameter_devirtualized":    d["ViaParameter"],
		"via_global_devirtualized":       d["ViaGlobal"],
		"from_constructor_devirtualized": d["FromConstructor"],
		"in_slice_devirtualized":         d["InSlice"],
	})
}
