package gbclient

import (
	"context"
	gbmap "github.com/Ghostbb-io/gb/container/gb_map"
	"github.com/Ghostbb-io/gb/internal/intlog"
	gbsel "github.com/Ghostbb-io/gb/net/gb_sel"
	gbsvc "github.com/Ghostbb-io/gb/net/gb_svc"
	"net/http"
)

type discoveryNode struct {
	service gbsvc.Service
	address string
}

// Service is the client discovery service.
func (n *discoveryNode) Service() gbsvc.Service {
	return n.service
}

// Address returns the address of the node.
func (n *discoveryNode) Address() string {
	return n.address
}

var clientSelectorMap = gbmap.New(true)

// internalMiddlewareDiscovery is a client middleware that enables service discovery feature for client.
func internalMiddlewareDiscovery(c *Client, r *http.Request) (response *Response, err error) {
	if c.discovery == nil {
		return c.Next(r)
	}
	var (
		ctx     = r.Context()
		service gbsvc.Service
	)
	service, err = gbsvc.GetAndWatchWithDiscovery(ctx, c.discovery, r.URL.Host, func(service gbsvc.Service) {
		intlog.Printf(ctx, `http client watching service "%s" changed`, service.GetPrefix())
		if v := clientSelectorMap.Get(service.GetPrefix()); v != nil {
			if err = updateSelectorNodesByService(ctx, v.(gbsel.Selector), service); err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
		}
	})
	if err != nil {
		return nil, err
	}
	if service == nil {
		return c.Next(r)
	}
	// Balancer.
	var (
		selectorMapKey   = service.GetPrefix()
		selectorMapValue = clientSelectorMap.GetOrSetFuncLock(selectorMapKey, func() interface{} {
			intlog.Printf(ctx, `http client create selector for service "%s"`, selectorMapKey)
			selector := c.builder.Build()
			// Update selector nodes.
			if err = updateSelectorNodesByService(ctx, selector, service); err != nil {
				return nil
			}
			return selector
		})
	)
	if err != nil {
		return nil, err
	}
	selector := selectorMapValue.(gbsel.Selector)
	// Pick one node from multiple addresses.
	node, done, err := selector.Pick(ctx)
	if err != nil {
		return nil, err
	}
	if done != nil {
		defer done(ctx, gbsel.DoneInfo{})
	}
	r.Host = node.Address()
	r.URL.Host = node.Address()
	return c.Next(r)
}

func updateSelectorNodesByService(ctx context.Context, selector gbsel.Selector, service gbsvc.Service) error {
	nodes := make(gbsel.Nodes, 0)
	for _, endpoint := range service.GetEndpoints() {
		nodes = append(nodes, &discoveryNode{
			service: service,
			address: endpoint.String(),
		})
	}
	return selector.Update(ctx, nodes)
}
