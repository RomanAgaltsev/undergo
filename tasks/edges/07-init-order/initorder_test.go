package initorder

import (
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	lines, err := Order()
	if err != nil {
		t.Fatalf("running testdata/prog: %v", err)
	}
	if len(lines) == 0 {
		t.Fatal("testdata/prog printed nothing")
	}

	// Which of main's two package-level variables was initialised first.
	mainFirst := ""
	for _, l := range lines {
		if strings.HasPrefix(l, "main.First") || strings.HasPrefix(l, "main.Second") {
			mainFirst = l
			break
		}
	}

	// Does a package's own init() run after that package's variables?
	alphaVar, alphaInit := indexOf(lines, "alpha.A"), indexOf(lines, "alpha.init")

	predict.Check(t, map[string]any{
		"first_line":                       lines[0],
		"main_first_var":                   mainFirst,
		"init_runs_after_its_package_vars": alphaVar < alphaInit,
		"full_order":                       strings.Join(lines, ","),
	})
}

func indexOf(lines []string, want string) int {
	for i, l := range lines {
		if l == want {
			return i
		}
	}
	return -1
}
