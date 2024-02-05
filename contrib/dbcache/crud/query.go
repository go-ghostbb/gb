package crud

import (
	"fmt"
	"ghostbb.io/gb/contrib/dbcache/cache"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/intlog"
	"ghostbb.io/gb/internal/json"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"reflect"
	"sync"
)

func (h *Handler) beforeQuery(db *gorm.DB) {
	var (
		tableName = h.getTableName(db)
		ctx       = db.Statement.Context
		level     = h.parseLevel(ctx)
	)
	callbacks.BuildQuerySQL(db)
	if level == CacheNone {
		return
	}

	// Use singleflight to avoid cache breakdown
	var sfKey = h.genSFKey(tableName, db.Statement.SQL.String(), db.Statement.Vars...)
	value, err, share := h.sf.Do(sfKey, func() (interface{}, error) {
		if h.searchCache(db) {
			return Result{Dest: db.Statement.Dest, RowsAffected: db.Statement.RowsAffected}, db.Error
		}
		h.query(db)
		h.doCache(db)
		return Result{Dest: db.Statement.Dest, RowsAffected: db.Statement.RowsAffected}, db.Error
	})
	if share {
		intlog.Print(ctx, "singleflight hit!!!")
		if result, ok := value.(Result); ok {
			if err1 := gbconv.Struct(result.Dest, db.Statement.Dest); err1 != nil {
				db.Error = err1
			}
			db.Statement.RowsAffected = result.RowsAffected
		} else {
			db.Error = gberror.NewCode(gbcode.CodeInternalError, "type conversion error")
			return
		}
		if err != nil {
			db.Error = gberror.Wrap(db.Error, err.Error())
		}
	}
}

func (h *Handler) searchCache(db *gorm.DB) (hit bool) {
	var (
		tableName = h.getTableName(db)
		ctx       = db.Statement.Context
		sql       = db.Statement.SQL.String()
		vars      = db.Statement.Vars
	)

	// search cache hit
	cacheValue, err := h.cache.GetSearchCache(ctx, tableName, sql, vars...)
	if err != nil {
		if !gberror.Is(err, cache.ErrCacheNotFound) {
			intlog.Errorf(ctx, "[BeforeQuery] get cache value for sql %s error: %v", sql, err)
		}
		return
	}

	// dealing with cache penetration
	if cacheValue == gorm.ErrRecordNotFound.Error() {
		db.Error = gorm.ErrRecordNotFound
		hit = true
		return
	}

	rowsAffectedPos := gbstr.Pos(cacheValue, "|")
	db.RowsAffected = gbconv.Int64(cacheValue[:rowsAffectedPos])
	err = json.Unmarshal([]byte(cacheValue[rowsAffectedPos+1:]), db.Statement.Dest)
	if err != nil {
		intlog.Errorf(ctx, "[BeforeQuery] unmarshal search cache error: %v", err)
		return
	}
	hit = true
	return
}

func (h *Handler) doCache(db *gorm.DB) {
	var (
		tableName = h.getTableName(db)
		ctx       = db.Statement.Context
		level     = h.parseLevel(ctx)
		sql       = db.Statement.SQL.String()
		vars      = db.Statement.Vars
	)

	if level == CacheNone {
		return
	}

	if db.Error != nil {
		// dealing with cache penetration
		if gberror.Is(db.Error, gorm.ErrRecordNotFound) {
			intlog.Printf(ctx, "[AfterQuery] set cache: %s", gorm.ErrRecordNotFound.Error())
			err := h.cache.SetSearchCache(ctx, gorm.ErrRecordNotFound.Error(), tableName, sql, vars...)
			if err != nil {
				intlog.Errorf(ctx, "[AfterQuery] set search cache for sql: %s error: %v", sql, err)
				return
			}
			intlog.Printf(ctx, "[AfterQuery] sql %s cached", sql)
		}
		return
	}

	var (
		object  = make([]interface{}, 0)
		destRef = reflect.Indirect(reflect.ValueOf(db.Statement.Dest))
	)

	switch destRef.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < destRef.Len(); i++ {
			object = append(object, destRef.Index(i).Interface())
		}
	case reflect.Struct:
		object = append(object, destRef.Interface())
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		if h.cache.Config().MaxItem != 0 && len(object) > h.cache.Config().MaxItem {
			return
		}
		cacheBytes, err := json.Marshal(db.Statement.Dest)
		if err != nil {
			intlog.Errorf(ctx, "[AfterQuery] cannot marshal cache for sql: %s, not cached", sql)
			return
		}
		err = h.cache.SetSearchCache(ctx, fmt.Sprintf("%d|%s", db.RowsAffected, string(cacheBytes)), tableName, sql, vars...)
		if err != nil {
			intlog.Errorf(ctx, "[AfterQuery] set search cache for sql: %s error: %v", sql, err)
			return
		}
	}()

	if !h.cache.Config().AsyncWrite {
		wg.Wait()
	}
}
