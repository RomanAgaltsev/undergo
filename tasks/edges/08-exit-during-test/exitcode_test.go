package exitcode

import (
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	out, passed, err := Run()
	if err != nil {
		t.Fatalf("running testdata/suite: %v", err)
	}

	predict.Check(t, map[string]any{
		"go_test_reports_failure": !passed,
		"output_names_os_exit":    strings.Contains(out, "os.Exit(0)"),
		"earlier_test_ran":        strings.Contains(out, "ALPHA RAN"),
		"later_test_ran":          strings.Contains(out, "GAMMA RAN"),
	})
}
