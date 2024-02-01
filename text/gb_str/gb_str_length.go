package gbstr

import "unicode/utf8"

// LenRune returns string length of unicode.
func LenRune(str string) int {
	return utf8.RuneCountInString(str)
}
