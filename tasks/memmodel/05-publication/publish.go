// Package publish compares five ways of lazily building one shared object and
// handing it to other goroutines.
//
// Config is filled field by field. A reader that sees a non-nil pointer but
// unfilled fields has observed a partially constructed object.
package publish

import (
	"sync"
	"sync/atomic"
)

// Config is the object being published. NewConfig fills every field with the
// same marker, so any reader seeing a mix — or a zero — has caught a partial
// publication.
type Config struct{ A, B, C int }

// NewConfig builds a Config with every field set to marker.
func NewConfig(marker int) *Config {
	c := new(Config)
	c.A = marker
	c.B = marker
	c.C = marker
	return c
}

// Complete reports whether every field agrees and is non-zero.
func (c *Config) Complete() bool {
	return c != nil && c.A != 0 && c.A == c.B && c.B == c.C
}

// Getter is one lazy-initialisation strategy.
type Getter interface{ Get(marker int) *Config }

// PlainPointer publishes through an ordinary pointer field, with no
// synchronisation at all.
type PlainPointer struct{ p *Config }

// Get returns the shared Config, building it on first use.
func (s *PlainPointer) Get(marker int) *Config {
	if s.p == nil {
		s.p = NewConfig(marker)
	}
	return s.p
}

// AtomicPointer publishes through atomic.Pointer.
type AtomicPointer struct{ p atomic.Pointer[Config] }

// Get returns the shared Config, building it on first use.
func (s *AtomicPointer) Get(marker int) *Config {
	if c := s.p.Load(); c != nil {
		return c
	}
	c := NewConfig(marker)
	s.p.CompareAndSwap(nil, c)
	return s.p.Load()
}

// OnceGuarded publishes through sync.Once.
type OnceGuarded struct {
	once sync.Once
	p    *Config
}

// Get returns the shared Config, building it on first use.
func (s *OnceGuarded) Get(marker int) *Config {
	s.once.Do(func() { s.p = NewConfig(marker) })
	return s.p
}

// MutexGuarded holds the lock across both the check and the build.
type MutexGuarded struct {
	mu sync.Mutex
	p  *Config
}

// Get returns the shared Config, building it on first use.
func (s *MutexGuarded) Get(marker int) *Config {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.p == nil {
		s.p = NewConfig(marker)
	}
	return s.p
}

// DoubleChecked reads the pointer without the lock first, then takes the lock.
type DoubleChecked struct {
	mu sync.Mutex
	p  *Config
}

// Get returns the shared Config, building it on first use.
func (s *DoubleChecked) Get(marker int) *Config {
	if s.p != nil {
		return s.p
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.p == nil {
		s.p = NewConfig(marker)
	}
	return s.p
}
