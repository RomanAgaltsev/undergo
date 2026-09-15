//go:build amd64

package abi

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions reads the compiler's own listing and grades what it shows.
func TestPredictions(t *testing.T) {
	listing, err := Listing()
	if err != nil {
		t.Fatalf("compiling for a listing: %v", err)
	}
	if bodyOf(listing, "Nine") == "" {
		t.Fatal("the listing does not contain Nine; the symbol match is wrong")
	}

	predict.Check(t, map[string]any{
		"nine_reads_stack":  ReadsArgumentFromStack(listing, "Nine"),
		"ten_reads_stack":   ReadsArgumentFromStack(listing, "Ten"),
		"pair_reads_stack":  ReadsArgumentFromStack(listing, "SumPair"),
		"mixed_uses_floats": UsesFloatRegisters(listing, "Mixed"),
	})
}
