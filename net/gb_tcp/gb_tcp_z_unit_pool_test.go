package gbtcp_test

import (
	gbtcp "ghostbb.io/gb/net/gb_tcp"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
	"time"
)

func Test_Pool_Basic1(t *testing.T) {
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
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.SendPkg(data)
		t.AssertNil(err)
		err = conn.SendPkgWithTimeout(data, time.Second)
		t.AssertNil(err)
	})

	gbtest.C(t, func(t *gbtest.T) {
		_, err := gbtcp.NewPoolConn("127.0.0.1:80")
		t.AssertNE(err, nil)
	})
}

func Test_Pool_Basic2(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		conn.Close()
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.SendPkg(data)
		t.AssertNil(err)
		//err = conn.SendPkgWithTimeout(data, time.Second)
		//t.AssertNil(err)

		_, err = conn.SendRecv(data, -1)
		t.AssertNE(err, nil)
	})
}

func Test_Pool_Send(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.Send(data)
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Pool_Recv(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.Send(data)
		t.AssertNil(err)
		time.Sleep(100 * time.Millisecond)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Pool_RecvLine(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999\n")
		err = conn.Send(data)
		t.AssertNil(err)
		time.Sleep(100 * time.Millisecond)
		result, err := conn.RecvLine()
		t.AssertNil(err)
		splitData := gbstr.Split(string(data), "\n")
		t.Assert(result, splitData[0])
	})
}

func Test_Pool_RecvTill(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999\n")
		err = conn.Send(data)
		t.AssertNil(err)
		time.Sleep(100 * time.Millisecond)
		result, err := conn.RecvTill([]byte("\n"))
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Pool_RecvWithTimeout(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.Send(data)
		t.AssertNil(err)
		time.Sleep(100 * time.Millisecond)
		result, err := conn.RecvWithTimeout(-1, time.Millisecond*500)
		t.AssertNil(err)
		t.Assert(data, result)
	})
}

func Test_Pool_SendWithTimeout(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.SendWithTimeout(data, time.Millisecond*500)
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(data, result)
	})
}

func Test_Pool_SendRecvWithTimeout(t *testing.T) {
	s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(s.GetListenedAddress())
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		result, err := conn.SendRecvWithTimeout(data, -1, time.Millisecond*500)
		t.AssertNil(err)
		t.Assert(data, result)
	})
}
