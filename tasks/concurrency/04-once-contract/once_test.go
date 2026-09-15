package oncecontract

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	do := RunDo()

	predict.Check(t, map[string]any{
		"do_panic_reaches_caller":     do.FirstPanicked,
		"do_second_call_runs_f":       do.SecondRanF,
		"do_second_call_panics":       do.SecondPanicked,
		"oncefunc_second_call_panics": OnceFuncSecondCallPanics(),
		"oncevalue_calls_f_twice":     OnceValueFCalls() == 2,
	})
}
