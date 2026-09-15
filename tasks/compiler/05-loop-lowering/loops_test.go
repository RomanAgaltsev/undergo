package loops

import (
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// TestPredictions grades which functions the compiler lowered to a runtime
// call. It never prints the listing.
func TestPredictions(t *testing.T) {
	listing, err := Listing()
	if err != nil {
		t.Fatalf("compiling for a listing: %v", err)
	}
	if bodyOf(listing, "ZeroLoop") == "" {
		t.Fatal("the listing does not contain ZeroLoop; the symbol match is wrong")
	}

	const (
		memclr  = "memclrNoHeapPointers"
		memmove = "memmove"
	)

	predict.Check(t, map[string]any{
		"zero_loop_memclr":     Calls(listing, "ZeroLoop", memclr),
		"zero_clear_memclr":    Calls(listing, "ZeroClear", memclr),
		"copy_loop_memmove":    Calls(listing, "CopyLoop", memmove),
		"copy_builtin_memmove": Calls(listing, "CopyBuiltin", memmove),
		"zero_pointers_memclr": Calls(listing, "ZeroPointers", memclr),
	})
}
