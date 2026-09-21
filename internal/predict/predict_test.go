package predict

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// The interface exists to make Check testable; this is the assertion that it
// still describes the thing it stands in for.
var _ TB = (*testing.T)(nil)

// recorder stands in for *testing.T so that a grading can be inspected instead
// of being obeyed. Its Fatalf ends the calling goroutine exactly as the real one
// does, because a Check that stops at the first fatal and a Check that carries on
// are different functions, and the one under test has to be the one that ships.
type recorder struct {
	errs  []string
	logs  []string
	fatal string
}

func (r *recorder) Helper() {}

func (r *recorder) Errorf(format string, args ...any) {
	r.errs = append(r.errs, fmt.Sprintf(format, args...))
}

func (r *recorder) Logf(format string, args ...any) {
	r.logs = append(r.logs, fmt.Sprintf(format, args...))
}

func (r *recorder) Fatalf(format string, args ...any) {
	r.fatal = fmt.Sprintf(format, args...)
	runtime.Goexit()
}

// failed reports what a real *testing.T would have reported: any Errorf, or a
// Fatalf, means the task did not pass.
func (r *recorder) failed() bool { return len(r.errs) > 0 || r.fatal != "" }

func (r *recorder) all() string {
	return strings.Join(append(append([]string{}, r.errs...), r.logs...), "\n") + "\n" + r.fatal
}

// grade writes a prediction file and runs Check against it, on its own goroutine
// so that a Fatalf can end that goroutine without ending the test.
func grade(t *testing.T, file string, measured map[string]any) *recorder {
	t.Helper()
	path := filepath.Join(t.TempDir(), DefaultFile)
	if err := os.WriteFile(path, []byte(file), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("UNDERGO_PREDICTION", path)

	r := &recorder{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		Check(r, measured)
	}()
	<-done
	return r
}

// measured is one slot of each kind Check has to compare: an integer, a boolean
// and a string. Between them they cover every branch of canonical.
func measured() map[string]any {
	return map[string]any{
		"sizeof_header": 24,
		"escapes":       true,
		"caps":          "1 2 4 8",
	}
}

// This is the test the repository did not have. Check grades 77 of the 96
// machine-graded tasks, and gate 2 cannot prove it works: gate 2 feeds it only
// reference answers, which are correct by construction, so a Check that accepted
// everything would leave every gate green and every task ungraded.
func TestCheckRejectsWrongAnswers(t *testing.T) {
	tests := []struct {
		name       string
		file       string
		wantFailed bool
		mentions   string
	}{
		{
			name: "every slot correct",
			file: "sizeof_header: 24\nescapes: true\ncaps: 1 2 4 8\n",
		},
		{
			name:       "one integer wrong",
			file:       "sizeof_header: 25\nescapes: true\ncaps: 1 2 4 8\n",
			wantFailed: true,
			mentions:   `slot "sizeof_header": incorrect`,
		},
		{
			name:       "boolean inverted",
			file:       "sizeof_header: 24\nescapes: false\ncaps: 1 2 4 8\n",
			wantFailed: true,
			mentions:   `slot "escapes": incorrect`,
		},
		{
			name:       "sequence off by one element",
			file:       "sizeof_header: 24\nescapes: true\ncaps: 1 2 4 16\n",
			wantFailed: true,
			mentions:   `slot "caps": incorrect`,
		},
		{
			name:       "a slot left out",
			file:       "sizeof_header: 24\nescapes: true\n",
			wantFailed: true,
			mentions:   `slot "caps": no prediction`,
		},
		{
			name:       "a slot left blank",
			file:       "sizeof_header: 24\nescapes: true\ncaps:\n",
			wantFailed: true,
			mentions:   `slot "caps": no prediction`,
		},
		{
			name:       "a slot this task does not have",
			file:       "sizeof_header: 24\nescapes: true\ncaps: 1 2 4 8\nbogus: 1\n",
			wantFailed: true,
			mentions:   `slot "bogus": not a slot in this task`,
		},
		{
			name:       "every slot wrong",
			file:       "sizeof_header: 0\nescapes: false\ncaps: nothing\n",
			wantFailed: true,
		},
		{
			// canonical exists so that a solver is graded on the value and not
			// on how they wrote it down.
			name: "hex, quotes and padding all accepted",
			file: "sizeof_header:   0x18   # a comment\nescapes: \"true\"\ncaps: '1 2 4 8'\n",
		},
		{
			name:       "a file that is not key: value",
			file:       "this line has no colon\n",
			wantFailed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := grade(t, tc.file, measured())
			if got := r.failed(); got != tc.wantFailed {
				t.Errorf("failed = %v, want %v\n%s", got, tc.wantFailed, r.all())
			}
			if tc.mentions != "" && !strings.Contains(r.all(), tc.mentions) {
				t.Errorf("output does not mention %q:\n%s", tc.mentions, r.all())
			}
		})
	}
}

