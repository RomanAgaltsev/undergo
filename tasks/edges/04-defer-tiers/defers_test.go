package defers

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

var functions = []string{"One", "Several", "InLoop", "InBranch", "Many", "WithRecover"}

// heapDeferred compiles this package and reports, per function, whether any
// defer fell back to the runtime's heap path.
//
// It looks for deferproc specifically. Every function containing a defer also
// calls deferreturn, open-coded or not, so matching on "defer" alone reports
// true for all of them — which is the first thing that goes wrong when writing
// this measurement.
func heapDeferred(t *testing.T) map[string]bool {
	t.Helper()

	out, err := exec.Command("go", "build", "-gcflags=-S", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=-S: %v\n%s", err, out)
	}

	heap := make(map[string]bool, len(functions))
	current := ""
	for line := range strings.SplitSeq(string(out), "\n") {
		if !strings.HasPrefix(line, "\t") {
			current = ""
			for _, name := range functions {
				if strings.Contains(line, "."+name+" STEXT") {
					current = name
					heap[name] = false
					break
				}
			}
			continue
		}
		if current != "" && strings.Contains(line, "CALL") && strings.Contains(line, "runtime.deferproc") {
			heap[current] = true
		}
	}

	for _, name := range functions {
		if _, ok := heap[name]; !ok {
			t.Fatalf("%s did not appear in the assembly listing", name)
		}
	}
	return heap
}

// TestPredictions asks the compiler which defers reached the runtime and
// compares with your answers. It never prints an instruction.
func TestPredictions(t *testing.T) {
	h := heapDeferred(t)

	predict.Check(t, map[string]any{
		"one_uses_heap_defer":          h["One"],
		"several_uses_heap_defer":      h["Several"],
		"in_loop_uses_heap_defer":      h["InLoop"],
		"in_branch_uses_heap_defer":    h["InBranch"],
		"many_uses_heap_defer":         h["Many"],
		"with_recover_uses_heap_defer": h["WithRecover"],
	})
}
