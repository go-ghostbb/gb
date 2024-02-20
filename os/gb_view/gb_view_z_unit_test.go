package gbview_test

import (
	"context"
	"fmt"
	gbhtml "ghostbb.io/gb/encoding/gb_html"
	"ghostbb.io/gb/frame/g"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbres "ghostbb.io/gb/os/gb_res"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbview "ghostbb.io/gb/os/gb_view"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbmode "ghostbb.io/gb/util/gb_mode"
	gbuid "ghostbb.io/gb/util/gb_uid"
	"os"
	"strings"
	"testing"
	"time"
)

func init() {
	os.Setenv("gb_GVIEW_ERRORPRINT", "false")
}

func Test_Basic(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		str := `hello {{.name}},version:{{.version}};hello {{GetName}},version:{{GetVersion}};{{.other}}`
		pwd := gbfile.Pwd()
		view := gbview.New()
		view.SetDelimiters("{{", "}}")
		view.AddPath(pwd)
		view.SetPath(pwd)
		view.Assign("name", "gb")
		view.Assigns(g.Map{"version": "1.7.0"})
		view.BindFunc("GetName", func() string { return "gb" })
		view.BindFuncMap(gbview.FuncMap{"GetVersion": func() string { return "1.7.0" }})
		result, err := view.ParseContent(context.TODO(), str, g.Map{"other": "that's all"})
		t.Assert(err != nil, false)
		t.Assert(result, "hello gb,version:1.7.0;hello gb,version:1.7.0;that's all")

		// 測試api方法
		str = `hello {{.name}}`
		result, err = gbview.ParseContent(context.TODO(), str, g.Map{"name": "gb"})
		t.Assert(err != nil, false)
		t.Assert(result, "hello gb")

		// 測試instance方法
		result, err = gbview.Instance().ParseContent(context.TODO(), str, g.Map{"name": "gb"})
		t.Assert(err != nil, false)
		t.Assert(result, "hello gb")
	})
}

