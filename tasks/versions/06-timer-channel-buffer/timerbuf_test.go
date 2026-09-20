package timerbuf

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	get := func(name, line string) string {
		t.Helper()
		got, err := Cap(name, line)
		if err != nil {
			t.Fatalf("%s at go %s: %v", name, line, err)
		}
		return got
	}

	predict.Check(t, map[string]any{
		"cap_on_local":                get(Local, "1.27"),
		"cap_under_go1_22":            get("go1.22.12", "1.19"),
		"cap_on_local_at_go_line_119": get(Local, "1.19"),
	})
}
