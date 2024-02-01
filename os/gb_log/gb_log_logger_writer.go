package gblog

import (
	"bytes"
	"context"
)

// Write implements the io.Writer interface.
// It just prints the content using Print.
func (l *Logger) Write(p []byte) (n int, err error) {
	l.Header(false).Print(context.TODO(), string(bytes.TrimRight(p, "\r\n")))
	return len(p), nil
}
