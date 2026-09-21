package pauses

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	local, err := Names(Local)
	if err != nil {
		t.Fatalf("local: %v", err)
	}
	old, err := Names("go1.21.13")
	if err != nil {
		t.Fatalf("go1.21.13: %v", err)
	}

	const (
		oldName = "/gc/pauses:seconds"
		newName = "/sched/pauses/total/gc:seconds"
	)

	predict.Check(t, map[string]any{
		"pause_metric_count_on_local":                len(local),
		"pause_metric_count_under_go1_21":            len(old),
		"gc_pauses_seconds_present_on_local":         Has(local, oldName),
		"sched_pauses_total_gc_present_on_local":     Has(local, newName),
		"sched_pauses_total_gc_present_under_go1_21": Has(old, newName),
	})
}
