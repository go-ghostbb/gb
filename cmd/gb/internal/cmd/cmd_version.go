package cmd

import (
	"bytes"
	"context"
	"fmt"
	"ghostbb.io/gb"
	"ghostbb.io/gb/cmd/gb/internal/utility/mlog"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/frame/g"
	gbbuild "ghostbb.io/gb/os/gb_build"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbproc "ghostbb.io/gb/os/gb_proc"
	gbregex "ghostbb.io/gb/text/gb_regex"
	gbstr "ghostbb.io/gb/text/gb_str"
	"runtime"
	"strings"
	"time"
)

var (
	Version = cVersion{}
)

const (
	defaultIndent = "{{indent}}"
)

type cVersion struct {
	g.Meta `name:"version" brief:"show version information of current binary"`
}

type cVersionInput struct {
	g.Meta `name:"version"`
}

type cVersionOutput struct{}

func (c cVersion) Index(ctx context.Context, in cVersionInput) (*cVersionOutput, error) {
	detailBuffer := &detailBuffer{}
	detailBuffer.WriteString(fmt.Sprintf("%s", gb.VERSION))

	detailBuffer.appendLine(0, "Welcome to GB!")

	detailBuffer.appendLine(0, "Env Detail:")
	goVersion, ok := getGoVersion()
	if ok {
		detailBuffer.appendLine(1, fmt.Sprintf("Go Version: %s", goVersion))
		detailBuffer.appendLine(1, fmt.Sprintf("GB Version(go.mod): %s", getGBVersion(2)))
	} else {
		v, err := getGBVersionOfCurrentProject()
		if err == nil {
			detailBuffer.appendLine(1, fmt.Sprintf("GB Version(go.mod): %s", v))
		} else {
			detailBuffer.appendLine(1, fmt.Sprintf("GB Version(go.mod): %s", err.Error()))
		}
	}

	detailBuffer.appendLine(0, "CLI Detail:")
	detailBuffer.appendLine(1, fmt.Sprintf("Installed At: %s", gbfile.SelfPath()))
	info := gbbuild.Info()
	if info.GoFrame == "" {
		detailBuffer.appendLine(1, fmt.Sprintf("Built Go Version: %s", runtime.Version()))
		detailBuffer.appendLine(1, fmt.Sprintf("Built GB Version: %s", gb.VERSION))
	} else {
		if info.Git == "" {
			info.Git = "none"
		}
		detailBuffer.appendLine(1, fmt.Sprintf("Built Go Version: %s", info.Golang))
		detailBuffer.appendLine(1, fmt.Sprintf("Built GB Version: %s", info.GoFrame))
		detailBuffer.appendLine(1, fmt.Sprintf("Git Commit: %s", info.Git))
		detailBuffer.appendLine(1, fmt.Sprintf("Built Time: %s", info.Time))
	}

	detailBuffer.appendLine(0, "Others Detail:")
	detailBuffer.appendLine(1, "Docs: https://ghostbb.io")
	detailBuffer.appendLine(1, fmt.Sprintf("Now : %s", time.Now().Format(time.RFC3339)))

	mlog.Print(detailBuffer.replaceAllIndent("  "))
	return nil, nil
}

// detailBuffer is a buffer for detail information.
type detailBuffer struct {
	bytes.Buffer
}

// appendLine appends a line to the buffer with given indent level.
func (d *detailBuffer) appendLine(indentLevel int, line string) {
	d.WriteString(fmt.Sprintf("\n%s%s", strings.Repeat(defaultIndent, indentLevel), line))
}

// replaceAllIndent replaces the tab with given indent string and prints the buffer content.
func (d *detailBuffer) replaceAllIndent(indentStr string) string {
	return strings.ReplaceAll(d.String(), defaultIndent, indentStr)
}

// getGoVersion returns the go version
func getGoVersion() (goVersion string, ok bool) {
	goVersion, err := gbproc.ShellExec(context.Background(), "go version")
	if err != nil {
		return "", false
	}
	goVersion = gbstr.TrimLeftStr(goVersion, "go version ")
	goVersion = gbstr.TrimRightStr(goVersion, "\n")
	return goVersion, true
}

// getGBVersion returns the gb version of current project using.
func getGBVersion(indentLevel int) (gfVersion string) {
	pkgInfo, err := gbproc.ShellExec(context.Background(), `go list -f "{{if (not .Main)}}{{.Path}}@{{.Version}}{{end}}" -m all`)
	if err != nil {
		return "cannot find go.mod"
	}
	pkgList := gbstr.Split(pkgInfo, "\n")
	for _, v := range pkgList {
		if strings.HasPrefix(v, "ghostbb.io/gb") {
			gfVersion += fmt.Sprintf("\n%s%s", strings.Repeat(defaultIndent, indentLevel), v)
		}
	}
	return
}

// getGBVersionOfCurrentProject checks and returns the GoFrame version current project using.
func getGBVersionOfCurrentProject() (string, error) {
	goModPath := gbfile.Join(gbfile.Pwd(), "go.mod")
	if gbfile.Exists(goModPath) {
		lines := gbstr.SplitAndTrim(gbfile.GetContents(goModPath), "\n")
		for _, line := range lines {
			line = gbstr.Trim(line)
			line = gbstr.TrimLeftStr(line, "require ")
			line = gbstr.Trim(line)
			// Version 1.
			match, err := gbregex.MatchString(`^ghostbb\.io/gb\s+(.+)$`, line)
			if err != nil {
				return "", err
			}
			if len(match) <= 1 {
				// Version > 1.
				match, err = gbregex.MatchString(`^ghostbb\.io/gb/v\d\s+(.+)$`, line)
				if err != nil {
					return "", err
				}
			}
			if len(match) > 1 {
				return gbstr.Trim(match[1]), nil
			}
		}

		return "", gberror.New("cannot find gb requirement in go.mod")
	} else {
		return "", gberror.New("cannot find go.mod")
	}
}
