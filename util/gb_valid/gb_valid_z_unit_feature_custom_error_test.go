package gbvalid_test

import (
	"context"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	"strings"
	"testing"
)

func Test_Map(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			rule = "ipv4"
			val  = "0.0.0"
			err  = g.Validator().Data(val).Rules(rule).Run(context.TODO())
			msg  = map[string]string{
				"ipv4": "The value `0.0.0` is not a valid IPv4 address",
			}
		)
		t.Assert(err.Map(), msg)
	})
}

func Test_FirstString(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			rule = "ipv4"
			val  = "0.0.0"
			err  = g.Validator().Data(val).Rules(rule).Run(context.TODO())
		)
		t.Assert(err.FirstError(), "The value `0.0.0` is not a valid IPv4 address")
	})
}

func Test_CustomError1(t *testing.T) {
	rule := "integer|length:6,16"
	msgs := map[string]string{
		"integer": "請輸入一個整數",
		"length":  "參數長度不對",
	}
	e := g.Validator().Data("6.66").Rules(rule).Messages(msgs).Run(context.TODO())
	if e == nil || len(e.Map()) != 2 {
		t.Error("規則校驗失敗")
	} else {
		if v, ok := e.Map()["integer"]; ok {
			if strings.Compare(v.Error(), msgs["integer"]) != 0 {
				t.Error("錯誤訊息不匹配")
			}
		}
		if v, ok := e.Map()["length"]; ok {
			if strings.Compare(v.Error(), msgs["length"]) != 0 {
				t.Error("錯誤訊息不匹配")
			}
		}
	}
}

func Test_CustomError2(t *testing.T) {
	rule := "integer|length:6,16"
	msgs := "請輸入一個整數|參數長度不對"
	e := g.Validator().Data("6.66").Rules(rule).Messages(msgs).Run(context.TODO())
	if e == nil || len(e.Map()) != 2 {
		t.Error("规则校验失败")
	} else {
		if v, ok := e.Map()["integer"]; ok {
			if strings.Compare(v.Error(), "請輸入一個整數") != 0 {
				t.Error("錯誤訊息不匹配")
			}
		}
		if v, ok := e.Map()["length"]; ok {
			if strings.Compare(v.Error(), "參數長度不對") != 0 {
				t.Error("錯誤訊息不匹配")
			}
		}
	}
}
