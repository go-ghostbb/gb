package gbproc

import (
	"context"
	"io"
)

// MustShell performs as Shell, but it panics if any error occurs.
func MustShell(ctx context.Context, cmd string, out io.Writer, in io.Reader) {
	if err := Shell(ctx, cmd, out, in); err != nil {
		panic(err)
	}
}

// MustShellRun performs as ShellRun, but it panics if any error occurs.
func MustShellRun(ctx context.Context, cmd string) {
	if err := ShellRun(ctx, cmd); err != nil {
		panic(err)
	}
}

// MustShellExec performs as ShellExec, but it panics if any error occurs.
func MustShellExec(ctx context.Context, cmd string, environment ...[]string) string {
	result, err := ShellExec(ctx, cmd, environment...)
	if err != nil {
		panic(err)
	}
	return result
}
