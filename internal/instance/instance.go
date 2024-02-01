// Package instance provides instances management.
//
// Note that this package is not used for cache, as it has no cache expiration.
package instance

import (
	gbmap "github.com/Ghostbb-io/gb/container/gb_map"
	gbhash "github.com/Ghostbb-io/gb/encoding/gb_hash"
)

const (
	groupNumber = 64
)

var (
	groups = make([]*gbmap.StrAnyMap, groupNumber)
)

func init() {
	for i := 0; i < groupNumber; i++ {
		groups[i] = gbmap.NewStrAnyMap(true)
	}
}

func getGroup(key string) *gbmap.StrAnyMap {
	return groups[int(gbhash.DJB([]byte(key))%groupNumber)]
}

// Get returns the instance by given name.
func Get(name string) interface{} {
	return getGroup(name).Get(name)
}

// Set sets an instance to the instance manager with given name.
func Set(name string, instance interface{}) {
	getGroup(name).Set(name, instance)
}

// GetOrSet returns the instance by name,
// or set instance to the instance manager if it does not exist and returns this instance.
func GetOrSet(name string, instance interface{}) interface{} {
	return getGroup(name).GetOrSet(name, instance)
}

// GetOrSetFunc returns the instance by name,
// or sets instance with returned value of callback function `f` if it does not exist
// and then returns this instance.
func GetOrSetFunc(name string, f func() interface{}) interface{} {
	return getGroup(name).GetOrSetFunc(name, f)
}

// GetOrSetFuncLock returns the instance by name,
// or sets instance with returned value of callback function `f` if it does not exist
// and then returns this instance.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function `f`
// with mutex.Lock of the hash map.
func GetOrSetFuncLock(name string, f func() interface{}) interface{} {
	return getGroup(name).GetOrSetFuncLock(name, f)
}

// SetIfNotExist sets `instance` to the map if the `name` does not exist, then returns true.
// It returns false if `name` exists, and `instance` would be ignored.
func SetIfNotExist(name string, instance interface{}) bool {
	return getGroup(name).SetIfNotExist(name, instance)
}

// Clear deletes all instances stored.
func Clear() {
	for i := 0; i < groupNumber; i++ {
		groups[i].Clear()
	}
}
