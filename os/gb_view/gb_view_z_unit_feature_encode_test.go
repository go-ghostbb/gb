package gbview_test

import (
	"context"
	"ghostbb.io/gb/frame/g"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbview "ghostbb.io/gb/os/gb_view"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_Encode_Parse(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		v.SetPath(gbtest.DataPath("tpl"))
		v.SetAutoEncode(true)
		result, err := v.Parse(context.TODO(), "encode.tpl", g.Map{
			"title": "<b>my title</b>",
		})
		t.AssertNil(err)
		t.Assert(result, "<div>&lt;b&gt;my title&lt;/b&gt;</div>")
	})
}

func Test_Encode_ParseContent(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		v := gbview.New()
		tplContent := gbfile.GetContents(gbtest.DataPath("tpl", "encode.tpl"))
		v.SetAutoEncode(true)
		result, err := v.ParseContent(context.TODO(), tplContent, g.Map{
			"title": "<b>my title</b>",
		})
		t.AssertNil(err)
		t.Assert(result, "<div>&lt;b&gt;my title&lt;/b&gt;</div>")
	})
}
