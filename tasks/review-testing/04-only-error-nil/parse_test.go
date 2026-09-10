package drill

import "testing"

func TestParse(t *testing.T) {
	_, err := Parse("42")
	if err != nil {
		t.Fatal(err)
	}
}
