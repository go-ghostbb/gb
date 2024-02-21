package file_test

import (
	"fmt"
	"ghostbb.io/gb/contrib/registry/file"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/frame/g"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbuid "ghostbb.io/gb/util/gb_uid"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

var ctx = gbctx.GetInitCtx()

func Test_HTTP_Registry(t *testing.T) {
	var (
		svcName = gbuid.S()
		dirPath = gbfile.Temp(svcName)
	)
	defer gbfile.Remove(dirPath)
	gbsvc.SetRegistry(file.New(dirPath))

	s := g.Server(svcName)
	s.GET("/http-registry", func(c *gin.Context) {
		c.String(200, svcName)
	})

	s.SetDumpRouterMap(false)
	s.SetTerminal(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gbtest.C(t, func(t *gbtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://%s", svcName))
		// GET
		t.Assert(client.GetContent(ctx, "/http-registry"), svcName)
	})
}

func Test_HTTP_Discovery_Disable(t *testing.T) {
	var (
		svcName = gbuid.S()
		dirPath = gbfile.Temp(svcName)
	)
	defer gbfile.Remove(dirPath)
	gbsvc.SetRegistry(file.New(dirPath))

	s := g.Server()
	s.GET("/http-registry", func(c *gin.Context) {
		c.String(200, svcName)
	})
	s.SetDumpRouterMap(false)
	s.SetTerminal(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gbtest.C(t, func(t *gbtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		result, err := client.Get(ctx, "/http-registry")
		defer result.Close()
		t.Assert(gberror.Code(err), gbcode.CodeNotFound)
	})
	gbtest.C(t, func(t *gbtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		result, err := client.Discovery(nil).Get(ctx, "/http-registry")
		defer result.Close()
		t.AssertNil(err)
		t.Assert(result.ReadAllString(), svcName)
	})
}

func Test_HTTP_Server_Endpoints(t *testing.T) {
	var (
		svcName = gbuid.S()
		dirPath = gbfile.Temp(svcName)
	)
	defer gbfile.Remove(dirPath)
	gbsvc.SetRegistry(file.New(dirPath))

	endpoints := []string{"10.0.0.1:8000", "10.0.0.2:8000"}
	s := g.Server(svcName)
	s.SetEndpoints(endpoints)
	s.GET("/http-registry", func(c *gin.Context) {
		c.String(200, svcName)
	})
	s.SetDumpRouterMap(false)
	s.SetTerminal(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gbtest.C(t, func(t *gbtest.T) {
		service, err := gbsvc.Get(ctx, svcName)
		t.AssertNil(err)
		t.Assert(service.GetEndpoints(), gbstr.Join(endpoints, ","))
	})
}
