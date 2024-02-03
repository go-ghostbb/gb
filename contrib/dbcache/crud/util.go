package crud

import (
	"context"
	"fmt"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"strings"
)

func (h *Handler) parseLevel(ctx context.Context) cacheLevel {
	value := ctx.Value(cacheCtxLevelKey)
	if value == nil {
		return cacheNone
	}
	if v, ok := value.(cacheLevel); ok {
		return v
	}
	return cacheNone
}

func (h *Handler) genSFKey(tableName string, sql string, vars ...interface{}) string {
	var buf = strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		buf.WriteString(fmt.Sprintf(":%s", gbconv.String(v)))
	}
	return fmt.Sprintf("%s:%s", tableName, buf.String())
}
