// Package timerbuf asks what time.After's channel used to be.
//
// Before Go 1.23 the channel time.After returns had a one-element buffer, so
// the runtime could deliver the time without a waiting receiver. From 1.23 it
// is unbuffered, and asynctimerchan=1 restored the old behaviour for modules
// that needed it. Go 1.27 removed that setting.
package timerbuf

import "github.com/RomanAgaltsev/undergo/internal/toolchain"

// Program prints the capacity of the channel time.After returns.
const Program = `package main

import (
	"fmt"
	"time"
)

func main() { fmt.Print(cap(time.After(time.Hour))) }
`

// Local is the toolchain running this test.
const Local = toolchain.Local

// Cap runs Program under the named toolchain at the given go line and returns
// the printed capacity.
func Cap(name, goLine string) (string, error) {
	res, err := toolchain.Run(name, goLine, Program)
	if err != nil {
		return "", err
	}
	return res.Output, nil
}
