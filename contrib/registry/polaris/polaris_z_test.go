package polaris_test

import (
	"context"
	"ghostbb.io/gb/contrib/registry/polaris"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbstr "ghostbb.io/gb/text/gb_str"
	"os"
	"testing"
	"time"

	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
)

// TestRegistry_Register TestRegistryManyService
func TestRegistry_Register(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-registry/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-registry/log"); err != nil {
		t.Fatal(err)
	}

	r := polaris.NewWithConfig(
		conf,
		polaris.WithTimeout(time.Second*10),
		polaris.WithTTL(100),
	)

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-register-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s); err != nil {
		t.Fatal(err)
	}
}

// TestRegistry_Deregister TestRegistryManyService
func TestRegistry_Deregister(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-registry/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-registry/log"); err != nil {
		t.Fatal(err)
	}

	r := polaris.NewWithConfig(
		conf,
		polaris.WithTimeout(time.Second*10),
		polaris.WithTTL(100),
	)

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-deregister-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s); err != nil {
		t.Fatal(err)
	}
}

// TestRegistryMany TestRegistryManyService
func TestRegistryMany(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-registry-many/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-registry-many/log"); err != nil {
		t.Fatal(err)
	}

	r := polaris.NewWithConfig(
		conf,
		polaris.WithTimeout(time.Second*10),
		polaris.WithTTL(100),
	)

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

	if err = r.Deregister(context.Background(), s0); err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s1); err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s2); err != nil {
		t.Fatal(err)
	}
}

// TestRegistry_Search Test GetService
func TestRegistry_Search(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-get-service/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-get-service/log"); err != nil {
		t.Fatal(err)
	}

	r := polaris.NewWithConfig(
		conf,
		polaris.WithTimeout(time.Second*10),
		polaris.WithTTL(100),
	)

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-4-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	serviceInstances, err := r.Search(context.Background(), gbsvc.SearchInput{
		Prefix:   s.GetPrefix(),
		Name:     svc.Name,
		Version:  svc.Version,
		Metadata: svc.Metadata,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range serviceInstances {
		t.Log(instance)
	}

	if err = r.Deregister(context.Background(), s); err != nil {
		t.Fatal(err)
	}
}

// TestRegistry_Watch Test Watch
func TestRegistry_Watch(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-watch/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-watch/log"); err != nil {
		t.Fatal(err)
	}
	r := polaris.NewWithConfig(
		conf,
		polaris.WithTimeout(time.Second*10),
		polaris.WithTTL(100),
	)

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-5-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s := &polaris.Service{
		Service: svc,
	}

	watch, err := r.Watch(context.Background(), s.GetPrefix())
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Register service success svc instance id:", s1.(*polaris.Service).ID)
	// watch svc
	time.Sleep(time.Second * 1)

	// svc register, AddEvent
	next, err := watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output one instance
		t.Log("Register Proceed service: ", instance.GetEndpoints().String())
	}

	if err = r.Deregister(context.Background(), s1); err != nil {
		t.Fatal(err)
	}

	// svc deregister, DeleteEvent
	next, err = watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		t.Log("Deregister Proceed first delete service: ", instance.GetEndpoints().String(), ", instance id: ", instance.(*polaris.Service).ID)
	}

	if err = watch.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err = watch.Proceed(); err == nil {
		// if nil, stop failed
		t.Fatal()
	}
	t.Log("Watch close success")
}

