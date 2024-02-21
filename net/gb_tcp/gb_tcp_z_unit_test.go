package gbtcp_test

import (
	"crypto/tls"
	"fmt"
	gbdebug "ghostbb.io/gb/debug/gb_debug"
	gbtcp "ghostbb.io/gb/net/gb_tcp"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
	"time"
)

var (
	simpleTimeout = time.Millisecond * 100
	sendData      = []byte("hello")
	invalidAddr   = "127.0.0.1:99999"
	crtFile       = gbfile.Dir(gbdebug.CallerFilePath()) + gbfile.Separator + "testdata/server.crt"
	keyFile       = gbfile.Dir(gbdebug.CallerFilePath()) + gbfile.Separator + "testdata/server.key"
)

func startTCPServer(addr string) *gbtcp.Server {
	s := gbtcp.NewServer(addr, func(conn *gbtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	time.Sleep(simpleTimeout)
	return s
}

func startTCPPkgServer(addr string) *gbtcp.Server {
	s := gbtcp.NewServer(addr, func(conn *gbtcp.Conn) {
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
	time.Sleep(simpleTimeout)
	return s
}

func startTCPTLSServer(addr string) *gbtcp.Server {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{
			{},
		},
	}
	s := gbtcp.NewServerTLS(addr, tlsConfig, func(conn *gbtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	time.Sleep(simpleTimeout)
	return s
}

func startTCPKeyCrtServer(addr string) *gbtcp.Server {
	s, _ := gbtcp.NewServerKeyCrt(addr, crtFile, keyFile, func(conn *gbtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	time.Sleep(simpleTimeout)
	return s
}

func TestGetFreePorts(t *testing.T) {
	ports, _ := gbtcp.GetFreePorts(2)
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertGT(ports[0], 0)
		t.AssertGT(ports[1], 0)
	})

	startTCPServer(fmt.Sprintf("%s:%d", "127.0.0.1", ports[0]))

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", ports[0]))
		t.AssertNil(err)
		defer conn.Close()
		result, err := conn.SendRecv(sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", 80))
		t.AssertNE(err, nil)
		t.AssertNil(conn)
	})
}

func TestMustGetFreePort(t *testing.T) {
	port := gbtcp.MustGetFreePort()
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", port)
	startTCPServer(addr)

	gbtest.C(t, func(t *gbtest.T) {
		result, err := gbtcp.SendRecv(addr, sendData, -1)
		t.AssertNil(err)
		t.Assert(sendData, result)
	})
}

func TestNewConn(t *testing.T) {
	addr := gbtcp.FreePortAddress

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(addr, simpleTimeout)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		s := startTCPServer(gbtcp.FreePortAddress)

		conn, err := gbtcp.NewConn(s.GetListenedAddress(), simpleTimeout)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		defer conn.Close()
		result, err := conn.SendRecv(sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

// TODO
func TestNewConnTLS(t *testing.T) {
	addr := gbtcp.FreePortAddress

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConnTLS(addr, &tls.Config{})
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		s := startTCPTLSServer(addr)

		conn, err := gbtcp.NewConnTLS(s.GetListenedAddress(), &tls.Config{
			InsecureSkipVerify: true,
			Certificates: []tls.Certificate{
				{},
			},
		})
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func TestNewConnKeyCrt(t *testing.T) {
	addr := gbtcp.FreePortAddress

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConnKeyCrt(addr, crtFile, keyFile)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		s := startTCPKeyCrtServer(addr)

		conn, err := gbtcp.NewConnKeyCrt(s.GetListenedAddress(), crtFile, keyFile)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func TestConn_Send(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		err = conn.Send(sendData, gbtcp.Retry{Count: 1})
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SendWithTimeout(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		err = conn.SendWithTimeout(sendData, time.Second, gbtcp.Retry{Count: 1})
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SendRecv(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		result, err := conn.SendRecv(sendData, -1, gbtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SendRecvWithTimeout(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		result, err := conn.SendRecvWithTimeout(sendData, -1, time.Second, gbtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_RecvWithTimeout(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.Send(sendData)
		result, err := conn.RecvWithTimeout(-1, time.Second, gbtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_RecvLine(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		data := []byte("hello\n")
		conn.Send(data)
		result, err := conn.RecvLine(gbtcp.Retry{Count: 1})
		t.AssertNil(err)
		splitData := gbstr.Split(string(data), "\n")
		t.Assert(result, splitData[0])
	})
}

func TestConn_RecvTill(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.Send(sendData)
		result, err := conn.RecvTill([]byte("hello"), gbtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SetDeadline(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.SetDeadline(time.Time{})
		err = conn.Send(sendData, gbtcp.Retry{Count: 1})
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SetReceiveBufferWait(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.SetBufferWaitRecv(time.Millisecond * 100)
		err = conn.Send(sendData, gbtcp.Retry{Count: 1})
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestNewNetConnKeyCrt(t *testing.T) {
	addr := gbtcp.FreePortAddress

	startTCPKeyCrtServer(addr)

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewNetConnKeyCrt(addr, "crtFile", keyFile, time.Second)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		conn, err := gbtcp.NewNetConnKeyCrt(addr, crtFile, keyFile, time.Second)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func TestSend(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.Send(invalidAddr, sendData, gbtcp.Retry{Count: 1})
		t.AssertNE(err, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.Send(s.GetListenedAddress(), sendData, gbtcp.Retry{Count: 1})
		t.AssertNil(err)
	})
}

func TestSendRecv(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		result, err := gbtcp.SendRecv(invalidAddr, sendData, -1)
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		result, err := gbtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestSendWithTimeout(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.SendWithTimeout(invalidAddr, sendData, time.Millisecond*500)
		t.AssertNE(err, nil)
		err = gbtcp.SendWithTimeout(s.GetListenedAddress(), sendData, time.Millisecond*500)
		t.AssertNil(err)
	})
}

func TestSendRecvWithTimeout(t *testing.T) {
	s := startTCPServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		result, err := gbtcp.SendRecvWithTimeout(invalidAddr, sendData, -1, time.Millisecond*500)
		t.AssertNil(result)
		t.AssertNE(err, nil)
		result, err = gbtcp.SendRecvWithTimeout(s.GetListenedAddress(), sendData, -1, time.Millisecond*500)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestSendPkg(t *testing.T) {
	s := startTCPPkgServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		err = gbtcp.SendPkg(invalidAddr, sendData)
		t.AssertNE(err, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.SendPkg(s.GetListenedAddress(), sendData, gbtcp.PkgOption{Retry: gbtcp.Retry{Count: 3}})
		t.AssertNil(err)
		err = gbtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
	})
}

func TestSendRecvPkg(t *testing.T) {
	s := startTCPPkgServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		_, err = gbtcp.SendRecvPkg(invalidAddr, sendData)
		t.AssertNE(err, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		result, err := gbtcp.SendRecvPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestSendPkgWithTimeout(t *testing.T) {
	s := startTCPPkgServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		err = gbtcp.SendPkgWithTimeout(invalidAddr, sendData, time.Second)
		t.AssertNE(err, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		err = gbtcp.SendPkgWithTimeout(s.GetListenedAddress(), sendData, time.Second)
		t.AssertNil(err)
	})
}

func TestSendRecvPkgWithTimeout(t *testing.T) {
	s := startTCPPkgServer(gbtcp.FreePortAddress)

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		_, err = gbtcp.SendRecvPkgWithTimeout(invalidAddr, sendData, time.Second)
		t.AssertNE(err, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := gbtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		result, err := gbtcp.SendRecvPkgWithTimeout(s.GetListenedAddress(), sendData, time.Second)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestNewServer(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		}, "NewServer")
		defer s.Close()
		go s.Run()

		time.Sleep(simpleTimeout)

		result, err := gbtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestGetServer(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbtcp.GetServer("GetServer")
		defer s.Close()
		go s.Run()

		t.Assert(s.GetAddress(), "")
	})

	gbtest.C(t, func(t *gbtest.T) {
		gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		}, "NewServer")

		s := gbtcp.GetServer("NewServer")
		defer s.Close()
		go s.Run()

		time.Sleep(simpleTimeout)

		result, err := gbtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestServer_SetAddress(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbtcp.NewServer("", func(conn *gbtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		})
		defer s.Close()
		t.Assert(s.GetAddress(), "")
		s.SetAddress(gbtcp.FreePortAddress)
		go s.Run()

		time.Sleep(simpleTimeout)

		result, err := gbtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestServer_SetHandler(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbtcp.NewServer(gbtcp.FreePortAddress, nil)
		defer s.Close()
		s.SetHandler(func(conn *gbtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		})
		go s.Run()

		time.Sleep(simpleTimeout)

		result, err := gbtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestServer_Run(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbtcp.NewServer(gbtcp.FreePortAddress, func(conn *gbtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		})
		defer s.Close()
		go s.Run()

		time.Sleep(simpleTimeout)

		result, err := gbtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})

	gbtest.C(t, func(t *gbtest.T) {
		s := gbtcp.NewServer(gbtcp.FreePortAddress, nil)
		defer s.Close()
		go func() {
			err := s.Run()
			t.AssertNE(err, nil)
		}()
	})
}