func Test_Func(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		str := `{{eq 1 1}};{{eq 1 2}};{{eq "A" "B"}}`
		result, err := gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `true;false;false`)

		str = `{{ne 1 2}};{{ne 1 1}};{{ne "A" "B"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `true;false;true`)

		str = `{{lt 1 2}};{{lt 1 1}};{{lt 1 0}};{{lt "A" "B"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `true;false;false;true`)

		str = `{{le 1 2}};{{le 1 1}};{{le 1 0}};{{le "A" "B"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `true;true;false;true`)

		str = `{{gt 1 2}};{{gt 1 1}};{{gt 1 0}};{{gt "A" "B"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `false;false;true;false`)

		str = `{{ge 1 2}};{{ge 1 1}};{{ge 1 0}};{{ge "A" "B"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `false;true;true;false`)

		str = `{{"<div>測試</div>"|text}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `測試`)

		str = `{{"<div>測試</div>"|html}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `&lt;div&gt;測試&lt;/div&gt;`)

		str = `{{"<div>測試</div>"|htmlencode}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `&lt;div&gt;測試&lt;/div&gt;`)

		str = `{{"&lt;div&gt;測試&lt;/div&gt;"|htmldecode}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `<div>測試</div>`)

		str = `{{"https://ghostbb.io"|url}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `https%3A%2F%2Fghostbb.io`)

		str = `{{"https://ghostbb.io"|urlencode}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `https%3A%2F%2Fghostbb.io`)

		str = `{{"https%3A%2F%2Fghostbb.io"|urldecode}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `https://ghostbb.io`)
		str = `{{"https%3NA%2F%2Fghostbb.io"|urldecode}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(gbstr.Contains(result, "invalid URL escape"), true)

		str = `{{1540822968 | date "Y-m-d"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `2018-10-29`)
		str = `{{date "Y-m-d"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)

		str = `{{"我是台灣人" | substr 2 -1}};{{"我是台灣人" | substr 2  2}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `台灣;台灣`)

		str = `{{"我是台灣人" | strlimit 2  "..."}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `我是...`)

		str = `{{"I'm台灣人" | replace "I'm" "我是"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `我是台灣人`)

		str = `{{compare "A" "B"}};{{compare "1" "2"}};{{compare 2 1}};{{compare 1 1}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `-1;-1;1;0`)

		str = `{{"熱愛GB熱愛生活" | hidestr 20  "*"}};{{"熱愛GB熱愛生活" | hidestr 50  "*"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `熱愛GB*愛生活;熱愛****生活`)

		str = `{{"熱愛GB熱愛生活" | highlight "GB" "red"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `熱愛<span style="color:red;">GB</span>熱愛生活`)

		str = `{{"gb" | toupper}};{{"GB" | tolower}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `GB;gb`)

		str = `{{concat "I" "Love" "Ghostbb"}}`
		result, err = gbview.ParseContent(context.TODO(), str, nil)
		t.AssertNil(err)
		t.Assert(result, `ILoveGhostbb`)
	})
	// eq: multiple values.
	gbtest.C(t, func(t *gbtest.T) {
		str := `{{eq 1 2 1 3 4 5}}`
		result, err := gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `true`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		str := `{{eq 6 2 1 3 4 5}}`
		result, err := gbview.ParseContent(context.TODO(), str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `false`)
	})
}

func Test_FuncNl2Br(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		str := `{{"Ghost\nbb" | nl2br}}`
		result, err := gbview.ParseContent(context.TODO(), str, nil)
		t.AssertNil(err)
		t.Assert(result, `Ghost<br>bb`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := ""
		for i := 0; i < 3000; i++ {
			s += "Ghost\nbb\n中文"
		}
		str := `{{.content | nl2br}}`
		result, err := gbview.ParseContent(context.TODO(), str, g.Map{
			"content": s,
		})
		t.AssertNil(err)
		t.Assert(result, strings.Replace(s, "\n", "<br>", -1))
	})
}

func Test_FuncInclude(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			header = `<h1>HEADER</h1>`
			main   = `<h1>hello gb</h1>`
			footer = `<h1>FOOTER</h1>`
			layout = `{{include "header.html" .}}
{{include "main.html" .}}
{{include "footer.html" .}}
{{include "footer_not_exist.html" .}}
{{include "" .}}`
			templatePath = gbfile.Temp(gbuid.S())
		)

		gbfile.Mkdir(templatePath)
		defer gbfile.Remove(templatePath)

		t.AssertNil(gbfile.PutContents(gbfile.Join(templatePath, `header.html`), header))
		t.AssertNil(gbfile.PutContents(gbfile.Join(templatePath, `main.html`), main))
		t.AssertNil(gbfile.PutContents(gbfile.Join(templatePath, `footer.html`), footer))
		t.AssertNil(gbfile.PutContents(gbfile.Join(templatePath, `layout.html`), layout))

		view := gbview.New(templatePath)
		result, err := view.Parse(context.TODO(), "notfound.html")
		t.AssertNE(err, nil)
		t.Assert(result, ``)

		result, err = view.Parse(context.TODO(), "layout.html")
		t.AssertNil(err)
		t.Assert(result, `<h1>HEADER</h1>
<h1>hello gb</h1>
<h1>FOOTER</h1>
template file "footer_not_exist.html" not found
`)

		t.AssertNil(gbfile.PutContents(gbfile.Join(templatePath, `notfound.html`), "notfound"))
		result, err = view.Parse(context.TODO(), "notfound.html")
		t.AssertNil(err)
		t.Assert(result, `notfound`)
	})
}

//func Test_SetPath(t *testing.T) {
//	gbtest.C(t, func(t *gbtest.T) {
//		view := gbview.Instance("addpath")
//		err := view.AddPath("tmp")
//		t.AssertNE(err, nil)
//
//		err = view.AddPath("gb_view.go")
//		t.AssertNE(err, nil)
//
//		os.Setenv("GB_VIEW_PATH", "tmp")
//		view = gbview.Instance("setpath")
//		err = view.SetPath("tmp")
//		t.AssertNE(err, nil)
//
//		err = view.SetPath("gb_view.go")
//		t.AssertNE(err, nil)
//
//		view = gbview.New(gbfile.Pwd())
//		err = view.SetPath("tmp")
//		t.AssertNE(err, nil)
//
//		err = view.SetPath("gb_view.go")
//		t.AssertNE(err, nil)
//
//		os.Setenv("GB_VIEW_PATH", "template")
//		gbfile.Mkdir(gbfile.Pwd() + gbfile.Separator + "template")
//		view = gbview.New()
//	})
//}

func Test_ParseContent(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		str := `{{.name}}`
		view := gbview.New()
		result, err := view.ParseContent(context.TODO(), str, g.Map{"name": func() {}})
		t.Assert(err != nil, true)
		t.Assert(result, ``)
	})
}

