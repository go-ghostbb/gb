package gbdb

import (
	"gorm.io/gorm"
	"sync"
)

func AfterUpdate(cache *Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.RowsAffected == 0 {
			return // no rows affected, no need to invalidate cache
		}

		tableName := ""
		if db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		} else {
			tableName = db.Statement.Table
		}
		ctx := db.Statement.Context

		do := db.Error == nil && cache.Config.InvalidateWhenUpdate && shouldCache(tableName, cache.Config.Tables) && cacheCtxCheck(ctx)
		if do {
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()

				if cache.Config.Level == CacheLevelAll || cache.Config.Level == CacheLevelOnlyPrimary {
					primaryKeys := getPrimaryKeysFromWhereClause(db)
					cache.logger.Info(ctx, "[AfterUpdate] parse primary keys = %v", primaryKeys)

					if len(primaryKeys) > 0 {
						cache.logger.Info(ctx, "[AfterUpdate] now start to invalidate cache for primary keys: %+v", primaryKeys)

						err := cache.BatchInvalidatePrimaryCache(ctx, tableName, primaryKeys)
						if err != nil {
							cache.logger.Error(ctx, "[AfterUpdate] invalidating primary cache for key %v error: %v", primaryKeys, err)
							return
						}
						cache.logger.Info(ctx, "[AfterUpdate] invalidating cache for primary keys: %+v finished.", primaryKeys)
					} else {
						cache.logger.Info(ctx, "[AfterUpdate] now start to invalidate all primary cache for model: %s", tableName)
						err := cache.InvalidateAllPrimaryCache(ctx, tableName)
						if err != nil {
							cache.logger.Error(ctx, "[AfterUpdate] invalidating primary cache for model %s error: %v", tableName, err)
							return
						}
						cache.logger.Info(ctx, "[AfterUpdate] invalidating all primary cache for model: %s finished.", tableName)
					}
				}
			}()

			go func() {
				defer wg.Done()

				if cache.Config.Level == CacheLevelAll || cache.Config.Level == CacheLevelOnlySearch {
					cache.logger.Info(ctx, "[AfterUpdate] now start to invalidate search cache for model: %s", tableName)
					err := cache.InvalidateSearchCache(ctx, tableName)
					if err != nil {
						cache.logger.Error(ctx, "[AfterUpdate] invalidating search cache for model %s error: %v", tableName, err)
						return
					}
					cache.logger.Info(ctx, "[AfterUpdate] invalidating search cache for model: %s finished.", tableName)
				}
			}()

			if !cache.Config.AsyncWrite {
				wg.Wait()
			}
		}
	}
}
