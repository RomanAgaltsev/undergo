package abandoned

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	got := map[string]any{}

	ch1 := make(chan int)
	go UnbufferedWorker(ch1)
	got["unbuffered_leaks"] = Parked("UnbufferedWorker")

	ch2 := make(chan int, 1)
	go CapOneWorker(ch2)
	got["capacity_one_leaks"] = Parked("CapOneWorker")

	ch3 := make(chan int)
	go DefaultWorker(ch3)
	got["select_with_default_leaks"] = Parked("DefaultWorker")

	ch4 := make(chan int, 1)
	go CapOneTwiceWorker(ch4)
	got["capacity_one_two_results_leaks"] = Parked("CapOneTwiceWorker")

	ch5 := make(chan int)
	go LateDrainWorker(ch5)
	go func() { <-ch5 }()
	got["late_drain_leaks"] = Parked("LateDrainWorker")

	predict.Check(t, got)
}
