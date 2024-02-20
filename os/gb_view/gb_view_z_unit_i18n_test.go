package gbview_test

import (
	"context"
	gbdebug "ghostbb.io/gb/debug/gb_debug"
	"ghostbb.io/gb/frame/g"
	gbi18n "ghostbb.io/gb/i18n/gb_i18n"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbview "ghostbb.io/gb/os/gb_view"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_I18n(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		content := `{{.name}} says "{#hello}{#world}!"`
		expect1 := `john says "你好世界!"`
		expect2 := `john says "こんにちは世界!"`
		expect3 := `john says "{#hello}{#world}!"`

		g.I18n().SetPath(gbtest.DataPath("i18n"))

		g.I18n().SetLanguage("zh-TW")
		result1, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name": "john",
		})
		t.AssertNil(err)
		t.Assert(result1, expect1)

		g.I18n().SetLanguage("ja")
		result2, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name": "john",
		})
		t.AssertNil(err)
		t.Assert(result2, expect2)

		g.I18n().SetLanguage("none")
		result3, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name": "john",
		})
		t.AssertNil(err)
		t.Assert(result3, expect3)
	})
	gbtest.C(t, func(t *gbtest.T) {
		content := `{{.name}} says "{#hello}{#world}!"`
		expect1 := `john says "你好世界!"`
		expect2 := `john says "こんにちは世界!"`
		expect3 := `john says "{#hello}{#world}!"`

		g.I18n().SetPath(gbdebug.CallerDirectory() + gbfile.Separator + "testdata" + gbfile.Separator + "i18n")

		result1, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name":         "john",
			"I18nLanguage": "zh-TW",
		})
		t.AssertNil(err)
		t.Assert(result1, expect1)

		result2, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name":         "john",
			"I18nLanguage": "ja",
		})
		t.AssertNil(err)
		t.Assert(result2, expect2)

		result3, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name":         "john",
			"I18nLanguage": "none",
		})
		t.AssertNil(err)
		t.Assert(result3, expect3)
	})
	// gbi18n manager is nil
	gbtest.C(t, func(t *gbtest.T) {
		content := `{{.name}} says "{#hello}{#world}!"`
		expect1 := `john says "{#hello}{#world}!"`

		g.I18n().SetPath(gbdebug.CallerDirectory() + gbfile.Separator + "testdata" + gbfile.Separator + "i18n")

		view := gbview.New()
		view.SetI18n(nil)
		result1, err := view.ParseContent(context.TODO(), content, g.Map{
			"name":         "john",
			"I18nLanguage": "zh-TW",
		})
		t.AssertNil(err)
		t.Assert(result1, expect1)
	})
	// SetLanguage in context
	gbtest.C(t, func(t *gbtest.T) {
		content := `{{.name}} says "{#hello}{#world}!"`
		expect1 := `john says "你好世界!"`
		ctx := gbctx.New()
		g.I18n().SetPath(gbdebug.CallerDirectory() + gbfile.Separator + "testdata" + gbfile.Separator + "i18n")
		ctx = gbi18n.WithLanguage(ctx, "zh-TW")
		t.Log(gbi18n.LanguageFromCtx(ctx))

		view := gbview.New()

		result1, err := view.ParseContent(ctx, content, g.Map{
			"name": "john",
		})
		t.AssertNil(err)
		t.Assert(result1, expect1)
	})

}
