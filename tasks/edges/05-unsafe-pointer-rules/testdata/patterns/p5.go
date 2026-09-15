//go:build ignore

// Pattern 4: conversion of a slice's data pointer, used immediately.
package patterns

import "unsafe"

func P5(b []byte) string { return unsafe.String(&b[0], len(b)) }
