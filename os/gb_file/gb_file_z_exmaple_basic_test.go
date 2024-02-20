package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
	"os"
)

func ExampleMkdir() {
	// init
	var (
		path = gbfile.Temp("gbfile_example_basic_dir")
	)

	// Creates directory
	gbfile.Mkdir(path)

	// Check if directory exists
	fmt.Println(gbfile.IsDir(path))

	// Output:
	// true
}

func ExampleCreate() {
	// init
	var (
		path     = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
		dataByte = make([]byte, 50)
	)
	// Check whether the file exists
	isFile := gbfile.IsFile(path)

	fmt.Println(isFile)

	// Creates file with given `path` recursively
	fileHandle, _ := gbfile.Create(path)
	defer fileHandle.Close()

	// Write some content to file
	n, _ := fileHandle.WriteString("hello ghostbb")

	// Check whether the file exists
	isFile = gbfile.IsFile(path)

	fmt.Println(isFile)

	// Reads len(b) bytes from the File
	fileHandle.ReadAt(dataByte, 0)

	fmt.Println(string(dataByte[:n]))

	// Output:
	// false
	// true
	// hello ghostbb
}

func ExampleOpen() {
	// init
	var (
		path     = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
		dataByte = make([]byte, 4096)
	)
	// Open file or directory with READONLY model
	file, _ := gbfile.Open(path)
	defer file.Close()

	// Read data
	n, _ := file.Read(dataByte)

	fmt.Println(string(dataByte[:n]))

	// Output:
	// hello ghostbb
}

func ExampleOpenFile() {
	// init
	var (
		path     = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
		dataByte = make([]byte, 4096)
	)
	// Opens file/directory with custom `flag` and `perm`
	// Create if file does not exist,it is created in a readable and writable mode,prem 0777
	openFile, _ := gbfile.OpenFile(path, os.O_CREATE|os.O_RDWR, gbfile.DefaultPermCopy)
	defer openFile.Close()

	// Write some content to file
	writeLength, _ := openFile.WriteString("hello ghostbb test open file")

	fmt.Println(writeLength)

	// Read data
	n, _ := openFile.ReadAt(dataByte, 0)

	fmt.Println(string(dataByte[:n]))

	// Output:
	// 28
	// hello ghostbb test open file
}

func ExampleOpenWithFlag() {
	// init
	var (
		path     = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
		dataByte = make([]byte, 4096)
	)

	// Opens file/directory with custom `flag`
	// Create if file does not exist,it is created in a readable and writable mode with default `perm` is 0666
	openFile, _ := gbfile.OpenWithFlag(path, os.O_CREATE|os.O_RDWR)
	defer openFile.Close()

	// Write some content to file
	writeLength, _ := openFile.WriteString("hello ghostbb test open file with flag")

	fmt.Println(writeLength)

	// Read data
	n, _ := openFile.ReadAt(dataByte, 0)

	fmt.Println(string(dataByte[:n]))

	// Output:
	// 38
	// hello ghostbb test open file with flag
}

func ExampleJoin() {
	// init
	var (
		dirPath  = gbfile.Temp("gbfile_example_basic_dir")
		filePath = "file1"
	)

	// Joins string array paths with file separator of current system.
	joinString := gbfile.Join(dirPath, filePath)

	fmt.Println(joinString)

	// May Output:
	// /tmp/gbfile_example_basic_dir/file1
}

func ExampleExists() {
	// init
	var (
		path = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)
	// Checks whether given `path` exist.
	joinString := gbfile.Exists(path)

	fmt.Println(joinString)

	// Output:
	// true
}

func ExampleIsDir() {
	// init
	var (
		path     = gbfile.Temp("gbfile_example_basic_dir")
		filePath = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)
	// Checks whether given `path` a directory.
	fmt.Println(gbfile.IsDir(path))
	fmt.Println(gbfile.IsDir(filePath))

	// Output:
	// true
	// false
}

func ExamplePwd() {
	// Get absolute path of current working directory.
	fmt.Println(gbfile.Pwd())

	// May Output:
	// xxx/gf/os/gbfile
}