func Test_HotReload(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		dirPath := gbfile.Join(
			gbfile.Temp(),
			"testdata",
			"template-"+gbconv.String(gbtime.TimestampNano()),
		)
		defer gbfile.Remove(dirPath)
		filePath := gbfile.Join(dirPath, "test.html")

		// Initialize data.
		err := gbfile.PutContents(filePath, "test:{{.var}}")
		t.AssertNil(err)

		view := gbview.New(dirPath)

		time.Sleep(100 * time.Millisecond)
		result, err := view.Parse(context.TODO(), "test.html", g.Map{
			"var": "1",
		})
		t.AssertNil(err)
		t.Assert(result, `test:1`)

		// Update data.
		err = gbfile.PutContents(filePath, "test2:{{.var}}")
		t.AssertNil(err)

		time.Sleep(100 * time.Millisecond)
		result, err = view.Parse(context.TODO(), "test.html", g.Map{
			"var": "2",
		})
		t.AssertNil(err)
		t.Assert(result, `test2:2`)
	})
}

func Test_XSS(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		s := "<br>"
		r, err := v.ParseContent(context.TODO(), "{{.v}}", g.Map{
			"v": s,
		})
		t.AssertNil(err)
		t.Assert(r, s)
	})
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.SetAutoEncode(true)
		s := "<br>"
		r, err := v.ParseContent(context.TODO(), "{{.v}}", g.Map{
			"v": s,
		})
		t.AssertNil(err)
		t.Assert(r, gbhtml.Entities(s))
	})
	// Tag "if".
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.SetAutoEncode(true)
		s := "<br>"
		r, err := v.ParseContent(context.TODO(), "{{if eq 1 1}}{{.v}}{{end}}", g.Map{
			"v": s,
		})
		t.AssertNil(err)
		t.Assert(r, gbhtml.Entities(s))
	})
}

type TypeForBuildInFuncMap struct {
	Name  string
	Score float32
}

func (t *TypeForBuildInFuncMap) Test() (*TypeForBuildInFuncMap, error) {
	return &TypeForBuildInFuncMap{"john", 99.9}, nil
}

func Test_BuildInFuncMap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", new(TypeForBuildInFuncMap))
		r, err := v.ParseContent(context.TODO(), "{{range $k, $v := map .v.Test}} {{$k}}:{{$v}} {{end}}")
		t.AssertNil(err)
		t.Assert(gbstr.Contains(r, "Name:john"), true)
		t.Assert(gbstr.Contains(r, "Score:99.9"), true)
	})

	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(context.TODO(), "{{range $k, $v := map }} {{$k}}:{{$v}} {{end}}")
		t.AssertNil(err)
		t.Assert(gbstr.Contains(r, "Name:john"), false)
		t.Assert(gbstr.Contains(r, "Score:99.9"), false)
	})
}

type TypeForBuildInFuncMaps struct {
	Name  string
	Score float32
}

func (t *TypeForBuildInFuncMaps) Test() ([]*TypeForBuildInFuncMaps, error) {
	return []*TypeForBuildInFuncMaps{
		{"john", 99.9},
		{"smith", 100},
	}, nil
}

func Test_BuildInFuncMaps(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", new(TypeForBuildInFuncMaps))
		r, err := v.ParseContent(context.TODO(), "{{range $k, $v := maps .v.Test}} {{$k}}:{{$v.Name}} {{$v.Score}} {{end}}")
		t.AssertNil(err)
		t.Assert(r, ` 0:john 99.9  1:smith 100 `)
	})

	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", new(TypeForBuildInFuncMaps))
		r, err := v.ParseContent(context.TODO(), "{{range $k, $v := maps }} {{$k}}:{{$v.Name}} {{$v.Score}} {{end}}")
		t.AssertNil(err)
		t.Assert(r, ``)
	})
}

func Test_BuildInFuncDump(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", g.Map{
			"name":  "john",
			"score": 100,
		})
		r, err := v.ParseContent(context.TODO(), "{{dump .}}")
		t.AssertNil(err)
		fmt.Println(r)
		t.Assert(gbstr.Contains(r, `"name":  "john"`), true)
		t.Assert(gbstr.Contains(r, `"score": 100`), true)
	})

	gbtest.C(t, func(t *gbtest.T) {
		mode := gbmode.Mode()
		gbmode.SetTesting()
		defer gbmode.Set(mode)
		v := gbview.New()
		v.Assign("v", g.Map{
			"name":  "john",
			"score": 100,
		})
		r, err := v.ParseContent(context.TODO(), "{{dump .}}")
		t.AssertNil(err)
		fmt.Println(r)
		t.Assert(gbstr.Contains(r, `"name":  "john"`), false)
		t.Assert(gbstr.Contains(r, `"score": 100`), false)
	})
}

