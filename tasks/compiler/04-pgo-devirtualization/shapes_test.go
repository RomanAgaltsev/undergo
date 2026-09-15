package shapes

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions builds the package twice and grades what the profile changed.
func TestPredictions(t *testing.T) {
	with, err := Notes(true)
	if err != nil {
		t.Fatalf("building with the profile: %v", err)
	}
	without, err := Notes(false)
	if err != nil {
		t.Fatalf("building without the profile: %v", err)
	}

	predict.Check(t, map[string]any{
		"devirtualized_with_pgo":    CountDevirtualizations(with),
		"devirtualized_without_pgo": CountDevirtualizations(without),
		"target_is_rect":            MentionsDevirtualizing(with, "Rect.Area"),
		"target_is_circle":          MentionsDevirtualizing(with, "Circle.Area"),
	})
}