func ExampleChdir() {
	// init
	var (
		path = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)
	// Get current working directory
	fmt.Println(gbfile.Pwd())

	// Changes the current working directory to the named directory.
	gbfile.Chdir(path)

	// Get current working directory
	fmt.Println(gbfile.Pwd())

	// May Output:
	// xxx/gf/os/gbfile
	// /tmp/gbfile_example_basic_dir/file1
}

func ExampleIsFile() {
	// init
	var (
		filePath = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
		dirPath  = gbfile.Temp("gbfile_example_basic_dir")
	)
	// Checks whether given `path` a file, which means it's not a directory.
	fmt.Println(gbfile.IsFile(filePath))
	fmt.Println(gbfile.IsFile(dirPath))

	// Output:
	// true
	// false
}

func ExampleStat() {
	// init
	var (
		path = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)
	// Get a FileInfo describing the named file.
	stat, _ := gbfile.Stat(path)

	fmt.Println(stat.Name())
	fmt.Println(stat.IsDir())
	fmt.Println(stat.Mode())
	fmt.Println(stat.ModTime())
	fmt.Println(stat.Size())
	fmt.Println(stat.Sys())

	// May Output:
	// file1
	// false
	// -rwxr-xr-x
	// 2021-12-02 11:01:27.261441694 +0800 CST
	// &{16777220 33261 1 8597857090 501 20 0 [0 0 0 0] {1638414088 192363490} {1638414087 261441694} {1638414087 261441694} {1638413480 485068275} 38 8 4096 0 0 0 [0 0]}
}

func ExampleMove() {
	// init
	var (
		srcPath = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
		dstPath = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file2")
	)
	// Check is file
	fmt.Println(gbfile.IsFile(dstPath))

	//  Moves `src` to `dst` path.
	// If `dst` already exists and is not a directory, it'll be replaced.
	gbfile.Move(srcPath, dstPath)

	fmt.Println(gbfile.IsFile(srcPath))
	fmt.Println(gbfile.IsFile(dstPath))

	// Output:
	// false
	// false
	// true
}

func ExampleRename() {
	// init
	var (
		srcPath = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file2")
		dstPath = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)
	// Check is file
	fmt.Println(gbfile.IsFile(dstPath))

	//  renames (moves) `src` to `dst` path.
	// If `dst` already exists and is not a directory, it'll be replaced.
	gbfile.Rename(srcPath, dstPath)

	fmt.Println(gbfile.IsFile(srcPath))
	fmt.Println(gbfile.IsFile(dstPath))

	// Output:
	// false
	// false
	// true
}

func ExampleDirNames() {
	// init
	var (
		path = gbfile.Temp("gbfile_example_basic_dir")
	)
	// Get sub-file names of given directory `path`.
	dirNames, _ := gbfile.DirNames(path)

	fmt.Println(dirNames)

	// May Output:
	// [file1]
}

func ExampleGlob() {
	// init
	var (
		path = gbfile.Pwd() + gbfile.Separator + "*_example_basic_test.go"
	)
	// Get sub-file names of given directory `path`.
	// Only show file name
	matchNames, _ := gbfile.Glob(path, true)

	fmt.Println(matchNames)

	// Show full path of the file
	matchNames, _ = gbfile.Glob(path, false)

	fmt.Println(matchNames)

	// May Output:
	// [gbfile_z_example_basic_test.go]
	// [xxx/gf/os/gbfile/gbfile_z_example_basic_test.go]
}

func ExampleIsReadable() {
	// init
	var (
		path = gbfile.Pwd() + gbfile.Separator + "testdata/readline/file.log"
	)

	// Checks whether given `path` is readable.
	fmt.Println(gbfile.IsReadable(path))

	// Output:
	// true
}

func ExampleIsWritable() {
	// init
	var (
		path = gbfile.Pwd() + gbfile.Separator + "testdata/readline/"
		file = "file.log"
	)

	// Checks whether given `path` is writable.
	fmt.Println(gbfile.IsWritable(path))
	fmt.Println(gbfile.IsWritable(path + file))

	// Output:
	// true
	// true
}

