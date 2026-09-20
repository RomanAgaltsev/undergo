package compiles

import (
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	got := map[string]any{}
	var lastDiagnostic string

	for _, c := range []struct {
		slot string
		line string
		src  string
	}{
		{"clear_min_max_at_go120", "1.20", ClearMinMax},
		{"range_int_at_go121", "1.21", RangeInt},
		{"range_func_at_go122", "1.22", RangeFunc},
		{"new_expr_at_go125", "1.25", NewExpr},
		{"generic_alias_at_go123", "1.23", GenericAlias},
	} {
		ok, diag, err := Compiles(c.line, c.src)
		if err != nil {
			t.Fatalf("%s: %v", c.slot, err)
		}
		got[c.slot] = ok
		if !ok {
			lastDiagnostic = diag
		}
	}

	// The diagnostic names the mechanism, which is the part worth knowing.
	got["diagnostic_names_lang_flag"] = strings.Contains(lastDiagnostic, "-lang was set to")

	predict.Check(t, got)
}
