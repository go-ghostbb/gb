package gbdb

import (
	"context"
	gbmap "ghostbb.io/container/gb_map"
	gbcode "ghostbb.io/errors/gb_code"
	gberror "ghostbb.io/errors/gb_error"
	"ghostbb.io/internal/intlog"
	gbcache "ghostbb.io/os/gb_cache"
	gbctx "ghostbb.io/os/gb_ctx"
	gblog "ghostbb.io/os/gb_log"
	gbrand "ghostbb.io/util/gb_rand"
	"gorm.io/gorm"
)

const (
	defaultCharset  = `utf8`
	defaultProtocol = `tcp`
	dbRoleSlave     = `slave`

	// type:[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	linkPattern = `(\w+):([\w\-\$]*):(.*?)@(\w+?)\((.+?)\)/{0,1}([^\?]*)\?{0,1}(.*)`
)

var (
	// instances is the management map for instances.
	instances = gbmap.NewStrAnyMap(true)

	// driverMap manages all custom registered driver.
	driverMap = map[string]Driver{}

	// global logger
	globalLogger = gblog.New()
)

type DecodeFunc func(*ConfigNode)

var decode DecodeFunc = nil

// Core is the base struct for database management.
type Core struct {
	*gorm.DB
	group         string           // Configuration group name.
	links         *gbmap.StrAnyMap // links caches all created links by node.
	logger        ILogger          // Logger for logging functionality.
	cache         *gbcache.Cache   // cache
	config        *ConfigNode      // Current config node.
	dynamicConfig dynamicConfig    // Dynamic configurations, which can be changed in runtime.
}

type dynamicConfig struct {
	MaxIdleConnCount int
	MaxOpenConnCount int
}

// Driver is the interface for integrating sql drivers into package gdb.
type Driver interface {
	// New creates and returns a database object for specified database server.
	New(core *Core, node *ConfigNode) (*gorm.DB, error)
}

// Register registers custom database driver to gbdb.
func Register(name string, driver Driver) error {
	driverMap[name] = newDriverWrapper(driver)
	return nil
}

func SetDecodeFunc(fn DecodeFunc) {
	decode = fn
}

func SetGlobalLoggerConfigWithMap(m map[string]interface{}) error {
	return globalLogger.SetConfigWithMap(m)
}

// NewByGroup creates and returns an ORM object with global configurations.
// The parameter `name` specifies the configuration group name,
// which is DefaultGroupName in default.
func NewByGroup(group ...string) (core *Core, err error) {
	groupName := configs.group
	if len(group) > 0 && group[0] != "" {
		groupName = group[0]
	}
	configs.RLock()
	defer configs.RUnlock()

	if len(configs.config) < 1 {
		return nil, gberror.NewCode(
			gbcode.CodeInvalidConfiguration,
			"database configuration is empty, please set the database configuration before using",
		)
	}
	if _, ok := configs.config[groupName]; ok {
		var node *ConfigNode
		if node, err = getConfigNodeByGroup(groupName, true); err == nil {
			return newDBByConfigNode(node, groupName)
		}
		return nil, err
	}
	return nil, gberror.NewCodef(
		gbcode.CodeInvalidConfiguration,
		`database configuration node "%s" is not found, did you misspell group name "%s" or miss the database configuration?`,
		groupName, groupName,
	)
}

// newDBByConfigNode creates and returns an ORM object with given configuration node and group name.
func newDBByConfigNode(node *ConfigNode, group string) (core *Core, err error) {
	if decode != nil {
		decode(node)
	}

	if node.Link != "" {
		node = parseConfigNodeLink(node)
	}
	c := &Core{
		group:  group,
		links:  gbmap.NewStrAnyMap(true),
		logger: NewLogger(group, globalLogger.Clone()),
		config: node,
		dynamicConfig: dynamicConfig{
			MaxIdleConnCount: node.MaxIdleConnCount,
			MaxOpenConnCount: node.MaxOpenConnCount,
		},
	}

	c.logger.SetSlowThreshold(node.SlowThreshold)
	c.logger.SetIgnoreRecordNotFoundError(node.IgnoreRecordNotFoundError)
	c.logger.SetLogCat(node.LogCat)
	c.logger.SetLogStdout(node.LogStdout)

	if v, ok := driverMap[node.Type]; ok {
		if c.DB, err = v.New(c, node); err != nil {
			return nil, err
		}
		intlog.Printf(gbctx.New(), "%s | %s | database connection successful.", group, node.Type)
		return c, nil
	}
	errorMsg := `cannot find database driver for specified database type "%s"`
	errorMsg += `, did you misspell type name "%s" or forget importing the database driver? `
	return nil, gberror.NewCodef(gbcode.CodeInvalidConfiguration, errorMsg, node.Type, node.Type)
}

// getConfigNodeByGroup calculates and returns a configuration node of given group. It
// calculates the value internally using weight algorithm for load balance.
//
// The parameter `master` specifies whether retrieving a master node, or else a slave node
// if master-slave configured.
func getConfigNodeByGroup(group string, master bool) (*ConfigNode, error) {
	if list, ok := configs.config[group]; ok {
		// Separates master and slave configuration nodes array.
		var (
			masterList = make(ConfigGroup, 0)
			slaveList  = make(ConfigGroup, 0)
		)
		for i := 0; i < len(list); i++ {
			if list[i].Role == dbRoleSlave {
				slaveList = append(slaveList, list[i])
			} else {
				masterList = append(masterList, list[i])
			}
		}
		if len(masterList) < 1 {
			return nil, gberror.NewCode(
				gbcode.CodeInvalidConfiguration,
				"at least one master node configuration's need to make sense",
			)
		}
		if len(slaveList) < 1 {
			slaveList = masterList
		}
		if master {
			return getConfigNodeByWeight(masterList), nil
		} else {
			return getConfigNodeByWeight(slaveList), nil
		}
	}
	return nil, gberror.NewCodef(
		gbcode.CodeInvalidConfiguration,
		"empty database configuration for item name '%s'",
		group,
	)
}

// getConfigNodeByWeight calculates the configuration weights and randomly returns a node.
//
// Calculation algorithm brief:
// 1. If we have 2 nodes, and their weights are both 1, then the weight range is [0, 199];
// 2. Node1 weight range is [0, 99], and node2 weight range is [100, 199], ratio is 1:1;
// 3. If the random number is 99, it then chooses and returns node1;.
func getConfigNodeByWeight(cg ConfigGroup) *ConfigNode {
	if len(cg) < 2 {
		return &cg[0]
	}
	var total int
	for i := 0; i < len(cg); i++ {
		total += cg[i].Weight * 100
	}
	// If total is 0 means all the nodes have no weight attribute configured.
	// It then defaults each node's weight attribute to 1.
	if total == 0 {
		for i := 0; i < len(cg); i++ {
			cg[i].Weight = 1
			total += cg[i].Weight * 100
		}
	}
	// Exclude the right border value.
	var (
		min    = 0
		max    = 0
		random = gbrand.N(0, total-1)
	)
	for i := 0; i < len(cg); i++ {
		max = min + cg[i].Weight*100
		if random >= min && random < max {
			// ====================================================
			// Return a COPY of the ConfigNode.
			// ====================================================
			node := ConfigNode{}
			node = cg[i]
			return &node
		}
		min = max
	}
	return nil
}

func WithCacheCtx() context.Context {
	return context.WithValue(gbctx.New(), DBCacheName, struct{}{})
}
