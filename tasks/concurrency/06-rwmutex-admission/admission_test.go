package rwadmission

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	a := Measure()
	if !a.SetUp {
		t.Fatal("the writer never blocked; the experiment did not set up")
	}

	predict.Check(t, map[string]any{
		"second_reader_parked":                        a.SecondReaderParked,
		"second_reader_admitted_while_writer_pending": a.SecondReaderAdmitted,
		"acquisition_order":                           a.Order,
		"writer_waits_for_first_reader":               a.WriterWaitedForFirst,
	})
}