func Test_BuildInFuncJson(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", g.Map{
			"name": "john",
		})
		r, err := v.ParseContent(context.TODO(), "{{json .v}}")
		t.AssertNil(err)
		t.Assert(r, `{"name":"john"}`)
	})
}

func Test_BuildInFuncXml(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", g.Map{
			"name": "john",
		})
		r, err := v.ParseContent(context.TODO(), "{{xml .v}}")
		t.AssertNil(err)
		t.Assert(r, `<name>john</name>`)
	})
}

func Test_BuildInFuncIni(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", g.Map{
			"name": "john",
		})
		r, err := v.ParseContent(context.TODO(), "{{ini .v}}")
		t.AssertNil(err)
		t.Assert(r, `name=john
`)
	})
}

func Test_BuildInFuncYaml(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", g.Map{
			"name": "john",
		})
		r, err := v.ParseContent(context.TODO(), "{{yaml .v}}")
		t.AssertNil(err)
		t.Assert(r, `name: john
`)
	})
}

func Test_BuildInFuncYamlIndent(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", g.Map{
			"name": "john",
		})
		r, err := v.ParseContent(context.TODO(), `{{yamli .v "####"}}`)
		t.AssertNil(err)
		t.Assert(r, `####name: john
`)
	})
}

func Test_BuildInFuncToml(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.Assign("v", g.Map{
			"name": "john",
		})
		r, err := v.ParseContent(context.TODO(), "{{toml .v}}")
		t.AssertNil(err)
		t.Assert(r, `name = "john"
`)
	})
}

func Test_BuildInFuncPlus(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(gbctx.New(), "{{plus 1 2 3}}")
		t.AssertNil(err)
		t.Assert(r, `6`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(gbctx.New(), "{{1| plus 2}}")
		t.AssertNil(err)
		t.Assert(r, `3`)
	})
}

func Test_BuildInFuncMinus(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(gbctx.New(), "{{minus 1 2 3}}")
		t.AssertNil(err)
		t.Assert(r, `-4`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(gbctx.New(), "{{2 | minus 3}}")
		t.AssertNil(err)
		t.Assert(r, `1`)
	})
}

func Test_BuildInFuncTimes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(gbctx.New(), "{{times 1 2 3 4}}")
		t.AssertNil(err)
		t.Assert(r, `24`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(gbctx.New(), "{{2 | times 3}}")
		t.AssertNil(err)
		t.Assert(r, `6`)
	})
}

func Test_BuildInFuncDivide(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(gbctx.New(), "{{divide 8 2 2}}")
		t.AssertNil(err)
		t.Assert(r, `2`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(gbctx.New(), "{{2 | divide 4}}")
		t.AssertNil(err)
		t.Assert(r, `2`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		r, err := v.ParseContent(gbctx.New(), "{{divide 8 0}}")
		t.AssertNil(err)
		t.Assert(r, `0`)
	})
}

func Test_Custom1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		err := v.SetPath(gbtest.DataPath("custom1"))
		t.AssertNil(err)
		r, err := v.ParseOption(context.TODO(), gbview.Option{
			File:   "gbview.tpl",
			Orphan: true,
			Params: map[string]interface{}{
				"hello": "world",
			},
		})
		t.AssertNil(err)
		t.Assert(r, `test.tpl content, vars: world`)
	})
}

// template/gbview_test.html
// name:{{.name}}
func init() {
	if err := gbres.Add("H4sIAAAAAAAC/wrwZmYRYeBg4GCYzRMSwYAEJBk4GUpScwtyEktS9dOTyjJTy+NLUotL9DJKcnNCQ1gZGM2ZrqSe8j7ju3XVumvMF7dubWBgYPj/P8CbnUM/e4GIEAMDAx8DAwPMBgYGbl5UGziQbAAbyM58JRWkHVlRgDcjkwgzwonIRoOcCAPbGkEkYQcjzMPuIAgQYPjvmIRkHpLzWNlACpgYmBj6GRgYFoKVAwIAAP//heoHdkYBAAA="); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}

func Test_GBviewInGBres(t *testing.T) {
	gbres.Dump()
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.SetPath("template")
		result, err := v.Parse(context.TODO(), "gbview_test.html", g.Map{
			"name": "john",
		})
		t.AssertNil(err)
		t.Assert(result, "name:john")
	})
}
