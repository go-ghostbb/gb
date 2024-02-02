package gbproc

import (
	"context"
	"fmt"
	gbqueue "ghostbb.io/gb/container/gb_queue"
	gbtype "ghostbb.io/gb/container/gb_type"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/json"
	gbtcp "ghostbb.io/gb/net/gb_tcp"
	gbfile "ghostbb.io/gb/os/gb_file"
	gblog "ghostbb.io/gb/os/gb_log"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"net"
)

var (
	// tcpListened marks whether the receiving listening service started.
	tcpListened = gbtype.NewBool()
)

// Receive blocks and receives message from other process using local TCP listening.
// Note that, it only enables the TCP listening service when this function called.
func Receive(group ...string) *MsgRequest {
	// Use atomic operations to guarantee only one receiver goroutine listening.
	if tcpListened.Cas(false, true) {
		go receiveTcpListening()
	}
	var groupName string
	if len(group) > 0 {
		groupName = group[0]
	} else {
		groupName = defaultGroupNameForProcComm
	}
	queue := commReceiveQueues.GetOrSetFuncLock(groupName, func() interface{} {
		return gbqueue.New(maxLengthForProcMsgQueue)
	}).(*gbqueue.Queue)

	// Blocking receiving.
	if v := queue.Pop(); v != nil {
		return v.(*MsgRequest)
	}
	return nil
}

// receiveTcpListening scans local for available port and starts listening.
func receiveTcpListening() {
	var (
		listen  *net.TCPListener
		conn    net.Conn
		port    = gbtcp.MustGetFreePort()
		address = fmt.Sprintf("127.0.0.1:%d", port)
	)
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(gberror.Wrap(err, `net.ResolveTCPAddr failed`))
	}
	listen, err = net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		panic(gberror.Wrapf(err, `net.ListenTCP failed for address "%s"`, address))
	}
	// Save the port to the pid file.
	if err = gbfile.PutContents(getCommFilePath(Pid()), gbconv.String(port)); err != nil {
		panic(err)
	}
	// Start listening.
	for {
		if conn, err = listen.Accept(); err != nil {
			gblog.Error(context.TODO(), err)
		} else if conn != nil {
			go receiveTcpHandler(gbtcp.NewConnByNetConn(conn))
		}
	}
}

// receiveTcpHandler is the connection handler for receiving data.
func receiveTcpHandler(conn *gbtcp.Conn) {
	var (
		ctx      = context.TODO()
		result   []byte
		response MsgResponse
	)
	for {
		response.Code = 0
		response.Message = ""
		response.Data = nil
		buffer, err := conn.RecvPkg()
		if len(buffer) > 0 {
			// Package decoding.
			msg := new(MsgRequest)
			if err = json.UnmarshalUseNumber(buffer, msg); err != nil {
				continue
			}
			if msg.ReceiverPid != Pid() {
				// Not mine package.
				response.Message = fmt.Sprintf(
					"receiver pid not match, target: %d, current: %d",
					msg.ReceiverPid, Pid(),
				)
			} else if v := commReceiveQueues.Get(msg.Group); v == nil {
				// Group check.
				response.Message = fmt.Sprintf("group [%s] does not exist", msg.Group)
			} else {
				// Push to buffer queue.
				response.Code = 1
				v.(*gbqueue.Queue).Push(msg)
			}
		} else {
			// Empty package.
			response.Message = "empty package"
		}
		if err == nil {
			result, err = json.Marshal(response)
			if err != nil {
				gblog.Error(ctx, err)
			}
			if err = conn.SendPkg(result); err != nil {
				gblog.Error(ctx, err)
			}
		} else {
			// Just close the connection if any error occurs.
			if err = conn.Close(); err != nil {
				gblog.Error(ctx, err)
			}
			break
		}
	}
}
