// Package drill — C8/06 timing-compare.
package drill

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Verify reports whether signature is a valid HMAC-SHA256 of message under key.
func Verify(message, signature string, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))
	return expected == signature
}
