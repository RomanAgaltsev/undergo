// Package drill — C6/05 bool-params.
package drill

// Connect opens a connection configured by the given flags.
func Connect(tls, retry, verbose bool) string {
	mode := ""
	if tls {
		mode += "tls "
	}
	if retry {
		mode += "retry "
	}
	if verbose {
		mode += "verbose"
	}
	return mode
}

// Enable toggles a single feature.
func Enable(on bool) bool {
	return on
}
