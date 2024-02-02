package gbdb

import (
	"fmt"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/deepcopy"
	"ghostbb.io/gb/internal/json"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"reflect"
	"sync"
)

func newQueryHandler(c *Cache) *queryHandler {
	return &queryHandler{cache: c}
}

type queryHandler struct {
	cache        *Cache
	singleFlight Group
}

func (h *queryHandler) Bind(db *gorm.DB) error {
	err := db.Callback().Query().Before("gorm:query").Register("gorm:cache:before_query", h.BeforeQuery())
	if err != nil {
		return err
	}
	err = db.Callback().Query().After("gorm:after_query").Register("gorm:cache:after_query", h.AfterQuery())
	if err != nil {
		return err
	}
	return nil
}

func (h *queryHandler) BeforeQuery() func(db *gorm.DB) {
	cache := h.cache
	return func(db *gorm.DB) {
		callbacks.BuildQuerySQL(db)
		tableName := ""
		if db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		} else {
			tableName = db.Statement.Table
		}
		ctx := db.Statement.Context

		sql := db.Statement.SQL.String()
		db.InstanceSet("gorm:cache:sql", sql)
		db.InstanceSet("gorm:cache:vars", db.Statement.Vars)

		if shouldCache(tableName, cache.Config.Tables) && cacheCtxCheck(ctx) {
			hit := false
			defer func() {
				if hit {
					cache.IncrHitCount()
				} else {
					cache.IncrMissCount()
				}
			}()

			// singleFlight Check
			singleFlightKey := genSingleFlightKey(tableName, sql, db.Statement.Vars...)
			h.singleFlight.mu.Lock()
			if h.singleFlight.m == nil {
				h.singleFlight.m = make(map[string]*call)
			}
			if c, ok := h.singleFlight.m[singleFlightKey]; ok {
				c.dups++
				h.singleFlight.mu.Unlock()
				c.wg.Wait()

				db.Statement.Dest = deepcopy.Copy(c.dest)
				hit = true
				db.RowsAffected = c.rowsAffected
				db.Error = singleFlightHit // 為保證後續流程不走，必須設一個error
				if c.err != nil {
					// 紀錄原本的error
					db.Error = gberror.Wrap(db.Error, c.err.Error())
				}
				cache.logger.Info(ctx, "[BeforeQuery] single flight hit for key %v", singleFlightKey)

				return
			}
			c := &call{key: singleFlightKey}
			c.wg.Add(1)
			h.singleFlight.m[singleFlightKey] = c
			h.singleFlight.mu.Unlock()
			db.InstanceSet("gorm:cache:query:single_flight_call", c)

			if cache.Config.Level == CacheLevelAll || cache.Config.Level == CacheLevelOnlyPrimary {
				tryPrimaryCache := func() (hit bool) {
					primaryKeys := getPrimaryKeysFromWhereClause(db)

					cache.logger.Info(ctx, "[BeforeQuery] parse primary keys = %v", primaryKeys)

					if len(primaryKeys) == 0 {
						return
					}

					// if (IN primaryKeys)/(Eq primaryKey) are the only clauses
					hasOtherClauseInWhere := hasOtherClauseExceptPrimaryField(db)
					if hasOtherClauseInWhere {
						// if query has other clauses, it can only query the database
						return
					}

					// primary cache hit
					cacheValues, err := cache.BatchGetPrimaryCache(ctx, tableName, primaryKeys)
					if err != nil {
						cache.logger.Error(ctx, "[BeforeQuery] get primary cache value for key %v error: %v", primaryKeys, err)
						db.Error = nil
						return
					}
					if len(cacheValues) != len(primaryKeys) {
						db.Error = nil
						return
					}
					finalValue := ""

					destKind := reflect.Indirect(reflect.ValueOf(db.Statement.Dest)).Kind()
					if destKind == reflect.Struct && len(cacheValues) == 1 {
						finalValue = cacheValues[0]
					} else if (destKind == reflect.Array || destKind == reflect.Slice) && len(cacheValues) >= 1 {
						finalValue = "[" + gbstr.Join(cacheValues, ",") + "]"
					}
					if len(finalValue) == 0 {
						cache.logger.Error(ctx, "[BeforeQuery] length of cache values and dest not matched")
						db.Error = errCacheUnmarshal
						return
					}

					err = json.Unmarshal([]byte(finalValue), db.Statement.Dest)
					if err != nil {
						cache.logger.Error(ctx, "[BeforeQuery] unmarshal final value error: %v", err)
						db.Error = errCacheUnmarshal
						return
					}
					db.Error = primaryCacheHit
					hit = true
					return
				}
				if tryPrimaryCache() {
					hit = true
					return
				}
			}
			if cache.Config.Level == CacheLevelAll || cache.Config.Level == CacheLevelOnlySearch {
				trySearchCache := func() (hit bool) {
					// search cache hit
					cacheValue, err := cache.GetSearchCache(ctx, tableName, sql, db.Statement.Vars...)
					if err != nil {
						if !gberror.Is(err, ErrCacheNotFound) {
							cache.logger.Error(ctx, "[BeforeQuery] get cache value for sql %s error: %v", sql, err)
						}
						db.Error = nil
						return
					}
					cache.logger.Info(ctx, "[BeforeQuery] get value: %s", cacheValue)
					if cacheValue == "recordNotFound" { // 應對緩存穿透
						db.Error = recordNotFoundCacheHit
						hit = true
						return
					}

					rowsAffectedPos := gbstr.Pos(cacheValue, "|")
					db.RowsAffected = gbconv.Int64(cacheValue[:rowsAffectedPos])
					if err != nil {
						cache.logger.Error(ctx, "[BeforeQuery] unmarshal rows affected cache error: %v", err)
						db.Error = nil
						return
					}
					err = json.Unmarshal([]byte(cacheValue[rowsAffectedPos+1:]), db.Statement.Dest)
					if err != nil {
						cache.logger.Error(ctx, "[BeforeQuery] unmarshal search cache error: %v", err)
						db.Error = nil
						return
					}
					db.Error = searchCacheHit
					hit = true
					return
				}
				if !hit && trySearchCache() {
					hit = true
				}
			}
		}
	}
}

