package drill

import "testing"

func TestParseIDMatchesStd(t *testing.T) {
	for _, s := range []string{"000000", "000001", "123456", "999999"} {
		if got, want := ParseID(s), parseIDStd(s); got != want {
			t.Errorf("ParseID(%q) = %d, want %d", s, got, want)
		}
	}
}

func BenchmarkParseStd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseIDStd("123456")
	}
}

func BenchmarkParseID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseID("123456")
	}
}
