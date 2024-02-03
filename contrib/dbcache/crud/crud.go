package crud

import (
	gbtype "ghostbb.io/gb/container/gb_type"
	"ghostbb.io/gb/contrib/dbcache/cache"
	"ghostbb.io/gb/contrib/dbcache/singleflight"
	gberror "ghostbb.io/gb/errors/gb_error"
	"gorm.io/gorm"
)

type cacheLevel int

const (
	cacheCtxLevelKey = "gb:gorm:cache:ctx_level"
	afterCreateKey   = "gb:gorm:cache:after_create"
	afterDeleteKey   = "gb:gorm:cache:after_delete"
	afterQueryKey    = "gb:gorm:cache:after_query"
	afterUpdateKey   = "gb:gorm:cache:after_update"
	beforeQueryKey   = "gb:gorm:cache:before_query"

	querySFCallKey = "gb:gorm:cache:sf_call"

	cacheNone    cacheLevel = 0
	cacheAll     cacheLevel = 1
	cachePrimary cacheLevel = 2
	cacheSearch  cacheLevel = 3
)

var (
	ErrSFHit = gberror.New("single flight hit")
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

	err = db.Callback().Query().Before("gorm:query").Register(beforeQueryKey, h.beforeQuery)
	if err != nil {
		return err
	}

	err = db.Callback().Query().After("gorm:after_query").Register(afterQueryKey, h.afterQuery)
	if err != nil {
		return err
	}

	return nil
}
