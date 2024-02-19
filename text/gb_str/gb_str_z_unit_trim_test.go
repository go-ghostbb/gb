package gbstr_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

func Test_Trim(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Trim(" 123456\n "), "123456")
		t.Assert(gbstr.Trim("#123456#;", "#;"), "123456")
	})
}

func Test_TrimStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimStr("gogo我愛gogo", "go"), "我愛")
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimStr("gogo我愛gogo", "go", 1), "go我愛go")
		t.Assert(gbstr.TrimStr("gogo我愛gogo", "go", 2), "我愛")
		t.Assert(gbstr.TrimStr("gogo我愛gogo", "go", -1), "我愛")
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimStr("啊我愛台灣人啊", "啊"), "我愛台灣人")
	})
}

func Test_TrimRight(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimRight(" 123456\n "), " 123456")
		t.Assert(gbstr.TrimRight("#123456#;", "#;"), "#123456")
	})
}

func Test_TrimRightStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimRightStr("gogo我愛gogo", "go"), "gogo我愛")
		t.Assert(gbstr.TrimRightStr("gogo我愛gogo", "go我愛gogo"), "go")
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimRightStr("gogo我愛gogo", "go", 1), "gogo我愛go")
		t.Assert(gbstr.TrimRightStr("gogo我愛gogo", "go", 2), "gogo我愛")
		t.Assert(gbstr.TrimRightStr("gogo我愛gogo", "go", -1), "gogo我愛")
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimRightStr("我愛台灣人", "人"), "我愛台灣")
		t.Assert(gbstr.TrimRightStr("我愛台灣人", "愛台灣人"), "我")
	})
}

func Test_TrimLeft(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimLeft(" \r123456\n "), "123456\n ")
		t.Assert(gbstr.TrimLeft("#;123456#;", "#;"), "123456#;")
	})
}

func Test_TrimLeftStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimLeftStr("gogo我愛gogo", "go"), "我愛gogo")
		t.Assert(gbstr.TrimLeftStr("gogo我愛gogo", "gogo我愛go"), "go")
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimLeftStr("gogo我愛gogo", "go", 1), "go我愛gogo")
		t.Assert(gbstr.TrimLeftStr("gogo我愛gogo", "go", 2), "我愛gogo")
		t.Assert(gbstr.TrimLeftStr("gogo我愛gogo", "go", -1), "我愛gogo")
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimLeftStr("我愛台灣人", "我愛"), "台灣人")
		t.Assert(gbstr.TrimLeftStr("我愛台灣人", "我愛台灣"), "人")
	})
}

func Test_TrimAll(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimAll("gogo我go\n愛gogo\n", "go"), "我愛")
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimAll("gogo\n我go愛gogo", "go"), "我愛")
		t.Assert(gbstr.TrimAll("gogo\n我go愛gogo\n", "go"), "我愛")
		t.Assert(gbstr.TrimAll("gogo\n我go\n愛gogo", "go"), "我愛")
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.TrimAll("啊我愛\n啊台灣\n人啊", "啊"), "我愛台灣人")
	})
}
