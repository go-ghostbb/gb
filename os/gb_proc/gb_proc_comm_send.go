package gbproc

import (
	gberror "github.com/Ghostbb-io/gb/errors/gb_error"
	"github.com/Ghostbb-io/gb/internal/json"
	gbtcp "github.com/Ghostbb-io/gb/net/gb_tcp"
	"io"
)

// Send sends data to specified process of given pid.
func Send(pid int, data []byte, group ...string) error {
	msg := MsgRequest{
		SenderPid:   Pid(),
		ReceiverPid: pid,
		Group:       defaultGroupNameForProcComm,
		Data:        data,
	}
	if len(group) > 0 {
		msg.Group = group[0]
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	var conn *gbtcp.PoolConn
	conn, err = getConnByPid(pid)
	if err != nil {
		return err
	}
	defer conn.Close()
	// Do the sending.
	var result []byte
	result, err = conn.SendRecvPkg(msgBytes, gbtcp.PkgOption{
		Retry: gbtcp.Retry{
			Count: 3,
		},
	})
	if len(result) > 0 {
		response := new(MsgResponse)
		if err = json.UnmarshalUseNumber(result, response); err == nil {
			if response.Code != 1 {
				err = gberror.New(response.Message)
			}
		}
	}
	// EOF is not really an error.
	if err == io.EOF {
		err = nil
	}
	return err
}
