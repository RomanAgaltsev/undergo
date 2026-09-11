// Package drill — C8/02 weak-token.
package drill

import (
	"encoding/hex"
	mrand "math/rand"
)

// Token returns a 16-byte session token as a hex string.
func Token() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(mrand.Intn(256))
	}
	return hex.EncodeToString(b)
}
