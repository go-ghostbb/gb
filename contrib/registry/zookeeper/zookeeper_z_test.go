package zookeeper_test

import (
	"context"
	"ghostbb.io/gb/contrib/registry/zookeeper"
	"ghostbb.io/gb/frame/g"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	"testing"
	"time"
)

// TestRegistry TestRegistryManyService
func TestRegistry(t *testing.T) {
	r := zookeeper.New([]string{"127.0.0.1:2181"}, zookeeper.WithRootPath("/gogb"))
	ctx := context.Background()

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-0-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s, err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(ctx, s)
	if err != nil {
		t.Fatal(err)
	}
}

// TestRegistryMany TestRegistryManyService
func TestRegistryMany(t *testing.T) {
	r := zookeeper.New([]string{"127.0.0.1:2181"}, zookeeper.WithRootPath("/gogb"))

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-1-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}
	svc1 := &gbsvc.LocalService{
		Name:      "ghostbb-provider-2-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9001"),
	}
	svc2 := &gbsvc.LocalService{
		Name:      "ghostbb-provider-3-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9002"),
	}

	s0, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}

	s2, err := r.Register(context.Background(), svc2)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), s0)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), s1)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), s2)
	if err != nil {
		t.Fatal(err)
	}
}

// TestGetService Test GetService
func TestGetService(t *testing.T) {
	r := zookeeper.New([]string{"127.0.0.1:2181"}, zookeeper.WithRootPath("/gogb"))
	ctx := context.Background()

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-4-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s, err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	serviceInstances, err := r.Search(ctx, gbsvc.SearchInput{
		Prefix:   s.GetPrefix(),
		Name:     svc.Name,
		Version:  svc.Version,
		Metadata: svc.Metadata,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range serviceInstances {
		g.Log().Info(ctx, instance)
	}

	err = r.Deregister(ctx, s)
	if err != nil {
		t.Fatal(err)
	}
}

// TestWatch Test Watch
func TestWatch(t *testing.T) {
	r := zookeeper.New([]string{"127.0.0.1:2181"}, zookeeper.WithRootPath("/gogb"))

	ctx := gbctx.New()

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-4-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}
	t.Log("watch")
	watch, err := r.Watch(context.Background(), svc.GetPrefix())
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}
	// watch svc
	// svc register, AddEvent
	next, err := watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output one instance
		g.Log().Info(ctx, "Register Proceed service: ", instance)
	}

	err = r.Deregister(context.Background(), s1)
	if err != nil {
		t.Fatal(err)
	}

	// svc deregister, DeleteEvent
	next, err = watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		g.Log().Info(ctx, "Deregister Proceed service: ", instance)
	}

	err = watch.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, err = watch.Proceed()
	if err == nil {
		// if nil, stop failed
		t.Fatal()
	}
}
