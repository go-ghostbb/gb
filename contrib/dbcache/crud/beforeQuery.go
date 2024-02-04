package crud

import (
	gbcode "ghostbb.io/gb/errors/gb_code"
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
	)
	callbacks.BuildQuerySQL(db)
	if level == CacheNone {
		return
	}

	if db.Statement.Schema != nil {
		tableName = db.Statement.Schema.Table
	} else {
		tableName = db.Statement.Table
	}

	// Use singleflight to avoid cache breakdown
	var sfKey = h.genSFKey(tableName, db.Statement.SQL.String(), db.Statement.Vars...)
	value, err, share := h.sf.Do(sfKey, func() (interface{}, error) {
		h.query(db)
		h.afterQuery(db)
		return nil, nil
	})
	if share {
		if result, ok := value.(Result); ok {
			db.Statement.Dest = deepcopy.Copy(result.Dest)
			db.Statement.RowsAffected = result.RowsAffected
		} else {
			db.Error = gberror.NewCode(gbcode.CodeInternalError, "type conversion error")
			return
		}
		// Use err to mark singleflight hit
		db.Error = ErrSFHit
		if err != nil {
			db.Error = gberror.Wrap(db.Error, err.Error())
		}
		return
	}
	db.InstanceSet(querySFCallKey, value)
}
