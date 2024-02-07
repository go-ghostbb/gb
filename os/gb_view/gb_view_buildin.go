package gbview

import (
	"bytes"
	"context"
	"fmt"
	gbhtml "ghostbb.io/gb/encoding/gb_html"
	gbjson "ghostbb.io/gb/encoding/gb_json"
	gburl "ghostbb.io/gb/encoding/gb_url"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbmode "ghostbb.io/gb/util/gb_mode"
	gbutil "ghostbb.io/gb/util/gb_util"
	htmltpl "html/template"
	"strings"
)

// buildInFuncDump implements build-in template function: dump
func (view *View) buildInFuncDump(values ...interface{}) string {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("\n")
	buffer.WriteString("<!--\n")
	if gbmode.IsDevelop() {
		for _, v := range values {
			gbutil.DumpTo(buffer, v, gbutil.DumpOption{})
			buffer.WriteString("\n")
		}
	} else {
		buffer.WriteString("dump feature is disabled as process is not running in develop mode\n")
	}
	buffer.WriteString("-->\n")
	return buffer.String()
}

// buildInFuncMap implements build-in template function: map
func (view *View) buildInFuncMap(value ...interface{}) map[string]interface{} {
	if len(value) > 0 {
		return gbconv.Map(value[0])
	}
	return map[string]interface{}{}
}

// buildInFuncMaps implements build-in template function: maps
func (view *View) buildInFuncMaps(value ...interface{}) []map[string]interface{} {
	if len(value) > 0 {
		return gbconv.Maps(value[0])
	}
	return []map[string]interface{}{}
}

// buildInFuncEq implements build-in template function: eq
func (view *View) buildInFuncEq(value interface{}, others ...interface{}) bool {
	s := gbconv.String(value)
	for _, v := range others {
		if strings.Compare(s, gbconv.String(v)) == 0 {
			return true
		}
	}
	return false
}

// buildInFuncNe implements build-in template function: ne
func (view *View) buildInFuncNe(value, other interface{}) bool {
	return strings.Compare(gbconv.String(value), gbconv.String(other)) != 0
}

// buildInFuncLt implements build-in template function: lt
func (view *View) buildInFuncLt(value, other interface{}) bool {
	s1 := gbconv.String(value)
	s2 := gbconv.String(other)
	if gbstr.IsNumeric(s1) && gbstr.IsNumeric(s2) {
		return gbconv.Int64(value) < gbconv.Int64(other)
	}
	return strings.Compare(s1, s2) < 0
}

// buildInFuncLe implements build-in template function: le
func (view *View) buildInFuncLe(value, other interface{}) bool {
	s1 := gbconv.String(value)
	s2 := gbconv.String(other)
	if gbstr.IsNumeric(s1) && gbstr.IsNumeric(s2) {
		return gbconv.Int64(value) <= gbconv.Int64(other)
	}
	return strings.Compare(s1, s2) <= 0
}

// buildInFuncGt implements build-in template function: gt
func (view *View) buildInFuncGt(value, other interface{}) bool {
	s1 := gbconv.String(value)
	s2 := gbconv.String(other)
	if gbstr.IsNumeric(s1) && gbstr.IsNumeric(s2) {
		return gbconv.Int64(value) > gbconv.Int64(other)
	}
	return strings.Compare(s1, s2) > 0
}

// buildInFuncGe implements build-in template function: ge
func (view *View) buildInFuncGe(value, other interface{}) bool {
	s1 := gbconv.String(value)
	s2 := gbconv.String(other)
	if gbstr.IsNumeric(s1) && gbstr.IsNumeric(s2) {
		return gbconv.Int64(value) >= gbconv.Int64(other)
	}
	return strings.Compare(s1, s2) >= 0
}

// buildInFuncInclude implements build-in template function: include
// Note that configuration AutoEncode does not affect the output of this function.
func (view *View) buildInFuncInclude(file interface{}, data ...map[string]interface{}) htmltpl.HTML {
	var m map[string]interface{} = nil
	if len(data) > 0 {
		m = data[0]
	}
	path := gbconv.String(file)
	if path == "" {
		return ""
	}
	// It will search the file internally.
	content, err := view.Parse(context.TODO(), path, m)
	if err != nil {
		return htmltpl.HTML(err.Error())
	}
	return htmltpl.HTML(content)
}

// buildInFuncText implements build-in template function: text
func (view *View) buildInFuncText(html interface{}) string {
	return gbhtml.StripTags(gbconv.String(html))
}

// buildInFuncHtmlEncode implements build-in template function: html
func (view *View) buildInFuncHtmlEncode(html interface{}) string {
	return gbhtml.Entities(gbconv.String(html))
}

// buildInFuncHtmlDecode implements build-in template function: htmldecode
func (view *View) buildInFuncHtmlDecode(html interface{}) string {
	return gbhtml.EntitiesDecode(gbconv.String(html))
}

// buildInFuncUrlEncode implements build-in template function: url
func (view *View) buildInFuncUrlEncode(url interface{}) string {
	return gburl.Encode(gbconv.String(url))
}

// buildInFuncUrlDecode implements build-in template function: urldecode
func (view *View) buildInFuncUrlDecode(url interface{}) string {
	if content, err := gburl.Decode(gbconv.String(url)); err == nil {
		return content
	} else {
		return err.Error()
	}
}

// buildInFuncDate implements build-in template function: date
func (view *View) buildInFuncDate(format interface{}, timestamp ...interface{}) string {
	t := int64(0)
	if len(timestamp) > 0 {
		t = gbconv.Int64(timestamp[0])
	}
	if t == 0 {
		t = gbtime.Timestamp()
	}
	return gbtime.NewFromTimeStamp(t).Format(gbconv.String(format))
}

