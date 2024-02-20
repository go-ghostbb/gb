package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
)

func ExampleScanDir() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_scan_dir")
		tempFile = gbfile.Join(tempDir, fileName)

		tempSubDir  = gbfile.Join(tempDir, "sub_dir")
		tempSubFile = gbfile.Join(tempSubDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")
	gbfile.PutContents(tempSubFile, "ghostbb example content")

	// scans directory recursively
	list, _ := gbfile.ScanDir(tempDir, "*", true)
	for _, v := range list {
		fmt.Println(gbfile.Basename(v))
	}

	// Output:
	// gbfile_example.txt
	// sub_dir
	// gbfile_example.txt
}

func ExampleScanDirFile() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_scan_dir_file")
		tempFile = gbfile.Join(tempDir, fileName)

		tempSubDir  = gbfile.Join(tempDir, "sub_dir")
		tempSubFile = gbfile.Join(tempSubDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")
	gbfile.PutContents(tempSubFile, "ghostbb example content")

	// scans directory recursively exclusive of directories
	list, _ := gbfile.ScanDirFile(tempDir, "*.txt", true)
	for _, v := range list {
		fmt.Println(gbfile.Basename(v))
	}

	// Output:
	// gbfile_example.txt
	// gbfile_example.txt
}

func ExampleScanDirFunc() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_scan_dir_func")
		tempFile = gbfile.Join(tempDir, fileName)

		tempSubDir  = gbfile.Join(tempDir, "sub_dir")
		tempSubFile = gbfile.Join(tempSubDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")
	gbfile.PutContents(tempSubFile, "ghostbb example content")

	// scans directory recursively
	list, _ := gbfile.ScanDirFunc(tempDir, "*", true, func(path string) string {
		// ignores some files
		if gbfile.Basename(path) == "gbfile_example.txt" {
			return ""
		}
		return path
	})
	for _, v := range list {
		fmt.Println(gbfile.Basename(v))
	}

	// Output:
	// sub_dir
}

func ExampleScanDirFileFunc() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_scan_dir_file_func")
		tempFile = gbfile.Join(tempDir, fileName)

		fileName1 = "gbfile_example_ignores.txt"
		tempFile1 = gbfile.Join(tempDir, fileName1)

		tempSubDir  = gbfile.Join(tempDir, "sub_dir")
		tempSubFile = gbfile.Join(tempSubDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")
	gbfile.PutContents(tempFile1, "ghostbb example content")
	gbfile.PutContents(tempSubFile, "ghostbb example content")

	// scans directory recursively exclusive of directories
	list, _ := gbfile.ScanDirFileFunc(tempDir, "*.txt", true, func(path string) string {
		// ignores some files
		if gbfile.Basename(path) == "gbfile_example_ignores.txt" {
			return ""
		}
		return path
	})
	for _, v := range list {
		fmt.Println(gbfile.Basename(v))
	}

	// Output:
	// gbfile_example.txt
	// gbfile_example.txt
}
