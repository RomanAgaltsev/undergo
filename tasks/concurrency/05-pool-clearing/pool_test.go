package poolclearing

import (
	"runtime"
	"sync"
	"testing"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

type buf struct{ b []byte }

// The pool's storage is per-P, so the whole measurement is pinned to one P.
// Constructions are counted rather than pointers compared: a pointer identity
// check would invite the allocator and the compiler into an answer that is
// supposed to be about sync.Pool.
func TestPredictions(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	got := map[string]any{}

	news := 0
	p := &sync.Pool{New: func() any { news++; return &buf{b: make([]byte, 64)} }}

	x := p.Get().(*buf)
	got["new_calls_before_any_gc"] = news
	p.Put(x)

	runtime.GC()
	before := news
	y := p.Get().(*buf)
	got["survives_one_gc"] = news == before
	p.Put(y)

	runtime.GC()
	runtime.GC()
	before = news
	_ = p.Get().(*buf)
	got["survives_two_gc"] = news == before
	got["new_calls_after_two_gc"] = news - before

	predict.Check(t, got)
}
