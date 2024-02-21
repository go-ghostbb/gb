package gbtcp_test

import (
	gbtcp "ghostbb.io/gb/net/gb_tcp"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
	"time"
)

func Test_Pool_Package_Basic(t *testing.T) {
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
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
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
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 65536)
		err = conn.SendPkg(data)
		t.AssertNE(err, nil)
	})
	// SendRecvPkg
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
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
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 65536)
		result, err := conn.SendRecvPkg(data)
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	// SendRecvPkg with big data - success.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
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

func Test_Pool_Package_Timeout(t *testing.T) {
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
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("10000")
		result, err := conn.SendRecvPkgWithTimeout(data, time.Millisecond*500)
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("10000")
		result, err := conn.SendRecvPkgWithTimeout(data, time.Second*2)
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Pool_Package_Option(t *testing.T) {
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
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF+1)
		result, err := conn.SendRecvPkg(data, gbtcp.PkgOption{HeaderSize: 1})
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	// SendRecvPkg with big data - success.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		result, err := conn.SendRecvPkg(data, gbtcp.PkgOption{HeaderSize: 1})
		t.AssertNil(err)
		t.Assert(result, data)
	})
	// SendRecvPkgWithTimeout with big data - failure.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF+1)
		result, err := conn.SendRecvPkgWithTimeout(data, time.Second, gbtcp.PkgOption{HeaderSize: 1})
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	// SendRecvPkgWithTimeout with big data - success.
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		result, err := conn.SendRecvPkgWithTimeout(data, time.Second, gbtcp.PkgOption{HeaderSize: 1})
		t.AssertNil(err)
		t.Assert(result, data)
	})
}
