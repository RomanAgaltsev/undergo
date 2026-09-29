// Package scan asks where a pointer can be kept so that the collector sees it.
//
// A pointer keeps its target alive only if the garbage collector finds it, and
// the collector finds a pointer only where it expects one to be. Each exported
// function stores the only reference to a fresh Payload somewhere, forces
// collections, and reports whether the Payload survived.
package scan

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"unsafe"
	"weak"
)

// Payload is an ordinary heap object with no pointers of its own.
type Payload struct{ data [64]byte }

// record is laid over the byte arena: a tag, then a pointer.
type record struct {
	tag uint64
	p   *Payload
}

var (
	byteArena = make([]byte, 64)
	ptrArena  = make([]unsafe.Pointer, 4)
	wordArena = make([]uintptr, 4)
	anyArena  = make([]any, 4)
)

// survives stores the only reference to a new Payload with store, collects
// twice, and reports whether the Payload is still alive.
//
// The weak pointer is the witness: it does not keep the Payload alive, and its
// Value is nil once the collector has reclaimed it. Two collections, so that a
// survivor has survived a whole cycle that began after the store.
func survives(store func(*Payload)) bool {
	w := func() weak.Pointer[Payload] {
		p := &Payload{}
		store(p)
		return weak.Make(p)
	}()
	runtime.GC()
	runtime.GC()
	return w.Value() != nil
}

// InBytes writes the pointer's bits into a []byte.
func InBytes() bool {
	return survives(func(p *Payload) {
		*(*unsafe.Pointer)(unsafe.Pointer(&byteArena[8])) = unsafe.Pointer(p)
	})
}

// ThroughView writes the pointer into the pointer-typed field of a record laid
// over the same []byte.
func ThroughView() bool {
	return survives(func(p *Payload) {
		r := (*record)(unsafe.Pointer(&byteArena[16]))
		r.p = p
	})
}

// InPointers stores it in a []unsafe.Pointer.
func InPointers() bool {
	return survives(func(p *Payload) { ptrArena[1] = unsafe.Pointer(p) })
}

// InWords stores its address in a []uintptr.
func InWords() bool {
	return survives(func(p *Payload) { wordArena[1] = uintptr(unsafe.Pointer(p)) })
}

// InAnys stores it in a []any.
func InAnys() bool {
	return survives(func(p *Payload) { anyArena[1] = p })
}

// VetFlags reports whether `go vet` finds anything to object to in this
// package, the stores above included.
func VetFlags() (bool, error) {
	cmd := exec.Command("go", "vet", ".")
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return false, nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) && len(out) > 0 {
		return true, nil
	}
	return false, err
}
