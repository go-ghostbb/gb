package gbredis

import (
	"context"
	gbvar "github.com/Ghostbb-io/gb/container/gb_var"
)

// IGroupScript manages redis script operations.
// Implements see redis.GroupScript.
type IGroupScript interface {
	Eval(ctx context.Context, script string, numKeys int64, keys []string, args []interface{}) (*gbvar.Var, error)
	EvalSha(ctx context.Context, sha1 string, numKeys int64, keys []string, args []interface{}) (*gbvar.Var, error)
	ScriptLoad(ctx context.Context, script string) (string, error)
	ScriptExists(ctx context.Context, sha1 string, sha1s ...string) (map[string]bool, error)
	ScriptFlush(ctx context.Context, option ...ScriptFlushOption) error
	ScriptKill(ctx context.Context) error
}

// ScriptFlushOption provides options for function ScriptFlush.
type ScriptFlushOption struct {
	SYNC  bool // SYNC  flushes the cache synchronously.
	ASYNC bool // ASYNC flushes the cache asynchronously.
}
