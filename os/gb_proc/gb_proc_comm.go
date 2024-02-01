package gbproc

import (
	"context"
	"fmt"
	gbmap "ghostbb.io/container/gb_map"
	gberror "ghostbb.io/errors/gb_error"
	"ghostbb.io/internal/intlog"
	gbtcp "ghostbb.io/net/gb_tcp"
	gbfile "ghostbb.io/os/gb_file"
	gbconv "ghostbb.io/util/gb_conv"
	"sync"
)

// MsgRequest is the request structure for process communication.
type MsgRequest struct {
	SenderPid   int    // Sender PID.
	ReceiverPid int    // Receiver PID.
	Group       string // Message group name.
	Data        []byte // Request data.
}

// MsgResponse is the response structure for process communication.
type MsgResponse struct {
	Code    int    // 1: OK; Other: Error.
	Message string // Response message.
	Data    []byte // Response data.
}

const (
	defaultFolderNameForProcComm = "gf_pid_port_mapping" // Default folder name for storing pid to port mapping files.
	defaultGroupNameForProcComm  = ""                    // Default group name.
	defaultTcpPortForProcComm    = 10000                 // Starting port number for receiver listening.
	maxLengthForProcMsgQueue     = 10000                 // Max size for each message queue of the group.
)

var (
	// commReceiveQueues is the group name to queue map for storing received data.
	// The value of the map is type of *gqueue.Queue.
	commReceiveQueues = gbmap.NewStrAnyMap(true)

	// commPidFolderPath specifies the folder path storing pid to port mapping files.
	commPidFolderPath string

	// commPidFolderPathOnce is used for lazy calculation for `commPidFolderPath` is necessary.
	commPidFolderPathOnce sync.Once
)

// getConnByPid creates and returns a TCP connection for specified pid.
func getConnByPid(pid int) (*gbtcp.PoolConn, error) {
	port := getPortByPid(pid)
	if port > 0 {
		if conn, err := gbtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", port)); err == nil {
			return conn, nil
		} else {
			return nil, err
		}
	}
	return nil, gberror.Newf(`could not find port for pid "%d"`, pid)
}

// getPortByPid returns the listening port for specified pid.
// It returns 0 if no port found for the specified pid.
func getPortByPid(pid int) int {
	path := getCommFilePath(pid)
	if path == "" {
		return 0
	}
	return gbconv.Int(gbfile.GetContentsWithCache(path))
}

// getCommFilePath returns the pid to port mapping file path for given pid.
func getCommFilePath(pid int) string {
	path, err := getCommPidFolderPath()
	if err != nil {
		intlog.Errorf(context.TODO(), `%+v`, err)
		return ""
	}
	return gbfile.Join(path, gbconv.String(pid))
}

// getCommPidFolderPath retrieves and returns the available directory for storing pid mapping files.
func getCommPidFolderPath() (folderPath string, err error) {
	commPidFolderPathOnce.Do(func() {
		availablePaths := []string{
			"/var/tmp",
			"/var/run",
		}
		if path, _ := gbfile.Home(".config"); path != "" {
			availablePaths = append(availablePaths, path)
		}
		availablePaths = append(availablePaths, gbfile.Temp())
		for _, availablePath := range availablePaths {
			checkPath := gbfile.Join(availablePath, defaultFolderNameForProcComm)
			if !gbfile.Exists(checkPath) && gbfile.Mkdir(checkPath) != nil {
				continue
			}
			if gbfile.IsWritable(checkPath) {
				commPidFolderPath = checkPath
				break
			}
		}
		if commPidFolderPath == "" {
			err = gberror.Newf(
				`cannot find available folder for storing pid to port mapping files in paths: %+v`,
				availablePaths,
			)
		}
	})
	folderPath = commPidFolderPath
	return
}
