// Package drill — C2/02 slice-alias.
package drill

// Split returns the 4-byte header and the remaining payload of b.
func Split(b []byte) (head, tail []byte) {
	head = b[:4]
	tail = b[4:]
	return head, tail
}

// Frame splits b, then appends a checksum byte to the header.
func Frame(b []byte, checksum byte) (head, tail []byte) {
	head, tail = Split(b)
	head = append(head, checksum)
	return head, tail
}

// CopyHeader returns an independent copy of the first 4 bytes of b.
func CopyHeader(b []byte) []byte {
	h := make([]byte, 4)
	copy(h, b[:4])
	return h
}
