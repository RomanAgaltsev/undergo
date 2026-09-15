package widget

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions builds three artifacts and grades how many of Widget's four
// methods survive in each.
func TestPredictions(t *testing.T) {
	archive, direct, byname, err := Counts()
	if err != nil {
		t.Fatalf("building the artifacts: %v", err)
	}

	predict.Check(t, map[string]any{
		"archive_retained":       archive,
		"direct_binary_retained": direct,
		"byname_binary_retained": byname,
	})
}
