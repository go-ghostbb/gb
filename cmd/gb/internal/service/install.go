package service

import (
	"context"
	"ghostbb.io/gb/cmd/gb/internal/utility/allyes"
	"ghostbb.io/gb/cmd/gb/internal/utility/mlog"
	gbarray "ghostbb.io/gb/container/gb_array"
	gbset "ghostbb.io/gb/container/gb_set"
	"ghostbb.io/gb/frame/g"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbenv "ghostbb.io/gb/os/gb_env"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"runtime"
	"strings"
)

var (
	Install = serviceInstall{}
)

type serviceInstall struct{}

type serviceInstallAvailablePath struct {
	dirPath   string
	filePath  string
	writable  bool
	installed bool
	IsSelf    bool
}

func (s serviceInstall) Run(ctx context.Context) (err error) {
	// Ask where to install.
	paths := s.getAvailablePaths()
	if len(paths) <= 0 {
		mlog.Printf("no path detected, you can manually install gb by copying the binary to path folder.")
		return
	}
	mlog.Printf("I found some installable paths for you(from $PATH): ")
	mlog.Printf("  %2s | %8s | %9s | %s", "Id", "Writable", "Installed", "Path")

	// Print all paths status and determine the default selectedID value.
	var (
		selectedID = -1
		newPaths   []serviceInstallAvailablePath
		pathSet    = gbset.NewStrSet() // Used for repeated items filtering.
	)
	for _, path := range paths {
		if !pathSet.AddIfNotExist(path.dirPath) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	paths = newPaths
	for id, path := range paths {
		mlog.Printf("  %2d | %8t | %9t | %s", id, path.writable, path.installed, path.dirPath)
		if selectedID == -1 {
			// Use the previously installed path as the most priority choice.
			if path.installed {
				selectedID = id
			}
		}
	}
	// If there's no previously installed path, use the first writable path.
	if selectedID == -1 {
		// Order by choosing priority.
		commonPaths := gbarray.NewStrArrayFrom(g.SliceStr{
			s.getGoPathBin(),
			`/usr/local/bin`,
			`/usr/bin`,
			`/usr/sbin`,
			`C:\Windows`,
			`C:\Windows\system32`,
			`C:\Go\bin`,
			`C:\Program Files`,
			`C:\Program Files (x86)`,
		})
		// Check the common installation directories.
		commonPaths.Iterator(func(k int, v string) bool {
			for id, aPath := range paths {
				if strings.EqualFold(aPath.dirPath, v) {
					selectedID = id
					return false
				}
			}
			return true
		})
		if selectedID == -1 {
			selectedID = 0
		}
	}

	if allyes.Check() {
		// Use the default selectedID.
		mlog.Printf("please choose one installation destination [default %d]: %d", selectedID, selectedID)
	} else {
		for {
			// Get input and update selectedID.
			var (
				inputID int
				input   = gbcmd.Scanf("please choose one installation destination [default %d]: ", selectedID)
			)
			if input != "" {
				inputID = gbconv.Int(input)
			} else {
				break
			}
			// Check if out of range.
			if inputID >= len(paths) || inputID < 0 {
				mlog.Printf("invalid install destination Id: %d", inputID)
				continue
			}
			selectedID = inputID
			break
		}
	}

	// Get selected destination path.
	dstPath := paths[selectedID]

	// Install the new binary.
	mlog.Debugf(`copy file from "%s" to "%s"`, gbfile.SelfPath(), dstPath.filePath)
	err = gbfile.CopyFile(gbfile.SelfPath(), dstPath.filePath)
	if err != nil {
		mlog.Printf("install gb binary to '%s' failed: %v", dstPath.dirPath, err)
		mlog.Printf("you can manually install gb by copying the binary to folder: %s", dstPath.dirPath)
	} else {
		mlog.Printf("gb binary is successfully installed to: %s", dstPath.filePath)
	}
	return
}

// IsInstalled checks and returns whether the binary is installed.
func (s serviceInstall) IsInstalled() (*serviceInstallAvailablePath, bool) {
	paths := s.getAvailablePaths()
	for _, aPath := range paths {
		if aPath.installed {
			return &aPath, true
		}
	}
	return nil, false
}

// getGoPathBin retrieves ad returns the GOPATH/bin path for binary.
func (s serviceInstall) getGoPathBin() string {
	if goPath := gbenv.Get(`GOPATH`).String(); goPath != "" {
		return gbfile.Join(goPath, "bin")
	}
	return ""
}

// getAvailablePaths returns the installation paths data for the binary.
func (s serviceInstall) getAvailablePaths() []serviceInstallAvailablePath {
	var (
		folderPaths    []serviceInstallAvailablePath
		binaryFileName = "gb" + gbfile.Ext(gbfile.SelfPath())
	)
	// $GOPATH/bin
	if goPathBin := s.getGoPathBin(); goPathBin != "" {
		folderPaths = s.checkAndAppendToAvailablePath(
			folderPaths, goPathBin, binaryFileName,
		)
	}
	switch runtime.GOOS {
	case "darwin":
		darwinInstallationCheckPaths := []string{"/usr/local/bin"}
		for _, v := range darwinInstallationCheckPaths {
			folderPaths = s.checkAndAppendToAvailablePath(
				folderPaths, v, binaryFileName,
			)
		}
		fallthrough

	default:
		// Search and find the writable directory path.
		envPath := gbenv.Get("PATH", gbenv.Get("Path").String()).String()
		if gbstr.Contains(envPath, ";") {
			// windows.
			for _, v := range gbstr.SplitAndTrim(envPath, ";") {
				if v == "." {
					continue
				}
				folderPaths = s.checkAndAppendToAvailablePath(
					folderPaths, v, binaryFileName,
				)
			}
		} else if gbstr.Contains(envPath, ":") {
			// *nix.
			for _, v := range gbstr.SplitAndTrim(envPath, ":") {
				if v == "." {
					continue
				}
				folderPaths = s.checkAndAppendToAvailablePath(
					folderPaths, v, binaryFileName,
				)
			}
		} else if envPath != "" {
			folderPaths = s.checkAndAppendToAvailablePath(
				folderPaths, envPath, binaryFileName,
			)
		} else {
			folderPaths = s.checkAndAppendToAvailablePath(
				folderPaths, "/usr/local/bin", binaryFileName,
			)
		}
	}
	return folderPaths
}

// checkAndAppendToAvailablePath checks if `path` is writable and already installed.
// It adds the `path` to `folderPaths` if it is writable or already installed, or else it ignores the `path`.
func (s serviceInstall) checkAndAppendToAvailablePath(folderPaths []serviceInstallAvailablePath, dirPath string, binaryFileName string) []serviceInstallAvailablePath {
	var (
		filePath  = gbfile.Join(dirPath, binaryFileName)
		writable  = gbfile.IsWritable(dirPath)
		installed = gbfile.Exists(filePath)
		self      = gbfile.SelfPath() == filePath
	)
	if !writable && !installed {
		return folderPaths
	}
	return append(
		folderPaths,
		serviceInstallAvailablePath{
			dirPath:   dirPath,
			writable:  writable,
			filePath:  filePath,
			installed: installed,
			IsSelf:    self,
		})
}
