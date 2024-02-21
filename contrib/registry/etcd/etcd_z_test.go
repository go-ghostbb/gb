package etcd_test

import (
	"ghostbb.io/gb/contrib/registry/etcd"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbuid "ghostbb.io/gb/util/gb_uid"
	"testing"
)

func TestRegistry(t *testing.T) {
	var (
		ctx      = gbctx.GetInitCtx()
		registry = etcd.New(`127.0.0.1:2379`)
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
		registry = etcd.New(`127.0.0.1:2379`)
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
		watcher, err := registry.Watch(ctx, svc1.GetPrefix())
		t.AssertNil(err)

		// Register another service.
		svc2 := &gbsvc.LocalService{
			Name:      svc1.Name,
			Endpoints: gbsvc.NewEndpoints("127.0.0.1:9999"),
		}
		registered, err := registry.Register(ctx, svc2)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc2.GetName())

		// Watch and retrieve the service changes:
		// svc1 and svc2 is the same service name, which has 2 endpoints.
		proceedResult, err := watcher.Proceed()
		t.AssertNil(err)
		t.Assert(len(proceedResult), 1)
		t.Assert(
			proceedResult[0].GetEndpoints(),
			gbsvc.Endpoints{svc1.GetEndpoints()[0], svc2.GetEndpoints()[0]},
		)

		// Watch and retrieve the service changes:
		// left only svc1, which means this service has only 1 endpoint.
		err = registry.Deregister(ctx, svc2)
		t.AssertNil(err)
		proceedResult, err = watcher.Proceed()
		t.AssertNil(err)
		t.Assert(
			proceedResult[0].GetEndpoints(),
			gbsvc.Endpoints{svc1.GetEndpoints()[0]},
		)
		t.AssertNil(watcher.Close())
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := registry.Deregister(ctx, svc1)
		t.AssertNil(err)
	})
}
