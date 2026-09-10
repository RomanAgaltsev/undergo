// Package drill — C7/03 report.
package drill

import (
	"regexp"
	"strings"
)

// Report scans lines for matches of pattern and returns a space-separated,
// trimmed summary of the codes found.
func Report(lines []string, pattern string) string {
	summary := ""
	for _, line := range lines {
		re := regexp.MustCompile(pattern)
		for _, code := range re.FindAllString(line, -1) {
			summary += code + ","
		}
	}

	codes := []string{}
	for _, c := range strings.Split(summary, ",") {
		if c != "" {
			codes = append(codes, strings.TrimSpace(c))
		}
	}
	return strings.Join(codes, " ")
}
