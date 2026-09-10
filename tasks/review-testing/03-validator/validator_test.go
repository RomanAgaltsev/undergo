package drill

import "testing"

func TestCheck(t *testing.T) {
	v := Validator{}

	_, err := v.Check("alice")
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"a", "b", "c"} {
		v.Check(name)
	}
}
