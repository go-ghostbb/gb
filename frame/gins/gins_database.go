package gins

import (
	"context"
	"fmt"
	gbdb "ghostbb.io/gb/database/gb_db"
	"ghostbb.io/gb/internal/consts"
	"ghostbb.io/gb/internal/instance"
	"ghostbb.io/gb/internal/intlog"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
)

func Database(name ...string) *gbdb.DB {
	var (
		err          error
		ctx          = context.Background()
		instanceName = gbdb.DefaultGroupName
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameDatabase, instanceName)
	return instance.GetOrSetFuncLock(instanceKey, func() interface{} {
		if !Config().Available(ctx) {
			return nil
		}

		var (
			configMap         map[string]interface{}
			dbConfigMap       map[string]interface{}
			dbLoggerConfigMap map[string]interface{}
			configNodeName    string
			db                *gbdb.DB
			dbConfig          gbdb.DatabaseConfig
		)

		if configMap, err = Config().Data(ctx); err != nil {
			intlog.Errorf(ctx, `retrieve config data map failed: %+v`, err)
		}
		// Find possible server configuration item by possible names.
		if len(configMap) > 0 {
			if v, _ := gbutil.MapPossibleItemByKey(configMap, consts.ConfigNodeNameDatabase); v != "" {
				configNodeName = v
			}
		}
		// Automatically retrieve configuration by instance name.
		dbConfigMap = Config().MustGet(
			ctx,
			fmt.Sprintf(`%s.%s`, configNodeName, instanceName),
		).Map()
		if len(dbConfigMap) == 0 {
			dbConfigMap = Config().MustGet(ctx, configNodeName).Map()
		}
		if len(dbConfigMap) > 0 {
			if dbConfig, err = parseDatabaseConfig(dbConfigMap); err != nil {
				panic(err)
			}

			// Database logger configuration checks.
			dbLoggerConfigMap = Config().MustGet(
				ctx,
				fmt.Sprintf(`%s.%s.%s`, configNodeName, instanceName, consts.ConfigNodeNameLogger),
			).Map()
			if len(dbLoggerConfigMap) == 0 && len(dbConfigMap) > 0 {
				dbLoggerConfigMap = gbconv.Map(dbConfigMap[consts.ConfigNodeNameLogger])
			}
			if len(dbLoggerConfigMap) > 0 {
				if err = dbConfig.Logger.SetConfigWithMap(dbLoggerConfigMap); err != nil {
					panic(err)
				}
			}

			// path
			if k, _ := gbutil.MapPossibleItemByKey(dbConfigMap, "LogPath"); k == "" {
				dbConfig.LogPath = dbConfig.Logger.GetPath()
			}

			if db, err = gbdb.NewDBByConfig(instanceName, dbConfig); err != nil {
				panic(err)
			}
		} else {
			// The configuration is not necessary, so it just prints internal logs.
			intlog.Printf(
				ctx,
				`missing configuration from configuration component for database "%s"`,
				instanceName,
			)
		}

		return db
	}).(*gbdb.DB)
}

func parseDatabaseConfig(m map[string]interface{}) (gbdb.DatabaseConfig, error) {
	// The m now is a shallow copy of m.
	// Any changes to m does not affect the original one.
	// A little tricky, isn't it?
	m = gbutil.MapCopy(m)

	config := gbdb.NewConfig()

	// Update the current configuration object.
	// It only updates the configured keys not all the object.
	if err := gbconv.Struct(m, &config); err != nil {
		return gbdb.DatabaseConfig{}, err
	}
	return config, nil
}
