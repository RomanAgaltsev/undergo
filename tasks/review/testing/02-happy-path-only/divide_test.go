package drill

import "testing"

func TestDivide(t *testing.T) {
	got, _ := Divide(6, 3)
	if got != 2 {
		t.Fatalf("Divide(6, 3) = %d, want 2", got)
	}
}
