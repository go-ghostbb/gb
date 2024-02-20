package crud

import (
	"ghostbb.io/gb/internal/intlog"
	"gorm.io/gorm"
	"sync"
)

func (h *Handler) afterUpdate(db *gorm.DB) {
	var (
		tableName = h.getTableName(db)
		ctx       = db.Statement.Context
		level     = h.parseLevel(ctx)
	)
	// no rows affected, no need to invalidate cache
	if db.RowsAffected == 0 || level == CacheNone {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		intlog.Printf(ctx, "[AfterUpdate] now start to invalidate search cache for model: %s", tableName)
		if err := h.cache.InvalidateSearchCache(ctx, tableName); err != nil {
			intlog.Errorf(ctx, "[AfterUpdate] invalidating search cache for model %s error: %v", tableName, err)
			return
		}
		intlog.Printf(ctx, "[AfterUpdate] invalidating search cache for model: %s finished.", tableName)
	}()

	if !h.cache.Config().AsyncWrite {
		wg.Wait()
	}
}
