package gbview

import (
	"context"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbi18n "ghostbb.io/gb/i18n/gb_i18n"
	"ghostbb.io/gb/internal/intlog"
	gbfile "ghostbb.io/gb/os/gb_file"
	gblog "ghostbb.io/gb/os/gb_log"
	gbres "ghostbb.io/gb/os/gb_res"
	gbspath "ghostbb.io/gb/os/gb_spath"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
)

// Config is the configuration object for template engine.
type Config struct {
	Paths       []string               `json:"paths"`       // Searching array for path, NOT concurrent-safe for performance purpose.
	Data        map[string]interface{} `json:"data"`        // Global template variables including configuration.
	DefaultFile string                 `json:"defaultFile"` // Default template file for parsing.
	Delimiters  []string               `json:"delimiters"`  // Custom template delimiters.
	AutoEncode  bool                   `json:"autoEncode"`  // Automatically encodes and provides safe html output, which is good for avoiding XSS.
	I18nManager *gbi18n.Manager        `json:"-"`           // I18n manager for the view.
}

const (
	// Default template file for parsing.
	defaultParsingFile = "index.html"
)

// DefaultConfig creates and returns a configuration object with default configurations.
func DefaultConfig() Config {
	return Config{
		DefaultFile: defaultParsingFile,
		I18nManager: gbi18n.Instance(),
		Delimiters:  make([]string, 2),
	}
}

// SetConfig sets the configuration for view.
func (view *View) SetConfig(config Config) error {
	var err error
	if len(config.Paths) > 0 {
		for _, v := range config.Paths {
			if err = view.AddPath(v); err != nil {
				return err
			}
		}
	}
	if len(config.Data) > 0 {
		view.Assigns(config.Data)
	}
	if config.DefaultFile != "" {
		view.SetDefaultFile(config.DefaultFile)
	}
	if len(config.Delimiters) > 1 {
		view.SetDelimiters(config.Delimiters[0], config.Delimiters[1])
	}
	view.config = config
	// Clear global template object cache.
	// It's just cache, do not hesitate clearing it.
	templates.Clear()

	intlog.Printf(context.TODO(), "SetConfig: %+v", view.config)
	return nil
}

// SetConfigWithMap set configurations with map for the view.
func (view *View) SetConfigWithMap(m map[string]interface{}) error {
	if len(m) == 0 {
		return gberror.NewCode(gbcode.CodeInvalidParameter, "configuration cannot be empty")
	}
	// The m now is a shallow copy of m.
	// Any changes to m does not affect the original one.
	// A little tricky, isn't it?
	m = gbutil.MapCopy(m)
	// Most common used configuration support for single view path.
	_, v1 := gbutil.MapPossibleItemByKey(m, "paths")
	_, v2 := gbutil.MapPossibleItemByKey(m, "path")
	if v1 == nil && v2 != nil {
		switch v2.(type) {
		case string:
			m["paths"] = []string{v2.(string)}
		case []string:
			m["paths"] = v2
		}
	}
	err := gbconv.Struct(m, &view.config)
	if err != nil {
		return err
	}
	return view.SetConfig(view.config)
}

