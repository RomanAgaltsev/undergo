// Package traps is a predict task about GODEBUG=checkfinalizers=1.
//
// There is nothing to read in this file. The four programs you are predicting
// live in testdata/, because each one has to run as its own process: when the
// runtime does diagnose a problem it does not warn, it throws.
//
// Read them in the README, or open testdata/trap1.go.txt through trap4.go.txt.
package traps
