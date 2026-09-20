package panicprint

import (
	"strings"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	line := func(v string) string {
		t.Helper()
		got, err := PanicLine(v)
		if err != nil {
			t.Fatalf("panicking with %s: %v", v, err)
		}
		return got
	}

	plain := line(PlainStruct)

	predict.Check(t, map[string]any{
		"named_int_line":                line(NamedInt),
		"named_string_line":             line(NamedString),
		"error_value_line":              line(ErrorValue),
		"plain_struct_line_has_address": strings.Contains(plain, "0x"),
	})
}
