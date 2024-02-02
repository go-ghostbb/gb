package gbdb

import (
	gbregex "ghostbb.io/gb/text/gb_regex"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"sync"
	"time"
)

func init() {
	configs.config = make(Config)
	configs.group = DefaultGroupName
}

// Config is the configuration management object.
type Config map[string]ConfigGroup

// ConfigGroup is a slice of configuration node for specified named group.
type ConfigGroup []ConfigNode

type ConfigNode struct {
	Host                      string        `json:"host"`                      // Host of server, ip or domain like: 127.0.0.1, localhost
	Port                      string        `json:"port"`                      // Port, it's commonly 3306.
	User                      string        `json:"user"`                      // Authentication username.
	Pass                      string        `json:"pass"`                      // Authentication password.
	Name                      string        `json:"name"`                      // Default used database name.
	Type                      string        `json:"type"`                      // Database type: mysql, sqlite, mssql, pgsql, oracle.
	Link                      string        `json:"link"`                      // (Optional) Custom link information for all configuration in one single string.
	Extra                     string        `json:"extra"`                     // (Optional) Extra configuration according the registered third-party database driver.
	Role                      string        `json:"role"`                      // (Optional, "master" in default) Node role, used for master-slave mode: master, slave.
	DryRun                    bool          `json:"dryRun"`                    // (Optional) Dry run, which does SELECT but no INSERT/UPDATE/DELETE statements.
	Weight                    int           `json:"weight"`                    // (Optional) Weight for load balance calculating, it's useless if there's just one node.
	Charset                   string        `json:"charset"`                   // (Optional, "utf8" in default) Custom charset when operating on database.
	Protocol                  string        `json:"protocol"`                  // (Optional, "tcp" in default) See net.Dial for more information which networks are available.
	Namespace                 string        `json:"namespace"`                 // (Optional) Namespace for some databases. Eg, in pgsql, the `Name` acts as the `catalog`, the `NameSpace` acts as the `schema`.
	TablePrefix               string        `json:"tablePrefix"`               // (Optional) Table name prefix.
	SingularTable             bool          `json:"singularTable"`             // (Optional) Use singular table name, table for `User` would be `user` with this option enabled
	NoLowerCase               bool          `json:"noLowerCase"`               // (Optional) Skip the snake_casing of names.
	MaxIdleConnCount          int           `json:"maxIdle"`                   // (Optional) Max idle connection configuration for underlying connection pool.
	MaxOpenConnCount          int           `json:"maxOpen"`                   // (Optional) Max open connection configuration for underlying connection pool.
	SlowThreshold             time.Duration `json:"slowThreshold"`             // (Optional) Slow threshold.
	IgnoreRecordNotFoundError bool          `json:"ignoreRecordNotFoundError"` // (Optional) Ignore record not found error.
	LogCat                    string        `json:"logCat"`                    // (Optional) Log cat.
	LogStdout                 bool          `json:"logStdout"`                 // (Optional) Stdout.
	CacheLevel                string        `json:"cacheLevel"`                // (Optional) Cache level
	CacheTTL                  time.Duration `json:"ttl"`
	CacheMaxItemCnt           int64         `json:"maxItemCnt"`
}

const (
	DefaultGroupName = "default" // Default group name.
)

// configs specifies internal used configuration object.
var configs struct {
	sync.RWMutex
	config Config // All configurations.
	group  string // Default configuration group.
}

// SetConfig sets the global configuration for package.
// It will overwrite the old configuration of package.
func SetConfig(config Config) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	for k, nodes := range config {
		for i, node := range nodes {
			nodes[i] = parseConfigNode(node)
		}
		config[k] = nodes
	}
	configs.config = config
}

// SetConfigGroup sets the configuration for given group.
func SetConfigGroup(group string, nodes ConfigGroup) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	for i, node := range nodes {
		nodes[i] = parseConfigNode(node)
	}
	configs.config[group] = nodes
}

// AddConfigNode adds one node configuration to configuration of given group.
func AddConfigNode(group string, node ConfigNode) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.config[group] = append(configs.config[group], parseConfigNode(node))
}

// GetConfig retrieves and returns the configuration of given group.
func GetConfig(group string) ConfigGroup {
	configs.RLock()
	defer configs.RUnlock()
	return configs.config[group]
}

// SetDefaultGroup sets the group name for default configuration.
func SetDefaultGroup(name string) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.group = name
}

// GetDefaultGroup returns the { name of default configuration.
func GetDefaultGroup() string {
	defer instances.Clear()
	configs.RLock()
	defer configs.RUnlock()
	return configs.group
}

// IsConfigured checks and returns whether the database configured.
// It returns true if any configuration exists.
func IsConfigured() bool {
	configs.RLock()
	defer configs.RUnlock()
	return len(configs.config) > 0
}

// parseConfigNode parses `Link` configuration syntax.
func parseConfigNode(node ConfigNode) ConfigNode {
	if node.Link != "" {
		node = *parseConfigNodeLink(&node)
	}
	if node.Link != "" && node.Type == "" {
		match, _ := gbregex.MatchString(`([a-z]+):(.+)`, node.Link)
		if len(match) == 3 {
			node.Type = gbstr.Trim(match[1])
			node.Link = gbstr.Trim(match[2])
		}
	}
	return node
}

func parseConfigNodeLink(node *ConfigNode) *ConfigNode {
	var match []string
	if node.Link != "" {
		match, _ = gbregex.MatchString(linkPattern, node.Link)
		if len(match) > 5 {
			node.Type = match[1]
			node.User = match[2]
			node.Pass = match[3]
			node.Protocol = match[4]
			array := gbstr.Split(match[5], ":")
			if len(array) == 2 && node.Protocol != "file" {
				node.Host = array[0]
				node.Port = array[1]
				node.Name = match[6]
			} else {
				node.Name = match[5]
			}
			if len(match) > 6 && match[7] != "" {
				node.Extra = match[7]
			}
			node.Link = ""
		}
	}
	if node.Extra != "" {
		if m, _ := gbstr.Parse(node.Extra); len(m) > 0 {
			_ = gbconv.Struct(m, &node)
		}
	}
	// Default value checks.
	if node.Charset == "" {
		node.Charset = defaultCharset
	}
	if node.Protocol == "" {
		node.Protocol = defaultProtocol
	}
	return node
}
