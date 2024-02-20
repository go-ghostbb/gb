package gbvalid_test

import (
	"context"
	"errors"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbuid "ghostbb.io/gb/util/gb_uid"
	gbvalid "ghostbb.io/gb/util/gb_valid"
	"testing"
)

func Test_CustomRule1(t *testing.T) {
	rule := "custom"
	gbvalid.RegisterRule(
		rule,
		func(ctx context.Context, in gbvalid.RuleFuncInput) error {
			pass := in.Value.String()
			if len(pass) != 6 {
				return errors.New(in.Message)
			}
			m := in.Data.Map()
			if m["data"] != pass {
				return errors.New(in.Message)
			}
			return nil
		},
	)

	gbtest.C(t, func(t *gbtest.T) {
		err := g.Validator().Data("123456").Rules(rule).Messages("custom message").Run(ctx)
		t.Assert(err.String(), "custom message")
		err = g.Validator().Data("123456").Assoc(g.Map{"data": "123456"}).Rules(rule).Messages("custom message").Run(ctx)
		t.AssertNil(err)
	})
	// Error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value string `v:"uid@custom#自定義錯誤"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123",
			Data:  "123456",
		}
		err := g.Validator().Data(st).Run(ctx)
		t.Assert(err.String(), "自定義錯誤")
	})
	// No error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value string `v:"uid@custom#自定義錯誤"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123456",
			Data:  "123456",
		}
		err := g.Validator().Data(st).Run(ctx)
		t.AssertNil(err)
	})
}

func Test_CustomRule2(t *testing.T) {
	rule := "required-map"
	gbvalid.RegisterRule(rule, func(ctx context.Context, in gbvalid.RuleFuncInput) error {
		m := in.Value.Map()
		if len(m) == 0 {
			return errors.New(in.Message)
		}
		return nil
	})
	// Check.
	gbtest.C(t, func(t *gbtest.T) {
		errStr := "data map should not be empty"
		t.Assert(g.Validator().Data(g.Map{}).Messages(errStr).Rules(rule).Run(ctx), errStr)
		t.Assert(g.Validator().Data(g.Map{"k": "v"}).Rules(rule).Messages(errStr).Run(ctx), nil)
	})
	// Error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value map[string]string `v:"uid@required-map#自定義錯誤"`
			Data  string            `p:"data"`
		}
		st := &T{
			Value: map[string]string{},
			Data:  "123456",
		}
		err := g.Validator().Data(st).Run(ctx)
		t.Assert(err.String(), "自定義錯誤")
	})
	// No error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value map[string]string `v:"uid@required-map#自定義錯誤"`
			Data  string            `p:"data"`
		}
		st := &T{
			Value: map[string]string{"k": "v"},
			Data:  "123456",
		}
		err := g.Validator().Data(st).Run(ctx)
		t.AssertNil(err)
	})
}

func Test_CustomRule_AllowEmpty(t *testing.T) {
	rule := "allow-empty-str"
	gbvalid.RegisterRule(rule, func(ctx context.Context, in gbvalid.RuleFuncInput) error {
		s := in.Value.String()
		if len(s) == 0 || s == "gb" {
			return nil
		}
		return errors.New(in.Message)
	})
	// Check.
	gbtest.C(t, func(t *gbtest.T) {
		errStr := "error"
		t.Assert(g.Validator().Data("").Rules(rule).Messages(errStr).Run(ctx), "")
		t.Assert(g.Validator().Data("gb").Rules(rule).Messages(errStr).Run(ctx), "")
		t.Assert(g.Validator().Data("gb2").Rules(rule).Messages(errStr).Run(ctx), errStr)
	})
	// Error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value string `v:"uid@allow-empty-str#自定義錯誤"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "",
			Data:  "123456",
		}
		err := g.Validator().Data(st).Run(ctx)
		t.AssertNil(err)
	})
	// No error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value string `v:"uid@allow-empty-str#自定義錯誤"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "john",
			Data:  "123456",
		}
		err := g.Validator().Data(st).Run(ctx)
		t.Assert(err.String(), "自定義錯誤")
	})
}

