package gins

import (
	"context"
	"fmt"
	gbredis "ghostbb.io/gb/database/gb_redis"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/consts"
	"ghostbb.io/gb/internal/instance"
	"ghostbb.io/gb/internal/intlog"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
)

// Redis returns an instance of redis client with specified configuration group name.
// Note that it panics if any error occurs duration instance creating.
func Redis(name ...string) *gbredis.Redis {
	var (
		err   error
		ctx   = context.Background()
		group = gbredis.DefaultGroupName
	)
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameRedis, group)
	result := instance.GetOrSetFuncLock(instanceKey, func() interface{} {
		// If already configured, it returns the redis instance.
		if _, ok := gbredis.GetConfig(group); ok {
			return gbredis.Instance(group)
		}
		if Config().Available(ctx) {
			var (
				configMap   map[string]interface{}
				redisConfig *gbredis.Config
				redisClient *gbredis.Redis
			)
			if configMap, err = Config().Data(ctx); err != nil {
				intlog.Errorf(ctx, `retrieve config data map failed: %+v`, err)
			}
			if _, v := gbutil.MapPossibleItemByKey(configMap, consts.ConfigNodeNameRedis); v != nil {
				configMap = gbconv.Map(v)
			}
			if len(configMap) > 0 {
				if v, ok := configMap[group]; ok {
					if redisConfig, err = gbredis.ConfigFromMap(gbconv.Map(v)); err != nil {
						panic(err)
					}
				} else {
					intlog.Printf(ctx, `missing configuration for redis group "%s"`, group)
				}
			} else {
				intlog.Print(ctx, `missing configuration for redis: "redis" node not found`)
			}
			if redisClient, err = gbredis.New(redisConfig); err != nil {
				panic(err)
			}
			return redisClient
		}
		panic(gberror.NewCode(
			gbcode.CodeMissingConfiguration,
			`no configuration found for creating redis client`,
		))
		return nil
	})
	if result != nil {
		return result.(*gbredis.Redis)
	}
	return nil
}
