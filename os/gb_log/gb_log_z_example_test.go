package gblog_test

import (
	"context"
	"ghostbb.io/gb/frame/g"
)

func ExampleContext() {
	ctx := context.WithValue(context.Background(), "Trace-Id", "123456789")
	g.Log().Error(ctx, "runtime error")

	// May Output:
	// 2020-06-08 20:17:03.630 [ERRO] {Trace-Id: 123456789} runtime error
	// Stack:
	// ...
}
