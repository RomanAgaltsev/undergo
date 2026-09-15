//go:build ignore

// Pattern 2: conversion of unsafe.Pointer to uintptr, but NOT back again. The
// value is only ever printed, so nothing depends on it still being an address.
package patterns

import (
	"fmt"
	"unsafe"
)

func P3(p *int) string {
	return fmt.Sprintf("%#x", uintptr(unsafe.Pointer(p)))
}
