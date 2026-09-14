package barrier

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// functions are the five stores this task asks about.
var functions = []string{"PointerStore", "ScalarStore", "NilStore", "LocalStore", "SliceStore"}

// barriers compiles this package with the assembly listing on and reports, per
// function, whether the compiler emitted a write barrier.
//
// It reports whether a barrier is present rather than how many symbols it took:
// the compiler emits more than one per barrier today, and a task asserting a
// count would break the moment that changed.
func barriers(t *testing.T) map[string]bool {
	t.Helper()

	out, err := exec.Command("go", "build", "-gcflags=-S", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=-S: %v\n%s", err, out)
	}

	found := make(map[string]bool, len(functions))
	current := ""
	for line := range strings.SplitSeq(string(out), "\n") {
		// A function's listing opens with its symbol at column zero; its
		// instructions are indented beneath it.
		if !strings.HasPrefix(line, "\t") {
			current = ""
			for _, name := range functions {
				// A function header reads "<import path>.<Name> STEXT ...".
				if strings.Contains(line, "."+name+" STEXT") {
					current = name
					found[name] = false
					break
				}
			}
			continue
		}
		if current != "" && strings.Contains(line, "gcWriteBarrier") {
			found[current] = true
		}
	}

	for _, name := range functions {
		if _, ok := found[name]; !ok {
			t.Fatalf("%s did not appear in the assembly listing", name)
		}
	}
	return found
}

// TestPredictions asks the compiler which stores are barriered and compares its
// answer with yours. It never prints the listing.
func TestPredictions(t *testing.T) {
	got := barriers(t)

	predict.Check(t, map[string]any{
		"pointer_store_barriered": got["PointerStore"],
		"scalar_store_barriered":  got["ScalarStore"],
		"nil_store_barriered":     got["NilStore"],
		"local_store_barriered":   got["LocalStore"],
		"slice_store_barriered":   got["SliceStore"],
	})
}
