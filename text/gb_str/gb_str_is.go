package gbstr

import "ghostbb.io/internal/utils"

// IsNumeric tests whether the given string s is numeric.
func IsNumeric(s string) bool {
	return utils.IsNumeric(s)
}