// TestWatcher_Proceed Test Watch
func TestWatcher_Proceed(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-watch/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-watch/log"); err != nil {
		t.Fatal(err)
	}
	r := polaris.NewWithConfig(
		conf,
		polaris.WithTimeout(time.Second*10),
		polaris.WithTTL(100),
	)

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-5-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s := &polaris.Service{
		Service: svc,
	}
	svc1 := &gbsvc.LocalService{
		Name:      "ghostbb-provider-5-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9001"),
	}

	watch, err := r.Watch(context.Background(), s.GetPrefix())
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Register service success svc instance id:", s1.(*polaris.Service).ID)
	s22, err := r.Register(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Register service success svc1 instance id:", s22.(*polaris.Service).ID)
	// watch svc
	time.Sleep(time.Second * 1)

	// svc register, AddEvent
	next, err := watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output one instance
		t.Log("Register Proceed service: ", instance.GetEndpoints().String())
	}

	if err = r.Deregister(context.Background(), s1); err != nil {
		t.Fatal(err)
	}

	// svc deregister, DeleteEvent
	next, err = watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		t.Log("Deregister Proceed first delete service: ", instance.GetEndpoints().String(), ", instance id: ", instance.(*polaris.Service).ID)
	}

	// ReRegister
	s1, err = r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Register service Regin register svc instance id:", s1.(*polaris.Service).ID)
	// svc deregister, DeleteEvent
	next, err = watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		t.Log("Deregister Proceed second register service: ", instance.GetEndpoints().String(), ", instance id: ", instance.(*polaris.Service).ID)
	}

	if err = r.Deregister(context.Background(), s22); err != nil {
		t.Fatal(err)
	}

	// svc deregister, DeleteEvent
	next, err = watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		t.Log("Deregister Proceed second delete service: ", instance.GetEndpoints().String(), ", instance id: ", instance.(*polaris.Service).ID)
	}

	// svc register, deleteEvent Deregister s1
	if err = r.Deregister(context.Background(), s1); err != nil {
		t.Fatal(err)
	}

	// svc deregister, DeleteEvent
	next, err = watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		t.Log("Deregister Proceed third delete service: ", instance.GetEndpoints().String(), ", instance id: ", instance.(*polaris.Service).ID)
	}

	if err = watch.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err = watch.Proceed(); err == nil {
		// if nil, stop failed
		t.Fatal()
	}
	t.Log("Watch close success")
}

// BenchmarkRegister
func BenchmarkRegister(b *testing.B) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-registry/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-registry/log"); err != nil {
		b.Fatal(err)
	}

	r := polaris.NewWithConfig(
		conf,
		polaris.WithTimeout(time.Second*10),
		polaris.WithTTL(100),
	)

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-0-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}
	for i := 0; i < b.N; i++ {
		s, err := r.Register(context.Background(), svc)
		if err != nil {
			b.Fatal(err)
		}

		if err = r.Deregister(context.Background(), s); err != nil {
			b.Fatal(err)
		}
	}
}

// TestRegistryManyForEndpoints TestRegistryManyForEndpointsService
func TestRegistryManyForEndpoints(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-registry-many/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-registry-many/log"); err != nil {
		t.Fatal(err)
	}

	r := polaris.NewWithConfig(
		conf,
		polaris.WithTimeout(time.Second*10),
		polaris.WithTTL(100),
	)

	var (
		serviceName   = "ghostbb-provider-tcp"
		version       = "latest"
		endpointOne   = "127.0.0.1:9000"
		endpointTwo   = "127.0.0.1:9001"
		endpointThree = "127.0.0.1:9002"
	)

	svc := &gbsvc.LocalService{
		Name:      serviceName,
		Version:   version,
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints(endpointOne),
	}

	svc1 := &gbsvc.LocalService{
		Name:      serviceName,
		Version:   version,
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints(endpointTwo),
	}

	svc2 := &gbsvc.LocalService{
		Name:      serviceName,
		Version:   version,
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints(endpointThree),
	}

	// svc register, AddEvent
	s0, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	// svc register, AddEvent
	s1, err := r.Register(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}

	// svc register, AddEvent
	s2, err := r.Register(context.Background(), svc2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Register service success sleep 1s")
	time.Sleep(time.Second * 2)

	// serviceName = "service-default-default-ghostbb-provider-tcp-latest"
	result, err := r.Search(context.Background(), gbsvc.SearchInput{
		Name: serviceName,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Search service success size:", len(result))
	for i := 0; i < len(result); i++ {
		t.Log("Endpoints:", result[i].GetEndpoints().String())
		if !gbstr.Contains(result[i].GetEndpoints().String(), endpointOne) {
			t.Fatal("endpointOne not found")
		}
		if !gbstr.Contains(result[i].GetEndpoints().String(), endpointTwo) {
			t.Fatal("endpointTwo not found")
		}
		if !gbstr.Contains(result[i].GetEndpoints().String(), endpointThree) {
			t.Fatal("endpointThree not found")
		}
	}
	t.Log("Search service success sleep 1s")
	time.Sleep(time.Second * 1)
	if err = r.Deregister(context.Background(), s0); err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s1); err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s2); err != nil {
		t.Fatal(err)
	}

	t.Log("Deregister success")
}

