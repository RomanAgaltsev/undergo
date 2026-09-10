// Package drill — C6/03 buffer-api.
package drill

import (
	"io"
	"os"
)

// Buffer accumulates bytes and tracks how many were written.
type Buffer struct {
	Data []byte
	size int
}

// Append adds p to the buffer.
func (b *Buffer) Append(p []byte) {
	b.Data = append(b.Data, p...)
	b.size += len(p)
}

// Len returns the number of bytes written.
func (b Buffer) Len() int {
	return b.size
}

// Open reads all of f into a new Buffer.
func Open(f *os.File) (*Buffer, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return &Buffer{Data: data, size: len(data)}, nil
}
