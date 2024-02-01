package gbjson

// SetSplitChar sets the separator char for hierarchical data access.
func (j *Json) SetSplitChar(char byte) {
	j.mu.Lock()
	j.c = char
	j.mu.Unlock()
}

// SetViolenceCheck enables/disables violence check for hierarchical data access.
func (j *Json) SetViolenceCheck(enabled bool) {
	j.mu.Lock()
	j.vc = enabled
	j.mu.Unlock()
}
