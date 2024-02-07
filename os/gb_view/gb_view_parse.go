package gbview

import (
	"bytes"
	"context"
	"fmt"
	gbmap "ghostbb.io/gb/container/gb_map"
	gbhash "ghostbb.io/gb/encoding/gb_hash"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/intlog"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbfsnotify "ghostbb.io/gb/os/gb_fsnotify"
	gblog "ghostbb.io/gb/os/gb_log"
	gbmlock "ghostbb.io/gb/os/gb_mlock"
	gbres "ghostbb.io/gb/os/gb_res"
	gbspath "ghostbb.io/gb/os/gb_spath"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbutil "ghostbb.io/gb/util/gb_util"
	htmltpl "html/template"
	"strconv"
	texttpl "text/template"
)

const (
	// Template name for content parsing.
	templateNameForContentParsing = "TemplateContent"
)

// fileCacheItem is the cache item for template file.
type fileCacheItem struct {
	path    string
	folder  string
	content string
}

var (
	// Templates cache map for template folder.
	// Note that there's no expiring logic for this map.
	templates = gbmap.NewStrAnyMap(true)

	// Try-folders for resource template file searching.
	resourceTryFolders = []string{
		"template/", "template", "/template", "/template/",
		"resource/template/", "resource/template", "/resource/template", "/resource/template/",
	}

	// Prefix array for trying searching in local system.
	localSystemTryFolders = []string{"", "template/", "resource/template"}
)

// Parse parses given template file `file` with given template variables `params`
// and returns the parsed template content.
func (view *View) Parse(ctx context.Context, file string, params ...Params) (result string, err error) {
	var usedParams Params
	if len(params) > 0 {
		usedParams = params[0]
	}
	return view.ParseOption(ctx, Option{
		File:    file,
		Content: "",
		Orphan:  false,
		Params:  usedParams,
	})
}

// ParseDefault parses the default template file with params.
func (view *View) ParseDefault(ctx context.Context, params ...Params) (result string, err error) {
	var usedParams Params
	if len(params) > 0 {
		usedParams = params[0]
	}
	return view.ParseOption(ctx, Option{
		File:    view.config.DefaultFile,
		Content: "",
		Orphan:  false,
		Params:  usedParams,
	})
}

// ParseContent parses given template content `content`  with template variables `params`
// and returns the parsed content in []byte.
func (view *View) ParseContent(ctx context.Context, content string, params ...Params) (string, error) {
	var usedParams Params
	if len(params) > 0 {
		usedParams = params[0]
	}
	return view.ParseOption(ctx, Option{
		Content: content,
		Orphan:  false,
		Params:  usedParams,
	})
}

// Option for template parsing.
type Option struct {
	File    string // Template file path in absolute or relative to searching paths.
	Content string // Template content, it ignores `File` if `Content` is given.
	Orphan  bool   // If true, the `File` is considered as a single file parsing without files recursively parsing from its folder.
	Params  Params // Template parameters map.
}

