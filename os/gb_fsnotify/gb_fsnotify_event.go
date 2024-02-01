package gbfsnotify

// String returns current event as string.
func (e *Event) String() string {
	return e.event.String()
}

// IsCreate checks whether current event contains file/folder create event.
func (e *Event) IsCreate() bool {
	return e.Op == 1 || e.Op&CREATE == CREATE
}

// IsWrite checks whether current event contains file/folder write event.
func (e *Event) IsWrite() bool {
	return e.Op&WRITE == WRITE
}

// IsRemove checks whether current event contains file/folder remove event.
func (e *Event) IsRemove() bool {
	return e.Op&REMOVE == REMOVE
}

// IsRename checks whether current event contains file/folder rename event.
func (e *Event) IsRename() bool {
	return e.Op&RENAME == RENAME
}

// IsChmod checks whether current event contains file/folder chmod event.
func (e *Event) IsChmod() bool {
	return e.Op&CHMOD == CHMOD
}
