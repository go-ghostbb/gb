package gbstr

import "strings"

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func Compare(a, b string) int {
	return strings.Compare(a, b)
}

// Equal reports whether `a` and `b`, interpreted as UTF-8 strings,
// are equal under Unicode case-folding, case-insensitively.
func Equal(a, b string) bool {
	return strings.EqualFold(a, b)
}