// TestWatcher_Close Test Close
func TestWatcher_Close(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-watch/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-watch/log"); err != nil {
		t.Fatal(err)
	}
	r := polaris.NewWithConfig(
		conf,
		polaris.WithTimeout(time.Second*10),
		polaris.WithTTL(100),
	)

	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-close-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s := &polaris.Service{
		Service: svc,
	}

	watch, err := r.Watch(context.Background(), s.GetPrefix())
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	// watch svc
	time.Sleep(time.Second * 1)
	if err = r.Deregister(context.Background(), s1); err != nil {
		t.Fatal(err)
	}

	if err = watch.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err = watch.Proceed(); err == nil {
		// if nil, stop failed
		t.Fatal()
	}
	t.Log("Watch close success")
}

// TestGetKey Test get key
func TestGetKey(t *testing.T) {
	svc := &gbsvc.LocalService{
		Name:      "ghostbb-provider-key-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
		Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s := &polaris.Service{
		Service: svc,
	}
	if s.GetKey() != "service-default-default-ghostbb-provider-key-tcp-test-127.0.0.1:9000" {
		t.Fatal("GetKey error key:", s.GetKey())
	}
	t.Log("GetKey success ")
}

// TestService_GetPrefix Test GetPrefix
func TestService_GetPrefix(t *testing.T) {
	type fields struct {
		Service gbsvc.Service
		ID      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestService_GetPrefix-0",
			fields: fields{
				Service: &gbsvc.LocalService{
					Name:      "ghostbb-provider-0-tcp",
					Version:   "test",
					Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
					Endpoints: gbsvc.NewEndpoints("127.0.0.1:9000"),
				},
				ID: "test",
			},
			want: "service-default-default-ghostbb-provider-0-tcp-test",
		},
		{
			name: "TestService_GetPrefix-1",
			fields: fields{
				Service: &gbsvc.LocalService{
					Name:      "ghostbb-provider-1-tcp",
					Version:   "test",
					Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
					Endpoints: gbsvc.NewEndpoints("127.0.0.1:9001"),
				},
				ID: "test",
			},
			want: "service-default-default-ghostbb-provider-1-tcp-test",
		},
		{
			name: "TestService_GetPrefix-2",
			fields: fields{
				Service: &gbsvc.LocalService{
					Name:      "ghostbb-provider-2-tcp",
					Version:   "latest",
					Metadata:  map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
					Endpoints: gbsvc.NewEndpoints("127.0.0.1:9002"),
				},
				ID: "latest",
			},
			want: "service-default-default-ghostbb-provider-2-tcp-latest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &polaris.Service{
				Service: tt.fields.Service,
				ID:      tt.fields.ID,
			}
			if got := s.GetPrefix(); got != tt.want {
				t.Errorf("GetPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestService_GetName Test GetName
func TestService_GetKey(t *testing.T) {
	type fields struct {
		Service gbsvc.Service
		ID      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestService_GetKey-0",
			fields: fields{
				Service: &gbsvc.LocalService{
					Namespace:  gbsvc.DefaultNamespace,
					Deployment: gbsvc.DefaultDeployment,
					Name:       "ghostbb-provider-0-tcp",
					Version:    "test",
					Metadata:   map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
					Endpoints:  gbsvc.NewEndpoints("127.0.0.1:9000"),
				},
				ID: "test",
			},
			want: "service-default-default-ghostbb-provider-0-tcp-test-127.0.0.1:9000",
		},
		{
			name: "TestService_GetKey-1",
			fields: fields{
				Service: &gbsvc.LocalService{
					Namespace:  gbsvc.DefaultNamespace,
					Deployment: gbsvc.DefaultDeployment,
					Name:       "ghostbb-provider-1-tcp",
					Version:    "latest",
					Metadata:   map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
					Endpoints:  gbsvc.NewEndpoints("127.0.0.1:9001"),
				},
				ID: "latest",
			},
			want: "service-default-default-ghostbb-provider-1-tcp-latest-127.0.0.1:9001",
		},
		{
			name: "TestService_GetKey-2",
			fields: fields{
				Service: &gbsvc.LocalService{
					Namespace:  gbsvc.DefaultNamespace,
					Deployment: gbsvc.DefaultDeployment,
					Name:       "ghostbb-provider-2-tcp",
					Version:    "latest",
					Metadata:   map[string]interface{}{"app": "ghostbb", gbsvc.MDProtocol: "tcp"},
					Endpoints:  gbsvc.NewEndpoints("127.0.0.1:9002"),
				},
				ID: "latest",
			},
			want: "service-default-default-ghostbb-provider-2-tcp-latest-127.0.0.1:9002",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &polaris.Service{
				Service: tt.fields.Service,
				ID:      tt.fields.ID,
			}
			if got := s.GetKey(); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
