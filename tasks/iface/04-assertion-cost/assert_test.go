package assert

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

var functions = []string{"ConcreteCommaOk", "ConcretePanicking", "EmptyToInterface", "InterfaceToInterface", "TypeSwitchConcrete"}

// runtimeHelped compiles this package and reports, per function, whether the
// assertion needed a call into the runtime.
//
// Symbols carry the full import path, so a function's listing opens with
// "<path>.<Name> STEXT" at column zero.
func runtimeHelped(t *testing.T) map[string]bool {
	t.Helper()

	out, err := exec.Command("go", "build", "-gcflags=-S", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=-S: %v\n%s", err, out)
	}

	helped := make(map[string]bool, len(functions))
	current := ""
	for line := range strings.SplitSeq(string(out), "\n") {
		if !strings.HasPrefix(line, "\t") {
			current = ""
			for _, name := range functions {
				if strings.Contains(line, "."+name+" STEXT") {
					current = name
					helped[name] = false
					break
				}
			}
			continue
		}
		if current == "" || !strings.Contains(line, "CALL") {
			continue
		}
		if strings.Contains(line, "runtime.typeAssert") ||
			strings.Contains(line, "runtime.assert") ||
			strings.Contains(line, "runtime.panicdottype") {
			helped[current] = true
		}
	}

	for _, name := range functions {
		if _, ok := helped[name]; !ok {
			t.Fatalf("%s did not appear in the assembly listing", name)
		}
	}
	return helped
}

// TestPredictions asks the compiler which assertions needed the runtime and
// compares with your answers. It never prints an instruction.
func TestPredictions(t *testing.T) {
	h := runtimeHelped(t)

	predict.Check(t, map[string]any{
		"concrete_comma_ok_calls_runtime":      h["ConcreteCommaOk"],
		"concrete_panicking_calls_runtime":     h["ConcretePanicking"],
		"empty_to_interface_calls_runtime":     h["EmptyToInterface"],
		"interface_to_interface_calls_runtime": h["InterfaceToInterface"],
		"type_switch_calls_runtime":            h["TypeSwitchConcrete"],
	})
}