func TestValidator_RuleFunc(t *testing.T) {
	ruleName := "custom_1"
	ruleFunc := func(ctx context.Context, in gbvalid.RuleFuncInput) error {
		pass := in.Value.String()
		if len(pass) != 6 {
			return errors.New(in.Message)
		}
		if m := in.Data.Map(); m["data"] != pass {
			return errors.New(in.Message)
		}
		return nil
	}
	gbtest.C(t, func(t *gbtest.T) {
		err := g.Validator().Rules(ruleName).
			Messages("custom message").
			RuleFunc(ruleName, ruleFunc).
			Data("123456").
			Run(ctx)
		t.Assert(err.String(), "custom message")
		err = g.Validator().
			Rules(ruleName).
			Messages("custom message").
			Data("123456").Assoc(g.Map{"data": "123456"}).
			RuleFunc(ruleName, ruleFunc).
			Run(ctx)
		t.AssertNil(err)
	})
	// Error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value string `v:"uid@custom_1#自定義錯誤"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123",
			Data:  "123456",
		}
		err := g.Validator().RuleFunc(ruleName, ruleFunc).Data(st).Run(ctx)
		t.Assert(err.String(), "自定義錯誤")
	})
	// No error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value string `v:"uid@custom_1#自定義錯誤"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123456",
			Data:  "123456",
		}
		err := g.Validator().RuleFunc(ruleName, ruleFunc).Data(st).Run(ctx)
		t.AssertNil(err)
	})
}

func TestValidator_RuleFuncMap(t *testing.T) {
	ruleName := "custom_1"
	ruleFunc := func(ctx context.Context, in gbvalid.RuleFuncInput) error {
		pass := in.Value.String()
		if len(pass) != 6 {
			return errors.New(in.Message)
		}
		if m := in.Data.Map(); m["data"] != pass {
			return errors.New(in.Message)
		}
		return nil
	}
	gbtest.C(t, func(t *gbtest.T) {
		err := g.Validator().
			Rules(ruleName).
			Messages("custom message").
			RuleFuncMap(map[string]gbvalid.RuleFunc{
				ruleName: ruleFunc,
			}).Data("123456").Run(ctx)
		t.Assert(err.String(), "custom message")
		err = g.Validator().
			Rules(ruleName).
			Messages("custom message").
			Data("123456").Assoc(g.Map{"data": "123456"}).
			RuleFuncMap(map[string]gbvalid.RuleFunc{
				ruleName: ruleFunc,
			}).Run(ctx)
		t.AssertNil(err)
	})
	// Error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value string `v:"uid@custom_1#自定義錯誤"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123",
			Data:  "123456",
		}
		err := g.Validator().
			RuleFuncMap(map[string]gbvalid.RuleFunc{
				ruleName: ruleFunc,
			}).Data(st).Run(ctx)
		t.Assert(err.String(), "自定義錯誤")
	})
	// No error with struct validation.
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value string `v:"uid@custom_1#自定義錯誤"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123456",
			Data:  "123456",
		}
		err := g.Validator().
			RuleFuncMap(map[string]gbvalid.RuleFunc{
				ruleName: ruleFunc,
			}).Data(st).Run(ctx)
		t.AssertNil(err)
	})
}

func Test_CustomRule_Overwrite(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var rule = "custom-" + gbuid.S()
		gbvalid.RegisterRule(rule, func(ctx context.Context, in gbvalid.RuleFuncInput) error {
			return gberror.New("1")
		})
		t.Assert(g.Validator().Rules(rule).Data(1).Run(ctx), "1")
		gbvalid.RegisterRule(rule, func(ctx context.Context, in gbvalid.RuleFuncInput) error {
			return gberror.New("2")
		})
		t.Assert(g.Validator().Rules(rule).Data(1).Run(ctx), "2")
	})
	g.Dump(gbvalid.GetRegisteredRuleMap())
}

func Test_CustomRule3(t *testing.T) {
	ruleName := "required"
	ruleFunc := func(ctx context.Context, in gbvalid.RuleFuncInput) error {
		return errors.New(in.Message)
	}
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			Value string `v:"uid@required"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "",
			Data:  "123456",
		}
		err := g.Validator().
			RuleFuncMap(map[string]gbvalid.RuleFunc{
				ruleName: ruleFunc,
			}).Data(st).Run(ctx)
		t.Assert(err.String(), `The uid field is required`)
	})
}
