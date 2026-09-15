// Package poolclearing measures how long an object put into a sync.Pool
// survives collection.
package poolclearing

import (
	"runtime"
	"sync"
)

type buf struct{ b []byte }

// Survival records what one pool did across forced collections.
type Survival struct {
	NewCallsBeforeAnyGC int
	SurvivesOneGC       bool
	SurvivesTwoGC       bool
	NewCallsAfterTwoGC  int
}

// Measure puts one object into a pool and forces collections around it.
//
// GOMAXPROCS is pinned to 1 for the duration: a Pool's storage is per-P, so a
// goroutine rescheduled onto a different P between the Put and the Get would
// miss the object for reasons that have nothing to do with the collector.
//
// Constructions are counted rather than pointers compared — a pointer identity
// check would invite the allocator and the compiler into an answer that is
// supposed to be about sync.Pool.
func Measure() Survival {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))

	var s Survival
	news := 0
	p := &sync.Pool{New: func() any { news++; return &buf{b: make([]byte, 64)} }}

	x := p.Get().(*buf)
	s.NewCallsBeforeAnyGC = news
	p.Put(x)

	runtime.GC()
	before := news
	y := p.Get().(*buf)
	s.SurvivesOneGC = news == before
	p.Put(y)

	runtime.GC()
	runtime.GC()
	before = news
	_ = p.Get().(*buf)
	s.SurvivesTwoGC = news == before
	s.NewCallsAfterTwoGC = news - before

	return s
}
