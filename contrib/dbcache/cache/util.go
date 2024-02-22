package cache

import (
	"fmt"
	gbsha256 "ghostbb.io/gb/crypto/gb_sha256"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"math/rand"
	"strings"
	"time"
)

func (c *Cache) genSearchCacheKey(tableName string, sql string, vars ...interface{}) string {
	buf := strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		buf.WriteString(fmt.Sprintf(":%s", gbconv.String(v)))
	}
	return fmt.Sprintf("%s:%s:s:%s:%s", cacheName, c.InstanceId, tableName, gbsha256.Encrypt256(buf.String()))
}

func (c *Cache) genCachePrefix(tableName string) string {
	return cacheName + ":" + c.InstanceId + ":s:" + tableName
}

func (c *Cache) TTL() time.Duration {
	randNum := rand.Float64()*0.2 + 0.9
	return time.Duration(float64(c.config.TTL) * randNum)
}
