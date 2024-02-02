package gblog

import gbmap "ghostbb.io/gb/container/gb_map"

const (
	// DefaultName is the default group name for instance usage.
	DefaultName = "default"
)

var (
	// Instances map.
	instances = gbmap.NewStrAnyMap(true)
)

// Instance returns an instance of Logger with default settings.
// The parameter `name` is the name for the instance.
func Instance(name ...string) *Logger {
	key := DefaultName
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		return New()
	}).(*Logger)
}
