package crud

import (
	"context"
	"fmt"
	gbsha256 "ghostbb.io/gb/crypto/gb_sha256"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"gorm.io/gorm"
	"strings"
)

func (h *Handler) parseLevel(ctx context.Context) cacheLevel {
	value := ctx.Value(CacheCtxKey)
	if value == nil {
		return CacheNone
	}
	if v, ok := value.(cacheLevel); ok {
		return v
	}
	return CacheNone
}

func (h *Handler) genSFKey(tableName string, sql string, vars ...interface{}) string {
	var buf = strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		buf.WriteString(fmt.Sprintf(":%s", gbconv.String(v)))
	}
	return gbsha256.Encrypt256(fmt.Sprintf("%s:%s", tableName, buf.String()))
}

func (h *Handler) getTableName(db *gorm.DB) string {
	if db.Statement.Schema != nil {
		return db.Statement.Schema.Table
	}
	return db.Statement.Table
}
