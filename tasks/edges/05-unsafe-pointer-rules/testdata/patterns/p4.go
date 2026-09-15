//go:build ignore

// Pattern 6: unsafe.Slice over a pointer the caller owns.
package patterns

import "unsafe"

func P4(p *byte, n int) []byte { return unsafe.Slice(p, n) }