func (h *queryHandler) AfterQuery() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		h.doCache(db)
		// 上面的cache完成後直接傳播給其他等待中的goroutine
		// 上面只處理非singleflight且無錯誤或記錄不存在的情況
		h.fillCallAfterQuery(db)

		// 下面處理命中了緩存的情況
		// 有以下幾種err是專門用來傳狀態的：正常的cacheHit 這種情況不存在error
		// RecordNotFoundCacheHit 這種情況只會在notfound之後出現
		// SingleFlightHit 這種情況下error中除了SingleFlightHit還可能會存在其他error來自gorm的error
		// 且遇到任何一種hit我們都可以認為是命中了緩存 同時只可能命中至多兩個hit（single+其他
		if gberror.HasError(db.Error, singleFlightHit) {
			if e, ok := db.Error.(*gberror.Error); ok {
				if e.Is(singleFlightHit) {
					db.Error = nil
				} else {
					db.Error = e.Current()
				}
			}
		}

		switch {
		case gberror.Is(db.Error, recordNotFoundCacheHit):
			db.Error = gorm.ErrRecordNotFound
		case gberror.Is(db.Error, searchCacheHit) || gberror.Is(db.Error, primaryCacheHit):
			db.Error = nil
		}
	}
}

func (h *queryHandler) doCache(db *gorm.DB) {
	var (
		cache     = h.cache
		tableName = ""
	)
	if db.Statement.Schema != nil {
		tableName = db.Statement.Schema.Table
	} else {
		tableName = db.Statement.Table
	}
	ctx := db.Statement.Context
	sqlObj, _ := db.InstanceGet("gorm:cache:sql")
	sql := sqlObj.(string)
	varObj, _ := db.InstanceGet("gorm:cache:vars")
	vars := varObj.([]interface{})

	if !(shouldCache(tableName, cache.Config.Tables) && cacheCtxCheck(ctx)) {
		return
	}

	// 這裡對應singleflight
	// 如果擊中db.Error裡面會被塞入一個singleflight hit error
	// 或是擊中緩存也會被塞入一個error
	// 就不用進行緩存
	if db.Error == nil {
		destValue := reflect.Indirect(reflect.ValueOf(db.Statement.Dest))
		// 如果是struct應該能提主鍵出來
		// 如果是slice需要判斷内部元素是不是struct，不是struct的都提不了主键
		if destValue.Kind() == reflect.Slice || destValue.Kind() == reflect.Array {
			if (destValue.Type().Elem().Kind() == reflect.Pointer && destValue.Type().Elem().Elem().Kind() != reflect.Struct) ||
				(destValue.Type().Elem().Kind() != reflect.Pointer && destValue.Type().Elem().Kind() != reflect.Struct) {
				return
			}
		}

		// error is nil -> cache not hit, we cache newly retrieved data
		primaryKeys, objects := getObjectsAfterLoad(db)

		var wg sync.WaitGroup
		wg.Add(2)

		// TODO: cache level ALL and SEARCH
		go func() {
			defer wg.Done()

			if cache.Config.Level == CacheLevelAll || cache.Config.Level == CacheLevelOnlySearch {
				// cache search data
				if cache.Config.MaxItemCnt != 0 && int64(len(objects)) > cache.Config.MaxItemCnt {
					return
				}

				cache.logger.Info(ctx, "[AfterQuery] start to set search cache for sql: %s", sql)
				cacheBytes, err := json.Marshal(db.Statement.Dest)
				if err != nil {
					cache.logger.Error(ctx, "[AfterQuery] cannot marshal cache for sql: %s, not cached", sql)
					return
				}
				cache.logger.Info(ctx, "[AfterQuery] set cache: %s", string(cacheBytes))
				err = cache.SetSearchCache(ctx, fmt.Sprintf("%d|", db.RowsAffected)+string(cacheBytes), tableName, sql, vars...)
				if err != nil {
					cache.logger.Error(ctx, "[AfterQuery] set search cache for sql: %s error: %v", sql, err)
					return
				}
				cache.logger.Info(ctx, "[AfterQuery] sql %s cached", sql)
			}
		}()

		// TODO: cache level ALL and PRIMARY
		go func() {
			defer wg.Done()

			if cache.Config.Level == CacheLevelAll || cache.Config.Level == CacheLevelOnlyPrimary {
				// cache primary cache data
				if len(primaryKeys) != len(objects) {
					return
				}
				if cache.Config.MaxItemCnt != 0 && int64(len(objects)) > cache.Config.MaxItemCnt {
					cache.logger.Info(ctx, "[AfterQuery] objects length is more than max item count, not cached")
					return
				}
				kvs := make([]Kv, 0, len(objects))
				for i := 0; i < len(objects); i++ {
					jsonStr, err := json.Marshal(objects[i])
					if err != nil {
						cache.logger.Error(ctx, "[AfterQuery] object %v cannot marshal, not cached", objects[i])
						continue
					}
					kvs = append(kvs, Kv{
						Key:   primaryKeys[i],
						Value: string(jsonStr),
					})
				}
				cache.logger.Info(ctx, "[AfterQuery] start to set primary cache for kvs: %+v", kvs)
				err := cache.BatchSetPrimaryKeyCache(ctx, tableName, kvs)
				if err != nil {
					cache.logger.Error(ctx, "[AfterQuery] batch set primary key cache for key %v error: %v", primaryKeys, err)
				}
			}
		}()

		if !cache.Config.AsyncWrite {
			wg.Wait()
		}
		return
	}

	// 對應緩存穿透
	if gberror.Is(db.Error, gorm.ErrRecordNotFound) && !cache.Config.DisableCachePenetrationProtect {
		cache.logger.Info(ctx, "[AfterQuery] set cache: %v", "recordNotFound")
		err := cache.SetSearchCache(ctx, "recordNotFound", tableName, sql, vars...)
		if err != nil {
			cache.logger.Error(ctx, "[AfterQuery] set search cache for sql: %s error: %v", sql, err)
			return
		}
		cache.logger.Info(ctx, "[AfterQuery] sql %s cached", sql)
		return
	}
}

func (h *queryHandler) fillCallAfterQuery(db *gorm.DB) {
	if singleFlightCallObj, exist := db.InstanceGet("gorm:cache:query:single_flight_call"); exist {
		c := singleFlightCallObj.(*call)
		c.dest = db.Statement.Dest
		c.rowsAffected = db.RowsAffected
		c.err = db.Error
		c.wg.Done()

		h.singleFlight.mu.Lock()
		if !c.forgotten {
			delete(h.singleFlight.m, c.key)
		}
		h.singleFlight.mu.Unlock()
	}
}
