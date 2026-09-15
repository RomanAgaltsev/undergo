//go:build ignore

// Pattern 1: conversion of *T1 to *T2 via unsafe.Pointer, same layout.
package patterns

import "unsafe"

func P1(f *float64) *uint64 { return (*uint64)(unsafe.Pointer(f)) }
