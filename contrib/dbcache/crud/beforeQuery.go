package crud

import (
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/deepcopy"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func (h *Handler) beforeQuery(db *gorm.DB) {
	var (
		tableName string
		ctx       = db.Statement.Context
		level     = h.parseLevel(ctx)
		sql       = db.Statement.SQL.String()
		vars      = db.Statement.Vars
	)

	callbacks.BuildQuerySQL(db)
	if db.Statement.Schema != nil {
		tableName = db.Statement.Schema.Table
	} else {
		tableName = db.Statement.Table
	}

	db.InstanceSet("gorm:cache:sql", sql)
	db.InstanceSet("gorm:cache:vars", vars)

	if level == cacheNone {
		return
	}

	var hit = false
	defer func() {
		if hit {
			h.hitCount.Add(1)
		} else {
			h.missCount.Add(1)
		}
	}()

	// Use singleflight to avoid cache penetration
	var sfKey = h.genSFKey(tableName, sql, vars...)
	value, err, share := h.sf.Do(sfKey, func() (interface{}, error) {
		return nil, nil
	})
	if share {
		hit = true
		db.Statement.Dest = deepcopy.Copy(value)
		db.Error = ErrSFHit
		if err != nil {
			db.Error = gberror.Wrap(db.Error, err.Error())
		}
	}
	db.InstanceSet(querySFCallKey, value)

}
