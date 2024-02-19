package gbvalid_test

import (
	"context"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbvalid "ghostbb.io/gb/util/gb_valid"
	"testing"
)

func Test_CheckMap1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		data := map[string]interface{}{
			"id":   "0",
			"name": "john",
		}
		rules := map[string]string{
			"id":   "required|between:1,100",
			"name": "required|length:6,16",
		}
		if m := g.Validator().Data(data).Rules(rules).Run(context.TODO()); m == nil {
			t.Error("CheckMap校驗失敗")
		} else {
			t.Assert(len(m.Maps()), 2)
			t.Assert(m.Maps()["id"]["between"], "The id value `0` must be between 1 and 100")
			t.Assert(m.Maps()["name"]["length"], "The name value `john` length must be between 6 and 16")
		}
	})
}

func Test_CheckMap2(t *testing.T) {
	var params interface{}
	gbtest.C(t, func(t *gbtest.T) {
		if err := g.Validator().Data(params).Run(context.TODO()); err == nil {
			t.AssertNil(err)
		}
	})

	kvmap := map[string]interface{}{
		"id":   "0",
		"name": "john",
	}
	rules := map[string]string{
		"id":   "required|between:1,100",
		"name": "required|length:6,16",
	}
	msgs := gbvalid.CustomMsg{
		"id": "ID不能為空|ID範圍應當為{min}到{max}",
		"name": map[string]string{
			"required": "名稱不能為空",
			"length":   "名稱長度為{min}到{max}個字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules).Messages(msgs).Run(context.TODO()); m == nil {
		t.Error("CheckMap校驗失敗")
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules = map[string]string{
		"id":   "required|between:1,100",
		"name": "required|length:4,16",
	}
	msgs = map[string]interface{}{
		"id": "ID不能為空|ID範圍應當為{min}到{max}",
		"name": map[string]string{
			"required": "名稱不能為空",
			"length":   "名稱長度為{min}到{max}個字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules = map[string]string{
		"id":   "",
		"name": "",
	}
	msgs = map[string]interface{}{
		"id": "ID不能為空|ID範圍應當為{min}到{max}",
		"name": map[string]string{
			"required": "名稱不能為空",
			"length":   "名稱長度為{min}到{max}個字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules2 := []string{
		"@required|between:1,100",
		"@required|length:4,16",
	}
	msgs = map[string]interface{}{
		"id": "ID不能為空|ID範圍應當為{min}到{max}",
		"name": map[string]string{
			"required": "名稱不能為空",
			"length":   "名稱長度為{min}到{max}個字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules2).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules2 = []string{
		"id@required|between:1,100",
		"name@required|length:4,16#名稱不能為空|",
	}
	msgs = map[string]interface{}{
		"id": "ID不能為空|ID範圍應當為{min}到{max}",
		"name": map[string]string{
			"required": "名稱不能為空",
			"length":   "名稱長度為{min}到{max}個字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules2).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}

	kvmap = map[string]interface{}{
		"id":   "1",
		"name": "john",
	}
	rules2 = []string{
		"id@required|between:1,100",
		"name@required|length:4,16#名稱不能為空",
	}
	msgs = map[string]interface{}{
		"id": "ID不能為空|ID範圍應當為{min}到{max}",
		"name": map[string]string{
			"required": "名稱不能為空",
			"length":   "名稱長度為{min}到{max}個字符",
		},
	}
	if m := g.Validator().Data(kvmap).Rules(rules2).Messages(msgs).Run(context.TODO()); m != nil {
		t.Error(m)
	}
}

func Test_CheckMapWithNilAndNotRequiredField(t *testing.T) {
	data := map[string]interface{}{
		"id": "1",
	}
	rules := map[string]string{
		"id":   "required",
		"name": "length:4,16",
	}
	if m := g.Validator().Data(data).Rules(rules).Run(context.TODO()); m != nil {
		t.Error(m)
	}
}

func Test_Sequence(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		params := map[string]interface{}{
			"passport":  "",
			"password":  "123456",
			"password2": "1234567",
		}
		rules := []string{
			"passport@required|length:6,16#帳號不能為空|帳號長度應當在{min}到{max}之間",
			"password@required|length:6,16|same:password2#密碼不能為空|密碼長度應當在{min}到{max}之間|兩次密碼输入不相等",
			"password2@required|length:6,16#",
		}
		err := g.Validator().Data(params).Rules(rules).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Map()), 2)
		t.Assert(err.Map()["required"], "帳號不能為空")
		t.Assert(err.Map()["length"], "帳號長度應當在6到16之間")
		t.Assert(len(err.Maps()), 2)

		t.Assert(len(err.Items()), 2)
		t.Assert(err.Items()[0]["passport"]["length"], "帳號長度應當在6到16之間")
		t.Assert(err.Items()[0]["passport"]["required"], "帳號不能為空")
		t.Assert(err.Items()[1]["password"]["same"], "兩次密碼输入不相等")

		t.Assert(err.String(), "帳號不能為空; 帳號長度應當在6到16之間; 兩次密碼输入不相等")
		t.Assert(err.Strings(), []string{"帳號不能為空", "帳號長度應當在6到16之間", "兩次密碼输入不相等"})

		k, m := err.FirstItem()
		t.Assert(k, "passport")
		t.Assert(m, err.Map())

		r, s := err.FirstRule()
		t.Assert(r, "required")
		t.Assert(s, "帳號不能為空")

		t.Assert(gberror.Current(err), "帳號不能為空")
	})
}

func Test_Map_Bail(t *testing.T) {
	// global bail
	gbtest.C(t, func(t *gbtest.T) {
		params := map[string]interface{}{
			"passport":  "",
			"password":  "123456",
			"password2": "1234567",
		}
		rules := []string{
			"passport@required|length:6,16#帳號不能為空|帳號長度應當在{min}到{max}之間",
			"password@required|length:6,16|same:password2#密碼不能為空|密碼長度應當在{min}到{max}之間|兩次密碼輸入不相等",
			"password2@required|length:6,16#",
		}
		err := g.Validator().Bail().Rules(rules).Data(params).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.String(), "帳號不能為空")
	})
	// global bail with rule bail
	gbtest.C(t, func(t *gbtest.T) {
		params := map[string]interface{}{
			"passport":  "",
			"password":  "123456",
			"password2": "1234567",
		}
		rules := []string{
			"passport@bail|required|length:6,16#|帳號不能為空|帳號長度應當在{min}到{max}之間",
			"password@required|length:6,16|same:password2#密碼不能為空|密碼長度應當在{min}到{max}之間|兩次密碼輸入不相等",
			"password2@required|length:6,16#",
		}
		err := g.Validator().Bail().Rules(rules).Data(params).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.String(), "帳號不能為空")
	})
}
