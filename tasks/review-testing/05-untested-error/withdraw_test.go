package drill

import "testing"

func TestWithdraw(t *testing.T) {
	got, err := Withdraw(100, 30)
	if err != nil {
		t.Fatal(err)
	}
	if got != 70 {
		t.Fatalf("got %d, want 70", got)
	}
}