func ExampleChmod() {
	// init
	var (
		path = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)

	// Get a FileInfo describing the named file.
	stat, err := gbfile.Stat(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Show original mode
	fmt.Println(stat.Mode())

	// Change file model
	gbfile.Chmod(path, gbfile.DefaultPermCopy)

	// Get a FileInfo describing the named file.
	stat, _ = gbfile.Stat(path)
	// Show the modified mode
	fmt.Println(stat.Mode())

	// Output:
	// -rw-r--r--
	// -rwxr-xr-x
}

func ExampleAbs() {
	// init
	var (
		path = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)

	// Get an absolute representation of path.
	fmt.Println(gbfile.Abs(path))

	// May Output:
	// /tmp/gbfile_example_basic_dir/file1
}

func ExampleRealPath() {
	// init
	var (
		realPath  = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
		worryPath = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "worryFile")
	)

	// fetch an absolute representation of path.
	fmt.Println(gbfile.RealPath(realPath))
	fmt.Println(gbfile.RealPath(worryPath))

	// May Output:
	// /tmp/gbfile_example_basic_dir/file1
}

func ExampleSelfPath() {

	// Get absolute file path of current running process
	fmt.Println(gbfile.SelfPath())

	// May Output:
	// xxx/___github_com_gogf_gf_v2_os_gbfile__ExampleSelfPath
}

func ExampleSelfName() {

	// Get file name of current running process
	fmt.Println(gbfile.SelfName())

	// May Output:
	// ___github_com_gogf_gf_v2_os_gbfile__ExampleSelfName
}

func ExampleSelfDir() {

	// Get absolute directory path of current running process
	fmt.Println(gbfile.SelfDir())

	// May Output:
	// /private/var/folders/p6/gc_9mm3j229c0mjrjp01gqn80000gn/T
}

func ExampleBasename() {
	// init
	var (
		path = gbfile.Pwd() + gbfile.Separator + "testdata/readline/file.log"
	)

	// Get the last element of path, which contains file extension.
	fmt.Println(gbfile.Basename(path))

	// Output:
	// file.log
}

func ExampleName() {
	// init
	var (
		path = gbfile.Pwd() + gbfile.Separator + "testdata/readline/file.log"
	)

	// Get the last element of path without file extension.
	fmt.Println(gbfile.Name(path))

	// Output:
	// file
}

func ExampleDir() {
	// init
	var (
		path = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)

	// Get all but the last element of path, typically the path's directory.
	fmt.Println(gbfile.Dir(path))

	// May Output:
	// /tmp/gbfile_example_basic_dir
}

func ExampleIsEmpty() {
	// init
	var (
		path = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)

	// Check whether the `path` is empty
	fmt.Println(gbfile.IsEmpty(path))

	// Truncate file
	gbfile.Truncate(path, 0)

	// Check whether the `path` is empty
	fmt.Println(gbfile.IsEmpty(path))

	// Output:
	// false
	// true
}

func ExampleExt() {
	// init
	var (
		path = gbfile.Pwd() + gbfile.Separator + "testdata/readline/file.log"
	)

	// Get the file name extension used by path.
	fmt.Println(gbfile.Ext(path))

	// Output:
	// .log
}

func ExampleExtName() {
	// init
	var (
		path = gbfile.Pwd() + gbfile.Separator + "testdata/readline/file.log"
	)

	// Get the file name extension used by path but the result does not contains symbol '.'.
	fmt.Println(gbfile.ExtName(path))

	// Output:
	// log
}

func ExampleTempDir() {
	// init
	var (
		fileName = "gbfile_example_basic_dir"
	)

	// fetch an absolute representation of path.
	path := gbfile.Temp(fileName)

	fmt.Println(path)

	// May Output:
	// /tmp/gbfile_example_basic_dir
}

func ExampleRemove() {
	// init
	var (
		path = gbfile.Join(gbfile.Temp("gbfile_example_basic_dir"), "file1")
	)

	// Checks whether given `path` a file, which means it's not a directory.
	fmt.Println(gbfile.IsFile(path))

	// deletes all file/directory with `path` parameter.
	gbfile.Remove(path)

	// Check again
	fmt.Println(gbfile.IsFile(path))

	// Output:
	// true
	// false
}
