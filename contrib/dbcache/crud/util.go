package crud

import (
	"context"
	"fmt"
	gbvar "ghostbb.io/gb/container/gb_var"
	gbsha256 "ghostbb.io/gb/crypto/gb_sha256"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"gorm.io/gorm"
	"reflect"
	"slices"
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

func (h *Handler) doFormat(dest any) any {
	destVar := gbvar.New(dest)

	if destVar.IsSlice() {
		return sliceFormat(destVar.Vars())
	}

	if destVar.IsStruct() {
		return structFormat(destVar)
	}

	return dest
}

func sliceFormat(destVars []*gbvar.Var) any {
	result := make([]any, 0)
	for _, v := range destVars {
		result = append(result, structFormat(v))
	}
	return result
}

func structFormat(destVar *gbvar.Var) any {
	var (
		newFields = make([]reflect.StructField, 0)
		newValues = make([]reflect.Value, 0)
		destRef   = reflect.Indirect(reflect.ValueOf(destVar.Interface()))
	)

	// 遍歷所有字段
	// 不緩存結構體和值是zero的字段
	// 以下例外:
	// 1.結構體為"gorm.Model"
	// 2.含有tag->dbcache:"true"
	// 3.gorm tag中含有embedded
	for i := 0; i < destRef.Type().NumField(); i++ {
		var (
			field      = destRef.Type().Field(i)
			fieldRef   = reflect.Indirect(destRef.Field(i))
			fieldPtr   = destRef.Field(i)
			gormTags   = gbstr.Split(gbstr.Replace(field.Tag.Get("gorm"), " ", ""), ";")
			dbcacheTag = gbstr.Replace(field.Tag.Get("dbcache"), " ", "")
		)
		slices.Sort(gormTags)
		_, embedded := slices.BinarySearch(gormTags, "embedded")

		if fieldRef.Kind() == reflect.Struct && fieldRef.Type().String() != "gorm.Model" && dbcacheTag != "true" && !embedded {
			continue
		}

		if fieldPtr.Kind() == reflect.Ptr && fieldPtr.IsZero() || fieldRef.IsZero() {
			continue
		}

		if dbcacheTag == "true" || embedded {
			// 遞迴檢查下一個結構體
			value := structFormat(gbvar.New(destRef.Field(i).Interface()))
			valueRef := reflect.Indirect(reflect.ValueOf(value))
			if valueRef.NumField() == 0 {
				continue
			}
			// 重寫字段結構體
			newFields = append(newFields, reflect.StructField{
				Name: field.Name,
				Type: valueRef.Type(),
				Tag:  field.Tag,
			})
			newValues = append(newValues, valueRef)
		} else {
			newFields = append(newFields, field)
			newValues = append(newValues, destRef.Field(i))
		}
	}

	// 創建一个包含新字段的結構體類型
	newStructType := reflect.StructOf(newFields)

	// 使用新結構體類型創建一个空的實例
	newStruct := reflect.New(newStructType).Elem()

	// 將原始結構體的值複製到新結構體中
	for i := 0; i < len(newFields); i++ {
		newStruct.Field(i).Set(newValues[i])
	}

	return newStruct.Interface()
}
