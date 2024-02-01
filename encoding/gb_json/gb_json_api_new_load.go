package gbjson

import (
	"bytes"
	gbini "github.com/Ghostbb-io/gb/encoding/gb_ini"
	gbproperties "github.com/Ghostbb-io/gb/encoding/gb_properties"
	gbtoml "github.com/Ghostbb-io/gb/encoding/gb_toml"
	gbxml "github.com/Ghostbb-io/gb/encoding/gb_xml"
	gbyaml "github.com/Ghostbb-io/gb/encoding/gb_yaml"
	gbcode "github.com/Ghostbb-io/gb/errors/gb_code"
	gberror "github.com/Ghostbb-io/gb/errors/gb_error"
	"github.com/Ghostbb-io/gb/internal/json"
	"github.com/Ghostbb-io/gb/internal/reflection"
	"github.com/Ghostbb-io/gb/internal/rwmutex"
	gbfile "github.com/Ghostbb-io/gb/os/gb_file"
	gbregex "github.com/Ghostbb-io/gb/text/gb_regex"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
	"reflect"
)

// New creates a Json object with any variable type of `data`, but `data` should be a map
// or slice for data access reason, or it will make no sense.
//
// The parameter `safe` specifies whether using this Json object in concurrent-safe context,
// which is false in default.
func New(data interface{}, safe ...bool) *Json {
	return NewWithTag(data, string(ContentTypeJson), safe...)
}

