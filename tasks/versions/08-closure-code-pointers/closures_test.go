package closures

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	localSame, localDistinct, err := Shares(Local)
	if err != nil {
		t.Fatalf("local: %v", err)
	}
	oldSame, oldDistinct, err := Shares("go1.26.0")
	if err != nil {
		t.Fatalf("go1.26.0: %v", err)
	}

	predict.Check(t, map[string]any{
		"same_site_shares_on_local":            localSame == "true",
		"same_site_shares_under_go1_26":        oldSame == "true",
		"distinct_literals_share_on_local":     localDistinct == "true",
		"distinct_literals_share_under_go1_26": oldDistinct == "true",
	})
}
