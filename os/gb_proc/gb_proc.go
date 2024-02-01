// Package gbproc implements management and communication for processes.
package gbproc

import (
	gbenv "github.com/Ghostbb-io/gb/os/gb_env"
	gbfile "github.com/Ghostbb-io/gb/os/gb_file"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
	"os"
	"runtime"
	"time"
)

const (
	envKeyPPid            = "GBPROC_PPID"
	tracingInstrumentName = "github.com/Ghostbb-io/gb/os/gbproc.Process"
)

var (
	processPid       = os.Getpid() // processPid is the pid of current process.
	processStartTime = time.Now()  // processStartTime is the start time of current process.
)

// Pid returns the pid of current process.
func Pid() int {
	return processPid
}

// PPid returns the custom parent pid if exists, or else it returns the system parent pid.
func PPid() int {
	if !IsChild() {
		return Pid()
	}
	ppidValue := os.Getenv(envKeyPPid)
	if ppidValue != "" && ppidValue != "0" {
		return gbconv.Int(ppidValue)
	}
	return PPidOS()
}

// PPidOS returns the system parent pid of current process.
// Note that the difference between PPidOS and PPid function is that the PPidOS returns
// the system ppid, but the PPid functions may return the custom pid by gproc if the custom
// ppid exists.
func PPidOS() int {
	return os.Getppid()
}

// IsChild checks and returns whether current process is a child process.
// A child process is forked by another gproc process.
func IsChild() bool {
	ppidValue := os.Getenv(envKeyPPid)
	return ppidValue != "" && ppidValue != "0"
}

// SetPPid sets custom parent pid for current process.
func SetPPid(ppid int) error {
	if ppid > 0 {
		return os.Setenv(envKeyPPid, gbconv.String(ppid))
	} else {
		return os.Unsetenv(envKeyPPid)
	}
}

// StartTime returns the start time of current process.
func StartTime() time.Time {
	return processStartTime
}

// Uptime returns the duration which current process has been running
func Uptime() time.Duration {
	return time.Since(processStartTime)
}

// SearchBinary searches the binary `file` in current working folder and PATH environment.
func SearchBinary(file string) string {
	// Check if it is absolute path of exists at current working directory.
	if gbfile.Exists(file) {
		return file
	}
	return SearchBinaryPath(file)
}

// SearchBinaryPath searches the binary `file` in PATH environment.
func SearchBinaryPath(file string) string {
	array := ([]string)(nil)
	switch runtime.GOOS {
	case "windows":
		envPath := gbenv.Get("PATH", gbenv.Get("Path")).String()
		if gbstr.Contains(envPath, ";") {
			array = gbstr.SplitAndTrim(envPath, ";")
		} else if gbstr.Contains(envPath, ":") {
			array = gbstr.SplitAndTrim(envPath, ":")
		}
		if gbfile.Ext(file) != ".exe" {
			file += ".exe"
		}

	default:
		array = gbstr.SplitAndTrim(gbenv.Get("PATH").String(), ":")
	}
	if len(array) > 0 {
		path := ""
		for _, v := range array {
			path = v + gbfile.Separator + file
			if gbfile.Exists(path) && gbfile.IsFile(path) {
				return path
			}
		}
	}
	return ""
}