// ParseOption implements template parsing using Option.
func (view *View) ParseOption(ctx context.Context, option Option) (result string, err error) {
	if option.Content != "" {
		return view.doParseContent(ctx, option.Content, option.Params)
	}
	if option.File == "" {
		return "", gberror.New(`template file cannot be empty`)
	}
	// It caches the file, folder and content to enhance performance.
	r := view.fileCacheMap.GetOrSetFuncLock(option.File, func() interface{} {
		var (
			path     string
			folder   string
			content  string
			resource *gbres.File
		)
		// Searching the absolute file path for `file`.
		path, folder, resource, err = view.searchFile(ctx, option.File)
		if err != nil {
			return nil
		}
		if resource != nil {
			content = string(resource.Content())
		} else {
			content = gbfile.GetContentsWithCache(path)
		}
		// Monitor template files changes using fsnotify asynchronously.
		if resource == nil {
			if _, err = gbfsnotify.AddOnce("gbview.Parse:"+folder, folder, func(event *gbfsnotify.Event) {
				// CLEAR THEM ALL.
				view.fileCacheMap.Clear()
				templates.Clear()
				gbfsnotify.Exit()
			}); err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
		}
		return &fileCacheItem{
			path:    path,
			folder:  folder,
			content: content,
		}
	})
	if r == nil {
		return
	}
	item := r.(*fileCacheItem)
	// It's not necessary continuing parsing if template content is empty.
	if item.content == "" {
		return "", nil
	}
	// If it's Orphan option, it just parses the single file by ParseContent.
	if option.Orphan {
		return view.doParseContent(ctx, item.content, option.Params)
	}
	// Get the template object instance for `folder`.
	var tpl interface{}
	tpl, err = view.getTemplate(item.path, item.folder, fmt.Sprintf(`*%s`, gbfile.Ext(item.path)))
	if err != nil {
		return "", err
	}
	// Using memory lock to ensure concurrent safety for template parsing.
	gbmlock.LockFunc("gview.Parse:"+item.path, func() {
		if view.config.AutoEncode {
			tpl, err = tpl.(*htmltpl.Template).Parse(item.content)
		} else {
			tpl, err = tpl.(*texttpl.Template).Parse(item.content)
		}
		if err != nil && item.path != "" {
			err = gberror.Wrap(err, item.path)
		}
	})
	if err != nil {
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing `params` or view.data because both variables are pointers.
	// It needs to merge the values of the two maps into a new map.
	variables := gbutil.MapMergeCopy(option.Params)
	if len(view.data) > 0 {
		gbutil.MapMerge(variables, view.data)
	}
	view.setI18nLanguageFromCtx(ctx, variables)

	buffer := bytes.NewBuffer(nil)
	if view.config.AutoEncode {
		newTpl, err := tpl.(*htmltpl.Template).Clone()
		if err != nil {
			return "", err
		}
		if err = newTpl.Execute(buffer, variables); err != nil {
			return "", err
		}
	} else {
		if err = tpl.(*texttpl.Template).Execute(buffer, variables); err != nil {
			return "", err
		}
	}

	// TODO any graceful plan to replace "<no value>"?
	result = gbstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(ctx, result, variables)
	return result, nil
}

// doParseContent parses given template content `content`  with template variables `params`
// and returns the parsed content in []byte.
func (view *View) doParseContent(ctx context.Context, content string, params Params) (string, error) {
	// It's not necessary continuing parsing if template content is empty.
	if content == "" {
		return "", nil
	}
	var (
		err error
		key = fmt.Sprintf("%s_%v_%v", templateNameForContentParsing, view.config.Delimiters, view.config.AutoEncode)
		tpl = templates.GetOrSetFuncLock(key, func() interface{} {
			if view.config.AutoEncode {
				return htmltpl.New(templateNameForContentParsing).Delims(
					view.config.Delimiters[0],
					view.config.Delimiters[1],
				).Funcs(view.funcMap)
			}
			return texttpl.New(templateNameForContentParsing).Delims(
				view.config.Delimiters[0],
				view.config.Delimiters[1],
			).Funcs(view.funcMap)
		})
	)
	// Using memory lock to ensure concurrent safety for content parsing.
	hash := strconv.FormatUint(gbhash.DJB64([]byte(content)), 10)
	gbmlock.LockFunc("gview.ParseContent:"+hash, func() {
		if view.config.AutoEncode {
			tpl, err = tpl.(*htmltpl.Template).Parse(content)
		} else {
			tpl, err = tpl.(*texttpl.Template).Parse(content)
		}
	})
	if err != nil {
		err = gberror.Wrapf(err, `template parsing failed`)
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing `params` or view.data because both variables are pointers.
	// It needs to merge the values of the two maps into a new map.
	variables := gbutil.MapMergeCopy(params)
	if len(view.data) > 0 {
		gbutil.MapMerge(variables, view.data)
	}
	view.setI18nLanguageFromCtx(ctx, variables)

	buffer := bytes.NewBuffer(nil)
	if view.config.AutoEncode {
		var newTpl *htmltpl.Template
		newTpl, err = tpl.(*htmltpl.Template).Clone()
		if err != nil {
			err = gberror.Wrapf(err, `template clone failed`)
			return "", err
		}
		if err = newTpl.Execute(buffer, variables); err != nil {
			err = gberror.Wrapf(err, `template parsing failed`)
			return "", err
		}
	} else {
		if err = tpl.(*texttpl.Template).Execute(buffer, variables); err != nil {
			err = gberror.Wrapf(err, `template parsing failed`)
			return "", err
		}
	}
	// TODO any graceful plan to replace "<no value>"?
	result := gbstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(ctx, result, variables)
	return result, nil
}

// getTemplate returns the template object associated with given template file `path`.
// It uses template cache to enhance performance, that is, it will return the same template object
// with the same given `path`. It will also automatically refresh the template cache
// if the template files under `path` changes (recursively).
func (view *View) getTemplate(filePath, folderPath, pattern string) (tpl interface{}, err error) {
	var (
		mapKey  = fmt.Sprintf("%s_%v", filePath, view.config.Delimiters)
		mapFunc = func() interface{} {
			tplName := filePath
			if view.config.AutoEncode {
				tpl = htmltpl.New(tplName).Delims(
					view.config.Delimiters[0],
					view.config.Delimiters[1],
				).Funcs(view.funcMap)
			} else {
				tpl = texttpl.New(tplName).Delims(
					view.config.Delimiters[0],
					view.config.Delimiters[1],
				).Funcs(view.funcMap)
			}
			// Firstly checking the resource manager.
			if !gbres.IsEmpty() {
				if files := gbres.ScanDirFile(folderPath, pattern, true); len(files) > 0 {
					if view.config.AutoEncode {
						var t = tpl.(*htmltpl.Template)
						for _, v := range files {
							_, err = t.New(v.FileInfo().Name()).Parse(string(v.Content()))
							if err != nil {
								err = view.formatTemplateObjectCreatingError(v.Name(), tplName, err)
								return nil
							}
						}
					} else {
						var t = tpl.(*texttpl.Template)
						for _, v := range files {
							_, err = t.New(v.FileInfo().Name()).Parse(string(v.Content()))
							if err != nil {
								err = view.formatTemplateObjectCreatingError(v.Name(), tplName, err)
								return nil
							}
						}
					}
					return tpl
				}
			}

			// Secondly checking the file system,
			// and then automatically parsing all its sub-files recursively.
			var files []string
			files, err = gbfile.ScanDir(folderPath, pattern, true)
			if err != nil {
				return nil
			}
			if view.config.AutoEncode {
				t := tpl.(*htmltpl.Template)
				for _, file := range files {
					if _, err = t.Parse(gbfile.GetContents(file)); err != nil {
						err = view.formatTemplateObjectCreatingError(file, tplName, err)
						return nil
					}
				}
			} else {
				t := tpl.(*texttpl.Template)
				for _, file := range files {
					if _, err = t.Parse(gbfile.GetContents(file)); err != nil {
						err = view.formatTemplateObjectCreatingError(file, tplName, err)
						return nil
					}
				}
			}
			return tpl
		}
	)
	result := templates.GetOrSetFuncLock(mapKey, mapFunc)
	if result != nil {
		return result, nil
	}
	return
}

// formatTemplateObjectCreatingError formats the error that created from creating template object.
func (view *View) formatTemplateObjectCreatingError(filePath, tplName string, err error) error {
	if err != nil {
		return gberror.NewSkip(1, gbstr.Replace(err.Error(), tplName, filePath))
	}
	return nil
}

// searchFile returns the found absolute path for `file` and its template folder path.
// Note that, the returned `folder` is the template folder path, but not the folder of
// the returned template file `path`.
func (view *View) searchFile(ctx context.Context, file string) (path string, folder string, resource *gbres.File, err error) {
	var tempPath string
	// Firstly checking the resource manager.
	if !gbres.IsEmpty() {
		// Try folders.
		for _, tryFolder := range resourceTryFolders {
			tempPath = tryFolder + file
			if resource = gbres.Get(tempPath); resource != nil {
				path = resource.Name()
				folder = tryFolder
				return
			}
		}
		// Search folders.
		view.searchPaths.RLockFunc(func(array []string) {
			for _, searchPath := range array {
				for _, tryFolder := range resourceTryFolders {
					tempPath = searchPath + tryFolder + file
					if resFile := gbres.Get(tempPath); resFile != nil {
						path = resFile.Name()
						folder = searchPath + tryFolder
						return
					}
				}
			}
		})
	}

	// Secondly checking the file system.
	if path == "" {
		// Absolute path.
		path = gbfile.RealPath(file)
		if path != "" {
			folder = gbfile.Dir(path)
			return
		}
		// In search paths.
		view.searchPaths.RLockFunc(func(array []string) {
			for _, searchPath := range array {
				searchPath = gbstr.TrimRight(searchPath, `\/`)
				for _, tryFolder := range localSystemTryFolders {
					relativePath := gbstr.TrimRight(
						gbfile.Join(tryFolder, file),
						`\/`,
					)
					if path, _ = gbspath.Search(searchPath, relativePath); path != "" {
						folder = gbfile.Join(searchPath, tryFolder)
						return
					}
				}
			}
		})
	}

	// Error checking.
	if path == "" {
		buffer := bytes.NewBuffer(nil)
		if view.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("cannot find template file \"%s\" in following paths:", file))
			view.searchPaths.RLockFunc(func(array []string) {
				index := 1
				for _, searchPath := range array {
					searchPath = gbstr.TrimRight(searchPath, `\/`)
					for _, tryFolder := range localSystemTryFolders {
						buffer.WriteString(fmt.Sprintf(
							"\n%d. %s",
							index, gbfile.Join(searchPath, tryFolder),
						))
						index++
					}
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf("cannot find template file \"%s\" with no path set/add", file))
		}
		if errorPrint() {
			gblog.Error(ctx, buffer.String())
		}
		err = gberror.NewCodef(gbcode.CodeInvalidParameter, `template file "%s" not found`, file)
	}
	return
}
