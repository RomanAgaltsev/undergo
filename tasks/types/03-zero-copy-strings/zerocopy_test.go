package zerocopy

import (
	"strings"
	"testing"
)

func TestRoundTripsPreserveContent(t *testing.T) {
	for _, s := range []string{"", "a", "hello, world", strings.Repeat("x", 4096)} {
		if got := string(StringToBytes(s)); got != s {
			t.Errorf("StringToBytes(%q) round-tripped to %q", s, got)
		}
		if got := BytesToString([]byte(s)); got != s {
			t.Errorf("BytesToString(%q) = %q", s, got)
		}
	}
}

func TestEmptyAndNilCases(t *testing.T) {
	if b := StringToBytes(""); b != nil {
		t.Errorf("StringToBytes(\"\") = %v, want nil", b)
	}
	if s := BytesToString(nil); s != "" {
		t.Errorf("BytesToString(nil) = %q, want empty", s)
	}
	if s := BytesToString([]byte{}); s != "" {
		t.Errorf("BytesToString([]byte{}) = %q, want empty", s)
	}
}

// TestConversionsDoNotAllocate is the point of the task: a correct-but-copying
// implementation passes every test above and fails this one.
func TestConversionsDoNotAllocate(t *testing.T) {
	s := strings.Repeat("x", 4096)
	if n := testing.AllocsPerRun(100, func() { Sink = StringToBytes(s) }); n != 0 {
		t.Errorf("StringToBytes allocated %v times per call, want 0", n)
	}

	b := make([]byte, 4096)
	if n := testing.AllocsPerRun(100, func() { SinkString = BytesToString(b) }); n != 0 {
		t.Errorf("BytesToString allocated %v times per call, want 0", n)
	}
}

// TestSharingIsObservable proves the conversion really shares rather than
// copying quickly: mutating the slice changes the string.
//
// This deliberately performs the undefined behaviour the doc comment warns
// about, in the one place where it is the thing under test.
func TestSharingIsObservable(t *testing.T) {
	b := []byte("hello")
	s := BytesToString(b)
	b[0] = 'j'
	if s != "jello" {
		t.Errorf("s = %q after mutating the slice — BytesToString copied instead of sharing", s)
	}
}
