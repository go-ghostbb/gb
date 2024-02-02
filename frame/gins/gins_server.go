package gins

import (
	"context"
	"fmt"
	"ghostbb.io/gb/internal/consts"
	"ghostbb.io/gb/internal/instance"
	"ghostbb.io/gb/internal/intlog"
	gbhttp "ghostbb.io/gb/net/gb_http"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
)

func Server(name ...interface{}) *gbhttp.Server {
	var (
		err          error
		ctx          = context.Background()
		instanceName = gbhttp.DefaultServerName
		instanceKey  = fmt.Sprintf("%s.%v", frameCoreComponentNameServer, name)
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = gbconv.String(name[0])
	}

	return instance.GetOrSetFuncLock(instanceKey, func() interface{} {
		server := gbhttp.GetServer(instanceName)
		if Config().Available(ctx) {
			// Server initialization from configuration.
			var (
				configMap             map[string]interface{}
				serverConfigMap       map[string]interface{}
				serverLoggerConfigMap map[string]interface{}
				configNodeName        string
			)
			if configMap, err = Config().Data(ctx); err != nil {
				intlog.Errorf(ctx, `retrieve config data map failed: %+v`, err)
			}
			// Find possible server configuration item by possible names.
			if len(configMap) > 0 {
				if v, _ := gbutil.MapPossibleItemByKey(configMap, consts.ConfigNodeNameServer); v != "" {
					configNodeName = v
				}
				if configNodeName == "" {
					if v, _ := gbutil.MapPossibleItemByKey(configMap, consts.ConfigNodeNameServerSecondary); v != "" {
						configNodeName = v
					}
				}
			}
			// Automatically retrieve configuration by instance name.
			serverConfigMap = Config().MustGet(
				ctx,
				fmt.Sprintf(`%s.%s`, configNodeName, instanceName),
			).Map()
			if len(serverConfigMap) == 0 {
				serverConfigMap = Config().MustGet(ctx, configNodeName).Map()
			}
			if len(serverConfigMap) > 0 {
				if err = server.SetConfigWithMap(serverConfigMap); err != nil {
					panic(err)
				}
			} else {
				// The configuration is not necessary, so it just prints internal logs.
				intlog.Printf(
					ctx,
					`missing configuration from configuration component for HTTP server "%s"`,
					instanceName,
				)
			}
			// Server logger configuration checks.
			serverLoggerConfigMap = Config().MustGet(
				ctx,
				fmt.Sprintf(`%s.%s`, configNodeName, consts.ConfigNodeNameLogger),
			).Map()
			if len(serverLoggerConfigMap) > 0 {
				if err = server.Logger().SetConfigWithMap(serverLoggerConfigMap); err != nil {
					panic(err)
				}
			}
		}
		return server
	}).(*gbhttp.Server)
}
