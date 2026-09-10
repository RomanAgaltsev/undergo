package predict

import "testing"

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