// NewWithTag creates a Json object with any variable type of `data`, but `data` should be a map
// or slice for data access reason, or it will make no sense.
//
// The parameter `tags` specifies priority tags for struct conversion to map, multiple tags joined
// with char ','.
//
// The parameter `safe` specifies whether using this Json object in concurrent-safe context, which
// is false in default.
func NewWithTag(data interface{}, tags string, safe ...bool) *Json {
	option := Options{
		Tags: tags,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return NewWithOptions(data, option)
}

// NewWithOptions creates a Json object with any variable type of `data`, but `data` should be a map
// or slice for data access reason, or it will make no sense.
func NewWithOptions(data interface{}, options Options) *Json {
	var j *Json
	switch data.(type) {
	case string, []byte:
		if r, err := loadContentWithOptions(data, options); err == nil {
			j = r
		} else {
			j = &Json{
				p:  &data,
				c:  byte(defaultSplitChar),
				vc: false,
			}
		}
	default:
		var (
			pointedData interface{}
			reflectInfo = reflection.OriginValueAndKind(data)
		)
		switch reflectInfo.OriginKind {
		case reflect.Slice, reflect.Array:
			pointedData = gbconv.Interfaces(data)

		case reflect.Map:
			pointedData = gbconv.MapDeep(data, options.Tags)

		case reflect.Struct:
			if v, ok := data.(iVal); ok {
				return NewWithOptions(v.Val(), options)
			}
			pointedData = gbconv.MapDeep(data, options.Tags)

		default:
			pointedData = data
		}
		j = &Json{
			p:  &pointedData,
			c:  byte(defaultSplitChar),
			vc: false,
		}
	}
	j.mu = rwmutex.Create(options.Safe)
	return j
}

// Load loads content from specified file `path`, and creates a Json object from its content.
func Load(path string, safe ...bool) (*Json, error) {
	if p, err := gbfile.Search(path); err != nil {
		return nil, err
	} else {
		path = p
	}
	options := Options{
		Type: ContentType(gbfile.Ext(path)),
	}
	if len(safe) > 0 && safe[0] {
		options.Safe = true
	}
	return doLoadContentWithOptions(gbfile.GetBytesWithCache(path), options)
}

// LoadWithOptions creates a Json object from given JSON format content and options.
func LoadWithOptions(data interface{}, options Options) (*Json, error) {
	return doLoadContentWithOptions(gbconv.Bytes(data), options)
}

// LoadJson creates a Json object from given JSON format content.
func LoadJson(data interface{}, safe ...bool) (*Json, error) {
	option := Options{
		Type: ContentTypeJson,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(gbconv.Bytes(data), option)
}

// LoadXml creates a Json object from given XML format content.
func LoadXml(data interface{}, safe ...bool) (*Json, error) {
	option := Options{
		Type: ContentTypeXml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(gbconv.Bytes(data), option)
}

// LoadIni creates a Json object from given INI format content.
func LoadIni(data interface{}, safe ...bool) (*Json, error) {
	option := Options{
		Type: ContentTypeIni,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(gbconv.Bytes(data), option)
}

// LoadYaml creates a Json object from given YAML format content.
func LoadYaml(data interface{}, safe ...bool) (*Json, error) {
	option := Options{
		Type: ContentTypeYaml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(gbconv.Bytes(data), option)
}

// LoadToml creates a Json object from given TOML format content.
func LoadToml(data interface{}, safe ...bool) (*Json, error) {
	option := Options{
		Type: ContentTypeToml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(gbconv.Bytes(data), option)
}

// LoadProperties creates a Json object from given TOML format content.
func LoadProperties(data interface{}, safe ...bool) (*Json, error) {
	option := Options{
		Type: ContentTypeProperties,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(gbconv.Bytes(data), option)
}

// LoadContent creates a Json object from given content, it checks the data type of `content`
// automatically, supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContent(data interface{}, safe ...bool) (*Json, error) {
	content := gbconv.Bytes(data)
	if len(content) == 0 {
		return New(nil, safe...), nil
	}
	return LoadContentType(checkDataType(content), content, safe...)
}

// LoadContentType creates a Json object from given type and content,
// supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContentType(dataType ContentType, data interface{}, safe ...bool) (*Json, error) {
	content := gbconv.Bytes(data)
	if len(content) == 0 {
		return New(nil, safe...), nil
	}
	// ignore UTF8-BOM
	if content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}
	options := Options{
		Type:      dataType,
		StrNumber: true,
	}
	if len(safe) > 0 && safe[0] {
		options.Safe = true
	}
	return doLoadContentWithOptions(content, options)
}

// IsValidDataType checks and returns whether given `dataType` a valid data type for loading.
func IsValidDataType(dataType ContentType) bool {
	if dataType == "" {
		return false
	}
	if dataType[0] == '.' {
		dataType = dataType[1:]
	}
	switch dataType {
	case
		ContentTypeJson,
		ContentTypeJs,
		ContentTypeXml,
		ContentTypeYaml,
		ContentTypeYml,
		ContentTypeToml,
		ContentTypeIni,
		ContentTypeProperties:
		return true
	}
	return false
}

func loadContentWithOptions(data interface{}, options Options) (*Json, error) {
	content := gbconv.Bytes(data)
	if len(content) == 0 {
		return NewWithOptions(nil, options), nil
	}
	if options.Type == "" {
		options.Type = checkDataType(content)
	}
	return loadContentTypeWithOptions(content, options)
}

func loadContentTypeWithOptions(data interface{}, options Options) (*Json, error) {
	content := gbconv.Bytes(data)
	if len(content) == 0 {
		return NewWithOptions(nil, options), nil
	}
	// ignore UTF8-BOM
	if content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}
	return doLoadContentWithOptions(content, options)
}

// doLoadContent creates a Json object from given content.
// It supports data content type as follows:
// JSON, XML, INI, YAML and TOML.
func doLoadContentWithOptions(data []byte, options Options) (*Json, error) {
	var (
		err    error
		result interface{}
	)
	if len(data) == 0 {
		return NewWithOptions(nil, options), nil
	}
	if options.Type == "" {
		options.Type = checkDataType(data)
	}
	options.Type = ContentType(gbstr.TrimLeft(
		string(options.Type), "."),
	)
	switch options.Type {
	case ContentTypeJson, ContentTypeJs:

	case ContentTypeXml:
		if data, err = gbxml.ToJson(data); err != nil {
			return nil, err
		}

	case ContentTypeYaml, ContentTypeYml:
		if data, err = gbyaml.ToJson(data); err != nil {
			return nil, err
		}

	case ContentTypeToml:
		if data, err = gbtoml.ToJson(data); err != nil {
			return nil, err
		}

	case ContentTypeIni:
		if data, err = gbini.ToJson(data); err != nil {
			return nil, err
		}
	case ContentTypeProperties:
		if data, err = gbproperties.ToJson(data); err != nil {
			return nil, err
		}

	default:
		err = gberror.NewCodef(
			gbcode.CodeInvalidParameter,
			`unsupported type "%s" for loading`,
			options.Type,
		)
	}
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if options.StrNumber {
		decoder.UseNumber()
	}
	if err = decoder.Decode(&result); err != nil {
		return nil, err
	}
	switch result.(type) {
	case string, []byte:
		return nil, gberror.Newf(`json decoding failed for content: %s`, data)
	}
	return NewWithOptions(result, options), nil
}

// checkDataType automatically checks and returns the data type for `content`.
// Note that it uses regular expression for loose checking, you can use LoadXXX/LoadContentType
// functions to load the content for certain content type.
func checkDataType(content []byte) ContentType {
	if json.Valid(content) {
		return ContentTypeJson
	} else if gbregex.IsMatch(`^<.+>[\S\s]+<.+>\s*$`, content) {
		return ContentTypeXml
	} else if !gbregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*"""[\s\S]+"""`, content) &&
		!gbregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*'''[\s\S]+'''`, content) &&
		((gbregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*".+"`, content) || gbregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*\w+`, content)) ||
			(gbregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*".+"`, content) || gbregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*\w+`, content))) {
		return ContentTypeYaml
	} else if !gbregex.IsMatch(`^[\s\t\n\r]*;.+`, content) &&
		!gbregex.IsMatch(`[\s\t\n\r]+;.+`, content) &&
		!gbregex.IsMatch(`[\n\r]+[\s\t\w\-]+\.[\s\t\w\-]+\s*=\s*.+`, content) &&
		(gbregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, content) || gbregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, content)) {
		return ContentTypeToml
	} else if gbregex.IsMatch(`\[[\w\.]+\]`, content) &&
		(gbregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, content) || gbregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, content)) {
		// Must contain "[xxx]" section.
		return ContentTypeIni
	} else if gbregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, content) {
		return ContentTypeProperties
	} else {
		return ""
	}
}
