package drill

import (
	"math/rand/v2"
	"testing"
)

// TestAgainstModel drives the store and a plain map with the same ten thousand
// random operations and checks that they agree after every step.
func TestAgainstModel(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	s := New(64)
	model := map[int]int{}

	for step := range 10_000 {
		k, v := rng.IntN(100), rng.Int()
		switch rng.IntN(3) {
		case 0, 1:
			_ = s.Put(k, v)
			model[k] = v
		case 2:
			s.Delete(k)
			delete(model, k)
		}

		for k, want := range model {
			if got, ok := s.Get(k); ok && got != want {
				t.Fatalf("step %d: Get(%d) = %d, want %d", step, k, got, want)
			}
		}
		if s.Len() > len(model) {
			t.Fatalf("step %d: store holds %d keys, model %d", step, s.Len(), len(model))
		}
	}
}
