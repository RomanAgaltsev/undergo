package drill

import "testing"

func TestClassify(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{10, "low"},
		{60, "mid"},
		{90, "high"},
	}
	for _, c := range cases {
		got, _ := Classify(c.score)
		if got != c.want {
			t.Errorf("score %d: got %s want %s", c.score, got, c.want)
		}
	}
}
