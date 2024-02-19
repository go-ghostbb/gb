package gbvalid_test

import (
	"context"
	gbi18n "ghostbb.io/gb/i18n/gb_i18n"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbvalid "ghostbb.io/gb/util/gb_valid"
	"testing"
)

func TestValidator_I18n(t *testing.T) {
	var (
		err         gbvalid.Error
		i18nManager = gbi18n.New(gbi18n.Options{Path: gbtest.DataPath("i18n")})
		ctxTW       = gbi18n.WithLanguage(context.TODO(), "tw")
		validator   = gbvalid.New().I18n(i18nManager)
	)
	gbtest.C(t, func(t *gbtest.T) {
		err = validator.Rules("required").Data("").Run(ctx)
		t.Assert(err.String(), "The field is required")

		err = validator.Rules("required").Data("").Run(ctxTW)
		t.Assert(err.String(), "字段不能為空")
	})
	gbtest.C(t, func(t *gbtest.T) {
		err = validator.Rules("required").Messages("CustomMessage").Data("").Run(ctxTW)
		t.Assert(err.String(), "自定義錯誤")
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Params struct {
			Page      int `v:"required|min:1         # page is required"`
			Size      int `v:"required|between:1,100 # size is required"`
			ProjectId int `v:"between:1,10000        # project id must between {min}, {max}"`
		}
		obj := &Params{
			Page: 1,
			Size: 10,
		}
		err = validator.Data(obj).Run(ctxTW)
		t.Assert(err.String(), "項目ID必須大於等於1並且要小於等於10000")
	})
}
