package gins

import (
	"context"
	"fmt"
	"ghostbb.io/gb/internal/consts"
	"ghostbb.io/gb/internal/instance"
	gblog "ghostbb.io/gb/os/gb_log"
	gbutil "ghostbb.io/gb/util/gb_util"
)

// Log returns an instance of gblog.Logger.
// The parameter `name` is the name for the instance.
// Note that it panics if any error occurs duration instance creating.
func Log(name ...string) *gblog.Logger {
	var (
		ctx          = context.Background()
		instanceName = gblog.DefaultName
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameLogger, instanceName)
	return instance.GetOrSetFuncLock(instanceKey, func() interface{} {
		logger := gblog.Instance(instanceName)
		// To avoid file no found error while it's not necessary.
		var (
			configMap      map[string]interface{}
			loggerNodeName = consts.ConfigNodeNameLogger
		)
		// Try to find possible `loggerNodeName` in case-insensitive way.
		if configData, _ := Config().Data(ctx); len(configData) > 0 {
			if v, _ := gbutil.MapPossibleItemByKey(configData, consts.ConfigNodeNameLogger); v != "" {
				loggerNodeName = v
			}
		}
		// Retrieve certain logger configuration by logger name.
		certainLoggerNodeName := fmt.Sprintf(`%s.%s`, loggerNodeName, instanceName)
		if v, _ := Config().Get(ctx, certainLoggerNodeName); !v.IsEmpty() {
			configMap = v.Map()
		}
		// Retrieve global logger configuration if configuration for certain logger name does not exist.
		if len(configMap) == 0 {
			if v, _ := Config().Get(ctx, loggerNodeName); !v.IsEmpty() {
				configMap = v.Map()
			}
		}
		// Set logger config if config map is not empty.
		if len(configMap) > 0 {
			if err := logger.SetConfigWithMap(configMap); err != nil {
				panic(err)
			}
		}
		return logger
	}).(*gblog.Logger)
}
