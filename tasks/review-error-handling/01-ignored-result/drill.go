// Package drill — C3/01 ignored-result.
package drill

import "strconv"

// Total parses each field as an integer and returns the sum.
func Total(fields []string) int {
	sum := 0
	for _, f := range fields {
		n, _ := strconv.Atoi(f)
		sum += n
	}
	return sum
}

// FirstValid returns the first field that parses as an integer.
func FirstValid(fields []string) (int, error) {
	for _, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil {
			continue
		}
		return n, nil
	}
	return 0, strconv.ErrSyntax
}