// Flipping any single slot must fail the task. A grader that catches some slots
// and not others is the same defect as one that catches none, only harder to
// notice — so this asserts the property over every slot rather than sampling it.
func TestCheckFailsOnEverySingleWrongSlot(t *testing.T) {
	correct := map[string]string{
		"sizeof_header": "24",
		"escapes":       "true",
		"caps":          "1 2 4 8",
	}
	wrong := map[string]string{
		"sizeof_header": "999",
		"escapes":       "false",
		"caps":          "nothing like it",
	}

	for slot := range correct {
		t.Run(slot, func(t *testing.T) {
			var b strings.Builder
			for name, value := range correct {
				if name == slot {
					value = wrong[name]
				}
				fmt.Fprintf(&b, "%s: %s\n", name, value)
			}
			r := grade(t, b.String(), measured())
			if !r.failed() {
				t.Errorf("a wrong %q was graded as correct:\n%s", slot, r.all())
			}
		})
	}
}

// Check must never print a measured value: a solver who got a slot wrong would
// otherwise be handed the answer on their first attempt, which is the whole
// reason a predict task stores no answer.
func TestCheckNeverPrintsTheMeasuredValue(t *testing.T) {
	m := map[string]any{"sizeof_header": 24, "escapes": true, "caps": "1 2 4 8"}
	r := grade(t, "sizeof_header: 999\nescapes: false\ncaps: wrong\n", m)
	if !r.failed() {
		t.Fatal("expected a failure to inspect")
	}
	for _, secret := range []string{"24", "true", "1 2 4 8"} {
		if strings.Contains(r.all(), secret) {
			t.Errorf("grading leaked the measured value %q:\n%s", secret, r.all())
		}
	}
}

func TestCheckFatalsWhenThePredictionFileIsMissing(t *testing.T) {
	t.Setenv("UNDERGO_PREDICTION", filepath.Join(t.TempDir(), "absent.yaml"))
	r := &recorder{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		Check(r, measured())
	}()
	<-done

	if r.fatal == "" {
		t.Fatal("a missing prediction file must be fatal, not a silent pass")
	}
	if len(r.errs) > 0 {
		t.Errorf("Fatalf must stop the grading, but it continued: %v", r.errs)
	}
}

func TestParseFlatFile(t *testing.T) {
	got, err := parse([]byte("# a comment\nsizeof_header: 24\nname: \"padded\"\n\nflag: true   # trailing\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"sizeof_header": "24", "name": "padded", "flag": "true"}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}

func TestParseRejectsGarbage(t *testing.T) {
	if _, err := parse([]byte("this line has no colon\n")); err == nil {
		t.Fatal("parse accepted a line with no key")
	}
}

func TestCanonicalNormalisesAcrossTypes(t *testing.T) {
	tests := []struct{ a, b any }{
		{"24", uintptr(24)},
		{"24", 24},
		{"0x18", 24},
		{"true", true},
		{"1.5", 1.5},
	}
	for _, tc := range tests {
		if canonical(tc.a) != canonical(tc.b) {
			t.Errorf("canonical(%v)=%q != canonical(%v)=%q",
				tc.a, canonical(tc.a), tc.b, canonical(tc.b))
		}
	}
	if canonical("25") == canonical(24) {
		t.Error("canonical collapsed two different values")
	}
}
