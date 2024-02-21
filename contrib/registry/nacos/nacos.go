// Package nacos implements service Registry and Discovery using nacos.
package nacos

import (
	"context"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/frame/g"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"

	"github.com/joy999/nacos-sdk-go/clients"
	"github.com/joy999/nacos-sdk-go/clients/naming_client"
	"github.com/joy999/nacos-sdk-go/common/constant"
	"github.com/joy999/nacos-sdk-go/vo"
)

const (
	cstServiceSeparator = "@@"
)

var (
	_ gbsvc.Registry = &Registry{}
)

// Registry is nacos registry.
type Registry struct {
	client      naming_client.INamingClient
	clusterName string
	groupName   string
}

// Config is the configuration object for nacos client.
type Config struct {
	ServerConfigs []constant.ServerConfig `v:"required"` // See constant.ServerConfig
	ClientConfig  *constant.ClientConfig  `v:"required"` // See constant.ClientConfig
}

// New a registry with address and opts
func New(address string, opts ...constant.ClientOption) (reg *Registry) {
	endpoints := gbstr.SplitAndTrim(address, ",")
	if len(endpoints) == 0 {
		panic(gberror.NewCodef(gbcode.CodeInvalidParameter, `invalid nacos address "%s"`, address))
	}

	clientConfig := constant.NewClientConfig(opts...)

	if len(clientConfig.NamespaceId) == 0 {
		clientConfig.NamespaceId = "public"
	}

	serverConfigs := make([]constant.ServerConfig, 0, len(endpoints))
	for _, endpoint := range endpoints {
		tmp := gbstr.Split(endpoint, ":")
		ip := tmp[0]
		port := gbconv.Uint64(tmp[1])
		if port == 0 {
			port = 8848
		}
		serverConfigs = append(serverConfigs, *constant.NewServerConfig(ip, port))
	}
	ctx := gbctx.New()
	reg, err := NewWithConfig(ctx, Config{
		ServerConfigs: serverConfigs,
		ClientConfig:  clientConfig,
	})

	if err != nil {
		panic(gberror.Wrap(err, `create nacos client failed`))
	}
	return
}

// NewWithConfig creates and returns registry with Config.
func NewWithConfig(ctx context.Context, config Config) (reg *Registry, err error) {
	// Data validation.
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		return nil, err
	}

	nameingClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  config.ClientConfig,
		ServerConfigs: config.ServerConfigs,
	})
	if err != nil {
		return
	}
	return NewWithClient(nameingClient), nil
}

// NewWithClient new the instance with INamingClient
func NewWithClient(client naming_client.INamingClient) *Registry {
	r := &Registry{
		client:      client,
		clusterName: "DEFAULT",
		groupName:   "DEFAULT_GROUP",
	}
	return r
}

// SetClusterName can set the clusterName. The default is 'DEFAULT'
func (reg *Registry) SetClusterName(clusterName string) *Registry {
	reg.clusterName = clusterName
	return reg
}

// SetGroupName can set the groupName. The default is 'DEFAULT_GROUP'
func (reg *Registry) SetGroupName(groupName string) *Registry {
	reg.groupName = groupName
	return reg
}
