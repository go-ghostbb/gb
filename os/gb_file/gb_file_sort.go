package gbfile

import (
	gbarray "github.com/Ghostbb-io/gb/container/gb_array"
	"strings"
)

// fileSortFunc is the comparison function for files.
// It sorts the array in order of: directory -> file.
// If `path1` and `path2` are the same type, it then sorts them as strings.
func fileSortFunc(path1, path2 string) int {
	isDirPath1 := IsDir(path1)
	isDirPath2 := IsDir(path2)
	if isDirPath1 && !isDirPath2 {
		return -1
	}
	if !isDirPath1 && isDirPath2 {
		return 1
	}
	if n := strings.Compare(path1, path2); n != 0 {
		return n
	} else {
		return -1
	}
}

// SortFiles sorts the `files` in order of: directory -> file.
// Note that the item of `files` should be absolute path.
func SortFiles(files []string) []string {
	array := gbarray.NewSortedStrArrayComparator(fileSortFunc)
	array.Add(files...)
	return array.Slice()
}