// buildInFuncCompare implements build-in template function: compare
func (view *View) buildInFuncCompare(value1, value2 interface{}) int {
	return strings.Compare(gbconv.String(value1), gbconv.String(value2))
}

// buildInFuncSubStr implements build-in template function: substr
func (view *View) buildInFuncSubStr(start, end, str interface{}) string {
	return gbstr.SubStrRune(gbconv.String(str), gbconv.Int(start), gbconv.Int(end))
}

// buildInFuncStrLimit implements build-in template function: strlimit
func (view *View) buildInFuncStrLimit(length, suffix, str interface{}) string {
	return gbstr.StrLimitRune(gbconv.String(str), gbconv.Int(length), gbconv.String(suffix))
}

// buildInFuncConcat implements build-in template function: concat
func (view *View) buildInFuncConcat(str ...interface{}) string {
	var s string
	for _, v := range str {
		s += gbconv.String(v)
	}
	return s
}

// buildInFuncReplace implements build-in template function: replace
func (view *View) buildInFuncReplace(search, replace, str interface{}) string {
	return gbstr.Replace(gbconv.String(str), gbconv.String(search), gbconv.String(replace), -1)
}

// buildInFuncHighlight implements build-in template function: highlight
func (view *View) buildInFuncHighlight(key, color, str interface{}) string {
	return gbstr.Replace(gbconv.String(str), gbconv.String(key), fmt.Sprintf(`<span style="color:%v;">%v</span>`, color, key))
}

// buildInFuncHideStr implements build-in template function: hidestr
func (view *View) buildInFuncHideStr(percent, hide, str interface{}) string {
	return gbstr.HideStr(gbconv.String(str), gbconv.Int(percent), gbconv.String(hide))
}

// buildInFuncToUpper implements build-in template function: toupper
func (view *View) buildInFuncToUpper(str interface{}) string {
	return gbstr.ToUpper(gbconv.String(str))
}

// buildInFuncToLower implements build-in template function: toupper
func (view *View) buildInFuncToLower(str interface{}) string {
	return gbstr.ToLower(gbconv.String(str))
}

// buildInFuncNl2Br implements build-in template function: nl2br
func (view *View) buildInFuncNl2Br(str interface{}) string {
	return gbstr.Nl2Br(gbconv.String(str))
}

// buildInFuncJson implements build-in template function: json ,
// which encodes and returns `value` as JSON string.
func (view *View) buildInFuncJson(value interface{}) (string, error) {
	b, err := gbjson.Marshal(value)
	return string(b), err
}

// buildInFuncXml implements build-in template function: xml ,
// which encodes and returns `value` as XML string.
func (view *View) buildInFuncXml(value interface{}, rootTag ...string) (string, error) {
	b, err := gbjson.New(value).ToXml(rootTag...)
	return string(b), err
}

// buildInFuncIni implements build-in template function: ini ,
// which encodes and returns `value` as XML string.
func (view *View) buildInFuncIni(value interface{}) (string, error) {
	b, err := gbjson.New(value).ToIni()
	return string(b), err
}

// buildInFuncYaml implements build-in template function: yaml ,
// which encodes and returns `value` as YAML string.
func (view *View) buildInFuncYaml(value interface{}) (string, error) {
	b, err := gbjson.New(value).ToYaml()
	return string(b), err
}

// buildInFuncYamlIndent implements build-in template function: yamli ,
// which encodes and returns `value` as YAML string with custom indent string.
func (view *View) buildInFuncYamlIndent(value, indent interface{}) (string, error) {
	b, err := gbjson.New(value).ToYamlIndent(gbconv.String(indent))
	return string(b), err
}

// buildInFuncToml implements build-in template function: toml ,
// which encodes and returns `value` as TOML string.
func (view *View) buildInFuncToml(value interface{}) (string, error) {
	b, err := gbjson.New(value).ToToml()
	return string(b), err
}

// buildInFuncPlus implements build-in template function: plus ,
// which returns the result that pluses all `deltas` to `value`.
func (view *View) buildInFuncPlus(value interface{}, deltas ...interface{}) string {
	result := gbconv.Float64(value)
	for _, v := range deltas {
		result += gbconv.Float64(v)
	}
	return gbconv.String(result)
}

// buildInFuncMinus implements build-in template function: minus ,
// which returns the result that subtracts all `deltas` from `value`.
func (view *View) buildInFuncMinus(value interface{}, deltas ...interface{}) string {
	result := gbconv.Float64(value)
	for _, v := range deltas {
		result -= gbconv.Float64(v)
	}
	return gbconv.String(result)
}

// buildInFuncTimes implements build-in template function: times ,
// which returns the result that multiplies `value` by all of `values`.
func (view *View) buildInFuncTimes(value interface{}, values ...interface{}) string {
	result := gbconv.Float64(value)
	for _, v := range values {
		result *= gbconv.Float64(v)
	}
	return gbconv.String(result)
}

// buildInFuncDivide implements build-in template function: divide ,
// which returns the result that divides `value` by all of `values`.
func (view *View) buildInFuncDivide(value interface{}, values ...interface{}) string {
	result := gbconv.Float64(value)
	for _, v := range values {
		value2Float64 := gbconv.Float64(v)
		if value2Float64 == 0 {
			// Invalid `value2`.
			return "0"
		}
		result /= value2Float64
	}
	return gbconv.String(result)
}
