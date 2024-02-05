package gins

import (
	"context"
	"fmt"
	gbdb "ghostbb.io/gb/database/gb_db"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/consts"
	"ghostbb.io/gb/internal/instance"
	"ghostbb.io/gb/internal/intlog"
	gbcfg "ghostbb.io/gb/os/gb_cfg"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
	"time"
)

func Database(name ...string) *gbdb.Core {
	var (
		ctx   = context.Background()
		group = gbdb.DefaultGroupName
	)
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameDatabase, group)
	db := instance.GetOrSetFuncLock(instanceKey, func() interface{} {
		// It ignores returned error to avoid file no found error while it's not necessary.
		var (
			configMap     map[string]interface{}
			configNodeKey = consts.ConfigNodeNameDatabase
		)
		// It firstly searches the configuration of the instance name.
		if configData, _ := Config().Data(ctx); len(configData) > 0 {
			if v, _ := gbutil.MapPossibleItemByKey(configData, consts.ConfigNodeNameDatabase); v != "" {
				configNodeKey = v
			}
		}
		if v, _ := Config().Get(ctx, configNodeKey); !v.IsEmpty() {
			configMap = v.Map()
		}
		// No configuration found, it formats and panics error.
		if len(configMap) == 0 && !gbdb.IsConfigured() {
			// File configuration object checks.
			var err error
			if fileConfig, ok := Config().GetAdapter().(*gbcfg.AdapterFile); ok {
				if _, err = fileConfig.GetFilePath(); err != nil {
					panic(gberror.WrapCode(gbcode.CodeMissingConfiguration, err,
						`configuration not found, did you miss the configuration file or misspell the configuration file name`,
					))
				}
			}
			// Panic if nothing found in Config object or in gdb configuration.
			if len(configMap) == 0 && !gbdb.IsConfigured() {
				panic(gberror.NewCodef(
					gbcode.CodeMissingConfiguration,
					`database initialization failed: configuration missing for database node "%s"`,
					consts.ConfigNodeNameDatabase,
				))
			}
		}

		if len(configMap) == 0 {
			configMap = make(map[string]interface{})
		}
		// Parse `m` as map-slice and adds it to global configurations for package gbdb.
		for g, groupConfig := range configMap {
			cg := gbdb.ConfigGroup{}
			switch value := groupConfig.(type) {
			case []interface{}:
				for _, v := range value {
					if node := parseDBConfigNode(v); node != nil {
						cg = append(cg, *node)
					}
				}
			case map[string]interface{}:
				if node := parseDBConfigNode(value); node != nil {
					cg = append(cg, *node)
				}
			}
			if len(cg) > 0 {
				if gbdb.GetConfig(group) == nil {
					intlog.Printf(ctx, "add configuration for group: %s, %#v", g, cg)
					gbdb.SetConfigGroup(g, cg)
				} else {
					intlog.Printf(ctx, "ignore configuration as it already exists for group: %s, %#v", g, cg)
					intlog.Printf(ctx, "%s, %#v", g, cg)
				}
			}
		}
		// Parse `m` as a single node configuration,
		// which is the default group configuration.
		if node := parseDBConfigNode(configMap); node != nil {
			cg := gbdb.ConfigGroup{}
			if node.Link != "" || node.Host != "" {
				cg = append(cg, *node)
			}
			if len(cg) > 0 {
				if gbdb.GetConfig(group) == nil {
					intlog.Printf(ctx, "add configuration for group: %s, %#v", gbdb.DefaultGroupName, cg)
					gbdb.SetConfigGroup(gbdb.DefaultGroupName, cg)
				} else {
					intlog.Printf(
						ctx,
						"ignore configuration as it already exists for group: %s, %#v",
						gbdb.DefaultGroupName, cg,
					)
					intlog.Printf(ctx, "%s, %#v", gbdb.DefaultGroupName, cg)
				}
			}
		}

		// Initialize logger
		var (
			loggerConfigMap map[string]interface{}
			loggerNodeName  = fmt.Sprintf("%s.%s", configNodeKey, consts.ConfigNodeNameLogger)
		)
		if v, _ := Config().Get(ctx, loggerNodeName); !v.IsEmpty() {
			loggerConfigMap = v.Map()
		}
		if len(loggerConfigMap) == 0 {
			if v, _ := Config().Get(ctx, configNodeKey); !v.IsEmpty() {
				loggerConfigMap = v.Map()
			}
		}
		if len(loggerConfigMap) > 0 {
			if err := gbdb.SetGlobalLoggerConfigWithMap(loggerConfigMap); err != nil {
				panic(err)
			}
		}

		// Create a new ORM object with given configurations.
		if db, err := gbdb.NewByGroup(name...); err == nil {
			return db
		} else {
			// If panics, often because it does not find its configuration for given group.
			panic(err)
		}
		return nil
	})
	if db != nil {
		return db.(*gbdb.Core)
	}
	return nil
}

func parseDBConfigNode(value interface{}) *gbdb.ConfigNode {
	nodeMap, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	var (
		node = &gbdb.ConfigNode{}
		err  = gbconv.Struct(nodeMap, node)
	)
	if err != nil {
		panic(err)
	}
	// Find possible `Link` configuration content.
	if _, v := gbutil.MapPossibleItemByKey(nodeMap, "Link"); v != nil {
		node.Link = gbconv.String(v)
	}
	if _, v := gbutil.MapPossibleItemByKey(nodeMap, "SlowThreshold"); v == nil {
		node.SlowThreshold = 200 * time.Millisecond
	}
	if _, v := gbutil.MapPossibleItemByKey(nodeMap, "IgnoreRecordNotFoundError"); v == nil {
		node.IgnoreRecordNotFoundError = true
	}
	if _, v := gbutil.MapPossibleItemByKey(nodeMap, "LogStdout"); v == nil {
		node.LogStdout = true
	}
	return node
}
