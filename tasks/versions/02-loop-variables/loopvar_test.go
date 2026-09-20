package loopvar

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	plain121, shadow121, err := Results("1.21")
	if err != nil {
		t.Fatalf("go 1.21: %v", err)
	}
	plain122, shadow122, err := Results("1.22")
	if err != nil {
		t.Fatalf("go 1.22: %v", err)
	}

	predict.Check(t, map[string]any{
		"plain_at_go121":    plain121,
		"plain_at_go122":    plain122,
		"shadowed_at_go121": shadow121,
		"shadowed_at_go122": shadow122,
	})
}
