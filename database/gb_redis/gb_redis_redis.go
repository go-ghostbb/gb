package gbredis

import (
	"context"
	gbvar "ghostbb.io/container/gb_var"
	gbcode "ghostbb.io/errors/gb_code"
	gberror "ghostbb.io/errors/gb_error"
	gbstr "ghostbb.io/text/gb_str"
)

// Redis client.
type Redis struct {
	config *Config
	localAdapter
	localGroup
}

type (
	localGroup struct {
		localGroupGeneric
		localGroupHash
		localGroupList
		localGroupPubSub
		localGroupScript
		localGroupSet
		localGroupSortedSet
		localGroupString
	}
	localAdapter        = Adapter
	localGroupGeneric   = IGroupGeneric
	localGroupHash      = IGroupHash
	localGroupList      = IGroupList
	localGroupPubSub    = IGroupPubSub
	localGroupScript    = IGroupScript
	localGroupSet       = IGroupSet
	localGroupSortedSet = IGroupSortedSet
	localGroupString    = IGroupString
)

const (
	errorNilRedis = `the Redis object is nil`
)

var (
	errorNilAdapter = gbstr.Trim(gbstr.Replace(`
redis adapter is not set, missing configuration or adapter register?
`, "\n", ""))
)

// initGroup initializes the group object of redis.
func (r *Redis) initGroup() *Redis {
	r.localGroup = localGroup{
		localGroupGeneric:   r.localAdapter.GroupGeneric(),
		localGroupHash:      r.localAdapter.GroupHash(),
		localGroupList:      r.localAdapter.GroupList(),
		localGroupPubSub:    r.localAdapter.GroupPubSub(),
		localGroupScript:    r.localAdapter.GroupScript(),
		localGroupSet:       r.localAdapter.GroupSet(),
		localGroupSortedSet: r.localAdapter.GroupSortedSet(),
		localGroupString:    r.localAdapter.GroupString(),
	}
	return r
}

// SetAdapter changes the underlying adapter with custom adapter for current redis client.
func (r *Redis) SetAdapter(adapter Adapter) {
	if r == nil {
		panic(gberror.NewCode(gbcode.CodeInvalidParameter, errorNilRedis))
	}
	r.localAdapter = adapter
}

// GetAdapter returns the adapter that is set in current redis client.
func (r *Redis) GetAdapter() Adapter {
	if r == nil {
		return nil
	}
	return r.localAdapter
}

// Conn retrieves and returns a connection object for continuous operations.
// Note that you should call Close function manually if you do not use this connection any further.
func (r *Redis) Conn(ctx context.Context) (Conn, error) {
	if r == nil {
		return nil, gberror.NewCode(gbcode.CodeInvalidParameter, errorNilRedis)
	}
	if r.localAdapter == nil {
		return nil, gberror.NewCode(gbcode.CodeNecessaryPackageNotImport, errorNilAdapter)
	}
	return r.localAdapter.Conn(ctx)
}

// Do send a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (r *Redis) Do(ctx context.Context, command string, args ...interface{}) (*gbvar.Var, error) {
	if r == nil {
		return nil, gberror.NewCode(gbcode.CodeInvalidParameter, errorNilRedis)
	}
	if r.localAdapter == nil {
		return nil, gberror.NewCodef(gbcode.CodeMissingConfiguration, errorNilAdapter)
	}
	return r.localAdapter.Do(ctx, command, args...)
}

// MustConn performs as function Conn, but it panics if any error occurs internally.
func (r *Redis) MustConn(ctx context.Context) Conn {
	c, err := r.Conn(ctx)
	if err != nil {
		panic(err)
	}
	return c
}

// MustDo performs as function Do, but it panics if any error occurs internally.
func (r *Redis) MustDo(ctx context.Context, command string, args ...interface{}) *gbvar.Var {
	v, err := r.Do(ctx, command, args...)
	if err != nil {
		panic(err)
	}
	return v
}

// Close closes current redis client, closes its connection pool and releases all its related resources.
func (r *Redis) Close(ctx context.Context) error {
	if r == nil || r.localAdapter == nil {
		return nil
	}
	return r.localAdapter.Close(ctx)
}
