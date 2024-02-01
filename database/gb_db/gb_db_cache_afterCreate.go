package gbdb

import (
	"gorm.io/gorm"
)

func AfterCreate(cache *Cache) func(db *gorm.DB) {
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
			if cache.Config.Level == CacheLevelAll || cache.Config.Level == CacheLevelOnlySearch {
				invalidSearchCache := func() {
					// We invalidate search cache here,
					// because any newly created objects may cause search cache results to be outdated and invalid.
					cache.logger.Info(ctx, "[AfterCreate] now start to invalidate search cache for model: %s", tableName)

					err := cache.InvalidateSearchCache(ctx, tableName)
					if err != nil {
						cache.logger.Error(ctx, "[AfterCreate] invalidating search cache for model %s error: %v", tableName, err)
						return
					}
					cache.logger.Info(ctx, "[AfterCreate] invalidating search cache for model: %s finished.", tableName)
				}
				if cache.Config.AsyncWrite {
					go invalidSearchCache()
				} else {
					invalidSearchCache()
				}
			}
		}
	}
}
