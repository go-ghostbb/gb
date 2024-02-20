package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
)

func ExampleCopy() {
	// init
	var (
		srcFileName = "gbfile_example.txt"
		srcTempDir  = gbfile.Temp("gbfile_example_copy_src")
		srcTempFile = gbfile.Join(srcTempDir, srcFileName)

		// copy file
		dstFileName = "gbfile_example_copy.txt"
		dstTempFile = gbfile.Join(srcTempDir, dstFileName)

		// copy dir
		dstTempDir = gbfile.Temp("gbfile_example_copy_dst")
	)

	// write contents
	gbfile.PutContents(srcTempFile, "ghostbb example copy")

	// copy file
	gbfile.Copy(srcTempFile, dstTempFile)

	// read contents after copy file
	fmt.Println(gbfile.GetContents(dstTempFile))

	// copy dir
	gbfile.Copy(srcTempDir, dstTempDir)

	// list copy dir file
	fList, _ := gbfile.ScanDir(dstTempDir, "*", false)
	for _, v := range fList {
		fmt.Println(gbfile.Basename(v))
	}

	// Output:
	// ghostbb example copy
	// gbfile_example.txt
	// gbfile_example_copy.txt
}
