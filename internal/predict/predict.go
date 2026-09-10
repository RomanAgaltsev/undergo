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
	"testing"
)

// DefaultFile is where a prediction lives inside a work directory.
const DefaultFile = "prediction.yaml"

// Check compares each measured slot against the solver's prediction file.
func Check(t *testing.T, measured map[string]any) {
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

	correct := 0
	slots := slices.Sorted(maps.Keys(measured))
	for _, slot := range slots {
		answer, ok := got[slot]
		if !ok || answer == "" {
			t.Errorf("slot %q: no prediction", slot)
			continue
		}
		if canonical(answer) == canonical(measured[slot]) {
			t.Logf("slot %q: correct", slot)
			correct++
			continue
		}
		t.Errorf("slot %q: incorrect", slot)
	}
	for slot := range got {
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
