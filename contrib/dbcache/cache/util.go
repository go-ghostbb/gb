package cache

import (
	"fmt"
	gbsha256 "ghostbb.io/gb/crypto/gb_sha256"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"math/rand"
	"strings"
)

func (c *Cache) genPrimaryCacheKey(tableName string, primaryKey string) string {
	return fmt.Sprintf("%s:%s:p:%s:%s", cacheName, c.InstanceId, tableName, primaryKey)
}

func (c *Cache) genSearchCacheKey(tableName string, sql string, vars ...interface{}) string {
	buf := strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		buf.WriteString(fmt.Sprintf(":%s", gbconv.String(v)))
	}
	return fmt.Sprintf("%s:%s:s:%s:%s", cacheName, c.InstanceId, tableName, gbsha256.Encrypt256(buf.String()))
}

func (c *Cache) genCachePrefix(tableName string) string {
	return cacheName + ":" + c.InstanceId + ":p:" + tableName
}

func (c *Cache) TTL() int64 {
	randNum := rand.Float64()*0.2 + 0.9
	return int64(float64(c.config.TTL) * randNum)
}
