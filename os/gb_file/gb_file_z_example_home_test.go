package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
)

func ExampleHome() {
	// user's home directory
	homePath, _ := gbfile.Home()
	fmt.Println(homePath)

	// May Output:
	// C:\Users\hailaz
}
