package gbstr

import "strings"

// Repeat returns a new string consisting of multiplier copies of the string input.
//
// Example:
// Repeat("a", 3) -> "aaa"
func Repeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}
