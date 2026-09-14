package seqlock

import (
	"sync"
	"testing"
	"time"
)

// TestNoTornReads runs one writer and four readers, and fails if any reader
// ever assembles a Point from two different writes.
//
// The writer only ever publishes points whose three fields are equal, so a
// Point with unequal fields was torn.
func TestNoTornReads(t *testing.T) {
	var s Seqlock
	s.Write(Point{})

	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for n := int64(1); ; n++ {
			select {
			case <-stop:
				return
			default:
				s.Write(Point{X: n, Y: n, Z: n})
			}
		}
	}()

	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					if p := s.Read(); p.X != p.Y || p.Y != p.Z {
						t.Error("a reader observed a Point assembled from two different writes")
						return
					}
				}
			}
		}()
	}

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
}

// TestReadSeesLatest checks that a read after a completed write sees it.
func TestReadSeesLatest(t *testing.T) {
	var s Seqlock
	s.Write(Point{X: 7, Y: 7, Z: 7})
	if got := s.Read(); got != (Point{X: 7, Y: 7, Z: 7}) {
		t.Fatal("Read did not return the published Point")
	}
}

// TestReadsAreConcurrent checks that readers do not exclude one another: two
// readers must be able to be inside Read at the same time. A seqlock whose
// readers take a mutex would pass the tests above and fail this one.
func TestReadsAreConcurrent(t *testing.T) {
	var s Seqlock
	s.Write(Point{X: 1, Y: 1, Z: 1})

	const readers = 4
	inside := make(chan struct{}, readers)
	release := make(chan struct{})
	var wg sync.WaitGroup

	// Each reader reports that it is about to read, then reads repeatedly
	// until released. If Read serialised readers, fewer than all of them could
	// report before any finished.
	for range readers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			inside <- struct{}{}
			for {
				select {
				case <-release:
					return
				default:
					_ = s.Read()
				}
			}
		}()
	}

	deadline := time.After(2 * time.Second)
	for range readers {
		select {
		case <-inside:
		case <-deadline:
			close(release)
			wg.Wait()
			t.Fatal("not every reader reached Read; readers appear to exclude one another")
		}
	}
	close(release)
	wg.Wait()
}
