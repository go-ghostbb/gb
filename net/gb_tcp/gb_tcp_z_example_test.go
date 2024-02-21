package gbtcp_test

import (
	"fmt"
	gbtcp "ghostbb.io/gb/net/gb_tcp"
)

func ExampleGetFreePort() {
	fmt.Println(gbtcp.GetFreePort())

	// May Output:
	// 57429 <nil>
}

func ExampleGetFreePorts() {
	fmt.Println(gbtcp.GetFreePorts(2))

	// May Output:
	// [57743 57744] <nil>
}
