package nacos_test

import (
	"context"
	"ghostbb.io/gb/contrib/registry/nacos"
	"ghostbb.io/gb/frame/g"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbuid "ghostbb.io/gb/util/gb_uid"
	"sync/atomic"
	"testing"
	"time"

	"github.com/joy999/nacos-sdk-go/common/constant"
)

const (
	NACOS_ADDRESS   = `localhost:8848`
	NACOS_CACHE_DIR = `/tmp/nacos`
	NACOS_LOG_DIR   = `/tmp/nacos`
)

func TestRegistry(t *testing.T) {
	var (
		ctx      = gbctx.GetInitCtx()
		registry = nacos.New(NACOS_ADDRESS, func(cc *constant.ClientConfig) {
			cc.CacheDir = NACOS_CACHE_DIR
			cc.LogDir = NACOS_LOG_DIR
		})
	)
	svc := &gbsvc.LocalService{
		Name:      gbuid.S(),
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:8888"),
		Metadata: map[string]interface{}{
			"protocol": "https",
		},
	}
	gbtest.C(t, func(t *gbtest.T) {
		registered, err := registry.Register(ctx, svc)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc.GetName())
	})

	// Search by name.
	gbtest.C(t, func(t *gbtest.T) {
		result, err := registry.Search(ctx, gbsvc.SearchInput{
			Name: svc.Name,
		})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0].GetName(), svc.Name)
	})

	// Search by prefix.
	gbtest.C(t, func(t *gbtest.T) {
		result, err := registry.Search(ctx, gbsvc.SearchInput{
			Prefix: svc.GetPrefix(),
		})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0].GetName(), svc.Name)
	})

	// Search by metadata.
	gbtest.C(t, func(t *gbtest.T) {
		result, err := registry.Search(ctx, gbsvc.SearchInput{
			Name: svc.GetName(),
			Metadata: map[string]interface{}{
				"protocol": "https",
			},
		})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0].GetName(), svc.Name)
	})
	gbtest.C(t, func(t *gbtest.T) {
		result, err := registry.Search(ctx, gbsvc.SearchInput{
			Name: svc.GetName(),
			Metadata: map[string]interface{}{
				"protocol": "grpc",
			},
		})
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := registry.Deregister(ctx, svc)
		t.AssertNil(err)
	})
}

func TestWatch(t *testing.T) {
	var (
		ctx      = gbctx.GetInitCtx()
		registry = nacos.New(NACOS_ADDRESS, func(cc *constant.ClientConfig) {
			cc.CacheDir = NACOS_CACHE_DIR
			cc.LogDir = NACOS_LOG_DIR
		})
		registry2 = nacos.New(NACOS_ADDRESS, func(cc *constant.ClientConfig) {
			cc.CacheDir = NACOS_CACHE_DIR
			cc.LogDir = NACOS_LOG_DIR
		})
	)

	svc1 := &gbsvc.LocalService{
		Name:      gbuid.S(),
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:8888"),
		Metadata: map[string]interface{}{
			"protocol": "https",
		},
	}
	gbtest.C(t, func(t *gbtest.T) {
		registered, err := registry.Register(ctx, svc1)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc1.GetName())
	})

	gbtest.C(t, func(t *gbtest.T) {
		ctx := gbctx.New()
		watcher, err := registry.Watch(ctx, svc1.GetPrefix())
		t.AssertNil(err)

		var latestProceedResult atomic.Value
		g.Go(ctx, func(ctx context.Context) {
			var (
				err error
				res []gbsvc.Service
			)
			for err == nil {
				res, err = watcher.Proceed()
				t.AssertNil(err)
				latestProceedResult.Store(res)
			}
		}, func(ctx context.Context, exception error) {
			t.Fatal(exception)
		})

		// Register another service.
		svc2 := &gbsvc.LocalService{
			Name:      svc1.Name,
			Endpoints: gbsvc.NewEndpoints("127.0.0.2:9999"),
			Metadata: map[string]interface{}{
				"protocol": "https",
			},
		}
		registered, err := registry2.Register(ctx, svc2)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc2.GetName())

		time.Sleep(time.Second * 10)

		// Watch and retrieve the service changes:
		// svc1 and svc2 is the same service name, which has 2 endpoints.
		proceedResult, ok := latestProceedResult.Load().([]gbsvc.Service)
		t.Assert(ok, true)
		t.Assert(len(proceedResult), 1)
		t.Assert(
			allEndpoints(proceedResult),
			gbsvc.Endpoints{svc1.GetEndpoints()[0], svc2.GetEndpoints()[0]},
		)

		// Watch and retrieve the service changes:
		// left only svc1, which means this service has only 1 endpoint.
		err = registry2.Deregister(ctx, svc2)
		t.AssertNil(err)

		time.Sleep(time.Second * 10)
		proceedResult, ok = latestProceedResult.Load().([]gbsvc.Service)
		t.Assert(ok, true)
		t.Assert(len(proceedResult), 1)
		t.Assert(
			allEndpoints(proceedResult),
			gbsvc.Endpoints{svc1.GetEndpoints()[0]},
		)
		t.AssertNil(watcher.Close())
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := registry.Deregister(ctx, svc1)
		t.AssertNil(err)
	})
}

func allEndpoints(services []gbsvc.Service) gbsvc.Endpoints {
	m := map[gbsvc.Endpoint]struct{}{}
	for _, s := range services {
		for _, ep := range s.GetEndpoints() {
			m[ep] = struct{}{}
		}
	}
	var endpoints gbsvc.Endpoints
	for ep := range m {
		endpoints = append(endpoints, ep)
	}
	return sortEndpoints(endpoints)
}

func sortEndpoints(in gbsvc.Endpoints) gbsvc.Endpoints {
	var endpoints gbsvc.Endpoints
	endpoints = append(endpoints, in...)
	n := len(endpoints)
	for i := 0; i < n; i++ {
		for t := i; t < n; t++ {
			if endpoints[i].String() > endpoints[t].String() {
				endpoints[i], endpoints[t] = endpoints[t], endpoints[i]
			}
		}
	}
	return endpoints
}
