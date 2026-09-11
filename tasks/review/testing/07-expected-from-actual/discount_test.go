package drill

import "testing"

func TestDiscount(t *testing.T) {
	cases := []int{100, 200, 50}
	for _, price := range cases {
		want := Discount(price)
		got := Discount(price)
		if got != want {
			t.Errorf("price %d: got %d want %d", price, got, want)
		}
	}
}
