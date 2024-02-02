package gbredis

import (
	"context"
	gbmap "ghostbb.io/gb/container/gb_map"
	"ghostbb.io/gb/internal/intlog"
)

var (
	// localInstances for instance management of redis client.
	localInstances = gbmap.NewStrAnyMap(true)
)

// Instance returns an instance of redis client with specified group.
// The `name` param is unnecessary, if `name` is not passed,
// it returns a redis instance with default configuration group.
func Instance(name ...string) *Redis {
	group := DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	v := localInstances.GetOrSetFuncLock(group, func() interface{} {
		if config, ok := GetConfig(group); ok {
			r, err := New(config)
			if err != nil {
				intlog.Errorf(context.TODO(), `%+v`, err)
				return nil
			}
			return r
		}
		return nil
	})
	if v != nil {
		return v.(*Redis)
	}
	return nil
}
