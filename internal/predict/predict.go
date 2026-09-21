// Package predict grades a solver's written predictions against values the
// task measures at run time.
//
// A predict task stores no answer. The task's own test computes the truth and
// hands it here; this package compares it to prediction.yaml and reports only
// "correct" or "incorrect" per slot. It must never print a measured value,
// because that would hand over the answer on the first wrong attempt.
package predict

import (
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// DefaultFile is where a prediction lives inside a work directory.
const DefaultFile = "prediction.yaml"

// CIEnv is set by `undergo ci-verify` to say that the reader of this grading is
// a maintainer rather than a solver.
const CIEnv = "UNDERGO_CI"

// namesWrongSlots reports whether the grading may say WHICH slots are wrong.
//
// It may not, for a solver. `undergo verify` runs `go test -v`, so every t.Logf
// reaches the screen, and 34 of the 77 predict tasks grade booleans only — some
// of them nine or ten of them. Per-slot feedback turns that into a perfect
// oracle: answer everything true, read back which slots came out wrong, flip
// exactly those, and the task falls in two runs with no understanding at all.
// That is a different thing from the seal, which is deliberately only
// obfuscation: opening a seal is a deliberate act and is recorded as a peek,
// while hill-climbing is neither.
//
// Withholding the names costs the solver nothing they need. "You were wrong,
// and here is how many" is the signal; "slot 4 was wrong" is the ladder.
//
// ci-verify sets CIEnv because there the reader is the maintainer, the answers
// are already known to be correct, and naming the failing slots is what turned
// two blind gate-2 failures into a one-run diagnosis in M10. Slot names are
// public in task.yaml, and this package never prints a measured value either
// way.
func namesWrongSlots() bool { return os.Getenv(CIEnv) != "" }

// TB is the part of *testing.T that Check uses.
//
// It exists so that Check can be tested at all, which matters more here than the
// indirection costs. Check grades 77 of the 96 machine-graded tasks and is the
// only thing standing between a wrong prediction and a pass — and gate 2 cannot
// prove it works, because gate 2 only ever feeds it reference answers, which are
// correct by construction. A Check that accepted everything would leave every
// gate in this repository green.
//
// Taking *testing.T concretely made that test impossible twice over: testing.TB
// has an unexported method, so it cannot be implemented outside the testing
// package, and a failing subtest fails its parent, so t.Run's bool cannot express
// "this grading should have failed" either.
//
// *testing.T satisfies this interface, so every task's predict.Check(t, ...) is
// unchanged.
type TB interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Logf(format string, args ...any)
}

// Check compares each measured slot against the solver's prediction file.
func Check(t TB, measured map[string]any) {
	t.Helper()

	path := os.Getenv("UNDERGO_PREDICTION")
	if path == "" {
		path = DefaultFile
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read predictions from %s: %v", path, err)
	}
	got, err := parse(raw)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}

	named := namesWrongSlots()

	correct := 0
	slots := slices.Sorted(maps.Keys(measured))
	var wrong []string
	for _, slot := range slots {
		answer, ok := got[slot]
		if !ok || answer == "" {
			// Which slot is unanswered is not a hint about any answer, so this
			// one is always named.
			t.Errorf("slot %q: no prediction", slot)
			continue
		}
		if canonical(answer) == canonical(measured[slot]) {
			if named {
				t.Logf("slot %q: correct", slot)
			}
			correct++
			continue
		}
		wrong = append(wrong, slot)
		if named {
			t.Errorf("slot %q: incorrect", slot)
		}
	}
	if !named && len(wrong) > 0 {
		t.Errorf("%d of %d slots incorrect", len(wrong), len(slots))
	}
	for _, slot := range slices.Sorted(maps.Keys(got)) {
		if _, ok := measured[slot]; !ok {
			t.Errorf("slot %q: not a slot in this task", slot)
		}
	}
	t.Logf("%d/%d slots correct", correct, len(slots))
}

// parse reads the flat "key: value" prediction format. It is deliberately not
// a YAML parser: predictions are scalars only, and this package must stay
// dependency-free so a task never needs one to be graded.
func parse(b []byte) (map[string]string, error) {
	out := map[string]string{}
	for i, line := range strings.Split(string(b), "\n") {
		text := strings.TrimSpace(line)
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		key, value, ok := strings.Cut(text, ":")
		if !ok {
			return nil, fmt.Errorf("line %d: expected \"key: value\", got %q", i+1, text)
		}
		if cut := strings.Index(value, " #"); cut >= 0 {
			value = value[:cut]
		}
		out[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
	}
	return out, nil
}

// canonical renders a predicted string and a measured value in one comparable
// form, so that "24", 0x18 and uintptr(24) all agree.
func canonical(v any) string {
	if s, ok := v.(string); ok {
		s = strings.TrimSpace(s)
		if n, err := strconv.ParseInt(s, 0, 64); err == nil {
			return strconv.FormatInt(n, 10)
		}
		if b, err := strconv.ParseBool(s); err == nil {
			return strconv.FormatBool(b)
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return strconv.FormatFloat(f, 'g', -1, 64)
		}
		return s
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	default:
		return fmt.Sprint(v)
	}
}
