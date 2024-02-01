// Package fileinfo provides virtual os.FileInfo for given information.
package fileinfo

import (
	"os"
	"time"
)

type Info struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func New(name string, size int64, mode os.FileMode, modTime time.Time) *Info {
	return &Info{
		name:    name,
		size:    size,
		mode:    mode,
		modTime: modTime,
	}
}

func (i *Info) Name() string {
	return i.name
}

func (i *Info) Size() int64 {
	return i.size
}

func (i *Info) IsDir() bool {
	return i.mode.IsDir()
}

func (i *Info) Mode() os.FileMode {
	return i.mode
}

func (i *Info) ModTime() time.Time {
	return i.modTime
}

func (i *Info) Sys() interface{} {
	return nil
}
