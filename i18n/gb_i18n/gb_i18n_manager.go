package gbi18n

import (
	"context"
	"fmt"
	gbjson "ghostbb.io/gb/encoding/gb_json"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/intlog"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbfsnotify "ghostbb.io/gb/os/gb_fsnotify"
	gbres "ghostbb.io/gb/os/gb_res"
	gbregex "ghostbb.io/gb/text/gb_regex"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"strings"
	"sync"
)

// pathType is the type for i18n file path.
type pathType string

const (
	pathTypeNone   pathType = "none"
	pathTypeNormal pathType = "normal"
	pathTypeGBres  pathType = "gbres"
)

// Manager for i18n contents, it is concurrent safe, supporting hot reload.
type Manager struct {
	mu       sync.RWMutex
	data     map[string]map[string]string // Translating map.
	pattern  string                       // Pattern for regex parsing.
	pathType pathType                     // Path type for i18n files.
	options  Options                      // configuration options.
}

// Options is used for i18n object configuration.
type Options struct {
	Path       string          // I18n files storage path.
	Language   string          // Default local language.
	Delimiters []string        // Delimiters for variable parsing.
	Resource   *gbres.Resource // Resource for i18n files.
}

var (
	// defaultDelimiters defines the default language if user does not specify in options.
	defaultLanguage = "en"

	// defaultDelimiters defines the default key variable delimiters.
	defaultDelimiters = []string{"{#", "}"}

	// i18n files searching folders.
	searchFolders = []string{"manifest/i18n", "manifest/config/i18n", "i18n"}
)

// New creates and returns a new i18n manager.
// The optional parameter `option` specifies the custom options for i18n manager.
// It uses a default one if it's not passed.
func New(options ...Options) *Manager {
	var opts Options
	var pathType = pathTypeNone
	if len(options) > 0 {
		opts = options[0]
		pathType = opts.checkPathType(opts.Path)
	} else {
		opts = Options{}
		for _, folder := range searchFolders {
			pathType = opts.checkPathType(folder)
			if pathType != pathTypeNone {
				break
			}
		}
		if opts.Path != "" {
			// To avoid of the source path of GoFrame: github.com/gogf/i18n/gi18n
			if gbfile.Exists(opts.Path + gbfile.Separator + "gi18n") {
				opts.Path = ""
				pathType = pathTypeNone
			}
		}
	}
	if len(opts.Language) == 0 {
		opts.Language = defaultLanguage
	}
	if len(opts.Delimiters) == 0 {
		opts.Delimiters = defaultDelimiters
	}
	m := &Manager{
		options: opts,
		pattern: fmt.Sprintf(
			`%s(.+?)%s`,
			gbregex.Quote(opts.Delimiters[0]),
			gbregex.Quote(opts.Delimiters[1]),
		),
		pathType: pathType,
	}
	intlog.Printf(context.TODO(), `New: %#v`, m)
	return m
}

// checkPathType checks and returns the path type for given directory path.
func (o *Options) checkPathType(dirPath string) pathType {
	if dirPath == "" {
		return pathTypeNone
	}

	if o.Resource == nil {
		o.Resource = gbres.Instance()
	}

	if o.Resource.Contains(dirPath) {
		o.Path = dirPath
		return pathTypeGBres
	}

	realPath, _ := gbfile.Search(dirPath)
	if realPath != "" {
		o.Path = realPath
		return pathTypeNormal
	}

	return pathTypeNone
}

// SetPath sets the directory path storing i18n files.
func (m *Manager) SetPath(path string) error {
	pathType := m.options.checkPathType(path)
	if pathType == pathTypeNone {
		return gberror.NewCodef(gbcode.CodeInvalidParameter, `%s does not exist`, path)
	}

	m.pathType = pathType
	intlog.Printf(context.TODO(), `SetPath[%s]: %s`, m.pathType, m.options.Path)
	// Reset the manager after path changed.
	m.reset()
	return nil
}

// SetLanguage sets the language for translator.
func (m *Manager) SetLanguage(language string) {
	m.options.Language = language
	intlog.Printf(context.TODO(), `SetLanguage: %s`, m.options.Language)
}

// SetDelimiters sets the delimiters for translator.
func (m *Manager) SetDelimiters(left, right string) {
	m.pattern = fmt.Sprintf(`%s(.+?)%s`, gbregex.Quote(left), gbregex.Quote(right))
	intlog.Printf(context.TODO(), `SetDelimiters: %v`, m.pattern)
}

// T is alias of Translate for convenience.
func (m *Manager) T(ctx context.Context, content string) string {
	return m.Translate(ctx, content)
}

