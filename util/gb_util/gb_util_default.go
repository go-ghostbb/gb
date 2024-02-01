package gbutil

// GetOrDefaultStr checks and returns value according whether parameter `param` available.
// It returns `param[0]` if it is available, or else it returns `def`.
func GetOrDefaultStr(def string, param ...string) string {
	value := def
	if len(param) > 0 && param[0] != "" {
		value = param[0]
	}
	return value
}

// GetOrDefaultAny checks and returns value according whether parameter `param` available.
// It returns `param[0]` if it is available, or else it returns `def`.
func GetOrDefaultAny(def interface{}, param ...interface{}) interface{} {
	value := def
	if len(param) > 0 && param[0] != "" {
		value = param[0]
	}
	return value
}
