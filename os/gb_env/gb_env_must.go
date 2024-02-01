package gbenv

// MustSet performs as Set, but it panics if any error occurs.
func MustSet(key, value string) {
	if err := Set(key, value); err != nil {
		panic(err)
	}
}

// MustRemove performs as Remove, but it panics if any error occurs.
func MustRemove(key ...string) {
	if err := Remove(key...); err != nil {
		panic(err)
	}
}
