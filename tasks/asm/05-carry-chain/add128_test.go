//go:build amd64

package add128

import (
	"math/bits"
	"math/rand/v2"
	"testing"
)

// reference is what the answer must equal, bit for bit.
func reference(loA, hiA, loB, hiB uint64) (uint64, uint64) {
	lo, carry := bits.Add64(loA, loB, 0)
	hi, _ := bits.Add64(hiA, hiB, carry)
	return lo, hi
}

// checkOne compares one input against the reference. It never prints the
// expected value, so a failing run does not hand over the answer.
func checkOne(t *testing.T, loA, hiA, loB, hiB uint64) {
	t.Helper()
	wantLo, wantHi := reference(loA, hiA, loB, hiB)
	gotLo, gotHi := Add128(loA, hiA, loB, hiB)
	if gotLo != wantLo || gotHi != wantHi {
		t.Fatalf("Add128 disagrees with math/bits for lo=%#x hi=%#x + lo=%#x hi=%#x",
			loA, hiA, loB, hiB)
	}
}

// TestCarryPropagates is the case the stub gets wrong, named so the failure is
// legible before the random inputs start.
func TestCarryPropagates(t *testing.T) {
	lo, hi := Add128(^uint64(0), 0, 1, 0)
	if lo != 0 || hi != 1 {
		t.Fatal("a carry out of the low word did not reach the high word")
	}
}

// TestEdges covers the boundaries where a carry does and does not happen.
func TestEdges(t *testing.T) {
	edges := []uint64{0, 1, ^uint64(0), ^uint64(0) - 1, 1 << 63}
	for _, a := range edges {
		for _, b := range edges {
			checkOne(t, a, a, b, b)
			checkOne(t, a, b, b, a)
		}
	}
}

// TestAgainstReference runs ten thousand random inputs from a fixed seed, so a
// failure is reproducible.
func TestAgainstReference(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 10000 {
		checkOne(t, r.Uint64(), r.Uint64(), r.Uint64(), r.Uint64())
	}
}
