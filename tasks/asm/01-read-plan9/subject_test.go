package subject

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

var functions = []string{"AddInts", "LenOf", "FirstOf", "DivMod"}

// listing is one function's compiled form: its header line and its body.
type listing struct {
	header string
	body   []string
}

// listings compiles this package and splits the assembly by function.
//
// Symbols carry the full import path, so a function's listing opens with
// "<path>.<Name> STEXT ..." at column zero and its instructions are indented
// beneath it.
func listings(t *testing.T) map[string]listing {
	t.Helper()

	out, err := exec.Command("go", "build", "-gcflags=-S", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=-S: %v\n%s", err, out)
	}

	found := map[string]listing{}
	current := ""
	for line := range strings.SplitSeq(string(out), "\n") {
		if !strings.HasPrefix(line, "\t") {
			current = ""
			for _, name := range functions {
				if strings.Contains(line, "."+name+" STEXT") {
					current = name
					found[name] = listing{header: line}
					break
				}
			}
			continue
		}
		if current != "" {
			l := found[current]
			l.body = append(l.body, line)
			found[current] = l
		}
	}

	for _, name := range functions {
		if _, ok := found[name]; !ok {
			t.Fatalf("%s did not appear in the assembly listing", name)
		}
	}
	return found
}

var localsRE = regexp.MustCompile(`locals=0x([0-9a-f]+)`)

// frameBytes reads the local frame size out of a function's header line.
func frameBytes(t *testing.T, l listing) int {
	t.Helper()
	m := localsRE.FindStringSubmatch(l.header)
	if m == nil {
		t.Fatalf("no locals= in header: %s", l.header)
	}
	n, err := strconv.ParseInt(m[1], 16, 64)
	if err != nil {
		t.Fatalf("parsing locals: %v", err)
	}
	return int(n)
}

// uses reports whether any instruction in the body mentions text.
func uses(l listing, text string) bool {
	for _, line := range l.body {
		if strings.Contains(line, text) {
			return true
		}
	}
	return false
}

// TestPredictions reads the compiler's own assembly and compares it with your
// answers. It never prints an instruction.
func TestPredictions(t *testing.T) {
	ls := listings(t)

	predict.Check(t, map[string]any{
		"add_frame_bytes":      frameBytes(t, ls["AddInts"]),
		"add_is_nosplit":       strings.Contains(ls["AddInts"].header, "nosplit"),
		"add_uses_ax_and_bx":   uses(ls["AddInts"], "ADDQ\tBX, AX"),
		"lenof_moves_bx_to_ax": uses(ls["LenOf"], "MOVQ\tBX, AX"),
		"firstof_can_panic":    uses(ls["FirstOf"], "panicIndex"),
		"divmod_frame_is_zero": frameBytes(t, ls["DivMod"]) == 0,
	})
}
