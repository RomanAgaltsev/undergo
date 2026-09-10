// Package drill — C13/09 token-bucket.
package drill

import "time"

// Grant is an access grant with an expiry.
type Grant struct {
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// NewGrant issues a grant valid for 30 days from issuedAt.
func NewGrant(issuedAt time.Time) Grant {
	return Grant{
		IssuedAt:  issuedAt,
		ExpiresAt: issuedAt.Add(30 * 24 * time.Hour),
	}
}

// Expired reports whether the grant has expired as of now.
func (g Grant) Expired(now time.Time) bool {
	return now.Equal(g.ExpiresAt)
}

// Refill runs fn every interval to top up a bucket, until stop is signaled.
func Refill(interval time.Duration, fn func(), stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case <-time.Tick(interval):
			fn()
		}
	}
}
