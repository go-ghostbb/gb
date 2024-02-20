package gbproc_test

import (
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbproc "ghostbb.io/gb/os/gb_proc"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_ShellExec(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s, err := gbproc.ShellExec(gbctx.New(), `echo 123`)
		t.AssertNil(err)
		t.Assert(s, "123\r\n")
	})
	// error
	gbtest.C(t, func(t *gbtest.T) {
		_, err := gbproc.ShellExec(gbctx.New(), `NoneExistCommandCall`)
		t.AssertNE(err, nil)
	})
}
