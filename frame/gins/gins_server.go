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

// Server returns an instance of http server with specified name.
// Note that it panics if any error occurs duration instance creating.
func Server(name ...interface{}) *gbhttp.Server {
	var (
		err          error
		ctx          = context.Background()
		instanceName = gbhttp.DefaultServerName
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = gbconv.String(name[0])
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameServer, instanceName)
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

			// Server logger configuration checks.
			serverLoggerConfigMap = Config().MustGet(
				ctx,
				fmt.Sprintf(`%s.%s.%s`, configNodeName, instanceName, consts.ConfigNodeNameLogger),
			).Map()
			if len(serverLoggerConfigMap) == 0 && len(serverConfigMap) > 0 {
				serverLoggerConfigMap = gbconv.Map(serverConfigMap[consts.ConfigNodeNameLogger])
			}
			if len(serverLoggerConfigMap) > 0 {
				if err = server.Logger().SetConfigWithMap(serverLoggerConfigMap); err != nil {
					panic(err)
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
		}
		// The server name is necessary. It sets a default server name is it is not configured.
		if server.GetName() == "" || server.GetName() == gbhttp.DefaultServerName {
			server.SetName(instanceName)
		}
		return server
	}).(*gbhttp.Server)
}
