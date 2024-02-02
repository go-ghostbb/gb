package gbdb

import (
	"context"
	"fmt"
	gberror "ghostbb.io/gb/errors/gb_error"
	"math/rand"
	"reflect"
	"strings"
)

type Kv struct {
	Key   string
	Value string
}

var (
	recordNotFoundCacheHit = gberror.New("record not found cache hit")
	primaryCacheHit        = gberror.New("primary cache hit")
	searchCacheHit         = gberror.New("search cache hit")
	singleFlightHit        = gberror.New("single flight hit")

	errCacheUnmarshal = gberror.New("cache hit, but unmarshal error")
)

func genPrimaryCacheKey(instanceId string, tableName string, primaryKey string) string {
	return fmt.Sprintf("%s:%s:p:%s:%s", DBCacheName, instanceId, tableName, primaryKey)
}

func genPrimaryCachePrefix(instanceId string, tableName string) string {
	return DBCacheName + ":" + instanceId + ":p:" + tableName
}

func genSearchCacheKey(instanceId string, tableName string, sql string, vars ...interface{}) string {
	buf := strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		pv := reflect.ValueOf(v)
		if pv.Kind() == reflect.Ptr {
			buf.WriteString(fmt.Sprintf(":%v", pv.Elem()))
		} else {
			buf.WriteString(fmt.Sprintf(":%v", v))
		}
	}
	return fmt.Sprintf("%s:%s:s:%s:%s", DBCacheName, instanceId, tableName, buf.String())
}

func genSearchCachePrefix(instanceId string, tableName string) string {
	return DBCacheName + ":" + instanceId + ":s:" + tableName
}

func genSingleFlightKey(tableName string, sql string, vars ...interface{}) string {
	buf := strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		pv := reflect.ValueOf(v)
		if pv.Kind() == reflect.Ptr {
			buf.WriteString(fmt.Sprintf(":%v", pv.Elem()))
		} else {
			buf.WriteString(fmt.Sprintf(":%v", v))
		}
	}
	return fmt.Sprintf("%s:%s", tableName, buf.String())
}

func shouldCache(tableName string, tables []string) bool {
	if len(tables) == 0 {
		return true
	}
	return containString(tableName, tables)
}

func cacheCtxCheck(ctx context.Context) bool {
	value := ctx.Value(DBCacheName)
	switch value.(type) {
	case struct{}:
		return true
	}
	return true
}

func containString(target string, slice []string) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}

func randFloatingInt64(v int64) int64 {
	randNum := rand.Float64()*0.2 + 0.9
	return int64(float64(v) * randNum)
}
