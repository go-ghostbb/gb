package gbtcp_test

import (
	"fmt"
	gbdebug "ghostbb.io/gb/debug/gb_debug"
	gbtcp "ghostbb.io/gb/net/gb_tcp"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
	"time"
)

func Test_Package_Basic(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.RecvPkg()
			if err != nil {
				break
			}
			conn.SendPkg(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	// SendPkg
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		for i := 0; i < 100; i++ {
			err := conn.SendPkg([]byte(gbconv.String(i)))
			t.AssertNil(err)
		}
		for i := 0; i < 100; i++ {
			err := conn.SendPkgWithTimeout([]byte(gbconv.String(i)), time.Second)
			t.AssertNil(err)
		}
	})
	// SendPkg with big data - failure.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 65536)
		err = conn.SendPkg(data)
		t.AssertNE(err, nil)
	})
	// SendRecvPkg
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		for i := 100; i < 200; i++ {
			data := []byte(gbconv.String(i))
			result, err := conn.SendRecvPkg(data)
			t.AssertNil(err)
			t.Assert(result, data)
		}
		for i := 100; i < 200; i++ {
			data := []byte(gbconv.String(i))
			result, err := conn.SendRecvPkgWithTimeout(data, time.Second)
			t.AssertNil(err)
			t.Assert(result, data)
		}
	})
	// SendRecvPkg with big data - failure.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 65536)
		result, err := conn.SendRecvPkg(data)
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	// SendRecvPkg with big data - success.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 65500)
		data[100] = byte(65)
		data[65400] = byte(85)
		result, err := conn.SendRecvPkg(data)
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Package_Basic_HeaderSize1(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.RecvPkg(gbtcp.PkgOption{HeaderSize: 1})
			if err != nil {
				break
			}
			conn.SendPkg(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)

	// SendRecvPkg with empty data.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0)
		result, err := conn.SendRecvPkg(data)
		t.AssertNil(err)
		t.AssertNil(result)
	})
}

func Test_Package_Timeout(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.RecvPkg()
			if err != nil {
				break
			}
			time.Sleep(time.Second)
			gbtest.Assert(conn.SendPkg(data), nil)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("10000")
		result, err := conn.SendRecvPkgWithTimeout(data, time.Millisecond*500)
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("10000")
		result, err := conn.SendRecvPkgWithTimeout(data, time.Second*2)
		t.AssertNil(err)
		t.Assert(result, data)
	})
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("10000")
		result, err := conn.SendRecvPkgWithTimeout(data, time.Second*2, gbtcp.PkgOption{HeaderSize: 5})
		t.AssertNE(err, nil)
		t.AssertNil(result)
	})
}

func Test_Package_Option(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		defer conn.Close()
		option := gbtcp.PkgOption{HeaderSize: 1}
		for {
			data, err := conn.RecvPkg(option)
			if err != nil {
				break
			}
			gbtest.Assert(conn.SendPkg(data, option), nil)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	// SendRecvPkg with big data - failure.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF+1)
		result, err := conn.SendRecvPkg(data, gbtcp.PkgOption{HeaderSize: 1})
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	// SendRecvPkg with big data - success.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		result, err := conn.SendRecvPkg(data, gbtcp.PkgOption{HeaderSize: 1})
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Package_Option_HeadSize3(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		defer conn.Close()
		option := gbtcp.PkgOption{HeaderSize: 3}
		for {
			data, err := conn.RecvPkg(option)
			if err != nil {
				break
			}
			gbtest.Assert(conn.SendPkg(data, option), nil)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		result, err := conn.SendRecvPkg(data, gbtcp.PkgOption{HeaderSize: 3})
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Package_Option_HeadSize4(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		defer conn.Close()
		option := gbtcp.PkgOption{HeaderSize: 4}
		for {
			data, err := conn.RecvPkg(option)
			if err != nil {
				break
			}
			gbtest.Assert(conn.SendPkg(data, option), nil)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	// SendRecvPkg with big data - failure.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFFFF+1)
		_, err = conn.SendRecvPkg(data, gbtcp.PkgOption{HeaderSize: 4})
		t.Assert(err, nil)
	})
	// SendRecvPkg with big data - success.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		result, err := conn.SendRecvPkg(data, gbtcp.PkgOption{HeaderSize: 4})
		t.AssertNil(err)
		t.Assert(result, data)
	})
	// pkgOption.HeaderSize oversize
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		_, err = conn.SendRecvPkg(data, gbtcp.PkgOption{HeaderSize: 5})
		t.AssertNE(err, nil)
	})
}

func Test_Server_NewServerKeyCrt(t *testing.T) {
	var (
		noCrtFile = "noCrtFile"
		noKeyFile = "noKeyFile"
		crtFile   = gbfile.Dir(gbdebug.CallerFilePath()) + gbfile.Separator + "testdata/crtFile"
		keyFile   = gbfile.Dir(gbdebug.CallerFilePath()) + gbfile.Separator + "testdata/keyFile"
	)
	gbtest.C(t, func(t *gbtest.T) {
		addr := "127.0.0.1:%d"
		freePort, _ := gbtcp.GetFreePort()
		addr = fmt.Sprintf(addr, freePort)
		s, err := gbtcp.NewServerKeyCrt(addr, noCrtFile, noKeyFile, func(conn *gbtcp.Conn) {
		})
		if err != nil {
			t.AssertNil(s)
		}
	})
	gbtest.C(t, func(t *gbtest.T) {
		addr := "127.0.0.1:%d"
		freePort, _ := gbtcp.GetFreePort()
		addr = fmt.Sprintf(addr, freePort)
		s, err := gbtcp.NewServerKeyCrt(addr, crtFile, noKeyFile, func(conn *gbtcp.Conn) {
		})
		if err != nil {
			t.AssertNil(s)
		}
	})
	gbtest.C(t, func(t *gbtest.T) {
		addr := "127.0.0.1:%d"
		freePort, _ := gbtcp.GetFreePort()
		addr = fmt.Sprintf(addr, freePort)
		s, err := gbtcp.NewServerKeyCrt(addr, crtFile, keyFile, func(conn *gbtcp.Conn) {
		})
		if err != nil {
			t.AssertNil(s)
		}
	})
}

func Test_Conn_RecvPkgError(t *testing.T) {

	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		defer conn.Close()
		option := gbtcp.PkgOption{HeaderSize: 5}
		for {
			_, err := conn.RecvPkg(option)
			if err != nil {
				break
			}
		}
	})
	go s.Run()
	defer s.Close()

	time.Sleep(100 * time.Millisecond)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 65536)
		result, err := conn.SendRecvPkg(data)
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
}
