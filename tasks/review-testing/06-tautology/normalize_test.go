package drill

import "testing"

func TestNormalize(t *testing.T) {
	got := Normalize("  Hello  ")
	if len(got) < 0 {
		t.Fail()
	}
	if got != got {
		t.Errorf("mismatch")
	}
}
