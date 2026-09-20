// Package clocks is about the second clock inside a time.Time.
//
// time.Now returns a Time carrying two readings: a wall clock, which can be
// adjusted by NTP or by hand, and a monotonic reading, which cannot go
// backwards. Some operations carry the monotonic reading forward and some drop
// it, and two Times that name the same instant are not necessarily equal.
package clocks

import (
	"encoding/json"
	"strings"
	"time"
)

// Monotonic reports whether t still carries a monotonic reading. Time.String
// appends " m=±<seconds>" when one is present, which is the only thing in the
// standard library that will tell you.
func Monotonic(t time.Time) bool {
	return strings.Contains(t.String(), " m=")
}

// Reference is the instant every slot is asked about.
func Reference() time.Time { return time.Now() }

// JSONRoundTrip marshals t and unmarshals the result.
func JSONRoundTrip(t time.Time) (time.Time, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return time.Time{}, err
	}
	var out time.Time
	if err := json.Unmarshal(b, &out); err != nil {
		return time.Time{}, err
	}
	return out, nil
}

// Stripped returns t with its monotonic reading removed, which is what
// Round(0) is for.
func Stripped(t time.Time) time.Time { return t.Round(0) }
