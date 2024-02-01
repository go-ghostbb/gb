package gbfile

import (
	"bytes"
	gberror "ghostbb.io/errors/gb_error"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// Home returns absolute path of current user's home directory.
// The optional parameter `names` specifies the sub-folders/sub-files,
// which will be joined with current system separator and returned with the path.
func Home(names ...string) (string, error) {
	path, err := getHomePath()
	if err != nil {
		return "", err
	}
	for _, name := range names {
		path += Separator + name
	}
	return path, nil
}

// getHomePath returns absolute path of current user's home directory.
func getHomePath() (string, error) {
	u, err := user.Current()
	if nil == err {
		return u.HomeDir, nil
	}
	if runtime.GOOS == "windows" {
		return homeWindows()
	}
	return homeUnix()
}

// homeUnix retrieves and returns the home on unix system.
func homeUnix() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		err = gberror.Wrapf(err, `retrieve home directory failed`)
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", gberror.New("blank output when reading home directory")
	}

	return result, nil
}

// homeWindows retrieves and returns the home on windows system.
func homeWindows() (string, error) {
	var (
		drive = os.Getenv("HOMEDRIVE")
		path  = os.Getenv("HOMEPATH")
		home  = drive + path
	)
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", gberror.New("environment keys HOMEDRIVE, HOMEPATH and USERPROFILE are empty")
	}

	return home, nil
}
