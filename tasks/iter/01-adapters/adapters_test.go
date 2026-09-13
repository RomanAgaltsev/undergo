package adapters

import (
	"iter"
	"slices"
	"testing"
)

// source counts how many elements it actually produced, so a test can assert
// that it was not over-driven after its consumer stopped.
type source struct{ produced int }

// seq returns 0..n-1, counting each element as it goes.
func (s *source) seq(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			s.produced++
			if !yield(i) {
				return
			}
		}
	}
}

// naturals is an infinite sequence. Anything that drains it hangs.
func naturals() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func TestMapTransforms(t *testing.T) {
	src := (&source{}).seq(4)
	got := slices.Collect(Map(src, func(v int) int { return v * 10 }))
	if want := []int{0, 10, 20, 30}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFilterKeeps(t *testing.T) {
	src := (&source{}).seq(6)
	got := slices.Collect(Filter(src, func(v int) bool { return v%2 == 0 }))
	if want := []int{0, 2, 4}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTakeTakes(t *testing.T) {
	src := (&source{}).seq(10)
	got := slices.Collect(Take(src, 3))
	if want := []int{0, 1, 2}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestTakeStopsTheSource is why Take cannot filter after the fact: this hangs
// forever if Take drains its source.
func TestTakeStopsTheSource(t *testing.T) {
	got := slices.Collect(Take(naturals(), 3))
	if want := []int{0, 1, 2}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTakeNonPositiveNeverTouchesTheSource(t *testing.T) {
	for _, n := range []int{0, -1} {
		produced := &source{}
		src := produced.seq(5)
		if got := slices.Collect(Take(src, n)); len(got) != 0 {
			t.Errorf("Take(seq, %d) yielded %v, want nothing", n, got)
		}
		if produced.produced != 0 {
			t.Errorf("Take(seq, %d) drove the source %d times, want 0", n, produced.produced)
		}
	}
}

// TestAdaptersStopWhenTheConsumerBreaks is the yield-contract test: each adapter
// must return false up the chain rather than running its source to the end.
func TestAdaptersStopWhenTheConsumerBreaks(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		produced := &source{}
		src := produced.seq(100)
		for range Map(src, func(v int) int { return v }) {
			break
		}
		if produced.produced > 1 {
			t.Errorf("source produced %d elements after one break, want 1", produced.produced)
		}
	})

	t.Run("Filter", func(t *testing.T) {
		produced := &source{}
		src := produced.seq(100)
		for range Filter(src, func(int) bool { return true }) {
			break
		}
		if produced.produced > 1 {
			t.Errorf("source produced %d elements after one break, want 1", produced.produced)
		}
	})

	t.Run("Chunk", func(t *testing.T) {
		produced := &source{}
		src := produced.seq(100)
		for range Chunk(src, 4) {
			break
		}
		if produced.produced > 4 {
			t.Errorf("source produced %d elements after one break, want at most 4", produced.produced)
		}
	})
}

func TestZipPairs(t *testing.T) {
	a := (&source{}).seq(3)
	var b iter.Seq[string] = func(yield func(string) bool) {
		for _, s := range []string{"x", "y", "z"} {
			if !yield(s) {
				return
			}
		}
	}

	var keys []int
	var vals []string
	for k, v := range Zip(a, b) {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	if want := []int{0, 1, 2}; !slices.Equal(keys, want) {
		t.Errorf("keys = %v, want %v", keys, want)
	}
	if want := []string{"x", "y", "z"}; !slices.Equal(vals, want) {
		t.Errorf("values = %v, want %v", vals, want)
	}
}

func TestZipStopsAtTheShorter(t *testing.T) {
	short := (&source{}).seq(2)
	longProduced := &source{}
	long := longProduced.seq(100)

	pairs := 0
	for range Zip(short, long) {
		pairs++
	}
	if pairs != 2 {
		t.Errorf("got %d pairs, want 2", pairs)
	}
	if longProduced.produced > 3 {
		t.Errorf("the longer source produced %d elements, want at most 3", longProduced.produced)
	}
}

func TestChunkEmitsAShortFinalChunk(t *testing.T) {
	src := (&source{}).seq(7)

	var sizes []int
	for c := range Chunk(src, 3) {
		sizes = append(sizes, len(c))
	}
	if want := []int{3, 3, 1}; !slices.Equal(sizes, want) {
		t.Errorf("chunk sizes = %v, want %v", sizes, want)
	}
}

// TestChunkDoesNotAliasItsBuffer catches the optimisation of reusing one slice:
// every other chunk test passes, and the collected chunks all hold the last
// chunk's contents.
func TestChunkDoesNotAliasItsBuffer(t *testing.T) {
	src := (&source{}).seq(6)

	chunks := slices.Collect(Chunk(src, 2))
	if len(chunks) != 3 {
		t.Fatalf("got %d chunks, want 3", len(chunks))
	}
	want := [][]int{{0, 1}, {2, 3}, {4, 5}}
	for i := range want {
		if !slices.Equal(chunks[i], want[i]) {
			t.Errorf("chunk %d = %v, want %v — are the chunks sharing a buffer?", i, chunks[i], want[i])
		}
	}
}

func TestChunkNonPositive(t *testing.T) {
	src := (&source{}).seq(5)
	if got := slices.Collect(Chunk(src, 0)); len(got) != 0 {
		t.Errorf("Chunk(seq, 0) yielded %v, want nothing", got)
	}
}