// Tf is alias of TranslateFormat for convenience.
func (m *Manager) Tf(ctx context.Context, format string, values ...interface{}) string {
	return m.TranslateFormat(ctx, format, values...)
}

// TranslateFormat translates, formats and returns the `format` with configured language
// and given `values`.
func (m *Manager) TranslateFormat(ctx context.Context, format string, values ...interface{}) string {
	return fmt.Sprintf(m.Translate(ctx, format), values...)
}

// Translate translates `content` with configured language.
func (m *Manager) Translate(ctx context.Context, content string) string {
	m.init(ctx)
	m.mu.RLock()
	defer m.mu.RUnlock()
	transLang := m.options.Language
	if lang := LanguageFromCtx(ctx); lang != "" {
		transLang = lang
	}
	data := m.data[transLang]
	if data == nil {
		return content
	}
	// Parse content as name.
	if v, ok := data[content]; ok {
		return v
	}
	// Parse content as variables container.
	result, _ := gbregex.ReplaceStringFuncMatch(
		m.pattern, content,
		func(match []string) string {
			if v, ok := data[match[1]]; ok {
				return v
			}
			// return match[1] will return the content between delimiters
			// return match[0] will return the original content
			return match[0]
		})
	intlog.Printf(ctx, `Translate for language: %s`, transLang)
	return result
}

// GetContent retrieves and returns the configured content for given key and specified language.
// It returns an empty string if not found.
func (m *Manager) GetContent(ctx context.Context, key string) string {
	m.init(ctx)
	m.mu.RLock()
	defer m.mu.RUnlock()
	transLang := m.options.Language
	if lang := LanguageFromCtx(ctx); lang != "" {
		transLang = lang
	}
	if data, ok := m.data[transLang]; ok {
		return data[key]
	}
	return ""
}

// reset reset data of the manager.
func (m *Manager) reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = nil
}

// init initializes the manager for lazy initialization design.
// The i18n manager is only initialized once.
func (m *Manager) init(ctx context.Context) {
	m.mu.RLock()
	// If the data is not nil, means it's already initialized.
	if m.data != nil {
		m.mu.RUnlock()
		return
	}
	m.mu.RUnlock()

	defer func() {
		intlog.Printf(ctx, `Manager init finish: %#v`, m)
	}()

	intlog.Printf(ctx, `init path: %s`, m.options.Path)

	m.mu.Lock()
	defer m.mu.Unlock()
	switch m.pathType {
	case pathTypeGBres:
		files := m.options.Resource.ScanDirFile(m.options.Path, "*.*", true)
		if len(files) > 0 {
			var (
				path  string
				name  string
				lang  string
				array []string
			)
			m.data = make(map[string]map[string]string)
			for _, file := range files {
				name = file.Name()
				path = name[len(m.options.Path)+1:]
				array = strings.Split(path, "/")
				if len(array) > 1 {
					lang = array[0]
				} else if len(array) == 1 {
					lang = gbfile.Name(array[0])
				}
				if m.data[lang] == nil {
					m.data[lang] = make(map[string]string)
				}
				if j, err := gbjson.LoadContent(file.Content()); err == nil {
					for k, v := range j.Var().Map() {
						m.data[lang][k] = gbconv.String(v)
					}
				} else {
					intlog.Errorf(ctx, "load i18n file '%s' failed: %+v", name, err)
				}
			}
		}
	case pathTypeNormal:
		files, _ := gbfile.ScanDirFile(m.options.Path, "*.*", true)
		if len(files) == 0 {
			return
		}
		var (
			path  string
			lang  string
			array []string
		)
		m.data = make(map[string]map[string]string)
		for _, file := range files {
			path = file[len(m.options.Path)+1:]
			array = strings.Split(path, gbfile.Separator)
			if len(array) > 1 {
				lang = array[0]
			} else if len(array) == 1 {
				lang = gbfile.Name(array[0])
			}
			if m.data[lang] == nil {
				m.data[lang] = make(map[string]string)
			}
			if j, err := gbjson.LoadContent(gbfile.GetBytes(file)); err == nil {
				for k, v := range j.Var().Map() {
					m.data[lang][k] = gbconv.String(v)
				}
			} else {
				intlog.Errorf(ctx, "load i18n file '%s' failed: %+v", file, err)
			}
		}
		intlog.Printf(ctx, "i18n files loaded in path: %s", m.options.Path)
		// Monitor changes of i18n files for hot reload feature.
		_, _ = gbfsnotify.Add(m.options.Path, func(event *gbfsnotify.Event) {
			intlog.Printf(ctx, `i18n file changed: %s`, event.Path)
			// Any changes of i18n files, clear the data.
			m.reset()
			gbfsnotify.Exit()
		})
	}
}