// SetPath sets the template directory path for template file search.
// The parameter `path` can be absolute or relative path, but absolute path is suggested.
func (view *View) SetPath(path string) error {
	var (
		ctx      = context.TODO()
		isDir    = false
		realPath = ""
	)
	if file := gbres.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = gbfile.RealPath(path)
		if realPath == "" {
			// Relative path.
			view.searchPaths.RLockFunc(func(array []string) {
				for _, v := range array {
					if path, _ := gbspath.Search(v, path); path != "" {
						realPath = path
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = gbfile.IsDir(realPath)
		}
	}
	// Path not exist.
	if realPath == "" {
		err := gberror.NewCodef(gbcode.CodeInvalidParameter, `View.SetPath failed: path "%s" does not exist`, path)
		if errorPrint() {
			gblog.Error(ctx, err)
		}
		return err
	}
	// Should be a directory.
	if !isDir {
		err := gberror.NewCodef(gbcode.CodeInvalidParameter, `View.SetPath failed: path "%s" should be directory type`, path)
		if errorPrint() {
			gblog.Error(ctx, err)
		}
		return err
	}
	// Repeated path adding check.
	if view.searchPaths.Search(realPath) != -1 {
		return nil
	}
	view.searchPaths.Clear()
	view.searchPaths.Append(realPath)
	view.fileCacheMap.Clear()
	return nil
}

// AddPath adds an absolute or relative path to the search paths.
func (view *View) AddPath(path string) error {
	var (
		ctx      = context.TODO()
		isDir    = false
		realPath = ""
	)
	if file := gbres.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		if realPath = gbfile.RealPath(path); realPath == "" {
			// Relative path.
			view.searchPaths.RLockFunc(func(array []string) {
				for _, v := range array {
					if searchedPath, _ := gbspath.Search(v, path); searchedPath != "" {
						realPath = searchedPath
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = gbfile.IsDir(realPath)
		}
	}
	// Path not exist.
	if realPath == "" {
		err := gberror.NewCodef(gbcode.CodeInvalidParameter, `View.AddPath failed: path "%s" does not exist`, path)
		if errorPrint() {
			gblog.Error(ctx, err)
		}
		return err
	}
	// realPath should be type of folder.
	if !isDir {
		err := gberror.NewCodef(gbcode.CodeInvalidParameter, `View.AddPath failed: path "%s" should be directory type`, path)
		if errorPrint() {
			gblog.Error(ctx, err)
		}
		return err
	}
	// Repeated path adding check.
	if view.searchPaths.Search(realPath) != -1 {
		return nil
	}
	view.searchPaths.Append(realPath)
	view.fileCacheMap.Clear()
	return nil
}

// Assigns binds multiple global template variables to current view object.
// Note that it's not concurrent-safe, which means it would panic
// if it's called in multiple goroutines in runtime.
func (view *View) Assigns(data Params) {
	for k, v := range data {
		view.data[k] = v
	}
}

// Assign binds a global template variable to current view object.
// Note that it's not concurrent-safe, which means it would panic
// if it's called in multiple goroutines in runtime.
func (view *View) Assign(key string, value interface{}) {
	view.data[key] = value
}

// SetDefaultFile sets default template file for parsing.
func (view *View) SetDefaultFile(file string) {
	view.config.DefaultFile = file
}

// GetDefaultFile returns default template file for parsing.
func (view *View) GetDefaultFile() string {
	return view.config.DefaultFile
}

// SetDelimiters sets customized delimiters for template parsing.
func (view *View) SetDelimiters(left, right string) {
	view.config.Delimiters = []string{left, right}
}

// SetAutoEncode enables/disables automatically html encoding feature.
// When AutoEncode feature is enables, view engine automatically encodes and provides safe html output,
// which is good for avoid XSS.
func (view *View) SetAutoEncode(enable bool) {
	view.config.AutoEncode = enable
}

// BindFunc registers customized global template function named `name`
// with given function `function` to current view object.
// The `name` is the function name which can be called in template content.
func (view *View) BindFunc(name string, function interface{}) {
	view.funcMap[name] = function
	// Clear global template object cache.
	templates.Clear()
}

// BindFuncMap registers customized global template functions by map to current view object.
// The key of map is the template function name
// and the value of map is the address of customized function.
func (view *View) BindFuncMap(funcMap FuncMap) {
	for k, v := range funcMap {
		view.funcMap[k] = v
	}
	// Clear global template object cache.
	templates.Clear()
}

// SetI18n binds i18n manager to current view engine.
func (view *View) SetI18n(manager *gbi18n.Manager) {
	view.config.I18nManager = manager
}
