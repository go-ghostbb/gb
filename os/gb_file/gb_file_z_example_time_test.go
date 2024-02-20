package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
)

func ExampleMTime() {
	t := gbfile.MTime(gbfile.Temp())
	fmt.Println(t)

	// May Output:
	// 2021-11-02 15:18:43.901141 +0800 CST
}

func ExampleMTimestamp() {
	t := gbfile.MTimestamp(gbfile.Temp())
	fmt.Println(t)

	// May Output:
	// 1635838398
}

func ExampleMTimestampMilli() {
	t := gbfile.MTimestampMilli(gbfile.Temp())
	fmt.Println(t)

	// May Output:
	// 1635838529330
}
