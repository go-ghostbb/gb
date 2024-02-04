package crud

import (
	gbtype "ghostbb.io/gb/container/gb_type"
	"ghostbb.io/gb/contrib/dbcache/cache"
	"ghostbb.io/gb/contrib/dbcache/singleflight"
	gberror "ghostbb.io/gb/errors/gb_error"
	"gorm.io/gorm"
)

type (
	cacheLevel int
	exprType   string
)

const (
	CacheCtxLevelKey = "gb:gorm:cache:ctx_level"
	afterCreateKey   = "gb:gorm:cache:after_create"
	afterDeleteKey   = "gb:gorm:cache:after_delete"
	afterQueryKey    = "gb:gorm:cache:after_query"
	afterUpdateKey   = "gb:gorm:cache:after_update"
	beforeQueryKey   = "gb:gorm:cache:before_query"

	querySFCallKey = "gb:gorm:cache:sf_call"

	CacheNone    cacheLevel = 0
	CacheAll     cacheLevel = 1
	CachePrimary cacheLevel = 2
	CacheSearch  cacheLevel = 3

	exprEq    exprType = "eq"
	exprIn    exprType = "in"
	exprOther exprType = "other"
)

type Result struct {
	Dest         any
	RowsAffected int64
}

var (
	ErrSFHit           = gberror.New("single flight hit")
	ErrCacheUnmarshal  = gberror.New("cache hit, but unmarshal error")
	ErrPrimaryCacheHit = gberror.New("primary cache hit")
)

func New(cache *cache.Cache) *Handler {
	return &Handler{
		cache:     cache,
		hitCount:  gbtype.NewInt(),
		missCount: gbtype.NewInt(),
	}
}

type Handler struct {
	cache *cache.Cache
	sf    singleflight.Group

	query func(db *gorm.DB)

	hitCount  *gbtype.Int
	missCount *gbtype.Int
}

func (h *Handler) Bind(db *gorm.DB) (err error) {
	err = db.Callback().Create().After("gorm:create").Register(afterCreateKey, h.afterCreate)
	if err != nil {
		return err
	}

	err = db.Callback().Delete().After("gorm:delete").Register(afterDeleteKey, h.afterDelete)
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("gorm:update").Register(afterUpdateKey, h.afterUpdate)
	if err != nil {
		return err
	}

	h.query = db.Callback().Query().Get("gorm:query")

	err = db.Callback().Query().Replace("gorm:query", h.beforeQuery)
	if err != nil {
		return err
	}

	return nil
}
